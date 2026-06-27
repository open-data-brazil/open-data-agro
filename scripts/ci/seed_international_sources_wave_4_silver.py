#!/usr/bin/env python3
"""Seed minimal silver Delta for international sources wave 4 CI (Phase 49)."""

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
        "mart_cftc__cot_agricultural_futures",
        "mart_jrc__mars_crop_yield",
        "mart_fao__giews_crop_prospects",
        "mart_fao__amis_market_monitor",
        "mart_sagis__grain_supply_statistics",
        "mart_japan__maff_ag_trade",
        "mart_fred__commodity_indexes",
        "mart_nasa__power_agroclimatology",
        "mart_copernicus__era5_agroclimate",
    ]
    for table in gold_tables:
        (lake_root / "gold" / table).mkdir(parents=True, exist_ok=True)

    source = str(lake_root / "bronze/seed.parquet")
    ingested = "2026-06-26T12:00:00Z"
    meta = {"_ingested_at": [ingested], "_source_file": [source]}

    write_table(
        lake_root,
        "cftc",
        "cot_agricultural_futures",
        pa.table(
            {
                "report_date": ["2024-12-31"],
                "commodity_name": ["CORN"],
                "commodity_slug": ["milho"],
                "market_name": ["CORN - CHICAGO BOARD OF TRADE"],
                "open_interest_all": ["1347894"],
                "m_money_long": ["199417"],
                "m_money_short": ["81383"],
                "prod_merc_long": ["363396"],
                "prod_merc_short": ["762273"],
                "commodity_group": ["AGRICULTURE"],
                "futonly_or_combined": ["FutOnly"],
                "_dataset_id": ["cftc.cot-agricultural-futures"],
                **meta,
            }
        ),
    )

    write_table(
        lake_root,
        "jrc",
        "mars_crop_yield",
        pa.table(
            {
                "country": ["Colombia"],
                "crop": ["Maize (corn)"],
                "crop_slug": ["milho"],
                "forecast_yield_kg_ha": ["3855.56"],
                "five_yr_avg_kg_ha": ["3896.52"],
                "harvest_year": ["2024"],
                "forecast_timing": ["75p"],
                "region_name": ["Central America"],
                "_dataset_id": ["jrc.mars-crop-yield"],
                **meta,
            }
        ),
    )

    write_table(
        lake_root,
        "fao",
        "giews_crop_prospects",
        pa.table(
            {
                "country_code": ["BRA"],
                "country_name": ["Brazil"],
                "crop_slug": ["milho"],
                "marketing_year": ["2024"],
                "production_trend": ["stable"],
                "outlook_note": ["Favorable rains in Center-West"],
                "_dataset_id": ["fao.giews-crop-prospects"],
                **meta,
            }
        ),
    )

    write_table(
        lake_root,
        "fao",
        "amis_market_monitor",
        pa.table(
            {
                "commodity_slug": ["milho"],
                "refmonth": ["2024-01-01"],
                "indicator_slug": ["price_index"],
                "value": ["118.5"],
                "unit": ["index_2007_2008=100"],
                "_dataset_id": ["fao.amis-market-monitor"],
                **meta,
            }
        ),
    )

    write_table(
        lake_root,
        "sagis",
        "grain_supply_statistics",
        pa.table(
            {
                "commodity_slug": ["milho"],
                "marketing_year": ["2024"],
                "supply_t": ["14500000"],
                "demand_t": ["11200000"],
                "opening_stocks_t": ["4200000"],
                "closing_stocks_t": ["7500000"],
                "_dataset_id": ["sagis.grain-supply-statistics"],
                **meta,
            }
        ),
    )

    write_table(
        lake_root,
        "japan",
        "maff_ag_trade",
        pa.table(
            {
                "commodity_slug": ["arroz"],
                "refyear": ["2023"],
                "flow_code": ["export"],
                "value_jpy": ["12500000000"],
                "quantity_t": ["85000"],
                "_dataset_id": ["japan.maff-ag-trade"],
                **meta,
            }
        ),
    )

    write_table(
        lake_root,
        "fred",
        "commodity_indexes",
        pa.table(
            {
                "series_id": ["PALLFNFINDEXM"],
                "refmonth": ["2024-01-01"],
                "value": ["118.0"],
                "_dataset_id": ["fred.commodity-indexes"],
                **meta,
            }
        ),
    )

    write_table(
        lake_root,
        "nasa",
        "power_agroclimatology",
        pa.table(
            {
                "latitude": ["-15.8000"],
                "longitude": ["-47.9000"],
                "refdate": ["2024-01-01"],
                "parameter_slug": ["prectotcorr"],
                "value": ["22.4400"],
                "_dataset_id": ["nasa.power-agroclimatology"],
                **meta,
            }
        ),
    )

    write_table(
        lake_root,
        "copernicus",
        "era5_agroclimate",
        pa.table(
            {
                "latitude": ["-15.8"],
                "longitude": ["-47.9"],
                "refdate": ["2024-01-01"],
                "variable_slug": ["t2m"],
                "value": ["24.0"],
                "unit": ["degC"],
                "_dataset_id": ["copernicus.era5-agroclimate"],
                **meta,
            }
        ),
    )

    print(f"seeded international sources wave 4 silver under {lake_root / 'silver'}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
