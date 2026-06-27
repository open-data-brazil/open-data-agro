{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_ibge__censo_agro_estabelecimentos/mart.parquet',
    format='parquet'
) }}

select * from {{ ref('stg_ibge__censo_agro_estabelecimentos') }}
