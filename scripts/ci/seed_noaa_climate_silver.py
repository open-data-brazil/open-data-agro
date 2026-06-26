#!/usr/bin/env python3
"""Seed minimal silver Delta for NOAA climate indices CI (local-first)."""

from __future__ import annotations

import os
import sys
from pathlib import Path

import pyarrow as pa
from deltalake import write_deltalake


def main() -> int:
    lake_root = Path(os.environ.get("LAKE_LOCAL_ROOT", "/tmp/open-data-agro-lake"))
    lake_root.mkdir(parents=True, exist_ok=True)
    (lake_root / "gold" / "mart_noaa__enso_indices").mkdir(parents=True, exist_ok=True)
    (lake_root / "gold" / "mart_noaa__global_temp_anomaly").mkdir(parents=True, exist_ok=True)

    source = str(lake_root / "bronze/seed.parquet")
    ingested = "2026-06-26T12:00:00Z"

    enso = pa.table(
        {
            "refyear": ["2024", "2024"],
            "season_code": ["DJF", "JFM"],
            "sst_total": ["24.72", "24.85"],
            "anomaly": ["1.53", "1.61"],
            "index_name": ["oni", "oni"],
            "_dataset_id": ["noaa.enso-indices"] * 2,
            "_ingested_at": [ingested] * 2,
            "_source_file": [source] * 2,
        }
    )
    enso_path = lake_root / "silver" / "noaa" / "enso_indices"
    enso_path.parent.mkdir(parents=True, exist_ok=True)
    write_deltalake(str(enso_path), enso, mode="overwrite")

    temp = pa.table(
        {
            "refmonth": ["2024-01", "2024-02"],
            "anomaly": ["0.74", "0.81"],
            "unit": ["Degrees Celsius", "Degrees Celsius"],
            "base_period": ["1901-2000", "1901-2000"],
            "index_name": ["global_land_ocean", "global_land_ocean"],
            "_dataset_id": ["noaa.global-temp-anomaly"] * 2,
            "_ingested_at": [ingested] * 2,
            "_source_file": [source] * 2,
        }
    )
    temp_path = lake_root / "silver" / "noaa" / "global_temp_anomaly"
    write_deltalake(str(temp_path), temp, mode="overwrite")

    print(f"seeded NOAA climate silver under {lake_root / 'silver' / 'noaa'}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
