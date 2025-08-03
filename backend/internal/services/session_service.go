package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/cloud-consulting/backend/internal/domain"
	"github.com/cloud-consulting/backend/internal/interfaces"
)

// SessionServiceImpl implements the SessionService interface
type SessionServiceImpl struct {
	repository interfaces.ChatSessionRepository
	logger     *logrus.Logger

	// Configuration
	defaultSessionDuration time.Duration
	maxSessionDuration     time.Duration
	cleanupInterval        time.Duration
}

// NewSessionService creates a new session service instance
func NewSessionService(
	repository interfaces.ChatSessionRepository,
	logger *logrus.Logger,
) interfaces.SessionService {
	return &SessionServiceImpl{
		repository:             repository,
		logger:                 logger,
		defaultSessionDuration: 24 * time.Hour,     // 24 hours default
		maxSessionDuration:     7 * 24 * time.Hour, // 7 days maximum
		cleanupInterval:        time.Hour,          // Cleanup every hour
	}
}

// CreateSession creates a new chat session with proper validation and initialization
func (s *SessionServiceImpl) CreateSession(ctx context.Context, session *domain.ChatSession) error {
	// Validate the session
	if err := s.validateSessionForCreation(session); err != nil {
		s.logger.WithError(err).WithField("user_id", session.UserID).Error("Session validation failed")
		return &interfaces.ChatError{
			Operation: "CreateSession",
			Reason:    "Session validation failed",
			Code:      interfaces.ErrCodeValidationFailed,
			Cause:     err,
		}
	}

	// Generate ID if not provided
	if session.ID == "" {
		session.ID = uuid.New().String()
	}

	// Set timestamps
	now := time.Now()
	session.CreatedAt = now
	session.UpdatedAt = now
	session.LastActivity = now

	// Set default status if not provided
	if session.Status == "" {
		session.Status = domain.SessionStatusActive
	}

	// Set default expiration if not provided
	if session.ExpiresAt == nil {
		expiresAt := now.Add(s.defaultSessionDuration)
		session.ExpiresAt = &expiresAt
	}

	// Initialize metadata if nil
	if session.Metadata == nil {
		session.Metadata = make(map[string]interface{})
	}

	// Create the session in the repository
	if err := s.repository.Create(ctx, session); err != nil {
		s.logger.WithError(err).WithField("session_id", session.ID).Error("Failed to create session")
		return &interfaces.ChatError{
			Operation: "CreateSession",
			Reason:    "Failed to create session in database",
			Code:      interfaces.ErrCodeDatabaseError,
			Cause:     err,
		}
	}

	s.logger.WithFields(logrus.Fields{
		"session_id": session.ID,
		"user_id":    session.UserID,
		"expires_at": session.ExpiresAt,
	}).Info("Chat session created successfully")

	return nil
}

// GetSession retrieves a session by ID with validation
func (s *SessionServiceImpl) GetSession(ctx context.Context, sessionID string) (*domain.ChatSession, error) {
	if sessionID == "" {
		return nil, &interfaces.ChatError{
			Operation: "GetSession",
			Reason:    "Session ID is required",
			Code:      interfaces.ErrCodeValidationFailed,
		}
	}

	session, err := s.repository.GetByID(ctx, sessionID)
	if err != nil {
		s.logger.WithError(err).WithField("session_id", sessionID).Error("Failed to get session")
		return nil, &interfaces.ChatError{
			Operation: "GetSession",
			Reason:    "Failed to retrieve session",
			Code:      interfaces.ErrCodeDatabaseError,
			Cause:     err,
		}
	}

	if session == nil {
		return nil, &interfaces.SessionValidationError{
			SessionID: sessionID,
			Reason:    "Session not found",
			Code:      interfaces.ErrCodeSessionNotFound,
		}
	}

	return session, nil
}

// UpdateSession updates an existing session with validation
func (s *SessionServiceImpl) UpdateSession(ctx context.Context, session *domain.ChatSession) error {
	if err := s.validateSessionForUpdate(session); err != nil {
		s.logger.WithError(err).WithField("session_id", session.ID).Error("Session validation failed")
		return &interfaces.ChatError{
			Operation: "UpdateSession",
			Reason:    "Session validation failed",
			Code:      interfaces.ErrCodeValidationFailed,
			Cause:     err,
		}
	}

	// Update timestamp
	session.UpdatedAt = time.Now()

	if err := s.repository.Update(ctx, session); err != nil {
		s.logger.WithError(err).WithField("session_id", session.ID).Error("Failed to update session")
		return &interfaces.ChatError{
			Operation: "UpdateSession",
			Reason:    "Failed to update session in database",
			Code:      interfaces.ErrCodeDatabaseError,
			Cause:     err,
		}
	}

	s.logger.WithField("session_id", session.ID).Debug("Session updated successfully")
	return nil
}

// DeleteSession deletes a session by ID
func (s *SessionServiceImpl) DeleteSession(ctx context.Context, sessionID string) error {
	if sessionID == "" {
		return &interfaces.ChatError{
			Operation: "DeleteSession",
			Reason:    "Session ID is required",
			Code:      interfaces.ErrCodeValidationFailed,
		}
	}

	if err := s.repository.Delete(ctx, sessionID); err != nil {
		s.logger.WithError(err).WithField("session_id", sessionID).Error("Failed to delete session")
		return &interfaces.ChatError{
			Operation: "DeleteSession",
			Reason:    "Failed to delete session from database",
			Code:      interfaces.ErrCodeDatabaseError,
			Cause:     err,
		}
	}

	s.logger.WithField("session_id", sessionID).Info("Session deleted successfully")
	return nil
}

// GetUserSessions retrieves all sessions for a user
func (s *SessionServiceImpl) GetUserSessions(ctx context.Context, userID string) ([]*domain.ChatSession, error) {
	if userID == "" {
		return nil, &interfaces.ChatError{
			Operation: "GetUserSessions",
			Reason:    "User ID is required",
			Code:      interfaces.ErrCodeValidationFailed,
		}
	}

	sessions, err := s.repository.GetByUserID(ctx, userID)
	if err != nil {
		s.logger.WithError(err).WithField("user_id", userID).Error("Failed to get user sessions")
		return nil, &interfaces.ChatError{
			Operation: "GetUserSessions",
			Reason:    "Failed to retrieve user sessions",
			Code:      interfaces.ErrCodeDatabaseError,
			Cause:     err,
		}
	}

	return sessions, nil
}

// GetActiveSessions retrieves active sessions for a user
func (s *SessionServiceImpl) GetActiveSessions(ctx context.Context, userID string) ([]*domain.ChatSession, error) {
	if userID == "" {
		return nil, &interfaces.ChatError{
			Operation: "GetActiveSessions",
			Reason:    "User ID is required",
			Code:      interfaces.ErrCodeValidationFailed,
		}
	}

	sessions, err := s.repository.GetActiveByUserID(ctx, userID)
	if err != nil {
		s.logger.WithError(err).WithField("user_id", userID).Error("Failed to get active sessions")
		return nil, &interfaces.ChatError{
			Operation: "GetActiveSessions",
			Reason:    "Failed to retrieve active sessions",
			Code:      interfaces.ErrCodeDatabaseError,
			Cause:     err,
		}
	}

	// Filter out expired sessions
	activeSessions := make([]*domain.ChatSession, 0)
	for _, session := range sessions {
		if !session.IsExpired() {
			activeSessions = append(activeSessions, session)
		}
	}

	return activeSessions, nil
}

// ListSessions retrieves sessions based on filters
func (s *SessionServiceImpl) ListSessions(ctx context.Context, filters *domain.ChatSessionFilters) ([]*domain.ChatSession, error) {
	if filters == nil {
		filters = &domain.ChatSessionFilters{}
	}

	// Set default limit if not provided
	if filters.Limit <= 0 {
		filters.Limit = 50
	}

	sessions, err := s.repository.List(ctx, filters)
	if err != nil {
		s.logger.WithError(err).Error("Failed to list sessions")
		return nil, &interfaces.ChatError{
			Operation: "ListSessions",
			Reason:    "Failed to list sessions",
			Code:      interfaces.ErrCodeDatabaseError,
			Cause:     err,
		}
	}

	return sessions, nil
}

// ExpireSession marks a session as expired
func (s *SessionServiceImpl) ExpireSession(ctx context.Context, sessionID string) error {
	if sessionID == "" {
		return &interfaces.ChatError{
			Operation: "ExpireSession",
			Reason:    "Session ID is required",
			Code:      interfaces.ErrCodeValidationFailed,
		}
	}

	if err := s.repository.UpdateStatus(ctx, sessionID, domain.SessionStatusExpired); err != nil {
		s.logger.WithError(err).WithField("session_id", sessionID).Error("Failed to expire session")
		return &interfaces.ChatError{
			Operation: "ExpireSession",
			Reason:    "Failed to expire session",
			Code:      interfaces.ErrCodeDatabaseError,
			Cause:     err,
		}
	}

	s.logger.WithField("session_id", sessionID).Info("Session expired successfully")
	return nil
}

// TerminateSession marks a session as terminated
func (s *SessionServiceImpl) TerminateSession(ctx context.Context, sessionID string) error {
	if sessionID == "" {
		return &interfaces.ChatError{
			Operation: "TerminateSession",
			Reason:    "Session ID is required",
			Code:      interfaces.ErrCodeValidationFailed,
		}
	}

	if err := s.repository.UpdateStatus(ctx, sessionID, domain.SessionStatusTerminated); err != nil {
		s.logger.WithError(err).WithField("session_id", sessionID).Error("Failed to terminate session")
		return &interfaces.ChatError{
			Operation: "TerminateSession",
			Reason:    "Failed to terminate session",
			Code:      interfaces.ErrCodeDatabaseError,
			Cause:     err,
		}
	}

	s.logger.WithField("session_id", sessionID).Info("Session terminated successfully")
	return nil
}

// RefreshSession extends the session expiration time
func (s *SessionServiceImpl) RefreshSession(ctx context.Context, sessionID string, duration time.Duration) error {
	if sessionID == "" {
		return &interfaces.ChatError{
			Operation: "RefreshSession",
			Reason:    "Session ID is required",
			Code:      interfaces.ErrCodeValidationFailed,
		}
	}

	if duration <= 0 {
		duration = s.defaultSessionDuration
	}

	// Enforce maximum session duration
	if duration > s.maxSessionDuration {
		duration = s.maxSessionDuration
	}

	expiresAt := time.Now().Add(duration)
	if err := s.repository.SetExpiration(ctx, sessionID, expiresAt); err != nil {
		s.logger.WithError(err).WithField("session_id", sessionID).Error("Failed to refresh session")
		return &interfaces.ChatError{
			Operation: "RefreshSession",
			Reason:    "Failed to refresh session",
			Code:      interfaces.ErrCodeDatabaseError,
			Cause:     err,
		}
	}

	// Update last activity
	if err := s.repository.UpdateLastActivity(ctx, sessionID); err != nil {
		s.logger.WithError(err).WithField("session_id", sessionID).Warn("Failed to update last activity")
	}

	s.logger.WithFields(logrus.Fields{
		"session_id": sessionID,
		"expires_at": expiresAt,
	}).Debug("Session refreshed successfully")

	return nil
}

// ValidateSession validates a session for a specific user
func (s *SessionServiceImpl) ValidateSession(ctx context.Context, sessionID string, userID string) (*domain.ChatSession, error) {
	session, err := s.GetSession(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	// Check if session belongs to the user
	if session.UserID != userID {
		return nil, &interfaces.SessionValidationError{
			SessionID: sessionID,
			Reason:    "Session does not belong to the specified user",
			Code:      interfaces.ErrCodeUnauthorized,
		}
	}

	// Check if session is expired
	if session.IsExpired() {
		return nil, &interfaces.SessionValidationError{
			SessionID: sessionID,
			Reason:    "Session has expired",
			Code:      interfaces.ErrCodeSessionExpired,
		}
	}

	// Check if session is active
	if session.Status != domain.SessionStatusActive {
		return nil, &interfaces.SessionValidationError{
			SessionID: sessionID,
			Reason:    fmt.Sprintf("Session is not active (status: %s)", session.Status),
			Code:      interfaces.ErrCodeSessionInvalid,
		}
	}

	return session, nil
}

// IsSessionValid checks if a session is valid without returning the session
func (s *SessionServiceImpl) IsSessionValid(ctx context.Context, sessionID string) (bool, error) {
	session, err := s.repository.GetByID(ctx, sessionID)
	if err != nil {
		return false, err
	}

	if session == nil {
		return false, nil
	}

	return session.IsActive(), nil
}

// CleanupExpiredSessions removes expired sessions from the database
func (s *SessionServiceImpl) CleanupExpiredSessions(ctx context.Context) (int, error) {
	deletedCount, err := s.repository.DeleteExpiredSessions(ctx)
	if err != nil {
		s.logger.WithError(err).Error("Failed to cleanup expired sessions")
		return 0, &interfaces.ChatError{
			Operation: "CleanupExpiredSessions",
			Reason:    "Failed to cleanup expired sessions",
			Code:      interfaces.ErrCodeDatabaseError,
			Cause:     err,
		}
	}

	if deletedCount > 0 {
		s.logger.WithField("deleted_count", deletedCount).Info("Cleaned up expired sessions")
	}

	return deletedCount, nil
}

// CleanupInactiveSessions removes sessions that have been inactive for too long
func (s *SessionServiceImpl) CleanupInactiveSessions(ctx context.Context, inactiveThreshold time.Duration) (int, error) {
	if inactiveThreshold <= 0 {
		inactiveThreshold = 7 * 24 * time.Hour // Default to 7 days
	}

	deletedCount, err := s.repository.DeleteInactiveSessions(ctx, inactiveThreshold)
	if err != nil {
		s.logger.WithError(err).Error("Failed to cleanup inactive sessions")
		return 0, &interfaces.ChatError{
			Operation: "CleanupInactiveSessions",
			Reason:    "Failed to cleanup inactive sessions",
			Code:      interfaces.ErrCodeDatabaseError,
			Cause:     err,
		}
	}

	if deletedCount > 0 {
		s.logger.WithFields(logrus.Fields{
			"deleted_count":      deletedCount,
			"inactive_threshold": inactiveThreshold,
		}).Info("Cleaned up inactive sessions")
	}

	return deletedCount, nil
}

// GetSessionCount returns the count of sessions matching the filters
func (s *SessionServiceImpl) GetSessionCount(ctx context.Context, filters *domain.ChatSessionFilters) (int64, error) {
	if filters == nil {
		filters = &domain.ChatSessionFilters{}
	}

	count, err := s.repository.Count(ctx, filters)
	if err != nil {
		s.logger.WithError(err).Error("Failed to get session count")
		return 0, &interfaces.ChatError{
			Operation: "GetSessionCount",
			Reason:    "Failed to get session count",
			Code:      interfaces.ErrCodeDatabaseError,
			Cause:     err,
		}
	}

	return count, nil
}

// GetSessionStats returns statistics about chat sessions
func (s *SessionServiceImpl) GetSessionStats(ctx context.Context) (*interfaces.SessionStats, error) {
	// Get total sessions
	totalSessions, err := s.repository.Count(ctx, &domain.ChatSessionFilters{})
	if err != nil {
		return nil, &interfaces.ChatError{
			Operation: "GetSessionStats",
			Reason:    "Failed to get total session count",
			Code:      interfaces.ErrCodeDatabaseError,
			Cause:     err,
		}
	}

	// Get active sessions
	activeSessions, err := s.repository.Count(ctx, &domain.ChatSessionFilters{
		Status: domain.SessionStatusActive,
	})
	if err != nil {
		return nil, &interfaces.ChatError{
			Operation: "GetSessionStats",
			Reason:    "Failed to get active session count",
			Code:      interfaces.ErrCodeDatabaseError,
			Cause:     err,
		}
	}

	// Get expired sessions
	expiredSessions, err := s.repository.Count(ctx, &domain.ChatSessionFilters{
		Status: domain.SessionStatusExpired,
	})
	if err != nil {
		return nil, &interfaces.ChatError{
			Operation: "GetSessionStats",
			Reason:    "Failed to get expired session count",
			Code:      interfaces.ErrCodeDatabaseError,
			Cause:     err,
		}
	}

	// Build sessions by status map
	sessionsByStatus := make(map[domain.SessionStatus]int)
	for _, status := range []domain.SessionStatus{
		domain.SessionStatusActive,
		domain.SessionStatusInactive,
		domain.SessionStatusExpired,
		domain.SessionStatusTerminated,
	} {
		count, err := s.repository.Count(ctx, &domain.ChatSessionFilters{Status: status})
		if err != nil {
			s.logger.WithError(err).WithField("status", status).Warn("Failed to get session count by status")
			continue
		}
		sessionsByStatus[status] = int(count)
	}

	stats := &interfaces.SessionStats{
		TotalSessions:    totalSessions,
		ActiveSessions:   activeSessions,
		ExpiredSessions:  expiredSessions,
		SessionsByStatus: sessionsByStatus,
		// Note: AverageSessionDuration and TotalMessages would require more complex queries
		// These could be implemented later with additional repository methods
	}

	return stats, nil
}

// validateSessionForCreation validates a session before creation
func (s *SessionServiceImpl) validateSessionForCreation(session *domain.ChatSession) error {
	if session == nil {
		return fmt.Errorf("session cannot be nil")
	}

	if err := session.Validate(); err != nil {
		return err
	}

	// Additional creation-specific validations
	if session.ExpiresAt != nil && session.ExpiresAt.Before(time.Now()) {
		return fmt.Errorf("expiration time cannot be in the past")
	}

	return nil
}

// validateSessionForUpdate validates a session before update
func (s *SessionServiceImpl) validateSessionForUpdate(session *domain.ChatSession) error {
	if session == nil {
		return fmt.Errorf("session cannot be nil")
	}

	if session.ID == "" {
		return fmt.Errorf("session ID is required for update")
	}

	if err := session.Validate(); err != nil {
		return err
	}

	return nil
}
