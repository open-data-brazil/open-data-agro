#!/usr/bin/env python3
"""Seed minimal silver Delta for CEPEA indicators CI (local-first)."""

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


def indicador_rows(
    produto: str,
    praca: str,
    dataset_id: str,
    dates: list[str],
    precos: list[str],
    source: str,
    ingested: str,
) -> pa.Table:
    return pa.table(
        {
            "produto": [produto] * len(dates),
            "praca": [praca] * len(dates),
            "data": dates,
            "preco_rs_sc": precos,
            "variacao_dia_pct": ["0.15", "0.32"][: len(dates)],
            "preco_usd_sc": ["29.80", "25.82"][: len(dates)],
            "ano": [d[:4] for d in dates],
            "_dataset_id": [dataset_id] * len(dates),
            "_ingested_at": [ingested] * len(dates),
            "_source_file": [source] * len(dates),
        }
    )


def main() -> int:
    lake_root = Path(os.environ.get("LAKE_LOCAL_ROOT", "/tmp/open-data-agro-lake"))
    lake_root.mkdir(parents=True, exist_ok=True)
    (lake_root / "gold").mkdir(parents=True, exist_ok=True)
    for mart in [
        "mart_cepea__soja_paranagua",
        "mart_cepea__soja_parana",
        "mart_cepea__milho",
        "mart_cepea__boi_gordo",
    ]:
        (lake_root / "gold" / mart).mkdir(parents=True, exist_ok=True)

    source = str(lake_root / "bronze/seed.parquet")
    ingested = "2026-06-25T12:00:00Z"
    dates = ["2010-01-04", "2024-01-31"]
    precos = ["52.30", "124.58"]

    tables = [
        ("soja_paranagua", "soja", "Paranaguá", "cepea.soja-paranagua"),
        ("soja_parana", "soja", "Paraná", "cepea.soja-parana"),
        ("milho", "milho", "Campinas", "cepea.milho"),
        ("boi_gordo", "boi-gordo", "São Paulo", "cepea.boi-gordo"),
    ]
    for table, produto, praca, dataset_id in tables:
        write_table(
            lake_root,
            "cepea",
            table,
            indicador_rows(produto, praca, dataset_id, dates, precos, source, ingested),
        )

    print(f"seeded CEPEA silver under {lake_root / 'silver' / 'cepea'}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
