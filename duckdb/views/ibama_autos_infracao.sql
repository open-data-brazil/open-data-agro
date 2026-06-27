-- Published view: ibama.autos-infracao (dbt gold mart).

CREATE OR REPLACE VIEW analytics.ibama_autos_infracao AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_ibama__autos_infracao/mart.parquet');
