# Roadmap

## v0.1 — MVP (current)

Go engine + Python sidecar for Mem0, LoCoMo at 1x + 4x, voice verdict output.

- [x] Go engine: schema, adapter HTTP client, sidecar lifecycle, dataset
      loader, OpenAI answer/judge, concurrency runner, 5-stage pipeline,
      voice verdict + JSON report.
- [x] Python sidecar: `sidecars/mem0` (FastAPI over `mem0ai`).
- [x] Mem0 docker-compose (Postgres + pgvector).
- [x] LoCoMo loader with HuggingFace download + multi-session flattening.
- [ ] First published smoke-test results against Mem0.
- [ ] Hardened LoCoMo schema coverage (alternative upstream shapes).

## v0.2 — More providers, more load

- [ ] Memori sidecar (Postgres-backed).
- [ ] Graphiti sidecar (FalkorDB-backed).
- [ ] Cognee sidecar (SQLite + LanceDB + Kuzu).
- [ ] 16x concurrency level with back-pressure modelling.
- [ ] LongMemEval dataset loader.
- [ ] Custom JSONL dataset loader.
- [ ] `vbench compare` command for cross-provider reports.

## v0.3 — Voice-realistic runtime

- [ ] LiveKit Agents worker that drives the adapter during a live call.
- [ ] Deepgram STT / Cartesia TTS latency pass-through so the reported p95
      reflects the STT → memory → LLM → TTS loop, not just the memory step.
- [ ] Concurrent-session isolation audit.
- [ ] Voice-specific MemScore extensions (turn-level relevance, cold-caller
      warm-up latency).

## v0.4 — Cloud-mode + continuous

- [ ] Cloud-managed provider endpoints behind the same HTTP contract (config
      flip; no engine changes).
- [ ] GitHub Actions scheduled nightly benchmarks.
- [ ] Public leaderboard.
- [ ] Regression detection across provider SDK versions.
