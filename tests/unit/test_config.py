"""Unit tests for RunConfig validation."""

from __future__ import annotations

import pytest
from pydantic import ValidationError

from voice_memory_bench.core.config import RunConfig

MINIMAL_LLM = {
    "model": "gpt-4o",
    "api_key_env": "OPENAI_API_KEY",
}

MINIMAL_PROVIDER = {
    "name": "mem0",
    "config": {"postgres_url": "postgresql://localhost/test"},
}


@pytest.mark.unit
def test_valid_locomo_config() -> None:
    """A complete, valid LoCoMo config should parse without error."""
    config = RunConfig(
        run_name="test",
        dataset={"name": "locomo"},
        provider=MINIMAL_PROVIDER,
        answer_llm=MINIMAL_LLM,
        judge_llm=MINIMAL_LLM,
    )
    assert config.dataset.name == "locomo"
    assert config.concurrency == 1


@pytest.mark.unit
def test_custom_dataset_requires_path() -> None:
    """Custom dataset without path should raise ValidationError."""
    with pytest.raises(ValidationError, match=r"dataset\.path is required"):
        RunConfig(
            run_name="test",
            dataset={"name": "custom"},  # no path
            provider=MINIMAL_PROVIDER,
            answer_llm=MINIMAL_LLM,
            judge_llm=MINIMAL_LLM,
        )


@pytest.mark.unit
def test_custom_dataset_with_path_is_valid() -> None:
    """Custom dataset with path should parse successfully."""
    config = RunConfig(
        run_name="test",
        dataset={"name": "custom", "path": "/data/my_bench.jsonl"},
        provider=MINIMAL_PROVIDER,
        answer_llm=MINIMAL_LLM,
        judge_llm=MINIMAL_LLM,
    )
    assert config.dataset.path == "/data/my_bench.jsonl"


@pytest.mark.unit
def test_concurrency_must_be_positive() -> None:
    """concurrency=0 should fail validation."""
    with pytest.raises(ValidationError):
        RunConfig(
            run_name="test",
            dataset={"name": "locomo"},
            provider=MINIMAL_PROVIDER,
            answer_llm=MINIMAL_LLM,
            judge_llm=MINIMAL_LLM,
            concurrency=0,
        )
