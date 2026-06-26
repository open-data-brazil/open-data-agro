# Changelog

All notable changes to **Open Data Agro** are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

**License:** MIT — 100% open source. See [LICENSE](LICENSE) and [docs/OPEN-SOURCE.md](docs/OPEN-SOURCE.md).

---

## [Unreleased]

### Added

- **International sources wave 3 (Phase 45):** `oecd-fao.ag-outlook`, `fao.food-price-index`, `argentina.magyp-producion-granos` — full E2E pipelines, GE suites, dbt marts, DuckDB views, `make international-sources-wave-3-mvp` + `make ci-international-sources-wave-3-mvp`; IMF PCPS, Paraguay BCP, Uruguay INE, NOAA GPCC, China NBS, USDA AMS, Baltic BDI deferred (verified blockers)
- **BR sources wave 3 (Phase 44):** `dnit.snv-rodovias-federais`, `ipea.series-macro-regionais`, `ibge.pevs-producao-vegetal`, `ibge.ppm-producao-municipal`, `aneel.tarifas-energia`, `bndes.financiamento-agro`, `inmet.sequia-monitor` — full E2E pipelines, GE suites, dbt marts, DuckDB views, `make br-sources-wave-3-mvp` + `make ci-br-sources-wave-3-mvp`; `inpe.cptec-indices-climaticos`, `mapa.sigef-areas`, `embrapa.solos-brasil` remain deferred (no stable bulk URL)
- **Source discovery wave 3 (Phase 43):** `.local/SOURCE-DISCOVERY-CATALOG.md` wave 3 section (18 new candidates, 86 total); status sync for 28 implemented fichas; `DISCOVERY-REPORT-WAVE3.md` + `GAP-MATRIX-WAVE3.md`; routes to Phases 44–46 (ingestor only — no IA)
- **Ingestor signoff (Phase 42):** `make ingestor-signoff-mvp` + `make ci-ingestor-signoff-mvp`; `scripts/ci/spot_check_analytics.py` for PostgreSQL `analytics.*` row/date spot-checks
- **International sources wave 2 (Phase 41):** `igc.goi-index`, `usda.gats-trade`, `eurostat.ag-prices`, `argentina.bcra-cambio` — full E2E pipelines, GE suites, dbt marts, DuckDB views, `make international-sources-wave-2-mvp` + `make ci-international-sources-wave-2-mvp`
- **BR sources wave 2 (Phase 40):** `mapa.agrofit-produtos-formulados`, `mapa.agrofit-produtos-tecnicos`, `ana.hidrologia-series`, `antaq.movimentacao-carga-portuaria` — full E2E pipelines, GE suites, dbt marts, DuckDB views, `make br-sources-wave-2-mvp` + `make ci-br-sources-wave-2-mvp`; CONAB re-audit confirms no separate private-stocks bulk file on portal
- **Source discovery wave 2 (Phase 39):** `.local/SOURCE-DISCOVERY-CATALOG.md` wave 2 section (≥20 new candidates); `DISCOVERY-REPORT-WAVE2.md`; phases 40–42 planned

### Added (prior wave)

- **International new sources (Phase 38):** `eia.petroleum-prices` (P0), `usda.wasde`, `igc.goi-index`, `un.comtrade-bulk` (P1) — full ingest + GE + dbt + DuckDB; `make international-new-sources-mvp` + `make ci-international-new-sources-mvp`
- **BR new sources (Phase 37):** `ibge.lspa-area-producao`, `bcb.sgs-selic` — IBGE LSPA SIDRA 6588 monthly UF production + BCB Selic SGS 11, GE suites, dbt marts, DuckDB views, `make br-new-sources-mvp` + `make ci-br-new-sources-mvp`
- **International extended (Phase 36):** `fao.producao-agro`, `fao.comercio-agro`, `worldbank.ag-indices` — FAOSTAT annual bulk + Pink Sheet agriculture indices, GE suites, dbt marts, DuckDB views, `make international-extended-mvp` + `make ci-international-extended-mvp`
- **BR comex extended (Phase 35):** `mdic.comex-importacao-ncm-mes`, `mdic.comex-exportacao-uf-ncm`, `mdic.comex-importacao-diesel-ncm` — Comex Stat import/UF export API client, GE suites, dbt marts, DuckDB views, `make mdic-comex-extended-mvp` + `make ci-mdic-comex-extended-mvp`
- **BR logistics extended (Phase 34):** `antt.volume-trafego-pedagio` (P0) and `antt.receita-por-praca` (P1) — CKAN name/year resolver, GE suites, dbt marts, DuckDB views `analytics.antt_volume_trafego_pedagio` / `analytics.antt_receita_por_praca`, `make br-logistica-extended-mvp` + `make ci-br-logistica-extended-mvp`; `antt.tarifas-pedagio` deferred (CKAN package 404)
- **Collection hardening (Phase 33):** [docs/DATA-HISTORY-RANGES.md](docs/DATA-HISTORY-RANGES.md) — per-dataset source min years and `--from` backfill reference; CONAB portal audit; `make collection-hardening-mvp` + `make ci-collection-hardening-mvp`; `scripts/ci/check_data_history_ranges.py` wired into `ci-dbt`
- **BCB PTAX history:** catalog `start_date` / `period_start` aligned to **1984** per official SGS series limits

### Added (prior)

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
- **CI ANP combustíveis (Phase 12 P2):** `make ci-anp-mvp` in GitHub Actions `dbt` job — mirrors offline ANP pipeline with isolated `/tmp` lake
- **P1 collection sprint (Waves 0–2):** `make p1-collection-mvp` — IBGE localidades (UF/região/meso/micro) + CONAB preços municipais + frete/capacidade with `validate-codigo-ibge`
- **CI P1 collection (sprint Waves 0–2):** `make ci-p1-collection-mvp` in GitHub Actions `dbt` job — mirrors offline collection pipeline with isolated `/tmp` lake
- **Macro collection (Phases 17–19):** `make collection-macro-mvp` — INMET climate + BCB SGS + CEPEA indicators in one offline lake with DuckDB analytics smoke
- **CI macro collection (Phases 17–19):** `make ci-collection-macro-mvp` in GitHub Actions `dbt` job — mirrors offline INMET/BCB/CEPEA pipeline with isolated `/tmp` lake
- **Full collection sprint exit:** `make collection-full-mvp` — runs `p1-collection-mvp`, `collection-macro-mvp`, `ibge-pam-mvp`, and `anp-mvp` end-to-end
- **MDIC Comex exportação agro (Phase 21):** `mdic.comex-exportacao-ncm-mes` — Comex Stat API client, GE suite, dbt mart, DuckDB view `analytics.mdic_comex_exportacao_ncm_mes`, `make mdic-comex-mvp` + `make ci-mdic-comex-mvp`
- **ANTT praças de pedágio (Phase 22):** `antt.pracas-pedagio` — CKAN CSV resolver, GE suite, dbt mart, DuckDB view `analytics.antt_pracas_pedagio`, `make dnit-antt-logistica-mvp` + `make ci-dnit-antt-logistica-mvp`
- **MAPA ZARC tábua de risco (Phase 23):** `mapa.zarc-tabua-risco` — CKAN latest-safra CSV resolver, GE suite, dbt mart, DuckDB view `analytics.mapa_zarc_tabua_risco`, `make mapa-dados-mvp` + `make ci-mapa-dados-mvp`
- **B3 futuros agro (Phase 24):** `b3.futuro-soja`, `b3.futuro-milho`, `b3.futuro-boi` — SPRD BVBG.187 parser, GE suites, dbt marts, DuckDB views `analytics.b3_futuro_*`, `make b3-futuros-mvp` + `make ci-b3-futuros-mvp`
- **USDA FAS PSD (Phase 25):** `usda.psd-soja`, `usda.psd-milho`, `usda.psd-trigo` — AMIS SOAP parser, GE suites, dbt marts, DuckDB views `analytics.usda_psd_*`, `make usda-psd-mvp` + `make ci-usda-psd-mvp`
- **FAO FAOSTAT prices (Phase 26):** `fao.prices-agro` — bulk normalized CSV parser, GE suite, dbt mart, DuckDB view `analytics.fao_prices_agro`, `make fao-faostat-mvp` + `make ci-fao-faostat-mvp`
- **World Bank Pink Sheet (Phase 27):** `worldbank.pink-sheet-monthly` — XLSX unpivot parser, GE suite, dbt mart, DuckDB view `analytics.worldbank_pink_sheet_monthly`, `make worldbank-commodities-mvp` + `make ci-worldbank-commodities-mvp`
- **NOAA climate indices (Phase 28):** `noaa.enso-indices`, `noaa.global-temp-anomaly` — ONI ASCII + NCEI CSV parsers, GE suites, dbt marts, DuckDB views `analytics.noaa_enso_indices` / `analytics.noaa_global_temp_anomaly`, `make noaa-climate-mvp` + `make ci-noaa-climate-mvp`
- **Unified PostgreSQL (Phase 29):** `processor sync-postgres`, migration `000005_analytics_schema`, manifest tables, join-key indexes, `make unified-db-sync` + `make ci-unified-db-sync` — see [docs/POSTGRES-UNIFIED-SYNC.md](docs/POSTGRES-UNIFIED-SYNC.md)
- **Roadmap sync:** `docs/ROADMAP.md` updated to reflect Go local-first stack, phases 0–19 status, and CI/collection sprint exit targets
- **Phase TASKS bulk sync:** `.local/phases/*/TASKS.md` checkboxes aligned with repo reality (phases 10–13, 15–17, sprint progress table)
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
- **Roadmap sync:** `docs/ROADMAP.md` updated to reflect Go local-first stack, phases 0–19 status, and CI/collection sprint exit targets
- **CI dbt mirror:** `make ci-dbt` extended with `ci-validate-codigo-ibge` and `ci-collection-full-mvp` to match GitHub Actions `dbt` job locally
- **Phase TASKS bulk sync:** `.local/phases/*/TASKS.md` checkboxes aligned with repo reality (phases 10–13, 15–17, sprint progress table)
- **OFFICIAL-SOURCES status sync:** normalized all dataset rows to `**Pn — implemented**`; added `scripts/ci/check_official_sources_status.py` gate
- **Benchmark docs sync:** `.local/benchmark/README.md` paths → `scripts/benchmark/profiles/`; fast10=16 / fast10-stress=18; `make benchmark-ingestor-fast10-stress` documented; `scripts/ci/check_benchmark_profiles.py` gate
- **IMPLEMENTATION-PLAN sync:** Phases 15–19 documented as full E2E with `make *-mvp` gates; grãos DoD marked complete; `scripts/ci/check_implementation_plan.py` gate
- **README sprint exit:** status + quick start updated for Phases 0–19 collection complete; `scripts/ci/check_readme_status.py` gate
- **Phase 6 quality docs:** GE vs `validate_codigo_ibge.py` split in phase README + `scripts/quality/README.md`; `scripts/ci/check_phase6_quality_docs.py` gate
- **NEW-PROJECT-CHECKLIST sync:** updated for Go local-first post-collection state; `scripts/ci/check_new_project_checklist.py` gate
- **Prohort OFFICIAL-REFERENCE:** live portal column mapping for `conab.prohort-diario` + `conab.prohort-mensal`; `scripts/ci/check_prohort_official_reference.py` gate
- **Phase 20 scaffold:** `.local/phases/20-analytics-crossing/` for post-collection analytics crossing (DATA-CROSSING-VISION); `scripts/ci/check_phase20_scaffold.py` gate
- **MinIO CI integration (Phase 2):** `make ci-minio` — Docker MinIO, bronze S3 Put/List, DuckDB `s3://` smoke; GitHub Actions `go` job; `scripts/ci/check_minio_ci.py` gate
- **R2 production runbook (Phase 2):** `infra/r2/README.md` deploy runbook, `make validate-r2-env` / `make ci-validate-r2-env`, optional `R2_INTEGRATION` live smoke; `scripts/ci/check_r2_runbook.py` gate
- **Delta Lake silver versioning (Phase 3):** `DELTA_MIN_VERSIONS` wired in `scripts/delta/promote.py`, append + DuckDB time-travel tests, `make ci-delta-versioning`; `scripts/ci/check_delta_versioning.py` gate
- **New dataset checklist:** `docs/NEW-DATASET-CHECKLIST.md` for adding catalog datasets; `make ci-new-dataset-checklist`; `scripts/ci/check_new_dataset_checklist.py` gate

### Changed

- **docker-compose MinIO images:** pin to `minio/minio:latest` and `minio/mc:latest` (official release tags removed from Docker Hub)
- **DuckDB S3 smoke:** `FORCE INSTALL httpfs`, robust `parseCountCSV` for extension setup output; MinIO integration test seeds bronze Parquet before read

- **Collection sprint exit (2026-06-26):** Phases 0–19 E2E, `make ci-collection-full-mvp` in GitHub Actions, public docs synced (ROADMAP, OFFICIAL-SOURCES, NEW-PROJECT-CHECKLIST, README)
