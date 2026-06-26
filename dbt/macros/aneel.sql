{% macro aneel_fonte_oficial() %}
  https://dadosabertos.aneel.gov.br/dataset/bandeiras-tarifarias
{% endmacro %}

{% macro aneel_tarifas_energia_columns() %}
    DatGeracaoConjuntoDados,
    DatCompetencia,
    NomBandeiraAcionada,
    VlrAdicionalBandeira
{% endmacro %}
