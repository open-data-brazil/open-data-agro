{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_transportes__mtr_bit_malha_rodoviaria/mart.parquet',
    format='parquet'
) }}

select * from {{ ref('stg_transportes__mtr_bit_malha_rodoviaria') }}
