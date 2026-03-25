.PHONY: dev build docker test clean frontend backend migrate-up migrate-down migrate-create sqlc help

# Default target
help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# ---------------------------------------------------------------------------
# Development
# ---------------------------------------------------------------------------

dev: ## Start backend + frontend dev servers concurrently
	@echo "Starting ALMAZ dev servers..."
	@trap 'kill 0' EXIT; \
		(cd frontend && npm run dev) & \
		(cd backend && go run ./cmd/almaz) & \
		wait

# ---------------------------------------------------------------------------
# Build
# ---------------------------------------------------------------------------

frontend: ## Build the frontend
	cd frontend && npm run build

backend: frontend ## Build the Go binary with embedded frontend
	rm -rf backend/cmd/almaz/dist
	cp -r frontend/dist backend/cmd/almaz/dist
	cd backend && CGO_ENABLED=0 go build -o ../bin/almaz ./cmd/almaz
	rm -rf backend/cmd/almaz/dist

build: backend ## Build the full application (alias for backend)

# ---------------------------------------------------------------------------
# Docker
# ---------------------------------------------------------------------------

docker: ## Build the Docker image
	docker build -t almaz:latest .

# ---------------------------------------------------------------------------
# Database
# ---------------------------------------------------------------------------

migrate-up: ## Run all pending database migrations
	cd backend && go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate \
		-path ../migrations -database "$${DATABASE_URL}" up

migrate-down: ## Roll back the latest migration
	cd backend && go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate \
		-path ../migrations -database "$${DATABASE_URL}" down 1

migrate-create: ## Create a new migration (usage: make migrate-create NAME=xxx)
	@if [ -z "$(NAME)" ]; then echo "Usage: make migrate-create NAME=xxx"; exit 1; fi
	mkdir -p migrations
	@NUM=$$(printf "%06d" $$(($$(ls migrations/*.up.sql 2>/dev/null | wc -l) + 1))); \
	touch "migrations/$${NUM}_$(NAME).up.sql" "migrations/$${NUM}_$(NAME).down.sql"; \
	echo "Created migrations/$${NUM}_$(NAME).up.sql and .down.sql"

# ---------------------------------------------------------------------------
# Code generation
# ---------------------------------------------------------------------------

sqlc: ## Regenerate sqlc code
	cd backend && go run github.com/sqlc-dev/sqlc/cmd/sqlc generate

# ---------------------------------------------------------------------------
# Testing
# ---------------------------------------------------------------------------

test: ## Run all Go tests
	cd backend && go test ./...

# ---------------------------------------------------------------------------
# Cleanup
# ---------------------------------------------------------------------------

clean: ## Remove build artifacts
	rm -rf bin/
	rm -rf frontend/dist
	rm -rf frontend/node_modules/.vite/
	rm -rf backend/cmd/almaz/dist
