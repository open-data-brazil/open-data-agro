#!/usr/bin/env bash
# Git publish helper for the daily source health bot (mirrors doc-raiz data-refresh-publish.mjs).
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "${ROOT}"

configure_git() {
  git config user.name "github-actions[bot]"
  git config user.email "41898282+github-actions[bot]@users.noreply.github.com"
}

read_commit_message() {
  python3 - <<'PY'
import json
from pathlib import Path

latest = Path("data/source-health-reports/latest.json")
if not latest.exists():
    raise SystemExit("missing data/source-health-reports/latest.json — run source-health-bot first")
report = json.loads(latest.read_text(encoding="utf-8"))
print(report.get("commitMessage", "chore(source-health): daily probe"))
PY
}

try_push_main() {
  git push origin HEAD:main
}

create_branch_name() {
  date -u +"%Y-%m-%d"
}

open_pull_request() {
  local branch="$1"
  local title
  title="$(read_commit_message)"
  gh pr create \
    --base main \
    --head "${branch}" \
    --title "${title}" \
    --body "$(cat data/source-health-reports/pr-body.md 2>/dev/null || echo 'Automated source health probe — see data/source-health-reports/latest.json.')"
}

main() {
  if [[ ! -f data/source-health-reports/latest.json ]]; then
    echo "No reports to publish — run make source-health-bot first."
    exit 0
  fi

  configure_git
  git add data/source-health-reports docs/SOURCE-HEALTH.md

  if git diff --staged --quiet; then
    echo "No staged changes — nothing to publish."
    exit 0
  fi

  commit_msg="$(read_commit_message)"
  git commit -m "${commit_msg}"

  if try_push_main; then
    echo "Direct push to main succeeded."
    exit 0
  fi

  branch="bot/source-health-reports/$(create_branch_name)"
  git checkout -b "${branch}"
  git push -u origin "${branch}"
  echo "Branch pushed: ${branch}"

  if pr_url="$(open_pull_request "${branch}" 2>&1)"; then
    echo "Opened PR: ${pr_url}"
    gh pr merge "${branch}" --auto --merge || true
    exit 0
  fi

  echo "PR creation skipped (Actions token may lack pull-requests:write)."
  echo "Reports are on branch ${branch} — merge manually or set SOURCE_HEALTH_GITHUB_TOKEN."
  exit 0
}

main "$@"
