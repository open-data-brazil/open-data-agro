{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_worldbank__ag_indices/mart.parquet',
    format='parquet'
) }}

select
    refmonth,
    series_name,
    commodity_slug,
    unit,
    value,
    capturado_em,
    fonte_oficial,
    _dataset_id,
    _source_file
from {{ ref('stg_worldbank__ag_indices') }}
