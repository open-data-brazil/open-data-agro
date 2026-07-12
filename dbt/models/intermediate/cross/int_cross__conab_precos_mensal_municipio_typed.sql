{{ config(materialized='view') }}

select
    trim(produto) as produto,
    case
        when upper(trim(produto)) = 'SOJA' then 'soja'
        else lower(trim(produto))
    end as produto_slug,
    trim(classificacao_produto) as classificacao_produto,
    trim(id_produto) as id_produto,
    trim(municipio) as municipio,
    lpad(trim(cod_ibge), 7, '0') as cod_ibge,
    upper(trim(uf)) as uf,
    trim(regiao) as regiao,
    try_cast(ano as integer) as ano,
    try_cast(mes as integer) as mes,
    make_date(try_cast(ano as integer), try_cast(mes as integer), 1) as refmonth,
    trim(nivel_comercializacao) as nivel_comercializacao,
    {{ br_to_double('valor_produto_kg') }} as valor_produto_kg,
    (
        last_day(make_date(try_cast(ano as integer), try_cast(mes as integer), 1))
        + interval 30 day
    )::date as as_of,
    capturado_em,
    fonte_oficial,
    _dataset_id
from {{ read_gold_mart('mart_conab__precos_mensal_municipio') }}
where try_cast(ano as integer) is not null
  and try_cast(mes as integer) between 1 and 12
