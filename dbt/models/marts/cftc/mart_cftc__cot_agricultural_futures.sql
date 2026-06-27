{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_cftc__cot_agricultural_futures/mart.parquet',
    format='parquet'
) }}

select *
from {{ ref('stg_cftc__cot_agricultural_futures') }}
