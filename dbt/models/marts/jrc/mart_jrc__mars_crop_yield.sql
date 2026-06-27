{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_jrc__mars_crop_yield/mart.parquet',
    format='parquet'
) }}

select *
from {{ ref('stg_jrc__mars_crop_yield') }}
