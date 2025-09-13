package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cloud-consulting/backend/internal/domain"
	"github.com/cloud-consulting/backend/internal/handlers"
	"github.com/cloud-consulting/backend/internal/interfaces"
	"github.com/cloud-consulting/backend/internal/services"
	"github.com/cloud-consulting/backend/internal/storage"
)

// Test configuration
const (
	testJWTSecret = "test-jwt-secret-for-integration-tests"
	testDBURL     = "postgres://test:test@localhost:5432/chat_test?sslmode=disable"
)

// Integration test suite
type ChatAPIIntegrationTestSuite struct {
	server        *httptest.Server
	chatHandler   *handlers.ChatHandler
	sessionRepo   interfaces.ChatSessionRepository
	messageRepo   interfaces.ChatMessageRepository
	sessionSvc    interfaces.SessionService
	chatSvc       interfaces.ChatService
	bedrockSvc    interfaces.BedrockService
	knowledgeBase interfaces.KnowledgeBase
	logger        *logrus.Logger
}

// Mock Bedrock service for integration tests
type MockBedrockIntegration struct {
	responses map[string]string
}

func NewMockBedrockIntegration() *MockBedrockIntegration {
	return &MockBedrockIntegration{
		responses: map[string]string{
			"default":  "This is a mock AI response for integration testing.",
			"cost":     "For cost estimation, consider EC2 instances starting at $0.0116/hour for t3.nano.",
			"security": "AWS security best practices include enabling MFA, using IAM roles, and encrypting data at rest.",
		},
	}
}

func (m *MockBedrockIntegration) GenerateText(ctx context.Context, prompt string, options *interfaces.BedrockOptions) (*interfaces.BedrockResponse, error) {
	// Simulate processing time
	time.Sleep(50 * time.Millisecond)

	// Determine response based on prompt content
	response := m.responses["default"]
	promptLower := strings.ToLower(prompt)

	if strings.Contains(promptLower, "cost") {
		response = m.responses["cost"]
	} else if strings.Contains(promptLower, "security") {
		response = m.responses["security"]
	}

	return &interfaces.BedrockResponse{
		Content: response,
		Usage: interfaces.BedrockUsage{
			InputTokens:  len(prompt) / 4, // Rough token estimation
			OutputTokens: len(response) / 4,
		},
	}, nil
}

func (m *MockBedrockIntegration) GetModelInfo() interfaces.BedrockModelInfo {
	return interfaces.BedrockModelInfo{
		ModelID:     "mock-model-v1",
		MaxTokens:   4000,
		InputCost:   0.001,
		OutputCost:  0.002,
		Description: "Mock model for integration testing",
	}
}

func (m *MockBedrockIntegration) IsHealthy() bool {
	return true
}

// Mock Knowledge Base for integration tests
type MockKnowledgeBaseIntegration struct{}

func NewMockKnowledgeBaseIntegration() *MockKnowledgeBaseIntegration {
	return &MockKnowledgeBaseIntegration{}
}

func (m *MockKnowledgeBaseIntegration) GetServiceOfferings(ctx context.Context) ([]*interfaces.ServiceOffering, error) {
	return []*interfaces.ServiceOffering{
		{
			Name:            "Cloud Migration",
			Description:     "Comprehensive cloud migration services",
			Category:        "Migration",
			TypicalDuration: "3-6 months",
			TeamSize:        "3-5 consultants",
			KeyBenefits:     []string{"Reduced costs", "Improved scalability"},
			Deliverables:    []string{"Migration plan", "Implementation"},
		},
	}, nil
}

func (m *MockKnowledgeBaseIntegration) GetTeamExpertise(ctx context.Context) ([]*interfaces.TeamExpertise, error) {
	return []*interfaces.TeamExpertise{
		{
			ConsultantName:  "John Smith",
			Role:            "Senior Cloud Architect",
			ExperienceYears: 8,
			ExpertiseAreas:  []string{"AWS Migration", "Architecture Design"},
			CloudProviders:  []string{"AWS", "Azure"},
		},
	}, nil
}

func (m *MockKnowledgeBaseIntegration) GetPastSolutions(ctx context.Context, serviceType, industry string) ([]*interfaces.PastSolution, error) {
	return []*interfaces.PastSolution{
		{
			Title:            "E-commerce Migration",
			Industry:         "Retail",
			ProblemStatement: "Legacy infrastructure limitations",
			SolutionApproach: "Lift and shift with optimization",
			Technologies:     []string{"EC2", "RDS", "CloudFront"},
			TimeToValue:      "4 months",
			CostSavings:      50000,
		},
	}, nil
}

func (m *MockKnowledgeBaseIntegration) GetConsultingApproach(ctx context.Context, serviceType string) (*interfaces.ConsultingApproach, error) {
	return &interfaces.ConsultingApproach{
		Name:              "Migration Methodology",
		Philosophy:        "Phased approach with minimal disruption",
		EngagementModel:   "Fixed scope with milestones",
		KeyPrinciples:     []string{"Risk mitigation", "Business continuity"},
		ClientInvolvement: "Active collaboration required",
		KnowledgeTransfer: "Comprehensive training provided",
	}, nil
}

func (m *MockKnowledgeBaseIntegration) GetClientHistory(ctx context.Context, company string) ([]*interfaces.ClientEngagement, error) {
	return []*interfaces.ClientEngagement{
		{
			ProjectName:        "Previous Migration",
			StartDate:          time.Now().AddDate(-1, 0, 0),
			Status:             "Completed",
			ClientSatisfaction: 9.2,
		},
	}, nil
}

// Setup test suite
func setupIntegrationTestSuite(t *testing.T) *ChatAPIIntegrationTestSuite {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel) // Reduce noise in tests

	// Use in-memory repositories for testing
	sessionRepo := storage.NewMemoryChatSessionRepository()
	messageRepo := storage.NewMemoryChatMessageRepository()

	// Create services
	sessionSvc := services.NewSessionService(sessionRepo, logger)
	chatSvc := services.NewChatService(messageRepo, sessionSvc, nil, logger) // Bedrock will be set later

	// Create mock external services
	bedrockSvc := NewMockBedrockIntegration()
	knowledgeBase := NewMockKnowledgeBaseIntegration()

	// Create chat handler
	chatHandler := handlers.NewChatHandler(
		logger,
		bedrockSvc,
		knowledgeBase,
		sessionSvc,
		chatSvc,
		nil, // authHandler not needed for these tests
		testJWTSecret,
	)

	// Setup Gin router
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Add routes
	api := router.Group("/api/v1")
	{
		admin := api.Group("/admin")
		{
			chat := admin.Group("/chat")
			{
				chat.POST("/sessions", chatHandler.CreateChatSession)
				chat.GET("/sessions", chatHandler.GetChatSessions)
				chat.GET("/sessions/:sessionId", chatHandler.GetChatSession)
				chat.PUT("/sessions/:sessionId", chatHandler.UpdateChatSession)
				chat.DELETE("/sessions/:sessionId", chatHandler.DeleteChatSession)
				chat.GET("/sessions/:sessionId/history", chatHandler.GetSessionHistory)
				chat.POST("/sessions/:sessionId/messages", chatHandler.SendMessage)
				chat.GET("/ws", chatHandler.HandleWebSocket)
			}
		}
	}

	// Create test server
	server := httptest.NewServer(router)

	return &ChatAPIIntegrationTestSuite{
		server:        server,
		chatHandler:   chatHandler,
		sessionRepo:   sessionRepo,
		messageRepo:   messageRepo,
		sessionSvc:    sessionSvc,
		chatSvc:       chatSvc,
		bedrockSvc:    bedrockSvc,
		knowledgeBase: knowledgeBase,
		logger:        logger,
	}
}

func (suite *ChatAPIIntegrationTestSuite) tearDown() {
	suite.server.Close()
}

// Helper function to create authenticated request
func (suite *ChatAPIIntegrationTestSuite) createAuthenticatedRequest(method, path string, body interface{}) *http.Request {
	var reqBody *bytes.Buffer
	if body != nil {
		jsonBody, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(jsonBody)
	} else {
		reqBody = bytes.NewBuffer(nil)
	}

	req, _ := http.NewRequest(method, suite.server.URL+path, reqBody)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-token") // Mock token

	return req
}

// Test REST API endpoints
func TestChatAPI_CreateSession_Integration(t *testing.T) {
	suite := setupIntegrationTestSuite(t)
	defer suite.tearDown()

	sessionData := map[string]interface{}{
		"client_name": "Integration Test Client",
		"context":     "Integration testing session",
		"metadata": map[string]interface{}{
			"meeting_type": "consultation",
		},
	}

	req := suite.createAuthenticatedRequest("POST", "/api/v1/admin/chat/sessions", sessionData)

	// Mock user context (normally set by auth middleware)
	req = req.WithContext(context.WithValue(req.Context(), "user_id", "test-user-id"))

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	assert.True(t, response["success"].(bool))
	assert.NotNil(t, response["data"])

	sessionData := response["data"].(map[string]interface{})
	assert.NotEmpty(t, sessionData["id"])
	assert.Equal(t, "Integration Test Client", sessionData["client_name"])
	assert.Equal(t, "Integration testing session", sessionData["context"])
}

func TestChatAPI_SendMessage_Integration(t *testing.T) {
	suite := setupIntegrationTestSuite(t)
	defer suite.tearDown()

	ctx := context.Background()

	// Create a session first
	session := &domain.ChatSession{
		ID:           "test-session-integration",
		UserID:       "test-user-id",
		ClientName:   "Test Client",
		Context:      "Integration test",
		Status:       domain.SessionStatusActive,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		LastActivity: time.Now(),
		ExpiresAt:    timePtr(time.Now().Add(24 * time.Hour)),
	}

	err := suite.sessionSvc.CreateSession(ctx, session)
	require.NoError(t, err)

	// Send a message
	messageData := map[string]interface{}{
		"message":     "What are the best practices for AWS security?",
		"session_id":  session.ID,
		"client_name": "Test Client",
		"context":     "Security consultation",
	}

	req := suite.createAuthenticatedRequest("POST", fmt.Sprintf("/api/v1/admin/chat/sessions/%s/messages", session.ID), messageData)
	req = req.WithContext(context.WithValue(req.Context(), "user_id", "test-user-id"))

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	assert.True(t, response["success"].(bool))
	assert.NotNil(t, response["message"])

	messageResp := response["message"].(map[string]interface{})
	assert.Equal(t, "assistant", messageResp["type"])
	assert.Contains(t, messageResp["content"], "security")
	assert.Equal(t, session.ID, messageResp["session_id"])
}

func TestChatAPI_GetSessionHistory_Integration(t *testing.T) {
	suite := setupIntegrationTestSuite(t)
	defer suite.tearDown()

	ctx := context.Background()

	// Create a session with messages
	session := &domain.ChatSession{
		ID:           "test-session-history",
		UserID:       "test-user-id",
		ClientName:   "Test Client",
		Status:       domain.SessionStatusActive,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		LastActivity: time.Now(),
		ExpiresAt:    timePtr(time.Now().Add(24 * time.Hour)),
	}

	err := suite.sessionSvc.CreateSession(ctx, session)
	require.NoError(t, err)

	// Add some messages
	messages := []*domain.ChatMessage{
		{
			ID:        "msg-1",
			SessionID: session.ID,
			Type:      domain.MessageTypeUser,
			Content:   "Hello",
			Status:    domain.MessageStatusSent,
			CreatedAt: time.Now(),
		},
		{
			ID:        "msg-2",
			SessionID: session.ID,
			Type:      domain.MessageTypeAssistant,
			Content:   "Hi there! How can I help you?",
			Status:    domain.MessageStatusSent,
			CreatedAt: time.Now().Add(time.Second),
		},
	}

	for _, msg := range messages {
		err = suite.messageRepo.Create(ctx, msg)
		require.NoError(t, err)
	}

	// Get session history
	req := suite.createAuthenticatedRequest("GET", fmt.Sprintf("/api/v1/admin/chat/sessions/%s/history", session.ID), nil)
	req = req.WithContext(context.WithValue(req.Context(), "user_id", "test-user-id"))

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	assert.True(t, response["success"].(bool))
	assert.NotNil(t, response["messages"])

	messagesResp := response["messages"].([]interface{})
	assert.Len(t, messagesResp, 2)

	// Check message order (should be chronological)
	firstMsg := messagesResp[0].(map[string]interface{})
	assert.Equal(t, "user", firstMsg["type"])
	assert.Equal(t, "Hello", firstMsg["content"])

	secondMsg := messagesResp[1].(map[string]interface{})
	assert.Equal(t, "assistant", secondMsg["type"])
	assert.Equal(t, "Hi there! How can I help you?", secondMsg["content"])
}

func TestChatAPI_UpdateSession_Integration(t *testing.T) {
	suite := setupIntegrationTestSuite(t)
	defer suite.tearDown()

	ctx := context.Background()

	// Create a session
	session := &domain.ChatSession{
		ID:           "test-session-update",
		UserID:       "test-user-id",
		ClientName:   "Original Client",
		Context:      "Original context",
		Status:       domain.SessionStatusActive,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		LastActivity: time.Now(),
		ExpiresAt:    timePtr(time.Now().Add(24 * time.Hour)),
	}

	err := suite.sessionSvc.CreateSession(ctx, session)
	require.NoError(t, err)

	// Update session
	updateData := map[string]interface{}{
		"client_name": "Updated Client",
		"context":     "Updated context",
		"metadata": map[string]interface{}{
			"meeting_type": "follow-up",
		},
	}

	req := suite.createAuthenticatedRequest("PUT", fmt.Sprintf("/api/v1/admin/chat/sessions/%s", session.ID), updateData)
	req = req.WithContext(context.WithValue(req.Context(), "user_id", "test-user-id"))

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	assert.True(t, response["success"].(bool))

	// Verify session was updated
	updatedSession, err := suite.sessionSvc.GetSession(ctx, session.ID)
	require.NoError(t, err)
	assert.Equal(t, "Updated Client", updatedSession.ClientName)
	assert.Equal(t, "Updated context", updatedSession.Context)
}

func TestChatAPI_DeleteSession_Integration(t *testing.T) {
	suite := setupIntegrationTestSuite(t)
	defer suite.tearDown()

	ctx := context.Background()

	// Create a session
	session := &domain.ChatSession{
		ID:           "test-session-delete",
		UserID:       "test-user-id",
		ClientName:   "Test Client",
		Status:       domain.SessionStatusActive,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		LastActivity: time.Now(),
		ExpiresAt:    timePtr(time.Now().Add(24 * time.Hour)),
	}

	err := suite.sessionSvc.CreateSession(ctx, session)
	require.NoError(t, err)

	// Delete session
	req := suite.createAuthenticatedRequest("DELETE", fmt.Sprintf("/api/v1/admin/chat/sessions/%s", session.ID), nil)
	req = req.WithContext(context.WithValue(req.Context(), "user_id", "test-user-id"))

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	assert.True(t, response["success"].(bool))

	// Verify session was deleted
	_, err = suite.sessionSvc.GetSession(ctx, session.ID)
	assert.Error(t, err)
}

// WebSocket integration tests
func TestChatAPI_WebSocket_Integration(t *testing.T) {
	suite := setupIntegrationTestSuite(t)
	defer suite.tearDown()

	// Convert HTTP URL to WebSocket URL
	wsURL := strings.Replace(suite.server.URL, "http://", "ws://", 1) + "/api/v1/admin/chat/ws?token=test-token"

	// Connect to WebSocket
	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial(wsURL, nil)
	require.NoError(t, err)
	defer conn.Close()

	// Send a chat message
	chatRequest := map[string]interface{}{
		"message":     "What is AWS EC2?",
		"client_name": "WebSocket Test Client",
		"context":     "WebSocket integration test",
	}

	err = conn.WriteJSON(chatRequest)
	require.NoError(t, err)

	// Read response
	var response map[string]interface{}
	err = conn.ReadJSON(&response)
	require.NoError(t, err)

	assert.True(t, response["success"].(bool))
	assert.NotNil(t, response["message"])
	assert.NotEmpty(t, response["session_id"])

	messageResp := response["message"].(map[string]interface{})
	assert.Equal(t, "assistant", messageResp["type"])
	assert.NotEmpty(t, messageResp["content"])
	assert.NotEmpty(t, messageResp["id"])
}

func TestChatAPI_WebSocket_MultipleMessages_Integration(t *testing.T) {
	suite := setupIntegrationTestSuite(t)
	defer suite.tearDown()

	wsURL := strings.Replace(suite.server.URL, "http://", "ws://", 1) + "/api/v1/admin/chat/ws?token=test-token"

	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial(wsURL, nil)
	require.NoError(t, err)
	defer conn.Close()

	messages := []string{
		"What is AWS EC2?",
		"How much does it cost?",
		"What are the security best practices?",
	}

	var sessionID string

	for i, message := range messages {
		chatRequest := map[string]interface{}{
			"message":     message,
			"client_name": "Multi-Message Test Client",
			"context":     "Multiple message test",
		}

		if sessionID != "" {
			chatRequest["session_id"] = sessionID
		}

		err = conn.WriteJSON(chatRequest)
		require.NoError(t, err)

		var response map[string]interface{}
		err = conn.ReadJSON(&response)
		require.NoError(t, err)

		assert.True(t, response["success"].(bool), "Message %d failed", i+1)
		assert.NotNil(t, response["message"])

		if sessionID == "" {
			sessionID = response["session_id"].(string)
		} else {
			assert.Equal(t, sessionID, response["session_id"], "Session ID should remain consistent")
		}

		messageResp := response["message"].(map[string]interface{})
		assert.Equal(t, "assistant", messageResp["type"])
		assert.NotEmpty(t, messageResp["content"])
	}
}

// Database integration tests (if using real database)
func TestChatAPI_DatabaseIntegration(t *testing.T) {
	// Skip if no database available
	if testing.Short() {
		t.Skip("Skipping database integration test in short mode")
	}

	suite := setupIntegrationTestSuite(t)
	defer suite.tearDown()

	ctx := context.Background()

	// Test session persistence
	session := &domain.ChatSession{
		ID:           "db-test-session",
		UserID:       "db-test-user",
		ClientName:   "Database Test Client",
		Context:      "Database integration test",
		Status:       domain.SessionStatusActive,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		LastActivity: time.Now(),
		ExpiresAt:    timePtr(time.Now().Add(24 * time.Hour)),
		Metadata: map[string]interface{}{
			"test": "database_integration",
		},
	}

	// Create session
	err := suite.sessionSvc.CreateSession(ctx, session)
	require.NoError(t, err)

	// Retrieve session
	retrievedSession, err := suite.sessionSvc.GetSession(ctx, session.ID)
	require.NoError(t, err)
	assert.Equal(t, session.ID, retrievedSession.ID)
	assert.Equal(t, session.ClientName, retrievedSession.ClientName)

	// Test message persistence
	message := &domain.ChatMessage{
		ID:        "db-test-message",
		SessionID: session.ID,
		Type:      domain.MessageTypeUser,
		Content:   "Database test message",
		Status:    domain.MessageStatusSent,
		CreatedAt: time.Now(),
		Metadata: map[string]interface{}{
			"test": "database_integration",
		},
	}

	err = suite.messageRepo.Create(ctx, message)
	require.NoError(t, err)

	// Retrieve message
	retrievedMessage, err := suite.messageRepo.GetByID(ctx, message.ID)
	require.NoError(t, err)
	assert.Equal(t, message.ID, retrievedMessage.ID)
	assert.Equal(t, message.Content, retrievedMessage.Content)

	// Test session history
	messages, err := suite.messageRepo.GetBySessionID(ctx, session.ID, 10, 0)
	require.NoError(t, err)
	assert.Len(t, messages, 1)
	assert.Equal(t, message.ID, messages[0].ID)
}

// Performance integration tests
func TestChatAPI_Performance_Integration(t *testing.T) {
	suite := setupIntegrationTestSuite(t)
	defer suite.tearDown()

	// Test concurrent session creation
	t.Run("ConcurrentSessionCreation", func(t *testing.T) {
		concurrency := 10
		done := make(chan bool, concurrency)

		for i := 0; i < concurrency; i++ {
			go func(id int) {
				sessionData := map[string]interface{}{
					"client_name": fmt.Sprintf("Concurrent Client %d", id),
					"context":     fmt.Sprintf("Concurrent test %d", id),
				}

				req := suite.createAuthenticatedRequest("POST", "/api/v1/admin/chat/sessions", sessionData)
				req = req.WithContext(context.WithValue(req.Context(), "user_id", fmt.Sprintf("user-%d", id)))

				client := &http.Client{}
				resp, err := client.Do(req)
				assert.NoError(t, err)
				if resp != nil {
					resp.Body.Close()
					assert.Equal(t, http.StatusCreated, resp.StatusCode)
				}

				done <- true
			}(i)
		}

		// Wait for all goroutines to complete
		for i := 0; i < concurrency; i++ {
			<-done
		}
	})

	// Test message throughput
	t.Run("MessageThroughput", func(t *testing.T) {
		ctx := context.Background()

		// Create a session
		session := &domain.ChatSession{
			ID:           "perf-test-session",
			UserID:       "perf-test-user",
			ClientName:   "Performance Test Client",
			Status:       domain.SessionStatusActive,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
			LastActivity: time.Now(),
			ExpiresAt:    timePtr(time.Now().Add(24 * time.Hour)),
		}

		err := suite.sessionSvc.CreateSession(ctx, session)
		require.NoError(t, err)

		// Send multiple messages rapidly
		messageCount := 20
		start := time.Now()

		for i := 0; i < messageCount; i++ {
			messageData := map[string]interface{}{
				"message":    fmt.Sprintf("Performance test message %d", i),
				"session_id": session.ID,
			}

			req := suite.createAuthenticatedRequest("POST", fmt.Sprintf("/api/v1/admin/chat/sessions/%s/messages", session.ID), messageData)
			req = req.WithContext(context.WithValue(req.Context(), "user_id", "perf-test-user"))

			client := &http.Client{}
			resp, err := client.Do(req)
			assert.NoError(t, err)
			if resp != nil {
				resp.Body.Close()
				assert.Equal(t, http.StatusOK, resp.StatusCode)
			}
		}

		duration := time.Since(start)
		messagesPerSecond := float64(messageCount) / duration.Seconds()

		t.Logf("Processed %d messages in %v (%.2f messages/second)", messageCount, duration, messagesPerSecond)

		// Assert reasonable performance (adjust threshold as needed)
		assert.Greater(t, messagesPerSecond, 5.0, "Message processing should be at least 5 messages/second")
	})
}

// Error handling integration tests
func TestChatAPI_ErrorHandling_Integration(t *testing.T) {
	suite := setupIntegrationTestSuite(t)
	defer suite.tearDown()

	t.Run("InvalidSessionID", func(t *testing.T) {
		req := suite.createAuthenticatedRequest("GET", "/api/v1/admin/chat/sessions/invalid-session-id", nil)
		req = req.WithContext(context.WithValue(req.Context(), "user_id", "test-user"))

		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)
		assert.False(t, response["success"].(bool))
	})

	t.Run("MalformedJSON", func(t *testing.T) {
		req, _ := http.NewRequest("POST", suite.server.URL+"/api/v1/admin/chat/sessions", strings.NewReader("invalid json"))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer test-token")
		req = req.WithContext(context.WithValue(req.Context(), "user_id", "test-user"))

		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("UnauthorizedAccess", func(t *testing.T) {
		req, _ := http.NewRequest("GET", suite.server.URL+"/api/v1/admin/chat/sessions", nil)
		// No authorization header

		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

// Helper function
func timePtr(t time.Time) *time.Time {
	return &t
}

// Run all integration tests
func TestMain(m *testing.M) {
	// Setup test environment
	gin.SetMode(gin.TestMode)

	// Run tests
	m.Run()
}
