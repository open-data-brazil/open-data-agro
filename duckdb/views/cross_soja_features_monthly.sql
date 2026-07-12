-- Published view: monthly municipal soja crossing features (Phase 20).

CREATE OR REPLACE VIEW analytics.cross_soja_features_monthly AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_cross__soja_features_monthly/mart.parquet');
