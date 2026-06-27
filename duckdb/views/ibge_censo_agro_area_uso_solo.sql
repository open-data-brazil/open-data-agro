-- Published view: ibge.censo-agro-area-uso-solo (dbt gold mart).

CREATE OR REPLACE VIEW analytics.ibge_censo_agro_area_uso_solo AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_ibge__censo_agro_area_uso_solo/mart.parquet');
