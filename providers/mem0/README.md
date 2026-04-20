# Mem0 Provider

Mem0 OSS is the MVP reference provider for vbench. It backs `vbench-mem0`
(the Python sidecar) with PostgreSQL + pgvector.

## Infra

| Component | Version |
|-----------|---------|
| Postgres  | 16 (with pgvector) — `pgvector/pgvector:pg16` |
| Sidecar   | [`sidecars/mem0`](../../sidecars/mem0) (`uv sync`) |

## Start the backing store

```bash
export POSTGRES_PASSWORD=changeme
make providers-up-mem0
```

## Sync the sidecar env

```bash
make sidecar-sync
```

## Config block

```yaml
provider:
  name: mem0
  # optional: override the sidecar launch command
  # command: ["uv", "run", "vbench-mem0"]
  config:
    postgres_url: postgresql://mem0:${POSTGRES_PASSWORD}@localhost:5432/mem0
    llm_provider: openai
    llm_model: gpt-4o-mini
    # collection_name: vbench
    # embedding_model: text-embedding-3-small
```

## Supported retrieval modes

| Mode     | Supported |
|----------|-----------|
| SEMANTIC | yes       |
| HYBRID   | yes       |
| KEYWORD  | no        |
| TEMPORAL | no        |
| GRAPH    | no        |

Unsupported modes surface in the report as a SKIPPED entry with reason.

## Links

- Mem0 — https://github.com/mem0ai/mem0
- pgvector — https://github.com/pgvector/pgvector
