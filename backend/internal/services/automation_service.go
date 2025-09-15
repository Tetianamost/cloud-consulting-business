package services

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/cloud-consulting/backend/internal/domain"
	"github.com/cloud-consulting/backend/internal/interfaces"
	"github.com/google/uuid"
)

// AutomationService implements advanced automation and integration capabilities
type AutomationService struct {
	environmentDiscovery interfaces.EnvironmentDiscoveryService
	integrationService   interfaces.IntegrationService
	reportService        interfaces.ReportService
	recommendationEngine interfaces.ProactiveRecommendationEngine
	logger               *log.Logger
}

// NewAutomationService creates a new automation service
func NewAutomationService(
	envDiscovery interfaces.EnvironmentDiscoveryService,
	integrationSvc interfaces.IntegrationService,
	reportSvc interfaces.ReportService,
	recEngine interfaces.ProactiveRecommendationEngine,
	logger *log.Logger,
) *AutomationService {
	return &AutomationService{
		environmentDiscovery: envDiscovery,
		integrationService:   integrationSvc,
		reportService:        reportSvc,
		recommendationEngine: recEngine,
		logger:               logger,
	}
}

// DiscoverClientEnvironment performs automated discovery of client cloud environment
func (a *AutomationService) DiscoverClientEnvironment(ctx context.Context, clientID string, credentials *interfaces.ClientCredentials) (*interfaces.EnvironmentDiscovery, error) {
	a.logger.Printf("Starting environment discovery for client: %s", clientID)

	var snapshot *interfaces.EnvironmentSnapshot

	// Perform provider-specific discovery
	switch strings.ToLower(credentials.Provider) {
	case "aws":
		awsCreds := &interfaces.AWSCredentials{
			AccessKeyID:     credentials.Credentials["access_key_id"].(string),
			SecretAccessKey: credentials.Credentials["secret_access_key"].(string),
			Region:          credentials.Region,
		}
		if sessionToken, ok := credentials.Credentials["session_token"]; ok {
			awsCreds.SessionToken = sessionToken.(string)
		}
		awsSnapshot, err := a.environmentDiscovery.ScanAWSEnvironment(ctx, awsCreds)
		if err != nil {
			return nil, fmt.Errorf("failed to scan AWS environment: %w", err)
		}
		snapshot = awsSnapshot.EnvironmentSnapshot
	case "azure":
		azureCreds := &interfaces.AzureCredentials{
			ClientID:       credentials.Credentials["client_id"].(string),
			ClientSecret:   credentials.Credentials["client_secret"].(string),
			TenantID:       credentials.Credentials["tenant_id"].(string),
			SubscriptionID: credentials.Credentials["subscription_id"].(string),
		}
		azureSnapshot, err := a.environmentDiscovery.ScanAzureEnvironment(ctx, azureCreds)
		if err != nil {
			return nil, fmt.Errorf("failed to scan Azure environment: %w", err)
		}
		snapshot = azureSnapshot.EnvironmentSnapshot
	case "gcp":
		gcpCreds := &interfaces.GCPCredentials{
			ProjectID:           credentials.Credentials["project_id"].(string),
			ServiceAccountKey:   credentials.Credentials["service_account_key"].(string),
			ServiceAccountEmail: credentials.Credentials["service_account_email"].(string),
		}
		gcpSnapshot, err := a.environmentDiscovery.ScanGCPEnvironment(ctx, gcpCreds)
		if err != nil {
			return nil, fmt.Errorf("failed to scan GCP environment: %w", err)
		}
		snapshot = gcpSnapshot.EnvironmentSnapshot
	default:
		return nil, fmt.Errorf("unsupported cloud provider: %s", credentials.Provider)
	}

	// Convert snapshot to discovery format
	discovery := &interfaces.EnvironmentDiscovery{
		ClientID:     clientID,
		Provider:     credentials.Provider,
		DiscoveredAt: time.Now(),
		Resources:    a.convertResourcesToDiscovered(snapshot.Resources),
		Services:     a.extractServices(snapshot),
		Configurations: map[string]interface{}{
			"provider":   credentials.Provider,
			"region":     credentials.Region,
			"account_id": credentials.AccountID,
			"scan_time":  snapshot.SnapshotTime,
		},
		CostEstimate:     a.estimateCosts(snapshot),
		SecurityFindings: a.analyzeSecurityFindings(snapshot),
		Recommendations:  a.generateAutomatedRecommendations(snapshot),
	}

	a.logger.Printf("Environment discovery completed for client: %s, found %d resources", clientID, len(discovery.Resources))
	return discovery, nil
}

// AnalyzeEnvironmentChanges compares environment snapshots and analyzes changes
func (a *AutomationService) AnalyzeEnvironmentChanges(ctx context.Context, clientID string, previousSnapshot, currentSnapshot *interfaces.EnvironmentSnapshot) (*interfaces.ChangeAnalysis, error) {
	a.logger.Printf("Analyzing environment changes for client: %s", clientID)

	analysis := &interfaces.ChangeAnalysis{
		ClientID:     clientID,
		AnalysisTime: time.Now(),
		TimeRange: interfaces.TimeRange{
			StartDate: previousSnapshot.SnapshotTime,
			EndDate:   currentSnapshot.SnapshotTime,
		},
	}

	// Analyze resource changes
	analysis.AddedResources = a.findAddedResources(previousSnapshot.Resources, currentSnapshot.Resources)
	analysis.ModifiedResources = a.findModifiedResources(previousSnapshot.Resources, currentSnapshot.Resources)
	analysis.DeletedResources = a.findDeletedResources(previousSnapshot.Resources, currentSnapshot.Resources)

	// Analyze cost impact
	analysis.CostImpact = a.analyzeCostImpact(previousSnapshot.Costs, currentSnapshot.Costs)

	// Analyze security impact
	analysis.SecurityImpact = a.analyzeSecurityImpact(analysis.AddedResources, analysis.ModifiedResources, analysis.DeletedResources)

	// Generate change-based recommendations
	analysis.Recommendations = a.generateChangeRecommendations(analysis)

	a.logger.Printf("Change analysis completed: %d added, %d modified, %d deleted resources",
		len(analysis.AddedResources), len(analysis.ModifiedResources), len(analysis.DeletedResources))

	return analysis, nil
}

// RegisterIntegration registers a new third-party integration
func (a *AutomationService) RegisterIntegration(ctx context.Context, integration *interfaces.Integration) error {
	a.logger.Printf("Registering integration: %s for client: %s", integration.Name, integration.ClientID)

	integration.ID = uuid.New().String()
	integration.Status = interfaces.IntegrationStatusPending
	integration.CreatedAt = time.Now()
	integration.UpdatedAt = time.Now()

	// Test the integration before registering
	testResult, err := a.testIntegrationConfig(ctx, integration)
	if err != nil {
		integration.Status = interfaces.IntegrationStatusError
		return fmt.Errorf("integration test failed: %w", err)
	}

	if testResult.Success {
		integration.Status = interfaces.IntegrationStatusActive
	} else {
		integration.Status = interfaces.IntegrationStatusError
		return fmt.Errorf("integration test failed: %s", testResult.Message)
	}

	// Store integration configuration (in a real implementation, this would persist to database)
	a.logger.Printf("Integration registered successfully: %s", integration.ID)
	return nil
}

// GetIntegrations retrieves all integrations for a client
func (a *AutomationService) GetIntegrations(ctx context.Context, clientID string) ([]*interfaces.Integration, error) {
	// In a real implementation, this would query the database
	// For now, return mock data
	return []*interfaces.Integration{
		{
			ID:       uuid.New().String(),
			ClientID: clientID,
			Type:     interfaces.IntegrationTypeMonitoring,
			Name:     "CloudWatch Integration",
			Status:   interfaces.IntegrationStatusActive,
			Configuration: map[string]interface{}{
				"region":     "us-east-1",
				"log_groups": []string{"/aws/lambda/function1", "/aws/ec2/instance1"},
			},
			CreatedAt: time.Now().Add(-24 * time.Hour),
			UpdatedAt: time.Now().Add(-1 * time.Hour),
		},
	}, nil
}

// TestIntegration tests an existing integration
func (a *AutomationService) TestIntegration(ctx context.Context, integrationID string) (*interfaces.IntegrationTestResult, error) {
	a.logger.Printf("Testing integration: %s", integrationID)

	// In a real implementation, this would retrieve the integration and test it
	return &interfaces.IntegrationTestResult{
		IntegrationID: integrationID,
		Success:       true,
		Message:       "Integration test successful",
		TestedAt:      time.Now(),
		ResponseTime:  150,
		Details: map[string]interface{}{
			"connection_status": "healthy",
			"data_sync_status":  "up_to_date",
		},
	}, nil
}

// GenerateAutomatedReport generates reports based on triggers
func (a *AutomationService) GenerateAutomatedReport(ctx context.Context, trigger *interfaces.ReportTrigger) (*domain.Report, error) {
	a.logger.Printf("Generating automated report for trigger: %s", trigger.ID)

	// Create inquiry based on trigger conditions
	inquiry := &domain.Inquiry{
		ID:        uuid.New().String(),
		Company:   fmt.Sprintf("Client-%s", trigger.ClientID),
		Services:  []string{"automated-analysis"},
		Message:   a.buildTriggerMessage(trigger),
		CreatedAt: time.Now(),
	}

	// Generate report using existing report service
	report, err := a.reportService.GenerateReport(ctx, inquiry)
	if err != nil {
		return nil, fmt.Errorf("failed to generate automated report: %w", err)
	}

	// Add automation-specific metadata
	report.GeneratedBy = "automation-service"
	report.UpdatedAt = time.Now()

	a.logger.Printf("Automated report generated: %s", report.ID)
	return report, nil
}

// ScheduleReportGeneration schedules automated report generation
func (a *AutomationService) ScheduleReportGeneration(ctx context.Context, schedule *interfaces.ReportSchedule) error {
	a.logger.Printf("Scheduling report generation: %s for client: %s", schedule.ReportType, schedule.ClientID)

	schedule.ID = uuid.New().String()
	schedule.CreatedAt = time.Now()
	schedule.UpdatedAt = time.Now()

	// Calculate next run time based on cron expression
	schedule.NextRun = a.calculateNextRun(schedule.CronExpression)

	// In a real implementation, this would store the schedule and set up a cron job
	a.logger.Printf("Report schedule created: %s, next run: %s", schedule.ID, schedule.NextRun.Format(time.RFC3339))
	return nil
}

// GenerateProactiveRecommendations generates proactive recommendations based on usage patterns
func (a *AutomationService) GenerateProactiveRecommendations(ctx context.Context, clientID string, usagePatterns *interfaces.UsagePatterns) ([]*interfaces.ProactiveRecommendation, error) {
	a.logger.Printf("Generating proactive recommendations for client: %s", clientID)

	var recommendations []*interfaces.ProactiveRecommendation

	// Generate cost optimization recommendations
	if usagePatterns.CostTrends != nil {
		costRecommendations := a.generateCostOptimizationRecommendations(clientID, usagePatterns.CostTrends)
		recommendations = append(recommendations, costRecommendations...)
	}

	// Generate performance recommendations
	if usagePatterns.PerformanceMetrics != nil {
		perfRecommendations := a.generatePerformanceRecommendations(clientID, usagePatterns.PerformanceMetrics)
		recommendations = append(recommendations, perfRecommendations...)
	}

	// Generate security recommendations
	if len(usagePatterns.AutomationSecurityEvents) > 0 {
		secRecommendations := a.generateSecurityRecommendations(clientID, usagePatterns.AutomationSecurityEvents)
		recommendations = append(recommendations, secRecommendations...)
	}

	// Generate anomaly-based recommendations
	if len(usagePatterns.AnomaliesDetected) > 0 {
		anomalyRecommendations := a.generateAnomalyRecommendations(clientID, usagePatterns.AnomaliesDetected)
		recommendations = append(recommendations, anomalyRecommendations...)
	}

	a.logger.Printf("Generated %d proactive recommendations for client: %s", len(recommendations), clientID)
	return recommendations, nil
}

// AnalyzeUsagePatterns analyzes client usage patterns over a time range
func (a *AutomationService) AnalyzeUsagePatterns(ctx context.Context, clientID string, timeRange interfaces.AutomationTimeRange) (*interfaces.UsagePatterns, error) {
	a.logger.Printf("Analyzing usage patterns for client: %s from %s to %s",
		clientID, timeRange.StartTime.Format("2006-01-02"), timeRange.EndTime.Format("2006-01-02"))

	// In a real implementation, this would query metrics from various sources
	patterns := &interfaces.UsagePatterns{
		ClientID:  clientID,
		TimeRange: timeRange,
		CostTrends: &interfaces.AutomationCostTrends{
			TotalCost:     15000.50,
			PreviousCost:  14200.30,
			PercentChange: 5.6,
			Trend:         "increasing",
			ServiceBreakdown: map[string]float64{
				"EC2":    8500.25,
				"RDS":    3200.15,
				"S3":     1800.10,
				"Lambda": 1500.00,
			},
			DailyTrends: a.generateMockDailyCosts(timeRange),
		},
		ResourceUtilization: &interfaces.AutomationResourceUtilization{
			CPU: &interfaces.CPUUtilization{
				Average: 65.5,
				Maximum: 89.2,
				Minimum: 12.3,
			},
			Memory: &interfaces.MemoryUtilization{
				Average: 72.1,
				Maximum: 95.8,
				Minimum: 25.4,
			},
		},
		PerformanceMetrics: &interfaces.AutomationPerformanceMetrics{
			ResponseTime: &interfaces.ResponseTimeMetrics{
				Average: 245.5,
				P95:     890.2,
				P99:     1250.8,
			},
			Throughput: &interfaces.ThroughputMetrics{
				RequestsPerSecond:     1250.5,
				TransactionsPerSecond: 890.2,
			},
			ErrorRate:    0.02,
			Availability: 99.95,
		},
		AnomaliesDetected: a.generateMockAnomalies(clientID, timeRange),
	}

	return patterns, nil
}

// Helper methods

func (a *AutomationService) convertResourcesToDiscovered(resources []*interfaces.ResourceSnapshot) []*interfaces.DiscoveredResource {
	var discovered []*interfaces.DiscoveredResource
	for _, resource := range resources {
		discovered = append(discovered, &interfaces.DiscoveredResource{
			ID:            resource.ResourceID,
			Type:          resource.Type,
			Name:          resource.ResourceID,
			Configuration: resource.Configuration,
			Cost: &interfaces.ResourceCost{
				MonthlyCost: 100.0, // Mock cost
				Currency:    "USD",
			},
		})
	}
	return discovered
}

func (a *AutomationService) extractServices(snapshot *interfaces.EnvironmentSnapshot) []*interfaces.DiscoveredService {
	// Extract unique service types from resources
	serviceMap := make(map[string]*interfaces.DiscoveredService)
	for _, resource := range snapshot.Resources {
		if _, exists := serviceMap[resource.Type]; !exists {
			serviceMap[resource.Type] = &interfaces.DiscoveredService{
				Name: resource.Type,
				Type: "cloud-service",
				Configuration: map[string]interface{}{
					"provider": "aws", // This would be dynamic
				},
			}
		}
	}

	var services []*interfaces.DiscoveredService
	for _, service := range serviceMap {
		services = append(services, service)
	}
	return services
}

func (a *AutomationService) estimateCosts(snapshot *interfaces.EnvironmentSnapshot) *interfaces.CostEstimate {
	// Mock cost estimation based on resource count
	resourceCount := float64(len(snapshot.Resources))
	monthlyCost := resourceCount * 150.0 // $150 per resource average

	return &interfaces.CostEstimate{
		MonthlyCost: monthlyCost,
		AnnualCost:  monthlyCost * 12,
		Currency:    "USD",
		Breakdown: map[string]float64{
			"compute": monthlyCost * 0.6,
			"storage": monthlyCost * 0.2,
			"network": monthlyCost * 0.2,
		},
		LastUpdated: time.Now(),
	}
}

func (a *AutomationService) analyzeSecurityFindings(snapshot *interfaces.EnvironmentSnapshot) []*interfaces.SecurityFinding {
	// Mock security analysis
	return []*interfaces.SecurityFinding{
		{
			ID:          uuid.New().String(),
			Type:        "access_control",
			Severity:    "medium",
			Title:       "Overly Permissive IAM Policies",
			Description: "Some IAM policies grant broader permissions than necessary",
			Resource:    "iam-policy-123",
			Remediation: "Review and apply principle of least privilege",
		},
	}
}

func (a *AutomationService) generateAutomatedRecommendations(snapshot *interfaces.EnvironmentSnapshot) []*interfaces.AutomatedRecommendation {
	return []*interfaces.AutomatedRecommendation{
		{
			Type:        "cost_optimization",
			Title:       "Right-size EC2 Instances",
			Description: "Several EC2 instances are underutilized and can be downsized",
			Priority:    "medium",
			ActionItems: []string{
				"Analyze CPU and memory utilization over 30 days",
				"Identify instances with <20% average utilization",
				"Test workload on smaller instance types",
			},
			EstimatedSavings: 2500.00,
		},
	}
}

func (a *AutomationService) findAddedResources(previous, current []*interfaces.ResourceSnapshot) []*interfaces.ResourceChange {
	previousMap := make(map[string]*interfaces.ResourceSnapshot)
	for _, resource := range previous {
		previousMap[resource.ResourceID] = resource
	}

	var added []*interfaces.ResourceChange
	for _, resource := range current {
		if _, exists := previousMap[resource.ResourceID]; !exists {
			added = append(added, &interfaces.ResourceChange{
				ResourceID: resource.ResourceID,
				ChangeType: "added",
				NewState:   resource.Configuration,
				Impact:     "new_resource_added",
				Timestamp:  time.Now(),
			})
		}
	}
	return added
}

func (a *AutomationService) findModifiedResources(previous, current []*interfaces.ResourceSnapshot) []*interfaces.ResourceChange {
	previousMap := make(map[string]*interfaces.ResourceSnapshot)
	for _, resource := range previous {
		previousMap[resource.ResourceID] = resource
	}

	var modified []*interfaces.ResourceChange
	for _, resource := range current {
		if prevResource, exists := previousMap[resource.ResourceID]; exists {
			// Simple comparison - in reality, this would be more sophisticated
			if resource.State != prevResource.State {
				modified = append(modified, &interfaces.ResourceChange{
					ResourceID: resource.ResourceID,
					ChangeType: "modified",
					OldState:   prevResource.Configuration,
					NewState:   resource.Configuration,
					Impact:     "configuration_changed",
					Timestamp:  time.Now(),
				})
			}
		}
	}
	return modified
}

func (a *AutomationService) findDeletedResources(previous, current []*interfaces.ResourceSnapshot) []*interfaces.ResourceChange {
	currentMap := make(map[string]*interfaces.ResourceSnapshot)
	for _, resource := range current {
		currentMap[resource.ResourceID] = resource
	}

	var deleted []*interfaces.ResourceChange
	for _, resource := range previous {
		if _, exists := currentMap[resource.ResourceID]; !exists {
			deleted = append(deleted, &interfaces.ResourceChange{
				ResourceID: resource.ResourceID,
				ChangeType: "deleted",
				OldState:   resource.Configuration,
				Impact:     "resource_removed",
				Timestamp:  time.Now(),
			})
		}
	}
	return deleted
}

func (a *AutomationService) analyzeCostImpact(previousCosts, currentCosts *interfaces.CostSnapshot) *interfaces.AutomationCostImpactAnalysis {
	if previousCosts == nil || currentCosts == nil {
		return &interfaces.AutomationCostImpactAnalysis{
			TotalImpact: 0,
			Currency:    "USD",
			Trend:       "stable",
		}
	}

	impact := currentCosts.TotalCost - previousCosts.TotalCost
	var trend string
	if impact > 0 {
		trend = "increasing"
	} else if impact < 0 {
		trend = "decreasing"
	} else {
		trend = "stable"
	}

	return &interfaces.AutomationCostImpactAnalysis{
		TotalImpact: impact,
		Currency:    currentCosts.Currency,
		ImpactBreakdown: map[string]float64{
			"compute": impact * 0.6,
			"storage": impact * 0.3,
			"network": impact * 0.1,
		},
		Trend: trend,
		Recommendations: []string{
			"Monitor cost trends closely",
			"Review resource utilization",
			"Consider reserved instances for stable workloads",
		},
	}
}

func (a *AutomationService) analyzeSecurityImpact(added, modified, deleted []*interfaces.ResourceChange) *interfaces.SecurityImpactAnalysis {
	riskLevel := "low"
	if len(added) > 5 || len(modified) > 10 {
		riskLevel = "medium"
	}
	if len(deleted) > 3 {
		riskLevel = "high"
	}

	return &interfaces.SecurityImpactAnalysis{
		RiskLevel: riskLevel,
		NewVulnerabilities: []string{
			"New resources may not follow security baselines",
			"Modified configurations may introduce security gaps",
		},
		ResolvedIssues: []string{
			"Deleted resources reduce attack surface",
		},
		Recommendations: []string{
			"Review security configurations for new resources",
			"Validate compliance of modified resources",
			"Update security monitoring for environment changes",
		},
	}
}

func (a *AutomationService) generateChangeRecommendations(analysis *interfaces.ChangeAnalysis) []*interfaces.ChangeRecommendation {
	var recommendations []*interfaces.ChangeRecommendation

	if len(analysis.AddedResources) > 0 {
		recommendations = append(recommendations, &interfaces.ChangeRecommendation{
			Type:        "security",
			Title:       "Review New Resource Security",
			Description: "Ensure new resources follow security best practices",
			Priority:    "high",
			Actions: []string{
				"Review security groups and access policies",
				"Enable logging and monitoring",
				"Validate encryption settings",
			},
			Impact: "Reduces security risks for new resources",
		})
	}

	if analysis.CostImpact.TotalImpact > 1000 {
		recommendations = append(recommendations, &interfaces.ChangeRecommendation{
			Type:        "cost",
			Title:       "Monitor Cost Impact",
			Description: "Significant cost increase detected from recent changes",
			Priority:    "medium",
			Actions: []string{
				"Set up cost alerts",
				"Review resource sizing",
				"Consider cost optimization opportunities",
			},
			Impact: "Helps control and optimize costs",
		})
	}

	return recommendations
}

func (a *AutomationService) testIntegrationConfig(ctx context.Context, integration *interfaces.Integration) (*interfaces.IntegrationTestResult, error) {
	// Mock integration testing
	return &interfaces.IntegrationTestResult{
		IntegrationID: integration.ID,
		Success:       true,
		Message:       "Integration configuration is valid",
		TestedAt:      time.Now(),
		ResponseTime:  200,
		Details: map[string]interface{}{
			"config_valid": true,
			"connectivity": "ok",
		},
	}, nil
}

func (a *AutomationService) buildTriggerMessage(trigger *interfaces.ReportTrigger) string {
	switch trigger.TriggerType {
	case interfaces.TriggerTypeScheduled:
		return "Scheduled automated analysis and recommendations"
	case interfaces.TriggerTypeThreshold:
		return "Threshold-based analysis triggered by system conditions"
	case interfaces.TriggerTypeEnvironmentChange:
		return "Environment change analysis and impact assessment"
	case interfaces.TriggerTypeCostAnomaly:
		return "Cost anomaly detected - analysis and recommendations"
	case interfaces.TriggerTypeSecurityAlert:
		return "Security alert triggered - analysis and remediation recommendations"
	default:
		return "Automated analysis and recommendations"
	}
}

func (a *AutomationService) calculateNextRun(cronExpression string) time.Time {
	// Simple implementation - in reality, would use a proper cron parser
	switch cronExpression {
	case "@daily":
		return time.Now().Add(24 * time.Hour)
	case "@weekly":
		return time.Now().Add(7 * 24 * time.Hour)
	case "@monthly":
		return time.Now().Add(30 * 24 * time.Hour)
	default:
		return time.Now().Add(24 * time.Hour)
	}
}

func (a *AutomationService) generateCostOptimizationRecommendations(clientID string, costTrends *interfaces.AutomationCostTrends) []*interfaces.ProactiveRecommendation {
	var recommendations []*interfaces.ProactiveRecommendation

	if costTrends.PercentChange > 10 {
		recommendations = append(recommendations, &interfaces.ProactiveRecommendation{
			ID:               uuid.New().String(),
			ClientID:         clientID,
			Type:             interfaces.RecommendationTypeCostOptimization,
			Title:            "High Cost Increase Detected",
			Description:      fmt.Sprintf("Costs have increased by %.1f%% over the analysis period", costTrends.PercentChange),
			Priority:         interfaces.RecommendationPriorityHigh,
			Impact:           "Potential cost savings of $2,000-5,000 monthly",
			Effort:           "medium",
			PotentialSavings: costTrends.TotalCost * 0.15,
			ActionItems: []string{
				"Review resource utilization metrics",
				"Identify underutilized resources",
				"Consider reserved instances for stable workloads",
				"Implement automated scaling policies",
			},
			Resources: []string{
				"AWS Cost Explorer",
				"CloudWatch metrics",
				"Trusted Advisor recommendations",
			},
			CreatedAt: time.Now(),
			ExpiresAt: &[]time.Time{time.Now().Add(30 * 24 * time.Hour)}[0],
		})
	}

	return recommendations
}

func (a *AutomationService) generatePerformanceRecommendations(clientID string, perfMetrics *interfaces.AutomationPerformanceMetrics) []*interfaces.ProactiveRecommendation {
	var recommendations []*interfaces.ProactiveRecommendation

	if perfMetrics.ResponseTime.P95 > 1000 {
		recommendations = append(recommendations, &interfaces.ProactiveRecommendation{
			ID:          uuid.New().String(),
			ClientID:    clientID,
			Type:        interfaces.RecommendationTypePerformance,
			Title:       "High Response Time Detected",
			Description: fmt.Sprintf("95th percentile response time is %.1fms, exceeding recommended thresholds", perfMetrics.ResponseTime.P95),
			Priority:    interfaces.RecommendationPriorityMedium,
			Impact:      "Improved user experience and system performance",
			Effort:      "high",
			ActionItems: []string{
				"Analyze application performance bottlenecks",
				"Review database query performance",
				"Consider implementing caching strategies",
				"Optimize network latency",
			},
			Resources: []string{
				"Application Performance Monitoring tools",
				"Database performance insights",
				"CDN configuration",
			},
			CreatedAt: time.Now(),
		})
	}

	return recommendations
}

func (a *AutomationService) generateSecurityRecommendations(clientID string, securityEvents []*interfaces.AutomationSecurityEvent) []*interfaces.ProactiveRecommendation {
	var recommendations []*interfaces.ProactiveRecommendation

	if len(securityEvents) > 10 {
		recommendations = append(recommendations, &interfaces.ProactiveRecommendation{
			ID:          uuid.New().String(),
			ClientID:    clientID,
			Type:        interfaces.RecommendationTypeSecurity,
			Title:       "Increased Security Events",
			Description: fmt.Sprintf("Detected %d security events in the analysis period", len(securityEvents)),
			Priority:    interfaces.RecommendationPriorityHigh,
			Impact:      "Enhanced security posture and reduced risk",
			Effort:      "medium",
			ActionItems: []string{
				"Review security event patterns",
				"Update security policies and rules",
				"Implement additional monitoring",
				"Conduct security assessment",
			},
			Resources: []string{
				"Security monitoring dashboards",
				"Incident response procedures",
				"Security best practices documentation",
			},
			CreatedAt: time.Now(),
		})
	}

	return recommendations
}

func (a *AutomationService) generateAnomalyRecommendations(clientID string, anomalies []*interfaces.Anomaly) []*interfaces.ProactiveRecommendation {
	var recommendations []*interfaces.ProactiveRecommendation

	for _, anomaly := range anomalies {
		if anomaly.Severity == "high" || anomaly.Severity == "critical" {
			recommendations = append(recommendations, &interfaces.ProactiveRecommendation{
				ID:          uuid.New().String(),
				ClientID:    clientID,
				Type:        interfaces.RecommendationTypeOperational,
				Title:       fmt.Sprintf("Anomaly Detected: %s", anomaly.Type),
				Description: anomaly.Description,
				Priority:    interfaces.RecommendationPriorityHigh,
				Impact:      "Prevents potential system issues",
				Effort:      "low",
				ActionItems: []string{
					"Investigate anomaly root cause",
					"Review related system metrics",
					"Implement preventive measures",
				},
				Resources: []string{
					"System monitoring tools",
					"Historical performance data",
				},
				CreatedAt: time.Now(),
			})
		}
	}

	return recommendations
}

func (a *AutomationService) generateMockDailyCosts(timeRange interfaces.AutomationTimeRange) []interfaces.DailyCost {
	var dailyCosts []interfaces.DailyCost
	current := timeRange.StartTime
	baseCost := 450.0

	for current.Before(timeRange.EndTime) {
		// Add some variation to make it realistic
		variation := float64(current.Day()%7) * 20.0
		dailyCosts = append(dailyCosts, interfaces.DailyCost{
			Date: current,
			Cost: baseCost + variation,
		})
		current = current.Add(24 * time.Hour)
	}

	return dailyCosts
}

func (a *AutomationService) generateMockAnomalies(clientID string, timeRange interfaces.AutomationTimeRange) []*interfaces.Anomaly {
	return []*interfaces.Anomaly{
		{
			ID:          uuid.New().String(),
			Type:        "cost_spike",
			Severity:    "medium",
			Timestamp:   timeRange.EndTime.Add(-2 * time.Hour),
			Description: "Unusual cost increase detected in EC2 service",
			Value:       1250.0,
			Expected:    850.0,
			Deviation:   47.1,
			Metadata: map[string]interface{}{
				"service":     "EC2",
				"region":      "us-east-1",
				"instance_id": "i-1234567890abcdef0",
			},
		},
	}
}
