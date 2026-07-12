-- Published view: municipality dimension for analytics joins (Phase 20).

CREATE OR REPLACE VIEW analytics.dim_municipio AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/dim_municipio/mart.parquet');
