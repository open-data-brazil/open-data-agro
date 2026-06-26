{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_conab__precos_minimos/mart.parquet',
    format='parquet'
) }}

select
    produto,
    id_produto,
    uf,
    regionalizacao,
    ano_inicio_vigencia,
    mes_inicio_vigencia,
    ano_termino_vigencia,
    mes_termino_vigencia,
    preco,
    unidade_comercializacao,
    nome_normativo,
    url_normativo,
    capturado_em,
    fonte_oficial,
    _dataset_id,
    _source_file
from {{ ref('stg_conab__precos_minimos') }}
