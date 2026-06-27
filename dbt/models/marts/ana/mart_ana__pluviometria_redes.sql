{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_ana__pluviometria_redes/mart.parquet',
    format='parquet'
) }}

select * from {{ ref('stg_ana__pluviometria_redes') }}
