{% macro conab_silver_path(table_name) %}
  {{ var('lake_root') }}/silver/conab/{{ table_name }}
{% endmacro %}

{% macro conab_fonte_oficial() %}
  {{ return(var('conab_fonte_oficial')) }}
{% endmacro %}

{% macro conab_graos_core_columns() %}
  trim(produto) as produto,
  trim(uf) as uf,
  trim(coalesce(nullif(trim(safra), ''), trim(ano_agricola))) as safra,
  cast(area_plantada_mil_ha as varchar) as area_plantada_mil_ha,
  cast(producao_mil_t as varchar) as producao_mil_t
{% endmacro %}

{% macro conab_serie_core_columns() %}
  trim(produto) as produto,
  trim(uf) as uf,
  trim(ano_agricola) as ano,
  cast(area_plantada_mil_ha as varchar) as area_plantada_mil_ha,
  cast(producao_mil_t as varchar) as producao_mil_t
{% endmacro %}
