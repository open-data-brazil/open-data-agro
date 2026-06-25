# PostgreSQL (operational DB)

**Image:** `postgres:18.4-bookworm` (see [docker-compose.yml](../../docker-compose.yml))

Operational metadata lives in schema `catalog` — dataset registry, ingest jobs, file lineage, promotion audit.

## Local connection

```bash
docker compose up -d postgres
export DATABASE_URL=postgresql://open_data_agro:open_data_agro@localhost:${POSTGRES_HOST_PORT:-5432}/open_data_agro?sslmode=disable
```

Override host port when 5432 is taken: `POSTGRES_HOST_PORT=5433` in `.env` (see [.env.example](../../.env.example)).

No managed/cloud Postgres is required for MVP.

## Migrations (golang-migrate)

Versioned SQL lives in [migrations/](migrations/). Fresh Docker volumes apply them automatically via [init/000_run_migrations.sh](init/000_run_migrations.sh).

```bash
make migrate-install   # once: install golang-migrate CLI
make migrate-up        # apply pending migrations
make migrate-down      # roll back last migration
```

| Version | Migration | Purpose |
|---------|-----------|---------|
| 000001 | `catalog_schema` | `catalog` schema |
| 000002 | `ingest_schema` | registry, jobs, files (`uuidv7()`, `source_portal_url`) |
| 000003 | `promotion_schema` | promotion audit (`quality_failed` status) |
| 000004 | `monitoring_views` | `v_latest_successful_ingest`, `v_failed_jobs_last_7d` |

## Seed registry

Idempotent upsert from [configs/catalog/conab/registry.yaml](../../configs/catalog/conab/registry.yaml):

```bash
make seed        # all CONAB datasets
make seed-mvp    # MVP only: estimativa-graos, serie-historica-graos
```

`source_portal_url` is set to https://portaldeinformacoes.conab.gov.br/download-arquivos.html for all `conab.*` datasets.

## UUID policy — native UUIDv7

All primary keys use **PostgreSQL 18.4** `uuidv7()` (time-ordered). See [internal/ingest/fingerprint.go](../../internal/ingest/fingerprint.go) for bronze part filenames.

## Monitoring views

| View | Use |
|------|-----|
| `catalog.v_latest_successful_ingest` | Latest successful job per dataset |
| `catalog.v_failed_jobs_last_7d` | Failed jobs in the last 7 days |

```bash
./bin/ingestor status --latest
./bin/ingestor status --failed
```

## Verify

```sql
SELECT * FROM catalog.ingest_jobs ORDER BY started_at DESC LIMIT 5;
SELECT * FROM catalog.v_latest_successful_ingest;
```
