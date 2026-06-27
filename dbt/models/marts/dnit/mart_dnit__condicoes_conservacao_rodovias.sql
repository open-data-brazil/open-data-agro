{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_dnit__condicoes_conservacao_rodovias/mart.parquet',
    format='parquet'
) }}

select * from {{ ref('stg_dnit__condicoes_conservacao_rodovias') }}
