{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_ibge__ppm_ovinos_tosquiados/mart.parquet',
    format='parquet'
) }}

select
    *
from {{ ref('stg_ibge__ppm_ovinos_tosquiados') }}
