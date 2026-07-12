#!/usr/bin/env python3
"""Seed multi-year mart_ml__soy_daily_features for Phase 30 CI."""

from __future__ import annotations

import math
import os
from datetime import date, timedelta
from pathlib import Path

import pyarrow as pa
import pyarrow.parquet as pq


def _business_days(start: date, end: date) -> list[date]:
    days: list[date] = []
    d = start
    while d <= end:
        if d.weekday() < 5:
            days.append(d)
        d += timedelta(days=1)
    return days


def main() -> int:
    lake_root = Path(os.environ.get("LAKE_LOCAL_ROOT", "/tmp/ml-dataset-ci-lake"))
    out = lake_root / "gold" / "mart_ml__soy_daily_features"
    out.mkdir(parents=True, exist_ok=True)

    days = _business_days(date(2015, 1, 5), date(2024, 12, 31))
    n = len(days)
    prices = [100.0 + 10.0 * math.sin(i / 40.0) + 0.02 * i for i in range(n)]
    ptax = [3.5 + 0.5 * math.sin(i / 90.0) + 0.0005 * i for i in range(n)]

    ret_1d: list[float | None] = [None]
    ret_7d: list[float | None] = [None] * min(7, n)
    for i in range(1, n):
        ret_1d.append(prices[i] / prices[i - 1] - 1.0)
    for i in range(7, n):
        ret_7d.append(prices[i] / prices[i - 7] - 1.0)
    while len(ret_7d) < n:
        ret_7d.append(None)

    vol_7d: list[float | None] = []
    for i in range(n):
        window = [ret_1d[j] for j in range(max(0, i - 6), i + 1) if ret_1d[j] is not None]
        if len(window) < 2:
            vol_7d.append(None)
        else:
            m = sum(window) / len(window)
            vol_7d.append(math.sqrt(sum((x - m) ** 2 for x in window) / (len(window) - 1)))

    table = pa.table(
        {
            "produto_slug": ["soja"] * n,
            "data": days,
            "label_date": days,
            "preco_rs_sc": prices,
            "preco_usd_sc": [p / x for p, x in zip(prices, ptax, strict=True)],
            "variacao_dia_pct": [None if r is None else r * 100.0 for r in ret_1d],
            "cepea_ret_1d": ret_1d,
            "cepea_ret_7d": ret_7d,
            "cepea_vol_7d": vol_7d,
            "ptax_usd_venda": ptax,
            "b3_front_price": [12.0 + 0.01 * math.sin(i / 30.0) for i in range(n)],
            "b3_front_symbol": ["SOYH24"] * n,
            "basis_cepea_usd_minus_b3": [
                (prices[i] / ptax[i]) - (12.0 + 0.01 * math.sin(i / 30.0)) for i in range(n)
            ],
            "comex_soja_fob_usd": [1.0e9 + 1.0e6 * math.sin(i / 60.0) for i in range(n)],
            "comex_soja_kg": [2.0e9 + 2.0e6 * math.sin(i / 60.0) for i in range(n)],
            "frete_mt_pr_ton_avg": [400.0 + 20.0 * math.sin(i / 100.0) for i in range(n)],
            "frete_mt_pr_tkm_avg": [0.2 + 0.01 * math.sin(i / 100.0) for i in range(n)],
            "inmet_national_avg": [20.0 + 5.0 * math.sin(i / 50.0) for i in range(n)],
            "as_of_cepea": days,
            "as_of_ptax": days,
            "as_of_b3": days,
            "as_of_comex": [d - timedelta(days=45) for d in days],
            "as_of_frete": [d - timedelta(days=14) for d in days],
            "as_of_inmet": [d - timedelta(days=2) for d in days],
            "feature_as_of_max": days,
        }
    )
    pq.write_table(table, out / "mart.parquet")
    print(f"seeded {n} rows under {out / 'mart.parquet'}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
