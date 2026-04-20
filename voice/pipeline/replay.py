"""
Voice transcript replay pipeline stage.

This stage replays a text conversation as a sequence of simulated voice turns,
injecting realistic inter-turn gaps derived from the LoCoMo conversation timing
metadata. It drives the memory provider's write path exactly as a live voice
agent would: one add_message() call per turn, in order, with wall-clock timing
preserved.

Phase 2 implementation notes
-----------------------------
* Input: a BenchmarkItem with ConversationTurns.
* Output: a sequence of WriteResult objects with per-turn latencies.
* The replay speed can be scaled (e.g. 10x real-time) to stress-test the
  write path under higher-than-realtime load.
* For voice-shaped discourse, turns are shorter and more frequent than
  text-chat benchmarks. Disfluencies ("um", "uh") should be preserved
  rather than stripped, as they affect tokenisation and embedding quality.
"""

from __future__ import annotations


class TranscriptReplayStage:
    """Replays a conversation transcript against a memory provider."""

    async def run(self) -> None:
        """Execute the replay."""
        raise NotImplementedError("TODO(mahimai): Phase 2 — implement transcript replay")
