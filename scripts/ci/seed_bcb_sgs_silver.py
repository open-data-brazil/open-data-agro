#!/usr/bin/env python3
"""Seed minimal silver Delta for BCB SGS CI (local-first)."""

from __future__ import annotations

import os
import sys
from pathlib import Path

import pyarrow as pa
from deltalake import write_deltalake


def write_table(root: Path, agency: str, table: str, data: pa.Table) -> None:
    path = root / "silver" / agency / table
    path.parent.mkdir(parents=True, exist_ok=True)
    write_deltalake(str(path), data, mode="overwrite")


def sgs_rows(
    code: str,
    dataset_id: str,
    dates: list[str],
    valores: list[str],
    source: str,
    ingested: str,
) -> pa.Table:
    return pa.table(
        {
            "sgs_codigo": [code] * len(dates),
            "data": dates,
            "valor": valores,
            "ano": [d[:4] for d in dates],
            "_dataset_id": [dataset_id] * len(dates),
            "_ingested_at": [ingested] * len(dates),
            "_source_file": [source] * len(dates),
        }
    )


def main() -> int:
    lake_root = Path(os.environ.get("LAKE_LOCAL_ROOT", "/tmp/open-data-agro-lake"))
    lake_root.mkdir(parents=True, exist_ok=True)
    (lake_root / "gold").mkdir(parents=True, exist_ok=True)
    for mart in [
        "mart_bcb__sgs_ipca",
        "mart_bcb__sgs_ipca_12m",
        "mart_bcb__sgs_igpm",
        "mart_bcb__sgs_ptax_usd_venda",
        "mart_bcb__sgs_ptax_usd_compra",
    ]:
        (lake_root / "gold" / mart).mkdir(parents=True, exist_ok=True)

    source = str(lake_root / "bronze/seed.parquet")
    ingested = "2026-06-25T12:00:00Z"

    tables = [
        ("sgs_ipca", sgs_rows("433", "bcb.sgs-ipca", ["2020-01-01", "2020-02-01"], ["0.21", "0.25"], source, ingested)),
        ("sgs_ipca_12m", sgs_rows("13522", "bcb.sgs-ipca-12m", ["2020-01-01", "2020-02-01"], ["4.31", "4.56"], source, ingested)),
        ("sgs_igpm", sgs_rows("189", "bcb.sgs-igpm", ["2020-01-01", "2020-02-01"], ["0.89", "0.12"], source, ingested)),
        ("sgs_ptax_usd_venda", sgs_rows("1", "bcb.sgs-ptax-usd-venda", ["2024-01-02", "2024-01-03"], ["4.8916", "4.9212"], source, ingested)),
        ("sgs_ptax_usd_compra", sgs_rows("10813", "bcb.sgs-ptax-usd-compra", ["2024-01-02", "2024-01-03"], ["4.8900", "4.9200"], source, ingested)),
    ]
    for table, data in tables:
        write_table(lake_root, "bcb", table, data)

    print(f"seeded BCB SGS silver under {lake_root / 'silver' / 'bcb'}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
