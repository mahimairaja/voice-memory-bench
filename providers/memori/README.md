# Memori Provider (GibsonAI)

Memori is a **SQL-native memory layer** built by GibsonAI. It uses PostgreSQL as
its only backing store — no vector database, no graph database. This makes it
the **lowest-latency contender** in voice-memory-bench: with no vector index to
query, retrieval can hit the Postgres query planner directly, which is highly
predictable under concurrent load.

> **Voice agent relevance**: Memori is the provider most structurally suited to
> voice agent latency budgets. If your SLA requires sub-50ms p95 retrieval and
> you're already running Postgres, Memori is the first adapter to evaluate.

## Infrastructure requirements

| Component | Minimum version | Notes |
|-----------|----------------|-------|
| PostgreSQL | 15+ | Standard installation, no extensions required |
| Memori SDK | TBD | See installation note below |

### Installation note

The Memori Python SDK is not yet published to PyPI under a stable package name.
Once released, it will be installable via `uv sync --extra memori`. Until then,
follow the [GibsonAI documentation](https://gibson.ai) for manual installation
and update the `pyproject.toml` `memori` extra accordingly.

## Quick start

```bash
# Start PostgreSQL
docker compose -f providers/memori/docker-compose.yml up -d

# Set secrets
export POSTGRES_PASSWORD=changeme

# Run the benchmark
uv run vmb run providers/memori/config.example.yaml
```

## Configuration reference

```yaml
provider:
  name: memori
  config:
    # Required: PostgreSQL DSN.
    postgres_url: postgresql://memori:${POSTGRES_PASSWORD}@localhost:5433/memori

    # Optional: schema name for namespace isolation (default: public).
    # schema: vmb
```

## Supported retrieval modes

| Mode | Supported |
|------|-----------|
| SEMANTIC | ✅ |
| KEYWORD | ✅ |
| HYBRID | ✅ |
| TEMPORAL | ❌ |
| GRAPH | ❌ |

## Useful links

- [GibsonAI](https://gibson.ai)
