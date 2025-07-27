package services

import (
	"context"
	"testing"
	"time"

	"github.com/cloud-consulting/backend/internal/domain"
	"github.com/cloud-consulting/backend/internal/interfaces"
	"github.com/stretchr/testify/assert"
)

// MockBedrockService for testing
type MockBedrockService struct{}

func (m *MockBedrockService) GenerateText(ctx context.Context, prompt string, options *interfaces.BedrockOptions) (*interfaces.BedrockResponse, error) {
	return &interfaces.BedrockResponse{
		Content: "Mock generated content",
		Usage: interfaces.BedrockUsage{
			InputTokens:  50,
			OutputTokens: 50,
		},
		Metadata: map[string]string{
			"model": "mock-model",
		},
	}, nil
}

func (m *MockBedrockService) GetModelInfo() interfaces.BedrockModelInfo {
	return interfaces.BedrockModelInfo{
		ModelID:     "mock-model",
		ModelName:   "Mock Model",
		Provider:    "mock",
		MaxTokens:   4000,
		IsAvailable: true,
	}
}

func (m *MockBedrockService) IsHealthy() bool {
	return true
}

func TestNewRoadmapGeneratorService(t *testing.T) {
	mockBedrock := &MockBedrockService{}
	service := NewRoadmapGeneratorService(mockBedrock)
	
	assert.NotNil(t, service)
	assert.NotNil(t, service.bedrockService)
	assert.NotNil(t, service.templates)
	
	// Check that default templates are loaded
	assert.Contains(t, service.templates, "migration")
	assert.Contains(t, service.templates, "optimization")
	assert.Contains(t, service.templates, "assessment")
	assert.Contains(t, service.templates, "architecture")
}

func TestGenerateImplementationRoadmap(t *testing.T) {
	mockBedrock := &MockBedrockService{}
	service := NewRoadmapGeneratorService(mockBedrock)
	ctx := context.Background()

	tests := []struct {
		name        string
		inquiry     *domain.Inquiry
		report      *domain.Report
		expectError bool
	}{
		{
			name: "successful migration roadmap",
			inquiry: &domain.Inquiry{
				ID:       "test-1",
				Company:  "TestCorp",
				Services: []string{"migration"},
				Message:  "Need to migrate to AWS",
				CreatedAt: time.Now(),
			},
			report:      nil,
			expectError: false,
		},
		{
			name: "successful optimization roadmap",
			inquiry: &domain.Inquiry{
				ID:       "test-2",
				Company:  "TestCorp",
				Services: []string{"optimization"},
				Message:  "Need to optimize costs",
				CreatedAt: time.Now(),
			},
			report:      nil,
			expectError: false,
		},
		{
			name:        "nil inquiry",
			inquiry:     nil,
			report:      nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			roadmap, err := service.GenerateImplementationRoadmap(ctx, tt.inquiry, tt.report)
			
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, roadmap)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, roadmap)
				assert.NotEmpty(t, roadmap.ID)
				assert.NotEmpty(t, roadmap.Title)
				assert.NotEmpty(t, roadmap.ProjectType)
				assert.NotEmpty(t, roadmap.TotalDuration)
				assert.NotEmpty(t, roadmap.EstimatedCost)
				assert.NotEmpty(t, roadmap.Phases)
				assert.Equal(t, tt.inquiry.ID, roadmap.InquiryID)
			}
		})
	}
}

func TestGeneratePhases(t *testing.T) {
	mockBedrock := &MockBedrockService{}
	service := NewRoadmapGeneratorService(mockBedrock)
	ctx := context.Background()

	requirements := []string{"migration", "security"}
	constraints := &interfaces.ProjectConstraints{
		Budget:         "medium",
		Timeline:       "standard",
		TeamSize:       5,
		RiskTolerance:  "medium",
		CloudProviders: []string{"AWS"},
	}

	phases, err := service.GeneratePhases(ctx, requirements, constraints)
	
	assert.NoError(t, err)
	assert.NotEmpty(t, phases)
	
	for _, phase := range phases {
		assert.NotEmpty(t, phase.ID)
		assert.NotEmpty(t, phase.Name)
		assert.NotEmpty(t, phase.Description)
		assert.NotEmpty(t, phase.Duration)
		assert.NotEmpty(t, phase.Priority)
		assert.NotEmpty(t, phase.RiskLevel)
	}
}

func TestEstimateResources(t *testing.T) {
	mockBedrock := &MockBedrockService{}
	service := NewRoadmapGeneratorService(mockBedrock)
	ctx := context.Background()

	// Create test phases with tasks
	phases := []interfaces.RoadmapPhase{
		{
			ID:   "phase-1",
			Name: "Test Phase 1",
			Tasks: []interfaces.Task{
				{
					ID:             "task-1",
					Name:           "Test Task 1",
					EstimatedHours: 40,
					SkillsRequired: []string{"Cloud Architect"},
				},
				{
					ID:             "task-2",
					Name:           "Test Task 2",
					EstimatedHours: 20,
					SkillsRequired: []string{"DevOps Engineer"},
				},
			},
			EstimatedCost: "$10000.00",
		},
		{
			ID:   "phase-2",
			Name: "Test Phase 2",
			Tasks: []interfaces.Task{
				{
					ID:             "task-3",
					Name:           "Test Task 3",
					EstimatedHours: 30,
					SkillsRequired: []string{"Cloud Architect"},
				},
			},
			EstimatedCost: "$15000.00",
		},
	}

	estimate, err := service.EstimateResources(ctx, phases)
	
	assert.NoError(t, err)
	assert.NotNil(t, estimate)
	assert.Equal(t, 90, estimate.TotalHours) // 40 + 20 + 30
	assert.NotEmpty(t, estimate.TotalCost)
	assert.Contains(t, estimate.RoleBreakdown, "Cloud Architect")
	assert.Contains(t, estimate.RoleBreakdown, "DevOps Engineer")
	assert.Equal(t, 70, estimate.RoleBreakdown["Cloud Architect"]) // 40 + 30
	assert.Equal(t, 20, estimate.RoleBreakdown["DevOps Engineer"])
}

func TestCalculateDependencies(t *testing.T) {
	mockBedrock := &MockBedrockService{}
	service := NewRoadmapGeneratorService(mockBedrock)
	ctx := context.Background()

	phases := []interfaces.RoadmapPhase{
		{
			ID:   "phase-1",
			Name: "Phase 1",
			Tasks: []interfaces.Task{
				{
					ID:           "task-1",
					Name:         "Task 1",
					Dependencies: []string{},
				},
			},
		},
		{
			ID:   "phase-2",
			Name: "Phase 2",
			Tasks: []interfaces.Task{
				{
					ID:           "task-2",
					Name:         "Task 2",
					Dependencies: []string{"task-1"},
				},
			},
		},
	}

	dependencies, err := service.CalculateDependencies(ctx, phases)
	
	assert.NoError(t, err)
	assert.NotEmpty(t, dependencies)
	
	// Should have at least one phase dependency
	foundPhaseDep := false
	foundTaskDep := false
	
	for _, dep := range dependencies {
		assert.NotEmpty(t, dep.ID)
		assert.NotEmpty(t, dep.FromID)
		assert.NotEmpty(t, dep.ToID)
		assert.NotEmpty(t, dep.Type)
		assert.NotEmpty(t, dep.Description)
		
		if dep.FromID == "phase-1" && dep.ToID == "phase-2" {
			foundPhaseDep = true
		}
		if dep.FromID == "task-1" && dep.ToID == "task-2" {
			foundTaskDep = true
		}
	}
	
	assert.True(t, foundPhaseDep, "Should have phase dependency")
	assert.True(t, foundTaskDep, "Should have task dependency")
}

func TestGenerateMilestones(t *testing.T) {
	mockBedrock := &MockBedrockService{}
	service := NewRoadmapGeneratorService(mockBedrock)
	ctx := context.Background()

	phases := []interfaces.RoadmapPhase{
		{
			ID:   "phase-1",
			Name: "Test Phase 1",
			Milestones: []interfaces.Milestone{
				{
					ID:   "milestone-1",
					Name: "Phase Milestone",
				},
			},
		},
		{
			ID:   "phase-2",
			Name: "Test Phase 2",
		},
	}

	milestones, err := service.GenerateMilestones(ctx, phases)
	
	assert.NoError(t, err)
	assert.NotEmpty(t, milestones)
	
	// Should have phase completion milestones + existing milestones + project milestones
	assert.GreaterOrEqual(t, len(milestones), 4) // 2 phase completions + 1 existing + 2 project milestones
	
	foundKickoff := false
	foundGoLive := false
	
	for _, milestone := range milestones {
		assert.NotEmpty(t, milestone.ID)
		assert.NotEmpty(t, milestone.Name)
		assert.NotEmpty(t, milestone.Type)
		assert.NotEmpty(t, milestone.Importance)
		
		if milestone.Name == "Project Kickoff" {
			foundKickoff = true
		}
		if milestone.Name == "Go-Live" {
			foundGoLive = true
		}
	}
	
	assert.True(t, foundKickoff, "Should have project kickoff milestone")
	assert.True(t, foundGoLive, "Should have go-live milestone")
}

func TestValidateRoadmap(t *testing.T) {
	mockBedrock := &MockBedrockService{}
	service := NewRoadmapGeneratorService(mockBedrock)
	ctx := context.Background()

	tests := []struct {
		name     string
		roadmap  *interfaces.ImplementationRoadmap
		isValid  bool
		hasError bool
	}{
		{
			name: "valid roadmap",
			roadmap: &interfaces.ImplementationRoadmap{
				ID:    "test-roadmap",
				Title: "Test Roadmap",
				Phases: []interfaces.RoadmapPhase{
					{
						ID:   "phase-1",
						Name: "Test Phase",
						Tasks: []interfaces.Task{
							{
								ID:   "task-1",
								Name: "Test Task",
							},
						},
						Duration: "2 weeks",
					},
				},
				Dependencies: []interfaces.Dependency{},
			},
			isValid:  true,
			hasError: false,
		},
		{
			name: "roadmap without title",
			roadmap: &interfaces.ImplementationRoadmap{
				ID:    "test-roadmap",
				Title: "",
				Phases: []interfaces.RoadmapPhase{
					{
						ID:   "phase-1",
						Name: "Test Phase",
					},
				},
			},
			isValid:  false,
			hasError: false,
		},
		{
			name: "roadmap without phases",
			roadmap: &interfaces.ImplementationRoadmap{
				ID:     "test-roadmap",
				Title:  "Test Roadmap",
				Phases: []interfaces.RoadmapPhase{},
			},
			isValid:  false,
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validation, err := service.ValidateRoadmap(ctx, tt.roadmap)
			
			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, validation)
				assert.Equal(t, tt.isValid, validation.IsValid)
				assert.GreaterOrEqual(t, validation.QualityScore, 0.0)
				assert.LessOrEqual(t, validation.QualityScore, 1.0)
				assert.GreaterOrEqual(t, validation.CompletenessScore, 0.0)
				assert.LessOrEqual(t, validation.CompletenessScore, 1.0)
			}
		})
	}
}

func TestDetermineProjectType(t *testing.T) {
	mockBedrock := &MockBedrockService{}
	service := NewRoadmapGeneratorService(mockBedrock)

	tests := []struct {
		name        string
		inquiry     *domain.Inquiry
		expectedType string
	}{
		{
			name: "migration project",
			inquiry: &domain.Inquiry{
				Services: []string{"migration"},
				Message:  "We need to migrate to cloud",
			},
			expectedType: "migration",
		},
		{
			name: "optimization project",
			inquiry: &domain.Inquiry{
				Services: []string{"optimization"},
				Message:  "We need to optimize our costs",
			},
			expectedType: "optimization",
		},
		{
			name: "assessment project",
			inquiry: &domain.Inquiry{
				Services: []string{"assessment"},
				Message:  "We need an assessment of our infrastructure",
			},
			expectedType: "assessment",
		},
		{
			name: "architecture project",
			inquiry: &domain.Inquiry{
				Services: []string{"architecture"},
				Message:  "We need architecture review",
			},
			expectedType: "architecture",
		},
		{
			name: "general project",
			inquiry: &domain.Inquiry{
				Services: []string{"consulting"},
				Message:  "We need general consulting help",
			},
			expectedType: "general",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			projectType := service.determineProjectType(tt.inquiry)
			assert.Equal(t, tt.expectedType, projectType)
		})
	}
}

func TestExtractConstraints(t *testing.T) {
	mockBedrock := &MockBedrockService{}
	service := NewRoadmapGeneratorService(mockBedrock)

	tests := []struct {
		name        string
		inquiry     *domain.Inquiry
		expectedBudget string
		expectedTimeline string
		expectedProviders []string
	}{
		{
			name: "budget mentioned",
			inquiry: &domain.Inquiry{
				Message: "We have a tight budget for this project",
			},
			expectedBudget: "medium",
			expectedTimeline: "standard",
			expectedProviders: []string{"AWS"},
		},
		{
			name: "urgent timeline",
			inquiry: &domain.Inquiry{
				Message: "This is urgent and needs to be done ASAP",
			},
			expectedBudget: "medium",
			expectedTimeline: "aggressive",
			expectedProviders: []string{"AWS"},
		},
		{
			name: "multi-cloud",
			inquiry: &domain.Inquiry{
				Message: "We want to use AWS and Azure for this project",
			},
			expectedBudget: "medium",
			expectedTimeline: "standard",
			expectedProviders: []string{"AWS", "Azure"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			constraints := service.extractConstraints(tt.inquiry)
			assert.Equal(t, tt.expectedBudget, constraints.Budget)
			assert.Equal(t, tt.expectedTimeline, constraints.Timeline)
			assert.Equal(t, tt.expectedProviders, constraints.CloudProviders)
		})
	}
}

func TestExtractIndustryContext(t *testing.T) {
	mockBedrock := &MockBedrockService{}
	service := NewRoadmapGeneratorService(mockBedrock)

	tests := []struct {
		name        string
		inquiry     *domain.Inquiry
		expectedContext string
	}{
		{
			name: "healthcare industry",
			inquiry: &domain.Inquiry{
				Message: "We are a healthcare company and need HIPAA compliance",
			},
			expectedContext: "healthcare",
		},
		{
			name: "financial industry",
			inquiry: &domain.Inquiry{
				Message: "We are a banking institution",
			},
			expectedContext: "financial",
		},
		{
			name: "retail industry",
			inquiry: &domain.Inquiry{
				Message: "We run an ecommerce platform",
			},
			expectedContext: "retail",
		},
		{
			name: "general industry",
			inquiry: &domain.Inquiry{
				Message: "We are a technology company",
			},
			expectedContext: "general",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context := service.extractIndustryContext(tt.inquiry)
			assert.Equal(t, tt.expectedContext, context)
		})
	}
}

func TestCalculateTotalDuration(t *testing.T) {
	mockBedrock := &MockBedrockService{}
	service := NewRoadmapGeneratorService(mockBedrock)

	tests := []struct {
		name     string
		phases   []interfaces.RoadmapPhase
		expected string
	}{
		{
			name: "short project",
			phases: []interfaces.RoadmapPhase{
				{Duration: "2 weeks"},
				{Duration: "1 week"},
			},
			expected: "3 weeks",
		},
		{
			name: "medium project",
			phases: []interfaces.RoadmapPhase{
				{Duration: "2-4 weeks"},
				{Duration: "3-6 weeks"},
			},
			expected: "2 months 2 weeks",
		},
		{
			name: "long project",
			phases: []interfaces.RoadmapPhase{
				{Duration: "4 weeks"},
				{Duration: "4 weeks"},
				{Duration: "4 weeks"},
				{Duration: "4 weeks"},
			},
			expected: "4 months",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			duration := service.calculateTotalDuration(tt.phases)
			assert.Equal(t, tt.expected, duration)
		})
	}
}

func TestExtractWeeksFromDuration(t *testing.T) {
	mockBedrock := &MockBedrockService{}
	service := NewRoadmapGeneratorService(mockBedrock)

	tests := []struct {
		name     string
		duration string
		expected int
	}{
		{
			name:     "single week",
			duration: "1 week",
			expected: 1,
		},
		{
			name:     "multiple weeks",
			duration: "3 weeks",
			expected: 3,
		},
		{
			name:     "range weeks",
			duration: "2-4 weeks",
			expected: 4,
		},
		{
			name:     "no weeks mentioned",
			duration: "some time",
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			weeks := service.extractWeeksFromDuration(tt.duration)
			assert.Equal(t, tt.expected, weeks)
		})
	}
}

func TestParseCostString(t *testing.T) {
	mockBedrock := &MockBedrockService{}
	service := NewRoadmapGeneratorService(mockBedrock)

	tests := []struct {
		name     string
		costStr  string
		expected float64
	}{
		{
			name:     "simple cost",
			costStr:  "$1000.00",
			expected: 1000.00,
		},
		{
			name:     "cost with commas",
			costStr:  "$10,000.50",
			expected: 10000.50,
		},
		{
			name:     "invalid cost",
			costStr:  "invalid",
			expected: 0.0,
		},
		{
			name:     "empty cost",
			costStr:  "",
			expected: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cost := service.parseCostString(tt.costStr)
			assert.Equal(t, tt.expected, cost)
		})
	}
}

func TestCalculateQualityScore(t *testing.T) {
	mockBedrock := &MockBedrockService{}
	service := NewRoadmapGeneratorService(mockBedrock)

	tests := []struct {
		name     string
		roadmap  *interfaces.ImplementationRoadmap
		minScore float64
		maxScore float64
	}{
		{
			name: "high quality roadmap",
			roadmap: &interfaces.ImplementationRoadmap{
				Title:    "Complete Roadmap",
				Overview: "Detailed overview",
				Phases: []interfaces.RoadmapPhase{
					{
						Name:        "Phase 1",
						Description: "Detailed description",
						Tasks: []interfaces.Task{
							{Name: "Task 1"},
						},
						Deliverables: []interfaces.Deliverable{
							{Name: "Deliverable 1"},
						},
					},
				},
			},
			minScore: 0.8,
			maxScore: 1.0,
		},
		{
			name: "low quality roadmap",
			roadmap: &interfaces.ImplementationRoadmap{
				Title: "",
				Phases: []interfaces.RoadmapPhase{
					{
						Name: "",
					},
				},
			},
			minScore: 0.0,
			maxScore: 0.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := service.calculateQualityScore(tt.roadmap)
			assert.GreaterOrEqual(t, score, tt.minScore)
			assert.LessOrEqual(t, score, tt.maxScore)
		})
	}
}

func TestCalculateCompletenessScore(t *testing.T) {
	mockBedrock := &MockBedrockService{}
	service := NewRoadmapGeneratorService(mockBedrock)

	tests := []struct {
		name     string
		roadmap  *interfaces.ImplementationRoadmap
		minScore float64
		maxScore float64
	}{
		{
			name: "complete roadmap",
			roadmap: &interfaces.ImplementationRoadmap{
				TotalDuration:  "6 months",
				EstimatedCost:  "$100000",
				Dependencies:   []interfaces.Dependency{{ID: "dep1"}},
				Risks:          []string{"Risk 1"},
				SuccessMetrics: []string{"Metric 1"},
				Phases: []interfaces.RoadmapPhase{
					{
						Duration:      "2 weeks",
						EstimatedCost: "$10000",
						Tasks:         []interfaces.Task{{Name: "Task 1"}},
						Deliverables:  []interfaces.Deliverable{{Name: "Deliverable 1"}},
						Milestones:    []interfaces.Milestone{{Name: "Milestone 1"}},
					},
				},
			},
			minScore: 0.8,
			maxScore: 1.0,
		},
		{
			name: "incomplete roadmap",
			roadmap: &interfaces.ImplementationRoadmap{
				Phases: []interfaces.RoadmapPhase{
					{
						Name: "Phase 1",
					},
				},
			},
			minScore: 0.0,
			maxScore: 0.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := service.calculateCompletenessScore(tt.roadmap)
			assert.GreaterOrEqual(t, score, tt.minScore)
			assert.LessOrEqual(t, score, tt.maxScore)
		})
	}
}