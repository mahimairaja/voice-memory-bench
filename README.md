# voice-memory-bench

> *LongMemEval tells you which memory framework is smartest. This benchmark tells you which one your voice agent can actually afford.*

## Why this exists

Existing memory benchmarks (LoCoMo, LongMemEval, ConvoMem) measure answer quality on text chat transcripts. They do not measure what matters for real-time voice agents: p50/p95/p99 retrieval latency under concurrent load, token footprint of injected memory, write-path latency during mid-call persistence, cold vs warm recall for returning callers, concurrent session isolation, or turn-level relevance over voice-shaped discourse — short turns, disfluencies, interruptions, topic drift.

`voice-memory-bench` is a Python harness that measures all of the above across self-hostable, open-source memory frameworks. It is MIT-licensed, fully reproducible, and produces MemScore triples (quality, latency, cost) rather than a single scalar.

## Supported Providers

| Provider | Backend | Self-host | Status |
|----------|---------|-----------|--------|
| [Mem0 OSS](https://github.com/mem0ai/mem0) | pgvector + Postgres | ✅ | Scaffolded |
| [Memori (GibsonAI)](https://github.com/Gibson-AI/memori) | SQL + Postgres | ✅ | Scaffolded |
| [Graphiti](https://github.com/getzep/graphiti) | FalkorDB (knowledge graph) | ✅ | Scaffolded |
| [Cognee](https://github.com/topoteretes/cognee) | SQLite + LanceDB + Kuzu | ✅ | Scaffolded |

**Self-hostable only.** Any provider that cannot be run on your own hardware under an open-source license is out of scope.

## Quickstart

```bash
# 1. Install
pip install uv
uv sync

# 2. Download a dataset
uv run vmb datasets download locomo

# 3. Run a benchmark
uv run vmb run examples/configs/mem0-locomo.yaml
```

## Architecture

See [docs/architecture.md](docs/architecture.md) for the full pipeline, adapter interface, and artifact schema.

## Credits

This project is a sibling to [MemoryBench](https://github.com/supermemoryai/memorybench) by the [Supermemory](https://supermemory.ai) team. MemoryBench is a TypeScript benchmarking harness with the same pipeline shape (`ingest → index → search → answer → evaluate`) and the same MemScore philosophy. `voice-memory-bench` reimplements that design in Python for three concrete reasons:

1. Every self-hostable memory framework targeted here (Mem0, Memori, Graphiti, Cognee) ships a Python SDK as its canonical integration surface. A Python harness talks to providers directly; a TypeScript harness has to shell out or wrap RPC.
2. The voice agent runtime this harness simulates ([LiveKit Agents](https://github.com/livekit/agents)) is Python-first.
3. Mahimai's production stack (FastAPI, NeonDB, Hetzner, Fly.io workers) is Python-native, so this harness doubles as a deployment template.

Thank you to the Supermemory team for publishing MemoryBench openly and for the MemScore framing.

## Status

Phase 1 (text-chat benchmark) is scaffolded. Phase 2 (voice benchmark) is scaffolded but not implemented. See [ROADMAP.md](ROADMAP.md).

[![CI](https://github.com/mahimairaja/voice-memory-bench/actions/workflows/ci.yml/badge.svg)](https://github.com/mahimairaja/voice-memory-bench/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Python 3.11+](https://img.shields.io/badge/python-3.11+-blue.svg)](https://www.python.org/downloads/)
