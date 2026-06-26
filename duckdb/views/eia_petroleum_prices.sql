-- Published view: U.S. EIA petroleum spot prices (dbt gold mart).

CREATE OR REPLACE VIEW analytics.eia_petroleum_prices AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_eia__petroleum_prices/mart.parquet');
