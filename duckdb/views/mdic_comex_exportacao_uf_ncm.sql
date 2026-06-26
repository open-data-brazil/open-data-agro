-- Published view: MDIC Comex export by UF × NCM (dbt gold mart).
-- Path placeholder __LAKE_ROOT__ is substituted by duckdb/scripts/analytics-init.sh

CREATE OR REPLACE VIEW analytics.mdic_comex_exportacao_uf_ncm AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_mdic__comex_exportacao_uf_ncm/mart.parquet');
