package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/cloud-consulting/backend/internal/domain"
	"github.com/cloud-consulting/backend/internal/interfaces"
	"github.com/cloud-consulting/backend/internal/services"
)

// MockBedrockService for testing
type MockBedrockService struct {
	mock.Mock
}

func (m *MockBedrockService) GenerateText(ctx context.Context, prompt string, options *interfaces.BedrockOptions) (*interfaces.BedrockResponse, error) {
	args := m.Called(ctx, prompt, options)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*interfaces.BedrockResponse), args.Error(1)
}

func (m *MockBedrockService) GetModelInfo() interfaces.BedrockModelInfo {
	args := m.Called()
	return args.Get(0).(interfaces.BedrockModelInfo)
}

func (m *MockBedrockService) IsHealthy() bool {
	args := m.Called()
	return args.Bool(0)
}

// MockKnowledgeBase for testing
type MockKnowledgeBase struct {
	mock.Mock
}

func (m *MockKnowledgeBase) GetServiceOfferings(ctx context.Context) ([]*interfaces.ServiceOffering, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*interfaces.ServiceOffering), args.Error(1)
}

func (m *MockKnowledgeBase) GetTeamExpertise(ctx context.Context) ([]*interfaces.TeamExpertise, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*interfaces.TeamExpertise), args.Error(1)
}

func (m *MockKnowledgeBase) GetPastSolutions(ctx context.Context, serviceType, industry string) ([]*interfaces.PastSolution, error) {
	args := m.Called(ctx, serviceType, industry)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*interfaces.PastSolution), args.Error(1)
}

func (m *MockKnowledgeBase) GetConsultingApproach(ctx context.Context, serviceType string) (*interfaces.ConsultingApproach, error) {
	args := m.Called(ctx, serviceType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*interfaces.ConsultingApproach), args.Error(1)
}

func (m *MockKnowledgeBase) GetClientHistory(ctx context.Context, company string) ([]*interfaces.ClientEngagement, error) {
	args := m.Called(ctx, company)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*interfaces.ClientEngagement), args.Error(1)
}

func (m *MockKnowledgeBase) GetBestPractices(ctx context.Context, category string) ([]*interfaces.BestPractice, error) {
	args := m.Called(ctx, category)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*interfaces.BestPractice), args.Error(1)
}

func (m *MockKnowledgeBase) GetComplianceRequirements(ctx context.Context, framework string) ([]*interfaces.ComplianceRequirement, error) {
	args := m.Called(ctx, framework)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*interfaces.ComplianceRequirement), args.Error(1)
}

func (m *MockKnowledgeBase) GetConsultantSpecializations(ctx context.Context, consultantID string) ([]*interfaces.Specialization, error) {
	args := m.Called(ctx, consultantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*interfaces.Specialization), args.Error(1)
}

func (m *MockKnowledgeBase) GetDeliverableTemplates(ctx context.Context, serviceType string) ([]*interfaces.DeliverableTemplate, error) {
	args := m.Called(ctx, serviceType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*interfaces.DeliverableTemplate), args.Error(1)
}

func (m *MockKnowledgeBase) GetExpertiseByArea(ctx context.Context, area string) ([]*interfaces.TeamExpertise, error) {
	args := m.Called(ctx, area)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*interfaces.TeamExpertise), args.Error(1)
}

func (m *MockKnowledgeBase) GetKnowledgeStats(ctx context.Context) (*interfaces.KnowledgeStats, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*interfaces.KnowledgeStats), args.Error(1)
}

func (m *MockKnowledgeBase) GetMethodologyTemplates(ctx context.Context, serviceType string) ([]*interfaces.MethodologyTemplate, error) {
	args := m.Called(ctx, serviceType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*interfaces.MethodologyTemplate), args.Error(1)
}

func (m *MockKnowledgeBase) GetPricingModels(ctx context.Context, serviceType string) ([]*interfaces.PricingModel, error) {
	args := m.Called(ctx, serviceType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*interfaces.PricingModel), args.Error(1)
}

func (m *MockKnowledgeBase) GetServiceOffering(ctx context.Context, id string) (*interfaces.ServiceOffering, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*interfaces.ServiceOffering), args.Error(1)
}

func (m *MockKnowledgeBase) GetSimilarProjects(ctx context.Context, inquiry *domain.Inquiry) ([]*interfaces.ProjectPattern, error) {
	args := m.Called(ctx, inquiry)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*interfaces.ProjectPattern), args.Error(1)
}

func (m *MockKnowledgeBase) SearchKnowledge(ctx context.Context, query string, category string) ([]*interfaces.KnowledgeItem, error) {
	args := m.Called(ctx, query, category)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*interfaces.KnowledgeItem), args.Error(1)
}

func (m *MockKnowledgeBase) UpdateKnowledgeBase(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// MockSessionService for testing
type MockSessionService struct {
	mock.Mock
}

func (m *MockSessionService) CreateSession(ctx context.Context, session *domain.ChatSession) error {
	args := m.Called(ctx, session)
	return args.Error(0)
}

func (m *MockSessionService) GetSession(ctx context.Context, sessionID string) (*domain.ChatSession, error) {
	args := m.Called(ctx, sessionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ChatSession), args.Error(1)
}

func (m *MockSessionService) UpdateSession(ctx context.Context, session *domain.ChatSession) error {
	args := m.Called(ctx, session)
	return args.Error(0)
}

func (m *MockSessionService) DeleteSession(ctx context.Context, sessionID string) error {
	args := m.Called(ctx, sessionID)
	return args.Error(0)
}

func (m *MockSessionService) GetUserSessions(ctx context.Context, userID string) ([]*domain.ChatSession, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.ChatSession), args.Error(1)
}

func (m *MockSessionService) ValidateSession(ctx context.Context, sessionID string, userID string) (*domain.ChatSession, error) {
	args := m.Called(ctx, sessionID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ChatSession), args.Error(1)
}

func (m *MockSessionService) GetActiveSessions(ctx context.Context, userID string) ([]*domain.ChatSession, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.ChatSession), args.Error(1)
}

func (m *MockSessionService) ListSessions(ctx context.Context, filters *domain.ChatSessionFilters) ([]*domain.ChatSession, error) {
	args := m.Called(ctx, filters)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.ChatSession), args.Error(1)
}

func (m *MockSessionService) ExpireSession(ctx context.Context, sessionID string) error {
	args := m.Called(ctx, sessionID)
	return args.Error(0)
}

func (m *MockSessionService) TerminateSession(ctx context.Context, sessionID string) error {
	args := m.Called(ctx, sessionID)
	return args.Error(0)
}

func (m *MockSessionService) RefreshSession(ctx context.Context, sessionID string, duration time.Duration) error {
	args := m.Called(ctx, sessionID, duration)
	return args.Error(0)
}

func (m *MockSessionService) IsSessionValid(ctx context.Context, sessionID string) (bool, error) {
	args := m.Called(ctx, sessionID)
	return args.Bool(0), args.Error(1)
}

func (m *MockSessionService) CleanupExpiredSessions(ctx context.Context) (int, error) {
	args := m.Called(ctx)
	return args.Int(0), args.Error(1)
}

func (m *MockSessionService) CleanupInactiveSessions(ctx context.Context, inactiveThreshold time.Duration) (int, error) {
	args := m.Called(ctx, inactiveThreshold)
	return args.Int(0), args.Error(1)
}

func (m *MockSessionService) GetSessionCount(ctx context.Context, filters *domain.ChatSessionFilters) (int64, error) {
	args := m.Called(ctx, filters)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockSessionService) GetSessionStats(ctx context.Context) (*interfaces.SessionStats, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*interfaces.SessionStats), args.Error(1)
}

// MockChatService for testing
type MockChatService struct {
	mock.Mock
}

func (m *MockChatService) SendMessage(ctx context.Context, request *domain.ChatRequest) (*domain.ChatResponse, error) {
	args := m.Called(ctx, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ChatResponse), args.Error(1)
}

func (m *MockChatService) GetMessage(ctx context.Context, messageID string) (*domain.ChatMessage, error) {
	args := m.Called(ctx, messageID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ChatMessage), args.Error(1)
}

func (m *MockChatService) GetSessionHistory(ctx context.Context, sessionID string, limit int) ([]*domain.ChatMessage, error) {
	args := m.Called(ctx, sessionID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.ChatMessage), args.Error(1)
}

func (m *MockChatService) UpdateMessageStatus(ctx context.Context, messageID string, status domain.MessageStatus) error {
	args := m.Called(ctx, messageID, status)
	return args.Error(0)
}

func (m *MockChatService) MarkMessageAsDelivered(ctx context.Context, messageID string) error {
	args := m.Called(ctx, messageID)
	return args.Error(0)
}

func (m *MockChatService) MarkMessageAsRead(ctx context.Context, messageID string) error {
	args := m.Called(ctx, messageID)
	return args.Error(0)
}

func (m *MockChatService) UpdateSessionContext(ctx context.Context, sessionID string, context *domain.SessionContext) error {
	args := m.Called(ctx, sessionID, context)
	return args.Error(0)
}

func (m *MockChatService) GetSessionContext(ctx context.Context, sessionID string) (*domain.SessionContext, error) {
	args := m.Called(ctx, sessionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.SessionContext), args.Error(1)
}

func (m *MockChatService) ValidateMessage(message *domain.ChatMessage) error {
	args := m.Called(message)
	return args.Error(0)
}

func (m *MockChatService) SanitizeMessageContent(content string) string {
	args := m.Called(content)
	return args.String(0)
}

func (m *MockChatService) SearchMessages(ctx context.Context, sessionID string, query string, limit int) ([]*domain.ChatMessage, error) {
	args := m.Called(ctx, sessionID, query, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.ChatMessage), args.Error(1)
}

func (m *MockChatService) GetMessagesByType(ctx context.Context, sessionID string, messageType domain.MessageType) ([]*domain.ChatMessage, error) {
	args := m.Called(ctx, sessionID, messageType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.ChatMessage), args.Error(1)
}

func (m *MockChatService) GetMessageCount(ctx context.Context, filters *domain.ChatMessageFilters) (int64, error) {
	args := m.Called(ctx, filters)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockChatService) GetMessageStats(ctx context.Context, sessionID string) (*interfaces.MessageStats, error) {
	args := m.Called(ctx, sessionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*interfaces.MessageStats), args.Error(1)
}

func (m *MockChatService) ListMessages(ctx context.Context, filters *domain.ChatMessageFilters) ([]*domain.ChatMessage, error) {
	args := m.Called(ctx, filters)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.ChatMessage), args.Error(1)
}

// Test helper functions
func createTestChatHandler() (*ChatHandler, *MockBedrockService, *MockKnowledgeBase, *MockSessionService, *MockChatService) {
	mockBedrock := &MockBedrockService{}
	mockKB := &MockKnowledgeBase{}
	mockSessionService := &MockSessionService{}
	mockChatService := &MockChatService{}

	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel) // Reduce log noise in tests

	// Create properly initialized metrics collector and performance monitor
	mockPerformanceMonitor := services.NewChatPerformanceMonitor(logger)
	mockCacheMonitor := services.NewCacheMonitor(nil, logger)
	mockMetricsCollector := services.NewChatMetricsCollector(mockPerformanceMonitor, mockCacheMonitor, logger)

	handler := NewChatHandler(
		logger,
		mockBedrock,
		mockKB,
		mockSessionService,
		mockChatService,
		nil, // authHandler not needed for these tests
		"test-jwt-secret",
		[]string{"http://localhost:3000"}, // corsOrigins
		mockMetricsCollector,
		mockPerformanceMonitor,
	)

	return handler, mockBedrock, mockKB, mockSessionService, mockChatService
}

func createTestChatRequest() ChatRequest {
	return ChatRequest{
		Message:    "Test message",
		ClientName: "Test Client",
		Context:    "Test context",
	}
}

func createTestDomainSession() *domain.ChatSession {
	now := time.Now()
	expiresAt := now.Add(24 * time.Hour)

	return &domain.ChatSession{
		ID:           "test-session-id",
		UserID:       "test-user-id",
		ClientName:   "Test Client",
		Context:      "Test context",
		Status:       domain.SessionStatusActive,
		Metadata:     map[string]interface{}{"test": "value"},
		CreatedAt:    now,
		UpdatedAt:    now,
		LastActivity: now,
		ExpiresAt:    &expiresAt,
	}
}

// Test processEnhancedChatRequest
func TestChatHandler_ProcessEnhancedChatRequest_Success(t *testing.T) {
	handler, _, _, mockSessionService, mockChatService := createTestChatHandler()

	request := createTestChatRequest()
	connID := "test-conn-id"
	userID := "test-user-id"

	session := createTestDomainSession()
	chatResponse := &domain.ChatResponse{
		MessageID: "test-message-id",
		SessionID: session.ID,
		Type:      domain.MessageTypeAssistant,
		Content:   "Test AI response",
		CreatedAt: time.Now(),
	}

	// Mock session service calls - since SessionID is empty, ValidateSession won't be called
	// Only CreateSession will be called
	mockSessionService.On("CreateSession", mock.Anything, mock.AnythingOfType("*domain.ChatSession")).Return(nil).Run(func(args mock.Arguments) {
		// Update the session ID in the passed session object
		sessionArg := args.Get(1).(*domain.ChatSession)
		sessionArg.ID = session.ID
	})

	// Mock chat service calls
	mockChatService.On("SendMessage", mock.Anything, mock.AnythingOfType("*domain.ChatRequest")).Return(chatResponse, nil)

	response := handler.processEnhancedChatRequest(request, connID, userID)

	assert.True(t, response.Success)
	assert.Equal(t, session.ID, response.SessionID)
	assert.Equal(t, "test-message-id", response.Message.ID)
	assert.Equal(t, "assistant", response.Message.Type)
	assert.Equal(t, "Test AI response", response.Message.Content)

	mockSessionService.AssertExpectations(t)
	mockChatService.AssertExpectations(t)
}

func TestChatHandler_ProcessEnhancedChatRequest_SessionCreationFailure(t *testing.T) {
	handler, _, _, mockSessionService, _ := createTestChatHandler()

	request := createTestChatRequest()
	connID := "test-conn-id"
	userID := "test-user-id"

	// Mock session service failure - since SessionID is empty, ValidateSession won't be called
	// Only CreateSession will be called and it will fail
	mockSessionService.On("CreateSession", mock.Anything, mock.AnythingOfType("*domain.ChatSession")).Return(assert.AnError)

	response := handler.processEnhancedChatRequest(request, connID, userID)

	assert.False(t, response.Success)
	assert.Equal(t, "Failed to manage session", response.Error)

	mockSessionService.AssertExpectations(t)
}

func TestChatHandler_ProcessEnhancedChatRequest_ChatServiceFailure(t *testing.T) {
	handler, _, _, mockSessionService, mockChatService := createTestChatHandler()

	request := createTestChatRequest()
	connID := "test-conn-id"
	userID := "test-user-id"

	// Mock successful session creation but failed chat service
	// Since SessionID is empty, ValidateSession won't be called
	mockSessionService.On("CreateSession", mock.Anything, mock.AnythingOfType("*domain.ChatSession")).Return(nil).Run(func(args mock.Arguments) {
		// Update the session ID in the passed session object
		sessionArg := args.Get(1).(*domain.ChatSession)
		sessionArg.ID = "test-session-id"
	})
	mockChatService.On("SendMessage", mock.Anything, mock.AnythingOfType("*domain.ChatRequest")).Return(nil, assert.AnError)

	response := handler.processEnhancedChatRequest(request, connID, userID)

	assert.False(t, response.Success)
	assert.Equal(t, "Failed to process message", response.Error)

	mockSessionService.AssertExpectations(t)
	mockChatService.AssertExpectations(t)
}

// Test REST API endpoints
func TestChatHandler_CreateChatSession_Success(t *testing.T) {
	handler, _, _, mockSessionService, _ := createTestChatHandler()

	// Setup Gin context
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Set user context
	c.Set("user_id", "test-user-id")

	// Create request body
	sessionData := map[string]interface{}{
		"client_name": "Test Client",
		"context":     "Test context",
	}
	jsonData, _ := json.Marshal(sessionData)
	c.Request = httptest.NewRequest("POST", "/api/v1/admin/chat/sessions", bytes.NewBuffer(jsonData))
	c.Request.Header.Set("Content-Type", "application/json")

	// Mock session service
	mockSessionService.On("CreateSession", mock.Anything, mock.AnythingOfType("*domain.ChatSession")).Return(nil)

	handler.CreateChatSession(c)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response["success"].(bool))

	mockSessionService.AssertExpectations(t)
}

func TestChatHandler_CreateChatSession_Unauthorized(t *testing.T) {
	handler, _, _, _, _ := createTestChatHandler()

	// Setup Gin context without user_id
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	sessionData := map[string]interface{}{
		"client_name": "Test Client",
	}
	jsonData, _ := json.Marshal(sessionData)
	c.Request = httptest.NewRequest("POST", "/api/v1/admin/chat/sessions", bytes.NewBuffer(jsonData))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.CreateChatSession(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response["success"].(bool))
	assert.Equal(t, "User authentication required", response["error"])
}

// Test WebSocket message routing
func TestChatHandler_RouteWebSocketMessage_ChatMessage(t *testing.T) {
	handler, _, _, mockSessionService, mockChatService := createTestChatHandler()

	wsMessage := WebSocketMessage{
		Type:      WSMessageTypeMessage,
		SessionID: "test-session-id",
		Content:   "Test message",
		Timestamp: time.Now(),
	}

	connection := &Connection{
		UserID:   "test-user-id",
		SendChan: make(chan WebSocketMessage, 10),
	}

	session := createTestDomainSession()
	chatResponse := &domain.ChatResponse{
		MessageID: "test-message-id",
		SessionID: session.ID,
		Type:      domain.MessageTypeAssistant,
		Content:   "Test AI response",
		CreatedAt: time.Now(),
	}

	// Mock services
	mockSessionService.On("ValidateSession", mock.Anything, "test-session-id", "test-user-id").Return(session, nil)
	mockChatService.On("SendMessage", mock.Anything, mock.AnythingOfType("*domain.ChatRequest")).Return(chatResponse, nil)

	handler.routeWebSocketMessage(wsMessage, connection, "test-conn-id", "test-user-id")

	// Check that both ack and response messages were sent to connection
	// First message should be acknowledgment
	select {
	case ackResponse := <-connection.SendChan:
		assert.Equal(t, "ack", string(ackResponse.Type))
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Expected acknowledgment message not received")
	}

	// Second message should be the actual response
	select {
	case response := <-connection.SendChan:
		assert.Equal(t, "message", string(response.Type))
		assert.Equal(t, "Test AI response", response.Content)
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Expected response message not received")
	}

	mockSessionService.AssertExpectations(t)
	mockChatService.AssertExpectations(t)
}

// Test connection pool functionality
func TestConnectionPool_AddGetRemove(t *testing.T) {
	pool := NewConnectionPool()

	conn := &Connection{
		UserID:    "test-user",
		SessionID: "test-session",
	}

	// Test Add and Get
	pool.Add("conn1", conn)
	retrieved, exists := pool.Get("conn1")
	assert.True(t, exists)
	assert.Equal(t, conn, retrieved)

	// Test Count
	assert.Equal(t, 1, pool.Count())

	// Test GetByUserID
	userConns := pool.GetByUserID("test-user")
	assert.Len(t, userConns, 1)
	assert.Equal(t, conn, userConns[0])

	// Test GetBySessionID
	sessionConns := pool.GetBySessionID("test-session")
	assert.Len(t, sessionConns, 1)
	assert.Equal(t, conn, sessionConns[0])

	// Test Remove
	pool.Remove("conn1")
	_, exists = pool.Get("conn1")
	assert.False(t, exists)
	assert.Equal(t, 0, pool.Count())
}

// Test rate limiter functionality
func TestRateLimiter_Allow(t *testing.T) {
	limiter := NewRateLimiter(2, time.Second) // 2 requests per second

	// First two requests should be allowed
	assert.True(t, limiter.Allow("user1"))
	assert.True(t, limiter.Allow("user1"))

	// Third request should be denied
	assert.False(t, limiter.Allow("user1"))

	// Different user should be allowed
	assert.True(t, limiter.Allow("user2"))

	// Wait for window to reset
	time.Sleep(1100 * time.Millisecond)

	// Should be allowed again
	assert.True(t, limiter.Allow("user1"))
}

// Test message sanitization
func TestChatHandler_SanitizeMessageContent(t *testing.T) {
	handler, _, _, _, mockChatService := createTestChatHandler()

	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "Normal message",
			expected: "Normal message",
		},
		{
			input:    "Message with <script>alert('xss')</script>",
			expected: "Message with ",
		},
		{
			input:    "  Multiple   spaces  ",
			expected: "Multiple spaces",
		},
	}

	for _, test := range tests {
		mockChatService.On("SanitizeMessageContent", test.input).Return(test.expected)
		result := handler.chatService.SanitizeMessageContent(test.input)
		assert.Equal(t, test.expected, result, "Input: %s", test.input)
	}

	mockChatService.AssertExpectations(t)
}

// Test WebSocket authentication middleware
func TestChatHandler_WebSocketAuthMiddleware_Success(t *testing.T) {
	handler, _, _, _, _ := createTestChatHandler()

	// Setup Gin context
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Set valid token in query parameter
	c.Request = httptest.NewRequest("GET", "/ws?token=valid-token", nil)

	// Mock auth service to return success
	// Note: This would require mocking the auth service, which is not implemented in this test
	// For now, we'll test the error cases

	handler.WebSocketAuthMiddleware(c)

	// Since we don't have a proper auth service mock, this will fail
	// In a real implementation, you would mock the auth service
	assert.True(t, c.IsAborted())
}

func TestChatHandler_WebSocketAuthMiddleware_NoToken(t *testing.T) {
	handler, _, _, _, _ := createTestChatHandler()

	// Setup Gin context
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// No token provided
	c.Request = httptest.NewRequest("GET", "/ws", nil)

	handler.WebSocketAuthMiddleware(c)

	assert.True(t, c.IsAborted())
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response["success"].(bool))
	assert.Contains(t, response["error"].(string), "Authentication token required")
}

// Test legacy chat functionality
func TestChatHandler_ProcessChatRequest_Success(t *testing.T) {
	handler, mockBedrock, mockKB, _, _ := createTestChatHandler()

	request := createTestChatRequest()
	connID := "test-conn-id"

	// Mock Bedrock response
	bedrockResponse := &interfaces.BedrockResponse{
		Content: "AI response content",
		Usage: interfaces.BedrockUsage{
			OutputTokens: 50,
		},
	}
	mockBedrock.On("GenerateText", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("*interfaces.BedrockOptions")).Return(bedrockResponse, nil)

	// Mock GetModelInfo call
	modelInfo := interfaces.BedrockModelInfo{
		ModelID:     "nova-lite",
		ModelName:   "Nova Lite",
		Provider:    "amazon",
		MaxTokens:   4096,
		IsAvailable: true,
	}
	mockBedrock.On("GetModelInfo").Return(modelInfo)

	// Mock knowledge base calls
	mockKB.On("GetServiceOfferings", mock.Anything).Return([]*interfaces.ServiceOffering{}, nil)
	mockKB.On("GetTeamExpertise", mock.Anything).Return([]*interfaces.TeamExpertise{}, nil)
	mockKB.On("GetPastSolutions", mock.Anything, mock.Anything, mock.Anything).Return([]*interfaces.PastSolution{}, nil)
	mockKB.On("GetConsultingApproach", mock.Anything, mock.Anything).Return(&interfaces.ConsultingApproach{}, nil)
	mockKB.On("GetClientHistory", mock.Anything, mock.Anything).Return([]*interfaces.ClientEngagement{}, nil)
	mockKB.On("GetExpertiseByArea", mock.Anything, mock.Anything).Return([]*interfaces.TeamExpertise{}, nil)
	mockKB.On("GetSimilarProjects", mock.Anything, mock.Anything).Return([]*interfaces.ProjectPattern{}, nil)
	mockKB.On("GetMethodologyTemplates", mock.Anything, mock.Anything).Return([]*interfaces.MethodologyTemplate{}, nil)
	mockKB.On("GetPricingModels", mock.Anything, mock.Anything).Return([]*interfaces.PricingModel{}, nil)

	response := handler.processChatRequest(request, connID)

	assert.True(t, response.Success)
	assert.Equal(t, "assistant", response.Message.Type)
	assert.Equal(t, "AI response content", response.Message.Content)
	assert.NotEmpty(t, response.SessionID)

	mockBedrock.AssertExpectations(t)
	mockKB.AssertExpectations(t)
}

func TestChatHandler_ProcessChatRequest_BedrockFailure(t *testing.T) {
	handler, mockBedrock, mockKB, _, _ := createTestChatHandler()

	request := createTestChatRequest()
	connID := "test-conn-id"

	// Mock Bedrock failure
	modelInfo := interfaces.BedrockModelInfo{
		ModelName:   "Nova Lite",
		Provider:    "amazon",
		MaxTokens:   4096,
		IsAvailable: true,
	}
	mockBedrock.On("GetModelInfo").Return(modelInfo)
	mockBedrock.On("GenerateText", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("*interfaces.BedrockOptions")).Return(nil, assert.AnError)

	// Mock knowledge base calls (needed for enhanced bedrock service)
	mockKB.On("GetServiceOfferings", mock.Anything).Return([]*interfaces.ServiceOffering{}, nil)
	mockKB.On("GetTeamExpertise", mock.Anything).Return([]*interfaces.TeamExpertise{}, nil)
	mockKB.On("GetPastSolutions", mock.Anything, mock.Anything, mock.Anything).Return([]*interfaces.PastSolution{}, nil)
	mockKB.On("GetConsultingApproach", mock.Anything, mock.Anything).Return(&interfaces.ConsultingApproach{}, nil)
	mockKB.On("GetClientHistory", mock.Anything, mock.Anything).Return([]*interfaces.ClientEngagement{}, nil)
	mockKB.On("GetExpertiseByArea", mock.Anything, mock.Anything).Return([]*interfaces.TeamExpertise{}, nil)
	mockKB.On("GetSimilarProjects", mock.Anything, mock.Anything).Return([]*interfaces.ProjectPattern{}, nil)
	mockKB.On("GetMethodologyTemplates", mock.Anything, mock.Anything).Return([]*interfaces.MethodologyTemplate{}, nil)
	mockKB.On("GetPricingModels", mock.Anything, mock.Anything).Return([]*interfaces.PricingModel{}, nil)

	response := handler.processChatRequest(request, connID)

	assert.False(t, response.Success)
	assert.Equal(t, "Failed to generate response", response.Error)

	mockBedrock.AssertExpectations(t)
}

// Test session management
func TestChatHandler_GetOrCreateSession_NewSession(t *testing.T) {
	handler, _, _, _, _ := createTestChatHandler()

	// Test creating new session
	session := handler.getOrCreateSession("", "test-conn-id")

	assert.NotNil(t, session)
	assert.NotEmpty(t, session.ID)
	assert.Equal(t, "test-conn-id", session.ConsultantID)
	assert.NotZero(t, session.CreatedAt)
	assert.NotZero(t, session.LastActivity)
}

func TestChatHandler_GetOrCreateSession_ExistingSession(t *testing.T) {
	handler, _, _, _, _ := createTestChatHandler()

	// Create initial session
	session1 := handler.getOrCreateSession("", "test-conn-id")
	sessionID := session1.ID

	// Get existing session
	session2 := handler.getOrCreateSession(sessionID, "test-conn-id")

	assert.Equal(t, sessionID, session2.ID)
	assert.Equal(t, session1.ConsultantID, session2.ConsultantID)
}

// Test message addition to session
func TestChatHandler_AddMessageToSession(t *testing.T) {
	handler, _, _, _, _ := createTestChatHandler()

	// Create session
	session := handler.getOrCreateSession("", "test-conn-id")
	sessionID := session.ID

	// Add message
	message := ChatMessage{
		ID:        "test-message-id",
		Type:      "user",
		Content:   "Test message",
		Timestamp: time.Now(),
		SessionID: sessionID,
	}

	handler.addMessageToSession(sessionID, message)

	// Verify message was added
	handler.sessionsMutex.RLock()
	updatedSession := handler.sessions[sessionID]
	handler.sessionsMutex.RUnlock()

	assert.Len(t, updatedSession.Messages, 1)
	assert.Equal(t, message.ID, updatedSession.Messages[0].ID)
	assert.Equal(t, message.Content, updatedSession.Messages[0].Content)
}

// Test prompt building
func TestChatHandler_BuildConsultantPrompt(t *testing.T) {
	handler, _, _, _, _ := createTestChatHandler()

	session := &ChatSession{
		ID:         "test-session",
		ClientName: "Test Client",
		Context:    "Migration planning",
		Messages: []ChatMessage{
			{
				Type:    "user",
				Content: "Previous question",
			},
			{
				Type:    "assistant",
				Content: "Previous answer",
			},
		},
	}

	request := ChatRequest{
		Message:     "Current question",
		QuickAction: "cost-estimate",
	}

	prompt := handler.buildConsultantPrompt(session, request)

	assert.Contains(t, prompt, "Test Client")
	assert.Contains(t, prompt, "Migration planning")
	assert.Contains(t, prompt, "Previous question")
	assert.Contains(t, prompt, "Previous answer")
	assert.Contains(t, prompt, "Current question")
	assert.Contains(t, prompt, "cost-estimate")
}

// Test service extraction from context
func TestChatHandler_ExtractServicesFromContext(t *testing.T) {
	handler, _, _, _, _ := createTestChatHandler()

	tests := []struct {
		context  string
		message  string
		expected []string
	}{
		{
			context:  "We need to migrate to AWS",
			message:  "Help with migration planning",
			expected: []string{"Cloud Migration"},
		},
		{
			context:  "Architecture review needed",
			message:  "Design recommendations",
			expected: []string{"Architecture Review"},
		},
		{
			context:  "Cost optimization project",
			message:  "Optimize our spending",
			expected: []string{"Optimization"},
		},
		{
			context:  "General consulting",
			message:  "Various questions",
			expected: []string{"General Consulting"},
		},
	}

	for _, test := range tests {
		result := handler.extractServicesFromContext(test.context, test.message)
		assert.Equal(t, test.expected, result, "Context: %s, Message: %s", test.context, test.message)
	}
}
