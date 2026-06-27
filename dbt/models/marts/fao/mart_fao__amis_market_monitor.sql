{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_fao__amis_market_monitor/mart.parquet',
    format='parquet'
) }}

select *
from {{ ref('stg_fao__amis_market_monitor') }}
