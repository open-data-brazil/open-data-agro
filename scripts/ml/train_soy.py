#!/usr/bin/env python3
"""Walk-forward LightGBM point + quantile training for soy (Phase 31 Epic C1)."""

from __future__ import annotations

import argparse
import json
import os
import sys
from pathlib import Path
from typing import Any

import numpy as np
import pyarrow.parquet as pq

ROOT = Path(__file__).resolve().parents[2]
sys.path.insert(0, str(ROOT / "scripts"))

from ml.baselines_soy import evaluate_horizon  # noqa: E402
from ml.config import (  # noqa: E402
    FEATURE_COLS,
    HORIZONS,
    SUCCESS_HORIZON,
    SUCCESS_IMPROVEMENT_MIN,
)
from ml.metrics import gate_passes, mae, rmse  # noqa: E402

QUANTILES: tuple[float, ...] = (0.1, 0.5, 0.9)
POINT_ROUNDS = 600
QUANTILE_ROUNDS = 400
EARLY_STOP = 50


def _require_lightgbm() -> Any:
    try:
        import lightgbm as lgb
    except ImportError as exc:  # pragma: no cover
        raise SystemExit(
            "lightgbm required: pip install -r toolchain/python-requirements-ml.txt"
        ) from exc
    return lgb


def load_rows(parquet_path: Path) -> list[dict[str, object]]:
    return pq.read_table(parquet_path).to_pylist()


def split_xy(
    rows: list[dict[str, object]],
    target: str,
    split_name: str,
    feature_cols: tuple[str, ...] = FEATURE_COLS,
) -> tuple[np.ndarray, np.ndarray]:
    xs: list[list[float]] = []
    ys: list[float] = []
    for row in rows:
        if row.get("split") != split_name:
            continue
        y = row.get(target)
        if y is None:
            continue
        feats: list[float] = []
        for col in feature_cols:
            val = row.get(col)
            feats.append(float("nan") if val is None else float(val))  # type: ignore[arg-type]
        xs.append(feats)
        ys.append(float(y))  # type: ignore[arg-type]
    if not xs:
        return np.empty((0, len(feature_cols))), np.empty((0,))
    return np.asarray(xs, dtype=np.float64), np.asarray(ys, dtype=np.float64)


def _train_booster(
    lgb: Any,
    x_train: np.ndarray,
    y_train: np.ndarray,
    x_val: np.ndarray,
    y_val: np.ndarray,
    *,
    objective: str,
    alpha: float | None,
    num_rounds: int,
) -> Any:
    params: dict[str, object] = {
        "objective": objective,
        "metric": "mae" if objective == "regression" else "quantile",
        "learning_rate": 0.05,
        "num_leaves": 31,
        "feature_fraction": 0.85,
        "bagging_fraction": 0.85,
        "bagging_freq": 1,
        "min_data_in_leaf": 20,
        "verbosity": -1,
        "seed": 42,
    }
    if alpha is not None:
        params["alpha"] = alpha
    dtrain = lgb.Dataset(x_train, label=y_train, feature_name=list(FEATURE_COLS))
    callbacks = [lgb.log_evaluation(0)]
    valid_sets = []
    if len(y_val) > 0:
        dval = lgb.Dataset(x_val, label=y_val, reference=dtrain, feature_name=list(FEATURE_COLS))
        valid_sets = [dval]
        callbacks.append(lgb.early_stopping(EARLY_STOP))
    return lgb.train(
        params,
        dtrain,
        num_boost_round=num_rounds,
        valid_sets=valid_sets or None,
        callbacks=callbacks,
    )


def train_horizon(
    lgb: Any,
    rows: list[dict[str, object]],
    horizon: int,
) -> dict[str, object]:
    target = f"y_ret_{horizon}d"
    x_train, y_train = split_xy(rows, target, "train")
    x_val, y_val = split_xy(rows, target, "val")
    x_test, y_test = split_xy(rows, target, "test")
    if len(y_train) < 50 or len(y_test) < 10:
        raise ValueError(f"insufficient rows for {target}: train={len(y_train)} test={len(y_test)}")

    point = _train_booster(
        lgb, x_train, y_train, x_val, y_val,
        objective="regression", alpha=None, num_rounds=POINT_ROUNDS,
    )
    point_pred = point.predict(x_test)
    point_mae = mae(y_test.tolist(), point_pred.tolist())
    point_rmse = rmse(y_test.tolist(), point_pred.tolist())

    quantiles: dict[str, object] = {}
    q_preds: dict[str, np.ndarray] = {}
    for q in QUANTILES:
        booster = _train_booster(
            lgb, x_train, y_train, x_val, y_val,
            objective="quantile", alpha=q, num_rounds=QUANTILE_ROUNDS,
        )
        pred = booster.predict(x_test)
        q_preds[f"p{int(q * 100):02d}"] = pred
        quantiles[f"p{int(q * 100):02d}"] = {
            "alpha": q,
            "best_iteration": int(getattr(booster, "best_iteration", 0) or 0),
            "test_mae": mae(y_test.tolist(), pred.tolist()),
        }

    baseline = evaluate_horizon(rows, horizon, "test")
    seasonal_mae = baseline["y_ret"]["seasonal_naive"]["mae"]  # type: ignore[index]
    last_mae = baseline["y_ret"]["last_value"]["mae"]  # type: ignore[index]
    assert point_mae is not None and seasonal_mae is not None
    cleared = gate_passes(point_mae, float(seasonal_mae), SUCCESS_IMPROVEMENT_MIN)

    return {
        "horizon_trading_days": horizon,
        "target": target,
        "n_train": int(len(y_train)),
        "n_val": int(len(y_val)),
        "n_test": int(len(y_test)),
        "point": {
            "test_mae": point_mae,
            "test_rmse": point_rmse,
            "best_iteration": int(getattr(point, "best_iteration", 0) or 0),
        },
        "quantiles": quantiles,
        "baseline_test": {
            "seasonal_naive_mae": seasonal_mae,
            "last_value_mae": last_mae,
        },
        "success_gate": {
            "applies": horizon == SUCCESS_HORIZON,
            "min_relative_improvement": SUCCESS_IMPROVEMENT_MIN,
            "threshold_mae": float(seasonal_mae) * (1.0 - SUCCESS_IMPROVEMENT_MIN),
            "model_mae": point_mae,
            "passed": cleared if horizon == SUCCESS_HORIZON else None,
        },
        "_booster_point": point,
        "_q_preds": q_preds,
        "_y_test": y_test,
        "_x_test": x_test,
        "_x_train": x_train,
        "_y_train": y_train,
        "_x_val": x_val,
        "_y_val": y_val,
    }


def _public_result(result: dict[str, object]) -> dict[str, object]:
    return {k: v for k, v in result.items() if not k.startswith("_")}


def load_manifest_hash(export_dir: Path) -> str | None:
    path = export_dir / "manifest.json"
    if not path.is_file():
        return None
    manifest = json.loads(path.read_text(encoding="utf-8"))
    return manifest.get("schema_hash") or manifest.get("content_hash")


def write_registry(
    out_dir: Path,
    report: dict[str, object],
    manifest_hash: str | None,
) -> Path:
    out_dir.mkdir(parents=True, exist_ok=True)
    entry = {
        "model_id": "soy_lgbm_walkforward_v1",
        "dataset": "ml.soy-daily-training",
        "manifest_schema_hash": manifest_hash,
        "feature_cols": list(FEATURE_COLS),
        "horizons": list(HORIZONS),
        "quantiles": list(QUANTILES),
        "params": {
            "point_rounds": POINT_ROUNDS,
            "quantile_rounds": QUANTILE_ROUNDS,
            "early_stopping": EARLY_STOP,
            "learning_rate": 0.05,
            "num_leaves": 31,
            "seed": 42,
        },
        "metrics": report,
    }
    path = out_dir / "soy_lgbm_v1.json"
    path.write_text(json.dumps(entry, indent=2, sort_keys=True) + "\n", encoding="utf-8")
    return path


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument(
        "--export-dir",
        default=os.environ.get("ML_EXPORT_DIR", str(ROOT / ".local" / "ml" / "export")),
    )
    parser.add_argument(
        "--out-dir",
        default=os.environ.get("ML_TRAIN_DIR", str(ROOT / ".local" / "ml" / "train")),
    )
    parser.add_argument(
        "--registry-dir",
        default=os.environ.get("ML_REGISTRY_DIR", str(ROOT / ".local" / "ml" / "registry")),
    )
    parser.add_argument(
        "--require-gate",
        action="store_true",
        help="Exit 1 when SUCCESS_HORIZON point MAE fails B3 gate",
    )
    args = parser.parse_args()

    lgb = _require_lightgbm()
    parquet_path = Path(args.export_dir) / "soy_daily_training.parquet"
    if not parquet_path.is_file():
        print(f"missing export parquet: {parquet_path}", file=sys.stderr)
        return 1

    rows = load_rows(parquet_path)
    results: list[dict[str, object]] = []
    gate_ok = False
    for h in HORIZONS:
        raw = train_horizon(lgb, rows, h)
        results.append(_public_result(raw))
        if h == SUCCESS_HORIZON:
            gate_ok = bool(raw["success_gate"]["passed"])  # type: ignore[index]

    report: dict[str, object] = {
        "dataset": "ml.soy-daily-training",
        "export_parquet": str(parquet_path),
        "success_horizon": SUCCESS_HORIZON,
        "b3_gate_passed": gate_ok,
        "tft_chronos": "skipped — MVP uses LightGBM only; foundation TS deferred",
        "results": results,
    }

    out_dir = Path(args.out_dir)
    out_dir.mkdir(parents=True, exist_ok=True)
    out_json = out_dir / "soy_lgbm_metrics.json"
    out_json.write_text(json.dumps(report, indent=2, sort_keys=True) + "\n", encoding="utf-8")

    md = [
        "# Soy LightGBM walk-forward metrics",
        "",
        f"Source: `{parquet_path}`",
        "",
        f"**B3 gate (h={SUCCESS_HORIZON}):** {'PASS' if gate_ok else 'FAIL'}",
        "",
        "| h | n_test | point MAE | seasonal MAE | gate |",
        "|---|--------|-----------|--------------|------|",
    ]
    for r in results:
        sg = r["success_gate"]  # type: ignore[index]
        base = r["baseline_test"]  # type: ignore[index]
        point = r["point"]  # type: ignore[index]
        gate_cell = (
            ("PASS" if sg["passed"] else "FAIL") if sg["applies"] else "—"
        )
        md.append(
            f"| {r['horizon_trading_days']} | {r['n_test']} | {point['test_mae']} | "
            f"{base['seasonal_naive_mae']} | {gate_cell} |"
        )
    md.extend(["", "## Quantiles (test MAE)", ""])
    for r in results:
        qs = r["quantiles"]  # type: ignore[index]
        md.append(
            f"- h={r['horizon_trading_days']}: "
            + ", ".join(f"{k}={qs[k]['test_mae']}" for k in sorted(qs))
        )
    (out_dir / "soy_lgbm_metrics.md").write_text("\n".join(md) + "\n", encoding="utf-8")

    reg_path = write_registry(Path(args.registry_dir), report, load_manifest_hash(Path(args.export_dir)))
    print(f"wrote {out_json}")
    print(f"wrote {reg_path}")
    print(f"b3_gate_passed={gate_ok}")

    if args.require_gate and not gate_ok:
        return 1
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
