# Testing Framework

This directory contains the testing framework and utilities for the Cloud Consulting Platform backend.

## Structure

- `shared/` - Common test utilities, mocks, and fixtures
- `config/` - Environment-specific test configurations
- `integration/` - Integration test files
- `performance/` - Performance test files  
- `email/` - Email system test files

## Running Tests

### Using Test Categories Script (Recommended)

The project provides a convenient script for running different test categories:

```bash
# From backend directory
./scripts/run_test_categories.sh [CATEGORY] [OPTIONS]

# Show available test categories and files
./scripts/run_test_categories.sh list

# Run specific test categories
./scripts/run_test_categories.sh unit              # Unit tests only
./scripts/run_test_categories.sh integration      # Integration tests
./scripts/run_test_categories.sh email           # Email system tests
./scripts/run_test_categories.sh performance     # Performance tests
./scripts/run_test_categories.sh all             # All categories

# Run with options
./scripts/run_test_categories.sh unit -v          # Verbose output
./scripts/run_test_categories.sh email --coverage # With coverage report
./scripts/run_test_categories.sh all --timeout=60s # Custom timeout
```

### Manual Test Execution

```bash
# Run all tests
go test ./...

# Run specific test category
go test ./integration/...
go test ./performance/...
go test ./email/...

# Run with verbose output
go test -v ./...

# Run with coverage
go test -cover ./...

# Run tests for specific environment
TEST_ENVIRONMENT=unit go test ./shared/...
TEST_ENVIRONMENT=integration go test ./integration/...
```

## Test Configuration

### Environment-Based Configuration

The testing framework supports multiple environments with different configurations:

- `unit` - Fast unit tests with mocks and in-memory databases
- `integration` - Integration tests with real services
- `performance` - Performance and load testing
- `e2e` - End-to-end testing
- `ci` - Continuous integration environment
- `local` - Local development environment

### Configuration Files

Environment-specific configurations are stored in `config/`:

- `config/unit.json` - Unit test configuration
- `config/integration.json` - Integration test configuration
- `config/performance.json` - Performance test configuration
- `config/e2e.json` - End-to-end test configuration
- `config/ci.json` - CI environment configuration

### Environment Variables

Test configuration can be overridden with environment variables:

- `TEST_ENVIRONMENT` - Test environment (unit, integration, performance, e2e, ci, local)
- `TEST_DATABASE_URL` - Database connection string for tests
- `TEST_REDIS_URL` - Redis connection string for tests
- `TEST_LOG_LEVEL` - Log level for tests (debug, info, warn, error)
- `TEST_TIMEOUT` - Test timeout duration
- `TEST_PARALLEL` - Run tests in parallel (true/false)
- `TEST_SKIP_SLOW` - Skip slow tests (true/false)
- `TEST_VERBOSE` - Verbose output (true/false)
- `BEDROCK_USE_MOCK` - Use mock Bedrock service (true/false)
- `SES_USE_MOCK` - Use mock SES service (true/false)

## Shared Utilities

The `shared/` directory contains common utilities that can be used across all test categories:

### Test Data Builders and Fixtures

```go
import "testing/shared"

// Create test data builder
builder := shared.NewTestDataBuilder()

// Build test inquiry
inquiry := builder.BuildTestInquiry()

// Build test inquiry with custom options
inquiry := builder.BuildTestInquiryWithOptions(map[string]interface{}{
    "name": "Custom Name",
    "priority": "high",
})

// Load complete test fixtures
fixtures := shared.LoadTestFixtures()
```

### Mock Implementations

```go
import "testing/shared"

// Create mock services
mockBedrock := &shared.MockBedrockService{}
mockEmail := &shared.MockEmailService{}
mockInquiry := &shared.MockInquiryService{}

// Set up mock expectations
mockBedrock.On("GenerateText", mock.Anything, mock.Anything, mock.Anything).
    Return(&shared.TestBedrockResponse{Content: "Mock response"}, nil)
```

### HTTP Testing Helpers

```go
import "testing/shared"

// Create HTTP test helper
helper := shared.NewHTTPTestHelper()

// Make JSON request
response := helper.MakeJSONRequest(t, "POST", "/api/inquiries", requestBody)

// Assert JSON response
helper.AssertJSONResponse(t, response, 200, expectedResponse)

// Run test cases
helper.RunHTTPTestCases(t, []shared.HTTPTestCase{
    {
        Name: "successful creation",
        Request: shared.POST("/api/inquiries").Body(requestBody).Build(),
        ExpectedStatus: 201,
        ExpectedBody: expectedResponse,
    },
})
```

### Configuration Management

```go
import "testing/shared"

// Initialize global config manager
shared.InitGlobalConfigManager(logger)

// Get configuration for environment
config, err := shared.GetGlobalConfig(shared.TestEnvUnit)

// Detect current environment
env := shared.GetGlobalConfigUtils().DetectEnvironment()

// Check if should skip slow tests
if shared.ShouldSkipSlow() {
    t.Skip("Skipping slow test")
}
```

### Test Suites

```go
import "testing/shared"

func TestMyFeature(t *testing.T) {
    // Create test suite
    suite, err := shared.NewTestSuite("my-feature-tests")
    if err != nil {
        t.Skip("Database not available")
        return
    }
    defer suite.Cleanup()
    
    // Setup test
    suite.SetupTest(t)
    defer suite.TeardownTest(t)
    
    // Run test logic
    // ...
}
```

### Database Testing

```go
import "testing/shared"

// Create test database
config := shared.LoadSimpleTestConfig()
logger := shared.SetupSimpleTestLogger("error")
db, err := shared.NewTestDatabase(config, logger)
if err != nil {
    t.Skip("Database not available")
    return
}
defer db.Close()

// Clean database tables
db.TruncateAllTables(t)
```

### Test Utilities

```go
import "testing/shared"

// Create test context with timeout
ctx, cancel := shared.CreateTestContext(30 * time.Second)
defer cancel()

// Skip if running in short mode
shared.SkipIfShort(t, "requires external services")

// Assert eventually true
shared.AssertEventuallyTrue(t, func() bool {
    return someCondition()
}, 5*time.Second, "condition should become true")

// Test metrics
metrics := shared.NewTestMetrics()
defer func() {
    metrics.Finish()
    metrics.LogMetrics(logger, t.Name())
}()
```

## Best Practices

### Test Organization

1. **Use appropriate test environments**: Unit tests should use `TestEnvUnit`, integration tests should use `TestEnvIntegration`
2. **Leverage shared utilities**: Use the shared mocks, fixtures, and helpers to reduce code duplication
3. **Follow naming conventions**: Test files should be named `*_test.go`, test functions should start with `Test`
4. **Use table-driven tests**: For testing multiple scenarios, use table-driven test patterns

### Configuration Management

1. **Environment detection**: Let the framework auto-detect the environment when possible
2. **Override with environment variables**: Use environment variables to override configuration for specific test runs
3. **Use appropriate timeouts**: Set reasonable timeouts for different test environments
4. **Mock external services**: Use mocks for unit tests, real services for integration tests

### Performance Considerations

1. **Parallel execution**: Enable parallel execution for unit tests, disable for integration tests that share resources
2. **Skip slow tests**: Use `TEST_SKIP_SLOW=true` for quick feedback during development
3. **Clean up resources**: Always clean up database tables, close connections, and release resources
4. **Use in-memory databases**: Use SQLite in-memory databases for unit tests

### Error Handling

1. **Graceful degradation**: Tests should skip gracefully when external dependencies are not available
2. **Clear error messages**: Provide clear error messages that help diagnose test failures
3. **Proper cleanup**: Ensure cleanup happens even when tests fail
4. **Timeout handling**: Set appropriate timeouts and handle timeout errors gracefully

## Examples

See the test files in the `shared/` directory for examples of how to use the testing utilities:

- `shared/testutil_test.go` - Basic utility functions
- `shared/config_test.go` - Configuration management
- `shared/fixtures_test.go` - Test data builders and fixtures

## Adding New Test Files

### Guidelines for New Tests

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

### Step-by-Step Process

1. **Create the test file** in the appropriate directory
2. **Write the test** following project conventions
3. **Run the test** using the test categories script
4. **Verify** the test appears in the listing

Example for adding an integration test:
```bash
# 1. Create file
touch testing/integration/test_new_feature_integration.go

# 2. Write test (see examples in existing files)

# 3. Run test
./scripts/run_test_categories.sh integration -v

# 4. Verify listing
./scripts/run_test_categories.sh list
```

For comprehensive testing guidelines, see: **[TESTING_GUIDE.md](../TESTING_GUIDE.md)**

## Migration from Root Directory

Test files are being migrated from `backend/` root to this organized structure:

- `test_email_*.go` → `backend/testing/email/`
- `test_chat_*.go` → `backend/testing/integration/`
- `test_performance_*.go` → `backend/testing/performance/`
- `test_bedrock_*.go` → `backend/testing/integration/`
- `test_ai_*.go` → `backend/testing/integration/`

This migration will be handled in subsequent tasks to avoid breaking existing functionality.