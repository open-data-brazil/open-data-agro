# PostgreSQL (operational DB)

**Image:** `postgres:18.4-bookworm` (see [docker-compose.yml](../../docker-compose.yml))

Operational metadata lives in schema `catalog` — ingest jobs, file lineage, promotion audit.

## UUID policy — native UUIDv7

All primary keys and application-generated identifiers use **UUID version 7** (time-ordered).

| Layer | Mechanism |
|-------|-----------|
| **PostgreSQL 18.4** | `uuidv7()` — built-in default on `id` columns (`002_ingest_schema.sql`, `003_promotion_schema.sql`) |
| **Go ingestor** | `github.com/google/uuid` `NewV7()` — bronze Parquet part filenames (`part-{uuid}.parquet`) |

Do **not** use `gen_random_uuid()` (v4) or `uuid.New()` / `uuid.NewString()` (v4) for new IDs.

### Why UUIDv7

- Time-sortable keys improve index locality on `started_at` / `created_at` query patterns
- Native in PostgreSQL 18.4 — no extension required
- Same version in Go keeps lake object keys consistent with DB row IDs

### Fresh database

Init scripts under `init/` run automatically on first `docker compose up postgres`.

### Existing volume (upgraded from v4 defaults)

Apply migration `005_uuidv7_defaults.sql` (or recreate the volume):

```bash
docker compose exec -T postgres psql -U open_data_agro -d open_data_agro \
  < infra/postgres/init/005_uuidv7_defaults.sql
```

Existing rows keep their v4 IDs; only **new** inserts get `uuidv7()`.

### Verify

```sql
INSERT INTO catalog.promotion_jobs (dataset_id, status)
VALUES ('conab.estimativa-graos', 'running')
RETURNING id;
-- third group starts with 7, e.g. xxxxxxxx-xxxx-7xxx-xxxx-xxxxxxxxxxxx
```

## Init script order

| File | Purpose |
|------|---------|
| `001_placeholder.sql` | `catalog` schema |
| `002_ingest_schema.sql` | ingest jobs + files |
| `003_promotion_schema.sql` | promotion audit |
| `004_quality_failed_status.sql` | `quality_failed` status |
| `005_uuidv7_defaults.sql` | upgrade v4 → v7 defaults |
