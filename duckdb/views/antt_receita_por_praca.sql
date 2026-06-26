-- Published view: ANTT toll revenue per plaza (dbt gold mart).
-- Path placeholder __LAKE_ROOT__ is substituted by duckdb/scripts/analytics-init.sh

CREATE OR REPLACE VIEW analytics.antt_receita_por_praca AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_antt__receita_por_praca/mart.parquet');
