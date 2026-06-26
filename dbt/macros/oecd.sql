{% macro oecd_fonte_oficial() %}
  https://www.oecd.org/en/data/datasets/oecd-fao-agricultural-outlook.html
{% endmacro %}

{% macro oecd_ag_outlook_columns() %}
    ref_area,
    ref_area_name,
    commodity_code,
    commodity_name,
    measure_code,
    measure_name,
    unit,
    unit_mult,
    year,
    value,
    obs_status
{% endmacro %}
