# Benchmark Datasets

## LoCoMo (MVP)

- Full name: Long-Horizon Conversational Memory
- Citation: Maharana et al., ACL 2024 — https://arxiv.org/abs/2402.17753
- Source: https://huggingface.co/datasets/snap-research/locomo

LoCoMo contains multi-session conversations between two speakers spanning
months, each with QA pairs that require recalling facts from earlier sessions.

```bash
./vbench datasets download locomo
./vbench datasets info locomo
```

vbench loads LoCoMo by flattening its session buckets into a single turn list
while preserving `session_id` on each turn. Search at test time is scoped to
the last session, which matches the voice-agent scenario: a returning caller
asking about something from earlier in *their* history.

## Roadmap datasets (v0.2+)

- **LongMemEval** (Wu et al., ICLR 2025 — https://arxiv.org/abs/2410.10813).
  Exercises temporal reasoning, knowledge updates, and the abstention case
  (safety-critical for voice agents).
- **Custom JSONL**. Supply benchmark data that matches `internal/schema.BenchmarkItem`:
  ```jsonl
  {"item_id": "...", "dataset": "custom", "conversation": [...], "questions": [...]}
  ```

Neither is wired up in the MVP.
