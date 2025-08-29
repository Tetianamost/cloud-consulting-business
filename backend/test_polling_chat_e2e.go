package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cloud-consulting/backend/internal/domain"
	"github.com/cloud-consulting/backend/internal/handlers"
	"github.com/cloud-consulting/backend/internal/interfaces"
	"github.com/cloud-consulting/backend/internal/services"
	"github.com/cloud-consulting/backend/internal/storage"
)

// E2E Test Suite for Polling Chat System
type PollingChatE2ETestSuite struct {
	server    *httptest.Server
	client    *http.Client
	authToken string
	sessionID string
	logger    *logrus.Logger
	cleanup   func()
}

// Test configuration for E2E tests
type E2ETestConfig struct {
	ServerURL       string
	AuthToken       string
	SessionID       string
	ClientName      string
	TestDuration    time.Duration
	MessageCount    int
	ConcurrentUsers int
}

// Message for E2E testing
type E2ETestMessage struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
	SessionID string    `json:"session_id"`
	Status    string    `json:"status,omitempty"`
}

// Setup E2E test environment
func setupE2ETestSuite(t *testing.T) *PollingChatE2ETestSuite {
	// Create test logger
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel) // Reduce noise in tests

	// Setup in-memory storage for testing
	memoryStorage := storage.NewMemoryStorage()

	// Create services
	chatMessageRepo := storage.NewMemoryChatMessageRepository(memoryStorage)
	chatSessionRepo := storage.NewMemoryChatSessionRepository(memoryStorage)
	sessionService := services.NewSessionService(chatSessionRepo, logger)

	// Mock Bedrock service for testing
	mockBedrockService := &MockBedrockService{
		responses: map[string]string{
			"hello": "Hello! How can I help you today?",
			"test":  "This is a test response from the AI assistant.",
			"help":  "I'm here to help with your cloud consulting needs.",
		},
	}

	chatService := services.NewChatService(chatMessageRepo, sessionService, mockBedrockService, logger)

	// Create auth services (simplified for testing)
	authService := &MockAuthService{
		validTokens: map[string]*interfaces.AuthContext{
			"test-token": {
				UserID:    "test-user-id",
				Username:  "testuser",
				Roles:     []string{"admin"},
				IssuedAt:  time.Now(),
				ExpiresAt: time.Now().Add(time.Hour),
			},
		},
	}

	securityService := &MockSecurityService{}
	rateLimiter := &MockRateLimiter{}
	auditLogger := &MockAuditLogger{}

	// Create handlers
	authHandler := &handlers.AuthHandler{}
	pollingChatHandler := handlers.NewPollingChatHandler(
		logger,
		chatService,
		sessionService,
		authHandler,
		authService,
		securityService,
		rateLimiter,
		auditLogger,
	)

	// Setup Gin router
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Add middleware
	router.Use(gin.Recovery())

	// Setup routes
	api := router.Group("/api/v1")
	adminAPI := api.Group("/admin")

	// Add auth middleware
	adminAPI.Use(pollingChatHandler.AuthMiddleware)

	// Chat routes
	chatAPI := adminAPI.Group("/chat")
	chatAPI.POST("/messages", pollingChatHandler.SendMessage)
	chatAPI.GET("/messages", pollingChatHandler.GetMessages)

	// Create test server
	testServer := httptest.NewServer(router)

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Generate test session ID
	sessionID := fmt.Sprintf("e2e-test-session-%d", time.Now().Unix())

	suite := &PollingChatE2ETestSuite{
		server:    testServer,
		client:    client,
		authToken: "test-token",
		sessionID: sessionID,
		logger:    logger,
		cleanup: func() {
			testServer.Close()
		},
	}

	return suite
}

// Mock services for E2E testing
type MockBedrockService struct {
	responses map[string]string
	mutex     sync.RWMutex
}

func (m *MockBedrockService) GenerateText(ctx context.Context, prompt string, options *interfaces.BedrockOptions) (*interfaces.BedrockResponse, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	// Simple keyword-based response generation
	prompt = strings.ToLower(prompt)
	for keyword, response := range m.responses {
		if strings.Contains(prompt, keyword) {
			return &interfaces.BedrockResponse{
				Content: response,
				Usage: interfaces.BedrockUsage{
					InputTokens:  len(strings.Split(prompt, " ")),
					OutputTokens: len(strings.Split(response, " ")),
				},
			}, nil
		}
	}

	// Default response
	return &interfaces.BedrockResponse{
		Content: "I understand your message. How can I assist you further?",
		Usage: interfaces.BedrockUsage{
			InputTokens:  len(strings.Split(prompt, " ")),
			OutputTokens: 10,
		},
	}, nil
}

func (m *MockBedrockService) GetModelInfo() interfaces.BedrockModelInfo {
	return interfaces.BedrockModelInfo{
		ModelID:    "test-model",
		ModelName:  "Test Model",
		MaxTokens:  4096,
		InputCost:  0.001,
		OutputCost: 0.002,
	}
}

func (m *MockBedrockService) IsHealthy() bool {
	return true
}

type MockAuthService struct {
	validTokens map[string]*interfaces.AuthContext
	mutex       sync.RWMutex
}

func (m *MockAuthService) ValidateToken(ctx context.Context, token string) (*interfaces.AuthContext, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if authContext, exists := m.validTokens[token]; exists {
		return authContext, nil
	}
	return nil, fmt.Errorf("invalid token")
}

func (m *MockAuthService) GenerateToken(ctx context.Context, userID string, duration time.Duration) (string, error) {
	return "generated-token", nil
}

func (m *MockAuthService) RefreshToken(ctx context.Context, token string) (string, error) {
	return "refreshed-token", nil
}

func (m *MockAuthService) RevokeToken(ctx context.Context, token string) error {
	return nil
}

type MockSecurityService struct{}

func (m *MockSecurityService) CheckRateLimit(ctx context.Context, userID string, action string) (*interfaces.RateLimitResult, error) {
	return &interfaces.RateLimitResult{Allowed: true}, nil
}

func (m *MockSecurityService) ValidateMessageContent(content string) error {
	return nil
}

func (m *MockSecurityService) SanitizeInput(input string) string {
	return input
}

func (m *MockSecurityService) DetectSpam(content string) (bool, float64) {
	return false, 0.0
}

type MockRateLimiter struct{}

func (m *MockRateLimiter) CheckLimit(ctx context.Context, key string, limit int, window time.Duration) (bool, error) {
	return true, nil
}

func (m *MockRateLimiter) GetRemainingRequests(ctx context.Context, key string, limit int, window time.Duration) (int, error) {
	return limit, nil
}

func (m *MockRateLimiter) ResetLimit(ctx context.Context, key string) error {
	return nil
}

type MockAuditLogger struct{}

func (m *MockAuditLogger) LogMessage(ctx context.Context, userID string, sessionID string, action string, metadata map[string]interface{}) {
	// No-op for testing
}

func (m *MockAuditLogger) LogLogin(ctx context.Context, userID string, success bool, metadata map[string]interface{}) {
	// No-op for testing
}

func (m *MockAuditLogger) LogRateLimitExceeded(ctx context.Context, userID string, action string, metadata map[string]interface{}) {
	// No-op for testing
}

func (m *MockAuditLogger) LogSecurityEvent(ctx context.Context, userID string, event string, metadata map[string]interface{}) {
	// No-op for testing
}

// Helper methods for E2E testing
func (suite *PollingChatE2ETestSuite) sendMessage(content string) (*handlers.SendMessageResponse, error) {
	request := handlers.SendMessageRequest{
		Content:    content,
		SessionID:  suite.sessionID,
		ClientName: "E2E Test Client",
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", suite.server.URL+"/api/v1/admin/chat/messages", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+suite.authToken)

	resp, err := suite.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response handlers.SendMessageResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w, body: %s", err, string(body))
	}

	return &response, nil
}

func (suite *PollingChatE2ETestSuite) getMessages(since string) (*handlers.GetMessagesResponse, error) {
	url := fmt.Sprintf("%s/api/v1/admin/chat/messages?session_id=%s", suite.server.URL, suite.sessionID)
	if since != "" {
		url += "&since=" + since
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+suite.authToken)

	resp, err := suite.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response handlers.GetMessagesResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w, body: %s", err, string(body))
	}

	return &response, nil
}

// E2E Test Cases

func TestPollingChatE2E_BasicMessageFlow(t *testing.T) {
	suite := setupE2ETestSuite(t)
	defer suite.cleanup()

	// Test sending a message
	sendResp, err := suite.sendMessage("Hello, this is a test message")
	require.NoError(t, err)
	assert.True(t, sendResp.Success)
	assert.NotEmpty(t, sendResp.MessageID)

	// Wait a moment for processing
	time.Sleep(100 * time.Millisecond)

	// Test retrieving messages
	getResp, err := suite.getMessages("")
	require.NoError(t, err)
	assert.True(t, getResp.Success)
	assert.GreaterOrEqual(t, len(getResp.Messages), 2) // User message + AI response

	// Verify message content
	userMessage := getResp.Messages[0]
	assert.Equal(t, "user", string(userMessage.Type))
	assert.Equal(t, "Hello, this is a test message", userMessage.Content)
	assert.Equal(t, suite.sessionID, userMessage.SessionID)

	// Verify AI response
	aiMessage := getResp.Messages[1]
	assert.Equal(t, "assistant", string(aiMessage.Type))
	assert.NotEmpty(t, aiMessage.Content)
	assert.Equal(t, suite.sessionID, aiMessage.SessionID)
}

func TestPollingChatE2E_MessagePolling(t *testing.T) {
	suite := setupE2ETestSuite(t)
	defer suite.cleanup()

	// Send initial message
	_, err := suite.sendMessage("Initial message")
	require.NoError(t, err)

	// Get initial messages
	initialResp, err := suite.getMessages("")
	require.NoError(t, err)
	initialCount := len(initialResp.Messages)

	// Send another message
	_, err = suite.sendMessage("Second message")
	require.NoError(t, err)

	// Wait for processing
	time.Sleep(100 * time.Millisecond)

	// Poll for new messages using timestamp
	since := time.Now().Add(-1 * time.Minute).Format(time.RFC3339)
	newResp, err := suite.getMessages(since)
	require.NoError(t, err)

	assert.True(t, newResp.Success)
	assert.Greater(t, len(newResp.Messages), initialCount)
}

func TestPollingChatE2E_ConversationFlow(t *testing.T) {
	suite := setupE2ETestSuite(t)
	defer suite.cleanup()

	// Simulate a conversation
	messages := []string{
		"Hello, I need help with AWS migration",
		"What are the best practices for cloud migration?",
		"How much would it cost to migrate a 100-server environment?",
		"Thank you for the information",
	}

	for i, msg := range messages {
		// Send message
		sendResp, err := suite.sendMessage(msg)
		require.NoError(t, err, "Failed to send message %d", i+1)
		assert.True(t, sendResp.Success)

		// Wait for AI response
		time.Sleep(200 * time.Millisecond)

		// Verify messages were stored
		getResp, err := suite.getMessages("")
		require.NoError(t, err, "Failed to get messages after message %d", i+1)
		assert.True(t, getResp.Success)

		// Should have at least (i+1)*2 messages (user + AI response for each)
		expectedMinMessages := (i + 1) * 2
		assert.GreaterOrEqual(t, len(getResp.Messages), expectedMinMessages,
			"Expected at least %d messages after sending %d user messages", expectedMinMessages, i+1)
	}
}

func TestPollingChatE2E_ConcurrentUsers(t *testing.T) {
	suite := setupE2ETestSuite(t)
	defer suite.cleanup()

	const numUsers = 5
	const messagesPerUser = 3

	var wg sync.WaitGroup
	results := make(chan error, numUsers*messagesPerUser)

	// Simulate multiple concurrent users
	for userID := 0; userID < numUsers; userID++ {
		wg.Add(1)
		go func(uid int) {
			defer wg.Done()

			// Each user has their own session
			userSessionID := fmt.Sprintf("concurrent-session-%d-%d", uid, time.Now().Unix())

			for msgID := 0; msgID < messagesPerUser; msgID++ {
				// Create a temporary suite for this user
				userSuite := &PollingChatE2ETestSuite{
					server:    suite.server,
					client:    suite.client,
					authToken: suite.authToken,
					sessionID: userSessionID,
					logger:    suite.logger,
				}

				message := fmt.Sprintf("User %d message %d", uid, msgID)
				_, err := userSuite.sendMessage(message)
				if err != nil {
					results <- fmt.Errorf("user %d message %d failed: %w", uid, msgID, err)
					return
				}

				// Small delay between messages
				time.Sleep(50 * time.Millisecond)
			}

			results <- nil
		}(userID)
	}

	// Wait for all users to complete
	wg.Wait()
	close(results)

	// Check results
	errorCount := 0
	for err := range results {
		if err != nil {
			t.Errorf("Concurrent user error: %v", err)
			errorCount++
		}
	}

	assert.Equal(t, 0, errorCount, "All concurrent users should succeed")
}

func TestPollingChatE2E_MessageValidation(t *testing.T) {
	suite := setupE2ETestSuite(t)
	defer suite.cleanup()

	testCases := []struct {
		name          string
		content       string
		sessionID     string
		expectSuccess bool
		expectedError string
	}{
		{
			name:          "Valid message",
			content:       "This is a valid message",
			sessionID:     suite.sessionID,
			expectSuccess: true,
		},
		{
			name:          "Empty content",
			content:       "",
			sessionID:     suite.sessionID,
			expectSuccess: false,
			expectedError: "message content cannot be empty",
		},
		{
			name:          "Very long message",
			content:       strings.Repeat("x", 10001),
			sessionID:     suite.sessionID,
			expectSuccess: false,
			expectedError: "message content cannot exceed 10000 characters",
		},
		{
			name:          "Empty session ID",
			content:       "Valid content",
			sessionID:     "",
			expectSuccess: false,
			expectedError: "session ID cannot be empty",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Temporarily change session ID for this test
			originalSessionID := suite.sessionID
			suite.sessionID = tc.sessionID
			defer func() { suite.sessionID = originalSessionID }()

			resp, err := suite.sendMessage(tc.content)
			require.NoError(t, err, "HTTP request should not fail")

			if tc.expectSuccess {
				assert.True(t, resp.Success, "Message should be accepted")
				assert.NotEmpty(t, resp.MessageID, "Should return message ID")
			} else {
				assert.False(t, resp.Success, "Message should be rejected")
				assert.Contains(t, resp.Error, tc.expectedError, "Should contain expected error message")
			}
		})
	}
}

func TestPollingChatE2E_ErrorHandling(t *testing.T) {
	suite := setupE2ETestSuite(t)
	defer suite.cleanup()

	// Test with invalid auth token
	originalToken := suite.authToken
	suite.authToken = "invalid-token"

	resp, err := suite.sendMessage("Test message")
	require.NoError(t, err, "HTTP request should not fail")
	assert.False(t, resp.Success, "Should fail with invalid token")
	assert.Contains(t, resp.Error, "Invalid or expired authentication token")

	// Restore valid token
	suite.authToken = originalToken

	// Test successful message after fixing auth
	resp, err = suite.sendMessage("Test message after auth fix")
	require.NoError(t, err)
	assert.True(t, resp.Success, "Should succeed with valid token")
}

func TestPollingChatE2E_PerformanceUnderLoad(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	suite := setupE2ETestSuite(t)
	defer suite.cleanup()

	const numMessages = 50
	const maxDuration = 10 * time.Second

	start := time.Now()

	// Send multiple messages rapidly
	for i := 0; i < numMessages; i++ {
		message := fmt.Sprintf("Performance test message %d", i)
		resp, err := suite.sendMessage(message)
		require.NoError(t, err, "Message %d should send successfully", i)
		assert.True(t, resp.Success, "Message %d should be accepted", i)

		// Small delay to avoid overwhelming the system
		time.Sleep(10 * time.Millisecond)
	}

	duration := time.Since(start)

	// Verify performance
	assert.Less(t, duration, maxDuration, "Should complete within time limit")

	// Verify all messages were stored
	getResp, err := suite.getMessages("")
	require.NoError(t, err)
	assert.True(t, getResp.Success)

	// Should have at least numMessages * 2 (user + AI response for each)
	expectedMinMessages := numMessages * 2
	assert.GreaterOrEqual(t, len(getResp.Messages), expectedMinMessages,
		"Should have stored all messages and responses")

	suite.logger.Infof("Performance test completed: %d messages in %v", numMessages, duration)
}

func TestPollingChatE2E_MessageOrdering(t *testing.T) {
	suite := setupE2ETestSuite(t)
	defer suite.cleanup()

	// Send messages with specific content to verify ordering
	messages := []string{
		"First message",
		"Second message",
		"Third message",
	}

	for _, msg := range messages {
		resp, err := suite.sendMessage(msg)
		require.NoError(t, err)
		assert.True(t, resp.Success)

		// Small delay between messages
		time.Sleep(100 * time.Millisecond)
	}

	// Retrieve all messages
	getResp, err := suite.getMessages("")
	require.NoError(t, err)
	assert.True(t, getResp.Success)

	// Verify messages are in chronological order
	var userMessages []*domain.ChatMessage
	for _, msg := range getResp.Messages {
		if msg.Type == domain.MessageTypeUser {
			userMessages = append(userMessages, msg)
		}
	}

	require.Len(t, userMessages, 3, "Should have 3 user messages")

	// Verify content order
	assert.Equal(t, "First message", userMessages[0].Content)
	assert.Equal(t, "Second message", userMessages[1].Content)
	assert.Equal(t, "Third message", userMessages[2].Content)

	// Verify timestamp order
	for i := 1; i < len(userMessages); i++ {
		assert.True(t, userMessages[i].CreatedAt.After(userMessages[i-1].CreatedAt) ||
			userMessages[i].CreatedAt.Equal(userMessages[i-1].CreatedAt),
			"Messages should be in chronological order")
	}
}

func TestPollingChatE2E_SessionManagement(t *testing.T) {
	suite := setupE2ETestSuite(t)
	defer suite.cleanup()

	// Send message to create session
	resp, err := suite.sendMessage("Initial message to create session")
	require.NoError(t, err)
	assert.True(t, resp.Success)

	// Verify session was created by retrieving messages
	getResp, err := suite.getMessages("")
	require.NoError(t, err)
	assert.True(t, getResp.Success)
	assert.NotEmpty(t, getResp.Messages)

	// All messages should have the same session ID
	for _, msg := range getResp.Messages {
		assert.Equal(t, suite.sessionID, msg.SessionID, "All messages should belong to the same session")
	}

	// Test with different session ID
	originalSessionID := suite.sessionID
	suite.sessionID = fmt.Sprintf("different-session-%d", time.Now().Unix())

	// Send message to new session
	resp, err = suite.sendMessage("Message in different session")
	require.NoError(t, err)
	assert.True(t, resp.Success)

	// Verify new session has its own messages
	newSessionResp, err := suite.getMessages("")
	require.NoError(t, err)
	assert.True(t, newSessionResp.Success)

	// Should only have messages from the new session
	for _, msg := range newSessionResp.Messages {
		assert.Equal(t, suite.sessionID, msg.SessionID, "Messages should belong to new session")
	}

	// Restore original session and verify it still has its messages
	suite.sessionID = originalSessionID
	originalSessionResp, err := suite.getMessages("")
	require.NoError(t, err)
	assert.True(t, originalSessionResp.Success)
	assert.NotEmpty(t, originalSessionResp.Messages, "Original session should still have its messages")
}

// Run all E2E tests
func TestPollingChatE2E_AllTests(t *testing.T) {
	// This is a meta-test that runs all E2E tests in sequence
	// Useful for CI/CD pipelines

	t.Run("BasicMessageFlow", TestPollingChatE2E_BasicMessageFlow)
	t.Run("MessagePolling", TestPollingChatE2E_MessagePolling)
	t.Run("ConversationFlow", TestPollingChatE2E_ConversationFlow)
	t.Run("MessageValidation", TestPollingChatE2E_MessageValidation)
	t.Run("ErrorHandling", TestPollingChatE2E_ErrorHandling)
	t.Run("MessageOrdering", TestPollingChatE2E_MessageOrdering)
	t.Run("SessionManagement", TestPollingChatE2E_SessionManagement)

	// Run performance and concurrent tests only if not in short mode
	if !testing.Short() {
		t.Run("ConcurrentUsers", TestPollingChatE2E_ConcurrentUsers)
		t.Run("PerformanceUnderLoad", TestPollingChatE2E_PerformanceUnderLoad)
	}
}
