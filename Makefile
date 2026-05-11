.PHONY: run dev build test lint tidy docs docker-build docker-up docker-down clean help
.DEFAULT_GOAL := help

BIN := bin/server
PKG := ./cmd/server

run: ## Run the server with env from configs/.env.example
	@set -a; . configs/.env.example; set +a; go run $(PKG)

dev: ## Run with hot reload via air
	@./scripts/dev.sh

build: ## Build the server binary into bin/
	go build -o $(BIN) $(PKG)

test: ## Run all tests
	go test ./...

lint: ## Run golangci-lint
	golangci-lint run

tidy: ## Tidy go.mod / go.sum
	go mod tidy

docs: ## Render api/openapi.yaml to docs/index.html (Redoc) for GitHub Pages
	npx --yes @redocly/cli build-docs api/openapi.yaml -o docs/index.html

docker-build: ## Build the Docker image (tag: go-http-server:dev)
	docker build -t go-http-server:dev .

docker-up: ## Start app + postgres via docker compose
	docker compose up --build

docker-down: ## Stop and remove docker compose services
	docker compose down

clean: ## Remove build artifacts
	rm -rf bin/ tmp/

help: ## List targets
	@grep -E '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | awk 'BEGIN{FS=":.*?## "}{printf "  \033[36m%-14s\033[0m %s\n", $$1, $$2}'
