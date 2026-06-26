#!/usr/bin/env python3
"""Seed minimal silver Delta for CONAB Mercado CI (local-first, no R2)."""

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
    (lake_root / "gold" / "mart_conab__oferta_demanda").mkdir(parents=True, exist_ok=True)
    (lake_root / "gold" / "mart_conab__precos_minimos").mkdir(parents=True, exist_ok=True)
    (lake_root / "gold" / "mart_conab__precos_semanal_uf").mkdir(parents=True, exist_ok=True)
    (lake_root / "gold" / "mart_conab__precos_semanal_municipio").mkdir(parents=True, exist_ok=True)
    (lake_root / "gold" / "mart_conab__precos_mensal_uf").mkdir(parents=True, exist_ok=True)
    (lake_root / "gold" / "mart_conab__precos_mensal_municipio").mkdir(parents=True, exist_ok=True)

    source = str(lake_root / "bronze/seed.parquet")
    ingested = "2026-06-25T12:00:00Z"
    oferta = pa.table(
        {
            "produto": ["SOJA", "MILHO"],
            "id_produto": ["1", "2"],
            "dsc_safra": ["2024/25", "2024/25"],
            "estoque_inicial_1000t": ["100", "80"],
            "producao_1000t": ["150", "120"],
            "importacao_1000t": ["0", "0"],
            "consumo_1000t": ["130", "100"],
            "exportacao_1000t": ["20", "10"],
            "estoque_final_1000t": ["100", "90"],
            "_dataset_id": ["conab.oferta-demanda", "conab.oferta-demanda"],
            "_ingested_at": [ingested, ingested],
            "_source_file": [source, source],
        }
    )
    precos = pa.table(
        {
            "produto": ["SOJA", "SOJA"],
            "classificao_produto": ["EM GRAOS", "EM GRAOS"],
            "id_produto": ["4744", "4744"],
            "uf": ["MT", "MT"],
            "regiao": ["CENTRO-OESTE", "CENTRO-OESTE"],
            "ano": ["2025", "2025"],
            "mes": ["6", "6"],
            "data_inicial_final_semana": [
                "02-06-2025 - 06-06-2025",
                "09-06-2025 - 13-06-2025",
            ],
            "semana": ["1", "2"],
            "dsc_nivel_comercializacao": [
                "PRECO RECEBIDO P/ PR",
                "PRECO RECEBIDO P/ PR",
            ],
            "valor_produto_kg": ["1,84", "1,84"],
            "_dataset_id": [
                "conab.precos-agropecuarios-semanal-uf",
                "conab.precos-agropecuarios-semanal-uf",
            ],
            "_ingested_at": [ingested, ingested],
            "_source_file": [source, source],
        }
    )
    write_table(lake_root, "oferta_demanda", oferta)

    precos_minimos = pa.table(
        {
            "descricao_produto_preco_minimo": ["SOJA", "SOJA"],
            "id_produto": ["4744", "4744"],
            "uf": ["MT", "MT"],
            "regionalizacao": ["MT", "MT"],
            "ano_inicio_vigencia": ["2024", "2025"],
            "mes_incio_vigencia": ["01", "01"],
            "ano_termino_vigencia": ["2024", "2025"],
            "mes_termino_vigencia": ["12", "12"],
            "preco": ["45.24", "46.10"],
            "dsc_unidade_comercializacao": ["60 kg", "60 kg"],
            "nome_normativo": ["PORTARIA N 190", "PORTARIA N 190"],
            "url": ["NI", "NI"],
            "_dataset_id": ["conab.precos-minimos", "conab.precos-minimos"],
            "_ingested_at": [ingested, ingested],
            "_source_file": [source, source],
        }
    )
    write_table(lake_root, "precos_minimos", precos_minimos)
    write_table(lake_root, "precos_agropecuarios_semanal_uf", precos)

    precos_municipio = pa.table(
        {
            "produto": ["SOJA", "SOJA"],
            "classificao_produto": ["EM GRAOS", "EM GRAOS"],
            "id_produto": ["4744", "4744"],
            "nom_municipio": ["SORRISO-MT", "SORRISO-MT"],
            "cod_ibge": ["5107925", "5107925"],
            "uf": ["MT", "MT"],
            "regiao": ["CENTRO-OESTE", "CENTRO-OESTE"],
            "ano": ["2025", "2025"],
            "mes": ["6", "6"],
            "data_inicial_final_semana": [
                "02-06-2025 - 06-06-2025",
                "09-06-2025 - 13-06-2025",
            ],
            "semana": ["1", "2"],
            "dsc_nivel_comercializacao": [
                "PRECO RECEBIDO P/ PR",
                "PRECO RECEBIDO P/ PR",
            ],
            "valor_produto_kg": ["1,84", "1,84"],
            "_dataset_id": [
                "conab.precos-agropecuarios-semanal-municipio",
                "conab.precos-agropecuarios-semanal-municipio",
            ],
            "_ingested_at": [ingested, ingested],
            "_source_file": [source, source],
        }
    )
    write_table(lake_root, "precos_agropecuarios_semanal_municipio", precos_municipio)

    precos_mensal_uf = pa.table(
        {
            "produto": ["SOJA", "SOJA"],
            "classificao_produto": ["EM GRAOS", "EM GRAOS"],
            "id_produto": ["4744", "4744"],
            "uf": ["MT", "MT"],
            "regiao": ["CENTRO-OESTE", "CENTRO-OESTE"],
            "ano": ["2025", "2025"],
            "mes": ["6", "7"],
            "dsc_nivel_comercializacao": [
                "PRECO RECEBIDO P/ PRODUTOR",
                "PRECO RECEBIDO P/ PRODUTOR",
            ],
            "valor_produto_kg": ["1,84", "1,9"],
            "_dataset_id": [
                "conab.precos-agropecuarios-mensal-uf",
                "conab.precos-agropecuarios-mensal-uf",
            ],
            "_ingested_at": [ingested, ingested],
            "_source_file": [source, source],
        }
    )
    write_table(lake_root, "precos_agropecuarios_mensal_uf", precos_mensal_uf)

    precos_mensal_municipio = pa.table(
        {
            "produto": ["SOJA", "SOJA"],
            "classificao_produto": ["EM GRAOS", "EM GRAOS"],
            "id_produto": ["4744", "4744"],
            "nom_municipio": ["SORRISO-MT", "SORRISO-MT"],
            "cod_ibge": ["5107925", "5107925"],
            "uf": ["MT", "MT"],
            "regiao": ["CENTRO-OESTE", "CENTRO-OESTE"],
            "ano": ["2025", "2025"],
            "mes": ["6", "7"],
            "dsc_nivel_comercializacao": [
                "PRECO RECEBIDO P/ PRODUTOR",
                "PRECO RECEBIDO P/ PRODUTOR",
            ],
            "valor_produto_kg": ["1,83", "1,88"],
            "_dataset_id": [
                "conab.precos-agropecuarios-mensal-municipio",
                "conab.precos-agropecuarios-mensal-municipio",
            ],
            "_ingested_at": [ingested, ingested],
            "_source_file": [source, source],
        }
    )
    write_table(lake_root, "precos_agropecuarios_mensal_municipio", precos_mensal_municipio)
    print(f"seeded mercado silver under {lake_root / 'silver' / 'conab'}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
