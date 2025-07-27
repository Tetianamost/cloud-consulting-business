package services

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/cloud-consulting/backend/internal/domain"
	"github.com/cloud-consulting/backend/internal/interfaces"
)

// MockKnowledgeBase is a mock implementation of KnowledgeBase interface
type MockKnowledgeBase struct {
	mock.Mock
}

func (m *MockKnowledgeBase) GetCloudServiceInfo(provider, service string) (*interfaces.CloudServiceInfo, error) {
	args := m.Called(provider, service)
	return args.Get(0).(*interfaces.CloudServiceInfo), args.Error(1)
}

func (m *MockKnowledgeBase) GetBestPractices(category, provider string) ([]*interfaces.BestPractice, error) {
	args := m.Called(category, provider)
	return args.Get(0).([]*interfaces.BestPractice), args.Error(1)
}

func (m *MockKnowledgeBase) GetComplianceRequirements(industry string) ([]*interfaces.ComplianceRequirement, error) {
	args := m.Called(industry)
	return args.Get(0).([]*interfaces.ComplianceRequirement), args.Error(1)
}

func (m *MockKnowledgeBase) GetArchitecturalPatterns(useCase, provider string) ([]*interfaces.ArchitecturalPattern, error) {
	args := m.Called(useCase, provider)
	return args.Get(0).([]*interfaces.ArchitecturalPattern), args.Error(1)
}

func (m *MockKnowledgeBase) GetDocumentationLinks(provider, topic string) ([]*interfaces.DocumentationLink, error) {
	args := m.Called(provider, topic)
	return args.Get(0).([]*interfaces.DocumentationLink), args.Error(1)
}

func (m *MockKnowledgeBase) UpdateKnowledgeBase(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockKnowledgeBase) SearchServices(ctx context.Context, query string, providers []string) ([]*interfaces.CloudServiceInfo, error) {
	args := m.Called(ctx, query, providers)
	return args.Get(0).([]*interfaces.CloudServiceInfo), args.Error(1)
}

func (m *MockKnowledgeBase) GetServiceAlternatives(provider, service string) (map[string]string, error) {
	args := m.Called(provider, service)
	return args.Get(0).(map[string]string), args.Error(1)
}

func (m *MockKnowledgeBase) IsHealthy() bool {
	args := m.Called()
	return args.Bool(0)
}

// MockDocumentationLibrary is a mock implementation of DocumentationLibrary interface
type MockDocumentationLibrary struct {
	mock.Mock
}

func (m *MockDocumentationLibrary) GetDocumentationLinks(ctx context.Context, provider, topic string) ([]*interfaces.DocumentationLink, error) {
	args := m.Called(ctx, provider, topic)
	return args.Get(0).([]*interfaces.DocumentationLink), args.Error(1)
}

func (m *MockDocumentationLibrary) ValidateLinks(ctx context.Context, links []*interfaces.DocumentationLink) ([]*interfaces.LinkValidation, error) {
	args := m.Called(ctx, links)
	return args.Get(0).([]*interfaces.LinkValidation), args.Error(1)
}

func (m *MockDocumentationLibrary) UpdateDocumentationIndex(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockDocumentationLibrary) SearchDocumentation(ctx context.Context, query string, providers []string) ([]*interfaces.DocumentationLink, error) {
	args := m.Called(ctx, query, providers)
	return args.Get(0).([]*interfaces.DocumentationLink), args.Error(1)
}

func (m *MockDocumentationLibrary) AddDocumentationLink(ctx context.Context, link *interfaces.DocumentationLink) error {
	args := m.Called(ctx, link)
	return args.Error(0)
}

func (m *MockDocumentationLibrary) RemoveDocumentationLink(ctx context.Context, linkID string) error {
	args := m.Called(ctx, linkID)
	return args.Error(0)
}

func (m *MockDocumentationLibrary) GetLinksByCategory(ctx context.Context, category string) ([]*interfaces.DocumentationLink, error) {
	args := m.Called(ctx, category)
	return args.Get(0).([]*interfaces.DocumentationLink), args.Error(1)
}

func (m *MockDocumentationLibrary) GetLinksByProvider(ctx context.Context, provider string) ([]*interfaces.DocumentationLink, error) {
	args := m.Called(ctx, provider)
	return args.Get(0).([]*interfaces.DocumentationLink), args.Error(1)
}

func (m *MockDocumentationLibrary) GetLinksByType(ctx context.Context, linkType string) ([]*interfaces.DocumentationLink, error) {
	args := m.Called(ctx, linkType)
	return args.Get(0).([]*interfaces.DocumentationLink), args.Error(1)
}

func (m *MockDocumentationLibrary) GetLinkValidationStatus(ctx context.Context, linkID string) (*interfaces.LinkValidation, error) {
	args := m.Called(ctx, linkID)
	return args.Get(0).(*interfaces.LinkValidation), args.Error(1)
}

func (m *MockDocumentationLibrary) IsHealthy() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockDocumentationLibrary) GetStats() *interfaces.DocumentationLibraryStats {
	args := m.Called()
	return args.Get(0).(*interfaces.DocumentationLibraryStats)
}

func TestNewRiskAssessorService(t *testing.T) {
	mockKB := &MockKnowledgeBase{}
	mockDocLib := &MockDocumentationLibrary{}

	service := NewRiskAssessorService(mockKB, mockDocLib)

	assert.NotNil(t, service)
	assert.IsType(t, &RiskAssessorService{}, service)
}

func TestRiskAssessorService_AssessRisks(t *testing.T) {
	mockKB := &MockKnowledgeBase{}
	mockDocLib := &MockDocumentationLibrary{}
	service := NewRiskAssessorService(mockKB, mockDocLib)

	ctx := context.Background()
	inquiry := &domain.Inquiry{
		ID:      "test-inquiry-1",
		Company: "Test Healthcare Corp",
		Message: "Need to migrate patient data to AWS cloud with HIPAA compliance",
		Services: []string{"migration"},
	}

	solution := &interfaces.ProposedSolution{
		ID:        "test-solution-1",
		InquiryID: inquiry.ID,
		CloudProviders: []string{"AWS"},
		Services: []interfaces.CloudService{
			{
				Provider:     "AWS",
				ServiceName:  "RDS",
				ServiceType:  "database",
				CriticalPath: true,
			},
		},
		Architecture: &interfaces.Architecture{
			ID:               "test-arch-1",
			Type:             "microservices",
			HighAvailability: false,
			DisasterRecovery: false,
			DataStorage: []interfaces.DataStorageComponent{
				{
					Type:             "relational",
					Provider:         "AWS",
					ServiceName:      "RDS",
					DataType:         "personal",
					SensitivityLevel: "high",
					BackupStrategy:   "",
				},
			},
			NetworkTopology: interfaces.NetworkTopology{
				SecurityGroups: []interfaces.SecurityGroup{
					{
						Name: "web-sg",
						Rules: []interfaces.FirewallRule{
							{
								Direction: "inbound",
								Protocol:  "TCP",
								Port:      "80",
								Source:    "0.0.0.0/0",
								Action:    "allow",
							},
						},
					},
				},
			},
			SecurityLayers: []interfaces.SecurityLayer{},
		},
		DataFlow: []interfaces.DataFlowComponent{
			{
				Source:      "web-app",
				Destination: "database",
				DataType:    "personal",
				Volume:      "high",
				Encryption:  false,
			},
		},
		EstimatedCost: "",
		Timeline:      "",
	}

	assessment, err := service.AssessRisks(ctx, inquiry, solution)

	assert.NoError(t, err)
	assert.NotNil(t, assessment)
	assert.Equal(t, inquiry.ID, assessment.InquiryID)
	assert.NotEmpty(t, assessment.ID)

	// Should identify technical risks
	assert.NotEmpty(t, assessment.TechnicalRisks)
	
	// Should identify security risks
	assert.NotEmpty(t, assessment.SecurityRisks)
	
	// Should identify compliance risks (HIPAA for healthcare)
	assert.NotEmpty(t, assessment.ComplianceRisks)
	
	// Should identify business risks
	assert.NotEmpty(t, assessment.BusinessRisks)
	
	// Should generate mitigation strategies
	assert.NotEmpty(t, assessment.MitigationStrategies)
	
	// Should have recommended actions
	assert.NotEmpty(t, assessment.RecommendedActions)
	
	// Should calculate overall risk level
	assert.NotEmpty(t, assessment.OverallRiskLevel)
	assert.Contains(t, []string{"low", "medium", "high", "critical"}, assessment.OverallRiskLevel)
}

func TestRiskAssessorService_IdentifySecurityRisks(t *testing.T) {
	mockKB := &MockKnowledgeBase{}
	mockDocLib := &MockDocumentationLibrary{}
	service := NewRiskAssessorService(mockKB, mockDocLib)

	ctx := context.Background()
	architecture := &interfaces.Architecture{
		DataStorage: []interfaces.DataStorageComponent{
			{
				ServiceName:      "user-db",
				SensitivityLevel: "high",
				BackupStrategy:   "",
			},
		},
		NetworkTopology: interfaces.NetworkTopology{
			SecurityGroups: []interfaces.SecurityGroup{
				{
					Name: "permissive-sg",
					Rules: []interfaces.FirewallRule{
						{
							Source: "0.0.0.0/0",
							Action: "allow",
						},
					},
				},
			},
		},
		SecurityLayers: []interfaces.SecurityLayer{},
	}

	risks, err := service.IdentifySecurityRisks(ctx, architecture)

	assert.NoError(t, err)
	assert.NotEmpty(t, risks)

	// Should identify data encryption risk
	found := false
	for _, risk := range risks {
		if risk.ThreatType == "data_exposure" {
			found = true
			assert.Equal(t, "security", risk.Category)
			assert.Equal(t, "high", risk.Impact)
			assert.True(t, risk.EncryptionRequired)
			break
		}
	}
	assert.True(t, found, "Should identify data encryption risk")

	// Should identify network security risk
	found = false
	for _, risk := range risks {
		if risk.ThreatType == "network_intrusion" {
			found = true
			assert.Equal(t, "security", risk.Category)
			assert.Contains(t, risk.AttackVectors, "port_scanning")
			break
		}
	}
	assert.True(t, found, "Should identify network security risk")

	// Should identify missing access control risk
	found = false
	for _, risk := range risks {
		if risk.ThreatType == "unauthorized_access" {
			found = true
			assert.Equal(t, "security", risk.Category)
			assert.Contains(t, risk.AffectedComponents, "entire_system")
			break
		}
	}
	assert.True(t, found, "Should identify access control risk")
}

func TestRiskAssessorService_EvaluateComplianceRisks(t *testing.T) {
	mockKB := &MockKnowledgeBase{}
	mockDocLib := &MockDocumentationLibrary{}
	service := NewRiskAssessorService(mockKB, mockDocLib)

	ctx := context.Background()
	industry := "healthcare"
	solution := &interfaces.ProposedSolution{
		Architecture: &interfaces.Architecture{
			DataStorage: []interfaces.DataStorageComponent{
				{
					ServiceName: "patient-db",
					DataType:    "health",
				},
			},
		},
		DataFlow: []interfaces.DataFlowComponent{
			{
				DataType:   "payment",
				Encryption: false,
			},
		},
	}

	risks, err := service.EvaluateComplianceRisks(ctx, industry, solution)

	assert.NoError(t, err)
	assert.NotEmpty(t, risks)

	// Should identify HIPAA risks for healthcare industry
	found := false
	for _, risk := range risks {
		if risk.Framework == "HIPAA" {
			found = true
			assert.Equal(t, "compliance", risk.Category)
			assert.Equal(t, "US", risk.Jurisdiction)
			assert.NotEmpty(t, risk.AuditRequirements)
			break
		}
	}
	assert.True(t, found, "Should identify HIPAA compliance risk")
}

func TestRiskAssessorService_GenerateMitigationStrategies(t *testing.T) {
	mockKB := &MockKnowledgeBase{}
	mockDocLib := &MockDocumentationLibrary{}
	service := NewRiskAssessorService(mockKB, mockDocLib)

	ctx := context.Background()
	risks := []*interfaces.Risk{
		{
			ID:       "risk-1",
			Category: "technical",
			Title:    "Single Point of Failure - No High Availability",
			Impact:   "high",
			RiskScore: 12,
		},
		{
			ID:       "risk-2",
			Category: "security",
			Title:    "Data Encryption Risk",
			Impact:   "high",
			RiskScore: 12,
		},
		{
			ID:       "risk-3",
			Category: "business",
			Title:    "Cost Estimation Risk",
			Impact:   "medium",
			RiskScore: 6,
		},
	}

	strategies, err := service.GenerateMitigationStrategies(ctx, risks)

	assert.NoError(t, err)
	assert.Len(t, strategies, 3)

	// Check technical mitigation strategy
	techStrategy := findStrategyByRiskID(strategies, "risk-1")
	assert.NotNil(t, techStrategy)
	assert.Contains(t, techStrategy.Strategy, "high availability")
	assert.NotEmpty(t, techStrategy.ImplementationSteps)
	assert.Equal(t, "critical", techStrategy.Priority)

	// Check security mitigation strategy
	secStrategy := findStrategyByRiskID(strategies, "risk-2")
	assert.NotNil(t, secStrategy)
	assert.Contains(t, secStrategy.Strategy, "encryption")
	assert.NotEmpty(t, secStrategy.ImplementationSteps)

	// Check business mitigation strategy
	bizStrategy := findStrategyByRiskID(strategies, "risk-3")
	assert.NotNil(t, bizStrategy)
	assert.Contains(t, bizStrategy.Strategy, "cost")
	assert.Equal(t, "medium", bizStrategy.Priority)
}

func TestRiskAssessorService_CalculateRiskScore(t *testing.T) {
	service := &RiskAssessorService{}

	tests := []struct {
		impact      string
		probability string
		expected    int
	}{
		{"critical", "high", 16},
		{"high", "high", 12},
		{"high", "medium", 9},
		{"medium", "medium", 6},
		{"low", "low", 2},
		{"unknown", "unknown", 1},
	}

	for _, test := range tests {
		score := service.calculateRiskScore(test.impact, test.probability)
		assert.Equal(t, test.expected, score, 
			"Risk score for impact=%s, probability=%s should be %d", 
			test.impact, test.probability, test.expected)
	}
}

func TestRiskAssessorService_ExtractIndustryFromInquiry(t *testing.T) {
	service := &RiskAssessorService{}

	tests := []struct {
		inquiry  *domain.Inquiry
		expected string
	}{
		{
			inquiry: &domain.Inquiry{
				Company: "Regional Hospital",
				Message: "Need healthcare cloud solution",
			},
			expected: "healthcare",
		},
		{
			inquiry: &domain.Inquiry{
				Company: "First National Bank",
				Message: "Banking system migration",
			},
			expected: "financial",
		},
		{
			inquiry: &domain.Inquiry{
				Company: "Retail Corp",
				Message: "Ecommerce platform setup",
			},
			expected: "retail",
		},
		{
			inquiry: &domain.Inquiry{
				Company: "Tech Startup",
				Message: "General cloud infrastructure",
			},
			expected: "general",
		},
	}

	for _, test := range tests {
		industry := service.extractIndustryFromInquiry(test.inquiry)
		assert.Equal(t, test.expected, industry)
	}
}

func TestRiskAssessorService_GetComplianceFrameworksForIndustry(t *testing.T) {
	service := &RiskAssessorService{}

	tests := []struct {
		industry string
		expected []string
	}{
		{"healthcare", []string{"HIPAA", "HITECH"}},
		{"financial", []string{"PCI-DSS", "SOX", "GLBA"}},
		{"retail", []string{"PCI-DSS", "GDPR"}},
		{"government", []string{"FedRAMP", "FISMA"}},
		{"general", []string{"GDPR", "SOC2"}},
	}

	for _, test := range tests {
		frameworks := service.getComplianceFrameworksForIndustry(test.industry)
		assert.Equal(t, test.expected, frameworks)
	}
}

func TestRiskAssessorService_CalculateOverallRiskLevel(t *testing.T) {
	service := &RiskAssessorService{}

	tests := []struct {
		risks    []*interfaces.Risk
		expected string
	}{
		{
			risks:    []*interfaces.Risk{},
			expected: "low",
		},
		{
			risks: []*interfaces.Risk{
				{Impact: "critical", RiskScore: 16},
			},
			expected: "critical",
		},
		{
			risks: []*interfaces.Risk{
				{Impact: "high", RiskScore: 12},
				{Impact: "high", RiskScore: 12},
				{Impact: "high", RiskScore: 12},
			},
			expected: "critical",
		},
		{
			risks: []*interfaces.Risk{
				{Impact: "medium", RiskScore: 6},
				{Impact: "medium", RiskScore: 6},
			},
			expected: "medium",
		},
		{
			risks: []*interfaces.Risk{
				{Impact: "low", RiskScore: 2},
			},
			expected: "low",
		},
	}

	for _, test := range tests {
		level := service.calculateOverallRiskLevel(test.risks)
		assert.Equal(t, test.expected, level)
	}
}

func TestRiskAssessorService_IsMigrationProject(t *testing.T) {
	service := &RiskAssessorService{}

	tests := []struct {
		inquiry  *domain.Inquiry
		expected bool
	}{
		{
			inquiry: &domain.Inquiry{
				Services: []string{"migration"},
			},
			expected: true,
		},
		{
			inquiry: &domain.Inquiry{
				Message: "Need to migrate our systems to cloud",
			},
			expected: true,
		},
		{
			inquiry: &domain.Inquiry{
				Services: []string{"assessment"},
				Message:  "General cloud assessment",
			},
			expected: false,
		},
	}

	for _, test := range tests {
		result := service.isMigrationProject(test.inquiry)
		assert.Equal(t, test.expected, result)
	}
}

// Helper function to find strategy by risk ID
func findStrategyByRiskID(strategies []*interfaces.MitigationStrategy, riskID string) *interfaces.MitigationStrategy {
	for _, strategy := range strategies {
		if strategy.RiskID == riskID {
			return strategy
		}
	}
	return nil
}

// Benchmark tests
func BenchmarkRiskAssessorService_AssessRisks(b *testing.B) {
	mockKB := &MockKnowledgeBase{}
	mockDocLib := &MockDocumentationLibrary{}
	service := NewRiskAssessorService(mockKB, mockDocLib)

	ctx := context.Background()
	inquiry := &domain.Inquiry{
		ID:      "bench-inquiry",
		Company: "Test Corp",
		Message: "Cloud migration project",
		Services: []string{"migration"},
	}

	solution := &interfaces.ProposedSolution{
		ID:        "bench-solution",
		InquiryID: inquiry.ID,
		CloudProviders: []string{"AWS"},
		Services: []interfaces.CloudService{
			{
				Provider:     "AWS",
				ServiceName:  "EC2",
				ServiceType:  "compute",
				CriticalPath: true,
			},
		},
		Architecture: &interfaces.Architecture{
			ID:               "bench-arch",
			Type:             "microservices",
			HighAvailability: true,
			DataStorage: []interfaces.DataStorageComponent{
				{
					ServiceName:      "main-db",
					SensitivityLevel: "medium",
					BackupStrategy:   "daily",
				},
			},
			NetworkTopology: interfaces.NetworkTopology{
				SecurityGroups: []interfaces.SecurityGroup{
					{
						Name: "app-sg",
						Rules: []interfaces.FirewallRule{
							{
								Source: "10.0.0.0/8",
								Action: "allow",
							},
						},
					},
				},
			},
		},
		EstimatedCost: "10000",
		Timeline:      "3 months",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := service.AssessRisks(ctx, inquiry, solution)
		if err != nil {
			b.Fatal(err)
		}
	}
}