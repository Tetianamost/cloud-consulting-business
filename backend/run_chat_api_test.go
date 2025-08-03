package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/cloud-consulting/backend/internal/config"
	"github.com/cloud-consulting/backend/internal/handlers"
	"github.com/cloud-consulting/backend/internal/services"
	"github.com/cloud-consulting/backend/internal/storage"
)

func main() {
	fmt.Println("🚀 Starting Chat REST API Integration Test...")

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
	fmt.Println("📝 Getting authentication token...")
	token := getAuthToken(router)
	if token == "" {
		fmt.Println("❌ Failed to get authentication token")
		return
	}
	fmt.Println("✅ Authentication token obtained")

	// Test session creation
	fmt.Println("📝 Testing session creation...")
	sessionID := testCreateChatSession(router, token)
	if sessionID == "" {
		fmt.Println("❌ Failed to create session")
		return
	}
	fmt.Printf("✅ Session created successfully with ID: %s\n", sessionID)

	// Test listing sessions
	fmt.Println("📝 Testing session listing...")
	if !testListChatSessions(router, token, sessionID) {
		fmt.Println("❌ Failed to list sessions")
		return
	}
	fmt.Println("✅ Sessions listed successfully")

	// Test getting specific session
	fmt.Println("📝 Testing get specific session...")
	if !testGetChatSession(router, token, sessionID) {
		fmt.Println("❌ Failed to get session")
		return
	}
	fmt.Println("✅ Session retrieved successfully")

	// Test updating session
	fmt.Println("📝 Testing session update...")
	if !testUpdateChatSession(router, token, sessionID) {
		fmt.Println("❌ Failed to update session")
		return
	}
	fmt.Println("✅ Session updated successfully")

	// Test getting session history
	fmt.Println("📝 Testing session history...")
	if !testGetChatSessionHistory(router, token, sessionID) {
		fmt.Println("❌ Failed to get session history")
		return
	}
	fmt.Println("✅ Session history retrieved successfully")

	// Test deleting session
	fmt.Println("📝 Testing session deletion...")
	if !testDeleteChatSession(router, token, sessionID) {
		fmt.Println("❌ Failed to delete session")
		return
	}
	fmt.Println("✅ Session deleted successfully")

	// Test error cases
	fmt.Println("📝 Testing error cases...")
	if !testErrorCases(router, token) {
		fmt.Println("❌ Error cases test failed")
		return
	}
	fmt.Println("✅ Error cases handled correctly")

	fmt.Println("\n🎉 All Chat REST API Integration Tests Passed!")
	fmt.Println("\n📋 Test Summary:")
	fmt.Println("✅ POST /api/v1/admin/chat/sessions - Session creation")
	fmt.Println("✅ GET /api/v1/admin/chat/sessions - Session listing")
	fmt.Println("✅ GET /api/v1/admin/chat/sessions/{id} - Session retrieval")
	fmt.Println("✅ PUT /api/v1/admin/chat/sessions/{id} - Session update")
	fmt.Println("✅ DELETE /api/v1/admin/chat/sessions/{id} - Session deletion")
	fmt.Println("✅ GET /api/v1/admin/chat/sessions/{id}/history - Session history")
	fmt.Println("✅ Error handling and authentication")
}

// getAuthToken gets a valid JWT token for testing
func getAuthToken(router *gin.Engine) string {
	loginRequest := map[string]string{
		"username": "admin",
		"password": "admin123",
	}

	body, _ := json.Marshal(loginRequest)
	req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		fmt.Printf("❌ Login failed with status: %d\n", w.Code)
		return ""
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		fmt.Printf("❌ Failed to parse login response: %v\n", err)
		return ""
	}

	data, ok := response["data"].(map[string]interface{})
	if !ok {
		fmt.Println("❌ Invalid login response format")
		return ""
	}

	token, ok := data["token"].(string)
	if !ok || token == "" {
		fmt.Println("❌ No token in login response")
		return ""
	}

	return token
}

// testCreateChatSession tests the POST /api/v1/admin/chat/sessions endpoint
func testCreateChatSession(router *gin.Engine, token string) string {
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

	if w.Code != http.StatusCreated {
		fmt.Printf("❌ Create session failed with status: %d, body: %s\n", w.Code, w.Body.String())
		return ""
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		fmt.Printf("❌ Failed to parse create response: %v\n", err)
		return ""
	}

	success, ok := response["success"].(bool)
	if !ok || !success {
		fmt.Printf("❌ Create session response indicates failure: %v\n", response)
		return ""
	}

	data, ok := response["data"].(map[string]interface{})
	if !ok {
		fmt.Printf("❌ Invalid create response format: %v\n", response)
		return ""
	}

	sessionID, ok := data["id"].(string)
	if !ok || sessionID == "" {
		fmt.Printf("❌ No session ID in create response: %v\n", data)
		return ""
	}

	return sessionID
}

// testListChatSessions tests the GET /api/v1/admin/chat/sessions endpoint
func testListChatSessions(router *gin.Engine, token string, expectedSessionID string) bool {
	req := httptest.NewRequest("GET", "/api/v1/admin/chat/sessions", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		fmt.Printf("❌ List sessions failed with status: %d, body: %s\n", w.Code, w.Body.String())
		return false
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		fmt.Printf("❌ Failed to parse list response: %v\n", err)
		return false
	}

	success, ok := response["success"].(bool)
	if !ok || !success {
		fmt.Printf("❌ List sessions response indicates failure: %v\n", response)
		return false
	}

	data, ok := response["data"].(map[string]interface{})
	if !ok {
		fmt.Printf("❌ Invalid list response format: %v\n", response)
		return false
	}

	sessions, ok := data["sessions"].([]interface{})
	if !ok {
		fmt.Printf("❌ No sessions array in list response: %v\n", data)
		return false
	}

	if len(sessions) != 1 {
		fmt.Printf("❌ Expected 1 session, got %d\n", len(sessions))
		return false
	}

	session := sessions[0].(map[string]interface{})
	if session["id"] != expectedSessionID {
		fmt.Printf("❌ Session ID mismatch: expected %s, got %v\n", expectedSessionID, session["id"])
		return false
	}

	return true
}

// testGetChatSession tests the GET /api/v1/admin/chat/sessions/:id endpoint
func testGetChatSession(router *gin.Engine, token string, sessionID string) bool {
	req := httptest.NewRequest("GET", "/api/v1/admin/chat/sessions/"+sessionID, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		fmt.Printf("❌ Get session failed with status: %d, body: %s\n", w.Code, w.Body.String())
		return false
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		fmt.Printf("❌ Failed to parse get response: %v\n", err)
		return false
	}

	success, ok := response["success"].(bool)
	if !ok || !success {
		fmt.Printf("❌ Get session response indicates failure: %v\n", response)
		return false
	}

	data, ok := response["data"].(map[string]interface{})
	if !ok {
		fmt.Printf("❌ Invalid get response format: %v\n", response)
		return false
	}

	if data["id"] != sessionID {
		fmt.Printf("❌ Session ID mismatch: expected %s, got %v\n", sessionID, data["id"])
		return false
	}

	return true
}

// testUpdateChatSession tests the PUT /api/v1/admin/chat/sessions/:id endpoint
func testUpdateChatSession(router *gin.Engine, token string, sessionID string) bool {
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

	if w.Code != http.StatusOK {
		fmt.Printf("❌ Update session failed with status: %d, body: %s\n", w.Code, w.Body.String())
		return false
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		fmt.Printf("❌ Failed to parse update response: %v\n", err)
		return false
	}

	success, ok := response["success"].(bool)
	if !ok || !success {
		fmt.Printf("❌ Update session response indicates failure: %v\n", response)
		return false
	}

	return true
}

// testGetChatSessionHistory tests the GET /api/v1/admin/chat/sessions/:id/history endpoint
func testGetChatSessionHistory(router *gin.Engine, token string, sessionID string) bool {
	req := httptest.NewRequest("GET", "/api/v1/admin/chat/sessions/"+sessionID+"/history", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		fmt.Printf("❌ Get session history failed with status: %d, body: %s\n", w.Code, w.Body.String())
		return false
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		fmt.Printf("❌ Failed to parse history response: %v\n", err)
		return false
	}

	success, ok := response["success"].(bool)
	if !ok || !success {
		fmt.Printf("❌ Get session history response indicates failure: %v\n", response)
		return false
	}

	data, ok := response["data"].(map[string]interface{})
	if !ok {
		fmt.Printf("❌ Invalid history response format: %v\n", response)
		return false
	}

	messages, ok := data["messages"].([]interface{})
	if !ok {
		fmt.Printf("❌ No messages array in history response: %v\n", data)
		return false
	}

	// Should be empty since we haven't sent any messages
	if len(messages) != 0 {
		fmt.Printf("❌ Expected 0 messages, got %d\n", len(messages))
		return false
	}

	return true
}

// testDeleteChatSession tests the DELETE /api/v1/admin/chat/sessions/:id endpoint
func testDeleteChatSession(router *gin.Engine, token string, sessionID string) bool {
	req := httptest.NewRequest("DELETE", "/api/v1/admin/chat/sessions/"+sessionID, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		fmt.Printf("❌ Delete session failed with status: %d, body: %s\n", w.Code, w.Body.String())
		return false
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		fmt.Printf("❌ Failed to parse delete response: %v\n", err)
		return false
	}

	success, ok := response["success"].(bool)
	if !ok || !success {
		fmt.Printf("❌ Delete session response indicates failure: %v\n", response)
		return false
	}

	// Verify session is actually deleted
	req = httptest.NewRequest("GET", "/api/v1/admin/chat/sessions/"+sessionID, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		fmt.Printf("❌ Session should be deleted but still accessible with status: %d\n", w.Code)
		return false
	}

	return true
}

// testErrorCases tests various error scenarios
func testErrorCases(router *gin.Engine, token string) bool {
	// Test unauthorized access (no token)
	req := httptest.NewRequest("GET", "/api/v1/admin/chat/sessions", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		fmt.Printf("❌ Expected unauthorized status, got: %d\n", w.Code)
		return false
	}

	// Test invalid session ID
	req = httptest.NewRequest("GET", "/api/v1/admin/chat/sessions/invalid-id", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		fmt.Printf("❌ Expected not found status for invalid ID, got: %d\n", w.Code)
		return false
	}

	// Test invalid JSON in create request
	req = httptest.NewRequest("POST", "/api/v1/admin/chat/sessions", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		fmt.Printf("❌ Expected bad request status for invalid JSON, got: %d\n", w.Code)
		return false
	}

	return true
}
