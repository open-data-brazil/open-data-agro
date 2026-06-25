# Changelog

All notable changes to **Open Data Agro** are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

**License:** MIT — 100% open source. See [LICENSE](LICENSE) and [docs/OPEN-SOURCE.md](docs/OPEN-SOURCE.md).

---

## [Unreleased]

### Added

- Agent harness (`agent-rules/`, `agent-harness/`, `.cursor/rules/`)
- Project governance docs (MIT license, SECURITY, CONTRIBUTING, CODE_OF_CONDUCT)
- Basic documentation scaffold (`docs/`, `AGENTS.md`, `README.md`)
- **CONAB Produção Agrícola (Phase 10 MVP):** `producao-agricola.yaml` catalog, official column mapping, golden test vectors, GE suites for grãos, `mart_conab__serie_historica_graos`, `make conab-mvp` offline pipeline
- **CONAB Mercado (Phase 11 MVP):** `mercado.yaml` catalog (8 datasets), full pipeline for `conab.oferta-demanda`, `make conab-mercado-mvp`
- **CONAB Abastecimento + ANP combustíveis (Phase 12 MVP):** `abastecimento.yaml` (3 datasets), `anp/combustiveis.yaml` (2 datasets), full pipeline for `conab.estoques-publicos` and ANP LPC weekly prices, `make conab-abastecimento-mvp`
- **CONAB Armazenamento e Logística (Phase 13 MVP):** `armazenamento-logistica.yaml` (3 datasets), legacy `.xls` ingest, full pipeline for `conab.armazenagem`, `make conab-armazenamento-mvp`
- **CONAB Agricultura Familiar (Phase 14 MVP):** `agricultura-familiar.yaml` (2 PAA datasets), full pipeline for entregas and propostas, `make conab-agricultura-familiar-mvp`
- **IBGE Localidades (Phase 15 ingest):** `ibge/localidades.yaml` (5 datasets), JSON API client, bronze Parquet for municipalities/UFs/regions, `make ibge-localidades-mvp`
- **IBGE PAM (Phase 16 ingest):** `ibge/pam.yaml` (3 datasets), SIDRA API client with chunked UF/year/crop pulls, bronze Parquet, GE suites, dbt staging for area-quantidade, `make ibge-pam-mvp`
- **INMET Clima Histórico (Phase 17 ingest):** `inmet/clima.yaml` (5 datasets), station catalog + BDMEP annual ZIP client, daily/monthly long-format bronze, GE suites, `make inmet-clima-mvp`
- **BCB Séries Macro (Phase 18 ingest):** `bcb/sgs.yaml` (5 datasets), SGS API client with 10-year pagination, bronze Parquet, GE suites, dbt staging for IPCA/PTAX, `make bcb-sgs-mvp`
- **CEPEA Preços Agro (Phase 19 ingest):** `cepea/indicadores.yaml` (4 datasets), HTML indicator client with CEPEA/NA mirror fallback, bronze Parquet, GE suites, `--from` ISO date, `make cepea-indicadores-mvp`
