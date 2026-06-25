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
