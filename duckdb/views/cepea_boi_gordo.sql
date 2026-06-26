-- Published view: CEPEA boi gordo São Paulo (dbt gold mart).
-- Path placeholder __LAKE_ROOT__ is substituted by duckdb/scripts/analytics-init.sh

CREATE OR REPLACE VIEW analytics.cepea_boi_gordo AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_cepea__boi_gordo/mart.parquet');
