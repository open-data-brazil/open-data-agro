-- Published view: Argentina MAGyP grain production (dbt gold mart).
CREATE OR REPLACE VIEW analytics.argentina_magyp_producion_granos AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_argentina__magyp_producion_granos/mart.parquet');
