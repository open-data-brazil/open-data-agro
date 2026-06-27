#!/usr/bin/env python3
"""Seed minimal silver Delta for BR sources wave 5 IBGE CI (Phase 53)."""

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
        "mart_ibge__ppm_efetivo_rebanhos",
        "mart_ibge__ppm_vacas_ordenhadas",
        "mart_ibge__ppm_ovinos_tosquiados",
        "mart_ibge__ppm_aquicultura",
        "mart_ibge__pam_precos_produtor",
        "mart_ibge__pam_culturas_estendidas",
        "mart_ibge__lspa_rendimento_medio",
        "mart_ibge__censo_agro_area_uso_solo",
        "mart_ibge__censo_agro_maquinario",
        "mart_ibge__pnad_rural_renda_ocupacao",
    ]
    for table in gold_tables:
        (lake_root / "gold" / table).mkdir(parents=True, exist_ok=True)

    source = str(lake_root / "bronze/seed.parquet")
    ingested = "2026-06-27T12:00:00Z"

    ppm_uf = pa.table(
        {
            "sidra_tabela": ["3939", "94"],
            "codigo_uf": ["11", "12"],
            "uf": ["Rondônia", "Acre"],
            "ano": ["2023", "2023"],
            "variavel_codigo": ["105", "107"],
            "variavel": ["Efetivo dos rebanhos", "Vacas ordenhadas"],
            "categoria_codigo": ["0", ""],
            "categoria": ["Total", ""],
            "valor": ["1200000", "354715"],
            "unidade_codigo": ["24", "24"],
            "unidade": ["Cabeças", "Cabeças"],
            "_dataset_id": ["ibge.ppm-efetivo-rebanhos", "ibge.ppm-vacas-ordenhadas"],
            "_ingested_at": [ingested, ingested],
            "_source_file": [source, source],
        }
    )
    write_table(lake_root, "ibge", "ppm_efetivo_rebanhos", ppm_uf.slice(0, 1))
    write_table(lake_root, "ibge", "ppm_vacas_ordenhadas", ppm_uf.slice(1, 2))
    write_table(lake_root, "ibge", "ppm_ovinos_tosquiados", ppm_uf.slice(1, 2))
    write_table(lake_root, "ibge", "ppm_aquicultura", ppm_uf.slice(0, 1))

    pam = pa.table(
        {
            "sidra_tabela": ["1612", "1612"],
            "codigo_ibge": ["1100015", "1100023"],
            "municipio": ["Alta Floresta D'Oeste - RO", "Ariquemes - RO"],
            "ano": ["2023", "2023"],
            "variavel_codigo": ["214", "109"],
            "variavel": ["Valor da produção", "Quantidade produzida"],
            "produto_codigo": ["2713", "2713"],
            "produto": ["Soja (em grão)", "Soja (em grão)"],
            "valor": ["50000", "12000"],
            "unidade_codigo": ["40", "1000"],
            "unidade": ["Mil Reais", "Toneladas"],
            "_dataset_id": ["ibge.pam-precos-produtor", "ibge.pam-culturas-estendidas"],
            "_ingested_at": [ingested, ingested],
            "_source_file": [source, source],
        }
    )
    write_table(lake_root, "ibge", "pam_precos_produtor", pam.slice(0, 1))
    write_table(lake_root, "ibge", "pam_culturas_estendidas", pam)

    lspa = pa.table(
        {
            "sidra_tabela": ["6588"],
            "codigo_uf": ["11"],
            "uf": ["Rondônia"],
            "mes": ["202312"],
            "variavel_codigo": ["35"],
            "variavel": ["Rendimento médio da produção"],
            "produto_codigo": ["39443"],
            "produto": ["Soja (em grão)"],
            "produto_slug": ["soja"],
            "valor": ["3200"],
            "unidade_codigo": ["1000"],
            "unidade": ["Quilogramas por Hectare"],
            "_dataset_id": ["ibge.lspa-rendimento-medio"],
            "_ingested_at": [ingested],
            "_source_file": [source],
        }
    )
    write_table(lake_root, "ibge", "lspa_rendimento_medio", lspa)

    censo = pa.table(
        {
            "sidra_tabela": ["6879", "6880"],
            "codigo_uf": ["11", "11"],
            "uf": ["Rondônia", "Rondônia"],
            "ano": ["2017", "2017"],
            "variavel_codigo": ["184", "183"],
            "variavel": [
                "Área dos estabelecimentos agropecuários",
                "Número de estabelecimentos agropecuários",
            ],
            "condicao_produtor_codigo": ["46502", "46502"],
            "condicao_produtor": ["Total", "Total"],
            "tipologia_codigo": ["46302", "46302"],
            "tipologia": ["Total", "Total"],
            "atividade_codigo": ["113601", "113601"],
            "atividade": ["Total", "Total"],
            "sexo_produtor_codigo": ["41145", "41145"],
            "sexo_produtor": ["Total", "Total"],
            "idade_produtor_codigo": ["45951", "45951"],
            "idade_produtor": ["Total", "Total"],
            "valor": ["9219883", "91438"],
            "unidade_codigo": ["1006", "1020"],
            "unidade": ["Hectares", "Unidades"],
            "_dataset_id": [
                "ibge.censo-agro-area-uso-solo",
                "ibge.censo-agro-maquinario",
            ],
            "_ingested_at": [ingested, ingested],
            "_source_file": [source, source],
        }
    )
    write_table(lake_root, "ibge", "censo_agro_area_uso_solo", censo.slice(0, 1))
    write_table(lake_root, "ibge", "censo_agro_maquinario", censo.slice(1, 2))

    pnad = pa.table(
        {
            "sidra_tabela": ["6385", "6385"],
            "codigo_uf": ["11", "12"],
            "uf": ["Rondônia", "Acre"],
            "trimestre": ["202506", "202506"],
            "variavel_codigo": ["4099", "4096"],
            "variavel": ["Rendimento médio habitual", "Ocupação"],
            "valor": ["2500", "45.2"],
            "unidade_codigo": ["101", "102"],
            "unidade": ["Reais", "Percentual"],
            "_dataset_id": [
                "ibge.pnad-rural-renda-ocupacao",
                "ibge.pnad-rural-renda-ocupacao",
            ],
            "_ingested_at": [ingested, ingested],
            "_source_file": [source, source],
        }
    )
    write_table(lake_root, "ibge", "pnad_rural_renda_ocupacao", pnad)

    print(f"seeded wave 5 IBGE silver under {lake_root}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
