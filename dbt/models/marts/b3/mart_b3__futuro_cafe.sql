{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_b3__futuro_cafe/mart.parquet',
    format='parquet'
) }}

select * from {{ ref('stg_b3__futuro_cafe') }}
