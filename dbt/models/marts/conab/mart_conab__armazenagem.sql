{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_conab__armazenagem/mart.parquet',
    format='parquet'
) }}

select
    identificacao_armazem,
    especie_armazem,
    tipo_armazem,
    tipo_entidade,
    tipo_pessoa,
    municipio,
    cod_ibge,
    uf,
    qtd_capacidade_estatica_t,
    qtd_capacidade_expedicao_t,
    qtd_capacidade_recepcao_t,
    latitude,
    longitude,
    nome_armazenador,
    endereco,
    email,
    capturado_em,
    fonte_oficial,
    _dataset_id,
    _source_file
from {{ ref('stg_conab__armazenagem') }}
