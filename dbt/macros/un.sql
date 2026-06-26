{% macro un_fonte_oficial() %}
  https://comtradeplus.un.org/
{% endmacro %}

{% macro un_comtrade_columns() %}
    reporter_code,
    reporter_desc,
    partner_code,
    partner_desc,
    flow_code,
    flow_desc,
    period,
    hs_code,
    commodity_slug,
    trade_value_usd,
    netweight_kg,
    qty,
    qty_unit_abbr
{% endmacro %}
