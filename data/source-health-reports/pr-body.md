## Source health daily probe

### Source health probe

- Executed at: 2026-07-26T06:10:16Z
- Run date: 2026-07-26
- Datasets probed: 131
- OK: 125 · Warning: 2 · Critical: 4
- Updated samples: 18
- Deprecated (2+ days): 4

### Source health alerts

- **ana.hidrologia-series** (warning, day 1): Possible link deprecation — official source unreachable after retries. (probe failed after 3 attempts: Get "https://telemetriaws1.ana.gov.br/ServiceANA.asmx?WSDL": context deadline exceeded (Client.Timeout exceeded while awaiting headers))
  - https://telemetriaws1.ana.gov.br/ServiceANA.asmx?WSDL
  - https://www.gov.br/ana/pt-br/acesso-a-informacao/dados-abertos
- **ana.pluviometria-redes** (warning, day 1): Possible link deprecation — official source unreachable after retries. (probe failed after 3 attempts: Get "https://telemetriaws1.ana.gov.br/ServiceANA.asmx?WSDL": context deadline exceeded (Client.Timeout exceeded while awaiting headers))
  - https://telemetriaws1.ana.gov.br/ServiceANA.asmx?WSDL
  - https://www.gov.br/ana/pt-br/acesso-a-informacao/dados-abertos
- **dnit.condicoes-conservacao-rodovias** (critical, day 6): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: Get "https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=condicoes-do-pavimento": dial tcp 189.9.19.9:443: i/o timeout)
  - https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=condicoes-do-pavimento
  - https://servicos.dnit.gov.br/dadosabertos/dataset/condicoes-do-pavimento
- **dnit.snv-rodovias-federais** (critical, day 6): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: Get "https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=jurisdicao-de-vias": dial tcp 189.9.19.9:443: i/o timeout)
  - https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=jurisdicao-de-vias
  - https://servicos.dnit.gov.br/dadosabertos/dataset/jurisdicao-de-vias
- **transportes.mtr-bit-malha-rodoviaria** (critical, day 6): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: Get "https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=jurisdicao-de-vias": dial tcp 189.9.19.9:443: i/o timeout)
  - https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=jurisdicao-de-vias
  - https://www.gov.br/transportes/pt-br/assuntos/dados-de-transportes/bit/bit-mapas
- **transportes.mtr-bit-malha-shapefile** (critical, day 6): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: unexpected status 404 for https://www.gov.br/transportes/pt-br/assuntos/dados-de-transportes/bit/bit-mapas/Base-GEO/BaseFerro.zip)
  - https://www.gov.br/transportes/pt-br/assuntos/dados-de-transportes/bit/bit-mapas/Base-GEO/BaseFerro.zip
  - https://www.gov.br/transportes/pt-br/assuntos/dados-de-transportes/bit/bit-mapas


See `data/source-health-reports/latest.json` for full outcomes.
