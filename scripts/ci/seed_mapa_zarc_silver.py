#!/usr/bin/env python3
"""Seed minimal silver Delta for MAPA ZARC tábua de risco CI (local-first)."""

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
    (lake_root / "gold" / "mart_mapa__zarc_tabua_risco").mkdir(parents=True, exist_ok=True)

    source = str(lake_root / "bronze/seed.parquet")
    ingested = "2026-06-26T12:00:00Z"
    dataset_id = "mapa.zarc-tabua-risco"

    dec_cols = {f"dec{i}": ["20", "20", "0"] for i in range(1, 37)}

    data = pa.table(
        {
            "Nome_cultura": ["Soja", "Milho", "Trigo"],
            "SafraIni": ["2025", "2025", "2025"],
            "SafraFin": ["2026", "2026", "2026"],
            "Cod_Cultura": ["12016000000021", "12012000000011", "12017100000031"],
            "Cod_Ciclo": ["21", "21", "22"],
            "Cod_Solo": ["2", "2", "16"],
            "geocodigo": ["4106902", "4106902", "3157609"],
            "UF": ["PR", "PR", "MG"],
            "municipio": ["Curitiba", "Curitiba", "Santa Fé de Minas"],
            "Cod_Clima": ["0", "0", "0"],
            "Nome_Clima": ["Não se aplica", "Não se aplica", "Não se aplica"],
            "Cod_Outros_Manejos": ["1", "1", "2"],
            "Nome_Outros_Manejos": ["Sequeiro", "Sequeiro", "Irrigado"],
            "Produtividade": ["", "", ""],
            "Cod_NM": ["041001", "041001", "015284"],
            "Cod_Munic": ["01", "01", "02"],
            "Cod_Meso": ["001", "001", "006"],
            "Cod_Micro": ["001", "001", "006"],
            "Portaria": [
                "Port.447_de_29-10-2025",
                "Port.447_de_29-10-2025",
                "Port.447_de_29-10-2025",
            ],
            **dec_cols,
            "_dataset_id": [dataset_id, dataset_id, dataset_id],
            "_ingested_at": [ingested, ingested, ingested],
            "_source_file": [source, source, source],
        }
    )
    write_table(lake_root, "mapa", "zarc_tabua_risco", data)

    print(f"seeded MAPA ZARC silver under {lake_root / 'silver' / 'mapa'}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
