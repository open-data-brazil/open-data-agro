-- Published view: ibge.censo-agro-estabelecimentos (dbt gold mart).

CREATE OR REPLACE VIEW analytics.ibge_censo_agro_estabelecimentos AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_ibge__censo_agro_estabelecimentos/mart.parquet');
