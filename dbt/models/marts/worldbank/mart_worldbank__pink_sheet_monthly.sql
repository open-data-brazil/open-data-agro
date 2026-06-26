{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_worldbank__pink_sheet_monthly/mart.parquet',
    format='parquet'
) }}

select * from {{ ref('stg_worldbank__pink_sheet_monthly') }}
