-- Published view: CONAB production cost survey (dbt gold mart).
-- Path placeholder __LAKE_ROOT__ is substituted by duckdb/scripts/analytics-init.sh

CREATE OR REPLACE VIEW analytics.conab_custo_producao AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_conab__custo_producao/mart.parquet');
