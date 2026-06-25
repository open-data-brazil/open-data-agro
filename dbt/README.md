# dbt — local transform (Phase 5)

dbt reads **silver Delta** under `LAKE_LOCAL_ROOT` and writes **gold marts** to `lake/gold/`.

## Setup

```bash
python3 -m venv .venv && source .venv/bin/activate
pip install -r toolchain/python-requirements.txt
cp dbt/profiles.yml.example ~/.dbt/profiles.yml   # or dbt/profiles.yml
export LAKE_LOCAL_ROOT=./lake
```

## Build (local)

```bash
./bin/ingestor run conab.estimativa-graos
./bin/processor promote --dataset conab.estimativa-graos
cd dbt && dbt deps && dbt build --select 'stg_conab__serie_historica_graos stg_conab__estimativa_graos+'
ls ../lake/gold/mart_conab__estimativa_graos/
```

## Profiles

| Target | Purpose |
|--------|---------|
| `dev_local` | Default — DuckDB file + `lake/gold` external root |
| `prod_r2` | Same models; `external_root` on `s3://{R2_BUCKET}/gold` |

See [profiles.yml.example](profiles.yml.example).

## Layout

```text
models/staging/conab/     → views over silver Parquet
models/intermediate/conab/ → unions / shared logic
models/marts/conab/       → external Parquet under lake/gold/
```
