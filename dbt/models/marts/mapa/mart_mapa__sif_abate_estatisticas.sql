{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_mapa__sif_abate_estatisticas/mart.parquet',
    format='parquet'
) }}

select * from {{ ref('stg_mapa__sif_abate_estatisticas') }}
