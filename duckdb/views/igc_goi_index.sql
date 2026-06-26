-- Published view: IGC Grains & Oilseeds Index (dbt gold mart).
CREATE OR REPLACE VIEW analytics.igc_goi_index AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_igc__goi_index/mart.parquet');
