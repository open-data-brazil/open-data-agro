.PHONY: test lint build build-processor clean duckdb-install python-install dbt-deps dbt-build dbt-build-mercado dbt-build-abastecimento dbt-build-armazenamento dbt-build-agricultura-familiar dbt-build-ibge-pam ibge-localidades-mvp ibge-pam-mvp migrate-install migrate-up migrate-down seed analytics-init analytics-smoke conab-reference conab-mvp conab-mercado-mvp conab-abastecimento-mvp conab-armazenamento-mvp conab-agricultura-familiar-mvp

BIN_DIR := bin
DUCKDB_VERSION ?= 1.5.4
LAKE_LOCAL_ROOT ?= ./lake
DUCKDB_PATH ?= ./duckdb/open_data_agro.duckdb
POSTGRES_HOST_PORT ?= 5432
DATABASE_URL ?= postgresql://open_data_agro:open_data_agro@localhost:$(POSTGRES_HOST_PORT)/open_data_agro?sslmode=disable
MIGRATIONS_PATH := infra/postgres/migrations
MIGRATE ?= migrate

test:
	go test ./...

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
	LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT) DUCKDB_PATH=$(DUCKDB_PATH) ./duckdb/scripts/analytics-init.sh

analytics-smoke:
	@command -v duckdb >/dev/null 2>&1 || (echo "run: make duckdb-install" && exit 1)
	duckdb $(DUCKDB_PATH) -c "SELECT COUNT(*) AS rows FROM analytics.conab_estimativa_graos"

conab-reference:
	@chmod +x scripts/conab/fetch_reference_samples.sh
	./scripts/conab/fetch_reference_samples.sh

dbt-build-mercado: dbt-deps
	cd dbt && LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT) dbt build --profiles-dir . --select 'stg_conab__oferta_demanda+'

conab-mercado-mvp:
	LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT) python3 scripts/ci/seed_mercado_silver.py
	$(MAKE) dbt-build-mercado LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT)
	$(MAKE) analytics-init LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT) DUCKDB_PATH=$(DUCKDB_PATH)
	duckdb $(DUCKDB_PATH) -c "SELECT COUNT(*) AS rows FROM analytics.conab_oferta_demanda"

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

dbt-build-ibge-pam:
	cd dbt && LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT) dbt build --profiles-dir . --select 'stg_ibge__pam_area_quantidade+'

ibge-pam-mvp:
	go test ./internal/ibge/... ./internal/ingest/ -run 'PAM|SIDRA|pam|FlattenPAM'
	LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT) python3 scripts/ci/seed_ibge_pam_silver.py
	$(MAKE) dbt-build-ibge-pam LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT)

conab-mvp:
	LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT) python3 scripts/ci/seed_dbt_silver.py
	$(MAKE) dbt-build LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT)
	$(MAKE) analytics-init LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT) DUCKDB_PATH=$(DUCKDB_PATH)
	$(MAKE) analytics-smoke DUCKDB_PATH=$(DUCKDB_PATH)
	duckdb $(DUCKDB_PATH) -c "SELECT COUNT(*) FROM analytics.conab_serie_historica_graos"

clean:
	rm -rf $(BIN_DIR)
