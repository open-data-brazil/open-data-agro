{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_dnit__snv_rodovias_federais/mart.parquet',
    format='parquet'
) }}

select
    br,
    uf,
    tipo_trecho,
    codigo,
    local_inicio,
    local_fim,
    km_inicial,
    km_final,
    extensao,
    superficie_federal,
    obras,
    federal_coincidente,
    administracao,
    ato_legal,
    estadual_coincidente,
    superficie_est_coincidente,
    jurisdicao,
    superficie,
    unidade_local,
    capturado_em,
    fonte_oficial,
    _dataset_id,
    _source_file
from {{ ref('stg_dnit__snv_rodovias_federais') }}
