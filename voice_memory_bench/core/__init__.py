"""Core abstractions: adapter protocol, pipeline stages, config schema."""

from voice_memory_bench.core.adapter import (
    AdapterConfigError,
    AdapterHealthError,
    BackingStore,
    CapabilityDescriptor,
    CapabilityNotSupportedError,
    MemoryAdapter,
    MemoryItem,
    RetrievalMode,
    RetrievalResult,
    WriteResult,
)
from voice_memory_bench.core.config import RunConfig
from voice_memory_bench.core.schemas import BenchmarkItem, MemScore

__all__ = [
    "AdapterConfigError",
    "AdapterHealthError",
    "BackingStore",
    "BenchmarkItem",
    "CapabilityDescriptor",
    "CapabilityNotSupportedError",
    "MemScore",
    "MemoryAdapter",
    "MemoryItem",
    "RetrievalMode",
    "RetrievalResult",
    "RunConfig",
    "WriteResult",
]
