.PHONY: help run build test test-coverage clean docker-up docker-down migrate lint swagger

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

run: ## Run the application
	@go run cmd/api/main.go

build: ## Build the application
	@echo "Building..."
	@go build -o bin/api cmd/api/main.go

test: ## Run tests
	@echo "Running tests..."
	@go test -v -race ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html

docker-up: ## Start docker containers
	@echo "Starting docker containers..."
	@docker-compose up -d

docker-down: ## Stop docker containers
	@echo "Stopping docker containers..."
	@docker-compose down

docker-logs: ## View docker logs
	@docker-compose logs -f

lint: ## Run linter
	@echo "Running linter..."
	@golangci-lint run ./...

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download
	@go mod verify

tidy: ## Tidy dependencies
	@echo "Tidying dependencies..."
	@go mod tidy

install-tools: ## Install development tools
	@echo "Installing tools..."
	@go install github.com/air-verse/air@latest
	@go install github.com/swaggo/swag/cmd/swag@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "✓ Tools installed to $(shell go env GOPATH)/bin"
	@echo "💡 Add $(shell go env GOPATH)/bin to your PATH or use 'make dev' to run with hot reload"

dev: ## Run with hot reload (requires air)
	@echo "Starting development server with hot reload..."
	@air

swagger: ## Generate swagger documentation
	@echo "Generating swagger documentation..."
	@$(shell go env GOPATH)/bin/swag init -g cmd/api/main.go
	@echo "✓ Swagger documentation generated in docs/"
