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
- **CONAB Mercado preços semanal UF (Phase 11 P1):** full pipeline for `conab.precos-agropecuarios-semanal-uf`, ISO-8859-1 → UTF-8 for portal TXT, `make conab-mercado-precos-mvp`
- **CONAB Mercado preços semanal município (Phase 11 P1):** full pipeline for `conab.precos-agropecuarios-semanal-municipio`, `cod_ibge` zero-padded in dbt, `make conab-mercado-precos-mvp`
- **CONAB Mercado preços mensal UF (Phase 11 P1):** full pipeline for `conab.precos-agropecuarios-mensal-uf`, monthly grain without `semana`, `make conab-mercado-precos-mvp`
- **CONAB Mercado preços mensal município (Phase 11 P1):** full pipeline for `conab.precos-agropecuarios-mensal-municipio`, `cod_ibge` zero-padded in dbt, `make conab-mercado-precos-mvp`
- **CONAB Mercado preços mínimos (Phase 11 P2):** full pipeline for `conab.precos-minimos`, vigency-period grain, `make conab-mercado-precos-minimos-mvp`
- **CONAB Prohort (Phase 11 P3):** full pipeline for `conab.prohort-diario` and `conab.prohort-mensal`, CEASA horticulture prices/trade, `make conab-mercado-prohort-mvp`
- **CONAB Frete (Phase 13 P1):** full pipeline for `conab.frete`, origin/destination IBGE grain, `make conab-armazenamento-logistica-mvp`
- **CONAB Capacidade Estática (Phase 13 P1):** full pipeline for `conab.serie-historica-capacidade-estatica`, UF × year grain from legacy `.xls`, `make conab-armazenamento-logistica-mvp`
- **CONAB Abastecimento + ANP combustíveis (Phase 12 MVP):** `abastecimento.yaml` (3 datasets), `anp/combustiveis.yaml` (2 datasets), full pipeline for all five datasets including operações and vendas balcão, `make conab-abastecimento-mvp`
- **ANP combustíveis standalone (Phase 12 P2):** `make anp-mvp` for LPC médios/postos only — `scripts/ci/seed_anp_silver.py`, `dbt-build-anp`, DuckDB views `analytics.anp_combustiveis_*`
- **P1 collection sprint (Waves 0–2):** `make p1-collection-mvp` — IBGE localidades (UF/região/meso/micro) + CONAB preços municipais + frete/capacidade with `validate-codigo-ibge`
- **CI P1 collection (sprint Waves 0–2):** `make ci-p1-collection-mvp` in GitHub Actions `dbt` job — mirrors offline collection pipeline with isolated `/tmp` lake
- **Macro collection (Phases 17–19):** `make collection-macro-mvp` — INMET climate + BCB SGS + CEPEA indicators in one offline lake with DuckDB analytics smoke
- **CI macro collection (Phases 17–19):** `make ci-collection-macro-mvp` in GitHub Actions `dbt` job — mirrors offline INMET/BCB/CEPEA pipeline with isolated `/tmp` lake
- **Full collection sprint exit:** `make collection-full-mvp` — runs `p1-collection-mvp`, `collection-macro-mvp`, `ibge-pam-mvp`, and `anp-mvp` end-to-end
- **CONAB Armazenamento e Logística (Phase 13 MVP):** `armazenamento-logistica.yaml` (3 datasets), legacy `.xls` ingest, full pipeline for `conab.armazenagem`, `make conab-armazenamento-mvp`
- **CONAB Agricultura Familiar (Phase 14 MVP):** `agricultura-familiar.yaml` (2 PAA datasets), full pipeline for entregas and propostas, `make conab-agricultura-familiar-mvp`
- **IBGE Localidades (Phase 15):** full E2E for municipios + UFs + regiões + meso/micro — dbt marts, DuckDB views `analytics.ibge_localidades_*`, `make ibge-localidades-mvp`
- **IBGE Localidades live smoke (Phase 15 P2):** `make ibge-localidades-live-smoke` — live ingestor for all five localidades datasets + `scripts/ci/check_ibge_localidades_bronze.py` row-count gate
- **validate-codigo-ibge-lake:** `make validate-codigo-ibge-lake` — cross-check CONAB/PAM `cod_ibge` against full `./lake` municipios mart (~5.5k rows)
- **Ingestor stress benchmark:** `make benchmark-ingestor-fast10-stress` — fast10 plus large CONAB tables (`operacoes-comercializacao`, `prohort-diario`) via `scripts/benchmark/profiles/fast10-stress.json`
- **IBGE cod_ibge validation (Phase 15 P4):** `scripts/quality/validate_codigo_ibge.py` cross-checks CONAB gold marts against `mart_ibge__localidades_municipios`, `make validate-codigo-ibge`
- **CONAB Mercado cod_ibge validation (Phase 11 P4):** `validate-codigo-ibge` wired into `conab-mercado-full-mvp`, `conab-mercado-precos-mvp`, and `conab-mercado-prohort-mvp`; shared `scripts/ci/reference_municipios.py` for CI seeds
- **CONAB Armazenamento cod_ibge validation (Phase 13 P4):** `validate-codigo-ibge` wired into `conab-armazenamento-mvp` and `conab-armazenamento-logistica-mvp`; frete origin/destination and armazenagem `cod_ibge` checked against IBGE localidades
- **CONAB Abastecimento + PAA cod_ibge validation (Phase 12/14 P4):** `validate-codigo-ibge` wired into `conab-abastecimento-mvp` and `conab-agricultura-familiar-mvp`; estoques públicos and Alimenta Brasil propostas checked against IBGE localidades
- **CONAB Produção cod_ibge validation (Phase 10 P4):** `validate-codigo-ibge` wired into `conab-mvp`; custo de produção `cod_ibge` checked against IBGE localidades
- **CI cod_ibge validation (Phase 15 P5):** `make ci-validate-codigo-ibge` seeds all CONAB marts with `cod_ibge`, runs dbt gold build, and cross-checks in GitHub Actions `dbt` job
- **CI PAM codigo_ibge validation (Phase 16 P4):** `ci-validate-codigo-ibge` extended with IBGE PAM gold marts (`codigo_ibge` vs localidades reference)
- **CONAB Mercado CI consolidation (Phase 11 §7):** `dbt-build-mercado` covers all 8 mercado marts, `conab-mercado-full-mvp`, committed `scripts/benchmark/profiles/fast10.json` with `precos-semanal-uf` + `frete`
- **IBGE PAM (Phase 16):** full E2E for area-quantidade, rendimento-valor, estabelecimentos — dbt marts, DuckDB views `analytics.ibge_pam_*`, `make ibge-pam-mvp`
- **CI IBGE PAM (Phase 16):** `make ci-ibge-pam-mvp` in GitHub Actions `dbt` job — mirrors offline PAM pipeline with `validate-codigo-ibge` on isolated `/tmp` lake
- **INMET Clima Histórico (Phase 17):** full E2E for station catalogs, BDMEP diário/mensal, pacote-anual-automaticas — dbt marts, DuckDB views, `make inmet-clima-mvp`
- **BCB Séries Macro (Phase 18):** full E2E for IPCA, IPCA 12m, IGP-M, PTAX compra/venda — dbt marts, DuckDB views, `make bcb-sgs-mvp`
- **CEPEA Preços Agro (Phase 19):** full E2E for soja Paranaguá/PR, milho, boi gordo — dbt marts, DuckDB views, `make cepea-indicadores-mvp`
