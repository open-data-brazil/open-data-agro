-- Published view: IBGE PEVS plant extraction by UF (dbt gold mart).

CREATE OR REPLACE VIEW analytics.ibge_pevs_producao_vegetal AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_ibge__pevs_producao_vegetal/mart.parquet');
