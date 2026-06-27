-- Published view: nasa power_agroclimatology (dbt gold mart).
CREATE OR REPLACE VIEW analytics.nasa_power_agroclimatology AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_nasa__power_agroclimatology/mart.parquet');
