CREATE OR REPLACE VIEW analytics.conab_alimenta_brasil_propostas AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_conab__alimenta_brasil_propostas/mart.parquet');
