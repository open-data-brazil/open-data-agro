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
| `conab.precos-agropecuarios-mensal-uf` | Preços agropecuários Mensal UF | Catalog registered |
| `conab.precos-agropecuarios-mensal-municipio` | Preços agropecuários Mensal Município | Catalog registered |
| `conab.precos-agropecuarios-semanal-uf` | Preços agropecuários Semanal UF | Catalog registered (P1 next) |
| `conab.precos-agropecuarios-semanal-municipio` | Preços agropecuários Semanal Municipio | Catalog registered |
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
| `ibge.localidades-municipios` | `/api/v1/localidades/municipios` | **P0 — implemented** |
| `ibge.localidades-ufs` | `/api/v1/localidades/estados` | **P0 — implemented** |
| `ibge.localidades-regioes` | `/api/v1/localidades/regioes` | **P1 — implemented** |
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

---

## Rules

- Secondary blogs and unofficial mirrors are **not** acceptable as sole citation.
- Store resolved download URL + `discovered_at` per ingest — portal section link is mandatory.

---

## Related

- [VISION.md](VISION.md)
- [GLOSSARY.md](GLOSSARY.md)
- [ROADMAP.md](ROADMAP.md)
