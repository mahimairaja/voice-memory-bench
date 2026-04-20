"""Unit tests for MemScore and related schemas."""
from __future__ import annotations

import pytest
from pydantic import ValidationError

from voice_memory_bench.core.schemas import MemScore


@pytest.mark.unit
def test_memscore_has_three_axes() -> None:
    """MemScore must expose quality, latency, and cost as separate fields."""
    score = MemScore(
        quality=0.85,
        latency_p50_ms=42.0,
        latency_p95_ms=98.0,
        latency_p99_ms=150.0,
        cost_per_item=0.0,
        token_footprint_p50=256,
    )
    assert score.quality == 0.85
    assert score.latency_p50_ms == 42.0
    assert score.cost_per_item == 0.0


@pytest.mark.unit
def test_memscore_quality_bounded() -> None:
    """quality must be in [0, 1]."""
    with pytest.raises(ValidationError):
        MemScore(
            quality=1.5,  # out of range
            latency_p50_ms=10.0,
            latency_p95_ms=20.0,
            latency_p99_ms=30.0,
            cost_per_item=0.0,
            token_footprint_p50=100,
        )


@pytest.mark.unit
def test_memscore_serialises_to_dict() -> None:
    """MemScore must serialise to a dict with all three axes present."""
    score = MemScore(
        quality=0.9,
        latency_p50_ms=50.0,
        latency_p95_ms=100.0,
        latency_p99_ms=200.0,
        cost_per_item=0.01,
        token_footprint_p50=512,
    )
    d = score.model_dump()
    assert "quality" in d
    assert "latency_p50_ms" in d
    assert "cost_per_item" in d
