"""
Pydantic v2 configuration schema for voice-memory-bench run configs.

A run config is a YAML file validated against :class:`RunConfig`. Secrets
(API keys, DB passwords) are never stored in the config file directly; instead
they are referenced via env-var interpolation: ${MY_SECRET}.
"""

from __future__ import annotations

from typing import Any
from pydantic import BaseModel, Field, model_validator


class ProviderConfig(BaseModel):
    """Provider-specific configuration block."""
    name: str = Field(..., description="Provider name: 'mem0', 'memori', 'graphiti', or 'cognee'.")
    config: dict[str, Any] = Field(
        default_factory=dict,
        description="Provider-specific key-value config. "
                    "Secrets should be env-var references like ${POSTGRES_PASSWORD}.",
    )


class LLMConfig(BaseModel):
    """Configuration for an LLM used in the answer or judge stage."""
    model: str = Field(..., description="Model identifier, e.g. 'gpt-4o' or 'ollama/llama3'.")
    base_url: str | None = Field(None, description="Override base URL for self-hosted models.")
    api_key_env: str | None = Field(
        None,
        description="Name of the environment variable holding the API key.",
    )
    temperature: float = Field(0.0, description="Sampling temperature. Use 0 for determinism.")
    max_tokens: int = Field(1024)
    seed: int = Field(42, description="Random seed for reproducibility.")


class DatasetConfig(BaseModel):
    """Dataset selection."""
    name: str = Field(..., description="'locomo', 'longmemeval', or 'custom'.")
    subset: str | None = Field(None, description="Dataset subset or split, if applicable.")
    path: str | None = Field(
        None,
        description="Path to a custom JSONL file. Required when name='custom'.",
    )
    max_items: int | None = Field(
        None,
        description="Limit the number of benchmark items. Useful for smoke tests.",
    )


class RunConfig(BaseModel):
    """
    Top-level run configuration. One file per benchmark run.

    Example
    -------
    See examples/configs/ for commented example files.
    """
    run_name: str = Field(..., description="Human-readable name for this run.")
    dataset: DatasetConfig
    provider: ProviderConfig
    answer_llm: LLMConfig
    judge_llm: LLMConfig
    concurrency: int = Field(1, ge=1, description="Number of concurrent benchmark workers.")
    seed: int = Field(42, description="Global random seed. All sub-seeds are derived from this.")
    output_dir: str = Field("runs", description="Root directory for run output.")
    dry_run: bool = Field(False, description="If True, print what would happen but do not execute.")

    @model_validator(mode="after")
    def custom_dataset_needs_path(self) -> "RunConfig":
        if self.dataset.name == "custom" and self.dataset.path is None:
            raise ValueError("dataset.path is required when dataset.name='custom'")
        return self
