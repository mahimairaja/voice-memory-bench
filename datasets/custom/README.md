# Custom Dataset Format

`voice-memory-bench` supports user-supplied benchmark datasets as a **JSONL file**
(one JSON object per line). This lets you benchmark memory providers against your
own conversational data — call transcripts, support conversations, domain-specific
dialogue — without waiting for official dataset support.

## JSONL schema

Each line must be a valid JSON object matching this structure:

```json
{
  "item_id": "unique-string-id",
  "dataset": "custom",
  "subset": null,
  "conversation": [
    {
      "turn_id": "turn-001",
      "session_id": "session-abc",
      "user_id": "user-123",
      "role": "user",
      "content": "My name is Alice and I'm allergic to peanuts.",
      "timestamp": "2024-01-15T10:00:00Z",
      "metadata": {}
    },
    {
      "turn_id": "turn-002",
      "session_id": "session-abc",
      "user_id": "user-123",
      "role": "assistant",
      "content": "Got it, Alice. I've noted your peanut allergy.",
      "timestamp": "2024-01-15T10:00:05Z",
      "metadata": {}
    }
  ],
  "questions": [
    {
      "question_id": "q-001",
      "question": "What allergy does the user have?",
      "reference_answer": "The user is allergic to peanuts.",
      "question_type": "factual"
    }
  ],
  "metadata": {}
}
```

## Field reference

### Top-level fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `item_id` | string | ✅ | Unique identifier for this benchmark item. Must be unique within the file. |
| `dataset` | string | ✅ | Should be `"custom"` for custom datasets. |
| `subset` | string \| null | ❌ | Optional subset label (e.g. `"voice-shaped"`, `"medical"`). |
| `conversation` | array | ✅ | List of conversation turns. See below. |
| `questions` | array | ✅ | List of evaluation questions with reference answers. See below. |
| `metadata` | object | ❌ | Arbitrary key-value metadata. Passed through to artifacts. |

### `conversation[*]` fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `turn_id` | string | ✅ | Unique turn identifier within the item. |
| `session_id` | string | ✅ | Session identifier. Used to scope memory writes. |
| `user_id` | string | ✅ | User identifier. |
| `role` | string | ✅ | `"user"` or `"assistant"`. |
| `content` | string | ✅ | The text of the turn. |
| `timestamp` | string \| null | ❌ | ISO 8601 timestamp. Used by temporal retrieval. |
| `metadata` | object | ❌ | Arbitrary metadata. |

### `questions[*]` fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `question_id` | string | ✅ | Unique question identifier. |
| `question` | string | ✅ | The evaluation question. |
| `reference_answer` | string | ✅ | The ground-truth answer used by the judge LLM. |
| `question_type` | string \| null | ❌ | Informal category: `"factual"`, `"temporal"`, `"preference"`, `"reasoning"`. |

## Using your dataset

```yaml
# In your run config:
dataset:
  name: custom
  path: datasets/custom/my_conversations.jsonl
  max_items: 100   # optional limit for smoke tests
```

## Contributing datasets

This is the directory where community contributions of voice-shaped benchmarks
will land. To contribute, open a PR adding a JSONL file and a `README.md`
describing the source, curation method, and any licensing constraints.
