-- Published view: CONAB grain production estimates (dbt gold mart).
-- Path placeholder __LAKE_ROOT__ is substituted by duckdb/scripts/analytics-init.sh

CREATE OR REPLACE VIEW analytics.conab_estimativa_graos AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_conab__estimativa_graos/mart.parquet');
