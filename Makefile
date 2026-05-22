.PHONY: build build-server build-sync-db run dev refresh refresh-soft test tidy clean help
.DEFAULT_GOAL := help

build: build-server build-sync-db ## Build all binaries into bin/

build-server: ## Build the server binary into bin/server
	go build -o bin/server ./cmd/server

build-sync-db: ## Build the sync-db binary into bin/sync-db
	go build -o bin/sync-db ./cmd/sync-db

run: ## Run the server (sources .env if present)
	@if [ -f .env ]; then set -a; . ./.env; set +a; fi; go run ./cmd/server

dev: ## Run with hot reload via air
	@./scripts/dev.sh

refresh: ## Run the cron orchestrator end-to-end
	@./scripts/refresh.sh

refresh-soft: ## Re-run extract + build + sync + web against existing var/ (no upstream fetch)
	@./scripts/refresh.sh --soft

test: ## Run all tests
	go test ./...

tidy: ## Tidy go.mod / go.sum
	go mod tidy

clean: ## Remove build artifacts
	rm -rf bin/

help: ## List targets
	@grep -E '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | awk 'BEGIN{FS=":.*?## "}{printf "  \033[36m%-14s\033[0m %s\n", $$1, $$2}'
