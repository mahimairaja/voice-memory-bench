# voice/ — Phase 2: Voice Benchmark

> **Status: INTENTIONALLY STUBBED.** Phase 2 is not implemented. The directory
> structure, module stubs, and docstrings are here so a second engineer can
> implement from the signatures alone. See [ROADMAP.md](ROADMAP.md) for the
> implementation plan.

---

## What phase 2 will measure

Phase 1 benchmarks answer quality, retrieval latency, and cost on text chat
transcripts. Phase 2 adds the dimensions that actually matter for real-time
voice agents:

| Metric | Why it matters |
|--------|---------------|
| **Wall-clock latency** from end-of-utterance to first memory-injected LLM token | Determines whether the agent "feels" natural or laggy |
| **p50 / p95 / p99 retrieval latency under concurrent sessions** | Tail latency drives perceived quality at scale |
| **Token footprint of injected memory** | Inflates time-to-first-token and per-call cost |
| **Write-path latency** (mid-call persistence) | Common in voice agents; tail behavior matters |
| **Cold vs warm recall** | First turn of a returning caller vs tenth turn of active session |
| **Concurrent session isolation** | Does user A's memory bleed into user B under load? |
| **Turn-level relevance over voice-shaped discourse** | Short turns, disfluencies, interruptions, topic drift |

## Architecture

```
voice/
├── pipeline/
│   ├── replay.py        # Replays transcripts as synthesised voice turns
│   ├── latency.py       # Wall-clock latency profiler
│   └── isolation.py     # Cross-user memory isolation audit
├── workers/
│   └── livekit_worker.py  # LiveKit Agents worker stub
├── stt/
│   └── deepgram.py      # Deepgram STT adapter stub
├── tts/
│   └── cartesia.py      # Cartesia TTS adapter stub
├── docker-compose.yml   # Self-hosted LiveKit + SIP stack sketch
├── README.md            # This file
└── ROADMAP.md           # Phase 2 implementation plan
```

The voice layer shares **only the adapter layer** with phase 1 — no pipeline
stage code, no dataset loaders, no reporters are reused. The adapter interface
(`voice_memory_bench/core/adapter.py`) is the only shared surface.

## How the harness will work

1. **Transcript replay**: LoCoMo / LongMemEval transcripts are split into
   individual turns. Each turn is synthesised to audio via Cartesia TTS.

2. **LiveKit session**: A simulated caller joins a LiveKit room. A LiveKit
   Agents worker serves the room, using the configured memory adapter as a tool.

3. **Latency measurement**: The harness records:
   - `t0`: VAD endpoint (end of user utterance detected by Deepgram)
   - `t1`: First LLM token generated after memory retrieval and injection
   - `latency = t1 - t0`

4. **Concurrency sweep**: The harness replays sessions at 1, 10, and 50
   concurrent sessions and records per-session latency histograms.

5. **Isolation audit**: After each concurrent run, the harness checks whether
   any retrieval for user A returned content written by user B.

## Infrastructure

The self-hosted stack for phase 2 requires:

- **LiveKit Server** (open-source, Docker)
- **LiveKit SIP Bridge** (for realistic telephony simulation)
- **Deepgram self-hosted** (on-prem STT, requires enterprise agreement — or
  swap for `whisper.cpp` for a fully open-source stack)
- **Cartesia self-hosted** (on-prem TTS — or swap for `piper-tts`)
- The memory provider's stack (from phase 1)

See `docker-compose.yml` for a sketch of this stack.

## Contributing

Phase 2 is the highest-impact unimplemented feature in this repo. If you want
to contribute, start with `voice/workers/livekit_worker.py` — the docstrings
describe exactly what the implementation should do.

Before opening a PR, post in [Discussions](https://github.com/mahimairaja/voice-memory-bench/discussions)
to coordinate with Mahimai on which concurrency levels and providers to target first.
