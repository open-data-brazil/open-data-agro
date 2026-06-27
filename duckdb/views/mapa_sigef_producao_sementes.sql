-- Published view: mapa.sigef-producao-sementes (dbt gold mart).

CREATE OR REPLACE VIEW analytics.mapa_sigef_producao_sementes AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_mapa__sigef_producao_sementes/mart.parquet');
