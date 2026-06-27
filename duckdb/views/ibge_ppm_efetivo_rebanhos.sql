-- Published view: ibge.ppm-efetivo-rebanhos (dbt gold mart).

CREATE OR REPLACE VIEW analytics.ibge_ppm_efetivo_rebanhos AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_ibge__ppm_efetivo_rebanhos/mart.parquet');
