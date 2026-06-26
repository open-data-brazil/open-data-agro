-- Published view: BCB SGS IGP-M (dbt gold mart).
-- Path placeholder __LAKE_ROOT__ is substituted by duckdb/scripts/analytics-init.sh

CREATE OR REPLACE VIEW analytics.bcb_sgs_igpm AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_bcb__sgs_igpm/mart.parquet');
