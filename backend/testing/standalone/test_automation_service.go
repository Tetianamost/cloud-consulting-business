package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/cloud-consulting/backend/internal/interfaces"
	"github.com/cloud-consulting/backend/internal/services"
)

func main() {
	logger := log.New(os.Stdout, "[AUTOMATION-TEST] ", log.LstdFlags)

	// Create service instances
	envDiscovery := services.NewEnvironmentDiscoveryService(logger)
	integrationSvc := services.NewIntegrationService(logger)
	recEngine := services.NewProactiveRecommendationEngine(logger)

	// For this test, we'll use nil for report generator since we're focusing on automation
	automationSvc := services.NewAutomationService(envDiscovery, integrationSvc, nil, recEngine, logger)

	ctx := context.Background()

	fmt.Println("=== Testing Automation Service ===")

	// Test 1: Environment Discovery
	fmt.Println("\n1. Testing Environment Discovery")
	testEnvironmentDiscovery(ctx, automationSvc, logger)

	// Test 2: Integration Management
	fmt.Println("\n2. Testing Integration Management")
	testIntegrationManagement(ctx, automationSvc, logger)

	// Test 3: Proactive Recommendations
	fmt.Println("\n3. Testing Proactive Recommendations")
	testProactiveRecommendations(ctx, automationSvc, logger)

	// Test 4: Usage Pattern Analysis
	fmt.Println("\n4. Testing Usage Pattern Analysis")
	testUsagePatternAnalysis(ctx, automationSvc, logger)

	// Test 5: Environment Change Analysis
	fmt.Println("\n5. Testing Environment Change Analysis")
	testEnvironmentChangeAnalysis(ctx, automationSvc, logger)

	fmt.Println("\n=== All Automation Tests Completed ===")
}

func testEnvironmentDiscovery(ctx context.Context, automationSvc *services.AutomationService, logger *log.Logger) {
	// Test AWS environment discovery
	awsCredentials := &interfaces.ClientCredentials{
		ClientID:  "test-client-1",
		Provider:  "aws",
		Region:    "us-east-1",
		AccountID: "123456789012",
		Credentials: map[string]interface{}{
			"access_key_id":     "AKIAIOSFODNN7EXAMPLE",
			"secret_access_key": "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
		},
	}

	discovery, err := automationSvc.DiscoverClientEnvironment(ctx, "test-client-1", awsCredentials)
	if err != nil {
		logger.Printf("Error discovering AWS environment: %v", err)
		return
	}

	fmt.Printf("✓ AWS Environment Discovery completed\n")
	fmt.Printf("  - Resources discovered: %d\n", len(discovery.Resources))
	fmt.Printf("  - Services discovered: %d\n", len(discovery.Services))
	fmt.Printf("  - Estimated monthly cost: $%.2f\n", discovery.CostEstimate.MonthlyCost)
	fmt.Printf("  - Security findings: %d\n", len(discovery.SecurityFindings))
	fmt.Printf("  - Recommendations: %d\n", len(discovery.Recommendations))

	// Test Azure environment discovery
	azureCredentials := &interfaces.ClientCredentials{
		ClientID: "test-client-2",
		Provider: "azure",
		Region:   "East US",
		Credentials: map[string]interface{}{
			"client_id":       "12345678-1234-1234-1234-123456789012",
			"client_secret":   "example-secret",
			"tenant_id":       "87654321-4321-4321-4321-210987654321",
			"subscription_id": "11111111-2222-3333-4444-555555555555",
		},
	}

	discovery, err = automationSvc.DiscoverClientEnvironment(ctx, "test-client-2", azureCredentials)
	if err != nil {
		logger.Printf("Error discovering Azure environment: %v", err)
		return
	}

	fmt.Printf("✓ Azure Environment Discovery completed\n")
	fmt.Printf("  - Resources discovered: %d\n", len(discovery.Resources))
	fmt.Printf("  - Services discovered: %d\n", len(discovery.Services))
	fmt.Printf("  - Estimated monthly cost: $%.2f\n", discovery.CostEstimate.MonthlyCost)
}

func testIntegrationManagement(ctx context.Context, automationSvc *services.AutomationService, logger *log.Logger) {
	// Test integration registration
	integration := &interfaces.Integration{
		ClientID: "test-client-1",
		Type:     interfaces.IntegrationTypeMonitoring,
		Name:     "Test CloudWatch Integration",
		Configuration: map[string]interface{}{
			"region":            "us-east-1",
			"access_key_id":     "test-key",
			"secret_access_key": "test-secret",
			"log_groups":        []string{"/aws/lambda/test"},
		},
	}

	err := automationSvc.RegisterIntegration(ctx, integration)
	if err != nil {
		logger.Printf("Error registering integration: %v", err)
		return
	}

	fmt.Printf("✓ Integration registered successfully: %s\n", integration.ID)
	fmt.Printf("  - Type: %s\n", integration.Type)
	fmt.Printf("  - Status: %s\n", integration.Status)

	// Test getting integrations
	integrations, err := automationSvc.GetIntegrations(ctx, "test-client-1")
	if err != nil {
		logger.Printf("Error getting integrations: %v", err)
		return
	}

	fmt.Printf("✓ Retrieved %d integrations for client\n", len(integrations))
	for _, integ := range integrations {
		fmt.Printf("  - %s (%s): %s\n", integ.Name, integ.Type, integ.Status)
	}

	// Test integration testing
	if len(integrations) > 0 {
		testResult, err := automationSvc.TestIntegration(ctx, integrations[0].ID)
		if err != nil {
			logger.Printf("Error testing integration: %v", err)
			return
		}

		fmt.Printf("✓ Integration test completed\n")
		fmt.Printf("  - Success: %t\n", testResult.Success)
		fmt.Printf("  - Response time: %dms\n", testResult.ResponseTime)
		fmt.Printf("  - Message: %s\n", testResult.Message)
	}
}

func testProactiveRecommendations(ctx context.Context, automationSvc *services.AutomationService, logger *log.Logger) {
	// Create mock usage patterns
	usagePatterns := &interfaces.UsagePatterns{
		ClientID: "test-client-1",
		TimeRange: interfaces.AutomationTimeRange{
			StartTime: time.Now().Add(-30 * 24 * time.Hour),
			EndTime:   time.Now(),
		},
		CostTrends: &interfaces.AutomationCostTrends{
			TotalCost:     15000.50,
			PreviousCost:  13500.25,
			PercentChange: 11.1,
			Trend:         "increasing",
			ServiceBreakdown: map[string]float64{
				"EC2": 8500.25,
				"RDS": 3200.15,
				"S3":  1800.10,
			},
		},
		PerformanceMetrics: &interfaces.AutomationPerformanceMetrics{
			ResponseTime: &interfaces.ResponseTimeMetrics{
				Average: 350.5,
				P95:     1250.8,
				P99:     2100.2,
			},
			ErrorRate:    0.025,
			Availability: 99.85,
		},
		AutomationSecurityEvents: []*interfaces.AutomationSecurityEvent{
			{
				ID:          "sec-event-1",
				Type:        "unauthorized_access_attempt",
				Severity:    "high",
				Timestamp:   time.Now().Add(-2 * time.Hour),
				Description: "Multiple failed login attempts detected",
			},
		},
		AnomaliesDetected: []*interfaces.Anomaly{
			{
				ID:          "anomaly-1",
				Type:        "cost_spike",
				Severity:    "medium",
				Timestamp:   time.Now().Add(-6 * time.Hour),
				Description: "Unusual cost increase in EC2 service",
				Value:       1250.0,
				Expected:    850.0,
				Deviation:   47.1,
			},
		},
	}

	recommendations, err := automationSvc.GenerateProactiveRecommendations(ctx, "test-client-1", usagePatterns)
	if err != nil {
		logger.Printf("Error generating proactive recommendations: %v", err)
		return
	}

	fmt.Printf("✓ Generated %d proactive recommendations\n", len(recommendations))
	for i, rec := range recommendations {
		fmt.Printf("  %d. %s (%s priority)\n", i+1, rec.Title, rec.Priority)
		fmt.Printf("     Type: %s\n", rec.Type)
		fmt.Printf("     Impact: %s\n", rec.Impact)
		if rec.PotentialSavings > 0 {
			fmt.Printf("     Potential savings: $%.2f\n", rec.PotentialSavings)
		}
		fmt.Printf("     Actions: %d items\n", len(rec.ActionItems))
	}
}

func testUsagePatternAnalysis(ctx context.Context, automationSvc *services.AutomationService, logger *log.Logger) {
	timeRange := interfaces.AutomationTimeRange{
		StartTime: time.Now().Add(-30 * 24 * time.Hour),
		EndTime:   time.Now(),
	}

	patterns, err := automationSvc.AnalyzeUsagePatterns(ctx, "test-client-1", timeRange)
	if err != nil {
		logger.Printf("Error analyzing usage patterns: %v", err)
		return
	}

	fmt.Printf("✓ Usage pattern analysis completed\n")
	fmt.Printf("  - Analysis period: %s to %s\n",
		patterns.TimeRange.StartTime.Format("2006-01-02"),
		patterns.TimeRange.EndTime.Format("2006-01-02"))

	if patterns.CostTrends != nil {
		fmt.Printf("  - Total cost: $%.2f (%.1f%% change)\n",
			patterns.CostTrends.TotalCost, patterns.CostTrends.PercentChange)
		fmt.Printf("  - Cost trend: %s\n", patterns.CostTrends.Trend)
		fmt.Printf("  - Service breakdown: %d services\n", len(patterns.CostTrends.ServiceBreakdown))
	}

	if patterns.ResourceUtilization != nil {
		if patterns.ResourceUtilization.CPU != nil {
			fmt.Printf("  - CPU utilization: %.1f%% avg, %.1f%% max\n",
				patterns.ResourceUtilization.CPU.Average, patterns.ResourceUtilization.CPU.Maximum)
		}
		if patterns.ResourceUtilization.Memory != nil {
			fmt.Printf("  - Memory utilization: %.1f%% avg, %.1f%% max\n",
				patterns.ResourceUtilization.Memory.Average, patterns.ResourceUtilization.Memory.Maximum)
		}
	}

	if patterns.PerformanceMetrics != nil {
		fmt.Printf("  - Response time: %.1fms avg, %.1fms P95\n",
			patterns.PerformanceMetrics.ResponseTime.Average, patterns.PerformanceMetrics.ResponseTime.P95)
		fmt.Printf("  - Error rate: %.2f%%\n", patterns.PerformanceMetrics.ErrorRate*100)
		fmt.Printf("  - Availability: %.2f%%\n", patterns.PerformanceMetrics.Availability)
	}

	fmt.Printf("  - Anomalies detected: %d\n", len(patterns.AnomaliesDetected))
}

func testEnvironmentChangeAnalysis(ctx context.Context, automationSvc *services.AutomationService, logger *log.Logger) {
	// Create mock snapshots
	previousSnapshot := &interfaces.EnvironmentSnapshot{
		ID:           "snapshot-1",
		Provider:     "aws",
		SnapshotTime: time.Now().Add(-24 * time.Hour),
		Resources: []*interfaces.ResourceSnapshot{
			{
				ResourceID: "i-1234567890abcdef0",
				Type:       "EC2Instance",
				State:      "running",
				Configuration: map[string]interface{}{
					"instance_type": "t3.medium",
				},
			},
			{
				ResourceID: "db-prod-instance",
				Type:       "RDSInstance",
				State:      "available",
				Configuration: map[string]interface{}{
					"instance_class": "db.t3.medium",
				},
			},
		},
		Costs: &interfaces.CostSnapshot{
			TotalCost: 1200.50,
			Currency:  "USD",
			Breakdown: map[string]float64{
				"EC2": 800.25,
				"RDS": 400.25,
			},
		},
	}

	currentSnapshot := &interfaces.EnvironmentSnapshot{
		ID:           "snapshot-2",
		Provider:     "aws",
		SnapshotTime: time.Now(),
		Resources: []*interfaces.ResourceSnapshot{
			{
				ResourceID: "i-1234567890abcdef0",
				Type:       "EC2Instance",
				State:      "running",
				Configuration: map[string]interface{}{
					"instance_type": "t3.large", // Changed from medium to large
				},
			},
			{
				ResourceID: "db-prod-instance",
				Type:       "RDSInstance",
				State:      "available",
				Configuration: map[string]interface{}{
					"instance_class": "db.t3.medium",
				},
			},
			{
				ResourceID: "i-new-instance-123",
				Type:       "EC2Instance",
				State:      "running",
				Configuration: map[string]interface{}{
					"instance_type": "t3.small", // New instance
				},
			},
		},
		Costs: &interfaces.CostSnapshot{
			TotalCost: 1450.75,
			Currency:  "USD",
			Breakdown: map[string]float64{
				"EC2": 950.50,
				"RDS": 400.25,
				"S3":  100.00, // New service
			},
		},
	}

	analysis, err := automationSvc.AnalyzeEnvironmentChanges(ctx, "test-client-1", previousSnapshot, currentSnapshot)
	if err != nil {
		logger.Printf("Error analyzing environment changes: %v", err)
		return
	}

	fmt.Printf("✓ Environment change analysis completed\n")
	fmt.Printf("  - Analysis period: %s to %s\n",
		analysis.TimeRange.StartDate.Format("2006-01-02 15:04:05"),
		analysis.TimeRange.EndDate.Format("2006-01-02 15:04:05"))
	fmt.Printf("  - Added resources: %d\n", len(analysis.AddedResources))
	fmt.Printf("  - Modified resources: %d\n", len(analysis.ModifiedResources))
	fmt.Printf("  - Deleted resources: %d\n", len(analysis.DeletedResources))

	if analysis.CostImpact != nil {
		fmt.Printf("  - Cost impact: $%.2f (%s trend)\n",
			analysis.CostImpact.TotalImpact, analysis.CostImpact.Trend)
	}

	if analysis.SecurityImpact != nil {
		fmt.Printf("  - Security risk level: %s\n", analysis.SecurityImpact.RiskLevel)
		fmt.Printf("  - New vulnerabilities: %d\n", len(analysis.SecurityImpact.NewVulnerabilities))
		fmt.Printf("  - Security recommendations: %d\n", len(analysis.SecurityImpact.Recommendations))
	}

	fmt.Printf("  - Change recommendations: %d\n", len(analysis.Recommendations))
	for i, rec := range analysis.Recommendations {
		fmt.Printf("    %d. %s (%s priority)\n", i+1, rec.Title, rec.Priority)
		fmt.Printf("       Type: %s\n", rec.Type)
		fmt.Printf("       Actions: %d items\n", len(rec.Actions))
	}
}

// Helper function to pretty print JSON
func prettyPrintJSON(v interface{}) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		return
	}
	fmt.Println(string(b))
}
