-- Published view: ana.pluviometria-redes (dbt gold mart).

CREATE OR REPLACE VIEW analytics.ana_pluviometria_redes AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_ana__pluviometria_redes/mart.parquet');
