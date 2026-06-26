-- Published view: FAO FAOSTAT annual trade quantities (dbt gold mart).

CREATE OR REPLACE VIEW analytics.fao_comercio_agro AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_fao__comercio_agro/mart.parquet');
