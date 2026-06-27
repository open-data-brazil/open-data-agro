{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_fao__giews_crop_prospects/mart.parquet',
    format='parquet'
) }}

select *
from {{ ref('stg_fao__giews_crop_prospects') }}
