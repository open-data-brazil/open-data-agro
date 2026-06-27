{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_fred__commodity_indexes/mart.parquet',
    format='parquet'
) }}

select *
from {{ ref('stg_fred__commodity_indexes') }}
