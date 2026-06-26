with silver as (
    select * from {{ read_silver_parquet('financiamento_agro', agency='bndes') }}
),

renamed as (
    select
        {{ bndes_financiamento_agro_columns() }},
        cast(_ingested_at as varchar) as capturado_em,
        '{{ bndes_fonte_oficial() }}' as fonte_oficial,
        _dataset_id,
        _source_file
    from silver
)

select * exclude (rn)
from (
    select
        *,
        row_number() over (
            partition by ano, mes
            order by capturado_em desc, _source_file desc
        ) as rn
    from renamed
)
where rn = 1
