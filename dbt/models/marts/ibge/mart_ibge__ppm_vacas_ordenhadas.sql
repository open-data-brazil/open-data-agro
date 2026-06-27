{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_ibge__ppm_vacas_ordenhadas/mart.parquet',
    format='parquet'
) }}

select
    *
from {{ ref('stg_ibge__ppm_vacas_ordenhadas') }}
