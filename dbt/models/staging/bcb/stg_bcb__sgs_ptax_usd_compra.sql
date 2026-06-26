with silver as (
    select * from {{ read_silver_parquet('sgs_ptax_usd_compra', agency='bcb') }}
),

renamed as (
    select
        {{ bcb_sgs_columns() }},
        cast(_ingested_at as varchar) as capturado_em,
        '{{ bcb_fonte_oficial() }}' as fonte_oficial,
        _dataset_id,
        _source_file
    from silver
)

select * exclude (rn)
from (
    select
        *,
        row_number() over (
            partition by sgs_codigo, data
            order by capturado_em desc, _source_file desc
        ) as rn
    from renamed
)
where rn = 1
