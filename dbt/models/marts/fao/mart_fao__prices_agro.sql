{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_fao__prices_agro/mart.parquet',
    format='parquet'
) }}

select * from {{ ref('stg_fao__prices_agro') }}
