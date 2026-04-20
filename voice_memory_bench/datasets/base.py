"""
Base dataset loader protocol and canonical record schema.

Adding a new dataset is a matter of implementing :class:`DatasetLoader`
and registering it in :data:`DATASET_REGISTRY`.
"""

from __future__ import annotations

import pathlib
from collections.abc import Iterator
from typing import Protocol, runtime_checkable

from voice_memory_bench.core.schemas import BenchmarkItem

DATASET_CACHE_DIR = pathlib.Path.home() / ".cache" / "voice-memory-bench" / "datasets"


@runtime_checkable
class DatasetLoader(Protocol):
    """Protocol for dataset loaders."""

    def download(self, cache_dir: pathlib.Path = DATASET_CACHE_DIR) -> None:
        """
        Download the dataset from its canonical source and cache it locally.

        Implementations must verify the downloaded file's SHA-256 hash against
        a hardcoded expected value. Raise ValueError if the hash does not match.
        """
        ...

    def is_cached(self, cache_dir: pathlib.Path = DATASET_CACHE_DIR) -> bool:
        """Return True if the dataset is already downloaded and verified."""
        ...

    def load(
        self,
        subset: str | None = None,
        max_items: int | None = None,
        cache_dir: pathlib.Path = DATASET_CACHE_DIR,
    ) -> Iterator[BenchmarkItem]:
        """
        Yield benchmark items from the dataset.

        Parameters
        ----------
        subset:
            Dataset split or subset name. Dataset-specific.
        max_items:
            Stop after yielding this many items. None means all.
        cache_dir:
            Where cached data lives.
        """
        ...
