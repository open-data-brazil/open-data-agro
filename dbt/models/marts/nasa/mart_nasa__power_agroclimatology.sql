{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_nasa__power_agroclimatology/mart.parquet',
    format='parquet'
) }}

select *
from {{ ref('stg_nasa__power_agroclimatology') }}
