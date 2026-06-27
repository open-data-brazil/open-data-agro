{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_sagis__grain_supply_statistics/mart.parquet',
    format='parquet'
) }}

select *
from {{ ref('stg_sagis__grain_supply_statistics') }}
