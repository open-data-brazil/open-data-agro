#!/usr/bin/env bash
# Download official CONAB Produção Agrícola samples into .local/reference/conab/ (gitignored).
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
OUT="${ROOT}/.local/reference/conab"
mkdir -p "${OUT}"

fetch() {
  local name="$1"
  local url="$2"
  echo "fetching ${name}..."
  curl -fsSL --max-time 60 "${url}" -o "${OUT}/${name}"
}

fetch "LevantamentoGraos.txt" "https://portaldeinformacoes.conab.gov.br/downloads/arquivos/LevantamentoGraos.txt"
fetch "SerieHistoricaGraos.txt" "https://portaldeinformacoes.conab.gov.br/downloads/arquivos/SerieHistoricaGraos.txt"
fetch "LevantamentoCana.txt" "https://portaldeinformacoes.conab.gov.br/downloads/arquivos/LevantamentoCana.txt"
fetch "SerieHistoricaCana.txt" "https://portaldeinformacoes.conab.gov.br/downloads/arquivos/SerieHistoricaCana.txt"

echo "saved samples under ${OUT}"
