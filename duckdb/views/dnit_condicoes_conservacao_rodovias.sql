-- Published view: dnit.condicoes-conservacao-rodovias (dbt gold mart).

CREATE OR REPLACE VIEW analytics.dnit_condicoes_conservacao_rodovias AS
SELECT *
FROM read_parquet('__LAKE_ROOT__/gold/mart_dnit__condicoes_conservacao_rodovias/mart.parquet');
