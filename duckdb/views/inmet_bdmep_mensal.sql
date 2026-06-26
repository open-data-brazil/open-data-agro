-- Published view: INMET BDMEP monthly climate rollups (dbt gold mart).
-- Path placeholder __LAKE_ROOT__ is substituted by duckdb/scripts/analytics-init.sh

CREATE OR REPLACE VIEW analytics.inmet_bdmep_mensal AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_inmet__bdmep_mensal/mart.parquet');
