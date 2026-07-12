#!/usr/bin/env python3
"""Walk-forward LightGBM frete forecast (Phase 32 Track F)."""

from __future__ import annotations

import argparse
import json
import os
import sys
import time
from pathlib import Path
from typing import Any

import numpy as np
import pyarrow.parquet as pq

ROOT = Path(__file__).resolve().parents[2]
sys.path.insert(0, str(ROOT / "scripts"))

from ml.baselines_frete import evaluate_split  # noqa: E402
from ml.config_frete import (  # noqa: E402
    CATEGORICAL_FEATURE_COLS,
    NUMERIC_FEATURE_COLS,
    SUCCESS_IMPROVEMENT_MIN,
    TARGET_COL,
)
from ml.metrics import gate_passes, mae, rmse  # noqa: E402

QUANTILES = (0.1, 0.5, 0.9)
FEATURE_COLS = NUMERIC_FEATURE_COLS + CATEGORICAL_FEATURE_COLS


def _require_lgb() -> Any:
    try:
        import lightgbm as lgb
    except ImportError as exc:  # pragma: no cover
        raise SystemExit(
            "lightgbm required: pip install -r toolchain/python-requirements-ml.txt"
        ) from exc
    return lgb


def _uf_code(uf: object) -> int:
    s = str(uf or "").strip().upper()
    if len(s) != 2:
        return -1
    return (ord(s[0]) - 64) * 32 + (ord(s[1]) - 64)


def split_xy(
    rows: list[dict[str, object]], split_name: str
) -> tuple[np.ndarray, np.ndarray]:
    xs: list[list[float]] = []
    ys: list[float] = []
    for row in rows:
        if row.get("split") != split_name:
            continue
        y = row.get(TARGET_COL)
        if y is None:
            continue
        feats: list[float] = []
        for col in NUMERIC_FEATURE_COLS:
            val = row.get(col)
            feats.append(float("nan") if val is None else float(val))  # type: ignore[arg-type]
        feats.append(float(_uf_code(row.get("uf_origem"))))
        feats.append(float(_uf_code(row.get("uf_destino"))))
        fonte = row.get("fonte_code")
        feats.append(float(fonte) if fonte is not None else 3.0)
        xs.append(feats)
        ys.append(float(y))  # type: ignore[arg-type]
    if not xs:
        return np.empty((0, len(FEATURE_COLS))), np.empty((0,))
    return np.asarray(xs, dtype=np.float64), np.asarray(ys, dtype=np.float64)


def _train(
    lgb: Any,
    x_tr: np.ndarray,
    y_tr: np.ndarray,
    x_va: np.ndarray,
    y_va: np.ndarray,
    *,
    objective: str,
    alpha: float | None,
    rounds: int,
) -> Any:
    params: dict[str, object] = {
        "objective": objective,
        "metric": "mae" if objective == "regression" else "quantile",
        "learning_rate": 0.05,
        "num_leaves": 63,
        "feature_fraction": 0.9,
        "bagging_fraction": 0.9,
        "bagging_freq": 1,
        "min_data_in_leaf": 20,
        "verbosity": -1,
        "seed": 42,
        "num_threads": max(1, (os.cpu_count() or 4) - 1),
    }
    if alpha is not None:
        params["alpha"] = alpha
    # Mark UF / fonte as categorical by index
    cat_idx = [len(NUMERIC_FEATURE_COLS), len(NUMERIC_FEATURE_COLS) + 1, len(NUMERIC_FEATURE_COLS) + 2]
    dtr = lgb.Dataset(x_tr, label=y_tr, feature_name=list(FEATURE_COLS), categorical_feature=cat_idx)
    callbacks = [lgb.log_evaluation(0)]
    valid = []
    if len(y_va):
        dva = lgb.Dataset(
            x_va, label=y_va, reference=dtr, feature_name=list(FEATURE_COLS), categorical_feature=cat_idx
        )
        valid = [dva]
        callbacks.append(lgb.early_stopping(40))
    return lgb.train(params, dtr, num_boost_round=rounds, valid_sets=valid or None, callbacks=callbacks)


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument(
        "--export-dir",
        default=os.environ.get(
            "ML_FRETE_EXPORT_DIR", str(ROOT / ".local" / "ml" / "export_frete")
        ),
    )
    parser.add_argument(
        "--out-dir",
        default=os.environ.get("ML_FRETE_TRAIN_DIR", str(ROOT / ".local" / "ml" / "train_frete")),
    )
    parser.add_argument(
        "--registry-dir",
        default=os.environ.get(
            "ML_FRETE_REGISTRY_DIR", str(ROOT / ".local" / "ml" / "registry")
        ),
    )
    parser.add_argument("--require-gate", action="store_true")
    args = parser.parse_args()

    lgb = _require_lgb()
    parquet_path = Path(args.export_dir) / "frete_training.parquet"
    if not parquet_path.is_file():
        print(f"missing export: {parquet_path}", file=sys.stderr)
        return 1

    t0 = time.perf_counter()
    rows = pq.read_table(parquet_path).to_pylist()
    x_tr, y_tr = split_xy(rows, "train")
    x_va, y_va = split_xy(rows, "val")
    x_te, y_te = split_xy(rows, "test")
    if len(y_tr) < 100 or len(y_te) < 50:
        print(f"insufficient rows train={len(y_tr)} test={len(y_te)}", file=sys.stderr)
        return 1

    point = _train(lgb, x_tr, y_tr, x_va, y_va, objective="regression", alpha=None, rounds=800)
    point_pred = point.predict(x_te)
    point_mae = mae(y_te.tolist(), point_pred.tolist())
    point_rmse = rmse(y_te.tolist(), point_pred.tolist())

    quantiles: dict[str, object] = {}
    for q in QUANTILES:
        booster = _train(lgb, x_tr, y_tr, x_va, y_va, objective="quantile", alpha=q, rounds=500)
        pred = booster.predict(x_te)
        quantiles[f"p{int(q * 100):02d}"] = {
            "alpha": q,
            "test_mae": mae(y_te.tolist(), pred.tolist()),
            "best_iteration": int(getattr(booster, "best_iteration", 0) or 0),
        }

    baseline = evaluate_split(rows, "test")
    seasonal_mae = baseline["seasonal_naive"]["mae"]  # type: ignore[index]
    assert point_mae is not None and seasonal_mae is not None
    gate_ok = gate_passes(point_mae, float(seasonal_mae), SUCCESS_IMPROVEMENT_MIN)
    elapsed = time.perf_counter() - t0

    manifest_hash = None
    man_path = Path(args.export_dir) / "manifest.json"
    if man_path.is_file():
        manifest_hash = json.loads(man_path.read_text(encoding="utf-8")).get("schema_hash")

    report = {
        "dataset": "ml.frete-training",
        "target": TARGET_COL,
        "n_train": int(len(y_tr)),
        "n_val": int(len(y_va)),
        "n_test": int(len(y_te)),
        "point": {"test_mae": point_mae, "test_rmse": point_rmse},
        "quantiles": quantiles,
        "baseline_test": {
            "seasonal_naive_mae": seasonal_mae,
            "corridor_mean_mae": baseline["corridor_mean"]["mae"],  # type: ignore[index]
            "global_mean_mae": baseline["global_mean"]["mae"],  # type: ignore[index]
        },
        "success_gate": {
            "min_relative_improvement": SUCCESS_IMPROVEMENT_MIN,
            "threshold_mae": float(seasonal_mae) * (1.0 - SUCCESS_IMPROVEMENT_MIN),
            "model_mae": point_mae,
            "passed": gate_ok,
        },
        "hardware": {
            "cpu_count": os.cpu_count(),
            "elapsed_sec": round(elapsed, 3),
            "num_threads": max(1, (os.cpu_count() or 4) - 1),
        },
        "manifest_schema_hash": manifest_hash,
    }

    out_dir = Path(args.out_dir)
    out_dir.mkdir(parents=True, exist_ok=True)
    out_json = out_dir / "frete_lgbm_metrics.json"
    out_json.write_text(json.dumps(report, indent=2, sort_keys=True) + "\n", encoding="utf-8")

    md = [
        "# Frete LightGBM metrics (real hardware)",
        "",
        f"Source: `{parquet_path}`",
        "",
        f"**Gate:** {'PASS' if gate_ok else 'FAIL'} · elapsed **{elapsed:.1f}s** · CPUs={os.cpu_count()}",
        "",
        f"- point MAE: `{point_mae}`",
        f"- seasonal MAE: `{seasonal_mae}`",
        f"- threshold: `{float(seasonal_mae) * (1.0 - SUCCESS_IMPROVEMENT_MIN)}`",
        "",
        "## Quantiles (test MAE)",
        "",
    ]
    for k, v in quantiles.items():
        md.append(f"- {k}: {v['test_mae']}")  # type: ignore[index]
    (out_dir / "frete_lgbm_metrics.md").write_text("\n".join(md) + "\n", encoding="utf-8")

    reg_dir = Path(args.registry_dir)
    reg_dir.mkdir(parents=True, exist_ok=True)
    reg = {
        "model_id": "frete_lgbm_walkforward_v1",
        "metrics": report,
        "feature_cols": list(FEATURE_COLS),
    }
    reg_path = reg_dir / "frete_lgbm_v1.json"
    reg_path.write_text(json.dumps(reg, indent=2, sort_keys=True) + "\n", encoding="utf-8")

    print(f"wrote {out_json}")
    print(f"wrote {reg_path}")
    print(f"gate_passed={gate_ok} elapsed_sec={elapsed:.2f}")
    if args.require_gate and not gate_ok:
        return 1
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
