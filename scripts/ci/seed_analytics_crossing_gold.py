#!/usr/bin/env python3
"""Seed minimal gold Parquet fixtures for Phase 20 analytics-crossing CI."""

from __future__ import annotations

import os
from pathlib import Path

import pyarrow as pa
import pyarrow.parquet as pq


def write_mart(root: Path, name: str, table: pa.Table) -> None:
    path = root / "gold" / name
    path.mkdir(parents=True, exist_ok=True)
    pq.write_table(table, path / "mart.parquet")


def main() -> int:
    lake_root = Path(os.environ.get("LAKE_LOCAL_ROOT", "/tmp/analytics-crossing-lake"))
    lake_root.mkdir(parents=True, exist_ok=True)

    ingested = "2026-06-26T12:00:00Z"
    fonte = "https://example.local/seed"

    write_mart(
        lake_root,
        "mart_cepea__soja_paranagua",
        pa.table(
            {
                "produto": ["soja", "soja", "soja", "soja"],
                "praca": ["Paranaguá"] * 4,
                "data": ["2024-01-02", "2024-01-03", "2024-01-04", "2024-01-05"],
                "preco_rs_sc": ["120.00", "121.50", "119.80", "122.00"],
                "variacao_dia_pct": ["0.10", "1.25", "-1.40", "1.84"],
                "preco_usd_sc": ["24.00", "24.30", "", "24.40"],
                "ano": ["2024"] * 4,
                "capturado_em": [ingested] * 4,
                "fonte_oficial": [fonte] * 4,
                "_dataset_id": ["cepea.soja-paranagua"] * 4,
                "_source_file": ["seed.parquet"] * 4,
            }
        ),
    )

    write_mart(
        lake_root,
        "mart_bcb__sgs_ptax_usd_venda",
        pa.table(
            {
                "sgs_codigo": ["1", "1", "1", "1", "1"],
                "data": [
                    "2024-01-01",
                    "2024-01-02",
                    "2024-01-03",
                    "2024-01-04",
                    "2024-01-05",
                ],
                "valor": ["4.95", "5.00", "5.00", "5.05", "5.00"],
                "ano": ["2024"] * 5,
                "capturado_em": [ingested] * 5,
                "fonte_oficial": [fonte] * 5,
                "_dataset_id": ["bcb.sgs-ptax-usd-venda"] * 5,
                "_source_file": ["seed.parquet"] * 5,
            }
        ),
    )

    write_mart(
        lake_root,
        "mart_b3__futuro_soja",
        pa.table(
            {
                "refdate": ["2024-01-02", "2024-01-03", "2024-01-04", "2024-01-05"],
                "symbol": ["SOYH24", "SOYH24", "SOYH24", "SOYH24"],
                "commodity": ["SOY"] * 4,
                "maturity_code": ["H24"] * 4,
                "previous_price": ["12.0", "12.1", "12.2", "12.1"],
                "price": ["12.1", "12.2", "12.0", "12.3"],
                "currency": ["USD"] * 4,
                "price_change": ["0.1", "0.1", "-0.2", "0.3"],
                "capturado_em": [ingested] * 4,
                "fonte_oficial": [fonte] * 4,
                "_dataset_id": ["b3.futuro-soja"] * 4,
                "_source_file": ["seed.parquet"] * 4,
            }
        ),
    )

    write_mart(
        lake_root,
        "mart_mdic__comex_exportacao_uf_ncm",
        pa.table(
            {
                "co_ncm": ["12019000", "12019000"],
                "ncm_descricao": ["Soja", "Soja"],
                "produto_slug": ["soja", "soja"],
                "uf": ["MT", "PR"],
                # Oct-2023 → as_of = 2023-12-15 so values are knowable on Jan-2024 spine
                "data": ["2023-10-01", "2023-10-01"],
                "valor_fob_usd": ["1000000", "500000"],
                "quantidade_kg": ["2000000", "1000000"],
                "ano": ["2023", "2023"],
                "capturado_em": [ingested, ingested],
                "fonte_oficial": [fonte, fonte],
                "_dataset_id": ["mdic.comex-exportacao-uf-ncm"] * 2,
                "_source_file": ["seed.parquet"] * 2,
            }
        ),
    )

    write_mart(
        lake_root,
        "mart_conab__frete",
        pa.table(
            {
                "fonte": ["CONTRATO", "CONTRATO"],
                "municipio_origem": ["LUCAS DO RIO VERDE-MT", "SORRISO-MT"],
                "cod_ibge_origem": ["5105259", "5107925"],
                "uf_origem": ["MT", "MT"],
                "municipio_destino": ["PARANAGUA-PR", "PARANAGUA-PR"],
                "cod_ibge_destino": ["4118204", "4118204"],
                "uf_destino": ["PR", "PR"],
                "ano": ["2023", "2023"],
                "mes": ["10", "10"],
                "distancia_km": ["1800", "1750"],
                "valor_frete_tonelada": ["350,5", "340,0"],
                "valor_tonelada_km": ["0,1947", "0,1943"],
                "capturado_em": [ingested, ingested],
                "fonte_oficial": [fonte, fonte],
                "_dataset_id": ["conab.frete"] * 2,
                "_source_file": ["seed.parquet"] * 2,
            }
        ),
    )

    write_mart(
        lake_root,
        "mart_conab__precos_mensal_municipio",
        pa.table(
            {
                "produto": ["SOJA", "SOJA"],
                "classificacao_produto": ["NÃO INFORMADO", "NÃO INFORMADO"],
                "id_produto": ["1", "1"],
                "municipio": ["SORRISO-MT", "LUCAS DO RIO VERDE-MT"],
                "cod_ibge": ["5107925", "5105259"],
                "uf": ["MT", "MT"],
                "regiao": ["CENTRO-OESTE", "CENTRO-OESTE"],
                "ano": ["2023", "2023"],
                "mes": ["10", "10"],
                "nivel_comercializacao": ["PREÇO PAGO PELO PRODUTOR"] * 2,
                "valor_produto_kg": ["1,85", "1,90"],
                "capturado_em": [ingested, ingested],
                "fonte_oficial": [fonte, fonte],
                "_dataset_id": ["conab.precos-agropecuarios-mensal-municipio"] * 2,
                "_source_file": ["seed.parquet"] * 2,
            }
        ),
    )

    write_mart(
        lake_root,
        "mart_ibge__localidades_municipios",
        pa.table(
            {
                "codigo_ibge": ["5107925", "5105259", "4118204"],
                "nome": ["Sorriso", "Lucas do Rio Verde", "Paranaguá"],
                "sigla_uf": ["MT", "MT", "PR"],
                "codigo_uf": ["51", "51", "41"],
                "codigo_regiao": ["5", "5", "4"],
                "nome_regiao": ["Centro-Oeste", "Centro-Oeste", "Sul"],
                "capturado_em": [ingested] * 3,
                "fonte_oficial": [fonte] * 3,
                "_dataset_id": ["ibge.localidades-municipios"] * 3,
                "_source_file": ["seed.parquet"] * 3,
            }
        ),
    )

    # Empty climate stub — exercises NULL climate features without blocking the build.
    write_mart(
        lake_root,
        "mart_inmet__bdmep_diario",
        pa.table(
            {
                "cd_estacao": pa.array([], type=pa.string()),
                "data": pa.array([], type=pa.string()),
                "variavel": pa.array([], type=pa.string()),
                "valor": pa.array([], type=pa.string()),
                "uf": pa.array([], type=pa.string()),
                "ano": pa.array([], type=pa.string()),
                "capturado_em": pa.array([], type=pa.string()),
                "fonte_oficial": pa.array([], type=pa.string()),
                "_dataset_id": pa.array([], type=pa.string()),
                "_source_file": pa.array([], type=pa.string()),
            }
        ),
    )

    print(f"seeded analytics-crossing gold under {lake_root / 'gold'}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
