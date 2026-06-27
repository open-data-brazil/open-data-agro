{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_copernicus__era5_agroclimate/mart.parquet',
    format='parquet'
) }}

select *
from {{ ref('stg_copernicus__era5_agroclimate') }}
