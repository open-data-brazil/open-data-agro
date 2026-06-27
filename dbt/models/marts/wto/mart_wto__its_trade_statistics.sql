{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_wto__its_trade_statistics/mart.parquet',
    format='parquet'
) }}

select *
from {{ ref('stg_wto__its_trade_statistics') }}
