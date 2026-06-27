-- Published view: fred commodity_indexes (dbt gold mart).
CREATE OR REPLACE VIEW analytics.fred_commodity_indexes AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_fred__commodity_indexes/mart.parquet');
