{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_mapa__sigef_areas/mart.parquet',
    format='parquet'
) }}

select
    *
from {{ ref('stg_mapa__sigef_areas') }}
