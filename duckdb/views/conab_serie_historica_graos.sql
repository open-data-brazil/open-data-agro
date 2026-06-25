-- Published view: CONAB historical grain series (local silver Parquet).
-- Path placeholder __LAKE_ROOT__ is substituted by duckdb/scripts/analytics-init.sh

CREATE OR REPLACE VIEW analytics.conab_serie_historica_graos AS
SELECT
    "Produto" AS produto,
    "UF" AS uf,
    "Ano" AS ano,
    "Produção (mil t)" AS producao_mil_t,
    cast(_ingested_at AS VARCHAR) AS capturado_em,
    'https://portaldeinformacoes.conab.gov.br/download-arquivos.html' AS fonte_oficial,
    _dataset_id,
    _source_file
FROM read_parquet(
    '__LAKE_ROOT__/silver/conab/serie_historica_graos/**/*.parquet',
    union_by_name := true
);
