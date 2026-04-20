"""
Cognee provider adapter.

Cognee is the local-first, data-sovereignty contender in this benchmark. It uses
a combination of SQLite (metadata), LanceDB (vectors), and Kuzu (knowledge graph),
all running on the local filesystem with no network dependencies for storage.

This makes Cognee the recommended choice for:
* EU/Swiss on-prem deployments with strict data residency requirements.
* Air-gapped environments where cloud vector stores are unavailable.
* Development and testing where spinning up Postgres/FalkorDB is inconvenient.

Self-hosting requirements
-------------------------
* No external services — all storage is local.
* The ``cognee`` Python package (``pip install voice-memory-bench[cognee]``).
* A local directory with write access (defaults to ``./cognee_data``).

Retrieval modes
---------------
* SEMANTIC — LanceDB vector similarity search.
* GRAPH    — Kuzu knowledge graph traversal.

NOT supported: KEYWORD, TEMPORAL, HYBRID.
"""

from __future__ import annotations

import pathlib
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

_SUPPORTED_MODES = frozenset({RetrievalMode.SEMANTIC, RetrievalMode.GRAPH})


class CogneeAdapter:
    """
    Adapter wrapping Cognee.

    Parameters
    ----------
    config:
        Provider config block. Expected keys:

        * ``data_dir`` (str) — Path to local storage directory.
          All SQLite, LanceDB, and Kuzu data is written here.
          Defaults to ``"./cognee_data"`` if not provided.
        * ``llm_provider`` (str, optional) — LLM provider for Cognee's internal
          knowledge extraction. Defaults to ``"openai"``.
        * ``llm_model`` (str, optional) — Model for extraction.
    """

    def __init__(self, config: dict[str, Any]) -> None:
        self._config = config
        self._data_dir = pathlib.Path(config.get("data_dir", "./cognee_data"))
        self._client: Any = None

    async def _init_client(self) -> None:
        """Lazily initialise the Cognee client, creating data_dir if needed."""
        raise NotImplementedError(
            "TODO(mahimai): import cognee and initialise with self._data_dir"
        )

    async def capabilities(self) -> CapabilityDescriptor:
        """Return Cognee's capability descriptor."""
        raise NotImplementedError(
            "TODO(mahimai): return CapabilityDescriptor with SEMANTIC and GRAPH support"
        )

    async def health_check(self) -> None:
        """Verify the data directory is accessible and writable."""
        raise NotImplementedError(
            "TODO(mahimai): check self._data_dir exists and is writable; raise AdapterHealthError if not"
        )

    async def add_message(
        self,
        user_id: str,
        session_id: str,
        role: str,
        content: str,
        metadata: dict[str, Any] | None = None,
    ) -> WriteResult:
        """Write a conversation message to Cognee."""
        raise NotImplementedError("TODO(mahimai): call cognee.add() and measure latency")

    async def add_fact(
        self,
        user_id: str,
        session_id: str,
        fact: str,
        metadata: dict[str, Any] | None = None,
    ) -> WriteResult:
        """Write a pre-extracted fact to Cognee."""
        raise NotImplementedError("TODO(mahimai): call cognee.add() with fact content")

    async def search(
        self,
        user_id: str,
        session_id: str,
        query: str,
        mode: RetrievalMode = RetrievalMode.SEMANTIC,
        top_k: int = 10,
        filters: dict[str, Any] | None = None,
    ) -> RetrievalResult:
        """Retrieve memories via LanceDB (SEMANTIC) or Kuzu (GRAPH)."""
        if mode not in _SUPPORTED_MODES:
            raise CapabilityNotSupportedError(
                provider="cognee",
                capability=mode,
                reason=f"Cognee supports {sorted(m.value for m in _SUPPORTED_MODES)}; got {mode!r}.",
            )
        raise NotImplementedError("TODO(mahimai): call cognee.search() and wrap in RetrievalResult")

    async def enumerate_memories(
        self,
        user_id: str,
        session_id: str | None = None,
    ) -> list[MemoryItem]:
        """List all memory items for the user."""
        raise NotImplementedError("TODO(mahimai): query Cognee for all items belonging to user_id")

    async def reset(self, user_id: str, session_id: str | None = None) -> None:
        """Delete all memory for the user."""
        raise NotImplementedError("TODO(mahimai): call cognee.prune() or delete user data directory")
