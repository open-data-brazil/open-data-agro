-- Published view: jrc mars_crop_yield (dbt gold mart).
CREATE OR REPLACE VIEW analytics.jrc_mars_crop_yield AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_jrc__mars_crop_yield/mart.parquet');
