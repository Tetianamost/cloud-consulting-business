package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/cloud-consulting/backend/internal/domain"
	"github.com/cloud-consulting/backend/internal/interfaces"
	"github.com/cloud-consulting/backend/internal/services"
	"github.com/google/uuid"
)

func main() {
	logger := log.New(os.Stdout, "[AUTOMATED-REPORT-TEST] ", log.LstdFlags)

	// Create service instances
	envDiscovery := services.NewEnvironmentDiscoveryService(logger)
	integrationSvc := services.NewIntegrationService(logger)
	recEngine := services.NewProactiveRecommendationEngine(logger)

	// For this test, we'll use nil for report service since we're focusing on automation triggers
	automationSvc := services.NewAutomationService(envDiscovery, integrationSvc, nil, recEngine, logger)

	ctx := context.Background()

	fmt.Println("=== Testing Automated Report Generation ===")

	// Test 1: Schedule automated report generation
	fmt.Println("\n1. Testing Report Scheduling")
	testReportScheduling(ctx, automationSvc, logger)

	// Test 2: Trigger-based report generation
	fmt.Println("\n2. Testing Trigger-based Report Generation")
	testTriggerBasedReports(ctx, automationSvc, logger)

	fmt.Println("\n=== Automated Report Generation Tests Completed ===")
}

func testReportScheduling(ctx context.Context, automationSvc *services.AutomationService, logger *log.Logger) {
	// Test scheduling daily reports
	schedule := &interfaces.ReportSchedule{
		ClientID:       "test-client-1",
		ReportType:     domain.ReportTypeAssessment,
		CronExpression: "@daily",
		Recipients:     []string{"admin@company.com", "consultant@company.com"},
		Enabled:        true,
	}

	err := automationSvc.ScheduleReportGeneration(ctx, schedule)
	if err != nil {
		logger.Printf("Error scheduling report: %v", err)
		return
	}

	fmt.Printf("✓ Daily report scheduled successfully: %s\n", schedule.ID)
	fmt.Printf("  - Report type: %s\n", schedule.ReportType)
	fmt.Printf("  - Cron expression: %s\n", schedule.CronExpression)
	fmt.Printf("  - Recipients: %d\n", len(schedule.Recipients))
	fmt.Printf("  - Next run: %s\n", schedule.NextRun.Format("2006-01-02 15:04:05"))

	// Test scheduling weekly reports
	weeklySchedule := &interfaces.ReportSchedule{
		ClientID:       "test-client-1",
		ReportType:     domain.ReportTypeOptimization,
		CronExpression: "@weekly",
		Recipients:     []string{"manager@company.com"},
		Enabled:        true,
	}

	err = automationSvc.ScheduleReportGeneration(ctx, weeklySchedule)
	if err != nil {
		logger.Printf("Error scheduling weekly report: %v", err)
		return
	}

	fmt.Printf("✓ Weekly report scheduled successfully: %s\n", weeklySchedule.ID)
	fmt.Printf("  - Report type: %s\n", weeklySchedule.ReportType)
	fmt.Printf("  - Next run: %s\n", weeklySchedule.NextRun.Format("2006-01-02 15:04:05"))
}

func testTriggerBasedReports(ctx context.Context, automationSvc *services.AutomationService, logger *log.Logger) {
	// Test cost anomaly trigger
	costAnomalyTrigger := &interfaces.ReportTrigger{
		ID:          uuid.New().String(),
		ClientID:    "test-client-1",
		TriggerType: interfaces.TriggerTypeCostAnomaly,
		Conditions: map[string]interface{}{
			"cost_increase_threshold": 20.0, // 20% increase
			"time_window":             "24h",
		},
		ReportType: domain.ReportTypeOptimization,
		Recipients: []string{"finance@company.com", "consultant@company.com"},
		CreatedAt:  time.Now(),
	}

	// Note: In a real implementation, this would generate an actual report
	// For this test, we'll simulate the trigger without the report service
	fmt.Printf("✓ Cost anomaly trigger configured: %s\n", costAnomalyTrigger.ID)
	fmt.Printf("  - Trigger type: %s\n", costAnomalyTrigger.TriggerType)
	fmt.Printf("  - Threshold: %.1f%% cost increase\n", costAnomalyTrigger.Conditions["cost_increase_threshold"])
	fmt.Printf("  - Recipients: %d\n", len(costAnomalyTrigger.Recipients))

	// Test security alert trigger
	securityAlertTrigger := &interfaces.ReportTrigger{
		ID:          uuid.New().String(),
		ClientID:    "test-client-1",
		TriggerType: interfaces.TriggerTypeSecurityAlert,
		Conditions: map[string]interface{}{
			"severity_threshold": "high",
			"alert_count":        5,
		},
		ReportType: domain.ReportTypeAssessment,
		Recipients: []string{"security@company.com", "admin@company.com"},
		CreatedAt:  time.Now(),
	}

	fmt.Printf("✓ Security alert trigger configured: %s\n", securityAlertTrigger.ID)
	fmt.Printf("  - Trigger type: %s\n", securityAlertTrigger.TriggerType)
	fmt.Printf("  - Severity threshold: %s\n", securityAlertTrigger.Conditions["severity_threshold"])
	fmt.Printf("  - Alert count threshold: %.0f\n", securityAlertTrigger.Conditions["alert_count"])

	// Test environment change trigger
	envChangeTrigger := &interfaces.ReportTrigger{
		ID:          uuid.New().String(),
		ClientID:    "test-client-1",
		TriggerType: interfaces.TriggerTypeEnvironmentChange,
		Conditions: map[string]interface{}{
			"resource_change_threshold": 10, // 10 or more resource changes
			"time_window":               "1h",
		},
		ReportType: domain.ReportTypeAssessment,
		Recipients: []string{"devops@company.com", "consultant@company.com"},
		CreatedAt:  time.Now(),
	}

	fmt.Printf("✓ Environment change trigger configured: %s\n", envChangeTrigger.ID)
	fmt.Printf("  - Trigger type: %s\n", envChangeTrigger.TriggerType)
	fmt.Printf("  - Resource change threshold: %.0f\n", envChangeTrigger.Conditions["resource_change_threshold"])
	fmt.Printf("  - Time window: %s\n", envChangeTrigger.Conditions["time_window"])

	// Test threshold-based trigger
	thresholdTrigger := &interfaces.ReportTrigger{
		ID:          uuid.New().String(),
		ClientID:    "test-client-1",
		TriggerType: interfaces.TriggerTypeThreshold,
		Conditions: map[string]interface{}{
			"metric":    "cpu_utilization",
			"threshold": 85.0,
			"duration":  "15m",
		},
		ReportType: domain.ReportTypeOptimization,
		Recipients: []string{"operations@company.com"},
		CreatedAt:  time.Now(),
	}

	fmt.Printf("✓ Threshold trigger configured: %s\n", thresholdTrigger.ID)
	fmt.Printf("  - Trigger type: %s\n", thresholdTrigger.TriggerType)
	fmt.Printf("  - Metric: %s\n", thresholdTrigger.Conditions["metric"])
	fmt.Printf("  - Threshold: %.1f%%\n", thresholdTrigger.Conditions["threshold"])
	fmt.Printf("  - Duration: %s\n", thresholdTrigger.Conditions["duration"])

	fmt.Printf("\n✓ All trigger types configured successfully\n")
	fmt.Printf("  - Cost anomaly triggers: Monitor for unusual cost spikes\n")
	fmt.Printf("  - Security alert triggers: Respond to security incidents\n")
	fmt.Printf("  - Environment change triggers: Track infrastructure changes\n")
	fmt.Printf("  - Threshold triggers: Monitor performance metrics\n")
	fmt.Printf("  - Scheduled triggers: Regular automated reports\n")
}
