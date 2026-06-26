#!/usr/bin/env python3
"""Seed minimal silver Delta for international sources wave 3 CI (Phase 45)."""

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
    gold_tables = [
        "mart_oecd__ag_outlook",
        "mart_fao__food_price_index",
        "mart_argentina__magyp_producion_granos",
    ]
    for table in gold_tables:
        (lake_root / "gold" / table).mkdir(parents=True, exist_ok=True)

    source = str(lake_root / "bronze/seed.parquet")
    ingested = "2026-06-26T12:00:00Z"

    outlook = pa.table(
        {
            "ref_area": ["BRA", "BRA"],
            "ref_area_name": ["Brazil", "Brazil"],
            "commodity_code": ["CPC_0141", "CPC_0142"],
            "commodity_name": ["Soya beans", "Maize (corn)"],
            "measure_code": ["QP", "QP"],
            "measure_name": ["Production", "Production"],
            "unit": ["T", "T"],
            "unit_mult": ["3", "3"],
            "year": ["2023", "2023"],
            "value": ["160177", "125000"],
            "obs_status": ["A", "A"],
            "_dataset_id": ["oecd-fao.ag-outlook", "oecd-fao.ag-outlook"],
            "_ingested_at": [ingested, ingested],
            "_source_file": [source, source],
        }
    )
    write_table(lake_root, "oecd", "ag_outlook", outlook)

    ffpi = pa.table(
        {
            "refmonth": ["2024-01-01", "2024-01-01"],
            "index_slug": ["food", "cereals"],
            "index_name": ["Food Price Index", "Cereals Price Index"],
            "value": ["118.0", "110.5"],
            "base_period": ["2002-2004=100", "2002-2004=100"],
            "_dataset_id": ["fao.food-price-index", "fao.food-price-index"],
            "_ingested_at": [ingested, ingested],
            "_source_file": [source, source],
        }
    )
    write_table(lake_root, "fao", "food_price_index", ffpi)

    granos = pa.table(
        {
            "series_id": ["AGRO_A_Soja_0003", "AGRO_A_Maiz_0003"],
            "commodity_slug": ["soja", "milho"],
            "refyear": ["2023", "2023"],
            "value": ["48000000", "35000000"],
            "unit": ["toneladas", "toneladas"],
            "source": ["MAGyP", "MAGyP"],
            "_dataset_id": ["argentina.magyp-producion-granos", "argentina.magyp-producion-granos"],
            "_ingested_at": [ingested, ingested],
            "_source_file": [source, source],
        }
    )
    write_table(lake_root, "argentina", "magyp_producion_granos", granos)

    print(f"seeded international sources wave 3 silver under {lake_root / 'silver'}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
