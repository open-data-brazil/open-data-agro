-- Published view: B3 CCM corn futures (dbt gold mart).

CREATE OR REPLACE VIEW analytics.b3_futuro_milho AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_b3__futuro_milho/mart.parquet');
