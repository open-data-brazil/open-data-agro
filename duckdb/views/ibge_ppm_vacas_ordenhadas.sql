-- Published view: ibge.ppm-vacas-ordenhadas (dbt gold mart).

CREATE OR REPLACE VIEW analytics.ibge_ppm_vacas_ordenhadas AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_ibge__ppm_vacas_ordenhadas/mart.parquet');
