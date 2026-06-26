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
    for mart in (
        "mart_eia__petroleum_prices",
        "mart_usda__wasde",
        "mart_igc__goi_index",
        "mart_un__comtrade_bulk",
    ):
        (lake_root / "gold" / mart).mkdir(parents=True, exist_ok=True)

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

    wasde = pa.table(
        {
            "report_month": ["2026-06", "2026-06"],
            "commodity": ["Wheat", "Corn"],
            "market_year": ["2025/26 (Est.)", "2025/26 (Est.)"],
            "attribute": ["Production", "Exports"],
            "value": ["800.1", "190.2"],
            "unit": ["million metric tons", "million metric tons"],
            "_dataset_id": ["usda.wasde", "usda.wasde"],
            "_ingested_at": [ingested, ingested],
            "_source_file": [source, source],
        }
    )
    write_table(lake_root, "usda", "wasde", wasde)

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

    comtrade = pa.table(
        {
            "reporter_code": ["76", "76"],
            "reporter_desc": ["Brazil", "Brazil"],
            "partner_code": ["0", "0"],
            "partner_desc": ["World", "World"],
            "flow_code": ["X", "M"],
            "flow_desc": ["Export", "Import"],
            "period": ["2024", "2024"],
            "hs_code": ["1201", "1005"],
            "commodity_slug": ["soja", "milho"],
            "trade_value_usd": ["58990470", "1234567"],
            "netweight_kg": ["147973486", "987654"],
            "qty": ["147975000", "987654"],
            "qty_unit_abbr": ["kg", "kg"],
            "_dataset_id": ["un.comtrade-bulk", "un.comtrade-bulk"],
            "_ingested_at": [ingested, ingested],
            "_source_file": [source, source],
        }
    )
    write_table(lake_root, "un", "comtrade_bulk", comtrade)

    print(f"seeded international new sources silver under {lake_root / 'silver'}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
