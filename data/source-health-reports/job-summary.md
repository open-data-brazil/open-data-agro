### Source health probe

- Executed at: 2026-08-16T03:50:33Z
- Run date: 2026-08-16
- Datasets probed: 131
- OK: 122 · Warning: 4 · Critical: 5
- Updated samples: 14
- Deprecated (2+ days): 5

### Source health alerts

- **dnit.condicoes-conservacao-rodovias** (critical, day 27): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: Get "https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=condicoes-do-pavimento": dial tcp 189.9.19.9:443: i/o timeout)
  - https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=condicoes-do-pavimento
  - https://servicos.dnit.gov.br/dadosabertos/dataset/condicoes-do-pavimento
- **dnit.snv-rodovias-federais** (critical, day 27): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: Get "https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=jurisdicao-de-vias": dial tcp 189.9.19.9:443: i/o timeout)
  - https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=jurisdicao-de-vias
  - https://servicos.dnit.gov.br/dadosabertos/dataset/jurisdicao-de-vias
- **ibge.pevs-producao-vegetal** (warning, day 1): Possible link deprecation — official source unreachable after retries. (probe failed after 3 attempts: Get "https://apisidra.ibge.gov.br/values/t/289/n3/all/p/2023/v/144": dial tcp 170.84.40.190:443: i/o timeout)
  - https://apisidra.ibge.gov.br/values/t/289/n3/all/p/2023/v/144
  - https://sidra.ibge.gov.br/pesquisa/pevs
- **mdic.comex-exportacao-uf-ncm** (warning, day 1): Possible link deprecation — official source unreachable after retries. (probe failed after 3 attempts: unexpected status 429 for https://api-comexstat.mdic.gov.br/general)
  - https://api-comexstat.mdic.gov.br/general
  - https://comexstat.mdic.gov.br/
- **mdic.comex-importacao-diesel-ncm** (warning, day 1): Possible link deprecation — official source unreachable after retries. (probe failed after 3 attempts: unexpected status 429 for https://api-comexstat.mdic.gov.br/general)
  - https://api-comexstat.mdic.gov.br/general
  - https://comexstat.mdic.gov.br/
- **ons.carga-energetica** (critical, day 2): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: unexpected status 403 for https://dados.ons.org.br/api/3/action/package_show?id=carga-energia)
  - https://dados.ons.org.br/api/3/action/package_show?id=carga-energia
  - https://dados.ons.org.br/dataset/carga-energia
- **suframa.comercio-mercadorias-zfm** (warning, day 1): Possible link deprecation — official source unreachable after retries. (probe failed after 3 attempts: Get "https://www.gov.br/suframa/pt-br/acesso-a-informacao/dados-abertos/base-de-dados/sao/copy2_of_RelatriodeNotasFiscaisVistoriadasporregiodecontroledaSuframa2021.xlsx": dial tcp 161.148.164.31:443: i/o timeout)
  - https://www.gov.br/suframa/pt-br/acesso-a-informacao/dados-abertos/base-de-dados/sao/copy2_of_RelatriodeNotasFiscaisVistoriadasporregiodecontroledaSuframa2021.xlsx
  - https://www.gov.br/suframa/pt-br/acesso-a-informacao/dados-abertos/base-de-dados
- **transportes.mtr-bit-malha-rodoviaria** (critical, day 27): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: Get "https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=jurisdicao-de-vias": dial tcp 189.9.19.9:443: i/o timeout)
  - https://servicos.dnit.gov.br/dadosabertos/api/3/action/package_show?id=jurisdicao-de-vias
  - https://www.gov.br/transportes/pt-br/assuntos/dados-de-transportes/bit/bit-mapas
- **transportes.mtr-bit-malha-shapefile** (critical, day 27): Consultation link deprecated — official source unreachable for 2 or more consecutive days. (probe failed after 3 attempts: Get "https://www.gov.br/transportes/pt-br/assuntos/dados-de-transportes/bit/bit-mapas/Base-GEO/BaseFerro.zip": dial tcp 161.148.164.31:443: i/o timeout)
  - https://www.gov.br/transportes/pt-br/assuntos/dados-de-transportes/bit/bit-mapas/Base-GEO/BaseFerro.zip
  - https://www.gov.br/transportes/pt-br/assuntos/dados-de-transportes/bit/bit-mapas
