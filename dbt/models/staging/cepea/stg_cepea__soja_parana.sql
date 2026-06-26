with silver as (
    select * from {{ read_silver_parquet('soja_parana', agency='cepea') }}
),

renamed as (
    select
        {{ cepea_indicador_columns() }},
        cast(_ingested_at as varchar) as capturado_em,
        '{{ cepea_fonte_oficial() }}' as fonte_oficial,
        _dataset_id,
        _source_file
    from silver
)

select * exclude (rn)
from (
    select
        *,
        row_number() over (
            partition by produto, praca, data
            order by capturado_em desc, _source_file desc
        ) as rn
    from renamed
)
where rn = 1
