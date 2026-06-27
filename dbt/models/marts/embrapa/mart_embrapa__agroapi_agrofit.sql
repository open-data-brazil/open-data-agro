{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_embrapa__agroapi_agrofit/mart.parquet',
    format='parquet'
) }}

select * from {{ ref('stg_embrapa__agroapi_agrofit') }}
