# Voice Benchmark Design (Phase 2)

> **Status:** Scaffolded. Not yet implemented.

## The Problem with Text Benchmarks

LoCoMo and LongMemEval measure answer quality on text chat transcripts. This is necessary but not sufficient for voice agents because:

1. **Latency is on the critical path.** In voice agents, memory retrieval happens between ASR completion and TTS start. The budget is ~300-500 ms for natural-feeling interaction. A benchmark that ignores latency cannot predict production behaviour.

2. **Concurrent load changes everything.** A single-threaded benchmark hides contention. Real deployments serve hundreds of concurrent callers; tail latency under load (p99) is the metric that matters.

3. **Voice-shaped discourse is different.** Voice turns are shorter, more frequent, and contain disfluencies ("um", "uh"). Embedding models trained on prose may perform differently on transcribed speech.

4. **Session isolation is a correctness requirement.** Memories from one caller must never appear in another's context. Text benchmarks don't test this.

## Phase 2 Metrics

| Metric | Description |
|--------|-------------|
| Write p50/p95/p99 | Turn write latency under 1x, 4x, 16x load |
| Retrieval p50/p95/p99 | Query latency under concurrent load |
| Token footprint | p50 token count of injected memory payload |
| Isolation failures | Count of cross-user memory leaks |
| Cold vs warm recall | Quality delta on first call vs returning caller |

## Architecture

The voice benchmark layer lives in `voice/`. It drives the same provider adapters as Phase 1, but wraps them in a LiveKit Agents worker that simulates real-time call traffic.

See `voice/pipeline/` for the stage implementations.
