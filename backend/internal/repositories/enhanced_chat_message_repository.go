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
	"github.com/cloud-consulting/backend/internal/storage"
)

// EnhancedChatMessageRepository implements ChatMessageRepository with caching and connection pooling
type EnhancedChatMessageRepository struct {
	pool   *storage.DatabasePool
	cache  *storage.RedisCache
	logger *logrus.Logger
}

// NewEnhancedChatMessageRepository creates a new enhanced chat message repository
func NewEnhancedChatMessageRepository(
	pool *storage.DatabasePool,
	cache *storage.RedisCache,
	logger *logrus.Logger,
) interfaces.ChatMessageRepository {
	return &EnhancedChatMessageRepository{
		pool:   pool,
		cache:  cache,
		logger: logger,
	}
}

// Create creates a new chat message with caching
func (r *EnhancedChatMessageRepository) Create(ctx context.Context, message *domain.ChatMessage) error {
	db := r.pool.GetDB()

	query := `
		INSERT INTO chat_messages (id, session_id, type, content, metadata, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	metadataJSON, err := json.Marshal(message.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	_, err = db.ExecContext(ctx, query,
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

	// Invalidate session messages cache
	if r.cache != nil {
		if err := r.cache.InvalidateSessionMessages(ctx, message.SessionID); err != nil {
			r.logger.WithError(err).WithField("session_id", message.SessionID).Warn("Failed to invalidate session messages cache")
		}
	}

	r.logger.WithField("message_id", message.ID).Info("Chat message created successfully")
	return nil
}

// GetByID retrieves a chat message by ID
func (r *EnhancedChatMessageRepository) GetByID(ctx context.Context, id string) (*domain.ChatMessage, error) {
	db := r.pool.GetDB()
	query := `
		SELECT id, session_id, type, content, metadata, status, created_at
		FROM chat_messages
		WHERE id = $1`

	row := db.QueryRowContext(ctx, query, id)

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

// Update updates an existing chat message with cache invalidation
func (r *EnhancedChatMessageRepository) Update(ctx context.Context, message *domain.ChatMessage) error {
	db := r.pool.GetDB()

	query := `
		UPDATE chat_messages 
		SET session_id = $2, type = $3, content = $4, metadata = $5, status = $6
		WHERE id = $1`

	metadataJSON, err := json.Marshal(message.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	result, err := db.ExecContext(ctx, query,
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

	// Invalidate session messages cache
	if r.cache != nil {
		if err := r.cache.InvalidateSessionMessages(ctx, message.SessionID); err != nil {
			r.logger.WithError(err).WithField("session_id", message.SessionID).Warn("Failed to invalidate session messages cache")
		}
	}

	r.logger.WithField("message_id", message.ID).Info("Chat message updated successfully")
	return nil
}

// Delete deletes a chat message with cache invalidation
func (r *EnhancedChatMessageRepository) Delete(ctx context.Context, id string) error {
	// Get message first to know the session ID for cache invalidation
	message, err := r.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if message == nil {
		return fmt.Errorf("chat message not found")
	}

	db := r.pool.GetDB()
	query := `DELETE FROM chat_messages WHERE id = $1`

	result, err := db.ExecContext(ctx, query, id)
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

	// Invalidate session messages cache
	if r.cache != nil {
		if err := r.cache.InvalidateSessionMessages(ctx, message.SessionID); err != nil {
			r.logger.WithError(err).WithField("session_id", message.SessionID).Warn("Failed to invalidate session messages cache")
		}
	}

	r.logger.WithField("message_id", id).Info("Chat message deleted successfully")
	return nil
}

// GetBySessionID retrieves chat messages for a session with pagination and caching
func (r *EnhancedChatMessageRepository) GetBySessionID(ctx context.Context, sessionID string, limit int, offset int) ([]*domain.ChatMessage, error) {
	// For recent messages (offset 0), try cache first
	if offset == 0 && r.cache != nil {
		if messages, err := r.cache.GetSessionMessages(ctx, sessionID); err == nil && messages != nil {
			// Return cached messages up to the limit
			if len(messages) >= limit {
				r.logger.WithField("session_id", sessionID).Debug("Session messages retrieved from cache")
				return messages[:limit], nil
			}
		}
	}

	// Fallback to database
	db := r.pool.GetDB()
	query := `
		SELECT id, session_id, type, content, metadata, status, created_at
		FROM chat_messages
		WHERE session_id = $1
		ORDER BY created_at ASC
		LIMIT $2 OFFSET $3`

	messages, err := r.queryMessages(ctx, query, sessionID, limit, offset)
	if err != nil {
		return nil, err
	}

	// Cache recent messages (offset 0) for future requests
	if offset == 0 && r.cache != nil && len(messages) > 0 {
		// Cache up to 50 recent messages
		cacheLimit := 50
		if len(messages) > cacheLimit {
			messages = messages[:cacheLimit]
		}

		if err := r.cache.SetSessionMessages(ctx, sessionID, messages); err != nil {
			r.logger.WithError(err).WithField("session_id", sessionID).Warn("Failed to cache session messages")
		}
	}

	r.logger.WithFields(logrus.Fields{
		"session_id":    sessionID,
		"message_count": len(messages),
		"limit":         limit,
		"offset":        offset,
	}).Debug("Session messages retrieved from database")
	return messages, nil
}

// List retrieves chat messages based on filters with optimized queries
func (r *EnhancedChatMessageRepository) List(ctx context.Context, filters *domain.ChatMessageFilters) ([]*domain.ChatMessage, error) {
	db := r.pool.GetDB()
	query := `
		SELECT id, session_id, type, content, metadata, status, created_at
		FROM chat_messages`

	var conditions []string
	var args []interface{}
	argIndex := 1

	// Build WHERE conditions with optimized indexing
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

	// Add ORDER BY for optimal index usage
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
func (r *EnhancedChatMessageRepository) Count(ctx context.Context, filters *domain.ChatMessageFilters) (int64, error) {
	db := r.pool.GetDB()
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
	err := db.QueryRowContext(ctx, query, args...).Scan(&count)
	if err != nil {
		r.logger.WithError(err).Error("Failed to count chat messages")
		return 0, fmt.Errorf("failed to count chat messages: %w", err)
	}

	return count, nil
}

// GetByType retrieves messages of a specific type for a session
func (r *EnhancedChatMessageRepository) GetByType(ctx context.Context, sessionID string, messageType domain.MessageType) ([]*domain.ChatMessage, error) {
	db := r.pool.GetDB()
	query := `
		SELECT id, session_id, type, content, metadata, status, created_at
		FROM chat_messages
		WHERE session_id = $1 AND type = $2
		ORDER BY created_at ASC`

	return r.queryMessages(ctx, query, sessionID, messageType)
}

// GetByStatus retrieves messages with a specific status for a session
func (r *EnhancedChatMessageRepository) GetByStatus(ctx context.Context, sessionID string, status domain.MessageStatus) ([]*domain.ChatMessage, error) {
	db := r.pool.GetDB()
	query := `
		SELECT id, session_id, type, content, metadata, status, created_at
		FROM chat_messages
		WHERE session_id = $1 AND status = $2
		ORDER BY created_at ASC`

	return r.queryMessages(ctx, query, sessionID, status)
}

// UpdateStatus updates the status of a chat message with cache invalidation
func (r *EnhancedChatMessageRepository) UpdateStatus(ctx context.Context, messageID string, status domain.MessageStatus) error {
	// Get message first to know the session ID for cache invalidation
	message, err := r.GetByID(ctx, messageID)
	if err != nil {
		return err
	}
	if message == nil {
		return fmt.Errorf("chat message not found")
	}

	db := r.pool.GetDB()
	query := `UPDATE chat_messages SET status = $2 WHERE id = $1`

	result, err := db.ExecContext(ctx, query, messageID, status)
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

	// Invalidate session messages cache
	if r.cache != nil {
		if err := r.cache.InvalidateSessionMessages(ctx, message.SessionID); err != nil {
			r.logger.WithError(err).WithField("session_id", message.SessionID).Warn("Failed to invalidate session messages cache")
		}
	}

	return nil
}

// Search searches for messages containing the query text with optimized full-text search
func (r *EnhancedChatMessageRepository) Search(ctx context.Context, sessionID string, query string, limit int) ([]*domain.ChatMessage, error) {
	db := r.pool.GetDB()

	// Use PostgreSQL's ILIKE for case-insensitive search with index support
	sqlQuery := `
		SELECT id, session_id, type, content, metadata, status, created_at
		FROM chat_messages
		WHERE session_id = $1 AND content ILIKE $2
		ORDER BY created_at DESC
		LIMIT $3`

	searchPattern := "%" + query + "%"
	return r.queryMessages(ctx, sqlQuery, sessionID, searchPattern, limit)
}

// DeleteBySessionID deletes all messages for a session with cache cleanup
func (r *EnhancedChatMessageRepository) DeleteBySessionID(ctx context.Context, sessionID string) error {
	db := r.pool.GetDB()
	query := `DELETE FROM chat_messages WHERE session_id = $1`

	result, err := db.ExecContext(ctx, query, sessionID)
	if err != nil {
		r.logger.WithError(err).WithField("session_id", sessionID).Error("Failed to delete messages by session ID")
		return fmt.Errorf("failed to delete messages by session ID: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	// Invalidate session messages cache
	if r.cache != nil {
		if err := r.cache.InvalidateSessionMessages(ctx, sessionID); err != nil {
			r.logger.WithError(err).WithField("session_id", sessionID).Warn("Failed to invalidate session messages cache")
		}
	}

	r.logger.WithFields(logrus.Fields{
		"session_id":    sessionID,
		"rows_affected": rowsAffected,
	}).Info("Messages deleted by session ID successfully")

	return nil
}

// GetLatestBySessionID retrieves the latest messages for a session with caching
func (r *EnhancedChatMessageRepository) GetLatestBySessionID(ctx context.Context, sessionID string, limit int) ([]*domain.ChatMessage, error) {
	// Try cache first for recent messages
	if r.cache != nil {
		if messages, err := r.cache.GetSessionMessages(ctx, sessionID); err == nil && messages != nil {
			// Return the latest messages from cache
			if len(messages) >= limit {
				// Get the last 'limit' messages
				start := len(messages) - limit
				r.logger.WithField("session_id", sessionID).Debug("Latest session messages retrieved from cache")
				return messages[start:], nil
			}
		}
	}

	// Fallback to database
	db := r.pool.GetDB()
	query := `
		SELECT id, session_id, type, content, metadata, status, created_at
		FROM chat_messages
		WHERE session_id = $1
		ORDER BY created_at DESC
		LIMIT $2`

	messages, err := r.queryMessages(ctx, query, sessionID, limit)
	if err != nil {
		return nil, err
	}

	// Reverse the order to get chronological order
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	// Cache the messages for future requests
	if r.cache != nil && len(messages) > 0 {
		if err := r.cache.SetSessionMessages(ctx, sessionID, messages); err != nil {
			r.logger.WithError(err).WithField("session_id", sessionID).Warn("Failed to cache latest session messages")
		}
	}

	r.logger.WithFields(logrus.Fields{
		"session_id":    sessionID,
		"message_count": len(messages),
		"limit":         limit,
	}).Debug("Latest session messages retrieved from database")
	return messages, nil
}

// queryMessages is a helper method to execute message queries and scan results
func (r *EnhancedChatMessageRepository) queryMessages(ctx context.Context, query string, args ...interface{}) ([]*domain.ChatMessage, error) {
	db := r.pool.GetDB()
	rows, err := db.QueryContext(ctx, query, args...)
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
