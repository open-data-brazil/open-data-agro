{% macro mexico_siap_produccion_agricola_fonte_oficial() %}
  https://www.gob.mx/siap
{% endmacro %}

{% macro mexico_siap_produccion_agricola_columns() %}
    state_code, crop_slug, refyear, area_ha, production_t
{% endmacro %}

{% macro mexico_siap_produccion_agricola_dedupe_keys() %}
    state_code, crop_slug, refyear
{% endmacro %}
