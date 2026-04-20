"""
Custom JSONL dataset loader.

Users can supply their own benchmark data as a JSONL file where each line is a
JSON object matching the :class:`BenchmarkItem` schema.

See datasets/custom/README.md for the full schema specification and examples.
"""

from __future__ import annotations

import pathlib
from collections.abc import Iterator

from voice_memory_bench.core.schemas import BenchmarkItem


class CustomDatasetLoader:
    """Loader for user-supplied JSONL benchmark data."""

    def __init__(self, path: pathlib.Path) -> None:
        self._path = path

    def download(self, cache_dir: pathlib.Path | None = None) -> None:
        """No-op: custom datasets are already on disk."""

    def is_cached(self, cache_dir: pathlib.Path | None = None) -> bool:
        return self._path.exists()

    def load(
        self,
        subset: str | None = None,
        max_items: int | None = None,
        cache_dir: pathlib.Path | None = None,
    ) -> Iterator[BenchmarkItem]:
        """Yield BenchmarkItem objects from the JSONL file."""
        raise NotImplementedError(
            "TODO(mahimai): open self._path, parse each line as BenchmarkItem"
        )
