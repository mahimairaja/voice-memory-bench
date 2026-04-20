# Architecture

## Overview

`voice-memory-bench` is a five-stage pipeline:

```
ingest → index → search → answer → evaluate
```

Each stage reads from the previous stage's artifact directory and writes its own artifacts. Stages are independently resumable: re-running with the same `run_id` skips any stage whose output artifact already exists and is marked `COMPLETE`.

## Pipeline Stages

### 1. Ingest
**Input:** Dataset (LoCoMo, LongMemEval, or custom JSONL)
**Output:** `runs/<run_id>/ingest/*.jsonl` — one `BenchmarkItem` per line

Normalises dataset-specific formats into the canonical `BenchmarkItem` schema. Each item contains a conversation history and one or more evaluation questions with reference answers.

### 2. Index
**Input:** `BenchmarkItem` list from ingest
**Output:** `runs/<run_id>/index/*.json` — one `IndexArtifact` per item

Calls `adapter.add_message()` for each conversation turn (or `adapter.add_fact()` for pre-extracted facts). Records per-turn write latency. All writes must reach the backing store before the stage completes.

### 3. Search
**Input:** `BenchmarkItem` list + indexed memories
**Output:** `runs/<run_id>/search/*.json` — one `SearchArtifact` per question

Calls `adapter.search()` for each evaluation question. Records retrieval latency and the exact memory payload that would be injected into the prompt.

### 4. Answer
**Input:** `SearchArtifact` list
**Output:** `runs/<run_id>/answer/*.json` — one `AnswerArtifact` per question

Passes `memory_payload + question` to the answer LLM and records the completion, token counts, and latency.

### 5. Evaluate
**Input:** `AnswerArtifact` list + reference answers
**Output:** `runs/<run_id>/evaluate/*.json` — one `EvaluationArtifact` per item

Calls the judge LLM to score each completion against its reference answer. Computes `MemScore` triples (quality, latency, cost) across all items.

## Adapter Interface

Every provider adapter implements the `MemoryAdapter` protocol defined in `voice_memory_bench/core/adapter.py`. Key methods:

| Method | Description |
|--------|-------------|
| `capabilities()` | Returns provider metadata and supported retrieval modes |
| `health_check()` | Verifies the backing service is reachable |
| `add_message()` | Writes a conversation turn |
| `add_fact()` | Writes a pre-extracted fact |
| `search()` | Retrieves relevant memories |
| `enumerate_memories()` | Lists all memories for a user |
| `reset()` | Deletes all memories for a user |

## Artifact Schema

All artifacts are Pydantic models serialised as JSON. See `voice_memory_bench/core/schemas.py` for the full definitions.

## MemScore

MemScore is a triple, never a scalar:

- **Quality**: Normalised answer quality score [0, 1], computed by the judge LLM
- **Latency**: p50/p95/p99 retrieval latency in milliseconds
- **Cost**: Average cost per benchmark item in USD (0 for fully self-hosted providers)

The three axes are always reported side-by-side. Do not collapse to a single number.
