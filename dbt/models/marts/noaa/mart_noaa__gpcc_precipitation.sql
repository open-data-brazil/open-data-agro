{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_noaa__gpcc_precipitation/mart.parquet',
    format='parquet'
) }}

select *
from {{ ref('stg_noaa__gpcc_precipitation') }}
