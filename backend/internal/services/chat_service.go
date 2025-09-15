package services

import (
	"context"
	"fmt"
	"html"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/cloud-consulting/backend/internal/domain"
	"github.com/cloud-consulting/backend/internal/interfaces"
)

// ChatServiceImpl implements the ChatService interface
type ChatServiceImpl struct {
	messageRepository interfaces.ChatMessageRepository
	sessionService    interfaces.SessionService
	aiService         interfaces.BedrockService
	enhancedAIService *ChatAwareBedrockService
	logger            *logrus.Logger

	// Configuration
	maxMessageLength int
	maxHistoryLimit  int
}

// NewChatService creates a new chat service instance
func NewChatService(
	messageRepository interfaces.ChatMessageRepository,
	sessionService interfaces.SessionService,
	aiService interfaces.BedrockService,
	logger *logrus.Logger,
) interfaces.ChatService {
	return &ChatServiceImpl{
		messageRepository: messageRepository,
		sessionService:    sessionService,
		aiService:         aiService,
		logger:            logger,
		maxMessageLength:  10000, // 10KB max message length
		maxHistoryLimit:   1000,  // Max 1000 messages in history
	}
}

// NewChatServiceWithEnhancedAI creates a new chat service instance with enhanced AI capabilities
func NewChatServiceWithEnhancedAI(
	messageRepository interfaces.ChatMessageRepository,
	sessionService interfaces.SessionService,
	aiService interfaces.BedrockService,
	enhancedAIService *ChatAwareBedrockService,
	logger *logrus.Logger,
) interfaces.ChatService {
	return &ChatServiceImpl{
		messageRepository: messageRepository,
		sessionService:    sessionService,
		aiService:         aiService,
		enhancedAIService: enhancedAIService,
		logger:            logger,
		maxMessageLength:  10000, // 10KB max message length
		maxHistoryLimit:   1000,  // Max 1000 messages in history
	}
}

// SendMessage processes a chat request and generates an AI response
func (s *ChatServiceImpl) SendMessage(ctx context.Context, request *domain.ChatRequest) (*domain.ChatResponse, error) {
	// Validate the request
	if err := s.validateChatRequest(request); err != nil {
		s.logger.WithError(err).WithField("session_id", request.SessionID).Error("Chat request validation failed")
		return nil, &interfaces.ChatError{
			Operation: "SendMessage",
			Reason:    "Request validation failed",
			Code:      interfaces.ErrCodeValidationFailed,
			Cause:     err,
		}
	}

	// Validate and get the session
	session, err := s.sessionService.GetSession(ctx, request.SessionID)
	if err != nil {
		s.logger.WithError(err).WithField("session_id", request.SessionID).Error("Failed to get session")
		return nil, &interfaces.ChatError{
			Operation: "SendMessage",
			Reason:    "Failed to retrieve session",
			Code:      interfaces.ErrCodeSessionNotFound,
			Cause:     err,
		}
	}

	// Check if session is active
	if !session.IsActive() {
		return nil, &interfaces.SessionValidationError{
			SessionID: request.SessionID,
			Reason:    "Session is not active",
			Code:      interfaces.ErrCodeSessionInvalid,
		}
	}

	// Sanitize message content
	sanitizedContent := s.SanitizeMessageContent(request.Content)

	// Create user message
	userMessage := &domain.ChatMessage{
		ID:        uuid.New().String(),
		SessionID: request.SessionID,
		Type:      request.Type,
		Content:   sanitizedContent,
		Metadata:  request.Metadata,
		Status:    domain.MessageStatusSent,
		CreatedAt: time.Now(),
	}

	// Initialize metadata if nil
	if userMessage.Metadata == nil {
		userMessage.Metadata = make(map[string]interface{})
	}

	// Store user message
	if err := s.messageRepository.Create(ctx, userMessage); err != nil {
		s.logger.WithError(err).WithField("message_id", userMessage.ID).Error("Failed to store user message")
		return nil, &interfaces.ChatError{
			Operation: "SendMessage",
			Reason:    "Failed to store user message",
			Code:      interfaces.ErrCodeDatabaseError,
			Cause:     err,
		}
	}

	// Update session activity
	session.UpdateActivity()
	if err := s.sessionService.UpdateSession(ctx, session); err != nil {
		s.logger.WithError(err).WithField("session_id", session.ID).Warn("Failed to update session activity")
	}

	// Generate AI response if this is a user message
	var aiResponse *domain.ChatMessage
	if request.Type == domain.MessageTypeUser {
		aiResponse, err = s.generateAIResponse(ctx, session, userMessage, request.QuickAction)
		if err != nil {
			s.logger.WithError(err).WithField("session_id", request.SessionID).Error("Failed to generate AI response")
			// Return the user message even if AI response fails
			return &domain.ChatResponse{
					MessageID:   userMessage.ID,
					SessionID:   userMessage.SessionID,
					Content:     userMessage.Content,
					Type:        userMessage.Type,
					Metadata:    userMessage.Metadata,
					CreatedAt:   userMessage.CreatedAt,
					Status:      userMessage.Status,
					ProcessTime: 0,
				}, &interfaces.ChatError{
					Operation: "SendMessage",
					Reason:    "Failed to generate AI response",
					Code:      interfaces.ErrCodeAIServiceError,
					Cause:     err,
				}
		}

		// Store AI response
		if err := s.messageRepository.Create(ctx, aiResponse); err != nil {
			s.logger.WithError(err).WithField("message_id", aiResponse.ID).Error("Failed to store AI response")
			// Return user message even if storing AI response fails
			return &domain.ChatResponse{
				MessageID:   userMessage.ID,
				SessionID:   userMessage.SessionID,
				Content:     userMessage.Content,
				Type:        userMessage.Type,
				Metadata:    userMessage.Metadata,
				CreatedAt:   userMessage.CreatedAt,
				Status:      userMessage.Status,
				ProcessTime: 0,
			}, nil
		}

		// Mark user message as delivered
		if err := s.messageRepository.UpdateStatus(ctx, userMessage.ID, domain.MessageStatusDelivered); err != nil {
			s.logger.WithError(err).WithField("message_id", userMessage.ID).Warn("Failed to update user message status")
		}

		// Return AI response
		return &domain.ChatResponse{
			MessageID:   aiResponse.ID,
			SessionID:   aiResponse.SessionID,
			Content:     aiResponse.Content,
			Type:        aiResponse.Type,
			Metadata:    aiResponse.Metadata,
			CreatedAt:   aiResponse.CreatedAt,
			Status:      aiResponse.Status,
			ProcessTime: 0, // Could be calculated if needed
		}, nil
	}

	// For non-user messages, just return the stored message
	return &domain.ChatResponse{
		MessageID:   userMessage.ID,
		SessionID:   userMessage.SessionID,
		Content:     userMessage.Content,
		Type:        userMessage.Type,
		Metadata:    userMessage.Metadata,
		CreatedAt:   userMessage.CreatedAt,
		Status:      userMessage.Status,
		ProcessTime: 0,
	}, nil
}

// GetMessage retrieves a specific message by ID
func (s *ChatServiceImpl) GetMessage(ctx context.Context, messageID string) (*domain.ChatMessage, error) {
	if messageID == "" {
		return nil, &interfaces.ChatError{
			Operation: "GetMessage",
			Reason:    "Message ID is required",
			Code:      interfaces.ErrCodeValidationFailed,
		}
	}

	message, err := s.messageRepository.GetByID(ctx, messageID)
	if err != nil {
		s.logger.WithError(err).WithField("message_id", messageID).Error("Failed to get message")
		return nil, &interfaces.ChatError{
			Operation: "GetMessage",
			Reason:    "Failed to retrieve message",
			Code:      interfaces.ErrCodeDatabaseError,
			Cause:     err,
		}
	}

	if message == nil {
		return nil, &interfaces.ChatError{
			Operation: "GetMessage",
			Reason:    "Message not found",
			Code:      interfaces.ErrCodeValidationFailed,
		}
	}

	return message, nil
}

// GetSessionHistory retrieves message history for a session
func (s *ChatServiceImpl) GetSessionHistory(ctx context.Context, sessionID string, limit int) ([]*domain.ChatMessage, error) {
	if sessionID == "" {
		return nil, &interfaces.ChatError{
			Operation: "GetSessionHistory",
			Reason:    "Session ID is required",
			Code:      interfaces.ErrCodeValidationFailed,
		}
	}

	// Validate limit
	if limit <= 0 {
		limit = 50 // Default limit
	}
	if limit > s.maxHistoryLimit {
		limit = s.maxHistoryLimit
	}

	// Verify session exists and is accessible
	_, err := s.sessionService.GetSession(ctx, sessionID)
	if err != nil {
		return nil, &interfaces.ChatError{
			Operation: "GetSessionHistory",
			Reason:    "Session not found or inaccessible",
			Code:      interfaces.ErrCodeSessionNotFound,
			Cause:     err,
		}
	}

	messages, err := s.messageRepository.GetLatestBySessionID(ctx, sessionID, limit)
	if err != nil {
		s.logger.WithError(err).WithField("session_id", sessionID).Error("Failed to get session history")
		return nil, &interfaces.ChatError{
			Operation: "GetSessionHistory",
			Reason:    "Failed to retrieve session history",
			Code:      interfaces.ErrCodeDatabaseError,
			Cause:     err,
		}
	}

	// Reverse the order to show oldest first
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}

// ListMessages retrieves messages based on filters
func (s *ChatServiceImpl) ListMessages(ctx context.Context, filters *domain.ChatMessageFilters) ([]*domain.ChatMessage, error) {
	if filters == nil {
		filters = &domain.ChatMessageFilters{}
	}

	// Set default limit if not provided
	if filters.Limit <= 0 {
		filters.Limit = 50
	}
	if filters.Limit > s.maxHistoryLimit {
		filters.Limit = s.maxHistoryLimit
	}

	messages, err := s.messageRepository.List(ctx, filters)
	if err != nil {
		s.logger.WithError(err).Error("Failed to list messages")
		return nil, &interfaces.ChatError{
			Operation: "ListMessages",
			Reason:    "Failed to list messages",
			Code:      interfaces.ErrCodeDatabaseError,
			Cause:     err,
		}
	}

	return messages, nil
}

// UpdateMessageStatus updates the status of a message
func (s *ChatServiceImpl) UpdateMessageStatus(ctx context.Context, messageID string, status domain.MessageStatus) error {
	if messageID == "" {
		return &interfaces.ChatError{
			Operation: "UpdateMessageStatus",
			Reason:    "Message ID is required",
			Code:      interfaces.ErrCodeValidationFailed,
		}
	}

	if err := s.messageRepository.UpdateStatus(ctx, messageID, status); err != nil {
		s.logger.WithError(err).WithFields(logrus.Fields{
			"message_id": messageID,
			"status":     status,
		}).Error("Failed to update message status")
		return &interfaces.ChatError{
			Operation: "UpdateMessageStatus",
			Reason:    "Failed to update message status",
			Code:      interfaces.ErrCodeDatabaseError,
			Cause:     err,
		}
	}

	return nil
}

// MarkMessageAsDelivered marks a message as delivered
func (s *ChatServiceImpl) MarkMessageAsDelivered(ctx context.Context, messageID string) error {
	return s.UpdateMessageStatus(ctx, messageID, domain.MessageStatusDelivered)
}

// MarkMessageAsRead marks a message as read
func (s *ChatServiceImpl) MarkMessageAsRead(ctx context.Context, messageID string) error {
	return s.UpdateMessageStatus(ctx, messageID, domain.MessageStatusRead)
}

// UpdateSessionContext updates the context of a session
func (s *ChatServiceImpl) UpdateSessionContext(ctx context.Context, sessionID string, context *domain.SessionContext) error {
	if sessionID == "" {
		return &interfaces.ChatError{
			Operation: "UpdateSessionContext",
			Reason:    "Session ID is required",
			Code:      interfaces.ErrCodeValidationFailed,
		}
	}

	session, err := s.sessionService.GetSession(ctx, sessionID)
	if err != nil {
		return &interfaces.ChatError{
			Operation: "UpdateSessionContext",
			Reason:    "Session not found",
			Code:      interfaces.ErrCodeSessionNotFound,
			Cause:     err,
		}
	}

	// Update session fields from context
	if context.ClientName != "" {
		session.ClientName = context.ClientName
	}
	if context.ProjectContext != "" {
		session.Context = context.ProjectContext
	}

	// Store additional context in metadata
	if session.Metadata == nil {
		session.Metadata = make(map[string]interface{})
	}

	session.Metadata["meeting_type"] = context.MeetingType
	session.Metadata["service_types"] = context.ServiceTypes
	session.Metadata["cloud_providers"] = context.CloudProviders
	session.Metadata["custom_fields"] = context.CustomFields

	if err := s.sessionService.UpdateSession(ctx, session); err != nil {
		return &interfaces.ChatError{
			Operation: "UpdateSessionContext",
			Reason:    "Failed to update session context",
			Code:      interfaces.ErrCodeDatabaseError,
			Cause:     err,
		}
	}

	return nil
}

// GetSessionContext retrieves the context of a session
func (s *ChatServiceImpl) GetSessionContext(ctx context.Context, sessionID string) (*domain.SessionContext, error) {
	if sessionID == "" {
		return nil, &interfaces.ChatError{
			Operation: "GetSessionContext",
			Reason:    "Session ID is required",
			Code:      interfaces.ErrCodeValidationFailed,
		}
	}

	session, err := s.sessionService.GetSession(ctx, sessionID)
	if err != nil {
		return nil, &interfaces.ChatError{
			Operation: "GetSessionContext",
			Reason:    "Session not found",
			Code:      interfaces.ErrCodeSessionNotFound,
			Cause:     err,
		}
	}

	context := &domain.SessionContext{
		ClientName:     session.ClientName,
		ProjectContext: session.Context,
	}

	// Extract additional context from metadata
	if session.Metadata != nil {
		if meetingType, ok := session.Metadata["meeting_type"].(string); ok {
			context.MeetingType = meetingType
		}
		if serviceTypes, ok := session.Metadata["service_types"].([]string); ok {
			context.ServiceTypes = serviceTypes
		}
		if cloudProviders, ok := session.Metadata["cloud_providers"].([]string); ok {
			context.CloudProviders = cloudProviders
		}
		if customFields, ok := session.Metadata["custom_fields"].(map[string]string); ok {
			context.CustomFields = customFields
		}
	}

	return context, nil
}

// ValidateMessage validates a chat message
func (s *ChatServiceImpl) ValidateMessage(message *domain.ChatMessage) error {
	if message == nil {
		return &interfaces.MessageValidationError{
			Reason: "Message cannot be nil",
			Code:   interfaces.ErrCodeValidationFailed,
		}
	}

	if err := message.Validate(); err != nil {
		return &interfaces.MessageValidationError{
			MessageID: message.ID,
			Reason:    err.Error(),
			Code:      interfaces.ErrCodeValidationFailed,
		}
	}

	if len(message.Content) > s.maxMessageLength {
		return &interfaces.MessageValidationError{
			MessageID: message.ID,
			Field:     "content",
			Reason:    fmt.Sprintf("Message content exceeds maximum length of %d characters", s.maxMessageLength),
			Code:      interfaces.ErrCodeMessageTooLong,
		}
	}

	return nil
}

// SanitizeMessageContent sanitizes message content to prevent XSS and other attacks
func (s *ChatServiceImpl) SanitizeMessageContent(content string) string {
	// HTML escape the content
	sanitized := html.EscapeString(content)

	// Remove excessive whitespace
	sanitized = strings.TrimSpace(sanitized)

	// Replace multiple consecutive spaces with single space
	spaceRegex := regexp.MustCompile(`\s+`)
	sanitized = spaceRegex.ReplaceAllString(sanitized, " ")

	// Remove potentially dangerous patterns (basic protection)
	dangerousPatterns := []string{
		`<script[^>]*>.*?</script>`,
		`javascript:`,
		`vbscript:`,
		`onload=`,
		`onerror=`,
		`onclick=`,
	}

	for _, pattern := range dangerousPatterns {
		regex := regexp.MustCompile(`(?i)` + pattern)
		sanitized = regex.ReplaceAllString(sanitized, "")
	}

	return sanitized
}

// SearchMessages searches for messages containing the query text
func (s *ChatServiceImpl) SearchMessages(ctx context.Context, sessionID string, query string, limit int) ([]*domain.ChatMessage, error) {
	if sessionID == "" {
		return nil, &interfaces.ChatError{
			Operation: "SearchMessages",
			Reason:    "Session ID is required",
			Code:      interfaces.ErrCodeValidationFailed,
		}
	}

	if query == "" {
		return nil, &interfaces.ChatError{
			Operation: "SearchMessages",
			Reason:    "Search query is required",
			Code:      interfaces.ErrCodeValidationFailed,
		}
	}

	if limit <= 0 {
		limit = 50
	}
	if limit > s.maxHistoryLimit {
		limit = s.maxHistoryLimit
	}

	// Sanitize search query
	sanitizedQuery := s.SanitizeMessageContent(query)

	messages, err := s.messageRepository.Search(ctx, sessionID, sanitizedQuery, limit)
	if err != nil {
		s.logger.WithError(err).WithFields(logrus.Fields{
			"session_id": sessionID,
			"query":      sanitizedQuery,
		}).Error("Failed to search messages")
		return nil, &interfaces.ChatError{
			Operation: "SearchMessages",
			Reason:    "Failed to search messages",
			Code:      interfaces.ErrCodeDatabaseError,
			Cause:     err,
		}
	}

	return messages, nil
}

// GetMessagesByType retrieves messages of a specific type for a session
func (s *ChatServiceImpl) GetMessagesByType(ctx context.Context, sessionID string, messageType domain.MessageType) ([]*domain.ChatMessage, error) {
	if sessionID == "" {
		return nil, &interfaces.ChatError{
			Operation: "GetMessagesByType",
			Reason:    "Session ID is required",
			Code:      interfaces.ErrCodeValidationFailed,
		}
	}

	messages, err := s.messageRepository.GetByType(ctx, sessionID, messageType)
	if err != nil {
		s.logger.WithError(err).WithFields(logrus.Fields{
			"session_id":   sessionID,
			"message_type": messageType,
		}).Error("Failed to get messages by type")
		return nil, &interfaces.ChatError{
			Operation: "GetMessagesByType",
			Reason:    "Failed to retrieve messages by type",
			Code:      interfaces.ErrCodeDatabaseError,
			Cause:     err,
		}
	}

	return messages, nil
}

// GetMessageCount returns the count of messages matching the filters
func (s *ChatServiceImpl) GetMessageCount(ctx context.Context, filters *domain.ChatMessageFilters) (int64, error) {
	if filters == nil {
		filters = &domain.ChatMessageFilters{}
	}

	count, err := s.messageRepository.Count(ctx, filters)
	if err != nil {
		s.logger.WithError(err).Error("Failed to get message count")
		return 0, &interfaces.ChatError{
			Operation: "GetMessageCount",
			Reason:    "Failed to get message count",
			Code:      interfaces.ErrCodeDatabaseError,
			Cause:     err,
		}
	}

	return count, nil
}

// GetMessageStats returns statistics about messages in a session
func (s *ChatServiceImpl) GetMessageStats(ctx context.Context, sessionID string) (*interfaces.MessageStats, error) {
	if sessionID == "" {
		return nil, &interfaces.ChatError{
			Operation: "GetMessageStats",
			Reason:    "Session ID is required",
			Code:      interfaces.ErrCodeValidationFailed,
		}
	}

	// Get total message count
	totalMessages, err := s.messageRepository.Count(ctx, &domain.ChatMessageFilters{
		SessionID: sessionID,
	})
	if err != nil {
		return nil, &interfaces.ChatError{
			Operation: "GetMessageStats",
			Reason:    "Failed to get total message count",
			Code:      interfaces.ErrCodeDatabaseError,
			Cause:     err,
		}
	}

	// Build messages by type map
	messagesByType := make(map[domain.MessageType]int)
	for _, msgType := range []domain.MessageType{
		domain.MessageTypeUser,
		domain.MessageTypeAssistant,
		domain.MessageTypeSystem,
		domain.MessageTypeError,
	} {
		count, err := s.messageRepository.Count(ctx, &domain.ChatMessageFilters{
			SessionID: sessionID,
			Type:      msgType,
		})
		if err != nil {
			s.logger.WithError(err).WithFields(logrus.Fields{
				"session_id":   sessionID,
				"message_type": msgType,
			}).Warn("Failed to get message count by type")
			continue
		}
		messagesByType[msgType] = int(count)
	}

	// Build messages by status map
	messagesByStatus := make(map[domain.MessageStatus]int)
	for _, status := range []domain.MessageStatus{
		domain.MessageStatusSent,
		domain.MessageStatusDelivered,
		domain.MessageStatusRead,
		domain.MessageStatusFailed,
	} {
		count, err := s.messageRepository.Count(ctx, &domain.ChatMessageFilters{
			SessionID: sessionID,
			Status:    status,
		})
		if err != nil {
			s.logger.WithError(err).WithFields(logrus.Fields{
				"session_id": sessionID,
				"status":     status,
			}).Warn("Failed to get message count by status")
			continue
		}
		messagesByStatus[status] = int(count)
	}

	// Get first and last message timestamps
	messages, err := s.messageRepository.GetBySessionID(ctx, sessionID, 1, 0)
	var firstMessageAt *time.Time
	if err == nil && len(messages) > 0 {
		firstMessageAt = &messages[0].CreatedAt
	}

	latestMessages, err := s.messageRepository.GetLatestBySessionID(ctx, sessionID, 1)
	var lastMessageAt *time.Time
	if err == nil && len(latestMessages) > 0 {
		lastMessageAt = &latestMessages[0].CreatedAt
	}

	stats := &interfaces.MessageStats{
		TotalMessages:    totalMessages,
		MessagesByType:   messagesByType,
		MessagesByStatus: messagesByStatus,
		FirstMessageAt:   firstMessageAt,
		LastMessageAt:    lastMessageAt,
		// Note: AverageResponseTime would require more complex calculation
		// This could be implemented later with additional tracking
	}

	return stats, nil
}

// validateChatRequest validates a chat request
func (s *ChatServiceImpl) validateChatRequest(request *domain.ChatRequest) error {
	if request == nil {
		return fmt.Errorf("request cannot be nil")
	}

	if request.SessionID == "" {
		return fmt.Errorf("session ID is required")
	}

	if request.Content == "" {
		return fmt.Errorf("message content is required")
	}

	if len(request.Content) > s.maxMessageLength {
		return fmt.Errorf("message content exceeds maximum length of %d characters", s.maxMessageLength)
	}

	if request.Type == "" {
		return fmt.Errorf("message type is required")
	}

	return nil
}

// generateAIResponse generates an AI response for a user message
func (s *ChatServiceImpl) generateAIResponse(ctx context.Context, session *domain.ChatSession, userMessage *domain.ChatMessage, quickAction string) (*domain.ChatMessage, error) {
	var aiMessage *domain.ChatMessage
	var err error

	// Use enhanced AI service if available
	if s.enhancedAIService != nil {
		aiMessage, err = s.generateEnhancedAIResponse(ctx, session, userMessage, quickAction)
		if err != nil {
			s.logger.WithError(err).Warn("Enhanced AI service failed, falling back to basic service")
			// Fall back to basic AI service
			aiMessage, err = s.generateBasicAIResponse(ctx, session, userMessage, quickAction)
		}
	} else {
		// Use basic AI service
		aiMessage, err = s.generateBasicAIResponse(ctx, session, userMessage, quickAction)
	}

	return aiMessage, err
}

// generateEnhancedAIResponse generates an AI response using the enhanced AI service
func (s *ChatServiceImpl) generateEnhancedAIResponse(ctx context.Context, session *domain.ChatSession, userMessage *domain.ChatMessage, quickAction string) (*domain.ChatMessage, error) {
	// Get session context
	sessionContext, err := s.GetSessionContext(ctx, session.ID)
	if err != nil {
		s.logger.WithError(err).Warn("Failed to get session context, using basic context")
		sessionContext = &domain.SessionContext{
			ClientName:     session.ClientName,
			ProjectContext: session.Context,
		}
	}

	// Get recent conversation history for context
	history, err := s.messageRepository.GetLatestBySessionID(ctx, session.ID, 5)
	if err != nil {
		s.logger.WithError(err).Warn("Failed to get conversation history")
		history = []*domain.ChatMessage{}
	}

	// Create enhanced chat request
	chatRequest := &ChatRequest{
		Content:     userMessage.Content,
		Session:     session,
		Context:     sessionContext,
		QuickAction: quickAction,
		History:     history,
	}

	// Generate optimized response
	startTime := time.Now()
	chatResponse, err := s.enhancedAIService.OptimizeResponse(ctx, chatRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to generate enhanced AI response: %w", err)
	}
	processingTime := time.Since(startTime)

	// Create AI message
	aiMessage := &domain.ChatMessage{
		ID:        uuid.New().String(),
		SessionID: session.ID,
		Type:      domain.MessageTypeAssistant,
		Content:   chatResponse.Content,
		Metadata: map[string]interface{}{
			"tokens_used":       chatResponse.TokensUsed,
			"processing_time":   processingTime.Milliseconds(),
			"enhanced_ai":       true,
			"optimized":         true,
			"cache_hit":         chatResponse.Metadata["cache_hit"],
			"model_id":          chatResponse.Metadata["model_id"],
			"prompt_tokens":     chatResponse.Metadata["prompt_tokens"],
			"completion_tokens": chatResponse.Metadata["completion_tokens"],
		},
		Status:    domain.MessageStatusSent,
		CreatedAt: time.Now(),
	}

	return aiMessage, nil
}

// generateBasicAIResponse generates an AI response using the basic AI service (fallback)
func (s *ChatServiceImpl) generateBasicAIResponse(ctx context.Context, session *domain.ChatSession, userMessage *domain.ChatMessage, quickAction string) (*domain.ChatMessage, error) {
	// Build prompt based on session context and message history
	prompt, err := s.buildPrompt(ctx, session, userMessage, quickAction)
	if err != nil {
		return nil, fmt.Errorf("failed to build prompt: %w", err)
	}

	// Configure AI service options
	// Use the configured model ID from the Bedrock service
	modelInfo := s.aiService.GetModelInfo()
	options := &interfaces.BedrockOptions{
		ModelID:     modelInfo.ModelID,
		MaxTokens:   1000,
		Temperature: 0.7,
		TopP:        0.9,
	}

	// Generate AI response
	startTime := time.Now()
	response, err := s.aiService.GenerateText(ctx, prompt, options)
	if err != nil {
		return nil, fmt.Errorf("failed to generate AI response: %w", err)
	}
	processingTime := time.Since(startTime)

	// Create AI message
	aiMessage := &domain.ChatMessage{
		ID:        uuid.New().String(),
		SessionID: session.ID,
		Type:      domain.MessageTypeAssistant,
		Content:   response.Content,
		Metadata: map[string]interface{}{
			"tokens_used":     response.Usage.OutputTokens,
			"processing_time": processingTime.Milliseconds(),
			"model_id":        options.ModelID,
			"enhanced_ai":     false,
		},
		Status:    domain.MessageStatusSent,
		CreatedAt: time.Now(),
	}

	return aiMessage, nil
}

// buildPrompt builds a context-aware prompt for AI response generation
func (s *ChatServiceImpl) buildPrompt(ctx context.Context, session *domain.ChatSession, userMessage *domain.ChatMessage, quickAction string) (string, error) {
	basePrompt := `You are an expert AWS cloud consultant assistant providing real-time support during client meetings. Your role is to help consultants provide immediate, authoritative, and specific answers to client questions.

CONTEXT:
- This is a live client meeting situation
- Responses must be direct, professional, and actionable
- Focus on AWS services and best practices
- Provide specific service names, configurations, and cost estimates when possible
- Keep responses concise but comprehensive (2-3 sentences max unless technical details are needed)

CONSULTANT GUIDELINES:
- Always provide specific AWS service recommendations
- Include approximate costs when relevant
- Mention compliance considerations for regulated industries
- Suggest next steps or follow-up actions
- Reference AWS documentation or whitepapers when helpful`

	// Add session context
	if session.ClientName != "" {
		basePrompt += "\n\nCLIENT: " + session.ClientName
	}

	if session.Context != "" {
		basePrompt += "\n\nMEETING CONTEXT: " + session.Context
	}

	// Add recent conversation history (last 5 messages for context)
	history, err := s.messageRepository.GetLatestBySessionID(ctx, session.ID, 5)
	if err == nil && len(history) > 0 {
		basePrompt += "\n\nCONVERSATION HISTORY:"
		// Reverse to show chronological order
		for i := len(history) - 1; i >= 0; i-- {
			msg := history[i]
			if msg.Type == domain.MessageTypeUser {
				basePrompt += "\nConsultant: " + msg.Content
			} else if msg.Type == domain.MessageTypeAssistant {
				basePrompt += "\nAssistant: " + msg.Content
			}
		}
	}

	// Handle quick actions
	if quickAction != "" {
		basePrompt += "\n\nQUICK ACTION REQUESTED: " + quickAction
	}

	// Add current question
	basePrompt += "\n\nCURRENT QUESTION: " + userMessage.Content
	basePrompt += "\n\nProvide a direct, expert-level response that the consultant can immediately use with their client:"

	return basePrompt, nil
}
