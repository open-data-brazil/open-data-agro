{% macro argentina_granos_fonte_oficial() %}
  https://datos.magyp.gob.ar/
{% endmacro %}

{% macro argentina_granos_columns() %}
    series_id,
    commodity_slug,
    refyear,
    value,
    unit,
    source
{% endmacro %}
