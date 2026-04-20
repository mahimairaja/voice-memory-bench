"""
Memori (GibsonAI) provider adapter.

Memori is the low-latency, no-vector-DB contender in this benchmark. It stores
memories in a plain Postgres schema with SQL-based retrieval, making it the
strongest candidate for voice agents with strict sub-100 ms latency budgets.

Self-hosting requirements
-------------------------
* PostgreSQL >= 14 (no pgvector required).
* The ``gibson-memory`` Python package (``pip install voice-memory-bench[memori]``).

Retrieval modes
---------------
* SEMANTIC — full-text search via Postgres tsvector/tsquery.
* KEYWORD  — exact keyword matching.
* HYBRID   — combined semantic + keyword with RRF ranking.

NOT supported: TEMPORAL, GRAPH.
"""

from __future__ import annotations

from typing import Any
import structlog

from voice_memory_bench.core.adapter import (
    AdapterConfigError,
    AdapterHealthError,
    BackingStore,
    CapabilityDescriptor,
    CapabilityNotSupportedError,
    MemoryItem,
    RetrievalMode,
    RetrievalResult,
    WriteResult,
)

logger = structlog.get_logger(__name__)

_SUPPORTED_MODES = frozenset({RetrievalMode.SEMANTIC, RetrievalMode.KEYWORD, RetrievalMode.HYBRID})


class MemoriAdapter:
    """
    Adapter wrapping Memori (GibsonAI).

    Parameters
    ----------
    config:
        Provider config block. Expected keys:

        * ``postgres_url`` (str) — DSN for the Postgres instance.
        * ``collection_name`` (str, optional) — logical namespace for memories.
          Defaults to ``"voice_memory_bench"``.
    """

    def __init__(self, config: dict[str, Any]) -> None:
        self._config = config
        self._client: Any = None
        self._validate_config()

    def _validate_config(self) -> None:
        """Raise AdapterConfigError for missing required keys."""
        required = ["postgres_url"]
        missing = [k for k in required if k not in self._config]
        if missing:
            raise AdapterConfigError(f"Memori adapter missing required config keys: {missing}")

    async def _init_client(self) -> None:
        """Lazily initialise the gibson-memory client on first use."""
        raise NotImplementedError(
            "TODO(mahimai): import gibson_memory and construct the client with self._config"
        )

    async def capabilities(self) -> CapabilityDescriptor:
        """Return Memori's capability descriptor."""
        raise NotImplementedError(
            "TODO(mahimai): return CapabilityDescriptor with SEMANTIC, KEYWORD, HYBRID support"
        )

    async def health_check(self) -> None:
        """Ping the Postgres instance."""
        raise NotImplementedError(
            "TODO(mahimai): attempt a simple SELECT 1 against the postgres_url and raise AdapterHealthError on failure"
        )

    async def add_message(
        self,
        user_id: str,
        session_id: str,
        role: str,
        content: str,
        metadata: dict[str, Any] | None = None,
    ) -> WriteResult:
        """Write a conversation message to Memori."""
        raise NotImplementedError("TODO(mahimai): call Memori client add() and measure latency")

    async def add_fact(
        self,
        user_id: str,
        session_id: str,
        fact: str,
        metadata: dict[str, Any] | None = None,
    ) -> WriteResult:
        """Write a pre-extracted fact to Memori."""
        raise NotImplementedError("TODO(mahimai): call Memori client add() with fact content")

    async def search(
        self,
        user_id: str,
        session_id: str,
        query: str,
        mode: RetrievalMode = RetrievalMode.SEMANTIC,
        top_k: int = 10,
        filters: dict[str, Any] | None = None,
    ) -> RetrievalResult:
        """Retrieve memories using SQL-based search."""
        if mode not in _SUPPORTED_MODES:
            raise CapabilityNotSupportedError(
                provider="memori",
                capability=mode,
                reason=f"Memori supports {sorted(m.value for m in _SUPPORTED_MODES)}; got {mode!r}.",
            )
        raise NotImplementedError("TODO(mahimai): call Memori search() and wrap in RetrievalResult")

    async def enumerate_memories(
        self,
        user_id: str,
        session_id: str | None = None,
    ) -> list[MemoryItem]:
        """List all memories for the user."""
        raise NotImplementedError("TODO(mahimai): call Memori get_all()")

    async def reset(self, user_id: str, session_id: str | None = None) -> None:
        """Delete all memories for the user."""
        raise NotImplementedError("TODO(mahimai): call Memori delete_all()")
