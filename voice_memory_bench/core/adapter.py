"""
Core adapter protocol and supporting types for voice-memory-bench.

Every memory provider implements the :class:`MemoryAdapter` protocol. The adapter
is the only place where provider-specific logic lives; nothing above this layer
may import from a concrete provider module.

Changing this interface is a breaking change and requires a major version bump.
See CONTRIBUTING.md for the versioning policy.
"""

from __future__ import annotations

import datetime
import enum
from typing import Any, Protocol, runtime_checkable

from pydantic import BaseModel, Field


class RetrievalMode(str, enum.Enum):
    """Retrieval modes that a provider may or may not support."""

    SEMANTIC = "semantic"
    KEYWORD = "keyword"
    TEMPORAL = "temporal"
    GRAPH = "graph"
    HYBRID = "hybrid"


class BackingStore(str, enum.Enum):
    """Backing store type reported by the adapter."""

    POSTGRES = "postgres"
    SQLITE = "sqlite"
    VECTOR_DB = "vector_db"
    GRAPH_DB = "graph_db"
    HYBRID = "hybrid"
    UNKNOWN = "unknown"


class CapabilityDescriptor(BaseModel):
    """
    Describes the capabilities and cost model of a provider adapter.

    Returned by :meth:`MemoryAdapter.capabilities` at startup.
    """

    provider_name: str = Field(..., description="Human-readable provider name.")
    provider_version: str = Field(..., description="Version of the provider SDK/library in use.")
    supported_retrieval_modes: list[RetrievalMode] = Field(
        ...,
        description="List of retrieval modes the adapter implements.",
    )
    backing_store: BackingStore = Field(..., description="Primary backing store type.")
    supports_user_scoping: bool = Field(
        ...,
        description="True if the provider natively scopes memory to users.",
    )
    supports_session_scoping: bool = Field(
        ...,
        description="True if the provider natively scopes memory to sessions.",
    )
    declared_cost_model: str | None = Field(
        None,
        description="Free-text description of the provider's cost model (tokens, API calls, etc.). "
        "None if self-hosted with no metered cost.",
    )
    extra: dict[str, Any] = Field(
        default_factory=dict,
        description="Provider-specific capability metadata not captured above.",
    )


class WriteResult(BaseModel):
    """Result of a write operation."""

    provider_id: str | None = Field(None, description="Provider-assigned ID for the written item.")
    latency_ms: float = Field(..., description="Wall-clock write latency in milliseconds.")
    tokens_written: int | None = Field(
        None, description="Token count of the written content, if known."
    )
    extra: dict[str, Any] = Field(default_factory=dict)


class MemoryItem(BaseModel):
    """A single memory item as returned by a retrieval operation."""

    item_id: str = Field(..., description="Provider-assigned or synthesized item identifier.")
    content: str = Field(..., description="The memory text as it would be injected into a prompt.")
    score: float | None = Field(
        None, description="Retrieval relevance score (higher = more relevant)."
    )
    created_at: datetime.datetime | None = Field(None)
    metadata: dict[str, Any] = Field(default_factory=dict)


class RetrievalResult(BaseModel):
    """Result of a retrieval operation."""

    items: list[MemoryItem] = Field(
        ..., description="Retrieved memory items, ordered by relevance."
    )
    latency_ms: float = Field(..., description="Wall-clock retrieval latency in milliseconds.")
    retrieval_mode: RetrievalMode = Field(..., description="The retrieval mode that was used.")
    token_footprint: int | None = Field(
        None,
        description="Total token count of all retrieved items concatenated, if pre-computed.",
    )
    extra: dict[str, Any] = Field(default_factory=dict)


class CapabilityNotSupportedError(Exception):
    """
    Raised by an adapter method when the requested capability is not supported.

    The benchmark runner catches this exception and records the test as SKIPPED
    with the reason from this exception rather than as a failure. This allows
    the final report to say "Provider X does not support temporal queries" rather
    than silently omitting the result.

    Usage::

        raise CapabilityNotSupportedError(
            provider="graphiti",
            capability=RetrievalMode.KEYWORD,
            reason="Graphiti does not implement keyword (BM25) retrieval.",
        )
    """

    def __init__(self, provider: str, capability: str | RetrievalMode, reason: str) -> None:
        self.provider = provider
        self.capability = capability
        self.reason = reason
        super().__init__(f"[{provider}] {capability!r} not supported: {reason}")


class AdapterConfigError(Exception):
    """Raised when an adapter cannot be initialised due to a configuration error."""


class AdapterHealthError(Exception):
    """Raised when :meth:`MemoryAdapter.health_check` detects the backing service is unavailable."""


@runtime_checkable
class MemoryAdapter(Protocol):
    """
    Protocol that every memory provider adapter must implement.

    Design notes
    ------------
    * All methods are ``async``. Providers are network-bound; sync would block the
      event loop and give misleading latency readings.
    * Every operation is scoped to a ``(user_id, session_id)`` pair. If the provider
      does not support native scoping, the adapter must synthesise it (e.g. by
      prefixing keys with ``{user_id}/{session_id}/``).
    * Adapters must not cache or buffer writes. Each call must reach the backing
      store before the coroutine returns, so that latency measurements are accurate.
    * Adapters must not raise provider-specific exceptions outside this module.
      Translate all provider errors into one of the typed exceptions defined here.

    Extension points
    ----------------
    Adding a new retrieval mode: add a value to :class:`RetrievalMode`, add the
    corresponding method to this protocol, and add ``raise CapabilityNotSupportedError``
    as the default body in the base class adapters that do not yet support it.
    """

    async def capabilities(self) -> CapabilityDescriptor:
        """
        Return the capability descriptor for this adapter.

        Called once at benchmark startup. The result is recorded in the run
        manifest and used to decide which test variants to skip.

        Returns
        -------
        CapabilityDescriptor
            Full description of what this adapter can do.
        """
        ...

    async def health_check(self) -> None:
        """
        Verify that the backing service is reachable and operational.

        Raises
        ------
        AdapterHealthError
            If the service is down, misconfigured, or returns an unexpected response.
        """
        ...

    async def add_message(
        self,
        user_id: str,
        session_id: str,
        role: str,
        content: str,
        metadata: dict[str, Any] | None = None,
    ) -> WriteResult:
        """
        Write a single conversation message to the provider's memory store.

        Parameters
        ----------
        user_id:
            Stable identifier for the end user (caller).
        session_id:
            Identifier for the current conversation session.
        role:
            Typically ``"user"`` or ``"assistant"``.
        content:
            The raw text of the message.
        metadata:
            Optional key-value pairs the adapter may pass through to the provider.

        Returns
        -------
        WriteResult
            Includes write latency and the provider-assigned memory ID, if any.
        """
        ...

    async def add_fact(
        self,
        user_id: str,
        session_id: str,
        fact: str,
        metadata: dict[str, Any] | None = None,
    ) -> WriteResult:
        """
        Write a free-text fact about the user to the memory store.

        Facts are distinct from messages: they are pre-extracted, atomic pieces
        of information (e.g. "user is allergic to peanuts") rather than raw
        conversational turns.

        Parameters
        ----------
        fact:
            The fact as a complete sentence.
        """
        ...

    async def search(
        self,
        user_id: str,
        session_id: str,
        query: str,
        mode: RetrievalMode = RetrievalMode.SEMANTIC,
        top_k: int = 10,
        filters: dict[str, Any] | None = None,
    ) -> RetrievalResult:
        """
        Retrieve memories relevant to *query*.

        Parameters
        ----------
        query:
            The retrieval query, typically the current user utterance or a
            benchmark question.
        mode:
            Which retrieval strategy to use. If the adapter does not support
            the requested mode, it must raise :exc:`CapabilityNotSupportedError`.
        top_k:
            Maximum number of items to return.
        filters:
            Provider-specific filter dict.

        Returns
        -------
        RetrievalResult
            Ordered list of memory items plus latency metrics.

        Raises
        ------
        CapabilityNotSupportedError
            If ``mode`` is not supported by this adapter.
        """
        ...

    async def enumerate_memories(
        self,
        user_id: str,
        session_id: str | None = None,
    ) -> list[MemoryItem]:
        """
        Return all memory items for *user_id*, optionally filtered to *session_id*.

        Used by the evaluator to compute recall@k without an extra LLM round-trip.
        The adapter must not truncate or paginate internally; it must return the
        complete state.
        """
        ...

    async def reset(self, user_id: str, session_id: str | None = None) -> None:
        """
        Delete all memory for *user_id* (and optionally *session_id*).

        Must be atomic with respect to concurrent reads. Benchmark items are
        isolated from each other by calling this between items.

        Parameters
        ----------
        session_id:
            If ``None``, wipe all sessions for the user. If provided, wipe only
            that session.
        """
        ...
