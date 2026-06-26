{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_mapa__agrofit_produtos_formulados/mart.parquet',
    format='parquet'
) }}

select
    nr_registro,
    marca_comercial,
    formulacao,
    ingrediente_ativo,
    titular_de_registro,
    classe,
    modo_de_acao,
    cultura,
    praga_nome_cientifico,
    praga_nome_comum,
    empresa_pais_tipo,
    classe_toxicologica,
    classe_ambiental,
    organicos,
    situacao,
    capturado_em,
    fonte_oficial,
    _dataset_id,
    _source_file
from {{ ref('stg_mapa__agrofit_produtos_formulados') }}
