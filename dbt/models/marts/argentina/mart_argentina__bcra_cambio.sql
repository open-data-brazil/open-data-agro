{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_argentina__bcra_cambio/mart.parquet',
    format='parquet'
) }}

select
    currency_code,
    currency_name,
    refdate,
    exchange_rate,
    rate_type,
    capturado_em,
    fonte_oficial,
    _dataset_id,
    _source_file
from {{ ref('stg_argentina__bcra_cambio') }}
