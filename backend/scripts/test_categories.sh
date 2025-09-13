#!/bin/bash

# Test Categories Script
# This script allows running specific categories of tests

set -e

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

print_status() {
    local color=$1
    local message=$2
    echo -e "${color}${message}${NC}"
}

show_usage() {
    echo "Usage: $0 [category]"
    echo
    echo "Available categories:"
    echo "  unit         - Run unit tests (*_test.go files)"
    echo "  integration  - Run integration test executables"
    echo "  email        - Run email test executables"
    echo "  performance  - Run performance test executables"
    echo "  all          - Run all test categories"
    echo "  list         - List available test files by category"
    echo
    echo "Examples:"
    echo "  $0 unit"
    echo "  $0 integration"
    echo "  $0 list"
}

# Change to backend directory
cd "$(dirname "$0")/.."

case "${1:-}" in
    "unit")
        print_status $BLUE "Running unit tests..."
        go test ./internal/... -v
        ;;
    
    "integration")
        print_status $BLUE "Listing integration test executables..."
        if [ -d "testing/integration" ]; then
            find testing/integration -name "*.go" -type f | head -10
            print_status $YELLOW "Note: Integration tests are standalone executables"
            print_status $YELLOW "Run them individually with: go run testing/integration/[filename].go"
        else
            print_status $RED "Integration test directory not found"
        fi
        ;;
    
    "email")
        print_status $BLUE "Listing email test executables..."
        if [ -d "testing/email" ]; then
            find testing/email -name "*.go" -type f | head -10
            print_status $YELLOW "Note: Email tests are standalone executables"
            print_status $YELLOW "Run them individually with: go run testing/email/[filename].go"
        else
            print_status $RED "Email test directory not found"
        fi
        ;;
    
    "performance")
        print_status $BLUE "Listing performance test executables..."
        if [ -d "testing/performance" ]; then
            find testing/performance -name "*.go" -type f
            print_status $YELLOW "Note: Performance tests are standalone executables"
            print_status $YELLOW "Run them individually with: go run testing/performance/[filename].go"
        else
            print_status $RED "Performance test directory not found"
        fi
        ;;
    
    "all")
        print_status $BLUE "Running all available tests..."
        echo
        print_status $BLUE "1. Unit tests:"
        go test ./internal/... -v
        echo
        print_status $BLUE "2. Test categories available:"
        $0 list
        ;;
    
    "list")
        print_status $BLUE "Available test files by category:"
        echo
        
        print_status $BLUE "Unit Tests (*_test.go):"
        find ./internal -name "*_test.go" -type f | wc -l | xargs echo "  Found" | sed 's/$/ unit test files/'
        
        if [ -d "testing/integration" ]; then
            print_status $BLUE "Integration Tests:"
            count=$(find testing/integration -name "*.go" -type f | wc -l)
            echo "  Found $count integration test files"
        fi
        
        if [ -d "testing/email" ]; then
            print_status $BLUE "Email Tests:"
            count=$(find testing/email -name "*.go" -type f | wc -l)
            echo "  Found $count email test files"
        fi
        
        if [ -d "testing/performance" ]; then
            print_status $BLUE "Performance Tests:"
            count=$(find testing/performance -name "*.go" -type f | wc -l)
            echo "  Found $count performance test files"
        fi
        
        print_status $BLUE "Standalone Tests (root directory):"
        count=$(find . -maxdepth 1 -name "test_*.go" -type f | wc -l)
        echo "  Found $count standalone test files"
        ;;
    
    "help"|"-h"|"--help"|"")
        show_usage
        ;;
    
    *)
        print_status $RED "Unknown category: $1"
        echo
        show_usage
        exit 1
        ;;
esac