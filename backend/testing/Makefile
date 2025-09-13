# Testing Directory Makefile
# Provides convenient commands for running different test categories

.PHONY: help test test-unit test-integration test-performance test-email test-short clean setup

# Default target
help:
	@echo "Available targets:"
	@echo "  setup           - Initialize testing environment"
	@echo "  test            - Run all tests"
	@echo "  test-unit       - Run unit tests only"
	@echo "  test-integration - Run integration tests"
	@echo "  test-performance - Run performance tests"
	@echo "  test-email      - Run email system tests"
	@echo "  test-short      - Run tests in short mode (skip long-running tests)"
	@echo "  clean           - Clean test artifacts"

# Setup testing environment
setup:
	@echo "Setting up testing environment..."
	go mod tidy
	@echo "Testing environment ready"

# Run all tests
test: test-unit test-integration test-performance test-email

# Run unit tests (in parent directory)
test-unit:
	@echo "Running unit tests..."
	cd .. && go test ./internal/...

# Run integration tests
test-integration:
	@echo "Running integration tests..."
	@if [ -n "$$(find integration -name '*.go' -type f 2>/dev/null)" ]; then \
		for test in integration/*.go; do \
			if [ -f "$$test" ]; then \
				echo "Running $$test..."; \
				go run $$test; \
			fi \
		done \
	else \
		echo "No integration tests found (will be added in subsequent tasks)"; \
	fi

# Run performance tests
test-performance:
	@echo "Running performance tests..."
	@if [ -n "$$(find performance -name '*.go' -type f 2>/dev/null)" ]; then \
		for test in performance/*.go; do \
			if [ -f "$$test" ]; then \
				echo "Running $$test..."; \
				go run $$test; \
			fi \
		done \
	else \
		echo "No performance tests found (will be added in subsequent tasks)"; \
	fi

# Run email tests
test-email:
	@echo "Running email tests..."
	@if [ -n "$$(find email -name '*.go' -type f 2>/dev/null)" ]; then \
		for test in email/*.go; do \
			if [ -f "$$test" ]; then \
				echo "Running $$test..."; \
				go run $$test; \
			fi \
		done \
	else \
		echo "No email tests found (will be added in subsequent tasks)"; \
	fi

# Run tests in short mode
test-short:
	@echo "Running tests in short mode..."
	cd .. && go test -short ./...

# Clean test artifacts
clean:
	@echo "Cleaning test artifacts..."
	find . -name "*.prof" -delete
	find . -name "*.test" -delete
	find . -name "*.log" -delete
	@echo "Test artifacts cleaned"

# Verify test structure
verify:
	@echo "Verifying test directory structure..."
	@echo "Directories:"
	@ls -la
	@echo ""
	@echo "Shared utilities:"
	@ls -la shared/
	@echo ""
	@echo "Go module status:"
	@go mod verify