#!/usr/bin/env python3
"""Seed minimal silver Delta for industry sources wave 5 CI (Phase 56)."""

from __future__ import annotations

import os
from pathlib import Path

import pyarrow as pa
from deltalake import write_deltalake


def write_table(root: Path, agency: str, table: str, data: pa.Table) -> None:
    path = root / "silver" / agency / table
    path.parent.mkdir(parents=True, exist_ok=True)
    write_deltalake(str(path), data, mode="overwrite")


def seed_abiove_stat(
    root: Path,
    table: str,
    dataset_id: str,
    section: str,
    row_label: str,
    period: str,
    metric: str,
    value: str,
) -> None:
    source = str(root / "bronze/seed.parquet")
    ingested = "2026-06-27T12:00:00Z"
    data = pa.table(
        {
            "section": [section, section],
            "row_label": [row_label, "Fev"],
            "period": [period, "Fev"],
            "metric": [metric, metric],
            "value": [value, "120.5"],
            "report_updated_at": ["2026-05-01", "2026-05-01"],
            "_dataset_id": [dataset_id, dataset_id],
            "_ingested_at": [ingested, ingested],
            "_source_file": [source, source],
        }
    )
    write_table(root, "abiove", table, data)


def seed_b3_futuro(
    root: Path,
    table: str,
    dataset_id: str,
    commodity: str,
    symbols: list[str],
    prices: list[str],
) -> None:
    source = str(root / "bronze/seed.parquet")
    ingested = "2026-06-27T12:00:00Z"
    data = pa.table(
        {
            "refdate": ["2025-06-27"] * len(symbols),
            "symbol": symbols,
            "commodity": [commodity] * len(symbols),
            "maturity_code": [s[len(commodity) :] for s in symbols],
            "previous_price": prices,
            "price": prices,
            "currency": ["USD" if commodity == "ICF" else "BRL"] * len(symbols),
            "price_change": ["0"] * len(symbols),
            "_dataset_id": [dataset_id] * len(symbols),
            "_ingested_at": [ingested] * len(symbols),
            "_source_file": [source] * len(symbols),
        }
    )
    write_table(root, "b3", table, data)


def main() -> int:
    lake_root = Path(os.environ.get("LAKE_LOCAL_ROOT", "/tmp/open-data-agro-lake"))
    lake_root.mkdir(parents=True, exist_ok=True)
    gold_tables = [
        "mart_abiove__balanco_complexo_soja",
        "mart_abiove__exportacoes_complexo_soja",
        "mart_abiove__capacidade_instalada_esmagamento",
        "mart_b3__futuro_cafe",
        "mart_b3__futuro_acucar",
    ]
    for table in gold_tables:
        (lake_root / "gold" / table).mkdir(parents=True, exist_ok=True)

    seed_abiove_stat(
        lake_root,
        "balanco_complexo_soja",
        "abiove.balanco-complexo-soja",
        "exportacoes_soja_grao",
        "Jan",
        "Jan",
        "valor_fob_usd_mil",
        "433558.617",
    )
    seed_abiove_stat(
        lake_root,
        "exportacoes_complexo_soja",
        "abiove.exportacoes-complexo-soja",
        "matéria-prima",
        "Óleo de soja",
        "2024",
        "volume_m3",
        "6665849.291",
    )
    seed_abiove_stat(
        lake_root,
        "capacidade_instalada_esmagamento",
        "abiove.capacidade-instalada-esmagamento",
        "esmagamento_mil_t",
        "1.3. Processamento",
        "Ago (3)",
        "volume_mil_t",
        "3035.015",
    )
    seed_b3_futuro(
        lake_root,
        "futuro_cafe",
        "b3.futuro-cafe",
        "ICF",
        ["ICFH26", "ICFK26"],
        ["212.0", "215.5"],
    )
    seed_b3_futuro(
        lake_root,
        "futuro_acucar",
        "b3.futuro-acucar",
        "CNL",
        ["CNLF26", "CNLH26"],
        ["18.7", "19.2"],
    )

    print(f"seeded wave 5 industry silver under {lake_root / 'silver'}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
