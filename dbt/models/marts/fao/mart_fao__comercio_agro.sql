{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_fao__comercio_agro/mart.parquet',
    format='parquet'
) }}

select
    area_code,
    area_name,
    item_code,
    item_name,
    commodity_slug,
    element_code,
    element_name,
    year,
    unit,
    value,
    flag,
    capturado_em,
    fonte_oficial,
    _dataset_id,
    _source_file
from {{ ref('stg_fao__comercio_agro') }}
