#!/bin/sh
set -eu

MC_ALIAS=local
MC_HOST="http://minio:9000"
BUCKET="${MINIO_BUCKET:-open-data-agro}"

until mc alias set "$MC_ALIAS" "$MC_HOST" "$MINIO_ROOT_USER" "$MINIO_ROOT_PASSWORD"; do
  echo "waiting for MinIO..."
  sleep 2
done

mc mb --ignore-existing "$MC_ALIAS/$BUCKET"

# Prefix placeholders (S3 keys; objects created on first ingest)
for prefix in bronze silver gold; do
  echo "$prefix" | mc pipe "$MC_ALIAS/$BUCKET/$prefix/.keep" >/dev/null 2>&1 || true
done

echo "MinIO bucket $BUCKET ready (bronze/, silver/, gold/)"
