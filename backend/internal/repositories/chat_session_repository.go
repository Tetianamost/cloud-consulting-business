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
)

// ChatSessionRepositoryImpl implements the ChatSessionRepository interface
type ChatSessionRepositoryImpl struct {
	db     *sql.DB
	logger *logrus.Logger
}

// NewChatSessionRepository creates a new chat session repository
func NewChatSessionRepository(db *sql.DB, logger *logrus.Logger) interfaces.ChatSessionRepository {
	return &ChatSessionRepositoryImpl{
		db:     db,
		logger: logger,
	}
}

// Create creates a new chat session in the database
func (r *ChatSessionRepositoryImpl) Create(ctx context.Context, session *domain.ChatSession) error {
	query := `
		INSERT INTO chat_sessions (id, user_id, client_name, context, status, metadata, created_at, updated_at, last_activity, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	metadataJSON, err := json.Marshal(session.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query,
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

	return nil
}

// GetByID retrieves a chat session by ID
func (r *ChatSessionRepositoryImpl) GetByID(ctx context.Context, id string) (*domain.ChatSession, error) {
	query := `
		SELECT id, user_id, client_name, context, status, metadata, created_at, updated_at, last_activity, expires_at
		FROM chat_sessions
		WHERE id = $1`

	row := r.db.QueryRowContext(ctx, query, id)

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

	return session, nil
}

// Update updates an existing chat session
func (r *ChatSessionRepositoryImpl) Update(ctx context.Context, session *domain.ChatSession) error {
	query := `
		UPDATE chat_sessions 
		SET user_id = $2, client_name = $3, context = $4, status = $5, metadata = $6, updated_at = $7, last_activity = $8, expires_at = $9
		WHERE id = $1`

	metadataJSON, err := json.Marshal(session.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	result, err := r.db.ExecContext(ctx, query,
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

	return nil
}

// Delete deletes a chat session by ID
func (r *ChatSessionRepositoryImpl) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM chat_sessions WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
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

	return nil
}

// GetByUserID retrieves all chat sessions for a user
func (r *ChatSessionRepositoryImpl) GetByUserID(ctx context.Context, userID string) ([]*domain.ChatSession, error) {
	query := `
		SELECT id, user_id, client_name, context, status, metadata, created_at, updated_at, last_activity, expires_at
		FROM chat_sessions
		WHERE user_id = $1
		ORDER BY created_at DESC`

	return r.querySessions(ctx, query, userID)
}

// GetActiveByUserID retrieves active chat sessions for a user
func (r *ChatSessionRepositoryImpl) GetActiveByUserID(ctx context.Context, userID string) ([]*domain.ChatSession, error) {
	query := `
		SELECT id, user_id, client_name, context, status, metadata, created_at, updated_at, last_activity, expires_at
		FROM chat_sessions
		WHERE user_id = $1 AND status = $2
		ORDER BY last_activity DESC`

	return r.querySessions(ctx, query, userID, domain.SessionStatusActive)
}

// List retrieves chat sessions based on filters
func (r *ChatSessionRepositoryImpl) List(ctx context.Context, filters *domain.ChatSessionFilters) ([]*domain.ChatSession, error) {
	query := `
		SELECT id, user_id, client_name, context, status, metadata, created_at, updated_at, last_activity, expires_at
		FROM chat_sessions`

	var conditions []string
	var args []interface{}
	argIndex := 1

	// Build WHERE conditions
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

	// Add ORDER BY
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
func (r *ChatSessionRepositoryImpl) Count(ctx context.Context, filters *domain.ChatSessionFilters) (int64, error) {
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
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&count)
	if err != nil {
		r.logger.WithError(err).Error("Failed to count chat sessions")
		return 0, fmt.Errorf("failed to count chat sessions: %w", err)
	}

	return count, nil
}

// GetExpiredSessions retrieves sessions that have expired
func (r *ChatSessionRepositoryImpl) GetExpiredSessions(ctx context.Context) ([]*domain.ChatSession, error) {
	query := `
		SELECT id, user_id, client_name, context, status, metadata, created_at, updated_at, last_activity, expires_at
		FROM chat_sessions
		WHERE expires_at IS NOT NULL AND expires_at < NOW()
		ORDER BY expires_at ASC`

	return r.querySessions(ctx, query)
}

// GetInactiveSessions retrieves sessions that have been inactive for longer than the threshold
func (r *ChatSessionRepositoryImpl) GetInactiveSessions(ctx context.Context, threshold time.Duration) ([]*domain.ChatSession, error) {
	query := `
		SELECT id, user_id, client_name, context, status, metadata, created_at, updated_at, last_activity, expires_at
		FROM chat_sessions
		WHERE last_activity < $1
		ORDER BY last_activity ASC`

	cutoffTime := time.Now().Add(-threshold)
	return r.querySessions(ctx, query, cutoffTime)
}

// DeleteExpiredSessions deletes sessions that have expired
func (r *ChatSessionRepositoryImpl) DeleteExpiredSessions(ctx context.Context) (int, error) {
	query := `DELETE FROM chat_sessions WHERE expires_at IS NOT NULL AND expires_at < NOW()`

	result, err := r.db.ExecContext(ctx, query)
	if err != nil {
		r.logger.WithError(err).Error("Failed to delete expired chat sessions")
		return 0, fmt.Errorf("failed to delete expired chat sessions: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return int(rowsAffected), nil
}

// DeleteInactiveSessions deletes sessions that have been inactive for longer than the threshold
func (r *ChatSessionRepositoryImpl) DeleteInactiveSessions(ctx context.Context, threshold time.Duration) (int, error) {
	query := `DELETE FROM chat_sessions WHERE last_activity < $1`

	cutoffTime := time.Now().Add(-threshold)
	result, err := r.db.ExecContext(ctx, query, cutoffTime)
	if err != nil {
		r.logger.WithError(err).Error("Failed to delete inactive chat sessions")
		return 0, fmt.Errorf("failed to delete inactive chat sessions: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return int(rowsAffected), nil
}

// UpdateStatus updates the status of a chat session
func (r *ChatSessionRepositoryImpl) UpdateStatus(ctx context.Context, sessionID string, status domain.SessionStatus) error {
	query := `UPDATE chat_sessions SET status = $2, updated_at = NOW() WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, sessionID, status)
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

	return nil
}

// UpdateLastActivity updates the last activity timestamp of a session
func (r *ChatSessionRepositoryImpl) UpdateLastActivity(ctx context.Context, sessionID string) error {
	query := `UPDATE chat_sessions SET last_activity = NOW(), updated_at = NOW() WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, sessionID)
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

	return nil
}

// SetExpiration sets the expiration time for a session
func (r *ChatSessionRepositoryImpl) SetExpiration(ctx context.Context, sessionID string, expiresAt time.Time) error {
	query := `UPDATE chat_sessions SET expires_at = $2, updated_at = NOW() WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, sessionID, expiresAt)
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

	return nil
}

// querySessions is a helper method to execute session queries and scan results
func (r *ChatSessionRepositoryImpl) querySessions(ctx context.Context, query string, args ...interface{}) ([]*domain.ChatSession, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
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
