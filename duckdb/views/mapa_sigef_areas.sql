-- Published view: mapa.sigef-areas (dbt gold mart).

CREATE OR REPLACE VIEW analytics.mapa_sigef_areas AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_mapa__sigef_areas/mart.parquet');
