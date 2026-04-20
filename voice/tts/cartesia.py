"""
Cartesia TTS adapter stub.

Cartesia Sonic is the recommended TTS provider for voice-memory-bench because:
* It achieves < 100 ms time-to-first-audio (TTFA) for streaming synthesis.
* Its Sonic model produces natural prosody that does not telegraph AI origin.
* It has a self-hosted option via Cartesia On-Prem.

Phase 2 implementation notes
-----------------------------
* Use the cartesia Python package.
* Implement streaming synthesis (yield audio chunks as they arrive).
* Record TTFA (time from text input to first audio byte) as a separate metric.
* Support voice cloning via voice_id parameter for persona consistency.
"""

from __future__ import annotations


class CartesiaTTSAdapter:
    """Cartesia Sonic text-to-speech adapter."""

    async def synthesise(self, text: str, voice_id: str = "default") -> bytes:
        """Synthesise text to audio bytes."""
        raise NotImplementedError("TODO(mahimai): Phase 2 — implement Cartesia synthesis")
