# DuckDB — local analytics

DuckDB is used for ad-hoc processing (Phase 4) and analytics views (Phase 8).

## Path convention

| Context | Path |
|---------|------|
| Default dev file | `./duckdb/analytics.duckdb` (see `DUCKDB_PATH` in `.env.example`) |
| In-memory | `DUCKDB_PATH=:memory:` |
| Exported snapshots | `./duckdb/exports/{dataset_id}-{YYYY-MM-DD}.duckdb` |

## Environment

Set `DUCKDB_PATH` in `.env` (copy from `.env.example`). The root package lists `duckdb` as a dev dependency for Node bindings in later phases.

## Related lake paths

| Layer | Local folder | R2 prefix |
|-------|--------------|-----------|
| Bronze (Parquet) | `lake/bronze/` | `s3://{R2_BUCKET}/bronze/` |
| Silver (Delta) | `lake/silver/` | `s3://{R2_BUCKET}/silver/` |
| Gold (dbt marts) | `lake/gold/` | `s3://{R2_BUCKET}/gold/` |
