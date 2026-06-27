-- Published view: mapa.sipeagro-produtos (dbt gold mart).

CREATE OR REPLACE VIEW analytics.mapa_sipeagro_produtos AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_mapa__sipeagro_produtos/mart.parquet');
