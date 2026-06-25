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

{% macro conab_oferta_demanda_columns() %}
  trim(produto) as produto,
  trim(dsc_safra) as safra,
  cast(id_produto as varchar) as id_produto,
  cast(estoque_inicial_1000t as varchar) as estoque_inicial_1000t,
  cast(producao_1000t as varchar) as producao_1000t,
  cast(importacao_1000t as varchar) as importacao_1000t,
  cast(consumo_1000t as varchar) as consumo_1000t,
  cast(exportacao_1000t as varchar) as exportacao_1000t,
  cast(estoque_final_1000t as varchar) as estoque_final_1000t
{% endmacro %}

{% macro conab_estoques_publicos_columns() %}
  trim(produto) as produto,
  cast(id_produto as varchar) as id_produto,
  trim(nom_municipio) as municipio,
  cast(cod_ibge as varchar) as cod_ibge,
  trim(uf) as uf,
  cast(num_ano as varchar) as ano,
  cast(num_mes as varchar) as mes,
  trim(conta_operacional) as conta_operacional,
  cast(qtd_estoque_kg as varchar) as qtd_estoque_kg
{% endmacro %}
