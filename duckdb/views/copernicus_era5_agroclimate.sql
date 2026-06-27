-- Published view: copernicus era5_agroclimate (dbt gold mart).
CREATE OR REPLACE VIEW analytics.copernicus_era5_agroclimate AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_copernicus__era5_agroclimate/mart.parquet');
