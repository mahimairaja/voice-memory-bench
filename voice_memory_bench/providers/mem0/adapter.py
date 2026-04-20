"""
Mem0 OSS provider adapter.

Mem0 is the reference baseline for voice-memory-bench. It uses pgvector as
the vector store and Postgres as the metadata store. It is LLM-agnostic and
supports semantic retrieval out of the box.

Self-hosting requirements
-------------------------
* PostgreSQL >= 15 with the ``pgvector`` extension enabled.
* The ``mem0ai`` Python package (``pip install voice-memory-bench[mem0]``).

See providers/mem0/README.md and providers/mem0/docker-compose.yml for infra setup.
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


class Mem0Adapter:
    """
    Adapter wrapping Mem0 OSS.

    Parameters
    ----------
    config:
        Provider config block from the run config. Expected keys:

        * ``postgres_url`` (str) — DSN for the Postgres + pgvector instance.
        * ``llm_provider`` (str) — LLM provider for Mem0's internal extraction
          (e.g. ``"openai"``). Not used for answer generation.
        * ``llm_model`` (str) — Model to use for Mem0's extraction.
        * ``embedding_model`` (str, optional) — Embedding model name.
          Defaults to Mem0's built-in default.
        * ``collection_name`` (str, optional) — pgvector collection name.
          Defaults to ``"voice_memory_bench"``.
    """

    def __init__(self, config: dict[str, Any]) -> None:
        self._config = config
        self._client: Any = None  # mem0.Memory instance, set in _init_client
        self._validate_config()

    def _validate_config(self) -> None:
        """Raise AdapterConfigError for missing required keys."""
        required = ["postgres_url"]
        missing = [k for k in required if k not in self._config]
        if missing:
            raise AdapterConfigError(f"Mem0 adapter missing required config keys: {missing}")

    async def _init_client(self) -> None:
        """Lazily initialise the mem0 client on first use."""
        raise NotImplementedError(
            "TODO(mahimai): import mem0ai and construct the Memory client with self._config"
        )

    async def capabilities(self) -> CapabilityDescriptor:
        """Return Mem0's capability descriptor."""
        raise NotImplementedError(
            "TODO(mahimai): return CapabilityDescriptor with SEMANTIC retrieval support"
        )

    async def health_check(self) -> None:
        """Ping the Postgres+pgvector instance."""
        raise NotImplementedError(
            "TODO(mahimai): attempt SELECT 1 against postgres_url;"
            " raise AdapterHealthError on failure"
        )

    async def add_message(
        self,
        user_id: str,
        session_id: str,
        role: str,
        content: str,
        metadata: dict[str, Any] | None = None,
    ) -> WriteResult:
        """Write a conversation message via mem0.Memory.add()."""
        raise NotImplementedError("TODO(mahimai): call mem0 client add() and measure latency")

    async def add_fact(
        self,
        user_id: str,
        session_id: str,
        fact: str,
        metadata: dict[str, Any] | None = None,
    ) -> WriteResult:
        """Write a pre-extracted fact via mem0.Memory.add()."""
        raise NotImplementedError("TODO(mahimai): call mem0 client add() with fact content")

    async def search(
        self,
        user_id: str,
        session_id: str,
        query: str,
        mode: RetrievalMode = RetrievalMode.SEMANTIC,
        top_k: int = 10,
        filters: dict[str, Any] | None = None,
    ) -> RetrievalResult:
        """Retrieve memories via mem0.Memory.search()."""
        if mode not in (RetrievalMode.SEMANTIC, RetrievalMode.HYBRID):
            raise CapabilityNotSupportedError(
                provider="mem0",
                capability=mode,
                reason=f"Mem0 OSS supports SEMANTIC and HYBRID retrieval only; got {mode!r}.",
            )
        raise NotImplementedError(
            "TODO(mahimai): call mem0 client search() and wrap in RetrievalResult"
        )

    async def enumerate_memories(
        self,
        user_id: str,
        session_id: str | None = None,
    ) -> list[MemoryItem]:
        """List all memories for the user via mem0.Memory.get_all()."""
        raise NotImplementedError("TODO(mahimai): call mem0 client get_all()")

    async def reset(self, user_id: str, session_id: str | None = None) -> None:
        """Delete all memories for the user."""
        raise NotImplementedError("TODO(mahimai): call mem0 client delete_all()")
