#!/usr/bin/env python3
"""Verify wave 4 gold marts exist and map to expected analytics table names."""

from __future__ import annotations

import argparse
import os
import sys
from pathlib import Path

WAVE4_MARTS: list[tuple[str, str]] = [
    ("mart_ibge__censo_agro_estabelecimentos", "ibge_censo_agro_estabelecimentos"),
    ("mart_ibge__pnad_continua_rural", "ibge_pnad_continua_rural"),
    ("mart_suframa__comercio_mercadorias_zfm", "suframa_comercio_mercadorias_zfm"),
    ("mart_transportes__mtr_bit_malha_rodoviaria", "transportes_mtr_bit_malha_rodoviaria"),
    ("mart_mapa__sif_abate_estatisticas", "mapa_sif_abate_estatisticas"),
    ("mart_ons__carga_energetica", "ons_carga_energetica"),
    ("mart_inpe__deter_alertas_desmatamento", "inpe_deter_alertas_desmatamento"),
    ("mart_dnit__condicoes_conservacao_rodovias", "dnit_condicoes_conservacao_rodovias"),
    ("mart_cftc__cot_agricultural_futures", "cftc_cot_agricultural_futures"),
    ("mart_jrc__mars_crop_yield", "jrc_mars_crop_yield"),
    ("mart_wto__its_trade_statistics", "wto_its_trade_statistics"),
    ("mart_fao__giews_crop_prospects", "fao_giews_crop_prospects"),
    ("mart_fao__amis_market_monitor", "fao_amis_market_monitor"),
    ("mart_sagis__grain_supply_statistics", "sagis_grain_supply_statistics"),
    ("mart_japan__maff_ag_trade", "japan_maff_ag_trade"),
    ("mart_mexico__siap_produccion_agricola", "mexico_siap_produccion_agricola"),
    ("mart_fred__commodity_indexes", "fred_commodity_indexes"),
    ("mart_nasa__power_agroclimatology", "nasa_power_agroclimatology"),
    ("mart_copernicus__era5_agroclimate", "copernicus_era5_agroclimate"),
    ("mart_noaa__gpcc_precipitation", "noaa_gpcc_precipitation"),
]


def mart_table_name(dir_name: str) -> str:
    if not dir_name.startswith("mart_"):
        raise ValueError(dir_name)
    return dir_name.removeprefix("mart_").replace("__", "_")


def main() -> int:
    parser = argparse.ArgumentParser(description="Verify wave 4 gold mart manifest")
    parser.add_argument("--lake-root", default=os.environ.get("LAKE_LOCAL_ROOT", "./lake"))
    args = parser.parse_args()

    lake_root = Path(args.lake_root).resolve()
    gold = lake_root / "gold"
    if not gold.is_dir():
        print(f"missing gold dir: {gold}", file=sys.stderr)
        return 2

    errors: list[str] = []
    for mart_dir, want_table in WAVE4_MARTS:
        got_table = mart_table_name(mart_dir)
        if got_table != want_table:
            errors.append(f"{mart_dir}: table name {got_table!r} != {want_table!r}")
        parquet = gold / mart_dir / "mart.parquet"
        if not parquet.is_file():
            errors.append(f"missing {parquet}")

    if errors:
        print("wave 4 gold manifest failures:", file=sys.stderr)
        for err in errors:
            print(f"  - {err}", file=sys.stderr)
        return 1

    print(f"wave 4 gold manifest ok ({len(WAVE4_MARTS)} marts under {gold})")
    return 0


if __name__ == "__main__":
    sys.exit(main())
