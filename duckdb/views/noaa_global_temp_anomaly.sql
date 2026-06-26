-- Published view: NOAA global land+ocean monthly temperature anomaly (dbt gold mart).

CREATE OR REPLACE VIEW analytics.noaa_global_temp_anomaly AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_noaa__global_temp_anomaly/mart.parquet');
