package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/cloud-consulting/backend/internal/config"
	"github.com/cloud-consulting/backend/internal/handlers"
	"github.com/cloud-consulting/backend/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	fmt.Println("Testing SimpleChatHandler Direct Integration")
	fmt.Println("===========================================")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create logger with debug level
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	// Create Bedrock service
	bedrockService := services.NewBedrockService(&cfg.Bedrock)

	// Create simple chat handler
	handler := handlers.NewSimpleChatHandler(logger, bedrockService)

	// Set up Gin in test mode
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/messages", handler.SendMessage)
	router.GET("/messages", handler.GetMessages)

	fmt.Println("\n1. Testing SendMessage Endpoint")
	fmt.Println("-------------------------------")

	// Create request payload
	payload := map[string]string{
		"content":    "What are AWS security best practices?",
		"session_id": "debug-session-123",
	}

	jsonPayload, _ := json.Marshal(payload)
	fmt.Printf("Sending payload: %s\n", string(jsonPayload))

	// Create HTTP request
	req, _ := http.NewRequest("POST", "/messages", bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()

	// Execute request
	fmt.Println("Executing SendMessage request...")
	router.ServeHTTP(w, req)

	fmt.Printf("Response status: %d\n", w.Code)
	fmt.Printf("Response body: %s\n", w.Body.String())

	if w.Code != http.StatusOK {
		fmt.Println("❌ SendMessage failed!")
		return
	}

	fmt.Println("\n2. Testing GetMessages Endpoint")
	fmt.Println("-------------------------------")

	// Get messages for this session
	req2, _ := http.NewRequest("GET", "/messages?session_id=debug-session-123", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	fmt.Printf("Messages response status: %d\n", w2.Code)
	fmt.Printf("Messages response body: %s\n", w2.Body.String())

	// Parse messages response
	var messagesResp map[string]interface{}
	if err := json.Unmarshal(w2.Body.Bytes(), &messagesResp); err != nil {
		fmt.Printf("❌ Failed to parse messages response: %v\n", err)
		return
	}

	if !messagesResp["success"].(bool) {
		fmt.Println("❌ GetMessages failed!")
		return
	}

	messages, ok := messagesResp["messages"].([]interface{})
	if !ok {
		fmt.Println("❌ Messages not found in response!")
		return
	}

	fmt.Printf("Total messages: %d\n", len(messages))

	fmt.Println("\n3. Analyzing Messages")
	fmt.Println("---------------------")

	for i, msg := range messages {
		msgMap, ok := msg.(map[string]interface{})
		if !ok {
			fmt.Printf("❌ Message %d is not a valid object\n", i+1)
			continue
		}

		role := msgMap["role"].(string)
		content := msgMap["content"].(string)
		id := msgMap["id"].(string)

		fmt.Printf("Message %d:\n", i+1)
		fmt.Printf("  ID: %s\n", id)
		fmt.Printf("  Role: %s\n", role)
		fmt.Printf("  Content length: %d characters\n", len(content))

		if role == "user" {
			fmt.Printf("  Content: %s\n", content)
		} else if role == "assistant" {
			if content == "" {
				fmt.Println("  ❌ ISSUE FOUND: AI response is empty!")
				fmt.Println("  This confirms the bug is in SimpleChatHandler")
			} else {
				fmt.Printf("  ✅ AI response generated successfully\n")
				fmt.Printf("  Content preview: %.200s...\n", content)
			}
		}
		fmt.Println()
	}

	fmt.Println("4. Diagnosis")
	fmt.Println("------------")

	if len(messages) < 2 {
		fmt.Println("❌ Expected 2 messages (user + AI), but got", len(messages))
	} else {
		userMsg := messages[0].(map[string]interface{})
		aiMsg := messages[1].(map[string]interface{})

		userContent := userMsg["content"].(string)
		aiContent := aiMsg["content"].(string)

		if userContent != "" && aiContent == "" {
			fmt.Println("🔍 CONFIRMED BUG:")
			fmt.Println("  - User message was stored correctly")
			fmt.Println("  - AI message was stored with empty content")
			fmt.Println("  - This means the generateAIResponse method is not working as expected")
			fmt.Println("  - OR the error handling is not working correctly")
			fmt.Println("  - OR there's a race condition in the message storage")
		} else if userContent != "" && aiContent != "" {
			fmt.Println("✅ NO BUG FOUND:")
			fmt.Println("  - Both messages were stored correctly")
			fmt.Println("  - The issue might be specific to the live server environment")
		}
	}
}
