"""Markdown report generator."""

from __future__ import annotations

import pathlib

from voice_memory_bench.core.schemas import EvaluationArtifact


class MarkdownReporter:
    """Writes report.md to the run directory."""

    def write(self, artifacts: list[EvaluationArtifact], output_dir: pathlib.Path) -> pathlib.Path:
        """Write a human-readable MemScore summary as Markdown. Returns the path written."""
        raise NotImplementedError("TODO(mahimai): produce a Markdown table of MemScore triples")
