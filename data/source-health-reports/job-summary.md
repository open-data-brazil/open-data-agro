### Source health probe

- Executed at: 2026-08-25T03:52:52Z
- Run date: 2026-08-25
- Datasets probed: 131
- OK: 119 · Warning: 8 · Critical: 4
- Updated samples: 18
- Deprecated (2+ days): 4

### Source health alerts

- **dnit.condicoes-conservacao-rodovias** (critical, day 36): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: Get "https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=condicoes-do-pavimento": dial tcp 189.9.19.9:443: i/o timeout)
  - https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=condicoes-do-pavimento
  - https://servicos.dnit.gov.br/dadosabertos/dataset/condicoes-do-pavimento
- **dnit.snv-rodovias-federais** (critical, day 36): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: Get "https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=jurisdicao-de-vias": dial tcp 189.9.19.9:443: i/o timeout)
  - https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=jurisdicao-de-vias
  - https://servicos.dnit.gov.br/dadosabertos/dataset/jurisdicao-de-vias
- **ibge.censo-agro-estabelecimentos** (warning, day 1): Possible link deprecation — official source unreachable after retries. (probe failed after 3 attempts: Get "https://apisidra.ibge.gov.br/values/t/6878/n3/all/p/2017/v/all": dial tcp 170.84.40.190:443: i/o timeout)
  - https://apisidra.ibge.gov.br/values/t/6878/n3/all/p/2017/v/all
  - https://censoagro2017.ibge.gov.br/
- **ibge.lspa-area-producao** (warning, day 1): Possible link deprecation — official source unreachable after retries. (probe failed after 3 attempts: Get "https://apisidra.ibge.gov.br/values/t/6588/n3/in%20n3%2011/p/202512/v/109/c48/39443": dial tcp 170.84.40.190:443: i/o timeout)
  - https://apisidra.ibge.gov.br/values/t/6588/n3/in%20n3%2011/p/202512/v/109/c48/39443
  - https://sidra.ibge.gov.br/pesquisa/lspa
- **ibge.ppm-producao-municipal** (warning, day 1): Possible link deprecation — official source unreachable after retries. (probe failed after 3 attempts: Get "https://apisidra.ibge.gov.br/values/t/74/n6/in%20n3%2011/p/2023/v/106": dial tcp 170.84.40.190:443: i/o timeout)
  - https://apisidra.ibge.gov.br/values/t/74/n6/in%20n3%2011/p/2023/v/106
  - https://sidra.ibge.gov.br/pesquisa/pam
- **ibge.ppm-efetivo-rebanhos** (warning, day 1): Possible link deprecation — official source unreachable after retries. (probe failed after 3 attempts: Get "https://apisidra.ibge.gov.br/values/t/3939/n3/11/p/2023/v/all": dial tcp 170.84.40.190:443: i/o timeout)
  - https://apisidra.ibge.gov.br/values/t/3939/n3/11/p/2023/v/all
  - https://sidra.ibge.gov.br/pesquisa/ppm
- **ibge.ppm-vacas-ordenhadas** (warning, day 1): Possible link deprecation — official source unreachable after retries. (probe failed after 3 attempts: Get "https://apisidra.ibge.gov.br/values/t/94/n3/11/p/2023/v/all": dial tcp 170.84.40.190:443: i/o timeout)
  - https://apisidra.ibge.gov.br/values/t/94/n3/11/p/2023/v/all
  - https://sidra.ibge.gov.br/pesquisa/ppm
- **ibge.ppm-ovinos-tosquiados** (warning, day 1): Possible link deprecation — official source unreachable after retries. (probe failed after 3 attempts: Get "https://apisidra.ibge.gov.br/values/t/95/n3/11/p/2023/v/all": dial tcp 170.84.40.190:443: i/o timeout)
  - https://apisidra.ibge.gov.br/values/t/95/n3/11/p/2023/v/all
  - https://sidra.ibge.gov.br/pesquisa/ppm
- **ibge.ppm-aquicultura** (warning, day 1): Possible link deprecation — official source unreachable after retries. (probe failed after 3 attempts: Get "https://apisidra.ibge.gov.br/values/t/3940/n3/11/p/2023/v/all": dial tcp 170.84.40.190:443: i/o timeout)
  - https://apisidra.ibge.gov.br/values/t/3940/n3/11/p/2023/v/all
  - https://sidra.ibge.gov.br/pesquisa/ppm
- **ibge.pnad-rural-renda-ocupacao** (warning, day 1): Possible link deprecation — official source unreachable after retries. (probe failed after 3 attempts: Get "https://apisidra.ibge.gov.br/values/t/6385/n3/11,12,13,14,15,16,17/p/last%201/v/all": dial tcp 170.84.40.190:443: i/o timeout)
  - https://apisidra.ibge.gov.br/values/t/6385/n3/11,12,13,14,15,16,17/p/last%201/v/all
  - https://sidra.ibge.gov.br/pesquisa/pnad
- **transportes.mtr-bit-malha-rodoviaria** (critical, day 36): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: Get "https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=jurisdicao-de-vias": dial tcp 189.9.19.9:443: i/o timeout)
  - https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=jurisdicao-de-vias
  - https://www.gov.br/transportes/pt-br/assuntos/dados-de-transportes/bit/bit-mapas
- **transportes.mtr-bit-malha-shapefile** (critical, day 36): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: unexpected status 404 for https://www.gov.br/transportes/pt-br/assuntos/dados-de-transportes/bit/bit-mapas/Base-GEO/BaseFerro.zip)
  - https://www.gov.br/transportes/pt-br/assuntos/dados-de-transportes/bit/bit-mapas/Base-GEO/BaseFerro.zip
  - https://www.gov.br/transportes/pt-br/assuntos/dados-de-transportes/bit/bit-mapas
