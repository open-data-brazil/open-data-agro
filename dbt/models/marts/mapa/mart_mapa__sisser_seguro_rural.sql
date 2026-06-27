{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_mapa__sisser_seguro_rural/mart.parquet',
    format='parquet'
) }}

select
    *
from {{ ref('stg_mapa__sisser_seguro_rural') }}
