#!/usr/bin/env python3
"""Seed minimal silver Delta for ANTT pedágios CI (local-first)."""

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
    (lake_root / "gold" / "mart_antt__pracas_pedagio").mkdir(parents=True, exist_ok=True)

    source = str(lake_root / "bronze/seed.parquet")
    ingested = "2026-06-26T12:00:00Z"
    dataset_id = "antt.pracas-pedagio"

    data = pa.table(
        {
            "concessionaria": ["AUTOPISTA FERNÃO DIAS", "AUTOPISTA FERNÃO DIAS"],
            "praca_de_pedagio": ["3 (Cambuí)", "5 (Carmo da Cachoeira)"],
            "ano_do_pnv_snv": ["2007", "2007"],
            "rodovia": ["BR-381", "BR-381"],
            "uf": ["MG", "MG"],
            "km_m": ["900.9", "735.5"],
            "municipal": ["Cambuí", "Carmo da Cachoeira"],
            "tipo_de_pista": ["Principal", "Principal"],
            "sentido": ["Crescente/Decrescente", "Crescente/Decrescente"],
            "situacao": ["Ativo", "Ativo"],
            "data_da_inativacao": ["", ""],
            "latitude": ["-22.628487", "-21.5457"],
            "longitude": ["-46.07789", "-45.240203"],
            "_dataset_id": [dataset_id, dataset_id],
            "_ingested_at": [ingested, ingested],
            "_source_file": [source, source],
        }
    )
    write_table(lake_root, "antt", "pracas_pedagio", data)

    print(f"seeded ANTT pedágios silver under {lake_root / 'silver' / 'antt'}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
