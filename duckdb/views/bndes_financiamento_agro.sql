-- Published view: BNDES agro disbursements by CNAE sector (dbt gold mart).

CREATE OR REPLACE VIEW analytics.bndes_financiamento_agro AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_bndes__financiamento_agro/mart.parquet');
