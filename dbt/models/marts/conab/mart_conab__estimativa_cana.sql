{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_conab__estimativa_cana/mart.parquet',
    format='parquet'
) }}

select
    produto,
    uf,
    safra,
    safra_previsao,
    id_produto,
    levantamento,
    id_levantamento,
    area_plantada_mil_ha,
    producao_mil_t,
    producao_acucar_mil_t,
    producao_etanol_anidro_mil_l,
    producao_etanol_hidratado_mil_l,
    producao_etanol_total_mil_l,
    producao_atr_kg_t,
    capturado_em,
    fonte_oficial,
    _dataset_id,
    _source_file
from {{ ref('stg_conab__estimativa_cana') }}
