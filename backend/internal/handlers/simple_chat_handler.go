package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/cloud-consulting/backend/internal/interfaces"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// SimpleChatHandler - Enhanced chat handler with real AI integration
type SimpleChatHandler struct {
	logger         *logrus.Logger
	messages       []SimpleChatMessage
	bedrockService interfaces.BedrockService
}

type SimpleChatMessage struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	Role      string    `json:"role"`
	Timestamp time.Time `json:"timestamp"`
	SessionID string    `json:"session_id"`
}

type SimpleSendMessageRequest struct {
	Content   string `json:"content" binding:"required"`
	SessionID string `json:"session_id"`
}

type SimpleSendMessageResponse struct {
	Success   bool   `json:"success"`
	MessageID string `json:"message_id"`
	Error     string `json:"error,omitempty"`
}

type SimpleGetMessagesResponse struct {
	Success  bool                `json:"success"`
	Messages []SimpleChatMessage `json:"messages"`
}

// NewSimpleChatHandler creates a new enhanced chat handler
func NewSimpleChatHandler(logger *logrus.Logger, bedrockService interfaces.BedrockService) *SimpleChatHandler {
	return &SimpleChatHandler{
		logger:         logger,
		messages:       make([]SimpleChatMessage, 0),
		bedrockService: bedrockService,
	}
}

// SendMessage handles sending a message and immediately returns a mock AI response
func (h *SimpleChatHandler) SendMessage(c *gin.Context) {
	var req SimpleSendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Error("Failed to bind request")
		c.JSON(http.StatusBadRequest, SimpleSendMessageResponse{
			Success: false,
			Error:   "Invalid request format",
		})
		return
	}

	// Generate message ID
	userMsgID := fmt.Sprintf("msg-%d-user", time.Now().UnixNano())
	aiMsgID := fmt.Sprintf("msg-%d-ai", time.Now().UnixNano())

	// Store user message
	userMessage := SimpleChatMessage{
		ID:        userMsgID,
		Content:   req.Content,
		Role:      "user",
		Timestamp: time.Now(),
		SessionID: req.SessionID,
	}
	h.messages = append(h.messages, userMessage)

	// Generate AI response using Bedrock with fallback
	h.logger.Info("Attempting to generate AI response using Bedrock")
	aiResponse, err := h.generateAIResponse(c.Request.Context(), req.Content)
	if err != nil {
		h.logger.WithError(err).Error("Failed to generate AI response, using fallback")
		aiResponse = h.generateFallbackResponse(req.Content)
	} else if aiResponse == "" {
		h.logger.Warn("AI response was empty, using fallback")
		aiResponse = h.generateFallbackResponse(req.Content)
	} else {
		h.logger.Info("Successfully generated AI response using Bedrock")
	}

	aiMessage := SimpleChatMessage{
		ID:        aiMsgID,
		Content:   aiResponse,
		Role:      "assistant",
		Timestamp: time.Now().Add(time.Millisecond * 100), // Slight delay
		SessionID: req.SessionID,
	}
	h.messages = append(h.messages, aiMessage)

	h.logger.WithFields(logrus.Fields{
		"user_message_id": userMsgID,
		"ai_message_id":   aiMsgID,
		"session_id":      req.SessionID,
	}).Info("Messages stored successfully")

	c.JSON(http.StatusOK, SimpleSendMessageResponse{
		Success:   true,
		MessageID: userMsgID,
	})
}

// GetMessages retrieves all messages for a session
func (h *SimpleChatHandler) GetMessages(c *gin.Context) {
	sessionID := c.Query("session_id")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, SimpleGetMessagesResponse{
			Success: false,
		})
		return
	}

	// Filter messages by session ID
	var sessionMessages []SimpleChatMessage
	for _, msg := range h.messages {
		if msg.SessionID == sessionID {
			sessionMessages = append(sessionMessages, msg)
		}
	}

	h.logger.WithFields(logrus.Fields{
		"session_id":     sessionID,
		"message_count":  len(sessionMessages),
		"total_messages": len(h.messages),
	}).Info("Retrieved messages for session")

	c.JSON(http.StatusOK, SimpleGetMessagesResponse{
		Success:  true,
		Messages: sessionMessages,
	})
}

// generateAIResponse generates a real AI response using Bedrock
func (h *SimpleChatHandler) generateAIResponse(ctx context.Context, userMessage string) (string, error) {
	// Create a professional consulting prompt
	prompt := fmt.Sprintf(`You are an expert AWS cloud consultant with deep knowledge of cloud architecture, migration strategies, cost optimization, and best practices. 

A client has asked: "%s"

Please provide a professional, detailed, and actionable response that:
1. Addresses their specific question or concern
2. Provides concrete AWS service recommendations where appropriate
3. Includes best practices and considerations
4. Offers next steps or implementation guidance
5. Maintains a professional consulting tone

Response:`, userMessage)

	// Set up Bedrock options for optimal response
	options := &interfaces.BedrockOptions{
		ModelID:     "amazon.nova-lite-v1:0", // Use Nova Lite model
		MaxTokens:   1000,
		Temperature: 0.7,
		TopP:        0.9,
	}

	// Generate response using Bedrock
	response, err := h.bedrockService.GenerateText(ctx, prompt, options)
	if err != nil {
		return "", fmt.Errorf("failed to generate AI response: %w", err)
	}

	return response.Content, nil
}

// generateFallbackResponse creates intelligent fallback responses when Bedrock is unavailable
func (h *SimpleChatHandler) generateFallbackResponse(userMessage string) string {
	userMessageLower := strings.ToLower(userMessage)

	// Cost-related queries
	if strings.Contains(userMessageLower, "cost") || strings.Contains(userMessageLower, "price") || strings.Contains(userMessageLower, "budget") {
		return "For AWS cost optimization, I recommend starting with the AWS Cost Calculator to estimate your workload costs. Key strategies include: using Reserved Instances for predictable workloads, implementing auto-scaling to match demand, leveraging S3 storage classes for data lifecycle management, and regularly reviewing your AWS Cost and Usage Reports. Would you like me to elaborate on any of these cost optimization strategies?"
	}

	// Security-related queries
	if strings.Contains(userMessageLower, "security") || strings.Contains(userMessageLower, "secure") || strings.Contains(userMessageLower, "compliance") {
		return "AWS security best practices include: implementing the principle of least privilege with IAM, enabling AWS CloudTrail for audit logging, using AWS Config for compliance monitoring, encrypting data at rest and in transit, and implementing network security with VPCs and security groups. For compliance frameworks like SOC 2, HIPAA, or PCI DSS, AWS provides specific guidance and services. What specific security concerns would you like to address?"
	}

	// Migration-related queries
	if strings.Contains(userMessageLower, "migrat") || strings.Contains(userMessageLower, "move") || strings.Contains(userMessageLower, "cloud") {
		return "AWS migration typically follows the 6 R's strategy: Rehost (lift-and-shift), Replatform (lift-tinker-and-shift), Refactor/Re-architect, Repurchase, Retain, or Retire. I recommend starting with the AWS Migration Hub to assess your current environment, then using AWS Application Discovery Service to understand dependencies. The AWS Database Migration Service and AWS Server Migration Service can help with the actual migration process. What type of workload are you looking to migrate?"
	}

	// Architecture-related queries
	if strings.Contains(userMessageLower, "architect") || strings.Contains(userMessageLower, "design") || strings.Contains(userMessageLower, "structure") {
		return "AWS Well-Architected Framework provides guidance across five pillars: Operational Excellence, Security, Reliability, Performance Efficiency, and Cost Optimization. For a typical web application, I'd recommend a multi-tier architecture using Application Load Balancer, Auto Scaling Groups, RDS for databases, and CloudFront for content delivery. The specific services depend on your requirements for scalability, availability, and performance. What type of application architecture are you designing?"
	}

	// Performance-related queries
	if strings.Contains(userMessageLower, "performance") || strings.Contains(userMessageLower, "speed") || strings.Contains(userMessageLower, "latency") {
		return "AWS performance optimization strategies include: using CloudFront CDN for global content delivery, implementing ElastiCache for caching, choosing appropriate EC2 instance types, optimizing database performance with RDS Performance Insights, and using AWS X-Ray for application tracing. Auto Scaling ensures you have the right capacity to meet demand. What specific performance challenges are you experiencing?"
	}

	// Default professional response
	return "Thank you for your inquiry. As an AWS cloud consultant, I'm here to help you with cloud architecture, migration strategies, cost optimization, security best practices, and performance improvements. Based on your question, I'd recommend we start by understanding your current infrastructure and specific requirements. Could you provide more details about your use case, current environment, or specific challenges you're facing? This will help me provide more targeted recommendations for your AWS journey."
}
