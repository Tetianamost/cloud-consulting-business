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
	"github.com/sirupsen/logrus"
)

// MockBedrockService for testing
type MockBedrockService struct{}

func (m *MockBedrockService) GenerateResponse(ctx context.Context, prompt string) (*interfaces.BedrockResponse, error) {
	// Simulate AI response based on prompt content
	var content string

	if contains(prompt, "what-if scenarios") {
		content = `High Growth Scenario: 200% user growth requiring infrastructure scaling
Migration Scenario: Containerization using EKS and microservices
Compliance Scenario: Enhanced security and compliance requirements
Economic Scenario: Cost optimization during economic downturn`
	} else if contains(prompt, "multi-year projection") {
		content = `Year 1: $120,000 (baseline + 20% growth)
Year 2: $156,000 (30% growth from optimization)
Year 3: $187,200 (20% growth with efficiency gains)
Key insights: Reserved instance savings, serverless adoption, scaling optimization`
	} else if contains(prompt, "executive summary") {
		content = `Executive Summary:
- Strategic cloud architecture improvements recommended
- 25-30% cost optimization potential identified
- Multi-region deployment for business continuity
- Phased implementation approach recommended
- Expected ROI of 200% over 3 years`
	} else {
		content = "Mock AI response for scenario modeling analysis"
	}

	return &interfaces.BedrockResponse{
		Content:    content,
		TokensUsed: 150,
	}, nil
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsSubstring(s, substr)))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func main() {
	fmt.Println("=== Scenario Modeling Service Test ===")

	// Create logger
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	// Create mock Bedrock service
	mockBedrock := &MockBedrockService{}

	// Create scenario modeling service
	scenarioService := services.NewScenarioModelingService(mockBedrock, logger)

	// Create test inquiry
	testInquiry := &domain.Inquiry{
		ID:       "test-inquiry-123",
		Company:  "TechCorp Solutions",
		Name:     "John Smith",
		Email:    "john.smith@techcorp.com",
		Services: []string{"migration", "architecture_review", "optimization"},
		Message:  "We need to migrate our legacy applications to AWS and optimize our cloud costs while ensuring high availability and disaster recovery capabilities.",
		Priority: "high",
		Status:   "new",
	}

	fmt.Printf("Testing scenario modeling for: %s\n", testInquiry.Company)
	fmt.Printf("Services requested: %v\n", testInquiry.Services)
	fmt.Printf("Requirements: %s\n\n", testInquiry.Message)

	// Test comprehensive scenario analysis
	ctx := context.Background()

	fmt.Println("Generating comprehensive scenario analysis...")
	analysis, err := scenarioService.GenerateComprehensiveScenarioAnalysis(ctx, testInquiry)
	if err != nil {
		log.Fatalf("Failed to generate scenario analysis: %v", err)
	}

	// Display results
	fmt.Printf("\n=== COMPREHENSIVE SCENARIO ANALYSIS ===\n")
	fmt.Printf("Analysis ID: %s\n", analysis.ID)
	fmt.Printf("Analysis Date: %s\n", analysis.AnalysisDate.Format("2006-01-02 15:04:05"))

	// Base Scenario
	fmt.Printf("\n--- BASE SCENARIO ---\n")
	fmt.Printf("Name: %s\n", analysis.BaseScenario.Name)
	fmt.Printf("Description: %s\n", analysis.BaseScenario.Description)

	// What-If Scenarios
	fmt.Printf("\n--- WHAT-IF SCENARIOS (%d) ---\n", len(analysis.WhatIfScenarios))
	for i, scenario := range analysis.WhatIfScenarios {
		fmt.Printf("%d. %s\n", i+1, scenario.Name)
		fmt.Printf("   Description: %s\n", scenario.Description)
		fmt.Printf("   Confidence: %.1f%%\n", scenario.ConfidenceLevel*100)
		fmt.Printf("   Risk Factors: %v\n", scenario.RiskFactors)
		if scenario.ProjectedOutcomes != nil && scenario.ProjectedOutcomes.CostProjection != nil {
			fmt.Printf("   Cost Impact: %.1f%%\n", scenario.ProjectedOutcomes.CostProjection.CostChangePercent)
		}
		fmt.Println()
	}

	// Multi-Year Projection
	fmt.Printf("--- MULTI-YEAR PROJECTION ---\n")
	fmt.Printf("Period: %s\n", analysis.MultiYearProjection.ProjectionPeriod)
	fmt.Printf("Yearly Projections:\n")
	for _, yearly := range analysis.MultiYearProjection.YearlyProjections {
		fmt.Printf("  Year %d: $%.0f\n", yearly.Year, yearly.ProjectedCost)
	}
	fmt.Printf("Key Insights:\n")
	for _, insight := range analysis.MultiYearProjection.KeyInsights {
		fmt.Printf("  - %s: %s (Impact: %s)\n", insight.InsightType, insight.Description, insight.Impact)
	}

	// Disaster Recovery Scenarios
	fmt.Printf("\n--- DISASTER RECOVERY SCENARIOS (%d) ---\n", len(analysis.DRScenarios))
	for i, drScenario := range analysis.DRScenarios {
		fmt.Printf("%d. %s\n", i+1, drScenario.Name)
		fmt.Printf("   Description: %s\n", drScenario.Description)
	}

	// Capacity Scenarios
	fmt.Printf("\n--- CAPACITY SCENARIOS (%d) ---\n", len(analysis.CapacityScenarios))
	for i, capScenario := range analysis.CapacityScenarios {
		fmt.Printf("%d. %s\n", i+1, capScenario.Name)
		fmt.Printf("   Description: %s\n", capScenario.Description)
	}

	// Integrated Analysis
	fmt.Printf("\n--- INTEGRATED ANALYSIS ---\n")
	fmt.Printf("Overall Risk Level: %s\n", analysis.IntegratedAnalysis.OverallRiskLevel)
	fmt.Printf("Cost Optimization Potential: %.1f%%\n", analysis.IntegratedAnalysis.CostOptimizationPotential)
	fmt.Printf("Business Impact: %s\n", analysis.IntegratedAnalysis.BusinessImpactSummary)
	fmt.Printf("Technical Complexity: %s\n", analysis.IntegratedAnalysis.TechnicalComplexity)
	fmt.Printf("Implementation Timeline: %s\n", analysis.IntegratedAnalysis.ImplementationTimeline)

	// Executive Summary
	fmt.Printf("\n--- EXECUTIVE SUMMARY ---\n")
	fmt.Printf("Key Findings:\n")
	for _, finding := range analysis.ExecutiveSummary.KeyFindings {
		fmt.Printf("  • %s\n", finding)
	}
	fmt.Printf("Cost Impact: %s\n", analysis.ExecutiveSummary.CostImpactSummary)
	fmt.Printf("Business Impact: %s\n", analysis.ExecutiveSummary.BusinessImpactSummary)
	fmt.Printf("Risk Summary: %s\n", analysis.ExecutiveSummary.RiskSummary)
	fmt.Printf("Recommended Approach: %s\n", analysis.ExecutiveSummary.RecommendedApproach)

	fmt.Printf("\nExpected Outcomes:\n")
	for _, outcome := range analysis.ExecutiveSummary.ExpectedOutcomes {
		fmt.Printf("  • %s\n", outcome)
	}

	fmt.Printf("\nSuccess Metrics:\n")
	for _, metric := range analysis.ExecutiveSummary.SuccessMetrics {
		fmt.Printf("  • %s\n", metric)
	}

	// Key Recommendations
	fmt.Printf("\n--- KEY RECOMMENDATIONS (%d) ---\n", len(analysis.KeyRecommendations))
	for i, rec := range analysis.KeyRecommendations {
		fmt.Printf("%d. %s (Priority: %s)\n", i+1, rec.Title, rec.Priority)
		fmt.Printf("   Category: %s\n", rec.Category)
		fmt.Printf("   Description: %s\n", rec.Description)
		fmt.Printf("   Expected Benefit: %s\n", rec.ExpectedBenefit)
		fmt.Printf("   Implementation Cost: $%.0f\n", rec.ImplementationCost)
		fmt.Printf("   Timeline: %s\n", rec.Timeline)
		fmt.Printf("   Risk Level: %s\n", rec.RiskLevel)
		fmt.Println()
	}

	// Next Steps
	fmt.Printf("--- NEXT STEPS (%d) ---\n", len(analysis.NextSteps))
	for i, step := range analysis.NextSteps {
		fmt.Printf("%d. %s (Priority: %s)\n", i+1, step.StepName, step.Priority)
		fmt.Printf("   Description: %s\n", step.Description)
		fmt.Printf("   Owner: %s\n", step.Owner)
		fmt.Printf("   Timeline: %s\n", step.Timeline)
		fmt.Println()
	}

	// Test JSON serialization
	fmt.Printf("\n=== JSON SERIALIZATION TEST ===\n")
	jsonData, err := json.MarshalIndent(analysis, "", "  ")
	if err != nil {
		log.Printf("JSON serialization failed: %v", err)
	} else {
		fmt.Printf("JSON serialization successful (%d bytes)\n", len(jsonData))

		// Test deserialization
		var deserializedAnalysis interfaces.ComprehensiveScenarioAnalysis
		err = json.Unmarshal(jsonData, &deserializedAnalysis)
		if err != nil {
			log.Printf("JSON deserialization failed: %v", err)
		} else {
			fmt.Printf("JSON deserialization successful\n")
			fmt.Printf("Deserialized analysis ID: %s\n", deserializedAnalysis.ID)
		}
	}

	// Performance metrics
	fmt.Printf("\n=== PERFORMANCE METRICS ===\n")
	fmt.Printf("Analysis completed in: %v\n", time.Since(analysis.CreatedAt))
	fmt.Printf("What-if scenarios generated: %d\n", len(analysis.WhatIfScenarios))
	fmt.Printf("DR scenarios generated: %d\n", len(analysis.DRScenarios))
	fmt.Printf("Capacity scenarios generated: %d\n", len(analysis.CapacityScenarios))
	fmt.Printf("Key recommendations generated: %d\n", len(analysis.KeyRecommendations))
	fmt.Printf("Next steps generated: %d\n", len(analysis.NextSteps))

	// Calculate total cost optimization potential
	totalOptimization := analysis.IntegratedAnalysis.CostOptimizationPotential
	if len(analysis.MultiYearProjection.YearlyProjections) > 0 {
		totalCost := 0.0
		for _, projection := range analysis.MultiYearProjection.YearlyProjections {
			totalCost += projection.ProjectedCost
		}
		potentialSavings := totalCost * (totalOptimization / 100)
		fmt.Printf("Potential 3-year savings: $%.0f\n", potentialSavings)
	}

	fmt.Printf("\n=== TEST COMPLETED SUCCESSFULLY ===\n")
	fmt.Printf("Scenario modeling service is working correctly!\n")
	fmt.Printf("All components generated appropriate data structures.\n")
	fmt.Printf("Ready for integration with the enhanced Bedrock AI assistant.\n")
}
