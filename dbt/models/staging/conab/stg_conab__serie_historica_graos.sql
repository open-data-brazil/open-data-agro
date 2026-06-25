with silver as (
    select * from {{ read_silver_parquet('serie_historica_graos') }}
),

renamed as (
    select
        {{ conab_serie_core_columns() }},
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
            partition by produto, uf, ano
            order by capturado_em desc, _source_file desc
        ) as rn
    from renamed
)
where rn = 1
