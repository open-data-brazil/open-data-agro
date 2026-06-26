-- Published view: CONAB sugarcane estimate (dbt gold mart).
-- Path placeholder __LAKE_ROOT__ is substituted by duckdb/scripts/analytics-init.sh

CREATE OR REPLACE VIEW analytics.conab_estimativa_cana AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_conab__estimativa_cana/mart.parquet');
