#!/bin/bash

# Test Categories Runner Script
# This script runs specific test categories based on the organized test structure
# Requirement: 4.4 - Script for running specific test categories

set -e  # Exit on any error

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

print_status() {
    local color=$1
    local message=$2
    echo -e "${color}${message}${NC}"
}

print_help() {
    echo "Test Categories Runner"
    echo "Usage: $0 [CATEGORY] [OPTIONS]"
    echo
    echo "Categories:"
    echo "  unit           Run unit tests (*_test.go files)"
    echo "  integration    Run integration tests (testing/integration/)"
    echo "  email          Run email system tests (testing/email/)"
    echo "  performance    Run performance tests (testing/performance/)"
    echo "  all            Run all test categories"
    echo "  list           List available test files by category"
    echo
    echo "Options:"
    echo "  -v, --verbose  Verbose output"
    echo "  -c, --coverage Generate coverage report"
    echo "  --timeout=30s  Set test timeout (default: 30s)"
    echo "  --help         Show this help message"
    echo
    echo "Examples:"
    echo "  $0 unit                    # Run unit tests"
    echo "  $0 integration -v          # Run integration tests with verbose output"
    echo "  $0 email --coverage        # Run email tests with coverage"
    echo "  $0 all --timeout=60s       # Run all tests with 60s timeout"
}

# Change to backend directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BACKEND_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
cd "$BACKEND_DIR"

# Default options
VERBOSE=false
COVERAGE=false
TIMEOUT="30s"
CATEGORY=""

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        unit|integration|email|performance|all|list)
            CATEGORY="$1"
            shift
            ;;
        -v|--verbose)
            VERBOSE=true
            shift
            ;;
        -c|--coverage)
            COVERAGE=true
            shift
            ;;
        --timeout=*)
            TIMEOUT="${1#*=}"
            shift
            ;;
        --help)
            print_help
            exit 0
            ;;
        *)
            print_status $RED "Unknown option: $1"
            print_help
            exit 1
            ;;
    esac
done

# Check if category is provided
if [ -z "$CATEGORY" ]; then
    print_status $RED "Error: No test category specified"
    print_help
    exit 1
fi

# Build verbose flag
VERBOSE_FLAG=""
if [ "$VERBOSE" = true ]; then
    VERBOSE_FLAG="-v"
fi

# Build coverage flags
COVERAGE_FLAGS=""
COVERAGE_OUTPUT=""
if [ "$COVERAGE" = true ]; then
    mkdir -p coverage
    COVERAGE_FLAGS="-coverprofile=coverage/coverage.out -covermode=atomic"
    COVERAGE_OUTPUT="coverage/coverage.out"
fi

# Function to run tests with common options
run_tests() {
    local test_path="$1"
    local description="$2"
    
    print_status $BLUE "Running $description..."
    
    local cmd="go test $VERBOSE_FLAG $COVERAGE_FLAGS -timeout=$TIMEOUT $test_path"
    print_status $BLUE "Command: $cmd"
    
    if eval "$cmd"; then
        print_status $GREEN "✅ $description completed successfully"
        
        # Generate coverage report if requested
        if [ "$COVERAGE" = true ] && [ -f "$COVERAGE_OUTPUT" ]; then
            print_status $BLUE "Generating coverage report..."
            go tool cover -html="$COVERAGE_OUTPUT" -o "coverage/${description// /_}_coverage.html"
            
            # Show coverage summary
            COVERAGE_PERCENT=$(go tool cover -func="$COVERAGE_OUTPUT" | grep total | awk '{print $3}')
            print_status $BLUE "Coverage: $COVERAGE_PERCENT"
        fi
        
        return 0
    else
        print_status $RED "❌ $description failed"
        return 1
    fi
}

# Function to list test files
list_tests() {
    print_status $BLUE "=== Available Test Files by Category ==="
    echo
    
    print_status $BLUE "Unit Tests (*_test.go in internal/):"
    find internal -name "*_test.go" -type f | sort | sed 's/^/  /'
    echo
    
    print_status $BLUE "Integration Tests (testing/integration/):"
    if [ -d "testing/integration" ]; then
        find testing/integration -name "*.go" -type f | sort | sed 's/^/  /'
    else
        print_status $YELLOW "  No integration test directory found"
    fi
    echo
    
    print_status $BLUE "Email Tests (testing/email/):"
    if [ -d "testing/email" ]; then
        find testing/email -name "*.go" -type f | sort | sed 's/^/  /'
    else
        print_status $YELLOW "  No email test directory found"
    fi
    echo
    
    print_status $BLUE "Performance Tests (testing/performance/):"
    if [ -d "testing/performance" ]; then
        find testing/performance -name "*.go" -type f | sort | sed 's/^/  /'
    else
        print_status $YELLOW "  No performance test directory found"
    fi
    echo
    
    print_status $BLUE "Standalone Test Files (test_*.go in root):"
    find . -maxdepth 1 -name "test_*.go" -type f | sort | sed 's/^/  /'
}

# Main execution
print_status $BLUE "=== Test Categories Runner ==="
print_status $BLUE "Backend Directory: $BACKEND_DIR"
print_status $BLUE "Category: $CATEGORY"
print_status $BLUE "Timeout: $TIMEOUT"
print_status $BLUE "Verbose: $VERBOSE"
print_status $BLUE "Coverage: $COVERAGE"
echo

case $CATEGORY in
    "list")
        list_tests
        exit 0
        ;;
    
    "unit")
        print_status $BLUE "=== Running Unit Tests ==="
        run_tests "./internal/..." "Unit Tests"
        ;;
    
    "integration")
        print_status $BLUE "=== Running Integration Tests ==="
        if [ -d "testing/integration" ]; then
            run_tests "./testing/integration/..." "Integration Tests"
        else
            print_status $YELLOW "No integration test directory found, checking for standalone integration tests..."
            # Look for integration test files in root
            INTEGRATION_FILES=$(find . -maxdepth 1 -name "test_*integration*.go" -o -name "test_*bedrock*.go" -o -name "test_*chat*.go" -o -name "test_*ai*.go" | tr '\n' ' ')
            if [ -n "$INTEGRATION_FILES" ]; then
                run_tests "$INTEGRATION_FILES" "Standalone Integration Tests"
            else
                print_status $RED "No integration tests found"
                exit 1
            fi
        fi
        ;;
    
    "email")
        print_status $BLUE "=== Running Email Tests ==="
        if [ -d "testing/email" ]; then
            run_tests "./testing/email/..." "Email Tests"
        else
            print_status $YELLOW "No email test directory found, checking for standalone email tests..."
            # Look for email test files in root
            EMAIL_FILES=$(find . -maxdepth 1 -name "test_*email*.go" -o -name "test_*ses*.go" | tr '\n' ' ')
            if [ -n "$EMAIL_FILES" ]; then
                run_tests "$EMAIL_FILES" "Standalone Email Tests"
            else
                print_status $RED "No email tests found"
                exit 1
            fi
        fi
        ;;
    
    "performance")
        print_status $BLUE "=== Running Performance Tests ==="
        if [ -d "testing/performance" ]; then
            run_tests "./testing/performance/..." "Performance Tests"
        else
            print_status $YELLOW "No performance test directory found, checking for standalone performance tests..."
            # Look for performance test files in root
            PERF_FILES=$(find . -maxdepth 1 -name "test_*performance*.go" -o -name "test_*load*.go" -o -name "test_*benchmark*.go" | tr '\n' ' ')
            if [ -n "$PERF_FILES" ]; then
                run_tests "$PERF_FILES" "Standalone Performance Tests"
            else
                print_status $RED "No performance tests found"
                exit 1
            fi
        fi
        ;;
    
    "all")
        print_status $BLUE "=== Running All Test Categories ==="
        
        # Track overall success
        OVERALL_SUCCESS=true
        
        # Run unit tests
        print_status $BLUE "\n--- Unit Tests ---"
        if ! run_tests "./internal/..." "Unit Tests"; then
            OVERALL_SUCCESS=false
        fi
        
        # Run integration tests
        print_status $BLUE "\n--- Integration Tests ---"
        if [ -d "testing/integration" ]; then
            if ! run_tests "./testing/integration/..." "Integration Tests"; then
                OVERALL_SUCCESS=false
            fi
        else
            INTEGRATION_FILES=$(find . -maxdepth 1 -name "test_*integration*.go" -o -name "test_*bedrock*.go" -o -name "test_*chat*.go" -o -name "test_*ai*.go" | tr '\n' ' ')
            if [ -n "$INTEGRATION_FILES" ]; then
                if ! run_tests "$INTEGRATION_FILES" "Standalone Integration Tests"; then
                    OVERALL_SUCCESS=false
                fi
            fi
        fi
        
        # Run email tests
        print_status $BLUE "\n--- Email Tests ---"
        if [ -d "testing/email" ]; then
            if ! run_tests "./testing/email/..." "Email Tests"; then
                OVERALL_SUCCESS=false
            fi
        else
            EMAIL_FILES=$(find . -maxdepth 1 -name "test_*email*.go" -o -name "test_*ses*.go" | tr '\n' ' ')
            if [ -n "$EMAIL_FILES" ]; then
                if ! run_tests "$EMAIL_FILES" "Standalone Email Tests"; then
                    OVERALL_SUCCESS=false
                fi
            fi
        fi
        
        # Run performance tests
        print_status $BLUE "\n--- Performance Tests ---"
        if [ -d "testing/performance" ]; then
            if ! run_tests "./testing/performance/..." "Performance Tests"; then
                OVERALL_SUCCESS=false
            fi
        else
            PERF_FILES=$(find . -maxdepth 1 -name "test_*performance*.go" -o -name "test_*load*.go" -o -name "test_*benchmark*.go" | tr '\n' ' ')
            if [ -n "$PERF_FILES" ]; then
                if ! run_tests "$PERF_FILES" "Standalone Performance Tests"; then
                    OVERALL_SUCCESS=false
                fi
            fi
        fi
        
        # Final summary
        echo
        print_status $BLUE "=== Test Suite Summary ==="
        if [ "$OVERALL_SUCCESS" = true ]; then
            print_status $GREEN "✅ All test categories completed successfully!"
        else
            print_status $RED "❌ Some test categories failed"
            exit 1
        fi
        ;;
    
    *)
        print_status $RED "Unknown category: $CATEGORY"
        print_help
        exit 1
        ;;
esac

echo
print_status $GREEN "Test execution completed!"

# Show coverage report location if generated
if [ "$COVERAGE" = true ] && [ -f "coverage/coverage.out" ]; then
    print_status $BLUE "Coverage report available at: coverage/${CATEGORY// /_}_coverage.html"
fi