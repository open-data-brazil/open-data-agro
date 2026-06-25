.PHONY: test lint build build-processor clean duckdb-install

BIN_DIR := bin
DUCKDB_VERSION ?= 1.5.4

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

duckdb-install:
	curl -fsSL https://install.duckdb.org | DUCKDB_VERSION=$(DUCKDB_VERSION) sh

clean:
	rm -rf $(BIN_DIR)
