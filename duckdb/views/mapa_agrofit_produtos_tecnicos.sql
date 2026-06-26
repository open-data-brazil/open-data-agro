-- Published view: MAPA Agrofit technical products registry (dbt gold mart).

CREATE OR REPLACE VIEW analytics.mapa_agrofit_produtos_tecnicos AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_mapa__agrofit_produtos_tecnicos/mart.parquet');
