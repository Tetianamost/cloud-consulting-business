#!/bin/bash

# Test runner script for AI Consultant Live Chat
# This script runs all types of tests: unit, integration, and performance

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test configuration
COVERAGE_THRESHOLD=80
TIMEOUT=30m
PARALLEL_JOBS=4

# Directories
BACKEND_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
FRONTEND_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../../frontend" && pwd)"

echo -e "${BLUE}=== AI Consultant Live Chat Test Suite ===${NC}"
echo "Backend Directory: $BACKEND_DIR"
echo "Frontend Directory: $FRONTEND_DIR"
echo ""

# Function to print section headers
print_section() {
    echo -e "${BLUE}=== $1 ===${NC}"
}

# Function to print success messages
print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

# Function to print error messages
print_error() {
    echo -e "${RED}✗ $1${NC}"
}

# Function to print warning messages
print_warning() {
    echo -e "${YELLOW}⚠ $1${NC}"
}

# Function to run command with error handling
run_command() {
    local cmd="$1"
    local description="$2"
    
    echo "Running: $description"
    if eval "$cmd"; then
        print_success "$description completed"
        return 0
    else
        print_error "$description failed"
        return 1
    fi
}

# Check if Go is installed
check_go() {
    if ! command -v go &> /dev/null; then
        print_error "Go is not installed or not in PATH"
        exit 1
    fi
    
    GO_VERSION=$(go version | awk '{print $3}')
    print_success "Go version: $GO_VERSION"
}

# Check if Node.js is installed
check_node() {
    if ! command -v node &> /dev/null; then
        print_error "Node.js is not installed or not in PATH"
        exit 1
    fi
    
    NODE_VERSION=$(node --version)
    print_success "Node.js version: $NODE_VERSION"
}

# Install Go dependencies
install_go_deps() {
    print_section "Installing Go Dependencies"
    cd "$BACKEND_DIR"
    
    run_command "go mod download" "Download Go modules"
    run_command "go mod tidy" "Tidy Go modules"
}

# Install Node.js dependencies
install_node_deps() {
    print_section "Installing Node.js Dependencies"
    cd "$FRONTEND_DIR"
    
    if [ -f "package-lock.json" ]; then
        run_command "npm ci" "Install Node.js dependencies (ci)"
    else
        run_command "npm install" "Install Node.js dependencies"
    fi
}

# Run Go unit tests
run_go_unit_tests() {
    print_section "Running Go Unit Tests"
    cd "$BACKEND_DIR"
    
    # Create coverage directory
    mkdir -p coverage
    
    # Run tests with coverage
    run_command "go test -v -race -coverprofile=coverage/unit.out -covermode=atomic -timeout=$TIMEOUT ./internal/..." "Go unit tests"
    
    # Generate coverage report
    if [ -f "coverage/unit.out" ]; then
        run_command "go tool cover -html=coverage/unit.out -o coverage/unit.html" "Generate unit test coverage report"
        
        # Check coverage threshold
        COVERAGE=$(go tool cover -func=coverage/unit.out | grep total | awk '{print $3}' | sed 's/%//')
        echo "Unit test coverage: ${COVERAGE}%"
        
        if (( $(echo "$COVERAGE >= $COVERAGE_THRESHOLD" | bc -l) )); then
            print_success "Coverage threshold met: ${COVERAGE}% >= ${COVERAGE_THRESHOLD}%"
        else
            print_warning "Coverage below threshold: ${COVERAGE}% < ${COVERAGE_THRESHOLD}%"
        fi
    fi
}

# Run Go integration tests
run_go_integration_tests() {
    print_section "Running Go Integration Tests"
    cd "$BACKEND_DIR"
    
    # Run integration tests
    run_command "go test -v -race -tags=integration -timeout=$TIMEOUT ./test_*integration*.go" "Go integration tests"
}

# Run Go performance tests
run_go_performance_tests() {
    print_section "Running Go Performance Tests"
    cd "$BACKEND_DIR"
    
    # Run performance tests (not in short mode)
    run_command "go test -v -run=TestLoadTest -timeout=$TIMEOUT ./test_performance_load.go" "Go load tests"
    run_command "go test -v -run=TestStressTest -timeout=$TIMEOUT ./test_performance_load.go" "Go stress tests"
    
    # Run benchmarks
    run_command "go test -bench=. -benchmem -timeout=$TIMEOUT ./test_performance_load.go" "Go benchmarks"
}

# Run Node.js unit tests
run_node_unit_tests() {
    print_section "Running Node.js Unit Tests"
    cd "$FRONTEND_DIR"
    
    # Run Jest tests
    run_command "npm run test -- --coverage --watchAll=false" "Node.js unit tests"
    
    # Check if coverage meets threshold
    if [ -f "coverage/lcov-report/index.html" ]; then
        print_success "Frontend coverage report generated"
    fi
}

# Run Cypress E2E tests
run_cypress_tests() {
    print_section "Running Cypress E2E Tests"
    cd "$FRONTEND_DIR"
    
    # Check if Cypress is installed
    if ! npm list cypress &> /dev/null; then
        print_warning "Cypress not installed, skipping E2E tests"
        return 0
    fi
    
    # Start the application in background for E2E tests
    print_warning "Starting application for E2E tests..."
    npm run start &
    APP_PID=$!
    
    # Wait for application to start
    sleep 10
    
    # Run Cypress tests
    if run_command "npm run cypress:run" "Cypress E2E tests"; then
        E2E_SUCCESS=true
    else
        E2E_SUCCESS=false
    fi
    
    # Stop the application
    kill $APP_PID 2>/dev/null || true
    
    if [ "$E2E_SUCCESS" = true ]; then
        print_success "E2E tests completed successfully"
    else
        print_error "E2E tests failed"
        return 1
    fi
}

# Run linting
run_linting() {
    print_section "Running Code Linting"
    
    # Go linting
    cd "$BACKEND_DIR"
    if command -v golangci-lint &> /dev/null; then
        run_command "golangci-lint run ./..." "Go linting"
    else
        print_warning "golangci-lint not installed, skipping Go linting"
    fi
    
    # Node.js linting
    cd "$FRONTEND_DIR"
    if npm list eslint &> /dev/null; then
        run_command "npm run lint" "Node.js linting"
    else
        print_warning "ESLint not configured, skipping Node.js linting"
    fi
}

# Run security checks
run_security_checks() {
    print_section "Running Security Checks"
    
    # Go security check
    cd "$BACKEND_DIR"
    if command -v gosec &> /dev/null; then
        run_command "gosec ./..." "Go security check"
    else
        print_warning "gosec not installed, skipping Go security check"
    fi
    
    # Node.js security audit
    cd "$FRONTEND_DIR"
    run_command "npm audit --audit-level=moderate" "Node.js security audit"
}

# Generate test reports
generate_reports() {
    print_section "Generating Test Reports"
    
    cd "$BACKEND_DIR"
    
    # Combine coverage reports if they exist
    if [ -f "coverage/unit.out" ]; then
        # Create combined coverage report
        echo "mode: atomic" > coverage/combined.out
        tail -n +2 coverage/unit.out >> coverage/combined.out
        
        run_command "go tool cover -html=coverage/combined.out -o coverage/combined.html" "Generate combined coverage report"
        
        # Generate coverage summary
        go tool cover -func=coverage/combined.out > coverage/summary.txt
        print_success "Coverage reports generated in coverage/ directory"
    fi
    
    # Create test summary
    cat > test_summary.txt << EOF
AI Consultant Live Chat - Test Summary
======================================

Test Run Date: $(date)
Go Version: $(go version | awk '{print $3}')
Node Version: $(node --version)

Test Results:
- Unit Tests: $([ -f "coverage/unit.out" ] && echo "PASSED" || echo "SKIPPED")
- Integration Tests: $([ $? -eq 0 ] && echo "PASSED" || echo "FAILED")
- E2E Tests: $([ "$E2E_SUCCESS" = true ] && echo "PASSED" || echo "SKIPPED")
- Performance Tests: $([ $? -eq 0 ] && echo "PASSED" || echo "SKIPPED")

Coverage:
$([ -f "coverage/summary.txt" ] && cat coverage/summary.txt | tail -1 || echo "No coverage data available")

EOF
    
    print_success "Test summary generated: test_summary.txt"
}

# Main execution
main() {
    local run_unit=true
    local run_integration=true
    local run_e2e=true
    local run_performance=false
    local run_lint=true
    local run_security=true
    
    # Parse command line arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            --unit-only)
                run_integration=false
                run_e2e=false
                run_performance=false
                shift
                ;;
            --integration-only)
                run_unit=false
                run_e2e=false
                run_performance=false
                shift
                ;;
            --e2e-only)
                run_unit=false
                run_integration=false
                run_performance=false
                shift
                ;;
            --performance)
                run_performance=true
                shift
                ;;
            --no-lint)
                run_lint=false
                shift
                ;;
            --no-security)
                run_security=false
                shift
                ;;
            --help)
                echo "Usage: $0 [options]"
                echo "Options:"
                echo "  --unit-only      Run only unit tests"
                echo "  --integration-only Run only integration tests"
                echo "  --e2e-only       Run only E2E tests"
                echo "  --performance    Include performance tests"
                echo "  --no-lint        Skip linting"
                echo "  --no-security    Skip security checks"
                echo "  --help           Show this help message"
                exit 0
                ;;
            *)
                print_error "Unknown option: $1"
                exit 1
                ;;
        esac
    done
    
    # Check prerequisites
    check_go
    check_node
    
    # Install dependencies
    install_go_deps
    install_node_deps
    
    # Run linting if enabled
    if [ "$run_lint" = true ]; then
        run_linting || print_warning "Linting failed but continuing..."
    fi
    
    # Run security checks if enabled
    if [ "$run_security" = true ]; then
        run_security_checks || print_warning "Security checks failed but continuing..."
    fi
    
    # Run tests based on options
    local test_failed=false
    
    if [ "$run_unit" = true ]; then
        run_go_unit_tests || test_failed=true
        run_node_unit_tests || test_failed=true
    fi
    
    if [ "$run_integration" = true ]; then
        run_go_integration_tests || test_failed=true
    fi
    
    if [ "$run_e2e" = true ]; then
        run_cypress_tests || test_failed=true
    fi
    
    if [ "$run_performance" = true ]; then
        run_go_performance_tests || test_failed=true
    fi
    
    # Generate reports
    generate_reports
    
    # Final summary
    print_section "Test Suite Summary"
    
    if [ "$test_failed" = true ]; then
        print_error "Some tests failed. Check the output above for details."
        exit 1
    else
        print_success "All tests passed successfully!"
        echo ""
        echo "Reports generated:"
        echo "  - Backend coverage: coverage/combined.html"
        echo "  - Frontend coverage: coverage/lcov-report/index.html"
        echo "  - Test summary: test_summary.txt"
        exit 0
    fi
}

# Run main function with all arguments
main "$@"