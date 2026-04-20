"""
Wall-clock latency measurement utilities for voice benchmarks.

In voice agents, the memory retrieval call sits on the critical path between
the end of an utterance (ASR finishes) and the start of a response (TTS begins).
The total budget is approximately 300-500 ms for a natural-feeling interaction.

This module provides:
* A context manager that measures wall-clock latency of any async operation.
* A LatencyProfile dataclass that accumulates p50/p95/p99 across N calls.
* A ConcurrentLatencyHarness that fires N concurrent retrieval calls and
  measures tail latency under load — the key metric that flat benchmarks miss.

Phase 2 implementation notes
-----------------------------
* Use time.perf_counter_ns() for high-resolution timing.
* Measure latency from the last byte of the ASR transcript to the first byte
  of memory payload returned — not to the first byte of the LLM response.
* Report p50/p95/p99 across all questions in the benchmark, not just mean.
"""

from __future__ import annotations


class LatencyProfile:
    """Accumulates latency measurements and computes percentile statistics."""

    def record(self, latency_ms: float) -> None:
        """Record a single latency measurement."""
        raise NotImplementedError("TODO(mahimai): Phase 2 — accumulate latency_ms")

    def p50(self) -> float:
        """Return the 50th percentile latency in milliseconds."""
        raise NotImplementedError("TODO(mahimai): Phase 2 — compute p50")

    def p95(self) -> float:
        """Return the 95th percentile latency in milliseconds."""
        raise NotImplementedError("TODO(mahimai): Phase 2 — compute p95")

    def p99(self) -> float:
        """Return the 99th percentile latency in milliseconds."""
        raise NotImplementedError("TODO(mahimai): Phase 2 — compute p99")
