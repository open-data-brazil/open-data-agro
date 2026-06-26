-- Published view: INMET conventional weather station catalog (dbt gold mart).
-- Path placeholder __LAKE_ROOT__ is substituted by duckdb/scripts/analytics-init.sh

CREATE OR REPLACE VIEW analytics.inmet_estacoes_convencionais AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_inmet__estacoes_convencionais/mart.parquet');
