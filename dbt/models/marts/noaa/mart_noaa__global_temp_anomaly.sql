{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_noaa__global_temp_anomaly/mart.parquet',
    format='parquet'
) }}

select * from {{ ref('stg_noaa__global_temp_anomaly') }}
