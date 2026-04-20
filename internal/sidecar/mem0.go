package sidecar

// Mem0DefaultCommand is the argv used to launch the Mem0 sidecar when the
// config does not override provider.command. It resolves the `uv run`
// entrypoint declared in sidecars/mem0/pyproject.toml.
var Mem0DefaultCommand = []string{"uv", "run", "vbench-mem0"}

// Mem0DefaultWorkingDir is the sidecar package root, relative to repo root.
var Mem0DefaultWorkingDir = "sidecars/mem0"
