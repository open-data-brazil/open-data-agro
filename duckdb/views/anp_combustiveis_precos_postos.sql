CREATE OR REPLACE VIEW analytics.anp_combustiveis_precos_postos AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_anp__combustiveis_precos_postos/mart.parquet');
