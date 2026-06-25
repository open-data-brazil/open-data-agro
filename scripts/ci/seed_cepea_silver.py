#!/usr/bin/env python3
"""Seed minimal silver Delta for CEPEA indicators CI (local-first)."""

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

    soja = pa.table(
        {
            "produto": ["soja", "soja"],
            "praca": ["Paranaguá", "Paranaguá"],
            "data": ["2010-01-04", "2024-01-31"],
            "preco_rs_sc": ["52.30", "124.58"],
            "variacao_dia_pct": ["0.15", "0.32"],
            "preco_usd_sc": ["29.80", "25.82"],
            "ano": ["2010", "2024"],
            "_dataset_id": ["cepea.soja-paranagua", "cepea.soja-paranagua"],
            "_ingested_at": [ingested, ingested],
            "_source_file": [source, source],
        }
    )
    write_table(lake_root, "cepea", "soja_paranagua", soja)

    print(f"seeded CEPEA silver under {lake_root / 'silver' / 'cepea'}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
