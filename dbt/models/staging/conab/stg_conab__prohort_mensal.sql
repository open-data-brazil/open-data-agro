with silver as (
    select * from {{ read_silver_parquet('prohort_mensal') }}
),

renamed as (
    select
        {{ conab_prohort_mensal_columns() }},
        cast(_ingested_at as varchar) as capturado_em,
        '{{ conab_fonte_oficial() }}' as fonte_oficial,
        _dataset_id,
        _source_file
    from silver
)

select * exclude (rn)
from (
    select
        *,
        row_number() over (
            partition by
                ano,
                mes,
                cod_ibge_municipio_ceasa,
                ceasa,
                produto,
                cod_ibge_municipio_origem,
                pais_origem
            order by capturado_em desc, _source_file desc
        ) as rn
    from renamed
)
where rn = 1
