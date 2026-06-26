-- Published view: ANA HidroWeb daily hydrology series (dbt gold mart).

CREATE OR REPLACE VIEW analytics.ana_hidrologia_series AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_ana__hidrologia_series/mart.parquet');
