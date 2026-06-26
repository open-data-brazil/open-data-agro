-- Published view: CONAB sugarcane historical series (dbt gold mart).
-- Path placeholder __LAKE_ROOT__ is substituted by duckdb/scripts/analytics-init.sh

CREATE OR REPLACE VIEW analytics.conab_serie_historica_cana AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_conab__serie_historica_cana/mart.parquet');
