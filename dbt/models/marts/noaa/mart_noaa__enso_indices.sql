{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_noaa__enso_indices/mart.parquet',
    format='parquet'
) }}

select * from {{ ref('stg_noaa__enso_indices') }}
