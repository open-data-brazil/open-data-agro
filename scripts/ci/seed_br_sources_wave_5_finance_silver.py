#!/usr/bin/env python3
"""Seed minimal silver Delta for BR sources wave 5 finance CI (Phase 55)."""

from __future__ import annotations

import os
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
    gold_tables = [
        "mart_bcb__cim_agro_credito_rural",
        "mart_bndes__desembolsos_linhas_agro",
        "mart_anp__etanol_precos",
    ]
    for table in gold_tables:
        (lake_root / "gold" / table).mkdir(parents=True, exist_ok=True)

    source = str(lake_root / "bronze/seed.parquet")
    ingested = "2026-06-27T12:00:00Z"

    cim_agro = pa.table(
        {
            "sgs_codigo": ["21087", "21087"],
            "data": ["2024-01-01", "2024-02-01"],
            "valor": ["3.93", "4.42"],
            "ano": ["2024", "2024"],
            "_dataset_id": ["bcb.cim-agro-credito-rural", "bcb.cim-agro-credito-rural"],
            "_ingested_at": [ingested, ingested],
            "_source_file": [source, source],
        }
    )
    write_table(lake_root, "bcb", "cim_agro_credito_rural", cim_agro)

    desembolsos = pa.table(
        {
            "ano": ["2024", "2024"],
            "mes": ["1", "2"],
            "bndes_finem": ["100,5", "110,0"],
            "bndes_exim": ["20,1", "22,0"],
            "bndes_mercado_de_capitais": ["30,2", "28,5"],
            "bndes_nao_reembolsavel": ["0,0", "1,0"],
            "bndes_microcredito": ["0,0", "0,0"],
            "bndes_prestacao_de_garantia": ["0,0", "0,0"],
            "bndes_finame": ["5,0", "6,5"],
            "_dataset_id": ["bndes.desembolsos-linhas-agro", "bndes.desembolsos-linhas-agro"],
            "_ingested_at": [ingested, ingested],
            "_source_file": [source, source],
        }
    )
    write_table(lake_root, "bndes", "desembolsos_linhas_agro", desembolsos)

    etanol = pa.table(
        {
            "DATA INICIAL": ["2024-01-01", "2024-01-01"],
            "DATA FINAL": ["2024-01-07", "2024-01-07"],
            "ESTADO": ["SAO PAULO", "MINAS GERAIS"],
            "MUNICÍPIO": ["SAO PAULO", "UBERLANDIA"],
            "PRODUTO": ["ETANOL HIDRATADO", "ETANOL HIDRATADO"],
            "NÚMERO DE POSTOS PESQUISADOS": ["10", "8"],
            "UNIDADE DE MEDIDA": ["R$/l", "R$/l"],
            "PREÇO MÉDIO REVENDA": ["3.50", "3.45"],
            "DESVIO PADRÃO REVENDA": ["0.10", "0.08"],
            "PREÇO MÍNIMO REVENDA": ["3.30", "3.25"],
            "PREÇO MÁXIMO REVENDA": ["3.70", "3.60"],
            "COEF DE VARIAÇÃO REVENDA": ["0.03", "0.02"],
            "_dataset_id": ["anp.etanol-precos", "anp.etanol-precos"],
            "_ingested_at": [ingested, ingested],
            "_source_file": [source, source],
        }
    )
    write_table(lake_root, "anp", "etanol_precos", etanol)

    print(f"seeded wave 5 finance silver under {lake_root / 'silver'}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
