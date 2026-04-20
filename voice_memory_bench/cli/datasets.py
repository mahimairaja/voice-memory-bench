"""CLI subcommands for dataset management."""

from __future__ import annotations

import typer

app = typer.Typer(help="Download and inspect benchmark datasets.", no_args_is_help=True)


@app.command("download")
def download_dataset(
    dataset: str = typer.Argument(..., help="Dataset name: locomo or longmemeval."),
    cache_dir: str = typer.Option(None, "--cache-dir", help="Override cache directory."),
    dry_run: bool = typer.Option(False, "--dry-run"),
) -> None:
    """Download a benchmark dataset and verify its hash."""
    raise NotImplementedError(
        "TODO(mahimai): look up loader in DATASET_REGISTRY and call download()"
    )


@app.command("info")
def dataset_info(
    dataset: str = typer.Argument(..., help="Dataset name."),
) -> None:
    """Print metadata about a dataset (size, splits, citation)."""
    raise NotImplementedError("TODO(mahimai): print dataset metadata")
