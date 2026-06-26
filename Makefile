.PHONY: test lint build build-processor clean duckdb-install python-install dbt-deps dbt-build dbt-build-mercado dbt-build-mercado-precos dbt-build-abastecimento dbt-build-armazenamento dbt-build-agricultura-familiar dbt-build-ibge-localidades dbt-build-ibge-pam dbt-build-bcb-sgs ibge-localidades-mvp ibge-pam-mvp inmet-clima-mvp bcb-sgs-mvp cepea-indicadores-mvp ci-go ci-dbt benchmark-ingestor benchmark-ingestor-clean benchmark-ingestor-fast10 benchmark-ingestor-fast10-clean migrate-install migrate-up migrate-down seed analytics-init analytics-smoke conab-reference conab-mvp conab-mercado-mvp conab-mercado-precos-mvp conab-abastecimento-mvp conab-armazenamento-mvp conab-agricultura-familiar-mvp

BIN_DIR := bin
DUCKDB_VERSION ?= 1.5.4
LAKE_LOCAL_ROOT ?= ./lake
LAKE_ABS := $(abspath $(LAKE_LOCAL_ROOT))
DUCKDB_PATH ?= ./duckdb/open_data_agro.duckdb
POSTGRES_HOST_PORT ?= 5432
DATABASE_URL ?= postgresql://open_data_agro:open_data_agro@localhost:$(POSTGRES_HOST_PORT)/open_data_agro?sslmode=disable
MIGRATIONS_PATH := infra/postgres/migrations
MIGRATE ?= migrate

test:
	go test ./...

# Mirror GitHub Actions CI jobs locally (see .github/workflows/ci.yml).
ci-go: duckdb-install python-install
	go work sync
	go test ./...
	PATH="$(PWD)/.local/bin:$$PATH" DUCKDB_BIN="$(PWD)/.local/bin/duckdb" DUCKDB_INTEGRATION=1 go test ./internal/processor -run 'SmokeLocal|PreviewPromote' -count=1
	GE_INTEGRATION=1 go test ./internal/processor -run 'Quality' -count=1

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
	duckdb /tmp/open-data-agro-analytics.duckdb -c "SELECT COUNT(*) FROM analytics.conab_oferta_demanda"
	duckdb /tmp/open-data-agro-analytics.duckdb -c "SELECT * FROM analytics.conab_estimativa_graos LIMIT 10"

benchmark-ingestor:
	@test -f .env || (echo "copy .env.example to .env first" && exit 1)
	@set -a && . ./.env && set +a && python3 scripts/benchmark/ingestor_benchmark.py --json .local/benchmark/ingestor-latest.json

benchmark-ingestor-clean:
	@test -f .env || (echo "copy .env.example to .env first" && exit 1)
	@set -a && . ./.env && set +a && python3 scripts/benchmark/ingestor_benchmark.py --clean --all

benchmark-ingestor-fast10: build
	@test -f .env || (echo "copy .env.example to .env first" && exit 1)
	@set -a && . ./.env && set +a && python3 scripts/benchmark/ingestor_benchmark.py --profile .local/benchmark/profiles/fast10.json

benchmark-ingestor-fast10-clean: build
	@test -f .env || (echo "copy .env.example to .env first" && exit 1)
	@set -a && . ./.env && set +a && python3 scripts/benchmark/ingestor_benchmark.py --profile .local/benchmark/profiles/fast10.json --clean

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
	cd dbt && dbt deps --profiles-dir .

dbt-build: dbt-deps
	cd dbt && LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT) dbt build --profiles-dir . --select 'stg_conab__serie_historica_graos stg_conab__estimativa_graos+ mart_conab__serie_historica_graos'

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
	LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT) DUCKDB_PATH=$(DUCKDB_PATH) DUCKDB_BIN=$(DUCKDB_BIN) ./duckdb/scripts/analytics-init.sh

analytics-smoke:
	@command -v duckdb >/dev/null 2>&1 || (echo "run: make duckdb-install" && exit 1)
	duckdb $(DUCKDB_PATH) -c "SELECT COUNT(*) AS rows FROM analytics.conab_estimativa_graos"

conab-reference:
	@chmod +x scripts/conab/fetch_reference_samples.sh
	./scripts/conab/fetch_reference_samples.sh

dbt-build-mercado: dbt-deps
	cd dbt && LAKE_LOCAL_ROOT=$(LAKE_ABS) dbt build --profiles-dir . --select 'stg_conab__oferta_demanda+'

dbt-build-mercado-precos: dbt-deps
	cd dbt && LAKE_LOCAL_ROOT=$(LAKE_ABS) dbt build --profiles-dir . --select 'stg_conab__precos_semanal_uf+ stg_conab__precos_semanal_municipio+'

conab-mercado-mvp:
	LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT) python3 scripts/ci/seed_mercado_silver.py
	$(MAKE) dbt-build-mercado LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT)
	$(MAKE) analytics-init LAKE_LOCAL_ROOT=$(LAKE_ABS) DUCKDB_PATH=$(DUCKDB_PATH)
	duckdb $(DUCKDB_PATH) -c "SELECT COUNT(*) AS rows FROM analytics.conab_oferta_demanda"

conab-mercado-precos-mvp:
	go test ./internal/ingest/ -run 'PrecosSemanal'
	LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT) python3 scripts/ci/seed_mercado_silver.py
	$(MAKE) dbt-build-mercado-precos LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT)
	$(MAKE) analytics-init LAKE_LOCAL_ROOT=$(LAKE_ABS) DUCKDB_PATH=$(DUCKDB_PATH)
	duckdb $(DUCKDB_PATH) -c "SELECT COUNT(*) AS uf_rows FROM analytics.conab_precos_semanal_uf"
	duckdb $(DUCKDB_PATH) -c "SELECT COUNT(*) AS mun_rows FROM analytics.conab_precos_semanal_municipio"
	duckdb $(DUCKDB_PATH) -c "SELECT produto, municipio, cod_ibge, valor_produto_kg FROM analytics.conab_precos_semanal_municipio WHERE upper(trim(produto)) = 'SOJA' AND cod_ibge = '5107925' LIMIT 5"

dbt-build-abastecimento: dbt-deps
	cd dbt && LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT) dbt build --profiles-dir . --select 'stg_conab__estoques_publicos+ stg_anp__combustiveis_precos_medios_municipios+ stg_anp__combustiveis_precos_postos+'

conab-abastecimento-mvp:
	LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT) python3 scripts/ci/seed_abastecimento_silver.py
	$(MAKE) dbt-build-abastecimento LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT)
	$(MAKE) analytics-init LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT) DUCKDB_PATH=$(DUCKDB_PATH)
	duckdb $(DUCKDB_PATH) -c "SELECT COUNT(*) AS estoques FROM analytics.conab_estoques_publicos"
	duckdb $(DUCKDB_PATH) -c "SELECT COUNT(*) AS medios FROM analytics.anp_combustiveis_precos_medios_municipios"
	duckdb $(DUCKDB_PATH) -c "SELECT COUNT(*) AS postos FROM analytics.anp_combustiveis_precos_postos"

dbt-build-armazenamento: dbt-deps
	cd dbt && LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT) dbt build --profiles-dir . --select 'stg_conab__armazenagem+'

conab-armazenamento-mvp:
	LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT) python3 scripts/ci/seed_armazenamento_silver.py
	$(MAKE) dbt-build-armazenamento LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT)
	$(MAKE) analytics-init LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT) DUCKDB_PATH=$(DUCKDB_PATH)
	duckdb $(DUCKDB_PATH) -c "SELECT COUNT(*) AS armazens FROM analytics.conab_armazenagem"

dbt-build-agricultura-familiar: dbt-deps
	cd dbt && LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT) dbt build --profiles-dir . --select 'stg_conab__alimenta_brasil_entregas+ stg_conab__alimenta_brasil_propostas+'

conab-agricultura-familiar-mvp:
	LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT) python3 scripts/ci/seed_agricultura_familiar_silver.py
	$(MAKE) dbt-build-agricultura-familiar LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT)
	$(MAKE) analytics-init LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT) DUCKDB_PATH=$(DUCKDB_PATH)
	duckdb $(DUCKDB_PATH) -c "SELECT COUNT(*) AS entregas FROM analytics.conab_alimenta_brasil_entregas"
	duckdb $(DUCKDB_PATH) -c "SELECT COUNT(*) AS propostas FROM analytics.conab_alimenta_brasil_propostas"

ibge-localidades-mvp:
	go test ./internal/ibge/... ./internal/ingest/ -run 'IBGE|Localidades|Flatten|3550308'
	LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT) python3 scripts/ci/seed_ibge_localidades_silver.py
	$(MAKE) dbt-build-ibge-localidades LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT)
	$(MAKE) analytics-init LAKE_LOCAL_ROOT=$(LAKE_ABS) DUCKDB_PATH=$(DUCKDB_PATH)
	duckdb $(DUCKDB_PATH) -c "SELECT COUNT(*) AS municipios FROM analytics.ibge_localidades_municipios"
	duckdb $(DUCKDB_PATH) -c "SELECT codigo_ibge, nome, sigla_uf FROM analytics.ibge_localidades_municipios WHERE codigo_ibge = '3550308'"
	duckdb $(DUCKDB_PATH) -c "SELECT COUNT(*) AS ufs FROM analytics.ibge_localidades_ufs"

dbt-build-ibge-localidades: dbt-deps
	cd dbt && LAKE_LOCAL_ROOT=$(LAKE_ABS) dbt build --profiles-dir . --select 'stg_ibge__localidades_municipios+ stg_ibge__localidades_ufs+'

dbt-build-ibge-pam:
	cd dbt && LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT) dbt build --profiles-dir . --select 'stg_ibge__pam_area_quantidade+'

ibge-pam-mvp:
	go test ./internal/ibge/... ./internal/ingest/ -run 'PAM|SIDRA|pam|FlattenPAM'
	LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT) python3 scripts/ci/seed_ibge_pam_silver.py
	$(MAKE) dbt-build-ibge-pam LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT)

inmet-clima-mvp:
	go test ./internal/inmet/... ./internal/ingest/ -run 'INMET|FlattenEstacoes|BDMEP|Missing'

dbt-build-bcb-sgs:
	cd dbt && LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT) dbt build --profiles-dir . --select 'stg_bcb__sgs_ipca+ stg_bcb__sgs_ptax_usd_venda+'

bcb-sgs-mvp:
	go test ./internal/bcb/... ./internal/ingest/ -run 'BCB|SGS|FlattenSGS'
	LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT) python3 scripts/ci/seed_bcb_sgs_silver.py
	$(MAKE) dbt-build-bcb-sgs LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT)

cepea-indicadores-mvp:
	go test ./internal/cepea/... ./internal/ingest/ -run 'CEPA|Cepea|FlattenIndicador|ParseIndicator'
	LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT) python3 scripts/ci/seed_cepea_silver.py

conab-mvp:
	LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT) python3 scripts/ci/seed_dbt_silver.py
	$(MAKE) dbt-build LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT)
	$(MAKE) analytics-init LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT) DUCKDB_PATH=$(DUCKDB_PATH)
	$(MAKE) analytics-smoke DUCKDB_PATH=$(DUCKDB_PATH)
	duckdb $(DUCKDB_PATH) -c "SELECT COUNT(*) FROM analytics.conab_serie_historica_graos"

clean:
	rm -rf $(BIN_DIR)
