-- Published view: USDA FAS corn PSD (dbt gold mart).

CREATE OR REPLACE VIEW analytics.usda_psd_milho AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_usda__psd_milho/mart.parquet');
