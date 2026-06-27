#!/usr/bin/env python3
"""Seed minimal silver Delta for BR sources wave 4 CI (Phase 48)."""

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
        "mart_ibge__censo_agro_estabelecimentos",
        "mart_ibge__pnad_continua_rural",
        "mart_suframa__comercio_mercadorias_zfm",
        "mart_transportes__mtr_bit_malha_rodoviaria",
        "mart_mapa__sif_abate_estatisticas",
        "mart_ons__carga_energetica",
        "mart_inpe__deter_alertas_desmatamento",
        "mart_dnit__condicoes_conservacao_rodovias",
    ]
    for table in gold_tables:
        (lake_root / "gold" / table).mkdir(parents=True, exist_ok=True)

    source = str(lake_root / "bronze/seed.parquet")
    ingested = "2026-06-26T12:00:00Z"

    ibge_censo = pa.table(
        {
            "sidra_tabela": ["6878", "6878"],
            "codigo_uf": ["11", "12"],
            "uf": ["Rondônia", "Acre"],
            "ano": ["2017", "2017"],
            "variavel_codigo": ["183", "184"],
            "variavel": [
                "Número de estabelecimentos agropecuários",
                "Área dos estabelecimentos agropecuários",
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
            "valor": ["91438", "9219883"],
            "unidade_codigo": ["1020", "1006"],
            "unidade": ["Unidades", "Hectares"],
            "_dataset_id": [
                "ibge.censo-agro-estabelecimentos",
                "ibge.censo-agro-estabelecimentos",
            ],
            "_ingested_at": [ingested, ingested],
            "_source_file": [source, source],
        }
    )
    write_table(lake_root, "ibge", "censo_agro_estabelecimentos", ibge_censo)

    ibge_pnad = pa.table(
        {
            "sidra_tabela": ["6385", "6385"],
            "codigo_uf": ["11", "12"],
            "uf": ["Rondônia", "Acre"],
            "trimestre": ["202506", "202506"],
            "variavel_codigo": ["10605", "10606"],
            "variavel": [
                "Pessoas de 14 anos ou mais de idade ocupadas na semana de referência",
                "Coeficiente de variação - Pessoas de 14 anos ou mais de idade ocupadas na semana de referência",
            ],
            "valor": ["450", "2.1"],
            "unidade_codigo": ["1000", "2"],
            "unidade": ["Mil pessoas", "%"],
            "_dataset_id": ["ibge.pnad-continua-rural", "ibge.pnad-continua-rural"],
            "_ingested_at": [ingested, ingested],
            "_source_file": [source, source],
        }
    )
    write_table(lake_root, "ibge", "pnad_continua_rural", ibge_pnad)

    suframa = pa.table(
        {
            "column_1": ["NF", "NF"],
            "column_2": ["2024", "2024"],
            "_dataset_id": ["suframa.comercio-mercadorias-zfm", "suframa.comercio-mercadorias-zfm"],
            "_ingested_at": [ingested, ingested],
            "_source_file": [source, source],
        }
    )
    write_table(lake_root, "suframa", "comercio_mercadorias_zfm", suframa)

    transportes = pa.table(
        {
            "BR": ["010", "010"],
            "UF": ["DF", "DF"],
            "Código": ["010BDF0010", "010BDF0015"],
            "Jurisdição": ["Federal", "Federal"],
            "_dataset_id": [
                "transportes.mtr-bit-malha-rodoviaria",
                "transportes.mtr-bit-malha-rodoviaria",
            ],
            "_ingested_at": [ingested, ingested],
            "_source_file": [source, source],
        }
    )
    write_table(lake_root, "transportes", "mtr_bit_malha_rodoviaria", transportes)

    mapa_sif = pa.table(
        {
            "MES_ANO": ["01/2024", "01/2024"],
            "UF_PROCEDENCIA": ["SP", "PR"],
            "CATEGORIA": ["Bovino", "Suino"],
            "QTD_MACHO": ["1000", "500"],
            "QTD_FEMEA": ["900", "480"],
            "_dataset_id": ["mapa.sif-abate-estatisticas", "mapa.sif-abate-estatisticas"],
            "_ingested_at": [ingested, ingested],
            "_source_file": [source, source],
        }
    )
    write_table(lake_root, "mapa", "sif_abate_estatisticas", mapa_sif)

    ons_carga = pa.table(
        {
            "id_subsistema": ["SE", "S"],
            "nom_subsistema": ["Sudeste/Centro-Oeste", "Sul"],
            "din_instante": ["2024-01-01", "2024-01-01"],
            "val_cargaenergiamwmed": ["35089.38", "10472.40"],
            "_dataset_id": ["ons.carga-energetica", "ons.carga-energetica"],
            "_ingested_at": [ingested, ingested],
            "_source_file": [source, source],
        }
    )
    write_table(lake_root, "ons", "carga_energetica", ons_carga)

    inpe_deter = pa.table(
        {
            "view_date": ["2024-01-14", "2024-01-15"],
            "class_name": ["DESMATAMENTO_CR", "DEGRADACAO"],
            "uf": ["PA", "MT"],
            "municipality": ["Obidos", "Sinop"],
            "area_uc_km": ["0.5", "1.2"],
            "publish_month": ["2024-01", "2024-01"],
            "_dataset_id": [
                "inpe.deter-alertas-desmatamento",
                "inpe.deter-alertas-desmatamento",
            ],
            "_ingested_at": [ingested, ingested],
            "_source_file": [source, source],
        }
    )
    write_table(lake_root, "inpe", "deter_alertas_desmatamento", inpe_deter)

    dnit_cond = pa.table(
        {
            "id_malha": ["17475", "17475"],
            "UF": ["RR", "RR"],
            "Rodovia": ["BR-210", "BR-210"],
            "km": ["89", "90"],
            "ICM": ["51,25", "56,5"],
            "_dataset_id": [
                "dnit.condicoes-conservacao-rodovias",
                "dnit.condicoes-conservacao-rodovias",
            ],
            "_ingested_at": [ingested, ingested],
            "_source_file": [source, source],
        }
    )
    write_table(lake_root, "dnit", "condicoes_conservacao_rodovias", dnit_cond)

    print(f"seeded BR sources wave 4 silver under {lake_root / 'silver'}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
