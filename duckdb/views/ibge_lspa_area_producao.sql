-- Published view: IBGE LSPA monthly UF area/production (dbt gold mart).

CREATE OR REPLACE VIEW analytics.ibge_lspa_area_producao AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_ibge__lspa_area_producao/mart.parquet');
