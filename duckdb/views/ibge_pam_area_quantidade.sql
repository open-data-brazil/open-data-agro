-- Published view: IBGE PAM area and production quantity (dbt gold mart).
-- Path placeholder __LAKE_ROOT__ is substituted by duckdb/scripts/analytics-init.sh

CREATE OR REPLACE VIEW analytics.ibge_pam_area_quantidade AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_ibge__pam_area_quantidade/mart.parquet');
