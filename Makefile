.PHONY: help build run stop clean test docker-build docker-up docker-down docker-logs deploy-prod backup

# Variables
BINARY_NAME=suproxy-api
DOCKER_IMAGE=suproxy/backend
VERSION?=latest

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the application binary
	@echo "Building $(BINARY_NAME)..."
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o bin/$(BINARY_NAME) ./cmd/api

run: ## Run the application locally
	@echo "Running $(BINARY_NAME)..."
	@go run ./cmd/api

stop: ## Stop the application (development)
	@echo "Stopping $(BINARY_NAME)..."
	@pkill -f $(BINARY_NAME) || true

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f $(BINARY_NAME)
	@rm -f cmd/api/$(BINARY_NAME)
	@rm -f cmd/api/*.exe

test: ## Run unit tests only
	@echo "Running unit tests..."
	@go test -v -race -short -coverprofile=coverage.out ./...

test-unit: ## Run unit tests only (alias)
	@$(MAKE) test

test-integration: ## Run integration tests
	@echo "Running integration tests..."
	@INTEGRATION_TEST=true go test -v -race -coverprofile=coverage-integration.out ./...

test-all: ## Run all tests (unit + integration)
	@echo "Running all tests..."
	@INTEGRATION_TEST=true go test -v -race -coverprofile=coverage-all.out ./...

test-coverage: test ## Run tests with coverage report
	@go tool cover -html=coverage.out

test-coverage-integration: test-integration ## Run integration tests with coverage
	@go tool cover -html=coverage-integration.out

test-coverage-all: test-all ## Run all tests with coverage
	@go tool cover -html=coverage-all.out

test-verbose: ## Run tests with verbose output
	@echo "Running tests (verbose)..."
	@go test -v -race ./...

test-package: ## Run tests for specific package (usage: make test-package PKG=./internal/domain/user)
	@echo "Running tests for $(PKG)..."
	@go test -v -race $(PKG)

test-clean: ## Clean test cache and coverage files
	@echo "Cleaning test cache..."
	@go clean -testcache
	@rm -f coverage*.out

test-db-setup: ## Setup test database
	@echo "Setting up test database..."
	@docker-compose exec postgres psql -U suproxy -c "DROP DATABASE IF EXISTS suproxy_test;"
	@docker-compose exec postgres psql -U suproxy -c "CREATE DATABASE suproxy_test;"
	@docker-compose exec postgres psql -U suproxy -c "CREATE USER suproxy_test WITH PASSWORD 'suproxy_test';" || true
	@docker-compose exec postgres psql -U suproxy -c "GRANT ALL PRIVILEGES ON DATABASE suproxy_test TO suproxy_test;"

test-db-teardown: ## Teardown test database
	@echo "Tearing down test database..."
	@docker-compose exec postgres psql -U suproxy -c "DROP DATABASE IF EXISTS suproxy_test;"

docker-build: ## Build Docker image
	@echo "Building Docker image $(DOCKER_IMAGE):$(VERSION)..."
	@docker build -t $(DOCKER_IMAGE):$(VERSION) .

docker-up: ## Start development environment
	@echo "Starting development environment..."
	@docker-compose up -d

docker-down: ## Stop development environment
	@echo "Stopping development environment..."
	@docker-compose down

docker-logs: ## View Docker logs
	@docker-compose logs -f

deploy-prod: ## Deploy to production
	@echo "Deploying to production..."
	@bash scripts/deploy.sh

backup: ## Backup database
	@echo "Creating database backup..."
	@bash scripts/backup.sh

lint: ## Run linter
	@echo "Running linter..."
	@golangci-lint run ./...

fmt: ## Format code
	@echo "Formatting code..."
	@go fmt ./...

vet: ## Run go vet
	@echo "Running go vet..."
	@go vet ./...

migrate-up: ## Run database migrations up
	@echo "Running migrations..."
	@go run ./cmd/migrate up

migrate-down: ## Run database migrations down
	@echo "Rolling back migrations..."
	@go run ./cmd/migrate down

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download

tidy: ## Tidy dependencies
	@echo "Tidying dependencies..."
	@go mod tidy

security-scan: ## Run security scan
	@echo "Running security scan..."
	@gosec ./...
