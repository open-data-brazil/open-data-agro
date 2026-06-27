CREATE OR REPLACE VIEW analytics.abiove_exportacoes_complexo_soja AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_abiove__exportacoes_complexo_soja/mart.parquet');
