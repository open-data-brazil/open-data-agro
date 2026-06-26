{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_usda__psd_milho/mart.parquet',
    format='parquet'
) }}

select * from {{ ref('stg_usda__psd_milho') }}
