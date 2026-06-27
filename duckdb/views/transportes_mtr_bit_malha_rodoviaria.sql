-- Published view: transportes.mtr-bit-malha-rodoviaria (dbt gold mart).

CREATE OR REPLACE VIEW analytics.transportes_mtr_bit_malha_rodoviaria AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_transportes__mtr_bit_malha_rodoviaria/mart.parquet');
