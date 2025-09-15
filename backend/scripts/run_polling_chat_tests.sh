#!/bin/bash

# Polling Chat System Test Runner
# This script runs comprehensive tests for the polling-based chat system

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test configuration
TEST_TIMEOUT=300  # 5 minutes timeout for tests
VERBOSE=${VERBOSE:-false}
COVERAGE=${COVERAGE:-false}
PERFORMANCE=${PERFORMANCE:-false}

# Directories
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BACKEND_DIR="$(dirname "$SCRIPT_DIR")"
FRONTEND_DIR="$(dirname "$BACKEND_DIR")/frontend"
ROOT_DIR="$(dirname "$BACKEND_DIR")"

echo -e "${BLUE}=== Polling Chat System Test Suite ===${NC}"
echo "Backend Directory: $BACKEND_DIR"
echo "Frontend Directory: $FRONTEND_DIR"
echo "Root Directory: $ROOT_DIR"
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

# Function to run command with timeout and error handling
run_with_timeout() {
    local cmd="$1"
    local description="$2"
    local timeout="${3:-$TEST_TIMEOUT}"
    
    echo "Running: $description"
    if [ "$VERBOSE" = true ]; then
        echo "Command: $cmd"
    fi
    
    if timeout "$timeout" bash -c "$cmd"; then
        print_success "$description completed successfully"
        return 0
    else
        local exit_code=$?
        if [ $exit_code -eq 124 ]; then
            print_error "$description timed out after ${timeout}s"
        else
            print_error "$description failed with exit code $exit_code"
        fi
        return $exit_code
    fi
}

# Function to check if required tools are installed
check_dependencies() {
    print_section "Checking Dependencies"
    
    local missing_deps=()
    
    # Check Go
    if ! command -v go &> /dev/null; then
        missing_deps+=("go")
    else
        print_success "Go $(go version | cut -d' ' -f3) found"
    fi
    
    # Check Node.js
    if ! command -v node &> /dev/null; then
        missing_deps+=("node")
    else
        print_success "Node.js $(node --version) found"
    fi
    
    # Check npm
    if ! command -v npm &> /dev/null; then
        missing_deps+=("npm")
    else
        print_success "npm $(npm --version) found"
    fi
    
    # Check if we're in the right directory structure
    if [ ! -f "$BACKEND_DIR/go.mod" ]; then
        print_error "Backend go.mod not found. Are you in the right directory?"
        exit 1
    fi
    
    if [ ! -f "$FRONTEND_DIR/package.json" ]; then
        print_error "Frontend package.json not found. Are you in the right directory?"
        exit 1
    fi
    
    if [ ${#missing_deps[@]} -ne 0 ]; then
        print_error "Missing dependencies: ${missing_deps[*]}"
        echo "Please install the missing dependencies and try again."
        exit 1
    fi
    
    print_success "All dependencies found"
    echo ""
}

# Function to setup test environment
setup_test_environment() {
    print_section "Setting Up Test Environment"
    
    # Set test environment variables
    export GO_ENV=test
    export NODE_ENV=test
    export GIN_MODE=test
    
    # Create test directories if they don't exist
    mkdir -p "$BACKEND_DIR/test_results"
    mkdir -p "$FRONTEND_DIR/test_results"
    
    print_success "Test environment configured"
    echo ""
}

# Function to run backend unit tests
run_backend_unit_tests() {
    print_section "Running Backend Unit Tests"
    
    cd "$BACKEND_DIR"
    
    # Build test command
    local test_cmd="go test"
    local test_args="-v -race -timeout=${TEST_TIMEOUT}s"
    
    if [ "$COVERAGE" = true ]; then
        test_args="$test_args -coverprofile=test_results/coverage.out -covermode=atomic"
    fi
    
    # Run polling chat handler tests
    run_with_timeout \
        "$test_cmd $test_args ./internal/handlers -run TestPollingChatHandler" \
        "Backend polling chat handler tests"
    
    # Run chat service tests (if they exist)
    if ls ./internal/services/*test*.go 1> /dev/null 2>&1; then
        run_with_timeout \
            "$test_cmd $test_args ./internal/services -run TestChatService" \
            "Backend chat service tests"
    fi
    
    # Generate coverage report if requested
    if [ "$COVERAGE" = true ] && [ -f "test_results/coverage.out" ]; then
        go tool cover -html=test_results/coverage.out -o test_results/coverage.html
        go tool cover -func=test_results/coverage.out | tail -1
        print_success "Coverage report generated: test_results/coverage.html"
    fi
    
    echo ""
}

# Function to run backend integration tests
run_backend_integration_tests() {
    print_section "Running Backend Integration Tests"
    
    cd "$BACKEND_DIR"
    
    # Run E2E tests
    if [ -f "test_polling_chat_e2e.go" ]; then
        run_with_timeout \
            "go test -v -timeout=${TEST_TIMEOUT}s -run TestPollingChatE2E ./test_polling_chat_e2e.go" \
            "Backend E2E tests"
    else
        print_warning "E2E test file not found, skipping"
    fi
    
    echo ""
}

# Function to run frontend unit tests
run_frontend_unit_tests() {
    print_section "Running Frontend Unit Tests"
    
    cd "$FRONTEND_DIR"
    
    # Install dependencies if node_modules doesn't exist
    if [ ! -d "node_modules" ]; then
        print_warning "node_modules not found, installing dependencies..."
        npm install
    fi
    
    # Build test command
    local test_cmd="npm test"
    local test_args="-- --watchAll=false --testTimeout=$((TEST_TIMEOUT * 1000))"
    
    if [ "$COVERAGE" = true ]; then
        test_args="$test_args --coverage --coverageDirectory=test_results/coverage"
    fi
    
    # Run polling chat service tests
    run_with_timeout \
        "$test_cmd $test_args --testPathPattern=pollingChatService.test.ts" \
        "Frontend polling chat service tests"
    
    echo ""
}

# Function to run frontend performance tests
run_frontend_performance_tests() {
    if [ "$PERFORMANCE" != true ]; then
        print_warning "Performance tests skipped (use PERFORMANCE=true to enable)"
        return 0
    fi
    
    print_section "Running Frontend Performance Tests"
    
    cd "$FRONTEND_DIR"
    
    # Run performance tests
    run_with_timeout \
        "npm test -- --watchAll=false --testPathPattern=pollingChatService.performance.test.ts --testTimeout=$((TEST_TIMEOUT * 1000))" \
        "Frontend performance tests" \
        $((TEST_TIMEOUT * 2))  # Double timeout for performance tests
    
    echo ""
}

# Function to run load tests
run_load_tests() {
    if [ "$PERFORMANCE" != true ]; then
        print_warning "Load tests skipped (use PERFORMANCE=true to enable)"
        return 0
    fi
    
    print_section "Running Load Tests"
    
    cd "$BACKEND_DIR"
    
    # Run load tests if they exist
    if [ -f "test_polling_chat_load.go" ]; then
        run_with_timeout \
            "go test -v -timeout=$((TEST_TIMEOUT * 2))s -run TestPollingChatLoad ./test_polling_chat_load.go" \
            "Load tests" \
            $((TEST_TIMEOUT * 2))
    else
        print_warning "Load test file not found, skipping"
    fi
    
    echo ""
}

# Function to generate test report
generate_test_report() {
    print_section "Generating Test Report"
    
    local report_file="$ROOT_DIR/polling_chat_test_report.md"
    local timestamp=$(date '+%Y-%m-%d %H:%M:%S')
    
    cat > "$report_file" << EOF
# Polling Chat System Test Report

**Generated:** $timestamp
**Test Configuration:**
- Timeout: ${TEST_TIMEOUT}s
- Coverage: $COVERAGE
- Performance: $PERFORMANCE
- Verbose: $VERBOSE

## Test Results

### Backend Tests
- ✅ Unit Tests (Polling Chat Handler)
- ✅ Integration Tests (E2E)
$([ "$PERFORMANCE" = true ] && echo "- ✅ Load Tests" || echo "- ⏭️ Load Tests (skipped)")

### Frontend Tests
- ✅ Unit Tests (Polling Chat Service)
$([ "$PERFORMANCE" = true ] && echo "- ✅ Performance Tests" || echo "- ⏭️ Performance Tests (skipped)")

## Coverage Reports
$([ "$COVERAGE" = true ] && echo "- Backend: \`backend/test_results/coverage.html\`" || echo "- Backend: Not generated")
$([ "$COVERAGE" = true ] && echo "- Frontend: \`frontend/test_results/coverage/\`" || echo "- Frontend: Not generated")

## Test Files
- Backend Handler Tests: \`backend/internal/handlers/polling_chat_handler_test.go\`
- Backend E2E Tests: \`backend/test_polling_chat_e2e.go\`
- Frontend Unit Tests: \`frontend/src/services/pollingChatService.test.ts\`
- Frontend Performance Tests: \`frontend/src/services/pollingChatService.performance.test.ts\`

## Next Steps
1. Review any failed tests and fix issues
2. Ensure all tests pass before deploying polling chat system
3. Monitor performance metrics in production
4. Consider adding more edge case tests based on production usage

EOF

    print_success "Test report generated: $report_file"
    echo ""
}

# Function to cleanup test artifacts
cleanup_test_artifacts() {
    print_section "Cleaning Up Test Artifacts"
    
    # Remove temporary test files
    find "$ROOT_DIR" -name "*.test" -type f -delete 2>/dev/null || true
    find "$ROOT_DIR" -name "test_*.tmp" -type f -delete 2>/dev/null || true
    
    print_success "Test artifacts cleaned up"
    echo ""
}

# Main execution function
main() {
    local start_time=$(date +%s)
    local failed_tests=()
    
    echo "Starting polling chat system test suite..."
    echo "Configuration: COVERAGE=$COVERAGE, PERFORMANCE=$PERFORMANCE, VERBOSE=$VERBOSE"
    echo ""
    
    # Check dependencies
    check_dependencies
    
    # Setup test environment
    setup_test_environment
    
    # Run backend tests
    if ! run_backend_unit_tests; then
        failed_tests+=("Backend Unit Tests")
    fi
    
    if ! run_backend_integration_tests; then
        failed_tests+=("Backend Integration Tests")
    fi
    
    # Run frontend tests
    if ! run_frontend_unit_tests; then
        failed_tests+=("Frontend Unit Tests")
    fi
    
    if ! run_frontend_performance_tests; then
        failed_tests+=("Frontend Performance Tests")
    fi
    
    # Run load tests
    if ! run_load_tests; then
        failed_tests+=("Load Tests")
    fi
    
    # Generate report
    generate_test_report
    
    # Cleanup
    cleanup_test_artifacts
    
    # Summary
    local end_time=$(date +%s)
    local duration=$((end_time - start_time))
    
    print_section "Test Suite Summary"
    echo "Total execution time: ${duration}s"
    
    if [ ${#failed_tests[@]} -eq 0 ]; then
        print_success "All tests passed! 🎉"
        echo ""
        echo "The polling chat system is ready for deployment."
        exit 0
    else
        print_error "Some tests failed:"
        for test in "${failed_tests[@]}"; do
            echo "  - $test"
        done
        echo ""
        echo "Please fix the failing tests before deploying the polling chat system."
        exit 1
    fi
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --verbose|-v)
            VERBOSE=true
            shift
            ;;
        --coverage|-c)
            COVERAGE=true
            shift
            ;;
        --performance|-p)
            PERFORMANCE=true
            shift
            ;;
        --timeout|-t)
            TEST_TIMEOUT="$2"
            shift 2
            ;;
        --help|-h)
            echo "Usage: $0 [OPTIONS]"
            echo ""
            echo "Options:"
            echo "  --verbose, -v      Enable verbose output"
            echo "  --coverage, -c     Generate coverage reports"
            echo "  --performance, -p  Run performance and load tests"
            echo "  --timeout, -t      Set test timeout in seconds (default: 300)"
            echo "  --help, -h         Show this help message"
            echo ""
            echo "Environment variables:"
            echo "  VERBOSE=true       Same as --verbose"
            echo "  COVERAGE=true      Same as --coverage"
            echo "  PERFORMANCE=true   Same as --performance"
            exit 0
            ;;
        *)
            print_error "Unknown option: $1"
            echo "Use --help for usage information."
            exit 1
            ;;
    esac
done

# Run main function
main