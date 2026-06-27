{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_ibge__lspa_rendimento_medio/mart.parquet',
    format='parquet'
) }}

select
    *
from {{ ref('stg_ibge__lspa_rendimento_medio') }}
