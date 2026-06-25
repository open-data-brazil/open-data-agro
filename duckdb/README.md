# DuckDB — local analytics

DuckDB is the **offline analytic layer** over local gold marts and silver Parquet (Phase 8). Processing scripts for bronze/silver smoke tests live under `scripts/`.

**Local-first:** default views read `./lake/gold/` and `./lake/silver/` on disk — no R2 or `httpfs` required. See [.local/LOCAL-FIRST.md](../.local/LOCAL-FIRST.md).

## Path convention

| Context | Path |
|---------|------|
| Analytic catalog | `./duckdb/analytics.duckdb` (`DUCKDB_PATH` in `.env.example`) |
| In-memory | `DUCKDB_PATH=:memory:` |
| Published views | `duckdb/views/*.sql` → schema `analytics` |
| Portable exports | `duckdb/exports/{view}-{YYYY-MM-DD}.parquet` |
| Gold mart (dbt) | `lake/gold/mart_conab__estimativa_graos/mart.parquet` |
| Silver (serie) | `lake/silver/conab/serie_historica_graos/` |

## Quick start (full local pipeline)

```bash
make duckdb-install
docker compose up -d postgres
make seed-mvp
./bin/ingestor run conab.estimativa-graos    # network: CONAB download
./bin/processor promote --dataset conab.estimativa-graos
make dbt-build                              # writes gold mart
make analytics-init                         # creates analytics.duckdb + views

duckdb duckdb/analytics.duckdb -c "SELECT * FROM analytics.conab_estimativa_graos LIMIT 10"
```

CI shortcut (no CONAB download):

```bash
LAKE_LOCAL_ROOT=/tmp/open-data-agro-lake python3 scripts/ci/seed_dbt_silver.py
make dbt-build LAKE_LOCAL_ROOT=/tmp/open-data-agro-lake
make analytics-init LAKE_LOCAL_ROOT=/tmp/open-data-agro-lake DUCKDB_PATH=/tmp/analytics.duckdb
```

## Make targets

| Target | Purpose |
|--------|---------|
| `make analytics-init` | Create `analytics.duckdb` and apply views |
| `make analytics-smoke` | `SELECT COUNT(*)` on `analytics.conab_estimativa_graos` |
| `make duckdb-install` | Install DuckDB CLI 1.5.4 |

## Published views

| View | Source |
|------|--------|
| `analytics.conab_estimativa_graos` | `lake/gold/mart_conab__estimativa_graos/mart.parquet` |
| `analytics.conab_serie_historica_graos` | `lake/gold/mart_conab__serie_historica_graos/mart.parquet` |
| `analytics.conab_oferta_demanda` | `lake/gold/mart_conab__oferta_demanda/mart.parquet` |

SQL definitions: [views/](views/).

## Example queries (local-only)

**Total estimated production by crop and UF (latest season):**

```bash
duckdb duckdb/analytics.duckdb < duckdb/analyses/production_by_crop_uf.sql
```

**YoY change from historical series:**

```bash
duckdb duckdb/analytics.duckdb < duckdb/analyses/serie_historica_yoy.sql
```

Filter `produto` / `uf` / `safra` in `WHERE` to prune scans — gold mart is a single file; serie historica uses folder glob.

## Export workflow

```bash
./duckdb/export-mart.sh conab_estimativa_graos
```

Writes Parquet + `_metadata.json` with CONAB citation to [exports/](exports/). No cloud credentials required.

## Performance (laptop hardware)

| Tip | Detail |
|-----|--------|
| **Partition pruning** | Filter `safra` (estimativa) or `ano` + `produto` (serie) in `WHERE` before aggregations |
| **Memory limit** | `SET memory_limit = '4GB';` before large scans if DuckDB warns |
| **Threads** | `SET threads = 4;` on low-RAM machines |
| **Gold marts** | Single Parquet file per mart — fast for MVP volumes |

## Optional R2 attach (production parity)

Not used in the default profile. For object-store reads:

```sql
INSTALL httpfs; LOAD httpfs;
SET s3_endpoint = 'https://{account_id}.r2.cloudflarestorage.com';
-- set s3_access_key_id / s3_secret_access_key via env
SELECT * FROM read_parquet('s3://{bucket}/gold/mart_conab__estimativa_graos/mart.parquet');
```

Use `STORAGE_MODE=local` for development; keep `httpfs` out of `analytics-init`.

## Processing scripts (Phase 4)

| Script | Purpose |
|--------|---------|
| `scripts/smoke_read_parquet.sql` | Bronze row count (`processor smoke`) |
| `scripts/promote_bronze_to_silver.sql` | Promotion preview |
| `scripts/setup_delta.sql` | Optional `delta` extension for silver reads |

## Related lake paths

| Layer | Local folder | R2 prefix (production) |
|-------|--------------|------------------------|
| Bronze (Parquet) | `lake/bronze/` | `s3://{R2_BUCKET}/bronze/` |
| Silver (Delta) | `lake/silver/` | `s3://{R2_BUCKET}/silver/` |
| Gold (dbt marts) | `lake/gold/` | `s3://{R2_BUCKET}/gold/` |
