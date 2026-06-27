{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_ibge__pnad_rural_renda_ocupacao/mart.parquet',
    format='parquet'
) }}

select
    *
from {{ ref('stg_ibge__pnad_rural_renda_ocupacao') }}
