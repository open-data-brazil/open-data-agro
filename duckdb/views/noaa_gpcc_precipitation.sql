-- Published view: noaa gpcc_precipitation (dbt gold mart).
CREATE OR REPLACE VIEW analytics.noaa_gpcc_precipitation AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_noaa__gpcc_precipitation/mart.parquet');
