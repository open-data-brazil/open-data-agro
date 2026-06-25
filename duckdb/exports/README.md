# Portable analytics exports

Researchers can copy Parquet files from this directory without cloud credentials.

## Generate

```bash
make analytics-init
./duckdb/export-mart.sh conab_estimativa_graos
./duckdb/export-mart.sh conab_serie_historica_graos
```

Each export includes a `_metadata.json` sidecar with CONAB attribution.

## Attribution (mandatory)

Every export MUST cite:

- **Source:** CONAB — Portal de Informações Agropecuárias
- **URL:** https://portaldeinformacoes.conab.gov.br/download-arquivos.html

The sidecar JSON duplicates this for automated pipelines.

## Layout

```text
exports/
├── conab_estimativa_graos-YYYY-MM-DD.parquet
├── conab_estimativa_graos-YYYY-MM-DD_metadata.json
└── README.md
```

Exports are gitignored; regenerate locally after `dbt build`.
