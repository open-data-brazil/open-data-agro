#!/usr/bin/env python3
"""Seed minimal silver Delta for USDA PSD CI (local-first)."""

from __future__ import annotations

import os
import sys
from pathlib import Path

import pyarrow as pa
from deltalake import write_deltalake


def write_table(root: Path, table: str, data: pa.Table) -> None:
    path = root / "silver" / "usda" / table
    path.parent.mkdir(parents=True, exist_ok=True)
    write_deltalake(str(path), data, mode="overwrite")


def seed_dataset(
    root: Path,
    table: str,
    dataset_id: str,
    commodity_slug: str,
    commodity_code: str,
    commodity_name: str,
    countries: list[tuple[str, str]],
    attribute_id: str,
    attribute_name: str,
    values: list[str],
) -> None:
    source = str(root / "bronze/seed.parquet")
    ingested = "2026-06-26T12:00:00Z"
    n = len(countries)
    data = pa.table(
        {
            "commodity_code": [commodity_code] * n,
            "commodity_name": [commodity_name] * n,
            "commodity_slug": [commodity_slug] * n,
            "country_code": [c[0] for c in countries],
            "country_name": [c[1] for c in countries],
            "marketing_year": ["2024"] * n,
            "calendar_year": ["2024"] * n,
            "month": ["05"] * n,
            "attribute_id": [attribute_id] * n,
            "attribute_name": [attribute_name] * n,
            "unit_id": ["8"] * n,
            "unit_description": ["(1000 MT)"] * n,
            "value": values,
            "_dataset_id": [dataset_id] * n,
            "_ingested_at": [ingested] * n,
            "_source_file": [source] * n,
        }
    )
    write_table(root, table, data)


def main() -> int:
    lake_root = Path(os.environ.get("LAKE_LOCAL_ROOT", "/tmp/open-data-agro-lake"))
    lake_root.mkdir(parents=True, exist_ok=True)
    for mart in (
        "mart_usda__psd_soja",
        "mart_usda__psd_milho",
        "mart_usda__psd_trigo",
    ):
        (lake_root / "gold" / mart).mkdir(parents=True, exist_ok=True)

    seed_dataset(
        lake_root,
        "psd_soja",
        "usda.psd-soja",
        "soja",
        "2222000",
        "Oilseed, Soybean",
        [("BR", "Brazil"), ("US", "United States")],
        "28",
        "Production",
        ["165000.0000", "121000.0000"],
    )
    seed_dataset(
        lake_root,
        "psd_milho",
        "usda.psd-milho",
        "milho",
        "0440000",
        "Corn",
        [("BR", "Brazil"), ("US", "United States")],
        "28",
        "Production",
        ["125000.0000", "384000.0000"],
    )
    seed_dataset(
        lake_root,
        "psd_trigo",
        "usda.psd-trigo",
        "trigo",
        "0410000",
        "Wheat",
        [("BR", "Brazil"), ("AR", "Argentina")],
        "28",
        "Production",
        ["9500.0000", "15500.0000"],
    )

    print(f"seeded USDA PSD silver under {lake_root / 'silver' / 'usda'}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
