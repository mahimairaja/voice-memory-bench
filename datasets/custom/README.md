# Custom datasets

> **Status:** Custom JSONL loading lands in v0.2. This doc describes the target
> schema so contributors can start preparing datasets now. The MVP engine only
> registers `locomo`.

## Target JSONL schema

Each line is a JSON object that matches `internal/schema.BenchmarkItem`:

```json
{
  "item_id": "unique-id",
  "dataset": "custom",
  "conversation": [
    {
      "turn_id": "t1",
      "session_id": "s1",
      "user_id": "u1",
      "role": "user",
      "content": "My name is Alice and I'm allergic to peanuts.",
      "timestamp": "2024-01-15T10:00:00Z"
    }
  ],
  "questions": [
    {
      "question_id": "q1",
      "question": "What allergy does the user have?",
      "reference_answer": "peanuts"
    }
  ]
}
```

## Contributing a dataset

Open a PR that adds your JSONL file and a short README describing source,
curation method, and any licensing constraints.
