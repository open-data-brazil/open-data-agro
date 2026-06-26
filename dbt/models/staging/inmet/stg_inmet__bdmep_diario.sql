with silver as (
    select * from {{ read_silver_parquet('bdmep_diario', agency='inmet') }}
),

renamed as (
    select
        {{ inmet_bdmep_diario_columns() }},
        cast(_ingested_at as varchar) as capturado_em,
        '{{ inmet_bdmep_fonte_oficial() }}' as fonte_oficial,
        _dataset_id,
        _source_file
    from silver
)

select * exclude (rn)
from (
    select
        *,
        row_number() over (
            partition by cd_estacao, data, variavel
            order by capturado_em desc, _source_file desc
        ) as rn
    from renamed
)
where rn = 1
