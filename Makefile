.PHONY: help build run dev docker-build docker-up docker-down docker-logs clean test lint

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: ## Build the application
	@echo "Building application..."
	@go build -o bin/suproxy-api ./cmd/api

run: ## Run the application
	@echo "Running application..."
	@go run ./cmd/api

dev: ## Run the application with hot reload (requires air)
	@echo "Starting development server with hot reload..."
	@air

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	@docker-compose build

docker-up: ## Start Docker containers
	@echo "Starting Docker containers..."
	@docker-compose up -d
	@echo "Waiting for services to be ready..."
	@timeout /t 5 /nobreak > nul
	@echo "Services are up!"
	@echo "API: http://localhost:8080"
	@echo "Health check: http://localhost:8080/health"

docker-down: ## Stop Docker containers
	@echo "Stopping Docker containers..."
	@docker-compose down

docker-logs: ## Show Docker logs
	@docker-compose logs -f

docker-restart: docker-down docker-up ## Restart Docker containers

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@if exist bin rmdir /s /q bin
	@if exist tmp rmdir /s /q tmp
	@go clean

test: ## Run tests
	@echo "Running tests..."
	@go test -v -race -coverprofile=coverage.out ./...

lint: ## Run linter
	@echo "Running linter..."
	@golangci-lint run

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

.DEFAULT_GOAL := help
