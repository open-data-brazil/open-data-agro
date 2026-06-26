-- Published view: ANTAQ port cargo movement (dbt gold mart).

CREATE OR REPLACE VIEW analytics.antaq_movimentacao_carga_portuaria AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_antaq__movimentacao_carga_portuaria/mart.parquet');
