{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_oecd__ag_outlook/mart.parquet',
    format='parquet'
) }}

select
    ref_area,
    ref_area_name,
    commodity_code,
    commodity_name,
    measure_code,
    measure_name,
    unit,
    unit_mult,
    year,
    value,
    obs_status,
    capturado_em,
    fonte_oficial,
    _dataset_id,
    _source_file
from {{ ref('stg_oecd__ag_outlook') }}
