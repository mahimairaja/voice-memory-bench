"""JSON report generator."""

from __future__ import annotations

import pathlib

from voice_memory_bench.core.schemas import EvaluationArtifact


class JsonReporter:
    """Writes report.json to the run directory."""

    def write(self, artifacts: list[EvaluationArtifact], output_dir: pathlib.Path) -> pathlib.Path:
        """Write the full MemScore report as JSON. Returns the path written."""
        raise NotImplementedError("TODO(mahimai): serialise artifacts to report.json")
