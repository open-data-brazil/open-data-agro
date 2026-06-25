-- Published view: CONAB supply/demand balance (dbt gold mart).
-- Path placeholder __LAKE_ROOT__ is substituted by duckdb/scripts/analytics-init.sh

CREATE OR REPLACE VIEW analytics.conab_oferta_demanda AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_conab__oferta_demanda/mart.parquet');
