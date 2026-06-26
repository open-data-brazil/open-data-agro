{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_b3__futuro_boi/mart.parquet',
    format='parquet'
) }}

select * from {{ ref('stg_b3__futuro_boi') }}
