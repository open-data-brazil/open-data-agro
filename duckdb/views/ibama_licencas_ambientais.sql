-- Published view: ibama.licencas-ambientais (dbt gold mart).

CREATE OR REPLACE VIEW analytics.ibama_licencas_ambientais AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_ibama__licencas_ambientais/mart.parquet');
