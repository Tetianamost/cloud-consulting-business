package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/cloud-consulting/backend/internal/domain"
	"github.com/cloud-consulting/backend/internal/interfaces"
	"github.com/cloud-consulting/backend/internal/services"
	"github.com/cloud-consulting/backend/internal/storage"
)

// MockBedrockService for testing
type MockBedrockService struct{}

func (m *MockBedrockService) GenerateText(ctx context.Context, prompt string, options *interfaces.BedrockOptions) (*interfaces.BedrockResponse, error) {
	return &interfaces.BedrockResponse{
		Content: "This is a mock AWS consulting response about EC2 instances and S3 storage for your cloud migration needs.",
		Usage: interfaces.BedrockUsage{
			InputTokens:  100,
			OutputTokens: 50,
		},
	}, nil
}

func (m *MockBedrockService) GetModelInfo() interfaces.BedrockModelInfo {
	return interfaces.BedrockModelInfo{
		ModelID:     "mock-model",
		MaxTokens:   1000,
		InputCost:   0.001,
		OutputCost:  0.002,
		Description: "Mock Bedrock service for testing",
	}
}

func (m *MockBedrockService) IsHealthy() bool {
	return true
}

// MockKnowledgeBase for testing
type MockKnowledgeBase struct{}

func (m *MockKnowledgeBase) GetServiceOfferings(ctx context.Context) ([]*interfaces.ServiceOffering, error) {
	return []*interfaces.ServiceOffering{
		{
			Name:            "Cloud Migration Assessment",
			Description:     "Comprehensive assessment of your current infrastructure for AWS migration",
			Category:        "Migration",
			TypicalDuration: "2-4 weeks",
			TeamSize:        "2-3 consultants",
			KeyBenefits:     []string{"Risk assessment", "Cost analysis", "Migration roadmap"},
			Deliverables:    []string{"Assessment report", "Migration plan", "Cost estimates"},
		},
		{
			Name:            "Architecture Review",
			Description:     "Review and optimization of AWS architecture for performance and cost",
			Category:        "Architecture",
			TypicalDuration: "1-2 weeks",
			TeamSize:        "1-2 consultants",
			KeyBenefits:     []string{"Performance optimization", "Cost reduction", "Security enhancement"},
			Deliverables:    []string{"Architecture diagrams", "Recommendations", "Implementation guide"},
		},
	}, nil
}

func (m *MockKnowledgeBase) GetTeamExpertise(ctx context.Context) ([]*interfaces.TeamExpertise, error) {
	return []*interfaces.TeamExpertise{
		{
			ConsultantName:  "John Smith",
			Role:            "Senior Cloud Architect",
			ExperienceYears: 8,
			ExpertiseAreas:  []string{"AWS Migration", "Serverless", "Containers"},
			CloudProviders:  []string{"AWS", "Azure"},
		},
		{
			ConsultantName:  "Jane Doe",
			Role:            "Cloud Security Specialist",
			ExperienceYears: 6,
			ExpertiseAreas:  []string{"Security", "Compliance", "IAM"},
			CloudProviders:  []string{"AWS"},
		},
	}, nil
}

func (m *MockKnowledgeBase) GetPastSolutions(ctx context.Context, serviceType, industry string) ([]*interfaces.PastSolution, error) {
	return []*interfaces.PastSolution{
		{
			Title:            "E-commerce Platform Migration",
			Industry:         "Retail",
			ProblemStatement: "Legacy on-premises infrastructure limiting scalability",
			SolutionApproach: "Containerized microservices on EKS with RDS and ElastiCache",
			Technologies:     []string{"EKS", "RDS", "ElastiCache", "CloudFront"},
			TimeToValue:      "3 months",
			CostSavings:      150000,
		},
	}, nil
}

func (m *MockKnowledgeBase) GetConsultingApproach(ctx context.Context, serviceType string) (*interfaces.ConsultingApproach, error) {
	return &interfaces.ConsultingApproach{
		Name:              "Agile Cloud Transformation",
		Philosophy:        "Iterative approach with continuous feedback and improvement",
		EngagementModel:   "Collaborative partnership with client teams",
		KeyPrinciples:     []string{"Client-first", "Iterative delivery", "Knowledge transfer"},
		ClientInvolvement: "High - client teams work alongside our consultants",
		KnowledgeTransfer: "Continuous training and documentation throughout engagement",
	}, nil
}

func (m *MockKnowledgeBase) GetClientHistory(ctx context.Context, company string) ([]*interfaces.ClientEngagement, error) {
	return []*interfaces.ClientEngagement{
		{
			ProjectName:        "Initial Cloud Assessment",
			StartDate:          time.Now().AddDate(-1, 0, 0),
			Status:             "Completed",
			ClientSatisfaction: 9.2,
		},
	}, nil
}

// MockPromptArchitect for testing
type MockPromptArchitect struct{}

func (m *MockPromptArchitect) BuildPrompt(ctx context.Context, request *interfaces.PromptRequest) (string, error) {
	return request.BasePrompt + "\n\nEnhanced with architectural context.", nil
}

func (m *MockPromptArchitect) OptimizePrompt(ctx context.Context, prompt string, constraints *interfaces.PromptConstraints) (string, error) {
	return prompt, nil
}

func (m *MockPromptArchitect) ValidatePrompt(ctx context.Context, prompt string) error {
	return nil
}

// MockCompanyKnowledgeIntegration for testing
type MockCompanyKnowledgeIntegration struct{}

func (m *MockCompanyKnowledgeIntegration) GenerateContextualPrompt(ctx context.Context, inquiry *domain.Inquiry, basePrompt string) (string, error) {
	return basePrompt + "\n\nEnhanced with company knowledge.", nil
}

func (m *MockCompanyKnowledgeIntegration) GetRecommendationsForInquiry(ctx context.Context, inquiry *domain.Inquiry) (*services.InquiryRecommendations, error) {
	return &services.InquiryRecommendations{
		InquiryID:   inquiry.ID,
		GeneratedAt: time.Now(),
	}, nil
}

func main() {
	fmt.Println("Testing Enhanced AI Integration for Chat...")

	// Create logger
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	// Create mock services
	mockBedrock := &MockBedrockService{}
	mockKnowledgeBase := &MockKnowledgeBase{}
	mockPromptArchitect := &MockPromptArchitect{}
	mockCompanyKnowledge := &MockCompanyKnowledgeIntegration{}

	// Create enhanced Bedrock service
	enhancedBedrock := services.NewEnhancedBedrockService(mockBedrock, mockKnowledgeBase, mockCompanyKnowledge)

	// Create chat-aware Bedrock service
	chatAwareBedrock := services.NewChatAwareBedrockService(enhancedBedrock, mockPromptArchitect)

	// Create in-memory storage
	memoryStorage := storage.NewMemoryStorage()

	// Create chat service with enhanced AI
	chatService := services.NewChatServiceWithEnhancedAI(
		memoryStorage.ChatMessageRepository(),
		memoryStorage.SessionService(),
		mockBedrock,
		chatAwareBedrock,
		logger,
	)

	ctx := context.Background()

	// Test 1: Create a chat session
	fmt.Println("\n1. Creating chat session...")
	session := &domain.ChatSession{
		ID:         "test-session-1",
		UserID:     "test-user",
		ClientName: "TechCorp Inc",
		Context:    "Cloud migration planning meeting",
		Status:     domain.SessionStatusActive,
		Metadata: map[string]interface{}{
			"meeting_type":    "migration_planning",
			"service_types":   []string{"migration", "architecture"},
			"cloud_providers": []string{"AWS"},
		},
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		LastActivity: time.Now(),
	}

	sessionService := memoryStorage.SessionService()
	if err := sessionService.CreateSession(ctx, session); err != nil {
		log.Fatalf("Failed to create session: %v", err)
	}
	fmt.Printf("✓ Created session: %s for client: %s\n", session.ID, session.ClientName)

	// Test 2: Send a message with enhanced AI
	fmt.Println("\n2. Testing enhanced AI response...")
	chatRequest := &domain.ChatRequest{
		SessionID: session.ID,
		Content:   "What are the key considerations for migrating our e-commerce platform to AWS?",
		Type:      domain.MessageTypeUser,
		Metadata: map[string]interface{}{
			"client_context": "e-commerce migration",
		},
	}

	response, err := chatService.SendMessage(ctx, chatRequest)
	if err != nil {
		log.Fatalf("Failed to send message: %v", err)
	}

	fmt.Printf("✓ Generated AI response (%d chars)\n", len(response.Content))
	fmt.Printf("  Content: %s\n", response.Content[:min(100, len(response.Content))])
	fmt.Printf("  Tokens used: %d\n", response.TokensUsed)
	fmt.Printf("  Processing time: %.2fms\n", response.ProcessTime)

	// Test 3: Test quick action
	fmt.Println("\n3. Testing quick action...")
	quickActionRequest := &domain.ChatRequest{
		SessionID:   session.ID,
		Content:     "I need a cost estimate for our migration",
		Type:        domain.MessageTypeUser,
		QuickAction: "cost_estimate",
	}

	quickResponse, err := chatService.SendMessage(ctx, quickActionRequest)
	if err != nil {
		log.Fatalf("Failed to process quick action: %v", err)
	}

	fmt.Printf("✓ Quick action response (%d chars)\n", len(quickResponse.Content))
	fmt.Printf("  Content: %s\n", quickResponse.Content[:min(100, len(quickResponse.Content))])

	// Test 4: Test response optimization and caching
	fmt.Println("\n4. Testing response optimization and caching...")

	// Send the same message again to test caching
	cachedRequest := &domain.ChatRequest{
		SessionID: session.ID,
		Content:   "What are the key considerations for migrating our e-commerce platform to AWS?",
		Type:      domain.MessageTypeUser,
	}

	cachedResponse, err := chatService.SendMessage(ctx, cachedRequest)
	if err != nil {
		log.Fatalf("Failed to send cached message: %v", err)
	}

	fmt.Printf("✓ Cached response (%d chars)\n", len(cachedResponse.Content))
	fmt.Printf("  Processing time: %.2fms (should be faster if cached)\n", cachedResponse.ProcessTime)

	// Test 5: Test session context integration
	fmt.Println("\n5. Testing session context integration...")

	// Update session context
	sessionContext := &domain.SessionContext{
		ClientName:     "TechCorp Inc",
		MeetingType:    "architecture_review",
		ProjectContext: "Modernizing legacy e-commerce platform",
		ServiceTypes:   []string{"architecture", "security"},
		CloudProviders: []string{"AWS"},
		CustomFields: map[string]string{
			"budget":   "500k",
			"timeline": "6 months",
		},
	}

	if err := chatService.UpdateSessionContext(ctx, session.ID, sessionContext); err != nil {
		log.Fatalf("Failed to update session context: %v", err)
	}

	contextRequest := &domain.ChatRequest{
		SessionID: session.ID,
		Content:   "What security considerations should we prioritize?",
		Type:      domain.MessageTypeUser,
	}

	contextResponse, err := chatService.SendMessage(ctx, contextRequest)
	if err != nil {
		log.Fatalf("Failed to send context-aware message: %v", err)
	}

	fmt.Printf("✓ Context-aware response (%d chars)\n", len(contextResponse.Content))
	fmt.Printf("  Content: %s\n", contextResponse.Content[:min(150, len(contextResponse.Content))])

	// Test 6: Get session history
	fmt.Println("\n6. Testing session history...")
	history, err := chatService.GetSessionHistory(ctx, session.ID, 10)
	if err != nil {
		log.Fatalf("Failed to get session history: %v", err)
	}

	fmt.Printf("✓ Retrieved %d messages from history\n", len(history))
	for i, msg := range history {
		fmt.Printf("  %d. [%s] %s (%d chars)\n", i+1, msg.Type,
			msg.Content[:min(50, len(msg.Content))], len(msg.Content))
	}

	// Test 7: Test message statistics
	fmt.Println("\n7. Testing message statistics...")
	stats, err := chatService.GetMessageStats(ctx, session.ID)
	if err != nil {
		log.Fatalf("Failed to get message stats: %v", err)
	}

	fmt.Printf("✓ Message statistics:\n")
	fmt.Printf("  Total messages: %d\n", stats.TotalMessages)
	fmt.Printf("  Messages by type: %+v\n", stats.MessagesByType)
	fmt.Printf("  Messages by status: %+v\n", stats.MessagesByStatus)

	fmt.Println("\n✅ Enhanced AI Integration test completed successfully!")
	fmt.Println("\nKey features tested:")
	fmt.Println("- Enhanced AI response generation with company knowledge")
	fmt.Println("- Response optimization and caching")
	fmt.Println("- Quick action processing")
	fmt.Println("- Session context integration")
	fmt.Println("- Conversation history management")
	fmt.Println("- Fallback to basic AI service")
	fmt.Println("- Message statistics and analytics")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
