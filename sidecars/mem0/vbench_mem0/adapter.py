"""Mem0 adapter.

Wraps the mem0ai SDK behind the vbench sidecar contract. This is the only
place in the sidecar that imports mem0 directly; the HTTP layer stays
provider-agnostic so new adapters can be added without touching it.
"""

from __future__ import annotations

import importlib.metadata
import logging
import time
from typing import Any

logger = logging.getLogger(__name__)

try:
    from mem0 import Memory  # type: ignore[import-untyped]
except ImportError as exc:  # pragma: no cover
    raise RuntimeError(
        "mem0ai is not installed in this sidecar environment. Run `uv sync` in sidecars/mem0."
    ) from exc


SUPPORTED_MODES = {"semantic", "hybrid"}


class CapabilityNotSupported(Exception):
    def __init__(self, capability: str, reason: str) -> None:
        super().__init__(reason)
        self.capability = capability
        self.reason = reason


class Mem0Adapter:
    def __init__(self, config: dict[str, Any]) -> None:
        self._config = config or {}
        self._client = self._build_client()

    def _build_client(self) -> Memory:
        postgres_url = self._config.get("postgres_url")
        if not postgres_url:
            raise ValueError("Mem0 sidecar requires provider.config.postgres_url")

        # Best-effort parse. Mem0's vector_store.pgvector configuration wants
        # discrete fields (host, port, user, password, dbname, collection_name).
        from urllib.parse import urlparse

        parsed = urlparse(postgres_url)
        host = parsed.hostname or "localhost"
        port = parsed.port or 5432
        user = parsed.username or "postgres"
        password = parsed.password or ""
        dbname = (parsed.path or "/postgres").lstrip("/") or "postgres"
        collection = self._config.get("collection_name", "vbench")

        mem_config: dict[str, Any] = {
            "vector_store": {
                "provider": "pgvector",
                "config": {
                    "host": host,
                    "port": port,
                    "user": user,
                    "password": password,
                    "dbname": dbname,
                    "collection_name": collection,
                },
            }
        }
        if "llm_provider" in self._config and "llm_model" in self._config:
            mem_config["llm"] = {
                "provider": self._config["llm_provider"],
                "config": {"model": self._config["llm_model"]},
            }
        if "embedding_model" in self._config:
            mem_config["embedder"] = {
                "provider": self._config.get("embedding_provider", "openai"),
                "config": {"model": self._config["embedding_model"]},
            }

        return Memory.from_config(mem_config)

    # ---------- metadata ----------

    def capabilities(self) -> dict[str, Any]:
        try:
            version = importlib.metadata.version("mem0ai")
        except importlib.metadata.PackageNotFoundError:
            version = "unknown"
        return {
            "provider_name": "mem0",
            "provider_version": version,
            "supported_retrieval_modes": sorted(SUPPORTED_MODES),
            "backing_store": "postgres",
            "supports_user_scoping": True,
            # Writes tag memories with agent_id=session_id, but Mem0 OSS's
            # search() does not filter by agent_id, so end-to-end session
            # isolation is NOT enforced. We advertise False to prevent the
            # engine's isolation audits from trusting a guarantee the read
            # path does not provide.
            "supports_session_scoping": False,
            "declared_cost_model": None,
            "extra": {
                "session_scoping_on_write": True,
                "session_scoping_on_read": False,
            },
        }

    def health(self) -> None:
        # A round-trip that will raise if the backing store is unreachable.
        # The raised message is intentionally generic; the full exception
        # chain is preserved via `from exc` and the traceback is logged for
        # operators so connection/auth details are available in logs without
        # being embedded in the error string returned to callers.
        try:
            self._client.get_all(user_id="__vbench_healthcheck__")
        except Exception as exc:  # pragma: no cover - backing outage
            logger.exception("mem0 health check failed")
            raise RuntimeError("mem0 backing store unhealthy") from exc

    # ---------- writes ----------

    def add_message(
        self,
        *,
        user_id: str,
        session_id: str,
        role: str,
        content: str,
        metadata: dict[str, Any] | None = None,
    ) -> dict[str, Any]:
        start = time.perf_counter()
        result = self._client.add(
            messages=[{"role": role, "content": content}],
            user_id=user_id,
            agent_id=session_id,
            metadata=metadata or {},
        )
        elapsed_ms = (time.perf_counter() - start) * 1000.0
        provider_id = _extract_first_id(result)
        return {
            "provider_id": provider_id,
            "latency_ms": elapsed_ms,
            # Whitespace-split approximation; replaced with a model-matching
            # tokenizer as part of v0.1 finish (see TODO.md).
            "tokens_written": len(content.split()),
            "extra": {"raw": result} if isinstance(result, dict) else {},
        }

    def add_fact(
        self,
        *,
        user_id: str,
        session_id: str,
        fact: str,
        metadata: dict[str, Any] | None = None,
    ) -> dict[str, Any]:
        return self.add_message(
            user_id=user_id,
            session_id=session_id,
            role="user",
            content=fact,
            metadata=metadata,
        )

    # ---------- reads ----------

    def search(
        self,
        *,
        user_id: str,
        session_id: str,  # noqa: ARG002 - accepted to match the vbench HTTP contract; Mem0 OSS's search() does not scope by agent_id, and we keep cross-session recall as the voice-realistic default.
        query: str,
        mode: str,
        top_k: int,
        filters: dict[str, Any] | None = None,
    ) -> dict[str, Any]:
        if mode not in SUPPORTED_MODES:
            # Voice-contract: unsupported modes surface as 422 capability_not_supported
            # at the HTTP layer so the engine records them as SKIPPED.
            raise CapabilityNotSupported(
                capability=mode,
                reason=f"Mem0 OSS supports {sorted(SUPPORTED_MODES)}; got {mode!r}.",
            )
        start = time.perf_counter()
        kwargs: dict[str, Any] = {"query": query, "user_id": user_id, "limit": top_k}
        if filters:
            kwargs["filters"] = filters
        raw = self._client.search(**kwargs)
        elapsed_ms = (time.perf_counter() - start) * 1000.0
        items = _coerce_items(raw)
        payload = "\n".join(f"- {it['content']}" for it in items)
        return {
            "items": items,
            "latency_ms": elapsed_ms,
            "retrieval_mode": mode,
            "token_footprint": len(payload.split()),
            "extra": {},
        }

    def enumerate(
        self,
        *,
        user_id: str,
        session_id: str | None = None,  # noqa: ARG002 - accepted for API compatibility; Mem0 OSS does not scope get_all by agent_id.
    ) -> list[dict[str, Any]]:
        raw = self._client.get_all(user_id=user_id)
        return _coerce_items(raw)

    def reset(
        self,
        *,
        user_id: str,
        session_id: str | None = None,  # noqa: ARG002 - accepted for API compatibility; Mem0 OSS delete_all is user-scoped only.
    ) -> None:
        # Upstream mem0 can raise when the user has no memories to delete —
        # we log the failure at warning level so benign no-ops are visible
        # in sidecar logs without aborting the benchmark, but real
        # auth/network errors are not silently discarded.
        try:
            self._client.delete_all(user_id=user_id)
        except Exception as exc:  # noqa: BLE001 - mem0 does not expose typed errors
            logger.warning("mem0 reset(user_id=%s) failed: %s", user_id, exc)


def _coerce_items(raw: Any) -> list[dict[str, Any]]:
    if isinstance(raw, dict) and "results" in raw:
        raw = raw["results"]
    if not isinstance(raw, list):
        return []
    out: list[dict[str, Any]] = []
    for m in raw:
        if not isinstance(m, dict):
            continue
        out.append(
            {
                "item_id": str(m.get("id", "")),
                "content": str(m.get("memory") or m.get("text") or m.get("content") or ""),
                "score": m.get("score"),
                "created_at": m.get("created_at"),
                "metadata": m.get("metadata") or {},
            }
        )
    return out


def _extract_first_id(result: Any) -> str:
    if isinstance(result, dict):
        if "results" in result and isinstance(result["results"], list) and result["results"]:
            first = result["results"][0]
            if isinstance(first, dict) and "id" in first:
                return str(first["id"])
        if "id" in result:
            return str(result["id"])
    return ""
