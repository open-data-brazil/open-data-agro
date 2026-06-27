{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_ons__carga_energetica/mart.parquet',
    format='parquet'
) }}

select * from {{ ref('stg_ons__carga_energetica') }}
