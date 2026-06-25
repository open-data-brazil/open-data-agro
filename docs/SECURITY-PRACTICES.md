# Security practices

> Guidance for **maintainers**, **contributors**, and **integrators** of Open Data Agro.
> Reporting vulnerabilities: [SECURITY.md](../SECURITY.md) (private — not public issues).

---

## Threat model

This project combines **offline data embeds** and **maintainer fetch scripts**. Primary risks:

| Threat | Impact | Mitigation |
|--------|--------|------------|
| **Incorrect normalization** | Wrong municipality or crop codes → bad analytics | Official schemas + golden tests |
| **Stale data** | Decisions on outdated harvest figures | Metadata + refresh bot + drift logs |
| **Secrets in repo** | API key leak, supply chain compromise | `.env` gitignored, CI secrets only |
| **Tampered embeds** | Malicious data in published package | Checksums, PR review, reproducible fetch |
| **PII in tests** | Leaked farm owner data | Synthetic vectors only |
| **Misleading docs** | Insecure integration patterns | SECURITY-PRACTICES + API-CONTRACT |

---

## Maintainer checklist

### Every PR touching `packages/core` or `data/`

- [ ] Matches [OFFICIAL-SOURCES.md](OFFICIAL-SOURCES.md)
- [ ] Golden vector test added/updated
- [ ] No real PII in fixtures
- [ ] CHANGELOG entry (Security section if data-integrity fix)
- [ ] Version impact per [VERSIONING.md](VERSIONING.md)

### Every release

- [ ] `npm audit` / dependency review (when packages exist)
- [ ] Tag matches CHANGELOG version
- [ ] GitHub Release notes include security-relevant changes
- [ ] No secrets in published tarball

### Official schema updates

- [ ] Monitor agency notices (MAPA, IBGE, CONAB)
- [ ] Patch within SLA in [SECURITY.md](../SECURITY.md)
- [ ] Migration note if behavior changes

---

## Contributor checklist

- Never commit `.env`, tokens, or real PII
- Use `Signed-off-by` for commits (DCO recommended)
- Security bugs → private advisory first ([SECURITY.md](../SECURITY.md))
- MIT-compatible code only ([OPEN-SOURCE.md](OPEN-SOURCE.md))

---

## Integrator guidance

How to use this project **safely** in your application:

1. **Pin versions** — use exact semver in production until `1.0.0` API freeze
2. **Read metadata** — check `capturadoEm` before trusting embeds for compliance decisions
3. **Do not bypass core** — re-fetching government APIs at runtime duplicates drift risk unless you own SLAs
4. **Validate joins** — always use IBGE codes from this library, not ad-hoc string matching
5. **Report data bugs privately** — wrong official mapping is a security-class issue per [SECURITY.md](../SECURITY.md)

---

## Related

| Document | Purpose |
|----------|---------|
| [SECURITY.md](../SECURITY.md) | Vulnerability reporting |
| [OPEN-SOURCE.md](OPEN-SOURCE.md) | License and contribution policy |
| [OFFICIAL-SOURCES.md](OFFICIAL-SOURCES.md) | Source traceability |
