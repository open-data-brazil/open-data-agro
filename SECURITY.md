# Security Policy

> Open Data Agro is a **100% open-source** project (MIT).
> Security issues in data pipelines, published datasets, or integration adapters are **critical** — incorrect or tampered agricultural data undermines every downstream system.

---

## Supported versions

| Version | Supported |
|---------|-----------|
| Latest stable (`1.x`) | Yes |
| Previous minor (`1.(x-1)`) | Security fixes only |
| Pre-1.0 (`0.x`) | Best effort during alpha/beta |
| Unreleased `main` | Fix forward; no backports unless critical |

See [docs/VERSIONING.md](docs/VERSIONING.md) for release and support policy.

---

## Reporting a vulnerability

**Do not open a public GitHub issue** for security vulnerabilities.

### Preferred: GitHub Private Security Advisory

1. Go to [github.com/open-data-brazil/open-data-agro/security/advisories](https://github.com/open-data-brazil/open-data-agro/security/advisories)
2. Click **Report a vulnerability**
3. Describe the issue with reproduction steps and impact

### Alternative: email

If GitHub advisories are unavailable, email the maintainers with subject `[SECURITY] open-data-agro`:

- Include steps to reproduce
- Affected version(s)
- Impact assessment (e.g. dataset integrity, credential exposure, supply chain)

We aim to acknowledge within **48 hours** and provide an initial assessment within **5 business days**.

---

## What to report

| In scope | Out of scope |
|----------|--------------|
| Compromised fetch scripts or data refresh pipelines | Bugs in consumer apps using this project |
| Secrets or credentials committed to the repository | Social engineering |
| Supply-chain issues in published packages | Vulnerabilities in devDependencies only |
| Data integrity issues in embedded or published datasets | Missing features (use feature requests) |

---

## Severity (project-specific)

| Severity | Example | Target fix |
|----------|---------|------------|
| **Critical** | Leaked API keys, malicious data injection in published artifacts | Patch release ≤ 72h |
| **High** | Incorrect transformation of official agricultural data at scale | Patch release ≤ 7 days |
| **Medium** | Non-security correctness bug in data mapping | Next patch/minor |
| **Low** | Documentation or DX issue with security implication | Next minor |

When official agencies (MAPA, IBGE, CONAB, INMET) change schemas, we treat misalignment as **High** minimum until patched.

---

## Disclosure policy

- **Coordinated disclosure** — we request 90 days before public disclosure unless fix is released sooner.
- Credit given in [CHANGELOG.md](CHANGELOG.md) and GitHub advisory (unless reporter prefers anonymity).
- CVE requested for Critical/High when applicable.

---

## Security practices (project)

- Dependencies pinned and audited in CI (when scaffold exists)
- No secrets in committed data or scripts
- Golden test vectors from [official sources](docs/OFFICIAL-SOURCES.md)
- All contributions must be MIT-compatible — see [docs/OPEN-SOURCE.md](docs/OPEN-SOURCE.md)
- See [docs/SECURITY-PRACTICES.md](docs/SECURITY-PRACTICES.md) for integrator guidance

---

## Secure contribution

Contributors must **not** include:

- Real farm owner PII, CAR numbers tied to individuals, or production secrets in tests
- Secrets, API keys, or `.env` files
- Code under incompatible licenses

See [CONTRIBUTING.md](CONTRIBUTING.md#security-contributions).
