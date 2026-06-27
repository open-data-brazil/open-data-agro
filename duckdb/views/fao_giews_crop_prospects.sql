-- Published view: fao giews_crop_prospects (dbt gold mart).
CREATE OR REPLACE VIEW analytics.fao_giews_crop_prospects AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_fao__giews_crop_prospects/mart.parquet');
