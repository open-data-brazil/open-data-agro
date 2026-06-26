#!/usr/bin/env bash
# Create analytics.duckdb and apply versioned views from duckdb/views/*.sql
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
LAKE_LOCAL_ROOT="${LAKE_LOCAL_ROOT:-$ROOT/lake}"
LAKE_ABS="$(cd "$LAKE_LOCAL_ROOT" && pwd)"
DUCKDB_PATH="${DUCKDB_PATH:-$ROOT/duckdb/open_data_agro.duckdb}"
DUCKDB_BIN="${DUCKDB_BIN:-duckdb}"
VIEWS_DIR="$ROOT/duckdb/views"

if ! command -v "$DUCKDB_BIN" >/dev/null 2>&1; then
  echo "duckdb CLI not found — run: make duckdb-install" >&2
  exit 1
fi

mkdir -p "$(dirname "$DUCKDB_PATH")" "$ROOT/duckdb/exports"
mkdir -p "$LAKE_ABS/gold/mart_conab__estimativa_graos"
mkdir -p "$LAKE_ABS/gold/mart_conab__serie_historica_graos"
mkdir -p "$LAKE_ABS/gold/mart_conab__oferta_demanda"
mkdir -p "$LAKE_ABS/gold/mart_conab__precos_semanal_uf"
mkdir -p "$LAKE_ABS/gold/mart_conab__precos_semanal_municipio"
mkdir -p "$LAKE_ABS/gold/mart_conab__precos_mensal_uf"
mkdir -p "$LAKE_ABS/gold/mart_conab__precos_mensal_municipio"
mkdir -p "$LAKE_ABS/gold/mart_conab__estoques_publicos"
mkdir -p "$LAKE_ABS/gold/mart_anp__combustiveis_precos_medios_municipios"
mkdir -p "$LAKE_ABS/gold/mart_anp__combustiveis_precos_postos"
mkdir -p "$LAKE_ABS/gold/mart_conab__armazenagem"
mkdir -p "$LAKE_ABS/gold/mart_conab__alimenta_brasil_entregas"
mkdir -p "$LAKE_ABS/gold/mart_conab__alimenta_brasil_propostas"
mkdir -p "$LAKE_ABS/gold/mart_ibge__localidades_municipios"
mkdir -p "$LAKE_ABS/gold/mart_ibge__localidades_ufs"

"$DUCKDB_BIN" "$DUCKDB_PATH" -c "CREATE SCHEMA IF NOT EXISTS analytics;"

for view_file in "$VIEWS_DIR"/*.sql; do
  [[ -f "$view_file" ]] || continue
  sql="$(sed "s|__LAKE_ROOT__|${LAKE_ABS}|g" "$view_file")"
  parquet_path="$(printf '%s\n' "$sql" | sed -n "s/.*read_parquet('\([^']*\)').*/\1/p" | head -1)"
  if [[ -n "$parquet_path" && ! -f "$parquet_path" ]]; then
    echo "skip $(basename "$view_file") — missing $parquet_path"
    continue
  fi
  printf '%s\n' "$sql" | "$DUCKDB_BIN" "$DUCKDB_PATH"
done

echo "analytics catalog ready: $DUCKDB_PATH (lake=$LAKE_ABS)"
