with silver as (
    select * from {{ read_silver_parquet('estimativa_cafe') }}
),

renamed as (
    select
        {{ conab_estimativa_cafe_columns() }},
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
            partition by produto, uf, safra, id_levantamento
            order by capturado_em desc, _source_file desc
        ) as rn
    from renamed
)
where rn = 1
