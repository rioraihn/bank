.PHONY: help build run test clean deps fmt lint

# Variables
APP_NAME := wallet-service
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS := -ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)"

# Help
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Dependencies
deps: ## Install dependencies
	go mod download
	go mod tidy

# Build
build: deps ## Build the application
	CGO_ENABLED=0 go build -a -installsuffix cgo -o bin/$(APP_NAME) cmd/service/main.go $(LDFLAGS)

# Run
run: ## Run the application
	go run cmd/service/main.go

# Test
test: ## Run all tests
	go test -v ./...

test-coverage: ## Run tests with coverage
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Code Quality
fmt: ## Format code
	go fmt ./...

lint: ## Lint code
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Install it from https://golangci-lint.run/"; \
	fi

# Clean
clean: ## Clean build artifacts
	rm -rf bin/
	rm -f coverage.out coverage.html

# Development
dev-setup: ## Setup development environment
	@echo "Setting up development environment..."
	@if [ ! -f .env ]; then cp .env.example .env; echo "Created .env file from .env.example"; fi
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Docker
docker-build: ## Build Docker image
	docker build -t $(APP_NAME):$(VERSION) .

docker-run: ## Run Docker container
	docker run -p 8080:8080 $(APP_NAME):$(VERSION)

# Production
build-linux: ## Build for Linux
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o bin/$(APP_NAME)-linux-amd64 cmd/service/main.go $(LDFLAGS)

build-all: ## Build for multiple platforms
	@echo "Building for multiple platforms..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o bin/$(APP_NAME)-linux-amd64 cmd/service/main.go $(LDFLAGS)
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -a -installsuffix cgo -o bin/$(APP_NAME)-darwin-amd64 cmd/service/main.go $(LDFLAGS)
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -a -installsuffix cgo -o bin/$(APP_NAME)-windows-amd64.exe cmd/service/main.go $(LDFLAGS)

# Security
security-scan: ## Run security scan
	@if command -v gosec >/dev/null 2>&1; then \
		gosec ./...; \
	else \
		echo "gosec not installed. Install it with: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest"; \
	fi

# Database (for future use when adding persistent storage)
db-migrate-up: ## Run database migrations up
	@echo "Database migrations not yet implemented"

db-migrate-down: ## Run database migrations down
	@echo "Database migrations not yet implemented"