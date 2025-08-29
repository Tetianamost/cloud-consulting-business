package main

import (
	"context"
	"fmt"
	"log"

	"github.com/cloud-consulting/backend/internal/config"
	"github.com/cloud-consulting/backend/internal/handlers"
	"github.com/cloud-consulting/backend/internal/interfaces"
	"github.com/cloud-consulting/backend/internal/services"
	"github.com/sirupsen/logrus"
)

func main() {
	fmt.Println("Debugging SimpleChatHandler AI Response Issue")
	fmt.Println("============================================")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create logger
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	// Create Bedrock service
	bedrockService := services.NewBedrockService(&cfg.Bedrock)

	// Create simple chat handler
	_ = handlers.NewSimpleChatHandler(logger, bedrockService)

	fmt.Println("\n1. Testing Bedrock Service Directly")
	fmt.Println("-----------------------------------")

	ctx := context.Background()
	testMessage := "What are AWS security best practices?"

	// Test Bedrock service directly
	prompt := fmt.Sprintf(`You are an expert AWS cloud consultant with deep knowledge of cloud architecture, migration strategies, cost optimization, and best practices. 

A client has asked: "%s"

Please provide a professional, detailed, and actionable response that:
1. Addresses their specific question or concern
2. Provides concrete AWS service recommendations where appropriate
3. Includes best practices and considerations
4. Offers next steps or implementation guidance
5. Maintains a professional consulting tone

Response:`, testMessage)

	options := &interfaces.BedrockOptions{
		ModelID:     "amazon.nova-lite-v1:0",
		MaxTokens:   1000,
		Temperature: 0.7,
		TopP:        0.9,
	}

	fmt.Printf("Testing Bedrock service with prompt length: %d characters\n", len(prompt))
	response, err := bedrockService.GenerateText(ctx, prompt, options)
	if err != nil {
		fmt.Printf("❌ Bedrock service failed: %v\n", err)
		fmt.Println("This explains why the chat is returning empty responses!")
	} else {
		fmt.Printf("✅ Bedrock service succeeded!\n")
		fmt.Printf("Response content length: %d characters\n", len(response.Content))
		fmt.Printf("Response content preview: %.200s...\n", response.Content)

		if response.Content == "" {
			fmt.Println("❌ Bedrock returned empty content!")
		}
	}

	fmt.Println("\n2. Testing SimpleChatHandler generateAIResponse Method")
	fmt.Println("----------------------------------------------------")

	// We can't directly call generateAIResponse since it's private, but we can test the flow
	// by examining what happens when we send a message

	fmt.Println("The generateAIResponse method should:")
	fmt.Println("1. Call bedrockService.GenerateText with the prompt")
	fmt.Println("2. Return the response.Content if successful")
	fmt.Println("3. Return an error if Bedrock fails")
	fmt.Println("4. The SendMessage handler should catch errors and use fallback responses")

	fmt.Println("\n3. Testing Fallback Response Generation")
	fmt.Println("--------------------------------------")

	// Test fallback response for security question
	fmt.Printf("Testing fallback for: '%s'\n", testMessage)

	// We can't directly call generateFallbackResponse either, but we know it should work
	// based on the message content
	if contains(testMessage, "security") {
		fmt.Println("✅ Message contains 'security' - should trigger security fallback")
		expectedFallback := "AWS security best practices include: implementing the principle of least privilege with IAM..."
		fmt.Printf("Expected fallback preview: %.100s...\n", expectedFallback)
	}

	fmt.Println("\n4. Diagnosis")
	fmt.Println("------------")

	if err != nil {
		fmt.Println("🔍 ISSUE FOUND: Bedrock service is failing")
		fmt.Println("   - The generateAIResponse method returns an error")
		fmt.Println("   - The SendMessage handler should catch this and use fallback")
		fmt.Println("   - But the AI message is still showing empty content")
		fmt.Println("   - This suggests there's a bug in the error handling logic")
	} else if response != nil && response.Content == "" {
		fmt.Println("🔍 ISSUE FOUND: Bedrock service returns empty content")
		fmt.Println("   - The generateAIResponse method doesn't return an error")
		fmt.Println("   - But the content is empty")
		fmt.Println("   - The SendMessage handler should detect this and use fallback")
	} else {
		fmt.Println("🔍 ISSUE UNCLEAR: Bedrock service seems to work in isolation")
		fmt.Println("   - The issue might be in the SimpleChatHandler integration")
		fmt.Println("   - Or there might be a different error in the live system")
	}

	fmt.Println("\n5. Recommended Fix")
	fmt.Println("------------------")
	fmt.Println("Check the SimpleChatHandler.SendMessage method to ensure:")
	fmt.Println("1. Errors from generateAIResponse are properly caught")
	fmt.Println("2. Empty responses are properly detected")
	fmt.Println("3. Fallback responses are properly assigned to the AI message")
	fmt.Println("4. The AI message content is never left empty")
}

func contains(text, substr string) bool {
	return len(text) >= len(substr) &&
		(text == substr ||
			(len(text) > len(substr) &&
				(text[:len(substr)] == substr ||
					text[len(text)-len(substr):] == substr ||
					containsSubstring(text, substr))))
}

func containsSubstring(text, substr string) bool {
	for i := 0; i <= len(text)-len(substr); i++ {
		if text[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
