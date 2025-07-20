# Cloud Consulting Backend Makefile

.PHONY: help build run test test-coverage clean deps fmt lint docker-build docker-run migrate migrate-down

# Variables
BINARY_NAME=server
DOCKER_IMAGE=cloud-consulting-backend
DOCKER_TAG=latest

# Default target
help:
	@echo "Available targets:"
	@echo "  build         - Build the application"
	@echo "  run           - Run the application"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage"
	@echo "  clean         - Clean build artifacts"
	@echo "  deps          - Download dependencies"
	@echo "  fmt           - Format code"
	@echo "  lint          - Run linter"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-run    - Run Docker container"
	@echo "  migrate       - Run database migrations"
	@echo "  migrate-down  - Rollback migrations"

# Build the application
build:
	go build -o bin/server ./cmd/server

# Run the application
run:
	go run ./cmd/server

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Clean build artifacts
clean:
	rm -rf bin/
	go clean

# Download dependencies
deps:
	go mod download
	go mod tidy

# Format code
fmt:
	go fmt ./...

# Run linter (requires golangci-lint)
lint:
	golangci-lint run

# Build Docker image
docker-build:
	docker build -t cloud-consulting-backend .

# Run Docker container
docker-run:
	docker run -p 8080:8080 --env-file .env cloud-consulting-backend

# Development setup
dev-setup: deps
	@echo "Setting up development environment..."
	@if [ ! -f .env ]; then cp .env.example .env; echo "Created .env file from .env.example"; fi
	@echo "Development setup complete!"

# Run with live reload (requires air: go install github.com/cosmtrek/air@latest)
dev:
	air

# Database migrations (placeholder - will be implemented in later tasks)
migrate:
	@echo "Database migrations will be implemented in task 2.2"
	@echo "For now, you can use the init.sql script in scripts/"

migrate-down:
	@echo "Migration rollback will be implemented in task 2.2"

# Docker compose commands
docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f

# Quality checks
check: fmt lint test
	@echo "All quality checks passed!"

# Production build
build-prod:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o bin/$(BINARY_NAME) ./cmd/server