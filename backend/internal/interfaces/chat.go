package interfaces

import (
	"context"
	"time"

	"github.com/cloud-consulting/backend/internal/domain"
)

// SessionService defines the interface for chat session management
type SessionService interface {
	// Session CRUD operations
	CreateSession(ctx context.Context, session *domain.ChatSession) error
	GetSession(ctx context.Context, sessionID string) (*domain.ChatSession, error)
	UpdateSession(ctx context.Context, session *domain.ChatSession) error
	DeleteSession(ctx context.Context, sessionID string) error

	// Session querying
	GetUserSessions(ctx context.Context, userID string) ([]*domain.ChatSession, error)
	GetActiveSessions(ctx context.Context, userID string) ([]*domain.ChatSession, error)
	ListSessions(ctx context.Context, filters *domain.ChatSessionFilters) ([]*domain.ChatSession, error)

	// Session lifecycle management
	ExpireSession(ctx context.Context, sessionID string) error
	TerminateSession(ctx context.Context, sessionID string) error
	RefreshSession(ctx context.Context, sessionID string, duration time.Duration) error

	// Session validation and security
	ValidateSession(ctx context.Context, sessionID string, userID string) (*domain.ChatSession, error)
	IsSessionValid(ctx context.Context, sessionID string) (bool, error)

	// Session cleanup
	CleanupExpiredSessions(ctx context.Context) (int, error)
	CleanupInactiveSessions(ctx context.Context, inactiveThreshold time.Duration) (int, error)

	// Session statistics
	GetSessionCount(ctx context.Context, filters *domain.ChatSessionFilters) (int64, error)
	GetSessionStats(ctx context.Context) (*SessionStats, error)
}

// ChatService defines the interface for chat message handling
type ChatService interface {
	// Message operations
	SendMessage(ctx context.Context, request *domain.ChatRequest) (*domain.ChatResponse, error)
	GetMessage(ctx context.Context, messageID string) (*domain.ChatMessage, error)
	GetSessionHistory(ctx context.Context, sessionID string, limit int) ([]*domain.ChatMessage, error)
	ListMessages(ctx context.Context, filters *domain.ChatMessageFilters) ([]*domain.ChatMessage, error)

	// Message status tracking
	UpdateMessageStatus(ctx context.Context, messageID string, status domain.MessageStatus) error
	MarkMessageAsDelivered(ctx context.Context, messageID string) error
	MarkMessageAsRead(ctx context.Context, messageID string) error

	// Session context management
	UpdateSessionContext(ctx context.Context, sessionID string, context *domain.SessionContext) error
	GetSessionContext(ctx context.Context, sessionID string) (*domain.SessionContext, error)

	// Message validation and sanitization
	ValidateMessage(message *domain.ChatMessage) error
	SanitizeMessageContent(content string) string

	// Message search and filtering
	SearchMessages(ctx context.Context, sessionID string, query string, limit int) ([]*domain.ChatMessage, error)
	GetMessagesByType(ctx context.Context, sessionID string, messageType domain.MessageType) ([]*domain.ChatMessage, error)

	// Message statistics
	GetMessageCount(ctx context.Context, filters *domain.ChatMessageFilters) (int64, error)
	GetMessageStats(ctx context.Context, sessionID string) (*MessageStats, error)
}

// ChatSessionRepository defines the interface for chat session data access
type ChatSessionRepository interface {
	Create(ctx context.Context, session *domain.ChatSession) error
	GetByID(ctx context.Context, id string) (*domain.ChatSession, error)
	Update(ctx context.Context, session *domain.ChatSession) error
	Delete(ctx context.Context, id string) error

	GetByUserID(ctx context.Context, userID string) ([]*domain.ChatSession, error)
	GetActiveByUserID(ctx context.Context, userID string) ([]*domain.ChatSession, error)
	List(ctx context.Context, filters *domain.ChatSessionFilters) ([]*domain.ChatSession, error)
	Count(ctx context.Context, filters *domain.ChatSessionFilters) (int64, error)

	GetExpiredSessions(ctx context.Context) ([]*domain.ChatSession, error)
	GetInactiveSessions(ctx context.Context, threshold time.Duration) ([]*domain.ChatSession, error)
	DeleteExpiredSessions(ctx context.Context) (int, error)
	DeleteInactiveSessions(ctx context.Context, threshold time.Duration) (int, error)

	UpdateStatus(ctx context.Context, sessionID string, status domain.SessionStatus) error
	UpdateLastActivity(ctx context.Context, sessionID string) error
	SetExpiration(ctx context.Context, sessionID string, expiresAt time.Time) error
}

// ChatMessageRepository defines the interface for chat message data access
type ChatMessageRepository interface {
	Create(ctx context.Context, message *domain.ChatMessage) error
	GetByID(ctx context.Context, id string) (*domain.ChatMessage, error)
	Update(ctx context.Context, message *domain.ChatMessage) error
	Delete(ctx context.Context, id string) error

	GetBySessionID(ctx context.Context, sessionID string, limit int, offset int) ([]*domain.ChatMessage, error)
	List(ctx context.Context, filters *domain.ChatMessageFilters) ([]*domain.ChatMessage, error)
	Count(ctx context.Context, filters *domain.ChatMessageFilters) (int64, error)

	GetByType(ctx context.Context, sessionID string, messageType domain.MessageType) ([]*domain.ChatMessage, error)
	GetByStatus(ctx context.Context, sessionID string, status domain.MessageStatus) ([]*domain.ChatMessage, error)

	UpdateStatus(ctx context.Context, messageID string, status domain.MessageStatus) error
	Search(ctx context.Context, sessionID string, query string, limit int) ([]*domain.ChatMessage, error)

	DeleteBySessionID(ctx context.Context, sessionID string) error
	GetLatestBySessionID(ctx context.Context, sessionID string, limit int) ([]*domain.ChatMessage, error)
}

// Supporting types for chat services

// SessionStats represents statistics about chat sessions
type SessionStats struct {
	TotalSessions          int64                        `json:"total_sessions"`
	ActiveSessions         int64                        `json:"active_sessions"`
	ExpiredSessions        int64                        `json:"expired_sessions"`
	SessionsByStatus       map[domain.SessionStatus]int `json:"sessions_by_status"`
	AverageSessionDuration time.Duration                `json:"average_session_duration"`
	TotalMessages          int64                        `json:"total_messages"`
}

// MessageStats represents statistics about messages in a session
type MessageStats struct {
	TotalMessages       int64                        `json:"total_messages"`
	MessagesByType      map[domain.MessageType]int   `json:"messages_by_type"`
	MessagesByStatus    map[domain.MessageStatus]int `json:"messages_by_status"`
	AverageResponseTime time.Duration                `json:"average_response_time"`
	FirstMessageAt      *time.Time                   `json:"first_message_at"`
	LastMessageAt       *time.Time                   `json:"last_message_at"`
}

// SessionValidationError represents a session validation error
type SessionValidationError struct {
	SessionID string `json:"session_id"`
	Reason    string `json:"reason"`
	Code      string `json:"code"`
}

func (e SessionValidationError) Error() string {
	return e.Reason
}

// MessageValidationError represents a message validation error
type MessageValidationError struct {
	MessageID string `json:"message_id"`
	Field     string `json:"field"`
	Reason    string `json:"reason"`
	Code      string `json:"code"`
}

func (e MessageValidationError) Error() string {
	return e.Reason
}

// ChatError represents a general chat service error
type ChatError struct {
	Operation string `json:"operation"`
	Reason    string `json:"reason"`
	Code      string `json:"code"`
	Cause     error  `json:"cause,omitempty"`
}

func (e ChatError) Error() string {
	if e.Cause != nil {
		return e.Reason + ": " + e.Cause.Error()
	}
	return e.Reason
}

func (e ChatError) Unwrap() error {
	return e.Cause
}

// Error codes for chat operations
const (
	ErrCodeSessionNotFound = "SESSION_NOT_FOUND"
	ErrCodeSessionExpired  = "SESSION_EXPIRED"
	ErrCodeSessionInvalid  = "SESSION_INVALID"
	ErrCodeMessageTooLong  = "MESSAGE_TOO_LONG"
	ErrCodeDatabaseError   = "DATABASE_ERROR"
	ErrCodeAIServiceError  = "AI_SERVICE_ERROR"
)
