{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_usda__psd_soja/mart.parquet',
    format='parquet'
) }}

select * from {{ ref('stg_usda__psd_soja') }}
