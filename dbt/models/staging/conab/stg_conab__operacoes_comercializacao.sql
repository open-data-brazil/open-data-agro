with silver as (
    select * from {{ read_silver_parquet('operacoes_comercializacao') }}
),

renamed as (
    select
        {{ conab_operacoes_comercializacao_columns() }},
        cast(_ingested_at as varchar) as capturado_em,
        '{{ conab_fonte_oficial() }}' as fonte_oficial,
        _dataset_id,
        _source_file
    from silver
)

select * exclude (rn)
from (
    select
        *,
        row_number() over (
            partition by
                id_edital,
                num_lote,
                num_dco,
                produto,
                id_produto,
                operacao,
                ano_edital,
                mes_edital
            order by capturado_em desc, _source_file desc
        ) as rn
    from renamed
)
where rn = 1
