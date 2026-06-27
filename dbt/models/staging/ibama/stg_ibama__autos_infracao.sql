with silver as (
    select * from {{ read_silver_parquet('autos_infracao', agency='ibama') }}
),

renamed as (
    select
        *,
        cast(_ingested_at as varchar) as capturado_em,
        'https://dadosabertos.ibama.gov.br/' as fonte_oficial,
        _dataset_id,
        _source_file
    from silver
)

select * exclude (rn)
from (
    select
        *,
        row_number() over (
            partition by _dataset_id, SEQ_AUTO_INFRACAO
            order by capturado_em desc, _source_file desc
        ) as rn
    from renamed
)
where rn = 1
