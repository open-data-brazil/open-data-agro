-- Published view: NOAA Oceanic Niño Index (ONI) seasonal anomalies (dbt gold mart).

CREATE OR REPLACE VIEW analytics.noaa_enso_indices AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_noaa__enso_indices/mart.parquet');
