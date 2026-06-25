#!/usr/bin/env bash
# Fresh-database bootstrap: apply versioned migrations via golang-migrate.
set -euo pipefail

MIGRATE_VERSION="${MIGRATE_VERSION:-v4.18.2}"
MIGRATE_BIN="/tmp/migrate"
MIGRATIONS_PATH="${MIGRATIONS_PATH:-/migrations}"

if [[ ! -x "${MIGRATE_BIN}" ]]; then
  arch="$(uname -m)"
  case "${arch}" in
    x86_64) migrate_arch="amd64" ;;
    aarch64 | arm64) migrate_arch="arm64" ;;
    *)
      echo "unsupported architecture for migrate: ${arch}" >&2
      exit 1
      ;;
  esac
  url="https://github.com/golang-migrate/migrate/releases/download/${MIGRATE_VERSION}/migrate.linux-${migrate_arch}.tar.gz"
  if command -v curl >/dev/null 2>&1; then
    curl -fsSL "${url}" | tar xz -C /tmp migrate
  elif command -v wget >/dev/null 2>&1; then
    wget -qO- "${url}" | tar xz -C /tmp migrate
  else
    echo "curl or wget required to bootstrap golang-migrate" >&2
    exit 1
  fi
  chmod +x "${MIGRATE_BIN}"
fi

database_url="postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:5432/${POSTGRES_DB}?sslmode=disable"
"${MIGRATE_BIN}" -path "${MIGRATIONS_PATH}" -database "${database_url}" up
