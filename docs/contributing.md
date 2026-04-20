# Contributing

## Setup

```bash
git clone https://github.com/mahimairaja/voice-memory-bench
cd voice-memory-bench
pip install uv
uv sync --extra dev
pre-commit install
```

## Running Tests

```bash
# Unit tests only (no external services)
uv run pytest -m unit

# All tests (requires live services)
uv run pytest
```

## Adding a Provider

1. Create `voice_memory_bench/providers/<name>/adapter.py`
2. Implement all methods of the `MemoryAdapter` protocol
3. Add a `docker-compose.yml` in `providers/<name>/`
4. Add a `config.example.yaml` in `providers/<name>/`
5. Add integration tests in `tests/integration/test_<name>.py`
6. Update `docs/providers.md`

## Versioning Policy

The `MemoryAdapter` protocol is the public API. Changes to it are breaking changes and require a major version bump. Provider adapter internals are not public API.

## Code Style

- `ruff` for linting and formatting (line length 100)
- `mypy --strict` for type checking
- All public functions must have docstrings
- No provider-specific imports in `voice_memory_bench/core/`
