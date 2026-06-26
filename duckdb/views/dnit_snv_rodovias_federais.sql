-- Published view: DNIT SNV federal highway jurisdiction (dbt gold mart).

CREATE OR REPLACE VIEW analytics.dnit_snv_rodovias_federais AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_dnit__snv_rodovias_federais/mart.parquet');
