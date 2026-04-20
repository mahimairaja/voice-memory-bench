# Contributing

## Local dev

```bash
git clone https://github.com/mahimairaja/vbench
cd vbench
make build          # Go engine → ./vbench
make sidecar-sync   # uv sync for sidecars/mem0
make test           # go test ./...
make vet            # go vet ./...
```

## Adding a provider (v0.2)

New providers follow the sidecar shape: a Python package under `sidecars/<name>/`
exposing an HTTP server on `127.0.0.1:$VBENCH_SIDECAR_PORT`, with endpoints
`health`, `capabilities`, `add_message`, `add_fact`, `search`, `enumerate`,
`reset`. Unsupported capabilities return `422 capability_not_supported`.

1. Create `sidecars/<name>/` with a `pyproject.toml` that exposes a console
   script entrypoint (`[project.scripts]`).
2. Implement the adapter + FastAPI app (see `sidecars/mem0/` for shape).
3. Add `providers/<name>/docker-compose.yml` + README.
4. Register a default command in `internal/sidecar/<name>.go`.
5. Add an example config under `examples/configs/`.
6. Update `docs/providers.md`.

## Code style

- Go: `go vet ./...` must be clean. Use `go fmt`. No new lint dependencies
  in v0.1.
- Python: `ruff` (line length 100). Stay on Python 3.11+ features.
- Keep the public HTTP contract stable — any change to endpoint shape is a
  breaking change that bumps the minor version.
