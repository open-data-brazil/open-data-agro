{{ config(
    materialized='external',
    location=var('gold_root') ~ '/mart_ml__soy_daily_features/mart.parquet',
    format='parquet'
) }}

-- Daily soja Paranaguá feature panel (UC-ML-001). Grain: (produto_slug, data).
with cepea as (
    select *
    from {{ ref('int_cross__cepea_soja_paranagua_typed') }}
    where produto_slug = 'soja'
),

cepea_feat as (
    select
        produto_slug,
        data as label_date,
        data,
        preco_rs_sc,
        preco_usd_sc as preco_usd_sc_raw,
        variacao_dia_pct,
        as_of as as_of_cepea,
        lag(preco_rs_sc, 1) over (order by data) as preco_rs_sc_lag1,
        lag(preco_rs_sc, 7) over (order by data) as preco_rs_sc_lag7,
        stddev_samp(preco_rs_sc) over (
            order by data
            rows between 6 preceding and current row
        ) as cepea_vol_7d
    from cepea
),

ptax as (
    select data, ptax_usd_venda, as_of
    from {{ ref('int_cross__bcb_ptax_usd_venda_typed') }}
),

b3_front as (
    select
        refdate,
        price as b3_front_price,
        symbol as b3_front_symbol,
        as_of
    from (
        select
            *,
            row_number() over (
                partition by refdate
                order by maturity_code asc, symbol asc
            ) as rn
        from {{ ref('int_cross__b3_futuro_soja_typed') }}
    ) t
    where rn = 1
),

comex_br as (
    select
        refmonth,
        as_of,
        sum(valor_fob_usd) as comex_soja_fob_usd,
        sum(quantidade_kg) as comex_soja_kg
    from {{ ref('int_cross__mdic_comex_exportacao_uf_ncm_typed') }}
    where produto_slug = 'soja'
    group by 1, 2
),

frete_mt_pr as (
    select
        refmonth,
        as_of,
        valor_frete_ton_avg as frete_mt_pr_ton_avg,
        valor_tonelada_km_avg as frete_mt_pr_tkm_avg
    from {{ ref('int_cross__frete_uf_par') }}
    where uf_origem = 'MT'
      and uf_destino = 'PR'
),

inmet_national as (
    select
        data as obs_date,
        as_of,
        avg(valor) as inmet_valor_avg
    from {{ ref('int_cross__inmet_bdmep_diario_typed') }}
    group by 1, 2
)

select
    c.produto_slug,
    c.data,
    c.label_date,
    c.preco_rs_sc,
    coalesce(
        c.preco_usd_sc_raw,
        case
            when p.ptax_usd_venda is not null and p.ptax_usd_venda <> 0
                then c.preco_rs_sc / p.ptax_usd_venda
        end
    ) as preco_usd_sc,
    c.variacao_dia_pct,
    case
        when c.preco_rs_sc_lag1 is not null and c.preco_rs_sc_lag1 <> 0
            then (c.preco_rs_sc / c.preco_rs_sc_lag1) - 1
    end as cepea_ret_1d,
    case
        when c.preco_rs_sc_lag7 is not null and c.preco_rs_sc_lag7 <> 0
            then (c.preco_rs_sc / c.preco_rs_sc_lag7) - 1
    end as cepea_ret_7d,
    c.cepea_vol_7d,
    p.ptax_usd_venda,
    b.b3_front_price,
    b.b3_front_symbol,
    case
        when b.b3_front_price is not null and p.ptax_usd_venda is not null and p.ptax_usd_venda <> 0
            then c.preco_rs_sc / p.ptax_usd_venda - b.b3_front_price
    end as basis_cepea_usd_minus_b3,
    cx.comex_soja_fob_usd,
    cx.comex_soja_kg,
    f.frete_mt_pr_ton_avg,
    f.frete_mt_pr_tkm_avg,
    i.inmet_valor_avg as inmet_national_avg,
    c.as_of_cepea,
    p.as_of as as_of_ptax,
    b.as_of as as_of_b3,
    cx.as_of as as_of_comex,
    f.as_of as as_of_frete,
    i.as_of as as_of_inmet,
    greatest(
        c.as_of_cepea,
        coalesce(p.as_of, date '1900-01-01'),
        coalesce(b.as_of, date '1900-01-01'),
        coalesce(cx.as_of, date '1900-01-01'),
        coalesce(f.as_of, date '1900-01-01'),
        coalesce(i.as_of, date '1900-01-01')
    ) as feature_as_of_max
from cepea_feat c
asof left join ptax p
    on c.data >= p.data
asof left join b3_front b
    on c.data >= b.refdate
asof left join comex_br cx
    on c.data >= cx.as_of
asof left join frete_mt_pr f
    on c.data >= f.as_of
asof left join inmet_national i
    on c.data >= i.as_of
