# Great Expectations — local-first quality gates

GX project root: `expectations/gx/` (GE 1.18, config version 4.0).

Canonical suite sources live under `expectations/suites/`; copies synced to the GX store use flat names (e.g. `bronze.conab.estimativa_graos`).

## Data sources

| Layer | Path pattern | When validated |
|-------|--------------|----------------|
| Bronze | `{LAKE_LOCAL_ROOT}/bronze/conab/{slug}/ingest_date=*/part-*.parquet` | Before silver promotion |
| Silver | `{LAKE_LOCAL_ROOT}/silver/conab/{table}/` | Post-dbt (future) |

Default `LAKE_LOCAL_ROOT` is `./lake` (see `internal/config/lake.go`). Validation reads **local Parquet files only** — no R2 required.

For `STORAGE_MODE=r2` or `minio`, bronze files must still be materialized locally (or synced) before running checkpoints; object-store fluent datasources are optional and not wired in CI.

## Checkpoints

| Checkpoint | Suite | Dataset |
|------------|-------|---------|
| `bronze_conab_estimativa_graos` | `bronze.conab.estimativa_graos` | `conab.estimativa-graos` |

## Expectation list — `conab.estimativa-graos` (bronze)

Source: `expectations/suites/bronze/conab/estimativa_graos.json`

| Expectation | Target | Business rule |
|-------------|--------|---------------|
| `expect_table_row_count_to_be_between` | table | min ≥ 1 row |
| `expect_column_to_exist` | `Produto` | CONAB crop column |
| `expect_column_to_exist` | `UF` | Federative unit |
| `expect_column_to_exist` | `Safra` | Harvest season (grain key) |
| `expect_column_values_to_not_be_null` | `Produto`, `UF`, `Safra` | Grain keys must be present |

**Traceability**

- `dataset_id`: `conab.estimativa-graos`
- `fonte_oficial`: https://portaldeinformacoes.conab.gov.br/download-arquivos.html
- `conab_section`: Produção Agrícola > Estimativa Grãos

See also `.local/phases/06-quality-great-expectations/OFFICIAL-REFERENCE.md` and phase 10 business rules.

## Run locally

```bash
# Seed sample bronze (CI helper)
LAKE_LOCAL_ROOT=/tmp/open-data-agro-lake python3 scripts/ci/seed_ge_bronze.py

# Checkpoint (acceptance)
python3 scripts/quality/run_checkpoint.py \
  --checkpoint bronze_conab_estimativa_graos \
  --bronze-dir "$LAKE_LOCAL_ROOT/bronze/conab/estimativa-graos"

# Or via processor CLI (gates promotion)
processor quality --dataset conab.estimativa-graos
processor promote --dataset conab.estimativa-graos   # runs quality first
```

On quality failure, `processor promote` records `status=quality_failed` in `catalog.promotion_jobs` and does not write silver.
