-- Published view: ibge.pnad-continua-rural (dbt gold mart).

CREATE OR REPLACE VIEW analytics.ibge_pnad_continua_rural AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_ibge__pnad_continua_rural/mart.parquet');
