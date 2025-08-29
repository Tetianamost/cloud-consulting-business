package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/cloud-consulting/backend/internal/domain"
	"github.com/cloud-consulting/backend/internal/interfaces"
	"github.com/cloud-consulting/backend/internal/storage"
)

// EnhancedChatSessionRepository implements ChatSessionRepository with caching and connection pooling
type EnhancedChatSessionRepository struct {
	pool   *storage.DatabasePool
	cache  *storage.RedisCache
	logger *logrus.Logger
}

// NewEnhancedChatSessionRepository creates a new enhanced chat session repository
func NewEnhancedChatSessionRepository(
	pool *storage.DatabasePool,
	cache *storage.RedisCache,
	logger *logrus.Logger,
) interfaces.ChatSessionRepository {
	return &EnhancedChatSessionRepository{
		pool:   pool,
		cache:  cache,
		logger: logger,
	}
}

// Create creates a new chat session with caching
func (r *EnhancedChatSessionRepository) Create(ctx context.Context, session *domain.ChatSession) error {
	db := r.pool.GetDB()

	query := `
		INSERT INTO chat_sessions (id, user_id, client_name, context, status, metadata, created_at, updated_at, last_activity, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	metadataJSON, err := json.Marshal(session.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	_, err = db.ExecContext(ctx, query,
		session.ID,
		session.UserID,
		session.ClientName,
		session.Context,
		session.Status,
		metadataJSON,
		session.CreatedAt,
		session.UpdatedAt,
		session.LastActivity,
		session.ExpiresAt,
	)

	if err != nil {
		r.logger.WithError(err).WithField("session_id", session.ID).Error("Failed to create chat session")
		return fmt.Errorf("failed to create chat session: %w", err)
	}

	// Cache the session
	if r.cache != nil {
		if err := r.cache.SetSession(ctx, session); err != nil {
			r.logger.WithError(err).WithField("session_id", session.ID).Warn("Failed to cache session after creation")
		}

		// Invalidate user sessions cache
		if err := r.cache.InvalidateUserSessions(ctx, session.UserID); err != nil {
			r.logger.WithError(err).WithField("user_id", session.UserID).Warn("Failed to invalidate user sessions cache")
		}
	}

	r.logger.WithField("session_id", session.ID).Info("Chat session created successfully")
	return nil
}

// GetByID retrieves a chat session by ID with caching
func (r *EnhancedChatSessionRepository) GetByID(ctx context.Context, id string) (*domain.ChatSession, error) {
	// Try cache first
	if r.cache != nil {
		if session, err := r.cache.GetSession(ctx, id); err == nil && session != nil {
			r.logger.WithField("session_id", id).Debug("Session retrieved from cache")
			return session, nil
		}
	}

	// Fallback to database
	db := r.pool.GetDB()
	query := `
		SELECT id, user_id, client_name, context, status, metadata, created_at, updated_at, last_activity, expires_at
		FROM chat_sessions
		WHERE id = $1`

	row := db.QueryRowContext(ctx, query, id)

	session := &domain.ChatSession{}
	var metadataJSON []byte
	var clientName, context sql.NullString
	var expiresAt sql.NullTime

	err := row.Scan(
		&session.ID,
		&session.UserID,
		&clientName,
		&context,
		&session.Status,
		&metadataJSON,
		&session.CreatedAt,
		&session.UpdatedAt,
		&session.LastActivity,
		&expiresAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		r.logger.WithError(err).WithField("session_id", id).Error("Failed to get chat session")
		return nil, fmt.Errorf("failed to get chat session: %w", err)
	}

	// Handle nullable fields
	if clientName.Valid {
		session.ClientName = clientName.String
	}
	if context.Valid {
		session.Context = context.String
	}
	if expiresAt.Valid {
		session.ExpiresAt = &expiresAt.Time
	}

	// Unmarshal metadata
	if len(metadataJSON) > 0 {
		if err := json.Unmarshal(metadataJSON, &session.Metadata); err != nil {
			r.logger.WithError(err).WithField("session_id", id).Warn("Failed to unmarshal session metadata")
			session.Metadata = make(map[string]interface{})
		}
	} else {
		session.Metadata = make(map[string]interface{})
	}

	// Cache the session
	if r.cache != nil {
		if err := r.cache.SetSession(ctx, session); err != nil {
			r.logger.WithError(err).WithField("session_id", id).Warn("Failed to cache session after database retrieval")
		}
	}

	r.logger.WithField("session_id", id).Debug("Session retrieved from database")
	return session, nil
}

// Update updates an existing chat session with cache invalidation
func (r *EnhancedChatSessionRepository) Update(ctx context.Context, session *domain.ChatSession) error {
	db := r.pool.GetDB()

	query := `
		UPDATE chat_sessions 
		SET user_id = $2, client_name = $3, context = $4, status = $5, metadata = $6, updated_at = $7, last_activity = $8, expires_at = $9
		WHERE id = $1`

	metadataJSON, err := json.Marshal(session.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	result, err := db.ExecContext(ctx, query,
		session.ID,
		session.UserID,
		session.ClientName,
		session.Context,
		session.Status,
		metadataJSON,
		session.UpdatedAt,
		session.LastActivity,
		session.ExpiresAt,
	)

	if err != nil {
		r.logger.WithError(err).WithField("session_id", session.ID).Error("Failed to update chat session")
		return fmt.Errorf("failed to update chat session: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("chat session not found")
	}

	// Update cache
	if r.cache != nil {
		if err := r.cache.SetSession(ctx, session); err != nil {
			r.logger.WithError(err).WithField("session_id", session.ID).Warn("Failed to update session cache")
		}

		// Invalidate user sessions cache
		if err := r.cache.InvalidateUserSessions(ctx, session.UserID); err != nil {
			r.logger.WithError(err).WithField("user_id", session.UserID).Warn("Failed to invalidate user sessions cache")
		}
	}

	r.logger.WithField("session_id", session.ID).Info("Chat session updated successfully")
	return nil
}

// Delete deletes a chat session with cache invalidation
func (r *EnhancedChatSessionRepository) Delete(ctx context.Context, id string) error {
	// Get session first to know the user ID for cache invalidation
	session, err := r.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if session == nil {
		return fmt.Errorf("chat session not found")
	}

	db := r.pool.GetDB()
	query := `DELETE FROM chat_sessions WHERE id = $1`

	result, err := db.ExecContext(ctx, query, id)
	if err != nil {
		r.logger.WithError(err).WithField("session_id", id).Error("Failed to delete chat session")
		return fmt.Errorf("failed to delete chat session: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("chat session not found")
	}

	// Remove from cache
	if r.cache != nil {
		if err := r.cache.DeleteSession(ctx, id); err != nil {
			r.logger.WithError(err).WithField("session_id", id).Warn("Failed to delete session from cache")
		}

		// Invalidate user sessions cache
		if err := r.cache.InvalidateUserSessions(ctx, session.UserID); err != nil {
			r.logger.WithError(err).WithField("user_id", session.UserID).Warn("Failed to invalidate user sessions cache")
		}

		// Invalidate session messages cache
		if err := r.cache.InvalidateSessionMessages(ctx, id); err != nil {
			r.logger.WithError(err).WithField("session_id", id).Warn("Failed to invalidate session messages cache")
		}
	}

	r.logger.WithField("session_id", id).Info("Chat session deleted successfully")
	return nil
}

// GetByUserID retrieves all chat sessions for a user with caching
func (r *EnhancedChatSessionRepository) GetByUserID(ctx context.Context, userID string) ([]*domain.ChatSession, error) {
	// Try cache first
	if r.cache != nil {
		if sessions, err := r.cache.GetUserSessions(ctx, userID); err == nil && sessions != nil {
			r.logger.WithField("user_id", userID).Debug("User sessions retrieved from cache")
			return sessions, nil
		}
	}

	// Fallback to database
	query := `
		SELECT id, user_id, client_name, context, status, metadata, created_at, updated_at, last_activity, expires_at
		FROM chat_sessions
		WHERE user_id = $1
		ORDER BY created_at DESC`

	sessions, err := r.querySessions(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	// Cache the result
	if r.cache != nil {
		if err := r.cache.SetUserSessions(ctx, userID, sessions); err != nil {
			r.logger.WithError(err).WithField("user_id", userID).Warn("Failed to cache user sessions")
		}
	}

	r.logger.WithFields(logrus.Fields{
		"user_id":       userID,
		"session_count": len(sessions),
	}).Debug("User sessions retrieved from database")
	return sessions, nil
}

// GetActiveByUserID retrieves active chat sessions for a user
func (r *EnhancedChatSessionRepository) GetActiveByUserID(ctx context.Context, userID string) ([]*domain.ChatSession, error) {
	query := `
		SELECT id, user_id, client_name, context, status, metadata, created_at, updated_at, last_activity, expires_at
		FROM chat_sessions
		WHERE user_id = $1 AND status = $2
		ORDER BY last_activity DESC`

	return r.querySessions(ctx, query, userID, domain.SessionStatusActive)
}

// List retrieves chat sessions based on filters with optimized queries
func (r *EnhancedChatSessionRepository) List(ctx context.Context, filters *domain.ChatSessionFilters) ([]*domain.ChatSession, error) {
	query := `
		SELECT id, user_id, client_name, context, status, metadata, created_at, updated_at, last_activity, expires_at
		FROM chat_sessions`

	var conditions []string
	var args []interface{}
	argIndex := 1

	// Build WHERE conditions with optimized indexing
	if filters.UserID != "" {
		conditions = append(conditions, fmt.Sprintf("user_id = $%d", argIndex))
		args = append(args, filters.UserID)
		argIndex++
	}

	if filters.Status != "" {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, filters.Status)
		argIndex++
	}

	if filters.ClientName != "" {
		conditions = append(conditions, fmt.Sprintf("client_name ILIKE $%d", argIndex))
		args = append(args, "%"+filters.ClientName+"%")
		argIndex++
	}

	if filters.FromDate != nil {
		conditions = append(conditions, fmt.Sprintf("created_at >= $%d", argIndex))
		args = append(args, *filters.FromDate)
		argIndex++
	}

	if filters.ToDate != nil {
		conditions = append(conditions, fmt.Sprintf("created_at <= $%d", argIndex))
		args = append(args, *filters.ToDate)
		argIndex++
	}

	// Add WHERE clause if conditions exist
	if len(conditions) > 0 {
		query += " WHERE " + fmt.Sprintf("%s", conditions[0])
		for i := 1; i < len(conditions); i++ {
			query += " AND " + conditions[i]
		}
	}

	// Add ORDER BY for optimal index usage
	query += " ORDER BY created_at DESC"

	// Add LIMIT and OFFSET
	if filters.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argIndex)
		args = append(args, filters.Limit)
		argIndex++
	}

	if filters.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argIndex)
		args = append(args, filters.Offset)
	}

	return r.querySessions(ctx, query, args...)
}

// Count returns the count of chat sessions matching the filters
func (r *EnhancedChatSessionRepository) Count(ctx context.Context, filters *domain.ChatSessionFilters) (int64, error) {
	db := r.pool.GetDB()
	query := `SELECT COUNT(*) FROM chat_sessions`

	var conditions []string
	var args []interface{}
	argIndex := 1

	// Build WHERE conditions (same as List method)
	if filters.UserID != "" {
		conditions = append(conditions, fmt.Sprintf("user_id = $%d", argIndex))
		args = append(args, filters.UserID)
		argIndex++
	}

	if filters.Status != "" {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, filters.Status)
		argIndex++
	}

	if filters.ClientName != "" {
		conditions = append(conditions, fmt.Sprintf("client_name ILIKE $%d", argIndex))
		args = append(args, "%"+filters.ClientName+"%")
		argIndex++
	}

	if filters.FromDate != nil {
		conditions = append(conditions, fmt.Sprintf("created_at >= $%d", argIndex))
		args = append(args, *filters.FromDate)
		argIndex++
	}

	if filters.ToDate != nil {
		conditions = append(conditions, fmt.Sprintf("created_at <= $%d", argIndex))
		args = append(args, *filters.ToDate)
		argIndex++
	}

	// Add WHERE clause if conditions exist
	if len(conditions) > 0 {
		query += " WHERE " + fmt.Sprintf("%s", conditions[0])
		for i := 1; i < len(conditions); i++ {
			query += " AND " + conditions[i]
		}
	}

	var count int64
	err := db.QueryRowContext(ctx, query, args...).Scan(&count)
	if err != nil {
		r.logger.WithError(err).Error("Failed to count chat sessions")
		return 0, fmt.Errorf("failed to count chat sessions: %w", err)
	}

	return count, nil
}

// GetExpiredSessions retrieves sessions that have expired
func (r *EnhancedChatSessionRepository) GetExpiredSessions(ctx context.Context) ([]*domain.ChatSession, error) {
	query := `
		SELECT id, user_id, client_name, context, status, metadata, created_at, updated_at, last_activity, expires_at
		FROM chat_sessions
		WHERE expires_at IS NOT NULL AND expires_at < NOW()
		ORDER BY expires_at ASC`

	return r.querySessions(ctx, query)
}

// GetInactiveSessions retrieves sessions that have been inactive for longer than the threshold
func (r *EnhancedChatSessionRepository) GetInactiveSessions(ctx context.Context, threshold time.Duration) ([]*domain.ChatSession, error) {
	query := `
		SELECT id, user_id, client_name, context, status, metadata, created_at, updated_at, last_activity, expires_at
		FROM chat_sessions
		WHERE last_activity < $1
		ORDER BY last_activity ASC`

	cutoffTime := time.Now().Add(-threshold)
	return r.querySessions(ctx, query, cutoffTime)
}

// DeleteExpiredSessions deletes sessions that have expired with cache cleanup
func (r *EnhancedChatSessionRepository) DeleteExpiredSessions(ctx context.Context) (int, error) {
	// Get expired sessions first for cache cleanup
	expiredSessions, err := r.GetExpiredSessions(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get expired sessions: %w", err)
	}

	if len(expiredSessions) == 0 {
		return 0, nil
	}

	db := r.pool.GetDB()
	query := `DELETE FROM chat_sessions WHERE expires_at IS NOT NULL AND expires_at < NOW()`

	result, err := db.ExecContext(ctx, query)
	if err != nil {
		r.logger.WithError(err).Error("Failed to delete expired chat sessions")
		return 0, fmt.Errorf("failed to delete expired chat sessions: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	// Clean up cache for deleted sessions
	if r.cache != nil {
		userIDs := make(map[string]bool)
		for _, session := range expiredSessions {
			// Delete session from cache
			if err := r.cache.DeleteSession(ctx, session.ID); err != nil {
				r.logger.WithError(err).WithField("session_id", session.ID).Warn("Failed to delete expired session from cache")
			}

			// Delete session messages from cache
			if err := r.cache.InvalidateSessionMessages(ctx, session.ID); err != nil {
				r.logger.WithError(err).WithField("session_id", session.ID).Warn("Failed to invalidate session messages cache")
			}

			userIDs[session.UserID] = true
		}

		// Invalidate user sessions cache for affected users
		for userID := range userIDs {
			if err := r.cache.InvalidateUserSessions(ctx, userID); err != nil {
				r.logger.WithError(err).WithField("user_id", userID).Warn("Failed to invalidate user sessions cache")
			}
		}
	}

	r.logger.WithField("deleted_count", rowsAffected).Info("Expired chat sessions deleted successfully")
	return int(rowsAffected), nil
}

// DeleteInactiveSessions deletes sessions that have been inactive for longer than the threshold
func (r *EnhancedChatSessionRepository) DeleteInactiveSessions(ctx context.Context, threshold time.Duration) (int, error) {
	// Get inactive sessions first for cache cleanup
	inactiveSessions, err := r.GetInactiveSessions(ctx, threshold)
	if err != nil {
		return 0, fmt.Errorf("failed to get inactive sessions: %w", err)
	}

	if len(inactiveSessions) == 0 {
		return 0, nil
	}

	db := r.pool.GetDB()
	query := `DELETE FROM chat_sessions WHERE last_activity < $1`

	cutoffTime := time.Now().Add(-threshold)
	result, err := db.ExecContext(ctx, query, cutoffTime)
	if err != nil {
		r.logger.WithError(err).Error("Failed to delete inactive chat sessions")
		return 0, fmt.Errorf("failed to delete inactive chat sessions: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	// Clean up cache for deleted sessions
	if r.cache != nil {
		userIDs := make(map[string]bool)
		for _, session := range inactiveSessions {
			// Delete session from cache
			if err := r.cache.DeleteSession(ctx, session.ID); err != nil {
				r.logger.WithError(err).WithField("session_id", session.ID).Warn("Failed to delete inactive session from cache")
			}

			// Delete session messages from cache
			if err := r.cache.InvalidateSessionMessages(ctx, session.ID); err != nil {
				r.logger.WithError(err).WithField("session_id", session.ID).Warn("Failed to invalidate session messages cache")
			}

			userIDs[session.UserID] = true
		}

		// Invalidate user sessions cache for affected users
		for userID := range userIDs {
			if err := r.cache.InvalidateUserSessions(ctx, userID); err != nil {
				r.logger.WithError(err).WithField("user_id", userID).Warn("Failed to invalidate user sessions cache")
			}
		}
	}

	r.logger.WithFields(logrus.Fields{
		"deleted_count": rowsAffected,
		"threshold":     threshold,
	}).Info("Inactive chat sessions deleted successfully")
	return int(rowsAffected), nil
}

// UpdateStatus updates the status of a chat session with cache update
func (r *EnhancedChatSessionRepository) UpdateStatus(ctx context.Context, sessionID string, status domain.SessionStatus) error {
	db := r.pool.GetDB()
	query := `UPDATE chat_sessions SET status = $2, updated_at = NOW() WHERE id = $1`

	result, err := db.ExecContext(ctx, query, sessionID, status)
	if err != nil {
		r.logger.WithError(err).WithFields(logrus.Fields{
			"session_id": sessionID,
			"status":     status,
		}).Error("Failed to update session status")
		return fmt.Errorf("failed to update session status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("chat session not found")
	}

	// Update cache
	if r.cache != nil {
		// Get updated session and cache it
		if session, err := r.GetByID(ctx, sessionID); err == nil && session != nil {
			if err := r.cache.SetSession(ctx, session); err != nil {
				r.logger.WithError(err).WithField("session_id", sessionID).Warn("Failed to update session cache after status update")
			}

			// Invalidate user sessions cache
			if err := r.cache.InvalidateUserSessions(ctx, session.UserID); err != nil {
				r.logger.WithError(err).WithField("user_id", session.UserID).Warn("Failed to invalidate user sessions cache")
			}
		}
	}

	return nil
}

// UpdateLastActivity updates the last activity timestamp of a session
func (r *EnhancedChatSessionRepository) UpdateLastActivity(ctx context.Context, sessionID string) error {
	db := r.pool.GetDB()
	query := `UPDATE chat_sessions SET last_activity = NOW(), updated_at = NOW() WHERE id = $1`

	result, err := db.ExecContext(ctx, query, sessionID)
	if err != nil {
		r.logger.WithError(err).WithField("session_id", sessionID).Error("Failed to update last activity")
		return fmt.Errorf("failed to update last activity: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("chat session not found")
	}

	// Update cache
	if r.cache != nil {
		// Get updated session and cache it
		if session, err := r.GetByID(ctx, sessionID); err == nil && session != nil {
			if err := r.cache.SetSession(ctx, session); err != nil {
				r.logger.WithError(err).WithField("session_id", sessionID).Warn("Failed to update session cache after activity update")
			}
		}
	}

	return nil
}

// SetExpiration sets the expiration time for a session
func (r *EnhancedChatSessionRepository) SetExpiration(ctx context.Context, sessionID string, expiresAt time.Time) error {
	db := r.pool.GetDB()
	query := `UPDATE chat_sessions SET expires_at = $2, updated_at = NOW() WHERE id = $1`

	result, err := db.ExecContext(ctx, query, sessionID, expiresAt)
	if err != nil {
		r.logger.WithError(err).WithFields(logrus.Fields{
			"session_id": sessionID,
			"expires_at": expiresAt,
		}).Error("Failed to set session expiration")
		return fmt.Errorf("failed to set session expiration: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("chat session not found")
	}

	// Update cache
	if r.cache != nil {
		// Get updated session and cache it
		if session, err := r.GetByID(ctx, sessionID); err == nil && session != nil {
			if err := r.cache.SetSession(ctx, session); err != nil {
				r.logger.WithError(err).WithField("session_id", sessionID).Warn("Failed to update session cache after expiration update")
			}
		}
	}

	return nil
}

// querySessions is a helper method to execute session queries and scan results
func (r *EnhancedChatSessionRepository) querySessions(ctx context.Context, query string, args ...interface{}) ([]*domain.ChatSession, error) {
	db := r.pool.GetDB()
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		r.logger.WithError(err).Error("Failed to query chat sessions")
		return nil, fmt.Errorf("failed to query chat sessions: %w", err)
	}
	defer rows.Close()

	var sessions []*domain.ChatSession

	for rows.Next() {
		session := &domain.ChatSession{}
		var metadataJSON []byte
		var clientName, context sql.NullString
		var expiresAt sql.NullTime

		err := rows.Scan(
			&session.ID,
			&session.UserID,
			&clientName,
			&context,
			&session.Status,
			&metadataJSON,
			&session.CreatedAt,
			&session.UpdatedAt,
			&session.LastActivity,
			&expiresAt,
		)

		if err != nil {
			r.logger.WithError(err).Error("Failed to scan chat session row")
			return nil, fmt.Errorf("failed to scan chat session row: %w", err)
		}

		// Handle nullable fields
		if clientName.Valid {
			session.ClientName = clientName.String
		}
		if context.Valid {
			session.Context = context.String
		}
		if expiresAt.Valid {
			session.ExpiresAt = &expiresAt.Time
		}

		// Unmarshal metadata
		if len(metadataJSON) > 0 {
			if err := json.Unmarshal(metadataJSON, &session.Metadata); err != nil {
				r.logger.WithError(err).WithField("session_id", session.ID).Warn("Failed to unmarshal session metadata")
				session.Metadata = make(map[string]interface{})
			}
		} else {
			session.Metadata = make(map[string]interface{})
		}

		sessions = append(sessions, session)
	}

	if err := rows.Err(); err != nil {
		r.logger.WithError(err).Error("Error iterating over chat session rows")
		return nil, fmt.Errorf("error iterating over chat session rows: %w", err)
	}

	return sessions, nil
}
