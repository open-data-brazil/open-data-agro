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

{% macro conab_armazenagem_columns() %}
  trim(identificacao_armazem) as identificacao_armazem,
  trim(dsc_especie_armazem) as especie_armazem,
  trim(dsc_tipo_armazem) as tipo_armazem,
  trim(dsc_tipo_entidade) as tipo_entidade,
  trim(dsc_tipo_pessoa) as tipo_pessoa,
  trim(nom_municipio) as municipio,
  cast(cod_ibge as varchar) as cod_ibge,
  trim(uf) as uf,
  cast("qtd_capacidade_estatica(t)" as varchar) as qtd_capacidade_estatica_t,
  cast("qtd_capacidade_expedicao(t)" as varchar) as qtd_capacidade_expedicao_t,
  cast("qtd_capacidade_recepcao(t)" as varchar) as qtd_capacidade_recepcao_t,
  cast(latitude as varchar) as latitude,
  cast(longitude as varchar) as longitude,
  trim(nome_armazenador) as nome_armazenador,
  trim(endereco) as endereco,
  trim(email) as email
{% endmacro %}

{% macro conab_alimenta_brasil_entregas_columns() %}
  cast(ano_entrega as varchar) as ano_entrega,
  cast(mes_entrega as varchar) as mes_entrega,
  trim(municipio) as municipio,
  trim(uf) as uf,
  trim(regiao) as regiao,
  trim(sexo) as sexo,
  trim(ds_unidade_medida) as unidade_medida,
  cast(qtd_entregue as varchar) as qtd_entregue,
  cast(valor_entregue as varchar) as valor_entregue
{% endmacro %}

{% macro conab_alimenta_brasil_propostas_columns() %}
  cast(ano as varchar) as ano,
  cast(mes as varchar) as mes,
  trim(municipio) as municipio,
  cast(cod_ibge as varchar) as cod_ibge,
  trim(uf) as uf,
  trim(regiao) as regiao,
  cast(valor_formalizado as varchar) as valor_formalizado,
  cast(valor_executado as varchar) as valor_executado,
  cast(valor_devolvido as varchar) as valor_devolvido
{% endmacro %}

{% macro conab_precos_semanal_uf_columns() %}
  trim(produto) as produto,
  trim(classificao_produto) as classificacao_produto,
  cast(id_produto as varchar) as id_produto,
  trim(uf) as uf,
  trim(regiao) as regiao,
  cast(ano as varchar) as ano,
  cast(mes as varchar) as mes,
  trim(data_inicial_final_semana) as data_inicial_final_semana,
  cast(semana as varchar) as semana,
  trim(dsc_nivel_comercializacao) as nivel_comercializacao,
  cast(valor_produto_kg as varchar) as valor_produto_kg
{% endmacro %}

{% macro conab_precos_semanal_municipio_columns() %}
  trim(produto) as produto,
  trim(classificao_produto) as classificacao_produto,
  cast(id_produto as varchar) as id_produto,
  trim(nom_municipio) as municipio,
  lpad(trim(cod_ibge), 7, '0') as cod_ibge,
  trim(uf) as uf,
  trim(regiao) as regiao,
  cast(ano as varchar) as ano,
  cast(mes as varchar) as mes,
  trim(data_inicial_final_semana) as data_inicial_final_semana,
  cast(semana as varchar) as semana,
  trim(dsc_nivel_comercializacao) as nivel_comercializacao,
  cast(valor_produto_kg as varchar) as valor_produto_kg
{% endmacro %}

{% macro conab_precos_mensal_uf_columns() %}
  trim(produto) as produto,
  trim(classificao_produto) as classificacao_produto,
  cast(id_produto as varchar) as id_produto,
  trim(uf) as uf,
  trim(regiao) as regiao,
  cast(ano as varchar) as ano,
  cast(mes as varchar) as mes,
  trim(dsc_nivel_comercializacao) as nivel_comercializacao,
  cast(valor_produto_kg as varchar) as valor_produto_kg
{% endmacro %}

{% macro conab_precos_mensal_municipio_columns() %}
  trim(produto) as produto,
  trim(classificao_produto) as classificacao_produto,
  cast(id_produto as varchar) as id_produto,
  trim(nom_municipio) as municipio,
  lpad(trim(cod_ibge), 7, '0') as cod_ibge,
  trim(uf) as uf,
  trim(regiao) as regiao,
  cast(ano as varchar) as ano,
  cast(mes as varchar) as mes,
  trim(dsc_nivel_comercializacao) as nivel_comercializacao,
  cast(valor_produto_kg as varchar) as valor_produto_kg
{% endmacro %}
