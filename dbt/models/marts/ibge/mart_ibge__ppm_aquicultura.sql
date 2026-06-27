{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_ibge__ppm_aquicultura/mart.parquet',
    format='parquet'
) }}

select
    *
from {{ ref('stg_ibge__ppm_aquicultura') }}
