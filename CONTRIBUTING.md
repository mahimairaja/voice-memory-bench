# Contributing to vbench

## Quick start

```bash
git clone https://github.com/mahimairaja/vbench
cd vbench
make build
make test
```

## Development loop

1. Fork and create a feature branch.
2. Engine changes: `go build ./... && go vet ./... && go test ./...`.
3. Sidecar changes: `make sidecar-sync && make sidecar-lint`.
4. Update docs / example configs if the contract moves.
5. Submit a PR.

## Adding a provider

See [`docs/contributing.md`](docs/contributing.md). The MVP wires only Mem0;
Memori / Graphiti / Cognee land in v0.2 through the same sidecar shape.

## Reporting issues

Use the GitHub issue tracker. Security vulnerabilities: see `SECURITY.md`.

## Code of conduct

Contributor Covenant — see `CODE_OF_CONDUCT.md`.
