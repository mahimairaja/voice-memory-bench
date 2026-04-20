# Providers

## Mem0 OSS (MVP)

Backing store: Postgres 16 + pgvector.

Sidecar: [`sidecars/mem0`](../sidecars/mem0).
Backing infra: [`providers/mem0/docker-compose.yml`](../providers/mem0/docker-compose.yml).

### Supported retrieval modes

| Mode     | Supported |
|----------|-----------|
| semantic | yes       |
| hybrid   | yes       |
| keyword  | no        |
| temporal | no        |
| graph    | no        |

Unsupported modes return `422 capability_not_supported` from the sidecar, which
the engine surfaces as a SKIPPED entry in the report.

### Run it

```bash
export POSTGRES_PASSWORD=changeme
export OPENAI_API_KEY=sk-...
make providers-up-mem0
make sidecar-sync
make build
./vbench datasets download locomo
./vbench eval --config examples/configs/mem0-locomo.yaml --max-items 2
```

---

## Roadmap providers (v0.2+)

Each of these will land as a standalone sidecar package + docker-compose file,
wired through the same HTTP contract:

- **Memori** (SQL + tsvector, no embedding on the write path) — strong p99 candidate.
- **Graphiti** (FalkorDB-backed temporal graph) — good for multi-hop / time queries.
- **Cognee** (SQLite + LanceDB + Kuzu, fully local) — air-gapped / data-residency option.

None of these are in the MVP. Do not wire them up yet.
