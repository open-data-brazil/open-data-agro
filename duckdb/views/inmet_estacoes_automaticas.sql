-- Published view: INMET automatic weather station catalog (dbt gold mart).
-- Path placeholder __LAKE_ROOT__ is substituted by duckdb/scripts/analytics-init.sh

CREATE OR REPLACE VIEW analytics.inmet_estacoes_automaticas AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_inmet__estacoes_automaticas/mart.parquet');
