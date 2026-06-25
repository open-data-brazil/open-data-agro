CREATE OR REPLACE VIEW analytics.conab_estoques_publicos AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_conab__estoques_publicos/mart.parquet');
