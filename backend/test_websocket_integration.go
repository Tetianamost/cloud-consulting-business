package main

import (
	"fmt"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cloud-consulting/backend/internal/handlers"
	"github.com/cloud-consulting/backend/internal/services"
	"github.com/cloud-consulting/backend/internal/storage"
)

// WebSocket message types for testing
type WSTestMessage struct {
	Type      string                 `json:"type"`
	SessionID string                 `json:"session_id,omitempty"`
	MessageID string                 `json:"message_id,omitempty"`
	Content   string                 `json:"content,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
}

type WSTestResponse struct {
	Success   bool                   `json:"success"`
	SessionID string                 `json:"session_id,omitempty"`
	Message   WSTestMessageResponse  `json:"message,omitempty"`
	Error     string                 `json:"error,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

type WSTestMessageResponse struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
	SessionID string    `json:"session_id"`
}

// WebSocket test client
type WSTestClient struct {
	conn      *websocket.Conn
	responses chan WSTestResponse
	errors    chan error
	done      chan bool
	mutex     sync.RWMutex
}

func NewWSTestClient(serverURL string) (*WSTestClient, error) {
	wsURL := strings.Replace(serverURL, "http://", "ws://", 1) + "/api/v1/admin/chat/ws?token=test-token"

	dialer := websocket.Dialer{
		HandshakeTimeout: 5 * time.Second,
	}

	conn, _, err := dialer.Dial(wsURL, nil)
	if err != nil {
		return nil, err
	}

	client := &WSTestClient{
		conn:      conn,
		responses: make(chan WSTestResponse, 100),
		errors:    make(chan error, 10),
		done:      make(chan bool, 1),
	}

	// Start reading messages
	go client.readMessages()

	return client, nil
}

func (c *WSTestClient) readMessages() {
	defer close(c.responses)
	defer close(c.errors)

	for {
		select {
		case <-c.done:
			return
		default:
			var response WSTestResponse
			err := c.conn.ReadJSON(&response)
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					c.errors <- err
				}
				return
			}
			c.responses <- response
		}
	}
}

func (c *WSTestClient) SendMessage(message interface{}) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.conn.WriteJSON(message)
}

func (c *WSTestClient) WaitForResponse(timeout time.Duration) (*WSTestResponse, error) {
	select {
	case response := <-c.responses:
		return &response, nil
	case err := <-c.errors:
		return nil, err
	case <-time.After(timeout):
		return nil, fmt.Errorf("timeout waiting for response")
	}
}

func (c *WSTestClient) Close() error {
	c.done <- true
	return c.conn.Close()
}

// WebSocket integration test suite
type WSIntegrationTestSuite struct {
	server      *httptest.Server
	chatHandler *handlers.ChatHandler
	logger      *logrus.Logger
}

func setupWSIntegrationTestSuite(t *testing.T) *WSIntegrationTestSuite {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	// Create in-memory repositories
	sessionRepo := storage.NewMemoryChatSessionRepository()
	messageRepo := storage.NewMemoryChatMessageRepository()

	// Create services
	sessionSvc := services.NewSessionService(sessionRepo, logger)
	chatSvc := services.NewChatService(messageRepo, sessionSvc, nil, logger)

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
		nil,
		testJWTSecret,
	)

	// Setup Gin router
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Add WebSocket route
	api := router.Group("/api/v1")
	{
		admin := api.Group("/admin")
		{
			chat := admin.Group("/chat")
			{
				chat.GET("/ws", chatHandler.HandleWebSocket)
			}
		}
	}

	server := httptest.NewServer(router)

	return &WSIntegrationTestSuite{
		server:      server,
		chatHandler: chatHandler,
		logger:      logger,
	}
}

func (suite *WSIntegrationTestSuite) tearDown() {
	suite.server.Close()
}

// Test basic WebSocket connection
func TestWebSocket_BasicConnection(t *testing.T) {
	suite := setupWSIntegrationTestSuite(t)
	defer suite.tearDown()

	client, err := NewWSTestClient(suite.server.URL)
	require.NoError(t, err)
	defer client.Close()

	// Connection should be established successfully
	assert.NotNil(t, client.conn)
}

// Test single message exchange
func TestWebSocket_SingleMessage(t *testing.T) {
	suite := setupWSIntegrationTestSuite(t)
	defer suite.tearDown()

	client, err := NewWSTestClient(suite.server.URL)
	require.NoError(t, err)
	defer client.Close()

	// Send a message
	message := map[string]interface{}{
		"message":     "What is AWS Lambda?",
		"client_name": "WebSocket Test Client",
		"context":     "Single message test",
	}

	err = client.SendMessage(message)
	require.NoError(t, err)

	// Wait for response
	response, err := client.WaitForResponse(5 * time.Second)
	require.NoError(t, err)

	assert.True(t, response.Success)
	assert.NotEmpty(t, response.SessionID)
	assert.Equal(t, "assistant", response.Message.Type)
	assert.NotEmpty(t, response.Message.Content)
	assert.NotEmpty(t, response.Message.ID)
}

// Test conversation flow
func TestWebSocket_ConversationFlow(t *testing.T) {
	suite := setupWSIntegrationTestSuite(t)
	defer suite.tearDown()

	client, err := NewWSTestClient(suite.server.URL)
	require.NoError(t, err)
	defer client.Close()

	var sessionID string

	// First message
	message1 := map[string]interface{}{
		"message":     "What is AWS EC2?",
		"client_name": "Conversation Test Client",
		"context":     "Learning about AWS services",
	}

	err = client.SendMessage(message1)
	require.NoError(t, err)

	response1, err := client.WaitForResponse(5 * time.Second)
	require.NoError(t, err)

	assert.True(t, response1.Success)
	sessionID = response1.SessionID
	assert.NotEmpty(t, sessionID)

	// Second message in same session
	message2 := map[string]interface{}{
		"message":    "How much does it cost?",
		"session_id": sessionID,
	}

	err = client.SendMessage(message2)
	require.NoError(t, err)

	response2, err := client.WaitForResponse(5 * time.Second)
	require.NoError(t, err)

	assert.True(t, response2.Success)
	assert.Equal(t, sessionID, response2.SessionID)
	assert.Contains(t, strings.ToLower(response2.Message.Content), "cost")

	// Third message in same session
	message3 := map[string]interface{}{
		"message":    "What are the security best practices?",
		"session_id": sessionID,
	}

	err = client.SendMessage(message3)
	require.NoError(t, err)

	response3, err := client.WaitForResponse(5 * time.Second)
	require.NoError(t, err)

	assert.True(t, response3.Success)
	assert.Equal(t, sessionID, response3.SessionID)
	assert.Contains(t, strings.ToLower(response3.Message.Content), "security")
}

// Test multiple concurrent connections
func TestWebSocket_ConcurrentConnections(t *testing.T) {
	suite := setupWSIntegrationTestSuite(t)
	defer suite.tearDown()

	concurrency := 5
	var wg sync.WaitGroup
	results := make(chan bool, concurrency)

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(clientID int) {
			defer wg.Done()

			client, err := NewWSTestClient(suite.server.URL)
			if err != nil {
				results <- false
				return
			}
			defer client.Close()

			// Send a message
			message := map[string]interface{}{
				"message":     fmt.Sprintf("Concurrent test message from client %d", clientID),
				"client_name": fmt.Sprintf("Concurrent Client %d", clientID),
				"context":     "Concurrent connection test",
			}

			err = client.SendMessage(message)
			if err != nil {
				results <- false
				return
			}

			// Wait for response
			response, err := client.WaitForResponse(10 * time.Second)
			if err != nil {
				results <- false
				return
			}

			success := response.Success && response.Message.Type == "assistant" && response.Message.Content != ""
			results <- success
		}(i)
	}

	wg.Wait()
	close(results)

	// Check all connections succeeded
	successCount := 0
	for result := range results {
		if result {
			successCount++
		}
	}

	assert.Equal(t, concurrency, successCount, "All concurrent connections should succeed")
}

// Test WebSocket message types
func TestWebSocket_MessageTypes(t *testing.T) {
	suite := setupWSIntegrationTestSuite(t)
	defer suite.tearDown()

	client, err := NewWSTestClient(suite.server.URL)
	require.NoError(t, err)
	defer client.Close()

	testCases := []struct {
		name     string
		message  map[string]interface{}
		expected string
	}{
		{
			name: "Regular Message",
			message: map[string]interface{}{
				"type":    "message",
				"content": "What is AWS S3?",
			},
			expected: "assistant",
		},
		{
			name: "Typing Indicator",
			message: map[string]interface{}{
				"type":      "typing",
				"is_typing": true,
			},
			expected: "ack", // Should get acknowledgment
		},
		{
			name: "Heartbeat",
			message: map[string]interface{}{
				"type": "heartbeat",
			},
			expected: "ack",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err = client.SendMessage(tc.message)
			require.NoError(t, err)

			response, err := client.WaitForResponse(5 * time.Second)
			require.NoError(t, err)

			if tc.expected == "assistant" {
				assert.Equal(t, tc.expected, response.Message.Type)
			} else {
				// For non-message types, just verify we get a response
				assert.True(t, response.Success || response.Error == "")
			}
		})
	}
}

// Test WebSocket error handling
func TestWebSocket_ErrorHandling(t *testing.T) {
	suite := setupWSIntegrationTestSuite(t)
	defer suite.tearDown()

	client, err := NewWSTestClient(suite.server.URL)
	require.NoError(t, err)
	defer client.Close()

	testCases := []struct {
		name    string
		message map[string]interface{}
	}{
		{
			name: "Empty Message",
			message: map[string]interface{}{
				"message": "",
			},
		},
		{
			name: "Invalid Session ID",
			message: map[string]interface{}{
				"message":    "Test message",
				"session_id": "invalid-session-id",
			},
		},
		{
			name: "Malformed Request",
			message: map[string]interface{}{
				"invalid_field": "invalid_value",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err = client.SendMessage(tc.message)
			require.NoError(t, err)

			response, err := client.WaitForResponse(5 * time.Second)
			require.NoError(t, err)

			// Should get an error response or handle gracefully
			if !response.Success {
				assert.NotEmpty(t, response.Error)
			}
		})
	}
}

// Test WebSocket connection limits and rate limiting
func TestWebSocket_RateLimiting(t *testing.T) {
	suite := setupWSIntegrationTestSuite(t)
	defer suite.tearDown()

	client, err := NewWSTestClient(suite.server.URL)
	require.NoError(t, err)
	defer client.Close()

	// Send messages rapidly to trigger rate limiting
	messageCount := 100
	rateLimitHit := false

	for i := 0; i < messageCount; i++ {
		message := map[string]interface{}{
			"message": fmt.Sprintf("Rate limit test message %d", i),
		}

		err = client.SendMessage(message)
		if err != nil {
			break
		}

		// Try to get response quickly
		response, err := client.WaitForResponse(100 * time.Millisecond)
		if err != nil {
			continue // Timeout is expected under rate limiting
		}

		if !response.Success && strings.Contains(response.Error, "rate limit") {
			rateLimitHit = true
			break
		}
	}

	// Rate limiting should eventually kick in
	assert.True(t, rateLimitHit, "Rate limiting should be triggered with rapid messages")
}

// Test WebSocket connection persistence
func TestWebSocket_ConnectionPersistence(t *testing.T) {
	suite := setupWSIntegrationTestSuite(t)
	defer suite.tearDown()

	client, err := NewWSTestClient(suite.server.URL)
	require.NoError(t, err)
	defer client.Close()

	// Send messages over time to test connection persistence
	messageInterval := 500 * time.Millisecond
	messageCount := 5

	for i := 0; i < messageCount; i++ {
		message := map[string]interface{}{
			"message": fmt.Sprintf("Persistence test message %d", i),
		}

		err = client.SendMessage(message)
		require.NoError(t, err)

		response, err := client.WaitForResponse(5 * time.Second)
		require.NoError(t, err)
		assert.True(t, response.Success)

		if i < messageCount-1 {
			time.Sleep(messageInterval)
		}
	}
}

// Test WebSocket message ordering
func TestWebSocket_MessageOrdering(t *testing.T) {
	suite := setupWSIntegrationTestSuite(t)
	defer suite.tearDown()

	client, err := NewWSTestClient(suite.server.URL)
	require.NoError(t, err)
	defer client.Close()

	messageCount := 10
	responses := make([]*WSTestResponse, 0, messageCount)

	// Send multiple messages quickly
	for i := 0; i < messageCount; i++ {
		message := map[string]interface{}{
			"message": fmt.Sprintf("Order test message %d", i),
		}

		err = client.SendMessage(message)
		require.NoError(t, err)
	}

	// Collect all responses
	for i := 0; i < messageCount; i++ {
		response, err := client.WaitForResponse(10 * time.Second)
		require.NoError(t, err)
		responses = append(responses, response)
	}

	// Verify we got all responses
	assert.Len(t, responses, messageCount)

	// All responses should be successful
	for i, response := range responses {
		assert.True(t, response.Success, "Response %d should be successful", i)
		assert.NotEmpty(t, response.Message.Content, "Response %d should have content", i)
	}
}

// Test WebSocket session management
func TestWebSocket_SessionManagement(t *testing.T) {
	suite := setupWSIntegrationTestSuite(t)
	defer suite.tearDown()

	// Test multiple clients with different sessions
	client1, err := NewWSTestClient(suite.server.URL)
	require.NoError(t, err)
	defer client1.Close()

	client2, err := NewWSTestClient(suite.server.URL)
	require.NoError(t, err)
	defer client2.Close()

	// Client 1 sends a message
	message1 := map[string]interface{}{
		"message":     "Client 1 message",
		"client_name": "Client 1",
	}

	err = client1.SendMessage(message1)
	require.NoError(t, err)

	response1, err := client1.WaitForResponse(5 * time.Second)
	require.NoError(t, err)
	assert.True(t, response1.Success)
	sessionID1 := response1.SessionID

	// Client 2 sends a message
	message2 := map[string]interface{}{
		"message":     "Client 2 message",
		"client_name": "Client 2",
	}

	err = client2.SendMessage(message2)
	require.NoError(t, err)

	response2, err := client2.WaitForResponse(5 * time.Second)
	require.NoError(t, err)
	assert.True(t, response2.Success)
	sessionID2 := response2.SessionID

	// Sessions should be different
	assert.NotEqual(t, sessionID1, sessionID2, "Different clients should have different sessions")

	// Client 1 continues in same session
	message1_2 := map[string]interface{}{
		"message":    "Client 1 follow-up",
		"session_id": sessionID1,
	}

	err = client1.SendMessage(message1_2)
	require.NoError(t, err)

	response1_2, err := client1.WaitForResponse(5 * time.Second)
	require.NoError(t, err)
	assert.True(t, response1_2.Success)
	assert.Equal(t, sessionID1, response1_2.SessionID, "Client 1 should continue in same session")
}

// Benchmark WebSocket performance
func BenchmarkWebSocket_MessageThroughput(b *testing.B) {
	suite := setupWSIntegrationTestSuite(&testing.T{})
	defer suite.tearDown()

	client, err := NewWSTestClient(suite.server.URL)
	if err != nil {
		b.Fatal(err)
	}
	defer client.Close()

	message := map[string]interface{}{
		"message": "Benchmark test message",
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			err := client.SendMessage(message)
			if err != nil {
				b.Error(err)
				continue
			}

			_, err = client.WaitForResponse(5 * time.Second)
			if err != nil {
				b.Error(err)
			}
		}
	})
}

// Test WebSocket cleanup on disconnect
func TestWebSocket_CleanupOnDisconnect(t *testing.T) {
	suite := setupWSIntegrationTestSuite(t)
	defer suite.tearDown()

	client, err := NewWSTestClient(suite.server.URL)
	require.NoError(t, err)

	// Send a message to establish session
	message := map[string]interface{}{
		"message": "Test message before disconnect",
	}

	err = client.SendMessage(message)
	require.NoError(t, err)

	response, err := client.WaitForResponse(5 * time.Second)
	require.NoError(t, err)
	assert.True(t, response.Success)

	// Close connection
	err = client.Close()
	require.NoError(t, err)

	// Connection should be cleaned up (this is implicit - we're testing that no panics occur)
	time.Sleep(100 * time.Millisecond)
}
