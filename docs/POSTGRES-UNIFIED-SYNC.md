# Unified PostgreSQL sync (Stage G)

Gold dbt marts (`lake/gold/mart_*/mart.parquet`) are mirrored into PostgreSQL schema **`analytics`** for standard SQL access across all agencies and years.

**Decision record:** [adr/004-unified-postgresql-sync.md](adr/004-unified-postgresql-sync.md)

---

## Prerequisites

```bash
docker compose up -d postgres
export DATABASE_URL=postgresql://open_data_agro:open_data_agro@localhost:${POSTGRES_HOST_PORT:-5432}/open_data_agro?sslmode=disable
make migrate-up
make duckdb-install
```

Gold marts must exist (run phase MVPs or `make dbt-build` / collection targets first).

---

## Sync all marts

```bash
make unified-db-sync
```

This runs `migrate-up`, `processor sync-postgres`, and `scripts/ci/verify_unified_db_sync.py` (row-count parity).

Subset sync:

```bash
UNIFIED_DB_SYNC_MARTS=conab_estimativa_graos,conab_serie_historica_graos make unified-db-sync
```

---

## Table naming

| Gold path | PostgreSQL table | DuckDB view |
|-----------|-------------------|-------------|
| `gold/mart_conab__estimativa_graos/mart.parquet` | `analytics.conab_estimativa_graos` | `analytics.conab_estimativa_graos` |
| `gold/mart_ibge__localidades_municipios/mart.parquet` | `analytics.ibge_localidades_municipios` | `analytics.ibge_localidades_municipios` |

Rule: strip `mart_` prefix, replace `__` with `_`.

---

## Manifest

Each sync run writes:

| Object | Purpose |
|--------|---------|
| `analytics.sync_runs` | Run status, lake root, table count |
| `analytics.sync_tables` | Per-table row count, date range hints, gold path |
| `analytics.v_latest_sync_tables` | Latest successful sync per table |

```sql
SELECT table_name, row_count, min_date, max_date, synced_at
FROM analytics.v_latest_sync_tables
ORDER BY table_name;
```

---

## Join-key indexes

Created automatically when columns exist (aligned with [DATA-CROSSING-VISION](../.local/DATA-CROSSING-VISION.md)):

| Index suffix | Columns | Use |
|--------------|---------|-----|
| `_cod_ibge_idx` | `cod_ibge` | Municipal CONAB prices, frete, custo |
| `_codigo_ibge_idx` | `codigo_ibge` | IBGE dimension joins |
| `_produto_safra_idx` | `produto`, `safra` | Production / price season alignment |
| `_refmonth_idx` | `refmonth` | Monthly macro / commodity series |
| `_data_preco_idx` | `data_preco` | Daily / weekly price grain |
| `_capturado_em_idx` | `capturado_em` | Ingest lineage time |

All mart data columns are stored as **TEXT** — cast in queries as needed.

---

## CI smoke

```bash
make ci-unified-db-sync
```

Uses isolated `/tmp` lake, seeds gold subset, starts Postgres via Docker Compose, syncs, verifies parity.

---

## Related

- [infra/postgres/README.md](../infra/postgres/README.md) — operational DB + migrations
- [ROADMAP.md](ROADMAP.md) — Phase 29 status
