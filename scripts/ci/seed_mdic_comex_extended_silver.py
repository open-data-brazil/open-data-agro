#!/usr/bin/env python3
"""Seed minimal silver Delta for MDIC Comex extended CI (local-first)."""

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
        "mart_mdic__comex_importacao_ncm_mes",
        "mart_mdic__comex_exportacao_uf_ncm",
        "mart_mdic__comex_importacao_diesel_ncm",
    ):
        (lake_root / "gold" / mart).mkdir(parents=True, exist_ok=True)

    source = str(lake_root / "bronze/seed.parquet")
    ingested = "2026-06-26T12:00:00Z"

    import_data = pa.table(
        {
            "co_ncm": ["31021010", "31022100"],
            "ncm_descricao": ["Ureia", "Sulfato de amônio"],
            "produto_slug": ["ureia", "sulfato_amonia"],
            "data": ["2024-01-01", "2024-02-01"],
            "valor_cif_usd": ["125000000", "89000000"],
            "quantidade_kg": ["450000000", "320000000"],
            "valor_frete_usd": ["8500000", "5200000"],
            "valor_seguro_usd": ["120000", "80000"],
            "ano": ["2024", "2024"],
            "_dataset_id": ["mdic.comex-importacao-ncm-mes", "mdic.comex-importacao-ncm-mes"],
            "_ingested_at": [ingested, ingested],
            "_source_file": [source, source],
        }
    )
    write_table(lake_root, "mdic", "comex_importacao_ncm_mes", import_data)

    export_uf = pa.table(
        {
            "co_ncm": ["12019000", "12019000"],
            "ncm_descricao": ["Soja", "Soja"],
            "produto_slug": ["soja", "soja"],
            "uf": ["PR", "MT"],
            "data": ["2024-01-01", "2024-01-01"],
            "valor_fob_usd": ["498512571", "612000000"],
            "quantidade_kg": ["985317463", "1200000000"],
            "ano": ["2024", "2024"],
            "_dataset_id": ["mdic.comex-exportacao-uf-ncm", "mdic.comex-exportacao-uf-ncm"],
            "_ingested_at": [ingested, ingested],
            "_source_file": [source, source],
        }
    )
    write_table(lake_root, "mdic", "comex_exportacao_uf_ncm", export_uf)

    diesel = pa.table(
        {
            "co_ncm": ["27101921"],
            "ncm_descricao": ["Gasóleo (óleo diesel)"],
            "produto_slug": ["diesel"],
            "data": ["2024-01-01"],
            "valor_cif_usd": ["709854471"],
            "quantidade_kg": ["879308870"],
            "valor_frete_usd": ["71947963"],
            "valor_seguro_usd": ["79267"],
            "ano": ["2024"],
            "_dataset_id": ["mdic.comex-importacao-diesel-ncm"],
            "_ingested_at": [ingested],
            "_source_file": [source],
        }
    )
    write_table(lake_root, "mdic", "comex_importacao_diesel_ncm", diesel)

    print(f"seeded MDIC Comex extended silver under {lake_root / 'silver' / 'mdic'}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
