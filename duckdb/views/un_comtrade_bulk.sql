-- Published view: UN Comtrade Brazil ag trade (dbt gold mart).
CREATE OR REPLACE VIEW analytics.un_comtrade_bulk AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_un__comtrade_bulk/mart.parquet');
