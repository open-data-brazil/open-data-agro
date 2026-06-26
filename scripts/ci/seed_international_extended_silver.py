#!/usr/bin/env python3
"""Seed minimal silver Delta for international extended CI (local-first)."""

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
    for mart in (
        "mart_fao__producao_agro",
        "mart_fao__comercio_agro",
        "mart_worldbank__ag_indices",
    ):
        (lake_root / "gold" / mart).mkdir(parents=True, exist_ok=True)

    source = str(lake_root / "bronze/seed.parquet")
    ingested = "2026-06-26T12:00:00Z"

    producao = pa.table(
        {
            "area_code": ["9", "21"],
            "area_name": ["Argentina", "Brazil"],
            "item_code": ["56", "236"],
            "item_name": ["Maize (corn)", "Soya beans"],
            "commodity_slug": ["milho", "soja"],
            "element_code": ["5510", "5510"],
            "element_name": ["Production", "Production"],
            "year": ["2023", "2023"],
            "unit": ["t", "t"],
            "value": ["41409448.000000", "154000000.000000"],
            "flag": ["A", "A"],
            "_dataset_id": ["fao.producao-agro", "fao.producao-agro"],
            "_ingested_at": [ingested, ingested],
            "_source_file": [source, source],
        }
    )
    write_table(lake_root, "fao", "producao_agro", producao)

    comercio = pa.table(
        {
            "area_code": ["21", "231"],
            "area_name": ["Brazil", "United States of America"],
            "item_code": ["236", "56"],
            "item_name": ["Soya beans", "Maize (corn)"],
            "commodity_slug": ["soja", "milho"],
            "element_code": ["5922", "5922"],
            "element_name": ["Export quantity", "Export quantity"],
            "year": ["2023", "2023"],
            "unit": ["t", "t"],
            "value": ["98000000.000000", "52000000.000000"],
            "flag": ["A", "A"],
            "_dataset_id": ["fao.comercio-agro", "fao.comercio-agro"],
            "_ingested_at": [ingested, ingested],
            "_source_file": [source, source],
        }
    )
    write_table(lake_root, "fao", "comercio_agro", comercio)

    ag_indices = pa.table(
        {
            "refmonth": ["2024-01", "2024-01"],
            "series_name": ["Agriculture **", "Grains"],
            "commodity_slug": ["agriculture", "grains"],
            "unit": ["Index", "Index"],
            "value": ["118.5", "115.8"],
            "_dataset_id": ["worldbank.ag-indices", "worldbank.ag-indices"],
            "_ingested_at": [ingested, ingested],
            "_source_file": [source, source],
        }
    )
    write_table(lake_root, "worldbank", "ag_indices", ag_indices)

    print(f"seeded international extended silver under {lake_root / 'silver'}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
