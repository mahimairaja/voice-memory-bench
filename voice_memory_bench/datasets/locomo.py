"""
LoCoMo dataset loader.

LoCoMo (Long-Horizon Conversational Memory) is a benchmark for evaluating
memory systems in long-horizon conversations. It contains multi-session
conversations with associated QA pairs.

Citation
--------
Maharana et al., "Evaluating Very Long-Term Conversational Memory of LLM Agents",
ACL 2024. https://arxiv.org/abs/2402.17753

Dataset source
--------------
https://huggingface.co/datasets/snap-research/locomo
"""

from __future__ import annotations

import pathlib
from typing import Iterator

from voice_memory_bench.core.schemas import BenchmarkItem
from voice_memory_bench.datasets.base import DATASET_CACHE_DIR


LOCOMO_EXPECTED_SHA256 = "TODO(mahimai): fill in after verifying the canonical download"
LOCOMO_URL = "https://huggingface.co/datasets/snap-research/locomo/resolve/main/locomo10_test.json"


class LoCoMoLoader:
    """Loader for the LoCoMo benchmark dataset."""

    def download(self, cache_dir: pathlib.Path = DATASET_CACHE_DIR) -> None:
        """Download LoCoMo from HuggingFace and verify SHA-256."""
        raise NotImplementedError(
            "TODO(mahimai): download LOCOMO_URL to cache_dir/locomo/, verify SHA-256"
        )

    def is_cached(self, cache_dir: pathlib.Path = DATASET_CACHE_DIR) -> bool:
        """Return True if the LoCoMo data is already cached and verified."""
        raise NotImplementedError("TODO(mahimai): check for existence of cached file")

    def load(
        self,
        subset: str | None = None,
        max_items: int | None = None,
        cache_dir: pathlib.Path = DATASET_CACHE_DIR,
    ) -> Iterator[BenchmarkItem]:
        """Yield BenchmarkItem objects from the LoCoMo dataset."""
        raise NotImplementedError(
            "TODO(mahimai): parse the LoCoMo JSON format into BenchmarkItem instances"
        )
