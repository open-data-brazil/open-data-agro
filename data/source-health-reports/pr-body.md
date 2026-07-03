## Source health daily probe

### Source health probe

- Executed at: 2026-07-03T06:37:15Z
- Run date: 2026-07-03
- Datasets probed: 131
- OK: 126 · Warning: 5 · Critical: 0
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
- **transportes.mtr-bit-malha-rodoviaria** (warning, day 1): Possible link deprecation — official source unreachable after retries. (probe failed after 3 attempts: Get "https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=jurisdicao-de-vias": dial tcp 189.9.19.9:443: i/o timeout)
  - https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=jurisdicao-de-vias
  - https://www.gov.br/transportes/pt-br/assuntos/dados-de-transportes/bit/bit-mapas
- **transportes.mtr-bit-malha-shapefile** (warning, day 1): Possible link deprecation — official source unreachable after retries. (probe failed after 3 attempts: Get "https://www.gov.br/transportes/pt-br/assuntos/dados-de-transportes/bit/bit-mapas/Base-GEO/BaseFerro.zip": dial tcp 161.148.164.31:443: i/o timeout)
  - https://www.gov.br/transportes/pt-br/assuntos/dados-de-transportes/bit/bit-mapas/Base-GEO/BaseFerro.zip
  - https://www.gov.br/transportes/pt-br/assuntos/dados-de-transportes/bit/bit-mapas


See `data/source-health-reports/latest.json` for full outcomes.
