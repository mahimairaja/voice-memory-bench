# Phase 2 Roadmap — Voice Benchmark

> **This directory is intentionally stubbed.** Phase 2 will not be implemented
> until Phase 1 (text-chat benchmark) is validated end-to-end with at least the
> Mem0 reference baseline. See the main [ROADMAP.md](../ROADMAP.md) for the
> overall project roadmap.

---

## Why phase 2 is blocked on phase 1

Phase 2 requires a working adapter layer and a validated MemScore pipeline.
Implementing the LiveKit harness before the adapter interface is stable would
mean rewriting it every time the adapter interface changes. Phase 1 is the
stability foundation.

## Phase 2 implementation plan

When phase 1 is green (Mem0 adapter working end-to-end, MemScore triple
produced for LoCoMo), phase 2 will proceed in this order:

### Step 1: LiveKit worker skeleton

- [ ] Implement `voice/workers/livekit_worker.py` — a LiveKit Agents worker
  that accepts a `MemoryAdapter` as a constructor argument
- [ ] Wire up the worker so it calls `adapter.search()` on every user utterance
  and injects the result into the LLM context
- [ ] Verify the worker runs against a local LiveKit + Mem0 stack

### Step 2: Transcript replay

- [ ] Implement `voice/pipeline/replay.py` — splits a benchmark conversation
  into turns, synthesises each user turn to audio via Cartesia TTS, and
  submits it to the LiveKit room
- [ ] Add a `--replay-speed` flag (1x, 2x, real-time) to control pacing
- [ ] Validate that the worker correctly processes replayed turns

### Step 3: Latency profiler

- [ ] Implement `voice/pipeline/latency.py` — hooks into Deepgram's VAD
  endpoint event and the LLM first-token event to measure wall-clock latency
- [ ] Write per-turn latency to a JSONL artifact
- [ ] Compute p50/p95/p99 histograms from the artifact

### Step 4: Concurrency sweep

- [ ] Extend the replay harness to spawn N concurrent sessions (target: 1, 10, 50)
- [ ] Add session-level latency tracking (each session gets its own JSONL)
- [ ] Produce an aggregate latency histogram across all sessions

### Step 5: Isolation audit

- [ ] Implement `voice/pipeline/isolation.py` — after each concurrent run,
  checks whether any retrieval for user A returned content written by user B
- [ ] Produce an isolation report: 0 violations = PASS, any violation = FAIL

### Step 6: STT/TTS extension points

- [ ] Implement `voice/stt/deepgram.py` (primary) and a `whisper_cpp.py` stub
  for open-source alternative
- [ ] Implement `voice/tts/cartesia.py` (primary) and a `piper.py` stub for
  open-source alternative
- [ ] Document how to add a new STT or TTS provider

### Step 7: Reports

- [ ] Extend the `reporters/` module to produce voice-specific reports:
  latency histograms (HTML), isolation audit summary, concurrency comparison
- [ ] Add a `vmb voice run` CLI subcommand
- [ ] Add a `vmb voice compare` subcommand for comparing providers at different
  concurrency levels

---

## Tracking issue

Follow progress on the main GitHub issue tracker. Search for label `phase-2`.
