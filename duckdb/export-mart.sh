#!/usr/bin/env bash
# Export an analytics view to portable Parquet + CONAB _metadata.json sidecar.
set -euo pipefail

ROOT="$(cd "$(dirname "$0")" && pwd)"
VIEW_NAME="${1:-conab_estimativa_graos}"
LAKE_LOCAL_ROOT="${LAKE_LOCAL_ROOT:-$ROOT/../lake}"
DUCKDB_PATH="${DUCKDB_PATH:-$ROOT/duckdb/open_data_agro.duckdb}"
DUCKDB_BIN="${DUCKDB_BIN:-duckdb}"
EXPORT_DIR="$ROOT/exports"
DATE_UTC="$(date -u +%Y-%m-%d)"
TIMESTAMP_UTC="$(date -u +%Y-%m-%dT%H:%M:%SZ)"
OUT_PARQUET="$EXPORT_DIR/${VIEW_NAME}-${DATE_UTC}.parquet"
OUT_META="$EXPORT_DIR/${VIEW_NAME}-${DATE_UTC}_metadata.json"

if ! command -v "$DUCKDB_BIN" >/dev/null 2>&1; then
  echo "duckdb CLI not found — run: make duckdb-install" >&2
  exit 1
fi

if [[ ! -f "$DUCKDB_PATH" ]]; then
  echo "analytics catalog missing — run: make analytics-init" >&2
  exit 1
fi

mkdir -p "$EXPORT_DIR"

"$DUCKDB_BIN" "$DUCKDB_PATH" -c \
  "COPY (SELECT * FROM analytics.${VIEW_NAME}) TO '${OUT_PARQUET}' (FORMAT PARQUET);"

cat >"$OUT_META" <<EOF
{
  "source": "CONAB — Portal de Informações Agropecuárias",
  "source_url": "https://portaldeinformacoes.conab.gov.br/download-arquivos.html",
  "agencia": "CONAB",
  "view": "analytics.${VIEW_NAME}",
  "exported_at": "${TIMESTAMP_UTC}",
  "format": "parquet",
  "lake_local_root": "${LAKE_LOCAL_ROOT}"
}
EOF

echo "exported ${OUT_PARQUET}"
echo "metadata ${OUT_META}"
