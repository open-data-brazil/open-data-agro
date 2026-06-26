with silver as (
    select * from {{ read_silver_parquet('zarc_tabua_risco', agency='mapa') }}
),

renamed as (
    select
        Nome_cultura as nome_cultura,
        SafraIni as safra_ini,
        SafraFin as safra_fim,
        Cod_Cultura as cod_cultura,
        Cod_Ciclo as cod_ciclo,
        Cod_Solo as cod_solo,
        geocodigo,
        UF as uf,
        municipio,
        Cod_Clima as cod_clima,
        Nome_Clima as nome_clima,
        Cod_Outros_Manejos as cod_outros_manejos,
        Nome_Outros_Manejos as nome_outros_manejos,
        Produtividade as produtividade,
        Cod_NM as cod_nm,
        Cod_Munic as cod_munic,
        Cod_Meso as cod_meso,
        Cod_Micro as cod_micro,
        Portaria as portaria,
        dec1, dec2, dec3, dec4, dec5, dec6, dec7, dec8, dec9, dec10,
        dec11, dec12, dec13, dec14, dec15, dec16, dec17, dec18, dec19, dec20,
        dec21, dec22, dec23, dec24, dec25, dec26, dec27, dec28, dec29, dec30,
        dec31, dec32, dec33, dec34, dec35, dec36,
        cast(_ingested_at as varchar) as capturado_em,
        '{{ mapa_fonte_oficial() }}' as fonte_oficial,
        _dataset_id,
        _source_file
    from silver
)

select * exclude (rn)
from (
    select
        *,
        row_number() over (
            partition by geocodigo, cod_cultura, cod_ciclo, cod_solo, safra_ini, safra_fim, cod_outros_manejos
            order by capturado_em desc, _source_file desc
        ) as rn
    from renamed
)
where rn = 1
