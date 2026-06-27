# Source health reports

Daily output from `source-health-bot` (GitHub Action or `make source-health-bot`).

| File | Purpose |
|------|---------|
| `latest.json` | Most recent probe run (all outcomes + alerts) |
| `daily/YYYY-MM-DD.json` | Archived daily runs (90-day retention) |
| `job-summary.md` | CI / human summary |
| `CRITICAL-ALERTS.md` | Links unreachable 2+ consecutive days |
| `source-health/` | Per-dataset streak state |
| `probe-outcomes/` | Per-dataset last probe detail |

See [docs/SOURCE-HEALTH.md](../docs/SOURCE-HEALTH.md) for the latest status table.

## Automation

| Item | Value |
|------|-------|
| Schedule | Daily 03:00 UTC (00:00 America/Sao_Paulo) |
| Workflow | `.github/workflows/source-health-bot.yml` |
| Local run | `make source-health-bot` |
| Publish | `make source-health-publish` (or auto in CI) |
| Direct push | Set `SOURCE_HEALTH_GITHUB_TOKEN` repo secret (PAT with `main` bypass) |
| EIA probe | Optional `EIA_API_KEY` secret for `eia.petroleum-prices` sample probe |
