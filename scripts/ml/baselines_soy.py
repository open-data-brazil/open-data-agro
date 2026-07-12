#!/usr/bin/env python3
"""Naive baselines for soy daily targets (Phase 30 B1)."""

from __future__ import annotations

import argparse
import json
import math
import os
import sys
from pathlib import Path

import pyarrow.parquet as pq

ROOT = Path(__file__).resolve().parents[2]
sys.path.insert(0, str(ROOT / "scripts"))

from ml.config import HORIZONS, SUCCESS_HORIZON, SUCCESS_IMPROVEMENT_MIN  # noqa: E402


def _mae(y_true: list[float], y_pred: list[float]) -> float | None:
    if not y_true:
        return None
    return sum(abs(a - b) for a, b in zip(y_true, y_pred, strict=True)) / len(y_true)


def _rmse(y_true: list[float], y_pred: list[float]) -> float | None:
    if not y_true:
        return None
    return math.sqrt(sum((a - b) ** 2 for a, b in zip(y_true, y_pred, strict=True)) / len(y_true))


def evaluate_horizon(
    rows: list[dict[str, object]],
    horizon: int,
    split_name: str,
) -> dict[str, object]:
    """Last-value and seasonal-naive baselines on y_ret / y_level."""
    ret_key = f"y_ret_{horizon}d"
    level_key = f"y_level_{horizon}d"

    # Index by split for evaluation; need full series ordered by date for preds.
    ordered = sorted(rows, key=lambda r: str(r.get("data")))

    # Seasonal lag ≈ 252 trading days; fall back to horizon lag when short history.
    seasonal_lag = 252

    y_ret_true: list[float] = []
    y_ret_last: list[float] = []
    y_ret_seasonal: list[float] = []
    y_level_true: list[float] = []
    y_level_last: list[float] = []
    y_level_seasonal: list[float] = []

    for i, row in enumerate(ordered):
        if row.get("split") != split_name:
            continue
        y_ret = row.get(ret_key)
        y_level = row.get(level_key)
        price = row.get("preco_rs_sc")
        if y_ret is None or y_level is None or price is None:
            continue

        # Last-value: predict 0 return; predict current price for level.
        pred_ret_last = 0.0
        pred_level_last = float(price)

        # Seasonal: same return/level change as seasonal_lag ago; else lag=horizon.
        j = i - seasonal_lag
        if j < 0:
            j = i - horizon
        if j >= 0:
            past = ordered[j]
            past_ret = past.get(ret_key)
            past_level = past.get("preco_rs_sc")
            pred_ret_seasonal = float(past_ret) if past_ret is not None else 0.0
            pred_level_seasonal = float(past_level) if past_level is not None else pred_level_last
        else:
            pred_ret_seasonal = 0.0
            pred_level_seasonal = pred_level_last

        y_ret_true.append(float(y_ret))
        y_ret_last.append(pred_ret_last)
        y_ret_seasonal.append(pred_ret_seasonal)
        y_level_true.append(float(y_level))
        y_level_last.append(pred_level_last)
        y_level_seasonal.append(pred_level_seasonal)

    n = len(y_ret_true)
    return {
        "split": split_name,
        "horizon_trading_days": horizon,
        "n": n,
        "y_ret": {
            "last_value": {"mae": _mae(y_ret_true, y_ret_last), "rmse": _rmse(y_ret_true, y_ret_last)},
            "seasonal_naive": {
                "mae": _mae(y_ret_true, y_ret_seasonal),
                "rmse": _rmse(y_ret_true, y_ret_seasonal),
            },
        },
        "y_level": {
            "last_value": {
                "mae": _mae(y_level_true, y_level_last),
                "rmse": _rmse(y_level_true, y_level_last),
            },
            "seasonal_naive": {
                "mae": _mae(y_level_true, y_level_seasonal),
                "rmse": _rmse(y_level_true, y_level_seasonal),
            },
        },
    }


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument(
        "--export-dir",
        default=os.environ.get("ML_EXPORT_DIR", str(ROOT / ".local" / "ml" / "export")),
    )
    parser.add_argument(
        "--out-dir",
        default=os.environ.get("ML_BASELINES_DIR", str(ROOT / ".local" / "ml" / "baselines")),
    )
    args = parser.parse_args()

    parquet_path = Path(args.export_dir) / "soy_daily_training.parquet"
    if not parquet_path.is_file():
        print(f"missing export parquet: {parquet_path}", file=sys.stderr)
        return 1

    table = pq.read_table(parquet_path)
    rows = table.to_pylist()

    report: dict[str, object] = {
        "dataset": "ml.soy-daily-training",
        "success_gate": {
            "horizon": SUCCESS_HORIZON,
            "metric": "mae",
            "target": f"y_ret_{SUCCESS_HORIZON}d",
            "baseline": "seasonal_naive",
            "split": "test",
            "min_relative_improvement": SUCCESS_IMPROVEMENT_MIN,
            "rule": (
                "Phase 31 model MAE on test must be ≤ "
                f"(1 - {SUCCESS_IMPROVEMENT_MIN}) * seasonal_naive MAE "
                f"for y_ret_{SUCCESS_HORIZON}d"
            ),
        },
        "results": [],
    }

    for split_name in ("val", "test"):
        for h in HORIZONS:
            report["results"].append(evaluate_horizon(rows, h, split_name))

    out_dir = Path(args.out_dir)
    out_dir.mkdir(parents=True, exist_ok=True)
    out_path = out_dir / "soy_daily_baselines.json"
    out_path.write_text(json.dumps(report, indent=2, sort_keys=True) + "\n", encoding="utf-8")

    # Human-readable summary
    md_lines = [
        "# Soy daily baselines",
        "",
        f"Source: `{parquet_path}`",
        "",
        "## Success gate (B3)",
        "",
        report["success_gate"]["rule"],  # type: ignore[index]
        "",
        "| split | h | n | last MAE (ret) | seasonal MAE (ret) | last MAE (level) | seasonal MAE (level) |",
        "|-------|---|---|----------------|--------------------|------------------|----------------------|",
    ]
    for r in report["results"]:  # type: ignore[attr-defined]
        y_ret = r["y_ret"]  # type: ignore[index]
        y_level = r["y_level"]  # type: ignore[index]
        md_lines.append(
            f"| {r['split']} | {r['horizon_trading_days']} | {r['n']} | "
            f"{y_ret['last_value']['mae']} | {y_ret['seasonal_naive']['mae']} | "
            f"{y_level['last_value']['mae']} | {y_level['seasonal_naive']['mae']} |"
        )
    (out_dir / "soy_daily_baselines.md").write_text("\n".join(md_lines) + "\n", encoding="utf-8")
    print(f"wrote {out_path}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
