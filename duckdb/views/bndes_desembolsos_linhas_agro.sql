-- Published view: BNDES disbursements by product line (dbt gold mart).

CREATE OR REPLACE VIEW analytics.bndes_desembolsos_linhas_agro AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_bndes__desembolsos_linhas_agro/mart.parquet');
