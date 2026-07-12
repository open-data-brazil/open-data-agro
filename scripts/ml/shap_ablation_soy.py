#!/usr/bin/env python3
"""SHAP importance + feature-group ablation for soy LightGBM (Phase 31 C2)."""

from __future__ import annotations

import argparse
import json
import os
import sys
from pathlib import Path

import numpy as np

ROOT = Path(__file__).resolve().parents[2]
sys.path.insert(0, str(ROOT / "scripts"))

from ml.config import FEATURE_COLS, SUCCESS_HORIZON  # noqa: E402
from ml.metrics import mae  # noqa: E402
from ml.train_soy import (  # noqa: E402
    _require_lightgbm,
    _train_booster,
    load_rows,
    split_xy,
    POINT_ROUNDS,
)

ABLATION_GROUPS: dict[str, tuple[str, ...]] = {
    "cepea_level": ("preco_rs_sc", "preco_usd_sc", "variacao_dia_pct"),
    "cepea_momentum": ("cepea_ret_1d", "cepea_ret_7d", "cepea_vol_7d"),
    "fx": ("ptax_usd_venda",),
    "futures_basis": ("b3_front_price", "basis_cepea_usd_minus_b3"),
    "trade": ("comex_soja_fob_usd", "comex_soja_kg"),
    "logistics": ("frete_mt_pr_ton_avg", "frete_mt_pr_tkm_avg"),
    "climate": ("inmet_national_avg",),
}

SHAP_SAMPLE_CAP = 200


def _mean_abs_shap(booster: object, x_sample: np.ndarray) -> list[dict[str, float | str]]:
    import shap

    explainer = shap.TreeExplainer(booster)
    values = explainer.shap_values(x_sample)
    arr = np.asarray(values)
    if arr.ndim == 3:  # multi-output safety
        arr = arr[:, :, 0]
    means = np.mean(np.abs(arr), axis=0)
    ranked = sorted(
        (
            {"feature": FEATURE_COLS[i], "mean_abs_shap": float(means[i])}
            for i in range(len(FEATURE_COLS))
        ),
        key=lambda d: float(d["mean_abs_shap"]),
        reverse=True,
    )
    return ranked


def _ablation_delta(
    lgb: object,
    rows: list[dict[str, object]],
    target: str,
    drop_cols: tuple[str, ...],
    baseline_mae: float,
) -> dict[str, object]:
    keep = tuple(c for c in FEATURE_COLS if c not in drop_cols)
    x_train, y_train = split_xy(rows, target, "train", keep)
    x_val, y_val = split_xy(rows, target, "val", keep)
    x_test, y_test = split_xy(rows, target, "test", keep)
    # Temporary FEATURE_COLS override via feature_name in Dataset — retrain helper uses global.
    # Train inline with kept columns only.
    params = {
        "objective": "regression",
        "metric": "mae",
        "learning_rate": 0.05,
        "num_leaves": 31,
        "verbosity": -1,
        "seed": 42,
        "min_data_in_leaf": 20,
    }
    dtrain = lgb.Dataset(x_train, label=y_train, feature_name=list(keep))  # type: ignore[attr-defined]
    callbacks = [lgb.log_evaluation(0), lgb.early_stopping(40)]  # type: ignore[attr-defined]
    dval = lgb.Dataset(x_val, label=y_val, reference=dtrain, feature_name=list(keep))  # type: ignore[attr-defined]
    model = lgb.train(  # type: ignore[attr-defined]
        params, dtrain, num_boost_round=300, valid_sets=[dval], callbacks=callbacks
    )
    pred = model.predict(x_test)
    dropped_mae = mae(y_test.tolist(), pred.tolist())
    assert dropped_mae is not None
    return {
        "dropped_features": list(drop_cols),
        "test_mae": dropped_mae,
        "mae_delta_vs_full": dropped_mae - baseline_mae,
    }


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument(
        "--export-dir",
        default=os.environ.get("ML_EXPORT_DIR", str(ROOT / ".local" / "ml" / "export")),
    )
    parser.add_argument(
        "--out-dir",
        default=os.environ.get("ML_REPORTS_DIR", str(ROOT / ".local" / "ml" / "reports")),
    )
    args = parser.parse_args()

    lgb = _require_lightgbm()
    parquet_path = Path(args.export_dir) / "soy_daily_training.parquet"
    if not parquet_path.is_file():
        print(f"missing export parquet: {parquet_path}", file=sys.stderr)
        return 1

    rows = load_rows(parquet_path)
    target = f"y_ret_{SUCCESS_HORIZON}d"
    x_train, y_train = split_xy(rows, target, "train")
    x_val, y_val = split_xy(rows, target, "val")
    x_test, y_test = split_xy(rows, target, "test")
    if len(y_test) < 10:
        print("insufficient test rows for SHAP/ablation", file=sys.stderr)
        return 1

    booster = _train_booster(
        lgb, x_train, y_train, x_val, y_val,
        objective="regression", alpha=None, num_rounds=POINT_ROUNDS,
    )
    full_pred = booster.predict(x_test)
    full_mae = mae(y_test.tolist(), full_pred.tolist())
    assert full_mae is not None

    rng = np.random.default_rng(42)
    n_sample = min(SHAP_SAMPLE_CAP, len(x_test))
    idx = rng.choice(len(x_test), size=n_sample, replace=False)
    shap_ranked = _mean_abs_shap(booster, x_test[idx])

    ablation: dict[str, object] = {}
    for name, cols in ABLATION_GROUPS.items():
        ablation[name] = _ablation_delta(lgb, rows, target, cols, full_mae)

    payload = {
        "target": target,
        "full_model_test_mae": full_mae,
        "shap_sample_size": n_sample,
        "shap_mean_abs": shap_ranked,
        "ablation": ablation,
    }

    out_dir = Path(args.out_dir)
    out_dir.mkdir(parents=True, exist_ok=True)
    json_path = out_dir / "soy_shap_ablation.json"
    json_path.write_text(json.dumps(payload, indent=2, sort_keys=True) + "\n", encoding="utf-8")

    md = [
        "# Soy SHAP + ablation (Phase 31 C2)",
        "",
        f"Target: `{target}` · full-model test MAE: **{full_mae:.6g}**",
        "",
        f"SHAP sample size: {n_sample}",
        "",
        "## Mean |SHAP| (ranked)",
        "",
        "| feature | mean_abs_shap |",
        "|---------|---------------|",
    ]
    for row in shap_ranked:
        md.append(f"| {row['feature']} | {row['mean_abs_shap']:.6g} |")
    md.extend(["", "## Ablation (drop feature group, retrain)", "", "| group | test MAE | Δ vs full |", "|-------|----------|-----------|"])
    for name, stats in ablation.items():
        md.append(
            f"| {name} | {stats['test_mae']:.6g} | {stats['mae_delta_vs_full']:+.6g} |"  # type: ignore[index]
        )
    md_path = out_dir / "soy_shap_ablation.md"
    md_path.write_text("\n".join(md) + "\n", encoding="utf-8")
    print(f"wrote {json_path}")
    print(f"wrote {md_path}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
