with silver as (
    select * from {{ read_silver_parquet('volume_trafego_pedagio', agency='antt') }}
),

renamed as (
    select
        concessionaria,
        mes_ano,
        sentido,
        praca,
        tipo_cobranca,
        categoria_eixo,
        tipo_de_veiculo,
        volume_total,
        cast(_ingested_at as varchar) as capturado_em,
        '{{ antt_volume_trafego_fonte_oficial() }}' as fonte_oficial,
        _dataset_id,
        _source_file
    from silver
)

select * from renamed
