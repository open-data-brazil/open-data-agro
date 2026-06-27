{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_mapa__sigef_producao_sementes/mart.parquet',
    format='parquet'
) }}

select
    *
from {{ ref('stg_mapa__sigef_producao_sementes') }}
