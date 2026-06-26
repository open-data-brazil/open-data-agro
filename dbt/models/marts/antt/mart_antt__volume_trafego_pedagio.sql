{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_antt__volume_trafego_pedagio/mart.parquet',
    format='parquet'
) }}

select
    concessionaria,
    mes_ano,
    sentido,
    praca,
    tipo_cobranca,
    categoria_eixo,
    tipo_de_veiculo,
    volume_total,
    capturado_em,
    fonte_oficial,
    _dataset_id,
    _source_file
from {{ ref('stg_antt__volume_trafego_pedagio') }}
