-- Published view: BCB SGS PTAX USD sell (dbt gold mart).
-- Path placeholder __LAKE_ROOT__ is substituted by duckdb/scripts/analytics-init.sh

CREATE OR REPLACE VIEW analytics.bcb_sgs_ptax_usd_venda AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_bcb__sgs_ptax_usd_venda/mart.parquet');
