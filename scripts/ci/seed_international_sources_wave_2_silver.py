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
    gold_tables = [
        "mart_igc__goi_index",
        "mart_eurostat__ag_prices",
        "mart_argentina__bcra_cambio",
    ]
    for table in gold_tables:
        (lake_root / "gold" / table).mkdir(parents=True, exist_ok=True)

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

    eurostat = pa.table(
        {
            "dataset_code": ["apri_pi15_outa", "apri_pi15_outa"],
            "geo": ["EU27_2020", "EU27_2020"],
            "product_code": ["010000", "015000"],
            "product_name": ["Cereals (including seeds)", "Grain maize"],
            "year": ["2022", "2023"],
            "index_value": ["135.72", "148.31"],
            "base_period": ["2015=100", "2015=100"],
            "_dataset_id": ["eurostat.ag-prices", "eurostat.ag-prices"],
            "_ingested_at": [ingested, ingested],
            "_source_file": [source, source],
        }
    )
    write_table(lake_root, "eurostat", "ag_prices", eurostat)

    bcra = pa.table(
        {
            "currency_code": ["USD", "USD"],
            "currency_name": ["DOLAR E.E.U.U.", "DOLAR E.E.U.U."],
            "refdate": ["2026-06-24", "2026-06-25"],
            "exchange_rate": ["1479.00000000", "1477.00000000"],
            "rate_type": ["tipo_cotizacion", "tipo_cotizacion"],
            "_dataset_id": ["argentina.bcra-cambio", "argentina.bcra-cambio"],
            "_ingested_at": [ingested, ingested],
            "_source_file": [source, source],
        }
    )
    write_table(lake_root, "argentina", "bcra_cambio", bcra)

    print(f"seeded international sources wave 2 silver under {lake_root / 'silver'}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
