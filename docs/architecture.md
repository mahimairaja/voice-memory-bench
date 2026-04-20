# Architecture

## Shape

```
+---------------+       HTTP / 127.0.0.1        +------------------+
|   vbench      | ----------------------------> |  Python sidecar  |
|  (Go engine)  | <---------------------------- |  (e.g. mem0ai)   |
+---------------+                                +------------------+
      |                                                 |
      |                                                 |
      v                                                 v
+----------------+                             +----------------+
| runs/<run_id>/ |                             | Postgres +     |
|   artifacts    |                             | pgvector / ... |
+----------------+                             +----------------+
```

The Go engine is the process of record: it owns the CLI, the pipeline, the
clock used for latency numbers, and the output artifacts. The Python sidecar
is a thin HTTP wrapper around the provider SDK.

## Why Go + Python sidecars

1. **Distribution.** `vbench` is a single static binary. No interpreter,
   no venv, no version pinning on the consumer side. Agents can ship it the
   way they ship `kubectl`.
2. **p99 honesty.** Memory-provider latency is what voice agents feel. The
   Go scheduler has sub-ms GC pauses; the GIL does not. Measuring timing
   inside the Python SDK would leak interpreter jitter into the reported
   p99 and invalidate the voice verdict.
3. **Provider SDKs are Python-first.** Every target framework (Mem0, Memori,
   Graphiti, Cognee) ships a Python SDK as its canonical integration surface.
   The sidecar uses that surface directly rather than reimplementing it in Go.
4. **Cloud-mode parity.** The HTTP boundary mirrors what the cloud variants
   of these providers expose. Pointing `vbench` at a cloud endpoint later is a
   config change, not an engine rewrite.

**Latency is measured in the Go engine, never in the sidecar.** The sidecar's
self-reported `latency_ms` exists as a reference for debugging; the voice
verdict is driven by the Go-side wall clock around each HTTP call.

## Pipeline

```
ingest → index → search → answer → evaluate
```

Each stage writes a `.complete` sentinel. Re-running with `--run-id <id>`
picks up from the first stage without a sentinel, so a crash mid-index does
not reset ingestion.

| Stage     | Input                                  | Output                                                    |
|-----------|----------------------------------------|-----------------------------------------------------------|
| ingest    | Dataset (LoCoMo)                       | `runs/<id>/ingest/<item>.json`                            |
| index     | BenchmarkItems                         | `runs/<id>/index/<item>.json` (per-turn write latency)    |
| search    | BenchmarkItems + sidecar               | `runs/<id>/search/<level>x/<item>__<q>.json` per level    |
| answer    | Search artifacts                       | `runs/<id>/answer/<level>x/<item>__<q>.json`              |
| evaluate  | Search + Answer + judge LLM            | `runs/<id>/evaluate/<level>x/*.json` + `_memscore.json`   |

The search + answer + evaluate stages run once per concurrency level (MVP: 1x
and 4x). Index writes are serial — voice agents also write sequentially during
a call, and serial writes isolate the search measurement from write-path noise.

## HTTP contract (engine → sidecar)

- `GET /health` — liveness.
- `GET /capabilities` — provider metadata (supported retrieval modes, etc.).
- `POST /add_message` — write a conversation turn.
- `POST /add_fact` — write a pre-extracted fact.
- `POST /search` — retrieve memories. Unsupported modes return `422` with
  `error_type=capability_not_supported`, which the engine surfaces as a
  SKIPPED result.
- `POST /enumerate` — list all memories for a user (used for recall without
  an extra LLM round-trip).
- `POST /reset` — wipe memory for a user/session. Called once per item to
  isolate callers.

The sidecar binds `127.0.0.1:$VBENCH_SIDECAR_PORT` (a free port the engine
picks) and reads provider configuration from `$VBENCH_PROVIDER_CONFIG` (JSON).
The engine launches and tears down the subprocess via a process group, so
SIGINT on the engine cleans up the sidecar.

## MemScore

MemScore is a triple per (provider, concurrency) level. Never a scalar.

- **Quality** — mean judge-LLM score over all questions at this level.
- **Latency** — p50, p95, p99 over search-stage wall-clock latencies.
- **Cost** — USD per item (0 for self-hosted MVP).
- **Token footprint** — median token count of the injected memory payload.

The voice verdict is computed only from search-stage p95:

- `p95 < 300 ms` → EXCELLENT
- `300 ≤ p95 ≤ 500 ms` → ACCEPTABLE
- `p95 > 500 ms` → FAIL

## Layout

```
cmd/vbench/          Cobra CLI (root, eval, datasets, providers)
internal/
  schema/            RunConfig, BenchmarkItem, artifacts, MemScore
  adapter/           HTTP client, error envelope, retrieval types
  sidecar/           subprocess lifecycle, Mem0 defaults
  dataset/           loader interface, LoCoMo
  llm/               OpenAI answer + judge client
  concurrency/       worker pool, percentile helper
  engine/            pipeline + 5 stages
  report/            voice verdict + MemScore JSON writer
sidecars/mem0/       Python FastAPI sidecar (uv-managed)
providers/mem0/      docker-compose for Postgres + pgvector
examples/configs/    example run configs
```
