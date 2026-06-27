#!/usr/bin/env python3
"""Seed minimal silver Delta for BR sources wave 5 MAPA CI (Phase 52)."""

from __future__ import annotations

import os
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
        "mart_mapa__sipeagro_estabelecimentos",
        "mart_mapa__sipeagro_produtos",
        "mart_mapa__sigef_producao_sementes",
        "mart_mapa__sigef_areas",
        "mart_mapa__sisser_seguro_rural",
    ]
    for table in gold_tables:
        (lake_root / "gold" / table).mkdir(parents=True, exist_ok=True)

    source = str(lake_root / "bronze/seed.parquet")
    ingested = "2026-06-27T12:00:00Z"

    sipeagro = pa.table(
        {
            "linha_produto": ["Fertilizantes", "Qualidade Vegetal"],
            "uf": ["SC", "SP"],
            "municipio": ["Gaspar", "Campinas"],
            "numero_registro_estabelecimento": ["SC0001201", "SP0001001"],
            "status_registro": ["Ativo", "Ativo"],
            "cpf_cnpj": ["**.***.101/***-**", "**.***.200/***-**"],
            "razao_social": ["BUNGE ALIMENTOS S/A", "EXAMPLE AGRO LTDA"],
            "nome_fantasia": ["BUNGE ALIMENTOS S/A", "EXAMPLE AGRO"],
            "area_atuacao": ["FERTILIZANTE", "QUALIDADE VEGETAL"],
            "atividade": ["IMPORTADOR", "PRODUTOR"],
            "classificacao": ["PRODUTO IMPORTADO", "PRODUTO NACIONAL"],
            "caracteristica_adicional": ["", ""],
            "especie": ["", ""],
            "_dataset_id": [
                "mapa.sipeagro-estabelecimentos",
                "mapa.sipeagro-produtos",
            ],
            "_ingested_at": [ingested, ingested],
            "_source_file": [source, source],
        }
    )
    write_table(lake_root, "mapa", "sipeagro_estabelecimentos", sipeagro.slice(0, 1))
    write_table(lake_root, "mapa", "sipeagro_produtos", sipeagro.slice(1, 2))

    sigef_prod = pa.table(
        {
            "Safra": ["2024/2025", "2024/2025"],
            "Especie": ["SOJA", "MILHO"],
            "Categoria": ["S1", "S1"],
            "Cultivar": ["BRS 284", "AG 8088"],
            "Municipio": ["Sorriso", "Lucas do Rio Verde"],
            "UF": ["MT", "MT"],
            "Status": ["Colhido", "Colhido"],
            "Data do Plantio": ["01/10/2024", "15/09/2024"],
            "Data de Colheita": ["20/02/2025", "10/02/2025"],
            "Area": ["1200", "800"],
            "Producao bruta": ["3600", "6400"],
            "Producao estimada": ["3500", "6300"],
            "_dataset_id": ["mapa.sigef-producao-sementes", "mapa.sigef-producao-sementes"],
            "_ingested_at": [ingested, ingested],
            "_source_file": [source, source],
        }
    )
    write_table(lake_root, "mapa", "sigef_producao_sementes", sigef_prod)

    sigef_area = pa.table(
        {
            "TIPOPERIODO": ["SAFRINHA", "SAFRA"],
            "PERIODO": ["2024/2025", "2024/2025"],
            "AREATOTAL": ["500", "1200"],
            "MUNICIPIO": ["Cascavel", "Londrina"],
            "UF": ["PR", "PR"],
            "ESPECIE": ["SOJA", "MILHO"],
            "CULTIVAR": ["BRS 284", "AG 8088"],
            "AREAPLANTADA": ["450", "1100"],
            "AREAESTIMADA": ["460", "1150"],
            "QUANTRESERVADA": ["1000", "2000"],
            "DATAPLANTIO": ["01/02/2025", "01/10/2024"],
            "_dataset_id": ["mapa.sigef-areas", "mapa.sigef-areas"],
            "_ingested_at": [ingested, ingested],
            "_source_file": [source, source],
        }
    )
    write_table(lake_root, "mapa", "sigef_areas", sigef_area)

    sisser = pa.table(
        {
            "periodo_arquivo": ["PSR - 2025", "PSR - 2025"],
            "NM_RAZAO_SOCIAL": ["SEGURADORA EXEMPLO SA", "SEGURADORA EXEMPLO SA"],
            "CD_PROCESSO_SUSEP": ["12345", "12345"],
            "NR_PROPOSTA": ["1001", "1002"],
            "SG_UF_PROPRIEDADE": ["MT", "RS"],
            "_dataset_id": ["mapa.sisser-seguro-rural", "mapa.sisser-seguro-rural"],
            "_ingested_at": [ingested, ingested],
            "_source_file": [source, source],
        }
    )
    write_table(lake_root, "mapa", "sisser_seguro_rural", sisser)

    print(f"seeded wave 5 MAPA silver under {lake_root}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
