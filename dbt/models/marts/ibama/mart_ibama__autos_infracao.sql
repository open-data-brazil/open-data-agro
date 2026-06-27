{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_ibama__autos_infracao/mart.parquet',
    format='parquet'
) }}

select * from {{ ref('stg_ibama__autos_infracao') }}
