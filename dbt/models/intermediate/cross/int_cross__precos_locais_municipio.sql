{{ config(materialized='view') }}

-- CONAB municipal prices joined to IBGE municipality dim (soja focus retained as filter downstream).
select
    p.produto_slug,
    p.produto,
    p.classificacao_produto,
    p.id_produto,
    p.cod_ibge,
    coalesce(d.nome_municipio, p.municipio) as nome_municipio,
    coalesce(d.uf, p.uf) as uf,
    d.codigo_regiao,
    d.nome_regiao,
    p.refmonth,
    p.ano,
    p.mes,
    p.nivel_comercializacao,
    p.valor_produto_kg,
    p.as_of,
    p.capturado_em,
    p.fonte_oficial,
    p._dataset_id
from {{ ref('int_cross__conab_precos_mensal_municipio_typed') }} p
left join {{ ref('dim_municipio') }} d
    on p.cod_ibge = d.codigo_ibge
