{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_ibge__pam_culturas_estendidas/mart.parquet',
    format='parquet'
) }}

select
    *
from {{ ref('stg_ibge__pam_culturas_estendidas') }}
