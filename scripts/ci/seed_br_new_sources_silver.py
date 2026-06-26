#!/usr/bin/env python3
"""Seed minimal silver Delta for BR new sources CI (Phase 37)."""

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
    for mart in ("mart_ibge__lspa_area_producao", "mart_bcb__sgs_selic"):
        (lake_root / "gold" / mart).mkdir(parents=True, exist_ok=True)

    source = str(lake_root / "bronze/seed.parquet")
    ingested = "2026-06-26T12:00:00Z"

    lspa = pa.table(
        {
            "sidra_tabela": ["6588", "6588", "6588"],
            "codigo_uf": ["43", "43", "43"],
            "uf": ["Rio Grande do Sul", "Rio Grande do Sul", "Rio Grande do Sul"],
            "mes": ["202401", "202401", "202401"],
            "variavel_codigo": ["109", "216", "35"],
            "variavel": ["Área plantada", "Área colhida", "Produção"],
            "produto_codigo": ["39443", "39443", "39443"],
            "produto": ["1.17 Soja", "1.17 Soja", "1.17 Soja"],
            "produto_slug": ["soja", "soja", "soja"],
            "valor": ["6716122", "6500000", "13647103"],
            "unidade_codigo": ["1006", "1006", "1017"],
            "unidade": ["Hectares", "Hectares", "Toneladas"],
            "_dataset_id": [
                "ibge.lspa-area-producao",
                "ibge.lspa-area-producao",
                "ibge.lspa-area-producao",
            ],
            "_ingested_at": [ingested, ingested, ingested],
            "_source_file": [source, source, source],
        }
    )
    write_table(lake_root, "ibge", "lspa_area_producao", lspa)

    selic = pa.table(
        {
            "sgs_codigo": ["11", "11"],
            "data": ["2024-01-02", "2024-01-03"],
            "valor": ["11.75", "11.75"],
            "ano": ["2024", "2024"],
            "_dataset_id": ["bcb.sgs-selic", "bcb.sgs-selic"],
            "_ingested_at": [ingested, ingested],
            "_source_file": [source, source],
        }
    )
    write_table(lake_root, "bcb", "sgs_selic", selic)

    print(f"seeded BR new sources silver under {lake_root / 'silver'}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
