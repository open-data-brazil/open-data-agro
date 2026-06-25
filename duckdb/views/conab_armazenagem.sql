CREATE OR REPLACE VIEW analytics.conab_armazenagem AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_conab__armazenagem/mart.parquet');
