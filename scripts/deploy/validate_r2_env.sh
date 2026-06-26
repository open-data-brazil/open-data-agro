#!/usr/bin/env bash
# Validate STORAGE_MODE=r2 environment (offline fixture or live .env).
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$ROOT"

if [[ "${VALIDATE_R2_FIXTURE:-}" == "1" ]]; then
  export DATABASE_URL="postgresql://open_data_agro:open_data_agro@localhost:5432/open_data_agro?sslmode=disable"
  export STORAGE_MODE=r2
  export R2_ACCOUNT_ID=ci-fixture-account
  export R2_ACCESS_KEY_ID=ci-fixture-key
  export R2_SECRET_ACCESS_KEY=ci-fixture-secret
  export R2_BUCKET=open-data-agro
  unset R2_ENDPOINT
elif [[ -f .env ]]; then
  set -a
  # shellcheck disable=SC1091
  source .env
  set +a
fi

mode="${STORAGE_MODE:-local}"
if [[ "$mode" != "r2" ]]; then
  echo "validate_r2_env: STORAGE_MODE is ${mode} (not r2); set STORAGE_MODE=r2 in .env to validate production credentials"
  exit 0
fi

export VALIDATE_R2_ENV=1
echo "validate_r2_env: checking R2 config via internal/config..."
go test ./internal/config -run TestValidateR2EnvLive -count=1

if [[ "${R2_INTEGRATION:-}" == "1" ]]; then
  echo "validate_r2_env: live R2 Put/List integration..."
  go test ./internal/storage -run TestR2ListPrefixIntegration -count=1
fi

echo "validate_r2_env: PASS"
