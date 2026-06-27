#!/usr/bin/env python3
"""Live HTTP probes for Phase 51 wave 5 discovery — approved primary URLs only."""

from __future__ import annotations

import argparse
import json
import sys
import urllib.error
import urllib.request
from dataclasses import dataclass

USER_AGENT = "Mozilla/5.0 (compatible; OpenDataAgro/1.0; +https://github.com/open-data-brazil/open-data-agro)"
TIMEOUT_SEC = 30


@dataclass(frozen=True)
class ProbeSpec:
    dataset_id: str
    url: str
    expect_status: int = 200
    min_bytes: int = 100


# Approved wave 5 candidates — primary URLs verified 2026-06-27 (Phase 51).
APPROVED_PROBES: list[ProbeSpec] = [
    ProbeSpec(
        "mapa.sipeagro-estabelecimentos",
        "https://dados.agricultura.gov.br/api/3/action/package_show?id=sipeagro",
    ),
    ProbeSpec(
        "mapa.sipeagro-produtos",
        "https://dados.agricultura.gov.br/api/3/action/package_show?id=sipeagro",
    ),
    ProbeSpec(
        "mapa.sigef-producao-sementes",
        "https://dados.agricultura.gov.br/api/3/action/package_show?id=dados-referentes-ao-controle-da-producao-de-sementes-sigef",
    ),
    ProbeSpec(
        "mapa.sigef-areas",
        "https://dados.agricultura.gov.br/api/3/action/package_show?id=dados-referentes-ao-controle-da-producao-de-sementes-sigef",
    ),
    ProbeSpec(
        "mapa.sisser-seguro-rural",
        "https://dados.agricultura.gov.br/api/3/action/package_show?id=sisser3",
    ),
    ProbeSpec(
        "ibge.ppm-efetivo-rebanhos",
        "https://apisidra.ibge.gov.br/values/t/3939/n3/11/p/last%201/v/all",
    ),
    ProbeSpec(
        "ibge.ppm-vacas-ordenhadas",
        "https://apisidra.ibge.gov.br/values/t/94/n3/11/p/last%201/v/all",
    ),
    ProbeSpec(
        "ibge.ppm-ovinos-tosquiados",
        "https://apisidra.ibge.gov.br/values/t/95/n3/11/p/last%201/v/all",
    ),
    ProbeSpec(
        "ibge.ppm-aquicultura",
        "https://apisidra.ibge.gov.br/values/t/3940/n3/11/p/last%201/v/all",
    ),
    ProbeSpec(
        "ibge.pam-precos-produtor",
        "https://apisidra.ibge.gov.br/values/t/1612/n6/in%20n3%2011/p/2023/v/214/c81/2713",
    ),
    ProbeSpec(
        "ibge.pam-culturas-estendidas",
        "https://apisidra.ibge.gov.br/values/t/1612/n6/in%20n3%2011/p/2023/v/109,216,214/c81/2710",
    ),
    ProbeSpec(
        "ibge.lspa-rendimento-medio",
        "https://apisidra.ibge.gov.br/values/t/6588/n3/11/p/202312/v/35/c48/39443",
    ),
    ProbeSpec(
        "ibge.censo-agro-area-uso-solo",
        "https://apisidra.ibge.gov.br/values/t/6879/n3/11/p/2017/v/all",
    ),
    ProbeSpec(
        "ibge.censo-agro-maquinario",
        "https://apisidra.ibge.gov.br/values/t/6880/n3/11/p/2017/v/all",
    ),
    ProbeSpec(
        "ibge.pnad-rural-renda-ocupacao",
        "https://apisidra.ibge.gov.br/values/t/6385/n3/11,12,13,14,15,16,17/p/last%201/v/all",
    ),
    ProbeSpec(
        "ibama.sisfogo-incendios",
        "https://dadosabertos.ibama.gov.br/dados/SISFOGO/ROI.csv",
        min_bytes=1000,
    ),
    ProbeSpec(
        "ibama.licencas-ambientais",
        "https://dadosabertos.ibama.gov.br/api/3/action/package_show?id=licencas-ambientais-de-atividades-e-empreendimentos-licenciados-pelo-ibama",
    ),
    ProbeSpec(
        "embrapa.agroapi-agrofit",
        "https://www.embrapa.br/agroapi",
    ),
    ProbeSpec(
        "bcb.cim-agro-credito-rural",
        "https://api.bcb.gov.br/dados/serie/bcdata.sgs.21087/dados?formato=json&dataInicial=01/01/2024&dataFinal=01/02/2024",
        min_bytes=50,
    ),
    ProbeSpec(
        "abiove.balanco-complexo-soja",
        "https://abiove.org.br/abiove_content/Abiove/exp_202605.xlsx",
        min_bytes=10000,
    ),
    ProbeSpec(
        "b3.futuro-cafe",
        "https://www.b3.com.br/pesquisapregao/download?filelist=SPRD250627.zip",
        min_bytes=1000,
    ),
    ProbeSpec(
        "transportes.mtr-bit-malha-shapefile",
        "https://www.gov.br/transportes/pt-br/assuntos/dados-de-transportes/bit/Base-GEO/BaseFerro.zip",
        min_bytes=1000,
    ),
]

# Expected blockers — probe must NOT return 200 (confirms deferral still valid).
DEFERRED_BLOCKERS: list[ProbeSpec] = [
    ProbeSpec(
        "antaq.movimentacao-carga-portuaria",
        "https://web3.antaq.gov.br/ea/sense/download.html",
        expect_status=404,
        min_bytes=0,
    ),
    ProbeSpec(
        "fao.comercio-agro",
        "https://bulks-faostat.fao.org/production/Trade_Crops_Livestock_E_All_Data_(Normalized).zip",
        expect_status=403,
        min_bytes=0,
    ),
    ProbeSpec(
        "wto.its-trade-statistics",
        "https://api.wto.org/timeseries/v1/data",
        expect_status=401,
        min_bytes=0,
    ),
]


def probe_url(spec: ProbeSpec) -> tuple[int, int]:
    request = urllib.request.Request(spec.url, headers={"User-Agent": USER_AGENT})
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


def check_approved(spec: ProbeSpec) -> str | None:
    status, size = probe_url(spec)
    if status != spec.expect_status:
        return f"{spec.dataset_id}: expected HTTP {spec.expect_status}, got {status} ({spec.url})"
    if size < spec.min_bytes:
        return f"{spec.dataset_id}: body too small ({size} bytes, min {spec.min_bytes})"
    return None


def check_deferred(spec: ProbeSpec) -> str | None:
    status, _size = probe_url(spec)
    if status == 200:
        return (
            f"{spec.dataset_id}: expected blocker (not 200), got HTTP 200 — "
            "re-audit for Phase 57 re-enable"
        )
    return None


def main() -> int:
    parser = argparse.ArgumentParser(description="Verify Phase 51 wave 5 discovery probes")
    parser.add_argument("--json", action="store_true", help="Emit JSON report to stdout")
    args = parser.parse_args()

    failures: list[str] = []
    results: list[dict[str, object]] = []

    for spec in APPROVED_PROBES:
        status, size = probe_url(spec)
        ok = status == spec.expect_status and size >= spec.min_bytes
        results.append(
            {
                "dataset_id": spec.dataset_id,
                "url": spec.url,
                "status": status,
                "bytes": size,
                "kind": "approved",
                "ok": ok,
            }
        )
        err = check_approved(spec)
        if err:
            failures.append(err)

    for spec in DEFERRED_BLOCKERS:
        status, size = probe_url(spec)
        ok = status != 200
        results.append(
            {
                "dataset_id": spec.dataset_id,
                "url": spec.url,
                "status": status,
                "bytes": size,
                "kind": "deferred_blocker",
                "ok": ok,
            }
        )
        err = check_deferred(spec)
        if err:
            failures.append(err)

    if args.json:
        print(json.dumps({"failures": failures, "results": results}, indent=2))
    else:
        approved_ok = sum(1 for r in results if r["kind"] == "approved" and r["ok"])
        print(f"wave5 discovery probe: {approved_ok}/{len(APPROVED_PROBES)} approved URLs OK")
        for row in results:
            mark = "OK" if row["ok"] else "FAIL"
            print(f"  [{mark}] {row['dataset_id']}: HTTP {row['status']} ({row['bytes']} bytes)")
        for msg in failures:
            print(f"  ERROR: {msg}", file=sys.stderr)

    return 1 if failures else 0


if __name__ == "__main__":
    sys.exit(main())
