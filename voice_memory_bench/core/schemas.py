"""
Canonical record schemas shared across pipeline stages.

These are the data shapes that flow between stages via artifact files.
"""

from __future__ import annotations

from typing import Any
from pydantic import BaseModel, Field
import datetime


class ConversationTurn(BaseModel):
    """A single turn in a multi-turn conversation."""
    turn_id: str
    session_id: str
    user_id: str
    role: str = Field(..., description="'user' or 'assistant'")
    content: str
    timestamp: datetime.datetime | None = None
    metadata: dict[str, Any] = Field(default_factory=dict)


class EvaluationQuestion(BaseModel):
    """A question derived from the conversation, with a reference answer."""
    question_id: str
    question: str
    reference_answer: str
    question_type: str | None = Field(
        None,
        description="Question category, e.g. 'factual', 'temporal', 'preference'.",
    )


class BenchmarkItem(BaseModel):
    """
    A single benchmark item as normalised by the ingest stage.

    Each item contains a conversation history and one or more evaluation
    questions with reference answers.
    """
    item_id: str
    dataset: str
    subset: str | None = None
    conversation: list[ConversationTurn]
    questions: list[EvaluationQuestion]
    metadata: dict[str, Any] = Field(default_factory=dict)


class IndexArtifact(BaseModel):
    """Per-item artifact written by the index stage."""
    item_id: str
    provider: str
    write_results: list[dict[str, Any]]
    total_latency_ms: float
    p50_latency_ms: float
    p95_latency_ms: float


class SearchArtifact(BaseModel):
    """Per-question artifact written by the search stage."""
    item_id: str
    question_id: str
    provider: str
    retrieval_result: dict[str, Any]
    memory_payload: str = Field(..., description="The exact text that would be injected into the prompt.")


class AnswerArtifact(BaseModel):
    """Per-question artifact written by the answer stage."""
    item_id: str
    question_id: str
    provider: str
    prompt: str
    completion: str
    prompt_tokens: int
    completion_tokens: int
    latency_ms: float


class MemScore(BaseModel):
    """
    The MemScore triple: quality, latency, and cost.

    Never collapse to a scalar. All three axes are independent and reported
    side-by-side in every output.
    """
    quality: float = Field(..., ge=0.0, le=1.0, description="Normalised answer quality score [0, 1].")
    latency_p50_ms: float = Field(..., description="p50 retrieval latency across all questions.")
    latency_p95_ms: float = Field(..., description="p95 retrieval latency.")
    latency_p99_ms: float = Field(..., description="p99 retrieval latency.")
    cost_per_item: float = Field(..., description="Average cost per benchmark item in USD, or 0 if self-hosted.")
    token_footprint_p50: int = Field(..., description="p50 token count of injected memory payload.")


class EvaluationArtifact(BaseModel):
    """Per-item artifact written by the evaluate stage."""
    item_id: str
    provider: str
    dataset: str
    per_question_scores: list[dict[str, Any]]
    memscore: MemScore
