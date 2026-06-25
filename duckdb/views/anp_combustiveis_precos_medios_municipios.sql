CREATE OR REPLACE VIEW analytics.anp_combustiveis_precos_medios_municipios AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_anp__combustiveis_precos_medios_municipios/mart.parquet');
