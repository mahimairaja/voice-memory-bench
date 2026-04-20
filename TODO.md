# TODO

This file tracks implementation TODOs. See ROADMAP.md for the high-level plan.

## Critical Path (Phase 1)

- [ ] Implement `LoCoMoLoader.download()` and `LoCoMoLoader.load()`
- [ ] Implement `LongMemEvalLoader.download()` and `LongMemEvalLoader.load()`
- [ ] Implement `Mem0Adapter` — requires mem0ai package and pgvector
- [ ] Implement `MemoriAdapter` — requires gibson-memory package
- [ ] Implement `GraphitiAdapter` — requires graphiti-core and falkordb packages
- [ ] Implement `CogneeAdapter` — requires cognee package
- [ ] Implement ingest pipeline stage
- [ ] Implement index pipeline stage
- [ ] Implement search pipeline stage
- [ ] Implement answer pipeline stage (OpenAI-compatible LLM call)
- [ ] Implement evaluate pipeline stage (judge LLM scoring)
- [ ] Implement `JsonReporter.write()`
- [ ] Implement `MarkdownReporter.write()`
- [ ] Implement `vmb run` CLI command (wire up pipeline stages)
- [ ] Implement `vmb compare` CLI command
- [ ] Implement `vmb providers list` and `vmb providers check`
- [ ] Implement `vmb datasets download` and `vmb datasets info`
- [ ] Fill in SHA-256 hashes for dataset downloads

## Phase 2

- [ ] Implement `TranscriptReplayStage`
- [ ] Implement `LatencyProfile`
- [ ] Implement `IsolationAuditStage`
- [ ] Implement `VoiceMemoryWorker` (LiveKit Agents)
- [ ] Implement `DeepgramSTTAdapter`
- [ ] Implement `CartesiaTTSAdapter`
