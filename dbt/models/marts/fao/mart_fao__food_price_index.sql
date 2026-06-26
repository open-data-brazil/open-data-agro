{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_fao__food_price_index/mart.parquet',
    format='parquet'
) }}

select
    refmonth,
    index_slug,
    index_name,
    value,
    base_period,
    capturado_em,
    fonte_oficial,
    _dataset_id,
    _source_file
from {{ ref('stg_fao__food_price_index') }}
