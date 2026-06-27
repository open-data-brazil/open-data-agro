{% macro nasa_power_agroclimatology_fonte_oficial() %}
  https://power.larc.nasa.gov/
{% endmacro %}

{% macro nasa_power_agroclimatology_columns() %}
    latitude, longitude, refdate, parameter_slug, value
{% endmacro %}

{% macro nasa_power_agroclimatology_dedupe_keys() %}
    latitude, longitude, refdate, parameter_slug
{% endmacro %}
