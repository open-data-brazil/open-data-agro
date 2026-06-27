with silver as (
    select * from {{ read_silver_parquet('siap_produccion_agricola', agency='mexico') }}
),

renamed as (
    select
        {{ mexico_siap_produccion_agricola_columns() }},
        cast(_ingested_at as varchar) as capturado_em,
        '{{ mexico_siap_produccion_agricola_fonte_oficial() }}' as fonte_oficial,
        _dataset_id,
        _source_file
    from silver
)

select * exclude (rn)
from (
    select
        *,
        row_number() over (
            partition by state_code, crop_slug, refyear
            order by capturado_em desc, _source_file desc
        ) as rn
    from renamed
)
where rn = 1
