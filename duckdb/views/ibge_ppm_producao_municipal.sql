-- Published view: IBGE PPM municipal production (dbt gold mart).

CREATE OR REPLACE VIEW analytics.ibge_ppm_producao_municipal AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_ibge__ppm_producao_municipal/mart.parquet');
