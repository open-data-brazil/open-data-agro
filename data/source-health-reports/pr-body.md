## Source health daily probe

### Source health probe

- Executed at: 2026-07-15T05:35:20Z
- Run date: 2026-07-15
- Datasets probed: 131
- OK: 119 · Warning: 12 · Critical: 0
- Updated samples: 28
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
- **ibama.sisfogo-incendios** (warning, day 1): Possible link deprecation — official source unreachable after retries. (probe failed after 3 attempts: unexpected status 500 for https://dadosabertos.ibama.gov.br/dados/SISFOGO/ROI.csv)
  - https://dadosabertos.ibama.gov.br/dados/SISFOGO/ROI.csv
  - https://dadosabertos.ibama.gov.br/dataset/sisfogo-incendios-florestais
- **ibge.pam-culturas-estendidas** (warning, day 1): Possible link deprecation — official source unreachable after retries. (probe failed after 3 attempts: Get "https://apisidra.ibge.gov.br/values/t/1612/n6/in%20n3%2011/p/2024/v/109/c81/2720": dial tcp 170.84.40.190:443: i/o timeout)
  - https://apisidra.ibge.gov.br/values/t/1612/n6/in%20n3%2011/p/2024/v/109/c81/2720
  - https://sidra.ibge.gov.br/pesquisa/pam
- **ibge.lspa-rendimento-medio** (warning, day 1): Possible link deprecation — official source unreachable after retries. (probe failed after 3 attempts: Get "https://apisidra.ibge.gov.br/values/t/6588/n3/in%20n3%2011/p/202512/v/35/c48/39443": dial tcp 170.84.40.190:443: i/o timeout)
  - https://apisidra.ibge.gov.br/values/t/6588/n3/in%20n3%2011/p/202512/v/35/c48/39443
  - https://sidra.ibge.gov.br/pesquisa/lspa
- **ibge.censo-agro-area-uso-solo** (warning, day 1): Possible link deprecation — official source unreachable after retries. (probe failed after 3 attempts: Get "https://apisidra.ibge.gov.br/values/t/6879/n3/all/p/2017/v/all": dial tcp 170.84.40.190:443: i/o timeout)
  - https://apisidra.ibge.gov.br/values/t/6879/n3/all/p/2017/v/all
  - https://censoagro2017.ibge.gov.br/
- **ibge.censo-agro-maquinario** (warning, day 1): Possible link deprecation — official source unreachable after retries. (probe failed after 3 attempts: Get "https://apisidra.ibge.gov.br/values/t/6880/n3/all/p/2017/v/all": dial tcp 170.84.40.190:443: i/o timeout)
  - https://apisidra.ibge.gov.br/values/t/6880/n3/all/p/2017/v/all
  - https://censoagro2017.ibge.gov.br/
- **ibge.pnad-rural-renda-ocupacao** (warning, day 1): Possible link deprecation — official source unreachable after retries. (probe failed after 3 attempts: Get "https://apisidra.ibge.gov.br/values/t/6385/n3/11,12,13,14,15,16,17/p/last%201/v/all": dial tcp 170.84.40.190:443: i/o timeout)
  - https://apisidra.ibge.gov.br/values/t/6385/n3/11,12,13,14,15,16,17/p/last%201/v/all
  - https://sidra.ibge.gov.br/pesquisa/pnad
- **mdic.comex-importacao-diesel-ncm** (warning, day 1): Possible link deprecation — official source unreachable after retries. (probe failed after 3 attempts: unexpected status 429 for https://api-comexstat.mdic.gov.br/general)
  - https://api-comexstat.mdic.gov.br/general
  - https://comexstat.mdic.gov.br/
- **transportes.mtr-bit-malha-rodoviaria** (warning, day 1): Possible link deprecation — official source unreachable after retries. (probe failed after 3 attempts: Get "https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=jurisdicao-de-vias": dial tcp 189.9.19.9:443: i/o timeout)
  - https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=jurisdicao-de-vias
  - https://www.gov.br/transportes/pt-br/assuntos/dados-de-transportes/bit/bit-mapas
- **transportes.mtr-bit-malha-shapefile** (warning, day 1): Possible link deprecation — official source unreachable after retries. (probe failed after 3 attempts: unexpected status 404 for https://www.gov.br/transportes/pt-br/assuntos/dados-de-transportes/bit/bit-mapas/Base-GEO/BaseFerro.zip)
  - https://www.gov.br/transportes/pt-br/assuntos/dados-de-transportes/bit/bit-mapas/Base-GEO/BaseFerro.zip
  - https://www.gov.br/transportes/pt-br/assuntos/dados-de-transportes/bit/bit-mapas


See `data/source-health-reports/latest.json` for full outcomes.
