-- Published view: mapa.sisser-seguro-rural (dbt gold mart).

CREATE OR REPLACE VIEW analytics.mapa_sisser_seguro_rural AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_mapa__sisser_seguro_rural/mart.parquet');
