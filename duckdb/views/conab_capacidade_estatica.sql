-- Published view: CONAB static storage capacity by UF and year (dbt gold mart).
-- Path placeholder __LAKE_ROOT__ is substituted by duckdb/scripts/analytics-init.sh

CREATE OR REPLACE VIEW analytics.conab_capacidade_estatica AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_conab__capacidade_estatica/mart.parquet');
