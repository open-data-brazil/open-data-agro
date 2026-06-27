-- Published view: ibama.sisfogo-incendios (dbt gold mart).

CREATE OR REPLACE VIEW analytics.ibama_sisfogo_incendios AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_ibama__sisfogo_incendios/mart.parquet');
