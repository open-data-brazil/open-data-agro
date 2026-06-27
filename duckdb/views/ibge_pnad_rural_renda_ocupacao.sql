-- Published view: ibge.pnad-rural-renda-ocupacao (dbt gold mart).

CREATE OR REPLACE VIEW analytics.ibge_pnad_rural_renda_ocupacao AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_ibge__pnad_rural_renda_ocupacao/mart.parquet');
