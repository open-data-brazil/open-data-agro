#!/usr/bin/env python3
"""Seed minimal bronze Parquet for Great Expectations CI (local-first, no R2)."""

from __future__ import annotations

import os
import sys
from pathlib import Path

import pyarrow as pa
import pyarrow.parquet as pq


def main() -> int:
    lake_root = Path(os.environ.get("LAKE_LOCAL_ROOT", "/tmp/open-data-agro-lake"))
    bronze_dir = lake_root / "bronze" / "conab" / "estimativa-graos" / "ingest_date=2026-06-25"
    bronze_dir.mkdir(parents=True, exist_ok=True)

    table = pa.table(
        {
            "Produto": ["Soja", "Milho"],
            "UF": ["PR", "MT"],
            "Safra": ["2025/26", "2025/26"],
            "Região": ["Sul", "Centro-Oeste"],
            "Produção (mil t)": ["100", "120"],
        }
    )
    pq.write_table(table, bronze_dir / "part-seed.parquet")
    print(f"seeded bronze parquet under {bronze_dir}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
