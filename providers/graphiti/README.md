# Graphiti Provider (getzep + FalkorDB)

Graphiti is a **temporal knowledge graph** memory layer by [getzep](https://github.com/getzep),
licensed under Apache 2.0. It models memory as a property graph where entities,
relationships, and facts are time-stamped, allowing queries like "what did the
user say about X in the last 7 days?"

## Why FalkorDB (not Neo4j)

Graphiti supports multiple graph backends. This adapter uses **FalkorDB** rather
than Neo4j for three reasons:

1. **Single container**: FalkorDB ships as a single Docker image (`falkordb/falkordb`)
   with Redis-like ergonomics — no cluster, no separate coordinator process.
2. **Latency**: FalkorDB's in-memory graph query engine has materially lower
   query latency than Neo4j for read-heavy workloads typical of memory retrieval.
3. **License**: FalkorDB is MIT-licensed. Neo4j Community Edition has more
   restrictive terms.

If you want to use Neo4j, you can swap the backend in the config; the adapter
interface is the same, but the config block will differ. See the Graphiti docs.

## Infrastructure requirements

| Component | Minimum version | Notes |
|-----------|----------------|-------|
| FalkorDB | 4.x+ | Ships as `falkordb/falkordb:latest` Docker image |
| Python SDK | `graphiti-core>=0.3` | Install with `uv sync --extra graphiti` |
| FalkorDB client | `falkordb>=1.0` | Included in the graphiti extra |

## Quick start

```bash
# Start FalkorDB
docker compose -f providers/graphiti/docker-compose.yml up -d

# Set secrets
export OPENAI_API_KEY=sk-...   # Graphiti uses an LLM for entity extraction

# Run the benchmark
uv run vmb run providers/graphiti/config.example.yaml
```

## Configuration reference

```yaml
provider:
  name: graphiti
  config:
    # Required: FalkorDB connection string.
    falkordb_url: redis://localhost:6380

    # Required: LLM for Graphiti's entity/relationship extraction.
    llm_provider: openai
    llm_model: gpt-4o-mini

    # Optional: embedding model.
    # embedding_model: text-embedding-3-small

    # Optional: graph name prefix for namespace isolation.
    # graph_prefix: vmb
```

## Supported retrieval modes

| Mode | Supported |
|------|-----------|
| SEMANTIC | ✅ |
| TEMPORAL | ✅ |
| GRAPH | ✅ |
| KEYWORD | ❌ |
| HYBRID | ❌ |

Graphiti's temporal and graph retrieval are its differentiating capabilities.
Use the `graphiti-locomo-temporal.yaml` example config to benchmark them.

## Useful links

- [Graphiti GitHub](https://github.com/getzep/graphiti)
- [FalkorDB GitHub](https://github.com/FalkorDB/FalkorDB)
- [FalkorDB documentation](https://docs.falkordb.com)
