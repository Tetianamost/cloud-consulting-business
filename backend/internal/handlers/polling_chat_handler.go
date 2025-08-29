package handlers

import (
	"context"
	"crypto/md5"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/cloud-consulting/backend/internal/domain"
	"github.com/cloud-consulting/backend/internal/interfaces"
)

// SendMessageRequest represents the request structure for sending a message via polling
type SendMessageRequest struct {
	Content    string `json:"content" validate:"required,max=10000"`
	SessionID  string `json:"session_id" validate:"required"`
	ClientName string `json:"client_name,omitempty"`
}

// SendMessageResponse represents the response structure for sending a message
type SendMessageResponse struct {
	Success   bool   `json:"success"`
	MessageID string `json:"message_id,omitempty"`
	Content   string `json:"content,omitempty"`
	Type      string `json:"type,omitempty"`
	Error     string `json:"error,omitempty"`
}

// GetMessagesResponse represents the response structure for retrieving messages
type GetMessagesResponse struct {
	Success  bool                  `json:"success"`
	Messages []*domain.ChatMessage `json:"messages"`
	HasMore  bool                  `json:"has_more"`
	Error    string                `json:"error,omitempty"`
}

// PollingChatHandler handles HTTP-based polling chat requests
type PollingChatHandler struct {
	logger         *logrus.Logger
	chatService    interfaces.ChatService
	sessionService interfaces.SessionService
	authHandler    *AuthHandler

	// Security services
	chatAuthService     interfaces.ChatAuthService
	chatSecurityService interfaces.ChatSecurityService
	chatRateLimiter     interfaces.ChatRateLimiter
	chatAuditLogger     interfaces.ChatAuditLogger
}

// NewPollingChatHandler creates a new polling chat handler
func NewPollingChatHandler(
	logger *logrus.Logger,
	chatService interfaces.ChatService,
	sessionService interfaces.SessionService,
	authHandler *AuthHandler,
	chatAuthService interfaces.ChatAuthService,
	chatSecurityService interfaces.ChatSecurityService,
	chatRateLimiter interfaces.ChatRateLimiter,
	chatAuditLogger interfaces.ChatAuditLogger,
) *PollingChatHandler {
	return &PollingChatHandler{
		logger:              logger,
		chatService:         chatService,
		sessionService:      sessionService,
		authHandler:         authHandler,
		chatAuthService:     chatAuthService,
		chatSecurityService: chatSecurityService,
		chatRateLimiter:     chatRateLimiter,
		chatAuditLogger:     chatAuditLogger,
	}
}

// AuthMiddleware authenticates HTTP requests for polling chat
func (h *PollingChatHandler) AuthMiddleware(c *gin.Context) {
	// Extract token from Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		h.logger.Error("Polling chat authentication failed: no authorization header")
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Authorization header required",
		})
		c.Abort()
		return
	}

	// Extract Bearer token
	if !strings.HasPrefix(authHeader, "Bearer ") {
		h.logger.Error("Polling chat authentication failed: invalid authorization header format")
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Invalid authorization header format",
		})
		c.Abort()
		return
	}

	token := authHeader[7:] // Remove "Bearer " prefix

	// Validate token using enhanced auth service
	authContext, err := h.chatAuthService.ValidateToken(context.Background(), token)
	if err != nil {
		h.logger.WithError(err).WithField("token_preview", token[:min(len(token), 20)]+"...").Error("Polling chat authentication failed: invalid token")

		// Log failed authentication attempt
		metadata := map[string]interface{}{
			"ip_address": c.ClientIP(),
			"user_agent": c.GetHeader("User-Agent"),
			"error":      err.Error(),
		}
		h.chatAuditLogger.LogLogin(context.Background(), "unknown", false, metadata)

		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Invalid or expired authentication token",
		})
		c.Abort()
		return
	}

	// Check rate limit for API requests
	rateLimitResult, err := h.chatSecurityService.CheckRateLimit(context.Background(), authContext.UserID, "api_request")
	if err != nil {
		h.logger.WithError(err).Warn("Failed to check rate limit for polling chat request")
	} else if !rateLimitResult.Allowed {
		h.logger.WithField("user_id", authContext.UserID).Warn("Polling chat API rate limit exceeded")

		// Log rate limit exceeded
		metadata := map[string]interface{}{
			"ip_address": c.ClientIP(),
			"user_agent": c.GetHeader("User-Agent"),
			"action":     "api_request",
		}
		h.chatAuditLogger.LogRateLimitExceeded(context.Background(), authContext.UserID, "api_request", metadata)

		c.JSON(http.StatusTooManyRequests, gin.H{
			"success":     false,
			"error":       "API rate limit exceeded",
			"retry_after": int(rateLimitResult.RetryAfter.Seconds()),
		})
		c.Abort()
		return
	}

	// Debug logging
	h.logger.WithFields(logrus.Fields{
		"user_id":  authContext.UserID,
		"username": authContext.Username,
		"roles":    authContext.Roles,
	}).Info("Authentication context set successfully")

	// Set authentication context
	c.Set("auth_context", authContext)
	c.Set("user_id", authContext.UserID)
	c.Set("username", authContext.Username)

	// Debug logging
	h.logger.WithFields(logrus.Fields{
		"user_id":  authContext.UserID,
		"username": authContext.Username,
		"roles":    authContext.Roles,
	}).Info("Polling chat authentication successful")

	// Log successful authentication
	metadata := map[string]interface{}{
		"ip_address":      c.ClientIP(),
		"user_agent":      c.GetHeader("User-Agent"),
		"connection_type": "polling",
	}
	h.chatAuditLogger.LogLogin(context.Background(), authContext.UserID, true, metadata)

	c.Next()
}

// SendMessage handles POST /api/v1/admin/chat/messages - sends a new message
func (h *PollingChatHandler) SendMessage(c *gin.Context) {
	ctx := context.Background()

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		h.logger.Error("User ID not found in context")
		c.JSON(http.StatusInternalServerError, SendMessageResponse{
			Success: false,
			Error:   "Authentication context error",
		})
		return
	}

	userIDStr := userID.(string)

	// Parse request body
	var request SendMessageRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		h.logger.WithError(err).Error("Failed to parse send message request")
		c.JSON(http.StatusBadRequest, SendMessageResponse{
			Success: false,
			Error:   "Invalid request format",
		})
		return
	}

	// Validate message content
	if err := h.validateMessageContent(request.Content); err != nil {
		h.logger.WithError(err).Error("Message content validation failed")
		c.JSON(http.StatusBadRequest, SendMessageResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Validate session ID format
	if err := h.validateSessionID(request.SessionID); err != nil {
		h.logger.WithError(err).Error("Session ID validation failed")
		c.JSON(http.StatusBadRequest, SendMessageResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Validate or create session
	session, err := h.getOrCreateSession(ctx, request.SessionID, userIDStr, request.ClientName)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get or create session")
		c.JSON(http.StatusInternalServerError, SendMessageResponse{
			Success: false,
			Error:   "Failed to manage session",
		})
		return
	}

	// Create chat request
	chatRequest := &domain.ChatRequest{
		SessionID: session.ID,
		Content:   request.Content,
		Type:      domain.MessageTypeUser,
		Metadata: map[string]interface{}{
			"user_id":     userIDStr,
			"client_name": request.ClientName,
			"ip_address":  c.ClientIP(),
			"user_agent":  c.GetHeader("User-Agent"),
			"method":      "polling",
		},
	}

	// Send message through chat service
	response, err := h.chatService.SendMessage(ctx, chatRequest)
	if err != nil {
		h.logger.WithError(err).Error("Failed to send message through chat service")
		c.JSON(http.StatusInternalServerError, SendMessageResponse{
			Success: false,
			Error:   "Failed to process message",
		})
		return
	}

	// Log successful message send
	h.chatAuditLogger.LogMessage(ctx, userIDStr, session.ID, "sent", map[string]interface{}{
		"message_id": response.MessageID,
		"method":     "polling",
	})

	// Return success response with AI response content
	c.JSON(http.StatusOK, SendMessageResponse{
		Success:   true,
		MessageID: response.MessageID,
		Content:   response.Content,
		Type:      string(response.Type),
	})
}

// validateMessageContent validates the message content
func (h *PollingChatHandler) validateMessageContent(content string) error {
	if content == "" {
		return fmt.Errorf("message content cannot be empty")
	}
	if len(content) > 10000 {
		return fmt.Errorf("message content cannot exceed 10000 characters")
	}
	// Additional content validation can be added here
	return nil
}

// validateSessionID validates the session ID format
func (h *PollingChatHandler) validateSessionID(sessionID string) error {
	if sessionID == "" {
		return fmt.Errorf("session ID cannot be empty")
	}
	if len(sessionID) < 10 || len(sessionID) > 100 {
		return fmt.Errorf("session ID must be between 10 and 100 characters")
	}
	// Additional session ID format validation can be added here
	return nil
}

// getOrCreateSession gets an existing session or creates a new one
func (h *PollingChatHandler) getOrCreateSession(ctx context.Context, sessionID, userID, clientName string) (*domain.ChatSession, error) {
	// Try to get existing session
	session, err := h.sessionService.ValidateSession(ctx, sessionID, userID)
	if err == nil {
		// Update session context if client name is provided
		if clientName != "" && session.ClientName != clientName {
			sessionContext := &domain.SessionContext{
				ClientName: clientName,
			}
			if err := h.chatService.UpdateSessionContext(ctx, sessionID, sessionContext); err != nil {
				h.logger.WithError(err).Warn("Failed to update session context")
			}
		}
		return session, nil
	}

	h.logger.WithError(err).Info("Session not found or invalid, creating new session")

	// Create new session
	newSession := &domain.ChatSession{
		ID:           sessionID, // Use provided session ID
		UserID:       userID,
		ClientName:   clientName,
		Status:       domain.SessionStatusActive,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		LastActivity: time.Now(),
		ExpiresAt:    timePtr(time.Now().Add(24 * time.Hour)), // 24 hour expiration
		Metadata: map[string]interface{}{
			"created_via": "polling",
		},
	}

	if err := h.sessionService.CreateSession(ctx, newSession); err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	h.logger.WithFields(logrus.Fields{
		"session_id": newSession.ID,
		"user_id":    userID,
	}).Info("Created new chat session via polling")

	return newSession, nil
}

// GetMessages handles GET /api/v1/admin/chat/messages - retrieves messages for a session
// Supports conditional requests using ETags and timestamps for efficient caching
func (h *PollingChatHandler) GetMessages(c *gin.Context) {
	ctx := context.Background()

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		h.logger.Error("User ID not found in context")
		c.JSON(http.StatusInternalServerError, GetMessagesResponse{
			Success: false,
			Error:   "Authentication context error",
		})
		return
	}

	userIDStr := userID.(string)

	// Get query parameters
	sessionID := c.Query("session_id")
	since := c.Query("since") // Can be message ID or timestamp
	limitStr := c.DefaultQuery("limit", "50")
	offsetStr := c.DefaultQuery("offset", "0")

	// Get conditional request headers for caching
	ifNoneMatch := c.GetHeader("If-None-Match")
	ifModifiedSince := c.GetHeader("If-Modified-Since")

	// Validate required parameters
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, GetMessagesResponse{
			Success: false,
			Error:   "session_id parameter is required",
		})
		return
	}

	// Validate session ID format
	if err := h.validateSessionID(sessionID); err != nil {
		h.logger.WithError(err).Error("Session ID validation failed")
		c.JSON(http.StatusBadRequest, GetMessagesResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Parse limit and offset
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 1000 {
		c.JSON(http.StatusBadRequest, GetMessagesResponse{
			Success: false,
			Error:   "limit must be between 1 and 1000",
		})
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		c.JSON(http.StatusBadRequest, GetMessagesResponse{
			Success: false,
			Error:   "offset must be 0 or greater",
		})
		return
	}

	// Validate session access
	session, err := h.sessionService.ValidateSession(ctx, sessionID, userIDStr)
	if err != nil {
		h.logger.WithError(err).Error("Session validation failed")
		c.JSON(http.StatusForbidden, GetMessagesResponse{
			Success: false,
			Error:   "Invalid session or access denied",
		})
		return
	}

	// Build message filters
	filters := &domain.ChatMessageFilters{
		SessionID: sessionID,
		Limit:     limit,
		Offset:    offset,
	}

	// Handle 'since' parameter for efficient polling
	if since != "" {
		if err := h.applySinceFilter(filters, since); err != nil {
			h.logger.WithError(err).Error("Failed to apply since filter")
			c.JSON(http.StatusBadRequest, GetMessagesResponse{
				Success: false,
				Error:   "Invalid since parameter format",
			})
			return
		}
	}

	// Check conditional requests for caching optimization
	if ifNoneMatch != "" || ifModifiedSince != "" {
		// Generate current ETag based on session's last activity
		currentETag := h.generateETag(session)

		// Check If-None-Match header
		if ifNoneMatch != "" && ifNoneMatch == currentETag {
			h.logger.WithFields(logrus.Fields{
				"session_id": sessionID,
				"etag":       currentETag,
			}).Debug("ETag match, returning 304 Not Modified")

			c.Header("ETag", currentETag)
			c.Header("Cache-Control", "private, max-age=0, must-revalidate")
			c.Status(http.StatusNotModified)
			return
		}

		// Check If-Modified-Since header
		if ifModifiedSince != "" {
			if modifiedSince, err := time.Parse(http.TimeFormat, ifModifiedSince); err == nil {
				if session.LastActivity.Before(modifiedSince) || session.LastActivity.Equal(modifiedSince) {
					h.logger.WithFields(logrus.Fields{
						"session_id":        sessionID,
						"last_activity":     session.LastActivity,
						"if_modified_since": modifiedSince,
					}).Debug("Not modified since, returning 304 Not Modified")

					c.Header("ETag", currentETag)
					c.Header("Last-Modified", session.LastActivity.Format(http.TimeFormat))
					c.Header("Cache-Control", "private, max-age=0, must-revalidate")
					c.Status(http.StatusNotModified)
					return
				}
			}
		}
	}

	// Retrieve messages with optimized database queries
	messages, err := h.chatService.ListMessages(ctx, filters)
	if err != nil {
		h.logger.WithError(err).Error("Failed to retrieve messages")
		c.JSON(http.StatusInternalServerError, GetMessagesResponse{
			Success: false,
			Error:   "Failed to retrieve messages",
		})
		return
	}

	// Check if there are more messages available
	hasMore := len(messages) == limit

	// Generate ETag and set caching headers
	etag := h.generateETag(session)
	c.Header("ETag", etag)
	c.Header("Last-Modified", session.LastActivity.Format(http.TimeFormat))
	c.Header("Cache-Control", "private, max-age=0, must-revalidate")

	// Add compression support
	c.Header("Content-Encoding", "gzip")

	// Log successful message retrieval
	h.chatAuditLogger.LogMessage(ctx, userIDStr, sessionID, "retrieved", map[string]interface{}{
		"message_count": len(messages),
		"method":        "polling",
		"since":         since,
		"cached":        false,
	})

	// Return messages with caching headers
	c.JSON(http.StatusOK, GetMessagesResponse{
		Success:  true,
		Messages: messages,
		HasMore:  hasMore,
	})
}

// applySinceFilter applies the 'since' parameter to message filters
func (h *PollingChatHandler) applySinceFilter(filters *domain.ChatMessageFilters, since string) error {
	// Try to parse as timestamp first (RFC3339 format)
	if timestamp, err := time.Parse(time.RFC3339, since); err == nil {
		filters.FromDate = &timestamp
		return nil
	}

	// Try to parse as Unix timestamp
	if unixTime, err := strconv.ParseInt(since, 10, 64); err == nil {
		timestamp := time.Unix(unixTime, 0)
		filters.FromDate = &timestamp
		return nil
	}

	// If not a timestamp, treat as message ID (would require custom query logic)
	// For now, we'll return an error for unsupported format
	return fmt.Errorf("since parameter must be a valid timestamp (RFC3339 or Unix)")
}

// generateETag generates an ETag for conditional requests based on session state
func (h *PollingChatHandler) generateETag(session *domain.ChatSession) string {
	// Create ETag based on session ID and last activity timestamp
	data := fmt.Sprintf("%s-%d", session.ID, session.LastActivity.Unix())
	hash := md5.Sum([]byte(data))
	return fmt.Sprintf(`"%x"`, hash)
}

// timePtr function is defined in chat_handler.go in the same package

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
