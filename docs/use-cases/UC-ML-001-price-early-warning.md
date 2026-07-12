# UC-ML-001 — Soybean port price early-warning features

## Summary

Build a **daily soja feature panel** anchored on CEPEA Paranaguá (`preco_rs_sc` / `preco_usd_sc`) so Phase 30 can export walk-forward training sets and Phase 31 can train quantile early-warning models (horizons **7 / 30 / 90** days).

This use case covers **feature foundation only** (Phase 20). No model training.

## Actors

- Data engineer (dbt cross-source layer)
- ML engineer (Phase 30–31 consumer)

## Preconditions

- Gold marts exist for CEPEA soja Paranaguá, BCB PTAX, B3 SOY futures, MDIC Comex soja, CONAB frete / preços, IBGE localidades
- Cross-source dbt layer follows [ADR 005](../adr/005-cross-source-dbt-analytics.md)
- CEPEA CC BY-NC license accepted for research / non-commercial use (see License)

## Label (MVP)

| Field | Source | Notes |
|-------|--------|-------|
| `produto_slug` | `soja` | Canonical crop key |
| `label_date` | CEPEA `data` | Trading day of port reference |
| `preco_rs_sc` | CEPEA Paranaguá | R$ / sc 60 kg |
| `preco_usd_sc` | CEPEA or PTAX-derived | Prefer CEPEA when present; else `preco_rs_sc / ptax` |

**Forward targets (Phase 30):** return or level at `label_date + {7,30,90}` calendar days (trading-day alignment documented at export time).

## Join keys

| Key | Sources | Policy |
|-----|---------|--------|
| `produto_slug` | CEPEA `produto`, MDIC `produto_slug`, CONAB product filter | Lowercase slug; CONAB `SOJA` → `soja` |
| `data` / `refdate` | Daily series (CEPEA, PTAX, B3, INMET) | ISO date |
| `refmonth` | Monthly series (Comex, frete, CONAB municipal) | First day of month `YYYY-MM-01` |
| `uf` / `cod_ibge` | CONAB, IBGE `dim_municipio`, frete OD | Seven-digit IBGE; UF sigla |

## Publication lags (as_of policy)

Features MUST use only information with `as_of <= label_date`. Conservative publication lags:

| Source | Native grain | Assumed publish lag | `as_of` rule |
|--------|--------------|---------------------|--------------|
| CEPEA indicadores | Daily | 0 days (same session) | `as_of = data` |
| BCB PTAX | Daily | 0–1 business day | `as_of = data` (last ≤ label) |
| B3 futures | Daily | 0 days | `as_of = refdate` |
| MDIC Comex | Monthly | **45 days** after month end | `as_of = month_end + 45d` |
| CONAB frete | Monthly | **14 days** after month end | `as_of = month_end + 14d` |
| CONAB preços municipal | Monthly | **30 days** after month end | `as_of = month_end + 30d` |
| CONAB preços semanal UF | Weekly | **7 days** after week end | `as_of = week_end + 7d` |
| INMET BDMEP | Daily | **2 days** | `as_of = data + 2d` |
| USDA WASDE | Monthly | Release calendar (~mid-month) | Document at join; not in MVP panel |
| NOAA ENSO / sequia | Monthly | **10 days** after month end | `as_of = month_end + 10d` |

## Resampling policy (daily vs monthly)

See [GLOSSARY — Analytics resampling](../GLOSSARY.md#analytics-resampling-daily-vs-monthly).

- **Spine:** CEPEA daily trading days.
- **Monthly → daily:** forward-fill last known monthly feature whose `as_of <= label_date` (no interpolation of levels).
- **Weekly → daily:** same as-of forward-fill.
- **Never** use future-published months on the daily spine.

## Main flow

1. **GIVEN** gold CEPEA soja Paranaguá and supporting marts
2. **WHEN** `make analytics-crossing-mvp` builds `mart_ml__soy_daily_features`
3. **THEN** each row has grain `(produto_slug, data)`, typed numerics, and `feature_as_of_max <= label_date`
4. **AND** dbt leakage test passes
5. **AND** DuckDB view `analytics.ml_soy_daily_features` is queryable

## Alternate flows

| Case | Expected |
|------|----------|
| Missing optional climate gold (INMET empty) | Climate columns NULL; panel still builds |
| Missing B3 match for a day | `b3_front_price` NULL; basis NULL |
| CEPEA `preco_usd_sc` blank | Derive from PTAX when available |

## Municipality ↔ INMET station (optional later)

Spatial join options (not implemented in Phase 20 MVP):

1. **Nearest station** — haversine from `dim_municipio` centroid (requires lat/lon enrichment) to INMET station catalog
2. **UF aggregate** — mean/median of stations in the same UF (MVP-friendly; no geometry)
3. **IBGE micro/meso** — assign stations to IBGE regions via point-in-polygon (deferred)

Document choice in a follow-up ADR before municipal climate features ship.

## License — CEPEA CC BY-NC

CEPEA/ESALQ indicators are published under [CC BY-NC 4.0](https://www.cepea.org.br/br/licenca-de-uso-de-dados.aspx).

| Use | Allowed? |
|-----|----------|
| Research / open early-warning experiments in this repo | Yes, with attribution |
| Commercial SaaS / paid trading signals using CEPEA labels or features | **No** without a separate CEPEA commercial license |
| Redistribution of raw CEPEA series as a product | Follow CC BY-NC — non-commercial only |

Derived marts under `mart_ml__*` are **analytics artifacts**, not new official sources. See [OFFICIAL-SOURCES.md](../OFFICIAL-SOURCES.md) CEPEA section.

## Deliverable table

| Model | Role |
|-------|------|
| `dim_municipio` | IBGE municipality dimension (SCD Type 1) |
| `int_cross__precos_locais_municipio` | Typed CONAB municipal prices + dim |
| `int_cross__frete_uf_par` | UF-pair frete monthly aggregates |
| `mart_ml__soy_daily_features` | Daily soja MVP feature panel |

## Related

- [ADR 005 — Cross-source dbt analytics](../adr/005-cross-source-dbt-analytics.md)
- [ROADMAP.md](../ROADMAP.md) — Phase 20
- `.local/phases/20-analytics-crossing/TASKS.md`
