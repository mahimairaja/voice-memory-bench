"""CLI subcommand for comparing multiple runs."""
from __future__ import annotations
import pathlib
import typer

app = typer.Typer(help="Compare benchmark runs.", no_args_is_help=False)


@app.callback(invoke_without_command=True)
def compare_runs(
    run_dirs: list[pathlib.Path] = typer.Argument(..., help="Two or more run directories to compare."),
    output: pathlib.Path = typer.Option(
        pathlib.Path("comparison.md"),
        "--output",
        "-o",
        help="Output file for the comparison report.",
    ),
    dry_run: bool = typer.Option(False, "--dry-run"),
) -> None:
    """
    Produce a side-by-side MemScore comparison of two or more benchmark runs.

    Output is a Markdown table showing quality, latency, and cost for each run.
    """
    raise NotImplementedError("TODO(mahimai): read report.json from each run dir and produce comparison")
