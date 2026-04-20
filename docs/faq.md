# FAQ

**Q: Why Go for the engine and Python only for the sidecar?**
The engine ships as a single static binary, which is agent-friendly to
distribute. Go's scheduler has sub-ms GC pauses, which keeps the reported
p99 honest; the GIL would leak interpreter jitter into the measurement.
Provider SDKs are Python-first, so the sidecar uses them directly instead of
being reimplemented. The HTTP boundary mirrors what cloud variants of these
providers expose, so a config flip can later point the engine at a cloud API
without engine changes.

**Q: Why is latency measured in Go rather than in the sidecar?**
So the number is authoritative. The sidecar does self-report a latency for
debugging, but the voice verdict uses the Go-side wall clock around each
HTTP call. That is what a voice agent would feel.

**Q: Why only Mem0 in the MVP?**
To push an end-to-end headline verdict out the door. Memori, Graphiti, and
Cognee all sit behind the same HTTP contract — adding them is a v0.2 task.

**Q: Why no cloud providers?**
Self-hosted only in v0.1. Cloud-managed endpoints land in v0.4 via the same
HTTP contract.

**Q: What counts as ACCEPTABLE vs FAIL?**
Voice verdict is driven by search-stage p95: `<300 ms` EXCELLENT,
`300–500 ms` ACCEPTABLE, `>500 ms` FAIL.

**Q: How is this different from MemoryBench (Supermemory)?**
MemoryBench is a TypeScript harness that inspired the MemScore framing and
the five-stage pipeline shape. vbench inherits both, then narrows the focus
to self-hostable providers and voice-agent fitness metrics, and rebuilds the
engine in Go so measurements are clean.

**Q: What LLMs are supported for answer/judge stages?**
Anything that speaks the OpenAI chat-completions API. Set `base_url` in the
LLM config to point at a self-hosted endpoint (vLLM, Ollama, etc.).

**Q: How do I resume a crashed run?**
`./vbench eval --config … --run-id <id>`. Each stage writes a `.complete`
sentinel and the next run skips anything already marked complete.

**Q: Security?**
See `SECURITY.md`.
