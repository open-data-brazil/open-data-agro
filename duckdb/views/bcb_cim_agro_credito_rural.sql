-- Published view: BCB CIM-Agro rural credit rate (dbt gold mart).
-- Path placeholder __LAKE_ROOT__ is substituted by duckdb/scripts/analytics-init.sh

CREATE OR REPLACE VIEW analytics.bcb_cim_agro_credito_rural AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_bcb__cim_agro_credito_rural/mart.parquet');
