-- Published view: USDA FAS GATS trade statistics (dbt gold mart).

CREATE OR REPLACE VIEW analytics.usda_gats_trade AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_usda__gats_trade/mart.parquet');
