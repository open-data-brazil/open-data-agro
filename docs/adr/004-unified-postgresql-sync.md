# ADR 004 — Unified PostgreSQL sync (Stage G)

**Status:** Accepted  
**Date:** 2026-06-26  
**Context:** Phase 29 — sync gold dbt marts → PostgreSQL `analytics` schema

## Decision

`cmd/processor sync-postgres` performs a **full refresh** per mart:

1. Discover `lake/gold/mart_*/mart.parquet` files.
2. Export each mart to CSV via **DuckDB** `read_parquet` (already the project’s analytics runtime).
3. Load into PostgreSQL with **`COPY FROM STDIN`** through **pgx** (bulk, no row-by-row INSERT).
4. Store all mart columns as **`TEXT`** in Postgres — preserves full history without type coercion bugs across heterogeneous agencies.
5. Record a **manifest** in `analytics.sync_runs` + `analytics.sync_tables` (row counts, min/max date hints, parquet size).

Table names mirror DuckDB views: `mart_conab__estimativa_graos` → `analytics.conab_estimativa_graos`.

## Alternatives considered

| Option | Outcome |
|--------|---------|
| **dbt post-hook → Postgres** | Couples sync to dbt runs; harder to re-sync without rebuild; rejected for MVP |
| **INSERT SELECT from foreign table** | Requires `postgres_fdw` + file staging; more ops burden |
| **Typed schema per mart** | High maintenance (60+ marts); deferred — TEXT + indexes on join keys is sufficient for unified query layer |
| **Incremental by Delta version** | Silver is Delta; gold is external Parquet snapshot — full refresh of gold is cheap enough for MVP |

## Retention

Full refresh **replaces** table contents but never truncates history inside the source gold Parquet. No silent year filtering.

## Indexes

After each load, create indexes when join-key columns exist (see `docs/POSTGRES-UNIFIED-SYNC.md`):

- `(cod_ibge)`, `(codigo_ibge)`
- `(produto, safra)` — composite where both columns exist
- `(refmonth)`, `(data_preco)`, `(capturado_em)` — time keys

## Consequences

- Requires PostgreSQL 18.4 + migration `000005_analytics_schema`.
- Requires DuckDB CLI for CSV export (`make duckdb-install`).
- DuckDB remains the ad-hoc local analytics catalog; PostgreSQL is the **primary unified SQL store** for BI and future crossing.

## References

- [Phase 29 README](../../.local/phases/29-unified-postgresql/README.md)
- [DATA-CROSSING-VISION join keys](../../.local/DATA-CROSSING-VISION.md)
