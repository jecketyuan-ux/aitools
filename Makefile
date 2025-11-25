.PHONY: help build run test clean docker-build docker-up docker-down migrate-up migrate-down

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the application
	@echo "Building playedu..."
	@go build -o playedu cmd/api/main.go
	@echo "Build complete!"

run: ## Run the application
	@echo "Running playedu..."
	@go run cmd/api/main.go

test: ## Run tests
	@echo "Running tests..."
	@go test -v ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

lint: ## Run linter
	@echo "Running linter..."
	@golangci-lint run

fmt: ## Format code
	@echo "Formatting code..."
	@go fmt ./...
	@gofmt -s -w .

clean: ## Clean build files
	@echo "Cleaning..."
	@rm -f playedu
	@rm -f coverage.out coverage.html
	@echo "Clean complete!"

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	@docker-compose build

docker-up: ## Start Docker containers
	@echo "Starting Docker containers..."
	@docker-compose up -d
	@echo "Containers started!"

docker-down: ## Stop Docker containers
	@echo "Stopping Docker containers..."
	@docker-compose down
	@echo "Containers stopped!"

docker-logs: ## View Docker logs
	@docker-compose logs -f api

migrate-up: ## Run database migrations (up)
	@echo "Running database migrations..."
	@mysql -h 127.0.0.1 -u root -p < migrations/000001_init_schema.up.sql
	@echo "Migrations complete!"

migrate-down: ## Rollback database migrations
	@echo "Rolling back database migrations..."
	@mysql -h 127.0.0.1 -u root -p < migrations/000001_init_schema.down.sql
	@echo "Rollback complete!"

deps: ## Install dependencies
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy
	@echo "Dependencies installed!"

install-tools: ## Install development tools
	@echo "Installing development tools..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Tools installed!"
