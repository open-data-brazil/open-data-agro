-- Published view: CONAB monthly ag prices by municipality (dbt gold mart).
-- Path placeholder __LAKE_ROOT__ is substituted by duckdb/scripts/analytics-init.sh

CREATE OR REPLACE VIEW analytics.conab_precos_mensal_municipio AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_conab__precos_mensal_municipio/mart.parquet');
