# ADR 005 — Cross-source dbt analytics layer

**Status:** Accepted  
**Date:** 2026-07-12  
**Context:** Phase 20 — Epic A feature foundation for price early-warning

## Decision

Add a **cross-agency analytics layer** in dbt that reads **existing gold single-source marts** and writes:

| Prefix | Layer | Purpose |
|--------|-------|---------|
| `int_cross__*` | Intermediate | Typed casts, OD aggregates, local-price joins |
| `dim_municipio` | Intermediate dim | IBGE municipality keys for geo joins |
| `mart_ml__*` | Gold analytics | ML feature panels (e.g. `mart_ml__soy_daily_features`) |
| `mart_cross__*` | Gold analytics (optional) | Non-ML crossing marts (e.g. monthly municipal) |

### Hard rules

1. **No bronze changes** for crossing — do not alter GE bronze suites or ingest schemas.
2. **No mutation of single-source gold marts** (`mart_cepea__*`, `mart_conab__*`, …) — treat them as read-only inputs.
3. Cross models read gold via `read_gold_mart()` (Parquet path under `lake/gold/`), not by rebuilding agency staging unless a collection MVP is re-run.
4. Derived marts are **analytics**, not new official sources — do not add them to [OFFICIAL-SOURCES.md](../OFFICIAL-SOURCES.md) as agency datasets.
5. Every ML feature row carries `as_of` / `feature_as_of_max` so leakage tests can fail closed.

### Naming

```text
lake/gold/mart_ml__soy_daily_features/mart.parquet
lake/gold/mart_cross__soja_features_monthly/mart.parquet
lake/gold/dim_municipio/mart.parquet          # dimension (external)
lake/gold/int_cross__*                        # optional materialization; views OK
```

Intermediate `int_cross__*` default to **view**; `dim_municipio`, `mart_ml__*`, and `mart_cross__*` materialize as **external Parquet** under `gold_root`.

## Alternatives considered

| Option | Outcome |
|--------|---------|
| Join inside DuckDB views only (no dbt) | Weak lineage/tests; rejected |
| Mutate staging to typed columns in place | Breaks varchar gold contract + GE; rejected |
| Spark / separate warehouse | Violates local-first; deferred |
| `ref()` agency marts and rebuild upstream | Couples analytics CI to full silver; rejected for MVP gate |

## Consequences

- Phase 20 CI seeds **gold fixtures** (or reuses a local lake) then `dbt build --select` cross/ml models only.
- Phase 30 exports from `mart_ml__soy_daily_features` without touching ingest.
- CEPEA CC BY-NC still applies to any panel that includes CEPEA columns — see [UC-ML-001](../use-cases/UC-ML-001-price-early-warning.md).

## References

- [UC-ML-001](../use-cases/UC-ML-001-price-early-warning.md)
- [DATA-CROSSING-VISION](../../.local/DATA-CROSSING-VISION.md) (local)
- [ADR 004](004-unified-postgresql-sync.md) — unified DB remains TEXT; typed casts live in analytics dbt
