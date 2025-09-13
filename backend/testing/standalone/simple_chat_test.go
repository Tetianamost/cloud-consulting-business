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
	fmt.Println("🚀 Quick Chat API Test...")

	// Set up minimal test environment
	gin.SetMode(gin.TestMode)
	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel) // Minimal logging

	// Initialize services quickly
	cfg := &config.BedrockConfig{
		Region:  "us-east-1",
		APIKey:  "test-key",
		ModelID: "test-model",
		BaseURL: "https://test-bedrock.amazonaws.com",
		Timeout: 30,
	}
	bedrockService := services.NewBedrockService(cfg)
	knowledgeBase := services.NewInMemoryKnowledgeBase()

	// Initialize chat services
	chatSessionRepository := storage.NewInMemoryChatSessionRepository(logger)
	chatMessageRepository := storage.NewInMemoryChatMessageRepository(logger)
	sessionService := services.NewSessionService(chatSessionRepository, logger)
	chatService := services.NewChatService(chatMessageRepository, sessionService, bedrockService, logger)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler("test-jwt-secret")
	chatHandler := handlers.NewChatHandler(logger, bedrockService, knowledgeBase, sessionService, chatService, authHandler, "test-jwt-secret")

	// Set up minimal router
	router := gin.New()
	v1 := router.Group("/api/v1")

	// Auth route
	auth := v1.Group("/auth")
	auth.POST("/login", authHandler.Login)

	// Chat routes
	admin := v1.Group("/admin", authHandler.AuthMiddleware())
	admin.POST("/chat/sessions", chatHandler.CreateChatSession)
	admin.GET("/chat/sessions", chatHandler.ListChatSessions)

	fmt.Println("📝 Testing login...")

	// Test login
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
		fmt.Printf("❌ Login failed: %d - %s\n", w.Code, w.Body.String())
		return
	}

	var loginResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &loginResponse)
	data := loginResponse["data"].(map[string]interface{})
	token := data["token"].(string)
	fmt.Println("✅ Login successful")

	fmt.Println("📝 Testing session creation...")

	// Test session creation
	createRequest := map[string]interface{}{
		"client_name": "Test Client",
		"context":     "Test context",
	}
	body, _ = json.Marshal(createRequest)
	req = httptest.NewRequest("POST", "/api/v1/admin/chat/sessions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		fmt.Printf("❌ Session creation failed: %d - %s\n", w.Code, w.Body.String())
		return
	}
	fmt.Println("✅ Session creation successful")

	fmt.Println("📝 Testing session listing...")

	// Test session listing
	req = httptest.NewRequest("GET", "/api/v1/admin/chat/sessions", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		fmt.Printf("❌ Session listing failed: %d - %s\n", w.Code, w.Body.String())
		return
	}
	fmt.Println("✅ Session listing successful")

	fmt.Println("\n🎉 All basic tests passed!")
	fmt.Println("✅ POST /api/v1/admin/chat/sessions - Working")
	fmt.Println("✅ GET /api/v1/admin/chat/sessions - Working")
	fmt.Println("✅ Authentication middleware - Working")
}
