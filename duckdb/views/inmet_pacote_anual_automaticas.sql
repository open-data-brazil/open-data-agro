-- Published view: INMET BDMEP annual automatic-station package (dbt gold mart).
-- Path placeholder __LAKE_ROOT__ is substituted by duckdb/scripts/analytics-init.sh

CREATE OR REPLACE VIEW analytics.inmet_pacote_anual_automaticas AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_inmet__pacote_anual_automaticas/mart.parquet');
