### Source health probe

- Executed at: 2026-08-02T06:00:32Z
- Run date: 2026-08-02
- Datasets probed: 131
- OK: 124 · Warning: 3 · Critical: 4
- Updated samples: 15
- Deprecated (2+ days): 4

### Source health alerts

- **antt.volume-trafego-pedagio** (warning, day 1): Possible link deprecation — official source unreachable after retries. (parse ckan response: invalid character '<' looking for beginning of value)
  - https://dados.antt.gov.br/dataset/volume-trafego-praca-pedagio
- **antt.receita-por-praca** (warning, day 1): Possible link deprecation — official source unreachable after retries. (parse ckan response: invalid character '<' looking for beginning of value)
  - https://dados.antt.gov.br/dataset/receita-por-praca
- **antt.pracas-pedagio** (warning, day 1): Possible link deprecation — official source unreachable after retries. (parse ckan response: invalid character '<' looking for beginning of value)
  - https://dados.antt.gov.br/dataset/praca-de-pedagio
- **dnit.condicoes-conservacao-rodovias** (critical, day 13): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: Get "https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=condicoes-do-pavimento": dial tcp 189.9.19.9:443: i/o timeout)
  - https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=condicoes-do-pavimento
  - https://servicos.dnit.gov.br/dadosabertos/dataset/condicoes-do-pavimento
- **dnit.snv-rodovias-federais** (critical, day 13): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: Get "https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=jurisdicao-de-vias": dial tcp 189.9.19.9:443: i/o timeout)
  - https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=jurisdicao-de-vias
  - https://servicos.dnit.gov.br/dadosabertos/dataset/jurisdicao-de-vias
- **transportes.mtr-bit-malha-rodoviaria** (critical, day 13): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: Get "https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=jurisdicao-de-vias": dial tcp 189.9.19.9:443: i/o timeout)
  - https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=jurisdicao-de-vias
  - https://www.gov.br/transportes/pt-br/assuntos/dados-de-transportes/bit/bit-mapas
- **transportes.mtr-bit-malha-shapefile** (critical, day 13): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: unexpected status 404 for https://www.gov.br/transportes/pt-br/assuntos/dados-de-transportes/bit/bit-mapas/Base-GEO/BaseFerro.zip)
  - https://www.gov.br/transportes/pt-br/assuntos/dados-de-transportes/bit/bit-mapas/Base-GEO/BaseFerro.zip
  - https://www.gov.br/transportes/pt-br/assuntos/dados-de-transportes/bit/bit-mapas
