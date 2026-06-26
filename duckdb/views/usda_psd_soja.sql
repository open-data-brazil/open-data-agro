-- Published view: USDA FAS soybean PSD (dbt gold mart).

CREATE OR REPLACE VIEW analytics.usda_psd_soja AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_usda__psd_soja/mart.parquet');
