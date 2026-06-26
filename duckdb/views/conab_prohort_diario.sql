-- Published view: CONAB Prohort daily CEASA prices (dbt gold mart).
-- Path placeholder __LAKE_ROOT__ is substituted by duckdb/scripts/analytics-init.sh

CREATE OR REPLACE VIEW analytics.conab_prohort_diario AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_conab__prohort_diario/mart.parquet');
