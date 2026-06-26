-- Published view: FAO FAOSTAT producer prices (dbt gold mart).

CREATE OR REPLACE VIEW analytics.fao_prices_agro AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_fao__prices_agro/mart.parquet');
