# Testing Guide

This guide provides comprehensive instructions for running tests and adding new test files in the Cloud Consulting Platform.

## Table of Contents

1. [Test Organization](#test-organization)
2. [Running Tests](#running-tests)
3. [Test Categories](#test-categories)
4. [Adding New Tests](#adding-new-tests)
5. [Test Development Guidelines](#test-development-guidelines)
6. [Continuous Integration](#continuous-integration)
7. [Troubleshooting](#troubleshooting)

## Test Organization

The project follows a clean separation between different types of tests:

```
backend/
├── internal/               # Application code with unit tests
│   ├── domain/
│   │   └── *_test.go      # Domain model unit tests
│   ├── handlers/
│   │   └── *_test.go      # HTTP handler unit tests
│   ├── services/
│   │   └── *_test.go      # Service layer unit tests
│   └── repositories/
│       └── *_test.go      # Repository unit tests
├── testing/               # Organized test utilities
│   ├── integration/       # Integration test executables
│   ├── email/            # Email system test executables
│   ├── performance/      # Performance test executables
│   └── shared/           # Shared test utilities and mocks
└── scripts/              # Test execution scripts
```

### Test File Naming Conventions

- **Unit Tests**: `*_test.go` (alongside source code)
- **Integration Tests**: `test_*_integration.go` or descriptive names
- **Performance Tests**: `test_*_performance.go` or `test_*_load.go`
- **Email Tests**: `test_*_email.go` or `test_*_ses.go`
- **Standalone Executables**: `test_*.go` with `package main`

## Running Tests

### Using Test Categories Script (Recommended)

The project provides a convenient script for running different test categories:

```bash
# Show available test categories and files
./backend/scripts/run_test_categories.sh list

# Run specific test categories
./backend/scripts/run_test_categories.sh unit              # Unit tests only
./backend/scripts/run_test_categories.sh integration      # Integration tests
./backend/scripts/run_test_categories.sh email           # Email system tests
./backend/scripts/run_test_categories.sh performance     # Performance tests
./backend/scripts/run_test_categories.sh all             # All categories

# Run with options
./backend/scripts/run_test_categories.sh unit -v          # Verbose output
./backend/scripts/run_test_categories.sh email --coverage # With coverage report
./backend/scripts/run_test_categories.sh all --timeout=60s # Custom timeout
```

### Manual Test Execution

If you prefer manual control:

```bash
cd backend

# Unit tests (standard Go tests)
go test -v ./internal/...

# Integration tests (if organized in testing directory)
go test -v ./testing/integration/...

# Email tests (if organized in testing directory)
go test -v ./testing/email/...

# Performance tests (if organized in testing directory)
go test -v ./testing/performance/...

# Standalone test executables
go run test_specific_feature.go

# All tests with coverage
go test -v -coverprofile=coverage.out ./internal/...
go tool cover -html=coverage.out -o coverage.html
```

## Test Categories

### 1. Unit Tests

**Location**: `internal/**/*_test.go`
**Purpose**: Test individual functions, methods, and components in isolation
**Execution**: `./backend/scripts/run_test_categories.sh unit`

**Characteristics:**
- Fast execution (< 1 second per test)
- No external dependencies
- Use mocks for dependencies
- High code coverage target (>90%)

**Example Structure:**
```go
package services_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

func TestServiceMethod(t *testing.T) {
    // Arrange
    mockDep := &MockDependency{}
    service := NewService(mockDep)
    
    // Act
    result, err := service.Method(input)
    
    // Assert
    assert.NoError(t, err)
    assert.Equal(t, expected, result)
    mockDep.AssertExpectations(t)
}
```

### 2. Integration Tests

**Location**: `testing/integration/`
**Purpose**: Test component interactions, API endpoints, and system integration
**Execution**: `./backend/scripts/run_test_categories.sh integration`

**Characteristics:**
- Moderate execution time (1-10 seconds per test)
- May use real databases or external services
- Test complete workflows
- Focus on component boundaries

**Example Structure:**
```go
package main

import (
    "testing"
    "net/http/httptest"
)

func TestAPIEndpointIntegration(t *testing.T) {
    // Setup test server
    server := setupTestServer()
    defer server.Close()
    
    // Test complete API workflow
    response := makeAPIRequest(server.URL + "/api/endpoint")
    
    // Verify integration works
    assert.Equal(t, http.StatusOK, response.StatusCode)
}
```

### 3. Email Tests

**Location**: `testing/email/`
**Purpose**: Test email system functionality, SES integration, and template rendering
**Execution**: `./backend/scripts/run_test_categories.sh email`

**Characteristics:**
- May require AWS credentials for full testing
- Test email formatting and delivery
- Verify template rendering
- Test SES connectivity

**Example Structure:**
```go
package main

import (
    "testing"
    "context"
)

func TestEmailDelivery(t *testing.T) {
    // Setup email service
    emailService := setupEmailService()
    
    // Test email sending
    err := emailService.SendEmail(context.Background(), emailData)
    
    // Verify email was sent
    assert.NoError(t, err)
}
```

### 4. Performance Tests

**Location**: `testing/performance/`
**Purpose**: Test system performance, load handling, and resource usage
**Execution**: `./backend/scripts/run_test_categories.sh performance`

**Characteristics:**
- Longer execution time (10+ seconds)
- Test under load conditions
- Measure response times and throughput
- Identify performance bottlenecks

**Example Structure:**
```go
package main

import (
    "testing"
    "sync"
    "time"
)

func TestLoadHandling(t *testing.T) {
    // Setup load test
    concurrency := 100
    requests := 1000
    
    // Execute concurrent requests
    var wg sync.WaitGroup
    start := time.Now()
    
    for i := 0; i < concurrency; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            // Perform load test operations
        }()
    }
    
    wg.Wait()
    duration := time.Since(start)
    
    // Verify performance metrics
    assert.Less(t, duration, time.Second*30)
}
```

## Adding New Tests

### Guidelines for New Test Files

1. **Choose the Right Category:**
   - **Unit Test**: Testing a single function/method → `internal/package/*_test.go`
   - **Integration Test**: Testing component interaction → `testing/integration/test_*.go`
   - **Email Test**: Testing email functionality → `testing/email/test_*.go`
   - **Performance Test**: Testing performance/load → `testing/performance/test_*.go`

2. **Follow Naming Conventions:**
   ```bash
   # Unit tests (alongside source code)
   internal/services/chat_service_test.go
   
   # Integration tests
   testing/integration/test_chat_api_integration.go
   
   # Email tests
   testing/email/test_ses_connectivity.go
   
   # Performance tests
   testing/performance/test_chat_load.go
   ```

3. **Use Appropriate Package Declarations:**
   ```go
   // Unit tests (same package)
   package services
   
   // Unit tests (separate test package)
   package services_test
   
   // Standalone test executables
   package main
   ```

### Step-by-Step: Adding a New Unit Test

1. **Create the test file** alongside your source code:
   ```bash
   # If you have: internal/services/new_service.go
   # Create: internal/services/new_service_test.go
   ```

2. **Write the test:**
   ```go
   package services_test
   
   import (
       "testing"
       "github.com/stretchr/testify/assert"
       "your-project/internal/services"
   )
   
   func TestNewService_Method(t *testing.T) {
       // Test implementation
   }
   ```

3. **Run the test:**
   ```bash
   ./backend/scripts/run_test_categories.sh unit -v
   ```

### Step-by-Step: Adding a New Integration Test

1. **Create the test file** in the integration directory:
   ```bash
   touch backend/testing/integration/test_new_feature_integration.go
   ```

2. **Write the test:**
   ```go
   package main
   
   import (
       "testing"
       "fmt"
   )
   
   func TestNewFeatureIntegration(t *testing.T) {
       // Integration test implementation
   }
   
   func main() {
       fmt.Println("Run with: go test -v ./test_new_feature_integration.go")
   }
   ```

3. **Run the test:**
   ```bash
   ./backend/scripts/run_test_categories.sh integration -v
   ```

### Step-by-Step: Adding a New Email Test

1. **Create the test file** in the email directory:
   ```bash
   touch backend/testing/email/test_new_email_feature.go
   ```

2. **Write the test:**
   ```go
   package main
   
   import (
       "testing"
       "context"
   )
   
   func TestNewEmailFeature(t *testing.T) {
       // Email test implementation
   }
   
   func main() {
       // Standalone execution logic
   }
   ```

3. **Run the test:**
   ```bash
   ./backend/scripts/run_test_categories.sh email -v
   ```

### Step-by-Step: Adding a New Performance Test

1. **Create the test file** in the performance directory:
   ```bash
   touch backend/testing/performance/test_new_performance.go
   ```

2. **Write the test:**
   ```go
   package main
   
   import (
       "testing"
       "time"
   )
   
   func TestNewPerformance(t *testing.T) {
       start := time.Now()
       
       // Performance test implementation
       
       duration := time.Since(start)
       t.Logf("Operation took: %v", duration)
   }
   
   func BenchmarkNewFeature(b *testing.B) {
       for i := 0; i < b.N; i++ {
           // Benchmark implementation
       }
   }
   ```

3. **Run the test:**
   ```bash
   ./backend/scripts/run_test_categories.sh performance -v
   ```

## Test Development Guidelines

### Best Practices

1. **Test Structure (AAA Pattern):**
   ```go
   func TestFunction(t *testing.T) {
       // Arrange - Set up test data and dependencies
       input := "test input"
       expected := "expected output"
       
       // Act - Execute the function being tested
       result, err := FunctionUnderTest(input)
       
       // Assert - Verify the results
       assert.NoError(t, err)
       assert.Equal(t, expected, result)
   }
   ```

2. **Use Table-Driven Tests for Multiple Cases:**
   ```go
   func TestFunction(t *testing.T) {
       tests := []struct {
           name     string
           input    string
           expected string
           wantErr  bool
       }{
           {"valid input", "test", "result", false},
           {"invalid input", "", "", true},
       }
       
       for _, tt := range tests {
           t.Run(tt.name, func(t *testing.T) {
               result, err := FunctionUnderTest(tt.input)
               
               if tt.wantErr {
                   assert.Error(t, err)
               } else {
                   assert.NoError(t, err)
                   assert.Equal(t, tt.expected, result)
               }
           })
       }
   }
   ```

3. **Use Mocks for Dependencies:**
   ```go
   type MockDependency struct {
       mock.Mock
   }
   
   func (m *MockDependency) Method(arg string) (string, error) {
       args := m.Called(arg)
       return args.String(0), args.Error(1)
   }
   
   func TestWithMock(t *testing.T) {
       mockDep := &MockDependency{}
       mockDep.On("Method", "input").Return("output", nil)
       
       service := NewService(mockDep)
       result, err := service.UsesDependency("input")
       
       assert.NoError(t, err)
       assert.Equal(t, "output", result)
       mockDep.AssertExpectations(t)
   }
   ```

4. **Test Error Conditions:**
   ```go
   func TestErrorHandling(t *testing.T) {
       // Test various error conditions
       _, err := FunctionThatCanFail("")
       assert.Error(t, err)
       assert.Contains(t, err.Error(), "expected error message")
   }
   ```

5. **Use Subtests for Organization:**
   ```go
   func TestComplexFunction(t *testing.T) {
       t.Run("success case", func(t *testing.T) {
           // Test success scenario
       })
       
       t.Run("error case", func(t *testing.T) {
           // Test error scenario
       })
       
       t.Run("edge case", func(t *testing.T) {
           // Test edge cases
       })
   }
   ```

### Testing Utilities

The project provides shared testing utilities in `backend/testing/shared/`:

```go
// Use shared test utilities
import "your-project/testing/shared"

func TestWithSharedUtils(t *testing.T) {
    // Use shared fixtures
    user := shared.CreateTestUser()
    
    // Use shared HTTP helpers
    response := shared.MakeTestRequest("GET", "/api/endpoint", nil)
    
    // Use shared assertions
    shared.AssertValidResponse(t, response)
}
```

## Continuous Integration

### Running Tests in CI/CD

The test categories script is designed for CI/CD integration:

```bash
# In your CI/CD pipeline
./backend/scripts/run_test_categories.sh all --coverage --timeout=300s
```

### Coverage Requirements

- **Unit Tests**: >90% coverage
- **Integration Tests**: >80% coverage
- **Overall**: >85% coverage

### CI/CD Configuration Example

```yaml
# GitHub Actions example
- name: Run Unit Tests
  run: ./backend/scripts/run_test_categories.sh unit --coverage

- name: Run Integration Tests
  run: ./backend/scripts/run_test_categories.sh integration

- name: Upload Coverage
  uses: codecov/codecov-action@v1
  with:
    file: ./backend/coverage/coverage.out
```

## Troubleshooting

### Common Issues

1. **Tests Not Found:**
   ```bash
   # Check test organization
   ./backend/scripts/run_test_categories.sh list
   
   # Verify file locations match expected patterns
   ```

2. **Import Path Errors:**
   ```bash
   # Update Go modules
   cd backend
   go mod tidy
   
   # Check import paths in test files
   ```

3. **Test Timeouts:**
   ```bash
   # Increase timeout
   ./backend/scripts/run_test_categories.sh category --timeout=60s
   ```

4. **Coverage Issues:**
   ```bash
   # Generate detailed coverage report
   ./backend/scripts/run_test_categories.sh unit --coverage
   # Check coverage/unit_coverage.html
   ```

### Debug Test Execution

```bash
# Run with verbose output
./backend/scripts/run_test_categories.sh unit -v

# Run specific test file
go test -v ./internal/services/specific_service_test.go

# Run with race detection
go test -race -v ./internal/...

# Run with memory profiling
go test -memprofile=mem.prof -v ./internal/...
```

### Performance Test Debugging

```bash
# Run benchmarks
go test -bench=. -benchmem ./testing/performance/...

# Profile CPU usage
go test -cpuprofile=cpu.prof -bench=. ./testing/performance/...

# Analyze profiles
go tool pprof cpu.prof
```

## Test Maintenance

### Regular Tasks

1. **Update Test Dependencies:**
   ```bash
   cd backend
   go get -u ./...
   go mod tidy
   ```

2. **Review Test Coverage:**
   ```bash
   ./backend/scripts/run_test_categories.sh all --coverage
   # Review coverage reports
   ```

3. **Clean Up Obsolete Tests:**
   - Remove tests for deleted features
   - Update tests for changed APIs
   - Consolidate duplicate test logic

4. **Performance Baseline Updates:**
   - Update performance test expectations
   - Review benchmark results
   - Adjust timeout values if needed

### Best Practices for Maintenance

1. **Keep Tests Simple**: Each test should verify one specific behavior
2. **Avoid Test Dependencies**: Tests should be independent and runnable in any order
3. **Use Descriptive Names**: Test names should clearly describe what is being tested
4. **Regular Refactoring**: Keep test code clean and maintainable
5. **Documentation**: Document complex test setups and unusual requirements

This testing guide provides comprehensive instructions for working with the test suite. Keep it updated as new testing patterns and requirements emerge.