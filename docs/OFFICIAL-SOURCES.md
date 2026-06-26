# Official sources catalog

> Primary references per dataset. **No dataset without a cited official source.**

**Initial portal:** [CONAB — Downloads de Arquivos](https://portaldeinformacoes.conab.gov.br/download-arquivos.html)

---

## CONAB — Portal de Informações Agropecuárias

| Item | Value |
|------|-------|
| **Portal** | https://portaldeinformacoes.conab.gov.br/ |
| **Downloads** | https://portaldeinformacoes.conab.gov.br/download-arquivos.html |
| **Usage** | Reproduction allowed non-commercial with source citation; preserve data integrity |
| **Contact** | sutin@conab.gov.br |

Detailed per-dataset mapping: `.local/phases/10-conab-producao-agricola/OFFICIAL-REFERENCE.md` (local).

---

## Index

### Produção Agrícola (Phase 10 — MVP)

| Dataset ID | Portal label | Status |
|------------|--------------|--------|
| `conab.estimativa-graos` | Estimativa Grãos | **P0 MVP — implemented** |
| `conab.serie-historica-graos` | Série Histórica Grãos | **P0 MVP — implemented** |
| `conab.estimativa-cana` | Estimativa Cana-de-Açúcar | Planned |
| `conab.serie-historica-cana` | Série Histórica Cana-de-Açúcar | Planned |
| `conab.estimativa-cafe` | Estimativa Café | Planned |
| `conab.serie-historica-cafe` | Série Histórica Café | Planned |
| `conab.custo-producao` | Custo de Produção | Planned |

### Mercado (Phase 11)

| Dataset ID | Portal label | Status |
|------------|--------------|--------|
| `conab.oferta-demanda` | Oferta e Demanda | **P1 — implemented** |
| `conab.precos-minimos` | Preços Mínimos | Catalog registered |
| `conab.precos-agropecuarios-mensal-uf` | Preços agropecuários Mensal UF | **P1 — full pipeline** |
| `conab.precos-agropecuarios-mensal-municipio` | Preços agropecuários Mensal Município | **P1 — full pipeline** |
| `conab.precos-agropecuarios-semanal-uf` | Preços agropecuários Semanal UF | **P1 — full pipeline** |
| `conab.precos-agropecuarios-semanal-municipio` | Preços agropecuários Semanal Municipio | **P1 — full pipeline** |
| `conab.prohort-diario` | Prohort Diário | Catalog registered |
| `conab.prohort-mensal` | Prohort Mensal | Catalog registered |

### Abastecimento (Phase 12)

| Dataset ID | Portal label | Status |
|------------|--------------|--------|
| `conab.estoques-publicos` | Estoques Públicos | **P1 — implemented** |
| `conab.operacoes-comercializacao` | Operações de Comercialização | Catalog + GE + ingest |
| `conab.vendas-balcao` | Vendas em Balcão | Catalog + GE + ingest |

### ANP — Combustíveis (Phase 12 extension)

| Dataset ID | Portal label | Status |
|------------|--------------|--------|
| `anp.combustiveis-precos-medios-municipios` | LPC — preços médios por município | **P1 — implemented** |
| `anp.combustiveis-precos-postos` | LPC — preços por posto revendedor | **P1 — implemented** |

### Armazenamento e Logística (Phase 13)

| Dataset ID | Portal label | Status |
|------------|--------------|--------|
| `conab.armazenagem` | Armazenagem | **P1 — implemented** |
| `conab.frete` | Frete | Catalog + GE + ingest |
| `conab.serie-historica-capacidade-estatica` | Série Histórica da Capacidade Estática | Catalog + GE + `.xls` ingest |

### Agricultura Familiar (Phase 14)

| Dataset ID | Portal label | Status |
|------------|--------------|--------|
| `conab.alimenta-brasil-entregas` | Programa Alimenta Brasil - Entregas | **P1 — implemented** |
| `conab.alimenta-brasil-propostas` | Programa Alimenta Brasil - Propostas | **P1 — implemented** |

### IBGE — Localidades (Phase 15)

| Dataset ID | API resource | Status |
|------------|--------------|--------|
| `ibge.localidades-municipios` | `/api/v1/localidades/municipios` | **P0 — full pipeline** |
| `ibge.localidades-ufs` | `/api/v1/localidades/estados` | **P0 — full pipeline** |
| `ibge.localidades-regioes` | `/api/v1/localidades/regioes` | **P1 — ingest** |
| `ibge.localidades-mesorregioes` | `/api/v1/localidades/mesorregioes` | Catalog + ingest |
| `ibge.localidades-microrregioes` | `/api/v1/localidades/microrregioes` | Catalog + ingest |

**Fonte oficial:** [IBGE API de Localidades](https://servicodados.ibge.gov.br/api/docs/localidades)

### IBGE — PAM Produção Agrícola Municipal (Phase 16)

| Dataset ID | SIDRA table | Status |
|------------|-------------|--------|
| `ibge.pam-area-quantidade` | 1612 — área plantada, colhida, quantidade | **P0 — implemented** |
| `ibge.pam-rendimento-valor` | 1613 — rendimento médio, valor da produção | **P1 — implemented** |
| `ibge.pam-estabelecimentos` | 5457 — número de estabelecimentos | **P2 — implemented** |

**Fonte oficial:** [IBGE SIDRA — PAM](https://sidra.ibge.gov.br/pesquisa/pam) · API: [apisidra.ibge.gov.br](https://apisidra.ibge.gov.br/)

MVP crops (soja, milho, trigo) use SIDRA classification `c81` (1612), `c82` (1613), `c782` (5457) with product codes `2713`, `2711`, `2716`.

### INMET — Clima Histórico (Phase 17)

| Dataset ID | Source | Status |
|------------|--------|--------|
| `inmet.estacoes-automaticas` | Portal catálogo CSV | **P0 — implemented** |
| `inmet.bdmep-diario` | BDMEP annual ZIP (`/uploads/dadoshistoricos/{year}.zip`) | **P0 — implemented** |
| `inmet.estacoes-convencionais` | Portal catálogo CSV | **P1 — implemented** |
| `inmet.bdmep-mensal` | Monthly rollups from annual ZIP | **P2 — implemented** |
| `inmet.pacote-anual-automaticas` | BDMEP bulk annual ZIP | **P1 — implemented** |

**Fonte oficial:** [BDMEP — INMET](https://bdmep.inmet.gov.br/) · Portal: [portal.inmet.gov.br/dadoshistoricos](https://portal.inmet.gov.br/dadoshistoricos)

Timestamps in source files are **UTC**; missing values use sentinels `9999`, `Null`, or blank per INMET documentation.

### BCB — Séries Macroeconômicas SGS (Phase 18)

| Dataset ID | SGS code | Status |
|------------|----------|--------|
| `bcb.sgs-ipca` | 433 — IPCA variação mensual (%) | **P0 — implemented** |
| `bcb.sgs-ptax-usd-venda` | 1 — Dólar PTAX venda | **P0 — implemented** |
| `bcb.sgs-ipca-12m` | 13522 — IPCA acumulado 12 meses | **P1 — implemented** |
| `bcb.sgs-igpm` | 189 — IGP-M variação mensual | **P2 — implemented** |
| `bcb.sgs-ptax-usd-compra` | 10813 — Dólar PTAX compra | **P2 — implemented** |

**Fonte oficial:** [BCB Dados Abertos](https://dadosabertos.bcb.gov.br/) · API: [api.bcb.gov.br](https://api.bcb.gov.br/)

Historical backfill paginates `dataInicial`/`dataFinal` in ≤10-year chunks per BCB API limits.

### CEPEA — Preços Agro (Phase 19)

| Dataset ID | Indicator | Status |
|------------|-----------|--------|
| `cepea.soja-paranagua` | Soja — Paranaguá port (R$/sc 60 kg) | **P0 — implemented** |
| `cepea.soja-parana` | Soja — Paraná regional | **P1 — implemented** |
| `cepea.milho` | Milho — Campinas | **P1 — implemented** |
| `cepea.boi-gordo` | Boi gordo — São Paulo | **P2 — implemented** |

**Fonte oficial:** [CEPEA/ESALQ-USP](https://www.cepea.org.br/) · **License:** [CC BY-NC 4.0](https://www.cepea.org.br/br/licenca-de-uso-de-dados.aspx) — market reference (`fonte_tipo=referencia_mercado`), not `.gov.br`.

Programmatic ingest tries the CEPEA portal first; when Cloudflare blocks access, it falls back to the Notícias Agrícolas mirror (same CEPEA indicators). Full historical backfill from 2010 uses CEPEA “Consulta ao Banco de Dados” export; live ingest captures the latest published window.

Crossing with CONAB local prices (Phase 11) and BCB PTAX (Phase 18) is planned in analytics — see `.local/phases/DATA-CROSSING-VISION.md`.

---

## Rules

- Secondary blogs and unofficial mirrors are **not** acceptable as sole citation.
- Store resolved download URL + `discovered_at` per ingest — portal section link is mandatory.

---

## Related

- [VISION.md](VISION.md)
- [GLOSSARY.md](GLOSSARY.md)
- [ROADMAP.md](ROADMAP.md)
