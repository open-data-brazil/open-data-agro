-- Published view: CEPEA soja Paranaguá (dbt gold mart).
-- Path placeholder __LAKE_ROOT__ is substituted by duckdb/scripts/analytics-init.sh

CREATE OR REPLACE VIEW analytics.cepea_soja_paranagua AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_cepea__soja_paranagua/mart.parquet');
