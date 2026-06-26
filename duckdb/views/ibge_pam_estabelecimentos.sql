-- Published view: IBGE PAM farm establishment counts (dbt gold mart).
-- Path placeholder __LAKE_ROOT__ is substituted by duckdb/scripts/analytics-init.sh

CREATE OR REPLACE VIEW analytics.ibge_pam_estabelecimentos AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_ibge__pam_estabelecimentos/mart.parquet');
