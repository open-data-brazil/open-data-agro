-- Published view: CEPEA milho Campinas (dbt gold mart).
-- Path placeholder __LAKE_ROOT__ is substituted by duckdb/scripts/analytics-init.sh

CREATE OR REPLACE VIEW analytics.cepea_milho AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_cepea__milho/mart.parquet');
