-- Published view: sagis grain_supply_statistics (dbt gold mart).
CREATE OR REPLACE VIEW analytics.sagis_grain_supply_statistics AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_sagis__grain_supply_statistics/mart.parquet');
