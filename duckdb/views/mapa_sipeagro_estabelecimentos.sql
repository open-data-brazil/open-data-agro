-- Published view: mapa.sipeagro-estabelecimentos (dbt gold mart).

CREATE OR REPLACE VIEW analytics.mapa_sipeagro_estabelecimentos AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_mapa__sipeagro_estabelecimentos/mart.parquet');
