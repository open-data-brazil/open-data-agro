with silver as (
    select * from {{ read_silver_parquet('comex_exportacao_uf_ncm', agency='mdic') }}
),

renamed as (
    select
        {{ mdic_comex_export_uf_columns() }},
        cast(_ingested_at as varchar) as capturado_em,
        '{{ mdic_fonte_oficial() }}' as fonte_oficial,
        _dataset_id,
        _source_file
    from silver
)

select * exclude (rn)
from (
    select
        *,
        row_number() over (
            partition by co_ncm, uf, data
            order by capturado_em desc, _source_file desc
        ) as rn
    from renamed
)
where rn = 1
