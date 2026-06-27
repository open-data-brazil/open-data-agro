-- Published view: ibge.lspa-rendimento-medio (dbt gold mart).

CREATE OR REPLACE VIEW analytics.ibge_lspa_rendimento_medio AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_ibge__lspa_rendimento_medio/mart.parquet');
