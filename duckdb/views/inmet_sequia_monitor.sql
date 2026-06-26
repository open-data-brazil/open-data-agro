-- Published view: INMET/ANA drought monitor areas (dbt gold mart).

CREATE OR REPLACE VIEW analytics.inmet_sequia_monitor AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_inmet__sequia_monitor/mart.parquet');
