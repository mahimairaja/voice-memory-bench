"""
LiveKit Agents worker stub for voice-memory-bench.

This module will contain the LiveKit Agents worker that:
1. Receives audio from a LiveKit room (real or simulated).
2. Transcribes speech via the STT adapter (Deepgram or Whisper).
3. Queries the memory provider for relevant context.
4. Passes context + transcript to an LLM to generate a response.
5. Synthesises the response via the TTS adapter (Cartesia or ElevenLabs).
6. Writes the transcript turn to the memory provider.

Phase 2 implementation notes
-----------------------------
* Use livekit-agents >= 0.8 with the pipeline pattern (VoicePipelineAgent).
* The memory provider is injected at construction time, not hardcoded.
* All latency measurements are taken at the agent level, not the provider level,
  to capture end-to-end time including serialisation and network overhead.
* The worker must handle graceful shutdown: flush all pending memory writes
  before exiting, and emit a StageResult to the artifact directory.
"""

from __future__ import annotations


class VoiceMemoryWorker:
    """LiveKit Agents worker that integrates with a memory provider."""

    async def start(self) -> None:
        """Start the worker and connect to the LiveKit room."""
        raise NotImplementedError("TODO(mahimai): Phase 2 — implement LiveKit worker")
