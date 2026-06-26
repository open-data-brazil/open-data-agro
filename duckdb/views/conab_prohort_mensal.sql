-- Published view: CONAB Prohort monthly CEASA trade (dbt gold mart).
-- Path placeholder __LAKE_ROOT__ is substituted by duckdb/scripts/analytics-init.sh

CREATE OR REPLACE VIEW analytics.conab_prohort_mensal AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_conab__prohort_mensal/mart.parquet');
