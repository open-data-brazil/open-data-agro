#!/usr/bin/env bash
# Create analytics.duckdb and apply versioned views from duckdb/views/*.sql
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
LAKE_LOCAL_ROOT="${LAKE_LOCAL_ROOT:-$ROOT/lake}"
LAKE_ABS="$(cd "$LAKE_LOCAL_ROOT" && pwd)"
DUCKDB_PATH="${DUCKDB_PATH:-$ROOT/duckdb/analytics.duckdb}"
DUCKDB_BIN="${DUCKDB_BIN:-duckdb}"
VIEWS_DIR="$ROOT/duckdb/views"

if ! command -v "$DUCKDB_BIN" >/dev/null 2>&1; then
  echo "duckdb CLI not found — run: make duckdb-install" >&2
  exit 1
fi

mkdir -p "$(dirname "$DUCKDB_PATH")" "$ROOT/duckdb/exports"
mkdir -p "$LAKE_ABS/gold/mart_conab__estimativa_graos"
mkdir -p "$LAKE_ABS/gold/mart_conab__serie_historica_graos"

"$DUCKDB_BIN" "$DUCKDB_PATH" -c "CREATE SCHEMA IF NOT EXISTS analytics;"

for view_file in "$VIEWS_DIR"/*.sql; do
  [[ -f "$view_file" ]] || continue
  sql="$(sed "s|__LAKE_ROOT__|${LAKE_ABS}|g" "$view_file")"
  printf '%s\n' "$sql" | "$DUCKDB_BIN" "$DUCKDB_PATH"
done

echo "analytics catalog ready: $DUCKDB_PATH (lake=$LAKE_ABS)"
