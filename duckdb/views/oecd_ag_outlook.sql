-- Published view: OECD-FAO Agricultural Outlook (dbt gold mart).
CREATE OR REPLACE VIEW analytics.oecd_ag_outlook AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_oecd__ag_outlook/mart.parquet');
