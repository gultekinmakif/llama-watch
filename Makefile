.PHONY: build build-server build-refresh run dev refresh test tidy clean help
.DEFAULT_GOAL := help

build: build-server build-refresh ## Build both server and refresh binaries into bin/

build-server: ## Build the server binary into bin/server
	go build -o bin/server ./cmd/server

build-refresh: ## Build the refresh binary into bin/refresh
	go build -o bin/refresh ./cmd/refresh

run: ## Run the server (sources .env if present)
	@if [ -f .env ]; then set -a; . ./.env; set +a; fi; go run ./cmd/server

dev: ## Run with hot reload via air
	@./scripts/dev.sh

refresh: ## Run the cron orchestrator end-to-end
	@./scripts/refresh.sh

test: ## Run all tests
	go test ./...

tidy: ## Tidy go.mod / go.sum
	go mod tidy

clean: ## Remove build artifacts
	rm -rf bin/

help: ## List targets
	@grep -E '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | awk 'BEGIN{FS=":.*?## "}{printf "  \033[36m%-14s\033[0m %s\n", $$1, $$2}'
