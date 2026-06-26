with silver as (
    select * from {{ read_silver_parquet('ag_prices', agency='eurostat') }}
),

renamed as (
    select
        dataset_code,
        geo,
        product_code,
        product_name,
        year,
        index_value,
        base_period,
        cast(_ingested_at as varchar) as capturado_em,
        '{{ eurostat_fonte_oficial() }}' as fonte_oficial,
        _dataset_id,
        _source_file
    from silver
)

select * exclude (rn)
from (
    select
        *,
        row_number() over (
            partition by dataset_code, geo, product_code, year
            order by capturado_em desc, _source_file desc
        ) as rn
    from renamed
)
where rn = 1
