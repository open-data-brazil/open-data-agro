#!/usr/bin/env python3
"""Seed minimal silver Delta for IBGE PAM CI (local-first)."""

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

    source = str(lake_root / "bronze/seed.parquet")
    ingested = "2026-06-25T12:00:00Z"

    pam = pa.table(
        {
            "sidra_tabela": ["1612", "1612", "1612"],
            "codigo_ibge": ["5100102", "5100102", "5100102"],
            "municipio": ["Acorizal - MT", "Acorizal - MT", "Acorizal - MT"],
            "ano": ["2015", "2015", "2015"],
            "variavel_codigo": ["109", "216", "214"],
            "variavel": ["Área plantada", "Área colhida", "Quantidade produzida"],
            "produto_codigo": ["2713", "2713", "2713"],
            "produto": ["Soja (em grão)", "Soja (em grão)", "Soja (em grão)"],
            "valor": ["12000", "11500", "45000"],
            "unidade_codigo": ["1006", "1006", "1017"],
            "unidade": ["Hectares", "Hectares", "Toneladas"],
            "_dataset_id": [
                "ibge.pam-area-quantidade",
                "ibge.pam-area-quantidade",
                "ibge.pam-area-quantidade",
            ],
            "_ingested_at": [ingested, ingested, ingested],
            "_source_file": [source, source, source],
        }
    )
    write_table(lake_root, "ibge", "pam_area_quantidade", pam)

    print(f"seeded IBGE PAM silver under {lake_root / 'silver' / 'ibge'}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
