-- Published view: cftc cot_agricultural_futures (dbt gold mart).
CREATE OR REPLACE VIEW analytics.cftc_cot_agricultural_futures AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_cftc__cot_agricultural_futures/mart.parquet');
