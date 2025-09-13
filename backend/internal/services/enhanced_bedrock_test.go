package services

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/cloud-consulting/backend/internal/domain"
	"github.com/cloud-consulting/backend/internal/interfaces"
)

// MockBedrockServiceForEnhanced is a mock implementation for testing
type MockBedrockServiceForEnhanced struct {
	mock.Mock
}

func (m *MockBedrockServiceForEnhanced) GenerateText(ctx context.Context, prompt string, options *interfaces.BedrockOptions) (*interfaces.BedrockResponse, error) {
	args := m.Called(ctx, prompt, options)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*interfaces.BedrockResponse), args.Error(1)
}

func (m *MockBedrockServiceForEnhanced) GetModelInfo() interfaces.BedrockModelInfo {
	args := m.Called()
	return args.Get(0).(interfaces.BedrockModelInfo)
}

func (m *MockBedrockServiceForEnhanced) IsHealthy() bool {
	args := m.Called()
	return args.Bool(0)
}

// MockKnowledgeBaseForEnhanced is a mock implementation for testing
type MockKnowledgeBaseForEnhanced struct {
	mock.Mock
}

func (m *MockKnowledgeBaseForEnhanced) GetServiceOfferings(ctx context.Context) ([]*interfaces.ServiceOffering, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*interfaces.ServiceOffering), args.Error(1)
}

func (m *MockKnowledgeBaseForEnhanced) GetTeamExpertise(ctx context.Context) ([]*interfaces.TeamExpertise, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*interfaces.TeamExpertise), args.Error(1)
}

func (m *MockKnowledgeBaseForEnhanced) GetPastSolutions(ctx context.Context, serviceType, industry string) ([]*interfaces.PastSolution, error) {
	args := m.Called(ctx, serviceType, industry)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*interfaces.PastSolution), args.Error(1)
}

func (m *MockKnowledgeBaseForEnhanced) GetConsultingApproach(ctx context.Context, serviceType string) (*interfaces.ConsultingApproach, error) {
	args := m.Called(ctx, serviceType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*interfaces.ConsultingApproach), args.Error(1)
}

func (m *MockKnowledgeBaseForEnhanced) GetClientHistory(ctx context.Context, company string) ([]*interfaces.ClientEngagement, error) {
	args := m.Called(ctx, company)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*interfaces.ClientEngagement), args.Error(1)
}

func (m *MockKnowledgeBaseForEnhanced) GetBestPractices(ctx context.Context, category string) ([]*interfaces.BestPractice, error) {
	args := m.Called(ctx, category)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*interfaces.BestPractice), args.Error(1)
}

// Add other missing interface methods to satisfy the KnowledgeBase interface
func (m *MockKnowledgeBaseForEnhanced) GetServiceOffering(ctx context.Context, id string) (*interfaces.ServiceOffering, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*interfaces.ServiceOffering), args.Error(1)
}

func (m *MockKnowledgeBaseForEnhanced) GetPricingModels(ctx context.Context, serviceType string) ([]*interfaces.PricingModel, error) {
	args := m.Called(ctx, serviceType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*interfaces.PricingModel), args.Error(1)
}

func (m *MockKnowledgeBaseForEnhanced) GetConsultantSpecializations(ctx context.Context, consultantID string) ([]*interfaces.Specialization, error) {
	args := m.Called(ctx, consultantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*interfaces.Specialization), args.Error(1)
}

func (m *MockKnowledgeBaseForEnhanced) GetExpertiseByArea(ctx context.Context, area string) ([]*interfaces.TeamExpertise, error) {
	args := m.Called(ctx, area)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*interfaces.TeamExpertise), args.Error(1)
}

func (m *MockKnowledgeBaseForEnhanced) GetSimilarProjects(ctx context.Context, inquiry *domain.Inquiry) ([]*interfaces.ProjectPattern, error) {
	args := m.Called(ctx, inquiry)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*interfaces.ProjectPattern), args.Error(1)
}

func (m *MockKnowledgeBaseForEnhanced) GetMethodologyTemplates(ctx context.Context, serviceType string) ([]*interfaces.MethodologyTemplate, error) {
	args := m.Called(ctx, serviceType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*interfaces.MethodologyTemplate), args.Error(1)
}

func (m *MockKnowledgeBaseForEnhanced) GetDeliverableTemplates(ctx context.Context, serviceType string) ([]*interfaces.DeliverableTemplate, error) {
	args := m.Called(ctx, serviceType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*interfaces.DeliverableTemplate), args.Error(1)
}

func (m *MockKnowledgeBaseForEnhanced) UpdateKnowledgeBase(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockKnowledgeBaseForEnhanced) SearchKnowledge(ctx context.Context, query string, category string) ([]*interfaces.KnowledgeItem, error) {
	args := m.Called(ctx, query, category)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*interfaces.KnowledgeItem), args.Error(1)
}

func (m *MockKnowledgeBaseForEnhanced) GetKnowledgeStats(ctx context.Context) (*interfaces.KnowledgeStats, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*interfaces.KnowledgeStats), args.Error(1)
}

func (m *MockKnowledgeBaseForEnhanced) GetComplianceRequirements(ctx context.Context, framework string) ([]*interfaces.ComplianceRequirement, error) {
	args := m.Called(ctx, framework)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*interfaces.ComplianceRequirement), args.Error(1)
}

// MockCompanyKnowledgeIntegrationService for testing
type MockCompanyKnowledgeIntegrationService struct {
	mock.Mock
}

// Ensure it implements the interface
var _ CompanyKnowledgeIntegrator = (*MockCompanyKnowledgeIntegrationService)(nil)

func (m *MockCompanyKnowledgeIntegrationService) GenerateContextualPrompt(ctx context.Context, inquiry *domain.Inquiry, basePrompt string) (string, error) {
	args := m.Called(ctx, inquiry, basePrompt)
	return args.String(0), args.Error(1)
}

func (m *MockCompanyKnowledgeIntegrationService) GetRecommendationsForInquiry(ctx context.Context, inquiry *domain.Inquiry) (*InquiryRecommendations, error) {
	args := m.Called(ctx, inquiry)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*InquiryRecommendations), args.Error(1)
}

// Test helper functions
func createTestEnhancedBedrockService() (*EnhancedBedrockService, *MockBedrockServiceForEnhanced, *MockKnowledgeBaseForEnhanced, *MockCompanyKnowledgeIntegrationService) {
	mockBedrock := &MockBedrockServiceForEnhanced{}
	mockKB := &MockKnowledgeBaseForEnhanced{}
	mockCompanyInteg := &MockCompanyKnowledgeIntegrationService{}

	service := &EnhancedBedrockService{
		bedrockService:        mockBedrock,
		knowledgeBase:         mockKB,
		companyKnowledgeInteg: mockCompanyInteg,
	}

	return service, mockBedrock, mockKB, mockCompanyInteg
}

func createTestInquiry() *domain.Inquiry {
	return &domain.Inquiry{
		ID:        "test-inquiry-id",
		Name:      "John Doe",
		Company:   "Test Company",
		Message:   "We need help with AWS migration",
		Services:  []string{"Cloud Migration", "Architecture Review"},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func createTestBedrockOptions() *interfaces.BedrockOptions {
	return &interfaces.BedrockOptions{
		ModelID:     "anthropic.claude-3-sonnet-20240229-v1:0",
		MaxTokens:   1000,
		Temperature: 0.7,
		TopP:        0.9,
	}
}

// Test GenerateEnhancedResponse
func TestEnhancedBedrockService_GenerateEnhancedResponse_Success(t *testing.T) {
	service, mockBedrock, _, mockCompanyInteg := createTestEnhancedBedrockService()
	ctx := context.Background()

	inquiry := createTestInquiry()
	options := createTestBedrockOptions()

	enhancedPrompt := "Enhanced prompt with company knowledge"
	bedrockResponse := &interfaces.BedrockResponse{
		Content: "Enhanced AI response",
		Usage: interfaces.BedrockUsage{
			InputTokens:  100,
			OutputTokens: 200,
		},
	}

	// Mock company knowledge integration
	mockCompanyInteg.On("GenerateContextualPrompt", ctx, inquiry, mock.AnythingOfType("string")).Return(enhancedPrompt, nil)

	// Mock Bedrock service
	mockBedrock.On("GenerateText", ctx, enhancedPrompt, options).Return(bedrockResponse, nil)

	response, err := service.GenerateEnhancedResponse(ctx, inquiry, options)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "Enhanced AI response", response.Content)
	assert.Equal(t, 100, response.Usage.InputTokens)
	assert.Equal(t, 200, response.Usage.OutputTokens)

	mockCompanyInteg.AssertExpectations(t)
	mockBedrock.AssertExpectations(t)
}

func TestEnhancedBedrockService_GenerateEnhancedResponse_CompanyIntegrationFailure(t *testing.T) {
	service, mockBedrock, _, mockCompanyInteg := createTestEnhancedBedrockService()
	ctx := context.Background()

	inquiry := createTestInquiry()
	options := createTestBedrockOptions()

	bedrockResponse := &interfaces.BedrockResponse{
		Content: "Fallback AI response",
		Usage: interfaces.BedrockUsage{
			InputTokens:  80,
			OutputTokens: 150,
		},
	}

	// Mock company knowledge integration failure
	mockCompanyInteg.On("GenerateContextualPrompt", ctx, inquiry, mock.AnythingOfType("string")).Return("", assert.AnError)

	// Mock Bedrock service with fallback prompt
	mockBedrock.On("GenerateText", ctx, mock.AnythingOfType("string"), options).Return(bedrockResponse, nil)

	response, err := service.GenerateEnhancedResponse(ctx, inquiry, options)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "Fallback AI response", response.Content)

	mockCompanyInteg.AssertExpectations(t)
	mockBedrock.AssertExpectations(t)
}

func TestEnhancedBedrockService_GenerateEnhancedResponse_BedrockFailure(t *testing.T) {
	service, mockBedrock, _, mockCompanyInteg := createTestEnhancedBedrockService()
	ctx := context.Background()

	inquiry := createTestInquiry()
	options := createTestBedrockOptions()

	enhancedPrompt := "Enhanced prompt with company knowledge"

	// Mock company knowledge integration success
	mockCompanyInteg.On("GenerateContextualPrompt", ctx, inquiry, mock.AnythingOfType("string")).Return(enhancedPrompt, nil)

	// Mock Bedrock service failure
	mockBedrock.On("GenerateText", ctx, enhancedPrompt, options).Return(nil, assert.AnError)

	response, err := service.GenerateEnhancedResponse(ctx, inquiry, options)

	assert.Error(t, err)
	assert.Nil(t, response)

	mockCompanyInteg.AssertExpectations(t)
	mockBedrock.AssertExpectations(t)
}

// Test GenerateEnhancedResponseWithRecommendations
func TestEnhancedBedrockService_GenerateEnhancedResponseWithRecommendations_Success(t *testing.T) {
	service, mockBedrock, _, mockCompanyInteg := createTestEnhancedBedrockService()
	ctx := context.Background()

	inquiry := createTestInquiry()
	options := createTestBedrockOptions()

	enhancedPrompt := "Enhanced prompt with company knowledge"
	bedrockResponse := &interfaces.BedrockResponse{
		Content: "Enhanced AI response",
		Usage: interfaces.BedrockUsage{
			InputTokens:  100,
			OutputTokens: 200,
		},
	}

	recommendations := &InquiryRecommendations{
		InquiryID:   inquiry.ID,
		GeneratedAt: time.Now(),
		RecommendedServices: []*interfaces.ServiceOffering{
			{
				ID:          "aws-migration-hub",
				Name:        "AWS Migration Hub",
				Description: "Centralized migration tracking",
				Category:    "Migration",
				ServiceType: domain.ServiceTypeMigration,
			},
		},
	}

	// Mock company knowledge integration
	mockCompanyInteg.On("GenerateContextualPrompt", ctx, inquiry, mock.AnythingOfType("string")).Return(enhancedPrompt, nil)
	mockCompanyInteg.On("GetRecommendationsForInquiry", ctx, inquiry).Return(recommendations, nil)

	// Mock Bedrock service
	mockBedrock.On("GenerateText", ctx, enhancedPrompt, options).Return(bedrockResponse, nil)

	response, err := service.GenerateEnhancedResponseWithRecommendations(ctx, inquiry, options)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "Enhanced AI response", response.Content)
	assert.Equal(t, inquiry.ID, response.Recommendations.InquiryID)
	assert.Len(t, response.Recommendations.RecommendedServices, 1)
	assert.Equal(t, "AWS Migration Hub", response.Recommendations.RecommendedServices[0].Name)

	mockCompanyInteg.AssertExpectations(t)
	mockBedrock.AssertExpectations(t)
}

func TestEnhancedBedrockService_GenerateEnhancedResponseWithRecommendations_RecommendationsFailure(t *testing.T) {
	service, mockBedrock, _, mockCompanyInteg := createTestEnhancedBedrockService()
	ctx := context.Background()

	inquiry := createTestInquiry()
	options := createTestBedrockOptions()

	enhancedPrompt := "Enhanced prompt with company knowledge"
	bedrockResponse := &interfaces.BedrockResponse{
		Content: "Enhanced AI response",
		Usage: interfaces.BedrockUsage{
			InputTokens:  100,
			OutputTokens: 200,
		},
	}

	// Mock company knowledge integration
	mockCompanyInteg.On("GenerateContextualPrompt", ctx, inquiry, mock.AnythingOfType("string")).Return(enhancedPrompt, nil)
	mockCompanyInteg.On("GetRecommendationsForInquiry", ctx, inquiry).Return(nil, assert.AnError)

	// Mock Bedrock service
	mockBedrock.On("GenerateText", ctx, enhancedPrompt, options).Return(bedrockResponse, nil)

	response, err := service.GenerateEnhancedResponseWithRecommendations(ctx, inquiry, options)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "Enhanced AI response", response.Content)
	assert.Equal(t, inquiry.ID, response.Recommendations.InquiryID)
	assert.Empty(t, response.Recommendations.RecommendedServices)

	mockCompanyInteg.AssertExpectations(t)
	mockBedrock.AssertExpectations(t)
}

// Test buildEnhancedPrompt
func TestEnhancedBedrockService_BuildEnhancedPrompt_Success(t *testing.T) {
	service, _, mockKB, _ := createTestEnhancedBedrockService()
	ctx := context.Background()

	inquiry := createTestInquiry()

	// Mock knowledge base responses
	serviceOfferings := []*interfaces.ServiceOffering{
		{
			Name:            "Cloud Migration",
			Description:     "Comprehensive cloud migration services",
			Category:        "Migration",
			TypicalDuration: "3-6 months",
			TeamSize:        "3-5 consultants",
			KeyBenefits:     []string{"Reduced costs", "Improved scalability"},
			Deliverables:    []string{"Migration plan", "Implementation"},
		},
	}

	teamExpertise := []*interfaces.TeamExpertise{
		{
			ConsultantName:  "John Smith",
			Role:            "Senior Cloud Architect",
			ExperienceYears: 8,
			ExpertiseAreas:  []string{"AWS Migration", "Architecture Design"},
			CloudProviders:  []string{"AWS", "Azure"},
		},
	}

	pastSolutions := []*interfaces.PastSolution{
		{
			Title:            "E-commerce Migration",
			Industry:         "Retail",
			ProblemStatement: "Legacy infrastructure limitations",
			SolutionApproach: "Lift and shift with optimization",
			Technologies:     []string{"EC2", "RDS", "CloudFront"},
			TimeToValue:      "4 months",
			CostSavings:      50000,
		},
	}

	consultingApproach := &interfaces.ConsultingApproach{
		Name:              "Migration Methodology",
		Philosophy:        "Phased approach with minimal disruption",
		EngagementModel:   "Fixed scope with milestones",
		KeyPrinciples:     []string{"Risk mitigation", "Business continuity"},
		ClientInvolvement: "Active collaboration required",
		KnowledgeTransfer: "Comprehensive training provided",
	}

	clientHistory := []*interfaces.ClientEngagement{
		{
			ProjectName:        "Previous Migration",
			StartDate:          time.Now().AddDate(-1, 0, 0),
			Status:             "Completed",
			ClientSatisfaction: 9.2,
		},
	}

	// Mock all knowledge base calls
	mockKB.On("GetServiceOfferings", ctx).Return(serviceOfferings, nil)
	mockKB.On("GetTeamExpertise", ctx).Return(teamExpertise, nil)
	mockKB.On("GetPastSolutions", ctx, mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(pastSolutions, nil)
	mockKB.On("GetConsultingApproach", ctx, mock.AnythingOfType("string")).Return(consultingApproach, nil)
	mockKB.On("GetClientHistory", ctx, "Test Company").Return(clientHistory, nil)

	prompt, err := service.buildEnhancedPrompt(ctx, inquiry)

	assert.NoError(t, err)
	assert.NotEmpty(t, prompt)
	assert.Contains(t, prompt, "Cloud Migration")
	assert.Contains(t, prompt, "John Smith")
	assert.Contains(t, prompt, "E-commerce Migration")
	assert.Contains(t, prompt, "Migration Methodology")
	assert.Contains(t, prompt, "Previous Migration")
	assert.Contains(t, prompt, "Test Company")
	assert.Contains(t, prompt, "We need help with AWS migration")

	mockKB.AssertExpectations(t)
}

// Test helper methods
func TestEnhancedBedrockService_FilterRelevantOfferings(t *testing.T) {
	service, _, _, _ := createTestEnhancedBedrockService()

	offerings := []*interfaces.ServiceOffering{
		{
			Name:        "Cloud Migration",
			Description: "Comprehensive migration services",
			Category:    "Migration",
		},
		{
			Name:        "Security Assessment",
			Description: "Security review and recommendations",
			Category:    "Security",
		},
		{
			Name:        "Cost Optimization",
			Description: "Cost analysis and optimization",
			Category:    "Optimization",
		},
	}

	services := []string{"Cloud Migration", "Architecture Review"}

	relevant := service.filterRelevantOfferings(offerings, services)

	assert.Len(t, relevant, 1)
	assert.Equal(t, "Cloud Migration", relevant[0].Name)
}

func TestEnhancedBedrockService_FilterRelevantOfferings_NoMatches(t *testing.T) {
	service, _, _, _ := createTestEnhancedBedrockService()

	offerings := []*interfaces.ServiceOffering{
		{
			Name:        "Security Assessment",
			Description: "Security review and recommendations",
			Category:    "Security",
		},
		{
			Name:        "Cost Optimization",
			Description: "Cost analysis and optimization",
			Category:    "Optimization",
		},
	}

	services := []string{"Data Analytics"}

	relevant := service.filterRelevantOfferings(offerings, services)

	// Should return up to 3 offerings when no specific matches
	assert.Len(t, relevant, 2)
}

func TestEnhancedBedrockService_FilterRelevantExpertise(t *testing.T) {
	service, _, _, _ := createTestEnhancedBedrockService()

	expertise := []*interfaces.TeamExpertise{
		{
			ConsultantName: "John Smith",
			ExpertiseAreas: []string{"AWS Migration", "Architecture Design"},
		},
		{
			ConsultantName: "Jane Doe",
			ExpertiseAreas: []string{"Security", "Compliance"},
		},
		{
			ConsultantName: "Bob Johnson",
			ExpertiseAreas: []string{"Cost Optimization", "Performance"},
		},
	}

	services := []string{"Cloud Migration"}

	relevant := service.filterRelevantExpertise(expertise, services)

	assert.Len(t, relevant, 1)
	assert.Equal(t, "John Smith", relevant[0].ConsultantName)
}

func TestEnhancedBedrockService_InferIndustry(t *testing.T) {
	service, _, _, _ := createTestEnhancedBedrockService()

	tests := []struct {
		company  string
		expected string
	}{
		{"First National Bank", "Financial Services"},
		{"TechCorp Software", "Technology"},
		{"HealthCare Solutions", "Healthcare"},
		{"Retail Store Inc", "Retail"},
		{"Generic Company", ""},
	}

	for _, test := range tests {
		result := service.inferIndustry(test.company)
		assert.Equal(t, test.expected, result, "Company: %s", test.company)
	}
}

func TestEnhancedBedrockService_InferServiceType(t *testing.T) {
	service, _, _, _ := createTestEnhancedBedrockService()

	tests := []struct {
		services []string
		expected string
	}{
		{[]string{"Cloud Migration"}, string(domain.ReportTypeMigration)},
		{[]string{"Architecture Review"}, string(domain.ReportTypeArchitectureReview)},
		{[]string{"Assessment"}, string(domain.ReportTypeAssessment)},
		{[]string{"Optimization"}, string(domain.ReportTypeOptimization)},
		{[]string{"General Consulting"}, string(domain.ReportTypeGeneral)},
	}

	for _, test := range tests {
		result := service.inferServiceType(test.services)
		assert.Equal(t, test.expected, result, "Services: %v", test.services)
	}
}

// Test delegation methods
func TestEnhancedBedrockService_GenerateText_Delegation(t *testing.T) {
	service, mockBedrock, _, _ := createTestEnhancedBedrockService()
	ctx := context.Background()

	prompt := "Test prompt"
	options := createTestBedrockOptions()
	expectedResponse := &interfaces.BedrockResponse{
		Content: "Test response",
	}

	mockBedrock.On("GenerateText", ctx, prompt, options).Return(expectedResponse, nil)

	response, err := service.GenerateText(ctx, prompt, options)

	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, response)

	mockBedrock.AssertExpectations(t)
}

func TestEnhancedBedrockService_GetModelInfo_Delegation(t *testing.T) {
	service, mockBedrock, _, _ := createTestEnhancedBedrockService()

	expectedInfo := interfaces.BedrockModelInfo{
		ModelID:     "test-model",
		ModelName:   "Test Model",
		Provider:    "amazon",
		MaxTokens:   4000,
		IsAvailable: true,
	}

	mockBedrock.On("GetModelInfo").Return(expectedInfo)

	info := service.GetModelInfo()

	assert.Equal(t, expectedInfo, info)

	mockBedrock.AssertExpectations(t)
}

func TestEnhancedBedrockService_IsHealthy_Delegation(t *testing.T) {
	service, mockBedrock, _, _ := createTestEnhancedBedrockService()

	mockBedrock.On("IsHealthy").Return(true)

	healthy := service.IsHealthy()

	assert.True(t, healthy)

	mockBedrock.AssertExpectations(t)
}

// Test service matching logic
func TestEnhancedBedrockService_IsServiceMatch(t *testing.T) {
	service, _, _, _ := createTestEnhancedBedrockService()

	offering := &interfaces.ServiceOffering{
		Name:        "Cloud Migration Services",
		Description: "Comprehensive migration to AWS cloud",
		Category:    "Migration",
	}

	tests := []struct {
		service  string
		expected bool
	}{
		{"Cloud Migration", true},
		{"migration", true},
		{"AWS", true},
		{"Security", false},
		{"", false},
	}

	for _, test := range tests {
		result := service.isServiceMatch(offering, test.service)
		assert.Equal(t, test.expected, result, "Service: %s", test.service)
	}
}

func TestEnhancedBedrockService_IsExpertiseMatch(t *testing.T) {
	service, _, _, _ := createTestEnhancedBedrockService()

	expertise := &interfaces.TeamExpertise{
		ConsultantName: "John Smith",
		ExpertiseAreas: []string{"AWS Migration", "Cloud Architecture", "DevOps"},
	}

	tests := []struct {
		service  string
		expected bool
	}{
		{"Migration", true},
		{"Architecture", true},
		{"DevOps", true},
		{"Security", false},
		{"", false},
	}

	for _, test := range tests {
		result := service.isExpertiseMatch(expertise, test.service)
		assert.Equal(t, test.expected, result, "Service: %s", test.service)
	}
}
