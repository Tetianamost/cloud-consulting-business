package handlers

import (
	"context"
	"net/http"
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

// WebSocket upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// In production, implement proper origin checking
		return true
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// ChatHandler handles real-time consultant chat
type ChatHandler struct {
	logger           *logrus.Logger
	bedrockService   interfaces.BedrockService
	enhancedBedrock  *services.EnhancedBedrockService
	knowledgeBase    interfaces.KnowledgeBase
	sessions         map[string]*ChatSession
	connections      map[string]*websocket.Conn
	sessionsMutex    sync.RWMutex
	connectionsMutex sync.RWMutex
}

// NewChatHandler creates a new chat handler
func NewChatHandler(logger *logrus.Logger, bedrockService interfaces.BedrockService, knowledgeBase interfaces.KnowledgeBase) *ChatHandler {
	clientHistoryService := services.NewClientHistoryService(knowledgeBase)
	companyKnowledgeIntegrationService := services.NewCompanyKnowledgeIntegrationService(knowledgeBase, clientHistoryService)
	enhancedBedrock := services.NewEnhancedBedrockService(bedrockService, knowledgeBase, companyKnowledgeIntegrationService)

	return &ChatHandler{
		logger:          logger,
		bedrockService:  bedrockService,
		enhancedBedrock: enhancedBedrock,
		knowledgeBase:   knowledgeBase,
		sessions:        make(map[string]*ChatSession),
		connections:     make(map[string]*websocket.Conn),
	}
}

// HandleWebSocket handles WebSocket connections for real-time chat
func (h *ChatHandler) HandleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		h.logger.WithError(err).Error("Failed to upgrade WebSocket connection")
		return
	}
	defer conn.Close()

	// Generate connection ID
	connID := generateID()

	h.connectionsMutex.Lock()
	h.connections[connID] = conn
	h.connectionsMutex.Unlock()

	defer func() {
		h.connectionsMutex.Lock()
		delete(h.connections, connID)
		h.connectionsMutex.Unlock()
	}()

	h.logger.WithField("connection_id", connID).Info("New WebSocket connection established")

	// Handle incoming messages
	for {
		var request ChatRequest
		err := conn.ReadJSON(&request)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				h.logger.WithError(err).Error("WebSocket error")
			}
			break
		}

		// Process the chat request
		response := h.processChatRequest(request, connID)

		// Send response back to client
		if err := conn.WriteJSON(response); err != nil {
			h.logger.WithError(err).Error("Failed to send WebSocket response")
			break
		}
	}
}

// processChatRequest processes an incoming chat request
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
		options := &interfaces.BedrockOptions{
			ModelID:     "anthropic.claude-3-sonnet-20240229-v1:0",
			MaxTokens:   1000,
			Temperature: 0.7,
			TopP:        0.9,
		}

		// Try enhanced response first
		response, err := h.enhancedBedrock.GenerateEnhancedResponse(context.Background(), mockInquiry, options)
		if err == nil {
			return response.Content, nil
		}

		// Log the error but continue with fallback
		h.logger.WithError(err).Warn("Enhanced response failed, falling back to standard response")
	}

	// Fallback to standard consultant prompt
	prompt := h.buildConsultantPrompt(session, request)

	// Configure Bedrock options for chat responses
	options := &interfaces.BedrockOptions{
		ModelID:     "anthropic.claude-3-sonnet-20240229-v1:0",
		MaxTokens:   1000,
		Temperature: 0.7,
		TopP:        0.9,
	}

	// Generate response using standard Bedrock
	response, err := h.bedrockService.GenerateText(context.Background(), prompt, options)
	if err != nil {
		return "", err
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
