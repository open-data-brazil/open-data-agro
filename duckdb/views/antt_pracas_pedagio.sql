-- Published view: ANTT praças de pedágio (dbt gold mart).
-- Path placeholder __LAKE_ROOT__ is substituted by duckdb/scripts/analytics-init.sh

CREATE OR REPLACE VIEW analytics.antt_pracas_pedagio AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_antt__pracas_pedagio/mart.parquet');
