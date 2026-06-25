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


def main() -> int:
    lake_root = Path(os.environ.get("LAKE_LOCAL_ROOT", "/tmp/open-data-agro-lake"))
    lake_root.mkdir(parents=True, exist_ok=True)

    source = str(lake_root / "bronze/seed.parquet")
    ingested = "2026-06-25T12:00:00Z"

    ipca = pa.table(
        {
            "sgs_codigo": ["433", "433"],
            "data": ["2020-01-01", "2020-02-01"],
            "valor": ["0.21", "0.25"],
            "ano": ["2020", "2020"],
            "_dataset_id": ["bcb.sgs-ipca", "bcb.sgs-ipca"],
            "_ingested_at": [ingested, ingested],
            "_source_file": [source, source],
        }
    )
    write_table(lake_root, "bcb", "sgs_ipca", ipca)

    ptax = pa.table(
        {
            "sgs_codigo": ["1", "1"],
            "data": ["2024-01-02", "2024-01-03"],
            "valor": ["4.8916", "4.9212"],
            "ano": ["2024", "2024"],
            "_dataset_id": ["bcb.sgs-ptax-usd-venda", "bcb.sgs-ptax-usd-venda"],
            "_ingested_at": [ingested, ingested],
            "_source_file": [source, source],
        }
    )
    write_table(lake_root, "bcb", "sgs_ptax_usd_venda", ptax)

    print(f"seeded BCB SGS silver under {lake_root / 'silver' / 'bcb'}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
