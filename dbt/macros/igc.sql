{% macro igc_fonte_oficial() %}
  https://igc.int/en/public-site/markets/marketinfo-goi.aspx
{% endmacro %}

{% macro igc_goi_columns() %}
    refdate,
    index_slug,
    index_name,
    value,
    base_period,
    frequency
{% endmacro %}
