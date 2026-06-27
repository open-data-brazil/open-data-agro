CREATE OR REPLACE VIEW analytics.abiove_capacidade_instalada_esmagamento AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_abiove__capacidade_instalada_esmagamento/mart.parquet');
