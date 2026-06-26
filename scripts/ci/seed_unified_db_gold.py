#!/usr/bin/env python3
"""Seed minimal gold Parquet marts for unified PostgreSQL CI (local-first)."""

from __future__ import annotations

import os
import sys
from pathlib import Path

import pyarrow as pa
import pyarrow.parquet as pq


def write_mart(lake_root: Path, mart_dir: str, table: pa.Table) -> None:
    path = lake_root / "gold" / mart_dir
    path.mkdir(parents=True, exist_ok=True)
    pq.write_table(table, path / "mart.parquet")


def main() -> int:
    lake_root = Path(os.environ.get("LAKE_LOCAL_ROOT", "/tmp/unified-db-ci-lake"))
    lake_root.mkdir(parents=True, exist_ok=True)

    estimativa = pa.table(
        {
            "ano_agricola": ["2025/26", "2025/26"],
            "safra": ["UNICA", "UNICA"],
            "uf": ["PR", "MT"],
            "produto": ["SOJA", "MILHO"],
            "producao_mil_t": ["100", "120"],
            "capturado_em": ["2026-06-26T12:00:00Z", "2026-06-26T12:00:00Z"],
            "fonte_oficial": [
                "https://portaldeinformacoes.conab.gov.br/download-arquivos.html",
                "https://portaldeinformacoes.conab.gov.br/download-arquivos.html",
            ],
        }
    )
    write_mart(lake_root, "mart_conab__estimativa_graos", estimativa)

    municipios = pa.table(
        {
            "codigo_ibge": ["5107925", "3170701"],
            "nome": ["SORRISO", "VARGINHA"],
            "sigla_uf": ["MT", "MG"],
            "capturado_em": ["2026-06-26T12:00:00Z", "2026-06-26T12:00:00Z"],
            "fonte_oficial": [
                "https://servicodados.ibge.gov.br/api/docs/localidades",
                "https://servicodados.ibge.gov.br/api/docs/localidades",
            ],
        }
    )
    write_mart(lake_root, "mart_ibge__localidades_municipios", municipios)

    print(f"seeded unified-db gold marts under {lake_root / 'gold'}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
