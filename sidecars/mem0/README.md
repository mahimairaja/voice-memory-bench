# vbench-mem0

Python sidecar for [vbench](../../README.md). Wraps the Mem0 OSS SDK behind
an HTTP contract the vbench Go engine drives.

## Contract

- Listens on `127.0.0.1:$VBENCH_SIDECAR_PORT`.
- Reads provider config from `$VBENCH_PROVIDER_CONFIG` (JSON).
- Endpoints: `GET /health`, `GET /capabilities`, `POST /add_message`,
  `POST /add_fact`, `POST /search`, `POST /enumerate`, `POST /reset`.
- Unsupported capabilities return `422` with body
  `{"error_type":"capability_not_supported", "capability":"...", "message":"..."}`.

Latency is measured by the Go engine around each HTTP call; this sidecar does
not try to be authoritative about timing.

## Dev

```bash
cd sidecars/mem0
uv sync
VBENCH_SIDECAR_PORT=8765 VBENCH_PROVIDER_CONFIG='{"postgres_url":"postgresql://..."}' uv run vbench-mem0
```
