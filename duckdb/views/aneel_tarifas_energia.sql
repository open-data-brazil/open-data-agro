-- Published view: ANEEL tariff flag activation (dbt gold mart).

CREATE OR REPLACE VIEW analytics.aneel_tarifas_energia AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_aneel__tarifas_energia/mart.parquet');
