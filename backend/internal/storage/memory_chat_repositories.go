package storage

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/cloud-consulting/backend/internal/domain"
	"github.com/cloud-consulting/backend/internal/interfaces"
)

// InMemoryChatSessionRepository implements ChatSessionRepository using in-memory storage
type InMemoryChatSessionRepository struct {
	sessions map[string]*domain.ChatSession
	mutex    sync.RWMutex
	logger   *logrus.Logger
}

// NewInMemoryChatSessionRepository creates a new in-memory chat session repository
func NewInMemoryChatSessionRepository(logger *logrus.Logger) interfaces.ChatSessionRepository {
	return &InMemoryChatSessionRepository{
		sessions: make(map[string]*domain.ChatSession),
		logger:   logger,
	}
}

// Create creates a new chat session
func (r *InMemoryChatSessionRepository) Create(ctx context.Context, session *domain.ChatSession) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.sessions[session.ID]; exists {
		return fmt.Errorf("session with ID %s already exists", session.ID)
	}

	// Deep copy to avoid reference issues
	sessionCopy := *session
	if session.Metadata != nil {
		sessionCopy.Metadata = make(map[string]interface{})
		for k, v := range session.Metadata {
			sessionCopy.Metadata[k] = v
		}
	}

	r.sessions[session.ID] = &sessionCopy
	return nil
}

// GetByID retrieves a chat session by ID
func (r *InMemoryChatSessionRepository) GetByID(ctx context.Context, id string) (*domain.ChatSession, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	session, exists := r.sessions[id]
	if !exists {
		return nil, nil
	}

	// Return a copy to avoid reference issues
	sessionCopy := *session
	if session.Metadata != nil {
		sessionCopy.Metadata = make(map[string]interface{})
		for k, v := range session.Metadata {
			sessionCopy.Metadata[k] = v
		}
	}

	return &sessionCopy, nil
}

// Update updates an existing chat session
func (r *InMemoryChatSessionRepository) Update(ctx context.Context, session *domain.ChatSession) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.sessions[session.ID]; !exists {
		return fmt.Errorf("session with ID %s not found", session.ID)
	}

	// Deep copy to avoid reference issues
	sessionCopy := *session
	if session.Metadata != nil {
		sessionCopy.Metadata = make(map[string]interface{})
		for k, v := range session.Metadata {
			sessionCopy.Metadata[k] = v
		}
	}

	r.sessions[session.ID] = &sessionCopy
	return nil
}

// Delete deletes a chat session by ID
func (r *InMemoryChatSessionRepository) Delete(ctx context.Context, id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.sessions[id]; !exists {
		return fmt.Errorf("session with ID %s not found", id)
	}

	delete(r.sessions, id)
	return nil
}

// GetByUserID retrieves all chat sessions for a user
func (r *InMemoryChatSessionRepository) GetByUserID(ctx context.Context, userID string) ([]*domain.ChatSession, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var sessions []*domain.ChatSession
	for _, session := range r.sessions {
		if session.UserID == userID {
			sessionCopy := *session
			if session.Metadata != nil {
				sessionCopy.Metadata = make(map[string]interface{})
				for k, v := range session.Metadata {
					sessionCopy.Metadata[k] = v
				}
			}
			sessions = append(sessions, &sessionCopy)
		}
	}

	// Sort by created_at descending
	sort.Slice(sessions, func(i, j int) bool {
		return sessions[i].CreatedAt.After(sessions[j].CreatedAt)
	})

	return sessions, nil
}

// GetActiveByUserID retrieves active chat sessions for a user
func (r *InMemoryChatSessionRepository) GetActiveByUserID(ctx context.Context, userID string) ([]*domain.ChatSession, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var sessions []*domain.ChatSession
	for _, session := range r.sessions {
		if session.UserID == userID && session.Status == domain.SessionStatusActive {
			sessionCopy := *session
			if session.Metadata != nil {
				sessionCopy.Metadata = make(map[string]interface{})
				for k, v := range session.Metadata {
					sessionCopy.Metadata[k] = v
				}
			}
			sessions = append(sessions, &sessionCopy)
		}
	}

	// Sort by last_activity descending
	sort.Slice(sessions, func(i, j int) bool {
		return sessions[i].LastActivity.After(sessions[j].LastActivity)
	})

	return sessions, nil
}

// List retrieves chat sessions based on filters
func (r *InMemoryChatSessionRepository) List(ctx context.Context, filters *domain.ChatSessionFilters) ([]*domain.ChatSession, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var sessions []*domain.ChatSession
	for _, session := range r.sessions {
		// Apply filters
		if filters.UserID != "" && session.UserID != filters.UserID {
			continue
		}
		if filters.Status != "" && session.Status != filters.Status {
			continue
		}
		if filters.ClientName != "" && !strings.Contains(strings.ToLower(session.ClientName), strings.ToLower(filters.ClientName)) {
			continue
		}
		if filters.FromDate != nil && session.CreatedAt.Before(*filters.FromDate) {
			continue
		}
		if filters.ToDate != nil && session.CreatedAt.After(*filters.ToDate) {
			continue
		}

		sessionCopy := *session
		if session.Metadata != nil {
			sessionCopy.Metadata = make(map[string]interface{})
			for k, v := range session.Metadata {
				sessionCopy.Metadata[k] = v
			}
		}
		sessions = append(sessions, &sessionCopy)
	}

	// Sort by created_at descending
	sort.Slice(sessions, func(i, j int) bool {
		return sessions[i].CreatedAt.After(sessions[j].CreatedAt)
	})

	// Apply pagination
	start := filters.Offset
	if start > len(sessions) {
		return []*domain.ChatSession{}, nil
	}

	end := start + filters.Limit
	if end > len(sessions) {
		end = len(sessions)
	}

	return sessions[start:end], nil
}

// Count returns the count of chat sessions matching the filters
func (r *InMemoryChatSessionRepository) Count(ctx context.Context, filters *domain.ChatSessionFilters) (int64, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	count := int64(0)
	for _, session := range r.sessions {
		// Apply filters
		if filters.UserID != "" && session.UserID != filters.UserID {
			continue
		}
		if filters.Status != "" && session.Status != filters.Status {
			continue
		}
		if filters.ClientName != "" && !strings.Contains(strings.ToLower(session.ClientName), strings.ToLower(filters.ClientName)) {
			continue
		}
		if filters.FromDate != nil && session.CreatedAt.Before(*filters.FromDate) {
			continue
		}
		if filters.ToDate != nil && session.CreatedAt.After(*filters.ToDate) {
			continue
		}

		count++
	}

	return count, nil
}

// GetExpiredSessions retrieves sessions that have expired
func (r *InMemoryChatSessionRepository) GetExpiredSessions(ctx context.Context) ([]*domain.ChatSession, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var sessions []*domain.ChatSession
	now := time.Now()

	for _, session := range r.sessions {
		if session.ExpiresAt != nil && session.ExpiresAt.Before(now) {
			sessionCopy := *session
			if session.Metadata != nil {
				sessionCopy.Metadata = make(map[string]interface{})
				for k, v := range session.Metadata {
					sessionCopy.Metadata[k] = v
				}
			}
			sessions = append(sessions, &sessionCopy)
		}
	}

	// Sort by expires_at ascending
	sort.Slice(sessions, func(i, j int) bool {
		if sessions[i].ExpiresAt == nil {
			return false
		}
		if sessions[j].ExpiresAt == nil {
			return true
		}
		return sessions[i].ExpiresAt.Before(*sessions[j].ExpiresAt)
	})

	return sessions, nil
}

// GetInactiveSessions retrieves sessions that have been inactive for longer than the threshold
func (r *InMemoryChatSessionRepository) GetInactiveSessions(ctx context.Context, threshold time.Duration) ([]*domain.ChatSession, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var sessions []*domain.ChatSession
	cutoffTime := time.Now().Add(-threshold)

	for _, session := range r.sessions {
		if session.LastActivity.Before(cutoffTime) {
			sessionCopy := *session
			if session.Metadata != nil {
				sessionCopy.Metadata = make(map[string]interface{})
				for k, v := range session.Metadata {
					sessionCopy.Metadata[k] = v
				}
			}
			sessions = append(sessions, &sessionCopy)
		}
	}

	// Sort by last_activity ascending
	sort.Slice(sessions, func(i, j int) bool {
		return sessions[i].LastActivity.Before(sessions[j].LastActivity)
	})

	return sessions, nil
}

// DeleteExpiredSessions deletes sessions that have expired
func (r *InMemoryChatSessionRepository) DeleteExpiredSessions(ctx context.Context) (int, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	count := 0
	now := time.Now()

	for id, session := range r.sessions {
		if session.ExpiresAt != nil && session.ExpiresAt.Before(now) {
			delete(r.sessions, id)
			count++
		}
	}

	return count, nil
}

// DeleteInactiveSessions deletes sessions that have been inactive for longer than the threshold
func (r *InMemoryChatSessionRepository) DeleteInactiveSessions(ctx context.Context, threshold time.Duration) (int, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	count := 0
	cutoffTime := time.Now().Add(-threshold)

	for id, session := range r.sessions {
		if session.LastActivity.Before(cutoffTime) {
			delete(r.sessions, id)
			count++
		}
	}

	return count, nil
}

// UpdateStatus updates the status of a chat session
func (r *InMemoryChatSessionRepository) UpdateStatus(ctx context.Context, sessionID string, status domain.SessionStatus) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	session, exists := r.sessions[sessionID]
	if !exists {
		return fmt.Errorf("session with ID %s not found", sessionID)
	}

	session.Status = status
	session.UpdatedAt = time.Now()
	return nil
}

// UpdateLastActivity updates the last activity timestamp of a session
func (r *InMemoryChatSessionRepository) UpdateLastActivity(ctx context.Context, sessionID string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	session, exists := r.sessions[sessionID]
	if !exists {
		return fmt.Errorf("session with ID %s not found", sessionID)
	}

	now := time.Now()
	session.LastActivity = now
	session.UpdatedAt = now
	return nil
}

// SetExpiration sets the expiration time for a session
func (r *InMemoryChatSessionRepository) SetExpiration(ctx context.Context, sessionID string, expiresAt time.Time) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	session, exists := r.sessions[sessionID]
	if !exists {
		return fmt.Errorf("session with ID %s not found", sessionID)
	}

	session.ExpiresAt = &expiresAt
	session.UpdatedAt = time.Now()
	return nil
}

// InMemoryChatMessageRepository implements ChatMessageRepository using in-memory storage
type InMemoryChatMessageRepository struct {
	messages map[string]*domain.ChatMessage
	mutex    sync.RWMutex
	logger   *logrus.Logger
}

// NewInMemoryChatMessageRepository creates a new in-memory chat message repository
func NewInMemoryChatMessageRepository(logger *logrus.Logger) interfaces.ChatMessageRepository {
	return &InMemoryChatMessageRepository{
		messages: make(map[string]*domain.ChatMessage),
		logger:   logger,
	}
}

// Create creates a new chat message
func (r *InMemoryChatMessageRepository) Create(ctx context.Context, message *domain.ChatMessage) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.messages[message.ID]; exists {
		return fmt.Errorf("message with ID %s already exists", message.ID)
	}

	// Deep copy to avoid reference issues
	messageCopy := *message
	if message.Metadata != nil {
		messageCopy.Metadata = make(map[string]interface{})
		for k, v := range message.Metadata {
			messageCopy.Metadata[k] = v
		}
	}

	r.messages[message.ID] = &messageCopy
	return nil
}

// GetByID retrieves a chat message by ID
func (r *InMemoryChatMessageRepository) GetByID(ctx context.Context, id string) (*domain.ChatMessage, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	message, exists := r.messages[id]
	if !exists {
		return nil, nil
	}

	// Return a copy to avoid reference issues
	messageCopy := *message
	if message.Metadata != nil {
		messageCopy.Metadata = make(map[string]interface{})
		for k, v := range message.Metadata {
			messageCopy.Metadata[k] = v
		}
	}

	return &messageCopy, nil
}

// Update updates an existing chat message
func (r *InMemoryChatMessageRepository) Update(ctx context.Context, message *domain.ChatMessage) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.messages[message.ID]; !exists {
		return fmt.Errorf("message with ID %s not found", message.ID)
	}

	// Deep copy to avoid reference issues
	messageCopy := *message
	if message.Metadata != nil {
		messageCopy.Metadata = make(map[string]interface{})
		for k, v := range message.Metadata {
			messageCopy.Metadata[k] = v
		}
	}

	r.messages[message.ID] = &messageCopy
	return nil
}

// Delete deletes a chat message by ID
func (r *InMemoryChatMessageRepository) Delete(ctx context.Context, id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.messages[id]; !exists {
		return fmt.Errorf("message with ID %s not found", id)
	}

	delete(r.messages, id)
	return nil
}

// GetBySessionID retrieves chat messages for a session with pagination
func (r *InMemoryChatMessageRepository) GetBySessionID(ctx context.Context, sessionID string, limit int, offset int) ([]*domain.ChatMessage, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var messages []*domain.ChatMessage
	for _, message := range r.messages {
		if message.SessionID == sessionID {
			messageCopy := *message
			if message.Metadata != nil {
				messageCopy.Metadata = make(map[string]interface{})
				for k, v := range message.Metadata {
					messageCopy.Metadata[k] = v
				}
			}
			messages = append(messages, &messageCopy)
		}
	}

	// Sort by created_at ascending
	sort.Slice(messages, func(i, j int) bool {
		return messages[i].CreatedAt.Before(messages[j].CreatedAt)
	})

	// Apply pagination
	start := offset
	if start > len(messages) {
		return []*domain.ChatMessage{}, nil
	}

	end := start + limit
	if end > len(messages) {
		end = len(messages)
	}

	return messages[start:end], nil
}

// List retrieves chat messages based on filters
func (r *InMemoryChatMessageRepository) List(ctx context.Context, filters *domain.ChatMessageFilters) ([]*domain.ChatMessage, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var messages []*domain.ChatMessage
	for _, message := range r.messages {
		// Apply filters
		if filters.SessionID != "" && message.SessionID != filters.SessionID {
			continue
		}
		if filters.Type != "" && message.Type != filters.Type {
			continue
		}
		if filters.Status != "" && message.Status != filters.Status {
			continue
		}
		if filters.FromDate != nil && message.CreatedAt.Before(*filters.FromDate) {
			continue
		}
		if filters.ToDate != nil && message.CreatedAt.After(*filters.ToDate) {
			continue
		}

		messageCopy := *message
		if message.Metadata != nil {
			messageCopy.Metadata = make(map[string]interface{})
			for k, v := range message.Metadata {
				messageCopy.Metadata[k] = v
			}
		}
		messages = append(messages, &messageCopy)
	}

	// Sort by created_at ascending
	sort.Slice(messages, func(i, j int) bool {
		return messages[i].CreatedAt.Before(messages[j].CreatedAt)
	})

	// Apply pagination
	start := filters.Offset
	if start > len(messages) {
		return []*domain.ChatMessage{}, nil
	}

	end := start + filters.Limit
	if filters.Limit > 0 && end > len(messages) {
		end = len(messages)
	} else if filters.Limit <= 0 {
		end = len(messages)
	}

	return messages[start:end], nil
}

// Count returns the count of chat messages matching the filters
func (r *InMemoryChatMessageRepository) Count(ctx context.Context, filters *domain.ChatMessageFilters) (int64, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	count := int64(0)
	for _, message := range r.messages {
		// Apply filters
		if filters.SessionID != "" && message.SessionID != filters.SessionID {
			continue
		}
		if filters.Type != "" && message.Type != filters.Type {
			continue
		}
		if filters.Status != "" && message.Status != filters.Status {
			continue
		}
		if filters.FromDate != nil && message.CreatedAt.Before(*filters.FromDate) {
			continue
		}
		if filters.ToDate != nil && message.CreatedAt.After(*filters.ToDate) {
			continue
		}

		count++
	}

	return count, nil
}

// GetByType retrieves messages of a specific type for a session
func (r *InMemoryChatMessageRepository) GetByType(ctx context.Context, sessionID string, messageType domain.MessageType) ([]*domain.ChatMessage, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var messages []*domain.ChatMessage
	for _, message := range r.messages {
		if message.SessionID == sessionID && message.Type == messageType {
			messageCopy := *message
			if message.Metadata != nil {
				messageCopy.Metadata = make(map[string]interface{})
				for k, v := range message.Metadata {
					messageCopy.Metadata[k] = v
				}
			}
			messages = append(messages, &messageCopy)
		}
	}

	// Sort by created_at ascending
	sort.Slice(messages, func(i, j int) bool {
		return messages[i].CreatedAt.Before(messages[j].CreatedAt)
	})

	return messages, nil
}

// GetByStatus retrieves messages with a specific status for a session
func (r *InMemoryChatMessageRepository) GetByStatus(ctx context.Context, sessionID string, status domain.MessageStatus) ([]*domain.ChatMessage, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var messages []*domain.ChatMessage
	for _, message := range r.messages {
		if message.SessionID == sessionID && message.Status == status {
			messageCopy := *message
			if message.Metadata != nil {
				messageCopy.Metadata = make(map[string]interface{})
				for k, v := range message.Metadata {
					messageCopy.Metadata[k] = v
				}
			}
			messages = append(messages, &messageCopy)
		}
	}

	// Sort by created_at ascending
	sort.Slice(messages, func(i, j int) bool {
		return messages[i].CreatedAt.Before(messages[j].CreatedAt)
	})

	return messages, nil
}

// UpdateStatus updates the status of a chat message
func (r *InMemoryChatMessageRepository) UpdateStatus(ctx context.Context, messageID string, status domain.MessageStatus) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	message, exists := r.messages[messageID]
	if !exists {
		return fmt.Errorf("message with ID %s not found", messageID)
	}

	message.Status = status
	return nil
}

// Search searches for messages containing the query text
func (r *InMemoryChatMessageRepository) Search(ctx context.Context, sessionID string, query string, limit int) ([]*domain.ChatMessage, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var messages []*domain.ChatMessage
	queryLower := strings.ToLower(query)

	for _, message := range r.messages {
		if message.SessionID == sessionID && strings.Contains(strings.ToLower(message.Content), queryLower) {
			messageCopy := *message
			if message.Metadata != nil {
				messageCopy.Metadata = make(map[string]interface{})
				for k, v := range message.Metadata {
					messageCopy.Metadata[k] = v
				}
			}
			messages = append(messages, &messageCopy)
		}
	}

	// Sort by created_at descending (most recent first)
	sort.Slice(messages, func(i, j int) bool {
		return messages[i].CreatedAt.After(messages[j].CreatedAt)
	})

	// Apply limit
	if limit > 0 && len(messages) > limit {
		messages = messages[:limit]
	}

	return messages, nil
}

// DeleteBySessionID deletes all messages for a session
func (r *InMemoryChatMessageRepository) DeleteBySessionID(ctx context.Context, sessionID string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	count := 0
	for id, message := range r.messages {
		if message.SessionID == sessionID {
			delete(r.messages, id)
			count++
		}
	}

	r.logger.WithFields(logrus.Fields{
		"session_id":    sessionID,
		"rows_affected": count,
	}).Debug("Deleted messages by session ID")

	return nil
}

// GetLatestBySessionID retrieves the latest messages for a session
func (r *InMemoryChatMessageRepository) GetLatestBySessionID(ctx context.Context, sessionID string, limit int) ([]*domain.ChatMessage, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var messages []*domain.ChatMessage
	for _, message := range r.messages {
		if message.SessionID == sessionID {
			messageCopy := *message
			if message.Metadata != nil {
				messageCopy.Metadata = make(map[string]interface{})
				for k, v := range message.Metadata {
					messageCopy.Metadata[k] = v
				}
			}
			messages = append(messages, &messageCopy)
		}
	}

	// Sort by created_at descending (most recent first)
	sort.Slice(messages, func(i, j int) bool {
		return messages[i].CreatedAt.After(messages[j].CreatedAt)
	})

	// Apply limit
	if limit > 0 && len(messages) > limit {
		messages = messages[:limit]
	}

	return messages, nil
}
