#!/usr/bin/env python3
"""Seed minimal silver Delta for FAO FAOSTAT prices CI (local-first)."""

from __future__ import annotations

import os
import sys
from pathlib import Path

import pyarrow as pa
from deltalake import write_deltalake


def main() -> int:
    lake_root = Path(os.environ.get("LAKE_LOCAL_ROOT", "/tmp/open-data-agro-lake"))
    lake_root.mkdir(parents=True, exist_ok=True)
    (lake_root / "gold" / "mart_fao__prices_agro").mkdir(parents=True, exist_ok=True)

    source = str(lake_root / "bronze/seed.parquet")
    ingested = "2026-06-26T12:00:00Z"
    data = pa.table(
        {
            "area_code": ["21", "21", "231"],
            "area_name": ["Brazil", "Brazil", "United States of America"],
            "item_code": ["236", "56", "15"],
            "item_name": ["Soya beans", "Maize (corn)", "Wheat"],
            "commodity_slug": ["soja", "milho", "trigo"],
            "element_code": ["5532", "5532", "5539"],
            "element_name": [
                "Producer Price (USD/tonne)",
                "Producer Price (USD/tonne)",
                "Producer Price Index (2014-2016 = 100)",
            ],
            "year": ["2023", "2023", "2023"],
            "months_code": ["7021", "7021", "7021"],
            "months": ["Annual value", "Annual value", "Annual value"],
            "unit": ["USD", "USD", "Index"],
            "value": ["450.000000", "180.000000", "125.000000"],
            "flag": ["E", "E", "E"],
            "_dataset_id": ["fao.prices-agro"] * 3,
            "_ingested_at": [ingested] * 3,
            "_source_file": [source] * 3,
        }
    )

    path = lake_root / "silver" / "fao" / "prices_agro"
    path.parent.mkdir(parents=True, exist_ok=True)
    write_deltalake(str(path), data, mode="overwrite")

    print(f"seeded FAO prices silver under {lake_root / 'silver' / 'fao'}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
