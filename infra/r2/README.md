# Cloudflare R2 (production bronze)

Optional production backend when `STORAGE_MODE=r2`. Local development uses `STORAGE_MODE=local` by default; use MinIO (`STORAGE_MODE=minio`, `make ci-minio`) for S3-compatible integration tests without Cloudflare credentials.

Object keys match the local lake layout under `LAKE_LOCAL_ROOT` (path parity):

```text
bronze/{agency}/{dataset_slug}/ingest_date={YYYY-MM-DD}/part-{uuidv7}.parquet
bronze/{agency}/{dataset_slug}/ingest_date={YYYY-MM-DD}/_metadata.json
silver/...
gold/...
```

---

## Production deploy runbook

### 1. Prerequisites

- Cloudflare account with R2 enabled
- PostgreSQL and Redis for the ingestor (see root `docker-compose.yml` or managed equivalents)
- `DATABASE_URL` configured for the target environment
- Go toolchain matching `.local/STACK-VERSIONS.md` on ingestor/processor hosts

### 2. Provision the bucket

1. Open the [Cloudflare R2 dashboard](https://developers.cloudflare.com/r2/buckets/) and create a bucket (default name: `open-data-agro`).
2. Create an **R2 API token** scoped to **Object Read & Write** on that bucket.
3. Note the **account ID**, **access key ID**, and **secret access key**.

Optional: create placeholder prefixes `bronze/`, `silver/`, `gold/` (objects are created on first ingest).

### 3. Environment variables

Copy `.env.example` to `.env` on the deploy host and set:

| Variable | Required | Notes |
|----------|----------|-------|
| `STORAGE_MODE` | yes | `r2` |
| `R2_ACCOUNT_ID` | yes* | Used to derive endpoint when `R2_ENDPOINT` is unset |
| `R2_ACCESS_KEY_ID` | yes | API token access key |
| `R2_SECRET_ACCESS_KEY` | yes | API token secret |
| `R2_BUCKET` | yes | e.g. `open-data-agro` |
| `R2_ENDPOINT` | optional | `https://{account_id}.r2.cloudflarestorage.com` (auto from account ID) |
| `DATABASE_URL` | yes | Postgres catalog / dedup |
| `LAKE_LOCAL_ROOT` | optional | Still used for silver/gold Parquet + dbt locally; bronze goes to R2 |

\* Either `R2_ACCOUNT_ID` or an explicit `R2_ENDPOINT` with a real account id substituted for `{account_id}`.

Validation is enforced by `internal/config.LoadFromEnv()` when `STORAGE_MODE=r2`.

### 4. Deploy services

1. Build binaries: `make build` (ingestor + processor).
2. Run database migrations: `make migrate-up`.
3. Start ingestor with `.env` loaded (`STORAGE_MODE=r2`).
4. Run processor/GE/dbt on a host that can reach R2 (DuckDB uses `httpfs` + S3 secrets — see `internal/processor/duckdb.go`).

### 5. Verification

**Offline config check** (no Cloudflare network; CI uses fixture credentials):

```bash
make ci-validate-r2-env
```

**Production `.env` check** (on the deploy host with real secrets):

```bash
make validate-r2-env
```

**Live object-store smoke** (writes + deletes a test object under `bronze/integration-test/`):

```bash
R2_INTEGRATION=1 make validate-r2-env-live
```

**End-to-end ingest smoke** (after catalog seed):

```bash
# ingestor writes bronze to s3://{R2_BUCKET}/bronze/...
./bin/ingestor --dataset conab.estimativa-graos --from sample
```

Confirm objects in the R2 dashboard or:

```bash
aws s3 ls s3://open-data-agro/bronze/ --endpoint-url "https://{account_id}.r2.cloudflarestorage.com"
```

### 6. Rollback

1. Set `STORAGE_MODE=local` (or `minio` for staging).
2. Point `LAKE_LOCAL_ROOT` at a filesystem lake with existing bronze copies.
3. Redeploy ingestor/processor; bronze reads use local paths again.
4. R2 bucket can remain for archival; no automatic migration from R2 → local is provided in MVP.

### 7. Security notes

- Never commit `.env` or API tokens; `.env` is gitignored.
- Scope R2 tokens to the single bucket; rotate on deploy host compromise.
- CORS is **not** required for server-side ingestor/processor (enable only for browser clients).
- Great Expectations bronze checkpoints still expect local bronze paths unless synced — see `expectations/README.md`.

---

## S3 client settings

| Setting | Value |
|---------|-------|
| Region | `auto` |
| Path style | `false` (virtual-hosted) |
| SDK | AWS SDK v2 (`internal/storage/bronze.go`) |

## CORS (optional)

Enable CORS on the bucket only if browser clients need direct object access.

## Terraform

Infrastructure-as-code for R2 buckets is out of scope for the MVP; use the dashboard or add Terraform under `infra/r2/terraform/` in a later phase.

## Related

- `make validate-r2-env` / `make ci-validate-r2-env` — env validation
- `scripts/deploy/validate_r2_env.sh` — operator entrypoint
- `scripts/ci/check_r2_runbook.py` — CI gate for this runbook
- MinIO local substitute: [infra/minio/README.md](../minio/README.md)
