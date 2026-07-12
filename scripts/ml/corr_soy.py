#!/usr/bin/env python3
"""Lagged correlation report for soy daily features vs forward return (Phase 30 B2)."""

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

from ml.config import CORR_FEATURE_CANDIDATES, CORR_LAGS, CORR_TARGET  # noqa: E402


def _pearson(xs: list[float], ys: list[float]) -> float | None:
    n = len(xs)
    if n < 5:
        return None
    mx = sum(xs) / n
    my = sum(ys) / n
    num = sum((x - mx) * (y - my) for x, y in zip(xs, ys, strict=True))
    den_x = math.sqrt(sum((x - mx) ** 2 for x in xs))
    den_y = math.sqrt(sum((y - my) ** 2 for y in ys))
    if den_x == 0 or den_y == 0:
        return None
    return num / (den_x * den_y)


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument(
        "--export-dir",
        default=os.environ.get("ML_EXPORT_DIR", str(ROOT / ".local" / "ml" / "export")),
    )
    parser.add_argument(
        "--out-dir",
        default=os.environ.get("ML_CORR_DIR", str(ROOT / ".local" / "ml" / "corr")),
    )
    args = parser.parse_args()

    parquet_path = Path(args.export_dir) / "soy_daily_training.parquet"
    if not parquet_path.is_file():
        print(f"missing export parquet: {parquet_path}", file=sys.stderr)
        return 1

    rows = sorted(pq.read_table(parquet_path).to_pylist(), key=lambda r: str(r.get("data")))
    n = len(rows)

    results: list[dict[str, object]] = []
    for feat in CORR_FEATURE_CANDIDATES:
        for lag in CORR_LAGS:
            xs: list[float] = []
            ys: list[float] = []
            for i in range(n):
                j = i - lag
                if j < 0:
                    continue
                y = rows[i].get(CORR_TARGET)
                x = rows[j].get(feat)
                if y is None or x is None:
                    continue
                xs.append(float(x))
                ys.append(float(y))
            results.append(
                {
                    "feature": feat,
                    "lag_trading_days": lag,
                    "target": CORR_TARGET,
                    "n": len(xs),
                    "pearson": _pearson(xs, ys),
                }
            )

    # Rank by |corr| descending among non-null
    ranked = sorted(
        [r for r in results if r["pearson"] is not None],
        key=lambda r: abs(float(r["pearson"])),  # type: ignore[arg-type]
        reverse=True,
    )

    report = {
        "dataset": "ml.soy-daily-training",
        "target": CORR_TARGET,
        "lags": list(CORR_LAGS),
        "features": list(CORR_FEATURE_CANDIDATES),
        "pairs": results,
        "top_abs_corr": ranked[:20],
    }

    out_dir = Path(args.out_dir)
    out_dir.mkdir(parents=True, exist_ok=True)
    out_json = out_dir / "soy_daily_corr.json"
    out_json.write_text(json.dumps(report, indent=2, sort_keys=True) + "\n", encoding="utf-8")

    md = [
        "# Soy daily lagged correlations",
        "",
        f"Target: `{CORR_TARGET}`",
        "",
        "| feature | lag | n | pearson |",
        "|---------|-----|---|---------|",
    ]
    for r in ranked[:30]:
        md.append(
            f"| {r['feature']} | {r['lag_trading_days']} | {r['n']} | {r['pearson']:.4f} |"
        )
    (out_dir / "soy_daily_corr.md").write_text("\n".join(md) + "\n", encoding="utf-8")
    print(f"wrote {out_json}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
