{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_un__comtrade_bulk/mart.parquet',
    format='parquet'
) }}

select
    reporter_code,
    reporter_desc,
    partner_code,
    partner_desc,
    flow_code,
    flow_desc,
    period,
    hs_code,
    commodity_slug,
    trade_value_usd,
    netweight_kg,
    qty,
    qty_unit_abbr,
    capturado_em,
    fonte_oficial,
    _dataset_id,
    _source_file
from {{ ref('stg_un__comtrade_bulk') }}
