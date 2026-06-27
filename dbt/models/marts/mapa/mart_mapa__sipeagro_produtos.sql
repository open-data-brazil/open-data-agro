{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_mapa__sipeagro_produtos/mart.parquet',
    format='parquet'
) }}

select
    *
from {{ ref('stg_mapa__sipeagro_produtos') }}
