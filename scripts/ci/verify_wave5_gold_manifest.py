#!/usr/bin/env python3
"""Verify wave 5 gold marts exist and map to expected analytics table names."""

from __future__ import annotations

import argparse
import os
import sys
from pathlib import Path

# Phases 52–56 — 29 implemented dataset_ids (32 approved minus 3 Phase 51 rejections).
WAVE5_MARTS: list[tuple[str, str]] = [
    ("mart_mapa__sipeagro_estabelecimentos", "mapa_sipeagro_estabelecimentos"),
    ("mart_mapa__sipeagro_produtos", "mapa_sipeagro_produtos"),
    ("mart_mapa__sigef_producao_sementes", "mapa_sigef_producao_sementes"),
    ("mart_mapa__sigef_areas", "mapa_sigef_areas"),
    ("mart_mapa__sisser_seguro_rural", "mapa_sisser_seguro_rural"),
    ("mart_ibge__ppm_efetivo_rebanhos", "ibge_ppm_efetivo_rebanhos"),
    ("mart_ibge__ppm_vacas_ordenhadas", "ibge_ppm_vacas_ordenhadas"),
    ("mart_ibge__ppm_ovinos_tosquiados", "ibge_ppm_ovinos_tosquiados"),
    ("mart_ibge__ppm_aquicultura", "ibge_ppm_aquicultura"),
    ("mart_ibge__pam_precos_produtor", "ibge_pam_precos_produtor"),
    ("mart_ibge__pam_culturas_estendidas", "ibge_pam_culturas_estendidas"),
    ("mart_ibge__lspa_rendimento_medio", "ibge_lspa_rendimento_medio"),
    ("mart_ibge__censo_agro_area_uso_solo", "ibge_censo_agro_area_uso_solo"),
    ("mart_ibge__censo_agro_maquinario", "ibge_censo_agro_maquinario"),
    ("mart_ibge__pnad_rural_renda_ocupacao", "ibge_pnad_rural_renda_ocupacao"),
    ("mart_ibama__sisfogo_incendios", "ibama_sisfogo_incendios"),
    ("mart_ibama__licencas_ambientais", "ibama_licencas_ambientais"),
    ("mart_ibama__autos_infracao", "ibama_autos_infracao"),
    ("mart_ana__pluviometria_redes", "ana_pluviometria_redes"),
    ("mart_embrapa__agroapi_agrofit", "embrapa_agroapi_agrofit"),
    ("mart_transportes__mtr_bit_malha_shapefile", "transportes_mtr_bit_malha_shapefile"),
    ("mart_bcb__cim_agro_credito_rural", "bcb_cim_agro_credito_rural"),
    ("mart_bndes__desembolsos_linhas_agro", "bndes_desembolsos_linhas_agro"),
    ("mart_anp__etanol_precos", "anp_etanol_precos"),
    ("mart_abiove__balanco_complexo_soja", "abiove_balanco_complexo_soja"),
    ("mart_abiove__exportacoes_complexo_soja", "abiove_exportacoes_complexo_soja"),
    ("mart_abiove__capacidade_instalada_esmagamento", "abiove_capacidade_instalada_esmagamento"),
    ("mart_b3__futuro_cafe", "b3_futuro_cafe"),
    ("mart_b3__futuro_acucar", "b3_futuro_acucar"),
]


def mart_table_name(dir_name: str) -> str:
    if not dir_name.startswith("mart_"):
        raise ValueError(dir_name)
    return dir_name.removeprefix("mart_").replace("__", "_")


def main() -> int:
    parser = argparse.ArgumentParser(description="Verify wave 5 gold mart manifest")
    parser.add_argument("--lake-root", default=os.environ.get("LAKE_LOCAL_ROOT", "./lake"))
    args = parser.parse_args()

    lake_root = Path(args.lake_root).resolve()
    gold = lake_root / "gold"
    if not gold.is_dir():
        print(f"missing gold dir: {gold}", file=sys.stderr)
        return 2

    errors: list[str] = []
    for mart_dir, want_table in WAVE5_MARTS:
        got_table = mart_table_name(mart_dir)
        if got_table != want_table:
            errors.append(f"{mart_dir}: table name {got_table!r} != {want_table!r}")
        parquet = gold / mart_dir / "mart.parquet"
        if not parquet.is_file():
            errors.append(f"missing {parquet}")

    if errors:
        print("wave 5 gold manifest failures:", file=sys.stderr)
        for err in errors:
            print(f"  - {err}", file=sys.stderr)
        return 1

    print(f"wave 5 gold manifest ok ({len(WAVE5_MARTS)} marts under {gold})")
    return 0


if __name__ == "__main__":
    sys.exit(main())
