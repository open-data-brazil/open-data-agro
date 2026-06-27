-- Published view: mexico siap_produccion_agricola (dbt gold mart).
CREATE OR REPLACE VIEW analytics.mexico_siap_produccion_agricola AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_mexico__siap_produccion_agricola/mart.parquet');
