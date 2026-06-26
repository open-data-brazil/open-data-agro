with silver as (
    select * from {{ read_silver_parquet('producao_agro', agency='fao') }}
),

renamed as (
    select
        {{ fao_annual_bulk_columns() }},
        cast(_ingested_at as varchar) as capturado_em,
        '{{ fao_producao_fonte_oficial() }}' as fonte_oficial,
        _dataset_id,
        _source_file
    from silver
)

select * exclude (rn)
from (
    select
        *,
        row_number() over (
            partition by area_code, item_code, element_code, year
            order by capturado_em desc, _source_file desc
        ) as rn
    from renamed
)
where rn = 1
