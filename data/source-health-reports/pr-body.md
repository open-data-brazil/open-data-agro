## Source health daily probe

### Source health probe

- Executed at: 2026-07-09T06:53:55Z
- Run date: 2026-07-09
- Datasets probed: 131
- OK: 123 · Warning: 8 · Critical: 0
- Updated samples: 27
- Deprecated (2+ days): 0

### Source health alerts

- **anp.combustiveis-precos-medios-municipios** (warning, day 1): Possible link deprecation — official source unreachable after retries. (fetch LPC listing: Get "https://www.gov.br/anp/pt-br/assuntos/precos-e-defesa-da-concorrencia/precos/levantamento-de-precos-de-combustiveis-ultimas-semanas-pesquisadas": dial tcp 161.148.164.31:443: i/o timeout)
  - https://www.gov.br/anp/pt-br/assuntos/precos-e-defesa-da-concorrencia/precos/levantamento-de-precos-de-combustiveis-ultimas-semanas-pesquisadas
- **anp.etanol-precos** (warning, day 1): Possible link deprecation — official source unreachable after retries. (fetch LPC listing: Get "https://www.gov.br/anp/pt-br/assuntos/precos-e-defesa-da-concorrencia/precos/levantamento-de-precos-de-combustiveis-ultimas-semanas-pesquisadas": dial tcp 161.148.164.31:443: i/o timeout)
  - https://www.gov.br/anp/pt-br/assuntos/precos-e-defesa-da-concorrencia/precos/levantamento-de-precos-de-combustiveis-ultimas-semanas-pesquisadas
  - https://www.gov.br/anp/pt-br/assuntos/precos-e-indices/precos-de-combustiveis
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
- **transportes.mtr-bit-malha-rodoviaria** (warning, day 1): Possible link deprecation — official source unreachable after retries. (probe failed after 3 attempts: Get "https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=jurisdicao-de-vias": dial tcp 189.9.19.9:443: i/o timeout)
  - https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=jurisdicao-de-vias
  - https://www.gov.br/transportes/pt-br/assuntos/dados-de-transportes/bit/bit-mapas
- **transportes.mtr-bit-malha-shapefile** (warning, day 1): Possible link deprecation — official source unreachable after retries. (probe failed after 3 attempts: Get "https://www.gov.br/transportes/pt-br/assuntos/dados-de-transportes/bit/bit-mapas/Base-GEO/BaseFerro.zip": dial tcp 161.148.164.31:443: i/o timeout)
  - https://www.gov.br/transportes/pt-br/assuntos/dados-de-transportes/bit/bit-mapas/Base-GEO/BaseFerro.zip
  - https://www.gov.br/transportes/pt-br/assuntos/dados-de-transportes/bit/bit-mapas


See `data/source-health-reports/latest.json` for full outcomes.
