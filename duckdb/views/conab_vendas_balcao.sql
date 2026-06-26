-- Published view: CONAB counter sales by municipality (dbt gold mart).
-- Path placeholder __LAKE_ROOT__ is substituted by duckdb/scripts/analytics-init.sh

CREATE OR REPLACE VIEW analytics.conab_vendas_balcao AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_conab__vendas_balcao/mart.parquet');
