{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_abiove__capacidade_instalada_esmagamento/mart.parquet',
    format='parquet'
) }}

select * from {{ ref('stg_abiove__capacidade_instalada_esmagamento') }}
