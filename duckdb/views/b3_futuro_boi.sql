-- Published view: B3 BGI cattle futures (dbt gold mart).

CREATE OR REPLACE VIEW analytics.b3_futuro_boi AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_b3__futuro_boi/mart.parquet');
