package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/cloud-consulting/backend/internal/domain"
	"github.com/cloud-consulting/backend/internal/interfaces"
	"github.com/cloud-consulting/backend/internal/services"
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

func main() {
	fmt.Println("Testing Implementation Roadmap Generator...")
	
	// Create mock services
	mockBedrock := &MockBedrockService{}
	
	// Create roadmap generator
	roadmapGen := services.NewRoadmapGeneratorService(mockBedrock)
	
	// Test scenarios
	testScenarios := []struct {
		name    string
		inquiry *domain.Inquiry
		report  *domain.Report
	}{
		{
			name: "Cloud Migration Project",
			inquiry: &domain.Inquiry{
				ID:       "test-inquiry-1",
				Name:     "John Doe",
				Email:    "john@example.com",
				Company:  "TechCorp Inc",
				Services: []string{"migration"},
				Message:  "We need to migrate our on-premises infrastructure to AWS cloud. We have about 50 servers and need HIPAA compliance.",
				Status:   "new",
				Priority: "high",
				CreatedAt: time.Now(),
			},
			report: &domain.Report{
				ID:        "test-report-1",
				InquiryID: "test-inquiry-1",
				Type:      domain.ReportTypeMigration,
				Title:     "Cloud Migration Assessment",
				Content:   "Migration assessment content...",
				Status:    domain.ReportStatusGenerated,
				CreatedAt: time.Now(),
			},
		},
		{
			name: "Cloud Optimization Project",
			inquiry: &domain.Inquiry{
				ID:       "test-inquiry-2",
				Name:     "Jane Smith",
				Email:    "jane@example.com",
				Company:  "RetailCorp",
				Services: []string{"optimization"},
				Message:  "Our cloud costs are too high. We need to optimize our AWS spending and improve performance.",
				Status:   "new",
				Priority: "medium",
				CreatedAt: time.Now(),
			},
			report: &domain.Report{
				ID:        "test-report-2",
				InquiryID: "test-inquiry-2",
				Type:      domain.ReportTypeOptimization,
				Title:     "Cloud Optimization Analysis",
				Content:   "Optimization analysis content...",
				Status:    domain.ReportStatusGenerated,
				CreatedAt: time.Now(),
			},
		},
		{
			name: "Architecture Review Project",
			inquiry: &domain.Inquiry{
				ID:       "test-inquiry-3",
				Name:     "Bob Johnson",
				Email:    "bob@example.com",
				Company:  "FinanceCorp",
				Services: []string{"architecture_review"},
				Message:  "We need an architecture review for our multi-cloud setup. We're using AWS and Azure and want to ensure best practices.",
				Status:   "new",
				Priority: "high",
				CreatedAt: time.Now(),
			},
			report: &domain.Report{
				ID:        "test-report-3",
				InquiryID: "test-inquiry-3",
				Type:      domain.ReportTypeArchitectureReview,
				Title:     "Architecture Review Report",
				Content:   "Architecture review content...",
				Status:    domain.ReportStatusGenerated,
				CreatedAt: time.Now(),
			},
		},
	}
	
	ctx := context.Background()
	
	for _, scenario := range testScenarios {
		fmt.Printf("\n=== Testing %s ===\n", scenario.name)
		
		// Test roadmap generation
		roadmap, err := roadmapGen.GenerateImplementationRoadmap(ctx, scenario.inquiry, scenario.report)
		if err != nil {
			log.Printf("Error generating roadmap for %s: %v", scenario.name, err)
			continue
		}
		
		// Validate roadmap structure
		fmt.Printf("✓ Roadmap generated successfully\n")
		fmt.Printf("  - ID: %s\n", roadmap.ID)
		fmt.Printf("  - Title: %s\n", roadmap.Title)
		fmt.Printf("  - Project Type: %s\n", roadmap.ProjectType)
		fmt.Printf("  - Total Duration: %s\n", roadmap.TotalDuration)
		fmt.Printf("  - Estimated Cost: %s\n", roadmap.EstimatedCost)
		fmt.Printf("  - Number of Phases: %d\n", len(roadmap.Phases))
		fmt.Printf("  - Number of Dependencies: %d\n", len(roadmap.Dependencies))
		fmt.Printf("  - Number of Risks: %d\n", len(roadmap.Risks))
		fmt.Printf("  - Number of Success Metrics: %d\n", len(roadmap.SuccessMetrics))
		fmt.Printf("  - Cloud Providers: %v\n", roadmap.CloudProviders)
		fmt.Printf("  - Industry Context: %s\n", roadmap.IndustryContext)
		
		// Test phase details
		fmt.Printf("\n  Phases:\n")
		for i, phase := range roadmap.Phases {
			fmt.Printf("    %d. %s\n", i+1, phase.Name)
			fmt.Printf("       - Duration: %s\n", phase.Duration)
			fmt.Printf("       - Cost: %s\n", phase.EstimatedCost)
			fmt.Printf("       - Tasks: %d\n", len(phase.Tasks))
			fmt.Printf("       - Deliverables: %d\n", len(phase.Deliverables))
			fmt.Printf("       - Milestones: %d\n", len(phase.Milestones))
			fmt.Printf("       - Risk Level: %s\n", phase.RiskLevel)
			fmt.Printf("       - Priority: %s\n", phase.Priority)
			
			// Show first few tasks
			if len(phase.Tasks) > 0 {
				fmt.Printf("       - Sample Tasks:\n")
				for j, task := range phase.Tasks {
					if j >= 2 { // Show only first 2 tasks
						break
					}
					fmt.Printf("         * %s (%dh, %s priority)\n", task.Name, task.EstimatedHours, task.Priority)
				}
			}
		}
		
		// Test individual components
		fmt.Printf("\n  Testing individual components:\n")
		
		// Test phase generation
		constraints := &interfaces.ProjectConstraints{
			Budget:         "medium",
			Timeline:       "standard",
			TeamSize:       5,
			RiskTolerance:  "medium",
			CloudProviders: []string{"AWS"},
		}
		
		requirements := []string{"migration", "security", "compliance"}
		phases, err := roadmapGen.GeneratePhases(ctx, requirements, constraints)
		if err != nil {
			log.Printf("Error generating phases: %v", err)
		} else {
			fmt.Printf("    ✓ Phase generation: %d phases created\n", len(phases))
		}
		
		// Test resource estimation
		resourceEst, err := roadmapGen.EstimateResources(ctx, roadmap.Phases)
		if err != nil {
			log.Printf("Error estimating resources: %v", err)
		} else {
			fmt.Printf("    ✓ Resource estimation: %d total hours, %s total cost\n", 
				resourceEst.TotalHours, resourceEst.TotalCost)
		}
		
		// Test dependency calculation
		deps, err := roadmapGen.CalculateDependencies(ctx, roadmap.Phases)
		if err != nil {
			log.Printf("Error calculating dependencies: %v", err)
		} else {
			fmt.Printf("    ✓ Dependency calculation: %d dependencies identified\n", len(deps))
		}
		
		// Test milestone generation
		milestones, err := roadmapGen.GenerateMilestones(ctx, roadmap.Phases)
		if err != nil {
			log.Printf("Error generating milestones: %v", err)
		} else {
			fmt.Printf("    ✓ Milestone generation: %d milestones created\n", len(milestones))
		}
		
		// Test validation
		validation, err := roadmapGen.ValidateRoadmap(ctx, roadmap)
		if err != nil {
			log.Printf("Error validating roadmap: %v", err)
		} else {
			fmt.Printf("    ✓ Roadmap validation: Valid=%t, Quality=%.2f, Completeness=%.2f\n", 
				validation.IsValid, validation.QualityScore, validation.CompletenessScore)
			
			if len(validation.Errors) > 0 {
				fmt.Printf("      Errors: %v\n", validation.Errors)
			}
			if len(validation.Warnings) > 0 {
				fmt.Printf("      Warnings: %v\n", validation.Warnings)
			}
			if len(validation.Suggestions) > 0 {
				fmt.Printf("      Suggestions: %v\n", validation.Suggestions)
			}
		}
		
		// Test JSON serialization
		jsonData, err := json.MarshalIndent(roadmap, "", "  ")
		if err != nil {
			log.Printf("Error serializing roadmap to JSON: %v", err)
		} else {
			fmt.Printf("    ✓ JSON serialization: %d bytes\n", len(jsonData))
			
			// Save to file for inspection
			filename := fmt.Sprintf("test_roadmap_%s.json", scenario.inquiry.ID)
			if err := saveToFile(filename, jsonData); err != nil {
				log.Printf("Error saving roadmap to file: %v", err)
			} else {
				fmt.Printf("    ✓ Saved roadmap to %s\n", filename)
			}
		}
	}
	
	// Test error scenarios
	fmt.Printf("\n=== Testing Error Scenarios ===\n")
	
	// Test with nil inquiry
	_, err := roadmapGen.GenerateImplementationRoadmap(ctx, nil, nil)
	if err != nil {
		fmt.Printf("✓ Nil inquiry handling: %v\n", err)
	}
	
	// Test with empty inquiry
	emptyInquiry := &domain.Inquiry{
		ID: "empty-inquiry",
	}
	_, err = roadmapGen.GenerateImplementationRoadmap(ctx, emptyInquiry, nil)
	if err != nil {
		fmt.Printf("✓ Empty inquiry handling: %v\n", err)
	} else {
		fmt.Printf("✓ Empty inquiry handled gracefully\n")
	}
	
	fmt.Printf("\n=== Testing Template System ===\n")
	
	// Test different project types
	projectTypes := []string{"migration", "optimization", "assessment", "architecture", "unknown"}
	
	for _, projectType := range projectTypes {
		testInquiry := &domain.Inquiry{
			ID:       fmt.Sprintf("test-%s", projectType),
			Company:  "TestCorp",
			Services: []string{projectType},
			Message:  fmt.Sprintf("Test %s project", projectType),
		}
		
		roadmap, err := roadmapGen.GenerateImplementationRoadmap(ctx, testInquiry, nil)
		if err != nil {
			fmt.Printf("  %s: Error - %v\n", projectType, err)
		} else {
			fmt.Printf("  %s: ✓ Generated %d phases\n", projectType, len(roadmap.Phases))
		}
	}
	
	fmt.Printf("\n=== All Tests Completed ===\n")
}

func saveToFile(filename string, data []byte) error {
	// In a real implementation, you would save to a file
	// For this test, we'll just simulate success
	fmt.Printf("    [Simulated] Would save %d bytes to %s\n", len(data), filename)
	return nil
}