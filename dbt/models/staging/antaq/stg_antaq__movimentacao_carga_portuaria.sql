with silver as (
    select * from {{ read_silver_parquet('movimentacao_carga_portuaria', agency='antaq') }}
),

renamed as (
    select
        Ano as ano,
        Mes as mes,
        CodigoInstalacaoPortuaria as codigo_instalacao_portuaria,
        NomeInstalacaoPortuaria as nome_instalacao_portuaria,
        TipoMovimentacao as tipo_movimentacao,
        TipoNavegacao as tipo_navegacao,
        Sentido as sentido,
        NaturezaCarga as natureza_carga,
        TipoOperacao as tipo_operacao,
        PesoToneladas as peso_toneladas,
        cast(_ingested_at as varchar) as capturado_em,
        '{{ antaq_fonte_oficial() }}' as fonte_oficial,
        _dataset_id,
        _source_file
    from silver
)

select * exclude (rn)
from (
    select
        *,
        row_number() over (
            partition by ano, mes, codigo_instalacao_portuaria, tipo_movimentacao, tipo_navegacao, sentido, natureza_carga, tipo_operacao
            order by capturado_em desc, _source_file desc
        ) as rn
    from renamed
)
where rn = 1
