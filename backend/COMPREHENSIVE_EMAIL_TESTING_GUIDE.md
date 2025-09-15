# Comprehensive Email Event Tracking System Testing Guide

This document provides a complete guide to the comprehensive test suite for the email event tracking system, covering all components from repository layer to frontend data adapters.

## Overview

The email event tracking system test suite provides comprehensive coverage for:

- **Email Event Repository** - Database operations, CRUD, metrics calculation
- **Email Event Recorder Service** - Async event recording, retry logic, health checks
- **Email Metrics Service** - Real-time metrics calculation, filtering, validation
- **Admin Handler Email Endpoints** - API endpoints for metrics, status, and history
- **Email Service Integration** - End-to-end email sending with event recording
- **Frontend Data Adapter** - Real data prioritization over mock data

## Test Files Structure

```
backend/
├── test_email_event_repository_comprehensive.go      # Repository unit tests
├── test_email_event_recorder_comprehensive.go        # Recorder service unit tests
├── test_email_metrics_service_comprehensive.go       # Metrics service unit tests
├── test_admin_handler_email_comprehensive.go         # Admin handler API tests
├── test_email_service_integration_comprehensive.go   # Integration tests
├── run_comprehensive_email_tests.sh                  # Test runner script
└── COMPREHENSIVE_EMAIL_TESTING_GUIDE.md             # This guide

frontend/src/components/admin/
└── V0DataAdapter.test.ts                             # Frontend adapter tests
```

## Running Tests

### Quick Start

```bash
# Run all tests
cd backend
./run_comprehensive_email_tests.sh

# Run specific test suite
./run_comprehensive_email_tests.sh repository
./run_comprehensive_email_tests.sh recorder
./run_comprehensive_email_tests.sh metrics
./run_comprehensive_email_tests.sh handler
./run_comprehensive_email_tests.sh integration

# Run with verbose output and coverage
./run_comprehensive_email_tests.sh -v -c

# Run with custom timeout
./run_comprehensive_email_tests.sh -t 600s integration
```

### Individual Test Execution

```bash
# Repository tests
go test -v ./test_email_event_repository_comprehensive.go

# Recorder service tests
go test -v ./test_email_event_recorder_comprehensive.go

# Metrics service tests
go test -v ./test_email_metrics_service_comprehensive.go

# Admin handler tests
go test -v ./test_admin_handler_email_comprehensive.go

# Integration tests
go test -v ./test_email_service_integration_comprehensive.go
```

### Frontend Tests

```bash
cd frontend
npm test -- V0DataAdapter.test.ts
```

## Test Coverage

### 1. Email Event Repository Tests (`test_email_event_repository_comprehensive.go`)

**Coverage**: Database operations, CRUD operations, metrics calculation, filtering

**Test Cases**:
- ✅ **Create Operations**
  - Valid event creation
  - Event with all fields populated
  - Duplicate ID handling
  - Timestamp auto-generation

- ✅ **Update Operations**
  - Valid status updates
  - Non-existent event handling
  - Failed status updates with error messages
  - Bounced status updates with bounce types

- ✅ **Query Operations**
  - Get events by inquiry ID
  - Get events by SES message ID
  - Multiple events ordering (most recent first)
  - Empty result handling

- ✅ **Metrics Calculation**
  - All-time metrics with different statuses
  - Filtering by email type
  - Filtering by status
  - Filtering by inquiry ID
  - Rate calculations (delivery, bounce, spam)
  - Empty time range handling

- ✅ **List Operations**
  - Pagination support
  - Filtering combinations
  - Time range filtering
  - Ordering verification

- ✅ **Error Handling**
  - Null event handling
  - Invalid filter parameters
  - Database connection errors

- ✅ **Edge Cases**
  - Empty optional fields
  - Very long field values
  - Concurrent operations

**Requirements Covered**: 4.1, 4.2, 4.3 (Database schema and operations)

### 2. Email Event Recorder Tests (`test_email_event_recorder_comprehensive.go`)

**Coverage**: Async event recording, retry logic, health checks, non-blocking behavior

**Test Cases**:
- ✅ **Record Email Sent**
  - Valid event recording
  - ID and timestamp generation
  - Existing ID preservation
  - Status defaulting
  - Invalid event validation

- ✅ **Update Email Status**
  - Valid status updates with delivery timestamps
  - Error message recording
  - Non-existent message ID handling
  - Empty message ID validation

- ✅ **Get Events by Inquiry**
  - Valid inquiry ID retrieval
  - Empty inquiry ID handling
  - No events found scenarios
  - Repository error handling

- ✅ **Synchronous Methods**
  - Sync event recording
  - Sync status updates
  - Error propagation in sync mode

- ✅ **Non-blocking Behavior**
  - Async operation verification
  - Performance timing tests
  - Multiple concurrent operations

- ✅ **Retry Logic**
  - Transient failure recovery
  - Exponential backoff verification
  - Max retries exceeded handling

- ✅ **Health Checks**
  - Service health verification
  - Context-based health checks
  - Connection error detection

- ✅ **Concurrent Operations**
  - Thread safety verification
  - Race condition prevention

**Requirements Covered**: 2.1, 2.2, 2.3, 7.1 (Event recording and error handling)

### 3. Email Metrics Service Tests (`test_email_metrics_service_comprehensive.go`)

**Coverage**: Real-time metrics calculation, filtering, validation, extended methods

**Test Cases**:
- ✅ **Get Email Metrics**
  - All-time metrics calculation
  - Empty time range handling
  - Repository error handling
  - Rate calculations verification

- ✅ **Get Email Status by Inquiry**
  - Multiple email types handling
  - No emails found scenarios
  - Failed email status tracking
  - Most recent email identification

- ✅ **Get Email Event History**
  - Valid filter application
  - Email type filtering
  - Status filtering
  - Pagination support
  - Repository error handling

- ✅ **Filter Validation**
  - Invalid time ranges
  - Future time range limits
  - Negative limits and offsets
  - Excessive limits
  - Invalid email types and statuses

- ✅ **Extended Methods**
  - Metrics by email type
  - Recent email activity
  - Invalid hour ranges

- ✅ **Error Handling**
  - Repository connection failures
  - Service timeout handling
  - Invalid parameter handling

- ✅ **Edge Cases**
  - Empty inquiry IDs
  - Zero limit filters
  - Duplicate email types per inquiry

- ✅ **Health Checks**
  - Service health verification
  - Repository connectivity testing

**Requirements Covered**: 3.1, 3.2, 3.3 (Email metrics and calculations)

### 4. Admin Handler Email Tests (`test_admin_handler_email_comprehensive.go`)

**Coverage**: API endpoints, real data prioritization, error responses, health checks

**Test Cases**:
- ✅ **Get System Metrics**
  - Real email metrics integration
  - Unhealthy service handling
  - Missing service configuration
  - Invalid time range parameters
  - Service error handling

- ✅ **Get Email Status**
  - Successful status retrieval
  - Inquiry not found handling
  - No email events scenarios
  - Service unavailability
  - Service unhealthy states
  - Missing inquiry ID validation

- ✅ **Get Email Event History**
  - Valid filter processing
  - Email type filtering
  - Status filtering
  - Pagination implementation
  - Invalid parameter handling
  - Service availability checks

- ✅ **Error Handling**
  - Service timeout scenarios
  - Repository connection failures
  - Invalid API parameters
  - Missing dependencies

- ✅ **Edge Cases**
  - Large offset handling
  - Maximum limit capping
  - Concurrent request handling

**Requirements Covered**: 3.1, 3.2, 3.3, 6.4, 7.2, 7.3, 7.4 (API endpoints and error handling)

### 5. Email Service Integration Tests (`test_email_service_integration_comprehensive.go`)

**Coverage**: End-to-end email sending with event recording, template integration, SES integration

**Test Cases**:
- ✅ **Send Customer Confirmation**
  - Successful email with event recording
  - Template rendering failures
  - SES failures with event recording
  - Event recording failures (non-blocking)

- ✅ **Send Report Email**
  - Successful consultant notifications
  - High priority detection
  - SES message ID handling
  - Template data preparation

- ✅ **Send Inquiry Notification**
  - Successful notifications
  - High priority inquiry detection
  - Template integration

- ✅ **Error Handling**
  - Invalid inquiry data
  - Template service failures
  - SES service failures
  - Context cancellation
  - Context timeouts

- ✅ **Event Recording Failures**
  - Non-blocking email delivery
  - Partial failure handling
  - Invalid event data handling

- ✅ **Health Checks**
  - Service health verification
  - Dependency validation
  - Configuration validation

**Requirements Covered**: 2.1, 2.2, 2.3, 6.1, 6.2, 6.3, 7.1 (Email service integration)

### 6. Frontend Data Adapter Tests (`V0DataAdapter.test.ts`)

**Coverage**: Real data prioritization, mock data replacement, error handling

**Test Cases**:
- ✅ **Safe Adapt Email Metrics**
  - Null data handling
  - Real system metrics adaptation
  - Email status prioritization
  - Bounce/spam categorization
  - Mixed status handling
  - Open/click rate estimation
  - Error handling with null return

- ✅ **Has Real Email Data**
  - Data availability detection
  - System metrics validation
  - Email status validation
  - Partial data scenarios

- ✅ **Get Email Data Error Message**
  - Specific error code handling
  - Generic error messages
  - No data scenarios
  - Partial data availability

- ✅ **Adapt Email Statuses to Metrics**
  - Empty array handling
  - Null input handling
  - Metrics calculation
  - Error categorization

- ✅ **Integration with Real API Responses**
  - Typical successful responses
  - API error scenarios
  - Partial API failures

**Requirements Covered**: 6.1, 6.2, 6.3, 7.2, 7.3 (Frontend integration and error handling)

## Test Environment Setup

### Prerequisites

1. **Go 1.24+** - Required for backend tests
2. **PostgreSQL** - Optional for integration tests (tests will skip if unavailable)
3. **Node.js 18+** - Required for frontend tests
4. **Dependencies** - Automatically downloaded by test runner

### Environment Variables

```bash
# Test database (optional - tests will use in-memory if unavailable)
export TEST_DATABASE_URL="postgres://test:test@localhost/test_email_events?sslmode=disable"

# Email configuration for tests
export SES_SENDER_EMAIL="info@cloudpartner.pro"
export SES_REPLY_TO_EMAIL="info@cloudpartner.pro"

# Test configuration
export LOG_LEVEL="error"
export VERBOSE="false"
export COVERAGE="true"
export TIMEOUT="300s"
```

### Database Setup (Optional)

```sql
-- Create test database
CREATE DATABASE test_email_events;
CREATE USER test WITH PASSWORD 'test';
GRANT ALL PRIVILEGES ON DATABASE test_email_events TO test;
```

## Test Patterns and Best Practices

### 1. Mock Usage

All tests use comprehensive mocks to isolate components:

```go
// Repository mocks for service tests
type MockEmailEventRepository struct {
    mock.Mock
    events map[string]*domain.EmailEvent
    mutex  sync.RWMutex
}

// Service mocks for handler tests
type MockEmailMetricsService struct {
    mock.Mock
}
```

### 2. Test Data Creation

Consistent test data creation with helpers:

```go
func createTestEmailEvent() *domain.EmailEvent {
    return &domain.EmailEvent{
        ID:             uuid.New().String(),
        InquiryID:      "inquiry-" + uuid.New().String(),
        EmailType:      domain.EmailTypeCustomerConfirmation,
        RecipientEmail: "test@example.com",
        SenderEmail:    "info@cloudpartner.pro",
        Status:         domain.EmailStatusSent,
        SentAt:         time.Now(),
    }
}
```

### 3. Error Scenario Testing

Comprehensive error handling verification:

```go
t.Run("RepositoryError", func(t *testing.T) {
    mockRepo.On("GetMetrics", ctx, mock.Anything).Return(nil, fmt.Errorf("database error")).Once()
    
    metrics, err := service.GetEmailMetrics(ctx, timeRange)
    assert.Error(t, err)
    assert.Nil(t, metrics)
    assert.Contains(t, err.Error(), "database error")
})
```

### 4. Async Operation Testing

Non-blocking behavior verification:

```go
t.Run("NonBlockingBehavior", func(t *testing.T) {
    start := time.Now()
    
    for i := 0; i < 5; i++ {
        err := recorder.RecordEmailSent(ctx, event)
        assert.NoError(t, err)
    }
    
    elapsed := time.Since(start)
    assert.Less(t, elapsed, 100*time.Millisecond, "Should be non-blocking")
})
```

### 5. Integration Testing

End-to-end workflow verification:

```go
t.Run("SuccessfulEmailWithEventRecording", func(t *testing.T) {
    // Setup all mocks in chain
    mockTemplate.On("PrepareCustomerConfirmationData", inquiry).Return(templateData)
    mockSES.On("SendEmail", ctx, mock.Anything).Return(nil)
    mockEventRecorder.On("RecordEmailSent", ctx, mock.Anything).Return(nil)
    
    // Execute full workflow
    err := emailService.SendCustomerConfirmation(ctx, inquiry)
    
    // Verify all components called
    assert.NoError(t, err)
    mockTemplate.AssertExpectations(t)
    mockSES.AssertExpectations(t)
    mockEventRecorder.AssertExpectations(t)
})
```

## Coverage Goals

- **Unit Tests**: >90% code coverage
- **Integration Tests**: >80% workflow coverage
- **Error Scenarios**: 100% error path coverage
- **Edge Cases**: Comprehensive boundary testing

## Continuous Integration

### GitHub Actions Integration

```yaml
name: Email Event Tracking Tests
on: [push, pull_request]

jobs:
  backend-tests:
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:13
        env:
          POSTGRES_PASSWORD: test
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5

    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-go@v2
        with:
          go-version: 1.24
      
      - name: Run comprehensive email tests
        run: |
          cd backend
          ./run_comprehensive_email_tests.sh -v -c
      
      - name: Upload coverage
        uses: codecov/codecov-action@v1
        with:
          file: ./backend/combined_coverage.out

  frontend-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-node@v2
        with:
          node-version: '18'
      
      - name: Install dependencies
        run: |
          cd frontend
          npm ci
      
      - name: Run frontend tests
        run: |
          cd frontend
          npm test -- V0DataAdapter.test.ts --coverage
```

## Troubleshooting

### Common Issues

1. **Database Connection Failures**
   ```bash
   # Tests will skip database tests if PostgreSQL unavailable
   # Set TEST_DATABASE_URL to use specific database
   export TEST_DATABASE_URL="postgres://user:pass@localhost/testdb"
   ```

2. **Timeout Issues**
   ```bash
   # Increase timeout for slow systems
   ./run_comprehensive_email_tests.sh -t 600s
   ```

3. **Coverage Report Issues**
   ```bash
   # Disable coverage if causing issues
   ./run_comprehensive_email_tests.sh --no-coverage
   ```

4. **Mock Expectation Failures**
   ```bash
   # Run with verbose output to see detailed mock failures
   ./run_comprehensive_email_tests.sh -v
   ```

### Debug Mode

```bash
# Enable debug logging
export LOG_LEVEL="debug"
export VERBOSE="true"

# Run specific failing test
go test -v -run TestSpecificFailingTest ./test_file.go
```

## Performance Benchmarks

### Expected Performance

- **Repository Operations**: <10ms per operation
- **Service Operations**: <50ms per operation
- **API Endpoints**: <200ms response time
- **Event Recording**: <1ms (non-blocking)

### Benchmark Tests

```bash
# Run performance benchmarks
go test -bench=. -benchmem ./test_email_event_repository_comprehensive.go
```

## Maintenance

### Adding New Tests

1. **Follow Naming Convention**: `test_component_comprehensive.go`
2. **Use Established Patterns**: Mock setup, helper functions, error scenarios
3. **Update Test Runner**: Add new test to `run_comprehensive_email_tests.sh`
4. **Document Coverage**: Update this guide with new test cases

### Updating Existing Tests

1. **Maintain Backward Compatibility**: Don't break existing test patterns
2. **Update Documentation**: Reflect changes in this guide
3. **Verify Coverage**: Ensure coverage goals are maintained
4. **Test Dependencies**: Update mocks if interfaces change

## Conclusion

This comprehensive test suite ensures the email event tracking system meets all requirements with high reliability, performance, and maintainability. The tests cover all layers from database operations to frontend data handling, providing confidence in the system's ability to replace mock data with real email metrics.

For questions or issues with the test suite, refer to the troubleshooting section or check the individual test files for specific implementation details.