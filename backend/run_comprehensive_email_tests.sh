#!/bin/bash

# Comprehensive Email Event Tracking System Test Runner
# This script runs all comprehensive tests for the email event tracking system

set -e

echo "=== Comprehensive Email Event Tracking System Tests ==="
echo "Starting test execution at $(date)"
echo

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test configuration
TEST_DATABASE_URL=${TEST_DATABASE_URL:-"postgres://test:test@localhost/test_email_events?sslmode=disable"}
VERBOSE=${VERBOSE:-false}
COVERAGE=${COVERAGE:-true}
TIMEOUT=${TIMEOUT:-300s}

# Function to print colored output
print_status() {
    local color=$1
    local message=$2
    echo -e "${color}${message}${NC}"
}

# Function to run a test with error handling
run_test() {
    local test_name=$1
    local test_file=$2
    local description=$3
    
    print_status $BLUE "Running: $test_name"
    print_status $YELLOW "Description: $description"
    echo
    
    if [ "$VERBOSE" = "true" ]; then
        go test -v -timeout=$TIMEOUT "$test_file"
    else
        go test -timeout=$TIMEOUT "$test_file"
    fi
    
    if [ $? -eq 0 ]; then
        print_status $GREEN "✓ $test_name PASSED"
    else
        print_status $RED "✗ $test_name FAILED"
        return 1
    fi
    echo
}

# Function to run tests with coverage
run_test_with_coverage() {
    local test_name=$1
    local test_file=$2
    local description=$3
    local coverage_file="coverage_${test_name}.out"
    
    print_status $BLUE "Running: $test_name (with coverage)"
    print_status $YELLOW "Description: $description"
    echo
    
    if [ "$VERBOSE" = "true" ]; then
        go test -v -timeout=$TIMEOUT -coverprofile="$coverage_file" "$test_file"
    else
        go test -timeout=$TIMEOUT -coverprofile="$coverage_file" "$test_file"
    fi
    
    if [ $? -eq 0 ]; then
        print_status $GREEN "✓ $test_name PASSED"
        if [ -f "$coverage_file" ]; then
            coverage=$(go tool cover -func="$coverage_file" | grep total | awk '{print $3}')
            print_status $BLUE "  Coverage: $coverage"
        fi
    else
        print_status $RED "✗ $test_name FAILED"
        return 1
    fi
    echo
}

# Function to check prerequisites
check_prerequisites() {
    print_status $BLUE "Checking prerequisites..."
    
    # Check if Go is installed
    if ! command -v go &> /dev/null; then
        print_status $RED "Go is not installed or not in PATH"
        exit 1
    fi
    
    # Check Go version
    go_version=$(go version | awk '{print $3}' | sed 's/go//')
    print_status $GREEN "Go version: $go_version"
    
    # Check if required Go modules are available
    if [ ! -f "go.mod" ]; then
        print_status $RED "go.mod not found. Please run from the backend directory."
        exit 1
    fi
    
    # Download dependencies
    print_status $BLUE "Downloading dependencies..."
    go mod download
    
    # Check if PostgreSQL is available for integration tests
    if command -v psql &> /dev/null; then
        print_status $GREEN "PostgreSQL client available for integration tests"
    else
        print_status $YELLOW "PostgreSQL client not available - some integration tests may be skipped"
    fi
    
    echo
}

# Function to setup test environment
setup_test_environment() {
    print_status $BLUE "Setting up test environment..."
    
    # Set test environment variables
    export TEST_DATABASE_URL="$TEST_DATABASE_URL"
    export SES_SENDER_EMAIL="info@cloudpartner.pro"
    export SES_REPLY_TO_EMAIL="info@cloudpartner.pro"
    export LOG_LEVEL="error"
    
    # Create test results directory
    mkdir -p test_results
    
    print_status $GREEN "Test environment configured"
    echo
}

# Function to cleanup test environment
cleanup_test_environment() {
    print_status $BLUE "Cleaning up test environment..."
    
    # Remove coverage files
    rm -f coverage_*.out
    
    # Remove temporary test files
    rm -f test_*.tmp
    
    print_status $GREEN "Cleanup completed"
    echo
}

# Function to generate combined coverage report
generate_coverage_report() {
    if [ "$COVERAGE" = "true" ]; then
        print_status $BLUE "Generating combined coverage report..."
        
        # Combine all coverage files
        echo "mode: set" > combined_coverage.out
        for coverage_file in coverage_*.out; do
            if [ -f "$coverage_file" ]; then
                tail -n +2 "$coverage_file" >> combined_coverage.out
            fi
        done
        
        if [ -f "combined_coverage.out" ]; then
            # Generate HTML coverage report
            go tool cover -html=combined_coverage.out -o test_results/coverage.html
            
            # Calculate total coverage
            total_coverage=$(go tool cover -func=combined_coverage.out | grep total | awk '{print $3}')
            print_status $GREEN "Total coverage: $total_coverage"
            print_status $BLUE "HTML coverage report: test_results/coverage.html"
        fi
        
        echo
    fi
}

# Function to run all tests
run_all_tests() {
    local failed_tests=0
    local total_tests=0
    
    print_status $BLUE "Starting comprehensive test suite..."
    echo
    
    # Test 1: Email Event Repository Tests
    total_tests=$((total_tests + 1))
    if [ "$COVERAGE" = "true" ]; then
        run_test_with_coverage "email_event_repository" \
            "./test_email_event_repository_comprehensive.go" \
            "Unit tests for email event repository operations (CRUD, metrics, filtering)" || failed_tests=$((failed_tests + 1))
    else
        run_test "email_event_repository" \
            "./test_email_event_repository_comprehensive.go" \
            "Unit tests for email event repository operations (CRUD, metrics, filtering)" || failed_tests=$((failed_tests + 1))
    fi
    
    # Test 2: Email Event Recorder Tests
    total_tests=$((total_tests + 1))
    if [ "$COVERAGE" = "true" ]; then
        run_test_with_coverage "email_event_recorder" \
            "./test_email_event_recorder_comprehensive.go" \
            "Unit tests for email event recorder service (async recording, retry logic, health checks)" || failed_tests=$((failed_tests + 1))
    else
        run_test "email_event_recorder" \
            "./test_email_event_recorder_comprehensive.go" \
            "Unit tests for email event recorder service (async recording, retry logic, health checks)" || failed_tests=$((failed_tests + 1))
    fi
    
    # Test 3: Email Metrics Service Tests
    total_tests=$((total_tests + 1))
    if [ "$COVERAGE" = "true" ]; then
        run_test_with_coverage "email_metrics_service" \
            "./test_email_metrics_service_comprehensive.go" \
            "Unit tests for email metrics service (calculations, filtering, validation)" || failed_tests=$((failed_tests + 1))
    else
        run_test "email_metrics_service" \
            "./test_email_metrics_service_comprehensive.go" \
            "Unit tests for email metrics service (calculations, filtering, validation)" || failed_tests=$((failed_tests + 1))
    fi
    
    # Test 4: Admin Handler Email Tests
    total_tests=$((total_tests + 1))
    if [ "$COVERAGE" = "true" ]; then
        run_test_with_coverage "admin_handler_email" \
            "./test_admin_handler_email_comprehensive.go" \
            "Unit tests for admin handler email endpoints (metrics, status, history)" || failed_tests=$((failed_tests + 1))
    else
        run_test "admin_handler_email" \
            "./test_admin_handler_email_comprehensive.go" \
            "Unit tests for admin handler email endpoints (metrics, status, history)" || failed_tests=$((failed_tests + 1))
    fi
    
    # Test 5: Email Service Integration Tests
    total_tests=$((total_tests + 1))
    if [ "$COVERAGE" = "true" ]; then
        run_test_with_coverage "email_service_integration" \
            "./test_email_service_integration_comprehensive.go" \
            "Integration tests for email service with event recording" || failed_tests=$((failed_tests + 1))
    else
        run_test "email_service_integration" \
            "./test_email_service_integration_comprehensive.go" \
            "Integration tests for email service with event recording" || failed_tests=$((failed_tests + 1))
    fi
    
    # Generate coverage report
    if [ "$COVERAGE" = "true" ]; then
        generate_coverage_report
    fi
    
    # Print summary
    print_status $BLUE "=== Test Summary ==="
    print_status $BLUE "Total tests: $total_tests"
    print_status $GREEN "Passed: $((total_tests - failed_tests))"
    if [ $failed_tests -gt 0 ]; then
        print_status $RED "Failed: $failed_tests"
    else
        print_status $GREEN "Failed: $failed_tests"
    fi
    echo
    
    if [ $failed_tests -eq 0 ]; then
        print_status $GREEN "🎉 All tests passed!"
        return 0
    else
        print_status $RED "❌ Some tests failed!"
        return 1
    fi
}

# Function to run specific test
run_specific_test() {
    local test_name=$1
    
    case $test_name in
        "repository"|"repo")
            run_test "email_event_repository" \
                "./test_email_event_repository_comprehensive.go" \
                "Unit tests for email event repository operations"
            ;;
        "recorder"|"record")
            run_test "email_event_recorder" \
                "./test_email_event_recorder_comprehensive.go" \
                "Unit tests for email event recorder service"
            ;;
        "metrics"|"metric")
            run_test "email_metrics_service" \
                "./test_email_metrics_service_comprehensive.go" \
                "Unit tests for email metrics service"
            ;;
        "handler"|"admin")
            run_test "admin_handler_email" \
                "./test_admin_handler_email_comprehensive.go" \
                "Unit tests for admin handler email endpoints"
            ;;
        "integration"|"int")
            run_test "email_service_integration" \
                "./test_email_service_integration_comprehensive.go" \
                "Integration tests for email service"
            ;;
        *)
            print_status $RED "Unknown test: $test_name"
            print_status $BLUE "Available tests: repository, recorder, metrics, handler, integration"
            exit 1
            ;;
    esac
}

# Function to show help
show_help() {
    echo "Comprehensive Email Event Tracking System Test Runner"
    echo
    echo "Usage: $0 [OPTIONS] [TEST_NAME]"
    echo
    echo "Options:"
    echo "  -h, --help              Show this help message"
    echo "  -v, --verbose           Enable verbose output"
    echo "  -c, --coverage          Enable coverage reporting (default: true)"
    echo "  --no-coverage           Disable coverage reporting"
    echo "  -t, --timeout DURATION  Set test timeout (default: 300s)"
    echo "  --db-url URL            Set test database URL"
    echo
    echo "Test Names:"
    echo "  repository, repo        Run email event repository tests"
    echo "  recorder, record        Run email event recorder tests"
    echo "  metrics, metric         Run email metrics service tests"
    echo "  handler, admin          Run admin handler email tests"
    echo "  integration, int        Run email service integration tests"
    echo "  (no test name)          Run all tests"
    echo
    echo "Examples:"
    echo "  $0                      # Run all tests"
    echo "  $0 repository           # Run only repository tests"
    echo "  $0 -v --no-coverage     # Run all tests with verbose output, no coverage"
    echo "  $0 -t 600s integration  # Run integration tests with 10 minute timeout"
    echo
    echo "Environment Variables:"
    echo "  TEST_DATABASE_URL       PostgreSQL connection string for tests"
    echo "  VERBOSE                 Enable verbose output (true/false)"
    echo "  COVERAGE                Enable coverage reporting (true/false)"
    echo "  TIMEOUT                 Test timeout duration"
}

# Main execution
main() {
    local test_name=""
    
    # Parse command line arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h|--help)
                show_help
                exit 0
                ;;
            -v|--verbose)
                VERBOSE=true
                shift
                ;;
            -c|--coverage)
                COVERAGE=true
                shift
                ;;
            --no-coverage)
                COVERAGE=false
                shift
                ;;
            -t|--timeout)
                TIMEOUT="$2"
                shift 2
                ;;
            --db-url)
                TEST_DATABASE_URL="$2"
                shift 2
                ;;
            -*)
                print_status $RED "Unknown option: $1"
                show_help
                exit 1
                ;;
            *)
                if [ -z "$test_name" ]; then
                    test_name="$1"
                else
                    print_status $RED "Multiple test names specified"
                    exit 1
                fi
                shift
                ;;
        esac
    done
    
    # Setup
    check_prerequisites
    setup_test_environment
    
    # Trap to ensure cleanup
    trap cleanup_test_environment EXIT
    
    # Run tests
    if [ -z "$test_name" ]; then
        run_all_tests
    else
        run_specific_test "$test_name"
    fi
}

# Run main function with all arguments
main "$@"