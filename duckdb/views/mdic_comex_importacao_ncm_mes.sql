-- Published view: MDIC Comex import by NCM (dbt gold mart).
-- Path placeholder __LAKE_ROOT__ is substituted by duckdb/scripts/analytics-init.sh

CREATE OR REPLACE VIEW analytics.mdic_comex_importacao_ncm_mes AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_mdic__comex_importacao_ncm_mes/mart.parquet');
