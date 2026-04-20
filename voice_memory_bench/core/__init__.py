"""Core abstractions: adapter protocol, pipeline stages, config schema."""

from voice_memory_bench.core.adapter import (
    MemoryAdapter,
    CapabilityDescriptor,
    CapabilityNotSupportedError,
    AdapterConfigError,
    AdapterHealthError,
    RetrievalMode,
    BackingStore,
    WriteResult,
    MemoryItem,
    RetrievalResult,
)
from voice_memory_bench.core.config import RunConfig
from voice_memory_bench.core.schemas import MemScore, BenchmarkItem

__all__ = [
    "MemoryAdapter",
    "CapabilityDescriptor",
    "CapabilityNotSupportedError",
    "AdapterConfigError",
    "AdapterHealthError",
    "RetrievalMode",
    "BackingStore",
    "WriteResult",
    "MemoryItem",
    "RetrievalResult",
    "RunConfig",
    "MemScore",
    "BenchmarkItem",
]
