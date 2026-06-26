#!/usr/bin/env python3
"""Seed minimal silver Delta for MDIC Comex CI (local-first)."""

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


def comex_rows(
    dataset_id: str,
    dates: list[str],
    co_ncms: list[str],
    slugs: list[str],
    fobs: list[str],
    kgs: list[str],
    source: str,
    ingested: str,
) -> pa.Table:
    return pa.table(
        {
            "co_ncm": co_ncms,
            "ncm_descricao": ["Soja"] * len(dates),
            "produto_slug": slugs,
            "data": dates,
            "valor_fob_usd": fobs,
            "quantidade_kg": kgs,
            "ano": [d[:4] for d in dates],
            "_dataset_id": [dataset_id] * len(dates),
            "_ingested_at": [ingested] * len(dates),
            "_source_file": [source] * len(dates),
        }
    )


def main() -> int:
    lake_root = Path(os.environ.get("LAKE_LOCAL_ROOT", "/tmp/open-data-agro-lake"))
    lake_root.mkdir(parents=True, exist_ok=True)
    (lake_root / "gold" / "mart_mdic__comex_exportacao_ncm_mes").mkdir(parents=True, exist_ok=True)

    source = str(lake_root / "bronze/seed.parquet")
    ingested = "2026-06-26T12:00:00Z"
    dataset_id = "mdic.comex-exportacao-ncm-mes"

    data = comex_rows(
        dataset_id,
        ["2024-01-01", "2024-02-01"],
        ["12019000", "12019000"],
        ["soja", "soja"],
        ["1454912473", "2919570191"],
        ["2854860075", "6608133608"],
        source,
        ingested,
    )
    write_table(lake_root, "mdic", "comex_exportacao_ncm_mes", data)

    print(f"seeded MDIC Comex silver under {lake_root / 'silver' / 'mdic'}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
