{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_ibama__sisfogo_incendios/mart.parquet',
    format='parquet'
) }}

select * from {{ ref('stg_ibama__sisfogo_incendios') }}
