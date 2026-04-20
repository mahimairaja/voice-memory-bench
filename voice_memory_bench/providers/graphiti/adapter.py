"""
Graphiti (getzep) provider adapter with FalkorDB backend.

Graphiti is a temporal knowledge graph framework from the getzep team. It models
memories as a property graph with time-encoded edges, enabling temporally-aware
queries such as "what did the user say last week about X?" that are difficult or
impossible to answer with flat vector stores.

Why FalkorDB (not Neo4j)?
--------------------------
FalkorDB is a Redis module that speaks the Cypher query language over a Redis
wire protocol. Compared to Neo4j it offers:
* **Single Docker container**: no separate JVM process, no heap tuning.
* **Redis-like ergonomics**: RESP3 protocol, works with standard Redis clients.
* **Lower query latency**: sparse adjacency matrix representation gives sub-5 ms
  hop latency on graphs with < 1M nodes, which is typical for per-user memory.
* **Apache-2.0 license**: no AGPL concerns for self-hosting.

Self-hosting requirements
-------------------------
* FalkorDB running (see providers/graphiti/docker-compose.yml).
* The ``graphiti-core`` and ``falkordb`` Python packages
  (``pip install voice-memory-bench[graphiti]``).

Retrieval modes
---------------
* GRAPH    — multi-hop graph traversal.
* TEMPORAL — time-filtered graph queries.
* SEMANTIC — embedding-based edge/node search.

NOT supported: KEYWORD.
"""

from __future__ import annotations

from typing import Any

import structlog

from voice_memory_bench.core.adapter import (
    AdapterConfigError,
    CapabilityDescriptor,
    CapabilityNotSupportedError,
    MemoryItem,
    RetrievalMode,
    RetrievalResult,
    WriteResult,
)

logger = structlog.get_logger(__name__)

_SUPPORTED_MODES = frozenset({RetrievalMode.GRAPH, RetrievalMode.TEMPORAL, RetrievalMode.SEMANTIC})


class GraphitiAdapter:
    """
    Adapter wrapping Graphiti with a FalkorDB backend.

    Parameters
    ----------
    config:
        Provider config block. Expected keys:

        * ``falkordb_url`` (str) — Redis-compatible URL for FalkorDB,
          e.g. ``"redis://localhost:6379"``.
        * ``llm_provider`` (str) — LLM provider for Graphiti's internal
          entity/relation extraction.
        * ``llm_model`` (str) — Model for extraction.
        * ``graph_name`` (str, optional) — FalkorDB graph name.
          Defaults to ``"voice_memory_bench"``.
    """

    def __init__(self, config: dict[str, Any]) -> None:
        self._config = config
        self._client: Any = None
        self._validate_config()

    def _validate_config(self) -> None:
        """Raise AdapterConfigError for missing required keys."""
        required = ["falkordb_url", "llm_provider", "llm_model"]
        missing = [k for k in required if k not in self._config]
        if missing:
            raise AdapterConfigError(f"Graphiti adapter missing required config keys: {missing}")

    async def _init_client(self) -> None:
        """Lazily initialise the Graphiti client with a FalkorDB driver."""
        raise NotImplementedError(
            "TODO(mahimai): import graphiti_core and falkordb; construct Graphiti client"
        )

    async def capabilities(self) -> CapabilityDescriptor:
        """Return Graphiti's capability descriptor."""
        raise NotImplementedError(
            "TODO(mahimai): return CapabilityDescriptor with GRAPH, TEMPORAL, SEMANTIC support"
        )

    async def health_check(self) -> None:
        """Ping the FalkorDB instance."""
        raise NotImplementedError(
            "TODO(mahimai): send a PING command to FalkorDB and raise AdapterHealthError on failure"
        )

    async def add_message(
        self,
        user_id: str,
        session_id: str,
        role: str,
        content: str,
        metadata: dict[str, Any] | None = None,
    ) -> WriteResult:
        """Ingest a conversation message as a graph episode."""
        raise NotImplementedError("TODO(mahimai): call graphiti.add_episode() and measure latency")

    async def add_fact(
        self,
        user_id: str,
        session_id: str,
        fact: str,
        metadata: dict[str, Any] | None = None,
    ) -> WriteResult:
        """Ingest a pre-extracted fact as a graph node."""
        raise NotImplementedError("TODO(mahimai): add fact as a graph node/episode")

    async def search(
        self,
        user_id: str,
        session_id: str,
        query: str,
        mode: RetrievalMode = RetrievalMode.SEMANTIC,
        top_k: int = 10,
        filters: dict[str, Any] | None = None,
    ) -> RetrievalResult:
        """Retrieve memories via graph traversal or temporal query."""
        if mode not in _SUPPORTED_MODES:
            raise CapabilityNotSupportedError(
                provider="graphiti",
                capability=mode,
                reason=(
                    f"Graphiti supports {sorted(m.value for m in _SUPPORTED_MODES)}; got {mode!r}."
                ),
            )
        raise NotImplementedError(
            "TODO(mahimai): call graphiti.search() and wrap in RetrievalResult"
        )

    async def enumerate_memories(
        self,
        user_id: str,
        session_id: str | None = None,
    ) -> list[MemoryItem]:
        """List all graph nodes/edges for the user."""
        raise NotImplementedError(
            "TODO(mahimai): query FalkorDB for all nodes belonging to user_id"
        )

    async def reset(self, user_id: str, session_id: str | None = None) -> None:
        """Delete all graph data for the user."""
        raise NotImplementedError("TODO(mahimai): delete all nodes/edges for user_id from FalkorDB")
