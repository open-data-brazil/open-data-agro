#!/usr/bin/env python3
"""Seed minimal silver Delta for CONAB Armazenamento e Logística CI (local-first)."""

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
    (lake_root / "gold" / "mart_conab__armazenagem").mkdir(parents=True, exist_ok=True)

    source = str(lake_root / "bronze/seed.parquet")
    ingested = "2026-06-25T12:00:00Z"

    armazenagem = pa.table(
        {
            "identificacao_armazem": ["35.0277.0001-4", "35.0277.0002-2"],
            "dsc_especie_armazem": ["CONVENCIONAL", "CONVENCIONAL"],
            "dsc_tipo_armazem": ["CONVENCIONAL", "CONVENCIONAL"],
            "dsc_tipo_entidade": ["OFICIAL", "OFICIAL"],
            "dsc_tipo_pessoa": ["PESSOA JURIDICA", "PESSOA JURIDICA"],
            "nom_municipio": ["CRUZEIRO DO SUL-AC", "CRUZEIRO DO SUL-AC"],
            "cod_ibge": ["1200203", "1200203"],
            "uf": ["AC", "AC"],
            "qtd_capacidade_estatica(t)": ["3861,0", "3861,0"],
            "qtd_capacidade_expedicao(t)": ["10", "30"],
            "qtd_capacidade_recepcao(t)": ["10", "30"],
            "latitude": ["-7.6573583", "-7.6573583"],
            "longitude": ["-72.6467501", "-72.6467501"],
            "nome_armazenador": ["COMPANHIA DE ARMAZENS GERAIS", "COMPANHIA DE ARMAZENS GERAIS"],
            "endereco": ["AVENIDA 25 DE AGOSTO", "AVENIDA 25 DE AGOSTO"],
            "email": ["", ""],
            "_dataset_id": ["conab.armazenagem", "conab.armazenagem"],
            "_ingested_at": [ingested, ingested],
            "_source_file": [source, source],
        }
    )
    write_table(lake_root, "conab", "armazenagem", armazenagem)
    print(f"seeded armazenagem silver under {lake_root / 'silver' / 'conab' / 'armazenagem'}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
