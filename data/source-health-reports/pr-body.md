## Source health daily probe

### Source health probe

- Executed at: 2026-07-19T05:57:07Z
- Run date: 2026-07-19
- Datasets probed: 131
- OK: 120 · Warning: 11 · Critical: 0
- Updated samples: 31
- Deprecated (2+ days): 0

### Source health alerts

- **dnit.condicoes-conservacao-rodovias** (warning, day 1): Possible link deprecation — official source unreachable after retries. (probe failed after 3 attempts: Get "https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=condicoes-do-pavimento": dial tcp 189.9.19.9:443: i/o timeout)
  - https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=condicoes-do-pavimento
  - https://servicos.dnit.gov.br/dadosabertos/dataset/condicoes-do-pavimento
- **dnit.snv-rodovias-federais** (warning, day 1): Possible link deprecation — official source unreachable after retries. (probe failed after 3 attempts: Get "https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=jurisdicao-de-vias": dial tcp 189.9.19.9:443: i/o timeout)
  - https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=jurisdicao-de-vias
  - https://servicos.dnit.gov.br/dadosabertos/dataset/jurisdicao-de-vias
- **eia.petroleum-prices** (warning, day 1): Possible link deprecation — official source unreachable after retries. (probe failed after 3 attempts: unexpected status 403 for https://api.eia.gov/v2/petroleum/pri/spt/data)
  - https://api.eia.gov/v2/petroleum/pri/spt/data
  - https://www.eia.gov/opendata/
- **eurostat.ag-prices** (warning, day 1): Possible link deprecation — official source unreachable after retries. (probe failed after 3 attempts: unexpected status 404 for https://ec.europa.eu/eurostat/api/dissemination/statistics/1.0/data/apri_pi15_outa?format=JSON&geo=EU27_2020&lang=en&product=010000&product=011000&product=015000&sinceTimePeriod=2020)
  - https://ec.europa.eu/eurostat/api/dissemination/statistics/1.0/data/apri_pi15_outa?format=JSON&geo=EU27_2020&lang=en&product=010000&product=011000&product=015000&sinceTimePeriod=2020
  - https://ec.europa.eu/eurostat/web/agriculture/database
- **ibge.localidades-regioes** (warning, day 1): Possible link deprecation — official source unreachable after retries. (probe failed after 3 attempts: Get "https://servicodados.ibge.gov.br/api/v1/localidades/regioes?orderBy=nome": dial tcp 170.84.40.205:443: i/o timeout)
  - https://servicodados.ibge.gov.br/api/v1/localidades/regioes?orderBy=nome
  - https://servicodados.ibge.gov.br/api/docs/localidades
- **ibge.localidades-mesorregioes** (warning, day 1): Possible link deprecation — official source unreachable after retries. (probe failed after 3 attempts: Get "https://servicodados.ibge.gov.br/api/v1/localidades/mesorregioes?orderBy=nome": dial tcp 170.84.40.205:443: i/o timeout)
  - https://servicodados.ibge.gov.br/api/v1/localidades/mesorregioes?orderBy=nome
  - https://servicodados.ibge.gov.br/api/docs/localidades
- **ibge.localidades-microrregioes** (warning, day 1): Possible link deprecation — official source unreachable after retries. (probe failed after 3 attempts: Get "https://servicodados.ibge.gov.br/api/v1/localidades/microrregioes?orderBy=nome": dial tcp 170.84.40.205:443: i/o timeout)
  - https://servicodados.ibge.gov.br/api/v1/localidades/microrregioes?orderBy=nome
  - https://servicodados.ibge.gov.br/api/docs/localidades
- **ibge.lspa-area-producao** (warning, day 1): Possible link deprecation — official source unreachable after retries. (probe failed after 3 attempts: Get "https://apisidra.ibge.gov.br/values/t/6588/n3/in%20n3%2011/p/202512/v/109/c48/39443": dial tcp 170.84.40.190:443: i/o timeout)
  - https://apisidra.ibge.gov.br/values/t/6588/n3/in%20n3%2011/p/202512/v/109/c48/39443
  - https://sidra.ibge.gov.br/pesquisa/lspa
- **ibge.censo-agro-maquinario** (warning, day 1): Possible link deprecation — official source unreachable after retries. (probe failed after 3 attempts: Get "https://apisidra.ibge.gov.br/values/t/6880/n3/all/p/2017/v/all": dial tcp 170.84.40.190:443: i/o timeout)
  - https://apisidra.ibge.gov.br/values/t/6880/n3/all/p/2017/v/all
  - https://censoagro2017.ibge.gov.br/
- **transportes.mtr-bit-malha-rodoviaria** (warning, day 1): Possible link deprecation — official source unreachable after retries. (probe failed after 3 attempts: Get "https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=jurisdicao-de-vias": dial tcp 189.9.19.9:443: i/o timeout)
  - https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=jurisdicao-de-vias
  - https://www.gov.br/transportes/pt-br/assuntos/dados-de-transportes/bit/bit-mapas
- **transportes.mtr-bit-malha-shapefile** (warning, day 1): Possible link deprecation — official source unreachable after retries. (probe failed after 3 attempts: unexpected status 404 for https://www.gov.br/transportes/pt-br/assuntos/dados-de-transportes/bit/bit-mapas/Base-GEO/BaseFerro.zip)
  - https://www.gov.br/transportes/pt-br/assuntos/dados-de-transportes/bit/bit-mapas/Base-GEO/BaseFerro.zip
  - https://www.gov.br/transportes/pt-br/assuntos/dados-de-transportes/bit/bit-mapas


See `data/source-health-reports/latest.json` for full outcomes.
