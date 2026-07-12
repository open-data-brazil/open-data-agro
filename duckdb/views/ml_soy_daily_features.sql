-- Published view: ML soy daily features (Phase 20 analytics).
-- Path placeholder __LAKE_ROOT__ is substituted by duckdb/scripts/analytics-init.sh

CREATE OR REPLACE VIEW analytics.ml_soy_daily_features AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_ml__soy_daily_features/mart.parquet');
