-- Published view: USDA WASDE monthly supply/demand (dbt gold mart).
-- Path placeholder __LAKE_ROOT__ is substituted by duckdb/scripts/analytics-init.sh
CREATE OR REPLACE VIEW analytics.usda_wasde AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_usda__wasde/mart.parquet');
