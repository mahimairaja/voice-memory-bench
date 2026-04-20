"""Entrypoint for the vbench-mem0 sidecar."""

from __future__ import annotations

import os
import sys

import uvicorn

from .app import build_app


def run() -> None:
    port_raw = os.environ.get("VBENCH_SIDECAR_PORT")
    if not port_raw:
        print("VBENCH_SIDECAR_PORT is required", file=sys.stderr)
        sys.exit(2)
    try:
        port = int(port_raw)
    except ValueError:
        print(f"VBENCH_SIDECAR_PORT must be an integer, got {port_raw!r}", file=sys.stderr)
        sys.exit(2)

    app = build_app()
    uvicorn.run(app, host="127.0.0.1", port=port, log_level="warning")


if __name__ == "__main__":
    run()
