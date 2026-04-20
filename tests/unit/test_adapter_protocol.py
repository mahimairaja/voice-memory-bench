"""Unit tests for adapter protocol and error types."""

from __future__ import annotations

import pytest

from voice_memory_bench.core.adapter import (
    AdapterConfigError,
    AdapterHealthError,
    CapabilityNotSupportedError,
    RetrievalMode,
)


@pytest.mark.unit
def test_capability_not_supported_error_fields() -> None:
    """CapabilityNotSupportedError must carry provider and capability."""
    err = CapabilityNotSupportedError(
        provider="mem0",
        capability=RetrievalMode.TEMPORAL,
        reason="Mem0 OSS does not support temporal queries.",
    )
    assert err.provider == "mem0"
    assert err.capability == RetrievalMode.TEMPORAL
    assert "temporal" in err.reason.lower()


@pytest.mark.unit
def test_capability_not_supported_error_message() -> None:
    """Error message should include provider name and capability."""
    err = CapabilityNotSupportedError(
        provider="cognee",
        capability=RetrievalMode.KEYWORD,
        reason="No BM25 index.",
    )
    msg = str(err)
    assert "cognee" in msg
    assert "keyword" in msg.lower()


@pytest.mark.unit
def test_adapter_config_error_is_exception() -> None:
    """AdapterConfigError must be a subclass of Exception."""
    err = AdapterConfigError("missing postgres_url")
    assert isinstance(err, Exception)
    assert "postgres_url" in str(err)


@pytest.mark.unit
def test_adapter_health_error_is_exception() -> None:
    """AdapterHealthError must be a subclass of Exception."""
    err = AdapterHealthError("connection refused")
    assert isinstance(err, Exception)
