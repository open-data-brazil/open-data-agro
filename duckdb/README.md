# DuckDB — local analytics

DuckDB is used for ad-hoc processing (Phase 4) and analytics views (Phase 8).

**Local-first:** all default scripts read from `./lake/` on disk — no R2 or `httpfs` required. See [.local/LOCAL-FIRST.md](.local/LOCAL-FIRST.md).

## Path convention

| Context | Path |
|---------|------|
| Default dev file | `./duckdb/analytics.duckdb` (see `DUCKDB_PATH` in `.env.example`) |
| In-memory | `DUCKDB_PATH=:memory:` |
| Bronze reads | `read_parquet('lake/bronze/conab/{slug}/**/*.parquet')` |
| Silver Delta | `delta_scan('lake/silver/conab/{table}/')` |
| Exported snapshots | `./duckdb/exports/{dataset_id}-{YYYY-MM-DD}.parquet` |

## Environment

Set `DUCKDB_PATH` and `LAKE_LOCAL_ROOT` in `.env` (copy from `.env.example`). Phase 4 wires DuckDB via `cmd/processor`.

## Related lake paths

| Layer | Local folder | R2 prefix (production) |
|-------|--------------|------------------------|
| Bronze (Parquet) | `lake/bronze/` | `s3://{R2_BUCKET}/bronze/` |
| Silver (Delta) | `lake/silver/` | `s3://{R2_BUCKET}/silver/` |
| Gold (dbt marts) | `lake/gold/` | `s3://{R2_BUCKET}/gold/` |

## Quick start (local)

```bash
docker compose up -d postgres
./bin/ingestor run conab.estimativa-graos
duckdb -c "SELECT count(*) FROM read_parquet('lake/bronze/conab/estimativa-graos/**/*.parquet')"
```
