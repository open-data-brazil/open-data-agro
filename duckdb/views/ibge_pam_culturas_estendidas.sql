-- Published view: ibge.pam-culturas-estendidas (dbt gold mart).

CREATE OR REPLACE VIEW analytics.ibge_pam_culturas_estendidas AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_ibge__pam_culturas_estendidas/mart.parquet');
