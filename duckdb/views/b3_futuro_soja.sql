-- Published view: B3 SOY futures (dbt gold mart).

CREATE OR REPLACE VIEW analytics.b3_futuro_soja AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_b3__futuro_soja/mart.parquet');
