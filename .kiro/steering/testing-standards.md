# Testing Standards and Guidelines

## Overview

This document establishes comprehensive testing standards for the Cloud Consulting Platform, covering unit testing, integration testing, end-to-end testing, and performance testing for both Go backend and React frontend components.

## Testing Philosophy

### Core Principles
1. **Test-Driven Development (TDD)**: Write tests before implementation when possible
2. **Comprehensive Coverage**: Aim for >90% code coverage on new features
3. **Fast Feedback**: Tests should run quickly and provide clear failure messages
4. **Reliable Tests**: Tests should be deterministic and not flaky
5. **Maintainable Tests**: Tests should be easy to understand and modify

### Testing Pyramid
1. **Unit Tests (70%)**: Fast, isolated tests for individual functions/components
2. **Integration Tests (20%)**: Tests for component interactions and API endpoints
3. **End-to-End Tests (10%)**: Full user workflow tests

## Backend Testing Standards (Go)

### Unit Testing Patterns

#### Service Layer Testing
```go
func TestChatService_SendMessage(t *testing.T) {
    // Arrange
    mockRepo := &MockChatRepository{}
    mockAI := &MockAIService{}
    service := NewChatService(mockRepo, mockAI)
    
    ctx := context.Background()
    request := &SendMessageRequest{
        SessionID: "test-session",
        Message:   "Hello, AI",
        UserID:    "user-123",
    }
    
    // Setup mocks
    mockRepo.On("SaveMessage", mock.Anything).Return(nil)
    mockAI.On("GenerateResponse", mock.Anything).Return(&AIResponse{
        Content: "Hello! How can I help you?",
    }, nil)
    
    // Act
    response, err := service.SendMessage(ctx, request)
    
    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, response)
    assert.Equal(t, "Hello! How can I help you?", response.AIResponse.Content)
    mockRepo.AssertExpectations(t)
    mockAI.AssertExpectations(t)
}
```

#### Repository Testing
```go
func TestChatRepository_SaveMessage(t *testing.T) {
    // Use test database or in-memory storage
    db := setupTestDB(t)
    defer cleanupTestDB(t, db)
    
    repo := NewChatRepository(db)
    
    message := &ChatMessage{
        ID:        "msg-123",
        SessionID: "session-123",
        Content:   "Test message",
        Type:      "user",
        CreatedAt: time.Now(),
    }
    
    err := repo.SaveMessage(context.Background(), message)
    
    assert.NoError(t, err)
    
    // Verify message was saved
    saved, err := repo.GetMessage(context.Background(), "msg-123")
    assert.NoError(t, err)
    assert.Equal(t, message.Content, saved.Content)
}
```

#### Handler Testing
```go
func TestChatHandler_SendMessage(t *testing.T) {
    // Setup
    mockService := &MockChatService{}
    handler := NewChatHandler(mockService)
    
    tests := []struct {
        name           string
        requestBody    string
        setupMocks     func()
        expectedStatus int
        expectedBody   string
    }{
        {
            name: "successful message send",
            requestBody: `{
                "session_id": "test-session",
                "message": "Hello",
                "user_id": "user-123"
            }`,
            setupMocks: func() {
                mockService.On("SendMessage", mock.Anything, mock.Anything).Return(&SendMessageResponse{
                    MessageID: "msg-123",
                    AIResponse: &AIResponse{Content: "Hi there!"},
                }, nil)
            },
            expectedStatus: http.StatusOK,
            expectedBody:   `{"message_id":"msg-123","ai_response":{"content":"Hi there!"}}`,
        },
        {
            name:           "invalid request body",
            requestBody:    `{"invalid": "json"}`,
            setupMocks:     func() {},
            expectedStatus: http.StatusBadRequest,
            expectedBody:   `{"error":"Invalid request format"}`,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            tt.setupMocks()
            
            req := httptest.NewRequest("POST", "/api/chat/send", strings.NewReader(tt.requestBody))
            req.Header.Set("Content-Type", "application/json")
            w := httptest.NewRecorder()
            
            handler.SendMessage(w, req)
            
            assert.Equal(t, tt.expectedStatus, w.Code)
            assert.JSONEq(t, tt.expectedBody, w.Body.String())
        })
    }
}
```

### Integration Testing

#### API Integration Tests
```go
func TestChatAPI_Integration(t *testing.T) {
    // Setup test server
    server := setupTestServer(t)
    defer server.Close()
    
    client := &http.Client{Timeout: 10 * time.Second}
    
    // Test complete chat flow
    t.Run("complete chat flow", func(t *testing.T) {
        // 1. Create session
        sessionResp := createTestSession(t, client, server.URL)
        sessionID := sessionResp.SessionID
        
        // 2. Send message
        messageResp := sendTestMessage(t, client, server.URL, sessionID, "Hello")
        assert.NotEmpty(t, messageResp.MessageID)
        assert.NotEmpty(t, messageResp.AIResponse.Content)
        
        // 3. Get message history
        history := getMessageHistory(t, client, server.URL, sessionID)
        assert.Len(t, history.Messages, 2) // User message + AI response
    })
}
```

#### Database Integration Tests
```go
func TestDatabaseIntegration(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping database integration test")
    }
    
    db := setupTestDatabase(t)
    defer cleanupTestDatabase(t, db)
    
    // Test repository operations
    repo := NewChatRepository(db)
    
    // Test transaction handling
    err := repo.WithTransaction(func(tx *sql.Tx) error {
        // Perform multiple operations
        return nil
    })
    
    assert.NoError(t, err)
}
```

### Performance Testing

#### Load Testing
```go
func TestChatHandler_Load(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping load test")
    }
    
    handler := setupTestHandler(t)
    
    // Simulate concurrent requests
    concurrency := 100
    requests := 1000
    
    var wg sync.WaitGroup
    errors := make(chan error, requests)
    
    for i := 0; i < concurrency; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            
            for j := 0; j < requests/concurrency; j++ {
                err := sendTestRequest(handler)
                if err != nil {
                    errors <- err
                }
            }
        }()
    }
    
    wg.Wait()
    close(errors)
    
    errorCount := 0
    for err := range errors {
        t.Logf("Request error: %v", err)
        errorCount++
    }
    
    // Assert acceptable error rate (< 1%)
    assert.Less(t, errorCount, requests/100)
}
```

#### Benchmark Tests
```go
func BenchmarkChatService_SendMessage(b *testing.B) {
    service := setupBenchmarkService(b)
    request := &SendMessageRequest{
        SessionID: "bench-session",
        Message:   "Benchmark message",
        UserID:    "bench-user",
    }
    
    b.ResetTimer()
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            _, err := service.SendMessage(context.Background(), request)
            if err != nil {
                b.Fatal(err)
            }
        }
    })
}
```

## Frontend Testing Standards (React)

### Component Unit Testing

#### Basic Component Testing
```tsx
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { ChatMessage } from './ChatMessage';

describe('ChatMessage', () => {
  const mockMessage = {
    id: 'msg-1',
    content: 'Hello, world!',
    type: 'user' as const,
    createdAt: new Date('2024-01-01T10:00:00Z'),
    status: 'delivered' as const,
  };

  it('renders user message correctly', () => {
    render(<ChatMessage message={mockMessage} isOwn={true} />);
    
    expect(screen.getByText('Hello, world!')).toBeInTheDocument();
    expect(screen.getByText('10:00 AM')).toBeInTheDocument();
  });

  it('shows retry button for failed messages', () => {
    const failedMessage = { ...mockMessage, status: 'failed' as const };
    const mockRetry = jest.fn();
    
    render(<ChatMessage message={failedMessage} isOwn={true} onRetry={mockRetry} />);
    
    const retryButton = screen.getByRole('button', { name: /retry/i });
    fireEvent.click(retryButton);
    
    expect(mockRetry).toHaveBeenCalledTimes(1);
  });
});
```

#### Hook Testing
```tsx
import { renderHook, act } from '@testing-library/react';
import { useChat } from './useChat';

describe('useChat', () => {
  it('should send message and update state', async () => {
    const { result } = renderHook(() => useChat('session-123'));
    
    act(() => {
      result.current.sendMessage('Hello');
    });
    
    expect(result.current.messages).toHaveLength(1);
    expect(result.current.messages[0].content).toBe('Hello');
    expect(result.current.messages[0].status).toBe('sending');
    
    // Wait for AI response
    await waitFor(() => {
      expect(result.current.messages).toHaveLength(2);
    });
    
    expect(result.current.messages[1].type).toBe('ai');
  });
});
```

#### Redux Testing
```tsx
import { configureStore } from '@reduxjs/toolkit';
import chatReducer, { messageReceived, messageSent } from './chatSlice';

describe('chatSlice', () => {
  let store: ReturnType<typeof configureStore>;

  beforeEach(() => {
    store = configureStore({
      reducer: {
        chat: chatReducer,
      },
    });
  });

  it('should handle message sent', () => {
    const message = {
      id: 'msg-1',
      content: 'Hello',
      type: 'user' as const,
      createdAt: new Date(),
      status: 'sending' as const,
    };

    store.dispatch(messageSent(message));

    const state = store.getState().chat;
    expect(state.messages).toHaveLength(1);
    expect(state.messages[0]).toEqual(message);
  });
});
```

### Integration Testing

#### Component Integration
```tsx
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { Provider } from 'react-redux';
import { ChatContainer } from './ChatContainer';
import { store } from '../store';

// Mock API calls
jest.mock('../services/chatService', () => ({
  sendMessage: jest.fn().mockResolvedValue({
    messageId: 'msg-1',
    aiResponse: { content: 'AI response' },
  }),
}));

describe('ChatContainer Integration', () => {
  it('should handle complete message flow', async () => {
    render(
      <Provider store={store}>
        <ChatContainer sessionId="test-session" />
      </Provider>
    );

    const input = screen.getByPlaceholderText('Type a message...');
    const sendButton = screen.getByRole('button', { name: 'Send' });

    // Send message
    fireEvent.change(input, { target: { value: 'Hello AI' } });
    fireEvent.click(sendButton);

    // Check user message appears
    expect(screen.getByText('Hello AI')).toBeInTheDocument();

    // Wait for AI response
    await waitFor(() => {
      expect(screen.getByText('AI response')).toBeInTheDocument();
    });
  });
});
```

### End-to-End Testing (Cypress)

#### Chat Flow E2E Test
```typescript
// cypress/e2e/chat-flow.cy.ts
describe('Chat Flow', () => {
  beforeEach(() => {
    cy.login('admin', 'password');
    cy.visit('/admin/dashboard');
  });

  it('should complete full chat conversation', () => {
    // Open chat widget
    cy.get('[data-testid="chat-widget-toggle"]').click();
    
    // Send message
    cy.get('[data-testid="chat-input"]').type('Hello, I need help with AWS migration');
    cy.get('[data-testid="send-button"]').click();
    
    // Verify message appears
    cy.get('[data-testid="chat-messages"]')
      .should('contain', 'Hello, I need help with AWS migration');
    
    // Wait for AI response
    cy.get('[data-testid="ai-message"]', { timeout: 10000 })
      .should('be.visible')
      .and('contain.text', 'migration');
    
    // Send follow-up message
    cy.get('[data-testid="chat-input"]').type('What about costs?');
    cy.get('[data-testid="send-button"]').click();
    
    // Verify conversation continues
    cy.get('[data-testid="chat-messages"]')
      .should('contain', 'What about costs?');
  });

  it('should handle connection failures gracefully', () => {
    // Simulate network failure
    cy.intercept('POST', '/api/chat/send', { forceNetworkError: true });
    
    cy.get('[data-testid="chat-widget-toggle"]').click();
    cy.get('[data-testid="chat-input"]').type('Test message');
    cy.get('[data-testid="send-button"]').click();
    
    // Should show error state
    cy.get('[data-testid="message-error"]').should('be.visible');
    cy.get('[data-testid="retry-button"]').should('be.visible');
  });
});
```

## Test Data Management

### Test Fixtures
```go
// testdata/fixtures.go
type TestFixtures struct {
    Users    []*User
    Sessions []*ChatSession
    Messages []*ChatMessage
}

func LoadTestFixtures() *TestFixtures {
    return &TestFixtures{
        Users: []*User{
            {ID: "user-1", Email: "test@example.com", Role: "admin"},
            {ID: "user-2", Email: "user@example.com", Role: "user"},
        },
        Sessions: []*ChatSession{
            {ID: "session-1", UserID: "user-1", CreatedAt: time.Now()},
        },
        Messages: []*ChatMessage{
            {ID: "msg-1", SessionID: "session-1", Content: "Hello", Type: "user"},
        },
    }
}
```

### Database Seeding
```go
func setupTestDB(t *testing.T) *sql.DB {
    db, err := sql.Open("postgres", testDatabaseURL)
    require.NoError(t, err)
    
    // Run migrations
    err = runMigrations(db)
    require.NoError(t, err)
    
    // Seed test data
    fixtures := LoadTestFixtures()
    err = seedTestData(db, fixtures)
    require.NoError(t, err)
    
    return db
}

func cleanupTestDB(t *testing.T, db *sql.DB) {
    _, err := db.Exec("TRUNCATE TABLE chat_messages, chat_sessions, users CASCADE")
    require.NoError(t, err)
    
    db.Close()
}
```

## Mock and Stub Patterns

### Service Mocks
```go
type MockChatService struct {
    mock.Mock
}

func (m *MockChatService) SendMessage(ctx context.Context, req *SendMessageRequest) (*SendMessageResponse, error) {
    args := m.Called(ctx, req)
    return args.Get(0).(*SendMessageResponse), args.Error(1)
}

func (m *MockChatService) GetHistory(ctx context.Context, sessionID string) ([]*ChatMessage, error) {
    args := m.Called(ctx, sessionID)
    return args.Get(0).([]*ChatMessage), args.Error(1)
}
```

### HTTP Mocks
```go
func setupMockServer(t *testing.T) *httptest.Server {
    mux := http.NewServeMux()
    
    mux.HandleFunc("/api/bedrock/generate", func(w http.ResponseWriter, r *http.Request) {
        response := map[string]interface{}{
            "content": "Mocked AI response",
            "confidence": 0.95,
        }
        
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(response)
    })
    
    return httptest.NewServer(mux)
}
```

## Test Configuration

### Environment Setup
```go
// config/test.go
func LoadTestConfig() *Config {
    return &Config{
        Database: DatabaseConfig{
            URL: getEnvOrDefault("TEST_DATABASE_URL", "postgres://test:test@localhost/test_db"),
        },
        Redis: RedisConfig{
            URL: getEnvOrDefault("TEST_REDIS_URL", "redis://localhost:6379/1"),
        },
        Bedrock: BedrockConfig{
            BaseURL: "http://localhost:8080", // Mock server
        },
    }
}
```

### Test Utilities
```go
// testutil/helpers.go
func CreateTestUser(t *testing.T, db *sql.DB) *User {
    user := &User{
        ID:    uuid.New().String(),
        Email: fmt.Sprintf("test-%d@example.com", time.Now().Unix()),
        Role:  "admin",
    }
    
    err := saveUser(db, user)
    require.NoError(t, err)
    
    return user
}

func CreateTestSession(t *testing.T, db *sql.DB, userID string) *ChatSession {
    session := &ChatSession{
        ID:        uuid.New().String(),
        UserID:    userID,
        CreatedAt: time.Now(),
        ExpiresAt: time.Now().Add(30 * time.Minute),
    }
    
    err := saveSession(db, session)
    require.NoError(t, err)
    
    return session
}
```

## Continuous Integration

### Test Pipeline
```yaml
# .github/workflows/test.yml
name: Test Suite

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
      redis:
        image: redis:6
        options: >-
          --health-cmd "redis-cli ping"
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5

    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-go@v2
        with:
          go-version: 1.24
      
      - name: Run unit tests
        run: |
          cd backend
          go test -v -race -coverprofile=coverage.out ./...
      
      - name: Run integration tests
        run: |
          cd backend
          go test -v -tags=integration ./...
      
      - name: Upload coverage
        uses: codecov/codecov-action@v1
        with:
          file: ./backend/coverage.out

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
      
      - name: Run unit tests
        run: |
          cd frontend
          npm test -- --coverage --watchAll=false
      
      - name: Run E2E tests
        run: |
          cd frontend
          npm run test:e2e
```

## Quality Gates

### Coverage Requirements
- Unit tests: >90% coverage
- Integration tests: >80% coverage
- E2E tests: Cover all critical user paths

### Performance Benchmarks
- API response time: <200ms (95th percentile)
- Frontend render time: <100ms
- Database query time: <50ms

### Code Quality
- All tests must pass
- No flaky tests allowed
- Test code must follow same quality standards as production code
- Tests must be maintainable and readable