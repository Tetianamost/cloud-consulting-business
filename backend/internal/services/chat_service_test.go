package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/cloud-consulting/backend/internal/domain"
	"github.com/cloud-consulting/backend/internal/interfaces"
)

// MockChatMessageRepository is a mock implementation of ChatMessageRepository
type MockChatMessageRepository struct {
	mock.Mock
}

func (m *MockChatMessageRepository) Create(ctx context.Context, message *domain.ChatMessage) error {
	args := m.Called(ctx, message)
	return args.Error(0)
}

func (m *MockChatMessageRepository) GetByID(ctx context.Context, id string) (*domain.ChatMessage, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ChatMessage), args.Error(1)
}

func (m *MockChatMessageRepository) Update(ctx context.Context, message *domain.ChatMessage) error {
	args := m.Called(ctx, message)
	return args.Error(0)
}

func (m *MockChatMessageRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockChatMessageRepository) GetBySessionID(ctx context.Context, sessionID string, limit int, offset int) ([]*domain.ChatMessage, error) {
	args := m.Called(ctx, sessionID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.ChatMessage), args.Error(1)
}

func (m *MockChatMessageRepository) List(ctx context.Context, filters *domain.ChatMessageFilters) ([]*domain.ChatMessage, error) {
	args := m.Called(ctx, filters)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.ChatMessage), args.Error(1)
}

func (m *MockChatMessageRepository) Count(ctx context.Context, filters *domain.ChatMessageFilters) (int64, error) {
	args := m.Called(ctx, filters)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockChatMessageRepository) GetByType(ctx context.Context, sessionID string, messageType domain.MessageType) ([]*domain.ChatMessage, error) {
	args := m.Called(ctx, sessionID, messageType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.ChatMessage), args.Error(1)
}

func (m *MockChatMessageRepository) GetByStatus(ctx context.Context, sessionID string, status domain.MessageStatus) ([]*domain.ChatMessage, error) {
	args := m.Called(ctx, sessionID, status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.ChatMessage), args.Error(1)
}

func (m *MockChatMessageRepository) UpdateStatus(ctx context.Context, messageID string, status domain.MessageStatus) error {
	args := m.Called(ctx, messageID, status)
	return args.Error(0)
}

func (m *MockChatMessageRepository) Search(ctx context.Context, sessionID string, query string, limit int) ([]*domain.ChatMessage, error) {
	args := m.Called(ctx, sessionID, query, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.ChatMessage), args.Error(1)
}

func (m *MockChatMessageRepository) DeleteBySessionID(ctx context.Context, sessionID string) error {
	args := m.Called(ctx, sessionID)
	return args.Error(0)
}

func (m *MockChatMessageRepository) GetLatestBySessionID(ctx context.Context, sessionID string, limit int) ([]*domain.ChatMessage, error) {
	args := m.Called(ctx, sessionID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.ChatMessage), args.Error(1)
}

// MockSessionService is a mock implementation of SessionService
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

func (m *MockSessionService) ValidateSession(ctx context.Context, sessionID string, userID string) (*domain.ChatSession, error) {
	args := m.Called(ctx, sessionID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ChatSession), args.Error(1)
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

// MockBedrockServiceForChat is a mock implementation of BedrockService for chat tests
type MockBedrockServiceForChat struct {
	mock.Mock
}

func (m *MockBedrockServiceForChat) GenerateText(ctx context.Context, prompt string, options *interfaces.BedrockOptions) (*interfaces.BedrockResponse, error) {
	args := m.Called(ctx, prompt, options)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*interfaces.BedrockResponse), args.Error(1)
}

func (m *MockBedrockServiceForChat) GetModelInfo() interfaces.BedrockModelInfo {
	args := m.Called()
	return args.Get(0).(interfaces.BedrockModelInfo)
}

func (m *MockBedrockServiceForChat) IsHealthy() bool {
	args := m.Called()
	return args.Bool(0)
}

// Test helper functions
func createTestMessage() *domain.ChatMessage {
	return &domain.ChatMessage{
		ID:        "test-message-id",
		SessionID: "test-session-id",
		Type:      domain.MessageTypeUser,
		Content:   "Test message content",
		Metadata:  map[string]interface{}{"test": "value"},
		Status:    domain.MessageStatusSent,
		CreatedAt: time.Now(),
	}
}

func createTestChatRequest() *domain.ChatRequest {
	return &domain.ChatRequest{
		SessionID: "test-session-id",
		Content:   "Test message content",
		Type:      domain.MessageTypeUser,
		Metadata:  map[string]interface{}{"test": "value"},
	}
}

func createTestChatService() (interfaces.ChatService, *MockChatMessageRepository, *MockSessionService, *MockBedrockServiceForChat) {
	mockMessageRepo := &MockChatMessageRepository{}
	mockSessionService := &MockSessionService{}
	mockBedrockService := &MockBedrockServiceForChat{}
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel) // Reduce log noise in tests

	service := NewChatService(mockMessageRepo, mockSessionService, mockBedrockService, logger)
	return service, mockMessageRepo, mockSessionService, mockBedrockService
}

// Test SendMessage
func TestChatService_SendMessage_Success(t *testing.T) {
	service, mockMessageRepo, mockSessionService, mockBedrockService := createTestChatService()
	ctx := context.Background()

	request := createTestChatRequest()
	session := createTestSession()

	// Mock session retrieval
	mockSessionService.On("GetSession", ctx, "test-session-id").Return(session, nil)

	// Mock message creation (user message)
	mockMessageRepo.On("Create", ctx, mock.AnythingOfType("*domain.ChatMessage")).Return(nil).Once()

	// Mock session update
	mockSessionService.On("UpdateSession", ctx, mock.AnythingOfType("*domain.ChatSession")).Return(nil)

	// Mock AI response generation
	aiResponse := &interfaces.BedrockResponse{
		Content: "AI response content",
		Usage: interfaces.BedrockUsage{
			OutputTokens: 50,
		},
	}
	mockBedrockService.On("GenerateText", ctx, mock.AnythingOfType("string"), mock.AnythingOfType("*interfaces.BedrockOptions")).Return(aiResponse, nil)

	// Mock getting conversation history for prompt building
	mockMessageRepo.On("GetLatestBySessionID", ctx, "test-session-id", 5).Return([]*domain.ChatMessage{}, nil)

	// Mock AI message creation
	mockMessageRepo.On("Create", ctx, mock.AnythingOfType("*domain.ChatMessage")).Return(nil).Once()

	// Mock user message status update
	mockMessageRepo.On("UpdateStatus", ctx, mock.AnythingOfType("string"), domain.MessageStatusDelivered).Return(nil)

	response, err := service.SendMessage(ctx, request)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "test-session-id", response.SessionID)
	assert.Equal(t, domain.MessageTypeAssistant, response.Type)
	assert.Equal(t, "AI response content", response.Content)

	mockMessageRepo.AssertExpectations(t)
	mockSessionService.AssertExpectations(t)
	mockBedrockService.AssertExpectations(t)
}

func TestChatService_SendMessage_ValidationError(t *testing.T) {
	service, _, _, _ := createTestChatService()
	ctx := context.Background()

	// Test with invalid request (empty session ID)
	request := &domain.ChatRequest{
		Content: "Test message",
		Type:    domain.MessageTypeUser,
	}

	response, err := service.SendMessage(ctx, request)

	assert.Error(t, err)
	assert.Nil(t, response)
	var chatErr *interfaces.ChatError
	assert.ErrorAs(t, err, &chatErr)
	assert.Equal(t, interfaces.ErrCodeValidationFailed, chatErr.Code)
}

func TestChatService_SendMessage_SessionNotFound(t *testing.T) {
	service, _, mockSessionService, _ := createTestChatService()
	ctx := context.Background()

	request := createTestChatRequest()

	mockSessionService.On("GetSession", ctx, "test-session-id").Return(nil, errors.New("session not found"))

	response, err := service.SendMessage(ctx, request)

	assert.Error(t, err)
	assert.Nil(t, response)
	var chatErr *interfaces.ChatError
	assert.ErrorAs(t, err, &chatErr)
	assert.Equal(t, interfaces.ErrCodeSessionNotFound, chatErr.Code)

	mockSessionService.AssertExpectations(t)
}

func TestChatService_SendMessage_InactiveSession(t *testing.T) {
	service, _, mockSessionService, _ := createTestChatService()
	ctx := context.Background()

	request := createTestChatRequest()
	session := createTestSession()
	session.Status = domain.SessionStatusInactive

	mockSessionService.On("GetSession", ctx, "test-session-id").Return(session, nil)

	response, err := service.SendMessage(ctx, request)

	assert.Error(t, err)
	assert.Nil(t, response)
	var validationErr *interfaces.SessionValidationError
	assert.ErrorAs(t, err, &validationErr)
	assert.Equal(t, interfaces.ErrCodeSessionInvalid, validationErr.Code)

	mockSessionService.AssertExpectations(t)
}

// Test GetMessage
func TestChatService_GetMessage_Success(t *testing.T) {
	service, mockMessageRepo, _, _ := createTestChatService()
	ctx := context.Background()

	expectedMessage := createTestMessage()
	mockMessageRepo.On("GetByID", ctx, "test-message-id").Return(expectedMessage, nil)

	message, err := service.GetMessage(ctx, "test-message-id")

	assert.NoError(t, err)
	assert.Equal(t, expectedMessage, message)

	mockMessageRepo.AssertExpectations(t)
}

func TestChatService_GetMessage_NotFound(t *testing.T) {
	service, mockMessageRepo, _, _ := createTestChatService()
	ctx := context.Background()

	mockMessageRepo.On("GetByID", ctx, "nonexistent-id").Return(nil, nil)

	message, err := service.GetMessage(ctx, "nonexistent-id")

	assert.Error(t, err)
	assert.Nil(t, message)
	var chatErr *interfaces.ChatError
	assert.ErrorAs(t, err, &chatErr)
	assert.Equal(t, interfaces.ErrCodeValidationFailed, chatErr.Code)

	mockMessageRepo.AssertExpectations(t)
}

func TestChatService_GetMessage_EmptyID(t *testing.T) {
	service, _, _, _ := createTestChatService()
	ctx := context.Background()

	message, err := service.GetMessage(ctx, "")

	assert.Error(t, err)
	assert.Nil(t, message)
	var chatErr *interfaces.ChatError
	assert.ErrorAs(t, err, &chatErr)
	assert.Equal(t, interfaces.ErrCodeValidationFailed, chatErr.Code)
}

// Test GetSessionHistory
func TestChatService_GetSessionHistory_Success(t *testing.T) {
	service, mockMessageRepo, mockSessionService, _ := createTestChatService()
	ctx := context.Background()

	session := createTestSession()
	messages := []*domain.ChatMessage{createTestMessage()}

	mockSessionService.On("GetSession", ctx, "test-session-id").Return(session, nil)
	mockMessageRepo.On("GetLatestBySessionID", ctx, "test-session-id", 50).Return(messages, nil)

	result, err := service.GetSessionHistory(ctx, "test-session-id", 50)

	assert.NoError(t, err)
	assert.Equal(t, messages, result)

	mockMessageRepo.AssertExpectations(t)
	mockSessionService.AssertExpectations(t)
}

func TestChatService_GetSessionHistory_EmptySessionID(t *testing.T) {
	service, _, _, _ := createTestChatService()
	ctx := context.Background()

	result, err := service.GetSessionHistory(ctx, "", 50)

	assert.Error(t, err)
	assert.Nil(t, result)
	var chatErr *interfaces.ChatError
	assert.ErrorAs(t, err, &chatErr)
	assert.Equal(t, interfaces.ErrCodeValidationFailed, chatErr.Code)
}

// Test UpdateMessageStatus
func TestChatService_UpdateMessageStatus_Success(t *testing.T) {
	service, mockMessageRepo, _, _ := createTestChatService()
	ctx := context.Background()

	mockMessageRepo.On("UpdateStatus", ctx, "test-message-id", domain.MessageStatusDelivered).Return(nil)

	err := service.UpdateMessageStatus(ctx, "test-message-id", domain.MessageStatusDelivered)

	assert.NoError(t, err)

	mockMessageRepo.AssertExpectations(t)
}

func TestChatService_UpdateMessageStatus_EmptyID(t *testing.T) {
	service, _, _, _ := createTestChatService()
	ctx := context.Background()

	err := service.UpdateMessageStatus(ctx, "", domain.MessageStatusDelivered)

	assert.Error(t, err)
	var chatErr *interfaces.ChatError
	assert.ErrorAs(t, err, &chatErr)
	assert.Equal(t, interfaces.ErrCodeValidationFailed, chatErr.Code)
}

// Test MarkMessageAsDelivered
func TestChatService_MarkMessageAsDelivered_Success(t *testing.T) {
	service, mockMessageRepo, _, _ := createTestChatService()
	ctx := context.Background()

	mockMessageRepo.On("UpdateStatus", ctx, "test-message-id", domain.MessageStatusDelivered).Return(nil)

	err := service.MarkMessageAsDelivered(ctx, "test-message-id")

	assert.NoError(t, err)

	mockMessageRepo.AssertExpectations(t)
}

// Test MarkMessageAsRead
func TestChatService_MarkMessageAsRead_Success(t *testing.T) {
	service, mockMessageRepo, _, _ := createTestChatService()
	ctx := context.Background()

	mockMessageRepo.On("UpdateStatus", ctx, "test-message-id", domain.MessageStatusRead).Return(nil)

	err := service.MarkMessageAsRead(ctx, "test-message-id")

	assert.NoError(t, err)

	mockMessageRepo.AssertExpectations(t)
}

// Test UpdateSessionContext
func TestChatService_UpdateSessionContext_Success(t *testing.T) {
	service, _, mockSessionService, _ := createTestChatService()
	ctx := context.Background()

	session := createTestSession()
	context := &domain.SessionContext{
		ClientName:     "Updated Client",
		MeetingType:    "consultation",
		ProjectContext: "Updated context",
		ServiceTypes:   []string{"migration", "optimization"},
	}

	mockSessionService.On("GetSession", ctx, "test-session-id").Return(session, nil)
	mockSessionService.On("UpdateSession", ctx, mock.AnythingOfType("*domain.ChatSession")).Return(nil)

	err := service.UpdateSessionContext(ctx, "test-session-id", context)

	assert.NoError(t, err)

	mockSessionService.AssertExpectations(t)
}

// Test GetSessionContext
func TestChatService_GetSessionContext_Success(t *testing.T) {
	service, _, mockSessionService, _ := createTestChatService()
	ctx := context.Background()

	session := createTestSession()
	session.Metadata = map[string]interface{}{
		"meeting_type":    "consultation",
		"service_types":   []string{"migration"},
		"cloud_providers": []string{"aws"},
	}

	mockSessionService.On("GetSession", ctx, "test-session-id").Return(session, nil)

	context, err := service.GetSessionContext(ctx, "test-session-id")

	assert.NoError(t, err)
	assert.NotNil(t, context)
	assert.Equal(t, session.ClientName, context.ClientName)
	assert.Equal(t, session.Context, context.ProjectContext)

	mockSessionService.AssertExpectations(t)
}

// Test ValidateMessage
func TestChatService_ValidateMessage_Success(t *testing.T) {
	service, _, _, _ := createTestChatService()

	message := createTestMessage()

	err := service.ValidateMessage(message)

	assert.NoError(t, err)
}

func TestChatService_ValidateMessage_NilMessage(t *testing.T) {
	service, _, _, _ := createTestChatService()

	err := service.ValidateMessage(nil)

	assert.Error(t, err)
	var validationErr *interfaces.MessageValidationError
	assert.ErrorAs(t, err, &validationErr)
	assert.Equal(t, interfaces.ErrCodeValidationFailed, validationErr.Code)
}

// Test SanitizeMessageContent
func TestChatService_SanitizeMessageContent(t *testing.T) {
	service, _, _, _ := createTestChatService()

	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "Hello <script>alert('xss')</script> world",
			expected: "Hello  world",
		},
		{
			input:    "Normal message content",
			expected: "Normal message content",
		},
		{
			input:    "  Multiple   spaces   ",
			expected: "Multiple spaces",
		},
		{
			input:    "HTML <b>bold</b> text",
			expected: "HTML &lt;b&gt;bold&lt;/b&gt; text",
		},
	}

	for _, test := range tests {
		result := service.SanitizeMessageContent(test.input)
		assert.Equal(t, test.expected, result, "Input: %s", test.input)
	}
}

// Test SearchMessages
func TestChatService_SearchMessages_Success(t *testing.T) {
	service, mockMessageRepo, _, _ := createTestChatService()
	ctx := context.Background()

	expectedMessages := []*domain.ChatMessage{createTestMessage()}
	mockMessageRepo.On("Search", ctx, "test-session-id", "test query", 50).Return(expectedMessages, nil)

	messages, err := service.SearchMessages(ctx, "test-session-id", "test query", 50)

	assert.NoError(t, err)
	assert.Equal(t, expectedMessages, messages)

	mockMessageRepo.AssertExpectations(t)
}

func TestChatService_SearchMessages_EmptyQuery(t *testing.T) {
	service, _, _, _ := createTestChatService()
	ctx := context.Background()

	messages, err := service.SearchMessages(ctx, "test-session-id", "", 50)

	assert.Error(t, err)
	assert.Nil(t, messages)
	var chatErr *interfaces.ChatError
	assert.ErrorAs(t, err, &chatErr)
	assert.Equal(t, interfaces.ErrCodeValidationFailed, chatErr.Code)
}

// Test GetMessagesByType
func TestChatService_GetMessagesByType_Success(t *testing.T) {
	service, mockMessageRepo, _, _ := createTestChatService()
	ctx := context.Background()

	expectedMessages := []*domain.ChatMessage{createTestMessage()}
	mockMessageRepo.On("GetByType", ctx, "test-session-id", domain.MessageTypeUser).Return(expectedMessages, nil)

	messages, err := service.GetMessagesByType(ctx, "test-session-id", domain.MessageTypeUser)

	assert.NoError(t, err)
	assert.Equal(t, expectedMessages, messages)

	mockMessageRepo.AssertExpectations(t)
}

// Test GetMessageCount
func TestChatService_GetMessageCount_Success(t *testing.T) {
	service, mockMessageRepo, _, _ := createTestChatService()
	ctx := context.Background()

	filters := &domain.ChatMessageFilters{SessionID: "test-session-id"}
	mockMessageRepo.On("Count", ctx, filters).Return(int64(10), nil)

	count, err := service.GetMessageCount(ctx, filters)

	assert.NoError(t, err)
	assert.Equal(t, int64(10), count)

	mockMessageRepo.AssertExpectations(t)
}

// Test GetMessageStats
func TestChatService_GetMessageStats_Success(t *testing.T) {
	service, mockMessageRepo, _, _ := createTestChatService()
	ctx := context.Background()

	// Mock the various count calls
	mockMessageRepo.On("Count", ctx, &domain.ChatMessageFilters{SessionID: "test-session-id"}).Return(int64(100), nil)
	mockMessageRepo.On("Count", ctx, &domain.ChatMessageFilters{SessionID: "test-session-id", Type: domain.MessageTypeUser}).Return(int64(50), nil)
	mockMessageRepo.On("Count", ctx, &domain.ChatMessageFilters{SessionID: "test-session-id", Type: domain.MessageTypeAssistant}).Return(int64(45), nil)
	mockMessageRepo.On("Count", ctx, &domain.ChatMessageFilters{SessionID: "test-session-id", Type: domain.MessageTypeSystem}).Return(int64(5), nil)
	mockMessageRepo.On("Count", ctx, &domain.ChatMessageFilters{SessionID: "test-session-id", Type: domain.MessageTypeError}).Return(int64(0), nil)

	mockMessageRepo.On("Count", ctx, &domain.ChatMessageFilters{SessionID: "test-session-id", Status: domain.MessageStatusSent}).Return(int64(80), nil)
	mockMessageRepo.On("Count", ctx, &domain.ChatMessageFilters{SessionID: "test-session-id", Status: domain.MessageStatusDelivered}).Return(int64(15), nil)
	mockMessageRepo.On("Count", ctx, &domain.ChatMessageFilters{SessionID: "test-session-id", Status: domain.MessageStatusRead}).Return(int64(5), nil)
	mockMessageRepo.On("Count", ctx, &domain.ChatMessageFilters{SessionID: "test-session-id", Status: domain.MessageStatusFailed}).Return(int64(0), nil)

	// Mock first and last message queries
	firstMessage := createTestMessage()
	mockMessageRepo.On("GetBySessionID", ctx, "test-session-id", 1, 0).Return([]*domain.ChatMessage{firstMessage}, nil)

	lastMessage := createTestMessage()
	mockMessageRepo.On("GetLatestBySessionID", ctx, "test-session-id", 1).Return([]*domain.ChatMessage{lastMessage}, nil)

	stats, err := service.GetMessageStats(ctx, "test-session-id")

	assert.NoError(t, err)
	assert.NotNil(t, stats)
	assert.Equal(t, int64(100), stats.TotalMessages)
	assert.Equal(t, 50, stats.MessagesByType[domain.MessageTypeUser])
	assert.Equal(t, 45, stats.MessagesByType[domain.MessageTypeAssistant])
	assert.Equal(t, 80, stats.MessagesByStatus[domain.MessageStatusSent])
	assert.NotNil(t, stats.FirstMessageAt)
	assert.NotNil(t, stats.LastMessageAt)

	mockMessageRepo.AssertExpectations(t)
}

func TestChatService_GetMessageStats_EmptySessionID(t *testing.T) {
	service, _, _, _ := createTestChatService()
	ctx := context.Background()

	stats, err := service.GetMessageStats(ctx, "")

	assert.Error(t, err)
	assert.Nil(t, stats)
	var chatErr *interfaces.ChatError
	assert.ErrorAs(t, err, &chatErr)
	assert.Equal(t, interfaces.ErrCodeValidationFailed, chatErr.Code)
}
