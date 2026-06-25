# Open source commitment

> **Open Data Agro is and will remain 100% open source.**
> Core data modules and fetch pipelines are MIT-licensed forever — no paid tier, no proprietary fork, no "open core" trap.

---

## License

| Item | Policy |
|------|--------|
| **License** | [MIT](../LICENSE) |
| **SPDX identifier** | `MIT` |
| **Copyright** | Open Data Agro contributors |
| **Commercial use** | Allowed |
| **Modification** | Allowed |
| **Distribution** | Allowed |
| **Private use** | Allowed |
| **Patent grant** | Not explicit (standard MIT) |
| **Warranty** | None — "AS IS" |

The full license text is in [LICENSE](../LICENSE). A copy must ship with every distribution (npm tarball, fork, bundle).

---

## What "100% open source" means here

### Always free and open (this repository)

- All **embedded datasets** and normalization logic
- All **fetch and refresh scripts**
- All **tests and golden vectors**
- All **documentation** in `docs/`
- All **CI/CD configuration**

### Never in this repo

- Paywalled dataset modules
- "Enterprise edition" with exclusive agricultural data
- Source-available but non-OSI licenses for core code
- Time-limited or seat-licensed data features

### Optional future packages (still OSS)

Separate npm packages under the same org, **also MIT**:

- `@open-data-agro/adapters-*` — optional HTTP live lookups
- UI or dashboard packages

Each must have its own LICENSE (MIT) and live in public GitHub repos.

---

## Contribution licensing

| Rule | Detail |
|------|--------|
| **Inbound license** | MIT (same as project) |
| **CLA** | Not required |
| **DCO** | Recommended — sign commits with `Signed-off-by` for traceability |
| **Third-party code** | Must be MIT-compatible or public domain; document in PR |
| **Copied schemas** | Implement from official specs — do not copy GPL source code verbatim |

By submitting a PR, you grant the project permission to include your work under MIT.

---

## Dependencies policy

| Dependency type | Requirement |
|-----------------|-------------|
| **Runtime** | Minimal dependencies; permissive license only (MIT, BSD, Apache-2.0, ISC) |
| **Dev** | Same; audited in CI |
| **Forbidden** | GPL/AGPL in runtime path (copyleft conflicts with maximal adoption) |

---

## Trademarks

- Project name **"Open Data Agro"** — use freely for attribution
- Do not imply official endorsement by MAPA, IBGE, CONAB, INMET, or government agencies

---

## Forking

You may fork freely under MIT terms. We encourage:

- Attribution to this project
- Upstream PRs for data integrity fixes (especially security)
- Distinct naming if fork diverges significantly

---

## Related documents

| Document | Purpose |
|----------|---------|
| [LICENSE](../LICENSE) | Legal text |
| [CONTRIBUTING.md](../CONTRIBUTING.md) | How to contribute |
| [SECURITY.md](../SECURITY.md) | Vulnerability reporting |
| [VERSIONING.md](VERSIONING.md) | Release and support policy |
| [VISION.md](VISION.md) | Product mission |

---

## FAQ

**Can I use this in a commercial agtech product?**  
Yes. MIT allows commercial use without payment to us.

**Will a "pro" version with more datasets exist?**  
Not from this project. All planned datasets ship in open source.

**What if CONAB or IBGE changes schemas?**  
We patch open source per [VERSIONING.md](VERSIONING.md). Fixes are never paywalled.

**Can I dual-license my fork?**  
Your fork can add licenses only for **your new code**; MIT portions remain under MIT.
