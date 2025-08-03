package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cloud-consulting/backend/internal/config"
	"github.com/cloud-consulting/backend/internal/domain"
	"github.com/cloud-consulting/backend/internal/handlers"
	"github.com/cloud-consulting/backend/internal/services"
	"github.com/cloud-consulting/backend/internal/storage"
)

// TestChatRESTAPIIntegration tests all REST API endpoints for chat management
func TestChatRESTAPIIntegration(t *testing.T) {
	// Set up test environment
	gin.SetMode(gin.TestMode)
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel) // Reduce noise in tests

	// Initialize services
	cfg := &config.BedrockConfig{
		Region:  "us-east-1",
		APIKey:  "test-key",
		ModelID: "test-model",
		BaseURL: "https://test-bedrock.amazonaws.com",
		Timeout: 30,
	}
	bedrockService := services.NewBedrockService(cfg)
	knowledgeBase := services.NewInMemoryKnowledgeBase()

	// Initialize chat repositories and services
	chatSessionRepository := storage.NewInMemoryChatSessionRepository(logger)
	chatMessageRepository := storage.NewInMemoryChatMessageRepository(logger)
	sessionService := services.NewSessionService(chatSessionRepository, logger)
	chatService := services.NewChatService(chatMessageRepository, sessionService, bedrockService, logger)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler("test-jwt-secret")
	chatHandler := handlers.NewChatHandler(logger, bedrockService, knowledgeBase, sessionService, chatService, authHandler, "test-jwt-secret")

	// Set up router
	router := gin.New()
	router.Use(gin.Recovery())

	// Add auth and chat routes
	v1 := router.Group("/api/v1")
	auth := v1.Group("/auth")
	auth.POST("/login", authHandler.Login)

	admin := v1.Group("/admin", authHandler.AuthMiddleware())
	admin.POST("/chat/sessions", chatHandler.CreateChatSession)
	admin.GET("/chat/sessions", chatHandler.ListChatSessions)
	admin.GET("/chat/sessions/:id", chatHandler.GetChatSessionByID)
	admin.PUT("/chat/sessions/:id", chatHandler.UpdateChatSession)
	admin.DELETE("/chat/sessions/:id", chatHandler.DeleteChatSession)
	admin.GET("/chat/sessions/:id/history", chatHandler.GetChatSessionHistory)

	// Get auth token
	token := getAuthToken(t, router)

	// Test session creation
	sessionID := testCreateChatSession(t, router, token)

	// Test listing sessions
	testListChatSessions(t, router, token, sessionID)

	// Test getting specific session
	testGetChatSession(t, router, token, sessionID)

	// Test updating session
	testUpdateChatSession(t, router, token, sessionID)

	// Test getting session history
	testGetChatSessionHistory(t, router, token, sessionID)

	// Test deleting session
	testDeleteChatSession(t, router, token, sessionID)

	// Test error cases
	testErrorCases(t, router, token)

	fmt.Println("✅ All chat REST API integration tests passed!")
}

// getAuthToken gets a valid JWT token for testing
func getAuthToken(t *testing.T, router *gin.Engine) string {
	loginRequest := map[string]string{
		"username": "admin",
		"password": "admin123",
	}

	body, _ := json.Marshal(loginRequest)
	req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	data, ok := response["data"].(map[string]interface{})
	require.True(t, ok)

	token, ok := data["token"].(string)
	require.True(t, ok)
	require.NotEmpty(t, token)

	return token
}

// testCreateChatSession tests the POST /api/v1/admin/chat/sessions endpoint
func testCreateChatSession(t *testing.T, router *gin.Engine, token string) string {
	fmt.Println("Testing chat session creation...")

	createRequest := map[string]interface{}{
		"client_name": "Test Client",
		"context":     "Test meeting context",
		"metadata": map[string]interface{}{
			"test_key": "test_value",
		},
	}

	body, _ := json.Marshal(createRequest)
	req := httptest.NewRequest("POST", "/api/v1/admin/chat/sessions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.True(t, response["success"].(bool))

	data, ok := response["data"].(map[string]interface{})
	require.True(t, ok)

	sessionID, ok := data["id"].(string)
	require.True(t, ok)
	require.NotEmpty(t, sessionID)

	assert.Equal(t, "Test Client", data["client_name"])
	assert.Equal(t, "Test meeting context", data["context"])
	assert.Equal(t, string(domain.SessionStatusActive), data["status"])

	fmt.Printf("✅ Session created successfully with ID: %s\n", sessionID)
	return sessionID
}

// testListChatSessions tests the GET /api/v1/admin/chat/sessions endpoint
func testListChatSessions(t *testing.T, router *gin.Engine, token string, expectedSessionID string) {
	fmt.Println("Testing chat sessions listing...")

	req := httptest.NewRequest("GET", "/api/v1/admin/chat/sessions", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.True(t, response["success"].(bool))

	data, ok := response["data"].(map[string]interface{})
	require.True(t, ok)

	sessions, ok := data["sessions"].([]interface{})
	require.True(t, ok)
	assert.Len(t, sessions, 1)

	session := sessions[0].(map[string]interface{})
	assert.Equal(t, expectedSessionID, session["id"])

	fmt.Println("✅ Sessions listed successfully")
}

// testGetChatSession tests the GET /api/v1/admin/chat/sessions/:id endpoint
func testGetChatSession(t *testing.T, router *gin.Engine, token string, sessionID string) {
	fmt.Println("Testing get specific chat session...")

	req := httptest.NewRequest("GET", "/api/v1/admin/chat/sessions/"+sessionID, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.True(t, response["success"].(bool))

	data, ok := response["data"].(map[string]interface{})
	require.True(t, ok)

	assert.Equal(t, sessionID, data["id"])
	assert.Equal(t, "Test Client", data["client_name"])

	fmt.Println("✅ Session retrieved successfully")
}

// testUpdateChatSession tests the PUT /api/v1/admin/chat/sessions/:id endpoint
func testUpdateChatSession(t *testing.T, router *gin.Engine, token string, sessionID string) {
	fmt.Println("Testing chat session update...")

	updateRequest := map[string]interface{}{
		"client_name":     "Updated Client",
		"context":         "Updated context",
		"meeting_type":    "discovery",
		"service_types":   []string{"migration", "optimization"},
		"cloud_providers": []string{"aws", "azure"},
		"custom_fields": map[string]string{
			"priority": "high",
		},
	}

	body, _ := json.Marshal(updateRequest)
	req := httptest.NewRequest("PUT", "/api/v1/admin/chat/sessions/"+sessionID, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.True(t, response["success"].(bool))

	data, ok := response["data"].(map[string]interface{})
	require.True(t, ok)

	assert.Equal(t, sessionID, data["id"])
	assert.Equal(t, "Updated Client", data["client_name"])
	assert.Equal(t, "Updated context", data["context"])

	fmt.Println("✅ Session updated successfully")
}

// testGetChatSessionHistory tests the GET /api/v1/admin/chat/sessions/:id/history endpoint
func testGetChatSessionHistory(t *testing.T, router *gin.Engine, token string, sessionID string) {
	fmt.Println("Testing get chat session history...")

	req := httptest.NewRequest("GET", "/api/v1/admin/chat/sessions/"+sessionID+"/history", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.True(t, response["success"].(bool))

	data, ok := response["data"].(map[string]interface{})
	require.True(t, ok)

	messages, ok := data["messages"].([]interface{})
	require.True(t, ok)
	// Should be empty since we haven't sent any messages
	assert.Len(t, messages, 0)

	fmt.Println("✅ Session history retrieved successfully")
}

// testDeleteChatSession tests the DELETE /api/v1/admin/chat/sessions/:id endpoint
func testDeleteChatSession(t *testing.T, router *gin.Engine, token string, sessionID string) {
	fmt.Println("Testing chat session deletion...")

	req := httptest.NewRequest("DELETE", "/api/v1/admin/chat/sessions/"+sessionID, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.True(t, response["success"].(bool))
	assert.Equal(t, "Session deleted successfully", response["message"])

	// Verify session is actually deleted
	req = httptest.NewRequest("GET", "/api/v1/admin/chat/sessions/"+sessionID, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	fmt.Println("✅ Session deleted successfully")
}

// testErrorCases tests various error scenarios
func testErrorCases(t *testing.T, router *gin.Engine, token string) {
	fmt.Println("Testing error cases...")

	// Test unauthorized access (no token)
	req := httptest.NewRequest("GET", "/api/v1/admin/chat/sessions", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// Test invalid session ID
	req = httptest.NewRequest("GET", "/api/v1/admin/chat/sessions/invalid-id", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)

	// Test invalid JSON in create request
	req = httptest.NewRequest("POST", "/api/v1/admin/chat/sessions", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	fmt.Println("✅ Error cases handled correctly")
}

// Helper function to read response body
func readResponseBody(t *testing.T, resp *http.Response) []byte {
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	defer resp.Body.Close()
	return body
}

// This file contains integration tests for the chat REST API endpoints
// Run with: go test -run TestChatRESTAPIIntegration
