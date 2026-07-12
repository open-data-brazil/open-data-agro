{{ config(materialized='view') }}

select
    trim(fonte) as fonte,
    trim(municipio_origem) as municipio_origem,
    lpad(trim(cod_ibge_origem), 7, '0') as cod_ibge_origem,
    upper(trim(uf_origem)) as uf_origem,
    trim(municipio_destino) as municipio_destino,
    lpad(trim(cod_ibge_destino), 7, '0') as cod_ibge_destino,
    upper(trim(uf_destino)) as uf_destino,
    try_cast(ano as integer) as ano,
    try_cast(mes as integer) as mes,
    make_date(try_cast(ano as integer), try_cast(mes as integer), 1) as refmonth,
    {{ br_to_double('distancia_km') }} as distancia_km,
    {{ br_to_double('valor_frete_tonelada') }} as valor_frete_tonelada,
    {{ br_to_double('valor_tonelada_km') }} as valor_tonelada_km,
    (
        last_day(make_date(try_cast(ano as integer), try_cast(mes as integer), 1))
        + interval 14 day
    )::date as as_of,
    capturado_em,
    fonte_oficial,
    _dataset_id
from {{ read_gold_mart('mart_conab__frete') }}
where try_cast(ano as integer) is not null
  and try_cast(mes as integer) between 1 and 12
