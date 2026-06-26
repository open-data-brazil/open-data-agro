# MinIO (local S3 API)

Optional substitute for Cloudflare R2 during development and integration tests.

## Start

```bash
docker compose up -d minio minio-init
```

Set in `.env`:

```text
STORAGE_MODE=minio
MINIO_ENDPOINT=http://localhost:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin
MINIO_BUCKET=open-data-agro
```

Console: http://localhost:9001

## Integration test

```bash
make ci-minio
```

Or manually:

```bash
docker compose up -d minio minio-init
MINIO_INTEGRATION=1 go test ./internal/storage -run TestMinIOListPrefixIntegration -count=1
MINIO_INTEGRATION=1 go test ./internal/processor -run TestSmokeMinIOPathSwap -count=1
```

## Path parity

Object keys match local filesystem layout under `LAKE_LOCAL_ROOT`:

```text
bronze/conab/{dataset_slug}/ingest_date={YYYY-MM-DD}/part-{uuidv7}.parquet
bronze/conab/{dataset_slug}/ingest_date={YYYY-MM-DD}/_metadata.json
```
