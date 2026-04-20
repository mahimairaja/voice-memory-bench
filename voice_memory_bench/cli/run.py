"""CLI subcommand for running a full benchmark."""

from __future__ import annotations

import pathlib

import typer

app = typer.Typer(help="Run a benchmark.", no_args_is_help=False)


@app.callback(invoke_without_command=True)
def run_benchmark(
    config: pathlib.Path = typer.Argument(..., help="Path to run config YAML."),
    run_id: str | None = typer.Option(
        None, "--run-id", help="Override run ID for reproducibility."
    ),
    stages: str | None = typer.Option(
        None,
        "--stages",
        help="Comma-separated list of stages to run (default: all). E.g. 'index,search,answer'.",
    ),
    resume: bool = typer.Option(
        True, "--resume/--no-resume", help="Resume from last completed stage."
    ),
    dry_run: bool = typer.Option(False, "--dry-run"),
    verbose: bool = typer.Option(False, "--verbose", "-v"),
) -> None:
    """
    Run a full benchmark pipeline from a config file.

    Stages: ingest -> index -> search -> answer -> evaluate

    Each stage writes artifacts to runs/<run_id>/. Re-running with the same
    run ID resumes from the first incomplete stage.
    """
    raise NotImplementedError("TODO(mahimai): load config, build pipeline, execute stages")
