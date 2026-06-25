.PHONY: test lint build build-processor clean

BIN_DIR := bin

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

clean:
	rm -rf $(BIN_DIR)
