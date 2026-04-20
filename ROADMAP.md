# Roadmap

## Phase 1: Text-Chat Benchmark (In Progress)

- [ ] Core adapter protocol (`voice_memory_bench/core/adapter.py`) ✅ Scaffolded
- [ ] Pydantic config schema (`voice_memory_bench/core/config.py`) ✅ Scaffolded
- [ ] Pipeline stage definitions ✅ Scaffolded
- [ ] LoCoMo dataset loader
- [ ] LongMemEval dataset loader
- [ ] Mem0 adapter implementation
- [ ] Memori adapter implementation
- [ ] Graphiti adapter implementation
- [ ] Cognee adapter implementation
- [ ] Ingest stage
- [ ] Index stage
- [ ] Search stage
- [ ] Answer stage
- [ ] Evaluate stage (judge LLM)
- [ ] JSON reporter
- [ ] Markdown reporter
- [ ] CLI (`vmb run`, `vmb compare`)
- [ ] First benchmark results published

## Phase 2: Voice Benchmark (Planned)

- [ ] LiveKit Agents worker
- [ ] Deepgram STT adapter
- [ ] Cartesia TTS adapter
- [ ] Transcript replay stage
- [ ] Latency profiling harness
- [ ] Concurrent session isolation audit
- [ ] Voice-specific MemScore metrics
- [ ] First voice benchmark results published

## Phase 3: Continuous Benchmarking (Future)

- [ ] GitHub Actions scheduled benchmark runs
- [ ] Public leaderboard
- [ ] Provider version tracking
- [ ] Regression detection
