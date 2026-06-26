#!/usr/bin/env python3
"""Seed minimal silver Delta for IBGE Localidades CI (local-first)."""

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
    (lake_root / "gold" / "mart_ibge__localidades_municipios").mkdir(parents=True, exist_ok=True)
    (lake_root / "gold" / "mart_ibge__localidades_ufs").mkdir(parents=True, exist_ok=True)

    source = str(lake_root / "bronze/seed.parquet")
    ingested = "2026-06-25T12:00:00Z"

    municipios = pa.table(
        {
            "codigo_ibge": ["3550308", "5100102"],
            "nome": ["São Paulo", "Abadia dos Dourados"],
            "sigla_uf": ["SP", "MG"],
            "codigo_uf": ["35", "31"],
            "codigo_regiao": ["3", "3"],
            "nome_regiao": ["Sudeste", "Sudeste"],
            "_dataset_id": [
                "ibge.localidades-municipios",
                "ibge.localidades-municipios",
            ],
            "_ingested_at": [ingested, ingested],
            "_source_file": [source, source],
        }
    )

    ufs = pa.table(
        {
            "codigo_uf": ["35", "51"],
            "sigla_uf": ["SP", "MT"],
            "nome": ["São Paulo", "Mato Grosso"],
            "codigo_regiao": ["3", "5"],
            "sigla_regiao": ["SE", "CO"],
            "nome_regiao": ["Sudeste", "Centro-Oeste"],
            "_dataset_id": ["ibge.localidades-ufs", "ibge.localidades-ufs"],
            "_ingested_at": [ingested, ingested],
            "_source_file": [source, source],
        }
    )

    write_table(lake_root, "ibge", "localidades_municipios", municipios)
    write_table(lake_root, "ibge", "localidades_ufs", ufs)
    print(f"seeded IBGE localidades silver under {lake_root / 'silver' / 'ibge'}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
