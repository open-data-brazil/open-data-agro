-- Published view: IBGE municipality reference (dbt gold mart).
-- Path placeholder __LAKE_ROOT__ is substituted by duckdb/scripts/analytics-init.sh

CREATE OR REPLACE VIEW analytics.ibge_localidades_municipios AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_ibge__localidades_municipios/mart.parquet');
