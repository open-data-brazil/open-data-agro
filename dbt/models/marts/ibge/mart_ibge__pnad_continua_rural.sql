{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_ibge__pnad_continua_rural/mart.parquet',
    format='parquet'
) }}

select * from {{ ref('stg_ibge__pnad_continua_rural') }}
