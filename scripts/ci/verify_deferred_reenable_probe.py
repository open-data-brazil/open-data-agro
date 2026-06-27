#!/usr/bin/env python3
"""Phase 57 live audit for deferred catalog entries in configs/catalog/_deferred/."""

from __future__ import annotations

import argparse
import json
import os
import sys
import urllib.error
import urllib.request
from dataclasses import dataclass
from pathlib import Path

import yaml

USER_AGENT = "Mozilla/5.0 (compatible; OpenDataAgro/1.0; +https://github.com/open-data-brazil/open-data-agro)"
TIMEOUT_SEC = 30
REPO_ROOT = Path(__file__).resolve().parents[2]
DEFERRED_YAML = REPO_ROOT / "configs/catalog/_deferred/unreachable_sources.yaml"


@dataclass(frozen=True)
class ProbeSpec:
    dataset_id: str
    url: str
    kind: str  # deferred_yaml | discovery_deferral
    min_bytes: int = 100
    require_key_env: str = ""
    portal_only: bool = False


DISCOVERY_DEFERRALS: list[ProbeSpec] = [
    ProbeSpec(
        "antt.tarifas-pedagio",
        "https://dados.antt.gov.br/api/3/action/package_show?id=tarifas-praca-pedagio",
        "discovery_deferral",
    ),
    ProbeSpec(
        "inpe.cptec-indices-climaticos",
        "https://ftp.cptec.inpe.br/modelos/tempo/sazonal/",
        "discovery_deferral",
    ),
    ProbeSpec(
        "embrapa.solos-brasil",
        "https://www.embrapa.br/solos-de-terras-brasileiras",
        "discovery_deferral",
        portal_only=True,
    ),
]

PSD_WSDL = "https://apps.fas.usda.gov/PSDExternalAPIService/svcPSD_AMIS.asmx?WSDL"


def probe_url(url: str) -> tuple[int, int]:
    request = urllib.request.Request(url, headers={"User-Agent": USER_AGENT})
    api_key_env = os.environ.get("USDA_FAS_API_KEY", "").strip()
    if api_key_env and "fas.usda.gov/OpenData" in url:
        request.add_header("API_KEY", api_key_env)
    wto_key = os.environ.get("WTO_API_KEY", "").strip()
    if wto_key and "api.wto.org" in url:
        request.add_header("Ocp-Apim-Subscription-Key", wto_key)

    try:
        with urllib.request.urlopen(request, timeout=TIMEOUT_SEC) as response:
            content_length = response.headers.get("Content-Length")
            if content_length is not None and content_length.isdigit():
                size = int(content_length)
            else:
                body = response.read(65536)
                size = len(body)
            return response.status, size
    except urllib.error.HTTPError as exc:
        chunk = exc.read(8192)
        return exc.code, len(chunk)
    except urllib.error.URLError:
        return 0, 0


def load_deferred_specs() -> list[ProbeSpec]:
    data = yaml.safe_load(DEFERRED_YAML.read_text(encoding="utf-8"))
    specs: list[ProbeSpec] = []
    for entry in data.get("entries", []):
        dataset_id = entry["dataset_id"]
        if dataset_id.startswith("usda.psd-"):
            continue
        url = entry.get("fao_bulk_url") or entry["source_url"]
        require_key = ""
        portal_only = False
        min_bytes = 100

        if dataset_id == "mexico.siap-produccion-agricola":
            portal_only = True
            min_bytes = 0
        if dataset_id == "noaa.gpcc-precipitation":
            portal_only = True
        if dataset_id == "usda.gats-trade":
            require_key = "USDA_FAS_API_KEY"
        if dataset_id == "wto.its-trade-statistics":
            require_key = "WTO_API_KEY"
        if dataset_id == "fao.comercio-agro":
            min_bytes = 10000

        specs.append(
            ProbeSpec(
                dataset_id=dataset_id,
                url=url,
                kind="deferred_yaml",
                min_bytes=min_bytes,
                require_key_env=require_key,
                portal_only=portal_only,
            )
        )
    return specs


def evaluate(spec: ProbeSpec, status: int, size: int) -> tuple[str, bool, str | None]:
    """Return (verdict, ok, error_message)."""
    if spec.portal_only:
        if status == 200:
            return "still_deferred_portal_only", True, None
        if status in (404, 403, 401, 0):
            return "still_blocked", True, None
        return "unexpected_status", False, f"{spec.dataset_id}: portal probe HTTP {status}"

    if spec.require_key_env:
        if os.environ.get(spec.require_key_env, "").strip():
            if status == 200 and size >= spec.min_bytes:
                return "re_enable_candidate", True, None
            if status == 200:
                return "re_enable_candidate", True, None
        if status in (401, 403):
            return "still_blocked_no_key", True, None
        if status == 200:
            return "re_enable_candidate", False, (
                f"{spec.dataset_id}: HTTP 200 without {spec.require_key_env} — verify key path"
            )
        return "still_blocked", True, None

    if spec.dataset_id.startswith("usda.psd-"):
        return "still_blocked", status != 200, (
            None
            if status != 200
            else f"{spec.dataset_id}: PSD WSDL now HTTP 200 — re-audit for Phase 57 re-enable"
        )

    if status == 200 and size >= spec.min_bytes:
        return "re_enable_candidate", False, (
            f"{spec.dataset_id}: bulk URL now HTTP 200 ({size} bytes) — run re-enable checklist"
        )

    if status == 200:
        return "re_enable_candidate", False, (
            f"{spec.dataset_id}: HTTP 200 but body too small ({size} bytes)"
        )

    return "still_blocked", True, None


def main() -> int:
    parser = argparse.ArgumentParser(description="Phase 57 deferred re-enable live audit")
    parser.add_argument("--json", action="store_true", help="Emit JSON report")
    args = parser.parse_args()

    specs = load_deferred_specs()
    for psd_id in ("usda.psd-soja", "usda.psd-milho", "usda.psd-trigo"):
        specs.append(
            ProbeSpec(
                dataset_id=psd_id,
                url=PSD_WSDL,
                kind="deferred_yaml",
                min_bytes=500,
            )
        )
    specs.extend(DISCOVERY_DEFERRALS)

    failures: list[str] = []
    results: list[dict[str, object]] = []
    re_enable_count = 0

    for spec in specs:
        status, size = probe_url(spec.url)
        verdict, ok, err = evaluate(spec, status, size)
        if verdict == "re_enable_candidate":
            re_enable_count += 1
        results.append(
            {
                "dataset_id": spec.dataset_id,
                "url": spec.url,
                "status": status,
                "bytes": size,
                "kind": spec.kind,
                "verdict": verdict,
                "ok": ok,
            }
        )
        if err:
            failures.append(err)

    if args.json:
        print(
            json.dumps(
                {
                    "re_enable_candidates": re_enable_count,
                    "failures": failures,
                    "results": results,
                },
                indent=2,
            )
        )
    else:
        blocked = sum(1 for r in results if r["verdict"].startswith("still_"))
        print(
            f"deferred re-enable audit: {blocked}/{len(results)} still blocked, "
            f"{re_enable_count} re-enable candidate(s)"
        )
        for row in results:
            mark = "OK" if row["ok"] else "FAIL"
            print(
                f"  [{mark}] {row['dataset_id']}: HTTP {row['status']} "
                f"({row['bytes']} bytes) — {row['verdict']}"
            )
        for msg in failures:
            print(f"  ACTION: {msg}", file=sys.stderr)

    return 1 if failures else 0


if __name__ == "__main__":
    sys.exit(main())
