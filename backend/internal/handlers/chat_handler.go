package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"

	"github.com/cloud-consulting/backend/internal/domain"
	"github.com/cloud-consulting/backend/internal/interfaces"
	"github.com/cloud-consulting/backend/internal/services"
)

// WebSocketMessageType defines the types of WebSocket messages
type WebSocketMessageType string

const (
	WSMessageTypeMessage   WebSocketMessageType = "message"
	WSMessageTypeTyping    WebSocketMessageType = "typing"
	WSMessageTypeStatus    WebSocketMessageType = "status"
	WSMessageTypeError     WebSocketMessageType = "error"
	WSMessageTypePresence  WebSocketMessageType = "presence"
	WSMessageTypeAck       WebSocketMessageType = "ack"
	WSMessageTypeHeartbeat WebSocketMessageType = "heartbeat"
)

// WebSocketMessage represents a WebSocket message following the protocol
type WebSocketMessage struct {
	Type      WebSocketMessageType   `json:"type"`
	SessionID string                 `json:"session_id"`
	MessageID string                 `json:"message_id,omitempty"`
	Content   string                 `json:"content,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
}

// TypingIndicator represents a typing indicator message
type TypingIndicator struct {
	UserID    string `json:"user_id"`
	SessionID string `json:"session_id"`
	IsTyping  bool   `json:"is_typing"`
}

// PresenceUpdate represents a presence update message
type PresenceUpdate struct {
	UserID    string `json:"user_id"`
	SessionID string `json:"session_id"`
	Status    string `json:"status"` // "online", "offline", "away"
}

// MessageAcknowledgment represents a message acknowledgment
type MessageAcknowledgment struct {
	MessageID string `json:"message_id"`
	Status    string `json:"status"` // "delivered", "read", "failed"
}

// ChatMessage represents a message in the consultant chat
type ChatMessage struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"` // "user", "assistant", "system"
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
	SessionID string    `json:"session_id"`
}

// ChatSession represents an active consultant chat session
type ChatSession struct {
	ID           string        `json:"id"`
	ConsultantID string        `json:"consultant_id"`
	ClientName   string        `json:"client_name,omitempty"`
	Context      string        `json:"context,omitempty"`
	Messages     []ChatMessage `json:"messages"`
	CreatedAt    time.Time     `json:"created_at"`
	LastActivity time.Time     `json:"last_activity"`
}

// ChatRequest represents an incoming chat request
type ChatRequest struct {
	Message     string `json:"message"`
	SessionID   string `json:"session_id,omitempty"`
	ClientName  string `json:"client_name,omitempty"`
	Context     string `json:"context,omitempty"`
	QuickAction string `json:"quick_action,omitempty"`
}

// ChatResponse represents a chat response
type ChatResponse struct {
	Message   ChatMessage `json:"message"`
	SessionID string      `json:"session_id"`
	Success   bool        `json:"success"`
	Error     string      `json:"error,omitempty"`
}

// WebSocket upgrader with proper origin checking
// Note: This will be configured per handler instance to use proper CORS settings

// PendingMessage represents a message waiting for acknowledgment
type PendingMessage struct {
	Message    WebSocketMessage
	Timestamp  time.Time
	Retries    int
	MaxRetries int
}

// Connection represents a WebSocket connection with metadata
type Connection struct {
	Conn            *websocket.Conn
	UserID          string
	SessionID       string
	LastPing        time.Time
	LastActivity    time.Time
	IsTyping        bool
	Metadata        map[string]interface{}
	SendChan        chan WebSocketMessage
	CloseChan       chan bool
	PendingMessages map[string]*PendingMessage
	PendingMutex    sync.RWMutex
}

// ConnectionPool manages WebSocket connections
type ConnectionPool struct {
	connections map[string]*Connection
	mutex       sync.RWMutex
}

// NewConnectionPool creates a new connection pool
func NewConnectionPool() *ConnectionPool {
	return &ConnectionPool{
		connections: make(map[string]*Connection),
	}
}

// Add adds a connection to the pool
func (cp *ConnectionPool) Add(connID string, conn *Connection) {
	cp.mutex.Lock()
	defer cp.mutex.Unlock()
	cp.connections[connID] = conn
}

// Remove removes a connection from the pool
func (cp *ConnectionPool) Remove(connID string) {
	cp.mutex.Lock()
	defer cp.mutex.Unlock()
	delete(cp.connections, connID)
}

// Get gets a connection from the pool
func (cp *ConnectionPool) Get(connID string) (*Connection, bool) {
	cp.mutex.RLock()
	defer cp.mutex.RUnlock()
	conn, exists := cp.connections[connID]
	return conn, exists
}

// GetByUserID gets all connections for a user
func (cp *ConnectionPool) GetByUserID(userID string) []*Connection {
	cp.mutex.RLock()
	defer cp.mutex.RUnlock()

	var userConnections []*Connection
	for _, conn := range cp.connections {
		if conn.UserID == userID {
			userConnections = append(userConnections, conn)
		}
	}
	return userConnections
}

// GetBySessionID gets all connections for a session
func (cp *ConnectionPool) GetBySessionID(sessionID string) []*Connection {
	cp.mutex.RLock()
	defer cp.mutex.RUnlock()

	var sessionConnections []*Connection
	for _, conn := range cp.connections {
		if conn.SessionID == sessionID {
			sessionConnections = append(sessionConnections, conn)
		}
	}
	return sessionConnections
}

// Count returns the number of active connections
func (cp *ConnectionPool) Count() int {
	cp.mutex.RLock()
	defer cp.mutex.RUnlock()
	return len(cp.connections)
}

// RateLimiter provides simple rate limiting functionality
type RateLimiter struct {
	requests map[string][]time.Time
	mutex    sync.RWMutex
	limit    int
	window   time.Duration
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
}

// Allow checks if a request is allowed for the given key
func (rl *RateLimiter) Allow(key string) bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()

	// Clean up old requests
	if requests, exists := rl.requests[key]; exists {
		var validRequests []time.Time
		for _, reqTime := range requests {
			if now.Sub(reqTime) < rl.window {
				validRequests = append(validRequests, reqTime)
			}
		}
		rl.requests[key] = validRequests
	}

	// Check if limit is exceeded
	if len(rl.requests[key]) >= rl.limit {
		return false
	}

	// Add current request
	rl.requests[key] = append(rl.requests[key], now)
	return true
}

// ChatHandler handles real-time consultant chat with enhanced session management
type ChatHandler struct {
	logger           *logrus.Logger
	bedrockService   interfaces.BedrockService
	enhancedBedrock  *services.EnhancedBedrockService
	knowledgeBase    interfaces.KnowledgeBase
	sessionService   interfaces.SessionService
	chatService      interfaces.ChatService
	authHandler      *AuthHandler
	connectionPool   *ConnectionPool
	rateLimiter      *RateLimiter
	jwtSecret        string
	corsOrigins      []string                   // CORS allowed origins for WebSocket
	upgrader         websocket.Upgrader         // WebSocket upgrader with custom CheckOrigin
	sessions         map[string]*ChatSession    // Legacy - will be removed
	connections      map[string]*websocket.Conn // Legacy - will be removed
	sessionsMutex    sync.RWMutex               // Legacy - will be removed
	connectionsMutex sync.RWMutex               // Legacy - will be removed

	// Enhanced security services
	chatAuthService     interfaces.ChatAuthService
	chatSecurityService interfaces.ChatSecurityService
	chatRateLimiter     interfaces.ChatRateLimiter
	chatAuditLogger     interfaces.ChatAuditLogger
	authMiddleware      interfaces.ChatAuthMiddleware

	// Metrics and monitoring
	metricsCollector   *services.ChatMetricsCollector
	performanceMonitor *services.ChatPerformanceMonitor
}

// NewChatHandler creates a new enhanced chat handler with session management
func NewChatHandler(
	logger *logrus.Logger,
	bedrockService interfaces.BedrockService,
	knowledgeBase interfaces.KnowledgeBase,
	sessionService interfaces.SessionService,
	chatService interfaces.ChatService,
	authHandler *AuthHandler,
	jwtSecret string,
	corsOrigins []string,
	metricsCollector *services.ChatMetricsCollector,
	performanceMonitor *services.ChatPerformanceMonitor,
) *ChatHandler {
	clientHistoryService := services.NewClientHistoryService(knowledgeBase)
	companyKnowledgeIntegrationService := services.NewCompanyKnowledgeIntegrationService(knowledgeBase, clientHistoryService)
	enhancedBedrock := services.NewEnhancedBedrockService(bedrockService, knowledgeBase, companyKnowledgeIntegrationService)

	// Initialize security services
	chatAuthService := services.NewChatAuthService(jwtSecret, logger)
	chatRateLimiter := services.NewChatRateLimiter(logger)
	chatAuditLogger := services.NewChatAuditLogger(logger)
	chatSecurityService := services.NewChatSecurityService(logger, jwtSecret, chatRateLimiter, chatAuditLogger)
	authMiddleware := services.NewChatAuthMiddleware(chatAuthService, chatAuditLogger, logger)

	// Custom CheckOrigin for WebSocket upgrader
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			origin := r.Header.Get("Origin")
			for _, allowed := range corsOrigins {
				if allowed == "*" || origin == allowed {
					return true
				}
			}
			return false
		},
	}

	return &ChatHandler{
		logger:          logger,
		bedrockService:  bedrockService,
		enhancedBedrock: enhancedBedrock,
		knowledgeBase:   knowledgeBase,
		sessionService:  sessionService,
		chatService:     chatService,
		authHandler:     authHandler,
		connectionPool:  NewConnectionPool(),
		rateLimiter:     NewRateLimiter(60, time.Minute), // 60 messages per minute
		jwtSecret:       jwtSecret,
		corsOrigins:     corsOrigins,
		upgrader:        upgrader,
		sessions:        make(map[string]*ChatSession),    // Legacy
		connections:     make(map[string]*websocket.Conn), // Legacy

		// Enhanced security services
		chatAuthService:     chatAuthService,
		chatSecurityService: chatSecurityService,
		chatRateLimiter:     chatRateLimiter,
		chatAuditLogger:     chatAuditLogger,
		authMiddleware:      authMiddleware,

		// Metrics and monitoring
		metricsCollector:   metricsCollector,
		performanceMonitor: performanceMonitor,
	}
}

// WebSocketAuthMiddleware authenticates WebSocket connections using enhanced security
func (h *ChatHandler) WebSocketAuthMiddleware(c *gin.Context) {
	fmt.Printf("[CHAT DEBUG] WebSocketAuthMiddleware called\n")
	// Extract token from query parameter or header
	token := c.Query("token")
	fmt.Printf("[CHAT DEBUG] Token from query: %s\n", token)
	if token == "" {
		authHeader := c.GetHeader("Authorization")
		fmt.Printf("[CHAT DEBUG] Authorization header: %s\n", authHeader)
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			token = authHeader[7:]
			fmt.Printf("[CHAT DEBUG] Token from header: %s\n", token)
		}
	}

	if token == "" {
		h.logger.Error("WebSocket authentication failed: no token provided")
		fmt.Printf("[CHAT DEBUG] No token provided, aborting\n")

		// Log authentication attempt
		metadata := map[string]interface{}{
			"ip_address": c.ClientIP(),
			"user_agent": c.GetHeader("User-Agent"),
			"error":      "no_token_provided",
		}
		h.chatAuditLogger.LogLogin(context.Background(), "unknown", false, metadata)

		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Authentication token required for WebSocket connection",
		})
		c.Abort()
		return
	}

	// Validate token using enhanced auth service
	authContext, err := h.chatAuthService.ValidateToken(context.Background(), token)
	if err != nil {
		h.logger.WithError(err).Error("WebSocket authentication failed: invalid token")
		fmt.Printf("[CHAT DEBUG] Token validation error: %v\n", err)

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

	// Check rate limit for WebSocket connections
	// --- for development only: allow unlimited connections from localhost ---
	clientIP := c.ClientIP()
	if clientIP == "127.0.0.1" || clientIP == "::1" {
		// Skip rate limiting for localhost during development
	} else {
		rateLimitResult, err := h.chatSecurityService.CheckRateLimit(context.Background(), authContext.UserID, "connection")
		if err != nil {
			h.logger.WithError(err).Warn("Failed to check rate limit for WebSocket connection")
		} else if !rateLimitResult.Allowed {
			h.logger.WithField("user_id", authContext.UserID).Warn("WebSocket connection rate limit exceeded")

			// Log rate limit exceeded
			metadata := map[string]interface{}{
				"ip_address": c.ClientIP(),
				"user_agent": c.GetHeader("User-Agent"),
				"action":     "connection",
			}
			h.chatAuditLogger.LogRateLimitExceeded(context.Background(), authContext.UserID, "connection", metadata)

			c.JSON(http.StatusTooManyRequests, gin.H{
				"success":     false,
				"error":       "Connection rate limit exceeded",
				"retry_after": int(rateLimitResult.RetryAfter.Seconds()),
			})
			c.Abort()
			return
		}
	}

	// Set enhanced authentication context
	c.Set("auth_context", authContext)
	c.Set("user_id", authContext.UserID)
	c.Set("username", authContext.Username)
	c.Set("user_roles", authContext.Roles)
	c.Set("user_permissions", authContext.Permissions)

	// Log successful authentication
	metadata := map[string]interface{}{
		"ip_address":      c.ClientIP(),
		"user_agent":      c.GetHeader("User-Agent"),
		"connection_type": "websocket",
	}
	h.chatAuditLogger.LogLogin(context.Background(), authContext.UserID, true, metadata)

	c.Next()
}

// HandleWebSocket handles WebSocket connections for real-time chat with enhanced session management
func (h *ChatHandler) HandleWebSocket(c *gin.Context) {
	fmt.Println("[CHAT DEBUG] HandleWebSocket called")
	// Authenticate the WebSocket connection
	h.WebSocketAuthMiddleware(c)
	if c.IsAborted() {
		return
	}

	// Get user information from context
	userID, exists := c.Get("user_id")
	if !exists {
		h.logger.Error("User ID not found in context after authentication")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Authentication context error",
		})
		return
	}

	userIDStr := userID.(string)

	// Upgrade to WebSocket
	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		h.logger.WithError(err).Error("Failed to upgrade WebSocket connection")
		// Record connection failure
		if h.metricsCollector != nil {
			h.metricsCollector.RecordConnection("failed")
			h.metricsCollector.RecordError("system")
		}
		return
	}
	// Generate connection ID
	connID := generateID()

	defer func() {
		h.logger.WithFields(logrus.Fields{
			"connection_id": connID,
			"user_id":       userIDStr,
		}).Info("WebSocket connection closed (deferred cleanup)")
		conn.Close()
	}()

	// Create connection object
	connection := &Connection{
		Conn:            conn,
		UserID:          userIDStr,
		LastPing:        time.Now(),
		LastActivity:    time.Now(),
		IsTyping:        false,
		Metadata:        make(map[string]interface{}),
		SendChan:        make(chan WebSocketMessage, 256),
		CloseChan:       make(chan bool, 1),
		PendingMessages: make(map[string]*PendingMessage),
	}

	// Add to connection pool
	h.connectionPool.Add(connID, connection)
	defer func() {
		h.connectionPool.Remove(connID)
		// Record connection closed
		if h.metricsCollector != nil {
			connectionDuration := time.Since(connection.LastActivity)
			h.metricsCollector.RecordConnection("closed", connectionDuration)
		}
		if h.performanceMonitor != nil {
			h.performanceMonitor.RecordConnectionClosed()
		}
	}()

	// Record successful connection
	if h.metricsCollector != nil {
		h.metricsCollector.RecordConnection("opened")
	}
	if h.performanceMonitor != nil {
		h.performanceMonitor.RecordConnectionOpened()
	}

	h.logger.WithFields(logrus.Fields{
		"connection_id": connID,
		"user_id":       userIDStr,
	}).Info("New authenticated WebSocket connection established")

	// Set up ping/pong handlers for connection health
	conn.SetPongHandler(func(string) error {
		connection.LastPing = time.Now()
		return nil
	})

	// Start ping routine and connection monitor
	go h.pingRoutine(conn, connID)
	go h.connectionMonitor(connection, connID)

	// Start message sender goroutine
	go h.messageSender(connection, connID)

	// Handle incoming messages
	fmt.Printf("[CHAT DEBUG] WebSocket message loop started for connection %s\n", connID)
	for {
		var wsMessage WebSocketMessage
		err := conn.ReadJSON(&wsMessage)
		if err != nil {
			fmt.Printf("[CHAT DEBUG] ReadJSON error for connection %s: %v\n", connID, err)
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				h.logger.WithError(err).WithField("connection_id", connID).Error("WebSocket unexpected close error")
				// Record WebSocket error
				if h.metricsCollector != nil {
					h.metricsCollector.RecordConnection("error")
					h.metricsCollector.RecordError("system")
				}
			} else {
				h.logger.WithError(err).WithField("connection_id", connID).Info("WebSocket connection closed by client or normal closure")
			}
			break
		}
		fmt.Printf("[CHAT DEBUG] Received WebSocket message for connection %s: %+v\n", connID, wsMessage)

		// Record message received
		if h.metricsCollector != nil {
			h.metricsCollector.RecordMessage("received")
		}
		if h.performanceMonitor != nil {
			h.performanceMonitor.RecordMessageReceived()
		}

		// Update connection activity
		connection.LastActivity = time.Now()

		// Rate limiting check
		if !h.rateLimiter.Allow(userIDStr) {
			h.logger.WithField("user_id", userIDStr).Warn("Rate limit exceeded for user")
			// Record rate limit error
			if h.metricsCollector != nil {
				h.metricsCollector.RecordError("validation")
			}
			errorMessage := WebSocketMessage{
				Type:      WSMessageTypeError,
				Content:   "Rate limit exceeded. Please slow down your requests.",
				Timestamp: time.Now(),
			}
			select {
			case connection.SendChan <- errorMessage:
			default:
				h.logger.Warn("Failed to send rate limit error - channel full")
			}
			continue
		}

		// Route message based on type
		startTime := time.Now()
		h.routeWebSocketMessage(wsMessage, connection, connID, userIDStr)

		// Record message processing time
		processingTime := time.Since(startTime)
		if h.metricsCollector != nil {
			h.metricsCollector.RecordMessage("sent", processingTime)
		}
		if h.performanceMonitor != nil {
			h.performanceMonitor.RecordResponseTime(processingTime)
		}
	}

	// Signal connection close
	fmt.Printf("[CHAT DEBUG] WebSocket message loop exiting for connection %s\n", connID)
	h.logger.WithFields(logrus.Fields{
		"connection_id": connID,
		"user_id":       userIDStr,
	}).Info("WebSocket connection handler exiting, closing connection")
	close(connection.CloseChan)
}

// pingRoutine sends periodic ping messages to keep the connection alive
func (h *ChatHandler) pingRoutine(conn *websocket.Conn, connID string) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				h.logger.WithError(err).WithField("connection_id", connID).Error("Failed to send ping, closing connection")
				return
			}
		}
	}
}

// connectionMonitor monitors connection health and handles cleanup
func (h *ChatHandler) connectionMonitor(connection *Connection, connID string) {
	ticker := time.NewTicker(60 * time.Second) // Check every minute
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Check if connection is stale (no activity for 5 minutes)
			if time.Since(connection.LastActivity) > 5*time.Minute {
				h.logger.WithFields(logrus.Fields{
					"connection_id": connID,
					"user_id":       connection.UserID,
					"last_activity": connection.LastActivity,
				}).Info("Connection is stale, closing")

				// Close the connection
				h.logger.WithFields(logrus.Fields{
					"connection_id": connID,
					"user_id":       connection.UserID,
				}).Info("Closing stale WebSocket connection due to inactivity")
				connection.Conn.Close()
				return
			}

			// Check if ping is stale (no pong for 2 minutes)
			if time.Since(connection.LastPing) > 2*time.Minute {
				h.logger.WithFields(logrus.Fields{
					"connection_id": connID,
					"user_id":       connection.UserID,
					"last_ping":     connection.LastPing,
				}).Warn("Connection ping is stale")
			}

		case <-connection.CloseChan:
			return
		}
	}
}

// processEnhancedChatRequest processes an incoming chat request with enhanced session management
func (h *ChatHandler) processEnhancedChatRequest(request ChatRequest, connID, userID string) ChatResponse {
	ctx := context.Background()

	// Get or create session using the session service
	session, err := h.getOrCreateEnhancedSession(ctx, request.SessionID, userID, request)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get or create session")
		return ChatResponse{
			Success: false,
			Error:   "Failed to manage session",
		}
	}

	// Update connection with session ID
	if conn, exists := h.connectionPool.Get(connID); exists {
		conn.SessionID = session.ID
	}

	// Send message through chat service
	chatRequest := &domain.ChatRequest{
		SessionID:   session.ID,
		Content:     request.Message,
		Type:        domain.MessageTypeUser,
		QuickAction: request.QuickAction,
		Metadata: map[string]interface{}{
			"connection_id": connID,
			"user_id":       userID,
			"client_name":   request.ClientName,
			"context":       request.Context,
		},
	}

	response, err := h.chatService.SendMessage(ctx, chatRequest)
	if err != nil {
		h.logger.WithError(err).Error("Failed to send message through chat service")
		return ChatResponse{
			SessionID: session.ID,
			Success:   false,
			Error:     "Failed to process message",
		}
	}

	// Convert domain response to handler response
	return ChatResponse{
		Message: ChatMessage{
			ID:        response.MessageID,
			Type:      string(response.Type),
			Content:   response.Content,
			Timestamp: response.CreatedAt,
			SessionID: response.SessionID,
		},
		SessionID: session.ID,
		Success:   true,
		Error:     "",
	}
}

// getOrCreateEnhancedSession gets an existing session or creates a new one using session service
func (h *ChatHandler) getOrCreateEnhancedSession(ctx context.Context, sessionID, userID string, request ChatRequest) (*domain.ChatSession, error) {
	// If session ID is provided, try to get existing session
	if sessionID != "" {
		session, err := h.sessionService.ValidateSession(ctx, sessionID, userID)
		if err == nil {
			// Update session context if provided
			if request.ClientName != "" || request.Context != "" {
				sessionContext := &domain.SessionContext{
					ClientName:     request.ClientName,
					ProjectContext: request.Context,
				}
				if err := h.chatService.UpdateSessionContext(ctx, sessionID, sessionContext); err != nil {
					h.logger.WithError(err).Warn("Failed to update session context")
				}
			}
			return session, nil
		}
		h.logger.WithError(err).Warn("Failed to validate existing session, creating new one")
	}

	// Create new session
	newSession := &domain.ChatSession{
		ID:           generateID(),
		UserID:       userID,
		ClientName:   request.ClientName,
		Context:      request.Context,
		Status:       domain.SessionStatusActive,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		LastActivity: time.Now(),
		ExpiresAt:    timePtr(time.Now().Add(24 * time.Hour)), // 24 hour expiration
		Metadata: map[string]interface{}{
			"created_via": "websocket",
		},
	}

	if err := h.sessionService.CreateSession(ctx, newSession); err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	h.logger.WithFields(logrus.Fields{
		"session_id": newSession.ID,
		"user_id":    userID,
	}).Info("Created new chat session")

	return newSession, nil
}

// processChatRequest processes an incoming chat request (legacy method for backward compatibility)
func (h *ChatHandler) processChatRequest(request ChatRequest, connID string) ChatResponse {
	// Get or create session
	session := h.getOrCreateSession(request.SessionID, connID)

	// Add user message to session
	userMessage := ChatMessage{
		ID:        generateID(),
		Type:      "user",
		Content:   request.Message,
		Timestamp: time.Now(),
		SessionID: session.ID,
	}

	h.addMessageToSession(session.ID, userMessage)

	// Update session context if provided
	if request.ClientName != "" {
		session.ClientName = request.ClientName
	}
	if request.Context != "" {
		session.Context = request.Context
	}

	// Generate AI response
	aiResponse, err := h.generateAIResponse(session, request)
	if err != nil {
		h.logger.WithError(err).Error("Failed to generate AI response")
		return ChatResponse{
			SessionID: session.ID,
			Success:   false,
			Error:     "Failed to generate response",
		}
	}

	// Add AI message to session
	aiMessage := ChatMessage{
		ID:        generateID(),
		Type:      "assistant",
		Content:   aiResponse,
		Timestamp: time.Now(),
		SessionID: session.ID,
	}

	h.addMessageToSession(session.ID, aiMessage)

	return ChatResponse{
		Message:   aiMessage,
		SessionID: session.ID,
		Success:   true,
	}
}

// getOrCreateSession gets an existing session or creates a new one
func (h *ChatHandler) getOrCreateSession(sessionID, connID string) *ChatSession {
	h.sessionsMutex.Lock()
	defer h.sessionsMutex.Unlock()

	if sessionID != "" {
		if session, exists := h.sessions[sessionID]; exists {
			session.LastActivity = time.Now()
			return session
		}
	}

	// Create new session
	newSessionID := generateID()
	session := &ChatSession{
		ID:           newSessionID,
		ConsultantID: connID, // Using connection ID as consultant ID for now
		Messages:     make([]ChatMessage, 0),
		CreatedAt:    time.Now(),
		LastActivity: time.Now(),
	}

	h.sessions[newSessionID] = session
	return session
}

// addMessageToSession adds a message to a session
func (h *ChatHandler) addMessageToSession(sessionID string, message ChatMessage) {
	h.sessionsMutex.Lock()
	defer h.sessionsMutex.Unlock()

	if session, exists := h.sessions[sessionID]; exists {
		session.Messages = append(session.Messages, message)
		session.LastActivity = time.Now()
	}
}

// generateAIResponse generates an AI response using enhanced Bedrock with knowledge base
func (h *ChatHandler) generateAIResponse(session *ChatSession, request ChatRequest) (string, error) {
	// Try to use enhanced response if we can create a mock inquiry from the chat context
	if session.ClientName != "" && session.Context != "" {
		// Create a mock inquiry from chat context for knowledge base integration
		mockInquiry := h.createMockInquiryFromChat(session, request)

		// Configure Bedrock options for chat responses
		// Use the configured model ID from the Bedrock service
		modelInfo := h.bedrockService.GetModelInfo()
		options := &interfaces.BedrockOptions{
			ModelID:     modelInfo.ModelID,
			MaxTokens:   1000,
			Temperature: 0.7,
			TopP:        0.9,
		}

		// Try enhanced response first
		aiStartTime := time.Now()
		response, err := h.enhancedBedrock.GenerateEnhancedResponse(context.Background(), mockInquiry, options)
		aiResponseTime := time.Since(aiStartTime)

		if err == nil {
			// Record successful AI request
			if h.metricsCollector != nil {
				tokensUsed := int64(len(response.Content) / 4) // Rough estimate
				cost := float64(tokensUsed) * 0.00001          // Rough cost estimate
				h.metricsCollector.RecordAIRequest(true, aiResponseTime, tokensUsed, "claude-3-sonnet", cost)
			}
			return response.Content, nil
		}

		// Record failed AI request
		if h.metricsCollector != nil {
			h.metricsCollector.RecordAIRequest(false, aiResponseTime, 0, "claude-3-sonnet", 0)
			h.metricsCollector.RecordError("system")
		}

		// Log the error but continue with fallback
		h.logger.WithError(err).Warn("Enhanced response failed, falling back to standard response")
	}

	// Fallback to standard consultant prompt
	prompt := h.buildConsultantPrompt(session, request)

	// Configure Bedrock options for chat responses
	// Use the configured model ID from the Bedrock service
	modelInfo := h.bedrockService.GetModelInfo()
	options := &interfaces.BedrockOptions{
		ModelID:     modelInfo.ModelID,
		MaxTokens:   1000,
		Temperature: 0.7,
		TopP:        0.9,
	}

	// Generate response using standard Bedrock
	aiStartTime := time.Now()
	response, err := h.bedrockService.GenerateText(context.Background(), prompt, options)
	aiResponseTime := time.Since(aiStartTime)

	if err != nil {
		// Record failed AI request
		if h.metricsCollector != nil {
			h.metricsCollector.RecordAIRequest(false, aiResponseTime, 0, "claude-3-sonnet", 0)
			h.metricsCollector.RecordError("system")
		}
		return "", err
	}

	// Record successful AI request
	if h.metricsCollector != nil {
		tokensUsed := int64(len(response.Content) / 4) // Rough estimate
		cost := float64(tokensUsed) * 0.00001          // Rough cost estimate
		h.metricsCollector.RecordAIRequest(true, aiResponseTime, tokensUsed, "claude-3-sonnet", cost)
	}

	return response.Content, nil
}

// createMockInquiryFromChat creates a mock inquiry from chat session for knowledge base integration
func (h *ChatHandler) createMockInquiryFromChat(session *ChatSession, request ChatRequest) *domain.Inquiry {
	// Import domain package at the top of the file
	return &domain.Inquiry{
		ID:        session.ID,
		Name:      "Chat Session",
		Company:   session.ClientName,
		Message:   request.Message,
		Services:  h.extractServicesFromContext(session.Context, request.Message),
		CreatedAt: session.CreatedAt,
		UpdatedAt: time.Now(),
	}
}

// extractServicesFromContext attempts to extract service types from chat context
func (h *ChatHandler) extractServicesFromContext(context, message string) []string {
	services := make([]string, 0)

	// Simple keyword matching to infer services
	combinedText := strings.ToLower(context + " " + message)

	if strings.Contains(combinedText, "migration") || strings.Contains(combinedText, "migrate") {
		services = append(services, "Cloud Migration")
	}
	if strings.Contains(combinedText, "architecture") || strings.Contains(combinedText, "design") {
		services = append(services, "Architecture Review")
	}
	if strings.Contains(combinedText, "assessment") || strings.Contains(combinedText, "evaluate") {
		services = append(services, "Assessment")
	}
	if strings.Contains(combinedText, "optimization") || strings.Contains(combinedText, "optimize") {
		services = append(services, "Optimization")
	}
	if strings.Contains(combinedText, "security") || strings.Contains(combinedText, "secure") {
		services = append(services, "Security Review")
	}

	// Default to general consulting if no specific services identified
	if len(services) == 0 {
		services = append(services, "General Consulting")
	}

	return services
}

// buildConsultantPrompt builds a context-aware prompt for the consultant chat
func (h *ChatHandler) buildConsultantPrompt(session *ChatSession, request ChatRequest) string {
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

	// Add conversation history (last 5 messages for context)
	if len(session.Messages) > 0 {
		basePrompt += "\n\nCONVERSATION HISTORY:"
		start := len(session.Messages) - 5
		if start < 0 {
			start = 0
		}

		for _, msg := range session.Messages[start:] {
			if msg.Type == "user" {
				basePrompt += "\nConsultant: " + msg.Content
			} else if msg.Type == "assistant" {
				basePrompt += "\nAssistant: " + msg.Content
			}
		}
	}

	// Handle quick actions
	if request.QuickAction != "" {
		basePrompt += "\n\nQUICK ACTION REQUESTED: " + request.QuickAction
	}

	// Add current question
	basePrompt += "\n\nCURRENT QUESTION: " + request.Message
	basePrompt += "\n\nProvide a direct, expert-level response that the consultant can immediately use with their client:"

	return basePrompt
}

// GetChatSessions returns all active chat sessions (for admin monitoring)
func (h *ChatHandler) GetChatSessions(c *gin.Context) {
	h.sessionsMutex.RLock()
	defer h.sessionsMutex.RUnlock()

	sessions := make([]*ChatSession, 0, len(h.sessions))
	for _, session := range h.sessions {
		sessions = append(sessions, session)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    sessions,
	})
}

// GetChatSession returns a specific chat session
func (h *ChatHandler) GetChatSession(c *gin.Context) {
	sessionID := c.Param("sessionId")

	h.sessionsMutex.RLock()
	session, exists := h.sessions[sessionID]
	h.sessionsMutex.RUnlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Session not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    session,
	})
}

// REST API Endpoints for Chat Management

// CreateChatSession creates a new chat session with security validation
func (h *ChatHandler) CreateChatSession(c *gin.Context) {
	fmt.Println("[CHAT DEBUG] CreateChatSession called")
	ctx := context.Background()

	// Get user information from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "User authentication required",
		})
		return
	}

	userIDStr := userID.(string)

	// Check rate limit for session creation
	rateLimitResult, err := h.chatSecurityService.CheckRateLimit(ctx, userIDStr, "session_creation")
	if err != nil {
		h.logger.WithError(err).Warn("Failed to check session creation rate limit")
	} else if !rateLimitResult.Allowed {
		h.logger.WithField("user_id", userIDStr).Warn("Session creation rate limit exceeded")

		// Log rate limit exceeded
		metadata := map[string]interface{}{
			"action":     "session_creation",
			"ip_address": c.ClientIP(),
			"user_agent": c.GetHeader("User-Agent"),
		}
		h.chatAuditLogger.LogRateLimitExceeded(ctx, userIDStr, "session_creation", metadata)

		c.JSON(http.StatusTooManyRequests, gin.H{
			"success":     false,
			"error":       "Session creation rate limit exceeded",
			"retry_after": int(rateLimitResult.RetryAfter.Seconds()),
		})
		return
	}

	// Parse request body
	var request struct {
		ClientName string                 `json:"client_name"`
		Context    string                 `json:"context"`
		Metadata   map[string]interface{} `json:"metadata"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		h.logger.WithError(err).Error("Failed to parse create session request")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request format",
		})
		return
	}

	// Validate and sanitize input data
	if request.ClientName != "" {
		if err := h.chatSecurityService.ValidateMessageContent(request.ClientName); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "Invalid client name: " + err.Error(),
			})
			return
		}
		request.ClientName = h.chatSecurityService.SanitizeMessageContent(request.ClientName)
	}

	if request.Context != "" {
		if err := h.chatSecurityService.ValidateMessageContent(request.Context); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "Invalid context: " + err.Error(),
			})
			return
		}
		request.Context = h.chatSecurityService.SanitizeMessageContent(request.Context)
	}

	// Create new session
	session := &domain.ChatSession{
		ID:           generateID(),
		UserID:       userIDStr,
		ClientName:   request.ClientName,
		Context:      request.Context,
		Status:       domain.SessionStatusActive,
		Metadata:     request.Metadata,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		LastActivity: time.Now(),
		ExpiresAt:    timePtr(time.Now().Add(24 * time.Hour)), // 24 hour expiration
	}

	// Initialize metadata if nil
	if session.Metadata == nil {
		session.Metadata = make(map[string]interface{})
	}
	session.Metadata["created_via"] = "rest_api"

	// Create session using session service
	if err := h.sessionService.CreateSession(ctx, session); err != nil {
		h.logger.WithError(err).Error("Failed to create session")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to create session",
		})
		return
	}

	// Log session creation for audit
	auditMetadata := map[string]interface{}{
		"client_name": request.ClientName,
		"context":     request.Context,
		"ip_address":  c.ClientIP(),
		"user_agent":  c.GetHeader("User-Agent"),
		"created_via": "rest_api",
	}
	h.chatAuditLogger.LogSessionCreated(ctx, userIDStr, session.ID, auditMetadata)

	h.logger.WithFields(logrus.Fields{
		"session_id": session.ID,
		"user_id":    userIDStr,
	}).Info("Chat session created via REST API")

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    session,
	})
}

// ListChatSessions returns all chat sessions for the authenticated user
func (h *ChatHandler) ListChatSessions(c *gin.Context) {
	// Get user information from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "User authentication required",
		})
		return
	}

	userIDStr := userID.(string)

	// Parse query parameters
	status := c.Query("status")
	limitStr := c.DefaultQuery("limit", "50")
	offsetStr := c.DefaultQuery("offset", "0")

	limit := 50
	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
		limit = l
	}

	offset := 0
	if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
		offset = o
	}

	// Build filters
	filters := &domain.ChatSessionFilters{
		UserID: userIDStr,
		Limit:  limit,
		Offset: offset,
	}

	if status != "" {
		filters.Status = domain.SessionStatus(status)
	}

	// Get sessions using session service
	ctx := context.Background()
	sessions, err := h.sessionService.ListSessions(ctx, filters)
	if err != nil {
		h.logger.WithError(err).WithField("user_id", userIDStr).Error("Failed to list sessions")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to retrieve sessions",
		})
		return
	}

	// Get total count for pagination
	totalCount, err := h.sessionService.GetSessionCount(ctx, &domain.ChatSessionFilters{UserID: userIDStr})
	if err != nil {
		h.logger.WithError(err).WithField("user_id", userIDStr).Warn("Failed to get session count")
		totalCount = int64(len(sessions))
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"sessions":    sessions,
			"total_count": totalCount,
			"limit":       limit,
			"offset":      offset,
		},
	})
}

// GetChatSessionByID returns a specific chat session by ID
func (h *ChatHandler) GetChatSessionByID(c *gin.Context) {
	sessionID := c.Param("id")

	// Get user information from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "User authentication required",
		})
		return
	}

	userIDStr := userID.(string)

	// Validate session belongs to user and is accessible
	ctx := context.Background()
	session, err := h.sessionService.ValidateSession(ctx, sessionID, userIDStr)
	if err != nil {
		h.logger.WithError(err).WithFields(logrus.Fields{
			"session_id": sessionID,
			"user_id":    userIDStr,
		}).Error("Failed to validate session")

		// Check if it's a validation error vs other error
		if _, ok := err.(*interfaces.SessionValidationError); ok {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error":   "Session not found or access denied",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Failed to retrieve session",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    session,
	})
}

// UpdateChatSession updates a chat session's context
func (h *ChatHandler) UpdateChatSession(c *gin.Context) {
	sessionID := c.Param("id")

	// Get user information from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "User authentication required",
		})
		return
	}

	userIDStr := userID.(string)

	// Parse request body
	var request struct {
		ClientName     string            `json:"client_name"`
		Context        string            `json:"context"`
		MeetingType    string            `json:"meeting_type"`
		ServiceTypes   []string          `json:"service_types"`
		CloudProviders []string          `json:"cloud_providers"`
		CustomFields   map[string]string `json:"custom_fields"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		h.logger.WithError(err).Error("Failed to parse update session request")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request format",
		})
		return
	}

	// Validate session belongs to user
	ctx := context.Background()
	_, err := h.sessionService.ValidateSession(ctx, sessionID, userIDStr)
	if err != nil {
		h.logger.WithError(err).WithFields(logrus.Fields{
			"session_id": sessionID,
			"user_id":    userIDStr,
		}).Error("Failed to validate session for update")

		if _, ok := err.(*interfaces.SessionValidationError); ok {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error":   "Session not found or access denied",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Failed to retrieve session",
			})
		}
		return
	}

	// Update session context using chat service
	sessionContext := &domain.SessionContext{
		ClientName:     request.ClientName,
		MeetingType:    request.MeetingType,
		ProjectContext: request.Context,
		ServiceTypes:   request.ServiceTypes,
		CloudProviders: request.CloudProviders,
		CustomFields:   request.CustomFields,
	}

	if err := h.chatService.UpdateSessionContext(ctx, sessionID, sessionContext); err != nil {
		h.logger.WithError(err).WithField("session_id", sessionID).Error("Failed to update session context")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to update session context",
		})
		return
	}

	// Get updated session
	updatedSession, err := h.sessionService.GetSession(ctx, sessionID)
	if err != nil {
		h.logger.WithError(err).WithField("session_id", sessionID).Error("Failed to get updated session")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to retrieve updated session",
		})
		return
	}

	h.logger.WithField("session_id", sessionID).Info("Session context updated via REST API")

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    updatedSession,
	})
}

// DeleteChatSession deletes a chat session
func (h *ChatHandler) DeleteChatSession(c *gin.Context) {
	sessionID := c.Param("id")

	// Get user information from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "User authentication required",
		})
		return
	}

	userIDStr := userID.(string)

	// Validate session belongs to user
	ctx := context.Background()
	_, err := h.sessionService.ValidateSession(ctx, sessionID, userIDStr)
	if err != nil {
		h.logger.WithError(err).WithFields(logrus.Fields{
			"session_id": sessionID,
			"user_id":    userIDStr,
		}).Error("Failed to validate session for deletion")

		if _, ok := err.(*interfaces.SessionValidationError); ok {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error":   "Session not found or access denied",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Failed to retrieve session",
			})
		}
		return
	}

	// Delete the session
	if err := h.sessionService.DeleteSession(ctx, sessionID); err != nil {
		h.logger.WithError(err).WithField("session_id", sessionID).Error("Failed to delete session")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to delete session",
		})
		return
	}

	h.logger.WithFields(logrus.Fields{
		"session_id": sessionID,
		"user_id":    userIDStr,
	}).Info("Session deleted via REST API")

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Session deleted successfully",
	})
}

// GetChatSessionHistory returns the message history for a chat session
func (h *ChatHandler) GetChatSessionHistory(c *gin.Context) {
	sessionID := c.Param("id")

	// Get user information from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "User authentication required",
		})
		return
	}

	userIDStr := userID.(string)

	// Parse query parameters
	limitStr := c.DefaultQuery("limit", "50")

	limit := 50
	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 1000 {
		limit = l
	}

	// Validate session belongs to user
	ctx := context.Background()
	_, err := h.sessionService.ValidateSession(ctx, sessionID, userIDStr)
	if err != nil {
		h.logger.WithError(err).WithFields(logrus.Fields{
			"session_id": sessionID,
			"user_id":    userIDStr,
		}).Error("Failed to validate session for history retrieval")

		if _, ok := err.(*interfaces.SessionValidationError); ok {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error":   "Session not found or access denied",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Failed to retrieve session",
			})
		}
		return
	}

	// Get message history using chat service
	messages, err := h.chatService.GetSessionHistory(ctx, sessionID, limit)
	if err != nil {
		h.logger.WithError(err).WithField("session_id", sessionID).Error("Failed to get session history")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to retrieve message history",
		})
		return
	}

	// Get message statistics
	stats, err := h.chatService.GetMessageStats(ctx, sessionID)
	if err != nil {
		h.logger.WithError(err).WithField("session_id", sessionID).Warn("Failed to get message stats")
		stats = nil
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"messages": messages,
			"stats":    stats,
			"limit":    limit,
		},
	})
}

// generateID generates a simple ID for messages and sessions
func generateID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(6)
}

// randomString generates a random string of specified length
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}

// timePtr returns a pointer to the given time
func timePtr(t time.Time) *time.Time {
	return &t
}

// BroadcastToSession broadcasts a message to all connections in a session
func (h *ChatHandler) BroadcastToSession(sessionID string, message interface{}) error {
	connections := h.connectionPool.GetBySessionID(sessionID)
	if len(connections) == 0 {
		return fmt.Errorf("no active connections for session %s", sessionID)
	}

	var errors []string
	for _, conn := range connections {
		if err := conn.Conn.WriteJSON(message); err != nil {
			h.logger.WithError(err).WithFields(logrus.Fields{
				"session_id": sessionID,
				"user_id":    conn.UserID,
			}).Error("Failed to broadcast message to connection")
			errors = append(errors, err.Error())
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to broadcast to some connections: %v", errors)
	}

	return nil
}

// BroadcastToUser broadcasts a message to all connections for a user
func (h *ChatHandler) BroadcastToUser(userID string, message interface{}) error {
	connections := h.connectionPool.GetByUserID(userID)
	if len(connections) == 0 {
		return fmt.Errorf("no active connections for user %s", userID)
	}

	var errors []string
	for _, conn := range connections {
		if err := conn.Conn.WriteJSON(message); err != nil {
			h.logger.WithError(err).WithFields(logrus.Fields{
				"user_id": userID,
			}).Error("Failed to broadcast message to user connection")
			errors = append(errors, err.Error())
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to broadcast to some connections: %v", errors)
	}

	return nil
}

// messageSender handles outgoing messages for a connection
func (h *ChatHandler) messageSender(connection *Connection, connID string) {
	defer func() {
		if r := recover(); r != nil {
			h.logger.WithField("connection_id", connID).Error("Message sender panic recovered")
		}
	}()

	// Start retry routine
	retryTicker := time.NewTicker(30 * time.Second)
	defer retryTicker.Stop()

	for {
		select {
		case message := <-connection.SendChan:
			if err := connection.Conn.WriteJSON(message); err != nil {
				h.logger.WithError(err).WithField("connection_id", connID).Error("Failed to send message")

				// Add to pending messages for retry if it's an important message
				if message.Type == WSMessageTypeMessage && message.MessageID != "" {
					h.addToPendingMessages(connection, message)
				}
				return
			}

			// Add to pending messages if it requires acknowledgment
			if h.requiresAcknowledgment(message) {
				h.addToPendingMessages(connection, message)
			}

		case <-retryTicker.C:
			h.retryPendingMessages(connection, connID)

		case <-connection.CloseChan:
			return
		}
	}
}

// requiresAcknowledgment checks if a message requires acknowledgment
func (h *ChatHandler) requiresAcknowledgment(message WebSocketMessage) bool {
	return message.Type == WSMessageTypeMessage && message.MessageID != ""
}

// addToPendingMessages adds a message to the pending messages queue
func (h *ChatHandler) addToPendingMessages(connection *Connection, message WebSocketMessage) {
	if message.MessageID == "" {
		return
	}

	connection.PendingMutex.Lock()
	defer connection.PendingMutex.Unlock()

	connection.PendingMessages[message.MessageID] = &PendingMessage{
		Message:    message,
		Timestamp:  time.Now(),
		Retries:    0,
		MaxRetries: 3,
	}
}

// retryPendingMessages retries pending messages that haven't been acknowledged
func (h *ChatHandler) retryPendingMessages(connection *Connection, connID string) {
	connection.PendingMutex.Lock()
	defer connection.PendingMutex.Unlock()

	now := time.Now()
	for messageID, pending := range connection.PendingMessages {
		// Retry messages older than 30 seconds
		if now.Sub(pending.Timestamp) > 30*time.Second {
			if pending.Retries < pending.MaxRetries {
				// Retry the message
				if err := connection.Conn.WriteJSON(pending.Message); err != nil {
					h.logger.WithError(err).WithFields(logrus.Fields{
						"connection_id": connID,
						"message_id":    messageID,
					}).Error("Failed to retry message")
					delete(connection.PendingMessages, messageID)
				} else {
					pending.Retries++
					pending.Timestamp = now
					h.logger.WithFields(logrus.Fields{
						"connection_id": connID,
						"message_id":    messageID,
						"retry_count":   pending.Retries,
					}).Debug("Retried message")
				}
			} else {
				// Max retries reached, remove from pending
				h.logger.WithFields(logrus.Fields{
					"connection_id": connID,
					"message_id":    messageID,
				}).Warn("Message max retries reached, giving up")
				delete(connection.PendingMessages, messageID)
			}
		}
	}
}

// routeWebSocketMessage routes incoming WebSocket messages based on type
func (h *ChatHandler) routeWebSocketMessage(wsMessage WebSocketMessage, connection *Connection, connID, userID string) {
	switch wsMessage.Type {
	case WSMessageTypeMessage:
		h.handleChatMessage(wsMessage, connection, connID, userID)
	case WSMessageTypeTyping:
		h.handleTypingIndicator(wsMessage, connection, connID, userID)
	case WSMessageTypePresence:
		h.handlePresenceUpdate(wsMessage, connection, connID, userID)
	case WSMessageTypeHeartbeat:
		h.handleHeartbeat(wsMessage, connection, connID, userID)
	case WSMessageTypeAck:
		h.handleMessageAcknowledgment(wsMessage, connection, connID, userID)
	default:
		h.logger.WithFields(logrus.Fields{
			"message_type":  wsMessage.Type,
			"connection_id": connID,
		}).Warn("Unknown WebSocket message type")
	}
}

// handleMessageAcknowledgment processes message acknowledgments
func (h *ChatHandler) handleMessageAcknowledgment(wsMessage WebSocketMessage, connection *Connection, connID, userID string) {
	if wsMessage.MessageID == "" {
		h.logger.WithField("connection_id", connID).Warn("Received acknowledgment without message ID")
		return
	}

	// Remove from pending messages
	connection.PendingMutex.Lock()
	delete(connection.PendingMessages, wsMessage.MessageID)
	connection.PendingMutex.Unlock()

	h.logger.WithFields(logrus.Fields{
		"message_id":    wsMessage.MessageID,
		"connection_id": connID,
	}).Debug("Message acknowledged")
}

// handleChatMessage processes chat messages with security validation
func (h *ChatHandler) handleChatMessage(wsMessage WebSocketMessage, connection *Connection, connID, userID string) {
	ctx := context.Background()

	// Validate message content
	if err := h.chatSecurityService.ValidateMessageContent(wsMessage.Content); err != nil {
		h.logger.WithError(err).WithField("user_id", userID).Warn("Message content validation failed")

		// Send error response
		errorMessage := WebSocketMessage{
			Type:      WSMessageTypeError,
			Content:   "Message content violates security policy",
			Timestamp: time.Now(),
		}
		select {
		case connection.SendChan <- errorMessage:
		default:
			h.logger.Warn("Failed to send validation error - channel full")
		}
		return
	}

	// Sanitize message content
	sanitizedContent := h.chatSecurityService.SanitizeMessageContent(wsMessage.Content)

	// Check rate limit for messages
	rateLimitResult, err := h.chatSecurityService.CheckRateLimit(ctx, userID, "message")
	if err != nil {
		h.logger.WithError(err).Warn("Failed to check message rate limit")
	} else if !rateLimitResult.Allowed {
		h.logger.WithField("user_id", userID).Warn("Message rate limit exceeded")

		// Log rate limit exceeded
		metadata := map[string]interface{}{
			"action": "message",
		}
		h.chatAuditLogger.LogRateLimitExceeded(ctx, userID, "message", metadata)

		// Send rate limit error
		errorMessage := WebSocketMessage{
			Type:      WSMessageTypeError,
			Content:   fmt.Sprintf("Rate limit exceeded. Please wait %d seconds before sending another message.", int(rateLimitResult.RetryAfter.Seconds())),
			Timestamp: time.Now(),
		}
		select {
		case connection.SendChan <- errorMessage:
		default:
			h.logger.Warn("Failed to send rate limit error - channel full")
		}
		return
	}

	// Convert WebSocket message to ChatRequest with sanitized content
	request := ChatRequest{
		Message:   sanitizedContent,
		SessionID: wsMessage.SessionID,
	}

	// Extract metadata
	if wsMessage.Metadata != nil {
		if clientName, ok := wsMessage.Metadata["client_name"].(string); ok {
			request.ClientName = clientName
		}
		if context, ok := wsMessage.Metadata["context"].(string); ok {
			request.Context = context
		}
		if quickAction, ok := wsMessage.Metadata["quick_action"].(string); ok {
			request.QuickAction = quickAction
		}
	}

	// Process the chat request
	response := h.processEnhancedChatRequest(request, connID, userID)

	// Convert to WebSocket message and send
	responseMessage := WebSocketMessage{
		Type:      WSMessageTypeMessage,
		SessionID: response.SessionID,
		MessageID: response.Message.ID,
		Content:   response.Message.Content,
		Timestamp: response.Message.Timestamp,
		Metadata: map[string]interface{}{
			"message_type": response.Message.Type,
			"success":      response.Success,
			"error":        response.Error,
		},
	}

	// Send acknowledgment for the incoming message
	ackMessage := WebSocketMessage{
		Type:      WSMessageTypeAck,
		MessageID: wsMessage.MessageID,
		SessionID: wsMessage.SessionID,
		Timestamp: time.Now(),
		Metadata: map[string]interface{}{
			"status": "received",
		},
	}

	// Send both messages
	select {
	case connection.SendChan <- ackMessage:
	default:
		h.logger.Warn("Failed to send acknowledgment - channel full")
	}

	select {
	case connection.SendChan <- responseMessage:
	default:
		h.logger.Warn("Failed to send response message - channel full")
	}

	// Broadcast to other connections in the same session if needed
	if wsMessage.SessionID != "" {
		h.broadcastToSessionExcept(wsMessage.SessionID, responseMessage, connID)
	}
}

// handleTypingIndicator processes typing indicator messages
func (h *ChatHandler) handleTypingIndicator(wsMessage WebSocketMessage, connection *Connection, connID, userID string) {
	var typingData TypingIndicator
	if err := json.Unmarshal([]byte(wsMessage.Content), &typingData); err != nil {
		h.logger.WithError(err).Error("Failed to parse typing indicator")
		return
	}

	// Update connection typing status
	connection.IsTyping = typingData.IsTyping

	// Broadcast typing indicator to other connections in the session
	if wsMessage.SessionID != "" {
		typingMessage := WebSocketMessage{
			Type:      WSMessageTypeTyping,
			SessionID: wsMessage.SessionID,
			Content:   wsMessage.Content,
			Timestamp: time.Now(),
			Metadata: map[string]interface{}{
				"user_id": userID,
			},
		}
		h.broadcastToSessionExcept(wsMessage.SessionID, typingMessage, connID)
	}
}

// handlePresenceUpdate processes presence update messages
func (h *ChatHandler) handlePresenceUpdate(wsMessage WebSocketMessage, connection *Connection, connID, userID string) {
	var presenceData PresenceUpdate
	if err := json.Unmarshal([]byte(wsMessage.Content), &presenceData); err != nil {
		h.logger.WithError(err).Error("Failed to parse presence update")
		return
	}

	// Update connection metadata
	connection.Metadata["presence_status"] = presenceData.Status

	// Broadcast presence update to other connections for the same user
	presenceMessage := WebSocketMessage{
		Type:      WSMessageTypePresence,
		SessionID: wsMessage.SessionID,
		Content:   wsMessage.Content,
		Timestamp: time.Now(),
		Metadata: map[string]interface{}{
			"user_id": userID,
		},
	}
	h.broadcastToUserExcept(userID, presenceMessage, connID)
}

// handleHeartbeat processes heartbeat messages
func (h *ChatHandler) handleHeartbeat(wsMessage WebSocketMessage, connection *Connection, connID, userID string) {
	// Update last ping time
	connection.LastPing = time.Now()

	// Send heartbeat response
	heartbeatResponse := WebSocketMessage{
		Type:      WSMessageTypeHeartbeat,
		Timestamp: time.Now(),
		Metadata: map[string]interface{}{
			"status": "alive",
		},
	}

	select {
	case connection.SendChan <- heartbeatResponse:
	default:
		h.logger.Warn("Failed to send heartbeat response - channel full")
	}
}

// broadcastToSessionExcept broadcasts a message to all connections in a session except the sender
func (h *ChatHandler) broadcastToSessionExcept(sessionID string, message WebSocketMessage, excludeConnID string) {
	connections := h.connectionPool.GetBySessionID(sessionID)
	for _, conn := range connections {
		// Find the connection ID for this connection
		h.connectionPool.mutex.RLock()
		var currentConnID string
		for id, c := range h.connectionPool.connections {
			if c == conn {
				currentConnID = id
				break
			}
		}
		h.connectionPool.mutex.RUnlock()

		if currentConnID != excludeConnID {
			select {
			case conn.SendChan <- message:
			default:
				h.logger.WithField("session_id", sessionID).Warn("Failed to broadcast to session connection - channel full")
			}
		}
	}
}

// broadcastToUserExcept broadcasts a message to all connections for a user except the sender
func (h *ChatHandler) broadcastToUserExcept(userID string, message WebSocketMessage, excludeConnID string) {
	connections := h.connectionPool.GetByUserID(userID)
	for _, conn := range connections {
		// Find the connection ID for this connection
		h.connectionPool.mutex.RLock()
		var currentConnID string
		for id, c := range h.connectionPool.connections {
			if c == conn {
				currentConnID = id
				break
			}
		}
		h.connectionPool.mutex.RUnlock()

		if currentConnID != excludeConnID {
			select {
			case conn.SendChan <- message:
			default:
				h.logger.WithField("user_id", userID).Warn("Failed to broadcast to user connection - channel full")
			}
		}
	}
}

// GetConnectionStats returns statistics about active connections
func (h *ChatHandler) GetConnectionStats() map[string]interface{} {
	return map[string]interface{}{
		"total_connections": h.connectionPool.Count(),
		"timestamp":         time.Now(),
	}
}

// GetAuthMiddleware returns the authentication middleware
func (h *ChatHandler) GetAuthMiddleware() interfaces.ChatAuthMiddleware {
	return h.authMiddleware
}

// GetSecurityService returns the security service
func (h *ChatHandler) GetSecurityService() interfaces.ChatSecurityService {
	return h.chatSecurityService
}

// GetAuditLogger returns the audit logger
func (h *ChatHandler) GetAuditLogger() interfaces.ChatAuditLogger {
	return h.chatAuditLogger
}

// GetRateLimiter returns the rate limiter
func (h *ChatHandler) GetRateLimiter() interfaces.ChatRateLimiter {
	return h.chatRateLimiter
}

// GetConnectionPool returns the connection pool for health monitoring
func (h *ChatHandler) GetConnectionPool() *ConnectionPool {
	return h.connectionPool
}

// GetAuthService returns the authentication service
func (h *ChatHandler) GetAuthService() interfaces.ChatAuthService {
	return h.chatAuthService
}

// getEnvAsSlice parses environment variable as comma-separated slice
func getEnvAsSlice(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		return strings.Split(value, ",")
	}
	return defaultValue
}
