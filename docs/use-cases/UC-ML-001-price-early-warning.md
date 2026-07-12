# UC-ML-001 — Soybean port price early-warning features

## Summary

Build a **daily soja feature panel** anchored on CEPEA Paranaguá (`preco_rs_sc` / `preco_usd_sc`) and a **versioned walk-forward training export** so Phase 31 can train quantile early-warning models (horizons **7 / 30 / 90** trading days).

Phase 20 delivers features; Phase 30 freezes targets, splits, baselines, and export. No production trading API.

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

**Forward targets (Phase 30, frozen):** trading-day leads `y_ret_{h}d` / `y_level_{h}d` for `h ∈ {7,30,90}` — see [Phase 30 section](#phase-30--training-dataset-epic-b-frozen).

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
| `mart_cross__soja_features_monthly` | Monthly municipal CONAB grain panel |

### Monthly municipal panel

- **Grain:** `(produto_slug, cod_ibge, refmonth)`
- **Spine:** CONAB municipal soja prices (`preco_local_kg`)
- **`label_date`:** CONAB publication `as_of` (`month_end + 30d`)
- **Features:** month-aggregated CEPEA / PTAX, UF Comex, frete UF→PR — ASOF join with `as_of <= label_date`
- **DuckDB view:** `analytics.cross_soja_features_monthly`

---

## Phase 30 — Training dataset (Epic B, frozen)

> Folded former UC-002 into this use case. Spec source: `scripts/ml/config.py`.

### Target definition (B4)

Horizons `h ∈ {7, 30, 90}` are **trading-day leads** on the CEPEA spine (next `h` available observations — not calendar days).

| Column | Definition |
|--------|------------|
| `y_ret_{h}d` | `preco_rs_sc[t+h] / preco_rs_sc[t] - 1` |
| `y_level_{h}d` | `preco_rs_sc[t+h]` |

Primary early-warning target for Phase 31 gate: **`y_ret_30d`**.

### Feature list (B4, frozen)

| Feature | Theme |
|---------|-------|
| `preco_rs_sc`, `preco_usd_sc`, `variacao_dia_pct` | CEPEA level |
| `cepea_ret_1d`, `cepea_ret_7d`, `cepea_vol_7d` | CEPEA momentum / vol |
| `ptax_usd_venda` | FX |
| `b3_front_price`, `basis_cepea_usd_minus_b3` | Futures / basis |
| `comex_soja_fob_usd`, `comex_soja_kg` | Trade |
| `frete_mt_pr_ton_avg`, `frete_mt_pr_tkm_avg` | Logistics |
| `inmet_national_avg` | Climate proxy (nullable) |

Meta / audit only (not model inputs): `produto_slug`, `data`, `label_date`, `b3_front_symbol`, `as_of_*`, `feature_as_of_max`, `split`.

### Walk-forward split (B5)

**Never** random shuffle. Split on calendar `data`:

| Split | Rule |
|-------|------|
| `train` | `data <= 2018-12-31` |
| `val` | `2019-01-01` … `2021-12-31` |
| `test` | `data >= 2022-01-01` |

### Success gate (B3)

On the **test** split, for target `y_ret_30d`:

> Phase 31 model **MAE** MUST beat **seasonal_naive** MAE by **≥ 15%**  
> i.e. `model_MAE <= 0.85 * seasonal_naive_MAE`.

Baselines (last-value + seasonal naive) are produced by `make ml-baselines` → `.local/ml/baselines/`.

### Export (B6–B7)

```bash
make ml-dataset-export          # Parquet + manifest.json
make ml-baselines ml-corr       # baselines + lagged correlations
python3 scripts/ml/verify_manifest.py
```

Artifacts (local, gitignored under `.local/ml/`):

| Path | Content |
|------|---------|
| `.local/ml/export/soy_daily_training.parquet` | Features + targets + `split` |
| `.local/ml/export/manifest.json` | schema hash, row count, date range, git SHA |
| `.local/ml/baselines/` | JSON + Markdown baseline metrics |
| `.local/ml/corr/` | Lagged Pearson report |

---

## Phase 31 — Price early-warning models (Epic C, MVP)

Research tooling only (no production trading API). Code: `scripts/ml/`.

### Models (C1)

| Item | Spec |
|------|------|
| Point | LightGBM regression on `y_ret_{h}d`, `h ∈ {7,30,90}` |
| Quantiles | LightGBM quantile `p10` / `p50` / `p90` per horizon |
| Split | Same walk-forward calendar as Phase 30 (never shuffle) |
| Gate | Test MAE on `y_ret_30d` ≤ `0.85 * seasonal_naive MAE` (B3) |

```bash
make ml-train-soy          # train + SHAP/ablation + early-warning eval + registry
make ci-ml-train-soy       # seeded CI (requires toolchain/python-requirements-ml.txt)
```

Artifacts (local, gitignored under `.local/ml/`):

| Path | Content |
|------|---------|
| `.local/ml/train/soy_lgbm_metrics.json` | Point + quantile MAE / B3 pass-fail |
| `.local/ml/reports/soy_shap_ablation.md` | Mean \|SHAP\| + feature-group ablation |
| `.local/ml/reports/soy_early_warning_eval.md` | Spike rules precision/recall |
| `.local/ml/registry/soy_lgbm_v1.json` | Params, metrics, dataset manifest hash |

### Early-warning rules (C3)

Labels from train empirical percentiles of `y_ret_30d`; alerts from quantile forecasts on the test window (see report).

### Deferred

- **C4** TFT / Chronos — skipped while LightGBM clears B3 (optional follow-up)
- Epics **D** (pgvector) / **E** (multi-crop productization) — later

## Related

- [ADR 005 — Cross-source dbt analytics](../adr/005-cross-source-dbt-analytics.md)
- [ROADMAP.md](../ROADMAP.md) — Phase 20 / 30 / 31
- [REFRESH-POLICY.md](../REFRESH-POLICY.md) — data refresh ≠ ML retrain
- `.local/phases/30-ml-training-dataset/TASKS.md`
- `.local/phases/31-price-prediction-ia/TASKS.md`
