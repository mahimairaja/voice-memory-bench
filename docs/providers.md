# Provider Guide

## Mem0 OSS

**Backing store:** PostgreSQL 16 + pgvector

Mem0 is the reference baseline. It extracts facts from conversation turns using an LLM, stores them as vector embeddings in pgvector, and retrieves by semantic similarity.

### Setup

```bash
docker compose -f providers/mem0/docker-compose.yml up -d
export POSTGRES_PASSWORD=changeme
export OPENAI_API_KEY=sk-...
```

### Supported retrieval modes
- `SEMANTIC` ✅
- `HYBRID` ✅
- `KEYWORD` ❌
- `TEMPORAL` ❌
- `GRAPH` ❌

---

## Memori (GibsonAI)

**Backing store:** PostgreSQL 16 (no pgvector)

Memori stores memories in a plain SQL schema with tsvector-based full-text search. No embedding step on the write path makes it the fastest for p99 latency.

### Setup

```bash
docker compose -f providers/memori/docker-compose.yml up -d
export POSTGRES_PASSWORD=changeme
```

### Supported retrieval modes
- `SEMANTIC` ✅ (tsvector/tsquery)
- `KEYWORD` ✅
- `HYBRID` ✅ (RRF)
- `TEMPORAL` ❌
- `GRAPH` ❌

---

## Graphiti

**Backing store:** FalkorDB (knowledge graph)

Graphiti models memories as a temporal property graph. It excels at multi-hop queries and time-aware retrieval ("what did the user say last week about X?").

### Setup

```bash
docker compose -f providers/graphiti/docker-compose.yml up -d
export OPENAI_API_KEY=sk-...
```

### Supported retrieval modes
- `GRAPH` ✅
- `TEMPORAL` ✅
- `SEMANTIC` ✅
- `KEYWORD` ❌
- `HYBRID` ❌

---

## Cognee

**Backing store:** SQLite + LanceDB + Kuzu (all local)

Cognee is the local-first option. No network dependencies for storage — everything runs on the local filesystem. Ideal for air-gapped or data-residency-constrained deployments.

### Setup

```bash
# No docker needed — Cognee runs entirely locally
export OPENAI_API_KEY=sk-...  # Only needed for extraction LLM
```

### Supported retrieval modes
- `SEMANTIC` ✅ (LanceDB)
- `GRAPH` ✅ (Kuzu)
- `KEYWORD` ❌
- `TEMPORAL` ❌
- `HYBRID` ❌
