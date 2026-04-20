"""
voice-memory-bench CLI entry point.

Usage
-----
    uv run vmb --help
    uv run vmb providers list
    uv run vmb datasets download locomo
    uv run vmb run examples/configs/mem0-locomo.yaml
    uv run vmb compare runs/run-a runs/run-b

Every command supports --dry-run, --verbose / --quiet, and --run-id.
"""

from __future__ import annotations

import typer

from voice_memory_bench.cli import compare as compare_cmd
from voice_memory_bench.cli import datasets as datasets_cmd
from voice_memory_bench.cli import providers as providers_cmd
from voice_memory_bench.cli import run as run_cmd

app = typer.Typer(
    name="vmb",
    help="voice-memory-bench: benchmark self-hostable AI memory frameworks.",
    no_args_is_help=True,
)

app.add_typer(providers_cmd.app, name="providers")
app.add_typer(datasets_cmd.app, name="datasets")
app.add_typer(run_cmd.app, name="run")
app.add_typer(compare_cmd.app, name="compare")


if __name__ == "__main__":
    app()
