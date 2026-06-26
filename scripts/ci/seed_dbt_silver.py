#!/usr/bin/env python3
"""Seed minimal silver Delta tables for dbt CI (local-first, no R2)."""

from __future__ import annotations

import os
import sys
from pathlib import Path

import pyarrow as pa
from deltalake import write_deltalake


def write_table(root: Path, table: str, data: pa.Table) -> None:
    path = root / "silver" / "conab" / table
    path.parent.mkdir(parents=True, exist_ok=True)
    write_deltalake(str(path), data, mode="overwrite")


def main() -> int:
    lake_root = Path(os.environ.get("LAKE_LOCAL_ROOT", "/tmp/open-data-agro-lake"))
    lake_root.mkdir(parents=True, exist_ok=True)
    (lake_root / "gold").mkdir(parents=True, exist_ok=True)
    (lake_root / "gold" / "mart_conab__estimativa_graos").mkdir(parents=True, exist_ok=True)
    (lake_root / "gold" / "mart_conab__serie_historica_graos").mkdir(parents=True, exist_ok=True)
    (lake_root / "gold" / "mart_conab__estimativa_cana").mkdir(parents=True, exist_ok=True)
    (lake_root / "gold" / "mart_conab__serie_historica_cana").mkdir(parents=True, exist_ok=True)
    (lake_root / "gold" / "mart_conab__estimativa_cafe").mkdir(parents=True, exist_ok=True)
    (lake_root / "gold" / "mart_conab__serie_historica_cafe").mkdir(parents=True, exist_ok=True)
    (lake_root / "gold" / "mart_conab__custo_producao").mkdir(parents=True, exist_ok=True)

    meta = {
        "_dataset_id": ["conab.estimativa-graos"],
        "_ingested_at": ["2026-06-25T12:00:00Z"],
        "_source_file": [str(lake_root / "bronze/seed.parquet")],
    }

    estimativa = pa.table(
        {
            "ano_agricola": ["2025/26", "2025/26"],
            "safra": ["UNICA", "UNICA"],
            "uf": ["PR", "MT"],
            "produto": ["SOJA", "MILHO"],
            "id_produto": ["1", "2"],
            "id_levantamento": ["012", "012"],
            "dsc_levantamento": ["12 LEV", "12 LEV"],
            "area_plantada_mil_ha": ["100", "120"],
            "producao_mil_t": ["100", "120"],
            "produtividade_mil_ha_mil_t": ["1.0", "1.0"],
            "_dataset_id": ["conab.estimativa-graos", "conab.estimativa-graos"],
            "_ingested_at": ["2026-06-25T12:00:00Z", "2026-06-25T12:00:00Z"],
            "_source_file": [meta["_source_file"][0], meta["_source_file"][0]],
        }
    )
    write_table(lake_root, "estimativa_graos", estimativa)

    serie = pa.table(
        {
            "ano_agricola": ["2020/21", "2021/22"],
            "dsc_safra_previsao": ["UNICA", "UNICA"],
            "uf": ["PR", "RS"],
            "produto": ["SOJA", "SOJA"],
            "id_produto": ["1", "1"],
            "area_plantada_mil_ha": ["50", "55"],
            "producao_mil_t": ["50", "55"],
            "produtividade_mil_ha_mil_t": ["1.0", "1.0"],
            "_dataset_id": ["conab.serie-historica-graos", "conab.serie-historica-graos"],
            "_ingested_at": ["2026-06-25T12:00:00Z", "2026-06-25T12:00:00Z"],
            "_source_file": [meta["_source_file"][0], meta["_source_file"][0]],
        }
    )
    write_table(lake_root, "serie_historica_graos", serie)

    estimativa_cana = pa.table(
        {
            "ano_agricola": ["2025/26", "2025/26"],
            "dsc_safra_previsao": ["UNICA", "UNICA"],
            "uf": ["SP", "MT"],
            "produto": ["CANA DE ACUCAR", "CANA DE ACUCAR"],
            "id_produto": ["4238", "4238"],
            "dsc_levantamento": ["1º LEV", "2º LEV"],
            "id_levantamento": ["1", "2"],
            "area_plantada_mil_ha": ["100", "120"],
            "producao_mil_t": ["1000", "1200"],
            "producao_acucar_mil_t": ["50", "60"],
            "producao_etanol_anidro_mil_l": ["100", "110"],
            "producao_etanol_hidratado_mil_l": ["200", "220"],
            "producao_etanol_total_mil_l": ["300", "330"],
            "produtcao_atr_kg_t": ["140", "145"],
            "_dataset_id": ["conab.estimativa-cana", "conab.estimativa-cana"],
            "_ingested_at": ["2026-06-25T12:00:00Z", "2026-06-25T12:00:00Z"],
            "_source_file": [meta["_source_file"][0], meta["_source_file"][0]],
        }
    )
    write_table(lake_root, "estimativa_cana", estimativa_cana)

    serie_cana = pa.table(
        {
            "ano_agricola": ["2020/21", "2021/22"],
            "dsc_safra_previsao": ["UNICA", "UNICA"],
            "uf": ["SP", "MT"],
            "produto": ["CANA DE ACUCAR", "CANA DE ACUCAR"],
            "id_produto": ["4238", "4238"],
            "area_plantada_mil_ha": ["50", "55"],
            "producao_mil_t": ["500", "550"],
            "dsc_situacao_levantamento": ["FINAL", "FINAL"],
            "producao_acucar_mil_t": ["25", "28"],
            "producao_etanol_anidro_mil_l": ["80", "85"],
            "producao_etanol_hidratado_mil_l": ["90", "95"],
            "producao_etanol_total_mil_l": ["170", "180"],
            "produtcao_atr_kg_t": ["140", "142"],
            "_dataset_id": ["conab.serie-historica-cana", "conab.serie-historica-cana"],
            "_ingested_at": ["2026-06-25T12:00:00Z", "2026-06-25T12:00:00Z"],
            "_source_file": [meta["_source_file"][0], meta["_source_file"][0]],
        }
    )
    write_table(lake_root, "serie_historica_cana", serie_cana)

    estimativa_cafe = pa.table(
        {
            "ano_agricola": ["2025/26", "2025/26"],
            "safra": ["UNICA", "UNICA"],
            "uf": ["MG", "SP"],
            "produto": ["CAFE", "CAFE"],
            "id_produto": ["7498", "7498"],
            "id_levantamento": ["001", "002"],
            "dsc_levantamento": ["1º LEV", "2º LEV"],
            "area_plantada_mil_ha": ["100", "120"],
            "producao_mil_t": ["50", "60"],
            "produtividade_mil_ha_mil_t": ["0.5", "0.5"],
            "_dataset_id": ["conab.estimativa-cafe", "conab.estimativa-cafe"],
            "_ingested_at": ["2026-06-25T12:00:00Z", "2026-06-25T12:00:00Z"],
            "_source_file": [meta["_source_file"][0], meta["_source_file"][0]],
        }
    )
    write_table(lake_root, "estimativa_cafe", estimativa_cafe)

    serie_cafe = pa.table(
        {
            "ano_agricola": ["2020/21", "2021/22"],
            "dsc_safra_previsao": ["UNICA", "UNICA"],
            "uf": ["MG", "SP"],
            "produto": ["CAFE", "CAFE"],
            "id_produto": ["7498", "7498"],
            "area_plantada_mil_ha": ["90", "95"],
            "producao_mil_t": ["45", "48"],
            "produtividade_mil_ha_mil_t": ["0.5", "0.5"],
            "_dataset_id": ["conab.serie-historica-cafe", "conab.serie-historica-cafe"],
            "_ingested_at": ["2026-06-25T12:00:00Z", "2026-06-25T12:00:00Z"],
            "_source_file": [meta["_source_file"][0], meta["_source_file"][0]],
        }
    )
    write_table(lake_root, "serie_historica_cafe", serie_cafe)

    custo = pa.table(
        {
            "empreendimento": ["AGRICULTURA EMPRESARIAL", "AGRICULTURA FAMILIAR"],
            "ano": ["2024", "2024"],
            "mes": ["1", "3"],
            "ano_mes": ["202401", "202403"],
            "produto": ["CAFE", "SOJA"],
            "id_produto": ["7498", "1"],
            "safra": ["TODAS", "TODAS"],
            "uf": ["MG", "MT"],
            "municipio": ["VARGINHA-MG", "SORRISO-MT"],
            "cod_ibge": ["3170701", "5107925"],
            "unidade_comercializacao": ["60 kg", "60 kg"],
            "vlr_custo_variavel_ha": ["10000", "8000"],
            "vlr_custo_variavel_unidade": ["300", "250"],
            "vlr_custo_fixo_ha": ["1000", "900"],
            "vlr_custo_fixo_unidade": ["30", "28"],
            "vlr_renda_fator_ha": ["500", "400"],
            "vlr_renda_fator_unidade": ["15", "12"],
            "_dataset_id": ["conab.custo-producao", "conab.custo-producao"],
            "_ingested_at": ["2026-06-25T12:00:00Z", "2026-06-25T12:00:00Z"],
            "_source_file": [meta["_source_file"][0], meta["_source_file"][0]],
        }
    )
    write_table(lake_root, "custo_producao", custo)

    print(f"seeded silver tables under {lake_root / 'silver' / 'conab'}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
