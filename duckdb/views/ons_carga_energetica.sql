-- Published view: ons.carga-energetica (dbt gold mart).

CREATE OR REPLACE VIEW analytics.ons_carga_energetica AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_ons__carga_energetica/mart.parquet');
