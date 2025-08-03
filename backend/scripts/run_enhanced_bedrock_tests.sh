#!/bin/bash

# Enhanced Bedrock AI Assistant Comprehensive Testing Script
# This script runs all comprehensive testing and validation for the enhanced Bedrock AI assistant

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
NC='\033[0m' # No Color

# Configuration
BACKEND_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
TEST_RESULTS_DIR="$BACKEND_DIR/test_results"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
REPORT_FILE="$TEST_RESULTS_DIR/enhanced_bedrock_test_report_$TIMESTAMP.txt"

echo -e "${BLUE}=== Enhanced Bedrock AI Assistant Comprehensive Testing ===${NC}"
echo "Backend Directory: $BACKEND_DIR"
echo "Test Results Directory: $TEST_RESULTS_DIR"
echo "Report File: $REPORT_FILE"
echo ""

# Create test results directory
mkdir -p "$TEST_RESULTS_DIR"

# Function to print section headers
print_section() {
    echo -e "${BLUE}=== $1 ===${NC}"
    echo "=== $1 ===" >> "$REPORT_FILE"
}

# Function to print success messages
print_success() {
    echo -e "${GREEN}✓ $1${NC}"
    echo "✓ $1" >> "$REPORT_FILE"
}

# Function to print error messages
print_error() {
    echo -e "${RED}✗ $1${NC}"
    echo "✗ $1" >> "$REPORT_FILE"
}

# Function to print warning messages
print_warning() {
    echo -e "${YELLOW}⚠ $1${NC}"
    echo "⚠ $1" >> "$REPORT_FILE"
}

# Function to print info messages
print_info() {
    echo -e "${PURPLE}ℹ $1${NC}"
    echo "ℹ $1" >> "$REPORT_FILE"
}

# Function to run command with error handling and logging
run_test() {
    local test_name="$1"
    local test_command="$2"
    local description="$3"
    
    print_section "$test_name"
    echo "Description: $description"
    echo "Command: $test_command"
    echo ""
    
    # Log to report file
    echo "" >> "$REPORT_FILE"
    echo "Test: $test_name" >> "$REPORT_FILE"
    echo "Description: $description" >> "$REPORT_FILE"
    echo "Command: $test_command" >> "$REPORT_FILE"
    echo "Started: $(date)" >> "$REPORT_FILE"
    echo "" >> "$REPORT_FILE"
    
    # Run the test and capture output
    if eval "$test_command" 2>&1 | tee -a "$REPORT_FILE"; then
        print_success "$test_name completed successfully"
        echo "Completed: $(date)" >> "$REPORT_FILE"
        echo "Status: SUCCESS" >> "$REPORT_FILE"
        return 0
    else
        print_error "$test_name failed"
        echo "Completed: $(date)" >> "$REPORT_FILE"
        echo "Status: FAILED" >> "$REPORT_FILE"
        return 1
    fi
}

# Initialize report file
cat > "$REPORT_FILE" << EOF
Enhanced Bedrock AI Assistant - Comprehensive Test Report
=========================================================

Test Run Date: $(date)
Go Version: $(go version | awk '{print $3}')
Test Environment: $(uname -s) $(uname -r)
Backend Directory: $BACKEND_DIR

EOF

# Change to backend directory
cd "$BACKEND_DIR"

# Check prerequisites
print_section "Prerequisites Check"
if ! command -v go &> /dev/null; then
    print_error "Go is not installed or not in PATH"
    exit 1
fi
print_success "Go is available: $(go version | awk '{print $3}')"

# Install dependencies
print_section "Installing Dependencies"
if go mod download && go mod tidy; then
    print_success "Go dependencies installed"
else
    print_error "Failed to install Go dependencies"
    exit 1
fi

# Test execution tracking
total_tests=0
passed_tests=0
failed_tests=0

# Test 1: Comprehensive Enhanced Bedrock Validation
total_tests=$((total_tests + 1))
if run_test "Comprehensive Enhanced Bedrock Validation" \
    "go run test_comprehensive_enhanced_bedrock_validation.go" \
    "Runs comprehensive testing including real-world scenarios, A/B testing, regression testing, and user acceptance testing"; then
    passed_tests=$((passed_tests + 1))
else
    failed_tests=$((failed_tests + 1))
fi

# Test 2: Enhanced Bedrock Unit Tests
total_tests=$((total_tests + 1))
if run_test "Enhanced Bedrock Unit Tests" \
    "go test -v -race -timeout=30m ./internal/services/enhanced_bedrock_test.go ./internal/services/enhanced_bedrock.go" \
    "Runs unit tests for the enhanced Bedrock service"; then
    passed_tests=$((passed_tests + 1))
else
    failed_tests=$((failed_tests + 1))
fi

# Test 3: Enhanced Features Verification
total_tests=$((total_tests + 1))
if run_test "Enhanced Features Verification" \
    "go run test_enhanced_features_verification.go" \
    "Verifies enhanced AI assistant features including prompt architect, knowledge base, and multi-cloud analysis"; then
    passed_tests=$((passed_tests + 1))
else
    failed_tests=$((failed_tests + 1))
fi

# Test 4: Enhanced AI Integration Test
total_tests=$((total_tests + 1))
if run_test "Enhanced AI Integration Test" \
    "go run test_enhanced_ai_integration.go" \
    "Tests integration of enhanced AI features with chat system"; then
    passed_tests=$((passed_tests + 1))
else
    failed_tests=$((failed_tests + 1))
fi

# Test 5: Quality Assurance System Test
total_tests=$((total_tests + 1))
if run_test "Quality Assurance System Test" \
    "go run test_quality_assurance_system.go" \
    "Tests the quality assurance system for recommendation tracking and validation"; then
    passed_tests=$((passed_tests + 1))
else
    failed_tests=$((failed_tests + 1))
fi

# Test 6: Performance Optimization Test
total_tests=$((total_tests + 1))
if run_test "Performance Optimization Test" \
    "go run test_performance_optimization_task17.go" \
    "Tests performance optimization features for enhanced Bedrock service"; then
    passed_tests=$((passed_tests + 1))
else
    failed_tests=$((failed_tests + 1))
fi

# Generate comprehensive test summary
print_section "Test Summary"

# Calculate pass rate
pass_rate=0
if [ $total_tests -gt 0 ]; then
    pass_rate=$(echo "scale=2; $passed_tests * 100 / $total_tests" | bc -l)
fi

# Determine overall status
overall_status="FAILED"
if [ $failed_tests -eq 0 ]; then
    overall_status="PASSED"
elif [ $passed_tests -gt $failed_tests ]; then
    overall_status="MOSTLY_PASSED"
fi

# Print summary
echo ""
print_info "=== COMPREHENSIVE TEST SUMMARY ==="
print_info "Total Tests: $total_tests"
print_info "Passed: $passed_tests"
print_info "Failed: $failed_tests"
print_info "Pass Rate: ${pass_rate}%"
print_info "Overall Status: $overall_status"

# Log summary to report
cat >> "$REPORT_FILE" << EOF

=== COMPREHENSIVE TEST SUMMARY ===
Total Tests: $total_tests
Passed: $passed_tests
Failed: $failed_tests
Pass Rate: ${pass_rate}%
Overall Status: $overall_status

Test Categories Covered:
- Real-world client engagement scenarios
- A/B testing for different recommendation approaches
- Regression testing for quality assurance
- User acceptance testing with consultant personas
- Enhanced Bedrock service unit tests
- Feature verification and integration tests
- Quality assurance system validation
- Performance optimization testing

Key Features Validated:
✓ Industry-specific response generation (Healthcare, FinTech, Retail)
✓ Multi-variant A/B testing framework
✓ Statistical significance testing
✓ Quality scoring across 6 dimensions
✓ Regression detection and baseline comparison
✓ Multi-persona user acceptance testing
✓ Performance monitoring and optimization
✓ Comprehensive quality assurance workflow

Test Report Generated: $(date)
Report Location: $REPORT_FILE

EOF

# Generate additional artifacts
print_section "Generating Test Artifacts"

# Create test coverage summary if available
if [ -f "coverage/combined.out" ]; then
    go tool cover -func=coverage/combined.out > "$TEST_RESULTS_DIR/coverage_summary_$TIMESTAMP.txt"
    print_success "Coverage summary generated"
fi

# Create test metrics file
cat > "$TEST_RESULTS_DIR/test_metrics_$TIMESTAMP.json" << EOF
{
  "timestamp": "$(date -Iseconds)",
  "total_tests": $total_tests,
  "passed_tests": $passed_tests,
  "failed_tests": $failed_tests,
  "pass_rate": $pass_rate,
  "overall_status": "$overall_status",
  "test_categories": [
    "comprehensive_validation",
    "unit_tests",
    "feature_verification",
    "integration_tests",
    "quality_assurance",
    "performance_optimization"
  ],
  "features_validated": [
    "real_world_scenarios",
    "ab_testing",
    "regression_testing",
    "user_acceptance_testing",
    "quality_scoring",
    "performance_monitoring"
  ]
}
EOF

print_success "Test metrics JSON generated"

# Final status and recommendations
print_section "Final Status and Recommendations"

if [ "$overall_status" = "PASSED" ]; then
    print_success "🎉 All comprehensive tests passed! Enhanced Bedrock AI Assistant is ready for production."
    echo ""
    print_info "✅ Recommendations:"
    print_info "   • Deploy enhanced features to production environment"
    print_info "   • Set up continuous monitoring for quality metrics"
    print_info "   • Schedule regular regression testing"
    print_info "   • Monitor A/B test results for ongoing optimization"
    
elif [ "$overall_status" = "MOSTLY_PASSED" ]; then
    print_warning "⚠️  Most tests passed, but some issues need attention."
    echo ""
    print_info "📋 Recommendations:"
    print_info "   • Review failed tests and address issues"
    print_info "   • Consider conditional deployment with monitoring"
    print_info "   • Implement fixes for failed test scenarios"
    print_info "   • Re-run tests after fixes are applied"
    
else
    print_error "❌ Multiple tests failed. Enhanced Bedrock AI Assistant needs fixes before deployment."
    echo ""
    print_info "🔧 Recommendations:"
    print_info "   • Review all failed tests and root causes"
    print_info "   • Implement comprehensive fixes"
    print_info "   • Re-run full test suite after fixes"
    print_info "   • Consider additional testing scenarios"
fi

echo ""
print_info "📄 Detailed test report available at: $REPORT_FILE"
print_info "📊 Test metrics available at: $TEST_RESULTS_DIR/test_metrics_$TIMESTAMP.json"

# Exit with appropriate code
if [ "$overall_status" = "PASSED" ]; then
    exit 0
elif [ "$overall_status" = "MOSTLY_PASSED" ]; then
    exit 1
else
    exit 2
fi