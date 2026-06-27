-- Published view: japan maff_ag_trade (dbt gold mart).
CREATE OR REPLACE VIEW analytics.japan_maff_ag_trade AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_japan__maff_ag_trade/mart.parquet');
