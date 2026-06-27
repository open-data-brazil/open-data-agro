-- Published view: transportes.mtr-bit-malha-shapefile (dbt gold mart).

CREATE OR REPLACE VIEW analytics.transportes_mtr_bit_malha_shapefile AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_transportes__mtr_bit_malha_shapefile/mart.parquet');
