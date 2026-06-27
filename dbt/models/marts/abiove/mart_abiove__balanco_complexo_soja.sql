{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_abiove__balanco_complexo_soja/mart.parquet',
    format='parquet'
) }}

select * from {{ ref('stg_abiove__balanco_complexo_soja') }}
