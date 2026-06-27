-- Published view: ibge.censo-agro-maquinario (dbt gold mart).

CREATE OR REPLACE VIEW analytics.ibge_censo_agro_maquinario AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_ibge__censo_agro_maquinario/mart.parquet');
