#!/usr/bin/env python3
"""Verify wave 3 gold marts exist and map to expected analytics table names."""

from __future__ import annotations

import argparse
import os
import sys
from pathlib import Path

WAVE3_MARTS: list[tuple[str, str]] = [
    ("mart_dnit__snv_rodovias_federais", "dnit_snv_rodovias_federais"),
    ("mart_ipea__series_macro_regionais", "ipea_series_macro_regionais"),
    ("mart_ibge__pevs_producao_vegetal", "ibge_pevs_producao_vegetal"),
    ("mart_ibge__ppm_producao_municipal", "ibge_ppm_producao_municipal"),
    ("mart_aneel__tarifas_energia", "aneel_tarifas_energia"),
    ("mart_bndes__financiamento_agro", "bndes_financiamento_agro"),
    ("mart_inmet__sequia_monitor", "inmet_sequia_monitor"),
    ("mart_oecd__ag_outlook", "oecd_ag_outlook"),
    ("mart_fao__food_price_index", "fao_food_price_index"),
    ("mart_argentina__magyp_producion_granos", "argentina_magyp_producion_granos"),
]


def mart_table_name(dir_name: str) -> str:
    if not dir_name.startswith("mart_"):
        raise ValueError(dir_name)
    return dir_name.removeprefix("mart_").replace("__", "_")


def main() -> int:
    parser = argparse.ArgumentParser(description="Verify wave 3 gold mart manifest")
    parser.add_argument("--lake-root", default=os.environ.get("LAKE_LOCAL_ROOT", "./lake"))
    args = parser.parse_args()

    lake_root = Path(args.lake_root).resolve()
    gold = lake_root / "gold"
    if not gold.is_dir():
        print(f"missing gold dir: {gold}", file=sys.stderr)
        return 2

    errors: list[str] = []
    for mart_dir, want_table in WAVE3_MARTS:
        got_table = mart_table_name(mart_dir)
        if got_table != want_table:
            errors.append(f"{mart_dir}: table name {got_table!r} != {want_table!r}")
        parquet = gold / mart_dir / "mart.parquet"
        if not parquet.is_file():
            errors.append(f"missing {parquet}")

    if errors:
        print("wave 3 gold manifest failures:", file=sys.stderr)
        for err in errors:
            print(f"  - {err}", file=sys.stderr)
        return 1

    print(f"wave 3 gold manifest ok ({len(WAVE3_MARTS)} marts under {gold})")
    return 0


if __name__ == "__main__":
    sys.exit(main())
