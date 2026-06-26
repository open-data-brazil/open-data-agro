with silver as (
    select * from {{ read_silver_parquet('frete') }}
),

renamed as (
    select
        {{ conab_frete_columns() }},
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
                fonte,
                cod_ibge_origem,
                cod_ibge_destino,
                ano,
                mes
            order by capturado_em desc, _source_file desc
        ) as rn
    from renamed
)
where rn = 1
