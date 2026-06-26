#!/usr/bin/env python3
"""Assert minimum bronze row counts for IBGE localidades live ingest smoke."""

from __future__ import annotations

import argparse
import sys
from pathlib import Path

import pyarrow.parquet as pq

# Official IBGE localidades catalog sizes (2024 layout); allow small drift.
MIN_BRONZE_ROWS: dict[str, int] = {
    "localidades-municipios": 5500,
    "localidades-ufs": 27,
    "localidades-regioes": 5,
    "localidades-mesorregioes": 130,
    "localidades-microrregioes": 500,
}


def count_bronze_rows(lake_root: Path, dataset_dir: str) -> int:
    bronze_dir = lake_root / "bronze" / "ibge" / dataset_dir
    parquet_files = list(bronze_dir.rglob("*.parquet"))
    if not parquet_files:
        return 0
    return sum(pq.read_table(path).num_rows for path in parquet_files)


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument(
        "--lake-root",
        type=Path,
        required=True,
        help="Lake root (e.g. ./lake)",
    )
    args = parser.parse_args()
    lake_root = args.lake_root.resolve()

    failed = False
    for dataset_dir, min_rows in MIN_BRONZE_ROWS.items():
        rows = count_bronze_rows(lake_root, dataset_dir)
        status = "OK" if rows >= min_rows else "FAIL"
        print(f"{dataset_dir}: {rows} rows (min {min_rows}) [{status}]")
        if rows < min_rows:
            failed = True

    if failed:
        print("Bronze row check failed — run ingestor for missing datasets.", file=sys.stderr)
        return 1
    print("All IBGE localidades bronze datasets meet minimum row counts.")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
