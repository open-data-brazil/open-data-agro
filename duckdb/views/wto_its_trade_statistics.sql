-- Published view: wto its_trade_statistics (dbt gold mart).
CREATE OR REPLACE VIEW analytics.wto_its_trade_statistics AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_wto__its_trade_statistics/mart.parquet');
