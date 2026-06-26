-- Published view: IPEA regional macro series (dbt gold mart).

CREATE OR REPLACE VIEW analytics.ipea_series_macro_regionais AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_ipea__series_macro_regionais/mart.parquet');
