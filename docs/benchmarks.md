# Benchmark Datasets

## LoCoMo

**Full name:** Long-Horizon Conversational Memory  
**Citation:** Maharana et al., ACL 2024. https://arxiv.org/abs/2402.17753  
**Source:** https://huggingface.co/datasets/snap-research/locomo

LoCoMo contains multi-session conversations between two speakers spanning months. Each conversation includes QA pairs that require recalling facts from earlier sessions.

**Why it's useful for voice agents:** Real conversations span multiple calls. LoCoMo tests whether a memory system can recall facts from a call that happened weeks ago.

**Download:**
```bash
uv run vmb datasets download locomo
```

---

## LongMemEval

**Full name:** LongMemEval  
**Citation:** Wu et al., ICLR 2025. https://arxiv.org/abs/2410.10813  
**Source:** https://huggingface.co/datasets/xiaowu0162/longmemeval

LongMemEval tests five memory abilities:
1. Single-session preference recall
2. Cross-session preference tracking
3. Temporal reasoning over memory timelines
4. Knowledge updates when facts change
5. Abstention when memory is absent

**Why it's useful for voice agents:** LongMemEval's abstention category tests the critical case where a provider confidently retrieves wrong memories — a safety issue in voice agents.

**Download:**
```bash
uv run vmb datasets download longmemeval
```

---

## Custom JSONL

Supply your own benchmark data in the `BenchmarkItem` schema:

```jsonl
{"item_id": "item-001", "dataset": "custom", "conversation": [...], "questions": [...]}
```

See `voice_memory_bench/core/schemas.py` for the full schema definition.

**Usage:**
```yaml
dataset:
  name: custom
  path: /path/to/my_bench.jsonl
```
