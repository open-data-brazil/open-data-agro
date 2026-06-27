{% macro fao_giews_crop_prospects_fonte_oficial() %}
  https://www.fao.org/giews/
{% endmacro %}

{% macro fao_giews_crop_prospects_columns() %}
    country_code, country_name, crop_slug, marketing_year, production_trend, outlook_note
{% endmacro %}

{% macro fao_giews_crop_prospects_dedupe_keys() %}
    country_code, crop_slug, marketing_year
{% endmacro %}
