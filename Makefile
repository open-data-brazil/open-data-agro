.PHONY: test lint build build-processor clean duckdb-install python-install dbt-deps dbt-build dbt-build-mercado migrate-install migrate-up migrate-down seed analytics-init analytics-smoke conab-reference conab-mvp conab-mercado-mvp

BIN_DIR := bin
DUCKDB_VERSION ?= 1.5.4
LAKE_LOCAL_ROOT ?= ./lake
DUCKDB_PATH ?= ./duckdb/analytics.duckdb
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

conab-mvp:
	LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT) python3 scripts/ci/seed_dbt_silver.py
	$(MAKE) dbt-build LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT)
	$(MAKE) analytics-init LAKE_LOCAL_ROOT=$(LAKE_LOCAL_ROOT) DUCKDB_PATH=$(DUCKDB_PATH)
	$(MAKE) analytics-smoke DUCKDB_PATH=$(DUCKDB_PATH)
	duckdb $(DUCKDB_PATH) -c "SELECT COUNT(*) FROM analytics.conab_serie_historica_graos"

clean:
	rm -rf $(BIN_DIR)
