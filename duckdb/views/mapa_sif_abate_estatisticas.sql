-- Published view: mapa.sif-abate-estatisticas (dbt gold mart).

CREATE OR REPLACE VIEW analytics.mapa_sif_abate_estatisticas AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_mapa__sif_abate_estatisticas/mart.parquet');
