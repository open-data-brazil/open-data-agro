-- Published view: MDIC Comex diesel import (dbt gold mart).
-- Path placeholder __LAKE_ROOT__ is substituted by duckdb/scripts/analytics-init.sh

CREATE OR REPLACE VIEW analytics.mdic_comex_importacao_diesel_ncm AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_mdic__comex_importacao_diesel_ncm/mart.parquet');
