"""FastAPI app implementing the vbench sidecar contract for Mem0."""

from __future__ import annotations

import json
import os
from typing import Any

from fastapi import FastAPI, HTTPException
from fastapi.responses import JSONResponse
from pydantic import BaseModel, Field

from .adapter import CapabilityNotSupported, Mem0Adapter


class AddMessageBody(BaseModel):
    user_id: str
    session_id: str
    role: str
    content: str
    metadata: dict[str, Any] | None = None


class AddFactBody(BaseModel):
    user_id: str
    session_id: str
    fact: str
    metadata: dict[str, Any] | None = None


class SearchBody(BaseModel):
    user_id: str
    session_id: str
    query: str
    mode: str = Field(default="semantic")
    top_k: int = Field(default=10, ge=1, le=100)
    filters: dict[str, Any] | None = None


class EnumerateBody(BaseModel):
    user_id: str
    session_id: str | None = None


class ResetBody(BaseModel):
    user_id: str
    session_id: str | None = None


def build_app() -> FastAPI:
    raw_config = os.environ.get("VBENCH_PROVIDER_CONFIG", "{}")
    try:
        provider_config = json.loads(raw_config)
    except json.JSONDecodeError as exc:
        raise RuntimeError(f"invalid VBENCH_PROVIDER_CONFIG: {exc}") from exc

    adapter = Mem0Adapter(provider_config)
    app = FastAPI(title="vbench-mem0", version="0.1.0")

    @app.get("/health")
    def health() -> dict[str, str]:
        adapter.health()
        return {"status": "ok"}

    @app.get("/capabilities")
    def capabilities() -> dict[str, Any]:
        return adapter.capabilities()

    @app.post("/add_message")
    def add_message(body: AddMessageBody) -> dict[str, Any]:
        return adapter.add_message(
            user_id=body.user_id,
            session_id=body.session_id,
            role=body.role,
            content=body.content,
            metadata=body.metadata,
        )

    @app.post("/add_fact")
    def add_fact(body: AddFactBody) -> dict[str, Any]:
        return adapter.add_fact(
            user_id=body.user_id,
            session_id=body.session_id,
            fact=body.fact,
            metadata=body.metadata,
        )

    @app.post("/search")
    def search(body: SearchBody) -> Any:
        try:
            return adapter.search(
                user_id=body.user_id,
                session_id=body.session_id,
                query=body.query,
                mode=body.mode,
                top_k=body.top_k,
                filters=body.filters,
            )
        except CapabilityNotSupported as exc:
            return JSONResponse(
                status_code=422,
                content={
                    "error_type": "capability_not_supported",
                    "provider": "mem0",
                    "capability": exc.capability,
                    "message": exc.reason,
                },
            )

    @app.post("/enumerate")
    def enumerate_(body: EnumerateBody) -> dict[str, Any]:
        items = adapter.enumerate(user_id=body.user_id, session_id=body.session_id)
        return {"items": items}

    @app.post("/reset")
    def reset(body: ResetBody) -> dict[str, str]:
        adapter.reset(user_id=body.user_id, session_id=body.session_id)
        return {"status": "ok"}

    @app.exception_handler(Exception)
    async def generic_error_handler(_req, exc: Exception) -> JSONResponse:  # type: ignore[override]
        if isinstance(exc, HTTPException):
            raise exc
        return JSONResponse(
            status_code=500,
            content={
                "error_type": "internal_error",
                "provider": "mem0",
                "message": str(exc),
            },
        )

    return app
