-- Published view: ANTT toll traffic volume (dbt gold mart).
-- Path placeholder __LAKE_ROOT__ is substituted by duckdb/scripts/analytics-init.sh

CREATE OR REPLACE VIEW analytics.antt_volume_trafego_pedagio AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_antt__volume_trafego_pedagio/mart.parquet');
