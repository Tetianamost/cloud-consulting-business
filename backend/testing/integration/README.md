# Integration Tests

This directory contains integration test executables that test component interactions, API endpoints, and external service integrations.

## Test Files

Integration tests will be moved here from the backend root directory in subsequent tasks:

- Chat system integration tests (`test_chat_*.go`)
- AI/Bedrock integration tests (`test_bedrock_*.go`, `test_ai_*.go`)
- API endpoint tests (`test_*_api_*.go`)
- Database integration tests
- External service integration tests

## Running Integration Tests

```bash
# Run specific integration test
cd backend && go run testing/integration/test_name.go

# Run all integration tests (when available)
cd backend/testing/integration && find . -name "*.go" -exec go run {} \;
```

## Requirements

- Test database connection
- External service credentials (AWS, etc.)
- Proper environment configuration