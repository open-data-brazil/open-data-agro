"""Minimal IBGE municipios gold mart for validate_codigo_ibge in CI seeds."""

from __future__ import annotations

from pathlib import Path

import pyarrow as pa
import pyarrow.parquet as pq

# Codes referenced by mercado / PAM / armazenamento / abastecimento / PAA CI seeds.
REFERENCE_MUNICIPIOS: list[dict[str, str]] = [
    {
        "codigo_ibge": "1200203",
        "nome": "Cruzeiro do Sul",
        "sigla_uf": "AC",
        "codigo_uf": "12",
        "codigo_regiao": "1",
        "nome_regiao": "Norte",
    },
    {
        "codigo_ibge": "1500800",
        "nome": "Ananindeua",
        "sigla_uf": "PA",
        "codigo_uf": "15",
        "codigo_regiao": "1",
        "nome_regiao": "Norte",
    },
    {
        "codigo_ibge": "5100102",
        "nome": "Acorizal",
        "sigla_uf": "MT",
        "codigo_uf": "51",
        "codigo_regiao": "5",
        "nome_regiao": "Centro-Oeste",
    },
    {
        "codigo_ibge": "5102637",
        "nome": "Campo Novo do Parecis",
        "sigla_uf": "MT",
        "codigo_uf": "51",
        "codigo_regiao": "5",
        "nome_regiao": "Centro-Oeste",
    },
    {
        "codigo_ibge": "5103403",
        "nome": "Cuiabá",
        "sigla_uf": "MT",
        "codigo_uf": "51",
        "codigo_regiao": "5",
        "nome_regiao": "Centro-Oeste",
    },
    {
        "codigo_ibge": "5107925",
        "nome": "Sorriso",
        "sigla_uf": "MT",
        "codigo_uf": "51",
        "codigo_regiao": "5",
        "nome_regiao": "Centro-Oeste",
    },
    {
        "codigo_ibge": "5211800",
        "nome": "Jaraguá",
        "sigla_uf": "GO",
        "codigo_uf": "52",
        "codigo_regiao": "5",
        "nome_regiao": "Centro-Oeste",
    },
    {
        "codigo_ibge": "5300108",
        "nome": "Brasília",
        "sigla_uf": "DF",
        "codigo_uf": "53",
        "codigo_regiao": "5",
        "nome_regiao": "Centro-Oeste",
    },
    {
        "codigo_ibge": "3548500",
        "nome": "Santos",
        "sigla_uf": "SP",
        "codigo_uf": "35",
        "codigo_regiao": "3",
        "nome_regiao": "Sudeste",
    },
    {
        "codigo_ibge": "3550308",
        "nome": "São Paulo",
        "sigla_uf": "SP",
        "codigo_uf": "35",
        "codigo_regiao": "3",
        "nome_regiao": "Sudeste",
    },
]

CAPTURADO_EM = "2026-06-25T12:00:00Z"
FONTE_OFICIAL = "https://servicodados.ibge.gov.br/api/docs/localidades"


def write_reference_municipios(lake_root: Path) -> None:
    """Write minimal municipios gold mart for validate_codigo_ibge in CI MVPs."""
    mart_dir = lake_root / "gold" / "mart_ibge__localidades_municipios"
    mart_dir.mkdir(parents=True, exist_ok=True)
    n = len(REFERENCE_MUNICIPIOS)
    table = pa.table(
        {
            "codigo_ibge": [row["codigo_ibge"] for row in REFERENCE_MUNICIPIOS],
            "nome": [row["nome"] for row in REFERENCE_MUNICIPIOS],
            "sigla_uf": [row["sigla_uf"] for row in REFERENCE_MUNICIPIOS],
            "codigo_uf": [row["codigo_uf"] for row in REFERENCE_MUNICIPIOS],
            "codigo_regiao": [row["codigo_regiao"] for row in REFERENCE_MUNICIPIOS],
            "nome_regiao": [row["nome_regiao"] for row in REFERENCE_MUNICIPIOS],
            "capturado_em": [CAPTURADO_EM] * n,
            "fonte_oficial": [FONTE_OFICIAL] * n,
            "_dataset_id": ["ibge.localidades-municipios"] * n,
            "_source_file": ["seed"] * n,
        }
    )
    pq.write_table(table, mart_dir / "mart.parquet")
