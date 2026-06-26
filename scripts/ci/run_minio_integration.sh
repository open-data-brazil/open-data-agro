#!/usr/bin/env bash
# Start Docker MinIO, seed bucket, run STORAGE_MODE=minio Go integration tests.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$ROOT"

if ! command -v docker >/dev/null 2>&1; then
  echo "run_minio_integration: docker required" >&2
  exit 1
fi

echo "Starting MinIO (docker compose)..."
docker compose up -d minio

echo "Waiting for MinIO health..."
for _ in $(seq 1 30); do
  if curl -sf http://localhost:9000/minio/health/live >/dev/null 2>&1; then
    break
  fi
  sleep 2
done
if ! curl -sf http://localhost:9000/minio/health/live >/dev/null 2>&1; then
  echo "run_minio_integration: MinIO did not become healthy" >&2
  exit 1
fi

docker compose run --rm minio-init

export MINIO_INTEGRATION=1
export MINIO_ENDPOINT="${MINIO_ENDPOINT:-http://localhost:9000}"
export MINIO_ACCESS_KEY="${MINIO_ACCESS_KEY:-minioadmin}"
export MINIO_SECRET_KEY="${MINIO_SECRET_KEY:-minioadmin}"
export MINIO_BUCKET="${MINIO_BUCKET:-open-data-agro}"

echo "MinIO bronze store (Put/List)..."
go test ./internal/storage -run TestMinIOListPrefixIntegration -count=1

PATH="${ROOT}/.local/bin:${PATH}"
export DUCKDB_BIN="${DUCKDB_BIN:-${ROOT}/.local/bin/duckdb}"

echo "MinIO DuckDB smoke (s3:// bronze URI)..."
go test ./internal/processor -run TestSmokeMinIOPathSwap -count=1

echo "run_minio_integration: PASS"
