-- Published view: World Bank Pink Sheet agriculture sub-indices (dbt gold mart).

CREATE OR REPLACE VIEW analytics.worldbank_ag_indices AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_worldbank__ag_indices/mart.parquet');
