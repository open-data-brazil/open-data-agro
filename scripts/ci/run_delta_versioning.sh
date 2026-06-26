#!/usr/bin/env bash
# Run native Delta Lake silver versioning integration tests (promote append + time travel).
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$ROOT"

PATH="${ROOT}/.local/bin:${PATH}"
export DUCKDB_BIN="${DUCKDB_BIN:-${ROOT}/.local/bin/duckdb}"
export DELTA_INTEGRATION=1

echo "Delta promote + append versioning..."
go test ./internal/processor -run 'TestPromoteLocalIntegration|TestPromoteAppendSecondVersion' -count=1

echo "DuckDB delta_scan smoke..."
go test ./internal/processor -run TestDuckDBDeltaScanSmoke -count=1

echo "run_delta_versioning: PASS"
