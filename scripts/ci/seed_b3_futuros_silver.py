#!/usr/bin/env python3
"""Seed minimal silver Delta for B3 futuros CI (local-first)."""

from __future__ import annotations

import os
import sys
from pathlib import Path

import pyarrow as pa
from deltalake import write_deltalake


def write_table(root: Path, table: str, data: pa.Table) -> None:
    path = root / "silver" / "b3" / table
    path.parent.mkdir(parents=True, exist_ok=True)
    write_deltalake(str(path), data, mode="overwrite")


def seed_dataset(
    root: Path,
    table: str,
    dataset_id: str,
    commodity: str,
    symbols: list[str],
    prices: list[str],
) -> None:
    source = str(root / "bronze/seed.parquet")
    ingested = "2026-06-26T12:00:00Z"
    data = pa.table(
        {
            "refdate": ["2025-06-25"] * len(symbols),
            "symbol": symbols,
            "commodity": [commodity] * len(symbols),
            "maturity_code": [s[len(commodity) :] for s in symbols],
            "previous_price": prices,
            "price": prices,
            "currency": ["USD" if commodity == "SOY" else "BRL"] * len(symbols),
            "price_change": ["0"] * len(symbols),
            "_dataset_id": [dataset_id] * len(symbols),
            "_ingested_at": [ingested] * len(symbols),
            "_source_file": [source] * len(symbols),
        }
    )
    write_table(root, table, data)


def main() -> int:
    lake_root = Path(os.environ.get("LAKE_LOCAL_ROOT", "/tmp/open-data-agro-lake"))
    lake_root.mkdir(parents=True, exist_ok=True)
    for mart in (
        "mart_b3__futuro_soja",
        "mart_b3__futuro_milho",
        "mart_b3__futuro_boi",
    ):
        (lake_root / "gold" / mart).mkdir(parents=True, exist_ok=True)

    seed_dataset(
        lake_root,
        "futuro_soja",
        "b3.futuro-soja",
        "SOY",
        ["SOYH26", "SOYK26"],
        ["390", "396.4"],
    )
    seed_dataset(
        lake_root,
        "futuro_milho",
        "b3.futuro-milho",
        "CCM",
        ["CCMU25", "CCMX25"],
        ["63.08", "64.10"],
    )
    seed_dataset(
        lake_root,
        "futuro_boi",
        "b3.futuro-boi",
        "BGI",
        ["BGIV25", "BGIX25"],
        ["335.75", "339.30"],
    )

    print(f"seeded B3 futuros silver under {lake_root / 'silver' / 'b3'}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
