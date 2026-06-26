{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_conab__custo_producao/mart.parquet',
    format='parquet'
) }}

select
    empreendimento,
    ano,
    mes,
    ano_mes,
    produto,
    id_produto,
    safra,
    uf,
    municipio,
    cod_ibge,
    unidade_comercializacao,
    vlr_custo_variavel_ha,
    vlr_custo_variavel_unidade,
    vlr_custo_fixo_ha,
    vlr_custo_fixo_unidade,
    vlr_renda_fator_ha,
    vlr_renda_fator_unidade,
    capturado_em,
    fonte_oficial,
    _dataset_id,
    _source_file
from {{ ref('stg_conab__custo_producao') }}
