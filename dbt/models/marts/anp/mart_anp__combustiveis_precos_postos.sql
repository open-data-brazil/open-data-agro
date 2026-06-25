{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_anp__combustiveis_precos_postos/mart.parquet',
    format='parquet'
) }}

select
    cnpj,
    razao,
    fantasia,
    endereco,
    numero,
    complemento,
    bairro,
    cep,
    municipio,
    estado,
    bandeira,
    produto,
    unidade_medida,
    preco_revenda,
    data_coleta,
    capturado_em,
    fonte_oficial,
    _dataset_id,
    _source_file
from {{ ref('stg_anp__combustiveis_precos_postos') }}
