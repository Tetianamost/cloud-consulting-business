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

// MockChatSessionRepository is a mock implementation of ChatSessionRepository
type MockChatSessionRepository struct {
	mock.Mock
}

func (m *MockChatSessionRepository) Create(ctx context.Context, session *domain.ChatSession) error {
	args := m.Called(ctx, session)
	return args.Error(0)
}

func (m *MockChatSessionRepository) GetByID(ctx context.Context, id string) (*domain.ChatSession, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ChatSession), args.Error(1)
}

func (m *MockChatSessionRepository) Update(ctx context.Context, session *domain.ChatSession) error {
	args := m.Called(ctx, session)
	return args.Error(0)
}

func (m *MockChatSessionRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockChatSessionRepository) GetByUserID(ctx context.Context, userID string) ([]*domain.ChatSession, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.ChatSession), args.Error(1)
}

func (m *MockChatSessionRepository) GetActiveByUserID(ctx context.Context, userID string) ([]*domain.ChatSession, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.ChatSession), args.Error(1)
}

func (m *MockChatSessionRepository) List(ctx context.Context, filters *domain.ChatSessionFilters) ([]*domain.ChatSession, error) {
	args := m.Called(ctx, filters)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.ChatSession), args.Error(1)
}

func (m *MockChatSessionRepository) Count(ctx context.Context, filters *domain.ChatSessionFilters) (int64, error) {
	args := m.Called(ctx, filters)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockChatSessionRepository) GetExpiredSessions(ctx context.Context) ([]*domain.ChatSession, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.ChatSession), args.Error(1)
}

func (m *MockChatSessionRepository) GetInactiveSessions(ctx context.Context, threshold time.Duration) ([]*domain.ChatSession, error) {
	args := m.Called(ctx, threshold)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.ChatSession), args.Error(1)
}

func (m *MockChatSessionRepository) DeleteExpiredSessions(ctx context.Context) (int, error) {
	args := m.Called(ctx)
	return args.Int(0), args.Error(1)
}

func (m *MockChatSessionRepository) DeleteInactiveSessions(ctx context.Context, threshold time.Duration) (int, error) {
	args := m.Called(ctx, threshold)
	return args.Int(0), args.Error(1)
}

func (m *MockChatSessionRepository) UpdateStatus(ctx context.Context, sessionID string, status domain.SessionStatus) error {
	args := m.Called(ctx, sessionID, status)
	return args.Error(0)
}

func (m *MockChatSessionRepository) UpdateLastActivity(ctx context.Context, sessionID string) error {
	args := m.Called(ctx, sessionID)
	return args.Error(0)
}

func (m *MockChatSessionRepository) SetExpiration(ctx context.Context, sessionID string, expiresAt time.Time) error {
	args := m.Called(ctx, sessionID, expiresAt)
	return args.Error(0)
}

// Test helper functions
func createTestSession() *domain.ChatSession {
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

func createTestSessionService() (interfaces.SessionService, *MockChatSessionRepository) {
	mockRepo := &MockChatSessionRepository{}
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel) // Reduce log noise in tests

	service := NewSessionService(mockRepo, logger)
	return service, mockRepo
}

// Test CreateSession
func TestSessionService_CreateSession_Success(t *testing.T) {
	service, mockRepo := createTestSessionService()
	ctx := context.Background()

	session := &domain.ChatSession{
		UserID:     "test-user-id",
		ClientName: "Test Client",
		Context:    "Test context",
		Metadata:   map[string]interface{}{"test": "value"},
	}

	mockRepo.On("Create", ctx, mock.AnythingOfType("*domain.ChatSession")).Return(nil)

	err := service.CreateSession(ctx, session)

	assert.NoError(t, err)
	assert.NotEmpty(t, session.ID)
	assert.Equal(t, domain.SessionStatusActive, session.Status)
	assert.NotZero(t, session.CreatedAt)
	assert.NotZero(t, session.UpdatedAt)
	assert.NotZero(t, session.LastActivity)
	assert.NotNil(t, session.ExpiresAt)

	mockRepo.AssertExpectations(t)
}

func TestSessionService_CreateSession_ValidationError(t *testing.T) {
	service, _ := createTestSessionService()
	ctx := context.Background()

	// Test with invalid session (empty UserID)
	session := &domain.ChatSession{
		ClientName: "Test Client",
	}

	err := service.CreateSession(ctx, session)

	assert.Error(t, err)
	var chatErr *interfaces.ChatError
	assert.ErrorAs(t, err, &chatErr)
	assert.Equal(t, interfaces.ErrCodeValidationFailed, chatErr.Code)
}

func TestSessionService_CreateSession_RepositoryError(t *testing.T) {
	service, mockRepo := createTestSessionService()
	ctx := context.Background()

	session := &domain.ChatSession{
		UserID:     "test-user-id",
		ClientName: "Test Client",
	}

	mockRepo.On("Create", ctx, mock.AnythingOfType("*domain.ChatSession")).Return(errors.New("database error"))

	err := service.CreateSession(ctx, session)

	assert.Error(t, err)
	var chatErr *interfaces.ChatError
	assert.ErrorAs(t, err, &chatErr)
	assert.Equal(t, interfaces.ErrCodeDatabaseError, chatErr.Code)

	mockRepo.AssertExpectations(t)
}

// Test GetSession
func TestSessionService_GetSession_Success(t *testing.T) {
	service, mockRepo := createTestSessionService()
	ctx := context.Background()

	expectedSession := createTestSession()
	mockRepo.On("GetByID", ctx, "test-session-id").Return(expectedSession, nil)

	session, err := service.GetSession(ctx, "test-session-id")

	assert.NoError(t, err)
	assert.Equal(t, expectedSession, session)

	mockRepo.AssertExpectations(t)
}

func TestSessionService_GetSession_NotFound(t *testing.T) {
	service, mockRepo := createTestSessionService()
	ctx := context.Background()

	mockRepo.On("GetByID", ctx, "nonexistent-id").Return(nil, nil)

	session, err := service.GetSession(ctx, "nonexistent-id")

	assert.Error(t, err)
	assert.Nil(t, session)
	var validationErr *interfaces.SessionValidationError
	assert.ErrorAs(t, err, &validationErr)
	assert.Equal(t, interfaces.ErrCodeSessionNotFound, validationErr.Code)

	mockRepo.AssertExpectations(t)
}

func TestSessionService_GetSession_EmptyID(t *testing.T) {
	service, _ := createTestSessionService()
	ctx := context.Background()

	session, err := service.GetSession(ctx, "")

	assert.Error(t, err)
	assert.Nil(t, session)
	var chatErr *interfaces.ChatError
	assert.ErrorAs(t, err, &chatErr)
	assert.Equal(t, interfaces.ErrCodeValidationFailed, chatErr.Code)
}

// Test UpdateSession
func TestSessionService_UpdateSession_Success(t *testing.T) {
	service, mockRepo := createTestSessionService()
	ctx := context.Background()

	session := createTestSession()
	session.ClientName = "Updated Client"

	mockRepo.On("Update", ctx, session).Return(nil)

	err := service.UpdateSession(ctx, session)

	assert.NoError(t, err)
	assert.NotZero(t, session.UpdatedAt)

	mockRepo.AssertExpectations(t)
}

func TestSessionService_UpdateSession_ValidationError(t *testing.T) {
	service, _ := createTestSessionService()
	ctx := context.Background()

	// Test with session missing ID
	session := &domain.ChatSession{
		UserID: "test-user-id",
	}

	err := service.UpdateSession(ctx, session)

	assert.Error(t, err)
	var chatErr *interfaces.ChatError
	assert.ErrorAs(t, err, &chatErr)
	assert.Equal(t, interfaces.ErrCodeValidationFailed, chatErr.Code)
}

// Test DeleteSession
func TestSessionService_DeleteSession_Success(t *testing.T) {
	service, mockRepo := createTestSessionService()
	ctx := context.Background()

	mockRepo.On("Delete", ctx, "test-session-id").Return(nil)

	err := service.DeleteSession(ctx, "test-session-id")

	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestSessionService_DeleteSession_EmptyID(t *testing.T) {
	service, _ := createTestSessionService()
	ctx := context.Background()

	err := service.DeleteSession(ctx, "")

	assert.Error(t, err)
	var chatErr *interfaces.ChatError
	assert.ErrorAs(t, err, &chatErr)
	assert.Equal(t, interfaces.ErrCodeValidationFailed, chatErr.Code)
}

// Test GetUserSessions
func TestSessionService_GetUserSessions_Success(t *testing.T) {
	service, mockRepo := createTestSessionService()
	ctx := context.Background()

	expectedSessions := []*domain.ChatSession{createTestSession()}
	mockRepo.On("GetByUserID", ctx, "test-user-id").Return(expectedSessions, nil)

	sessions, err := service.GetUserSessions(ctx, "test-user-id")

	assert.NoError(t, err)
	assert.Equal(t, expectedSessions, sessions)

	mockRepo.AssertExpectations(t)
}

func TestSessionService_GetUserSessions_EmptyUserID(t *testing.T) {
	service, _ := createTestSessionService()
	ctx := context.Background()

	sessions, err := service.GetUserSessions(ctx, "")

	assert.Error(t, err)
	assert.Nil(t, sessions)
	var chatErr *interfaces.ChatError
	assert.ErrorAs(t, err, &chatErr)
	assert.Equal(t, interfaces.ErrCodeValidationFailed, chatErr.Code)
}

// Test GetActiveSessions
func TestSessionService_GetActiveSessions_Success(t *testing.T) {
	service, mockRepo := createTestSessionService()
	ctx := context.Background()

	// Create active and expired sessions
	activeSession := createTestSession()
	expiredSession := createTestSession()
	expiredSession.ID = "expired-session"
	pastTime := time.Now().Add(-1 * time.Hour)
	expiredSession.ExpiresAt = &pastTime

	allSessions := []*domain.ChatSession{activeSession, expiredSession}
	mockRepo.On("GetActiveByUserID", ctx, "test-user-id").Return(allSessions, nil)

	sessions, err := service.GetActiveSessions(ctx, "test-user-id")

	assert.NoError(t, err)
	assert.Len(t, sessions, 1) // Only the non-expired session should be returned
	assert.Equal(t, activeSession.ID, sessions[0].ID)

	mockRepo.AssertExpectations(t)
}

// Test ValidateSession
func TestSessionService_ValidateSession_Success(t *testing.T) {
	service, mockRepo := createTestSessionService()
	ctx := context.Background()

	session := createTestSession()
	mockRepo.On("GetByID", ctx, "test-session-id").Return(session, nil)

	validatedSession, err := service.ValidateSession(ctx, "test-session-id", "test-user-id")

	assert.NoError(t, err)
	assert.Equal(t, session, validatedSession)

	mockRepo.AssertExpectations(t)
}

func TestSessionService_ValidateSession_WrongUser(t *testing.T) {
	service, mockRepo := createTestSessionService()
	ctx := context.Background()

	session := createTestSession()
	mockRepo.On("GetByID", ctx, "test-session-id").Return(session, nil)

	validatedSession, err := service.ValidateSession(ctx, "test-session-id", "wrong-user-id")

	assert.Error(t, err)
	assert.Nil(t, validatedSession)
	var validationErr *interfaces.SessionValidationError
	assert.ErrorAs(t, err, &validationErr)
	assert.Equal(t, interfaces.ErrCodeUnauthorized, validationErr.Code)

	mockRepo.AssertExpectations(t)
}

func TestSessionService_ValidateSession_Expired(t *testing.T) {
	service, mockRepo := createTestSessionService()
	ctx := context.Background()

	session := createTestSession()
	pastTime := time.Now().Add(-1 * time.Hour)
	session.ExpiresAt = &pastTime

	mockRepo.On("GetByID", ctx, "test-session-id").Return(session, nil)

	validatedSession, err := service.ValidateSession(ctx, "test-session-id", "test-user-id")

	assert.Error(t, err)
	assert.Nil(t, validatedSession)
	var validationErr *interfaces.SessionValidationError
	assert.ErrorAs(t, err, &validationErr)
	assert.Equal(t, interfaces.ErrCodeSessionExpired, validationErr.Code)

	mockRepo.AssertExpectations(t)
}

func TestSessionService_ValidateSession_Inactive(t *testing.T) {
	service, mockRepo := createTestSessionService()
	ctx := context.Background()

	session := createTestSession()
	session.Status = domain.SessionStatusInactive

	mockRepo.On("GetByID", ctx, "test-session-id").Return(session, nil)

	validatedSession, err := service.ValidateSession(ctx, "test-session-id", "test-user-id")

	assert.Error(t, err)
	assert.Nil(t, validatedSession)
	var validationErr *interfaces.SessionValidationError
	assert.ErrorAs(t, err, &validationErr)
	assert.Equal(t, interfaces.ErrCodeSessionInvalid, validationErr.Code)

	mockRepo.AssertExpectations(t)
}

// Test ExpireSession
func TestSessionService_ExpireSession_Success(t *testing.T) {
	service, mockRepo := createTestSessionService()
	ctx := context.Background()

	mockRepo.On("UpdateStatus", ctx, "test-session-id", domain.SessionStatusExpired).Return(nil)

	err := service.ExpireSession(ctx, "test-session-id")

	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

// Test TerminateSession
func TestSessionService_TerminateSession_Success(t *testing.T) {
	service, mockRepo := createTestSessionService()
	ctx := context.Background()

	mockRepo.On("UpdateStatus", ctx, "test-session-id", domain.SessionStatusTerminated).Return(nil)

	err := service.TerminateSession(ctx, "test-session-id")

	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

// Test RefreshSession
func TestSessionService_RefreshSession_Success(t *testing.T) {
	service, mockRepo := createTestSessionService()
	ctx := context.Background()

	duration := 2 * time.Hour
	mockRepo.On("SetExpiration", ctx, "test-session-id", mock.AnythingOfType("time.Time")).Return(nil)
	mockRepo.On("UpdateLastActivity", ctx, "test-session-id").Return(nil)

	err := service.RefreshSession(ctx, "test-session-id", duration)

	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

// Test CleanupExpiredSessions
func TestSessionService_CleanupExpiredSessions_Success(t *testing.T) {
	service, mockRepo := createTestSessionService()
	ctx := context.Background()

	mockRepo.On("DeleteExpiredSessions", ctx).Return(5, nil)

	count, err := service.CleanupExpiredSessions(ctx)

	assert.NoError(t, err)
	assert.Equal(t, 5, count)

	mockRepo.AssertExpectations(t)
}

// Test CleanupInactiveSessions
func TestSessionService_CleanupInactiveSessions_Success(t *testing.T) {
	service, mockRepo := createTestSessionService()
	ctx := context.Background()

	threshold := 7 * 24 * time.Hour
	mockRepo.On("DeleteInactiveSessions", ctx, threshold).Return(3, nil)

	count, err := service.CleanupInactiveSessions(ctx, threshold)

	assert.NoError(t, err)
	assert.Equal(t, 3, count)

	mockRepo.AssertExpectations(t)
}

// Test GetSessionCount
func TestSessionService_GetSessionCount_Success(t *testing.T) {
	service, mockRepo := createTestSessionService()
	ctx := context.Background()

	filters := &domain.ChatSessionFilters{UserID: "test-user-id"}
	mockRepo.On("Count", ctx, filters).Return(int64(10), nil)

	count, err := service.GetSessionCount(ctx, filters)

	assert.NoError(t, err)
	assert.Equal(t, int64(10), count)

	mockRepo.AssertExpectations(t)
}

// Test GetSessionStats
func TestSessionService_GetSessionStats_Success(t *testing.T) {
	service, mockRepo := createTestSessionService()
	ctx := context.Background()

	// Mock the various count calls
	mockRepo.On("Count", ctx, &domain.ChatSessionFilters{}).Return(int64(100), nil)
	mockRepo.On("Count", ctx, &domain.ChatSessionFilters{Status: domain.SessionStatusActive}).Return(int64(50), nil)
	mockRepo.On("Count", ctx, &domain.ChatSessionFilters{Status: domain.SessionStatusExpired}).Return(int64(20), nil)
	mockRepo.On("Count", ctx, &domain.ChatSessionFilters{Status: domain.SessionStatusInactive}).Return(int64(20), nil)
	mockRepo.On("Count", ctx, &domain.ChatSessionFilters{Status: domain.SessionStatusTerminated}).Return(int64(10), nil)

	stats, err := service.GetSessionStats(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, stats)
	assert.Equal(t, int64(100), stats.TotalSessions)
	assert.Equal(t, int64(50), stats.ActiveSessions)
	assert.Equal(t, int64(20), stats.ExpiredSessions)
	assert.Equal(t, 50, stats.SessionsByStatus[domain.SessionStatusActive])
	assert.Equal(t, 20, stats.SessionsByStatus[domain.SessionStatusExpired])

	mockRepo.AssertExpectations(t)
}

// Test IsSessionValid
func TestSessionService_IsSessionValid_Success(t *testing.T) {
	service, mockRepo := createTestSessionService()
	ctx := context.Background()

	session := createTestSession()
	mockRepo.On("GetByID", ctx, "test-session-id").Return(session, nil)

	isValid, err := service.IsSessionValid(ctx, "test-session-id")

	assert.NoError(t, err)
	assert.True(t, isValid)

	mockRepo.AssertExpectations(t)
}

func TestSessionService_IsSessionValid_NotFound(t *testing.T) {
	service, mockRepo := createTestSessionService()
	ctx := context.Background()

	mockRepo.On("GetByID", ctx, "nonexistent-id").Return(nil, nil)

	isValid, err := service.IsSessionValid(ctx, "nonexistent-id")

	assert.NoError(t, err)
	assert.False(t, isValid)

	mockRepo.AssertExpectations(t)
}

func TestSessionService_IsSessionValid_Expired(t *testing.T) {
	service, mockRepo := createTestSessionService()
	ctx := context.Background()

	session := createTestSession()
	pastTime := time.Now().Add(-1 * time.Hour)
	session.ExpiresAt = &pastTime

	mockRepo.On("GetByID", ctx, "test-session-id").Return(session, nil)

	isValid, err := service.IsSessionValid(ctx, "test-session-id")

	assert.NoError(t, err)
	assert.False(t, isValid)

	mockRepo.AssertExpectations(t)
}
