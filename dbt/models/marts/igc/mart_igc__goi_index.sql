{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_igc__goi_index/mart.parquet',
    format='parquet'
) }}

select
    refdate,
    index_slug,
    index_name,
    value,
    base_period,
    frequency,
    capturado_em,
    fonte_oficial,
    _dataset_id,
    _source_file
from {{ ref('stg_igc__goi_index') }}
