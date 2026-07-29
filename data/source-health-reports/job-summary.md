### Source health probe

- Executed at: 2026-07-29T05:56:32Z
- Run date: 2026-07-29
- Datasets probed: 131
- OK: 121 · Warning: 5 · Critical: 5
- Updated samples: 25
- Deprecated (2+ days): 5

### Source health alerts

- **anp.combustiveis-precos-medios-municipios** (warning, day 1): Possible link deprecation — official source unreachable after retries. (fetch LPC listing: Get "https://www.gov.br/anp/pt-br/assuntos/precos-e-defesa-da-concorrencia/precos/levantamento-de-precos-de-combustiveis-ultimas-semanas-pesquisadas": read tcp 10.1.0.126:49252->161.148.164.31:443: read: connection reset by peer)
  - https://www.gov.br/anp/pt-br/assuntos/precos-e-defesa-da-concorrencia/precos/levantamento-de-precos-de-combustiveis-ultimas-semanas-pesquisadas
- **anp.combustiveis-precos-postos** (warning, day 1): Possible link deprecation — official source unreachable after retries. (fetch LPC listing: Get "https://www.gov.br/anp/pt-br/assuntos/precos-e-defesa-da-concorrencia/precos/levantamento-de-precos-de-combustiveis-ultimas-semanas-pesquisadas": read tcp 10.1.0.126:49238->161.148.164.31:443: read: connection reset by peer)
  - https://www.gov.br/anp/pt-br/assuntos/precos-e-defesa-da-concorrencia/precos/levantamento-de-precos-de-combustiveis-ultimas-semanas-pesquisadas
- **anp.etanol-precos** (warning, day 1): Possible link deprecation — official source unreachable after retries. (fetch LPC listing: Get "https://www.gov.br/anp/pt-br/assuntos/precos-e-defesa-da-concorrencia/precos/levantamento-de-precos-de-combustiveis-ultimas-semanas-pesquisadas": read tcp 10.1.0.126:49262->161.148.164.31:443: read: connection reset by peer)
  - https://www.gov.br/anp/pt-br/assuntos/precos-e-defesa-da-concorrencia/precos/levantamento-de-precos-de-combustiveis-ultimas-semanas-pesquisadas
  - https://www.gov.br/anp/pt-br/assuntos/precos-e-indices/precos-de-combustiveis
- **dnit.condicoes-conservacao-rodovias** (critical, day 9): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: Get "https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=condicoes-do-pavimento": dial tcp 189.9.19.9:443: i/o timeout)
  - https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=condicoes-do-pavimento
  - https://servicos.dnit.gov.br/dadosabertos/dataset/condicoes-do-pavimento
- **dnit.snv-rodovias-federais** (critical, day 9): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: Get "https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=jurisdicao-de-vias": dial tcp 189.9.19.9:443: i/o timeout)
  - https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=jurisdicao-de-vias
  - https://servicos.dnit.gov.br/dadosabertos/dataset/jurisdicao-de-vias
- **ibge.lspa-rendimento-medio** (critical, day 2): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: Get "https://apisidra.ibge.gov.br/values/t/6588/n3/in%20n3%2011/p/202512/v/35/c48/39443": dial tcp 170.84.40.190:443: i/o timeout)
  - https://apisidra.ibge.gov.br/values/t/6588/n3/in%20n3%2011/p/202512/v/35/c48/39443
  - https://sidra.ibge.gov.br/pesquisa/lspa
- **ibge.pnad-rural-renda-ocupacao** (warning, day 1): Possible link deprecation — official source unreachable after retries. (probe failed after 3 attempts: Get "https://apisidra.ibge.gov.br/values/t/6385/n3/11,12,13,14,15,16,17/p/last%201/v/all": dial tcp 170.84.40.190:443: i/o timeout)
  - https://apisidra.ibge.gov.br/values/t/6385/n3/11,12,13,14,15,16,17/p/last%201/v/all
  - https://sidra.ibge.gov.br/pesquisa/pnad
- **mdic.comex-exportacao-uf-ncm** (warning, day 1): Possible link deprecation — official source unreachable after retries. (probe failed after 3 attempts: unexpected status 429 for https://api-comexstat.mdic.gov.br/general)
  - https://api-comexstat.mdic.gov.br/general
  - https://comexstat.mdic.gov.br/
- **transportes.mtr-bit-malha-rodoviaria** (critical, day 9): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: Get "https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=jurisdicao-de-vias": dial tcp 189.9.19.9:443: i/o timeout)
  - https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=jurisdicao-de-vias
  - https://www.gov.br/transportes/pt-br/assuntos/dados-de-transportes/bit/bit-mapas
- **transportes.mtr-bit-malha-shapefile** (critical, day 9): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: unexpected status 404 for https://www.gov.br/transportes/pt-br/assuntos/dados-de-transportes/bit/bit-mapas/Base-GEO/BaseFerro.zip)
  - https://www.gov.br/transportes/pt-br/assuntos/dados-de-transportes/bit/bit-mapas/Base-GEO/BaseFerro.zip
  - https://www.gov.br/transportes/pt-br/assuntos/dados-de-transportes/bit/bit-mapas
