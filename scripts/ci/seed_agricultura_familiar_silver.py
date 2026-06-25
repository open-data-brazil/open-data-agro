#!/usr/bin/env python3
"""Seed minimal silver Delta for CONAB Agricultura Familiar CI (local-first)."""

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
    for mart in (
        "mart_conab__alimenta_brasil_entregas",
        "mart_conab__alimenta_brasil_propostas",
    ):
        (lake_root / "gold" / mart).mkdir(parents=True, exist_ok=True)

    source = str(lake_root / "bronze/seed.parquet")
    ingested = "2026-06-25T12:00:00Z"

    entregas = pa.table(
        {
            "ano_entrega": ["2023,0", "2023,0"],
            "mes_entrega": ["1,0", "1,0"],
            "municipio": ["BRASILIA-DF", "GOIANIA-GO"],
            "uf": ["DF", "GO"],
            "regiao": ["CENTRO-OESTE", "CENTRO-OESTE"],
            "sexo": ["FEMININO", "MASCULINO"],
            "ds_unidade_medida": ["kg", "kg"],
            "qtd_entregue": ["3854,8", "2100,0"],
            "valor_entregue": ["19621,143", "10500,0"],
            "_dataset_id": ["conab.alimenta-brasil-entregas", "conab.alimenta-brasil-entregas"],
            "_ingested_at": [ingested, ingested],
            "_source_file": [source, source],
        }
    )
    write_table(lake_root, "conab", "alimenta_brasil_entregas", entregas)

    propostas = pa.table(
        {
            "ano": ["2015,0", "2015,0"],
            "mes": ["1,0", "2,0"],
            "municipio": ["JARAGUA-GO", "BRASILIA-DF"],
            "cod_ibge": ["5211800", "5300108"],
            "uf": ["GO", "DF"],
            "regiao": ["CENTRO-OESTE", "CENTRO-OESTE"],
            "valor_formalizado": ["455715,0", "120000,0"],
            "valor_executado": ["0,0", "50000,0"],
            "valor_devolvido": ["0,0", "0,0"],
            "_dataset_id": ["conab.alimenta-brasil-propostas", "conab.alimenta-brasil-propostas"],
            "_ingested_at": [ingested, ingested],
            "_source_file": [source, source],
        }
    )
    write_table(lake_root, "conab", "alimenta_brasil_propostas", propostas)

    print(f"seeded agricultura familiar silver under {lake_root / 'silver' / 'conab'}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
