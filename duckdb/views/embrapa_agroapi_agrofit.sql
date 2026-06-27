-- Published view: embrapa.agroapi-agrofit (dbt gold mart).

CREATE OR REPLACE VIEW analytics.embrapa_agroapi_agrofit AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_embrapa__agroapi_agrofit/mart.parquet');
