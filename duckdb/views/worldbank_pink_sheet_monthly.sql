-- Published view: World Bank Pink Sheet monthly prices (dbt gold mart).

CREATE OR REPLACE VIEW analytics.worldbank_pink_sheet_monthly AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_worldbank__pink_sheet_monthly/mart.parquet');
