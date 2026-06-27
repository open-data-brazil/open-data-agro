CREATE OR REPLACE VIEW analytics.anp_etanol_precos AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_anp__etanol_precos/mart.parquet');
