package main

import (
	"context"
	"fmt"
	"log"

	"github.com/cloud-consulting/backend/internal/domain"
	"github.com/cloud-consulting/backend/internal/interfaces"
	"github.com/cloud-consulting/backend/internal/services"
)

// MockBedrockService for testing
type MockBedrockService struct{}

func (m *MockBedrockService) GenerateText(ctx context.Context, prompt string, options *interfaces.BedrockOptions) (*interfaces.BedrockResponse, error) {
	// Return mock JSON responses based on prompt content
	if contains(prompt, "executive summary") {
		return &interfaces.BedrockResponse{
			Content: `{
				"key_findings": ["Current infrastructure has scalability limitations", "Security posture needs enhancement", "Cost optimization opportunities identified"],
				"recommendations": ["Migrate to cloud-native architecture", "Implement zero-trust security model", "Establish automated cost monitoring"],
				"business_impact": "Expected 30% reduction in operational costs and 50% improvement in system scalability",
				"investment_summary": "Initial investment of $200K with projected annual savings of $150K",
				"timeline": "6-month implementation with value realization starting in month 3"
			}`,
		}, nil
	}

	if contains(prompt, "business case") {
		return &interfaces.BedrockResponse{
			Content: `{
				"problem_statement": "Legacy infrastructure limiting business growth and increasing operational costs",
				"proposed_solution": "Cloud transformation with modern architecture and automated operations",
				"expected_benefits": ["Reduced operational costs", "Improved scalability", "Enhanced security", "Faster time to market"],
				"cost_benefit_analysis": "Initial investment of $200K offset by $150K annual savings, achieving positive ROI within 18 months",
				"roi": 35.5,
				"payback_period": "16 months"
			}`,
		}, nil
	}

	if contains(prompt, "ROI analysis") {
		return &interfaces.BedrockResponse{
			Content: `{
				"initial_investment": 200000.0,
				"annual_savings": 150000.0,
				"payback_period": "16 months",
				"npv": 275000.0,
				"irr": 22.5,
				"roi_percentage": 35.5
			}`,
		}, nil
	}

	if contains(prompt, "technical specifications") {
		return &interfaces.BedrockResponse{
			Content: `{
				"architecture_overview": "Cloud-native microservices architecture with containerized deployment on Kubernetes",
				"technical_requirements": ["High availability (99.9% uptime)", "Auto-scaling capabilities", "Secure data encryption", "API-first design"],
				"component_details": {
					"compute": "Kubernetes cluster with auto-scaling node groups",
					"storage": "Distributed object storage with automated backup",
					"networking": "Virtual private cloud with security groups and load balancers",
					"security": "Identity and access management with multi-factor authentication"
				},
				"integration_points": ["External payment gateway", "CRM system integration", "Analytics platform"],
				"security_considerations": ["Data encryption at rest and in transit", "Network segmentation", "Regular security audits", "Compliance monitoring"]
			}`,
		}, nil
	}

	if contains(prompt, "implementation plan") {
		return &interfaces.BedrockResponse{
			Content: `{
				"phases": [
					{
						"name": "Assessment and Planning",
						"description": "Comprehensive assessment of current state and detailed planning",
						"duration": "4 weeks",
						"deliverables": ["Current state assessment", "Architecture design", "Migration plan"],
						"dependencies": ["Stakeholder alignment", "Resource allocation"]
					},
					{
						"name": "Infrastructure Setup",
						"description": "Setup cloud infrastructure and core services",
						"duration": "6 weeks",
						"deliverables": ["Cloud infrastructure", "Security framework", "Monitoring setup"],
						"dependencies": ["Architecture approval", "Security clearance"]
					}
				],
				"timeline": "16 weeks total implementation",
				"resource_requirements": ["Cloud architect", "DevOps engineer", "Security specialist", "Project manager"],
				"risk_mitigation": ["Comprehensive testing", "Phased rollout", "Rollback procedures", "24/7 monitoring"],
				"success_metrics": ["System uptime > 99.9%", "Response time < 200ms", "Cost reduction > 25%"]
			}`,
		}, nil
	}

	// Default response
	return &interfaces.BedrockResponse{
		Content: "Mock response for testing",
	}, nil
}

func (m *MockBedrockService) GetModelInfo() interfaces.BedrockModelInfo {
	return interfaces.BedrockModelInfo{
		ModelID:     "amazon.nova-lite-v1:0",
		ModelName:   "Nova Lite",
		Provider:    "Amazon",
		MaxTokens:   4000,
		IsAvailable: true,
	}
}

func (m *MockBedrockService) IsHealthy() bool {
	return true
}

func contains(text, substr string) bool {
	return len(text) >= len(substr) && text[:len(substr)] == substr ||
		len(text) > len(substr) && text[len(text)-len(substr):] == substr ||
		(len(text) > len(substr) && findInString(text, substr))
}

func findInString(text, substr string) bool {
	for i := 0; i <= len(text)-len(substr); i++ {
		if text[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func main() {
	fmt.Println("Testing Enhanced Report Generator...")

	// Create mock services
	mockBedrock := &MockBedrockService{}

	// Create report generator
	reportGen := services.NewBasicReportGenerator(mockBedrock, nil, nil)

	// Create test inquiry
	inquiry := &domain.Inquiry{
		ID:       "test-123",
		Name:     "John Smith",
		Email:    "john@example.com",
		Company:  "TechCorp Inc",
		Services: []string{"migration", "optimization"},
		Message:  "We need to migrate our legacy infrastructure to the cloud and optimize costs",
		Priority: "high",
	}

	// Generate enhanced report
	ctx := context.Background()
	report, err := reportGen.GenerateReport(ctx, inquiry)
	if err != nil {
		log.Fatalf("Failed to generate report: %v", err)
	}

	// Display results
	fmt.Printf("\n=== ENHANCED CONSULTANT REPORT ===\n")
	fmt.Printf("Title: %s\n", report.Title)
	fmt.Printf("Type: %s\n", report.Type)
	fmt.Printf("Generated By: %s\n", report.GeneratedBy)
	fmt.Printf("\n=== EXECUTIVE SUMMARY ===\n")
	if report.ExecutiveSummary != nil {
		fmt.Printf("Key Findings: %v\n", report.ExecutiveSummary.KeyFindings)
		fmt.Printf("Business Impact: %s\n", report.ExecutiveSummary.BusinessImpact)
		fmt.Printf("Timeline: %s\n", report.ExecutiveSummary.Timeline)
	}

	fmt.Printf("\n=== BUSINESS CASE ===\n")
	if report.BusinessCase != nil {
		fmt.Printf("ROI: %.1f%%\n", report.BusinessCase.ROI)
		fmt.Printf("Payback Period: %s\n", report.BusinessCase.PaybackPeriod)
		fmt.Printf("Expected Benefits: %v\n", report.BusinessCase.ExpectedBenefits)
	}

	fmt.Printf("\n=== ROI ANALYSIS ===\n")
	if report.ROIAnalysis != nil {
		fmt.Printf("Initial Investment: $%.0f\n", report.ROIAnalysis.InitialInvestment)
		fmt.Printf("Annual Savings: $%.0f\n", report.ROIAnalysis.AnnualSavings)
		fmt.Printf("ROI Percentage: %.1f%%\n", report.ROIAnalysis.ROIPercentage)
	}

	fmt.Printf("\n=== QUALITY METRICS ===\n")
	if report.QualityMetrics != nil {
		fmt.Printf("Completeness Score: %.1f%%\n", report.QualityMetrics.CompletenessScore)
		fmt.Printf("Overall Quality Score: %.1f%%\n", report.QualityMetrics.OverallQualityScore)
	}

	fmt.Printf("\n=== FULL REPORT CONTENT ===\n")
	fmt.Printf("%s\n", report.Content)

	fmt.Println("\n✅ Enhanced report generation test completed successfully!")
}
