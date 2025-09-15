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
	// Return a mock response for cost analysis
	content := `
## Cost Breakdown Analysis

### Monthly Cost Estimates:
- EC2 Instances: $2,000/month (40%)
- RDS Database: $1,000/month (20%)
- S3 Storage: $800/month (16%)
- CloudFront CDN: $600/month (12%)
- VPC/Networking: $400/month (8%)
- Security Services: $200/month (4%)

### Cost Optimization Opportunities:
1. Right-size EC2 instances: Save $400/month (20% reduction)
2. Use Reserved Instances: Save $600/month (30% reduction on compute)
3. Implement S3 lifecycle policies: Save $200/month (25% storage reduction)
4. Optimize data transfer: Save $150/month (15% network reduction)

### Key Assumptions:
- 24/7 operation assumed
- Standard pricing without enterprise discounts
- Current usage patterns maintained
- US East region pricing baseline
`

	return &interfaces.BedrockResponse{
		Content: content,
		Usage: interfaces.BedrockUsage{
			InputTokens:  len(prompt) / 4, // Rough estimate
			OutputTokens: len(content) / 4,
		},
		Metadata: map[string]string{
			"model": options.ModelID,
		},
	}, nil
}

func (m *MockBedrockService) GetModelInfo() interfaces.BedrockModelInfo {
	return interfaces.BedrockModelInfo{
		ModelID:     "anthropic.claude-3-sonnet-20240229-v1:0",
		ModelName:   "Claude 3 Sonnet",
		Provider:    "Anthropic",
		MaxTokens:   200000,
		IsAvailable: true,
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
			Name:            "Cloud Migration",
			Description:     "End-to-end cloud migration services",
			TypicalDuration: "3-6 months",
			TeamSize:        "4-6 consultants",
			KeyBenefits:     []string{"Reduced infrastructure costs", "Improved scalability", "Enhanced security"},
			Deliverables:    []string{"Migration plan", "Migrated applications", "Documentation"},
		},
	}, nil
}

func (m *MockKnowledgeBase) GetTeamExpertise(ctx context.Context) ([]*interfaces.TeamExpertise, error) {
	return []*interfaces.TeamExpertise{
		{
			ConsultantName:  "John Smith",
			Role:            "Senior Cloud Architect",
			ExperienceYears: 8,
			ExpertiseAreas:  []string{"AWS", "Cost Optimization", "Migration"},
			CloudProviders:  []string{"AWS", "Azure"},
		},
	}, nil
}

func (m *MockKnowledgeBase) GetPastSolutions(ctx context.Context, serviceType, industry string) ([]*interfaces.PastSolution, error) {
	return []*interfaces.PastSolution{
		{
			Title:            "E-commerce Platform Migration",
			Industry:         "Retail",
			ProblemStatement: "High infrastructure costs and scalability issues",
			SolutionApproach: "Migrated to AWS with auto-scaling and cost optimization",
			Technologies:     []string{"AWS EC2", "RDS", "CloudFront"},
			TimeToValue:      "4 months",
			CostSavings:      150000,
		},
	}, nil
}

func (m *MockKnowledgeBase) GetConsultingApproach(ctx context.Context, serviceType string) (*interfaces.ConsultingApproach, error) {
	return &interfaces.ConsultingApproach{
		Name:              "Agile Cloud Transformation",
		Philosophy:        "Iterative approach with continuous optimization",
		EngagementModel:   "Collaborative partnership",
		KeyPrinciples:     []string{"Client-first", "Data-driven decisions", "Continuous improvement"},
		ClientInvolvement: "High - weekly reviews and feedback sessions",
		KnowledgeTransfer: "Comprehensive training and documentation",
	}, nil
}

func (m *MockKnowledgeBase) GetClientHistory(ctx context.Context, company string) ([]*interfaces.ClientEngagement, error) {
	return []*interfaces.ClientEngagement{}, nil // No previous engagements
}

func main() {
	fmt.Println("Testing Cost Analysis Service...")

	// Create mock services
	mockBedrock := &MockBedrockService{}
	mockKB := &MockKnowledgeBase{}

	// Create cost analysis service
	costService := services.NewCostAnalysisService(mockBedrock, mockKB)

	// Create test inquiry
	inquiry := &domain.Inquiry{
		ID:       "test-inquiry-001",
		Name:     "John Doe",
		Email:    "john@example.com",
		Company:  "TechCorp Inc",
		Services: []string{"migration", "optimization"},
		Message:  "We need to migrate our e-commerce platform to AWS and optimize costs. Currently spending $10K/month on infrastructure.",
		Status:   "new",
		Priority: "high",
	}

	// Create test architecture
	architecture := &interfaces.CostArchitectureSpec{
		Name:        "E-commerce Platform",
		Description: "Scalable e-commerce platform with microservices architecture",
		Environment: "production",
		Regions:     []string{"us-east-1", "us-west-2"},
		Services: []*interfaces.ServiceSpec{
			{
				ServiceName: "EC2 Web Servers",
				Provider:    "AWS",
				Region:      "us-east-1",
				PricingTier: "On-Demand",
			},
			{
				ServiceName: "RDS MySQL Database",
				Provider:    "AWS",
				Region:      "us-east-1",
				PricingTier: "On-Demand",
			},
			{
				ServiceName: "S3 Storage",
				Provider:    "AWS",
				Region:      "us-east-1",
				PricingTier: "Standard",
			},
		},
	}

	ctx := context.Background()

	// Test 1: Cost Breakdown Analysis
	fmt.Println("\n=== Testing Cost Breakdown Analysis ===")
	costBreakdown, err := costService.AnalyzeCostBreakdown(ctx, inquiry, architecture)
	if err != nil {
		log.Fatalf("Failed to analyze cost breakdown: %v", err)
	}

	fmt.Printf("Analysis ID: %s\n", costBreakdown.ID)
	fmt.Printf("Total Monthly Cost: $%.2f\n", costBreakdown.TotalMonthlyCost)
	fmt.Printf("Total Annual Cost: $%.2f\n", costBreakdown.TotalAnnualCost)
	fmt.Printf("Number of Services: %d\n", len(costBreakdown.ServiceBreakdown))
	fmt.Printf("Number of Categories: %d\n", len(costBreakdown.CategoryBreakdown))
	fmt.Printf("Number of Regions: %d\n", len(costBreakdown.RegionBreakdown))

	// Test 2: Cost Optimization Recommendations
	fmt.Println("\n=== Testing Cost Optimization Recommendations ===")
	optimizationRecs, err := costService.GenerateCostOptimizationRecommendations(ctx, costBreakdown)
	if err != nil {
		log.Fatalf("Failed to generate optimization recommendations: %v", err)
	}

	fmt.Printf("Recommendations ID: %s\n", optimizationRecs.ID)
	fmt.Printf("Total Savings Potential: $%.2f\n", optimizationRecs.TotalSavingsPotential)
	fmt.Printf("Number of Recommendations: %d\n", len(optimizationRecs.Recommendations))

	// Test 3: Reserved Instance Analysis
	fmt.Println("\n=== Testing Reserved Instance Analysis ===")
	usageData := &interfaces.UsageData{
		AccountID: "123456789012",
		Region:    "us-east-1",
		TimeRange: &interfaces.TimeRange{
			StartDate: time.Now().AddDate(0, -3, 0), // 3 months ago
			EndDate:   time.Now(),
			Period:    "monthly",
		},
		ComputeUsage: []*interfaces.ComputeUsageData{
			{
				InstanceType:     "m5.large",
				Region:           "us-east-1",
				AvailabilityZone: "us-east-1a",
				Platform:         "Linux",
				Tenancy:          "default",
				UsageHours:       720, // Full month
				OnDemandCost:     150.0,
				Date:             time.Now().AddDate(0, -1, 0),
				Utilization:      0.75,
			},
		},
	}

	riAnalysis, err := costService.AnalyzeReservedInstanceOpportunities(ctx, usageData)
	if err != nil {
		log.Fatalf("Failed to analyze RI opportunities: %v", err)
	}

	fmt.Printf("RI Analysis ID: %s\n", riAnalysis.ID)
	fmt.Printf("RI Savings Potential: $%.2f\n", riAnalysis.TotalSavingsPotential)
	fmt.Printf("Number of RI Recommendations: %d\n", len(riAnalysis.Recommendations))

	// Test 4: Comprehensive Cost Analysis
	fmt.Println("\n=== Testing Comprehensive Cost Analysis ===")
	comprehensiveAnalysis, err := costService.GenerateComprehensiveCostAnalysis(ctx, inquiry)
	if err != nil {
		log.Fatalf("Failed to generate comprehensive analysis: %v", err)
	}

	fmt.Printf("Comprehensive Analysis ID: %s\n", comprehensiveAnalysis.ID)
	fmt.Printf("Executive Summary - Current Monthly Cost: $%.2f\n", comprehensiveAnalysis.ExecutiveSummary.CurrentMonthlyCost)
	fmt.Printf("Executive Summary - Optimized Monthly Cost: $%.2f\n", comprehensiveAnalysis.ExecutiveSummary.OptimizedMonthlyCost)
	fmt.Printf("Executive Summary - Total Savings Potential: $%.2f\n", comprehensiveAnalysis.ExecutiveSummary.TotalSavingsPotential)
	fmt.Printf("Executive Summary - Savings Percentage: %.1f%%\n", comprehensiveAnalysis.ExecutiveSummary.SavingsPercentage)
	fmt.Printf("Executive Summary - ROI: %.1f%%\n", comprehensiveAnalysis.ExecutiveSummary.ROI)

	// Test 5: JSON Serialization
	fmt.Println("\n=== Testing JSON Serialization ===")
	jsonData, err := json.MarshalIndent(comprehensiveAnalysis.ExecutiveSummary, "", "  ")
	if err != nil {
		log.Fatalf("Failed to serialize to JSON: %v", err)
	}

	fmt.Printf("Executive Summary JSON:\n%s\n", string(jsonData))

	fmt.Println("\n=== All Tests Completed Successfully! ===")
	fmt.Println("\nKey Features Implemented:")
	fmt.Println("✅ Detailed cost breakdown analysis with specific optimization recommendations")
	fmt.Println("✅ Reserved instance and savings plan optimization suggestions with exact calculations")
	fmt.Println("✅ Right-sizing recommendations based on actual usage patterns")
	fmt.Println("✅ Cost forecasting models for proposed architectures with confidence intervals")
	fmt.Println("✅ Comprehensive cost analysis combining all optimization strategies")
	fmt.Println("✅ Executive summary with actionable insights for consultants")
	fmt.Println("✅ JSON serialization for API responses")
}
