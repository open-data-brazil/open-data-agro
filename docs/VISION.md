# Vision — Open Data Agro

## Mission

Build a **100% open-source** toolkit that any developer can use to **embed, refresh, and query** Brazilian agricultural public data — with schemas and transformations aligned to **official primary sources**, not unofficial mirrors or blog posts.

Wrong geospatial joins or crop codes destroy trust in an agricultural data toolkit. This project treats **source traceability** and **test vectors from official examples** as first-class requirements.

## What we build

| Capability | Description |
|------------|-------------|
| **Embed** | Offline snapshots of official datasets with capture metadata |
| **Fetch** | Maintainer scripts to refresh from `.gov.br` APIs and open data portals |
| **Normalize** | Canonical codes (IBGE, MAPA, CONAB) and geospatial keys |
| **Query** (planned) | Typed lookups and joins for apps, research, and civic tech |

## Design principles

1. **Official sources only** — every dataset links to a primary reference in [OFFICIAL-SOURCES.md](OFFICIAL-SOURCES.md).
2. **Offline-first** — runtime library does not depend on live government APIs.
3. **Metadata transparency** — `capturadoEm`, source URL, row counts, and drift logs.
4. **Fail closed** — missing or stale critical data surfaces explicit errors, not silent defaults.
5. **Constants over magic** — crop codes, units, and CRS definitions in one updatable place.
6. **Reproducible refresh** — fetch scripts are versioned and testable.

## Target consumers

- Agtech dashboards and decision-support tools
- Research and academic pipelines
- Civic tech and transparency projects
- ETL and data engineering workflows
- Open-source frameworks needing BR agricultural data without vendor lock-in

## Non-goals (v1)

- **Real-time trading or market predictions**
- **Private farm operational data** (telemetry, ERP integrations)
- **Replacing official government portals** — we mirror and normalize, not compete
- **Legal interpretation** of rural property or environmental compliance (CAR/SICAR)

## License

**MIT** — confirmed. See [LICENSE](../LICENSE) and [OPEN-SOURCE.md](OPEN-SOURCE.md).

Permanently **100% open source**: no paid tier, no open core, no proprietary datasets in this repository.

## Success criteria

- Every dataset has tests with vectors from official documents or API samples
- Public API stable and documented in [API-CONTRACT.md](API-CONTRACT.md)
- Refresh pipelines log drift and source health
- Documentation index complete in [README.md](README.md)
