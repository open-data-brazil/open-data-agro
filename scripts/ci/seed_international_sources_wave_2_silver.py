#!/usr/bin/env python3
"""Seed minimal silver Delta for international sources wave 2 CI (Phase 41)."""

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
    (lake_root / "gold" / "mart_igc__goi_index").mkdir(parents=True, exist_ok=True)

    source = str(lake_root / "bronze/seed.parquet")
    ingested = "2026-06-26T12:00:00Z"

    goi = pa.table(
        {
            "refdate": ["2000-01-03", "2000-01-03"],
            "index_slug": ["goi", "wheat"],
            "index_name": ["IGC GOI", "Wheat"],
            "value": ["96.6658447381306", "96.95303997255832"],
            "base_period": ["2000-01=100", "2000-01=100"],
            "frequency": ["daily", "daily"],
            "_dataset_id": ["igc.goi-index", "igc.goi-index"],
            "_ingested_at": [ingested, ingested],
            "_source_file": [source, source],
        }
    )
    write_table(lake_root, "igc", "goi_index", goi)

    print(f"seeded international sources wave 2 silver under {lake_root / 'silver'}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
