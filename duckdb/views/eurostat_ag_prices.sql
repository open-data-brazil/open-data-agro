-- Published view: EUROSTAT agricultural output price indices (dbt gold mart).

CREATE OR REPLACE VIEW analytics.eurostat_ag_prices AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_eurostat__ag_prices/mart.parquet');
