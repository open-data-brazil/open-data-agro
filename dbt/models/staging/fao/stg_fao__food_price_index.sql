with silver as (
    select * from {{ read_silver_parquet('food_price_index', agency='fao') }}
),

renamed as (
    select
        {{ fao_ffpi_columns() }},
        cast(_ingested_at as varchar) as capturado_em,
        '{{ fao_ffpi_fonte_oficial() }}' as fonte_oficial,
        _dataset_id,
        _source_file
    from silver
)

select * exclude (rn)
from (
    select
        *,
        row_number() over (
            partition by index_slug, refmonth
            order by capturado_em desc, _source_file desc
        ) as rn
    from renamed
)
where rn = 1
