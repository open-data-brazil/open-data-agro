{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_japan__maff_ag_trade/mart.parquet',
    format='parquet'
) }}

select *
from {{ ref('stg_japan__maff_ag_trade') }}
