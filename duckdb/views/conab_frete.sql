-- Published view: CONAB freight routes (dbt gold mart).
-- Path placeholder __LAKE_ROOT__ is substituted by duckdb/scripts/analytics-init.sh

CREATE OR REPLACE VIEW analytics.conab_frete AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_conab__frete/mart.parquet');
