-- Published view: FAO Food Price Index (dbt gold mart).
CREATE OR REPLACE VIEW analytics.fao_food_price_index AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_fao__food_price_index/mart.parquet');
