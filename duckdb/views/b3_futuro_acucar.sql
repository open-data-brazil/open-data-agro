CREATE OR REPLACE VIEW analytics.b3_futuro_acucar AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_b3__futuro_acucar/mart.parquet');
