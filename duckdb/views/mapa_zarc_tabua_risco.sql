-- Published view: MAPA ZARC tábua de risco (dbt gold mart).
-- Path placeholder __LAKE_ROOT__ is substituted by duckdb/scripts/analytics-init.sh

CREATE OR REPLACE VIEW analytics.mapa_zarc_tabua_risco AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_mapa__zarc_tabua_risco/mart.parquet');
