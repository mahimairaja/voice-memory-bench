# Cognee Provider

Cognee is a **local-first, air-gapped memory layer** that runs entirely on the
local filesystem using SQLite, LanceDB, and Kuzu. There are no network
dependencies for storage: all data lives under a configurable directory on disk.

> **EU/Swiss on-prem use case**: Cognee is the recommended provider for
> deployments where data must not leave the host machine — e.g. healthcare,
> legal, or government contexts subject to GDPR or Swiss DPA. With Cognee, the
> benchmark can run with zero outbound network traffic (assuming a local LLM).

## Infrastructure requirements

| Component | Minimum version | Notes |
|-----------|----------------|-------|
| SQLite | 3.35+ | Ships with Python 3.11 |
| LanceDB | bundled | Installed as part of `cognee` package |
| Kuzu | bundled | Installed as part of `cognee` package |
| Cognee SDK | `cognee>=0.1` | Install with `uv sync --extra cognee` |

**No Docker required.** All storage is local.

## Quick start

```bash
# Install
uv sync --extra cognee

# Configure (no infrastructure to start)
export OPENAI_API_KEY=sk-...   # or point to a local LLM

# Run the benchmark
uv run vmb run providers/cognee/config.example.yaml
```

## Configuration reference

```yaml
provider:
  name: cognee
  config:
    # Required: local directory for all Cognee data.
    # Relative paths are resolved from the current working directory.
    data_dir: ~/.local/share/vmb-cognee

    # Required: LLM for Cognee's knowledge graph construction.
    llm_provider: openai
    llm_model: gpt-4o-mini
    # To use a local model, set llm_provider: ollama and llm_model: llama3

    # Optional: embedding model.
    # embedding_model: text-embedding-3-small
```

## Supported retrieval modes

| Mode | Supported |
|------|-----------|
| SEMANTIC | ✅ |
| GRAPH | ✅ |
| KEYWORD | ❌ |
| TEMPORAL | ❌ |
| HYBRID | ❌ |

## Air-gapped operation

To run fully offline:
1. Use an Ollama local model for both the LLM and embeddings.
2. Pre-download the LoCoMo or LongMemEval dataset while you have connectivity.
3. Set `llm_provider: ollama` in the config.
4. Disconnect from the network and run `uv run vmb run ...` normally.

## Useful links

- [Cognee GitHub](https://github.com/topoteretes/cognee)
- [LanceDB GitHub](https://github.com/lancedb/lancedb)
- [Kuzu GitHub](https://github.com/kuzudb/kuzu)
