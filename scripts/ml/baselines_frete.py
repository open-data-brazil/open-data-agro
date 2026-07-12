#!/usr/bin/env python3
"""Naive baselines for frete target (Phase 32 F2)."""

from __future__ import annotations

import argparse
import json
import os
import sys
from pathlib import Path

import pyarrow.parquet as pq

ROOT = Path(__file__).resolve().parents[2]
sys.path.insert(0, str(ROOT / "scripts"))

from ml.config_frete import SUCCESS_IMPROVEMENT_MIN, TARGET_COL  # noqa: E402
from ml.metrics import mae, rmse  # noqa: E402


def evaluate_split(rows: list[dict[str, object]], split_name: str) -> dict[str, object]:
    ordered = sorted(rows, key=lambda r: (int(r["ano"]), int(r["mes"])))  # type: ignore[arg-type]
    y_true: list[float] = []
    pred_global: list[float] = []
    pred_corridor: list[float] = []
    pred_seasonal: list[float] = []

    # Running stats for causal baselines (train-time only leakage-free on test via history).
    global_sum = 0.0
    global_n = 0
    corridor_sum: dict[str, float] = {}
    corridor_n: dict[str, int] = {}
    seasonal_sum: dict[tuple[str, int], float] = {}  # corridor, month
    seasonal_n: dict[tuple[str, int], int] = {}

    for row in ordered:
        key = f"{row['uf_origem']}->{row['uf_destino']}"
        ton = float(row[TARGET_COL])  # type: ignore[arg-type]
        mes = int(row["mes"])  # type: ignore[arg-type]
        season_key = (key, mes)

        if row.get("split") == split_name:
            # Predictions from history before this row is added.
            g = (global_sum / global_n) if global_n else ton
            c = (corridor_sum[key] / corridor_n[key]) if corridor_n.get(key) else g
            s = (
                (seasonal_sum[season_key] / seasonal_n[season_key])
                if seasonal_n.get(season_key)
                else c
            )
            y_true.append(ton)
            pred_global.append(g)
            pred_corridor.append(c)
            pred_seasonal.append(s)

        global_sum += ton
        global_n += 1
        corridor_sum[key] = corridor_sum.get(key, 0.0) + ton
        corridor_n[key] = corridor_n.get(key, 0) + 1
        seasonal_sum[season_key] = seasonal_sum.get(season_key, 0.0) + ton
        seasonal_n[season_key] = seasonal_n.get(season_key, 0) + 1

    return {
        "split": split_name,
        "n": len(y_true),
        "global_mean": {"mae": mae(y_true, pred_global), "rmse": rmse(y_true, pred_global)},
        "corridor_mean": {"mae": mae(y_true, pred_corridor), "rmse": rmse(y_true, pred_corridor)},
        "seasonal_naive": {"mae": mae(y_true, pred_seasonal), "rmse": rmse(y_true, pred_seasonal)},
    }


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
        default=os.environ.get(
            "ML_FRETE_BASELINES_DIR", str(ROOT / ".local" / "ml" / "baselines_frete")
        ),
    )
    args = parser.parse_args()

    parquet_path = Path(args.export_dir) / "frete_training.parquet"
    if not parquet_path.is_file():
        print(f"missing export: {parquet_path}", file=sys.stderr)
        return 1

    rows = pq.read_table(parquet_path).to_pylist()
    report = {
        "dataset": "ml.frete-training",
        "target": TARGET_COL,
        "success_gate": {
            "metric": "mae",
            "baseline": "seasonal_naive",
            "split": "test",
            "min_relative_improvement": SUCCESS_IMPROVEMENT_MIN,
            "rule": (
                f"Phase 32 frete model MAE on test must be ≤ "
                f"(1 - {SUCCESS_IMPROVEMENT_MIN}) * seasonal_naive MAE for {TARGET_COL}"
            ),
        },
        "results": [evaluate_split(rows, s) for s in ("val", "test")],
    }

    out_dir = Path(args.out_dir)
    out_dir.mkdir(parents=True, exist_ok=True)
    out_json = out_dir / "frete_baselines.json"
    out_json.write_text(json.dumps(report, indent=2, sort_keys=True) + "\n", encoding="utf-8")

    md = [
        "# Frete baselines",
        "",
        f"Source: `{parquet_path}`",
        "",
        report["success_gate"]["rule"],  # type: ignore[index]
        "",
        "| split | n | global MAE | corridor MAE | seasonal MAE |",
        "|-------|---|------------|--------------|--------------|",
    ]
    for r in report["results"]:  # type: ignore[attr-defined]
        md.append(
            f"| {r['split']} | {r['n']} | {r['global_mean']['mae']} | "
            f"{r['corridor_mean']['mae']} | {r['seasonal_naive']['mae']} |"
        )
    (out_dir / "frete_baselines.md").write_text("\n".join(md) + "\n", encoding="utf-8")
    print(f"wrote {out_json}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
