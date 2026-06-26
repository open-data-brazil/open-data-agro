#!/usr/bin/env python3
"""Seed minimal silver Delta for BR sources wave 3 CI (Phase 44)."""

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
        "mart_dnit__snv_rodovias_federais",
        "mart_ipea__series_macro_regionais",
        "mart_ibge__pevs_producao_vegetal",
        "mart_ibge__ppm_producao_municipal",
        "mart_aneel__tarifas_energia",
        "mart_bndes__financiamento_agro",
        "mart_inmet__sequia_monitor",
    ]
    for table in gold_tables:
        (lake_root / "gold" / table).mkdir(parents=True, exist_ok=True)

    source = str(lake_root / "bronze/seed.parquet")
    ingested = "2026-06-26T12:00:00Z"

    dnit_snv = pa.table(
        {
            "BR": ["010", "010"],
            "UF": ["DF", "DF"],
            "Tipo de trecho": ["Eixo Principal", "Eixo Principal"],
            "Código": ["010BDF0010", "010BDF0015"],
            "Local de Início": ["ENTR BR-020", "ENTR DF-440"],
            "Local de Fim": ["ENTR DF-440", "ACESSO I SOBRADINHO"],
            "km inicial": ["0", "2,4"],
            "km final": ["2,4", "6"],
            "Extensão": ["2,4", "3,6"],
            "Superfície Federal": ["DUP", "DUP"],
            "Obras": ["", ""],
            "Federal Coincidente": ["010BDF0010", "010BDF0015"],
            "Administração": ["Convênio de Administração", "Convênio de Administração"],
            "Ato legal": ["", ""],
            "Estadual Coincidente": ["", ""],
            "Superfície Est. Coincidente": ["", ""],
            "Jurisdição": ["Federal", "Federal"],
            "Superfície": ["PAV", "PAV"],
            "Unidade Local": ["Brasília", "Brasília"],
            "_dataset_id": ["dnit.snv-rodovias-federais", "dnit.snv-rodovias-federais"],
            "_ingested_at": [ingested, ingested],
            "_source_file": [source, source],
        }
    )
    write_table(lake_root, "dnit", "snv_rodovias_federais", dnit_snv)

    ipea_series = pa.table(
        {
            "series_code": ["ADH_P_AGRO_RUR", "ADH_P_COM_RUR"],
            "refdate": ["2010-01-01", "2010-01-01"],
            "value": ["65.18", "4.5"],
            "region_level": ["Brasil", "Brasil"],
            "territory_code": ["0", "0"],
            "_dataset_id": ["ipea.series-macro-regionais", "ipea.series-macro-regionais"],
            "_ingested_at": [ingested, ingested],
            "_source_file": [source, source],
        }
    )
    write_table(lake_root, "ipea", "series_macro_regionais", ipea_series)

    ibge_pevs = pa.table(
        {
            "sidra_tabela": ["289", "289"],
            "codigo_uf": ["11", "11"],
            "uf": ["Rondônia", "Rondônia"],
            "ano": ["2023", "2023"],
            "variavel_codigo": ["144", "145"],
            "variavel": [
                "Quantidade produzida na extração vegetal",
                "Valor da produção na extração vegetal",
            ],
            "produto_codigo": ["0", "0"],
            "produto": ["Total", "Total"],
            "valor": ["..", "143835"],
            "unidade_codigo": ["", "40"],
            "unidade": ["", "Mil Reais"],
            "_dataset_id": ["ibge.pevs-producao-vegetal", "ibge.pevs-producao-vegetal"],
            "_ingested_at": [ingested, ingested],
            "_source_file": [source, source],
        }
    )
    write_table(lake_root, "ibge", "pevs_producao_vegetal", ibge_pevs)

    ibge_ppm = pa.table(
        {
            "sidra_tabela": ["74", "74"],
            "codigo_ibge": ["1100015", "1100015"],
            "municipio": ["Alta Floresta D'Oeste - RO", "Alta Floresta D'Oeste - RO"],
            "ano": ["2023", "2023"],
            "variavel_codigo": ["106", "215"],
            "variavel": ["Produção de origem animal", "Valor da produção"],
            "produto_codigo": ["0", "0"],
            "produto": ["Total", "Total"],
            "valor": ["..", "11307"],
            "unidade_codigo": ["", "40"],
            "unidade": ["", "Mil Reais"],
            "_dataset_id": ["ibge.ppm-producao-municipal", "ibge.ppm-producao-municipal"],
            "_ingested_at": [ingested, ingested],
            "_source_file": [source, source],
        }
    )
    write_table(lake_root, "ibge", "ppm_producao_municipal", ibge_ppm)

    aneel_tarifas = pa.table(
        {
            "DatGeracaoConjuntoDados": ["2026-06-01", "2026-06-01"],
            "DatCompetencia": ["2015-01-01", "2015-02-01"],
            "NomBandeiraAcionada": ["Vermelha P1", "Vermelha P1"],
            "VlrAdicionalBandeira": ["30,00", "30,00"],
            "_dataset_id": ["aneel.tarifas-energia", "aneel.tarifas-energia"],
            "_ingested_at": [ingested, ingested],
            "_source_file": [source, source],
        }
    )
    write_table(lake_root, "aneel", "tarifas_energia", aneel_tarifas)

    bndes_fin = pa.table(
        {
            "ano": ["1995", "1995"],
            "mes": ["1", "2"],
            "agropecuaria": ["49,4775568700004", "35,8757577499998"],
            "_dataset_id": ["bndes.financiamento-agro", "bndes.financiamento-agro"],
            "_ingested_at": [ingested, ingested],
            "_source_file": [source, source],
        }
    )
    write_table(lake_root, "bndes", "financiamento_agro", bndes_fin)

    inmet_sequia = pa.table(
        {
            "mapa_id": ["304", "304"],
            "ano": ["2026", "2026"],
            "mes": ["5", "5"],
            "categoria_seca": ["S4", "S3"],
            "area_km2": ["0", "0"],
            "area_id": ["5", "5"],
            "tipo_area": ["8", "8"],
            "_dataset_id": ["inmet.sequia-monitor", "inmet.sequia-monitor"],
            "_ingested_at": [ingested, ingested],
            "_source_file": [source, source],
        }
    )
    write_table(lake_root, "inmet", "sequia_monitor", inmet_sequia)

    print(f"seeded BR sources wave 3 silver under {lake_root / 'silver'}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
