{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_ana__hidrologia_series/mart.parquet',
    format='parquet'
) }}

select
    station_code,
    consistency_level,
    data_type,
    observed_at,
    daily_mean,
    acquisition_method,
    max_value,
    min_value,
    mean_value,
    capturado_em,
    fonte_oficial,
    _dataset_id,
    _source_file
from {{ ref('stg_ana__hidrologia_series') }}
