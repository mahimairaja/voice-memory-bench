"""
Deepgram STT adapter stub.

Deepgram Nova-2 is the recommended STT provider for voice-memory-bench because:
* It supports streaming transcription with < 300 ms time-to-first-word latency.
* It produces word-level timestamps that can be used to detect disfluencies.
* It has a self-hosted (on-prem) option for data-sovereignty deployments.

Phase 2 implementation notes
-----------------------------
* Use the deepgram-sdk Python package.
* Implement both streaming (WebSocket) and batch (REST) transcription.
* Pass word_boost hints from the memory context to improve domain accuracy.
* Strip filler words ("um", "uh", "like") only when strip_fillers=True —
  the default is False because fillers affect sentence boundary detection.
"""

from __future__ import annotations


class DeepgramSTTAdapter:
    """Deepgram Nova-2 speech-to-text adapter."""

    async def transcribe(self, audio_bytes: bytes, sample_rate: int = 16000) -> str:
        """Transcribe audio bytes to text."""
        raise NotImplementedError("TODO(mahimai): Phase 2 — implement Deepgram transcription")
