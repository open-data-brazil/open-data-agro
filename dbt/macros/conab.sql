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

{% macro conab_estimativa_cana_columns() %}
  trim(produto) as produto,
  trim(uf) as uf,
  trim(ano_agricola) as safra,
  trim(dsc_safra_previsao) as safra_previsao,
  cast(id_produto as varchar) as id_produto,
  trim(dsc_levantamento) as levantamento,
  cast(id_levantamento as varchar) as id_levantamento,
  cast(area_plantada_mil_ha as varchar) as area_plantada_mil_ha,
  cast(producao_mil_t as varchar) as producao_mil_t,
  cast(producao_acucar_mil_t as varchar) as producao_acucar_mil_t,
  cast(producao_etanol_anidro_mil_l as varchar) as producao_etanol_anidro_mil_l,
  cast(producao_etanol_hidratado_mil_l as varchar) as producao_etanol_hidratado_mil_l,
  cast(producao_etanol_total_mil_l as varchar) as producao_etanol_total_mil_l,
  cast(produtcao_atr_kg_t as varchar) as producao_atr_kg_t
{% endmacro %}

{% macro conab_serie_historica_cana_columns() %}
  trim(produto) as produto,
  trim(uf) as uf,
  trim(ano_agricola) as ano,
  trim(dsc_safra_previsao) as safra_previsao,
  cast(id_produto as varchar) as id_produto,
  cast(area_plantada_mil_ha as varchar) as area_plantada_mil_ha,
  cast(producao_mil_t as varchar) as producao_mil_t,
  trim(dsc_situacao_levantamento) as situacao_levantamento,
  cast(producao_acucar_mil_t as varchar) as producao_acucar_mil_t,
  cast(producao_etanol_anidro_mil_l as varchar) as producao_etanol_anidro_mil_l,
  cast(producao_etanol_hidratado_mil_l as varchar) as producao_etanol_hidratado_mil_l,
  cast(producao_etanol_total_mil_l as varchar) as producao_etanol_total_mil_l,
  cast(produtcao_atr_kg_t as varchar) as producao_atr_kg_t
{% endmacro %}

{% macro conab_estimativa_cafe_columns() %}
  trim(produto) as produto,
  trim(uf) as uf,
  trim(ano_agricola) as safra,
  trim(safra) as tipo_safra,
  cast(id_produto as varchar) as id_produto,
  cast(id_levantamento as varchar) as id_levantamento,
  trim(dsc_levantamento) as levantamento,
  cast(area_plantada_mil_ha as varchar) as area_plantada_mil_ha,
  cast(producao_mil_t as varchar) as producao_mil_t,
  cast(produtividade_mil_ha_mil_t as varchar) as produtividade_mil_ha_mil_t
{% endmacro %}

{% macro conab_serie_historica_cafe_columns() %}
  trim(produto) as produto,
  trim(uf) as uf,
  trim(ano_agricola) as ano,
  trim(dsc_safra_previsao) as safra_previsao,
  cast(id_produto as varchar) as id_produto,
  cast(area_plantada_mil_ha as varchar) as area_plantada_mil_ha,
  cast(producao_mil_t as varchar) as producao_mil_t,
  cast(produtividade_mil_ha_mil_t as varchar) as produtividade_mil_ha_mil_t
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
  lpad(trim(cod_ibge), 7, '0') as cod_ibge,
  trim(uf) as uf,
  cast(num_ano as varchar) as ano,
  cast(num_mes as varchar) as mes,
  trim(conta_operacional) as conta_operacional,
  cast(qtd_estoque_kg as varchar) as qtd_estoque_kg
{% endmacro %}

{% macro conab_operacoes_comercializacao_columns() %}
  trim(id_edital) as id_edital,
  cast(num_lote as varchar) as num_lote,
  trim(num_dco) as num_dco,
  trim(dsc_dco) as dsc_dco,
  trim(dsc_situacao_dco) as situacao_dco,
  trim(produto) as produto,
  cast(id_produto as varchar) as id_produto,
  trim(dsc_tipo_operacao) as tipo_operacao,
  trim(dsc_operacao) as operacao,
  cast(ano_edital as varchar) as ano_edital,
  cast(mes_edital as varchar) as mes_edital,
  trim(uf_armazem_origem) as uf_armazem_origem,
  trim(dsc_unidade_comercializacao) as unidade_comercializacao,
  cast(qtd_ofertada as varchar) as qtd_ofertada,
  cast(qtd_negociada as varchar) as qtd_negociada,
  cast(vlr_operacao_s_icms as varchar) as vlr_operacao_s_icms,
  trim(dsc_unidade_medida_ofertada_negociada) as unidade_medida_ofertada_negociada
{% endmacro %}

{% macro conab_vendas_balcao_columns() %}
  cast(num_ano_gravacao as varchar) as ano,
  cast(num_mes_gravacao as varchar) as mes,
  trim(munipio_armazem_venda) as municipio_armazem_venda,
  trim(uf) as uf,
  cast(qtd_produto_kg as varchar) as qtd_produto_kg,
  cast(valor_comercializado as varchar) as valor_comercializado,
  cast(numero_atendimentos as varchar) as numero_atendimentos,
  cast(clientes_atendidos as varchar) as clientes_atendidos
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

{% macro conab_precos_minimos_columns() %}
  trim(descricao_produto_preco_minimo) as produto,
  cast(id_produto as varchar) as id_produto,
  trim(uf) as uf,
  trim(regionalizacao) as regionalizacao,
  cast(ano_inicio_vigencia as varchar) as ano_inicio_vigencia,
  cast(mes_incio_vigencia as varchar) as mes_inicio_vigencia,
  cast(ano_termino_vigencia as varchar) as ano_termino_vigencia,
  cast(mes_termino_vigencia as varchar) as mes_termino_vigencia,
  cast(preco as varchar) as preco,
  trim(dsc_unidade_comercializacao) as unidade_comercializacao,
  trim(nome_normativo) as nome_normativo,
  trim(url) as url_normativo
{% endmacro %}

{% macro conab_frete_columns() %}
  trim(dsc_fonte) as fonte,
  trim(municipio_origem) as municipio_origem,
  lpad(trim(cod_ibge_origem), 7, '0') as cod_ibge_origem,
  trim(uf_origem) as uf_origem,
  trim(municipio_destino) as municipio_destino,
  lpad(trim(cod_ibge_destino), 7, '0') as cod_ibge_destino,
  trim(uf_destino) as uf_destino,
  cast(ano as varchar) as ano,
  cast(mes as varchar) as mes,
  cast(distancia_km as varchar) as distancia_km,
  cast(valor_frete_tonelada as varchar) as valor_frete_tonelada,
  cast(valor_tonelada_km as varchar) as valor_tonelada_km
{% endmacro %}

{% macro conab_capacidade_estatica_columns() %}
  cast("Ano" as varchar) as ano,
  trim("UF") as uf,
  cast("Quantidade (mil t)" as varchar) as quantidade_mil_t
{% endmacro %}

{% macro conab_prohort_diario_columns() %}
  trim(municipio_ceasa) as municipio_ceasa,
  lpad(trim(cod_ibge_municipio), 7, '0') as cod_ibge_municipio,
  trim(uf_ceasa) as uf_ceasa,
  trim(dsc_ceasa) as ceasa,
  trim(dsc_produto) as produto,
  trim(sig_unidade_medida) as unidade_medida,
  trim(data_preco) as data_preco,
  cast(preco_diario as varchar) as preco_diario
{% endmacro %}

{% macro conab_prohort_mensal_columns() %}
  cast(id_ano_comercializacao as varchar) as ano,
  cast(id_mes_comercializacao as varchar) as mes,
  trim(municipio_origem_produto) as municipio_origem,
  lpad(trim(cod_ibge_municipio_origem_produto), 7, '0') as cod_ibge_municipio_origem,
  trim(uf_origem_produto) as uf_origem,
  trim(dsc_ceasa) as ceasa,
  trim(uf_ceasa) as uf_ceasa,
  trim(municipio_ceasa) as municipio_ceasa,
  lpad(trim(cod_ibge_municipio_ceasa), 7, '0') as cod_ibge_municipio_ceasa,
  trim(dsc_produto) as produto,
  cast(qtd_comercializada_kg as varchar) as qtd_comercializada_kg,
  cast(valor_comercializado as varchar) as valor_comercializado,
  trim(pais_origem) as pais_origem
{% endmacro %}
