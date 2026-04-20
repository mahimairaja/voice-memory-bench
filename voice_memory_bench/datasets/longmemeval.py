"""
LongMemEval dataset loader.

LongMemEval is a benchmark for evaluating long-term interactive memory in chat
assistants. It tests five key memory abilities: single-session preference recall,
cross-session preference tracking, temporal reasoning over memory timelines,
knowledge updates when facts change, and abstention when memory is absent.

Citation
--------
Wu et al., "LongMemEval: Benchmarking Chat Assistants on Long-Term Interactive Memory",
ICLR 2025. https://arxiv.org/abs/2410.10813

Dataset source
--------------
https://huggingface.co/datasets/xiaowu0162/longmemeval
"""

from __future__ import annotations

import pathlib
from collections.abc import Iterator

from voice_memory_bench.core.schemas import BenchmarkItem
from voice_memory_bench.datasets.base import DATASET_CACHE_DIR

LONGMEMEVAL_EXPECTED_SHA256 = "TODO(mahimai): fill in after verifying the canonical download"
LONGMEMEVAL_URL = (
    "https://huggingface.co/datasets/xiaowu0162/longmemeval/resolve/main/longmemeval_s.json"
)


class LongMemEvalLoader:
    """Loader for the LongMemEval benchmark dataset."""

    def download(self, cache_dir: pathlib.Path = DATASET_CACHE_DIR) -> None:
        """Download LongMemEval from HuggingFace and verify SHA-256."""
        raise NotImplementedError(
            "TODO(mahimai): download LONGMEMEVAL_URL to cache_dir/longmemeval/, verify SHA-256"
        )

    def is_cached(self, cache_dir: pathlib.Path = DATASET_CACHE_DIR) -> bool:
        """Return True if the LongMemEval data is already cached and verified."""
        raise NotImplementedError("TODO(mahimai): check for existence of cached file")

    def load(
        self,
        subset: str | None = None,
        max_items: int | None = None,
        cache_dir: pathlib.Path = DATASET_CACHE_DIR,
    ) -> Iterator[BenchmarkItem]:
        """Yield BenchmarkItem objects from the LongMemEval dataset."""
        raise NotImplementedError(
            "TODO(mahimai): parse the LongMemEval JSON format into BenchmarkItem instances"
        )
