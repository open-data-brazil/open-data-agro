CREATE OR REPLACE VIEW analytics.abiove_balanco_complexo_soja AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_abiove__balanco_complexo_soja/mart.parquet');
