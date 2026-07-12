#!/usr/bin/env python3
"""Ensure optional gold stubs exist so Phase 20 dbt reads never fail on missing files."""

from __future__ import annotations

import os
from pathlib import Path

import pyarrow as pa
import pyarrow.parquet as pq

OPTIONAL_EMPTY: dict[str, list[str]] = {
    "mart_inmet__bdmep_diario": [
        "cd_estacao",
        "data",
        "variavel",
        "valor",
        "uf",
        "ano",
        "capturado_em",
        "fonte_oficial",
        "_dataset_id",
        "_source_file",
    ],
}


def main() -> int:
    lake_root = Path(os.environ.get("LAKE_LOCAL_ROOT", "./lake"))
    for mart, cols in OPTIONAL_EMPTY.items():
        path = lake_root / "gold" / mart / "mart.parquet"
        if path.is_file():
            continue
        path.parent.mkdir(parents=True, exist_ok=True)
        table = pa.table({c: pa.array([], type=pa.string()) for c in cols})
        pq.write_table(table, path)
        print(f"wrote empty stub {path}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
