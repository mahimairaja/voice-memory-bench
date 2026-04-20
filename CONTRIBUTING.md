# Contributing to voice-memory-bench

Thank you for your interest in contributing!

## Quick Start

```bash
git clone https://github.com/mahimairaja/voice-memory-bench
cd voice-memory-bench
pip install uv
uv sync --extra dev
pre-commit install
```

## Development Workflow

1. Fork the repository and create a feature branch
2. Make your changes
3. Run tests: `uv run pytest -m unit`
4. Run linter: `uv run ruff check .`
5. Run type checker: `uv run mypy voice_memory_bench/`
6. Submit a pull request

## Implementing a Provider Adapter

The most valuable contribution is a fully-implemented provider adapter. See `docs/contributing.md` for the full guide.

## Reporting Issues

Please use the GitHub issue tracker. For security vulnerabilities, see SECURITY.md.

## Code of Conduct

This project follows the Contributor Covenant. See CODE_OF_CONDUCT.md.
