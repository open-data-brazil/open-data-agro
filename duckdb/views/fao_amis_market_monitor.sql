-- Published view: fao amis_market_monitor (dbt gold mart).
CREATE OR REPLACE VIEW analytics.fao_amis_market_monitor AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_fao__amis_market_monitor/mart.parquet');
