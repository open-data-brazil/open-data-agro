{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_cross__soja_features_monthly/mart.parquet',
    format='parquet'
) }}

-- Monthly municipal soja panel (CONAB grain). Analytics only — not an official source.
-- Grain: (produto_slug, cod_ibge, refmonth). label_date = CONAB publication as_of.
with local_soja as (
    select
        produto_slug,
        cod_ibge,
        nome_municipio,
        uf,
        codigo_regiao,
        nome_regiao,
        refmonth,
        ano,
        mes,
        avg(valor_produto_kg) as preco_local_kg,
        max(as_of) as as_of_local,
        max(as_of) as label_date
    from {{ ref('int_cross__precos_locais_municipio') }}
    where produto_slug = 'soja'
      and cod_ibge is not null
      and refmonth is not null
    group by 1, 2, 3, 4, 5, 6, 7, 8, 9
),

cepea_m as (
    select
        date_trunc('month', data)::date as refmonth,
        avg(preco_rs_sc) as cepea_preco_rs_sc_avg,
        avg(preco_usd_sc) as cepea_preco_usd_sc_avg,
        max(as_of) as as_of_cepea
    from {{ ref('int_cross__cepea_soja_paranagua_typed') }}
    where produto_slug = 'soja'
    group by 1
),

ptax_m as (
    select
        date_trunc('month', data)::date as refmonth,
        avg(ptax_usd_venda) as ptax_usd_venda_avg,
        max(as_of) as as_of_ptax
    from {{ ref('int_cross__bcb_ptax_usd_venda_typed') }}
    group by 1
),

comex_uf as (
    select
        uf,
        refmonth,
        as_of as as_of_comex,
        sum(valor_fob_usd) as comex_soja_fob_usd,
        sum(quantidade_kg) as comex_soja_kg
    from {{ ref('int_cross__mdic_comex_exportacao_uf_ncm_typed') }}
    where produto_slug = 'soja'
    group by 1, 2, 3
),

frete_to_pr as (
    select
        uf_origem as uf,
        refmonth,
        as_of as as_of_frete,
        valor_frete_ton_avg as frete_uf_pr_ton_avg,
        valor_tonelada_km_avg as frete_uf_pr_tkm_avg
    from {{ ref('int_cross__frete_uf_par') }}
    where uf_destino = 'PR'
)

select
    l.produto_slug,
    l.cod_ibge,
    l.nome_municipio,
    l.uf,
    l.codigo_regiao,
    l.nome_regiao,
    l.refmonth,
    l.ano,
    l.mes,
    l.label_date,
    l.preco_local_kg,
    c.cepea_preco_rs_sc_avg,
    c.cepea_preco_usd_sc_avg,
    p.ptax_usd_venda_avg,
    case
        when c.cepea_preco_rs_sc_avg is not null
            and p.ptax_usd_venda_avg is not null
            and p.ptax_usd_venda_avg <> 0
            and l.preco_local_kg is not null
            then (l.preco_local_kg * 60.0) / p.ptax_usd_venda_avg
                 - c.cepea_preco_rs_sc_avg / p.ptax_usd_venda_avg
    end as basis_local_sc_usd_minus_cepea,
    cx.comex_soja_fob_usd,
    cx.comex_soja_kg,
    f.frete_uf_pr_ton_avg,
    f.frete_uf_pr_tkm_avg,
    l.as_of_local,
    c.as_of_cepea,
    p.as_of_ptax,
    cx.as_of_comex,
    f.as_of_frete,
    greatest(
        coalesce(c.as_of_cepea, date '1900-01-01'),
        coalesce(p.as_of_ptax, date '1900-01-01'),
        coalesce(cx.as_of_comex, date '1900-01-01'),
        coalesce(f.as_of_frete, date '1900-01-01')
    ) as feature_as_of_max
from local_soja l
asof left join cepea_m c
    on l.label_date >= c.as_of_cepea
asof left join ptax_m p
    on l.label_date >= p.as_of_ptax
asof left join comex_uf cx
    on l.uf = cx.uf
   and l.label_date >= cx.as_of_comex
asof left join frete_to_pr f
    on l.uf = f.uf
   and l.label_date >= f.as_of_frete
