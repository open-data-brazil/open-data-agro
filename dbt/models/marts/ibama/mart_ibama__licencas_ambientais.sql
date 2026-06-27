{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_ibama__licencas_ambientais/mart.parquet',
    format='parquet'
) }}

select * from {{ ref('stg_ibama__licencas_ambientais') }}
