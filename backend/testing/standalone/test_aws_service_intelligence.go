package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/cloud-consulting/backend/internal/config"
	"github.com/cloud-consulting/backend/internal/interfaces"
	"github.com/cloud-consulting/backend/internal/services"
)

func main() {
	fmt.Println("=== AWS Service Intelligence Service Test ===")

	// Initialize Bedrock service
	bedrockConfig := &config.BedrockConfig{
		APIKey:  "test-api-key",
		BaseURL: "https://bedrock-runtime.us-east-1.amazonaws.com",
		ModelID: "amazon.nova-lite-v1:0",
		Timeout: 30,
	}

	bedrockService := services.NewBedrockService(bedrockConfig)

	// Initialize AWS Service Intelligence service
	intelligenceService := services.NewAWSServiceIntelligenceService(bedrockService)

	ctx := context.Background()

	// Test 1: Get Service Status
	fmt.Println("\n--- Test 1: Get Service Status ---")
	status, err := intelligenceService.GetServiceStatus(ctx, "us-east-1")
	if err != nil {
		log.Printf("Error getting service status: %v", err)
	} else {
		fmt.Printf("Region: %s\n", status.Region)
		fmt.Printf("Overall Health: %s\n", status.OverallHealth)
		fmt.Printf("Services Count: %d\n", len(status.Services))
		for serviceName, health := range status.Services {
			fmt.Printf("  %s: %s\n", serviceName, health.Status)
		}
	}

	// Test 2: Get Service Status History
	fmt.Println("\n--- Test 2: Get Service Status History ---")
	history, err := intelligenceService.GetServiceStatusHistory(ctx, "EC2", "us-east-1", 7*24*time.Hour)
	if err != nil {
		log.Printf("Error getting service history: %v", err)
	} else {
		fmt.Printf("History Events Count: %d\n", len(history))
		for _, event := range history {
			fmt.Printf("  %s - %s: %s\n", event.Timestamp.Format("2006-01-02 15:04"), event.Status, event.Message)
		}
	}

	// Test 3: Analyze Service Impact
	fmt.Println("\n--- Test 3: Analyze Service Impact ---")
	impact, err := intelligenceService.AnalyzeServiceImpact(ctx, "RDS", "us-west-2")
	if err != nil {
		log.Printf("Error analyzing service impact: %v", err)
	} else {
		fmt.Printf("Service: %s\n", impact.ServiceName)
		fmt.Printf("Impact Level: %s\n", impact.ImpactLevel)
		fmt.Printf("Action Required: %t\n", impact.ActionRequired)
		fmt.Printf("Recommended Actions: %d\n", len(impact.RecommendedActions))
	}

	// Test 4: Get New Services
	fmt.Println("\n--- Test 4: Get New Services ---")
	since := time.Now().Add(-90 * 24 * time.Hour) // Last 90 days
	newServices, err := intelligenceService.GetNewServices(ctx, since)
	if err != nil {
		log.Printf("Error getting new services: %v", err)
	} else {
		fmt.Printf("New Services Count: %d\n", len(newServices))
		for _, service := range newServices {
			fmt.Printf("  %s (%s): %s\n", service.ServiceName, service.Category, service.Description)
		}
	}

	// Test 5: Evaluate Service for Client
	fmt.Println("\n--- Test 5: Evaluate Service for Client ---")
	if len(newServices) > 0 {
		clientContext := &interfaces.ClientContext{
			IndustryVertical:       "Healthcare",
			CompanySize:            "Medium",
			CurrentServices:        []string{"EC2", "RDS", "S3"},
			PreferredRegions:       []string{"us-east-1", "us-west-2"},
			ComplianceRequirements: []string{"HIPAA", "SOC2"},
			TechnicalMaturity:      "intermediate",
			Workloads: []interfaces.WorkloadProfile{
				{
					Name:            "Web Application",
					Type:            "web",
					Scale:           "medium",
					PerformanceReqs: []string{"High availability", "Low latency"},
					SecurityReqs:    []string{"Encryption", "Access control"},
					ComplianceReqs:  []string{"HIPAA"},
				},
			},
		}

		evaluation, err := intelligenceService.EvaluateServiceForClient(ctx, newServices[0], clientContext)
		if err != nil {
			log.Printf("Error evaluating service: %v", err)
		} else {
			fmt.Printf("Service: %s\n", evaluation.ServiceName)
			fmt.Printf("Relevance Score: %.1f\n", evaluation.RelevanceScore)
			fmt.Printf("Adoption Priority: %s\n", evaluation.AdoptionPriority)
			fmt.Printf("Implementation Complexity: %s\n", evaluation.ImplementationComplexity)
			fmt.Printf("Benefits: %d\n", len(evaluation.Benefits))
			fmt.Printf("Challenges: %d\n", len(evaluation.Challenges))
		}
	}

	// Test 6: Get Service Recommendations
	fmt.Println("\n--- Test 6: Get Service Recommendations ---")
	clientContext := &interfaces.ClientContext{
		IndustryVertical:       "Financial Services",
		CompanySize:            "Large",
		CurrentServices:        []string{"EC2", "RDS", "S3", "Lambda"},
		PreferredRegions:       []string{"us-east-1", "eu-west-1"},
		ComplianceRequirements: []string{"PCI-DSS", "SOX"},
		TechnicalMaturity:      "advanced",
		Workloads: []interfaces.WorkloadProfile{
			{
				Name:            "Trading Platform",
				Type:            "web",
				Scale:           "large",
				PerformanceReqs: []string{"Ultra-low latency", "High throughput"},
				SecurityReqs:    []string{"End-to-end encryption", "Multi-factor auth"},
				ComplianceReqs:  []string{"PCI-DSS"},
			},
		},
	}

	recommendations, err := intelligenceService.GetServiceRecommendations(ctx, clientContext)
	if err != nil {
		log.Printf("Error getting recommendations: %v", err)
	} else {
		fmt.Printf("Recommendations Count: %d\n", len(recommendations))
		for _, rec := range recommendations {
			fmt.Printf("  %s (%s): %s - Priority: %s\n",
				rec.ServiceName, rec.RecommendationType, rec.Rationale, rec.Priority)
		}
	}

	// Test 7: Get Deprecation Alerts
	fmt.Println("\n--- Test 7: Get Deprecation Alerts ---")
	alerts, err := intelligenceService.GetDeprecationAlerts(ctx)
	if err != nil {
		log.Printf("Error getting deprecation alerts: %v", err)
	} else {
		fmt.Printf("Deprecation Alerts Count: %d\n", len(alerts))
		for _, alert := range alerts {
			fmt.Printf("  %s (%s): %s - Effective: %s\n",
				alert.ServiceName, alert.DeprecationType, alert.Severity,
				alert.EffectiveDate.Format("2006-01-02"))
		}
	}

	// Test 8: Generate Migration Plan
	fmt.Println("\n--- Test 8: Generate Migration Plan ---")
	if len(alerts) > 0 {
		plan, err := intelligenceService.GenerateMigrationPlan(ctx, alerts[0].ServiceName, clientContext)
		if err != nil {
			log.Printf("Error generating migration plan: %v", err)
		} else {
			fmt.Printf("Migration from %s to %s\n", plan.DeprecatedService, plan.TargetService)
			fmt.Printf("Strategy: %s\n", plan.MigrationStrategy)
			fmt.Printf("Duration: %s\n", plan.EstimatedDuration)
			fmt.Printf("Cost: $%.2f\n", plan.EstimatedCost)
			fmt.Printf("Risk Level: %s\n", plan.RiskLevel)
			fmt.Printf("Migration Steps: %d\n", len(plan.MigrationSteps))
		}
	}

	// Test 9: Analyze Pricing Changes
	fmt.Println("\n--- Test 9: Analyze Pricing Changes ---")
	since = time.Now().Add(-30 * 24 * time.Hour) // Last 30 days
	pricingChanges, err := intelligenceService.AnalyzePricingChanges(ctx, since)
	if err != nil {
		log.Printf("Error analyzing pricing changes: %v", err)
	} else {
		fmt.Printf("Pricing Changes Count: %d\n", len(pricingChanges))
		for _, change := range pricingChanges {
			fmt.Printf("  %s (%s): %s - Impact: %s\n",
				change.ServiceName, change.Region, change.ChangeType, change.ImpactLevel)
		}
	}

	// Test 10: Calculate Cost Impact
	fmt.Println("\n--- Test 10: Calculate Cost Impact ---")
	if len(pricingChanges) > 0 {
		costImpact, err := intelligenceService.CalculateCostImpact(ctx, pricingChanges[0], clientContext)
		if err != nil {
			log.Printf("Error calculating cost impact: %v", err)
		} else {
			fmt.Printf("Service: %s\n", costImpact.ServiceName)
			fmt.Printf("Current Monthly Cost: $%.2f\n", costImpact.CurrentMonthlyCost)
			fmt.Printf("New Monthly Cost: $%.2f\n", costImpact.NewMonthlyCost)
			fmt.Printf("Cost Difference: $%.2f (%.1f%%)\n",
				costImpact.CostDifference, costImpact.PercentageChange)
			fmt.Printf("Annual Impact: $%.2f\n", costImpact.AnnualImpact)
			fmt.Printf("Impact Category: %s\n", costImpact.ImpactCategory)
			fmt.Printf("Action Required: %t\n", costImpact.ActionRequired)
		}
	}

	// Test 11: Health Check
	fmt.Println("\n--- Test 11: Health Check ---")
	fmt.Printf("Service Healthy: %t\n", intelligenceService.IsHealthy())
	fmt.Printf("Last Update: %s\n", intelligenceService.GetLastUpdateTime().Format("2006-01-02 15:04:05"))

	// Test 12: Refresh Intelligence
	fmt.Println("\n--- Test 12: Refresh Intelligence ---")
	err = intelligenceService.RefreshServiceIntelligence(ctx)
	if err != nil {
		log.Printf("Error refreshing intelligence: %v", err)
	} else {
		fmt.Println("Intelligence refreshed successfully")
		fmt.Printf("New Last Update: %s\n", intelligenceService.GetLastUpdateTime().Format("2006-01-02 15:04:05"))
	}

	fmt.Println("\n=== AWS Service Intelligence Service Test Complete ===")
}
