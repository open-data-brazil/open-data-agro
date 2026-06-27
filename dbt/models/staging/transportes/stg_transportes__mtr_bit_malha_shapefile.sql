with silver as (
    select * from {{ read_silver_parquet('mtr_bit_malha_shapefile', agency='transportes') }}
),

renamed as (
    select
        *,
        cast(_ingested_at as varchar) as capturado_em,
        'https://www.gov.br/transportes/pt-br/assuntos/dados-de-transportes/bit/bit-mapas' as fonte_oficial,
        _dataset_id,
        _source_file
    from silver
)

select * exclude (rn)
from (
    select
        *,
        row_number() over (
            partition by _dataset_id, objectid_1
            order by capturado_em desc, _source_file desc
        ) as rn
    from renamed
)
where rn = 1
