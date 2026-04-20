# Voice-fitness notes

vbench is voice-shaped from day one even though the MVP only drives the
memory provider from text-derived queries (LoCoMo QA). This note captures
why we measure what we measure, and what Phase 3 adds.

## Why voice changes the metric

A voice agent runs a tight loop:

```
user speech → STT → LLM (turn plan) → memory retrieve → LLM (final) → TTS → user ear
```

Memory retrieval sits on the critical path between ASR completion and TTS
start. The end-to-end budget for "feels natural" is roughly 300–500 ms. A
benchmark that only measures answer quality cannot predict whether the
provider is shippable; a benchmark that only measures p50 latency hides what
happens under 4x load — which is normal for any telephony workload.

## MVP measurements (v0.1)

- **Search p50/p95/p99** in the Go engine (authoritative clock) at 1x and 4x
  concurrency. Voice verdict is driven by p95 only.
- **Injected memory token footprint** (p50) — proxy for how much of the
  LLM's context budget the memory step is consuming.
- **Index-stage write latency** per turn — the mid-call persistence cost.
- **Judge LLM quality** — reported beside latency, never collapsed into it.

## Concurrency levels

MVP runs 1x and 4x. 4x is the first level where tail behaviour reliably
diverges from the single-caller case for small self-hosted stacks.

16x is deferred to v0.2 and requires explicit queue / back-pressure
modelling — at 16 concurrent callers against a default Postgres, the
provider is no longer what's being measured.

## Phase 3 additions

- **LiveKit Agents worker** that drives the adapter live during a call so the
  measured p95 includes STT pass-through and TTS start.
- **Concurrent-session isolation audit** — assert that memories from caller
  A never surface in caller B's retrieval result.
- **Cold vs warm caller quality delta** — a returning caller should not lose
  quality relative to their most recent call.
- **Voice-shaped discourse test set** — short turns, disfluencies, barge-ins
  transcribed from real audio.

None of Phase 3 is implemented yet.
