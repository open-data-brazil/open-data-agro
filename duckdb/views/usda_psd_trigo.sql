-- Published view: USDA FAS wheat PSD (dbt gold mart).

CREATE OR REPLACE VIEW analytics.usda_psd_trigo AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_usda__psd_trigo/mart.parquet');
