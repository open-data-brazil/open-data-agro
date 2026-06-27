{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_anp__etanol_precos/mart.parquet',
    format='parquet'
) }}

select
    data_inicial,
    data_final,
    estado,
    municipio,
    produto,
    numero_postos_pesquisados,
    unidade_medida,
    preco_medio_revenda,
    desvio_padrao_revenda,
    preco_minimo_revenda,
    preco_maximo_revenda,
    coef_variacao_revenda,
    capturado_em,
    fonte_oficial,
    _dataset_id,
    _source_file
from {{ ref('stg_anp__etanol_precos') }}
