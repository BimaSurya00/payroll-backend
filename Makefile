.PHONY: run build test lint clean docker-up docker-down migrate seed help

# Application name
APP_NAME=go-fiber-boilerplate

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=$(APP_NAME)

# Colors for output
GREEN=\033[0;32m
NC=\033[0m # No Color

## help: Display this help message
help:
	@echo "Available commands:"
	@echo "  make run          - Run the application"
	@echo "  make build        - Build the application"
	@echo "  make test         - Run tests"
	@echo "  make test-cover   - Run tests with coverage"
	@echo "  make lint         - Run linter"
	@echo "  make clean        - Clean build files"
	@echo "  make docker-up    - Start Docker containers"
	@echo "  make docker-down  - Stop Docker containers"
	@echo "  make migrate      - Run database migrations"
	@echo "  make seed         - Run database seeders"
	@echo "  make tidy         - Tidy go modules"

## run: Run the application
run:
	@echo "$(GREEN)Running application...$(NC)"
	$(GOCMD) run main.go

## build: Build the application
build:
	@echo "$(GREEN)Building application...$(NC)"
	$(GOBUILD) -o bin/$(BINARY_NAME) -v main.go

## test: Run tests
test:
	@echo "$(GREEN)Running tests...$(NC)"
	$(GOTEST) -v ./...

## test-cover: Run tests with coverage
test-cover:
	@echo "$(GREEN)Running tests with coverage...$(NC)"
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)Coverage report generated: coverage.html$(NC)"

## lint: Run linter (requires golangci-lint)
lint:
	@echo "$(GREEN)Running linter...$(NC)"
	golangci-lint run ./...

## clean: Clean build files
clean:
	@echo "$(GREEN)Cleaning...$(NC)"
	$(GOCLEAN)
	rm -rf bin/
	rm -f coverage.out coverage.html

## docker-up: Start Docker containers
docker-up:
	@echo "$(GREEN)Starting Docker containers...$(NC)"
	docker compose up -d

## docker-down: Stop Docker containers
docker-down:
	@echo "$(GREEN)Stopping Docker containers...$(NC)"
	docker compose down

## docker-rebuild: Rebuild and restart Docker containers
docker-rebuild:
	@echo "$(GREEN)Rebuilding Docker containers...$(NC)"
	docker compose up -d --build

## tidy: Tidy go modules
tidy:
	@echo "$(GREEN)Tidying go modules...$(NC)"
	$(GOMOD) tidy

## deps: Download dependencies
deps:
	@echo "$(GREEN)Downloading dependencies...$(NC)"
	$(GOMOD) download

## migrate-up: Run all pending database migrations
migrate-up:
	@echo "$(GREEN)Running migrations (up)...$(NC)"
	@if [ -z "$$DATABASE_URL" ]; then \
		DATABASE_URL="postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DATABASE)?sslmode=disable"; \
	fi; \
	for f in database/migrations/*.up.sql; do \
		echo "Applying $$f ..."; \
		PGPASSWORD=$$(echo "$$DATABASE_URL" | sed -n 's/.*:\([^@]*\)@.*/\1/p') \
		psql -h $$(echo "$$DATABASE_URL" | sed -n 's/.*@\([^:]*\):.*/\1/p') \
			-p $$(echo "$$DATABASE_URL" | sed -n 's/.*@\([^/]*\)\/.*/\1/p' | cut -d: -f2) \
			-U $$(echo "$$DATABASE_URL" | sed -n 's/.*://;s/\/.*//p' | cut -d@ -f1) \
			-d $$(echo "$$DATABASE_URL" | sed -n 's/.*\/\([^?]*\).*/\1/p') \
			-f "$$f"; \
	done; \
	echo "$(GREEN)All migrations applied!$(NC)"

## migrate: Run database migrations (alias for migrate-up)
migrate: migrate-up

## migrate-docker: Run migrations inside Docker container
migrate-docker:
	@echo "$(GREEN)Running migrations via Docker...$(NC)"
	docker run --rm --network hris-network \
		-e DATABASE_URL="postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@postgres:5432/$(POSTGRES_DATABASE)?sslmode=disable" \
		-v $(CURDIR)/database/migrations:/migrations \
		golang:1.24-bookworm \
		go run /migrations/migrate.go

## seed: Run database seeders
seed:
	@echo "$(GREEN)Running seeders...$(NC)"
	$(GOCMD) run cmd/seeder/main.go

## dev: Run with hot reload (requires air)
dev:
	@echo "$(GREEN)Running with hot reload...$(NC)"
	air