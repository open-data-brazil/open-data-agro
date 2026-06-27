{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_mexico__siap_produccion_agricola/mart.parquet',
    format='parquet'
) }}

select *
from {{ ref('stg_mexico__siap_produccion_agricola') }}
