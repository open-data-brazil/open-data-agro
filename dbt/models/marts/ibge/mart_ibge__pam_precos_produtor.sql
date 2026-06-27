{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_ibge__pam_precos_produtor/mart.parquet',
    format='parquet'
) }}

select
    *
from {{ ref('stg_ibge__pam_precos_produtor') }}
