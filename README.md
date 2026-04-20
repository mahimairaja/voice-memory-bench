# vbench

> *LongMemEval tells you which memory framework is smartest. vbench tells you which one your voice agent can actually ship.*

`vbench` is a Go + Python benchmark harness that measures whether a self-hostable
memory framework (Mem0 today; Memori / Graphiti / Cognee on the roadmap) is
**voice-agent fit** — not just whether it answers the question right.

The headline output is one line per concurrency level, so a non-technical user
or an agent downstream can act on it without parsing tables:

```
mem0 @ 1x: EXCELLENT  (p95 = 210 ms, quality = 0.74, tokens = 390)
mem0 @ 4x: ACCEPTABLE (p95 = 380 ms, quality = 0.72, tokens = 420)
```

## Voice verdicts

The verdict is driven by search-stage p95 latency, which is what a voice agent
feels on the read path:

| p95 latency       | Verdict     |
|-------------------|-------------|
| &lt; 300 ms       | EXCELLENT   |
| 300 – 500 ms      | ACCEPTABLE  |
| &gt; 500 ms       | FAIL        |

Quality (judge LLM score on LoCoMo QA) and injected-memory token footprint are
reported alongside the verdict — never collapsed into a single scalar.

## Architecture

- **Engine: Go** — single static binary, no GIL jitter polluting p99 numbers,
  first-class Cobra CLI, trivial distribution for agents.
- **Providers: Python sidecars** over HTTP — every target provider (Mem0 today)
  ships a Python SDK as its canonical integration surface. The HTTP boundary
  mirrors the cloud API shape, so pointing at a cloud endpoint later is a
  config flip, not a rewrite.
- **Latency is authoritative in the Go engine.** The sidecar self-reports a
  latency for reference, but the voice verdict uses the Go-side wall clock.

More detail: [`docs/architecture.md`](docs/architecture.md).

## MVP scope

| Axis         | In scope                              | Deferred to v0.2                     |
|--------------|---------------------------------------|--------------------------------------|
| Providers    | Mem0 OSS                              | Memori, Graphiti, Cognee             |
| Datasets     | LoCoMo                                | LongMemEval, custom JSONL            |
| Concurrency  | 1x, 4x                                | 16x                                  |
| Hosting      | Self-hosted                           | Cloud-managed (config flip)          |

## Quickstart

```bash
# 1. Start the backing store for Mem0.
export POSTGRES_PASSWORD=changeme
make providers-up-mem0

# 2. Sync the Mem0 sidecar environment.
make sidecar-sync

# 3. Build the engine.
make build

# 4. Download LoCoMo.
./vbench datasets download locomo

# 5. Run a smoke test.
export OPENAI_API_KEY=sk-...
./vbench eval --config examples/configs/mem0-locomo.yaml --max-items 2
```

## Layout

```
cmd/vbench/         # Cobra CLI (root, eval, datasets, providers)
internal/
  schema/           # Benchmark, artifacts, MemScore, RunConfig
  adapter/          # HTTP client to the sidecar + error envelope
  sidecar/          # Subprocess lifecycle + Mem0 defaults
  dataset/          # Loader interface + LoCoMo
  llm/              # OpenAI chat-completion answer + judge client
  concurrency/      # Worker pool + percentile helper
  engine/           # Pipeline + 5 stages (ingest → index → search → answer → evaluate)
  report/           # Voice verdict + JSON report
sidecars/mem0/      # Python FastAPI sidecar (vbench-mem0 entrypoint)
providers/mem0/     # Docker compose for Postgres/pgvector + provider README
examples/configs/   # Example run configs
docs/               # Architecture + benchmarks + providers + voice-bench notes
```

## Status

MVP: end-to-end pipeline against Mem0 + LoCoMo at 1x and 4x. See
[`ROADMAP.md`](ROADMAP.md) and [`TODO.md`](TODO.md) for what's planned next.

## Credits

vbench inherits the MemScore framing (quality / latency / cost, never
collapsed) and the five-stage pipeline shape (ingest → index → search → answer
→ evaluate) from [MemoryBench](https://github.com/supermemoryai/memorybench).
The voice-fitness emphasis and the Go + Python sidecar split are specific to
vbench.

## License

MIT — see [`LICENSE`](LICENSE).
