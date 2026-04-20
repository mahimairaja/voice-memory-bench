.PHONY: help install lint typecheck test test-unit test-integration format clean

PYTHON := python3
UV := uv

help:  ## Show this help
@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

install:  ## Install all dependencies including dev extras
$(UV) sync --extra dev

lint:  ## Run ruff linter
$(UV) run ruff check voice_memory_bench/ tests/ voice/

typecheck:  ## Run mypy type checker
$(UV) run mypy voice_memory_bench/

format:  ## Auto-format code with ruff
$(UV) run ruff format voice_memory_bench/ tests/ voice/
$(UV) run ruff check --fix voice_memory_bench/ tests/ voice/

test-unit:  ## Run unit tests (no external services)
$(UV) run pytest -m unit --tb=short -q

test-integration:  ## Run integration tests (requires live services)
$(UV) run pytest -m integration --tb=short

test:  ## Run all tests
$(UV) run pytest --tb=short

clean:  ## Remove build artifacts and caches
rm -rf dist/ build/ *.egg-info/ .pytest_cache/ .mypy_cache/ .ruff_cache/
find . -type d -name __pycache__ -exec rm -rf {} + 2>/dev/null || true
find . -type f -name "*.pyc" -delete

datasets-download-all:  ## Download all benchmark datasets
$(UV) run vmb datasets download locomo
$(UV) run vmb datasets download longmemeval

providers-up-mem0:  ## Start Mem0 Postgres
docker compose -f providers/mem0/docker-compose.yml up -d

providers-up-memori:  ## Start Memori Postgres
docker compose -f providers/memori/docker-compose.yml up -d

providers-up-graphiti:  ## Start Graphiti FalkorDB
docker compose -f providers/graphiti/docker-compose.yml up -d

providers-up-all:  ## Start all provider infrastructure
$(MAKE) providers-up-mem0 providers-up-memori providers-up-graphiti

providers-down:  ## Stop all provider infrastructure
docker compose -f providers/mem0/docker-compose.yml down
docker compose -f providers/memori/docker-compose.yml down
docker compose -f providers/graphiti/docker-compose.yml down
