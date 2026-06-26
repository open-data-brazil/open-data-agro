-- Published view: BCRA official USD exchange-rate series (dbt gold mart).

CREATE OR REPLACE VIEW analytics.argentina_bcra_cambio AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_argentina__bcra_cambio/mart.parquet');
