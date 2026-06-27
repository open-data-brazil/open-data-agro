-- Published view: ibge.ppm-ovinos-tosquiados (dbt gold mart).

CREATE OR REPLACE VIEW analytics.ibge_ppm_ovinos_tosquiados AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_ibge__ppm_ovinos_tosquiados/mart.parquet');
