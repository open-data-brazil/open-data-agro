{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_inpe__deter_alertas_desmatamento/mart.parquet',
    format='parquet'
) }}

select * from {{ ref('stg_inpe__deter_alertas_desmatamento') }}
