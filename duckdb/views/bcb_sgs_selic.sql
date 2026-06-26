-- Published view: BCB SGS Selic target rate (dbt gold mart).

CREATE OR REPLACE VIEW analytics.bcb_sgs_selic AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_bcb__sgs_selic/mart.parquet');
