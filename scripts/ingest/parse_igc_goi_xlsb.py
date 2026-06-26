#!/usr/bin/env python3
"""Parse IGC GOI xlsb workbook to JSON rows for Go ingest flatten."""

from __future__ import annotations

import json
import sys
from datetime import datetime, timedelta
from pathlib import Path

from pyxlsb import open_workbook


def excel_date(value: float) -> str:
    base = datetime(1899, 12, 30)
    return (base + timedelta(days=float(value))).strftime("%Y-%m-%d")


def main() -> int:
    if len(sys.argv) != 2:
        print("usage: parse_igc_goi_xlsb.py <path.xlsb>", file=sys.stderr)
        return 2

    path = Path(sys.argv[1])
    series = [
        ("goi", "IGC GOI", 1),
        ("wheat", "Wheat", 3),
        ("maize", "Maize", 4),
        ("soyabeans", "Soyabeans", 5),
        ("rice", "Rice", 6),
        ("barley", "Barley", 7),
    ]
    out: list[dict[str, str]] = []

    with open_workbook(path) as workbook:
        with workbook.get_sheet("GOI & Indices") as sheet:
            started = False
            for row in sheet.rows():
                values = [cell.v for cell in row]
                if not started:
                    if values and values[0] == "DATE":
                        started = True
                    continue
                if not values or values[0] is None:
                    continue
                try:
                    refdate = excel_date(values[0])
                except (TypeError, ValueError):
                    continue
                for slug, name, idx in series:
                    if idx >= len(values) or values[idx] is None:
                        continue
                    out.append(
                        {
                            "refdate": refdate,
                            "index_slug": slug,
                            "index_name": name,
                            "value": str(values[idx]),
                            "base_period": "2000-01=100",
                            "frequency": "daily",
                        }
                    )

    json.dump(out, sys.stdout)
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
