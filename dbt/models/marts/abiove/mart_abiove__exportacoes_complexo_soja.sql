{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_abiove__exportacoes_complexo_soja/mart.parquet',
    format='parquet'
) }}

select * from {{ ref('stg_abiove__exportacoes_complexo_soja') }}
