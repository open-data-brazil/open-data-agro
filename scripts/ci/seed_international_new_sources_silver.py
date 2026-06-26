#!/usr/bin/env python3
"""Seed minimal silver Delta for international new sources CI (Phase 38)."""

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
    (lake_root / "gold" / "mart_eia__petroleum_prices").mkdir(parents=True, exist_ok=True)

    source = str(lake_root / "bronze/seed.parquet")
    ingested = "2026-06-26T12:00:00Z"

    petroleum = pa.table(
        {
            "series_id": ["PET.RWTC.D", "PET.RBRTE.D"],
            "series_name": ["WTI Cushing OK Spot", "Europe Brent Spot"],
            "commodity_slug": ["wti_spot", "brent_spot"],
            "refdate": ["2024-01-02", "2024-01-02"],
            "unit": ["dollars per barrel", "dollars per barrel"],
            "value": ["70.62", "76.69"],
            "frequency": ["daily", "daily"],
            "_dataset_id": ["eia.petroleum-prices", "eia.petroleum-prices"],
            "_ingested_at": [ingested, ingested],
            "_source_file": [source, source],
        }
    )
    write_table(lake_root, "eia", "petroleum_prices", petroleum)

    print(f"seeded international new sources silver under {lake_root / 'silver'}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
