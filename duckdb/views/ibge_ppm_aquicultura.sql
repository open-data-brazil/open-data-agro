-- Published view: ibge.ppm-aquicultura (dbt gold mart).

CREATE OR REPLACE VIEW analytics.ibge_ppm_aquicultura AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_ibge__ppm_aquicultura/mart.parquet');
