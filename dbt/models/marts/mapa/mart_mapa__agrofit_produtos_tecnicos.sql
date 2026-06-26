{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_mapa__agrofit_produtos_tecnicos/mart.parquet',
    format='parquet'
) }}

select
    numero_registro,
    produto_tecnico_marca_comercial,
    ingrediente_ativo,
    classe,
    titular_registro,
    empresa_pais_tipo,
    classificacao_toxicologica,
    classificacao_ambiental,
    capturado_em,
    fonte_oficial,
    _dataset_id,
    _source_file
from {{ ref('stg_mapa__agrofit_produtos_tecnicos') }}
