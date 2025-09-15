---
inclusion: fileMatch
fileMatchPattern: '*ai*|*bedrock*|*intelligence*'
---

# AI Integration Guidelines

## Overview

This document provides comprehensive guidelines for integrating AI services, particularly AWS Bedrock, into the Cloud Consulting Platform. It covers prompt engineering, response handling, error management, and best practices for AI-powered features.

## AWS Bedrock Integration Patterns

### Service Configuration
```go
type BedrockConfig struct {
    Region          string        `env:"BEDROCK_REGION" default:"us-east-1"`
    ModelID         string        `env:"BEDROCK_MODEL_ID" default:"amazon.nova-lite-v1:0"`
    BaseURL         string        `env:"BEDROCK_BASE_URL"`
    BearerToken     string        `env:"AWS_BEARER_TOKEN_BEDROCK"`
    TimeoutSeconds  int           `env:"BEDROCK_TIMEOUT_SECONDS" default:"30"`
    MaxRetries      int           `env:"BEDROCK_MAX_RETRIES" default:"3"`
    RetryDelay      time.Duration `env:"BEDROCK_RETRY_DELAY" default:"1s"`
}
```

### Service Implementation Pattern
```go
type BedrockService struct {
    config *BedrockConfig
    client *http.Client
    logger *logrus.Logger
}

func (s *BedrockService) GenerateResponse(ctx context.Context, req *AIRequest) (*AIResponse, error) {
    // Implement with proper error handling, retries, and logging
    // Include request/response tracking for monitoring
    // Handle rate limiting and quota management
}
```

## Prompt Engineering Standards

### Prompt Structure
All AI prompts should follow this structure:
1. **System Context**: Define the AI's role and behavior
2. **Task Description**: Clear description of what needs to be done
3. **Input Data**: Structured input data with clear formatting
4. **Output Format**: Specify expected response format
5. **Constraints**: Any limitations or requirements

### Example Prompt Template
```go
const ConsultingReportPrompt = `
You are an expert AWS cloud consultant with deep knowledge of cloud architecture, migration strategies, and cost optimization.

TASK: Generate a professional consulting report based on the client inquiry below.

CLIENT INQUIRY:
Company: {{.Company}}
Services Requested: {{.Services}}
Message: {{.Message}}
Industry: {{.Industry}}

OUTPUT FORMAT:
Please provide a structured report with the following sections:
1. Executive Summary (2-3 sentences)
2. Current State Assessment (based on provided information)
3. Recommendations (3-5 specific recommendations)
4. Next Steps (actionable items)
5. Timeline Estimate

CONSTRAINTS:
- Keep the report professional and actionable
- Focus on AWS services and best practices
- Provide specific, measurable recommendations
- Limit response to 800-1000 words
- Use clear, business-friendly language
`
```

### Prompt Versioning
- Version all prompts using semantic versioning (v1.0.0)
- Store prompts in configuration files or database
- Implement A/B testing for prompt optimization
- Track prompt performance metrics

```go
type PromptTemplate struct {
    ID          string            `json:"id"`
    Version     string            `json:"version"`
    Name        string            `json:"name"`
    Template    string            `json:"template"`
    Variables   []string          `json:"variables"`
    CreatedAt   time.Time         `json:"created_at"`
    Metadata    map[string]string `json:"metadata"`
}
```

## Context Management

### Conversation Context
For chat-based AI interactions, maintain context across messages:

```go
type ConversationContext struct {
    SessionID     string         `json:"session_id"`
    Messages      []ChatMessage  `json:"messages"`
    UserProfile   *UserProfile   `json:"user_profile,omitempty"`
    CompanyInfo   *CompanyInfo   `json:"company_info,omitempty"`
    MaxMessages   int            `json:"max_messages"` // Limit context size
}

func (c *ConversationContext) BuildPrompt(newMessage string) string {
    // Include last N messages for context
    // Add user/company information if available
    // Format for optimal AI understanding
}
```

### Knowledge Base Integration
```go
type KnowledgeService interface {
    SearchRelevantDocs(ctx context.Context, query string) ([]*Document, error)
    GetCompanyKnowledge(ctx context.Context, companyID string) (*CompanyKnowledge, error)
    GetServiceTemplates(ctx context.Context, serviceType string) ([]*Template, error)
}

// Include relevant knowledge in AI prompts
func (s *AIService) EnhancePromptWithKnowledge(prompt string, query string) (string, error) {
    docs, err := s.knowledgeService.SearchRelevantDocs(ctx, query)
    if err != nil {
        return prompt, nil // Graceful degradation
    }
    
    return prompt + "\n\nRELEVANT KNOWLEDGE:\n" + formatDocs(docs), nil
}
```

## Response Processing

### Response Validation
```go
type AIResponse struct {
    Content     string            `json:"content"`
    Confidence  float64           `json:"confidence,omitempty"`
    Sources     []string          `json:"sources,omitempty"`
    Metadata    map[string]string `json:"metadata,omitempty"`
    ProcessedAt time.Time         `json:"processed_at"`
}

func (s *BedrockService) ValidateResponse(response *AIResponse) error {
    // Check response length and format
    // Validate content quality
    // Ensure no harmful content
    // Verify required sections are present
}
```

### Response Enhancement
```go
func (s *AIService) ProcessResponse(rawResponse string) (*ProcessedResponse, error) {
    // Parse structured content
    // Add formatting and styling
    // Include relevant links and references
    // Generate summary if needed
    
    return &ProcessedResponse{
        FormattedContent: formatted,
        Summary:         summary,
        ActionItems:     actionItems,
        References:      references,
    }, nil
}
```

## Error Handling and Fallbacks

### Error Categories
1. **Network Errors**: Timeout, connection issues
2. **Authentication Errors**: Invalid credentials, expired tokens
3. **Rate Limiting**: API quota exceeded
4. **Content Errors**: Invalid response format, harmful content
5. **Service Errors**: Bedrock service unavailable

### Fallback Strategies
```go
func (s *AIService) GenerateWithFallback(ctx context.Context, req *AIRequest) (*AIResponse, error) {
    // Primary: AWS Bedrock
    response, err := s.bedrockService.Generate(ctx, req)
    if err == nil {
        return response, nil
    }
    
    // Fallback 1: Cached similar response
    if cached := s.getCachedResponse(req); cached != nil {
        return cached, nil
    }
    
    // Fallback 2: Template-based response
    if template := s.getTemplateResponse(req); template != nil {
        return template, nil
    }
    
    // Fallback 3: Human handoff
    return s.createHumanHandoffResponse(req), nil
}
```

### Graceful Degradation
```go
func (h *ChatHandler) SendMessage(c *gin.Context) {
    message, err := h.chatService.SendMessage(ctx, req)
    if err != nil {
        // Log error but don't fail the request
        log.WithError(err).Warn("AI service unavailable, using fallback")
        
        // Return acknowledgment with fallback message
        c.JSON(http.StatusOK, gin.H{
            "message_id": message.ID,
            "status":     "received",
            "ai_response": "I'm currently experiencing technical difficulties. A human consultant will respond shortly.",
        })
        return
    }
    
    // Normal AI response
    c.JSON(http.StatusOK, message)
}
```

## Performance Optimization

### Caching Strategies
```go
type AICache interface {
    Get(ctx context.Context, key string) (*AIResponse, error)
    Set(ctx context.Context, key string, response *AIResponse, ttl time.Duration) error
    InvalidatePattern(ctx context.Context, pattern string) error
}

func (s *AIService) getCacheKey(request *AIRequest) string {
    // Create deterministic cache key based on:
    // - Prompt template version
    // - Input parameters
    // - User context (if relevant)
    return fmt.Sprintf("ai:response:%s", hash)
}
```

### Async Processing
```go
func (s *AIService) GenerateAsync(ctx context.Context, req *AIRequest) (*AsyncResponse, error) {
    // Start async processing
    jobID := uuid.New().String()
    
    go func() {
        response, err := s.bedrockService.Generate(ctx, req)
        if err != nil {
            s.notifyError(jobID, err)
            return
        }
        
        s.notifyComplete(jobID, response)
    }()
    
    return &AsyncResponse{JobID: jobID, Status: "processing"}, nil
}
```

### Streaming Responses
```go
func (s *AIService) GenerateStream(ctx context.Context, req *AIRequest) (<-chan *StreamChunk, error) {
    chunks := make(chan *StreamChunk, 100)
    
    go func() {
        defer close(chunks)
        
        // Stream response chunks as they arrive
        for chunk := range s.bedrockService.StreamGenerate(ctx, req) {
            select {
            case chunks <- chunk:
            case <-ctx.Done():
                return
            }
        }
    }()
    
    return chunks, nil
}
```

## Monitoring and Analytics

### Metrics to Track
```go
type AIMetrics struct {
    RequestCount      int64         `json:"request_count"`
    SuccessRate       float64       `json:"success_rate"`
    AverageLatency    time.Duration `json:"average_latency"`
    TokensUsed        int64         `json:"tokens_used"`
    CacheHitRate      float64       `json:"cache_hit_rate"`
    ErrorsByType      map[string]int64 `json:"errors_by_type"`
    CostEstimate      float64       `json:"cost_estimate"`
}
```

### Logging Standards
```go
log.WithFields(log.Fields{
    "request_id":     requestID,
    "user_id":        userID,
    "model":          "bedrock-nova",
    "prompt_version": "v1.2.0",
    "tokens_used":    tokensUsed,
    "latency_ms":     latency.Milliseconds(),
    "cache_hit":      cacheHit,
}).Info("AI request processed")
```

### Cost Monitoring
```go
type CostTracker struct {
    TokensUsed    int64   `json:"tokens_used"`
    RequestCount  int64   `json:"request_count"`
    EstimatedCost float64 `json:"estimated_cost"`
    Period        string  `json:"period"` // daily, monthly
}

func (s *AIService) TrackUsage(tokens int64, model string) {
    // Track token usage and estimate costs
    // Implement alerts for cost thresholds
    // Generate usage reports
}
```

## Security and Compliance

### Input Sanitization
```go
func (s *AIService) SanitizeInput(input string) (string, error) {
    // Remove or escape potentially harmful content
    // Validate input length and format
    // Check for injection attempts
    // Filter sensitive information
    
    if len(input) > MaxInputLength {
        return "", errors.New("input too long")
    }
    
    return sanitized, nil
}
```

### Output Filtering
```go
func (s *AIService) FilterOutput(output string) (string, error) {
    // Check for harmful or inappropriate content
    // Validate business context appropriateness
    // Remove any leaked sensitive information
    // Ensure professional tone
    
    return filtered, nil
}
```

### Data Privacy
```go
type PrivacyConfig struct {
    LogUserData      bool `env:"AI_LOG_USER_DATA" default:"false"`
    RetainHistory    bool `env:"AI_RETAIN_HISTORY" default:"true"`
    RetentionDays    int  `env:"AI_RETENTION_DAYS" default:"90"`
    AnonymizeData    bool `env:"AI_ANONYMIZE_DATA" default:"true"`
}
```

## Testing Strategies

### Unit Testing
```go
func TestAIService_GenerateResponse(t *testing.T) {
    tests := []struct {
        name           string
        request        *AIRequest
        mockResponse   *AIResponse
        expectedError  error
    }{
        {
            name: "successful generation",
            request: &AIRequest{
                Prompt: "Generate a report for...",
                Context: map[string]string{"company": "TechCorp"},
            },
            mockResponse: &AIResponse{
                Content: "Executive Summary...",
            },
        },
        // More test cases
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

### Integration Testing
```go
func TestBedrockIntegration(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test")
    }
    
    service := NewBedrockService(testConfig)
    
    response, err := service.GenerateResponse(context.Background(), &AIRequest{
        Prompt: testPrompt,
    })
    
    assert.NoError(t, err)
    assert.NotEmpty(t, response.Content)
    assert.True(t, len(response.Content) > 100) // Reasonable response length
}
```

### Load Testing
```go
func BenchmarkAIService_GenerateResponse(b *testing.B) {
    service := NewAIService(testConfig)
    request := &AIRequest{Prompt: "Test prompt"}
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, err := service.GenerateResponse(context.Background(), request)
        if err != nil {
            b.Fatal(err)
        }
    }
}
```

## Best Practices Summary

1. **Always implement fallbacks** - Never let AI service failures break user experience
2. **Cache aggressively** - AI responses are expensive, cache when possible
3. **Monitor costs** - Track token usage and implement cost controls
4. **Version prompts** - Treat prompts as code with proper versioning
5. **Validate everything** - Sanitize inputs and validate outputs
6. **Log comprehensively** - Track all AI interactions for debugging and optimization
7. **Test thoroughly** - Include unit, integration, and load tests
8. **Handle async gracefully** - Use async processing for long-running AI tasks
9. **Implement streaming** - Stream responses for better user experience
10. **Plan for scale** - Design for high-volume AI usage from day one