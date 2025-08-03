package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/cloud-consulting/backend/internal/domain"
	"github.com/cloud-consulting/backend/internal/interfaces"
)

// ChatMessageRepositoryImpl implements the ChatMessageRepository interface
type ChatMessageRepositoryImpl struct {
	db     *sql.DB
	logger *logrus.Logger
}

// NewChatMessageRepository creates a new chat message repository
func NewChatMessageRepository(db *sql.DB, logger *logrus.Logger) interfaces.ChatMessageRepository {
	return &ChatMessageRepositoryImpl{
		db:     db,
		logger: logger,
	}
}

// Create creates a new chat message in the database
func (r *ChatMessageRepositoryImpl) Create(ctx context.Context, message *domain.ChatMessage) error {
	query := `
		INSERT INTO chat_messages (id, session_id, type, content, metadata, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	metadataJSON, err := json.Marshal(message.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query,
		message.ID,
		message.SessionID,
		message.Type,
		message.Content,
		metadataJSON,
		message.Status,
		message.CreatedAt,
	)

	if err != nil {
		r.logger.WithError(err).WithField("message_id", message.ID).Error("Failed to create chat message")
		return fmt.Errorf("failed to create chat message: %w", err)
	}

	return nil
}

// GetByID retrieves a chat message by ID
func (r *ChatMessageRepositoryImpl) GetByID(ctx context.Context, id string) (*domain.ChatMessage, error) {
	query := `
		SELECT id, session_id, type, content, metadata, status, created_at
		FROM chat_messages
		WHERE id = $1`

	row := r.db.QueryRowContext(ctx, query, id)

	message := &domain.ChatMessage{}
	var metadataJSON []byte

	err := row.Scan(
		&message.ID,
		&message.SessionID,
		&message.Type,
		&message.Content,
		&metadataJSON,
		&message.Status,
		&message.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		r.logger.WithError(err).WithField("message_id", id).Error("Failed to get chat message")
		return nil, fmt.Errorf("failed to get chat message: %w", err)
	}

	// Unmarshal metadata
	if len(metadataJSON) > 0 {
		if err := json.Unmarshal(metadataJSON, &message.Metadata); err != nil {
			r.logger.WithError(err).WithField("message_id", id).Warn("Failed to unmarshal message metadata")
			message.Metadata = make(map[string]interface{})
		}
	} else {
		message.Metadata = make(map[string]interface{})
	}

	return message, nil
}

// Update updates an existing chat message
func (r *ChatMessageRepositoryImpl) Update(ctx context.Context, message *domain.ChatMessage) error {
	query := `
		UPDATE chat_messages 
		SET session_id = $2, type = $3, content = $4, metadata = $5, status = $6
		WHERE id = $1`

	metadataJSON, err := json.Marshal(message.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	result, err := r.db.ExecContext(ctx, query,
		message.ID,
		message.SessionID,
		message.Type,
		message.Content,
		metadataJSON,
		message.Status,
	)

	if err != nil {
		r.logger.WithError(err).WithField("message_id", message.ID).Error("Failed to update chat message")
		return fmt.Errorf("failed to update chat message: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("chat message not found")
	}

	return nil
}

// Delete deletes a chat message by ID
func (r *ChatMessageRepositoryImpl) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM chat_messages WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		r.logger.WithError(err).WithField("message_id", id).Error("Failed to delete chat message")
		return fmt.Errorf("failed to delete chat message: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("chat message not found")
	}

	return nil
}

// GetBySessionID retrieves chat messages for a session with pagination (optimized)
func (r *ChatMessageRepositoryImpl) GetBySessionID(ctx context.Context, sessionID string, limit int, offset int) ([]*domain.ChatMessage, error) {
	// Use optimized query with proper index utilization
	query := `
		SELECT id, session_id, type, content, metadata, status, created_at
		FROM chat_messages
		WHERE session_id = $1
		ORDER BY created_at ASC, id
		LIMIT $2 OFFSET $3`

	return r.queryMessages(ctx, query, sessionID, limit, offset)
}

// List retrieves chat messages based on filters
func (r *ChatMessageRepositoryImpl) List(ctx context.Context, filters *domain.ChatMessageFilters) ([]*domain.ChatMessage, error) {
	query := `
		SELECT id, session_id, type, content, metadata, status, created_at
		FROM chat_messages`

	var conditions []string
	var args []interface{}
	argIndex := 1

	// Build WHERE conditions
	if filters.SessionID != "" {
		conditions = append(conditions, fmt.Sprintf("session_id = $%d", argIndex))
		args = append(args, filters.SessionID)
		argIndex++
	}

	if filters.Type != "" {
		conditions = append(conditions, fmt.Sprintf("type = $%d", argIndex))
		args = append(args, filters.Type)
		argIndex++
	}

	if filters.Status != "" {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, filters.Status)
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
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	// Add ORDER BY
	query += " ORDER BY created_at ASC"

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

	return r.queryMessages(ctx, query, args...)
}

// Count returns the count of chat messages matching the filters
func (r *ChatMessageRepositoryImpl) Count(ctx context.Context, filters *domain.ChatMessageFilters) (int64, error) {
	query := `SELECT COUNT(*) FROM chat_messages`

	var conditions []string
	var args []interface{}
	argIndex := 1

	// Build WHERE conditions (same as List method)
	if filters.SessionID != "" {
		conditions = append(conditions, fmt.Sprintf("session_id = $%d", argIndex))
		args = append(args, filters.SessionID)
		argIndex++
	}

	if filters.Type != "" {
		conditions = append(conditions, fmt.Sprintf("type = $%d", argIndex))
		args = append(args, filters.Type)
		argIndex++
	}

	if filters.Status != "" {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, filters.Status)
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
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	var count int64
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&count)
	if err != nil {
		r.logger.WithError(err).Error("Failed to count chat messages")
		return 0, fmt.Errorf("failed to count chat messages: %w", err)
	}

	return count, nil
}

// GetByType retrieves messages of a specific type for a session
func (r *ChatMessageRepositoryImpl) GetByType(ctx context.Context, sessionID string, messageType domain.MessageType) ([]*domain.ChatMessage, error) {
	query := `
		SELECT id, session_id, type, content, metadata, status, created_at
		FROM chat_messages
		WHERE session_id = $1 AND type = $2
		ORDER BY created_at ASC`

	return r.queryMessages(ctx, query, sessionID, messageType)
}

// GetByStatus retrieves messages with a specific status for a session
func (r *ChatMessageRepositoryImpl) GetByStatus(ctx context.Context, sessionID string, status domain.MessageStatus) ([]*domain.ChatMessage, error) {
	query := `
		SELECT id, session_id, type, content, metadata, status, created_at
		FROM chat_messages
		WHERE session_id = $1 AND status = $2
		ORDER BY created_at ASC`

	return r.queryMessages(ctx, query, sessionID, status)
}

// UpdateStatus updates the status of a chat message
func (r *ChatMessageRepositoryImpl) UpdateStatus(ctx context.Context, messageID string, status domain.MessageStatus) error {
	query := `UPDATE chat_messages SET status = $2 WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, messageID, status)
	if err != nil {
		r.logger.WithError(err).WithFields(logrus.Fields{
			"message_id": messageID,
			"status":     status,
		}).Error("Failed to update message status")
		return fmt.Errorf("failed to update message status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("chat message not found")
	}

	return nil
}

// Search searches for messages containing the query text (optimized with full-text search)
func (r *ChatMessageRepositoryImpl) Search(ctx context.Context, sessionID string, query string, limit int) ([]*domain.ChatMessage, error) {
	// Use full-text search for better performance
	sqlQuery := `
		SELECT id, session_id, type, content, metadata, status, created_at,
		       ts_rank(to_tsvector('english', content), plainto_tsquery('english', $2)) as rank
		FROM chat_messages
		WHERE session_id = $1 
		  AND to_tsvector('english', content) @@ plainto_tsquery('english', $2)
		ORDER BY rank DESC, created_at DESC
		LIMIT $3`

	rows, err := r.db.QueryContext(ctx, sqlQuery, sessionID, query, limit)
	if err != nil {
		r.logger.WithError(err).Error("Failed to search chat messages")
		return nil, fmt.Errorf("failed to search chat messages: %w", err)
	}
	defer rows.Close()

	var messages []*domain.ChatMessage

	for rows.Next() {
		message := &domain.ChatMessage{}
		var metadataJSON []byte
		var rank float64 // We select rank but don't use it in the struct

		err := rows.Scan(
			&message.ID,
			&message.SessionID,
			&message.Type,
			&message.Content,
			&metadataJSON,
			&message.Status,
			&message.CreatedAt,
			&rank,
		)

		if err != nil {
			r.logger.WithError(err).Error("Failed to scan search result row")
			return nil, fmt.Errorf("failed to scan search result row: %w", err)
		}

		// Unmarshal metadata
		if len(metadataJSON) > 0 {
			if err := json.Unmarshal(metadataJSON, &message.Metadata); err != nil {
				r.logger.WithError(err).WithField("message_id", message.ID).Warn("Failed to unmarshal message metadata")
				message.Metadata = make(map[string]interface{})
			}
		} else {
			message.Metadata = make(map[string]interface{})
		}

		messages = append(messages, message)
	}

	if err := rows.Err(); err != nil {
		r.logger.WithError(err).Error("Error iterating over search result rows")
		return nil, fmt.Errorf("error iterating over search result rows: %w", err)
	}

	return messages, nil
}

// DeleteBySessionID deletes all messages for a session
func (r *ChatMessageRepositoryImpl) DeleteBySessionID(ctx context.Context, sessionID string) error {
	query := `DELETE FROM chat_messages WHERE session_id = $1`

	result, err := r.db.ExecContext(ctx, query, sessionID)
	if err != nil {
		r.logger.WithError(err).WithField("session_id", sessionID).Error("Failed to delete messages by session ID")
		return fmt.Errorf("failed to delete messages by session ID: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	r.logger.WithFields(logrus.Fields{
		"session_id":    sessionID,
		"rows_affected": rowsAffected,
	}).Debug("Deleted messages by session ID")

	return nil
}

// GetLatestBySessionID retrieves the latest messages for a session
func (r *ChatMessageRepositoryImpl) GetLatestBySessionID(ctx context.Context, sessionID string, limit int) ([]*domain.ChatMessage, error) {
	query := `
		SELECT id, session_id, type, content, metadata, status, created_at
		FROM chat_messages
		WHERE session_id = $1
		ORDER BY created_at DESC
		LIMIT $2`

	return r.queryMessages(ctx, query, sessionID, limit)
}

// queryMessages is a helper method to execute message queries and scan results
func (r *ChatMessageRepositoryImpl) queryMessages(ctx context.Context, query string, args ...interface{}) ([]*domain.ChatMessage, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		r.logger.WithError(err).Error("Failed to query chat messages")
		return nil, fmt.Errorf("failed to query chat messages: %w", err)
	}
	defer rows.Close()

	var messages []*domain.ChatMessage

	for rows.Next() {
		message := &domain.ChatMessage{}
		var metadataJSON []byte

		err := rows.Scan(
			&message.ID,
			&message.SessionID,
			&message.Type,
			&message.Content,
			&metadataJSON,
			&message.Status,
			&message.CreatedAt,
		)

		if err != nil {
			r.logger.WithError(err).Error("Failed to scan chat message row")
			return nil, fmt.Errorf("failed to scan chat message row: %w", err)
		}

		// Unmarshal metadata
		if len(metadataJSON) > 0 {
			if err := json.Unmarshal(metadataJSON, &message.Metadata); err != nil {
				r.logger.WithError(err).WithField("message_id", message.ID).Warn("Failed to unmarshal message metadata")
				message.Metadata = make(map[string]interface{})
			}
		} else {
			message.Metadata = make(map[string]interface{})
		}

		messages = append(messages, message)
	}

	if err := rows.Err(); err != nil {
		r.logger.WithError(err).Error("Error iterating over chat message rows")
		return nil, fmt.Errorf("error iterating over chat message rows: %w", err)
	}

	return messages, nil
}
