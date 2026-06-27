{% macro japan_maff_ag_trade_fonte_oficial() %}
  https://www.maff.go.jp/e/
{% endmacro %}

{% macro japan_maff_ag_trade_columns() %}
    commodity_slug, refyear, flow_code, value_jpy, quantity_t
{% endmacro %}

{% macro japan_maff_ag_trade_dedupe_keys() %}
    commodity_slug, refyear, flow_code
{% endmacro %}
