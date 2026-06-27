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

### Wave 3 marts (Phases 44–45)

Discovered automatically by `processor sync-postgres` when gold parquets exist under `lake/gold/`:

| Gold path | PostgreSQL / DuckDB table | Source |
|-----------|---------------------------|--------|
| `gold/mart_dnit__snv_rodovias_federais/mart.parquet` | `analytics.dnit_snv_rodovias_federais` | DNIT SNV federal highways |
| `gold/mart_ipea__series_macro_regionais/mart.parquet` | `analytics.ipea_series_macro_regionais` | IPEA regional macro series |
| `gold/mart_ibge__pevs_producao_vegetal/mart.parquet` | `analytics.ibge_pevs_producao_vegetal` | IBGE PEVS vegetable production |
| `gold/mart_ibge__ppm_producao_municipal/mart.parquet` | `analytics.ibge_ppm_producao_municipal` | IBGE PPM municipal production |
| `gold/mart_aneel__tarifas_energia/mart.parquet` | `analytics.aneel_tarifas_energia` | ANEEL energy tariffs |
| `gold/mart_bndes__financiamento_agro/mart.parquet` | `analytics.bndes_financiamento_agro` | BNDES agro financing |
| `gold/mart_inmet__sequia_monitor/mart.parquet` | `analytics.inmet_sequia_monitor` | INMET/ANA drought monitor |
| `gold/mart_oecd__ag_outlook/mart.parquet` | `analytics.oecd_ag_outlook` | OECD-FAO agricultural outlook |
| `gold/mart_fao__food_price_index/mart.parquet` | `analytics.fao_food_price_index` | FAO food price index |
| `gold/mart_argentina__magyp_producion_granos/mart.parquet` | `analytics.argentina_magyp_producion_granos` | Argentina MAGyP grain production |

Subset sync for wave 3 only:

```bash
UNIFIED_DB_SYNC_MARTS=$(WAVE3_SYNC_MARTS) make unified-db-sync
```

(`WAVE3_SYNC_MARTS` is defined in the root `Makefile`.)

Manifest verification without PostgreSQL:

```bash
make verify-wave3-gold-manifest
```

### Wave 4 marts (Phases 48–49)

| Gold path | PostgreSQL / DuckDB table | Source |
|-----------|---------------------------|--------|
| `gold/mart_ibge__censo_agro_estabelecimentos/mart.parquet` | `analytics.ibge_censo_agro_estabelecimentos` | IBGE Censo Agro 2017 |
| `gold/mart_ibge__pnad_continua_rural/mart.parquet` | `analytics.ibge_pnad_continua_rural` | IBGE PNAD rural occupation |
| `gold/mart_suframa__comercio_mercadorias_zfm/mart.parquet` | `analytics.suframa_comercio_mercadorias_zfm` | SUFRAMA ZFM trade |
| `gold/mart_transportes__mtr_bit_malha_rodoviaria/mart.parquet` | `analytics.transportes_mtr_bit_malha_rodoviaria` | MTR BIT / DNIT SNV roads |
| `gold/mart_mapa__sif_abate_estatisticas/mart.parquet` | `analytics.mapa_sif_abate_estatisticas` | MAPA SIF slaughter |
| `gold/mart_ons__carga_energetica/mart.parquet` | `analytics.ons_carga_energetica` | ONS energy load |
| `gold/mart_inpe__deter_alertas_desmatamento/mart.parquet` | `analytics.inpe_deter_alertas_desmatamento` | INPE DETER alerts |
| `gold/mart_dnit__condicoes_conservacao_rodovias/mart.parquet` | `analytics.dnit_condicoes_conservacao_rodovias` | DNIT pavement condition |
| `gold/mart_cftc__cot_agricultural_futures/mart.parquet` | `analytics.cftc_cot_agricultural_futures` | CFTC Commitment of Traders |
| `gold/mart_jrc__mars_crop_yield/mart.parquet` | `analytics.jrc_mars_crop_yield` | JRC MARS crop yield |
| `gold/mart_fao__giews_crop_prospects/mart.parquet` | `analytics.fao_giews_crop_prospects` | FAO GIEWS |
| `gold/mart_fao__amis_market_monitor/mart.parquet` | `analytics.fao_amis_market_monitor` | FAO AMIS |
| `gold/mart_sagis__grain_supply_statistics/mart.parquet` | `analytics.sagis_grain_supply_statistics` | SAGIS grain supply |
| `gold/mart_japan__maff_ag_trade/mart.parquet` | `analytics.japan_maff_ag_trade` | MAFF Japan ag trade |
| `gold/mart_fred__commodity_indexes/mart.parquet` | `analytics.fred_commodity_indexes` | FRED commodity indexes |
| `gold/mart_nasa__power_agroclimatology/mart.parquet` | `analytics.nasa_power_agroclimatology` | NASA POWER agroclimate |
| `gold/mart_copernicus__era5_agroclimate/mart.parquet` | `analytics.copernicus_era5_agroclimate` | Copernicus ERA5 |

Deferred wave-4 international marts (`wto`, `mexico.siap`, `noaa.gpcc`) are documented in [DEFERRED-SOURCES.md](DEFERRED-SOURCES.md).

Manifest verification without PostgreSQL:

```bash
make verify-wave4-gold-manifest
make spot-check-wave4-duckdb DUCKDB_PATH=duckdb/open_data_agro.duckdb
```

### Wave 5 marts (Phases 52–56)

| Gold path | PostgreSQL / DuckDB table | Source |
|-----------|---------------------------|--------|
| `gold/mart_mapa__sipeagro_estabelecimentos/mart.parquet` | `analytics.mapa_sipeagro_estabelecimentos` | MAPA SIPEAGRO establishments |
| `gold/mart_mapa__sipeagro_produtos/mart.parquet` | `analytics.mapa_sipeagro_produtos` | MAPA SIPEAGRO products |
| `gold/mart_mapa__sigef_producao_sementes/mart.parquet` | `analytics.mapa_sigef_producao_sementes` | MAPA SIGEF seed production |
| `gold/mart_mapa__sigef_areas/mart.parquet` | `analytics.mapa_sigef_areas` | MAPA SIGEF areas |
| `gold/mart_mapa__sisser_seguro_rural/mart.parquet` | `analytics.mapa_sisser_seguro_rural` | MAPA SISSER rural insurance |
| `gold/mart_ibge__ppm_efetivo_rebanhos/mart.parquet` | `analytics.ibge_ppm_efetivo_rebanhos` | IBGE PPM herd |
| `gold/mart_ibge__ppm_vacas_ordenhadas/mart.parquet` | `analytics.ibge_ppm_vacas_ordenhadas` | IBGE PPM dairy cows |
| `gold/mart_ibge__ppm_ovinos_tosquiados/mart.parquet` | `analytics.ibge_ppm_ovinos_tosquiados` | IBGE PPM shorn sheep |
| `gold/mart_ibge__ppm_aquicultura/mart.parquet` | `analytics.ibge_ppm_aquicultura` | IBGE PPM aquaculture |
| `gold/mart_ibge__pam_precos_produtor/mart.parquet` | `analytics.ibge_pam_precos_produtor` | IBGE PAM producer prices |
| `gold/mart_ibge__pam_culturas_estendidas/mart.parquet` | `analytics.ibge_pam_culturas_estendidas` | IBGE PAM extended crops |
| `gold/mart_ibge__lspa_rendimento_medio/mart.parquet` | `analytics.ibge_lspa_rendimento_medio` | IBGE LSPA yield |
| `gold/mart_ibge__censo_agro_area_uso_solo/mart.parquet` | `analytics.ibge_censo_agro_area_uso_solo` | IBGE Censo Agro land use |
| `gold/mart_ibge__censo_agro_maquinario/mart.parquet` | `analytics.ibge_censo_agro_maquinario` | IBGE Censo Agro machinery |
| `gold/mart_ibge__pnad_rural_renda_ocupacao/mart.parquet` | `analytics.ibge_pnad_rural_renda_ocupacao` | IBGE PNAD rural income |
| `gold/mart_ibama__sisfogo_incendios/mart.parquet` | `analytics.ibama_sisfogo_incendios` | IBAMA SISFOGO fires |
| `gold/mart_ibama__licencas_ambientais/mart.parquet` | `analytics.ibama_licencas_ambientais` | IBAMA environmental licenses |
| `gold/mart_ibama__autos_infracao/mart.parquet` | `analytics.ibama_autos_infracao` | IBAMA enforcement notices |
| `gold/mart_ana__pluviometria_redes/mart.parquet` | `analytics.ana_pluviometria_redes` | ANA pluviometry networks |
| `gold/mart_embrapa__agroapi_agrofit/mart.parquet` | `analytics.embrapa_agroapi_agrofit` | Embrapa AgroAPI Agrofit |
| `gold/mart_transportes__mtr_bit_malha_shapefile/mart.parquet` | `analytics.transportes_mtr_bit_malha_shapefile` | MTR BIT rail shapefile |
| `gold/mart_bcb__cim_agro_credito_rural/mart.parquet` | `analytics.bcb_cim_agro_credito_rural` | BCB rural credit rate |
| `gold/mart_bndes__desembolsos_linhas_agro/mart.parquet` | `analytics.bndes_desembolsos_linhas_agro` | BNDES agro disbursements |
| `gold/mart_anp__etanol_precos/mart.parquet` | `analytics.anp_etanol_precos` | ANP ethanol prices |
| `gold/mart_abiove__balanco_complexo_soja/mart.parquet` | `analytics.abiove_balanco_complexo_soja` | Abiove soy complex balance |
| `gold/mart_abiove__exportacoes_complexo_soja/mart.parquet` | `analytics.abiove_exportacoes_complexo_soja` | Abiove soy complex exports |
| `gold/mart_abiove__capacidade_instalada_esmagamento/mart.parquet` | `analytics.abiove_capacidade_instalada_esmagamento` | Abiove crush capacity |
| `gold/mart_b3__futuro_cafe/mart.parquet` | `analytics.b3_futuro_cafe` | B3 coffee futures (ICF) |
| `gold/mart_b3__futuro_acucar/mart.parquet` | `analytics.b3_futuro_acucar` | B3 sugar futures (CNL) |

Subset sync for wave 5 only:

```bash
UNIFIED_DB_SYNC_MARTS=$(WAVE5_SYNC_MARTS) make unified-db-sync
```

(`WAVE5_SYNC_MARTS` is defined in the root `Makefile`.)

Manifest verification without PostgreSQL:

```bash
make verify-wave5-gold-manifest
make spot-check-wave5-duckdb DUCKDB_PATH=duckdb/open_data_agro.duckdb
make ingestor-signoff-wave-5-mvp
```

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
