{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_usda__psd_trigo/mart.parquet',
    format='parquet'
) }}

select * from {{ ref('stg_usda__psd_trigo') }}
