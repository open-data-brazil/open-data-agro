with silver as (
    select * from {{ read_silver_parquet('alimenta_brasil_entregas') }}
),

renamed as (
    select
        {{ conab_alimenta_brasil_entregas_columns() }},
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
            partition by ano_entrega, mes_entrega, municipio, uf, sexo
            order by capturado_em desc, _source_file desc
        ) as rn
    from renamed
)
where rn = 1
