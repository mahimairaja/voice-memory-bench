# Mem0 Provider

Mem0 OSS is the **reference baseline** for voice-memory-bench. It provides hybrid
semantic memory using **pgvector** as the vector store and **PostgreSQL** as the
metadata store. It is LLM-agnostic: you configure which LLM to use for memory
extraction separately from the LLM used to answer benchmark questions.

## Why Mem0 is the baseline

Mem0 is the most-deployed open-source agent memory layer as of 2024. It has a
well-documented Python SDK, an active community, and an explicit self-hosting path.
Every other adapter in voice-memory-bench is compared against Mem0's MemScore triple.

## Infrastructure requirements

| Component | Minimum version | Notes |
|-----------|----------------|-------|
| PostgreSQL | 15+ | With the `pgvector` extension enabled |
| pgvector extension | 0.5+ | Ships as `pgvector/pgvector:pg16` Docker image |
| Python SDK | `mem0ai>=0.1` | Install with `uv sync --extra mem0` |

## Quick start

### 1. Start the infrastructure

```bash
docker compose -f providers/mem0/docker-compose.yml up -d
```

This starts PostgreSQL 16 with pgvector pre-installed on port 5432.

### 2. Set your secrets

```bash
export POSTGRES_PASSWORD=changeme
export OPENAI_API_KEY=sk-...   # or your preferred LLM API key
```

### 3. Run the benchmark

```bash
uv run vmb run providers/mem0/config.example.yaml
```

## Configuration reference

```yaml
provider:
  name: mem0
  config:
    # Required: PostgreSQL DSN for the Postgres + pgvector instance.
    postgres_url: postgresql://mem0:${POSTGRES_PASSWORD}@localhost:5432/mem0

    # Required: LLM used by Mem0 for internal memory extraction.
    # This is NOT the LLM used to answer benchmark questions.
    llm_provider: openai
    llm_model: gpt-4o-mini

    # Optional: embedding model (defaults to Mem0's built-in default).
    # embedding_model: text-embedding-3-small

    # Optional: pgvector collection name (default: voice_memory_bench).
    # collection_name: voice_memory_bench
```

## Supported retrieval modes

| Mode | Supported |
|------|-----------|
| SEMANTIC | ✅ |
| HYBRID | ✅ |
| KEYWORD | ❌ |
| TEMPORAL | ❌ |
| GRAPH | ❌ |

Requests for unsupported modes are recorded as SKIPPED in the benchmark report
with an explanatory message.

## Known limitations

- The OSS version does not support temporal or graph queries. The commercial
  Mem0 Cloud adds graph memory, but cloud-managed deployments are out of scope
  for this benchmark.
- Write throughput is bounded by Postgres transaction throughput. Under high
  concurrency (phase 2), expect p95/p99 write latency to diverge significantly
  from p50.

## Useful links

- [Mem0 GitHub](https://github.com/mem0ai/mem0)
- [Mem0 documentation](https://docs.mem0.ai)
- [pgvector GitHub](https://github.com/pgvector/pgvector)
