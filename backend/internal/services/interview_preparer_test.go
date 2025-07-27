package services

import (
	"context"
	"testing"
	"time"

	"github.com/cloud-consulting/backend/internal/domain"
	"github.com/cloud-consulting/backend/internal/interfaces"
)

// MockBedrockService for testing
type mockBedrockService struct{}

func (m *mockBedrockService) GenerateText(ctx context.Context, prompt string, options *interfaces.BedrockOptions) (*interfaces.BedrockResponse, error) {
	return &interfaces.BedrockResponse{
		Content: `TITLE: Test Interview Guide
OBJECTIVE: Test objective
ESTIMATED DURATION: 60 minutes
PREPARATION NOTES:
- Test note 1
- Test note 2
SECTION 1: TEST SECTION
Objective: Test section objective
Expected Duration: 30 minutes
Questions:
1. [MUST-ASK] Test question 1? (Type: business, Expected: test answer)
2. [SHOULD-ASK] Test question 2? (Type: technical, Expected: test answer)
FOLLOW-UP ACTIONS:
- Test action 1
- Test action 2`,
		Usage: interfaces.BedrockUsage{
			InputTokens:  100,
			OutputTokens: 200,
		},
	}, nil
}

func (m *mockBedrockService) GetModelInfo() interfaces.BedrockModelInfo {
	return interfaces.BedrockModelInfo{
		ModelID:     "test-model",
		ModelName:   "Test Model",
		Provider:    "Test",
		MaxTokens:   4000,
		IsAvailable: true,
	}
}

func (m *mockBedrockService) IsHealthy() bool {
	return true
}

// MockPromptArchitect for testing
type mockPromptArchitect struct{}

func (m *mockPromptArchitect) BuildReportPrompt(ctx context.Context, inquiry *domain.Inquiry, options *interfaces.PromptOptions) (string, error) {
	return "test prompt", nil
}

func (m *mockPromptArchitect) BuildInterviewPrompt(ctx context.Context, inquiry *domain.Inquiry) (string, error) {
	return "test interview prompt", nil
}

func (m *mockPromptArchitect) BuildRiskAssessmentPrompt(ctx context.Context, inquiry *domain.Inquiry) (string, error) {
	return "test risk prompt", nil
}

func (m *mockPromptArchitect) BuildCompetitiveAnalysisPrompt(ctx context.Context, inquiry *domain.Inquiry) (string, error) {
	return "test competitive prompt", nil
}

func (m *mockPromptArchitect) ValidatePrompt(prompt string) error {
	return nil
}

func (m *mockPromptArchitect) GetTemplate(templateName string) (*interfaces.PromptTemplate, error) {
	return &interfaces.PromptTemplate{Name: templateName}, nil
}

func (m *mockPromptArchitect) RegisterTemplate(template *interfaces.PromptTemplate) error {
	return nil
}

func (m *mockPromptArchitect) ListTemplates() []string {
	return []string{"test"}
}

func TestInterviewPreparerService_GenerateInterviewGuide(t *testing.T) {
	// Setup
	bedrockService := &mockBedrockService{}
	promptArchitect := &mockPromptArchitect{}
	service := NewInterviewPreparerService(bedrockService, promptArchitect)

	inquiry := &domain.Inquiry{
		ID:       "test-id",
		Name:     "Test User",
		Email:    "test@example.com",
		Company:  "Test Company",
		Services: []string{"migration"},
		Message:  "Test message",
		CreatedAt: time.Now(),
	}

	ctx := context.Background()

	// Execute
	guide, err := service.GenerateInterviewGuide(ctx, inquiry)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if guide == nil {
		t.Fatal("Expected interview guide, got nil")
	}

	if guide.ID == "" {
		t.Error("Expected guide ID to be set")
	}

	if guide.InquiryID != inquiry.ID {
		t.Errorf("Expected inquiry ID %s, got %s", inquiry.ID, guide.InquiryID)
	}

	if guide.Title == "" {
		t.Error("Expected guide title to be set")
	}

	if len(guide.Sections) == 0 {
		t.Error("Expected at least one section")
	}
}

func TestInterviewPreparerService_GenerateQuestionSet(t *testing.T) {
	// Setup
	bedrockService := &mockBedrockService{}
	promptArchitect := &mockPromptArchitect{}
	service := NewInterviewPreparerService(bedrockService, promptArchitect)

	ctx := context.Background()

	// Execute
	questionSet, err := service.GenerateQuestionSet(ctx, "migration", "healthcare")

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if questionSet == nil {
		t.Fatal("Expected question set, got nil")
	}

	if questionSet.ID == "" {
		t.Error("Expected question set ID to be set")
	}

	if questionSet.Category != "migration" {
		t.Errorf("Expected category 'migration', got %s", questionSet.Category)
	}

	if questionSet.Industry != "healthcare" {
		t.Errorf("Expected industry 'healthcare', got %s", questionSet.Industry)
	}
}

func TestInterviewPreparerService_GenerateDiscoveryChecklist(t *testing.T) {
	// Setup
	bedrockService := &mockBedrockService{}
	promptArchitect := &mockPromptArchitect{}
	service := NewInterviewPreparerService(bedrockService, promptArchitect)

	ctx := context.Background()

	// Execute
	checklist, err := service.GenerateDiscoveryChecklist(ctx, "migration")

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if checklist == nil {
		t.Fatal("Expected discovery checklist, got nil")
	}

	if checklist.ServiceType != "migration" {
		t.Errorf("Expected service type 'migration', got %s", checklist.ServiceType)
	}
}

func TestInterviewPreparerService_GenerateFollowUpQuestions(t *testing.T) {
	// Setup
	bedrockService := &mockBedrockService{}
	promptArchitect := &mockPromptArchitect{}
	service := NewInterviewPreparerService(bedrockService, promptArchitect)

	responses := []interfaces.InterviewResponse{
		{
			QuestionID: "q1",
			Response:   "Test response",
			Notes:      "Test notes",
			Confidence: "high",
		},
	}

	ctx := context.Background()

	// Execute
	questions, err := service.GenerateFollowUpQuestions(ctx, responses)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if questions == nil {
		t.Fatal("Expected follow-up questions, got nil")
	}
}

func TestInterviewPreparerService_InferIndustryFromCompany(t *testing.T) {
	// Setup
	bedrockService := &mockBedrockService{}
	promptArchitect := &mockPromptArchitect{}
	service := NewInterviewPreparerService(bedrockService, promptArchitect)

	// Test cases
	testCases := []struct {
		company  string
		expected string
	}{
		{"Regional Hospital", "healthcare"},
		{"Tech Solutions Inc", "technology"},
		{"First National Bank", "financial"},
		{"Manufacturing Corp", "manufacturing"},
		{"Retail Store", "retail"},
		{"Generic Company", "general"},
	}

	for _, tc := range testCases {
		result := service.inferIndustryFromCompany(tc.company)
		if result != tc.expected {
			t.Errorf("For company '%s', expected industry '%s', got '%s'", tc.company, tc.expected, result)
		}
	}
}

func TestInterviewPreparerService_InferQuestionCategory(t *testing.T) {
	// Setup
	bedrockService := &mockBedrockService{}
	promptArchitect := &mockPromptArchitect{}
	service := NewInterviewPreparerService(bedrockService, promptArchitect)

	// Test cases
	testCases := []struct {
		question string
		expected string
	}{
		{"What is your budget for this project?", "business"},
		{"What security measures do you have in place?", "security"},
		{"Describe your current server infrastructure", "infrastructure"},
		{"How do you plan to migrate your data?", "migration"},
		{"What are your general requirements?", "general"},
	}

	for _, tc := range testCases {
		result := service.inferQuestionCategory(tc.question)
		if result != tc.expected {
			t.Errorf("For question '%s', expected category '%s', got '%s'", tc.question, tc.expected, result)
		}
	}
}