{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_suframa__comercio_mercadorias_zfm/mart.parquet',
    format='parquet'
) }}

select * from {{ ref('stg_suframa__comercio_mercadorias_zfm') }}
