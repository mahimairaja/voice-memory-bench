"""
Pipeline stage definitions for voice-memory-bench.

The benchmark pipeline is: ingest -> index -> search -> answer -> evaluate.

Each stage is independently resumable: it reads from the previous stage's
artifact directory and writes to its own. Re-running with the same run ID
resumes from the first missing artifact, so a crash mid-index does not restart
ingestion from scratch.

Stage implementations live in voice_memory_bench/pipeline/ (one module per stage).
This module defines only the contracts and shared types.
"""

from __future__ import annotations

import datetime
import enum
import pathlib
from typing import Any, Protocol

from pydantic import BaseModel, Field


class StageStatus(str, enum.Enum):
    """Execution status of a pipeline stage."""

    PENDING = "pending"
    RUNNING = "running"
    COMPLETE = "complete"
    FAILED = "failed"
    SKIPPED = "skipped"


class StageResult(BaseModel):
    """Written to disk at the end of each stage."""

    stage: str
    run_id: str
    status: StageStatus
    started_at: datetime.datetime
    finished_at: datetime.datetime | None = None
    artifact_dir: str = Field(
        ..., description="Path to the directory where artifacts were written."
    )
    items_processed: int = 0
    items_skipped: int = 0
    items_failed: int = 0
    error: str | None = None
    extra: dict[str, Any] = Field(default_factory=dict)


class PipelineStage(Protocol):
    """
    Protocol that every pipeline stage module must implement.

    A stage is callable with a run context and returns a :class:`StageResult`.
    Stages are composable — the runner builds a chain and executes each in order,
    passing the previous stage's artifact directory as input.
    """

    async def run(
        self,
        run_id: str,
        artifact_root: pathlib.Path,
        config: dict[str, Any],
        resume: bool = True,
    ) -> StageResult:
        """
        Execute the stage.

        Parameters
        ----------
        run_id:
            Unique identifier for this benchmark run.
        artifact_root:
            Root directory where all run artifacts are stored. Stage should
            write to ``artifact_root / run_id / <stage_name>/``.
        config:
            Full validated run configuration dict.
        resume:
            If True and the stage's output artifact already exists, skip
            re-running and return the existing result.

        Returns
        -------
        StageResult
            Summary of what happened. Always written to disk before returning.
        """
        ...

    def artifact_dir(self, artifact_root: pathlib.Path, run_id: str) -> pathlib.Path:
        """Return the path where this stage's artifacts are written."""
        ...

    def is_complete(self, artifact_root: pathlib.Path, run_id: str) -> bool:
        """Return True if this stage's artifacts exist and are marked complete."""
        ...
