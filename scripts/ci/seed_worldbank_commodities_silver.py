#!/usr/bin/env python3
"""Seed minimal silver Delta for World Bank Pink Sheet CI (local-first)."""

from __future__ import annotations

import os
import sys
from pathlib import Path

import pyarrow as pa
from deltalake import write_deltalake


def main() -> int:
    lake_root = Path(os.environ.get("LAKE_LOCAL_ROOT", "/tmp/open-data-agro-lake"))
    lake_root.mkdir(parents=True, exist_ok=True)
    (lake_root / "gold" / "mart_worldbank__pink_sheet_monthly").mkdir(parents=True, exist_ok=True)

    source = str(lake_root / "bronze/seed.parquet")
    ingested = "2026-06-26T12:00:00Z"
    data = pa.table(
        {
            "refmonth": ["2024-01", "2024-01", "2024-01"],
            "series_name": ["Soybeans", "Maize", "Crude oil, average"],
            "commodity_slug": ["soja", "milho", "petroleo"],
            "unit": ["($/mt)", "($/mt)", "($/bbl)"],
            "value": ["490.5", "210.2", "80.1"],
            "_dataset_id": ["worldbank.pink-sheet-monthly"] * 3,
            "_ingested_at": [ingested] * 3,
            "_source_file": [source] * 3,
        }
    )

    path = lake_root / "silver" / "worldbank" / "pink_sheet_monthly"
    path.parent.mkdir(parents=True, exist_ok=True)
    write_deltalake(str(path), data, mode="overwrite")

    print(f"seeded World Bank Pink Sheet silver under {lake_root / 'silver' / 'worldbank'}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
