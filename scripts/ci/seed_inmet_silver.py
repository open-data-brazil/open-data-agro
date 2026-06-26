#!/usr/bin/env python3
"""Seed minimal silver Delta for INMET climate CI (local-first)."""

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


def meta_cols(dataset_id: str, source: str, ingested: str, n: int) -> dict[str, list[str]]:
    return {
        "_dataset_id": [dataset_id] * n,
        "_ingested_at": [ingested] * n,
        "_source_file": [source] * n,
    }


def main() -> int:
    lake_root = Path(os.environ.get("LAKE_LOCAL_ROOT", "/tmp/open-data-agro-lake"))
    lake_root.mkdir(parents=True, exist_ok=True)
    (lake_root / "gold").mkdir(parents=True, exist_ok=True)
    for mart in [
        "mart_inmet__estacoes_automaticas",
        "mart_inmet__estacoes_convencionais",
        "mart_inmet__bdmep_diario",
        "mart_inmet__bdmep_mensal",
    ]:
        (lake_root / "gold" / mart).mkdir(parents=True, exist_ok=True)

    source = str(lake_root / "bronze/seed.parquet")
    ingested = "2026-06-25T12:00:00Z"

    write_table(
        lake_root,
        "inmet",
        "estacoes_automaticas",
        pa.table(
            {
                "cd_estacao": ["A901", "A801"],
                "nome": ["SINOP", "PORTO ALEGRE"],
                "latitude": ["-11.8644", "-30.0500"],
                "longitude": ["-55.5022", "-51.1600"],
                "uf": ["MT", "RS"],
                "situacao": ["OPERANTE", "OPERANTE"],
                **meta_cols("inmet.estacoes-automaticas", source, ingested, 2),
            }
        ),
    )

    write_table(
        lake_root,
        "inmet",
        "estacoes_convencionais",
        pa.table(
            {
                "cd_estacao": ["B407", "B488"],
                "nome": ["CUIABA", "CAMPO GRANDE"],
                "latitude": ["-15.65", "-20.45"],
                "longitude": ["-56.10", "-54.67"],
                "uf": ["MT", "MS"],
                "situacao": ["OPERANTE", "OPERANTE"],
                "regiao": ["Centro-Oeste", "Centro-Oeste"],
                "altitude": ["182", "530"],
                **meta_cols("inmet.estacoes-convencionais", source, ingested, 2),
            }
        ),
    )

    diario = {
        "cd_estacao": ["A901"] * 4,
        "data": ["2023-01-01", "2023-01-01", "2023-01-02", "2023-01-02"],
        "variavel": ["precipitacao", "temperatura_ar", "precipitacao", "temperatura_ar"],
        "valor": ["0.5", "25.1", "2.0", "26.05"],
        "uf": ["MT"] * 4,
        "ano": ["2023"] * 4,
        **meta_cols("inmet.bdmep-diario", source, ingested, 4),
    }
    write_table(lake_root, "inmet", "bdmep_diario", pa.table(diario))

    mensal = {
        "cd_estacao": ["A901", "A901"],
        "mes": ["2023-01", "2023-01"],
        "variavel": ["precipitacao", "temperatura_ar"],
        "valor": ["2.5", "25.5"],
        "uf": ["MT", "MT"],
        "ano": ["2023", "2023"],
        **meta_cols("inmet.bdmep-mensal", source, ingested, 2),
    }
    write_table(lake_root, "inmet", "bdmep_mensal", pa.table(mensal))

    print(f"seeded INMET silver under {lake_root / 'silver' / 'inmet'}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
