with silver as (
    select * from {{ read_silver_parquet('magyp_producion_granos', agency='argentina') }}
),

renamed as (
    select
        {{ argentina_granos_columns() }},
        cast(_ingested_at as varchar) as capturado_em,
        '{{ argentina_granos_fonte_oficial() }}' as fonte_oficial,
        _dataset_id,
        _source_file
    from silver
)

select * exclude (rn)
from (
    select
        *,
        row_number() over (
            partition by commodity_slug, refyear
            order by capturado_em desc, _source_file desc
        ) as rn
    from renamed
)
where rn = 1
