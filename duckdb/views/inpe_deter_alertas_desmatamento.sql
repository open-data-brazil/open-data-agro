-- Published view: inpe.deter-alertas-desmatamento (dbt gold mart).

CREATE OR REPLACE VIEW analytics.inpe_deter_alertas_desmatamento AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_inpe__deter_alertas_desmatamento/mart.parquet');
