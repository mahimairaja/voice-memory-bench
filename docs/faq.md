# FAQ

**Q: Why Python and not TypeScript?**  
Every self-hostable memory framework targeted here (Mem0, Memori, Graphiti, Cognee) ships a Python SDK as its canonical integration surface. The voice agent runtime (LiveKit Agents) is Python-first. TypeScript would require shelling out or wrapping RPC.

**Q: Why not include OpenAI Memory, ChatGPT Memory, etc.?**  
Scope is self-hostable only. Any provider that cannot be run on your own hardware under an open-source license is out of scope by design.

**Q: What is MemScore?**  
MemScore is a triple: (quality, latency, cost). It is never collapsed to a single scalar. The philosophy is borrowed from the [MemoryBench](https://github.com/supermemoryai/memorybench) project.

**Q: Is this production-ready?**  
No. Phase 1 (text-chat benchmark) is scaffolded but not implemented. Phase 2 (voice benchmark) is sketched. See ROADMAP.md.

**Q: Why FalkorDB instead of Neo4j for Graphiti?**  
FalkorDB runs as a single Docker container (no JVM), speaks the Redis wire protocol, and is Apache-2.0 licensed. For per-user memory graphs (< 1M nodes), it has lower query latency than Neo4j.

**Q: Can I add my own dataset?**  
Yes. Provide a JSONL file where each line matches the `BenchmarkItem` schema and set `dataset.name: custom` with `dataset.path` in your run config.

**Q: How do I report a security vulnerability?**  
See SECURITY.md.

**Q: How is this different from MemoryBench by Supermemory?**  
MemoryBench is TypeScript, targets cloud providers, and measures answer quality. voice-memory-bench is Python, targets self-hostable providers, and measures quality + latency + cost with explicit voice-agent concerns.

**Q: What LLMs are supported for answer/judge stages?**  
Any model accessible via an OpenAI-compatible API. Set `base_url` in the LLM config to point at a self-hosted model (e.g. Ollama, vLLM).

**Q: Can I run only specific pipeline stages?**  
Yes: `uv run vmb run config.yaml --stages index,search,answer`
