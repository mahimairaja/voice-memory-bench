.PHONY: help build vet test test-v tidy clean \
	sidecar-sync sidecar-lint sidecar-test \
	providers-up-mem0 providers-down-mem0

GO := go
UV := uv
SIDECAR_DIR := sidecars/mem0

help:  ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-24s\033[0m %s\n", $$1, $$2}'

build:  ## Compile the vbench CLI to ./vbench
	$(GO) build -o vbench ./cmd/vbench

vet:  ## Run go vet across all packages
	$(GO) vet ./...

test:  ## Run Go unit tests
	$(GO) test ./...

test-v:  ## Run Go unit tests, verbose
	$(GO) test -v ./...

tidy:  ## go mod tidy
	$(GO) mod tidy

clean:  ## Remove build outputs
	rm -f vbench coverage.txt coverage.html
	rm -rf dist/

sidecar-sync:  ## uv sync the Mem0 sidecar environment
	cd $(SIDECAR_DIR) && $(UV) sync

sidecar-lint:  ## Lint the Mem0 sidecar (ruff fetched on demand via uvx)
	cd $(SIDECAR_DIR) && $(UV) tool run ruff check vbench_mem0

sidecar-test:  ## Run sidecar tests (pytest exit 5 = no tests collected, treated as pass)
	@cd $(SIDECAR_DIR) && $(UV) run pytest -q; rc=$$?; \
		if [ $$rc -eq 0 ] || [ $$rc -eq 5 ]; then exit 0; else exit $$rc; fi

providers-up-mem0:  ## Start Mem0 Postgres (pgvector)
	docker compose -f providers/mem0/docker-compose.yml up -d

providers-down-mem0:  ## Stop Mem0 Postgres
	docker compose -f providers/mem0/docker-compose.yml down
