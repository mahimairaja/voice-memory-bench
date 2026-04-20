"""CLI subcommands for managing providers."""

from __future__ import annotations

import typer

app = typer.Typer(help="Manage and inspect memory providers.", no_args_is_help=True)


@app.command("list")
def list_providers(
    verbose: bool = typer.Option(False, "--verbose", "-v"),
) -> None:
    """List all installed provider adapters and their capabilities."""
    raise NotImplementedError("TODO(mahimai): iterate PROVIDER_REGISTRY and print capabilities")


@app.command("check")
def check_provider(
    provider: str = typer.Argument(..., help="Provider name: mem0, memori, graphiti, or cognee."),
    config: str = typer.Option(..., "--config", "-c", help="Path to run config YAML."),
    dry_run: bool = typer.Option(False, "--dry-run"),
) -> None:
    """Run a health check against a configured provider."""
    raise NotImplementedError(
        "TODO(mahimai): load config, instantiate adapter, call health_check()"
    )
