-- Published view: suframa.comercio-mercadorias-zfm (dbt gold mart).

CREATE OR REPLACE VIEW analytics.suframa_comercio_mercadorias_zfm AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_suframa__comercio_mercadorias_zfm/mart.parquet');
