-- Published view: CONAB coffee estimate (dbt gold mart).
-- Path placeholder __LAKE_ROOT__ is substituted by duckdb/scripts/analytics-init.sh

CREATE OR REPLACE VIEW analytics.conab_estimativa_cafe AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_conab__estimativa_cafe/mart.parquet');
