.PHONY: test lint build build-processor clean duckdb-install python-install dbt-deps dbt-build dbt-build-mercado dbt-build-mercado-precos dbt-build-mercado-prohort dbt-build-abastecimento dbt-build-anp dbt-build-armazenamento dbt-build-armazenamento-logistica dbt-build-agricultura-familiar dbt-build-ibge-localidades dbt-build-ibge-pam dbt-build-bcb-sgs dbt-build-cepea dbt-build-inmet-clima dbt-build-mdic-comex dbt-build-antt-pedagios dbt-build-mapa-zarc dbt-build-b3-futuros ibge-localidades-mvp ibge-localidades-live-smoke ibge-pam-mvp ci-ibge-pam-mvp inmet-clima-mvp bcb-sgs-mvp cepea-indicadores-mvp mdic-comex-mvp ci-mdic-comex-mvp dnit-antt-logistica-mvp ci-dnit-antt-logistica-mvp mapa-dados-mvp ci-mapa-dados-mvp b3-futuros-mvp ci-b3-futuros-mvp dbt-build-usda-psd usda-psd-mvp ci-usda-psd-mvp dbt-build-fao-faostat fao-faostat-mvp ci-fao-faostat-mvp dbt-build-worldbank-commodities worldbank-commodities-mvp ci-worldbank-commodities-mvp dbt-build-international-extended international-extended-mvp ci-international-extended-mvp dbt-build-br-new-sources br-new-sources-mvp ci-br-new-sources-mvp dbt-build-international-new-sources international-new-sources-mvp ci-international-new-sources-mvp dbt-build-br-sources-wave-2 br-sources-wave-2-mvp ci-br-sources-wave-2-mvp dbt-build-international-sources-wave-2 international-sources-wave-2-mvp ci-international-sources-wave-2-mvp ingestor-signoff-mvp ci-ingestor-signoff-mvp dbt-build-noaa-climate noaa-climate-mvp ci-noaa-climate-mvp unified-db-sync ci-unified-db-sync seed-unified-db-gold anp-mvp ci-anp-mvp p1-collection-mvp ci-p1-collection-mvp collection-macro-mvp ci-collection-macro-mvp collection-full-mvp ci-collection-full-mvp ci-go ci-minio ci-validate-r2-env validate-r2-env validate-r2-env-live ci-delta-versioning ci-new-dataset-checklist ci-dbt ci-validate-codigo-ibge validate-codigo-ibge validate-codigo-ibge-lake benchmark-ingestor benchmark-ingestor-clean benchmark-ingestor-fast10 benchmark-ingestor-clean benchmark-ingestor-fast10-stress benchmark-ingestor-fast10-stress-clean migrate-install migrate-up migrate-down seed analytics-init analytics-smoke conab-reference conab-mvp conab-mercado-mvp conab-mercado-full-mvp conab-mercado-precos-mvp conab-mercado-precos-minimos-mvp conab-mercado-prohort-mvp conab-abastecimento-mvp conab-armazenamento-mvp conab-armazenamento-logistica-mvp conab-agricultura-familiar-mvp

BIN_DIR := bin
DUCKDB_VERSION ?= 1.5.4
LAKE_LOCAL_ROOT ?= ./lake
LAKE_ABS := $(abspath $(LAKE_LOCAL_ROOT))
DUCKDB_PATH ?= ./duckdb/open_data_agro.duckdb
DUCKDB_BIN ?= $(CURDIR)/.local/bin/duckdb
POSTGRES_HOST_PORT ?= 5432
DATABASE_URL ?= postgresql://open_data_agro:open_data_agro@localhost:$(POSTGRES_HOST_PORT)/open_data_agro?sslmode=disable
MIGRATIONS_PATH := infra/postgres/migrations
MIGRATE ?= migrate
BENCHMARK_PROFILE ?= scripts/benchmark/profiles/fast10.json
BENCHMARK_STRESS_PROFILE ?= scripts/benchmark/profiles/fast10-stress.json
MERCADO_DBT_SELECT := stg_conab__oferta_demanda+ stg_conab__precos_semanal_uf+ stg_conab__precos_semanal_municipio+ stg_conab__precos_mensal_uf+ stg_conab__precos_mensal_municipio+ stg_conab__precos_minimos+ stg_conab__prohort_diario+ stg_conab__prohort_mensal+
CI_COD_IBGE_LAKE ?= /tmp/cod-ibge-ci-lake
COD_IBGE_DBT_SELECT := stg_conab__custo_producao+ stg_conab__precos_semanal_municipio+ stg_conab__precos_mensal_municipio+ stg_conab__frete+ stg_conab__armazenagem+ stg_conab__estoques_publicos+ stg_conab__alimenta_brasil_propostas+ stg_conab__prohort_diario+ stg_conab__prohort_mensal+ stg_ibge__pam_area_quantidade+ stg_ibge__pam_rendimento_valor+ stg_ibge__pam_estabelecimentos+
COLLECTION_P1_LAKE ?= /tmp/p1-collection-lake
COLLECTION_P1_DUCKDB ?= /tmp/p1-collection.duckdb
CI_P1_COLLECTION_LAKE ?= /tmp/p1-collection-ci-lake
CI_P1_COLLECTION_DUCKDB ?= /tmp/p1-collection-ci.duckdb
COLLECTION_P1_DBT_SELECT := stg_ibge__localidades_ufs+ stg_ibge__localidades_regioes+ stg_ibge__localidades_mesorregioes+ stg_ibge__localidades_microrregioes+ stg_conab__precos_semanal_municipio+ stg_conab__precos_mensal_municipio+ stg_conab__frete+ stg_conab__capacidade_estatica+
COLLECTION_MACRO_LAKE ?= /tmp/macro-collection-lake
COLLECTION_MACRO_DUCKDB ?= /tmp/macro-collection.duckdb
CI_COLLECTION_MACRO_LAKE ?= /tmp/macro-collection-ci-lake
CI_COLLECTION_MACRO_DUCKDB ?= /tmp/macro-collection-ci.duckdb
CI_IBGE_PAM_LAKE ?= /tmp/ibge-pam-ci-lake
CI_IBGE_PAM_DUCKDB ?= /tmp/ibge-pam-ci.duckdb
CI_ANP_LAKE ?= /tmp/anp-ci-lake
CI_ANP_DUCKDB ?= /tmp/anp-ci.duckdb
CI_MDIC_LAKE ?= /tmp/mdic-ci-lake
CI_MDIC_DUCKDB ?= /tmp/mdic-ci.duckdb
CI_MDIC_COMEX_EXTENDED_LAKE ?= /tmp/mdic-comex-extended-ci-lake
CI_MDIC_COMEX_EXTENDED_DUCKDB ?= /tmp/mdic-comex-extended-ci.duckdb
CI_ANTT_PEDAGIOS_LAKE ?= /tmp/antt-pedagios-ci-lake
CI_ANTT_PEDAGIOS_DUCKDB ?= /tmp/antt-pedagios-ci.duckdb
CI_ANTT_LOGISTICA_EXTENDED_LAKE ?= /tmp/antt-logistica-extended-ci-lake
CI_ANTT_LOGISTICA_EXTENDED_DUCKDB ?= /tmp/antt-logistica-extended-ci.duckdb
CI_MAPA_ZARC_LAKE ?= /tmp/mapa-zarc-ci-lake
CI_MAPA_ZARC_DUCKDB ?= /tmp/mapa-zarc-ci.duckdb
CI_B3_FUTUROS_LAKE ?= /tmp/b3-futuros-ci-lake
CI_B3_FUTUROS_DUCKDB ?= /tmp/b3-futuros-ci.duckdb
CI_USDA_PSD_LAKE ?= /tmp/usda-psd-ci-lake
CI_USDA_PSD_DUCKDB ?= /tmp/usda-psd-ci.duckdb
CI_FAO_FAOSTAT_LAKE ?= /tmp/fao-faostat-ci-lake
CI_FAO_FAOSTAT_DUCKDB ?= /tmp/fao-faostat-ci.duckdb
CI_WORLDBANK_COMMODITIES_LAKE ?= /tmp/worldbank-commodities-ci-lake
CI_WORLDBANK_COMMODITIES_DUCKDB ?= /tmp/worldbank-commodities-ci.duckdb
CI_NOAA_CLIMATE_LAKE ?= /tmp/noaa-climate-ci-lake
CI_NOAA_CLIMATE_DUCKDB ?= /tmp/noaa-climate-ci.duckdb
CI_INTERNATIONAL_EXTENDED_LAKE ?= /tmp/international-extended-ci-lake
CI_INTERNATIONAL_EXTENDED_DUCKDB ?= /tmp/international-extended-ci.duckdb
CI_BR_NEW_SOURCES_LAKE ?= /tmp/br-new-sources-ci-lake
CI_BR_NEW_SOURCES_DUCKDB ?= /tmp/br-new-sources-ci.duckdb
CI_INTERNATIONAL_NEW_SOURCES_LAKE ?= /tmp/international-new-sources-ci-lake
CI_INTERNATIONAL_NEW_SOURCES_DUCKDB ?= /tmp/international-new-sources-ci.duckdb
CI_UNIFIED_DB_LAKE ?= /tmp/unified-db-ci-lake
UNIFIED_DB_SYNC_MARTS ?=
COLLECTION_MACRO_DBT_SELECT := stg_inmet__estacoes_automaticas+ stg_inmet__estacoes_convencionais+ stg_inmet__bdmep_diario+ stg_inmet__bdmep_mensal+ stg_inmet__pacote_anual_automaticas+ stg_bcb__sgs_ipca+ stg_bcb__sgs_ipca_12m+ stg_bcb__sgs_igpm+ stg_bcb__sgs_ptax_usd_venda+ stg_bcb__sgs_ptax_usd_compra+ stg_cepea__soja_paranagua+ stg_cepea__soja_parana+ stg_cepea__milho+ stg_cepea__boi_gordo+

test:
	go test ./...

# Mirror GitHub Actions CI jobs locally (see .github/workflows/ci.yml).
ci-go: duckdb-install python-install
	go work sync
	go test ./...
	PATH="$(PWD)/.local/bin:$$PATH" DUCKDB_BIN="$(PWD)/.local/bin/duckdb" DUCKDB_INTEGRATION=1 go test ./internal/processor -run 'SmokeLocal|PreviewPromote' -count=1
	GE_INTEGRATION=1 go test ./internal/processor -run 'Quality' -count=1

# Docker MinIO + STORAGE_MODE=minio integration (see .github/workflows/ci.yml go job).
ci-minio: duckdb-install
	bash scripts/ci/run_minio_integration.sh

# Offline R2 runbook + env validation gate (no Cloudflare credentials).
ci-validate-r2-env:
	python3 scripts/ci/check_r2_runbook.py
	VALIDATE_R2_FIXTURE=1 bash scripts/deploy/validate_r2_env.sh
	go test ./internal/config -run 'R2|ValidateR2' -count=1

# Production: validate STORAGE_MODE=r2 credentials in .env (optional live Put/List with R2_INTEGRATION=1).
validate-r2-env:
	bash scripts/deploy/validate_r2_env.sh

validate-r2-env-live: validate-r2-env
	R2_INTEGRATION=1 bash scripts/deploy/validate_r2_env.sh

# Native Delta Lake silver versioning (promote append + DuckDB time travel).
ci-delta-versioning: duckdb-install python-install
	bash scripts/ci/run_delta_versioning.sh

ci-new-dataset-checklist:
	python3 scripts/ci/check_new_dataset_checklist.py

ci-data-history-ranges:
	python3 scripts/ci/check_data_history_ranges.py

ci-dbt: duckdb-install python-install
	rm -rf /tmp/open-data-agro-lake /tmp/open-data-agro-ci.duckdb /tmp/open-data-agro-analytics.duckdb
	LAKE_LOCAL_ROOT=/tmp/open-data-agro-lake python3 scripts/ci/seed_dbt_silver.py
	cp -f dbt/profiles.yml.example dbt/profiles.yml
	cd dbt && PATH="$(CURDIR)/.local/bin:$$PATH" DUCKDB_BIN="$(CURDIR)/.local/bin/duckdb" \
		LAKE_LOCAL_ROOT=/tmp/open-data-agro-lake DUCKDB_PATH=/tmp/open-data-agro-ci.duckdb \
		dbt deps --profiles-dir . && \
		LAKE_LOCAL_ROOT=/tmp/open-data-agro-lake DUCKDB_PATH=/tmp/open-data-agro-ci.duckdb \
		dbt build --profiles-dir . --select 'stg_conab__serie_historica_graos stg_conab__estimativa_graos+'
	cd dbt && PATH="$(CURDIR)/.local/bin:$$PATH" DUCKDB_BIN="$(CURDIR)/.local/bin/duckdb" \
		LAKE_LOCAL_ROOT=/tmp/open-data-agro-lake DUCKDB_PATH=/tmp/open-data-agro-ci.duckdb \
		dbt docs generate --profiles-dir .
	PATH="$(PWD)/.local/bin:$$PATH" DUCKDB_BIN="$(PWD)/.local/bin/duckdb" \
		LAKE_LOCAL_ROOT=/tmp/open-data-agro-lake python3 scripts/ci/seed_mercado_silver.py
	cd dbt && PATH="$(CURDIR)/.local/bin:$$PATH" DUCKDB_BIN="$(CURDIR)/.local/bin/duckdb" \
		LAKE_LOCAL_ROOT=/tmp/open-data-agro-lake DUCKDB_PATH=/tmp/open-data-agro-analytics.duckdb \
		dbt build --profiles-dir . --select 'stg_conab__oferta_demanda+'
	cd "$(CURDIR)" && PATH="$(CURDIR)/.local/bin:$$PATH" DUCKDB_BIN="$(CURDIR)/.local/bin/duckdb" \
		LAKE_LOCAL_ROOT=/tmp/open-data-agro-lake DUCKDB_PATH=/tmp/open-data-agro-analytics.duckdb \
		$(MAKE) analytics-init analytics-smoke
	$(DUCKDB_BIN) /tmp/open-data-agro-analytics.duckdb -c "SELECT COUNT(*) FROM analytics.conab_oferta_demanda"
	$(DUCKDB_BIN) /tmp/open-data-agro-analytics.duckdb -c "SELECT * FROM analytics.conab_estimativa_graos LIMIT 10"
	$(MAKE) ci-validate-codigo-ibge CI_COD_IBGE_LAKE=$(CI_COD_IBGE_LAKE)
	$(MAKE) ci-collection-full-mvp
	$(MAKE) ci-data-history-ranges

# Offline CI: seed CONAB + IBGE PAM marts with cod_ibge/codigo_ibge, build gold, cross-check.
ci-validate-codigo-ibge: duckdb-install python-install dbt-deps
	rm -rf $(CI_COD_IBGE_LAKE) /tmp/cod-ibge-ci.duckdb
	cp -f dbt/profiles.yml.example dbt/profiles.yml
	LAKE_LOCAL_ROOT=$(CI_COD_IBGE_LAKE) python3 scripts/ci/seed_dbt_silver.py
	LAKE_LOCAL_ROOT=$(CI_COD_IBGE_LAKE) python3 scripts/ci/seed_mercado_silver.py
	LAKE_LOCAL_ROOT=$(CI_COD_IBGE_LAKE) python3 scripts/ci/seed_armazenamento_silver.py
	LAKE_LOCAL_ROOT=$(CI_COD_IBGE_LAKE) python3 scripts/ci/seed_abastecimento_silver.py
	LAKE_LOCAL_ROOT=$(CI_COD_IBGE_LAKE) python3 scripts/ci/seed_agricultura_familiar_silver.py
	LAKE_LOCAL_ROOT=$(CI_COD_IBGE_LAKE) python3 scripts/ci/seed_ibge_pam_silver.py
	cd dbt && PATH="$(CURDIR)/.local/bin:$$PATH" DUCKDB_BIN="$(CURDIR)/.local/bin/duckdb" \
		LAKE_LOCAL_ROOT=$(CI_COD_IBGE_LAKE) DUCKDB_PATH=/tmp/cod-ibge-ci.duckdb \
		dbt build --profiles-dir . --select '$(COD_IBGE_DBT_SELECT)'
	python3 scripts/quality/validate_codigo_ibge.py --lake-root $(CI_COD_IBGE_LAKE)

# P1 collection sprint: Waves 0–2 (localidades + municipal prices + logistics) in one offline lake.
p1-collection-mvp: duckdb-install python-install dbt-deps
	go test ./internal/ibge/... ./internal/ingest/ -run 'IBGE|Localidades|Precos|Frete|Capacidade'
	rm -rf $(COLLECTION_P1_LAKE) $(COLLECTION_P1_DUCKDB)
	cp -f dbt/profiles.yml.example dbt/profiles.yml
	LAKE_LOCAL_ROOT=$(COLLECTION_P1_LAKE) python3 scripts/ci/seed_ibge_localidades_silver.py
	LAKE_LOCAL_ROOT=$(COLLECTION_P1_LAKE) python3 scripts/ci/seed_mercado_silver.py
	LAKE_LOCAL_ROOT=$(COLLECTION_P1_LAKE) python3 scripts/ci/seed_armazenamento_silver.py
	cd dbt && PATH="$(CURDIR)/.local/bin:$$PATH" DUCKDB_BIN="$(CURDIR)/.local/bin/duckdb" \
		LAKE_LOCAL_ROOT=$(COLLECTION_P1_LAKE) DUCKDB_PATH=$(COLLECTION_P1_DUCKDB) \
		dbt build --profiles-dir . --select '$(COLLECTION_P1_DBT_SELECT)'
	$(MAKE) analytics-init LAKE_LOCAL_ROOT=$(COLLECTION_P1_LAKE) DUCKDB_PATH=$(COLLECTION_P1_DUCKDB)
	$(DUCKDB_BIN) $(COLLECTION_P1_DUCKDB) -c "SELECT COUNT(*) AS mesorregioes FROM analytics.ibge_localidades_mesorregioes"
	$(DUCKDB_BIN) $(COLLECTION_P1_DUCKDB) -c "SELECT COUNT(*) AS microrregioes FROM analytics.ibge_localidades_microrregioes"
	$(DUCKDB_BIN) $(COLLECTION_P1_DUCKDB) -c "SELECT COUNT(*) AS precos_mun FROM analytics.conab_precos_semanal_municipio"
	$(DUCKDB_BIN) $(COLLECTION_P1_DUCKDB) -c "SELECT COUNT(*) AS frete_rows FROM analytics.conab_frete"
	$(DUCKDB_BIN) $(COLLECTION_P1_DUCKDB) -c "SELECT uf, ano, quantidade_mil_t FROM analytics.conab_capacidade_estatica WHERE uf = 'MT' LIMIT 3"
	$(MAKE) validate-codigo-ibge LAKE_LOCAL_ROOT=$(COLLECTION_P1_LAKE)

# Mirror GitHub Actions: offline P1 collection sprint (Waves 0–2).
ci-p1-collection-mvp:
	$(MAKE) p1-collection-mvp \
		COLLECTION_P1_LAKE=$(CI_P1_COLLECTION_LAKE) \
		COLLECTION_P1_DUCKDB=$(CI_P1_COLLECTION_DUCKDB)

# Phases 17–19: INMET climate + BCB macro + CEPEA port prices in one offline lake.
collection-macro-mvp: duckdb-install python-install dbt-deps
	go test ./internal/inmet/... ./internal/bcb/... ./internal/cepea/... ./internal/ingest/ -run 'INMET|BCB|SGS|CEPA|Cepea|Flatten'
	rm -rf $(COLLECTION_MACRO_LAKE) $(COLLECTION_MACRO_DUCKDB)
	cp -f dbt/profiles.yml.example dbt/profiles.yml
	LAKE_LOCAL_ROOT=$(COLLECTION_MACRO_LAKE) python3 scripts/ci/seed_inmet_silver.py
	LAKE_LOCAL_ROOT=$(COLLECTION_MACRO_LAKE) python3 scripts/ci/seed_bcb_sgs_silver.py
	LAKE_LOCAL_ROOT=$(COLLECTION_MACRO_LAKE) python3 scripts/ci/seed_cepea_silver.py
	cd dbt && PATH="$(CURDIR)/.local/bin:$$PATH" DUCKDB_BIN="$(CURDIR)/.local/bin/duckdb" \
		LAKE_LOCAL_ROOT=$(COLLECTION_MACRO_LAKE) DUCKDB_PATH=$(COLLECTION_MACRO_DUCKDB) \
		dbt build --profiles-dir . --select '$(COLLECTION_MACRO_DBT_SELECT)'
	$(MAKE) analytics-init LAKE_LOCAL_ROOT=$(COLLECTION_MACRO_LAKE) DUCKDB_PATH=$(COLLECTION_MACRO_DUCKDB)
	$(DUCKDB_BIN) $(COLLECTION_MACRO_DUCKDB) -c "SELECT COUNT(*) AS estacoes_auto FROM analytics.inmet_estacoes_automaticas"
	$(DUCKDB_BIN) $(COLLECTION_MACRO_DUCKDB) -c "SELECT COUNT(*) AS bdmep_diario FROM analytics.inmet_bdmep_diario"
	$(DUCKDB_BIN) $(COLLECTION_MACRO_DUCKDB) -c "SELECT COUNT(*) AS pacote_anual FROM analytics.inmet_pacote_anual_automaticas"
	$(DUCKDB_BIN) $(COLLECTION_MACRO_DUCKDB) -c "SELECT COUNT(*) AS ipca FROM analytics.bcb_sgs_ipca"
	$(DUCKDB_BIN) $(COLLECTION_MACRO_DUCKDB) -c "SELECT COUNT(*) AS ptax FROM analytics.bcb_sgs_ptax_usd_venda"
	$(DUCKDB_BIN) $(COLLECTION_MACRO_DUCKDB) -c "SELECT COUNT(*) AS soja_paranagua FROM analytics.cepea_soja_paranagua"
	$(DUCKDB_BIN) $(COLLECTION_MACRO_DUCKDB) -c "SELECT produto, praca, data, preco_rs_sc FROM analytics.cepea_milho LIMIT 3"

# Mirror GitHub Actions: offline macro collection (Phases 17–19).
ci-collection-macro-mvp:
	$(MAKE) collection-macro-mvp \
		COLLECTION_MACRO_LAKE=$(CI_COLLECTION_MACRO_LAKE) \
		COLLECTION_MACRO_DUCKDB=$(CI_COLLECTION_MACRO_DUCKDB)

# Sprint exit: P1 prices/logistics + macro indicators + PAM + ANP offline pipelines.
collection-full-mvp:
	$(MAKE) p1-collection-mvp
	$(MAKE) collection-macro-mvp
	LAKE_LOCAL_ROOT=/tmp/collection-pam-lake DUCKDB_PATH=/tmp/collection-pam.duckdb $(MAKE) ibge-pam-mvp
	LAKE_LOCAL_ROOT=/tmp/collection-anp-lake DUCKDB_PATH=/tmp/collection-anp.duckdb $(MAKE) anp-mvp

# Mirror GitHub Actions: full collection sprint exit (isolated /tmp lakes per sub-pipeline).
ci-collection-full-mvp:
	$(MAKE) ci-p1-collection-mvp
	$(MAKE) ci-collection-macro-mvp
	$(MAKE) ci-ibge-pam-mvp
	$(MAKE) ci-anp-mvp

# Phase 33 — collection hardening: full MVP + stress benchmark + history doc gate.
collection-hardening-mvp:
	$(MAKE) collection-full-mvp
	$(MAKE) benchmark-ingestor-fast10-stress
	python3 scripts/ci/check_data_history_ranges.py
	python3 scripts/ci/check_benchmark_profiles.py

# CI mirror: offline collection gates + Phase 33 documentation checks (no live backfill).
ci-collection-hardening-mvp:
	python3 scripts/ci/check_data_history_ranges.py
	python3 scripts/ci/check_benchmark_profiles.py
	$(MAKE) ci-collection-full-mvp

benchmark-ingestor:
	@test -f .env || (echo "copy .env.example to .env first" && exit 1)
	@set -a && . ./.env && set +a && python3 scripts/benchmark/ingestor_benchmark.py --json .local/benchmark/ingestor-latest.json

benchmark-ingestor-clean:
	@test -f .env || (echo "copy .env.example to .env first" && exit 1)
	@set -a && . ./.env && set +a && python3 scripts/benchmark/ingestor_benchmark.py --clean --all

benchmark-ingestor-fast10: build
	@test -f .env || (echo "copy .env.example to .env first" && exit 1)
	@set -a && . ./.env && set +a && python3 scripts/benchmark/ingestor_benchmark.py --profile $(BENCHMARK_PROFILE)

benchmark-ingestor-fast10-clean: build
	@test -f .env || (echo "copy .env.example to .env first" && exit 1)
	@set -a && . ./.env && set +a && python3 scripts/benchmark/ingestor_benchmark.py --profile $(BENCHMARK_PROFILE) --clean

benchmark-ingestor-fast10-stress: build
	@test -f .env || (echo "copy .env.example to .env first" && exit 1)
	@set -a && . ./.env && set +a && python3 scripts/benchmark/ingestor_benchmark.py --profile $(BENCHMARK_STRESS_PROFILE)

benchmark-ingestor-fast10-stress-clean: build
	@test -f .env || (echo "copy .env.example to .env first" && exit 1)
	@set -a && . ./.env && set +a && python3 scripts/benchmark/ingestor_benchmark.py --profile $(BENCHMARK_STRESS_PROFILE) --clean

lint:
	golangci-lint run ./...

build:
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/ingestor ./cmd/ingestor

build-processor:
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/processor ./cmd/processor

python-install:
	python3 -m pip install -r toolchain/python-requirements.txt

duckdb-install:
	curl -fsSL https://install.duckdb.org | DUCKDB_VERSION=$(DUCKDB_VERSION) sh
	@mkdir -p .local/bin
	@ln -sf "$(HOME)/.duckdb/cli/latest/duckdb" .local/bin/duckdb

dbt-deps:
	cp -f dbt/profiles.yml.example dbt/profiles.yml
	cd dbt && dbt deps --profiles-dir .

dbt-build: dbt-deps
	cd dbt && LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT) dbt build --profiles-dir . --select 'stg_conab__serie_historica_graos stg_conab__estimativa_graos+ mart_conab__serie_historica_graos stg_conab__estimativa_cana+ stg_conab__serie_historica_cana+ mart_conab__serie_historica_cana stg_conab__estimativa_cafe+ stg_conab__serie_historica_cafe+ mart_conab__serie_historica_cafe stg_conab__custo_producao+'

migrate-install:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.18.2

migrate-up:
	$(MIGRATE) -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" up

migrate-down:
	$(MIGRATE) -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" down 1

seed:
	go run ./cmd/seed

seed-mvp:
	go run ./cmd/seed --mvp

analytics-init:
	@chmod +x duckdb/scripts/analytics-init.sh duckdb/export-mart.sh
	LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT) DUCKDB_PATH=$(DUCKDB_PATH) DUCKDB_BIN="$(DUCKDB_BIN)" ./duckdb/scripts/analytics-init.sh

analytics-smoke:
	@test -x "$(DUCKDB_BIN)" || (echo "run: make duckdb-install" && exit 1)
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT COUNT(*) AS rows FROM analytics.conab_estimativa_graos"

conab-reference:
	@chmod +x scripts/conab/fetch_reference_samples.sh
	./scripts/conab/fetch_reference_samples.sh

dbt-build-mercado: dbt-deps
	cd dbt && LAKE_LOCAL_ROOT=$(LAKE_ABS) dbt build --profiles-dir . --select '$(MERCADO_DBT_SELECT)'

dbt-build-mercado-precos: dbt-deps
	cd dbt && LAKE_LOCAL_ROOT=$(LAKE_ABS) dbt build --profiles-dir . --select 'stg_conab__precos_semanal_uf+ stg_conab__precos_semanal_municipio+ stg_conab__precos_mensal_uf+ stg_conab__precos_mensal_municipio+ stg_conab__precos_minimos+'

dbt-build-mercado-prohort: dbt-deps
	cd dbt && LAKE_LOCAL_ROOT=$(LAKE_ABS) dbt build --profiles-dir . --select 'stg_conab__prohort_diario+ stg_conab__prohort_mensal+'

conab-mercado-full-mvp:
	go test ./internal/ingest/ -run 'OfertaDemanda|Precos|Prohort'
	LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT) python3 scripts/ci/seed_mercado_silver.py
	$(MAKE) dbt-build-mercado LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT)
	$(MAKE) analytics-init LAKE_LOCAL_ROOT=$(LAKE_ABS) DUCKDB_PATH=$(DUCKDB_PATH)
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT COUNT(*) AS oferta FROM analytics.conab_oferta_demanda"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT COUNT(*) AS precos_uf FROM analytics.conab_precos_semanal_uf"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT COUNT(*) AS precos_mun FROM analytics.conab_precos_semanal_municipio"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT COUNT(*) AS precos_min FROM analytics.conab_precos_minimos"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT COUNT(*) AS prohort_diario FROM analytics.conab_prohort_diario"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT COUNT(*) AS prohort_mensal FROM analytics.conab_prohort_mensal"
	$(MAKE) validate-codigo-ibge LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT)

conab-mercado-prohort-mvp:
	go test ./internal/ingest/ -run 'Prohort'
	LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT) python3 scripts/ci/seed_mercado_silver.py
	$(MAKE) dbt-build-mercado-prohort LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT)
	$(MAKE) analytics-init LAKE_LOCAL_ROOT=$(LAKE_ABS) DUCKDB_PATH=$(DUCKDB_PATH)
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT COUNT(*) AS diario_rows FROM analytics.conab_prohort_diario"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT COUNT(*) AS mensal_rows FROM analytics.conab_prohort_mensal"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT produto, municipio_ceasa, cod_ibge_municipio, preco_diario FROM analytics.conab_prohort_diario WHERE cod_ibge_municipio = '3550308' LIMIT 5"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT produto, municipio_ceasa, cod_ibge_municipio_ceasa, qtd_comercializada_kg FROM analytics.conab_prohort_mensal WHERE cod_ibge_municipio_ceasa = '3550308' LIMIT 5"
	$(MAKE) validate-codigo-ibge LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT)

conab-mercado-precos-minimos-mvp:
	go test ./internal/ingest/ -run 'PrecosMinimos'
	LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT) python3 scripts/ci/seed_mercado_silver.py
	$(MAKE) dbt-deps
	cd dbt && LAKE_LOCAL_ROOT=$(LAKE_ABS) dbt build --profiles-dir . --select 'stg_conab__precos_minimos+'
	$(MAKE) analytics-init LAKE_LOCAL_ROOT=$(LAKE_ABS) DUCKDB_PATH=$(DUCKDB_PATH)
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT COUNT(*) AS rows FROM analytics.conab_precos_minimos"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT produto, uf, ano_inicio_vigencia, preco, unidade_comercializacao FROM analytics.conab_precos_minimos WHERE upper(trim(produto)) = 'SOJA' AND uf = 'MT' ORDER BY ano_inicio_vigencia DESC LIMIT 5"

conab-mercado-mvp:
	LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT) python3 scripts/ci/seed_mercado_silver.py
	$(MAKE) dbt-build-mercado LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT)
	$(MAKE) analytics-init LAKE_LOCAL_ROOT=$(LAKE_ABS) DUCKDB_PATH=$(DUCKDB_PATH)
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT COUNT(*) AS rows FROM analytics.conab_oferta_demanda"

conab-mercado-precos-mvp:
	go test ./internal/ingest/ -run 'Precos'
	LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT) python3 scripts/ci/seed_mercado_silver.py
	$(MAKE) dbt-build-mercado-precos LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT)
	$(MAKE) analytics-init LAKE_LOCAL_ROOT=$(LAKE_ABS) DUCKDB_PATH=$(DUCKDB_PATH)
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT COUNT(*) AS uf_rows FROM analytics.conab_precos_semanal_uf"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT COUNT(*) AS mun_rows FROM analytics.conab_precos_semanal_municipio"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT COUNT(*) AS mensal_uf_rows FROM analytics.conab_precos_mensal_uf"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT COUNT(*) AS mensal_mun_rows FROM analytics.conab_precos_mensal_municipio"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT produto, municipio, cod_ibge, mes, valor_produto_kg FROM analytics.conab_precos_mensal_municipio WHERE upper(trim(produto)) = 'SOJA' AND cod_ibge = '5107925' ORDER BY ano, mes LIMIT 5"
	$(MAKE) validate-codigo-ibge LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT)

dbt-build-abastecimento: dbt-deps
	cd dbt && LAKE_LOCAL_ROOT=$(LAKE_ABS) dbt build --profiles-dir . --select 'stg_conab__estoques_publicos+ stg_conab__operacoes_comercializacao+ stg_conab__vendas_balcao+ stg_anp__combustiveis_precos_medios_municipios+ stg_anp__combustiveis_precos_postos+'

dbt-build-anp: dbt-deps
	cd dbt && LAKE_LOCAL_ROOT=$(LAKE_ABS) dbt build --profiles-dir . --select 'stg_anp__combustiveis_precos_medios_municipios+ stg_anp__combustiveis_precos_postos+'

anp-mvp:
	go test ./internal/ingest/ -run 'ANP'
	LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT) python3 scripts/ci/seed_anp_silver.py
	$(MAKE) dbt-build-anp LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT)
	$(MAKE) analytics-init LAKE_LOCAL_ROOT=$(LAKE_ABS) DUCKDB_PATH=$(DUCKDB_PATH)
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT COUNT(*) AS medios FROM analytics.anp_combustiveis_precos_medios_municipios"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT COUNT(*) AS postos FROM analytics.anp_combustiveis_precos_postos"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT estado, municipio, produto, preco_medio_revenda FROM analytics.anp_combustiveis_precos_medios_municipios ORDER BY estado, municipio LIMIT 5"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT municipio, estado, produto, preco_revenda FROM analytics.anp_combustiveis_precos_postos WHERE estado = 'SAO PAULO' LIMIT 5"

# Mirror GitHub Actions: offline ANP combustíveis pipeline (Phase 12 P2).
ci-anp-mvp:
	$(MAKE) anp-mvp \
		LAKE_LOCAL_ROOT=$(CI_ANP_LAKE) \
		DUCKDB_PATH=$(CI_ANP_DUCKDB)

conab-abastecimento-mvp:
	go test ./internal/ingest/ -run 'Estoques|Operacoes|VendasBalcao|ANP'
	LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT) python3 scripts/ci/seed_abastecimento_silver.py
	$(MAKE) dbt-build-abastecimento LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT)
	$(MAKE) analytics-init LAKE_LOCAL_ROOT=$(LAKE_ABS) DUCKDB_PATH=$(DUCKDB_PATH)
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT COUNT(*) AS estoques FROM analytics.conab_estoques_publicos"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT COUNT(*) AS operacoes FROM analytics.conab_operacoes_comercializacao"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT COUNT(*) AS vendas FROM analytics.conab_vendas_balcao"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT produto, uf_armazem_origem, qtd_negociada FROM analytics.conab_operacoes_comercializacao WHERE upper(trim(produto)) = 'SOJA' LIMIT 5"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT municipio_armazem_venda, uf, ano, mes, qtd_produto_kg FROM analytics.conab_vendas_balcao WHERE uf = 'MT' ORDER BY ano DESC, mes DESC LIMIT 5"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT COUNT(*) AS medios FROM analytics.anp_combustiveis_precos_medios_municipios"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT COUNT(*) AS postos FROM analytics.anp_combustiveis_precos_postos"
	$(MAKE) validate-codigo-ibge LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT)

dbt-build-armazenamento: dbt-deps
	cd dbt && LAKE_LOCAL_ROOT=$(LAKE_ABS) dbt build --profiles-dir . --select 'stg_conab__armazenagem+'

dbt-build-armazenamento-logistica: dbt-deps
	cd dbt && LAKE_LOCAL_ROOT=$(LAKE_ABS) dbt build --profiles-dir . --select 'stg_conab__frete+ stg_conab__capacidade_estatica+'

conab-armazenamento-mvp:
	LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT) python3 scripts/ci/seed_armazenamento_silver.py
	$(MAKE) dbt-build-armazenamento LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT)
	$(MAKE) analytics-init LAKE_LOCAL_ROOT=$(LAKE_ABS) DUCKDB_PATH=$(DUCKDB_PATH)
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT COUNT(*) AS armazens FROM analytics.conab_armazenagem"
	$(MAKE) validate-codigo-ibge LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT)

conab-armazenamento-logistica-mvp:
	go test ./internal/ingest/ -run 'Frete|Capacidade'
	LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT) python3 scripts/ci/seed_armazenamento_silver.py
	$(MAKE) dbt-build-armazenamento-logistica LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT)
	$(MAKE) analytics-init LAKE_LOCAL_ROOT=$(LAKE_ABS) DUCKDB_PATH=$(DUCKDB_PATH)
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT COUNT(*) AS frete_rows FROM analytics.conab_frete"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT COUNT(*) AS capacidade_rows FROM analytics.conab_capacidade_estatica"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT uf, ano, quantidade_mil_t FROM analytics.conab_capacidade_estatica WHERE uf IN ('MT', 'PR', 'RS') ORDER BY ano DESC LIMIT 5"
	$(MAKE) validate-codigo-ibge LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT)

dbt-build-agricultura-familiar: dbt-deps
	cd dbt && LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT) dbt build --profiles-dir . --select 'stg_conab__alimenta_brasil_entregas+ stg_conab__alimenta_brasil_propostas+'

conab-agricultura-familiar-mvp:
	LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT) python3 scripts/ci/seed_agricultura_familiar_silver.py
	$(MAKE) dbt-build-agricultura-familiar LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT)
	$(MAKE) analytics-init LAKE_LOCAL_ROOT=$(LAKE_ABS) DUCKDB_PATH=$(DUCKDB_PATH)
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT COUNT(*) AS entregas FROM analytics.conab_alimenta_brasil_entregas"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT COUNT(*) AS propostas FROM analytics.conab_alimenta_brasil_propostas"
	$(MAKE) validate-codigo-ibge LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT)

validate-codigo-ibge: python-install
	python3 scripts/quality/validate_codigo_ibge.py --lake-root $(LAKE_ABS)

validate-codigo-ibge-lake: python-install
	@test -f $(LAKE_ABS)/gold/mart_ibge__localidades_municipios/mart.parquet || \
		(echo "Missing gold mart — ingest + promote + dbt on ./lake first" && exit 1)
	$(MAKE) validate-codigo-ibge LAKE_LOCAL_ROOT=./lake

IBGE_LOCALIDADES_LIVE_DATASETS := ibge.localidades-municipios ibge.localidades-ufs ibge.localidades-regioes ibge.localidades-mesorregioes ibge.localidades-microrregioes

ibge-localidades-live-smoke: python-install
	@test -f .env || (echo "Missing .env — copy .env.example and set DATABASE_URL" && exit 1)
	@set -a && . ./.env && set +a && \
	for ds in $(IBGE_LOCALIDADES_LIVE_DATASETS); do \
		echo "==> ingest $$ds"; \
		go run ./cmd/ingestor run $$ds || exit 1; \
	done
	python3 scripts/ci/check_ibge_localidades_bronze.py --lake-root $(LAKE_ABS)

ibge-localidades-mvp:
	go test ./internal/ibge/... ./internal/ingest/ -run 'IBGE|Localidades|Flatten|3550308'
	LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT) python3 scripts/ci/seed_ibge_localidades_silver.py
	$(MAKE) dbt-build-ibge-localidades LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT)
	$(MAKE) analytics-init LAKE_LOCAL_ROOT=$(LAKE_ABS) DUCKDB_PATH=$(DUCKDB_PATH)
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT COUNT(*) AS municipios FROM analytics.ibge_localidades_municipios"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT codigo_ibge, nome, sigla_uf FROM analytics.ibge_localidades_municipios WHERE codigo_ibge = '3550308'"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT COUNT(*) AS ufs FROM analytics.ibge_localidades_ufs"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT COUNT(*) AS regioes FROM analytics.ibge_localidades_regioes"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT codigo_regiao, sigla_regiao, nome FROM analytics.ibge_localidades_regioes ORDER BY codigo_regiao"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT COUNT(*) AS mesorregioes FROM analytics.ibge_localidades_mesorregioes"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT codigo_mesorregiao, nome, sigla_uf FROM analytics.ibge_localidades_mesorregioes WHERE sigla_uf = 'MT' ORDER BY codigo_mesorregiao LIMIT 3"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT COUNT(*) AS microrregioes FROM analytics.ibge_localidades_microrregioes"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT codigo_microrregiao, nome, codigo_mesorregiao FROM analytics.ibge_localidades_microrregioes WHERE codigo_microrregiao = '51006'"

dbt-build-ibge-localidades: dbt-deps
	cd dbt && LAKE_LOCAL_ROOT=$(LAKE_ABS) dbt build --profiles-dir . --select 'stg_ibge__localidades_municipios+ stg_ibge__localidades_ufs+ stg_ibge__localidades_regioes+ stg_ibge__localidades_mesorregioes+ stg_ibge__localidades_microrregioes+'

dbt-build-ibge-pam: dbt-deps
	cd dbt && LAKE_LOCAL_ROOT=$(LAKE_ABS) dbt build --profiles-dir . --select 'stg_ibge__pam_area_quantidade+ stg_ibge__pam_rendimento_valor+ stg_ibge__pam_estabelecimentos+'

ibge-pam-mvp:
	go test ./internal/ibge/... ./internal/ingest/ -run 'PAM|SIDRA|pam|FlattenPAM'
	LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT) python3 scripts/ci/seed_ibge_pam_silver.py
	$(MAKE) dbt-build-ibge-pam LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT)
	$(MAKE) validate-codigo-ibge LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT)
	$(MAKE) analytics-init LAKE_LOCAL_ROOT=$(LAKE_ABS) DUCKDB_PATH=$(DUCKDB_PATH)
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT COUNT(*) FROM analytics.ibge_pam_area_quantidade"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT municipio, ano, variavel, valor FROM analytics.ibge_pam_rendimento_valor LIMIT 2"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT COUNT(*) FROM analytics.ibge_pam_estabelecimentos"

# Mirror GitHub Actions: offline IBGE PAM pipeline (Phase 16).
ci-ibge-pam-mvp:
	$(MAKE) ibge-pam-mvp \
		LAKE_LOCAL_ROOT=$(CI_IBGE_PAM_LAKE) \
		DUCKDB_PATH=$(CI_IBGE_PAM_DUCKDB)

inmet-clima-mvp:
	go test ./internal/inmet/... ./internal/ingest/ -run 'INMET|FlattenEstacoes|BDMEP|Missing'
	LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT) python3 scripts/ci/seed_inmet_silver.py
	$(MAKE) dbt-build-inmet-clima LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT)
	$(MAKE) analytics-init LAKE_LOCAL_ROOT=$(LAKE_ABS) DUCKDB_PATH=$(DUCKDB_PATH)
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT COUNT(*) FROM analytics.inmet_estacoes_automaticas"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT cd_estacao, data, variavel, valor FROM analytics.inmet_bdmep_diario LIMIT 2"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT COUNT(*) FROM analytics.inmet_pacote_anual_automaticas"

dbt-build-inmet-clima: dbt-deps
	cd dbt && LAKE_LOCAL_ROOT=$(LAKE_ABS) dbt build --profiles-dir . --select 'stg_inmet__estacoes_automaticas+ stg_inmet__estacoes_convencionais+ stg_inmet__bdmep_diario+ stg_inmet__bdmep_mensal+ stg_inmet__pacote_anual_automaticas+'

dbt-build-bcb-sgs: dbt-deps
	cd dbt && LAKE_LOCAL_ROOT=$(LAKE_ABS) dbt build --profiles-dir . --select 'stg_bcb__sgs_ipca+ stg_bcb__sgs_ipca_12m+ stg_bcb__sgs_igpm+ stg_bcb__sgs_ptax_usd_venda+ stg_bcb__sgs_ptax_usd_compra+'

bcb-sgs-mvp:
	go test ./internal/bcb/... ./internal/ingest/ -run 'BCB|SGS|FlattenSGS'
	LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT) python3 scripts/ci/seed_bcb_sgs_silver.py
	$(MAKE) dbt-build-bcb-sgs LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT)
	$(MAKE) analytics-init LAKE_LOCAL_ROOT=$(LAKE_ABS) DUCKDB_PATH=$(DUCKDB_PATH)
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT COUNT(*) FROM analytics.bcb_sgs_ipca"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT COUNT(*) FROM analytics.bcb_sgs_ptax_usd_venda"

cepea-indicadores-mvp:
	go test ./internal/cepea/... ./internal/ingest/ -run 'CEPA|Cepea|FlattenIndicador|ParseIndicator'
	LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT) python3 scripts/ci/seed_cepea_silver.py
	$(MAKE) dbt-build-cepea LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT)
	$(MAKE) analytics-init LAKE_LOCAL_ROOT=$(LAKE_ABS) DUCKDB_PATH=$(DUCKDB_PATH)
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT COUNT(*) FROM analytics.cepea_soja_paranagua"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT produto, praca, data, preco_rs_sc FROM analytics.cepea_milho LIMIT 2"

dbt-build-cepea: dbt-deps
	cd dbt && LAKE_LOCAL_ROOT=$(LAKE_ABS) dbt build --profiles-dir . --select 'stg_cepea__soja_paranagua+ stg_cepea__soja_parana+ stg_cepea__milho+ stg_cepea__boi_gordo+'

dbt-build-mdic-comex: dbt-deps
	cd dbt && LAKE_LOCAL_ROOT=$(LAKE_ABS) dbt build --profiles-dir . --select 'stg_mdic__comex_exportacao_ncm_mes+'

mdic-comex-mvp:
	go test ./internal/mdic/... ./internal/ingest/ -run 'MDIC|Comex|FlattenComex'
	LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT) python3 scripts/ci/seed_mdic_silver.py
	$(MAKE) dbt-build-mdic-comex LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT)
	$(MAKE) analytics-init LAKE_LOCAL_ROOT=$(LAKE_ABS) DUCKDB_PATH=$(DUCKDB_PATH)
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT COUNT(*) FROM analytics.mdic_comex_exportacao_ncm_mes"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT co_ncm, produto_slug, data, valor_fob_usd FROM analytics.mdic_comex_exportacao_ncm_mes LIMIT 3"

ci-mdic-comex-mvp:
	$(MAKE) mdic-comex-mvp \
		LAKE_LOCAL_ROOT=$(CI_MDIC_LAKE) \
		DUCKDB_PATH=$(CI_MDIC_DUCKDB)

dbt-build-mdic-comex-extended: dbt-deps
	cd dbt && LAKE_LOCAL_ROOT=$(LAKE_ABS) dbt build --profiles-dir . --select 'stg_mdic__comex_importacao_ncm_mes+ stg_mdic__comex_exportacao_uf_ncm+ stg_mdic__comex_importacao_diesel_ncm+'

mdic-comex-extended-mvp:
	go test ./internal/mdic/... ./internal/ingest/ -run 'MDIC|Comex|FlattenComex'
	LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT) python3 scripts/ci/seed_mdic_comex_extended_silver.py
	$(MAKE) dbt-build-mdic-comex-extended LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT)
	$(MAKE) analytics-init LAKE_LOCAL_ROOT=$(LAKE_ABS) DUCKDB_PATH=$(DUCKDB_PATH)
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT COUNT(*) FROM analytics.mdic_comex_importacao_ncm_mes"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT co_ncm, produto_slug, data, valor_cif_usd FROM analytics.mdic_comex_importacao_ncm_mes LIMIT 3"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT COUNT(*) FROM analytics.mdic_comex_exportacao_uf_ncm"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT uf, co_ncm, produto_slug, valor_fob_usd FROM analytics.mdic_comex_exportacao_uf_ncm LIMIT 3"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT COUNT(*) FROM analytics.mdic_comex_importacao_diesel_ncm"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT produto_slug, data, valor_cif_usd FROM analytics.mdic_comex_importacao_diesel_ncm LIMIT 3"

ci-mdic-comex-extended-mvp:
	$(MAKE) mdic-comex-extended-mvp \
		LAKE_LOCAL_ROOT=$(CI_MDIC_COMEX_EXTENDED_LAKE) \
		DUCKDB_PATH=$(CI_MDIC_COMEX_EXTENDED_DUCKDB)

dbt-build-antt-pedagios: dbt-deps
	cd dbt && LAKE_LOCAL_ROOT=$(LAKE_ABS) dbt build --profiles-dir . --select 'stg_antt__pracas_pedagio+'

dnit-antt-logistica-mvp:
	go test ./internal/antt/... ./internal/ingest/ -run 'ANTT|PracasPedagio|ResolveURLCKAN'
	LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT) python3 scripts/ci/seed_antt_pedagios_silver.py
	$(MAKE) dbt-build-antt-pedagios LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT)
	$(MAKE) analytics-init LAKE_LOCAL_ROOT=$(LAKE_ABS) DUCKDB_PATH=$(DUCKDB_PATH)
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT COUNT(*) FROM analytics.antt_pracas_pedagio"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT rodovia, uf, praca_de_pedagio, situacao FROM analytics.antt_pracas_pedagio LIMIT 3"

ci-dnit-antt-logistica-mvp:
	$(MAKE) dnit-antt-logistica-mvp \
		LAKE_LOCAL_ROOT=$(CI_ANTT_PEDAGIOS_LAKE) \
		DUCKDB_PATH=$(CI_ANTT_PEDAGIOS_DUCKDB)

dbt-build-antt-logistica-extended: dbt-deps
	cd dbt && LAKE_LOCAL_ROOT=$(LAKE_ABS) dbt build --profiles-dir . --select 'stg_antt__volume_trafego_pedagio+ stg_antt__receita_por_praca+'

br-logistica-extended-mvp:
	go test ./internal/antt/... ./internal/ingest/ -run 'ANTT|PracasPedagio|VolumeTrafego|ReceitaPorPraca|ResolveURLCKAN'
	LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT) python3 scripts/ci/seed_antt_logistica_extended_silver.py
	$(MAKE) dbt-build-antt-logistica-extended LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT)
	$(MAKE) analytics-init LAKE_LOCAL_ROOT=$(LAKE_ABS) DUCKDB_PATH=$(DUCKDB_PATH)
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT COUNT(*) FROM analytics.antt_volume_trafego_pedagio"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT concessionaria, praca, mes_ano, volume_total FROM analytics.antt_volume_trafego_pedagio LIMIT 3"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT COUNT(*) FROM analytics.antt_receita_por_praca"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT concessionaria, praca_de_pedagio, rodovia, mes_ano FROM analytics.antt_receita_por_praca LIMIT 3"

ci-br-logistica-extended-mvp:
	$(MAKE) br-logistica-extended-mvp \
		LAKE_LOCAL_ROOT=$(CI_ANTT_LOGISTICA_EXTENDED_LAKE) \
		DUCKDB_PATH=$(CI_ANTT_LOGISTICA_EXTENDED_DUCKDB)

dbt-build-mapa-zarc: dbt-deps
	cd dbt && LAKE_LOCAL_ROOT=$(LAKE_ABS) dbt build --profiles-dir . --select 'stg_mapa__zarc_tabua_risco+'

mapa-dados-mvp:
	go test ./internal/mapa/... ./internal/ingest/ -run 'MAPA|ZARC|ZarcTabua|ResolveURLLatestSafra'
	LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT) python3 scripts/ci/seed_mapa_zarc_silver.py
	$(MAKE) dbt-build-mapa-zarc LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT)
	$(MAKE) analytics-init LAKE_LOCAL_ROOT=$(LAKE_ABS) DUCKDB_PATH=$(DUCKDB_PATH)
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT COUNT(*) FROM analytics.mapa_zarc_tabua_risco"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT nome_cultura, uf, municipio, safra_ini, safra_fim FROM analytics.mapa_zarc_tabua_risco LIMIT 3"

ci-mapa-dados-mvp:
	$(MAKE) mapa-dados-mvp \
		LAKE_LOCAL_ROOT=$(CI_MAPA_ZARC_LAKE) \
		DUCKDB_PATH=$(CI_MAPA_ZARC_DUCKDB)

dbt-build-b3-futuros: dbt-deps
	cd dbt && LAKE_LOCAL_ROOT=$(LAKE_ABS) dbt build --profiles-dir . --select 'stg_b3__futuro_soja+ stg_b3__futuro_milho+ stg_b3__futuro_boi+'

b3-futuros-mvp:
	go test ./internal/b3/... ./internal/ingest/ -run 'B3|Futuro|FlattenFuturo'
	LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT) python3 scripts/ci/seed_b3_futuros_silver.py
	$(MAKE) dbt-build-b3-futuros LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT)
	$(MAKE) analytics-init LAKE_LOCAL_ROOT=$(LAKE_ABS) DUCKDB_PATH=$(DUCKDB_PATH)
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT COUNT(*) FROM analytics.b3_futuro_soja"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT COUNT(*) FROM analytics.b3_futuro_milho"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT COUNT(*) FROM analytics.b3_futuro_boi"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT refdate, symbol, price, currency FROM analytics.b3_futuro_soja LIMIT 3"

ci-b3-futuros-mvp:
	$(MAKE) b3-futuros-mvp \
		LAKE_LOCAL_ROOT=$(CI_B3_FUTUROS_LAKE) \
		DUCKDB_PATH=$(CI_B3_FUTUROS_DUCKDB)

dbt-build-usda-psd: dbt-deps
	cd dbt && LAKE_LOCAL_ROOT=$(LAKE_ABS) dbt build --profiles-dir . --select 'stg_usda__psd_soja+ stg_usda__psd_milho+ stg_usda__psd_trigo+'

usda-psd-mvp:
	go test ./internal/usda/... ./internal/ingest/ -run 'USDA|PSD|FlattenPSD'
	LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT) python3 scripts/ci/seed_usda_psd_silver.py
	$(MAKE) dbt-build-usda-psd LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT)
	$(MAKE) analytics-init LAKE_LOCAL_ROOT=$(LAKE_ABS) DUCKDB_PATH=$(DUCKDB_PATH)
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT COUNT(*) FROM analytics.usda_psd_soja"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT COUNT(*) FROM analytics.usda_psd_milho"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT COUNT(*) FROM analytics.usda_psd_trigo"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT country_code, marketing_year, attribute_name, value FROM analytics.usda_psd_soja LIMIT 3"

ci-usda-psd-mvp:
	$(MAKE) usda-psd-mvp \
		LAKE_LOCAL_ROOT=$(CI_USDA_PSD_LAKE) \
		DUCKDB_PATH=$(CI_USDA_PSD_DUCKDB)

dbt-build-fao-faostat: dbt-deps
	cd dbt && LAKE_LOCAL_ROOT=$(LAKE_ABS) dbt build --profiles-dir . --select 'stg_fao__prices_agro+'

fao-faostat-mvp:
	go test ./internal/fao/... ./internal/ingest/ -run 'FAO|Prices|FlattenPrices'
	LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT) python3 scripts/ci/seed_fao_faostat_silver.py
	$(MAKE) dbt-build-fao-faostat LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT)
	$(MAKE) analytics-init LAKE_LOCAL_ROOT=$(LAKE_ABS) DUCKDB_PATH=$(DUCKDB_PATH)
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT COUNT(*) FROM analytics.fao_prices_agro"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT area_name, commodity_slug, year, element_name, value FROM analytics.fao_prices_agro LIMIT 3"

ci-fao-faostat-mvp:
	$(MAKE) fao-faostat-mvp \
		LAKE_LOCAL_ROOT=$(CI_FAO_FAOSTAT_LAKE) \
		DUCKDB_PATH=$(CI_FAO_FAOSTAT_DUCKDB)

dbt-build-worldbank-commodities: dbt-deps
	cd dbt && LAKE_LOCAL_ROOT=$(LAKE_ABS) dbt build --profiles-dir . --select 'stg_worldbank__pink_sheet_monthly+'

worldbank-commodities-mvp:
	go test ./internal/worldbank/... ./internal/ingest/ -run 'WorldBank|PinkSheet|FlattenPinkSheet'
	LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT) python3 scripts/ci/seed_worldbank_commodities_silver.py
	$(MAKE) dbt-build-worldbank-commodities LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT)
	$(MAKE) analytics-init LAKE_LOCAL_ROOT=$(LAKE_ABS) DUCKDB_PATH=$(DUCKDB_PATH)
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT COUNT(*) FROM analytics.worldbank_pink_sheet_monthly"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT refmonth, commodity_slug, series_name, value FROM analytics.worldbank_pink_sheet_monthly LIMIT 3"

ci-worldbank-commodities-mvp:
	$(MAKE) worldbank-commodities-mvp \
		LAKE_LOCAL_ROOT=$(CI_WORLDBANK_COMMODITIES_LAKE) \
		DUCKDB_PATH=$(CI_WORLDBANK_COMMODITIES_DUCKDB)

dbt-build-international-extended: dbt-deps
	cd dbt && LAKE_LOCAL_ROOT=$(LAKE_ABS) dbt build --profiles-dir . --select 'stg_fao__producao_agro+ stg_fao__comercio_agro+ stg_worldbank__ag_indices+'

international-extended-mvp:
	go test ./internal/fao/... ./internal/worldbank/... ./internal/ingest/ -run 'Producao|Comercio|AgIndices|GoldenVector'
	LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT) python3 scripts/ci/seed_international_extended_silver.py
	$(MAKE) dbt-build-international-extended LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT)
	$(MAKE) analytics-init LAKE_LOCAL_ROOT=$(LAKE_ABS) DUCKDB_PATH=$(DUCKDB_PATH)
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT COUNT(*) FROM analytics.fao_producao_agro"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT area_name, commodity_slug, year, value FROM analytics.fao_producao_agro LIMIT 3"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT COUNT(*) FROM analytics.fao_comercio_agro"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT area_name, commodity_slug, element_name, year FROM analytics.fao_comercio_agro LIMIT 3"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT COUNT(*) FROM analytics.worldbank_ag_indices"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT refmonth, commodity_slug, series_name, value FROM analytics.worldbank_ag_indices LIMIT 3"

ci-international-extended-mvp:
	$(MAKE) international-extended-mvp \
		LAKE_LOCAL_ROOT=$(CI_INTERNATIONAL_EXTENDED_LAKE) \
		DUCKDB_PATH=$(CI_INTERNATIONAL_EXTENDED_DUCKDB)

dbt-build-br-new-sources: dbt-deps
	cd dbt && LAKE_LOCAL_ROOT=$(LAKE_ABS) dbt build --profiles-dir . --select 'stg_ibge__lspa_area_producao+ stg_bcb__sgs_selic+'

br-new-sources-mvp:
	go test ./internal/ibge/... ./internal/bcb/... ./internal/ingest/ -run 'LSPA|FlattenLSPA|SGS|GoldenVector'
	LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT) python3 scripts/ci/seed_br_new_sources_silver.py
	$(MAKE) dbt-build-br-new-sources LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT)
	$(MAKE) analytics-init LAKE_LOCAL_ROOT=$(LAKE_ABS) DUCKDB_PATH=$(DUCKDB_PATH)
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT COUNT(*) FROM analytics.ibge_lspa_area_producao"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT uf, mes, produto_slug, variavel, valor FROM analytics.ibge_lspa_area_producao LIMIT 3"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT COUNT(*) FROM analytics.bcb_sgs_selic"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT data, valor FROM analytics.bcb_sgs_selic LIMIT 3"

ci-br-new-sources-mvp:
	$(MAKE) br-new-sources-mvp \
		LAKE_LOCAL_ROOT=$(CI_BR_NEW_SOURCES_LAKE) \
		DUCKDB_PATH=$(CI_BR_NEW_SOURCES_DUCKDB)

dbt-build-international-new-sources: dbt-deps
	cd dbt && LAKE_LOCAL_ROOT=$(LAKE_ABS) dbt build --profiles-dir . --select 'stg_eia__petroleum_prices+ stg_usda__wasde+ stg_igc__goi_index+ stg_un__comtrade_bulk+'

international-new-sources-mvp:
	go test ./internal/eia/... ./internal/usda/... ./internal/igc/... ./internal/un/... ./internal/ingest/ -run 'Petroleum|FlattenPetroleum|WASDE|FlattenWASDE|FlattenGOI|Comtrade|GoldenVector'
	LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT) python3 scripts/ci/seed_international_new_sources_silver.py
	$(MAKE) dbt-build-international-new-sources LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT)
	$(MAKE) analytics-init LAKE_LOCAL_ROOT=$(LAKE_ABS) DUCKDB_PATH=$(DUCKDB_PATH)
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT COUNT(*) FROM analytics.eia_petroleum_prices"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT series_name, commodity_slug, refdate, value FROM analytics.eia_petroleum_prices LIMIT 3"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT COUNT(*) FROM analytics.usda_wasde"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT commodity, market_year, attribute, value FROM analytics.usda_wasde LIMIT 3"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT COUNT(*) FROM analytics.igc_goi_index"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT index_slug, refdate, value FROM analytics.igc_goi_index LIMIT 3"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT COUNT(*) FROM analytics.un_comtrade_bulk"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT flow_code, hs_code, commodity_slug, trade_value_usd FROM analytics.un_comtrade_bulk LIMIT 3"

ci-international-new-sources-mvp:
	$(MAKE) international-new-sources-mvp \
		LAKE_LOCAL_ROOT=$(CI_INTERNATIONAL_NEW_SOURCES_LAKE) \
		DUCKDB_PATH=$(CI_INTERNATIONAL_NEW_SOURCES_DUCKDB)

CI_BR_SOURCES_WAVE_2_LAKE ?= /tmp/br-sources-wave-2-ci-lake
CI_BR_SOURCES_WAVE_2_DUCKDB ?= /tmp/br-sources-wave-2-ci.duckdb
CI_INTERNATIONAL_SOURCES_WAVE_2_LAKE ?= /tmp/international-sources-wave-2-ci-lake
CI_INTERNATIONAL_SOURCES_WAVE_2_DUCKDB ?= /tmp/international-sources-wave-2-ci.duckdb

dbt-build-br-sources-wave-2: dbt-deps
	cd dbt && LAKE_LOCAL_ROOT=$(LAKE_ABS) dbt build --profiles-dir . --select 'stg_mapa__agrofit_produtos_formulados+ stg_mapa__agrofit_produtos_tecnicos+ stg_ana__hidrologia_series+ stg_antaq__movimentacao_carga_portuaria+'

br-sources-wave-2-mvp:
	go test ./internal/mapa/... ./internal/ana/... ./internal/antaq/... ./internal/ingest/ -run 'Agrofit|ANA|ANTAQ|Hidrologia|GoldenVector'
	LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT) python3 scripts/ci/seed_br_sources_wave_2_silver.py
	$(MAKE) dbt-build-br-sources-wave-2 LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT)
	$(MAKE) analytics-init LAKE_LOCAL_ROOT=$(LAKE_ABS) DUCKDB_PATH=$(DUCKDB_PATH)
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT COUNT(*) FROM analytics.mapa_agrofit_produtos_formulados"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT nr_registro, marca_comercial, cultura FROM analytics.mapa_agrofit_produtos_formulados LIMIT 3"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT COUNT(*) FROM analytics.mapa_agrofit_produtos_tecnicos"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT numero_registro, produto_tecnico_marca_comercial FROM analytics.mapa_agrofit_produtos_tecnicos LIMIT 3"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT COUNT(*) FROM analytics.ana_hidrologia_series"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT station_code, observed_at, daily_mean FROM analytics.ana_hidrologia_series LIMIT 3"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT COUNT(*) FROM analytics.antaq_movimentacao_carga_portuaria"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT ano, mes, nome_instalacao_portuaria, peso_toneladas FROM analytics.antaq_movimentacao_carga_portuaria LIMIT 3"

ci-br-sources-wave-2-mvp:
	$(MAKE) br-sources-wave-2-mvp \
		LAKE_LOCAL_ROOT=$(CI_BR_SOURCES_WAVE_2_LAKE) \
		DUCKDB_PATH=$(CI_BR_SOURCES_WAVE_2_DUCKDB)

dbt-build-international-sources-wave-2: dbt-deps
	cd dbt && LAKE_LOCAL_ROOT=$(LAKE_ABS) dbt build --profiles-dir . --select 'stg_igc__goi_index+ stg_usda__gats_trade+ stg_eurostat__ag_prices+ stg_argentina__bcra_cambio+'

international-sources-wave-2-mvp:
	go test ./internal/igc/... ./internal/eurostat/... ./internal/argentina/... ./internal/usda/... ./internal/ingest/ -run 'IGC|GOI|GATS|Eurostat|Argentina|BCRA|GoldenVector'
	LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT) python3 scripts/ci/seed_international_sources_wave_2_silver.py
	$(MAKE) dbt-build-international-sources-wave-2 LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT)
	$(MAKE) analytics-init LAKE_LOCAL_ROOT=$(LAKE_ABS) DUCKDB_PATH=$(DUCKDB_PATH)
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT COUNT(*) FROM analytics.igc_goi_index"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT index_name, index_slug, refdate, value FROM analytics.igc_goi_index LIMIT 3"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT COUNT(*) FROM analytics.usda_gats_trade"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT commodity_code, partner_name, year, value FROM analytics.usda_gats_trade LIMIT 3"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT COUNT(*) FROM analytics.eurostat_ag_prices"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT product_name, year, index_value FROM analytics.eurostat_ag_prices LIMIT 3"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT COUNT(*) FROM analytics.argentina_bcra_cambio"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT refdate, exchange_rate FROM analytics.argentina_bcra_cambio LIMIT 3"

ci-international-sources-wave-2-mvp:
	$(MAKE) international-sources-wave-2-mvp \
		LAKE_LOCAL_ROOT=$(CI_INTERNATIONAL_SOURCES_WAVE_2_LAKE) \
		DUCKDB_PATH=$(CI_INTERNATIONAL_SOURCES_WAVE_2_DUCKDB)

ingestor-signoff-mvp: br-sources-wave-2-mvp international-sources-wave-2-mvp
	@echo "ingestor wave 2 MVPs passed"

ci-ingestor-signoff-mvp:
	$(MAKE) ci-br-sources-wave-2-mvp
	$(MAKE) ci-international-sources-wave-2-mvp
	@echo "ingestor signoff CI passed"

dbt-build-noaa-climate: dbt-deps
	cd dbt && LAKE_LOCAL_ROOT=$(LAKE_ABS) dbt build --profiles-dir . --select 'stg_noaa__enso_indices+ stg_noaa__global_temp_anomaly+'

noaa-climate-mvp:
	go test ./internal/noaa/... ./internal/ingest/ -run 'NOAA|ENSO|GlobalTemp|FlattenENSO|FlattenGlobalTemp|FlattenClimate'
	LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT) python3 scripts/ci/seed_noaa_climate_silver.py
	$(MAKE) dbt-build-noaa-climate LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT)
	$(MAKE) analytics-init LAKE_LOCAL_ROOT=$(LAKE_ABS) DUCKDB_PATH=$(DUCKDB_PATH)
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT COUNT(*) FROM analytics.noaa_enso_indices"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT COUNT(*) FROM analytics.noaa_global_temp_anomaly"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT refyear, season_code, anomaly FROM analytics.noaa_enso_indices LIMIT 3"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT refmonth, anomaly FROM analytics.noaa_global_temp_anomaly LIMIT 3"

ci-noaa-climate-mvp:
	$(MAKE) noaa-climate-mvp \
		LAKE_LOCAL_ROOT=$(CI_NOAA_CLIMATE_LAKE) \
		DUCKDB_PATH=$(CI_NOAA_CLIMATE_DUCKDB)

unified-db-sync: build-processor migrate-up
	@test -n "$(DATABASE_URL)"
	LAKE_LOCAL_ROOT=$(LAKE_ABS) DATABASE_URL="$(DATABASE_URL)" DUCKDB_BIN="$(DUCKDB_BIN)" DUCKDB_PATH=$(DUCKDB_PATH) UNIFIED_DB_SYNC_MARTS="$(UNIFIED_DB_SYNC_MARTS)" ./$(BIN_DIR)/processor sync-postgres
	DATABASE_URL="$(DATABASE_URL)" python3 scripts/ci/verify_unified_db_sync.py --lake-root "$(LAKE_ABS)"
	@(command -v psql >/dev/null 2>&1 && psql "$(DATABASE_URL)" -c "\dt analytics.*") || docker compose exec -T postgres psql -U open_data_agro -d open_data_agro -c "\dt analytics.*"

ci-unified-db-sync: seed-unified-db-gold
	$(MAKE) unified-db-sync \
		LAKE_LOCAL_ROOT=$(CI_UNIFIED_DB_LAKE) \
		DATABASE_URL="$(DATABASE_URL)"

seed-unified-db-gold:
	LAKE_LOCAL_ROOT=$(CI_UNIFIED_DB_LAKE) python3 scripts/ci/seed_unified_db_gold.py

conab-mvp:
	LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT) python3 scripts/ci/seed_dbt_silver.py
	$(MAKE) dbt-build LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT)
	$(MAKE) analytics-init LAKE_LOCAL_ROOT=$(LAKE_ABS) DUCKDB_PATH=$(DUCKDB_PATH)
	$(MAKE) analytics-smoke DUCKDB_PATH=$(DUCKDB_PATH)
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT COUNT(*) FROM analytics.conab_serie_historica_graos"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT COUNT(*) FROM analytics.conab_estimativa_cana"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT COUNT(*) FROM analytics.conab_serie_historica_cana"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT COUNT(*) FROM analytics.conab_estimativa_cafe"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT COUNT(*) FROM analytics.conab_serie_historica_cafe"
	$(DUCKDB_BIN) $(DUCKDB_PATH) -c "SELECT COUNT(*) FROM analytics.conab_custo_producao"
	$(MAKE) validate-codigo-ibge LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT)

clean:
	rm -rf $(BIN_DIR)
