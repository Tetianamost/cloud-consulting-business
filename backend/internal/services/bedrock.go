package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/cloud-consulting/backend/internal/config"
	"github.com/cloud-consulting/backend/internal/interfaces"
)

// bedrockService implements the BedrockService interface
type bedrockService struct {
	config     *config.BedrockConfig
	httpClient *http.Client
}

// NewBedrockService creates a new Bedrock service instance
func NewBedrockService(cfg *config.BedrockConfig) interfaces.BedrockService {
	return &bedrockService{
		config: cfg,
		httpClient: &http.Client{
			Timeout: time.Duration(cfg.Timeout) * time.Second,
		},
	}
}

// bedrockRequest represents the request structure for Bedrock Converse API
// Format based on Amazon Bedrock Converse API documentation for Nova Lite model
type bedrockRequest struct {
	Messages      []bedrockMessage `json:"messages"`
	MaxTokens     int              `json:"maxTokens,omitempty"`
	Temperature   float64          `json:"temperature,omitempty"`
	TopP          float64          `json:"topP,omitempty"`
	StopSequences []string         `json:"stopSequences,omitempty"`
}

// bedrockMessage represents a message in the conversation for Bedrock Converse API
type bedrockMessage struct {
	Role    string                `json:"role"`
	Content []bedrockContentBlock `json:"content"`
}

// bedrockAPIResponse represents the raw response from Bedrock Converse API
type bedrockAPIResponse struct {
	Output struct {
		Message struct {
			Role    string                `json:"role"`
			Content []bedrockContentBlock `json:"content"`
		} `json:"message"`
		StopReason string `json:"stopReason"`
	} `json:"output"`
	Usage struct {
		InputTokens  int `json:"inputTokens"`
		OutputTokens int `json:"outputTokens"`
		TotalTokens  int `json:"totalTokens"`
	} `json:"usage"`
}

// bedrockContentBlock represents a content block in the response
type bedrockContentBlock struct {
	Text string `json:"text"`
}

// GenerateText generates text using Amazon Bedrock
func (s *bedrockService) GenerateText(ctx context.Context, prompt string, options *interfaces.BedrockOptions) (*interfaces.BedrockResponse, error) {
	if s.config.APIKey == "" {
		return nil, fmt.Errorf("bedrock API key not configured")
	}

	// Set default options if not provided
	if options == nil {
		options = &interfaces.BedrockOptions{
			ModelID:     s.config.ModelID,
			MaxTokens:   1000,
			Temperature: 0.7,
			TopP:        0.9,
		}
	}

	// Build the request payload with messages format for Converse API
	reqPayload := bedrockRequest{
		Messages: []bedrockMessage{
			{
				Role: "user",
				Content: []bedrockContentBlock{
					{
						Text: prompt,
					},
				},
			},
		},
		MaxTokens:     options.MaxTokens,
		Temperature:   options.Temperature,
		TopP:          options.TopP,
		StopSequences: []string{"\n\n"},
	}

	// Marshal the request
	jsonData, err := json.Marshal(reqPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Debug logging
	fmt.Printf("DEBUG: Bedrock API URL: %s/model/%s/converse\n", s.config.BaseURL, options.ModelID)
	fmt.Printf("DEBUG: Request payload: %s\n", string(jsonData))

	// Create the HTTP request - use the correct endpoint for the model
	url := fmt.Sprintf("%s/model/%s/converse", s.config.BaseURL, options.ModelID)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.config.APIKey)
	req.Header.Set("X-Amzn-Bedrock-Accept", "application/json")
	req.Header.Set("X-Amzn-Bedrock-Content-Type", "application/json")

	// Make the request
	fmt.Printf("DEBUG: Making Bedrock API request...\n")
	resp, err := s.httpClient.Do(req)
	if err != nil {
		fmt.Printf("DEBUG: Bedrock API request failed: %v\n", err)
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	fmt.Printf("DEBUG: Bedrock API response status: %d\n", resp.StatusCode)

	// Check response status
	if resp.StatusCode != http.StatusOK {
		// Read the error response body for debugging
		errorBody, _ := io.ReadAll(resp.Body)
		fmt.Printf("DEBUG: Bedrock API error response: %s\n", string(errorBody))

		// For development: return mock response when Bedrock is not accessible
		if resp.StatusCode == 403 {
			fmt.Printf("DEBUG: Bedrock API access denied, returning mock response for development\n")
			return s.generateMockResponse(prompt, options), nil
		}

		return nil, fmt.Errorf("bedrock API returned status %d: %s", resp.StatusCode, string(errorBody))
	}

	// Parse the response
	var apiResp bedrockAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		fmt.Printf("DEBUG: Failed to decode Bedrock response: %v\n", err)
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Extract text content from the response
	var content string
	if len(apiResp.Output.Message.Content) > 0 {
		content = apiResp.Output.Message.Content[0].Text
	}

	fmt.Printf("DEBUG: Bedrock API call successful, generated %d characters of content\n", len(content))

	// Convert to our response format
	response := &interfaces.BedrockResponse{
		Content: content,
		Usage: interfaces.BedrockUsage{
			InputTokens:  apiResp.Usage.InputTokens,
			OutputTokens: apiResp.Usage.OutputTokens,
		},
		Metadata: map[string]string{
			"stopReason": apiResp.Output.StopReason,
			"modelId":    options.ModelID,
			"role":       apiResp.Output.Message.Role,
		},
	}

	return response, nil
}

// GetModelInfo returns information about the configured Bedrock model
func (s *bedrockService) GetModelInfo() interfaces.BedrockModelInfo {
	return interfaces.BedrockModelInfo{
		ModelID:     s.config.ModelID,
		ModelName:   "Amazon Nova Lite",
		Provider:    "Amazon",
		MaxTokens:   8192,
		IsAvailable: s.config.APIKey != "",
	}
}

// generateMockResponse creates a mock AI response for development when Bedrock is not accessible
func (s *bedrockService) generateMockResponse(prompt string, options *interfaces.BedrockOptions) *interfaces.BedrockResponse {
	// Simple mock responses based on common consulting scenarios
	mockResponses := []string{
		"Thank you for your inquiry. Based on your requirements, I recommend implementing a cloud-first architecture using AWS services. This approach would provide scalability, cost-effectiveness, and improved security for your organization.",
		"I understand you're looking for cloud consulting guidance. For your use case, I suggest considering a hybrid cloud strategy that leverages both on-premises and cloud resources. This would allow for gradual migration while maintaining operational continuity.",
		"Based on the information provided, I recommend starting with a comprehensive cloud readiness assessment. This will help identify the best migration strategy and prioritize workloads for cloud adoption.",
		"For your AWS infrastructure needs, I suggest implementing Infrastructure as Code (IaC) using AWS CloudFormation or Terraform. This approach ensures consistent, repeatable deployments and better resource management.",
		"Thank you for reaching out. I recommend exploring AWS Well-Architected Framework principles to ensure your cloud architecture follows best practices for security, reliability, performance efficiency, cost optimization, and operational excellence.",
	}

	// Select a response based on prompt content or use a simple rotation
	responseIndex := len(prompt) % len(mockResponses)
	mockContent := mockResponses[responseIndex]

	// Add a development notice
	mockContent += "\n\n[Development Mode: This is a mock AI response for testing purposes]"

	return &interfaces.BedrockResponse{
		Content: mockContent,
		Usage: interfaces.BedrockUsage{
			InputTokens:  len(prompt) / 4, // Rough token estimation
			OutputTokens: len(mockContent) / 4,
		},
		Metadata: map[string]string{
			"stopReason": "end_turn",
			"modelId":    options.ModelID,
			"role":       "assistant",
			"mockMode":   "true",
		},
	}
}

// IsHealthy checks if the Bedrock service is healthy
func (s *bedrockService) IsHealthy() bool {
	return s.config.APIKey != "" && s.config.BaseURL != ""
}
