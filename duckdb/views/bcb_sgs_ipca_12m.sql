-- Published view: BCB SGS IPCA 12m (dbt gold mart).
-- Path placeholder __LAKE_ROOT__ is substituted by duckdb/scripts/analytics-init.sh

CREATE OR REPLACE VIEW analytics.bcb_sgs_ipca_12m AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_bcb__sgs_ipca_12m/mart.parquet');
