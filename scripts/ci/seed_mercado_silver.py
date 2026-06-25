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

    source = str(lake_root / "bronze/seed.parquet")
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
            "_ingested_at": ["2026-06-25T12:00:00Z", "2026-06-25T12:00:00Z"],
            "_source_file": [source, source],
        }
    )
    write_table(lake_root, "oferta_demanda", oferta)
    print(f"seeded mercado silver under {lake_root / 'silver' / 'conab' / 'oferta_demanda'}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
