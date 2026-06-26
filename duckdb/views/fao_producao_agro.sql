-- Published view: FAO FAOSTAT annual production (dbt gold mart).

CREATE OR REPLACE VIEW analytics.fao_producao_agro AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_fao__producao_agro/mart.parquet');
