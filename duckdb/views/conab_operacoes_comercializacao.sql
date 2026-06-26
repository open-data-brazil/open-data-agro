-- Published view: CONAB commercialization operations (dbt gold mart).
-- Path placeholder __LAKE_ROOT__ is substituted by duckdb/scripts/analytics-init.sh

CREATE OR REPLACE VIEW analytics.conab_operacoes_comercializacao AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_conab__operacoes_comercializacao/mart.parquet');
