{% macro inmet_fonte_oficial() %}
  https://portal.inmet.gov.br/
{% endmacro %}

{% macro inmet_bdmep_fonte_oficial() %}
  https://bdmep.inmet.gov.br/
{% endmacro %}

{% macro inmet_estacoes_automaticas_columns() %}
  trim(cd_estacao) as cd_estacao,
  trim(nome) as nome,
  cast(latitude as varchar) as latitude,
  cast(longitude as varchar) as longitude,
  trim(uf) as uf,
  trim(situacao) as situacao
{% endmacro %}

{% macro inmet_estacoes_convencionais_columns() %}
  trim(cd_estacao) as cd_estacao,
  trim(nome) as nome,
  cast(latitude as varchar) as latitude,
  cast(longitude as varchar) as longitude,
  trim(uf) as uf,
  trim(situacao) as situacao,
  trim(regiao) as regiao,
  cast(altitude as varchar) as altitude
{% endmacro %}

{% macro inmet_bdmep_diario_columns() %}
  trim(cd_estacao) as cd_estacao,
  trim(data) as data,
  trim(variavel) as variavel,
  cast(valor as varchar) as valor,
  trim(uf) as uf,
  cast(ano as varchar) as ano
{% endmacro %}

{% macro inmet_bdmep_mensal_columns() %}
  trim(cd_estacao) as cd_estacao,
  trim(mes) as mes,
  trim(variavel) as variavel,
  cast(valor as varchar) as valor,
  trim(uf) as uf,
  cast(ano as varchar) as ano
{% endmacro %}

{% macro inmet_sequia_monitor_fonte_oficial() %}
  https://www.gov.br/ana/pt-br/servicos/monitor-de-seca
{% endmacro %}

{% macro inmet_sequia_monitor_columns() %}
    mapa_id,
    ano,
    mes,
    categoria_seca,
    area_km2,
    area_id,
    tipo_area
{% endmacro %}
