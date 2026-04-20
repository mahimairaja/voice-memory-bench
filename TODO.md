# TODO

Working list for the vbench MVP. See `ROADMAP.md` for version planning.

## Finish v0.1

- [ ] Run the end-to-end smoke test (`vbench eval --config examples/configs/mem0-locomo.yaml --max-items 2`) against a real Postgres instance and capture the first headline line.
- [ ] Publish first MemScore JSON from a 1x + 4x run.
- [ ] Verify LoCoMo schema handling against the actual upstream file; add fixtures for alternate shapes if encountered.
- [ ] Replace whitespace-split token counting with a tokenizer that matches the answer LLM.
- [ ] Add a `--dry-run` that exercises the sidecar lifecycle without hitting the LLM APIs.

## v0.2 (tracked here until moved into Issues)

- [ ] `sidecars/memori` package + `providers/memori` docker-compose.
- [ ] `sidecars/graphiti` package + `providers/graphiti` docker-compose.
- [ ] `sidecars/cognee` package + `providers/cognee` docker-compose.
- [ ] 16x concurrency level with explicit queue/back-pressure modelling.
- [ ] LongMemEval loader.
- [ ] Custom JSONL loader.
- [ ] `vbench compare <run-id> <run-id>` command.
