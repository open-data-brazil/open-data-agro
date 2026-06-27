-- Published view: ibge.pam-precos-produtor (dbt gold mart).

CREATE OR REPLACE VIEW analytics.ibge_pam_precos_produtor AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_ibge__pam_precos_produtor/mart.parquet');
