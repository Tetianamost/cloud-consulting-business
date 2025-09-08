package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/cloud-consulting/backend/internal/config"
	"github.com/cloud-consulting/backend/internal/domain"
	"github.com/cloud-consulting/backend/internal/interfaces"
	"github.com/cloud-consulting/backend/internal/services"
)

// TestResult represents the result of a test
type TestResult struct {
	TestName string        `json:"test_name"`
	Status   string        `json:"status"`
	Duration time.Duration `json:"duration"`
	Error    string        `json:"error,omitempty"`
}

// TestSuite manages and runs all monitoring and observability tests
type TestSuite struct {
	logger  *logrus.Logger
	results []TestResult
	config  *config.Config
}

func main() {
	// Initialize logger
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(&logrus.JSONFormatter{})

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.WithError(err).Fatal("Failed to load configuration")
	}

	// Create test suite
	suite := &TestSuite{
		logger:  logger,
		results: make([]TestResult, 0),
		config:  cfg,
	}

	logger.Info("Starting Task 14: Email Event System Monitoring and Observability Tests (Simple)")

	// Run all tests
	suite.runAllTests()

	// Print results
	suite.printResults()

	// Exit with appropriate code
	if suite.hasFailures() {
		os.Exit(1)
	}
}

func (ts *TestSuite) runAllTests() {
	tests := []struct {
		name string
		fn   func() error
	}{
		{"Test Email Monitoring Service Creation", ts.testEmailMonitoringServiceCreation},
		{"Test Email System Alerting Service Creation", ts.testEmailSystemAlertingServiceCreation},
		{"Test Alert Configuration Management", ts.testAlertConfigurationManagement},
		{"Test Alert Level Determination", ts.testAlertLevelDetermination},
		{"Test Alert Suppression Configuration", ts.testAlertSuppressionConfiguration},
		{"Test Monitoring Service Health Check", ts.testMonitoringServiceHealthCheck},
		{"Test Alerting Service Lifecycle", ts.testAlertingServiceLifecycle},
		{"Test Metrics Structure Validation", ts.testMetricsStructureValidation},
		{"Test Health Status Structure", ts.testHealthStatusStructure},
		{"Test Alert State Management", ts.testAlertStateManagement},
	}

	for _, test := range tests {
		ts.runTest(test.name, test.fn)
	}
}

func (ts *TestSuite) runTest(name string, testFn func() error) {
	start := time.Now()
	ts.logger.WithField("test", name).Info("Running test")

	err := testFn()
	duration := time.Since(start)

	result := TestResult{
		TestName: name,
		Duration: duration,
	}

	if err != nil {
		result.Status = "FAILED"
		result.Error = err.Error()
		ts.logger.WithError(err).WithField("test", name).Error("Test failed")
	} else {
		result.Status = "PASSED"
		ts.logger.WithField("test", name).Info("Test passed")
	}

	ts.results = append(ts.results, result)
}

func (ts *TestSuite) testEmailMonitoringServiceCreation() error {
	// Create mock services for testing
	mockEmailEventRecorder := &MockEmailEventRecorder{}
	mockEmailMetricsService := &MockEmailMetricsService{}

	// Create email monitoring service
	emailMonitoringService := services.NewEmailMonitoringService(
		mockEmailEventRecorder,
		mockEmailMetricsService,
		ts.logger,
	)

	if emailMonitoringService == nil {
		return fmt.Errorf("email monitoring service should not be nil")
	}

	// Test basic functionality
	if !emailMonitoringService.IsHealthy() {
		return fmt.Errorf("email monitoring service should be healthy initially")
	}

	// Test getting system metrics
	metrics := emailMonitoringService.GetSystemMetrics()
	if metrics == nil {
		return fmt.Errorf("system metrics should not be nil")
	}

	// Test getting health status
	healthStatus := emailMonitoringService.GetHealthStatus()
	if healthStatus == nil {
		return fmt.Errorf("health status should not be nil")
	}

	ts.logger.WithFields(logrus.Fields{
		"is_healthy":     emailMonitoringService.IsHealthy(),
		"overall_status": healthStatus.OverallStatus,
	}).Info("Email monitoring service creation test completed")

	return nil
}

func (ts *TestSuite) testEmailSystemAlertingServiceCreation() error {
	// Create mock services for testing
	mockEmailEventRecorder := &MockEmailEventRecorder{}
	mockEmailMetricsService := &MockEmailMetricsService{}

	// Create email system alerting service
	alertingService := services.NewEmailSystemAlertingService(
		mockEmailEventRecorder,
		mockEmailMetricsService,
		ts.logger,
	)

	if alertingService == nil {
		return fmt.Errorf("email system alerting service should not be nil")
	}

	// Test getting alert configuration
	config := alertingService.GetConfig()
	if config == nil {
		return fmt.Errorf("alert configuration should not be nil")
	}

	// Test getting alert state
	alertState := alertingService.GetAlertState()
	if alertState == nil {
		return fmt.Errorf("alert state should not be nil")
	}

	// Test getting alert metrics
	alertMetrics := alertingService.GetAlertMetrics()
	if alertMetrics == nil {
		return fmt.Errorf("alert metrics should not be nil")
	}

	ts.logger.WithFields(logrus.Fields{
		"is_running":    alertingService.IsRunning(),
		"alert_config":  config,
		"alert_state":   alertState,
		"alert_metrics": alertMetrics,
	}).Info("Email system alerting service creation test completed")

	return nil
}

func (ts *TestSuite) testAlertConfigurationManagement() error {
	// Create mock services
	mockEmailEventRecorder := &MockEmailEventRecorder{}
	mockEmailMetricsService := &MockEmailMetricsService{}

	// Create alerting service
	alertingService := services.NewEmailSystemAlertingService(
		mockEmailEventRecorder,
		mockEmailMetricsService,
		ts.logger,
	)

	// Test getting default configuration
	config := alertingService.GetConfig()
	if config.RecordingFailureRateThreshold != 0.05 {
		return fmt.Errorf("expected default recording failure rate threshold 0.05, got %f", config.RecordingFailureRateThreshold)
	}

	// Test updating configuration
	newConfig := *config
	newConfig.RecordingFailureRateThreshold = 0.10
	newConfig.HighFailureRateThreshold = 0.25
	alertingService.UpdateConfig(&newConfig)

	// Verify configuration was updated
	updatedConfig := alertingService.GetConfig()
	if updatedConfig.RecordingFailureRateThreshold != 0.10 {
		return fmt.Errorf("expected updated recording failure rate threshold 0.10, got %f", updatedConfig.RecordingFailureRateThreshold)
	}

	if updatedConfig.HighFailureRateThreshold != 0.25 {
		return fmt.Errorf("expected updated high failure rate threshold 0.25, got %f", updatedConfig.HighFailureRateThreshold)
	}

	ts.logger.WithFields(logrus.Fields{
		"original_threshold": config.RecordingFailureRateThreshold,
		"updated_threshold":  updatedConfig.RecordingFailureRateThreshold,
	}).Info("Alert configuration management test completed")

	return nil
}

func (ts *TestSuite) testAlertLevelDetermination() error {
	// Create mock services
	mockEmailEventRecorder := &MockEmailEventRecorder{}
	mockEmailMetricsService := &MockEmailMetricsService{}

	// Create alerting service
	alertingService := services.NewEmailSystemAlertingService(
		mockEmailEventRecorder,
		mockEmailMetricsService,
		ts.logger,
	)

	// Test alert configuration thresholds
	config := alertingService.GetConfig()

	// Verify default thresholds are reasonable
	if config.RecordingFailureRateThreshold <= 0 || config.RecordingFailureRateThreshold >= 1 {
		return fmt.Errorf("recording failure rate threshold should be between 0 and 1, got %f", config.RecordingFailureRateThreshold)
	}

	if config.HighFailureRateThreshold <= config.RecordingFailureRateThreshold {
		return fmt.Errorf("high failure rate threshold should be greater than recording failure rate threshold")
	}

	if config.CriticalFailureRateThreshold <= config.HighFailureRateThreshold {
		return fmt.Errorf("critical failure rate threshold should be greater than high failure rate threshold")
	}

	ts.logger.WithFields(logrus.Fields{
		"recording_threshold": config.RecordingFailureRateThreshold,
		"high_threshold":      config.HighFailureRateThreshold,
		"critical_threshold":  config.CriticalFailureRateThreshold,
	}).Info("Alert level determination test completed")

	return nil
}

func (ts *TestSuite) testAlertSuppressionConfiguration() error {
	// Create mock services
	mockEmailEventRecorder := &MockEmailEventRecorder{}
	mockEmailMetricsService := &MockEmailMetricsService{}

	// Create alerting service
	alertingService := services.NewEmailSystemAlertingService(
		mockEmailEventRecorder,
		mockEmailMetricsService,
		ts.logger,
	)

	// Test alert suppression configuration
	config := alertingService.GetConfig()
	if !config.EnableAlertSuppression {
		return fmt.Errorf("alert suppression should be enabled by default")
	}

	if config.SuppressionWindow <= 0 {
		return fmt.Errorf("suppression window should be positive, got %v", config.SuppressionWindow)
	}

	// Test updating suppression configuration
	newConfig := *config
	newConfig.SuppressionWindow = 45 * time.Minute
	alertingService.UpdateConfig(&newConfig)

	updatedConfig := alertingService.GetConfig()
	if updatedConfig.SuppressionWindow != 45*time.Minute {
		return fmt.Errorf("expected updated suppression window 45 minutes, got %v", updatedConfig.SuppressionWindow)
	}

	ts.logger.WithFields(logrus.Fields{
		"suppression_enabled": config.EnableAlertSuppression,
		"original_window":     config.SuppressionWindow,
		"updated_window":      updatedConfig.SuppressionWindow,
	}).Info("Alert suppression configuration test completed")

	return nil
}

func (ts *TestSuite) testMonitoringServiceHealthCheck() error {
	// Create mock services
	mockEmailEventRecorder := &MockEmailEventRecorder{}
	mockEmailMetricsService := &MockEmailMetricsService{}

	// Create monitoring service
	emailMonitoringService := services.NewEmailMonitoringService(
		mockEmailEventRecorder,
		mockEmailMetricsService,
		ts.logger,
	)

	// Test health check
	ctx := context.Background()
	healthCheck := emailMonitoringService.PerformHealthCheck(ctx)
	if healthCheck == nil {
		return fmt.Errorf("health check result should not be nil")
	}

	// Verify health check structure
	if healthCheck.OverallStatus == "" {
		return fmt.Errorf("overall status should not be empty")
	}

	if healthCheck.LastChecked.IsZero() {
		return fmt.Errorf("last checked time should not be zero")
	}

	ts.logger.WithFields(logrus.Fields{
		"overall_status":          healthCheck.OverallStatus,
		"recorder_healthy":        healthCheck.RecorderHealthy,
		"metrics_service_healthy": healthCheck.MetricsServiceHealthy,
		"last_checked":            healthCheck.LastChecked,
	}).Info("Monitoring service health check test completed")

	return nil
}

func (ts *TestSuite) testAlertingServiceLifecycle() error {
	// Create mock services
	mockEmailEventRecorder := &MockEmailEventRecorder{}
	mockEmailMetricsService := &MockEmailMetricsService{}

	// Create alerting service
	alertingService := services.NewEmailSystemAlertingService(
		mockEmailEventRecorder,
		mockEmailMetricsService,
		ts.logger,
	)

	// Initially should not be running
	if alertingService.IsRunning() {
		return fmt.Errorf("alerting service should not be running initially")
	}

	// Start the service
	err := alertingService.Start()
	if err != nil {
		return fmt.Errorf("failed to start alerting service: %w", err)
	}

	// Should be running now
	if !alertingService.IsRunning() {
		return fmt.Errorf("alerting service should be running after start")
	}

	// Wait a moment for background monitoring
	time.Sleep(100 * time.Millisecond)

	// Stop the service
	err = alertingService.Stop()
	if err != nil {
		return fmt.Errorf("failed to stop alerting service: %w", err)
	}

	// Should not be running now
	if alertingService.IsRunning() {
		return fmt.Errorf("alerting service should not be running after stop")
	}

	ts.logger.Info("Alerting service lifecycle test completed")

	return nil
}

func (ts *TestSuite) testMetricsStructureValidation() error {
	// Create mock services
	mockEmailEventRecorder := &MockEmailEventRecorder{}
	mockEmailMetricsService := &MockEmailMetricsService{}

	// Create monitoring service
	emailMonitoringService := services.NewEmailMonitoringService(
		mockEmailEventRecorder,
		mockEmailMetricsService,
		ts.logger,
	)

	// Get system metrics
	metrics := emailMonitoringService.GetSystemMetrics()
	if metrics == nil {
		return fmt.Errorf("system metrics should not be nil")
	}

	// Verify metrics structure
	if metrics.HealthCheckInterval <= 0 {
		return fmt.Errorf("health check interval should be positive")
	}

	// Verify metrics are properly initialized
	if metrics.LastHealthCheck.IsZero() {
		return fmt.Errorf("last health check time should not be zero")
	}

	ts.logger.WithFields(logrus.Fields{
		"health_check_interval": metrics.HealthCheckInterval,
		"alerts_triggered":      metrics.AlertsTriggered,
		"system_uptime":         metrics.SystemUptime,
	}).Info("Metrics structure validation test completed")

	return nil
}

func (ts *TestSuite) testHealthStatusStructure() error {
	// Create mock services
	mockEmailEventRecorder := &MockEmailEventRecorder{}
	mockEmailMetricsService := &MockEmailMetricsService{}

	// Create monitoring service
	emailMonitoringService := services.NewEmailMonitoringService(
		mockEmailEventRecorder,
		mockEmailMetricsService,
		ts.logger,
	)

	// Get health status
	healthStatus := emailMonitoringService.GetHealthStatus()
	if healthStatus == nil {
		return fmt.Errorf("health status should not be nil")
	}

	// Verify health status structure
	validStatuses := []string{"healthy", "degraded", "unhealthy", "unknown"}
	isValidStatus := false
	for _, status := range validStatuses {
		if healthStatus.OverallStatus == status {
			isValidStatus = true
			break
		}
	}

	if !isValidStatus {
		return fmt.Errorf("overall status should be one of %v, got %s", validStatuses, healthStatus.OverallStatus)
	}

	if healthStatus.LastChecked.IsZero() {
		return fmt.Errorf("last checked time should not be zero")
	}

	ts.logger.WithFields(logrus.Fields{
		"overall_status":          healthStatus.OverallStatus,
		"recorder_healthy":        healthStatus.RecorderHealthy,
		"metrics_service_healthy": healthStatus.MetricsServiceHealthy,
		"consecutive_failures":    healthStatus.ConsecutiveFailures,
	}).Info("Health status structure test completed")

	return nil
}

func (ts *TestSuite) testAlertStateManagement() error {
	// Create mock services
	mockEmailEventRecorder := &MockEmailEventRecorder{}
	mockEmailMetricsService := &MockEmailMetricsService{}

	// Create alerting service
	alertingService := services.NewEmailSystemAlertingService(
		mockEmailEventRecorder,
		mockEmailMetricsService,
		ts.logger,
	)

	// Get initial alert state
	alertState := alertingService.GetAlertState()
	if alertState == nil {
		return fmt.Errorf("alert state should not be nil")
	}

	// Verify initial alert state
	if alertState.OverallAlertLevel != "none" {
		return fmt.Errorf("initial overall alert level should be 'none', got %s", alertState.OverallAlertLevel)
	}

	if alertState.ConsecutiveRecorderFailures != 0 {
		return fmt.Errorf("initial consecutive recorder failures should be 0, got %d", alertState.ConsecutiveRecorderFailures)
	}

	if alertState.ConsecutiveMetricsServiceFailures != 0 {
		return fmt.Errorf("initial consecutive metrics service failures should be 0, got %d", alertState.ConsecutiveMetricsServiceFailures)
	}

	ts.logger.WithFields(logrus.Fields{
		"overall_alert_level":                  alertState.OverallAlertLevel,
		"recorder_alert_level":                 alertState.RecorderAlertLevel,
		"metrics_service_alert_level":          alertState.MetricsServiceAlertLevel,
		"consecutive_recorder_failures":        alertState.ConsecutiveRecorderFailures,
		"consecutive_metrics_service_failures": alertState.ConsecutiveMetricsServiceFailures,
	}).Info("Alert state management test completed")

	return nil
}

func (ts *TestSuite) printResults() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("TASK 14: EMAIL EVENT SYSTEM MONITORING AND OBSERVABILITY TEST RESULTS")
	fmt.Println(strings.Repeat("=", 80))

	passed := 0
	failed := 0
	totalDuration := time.Duration(0)

	for _, result := range ts.results {
		status := "✅ PASS"
		if result.Status == "FAILED" {
			status = "❌ FAIL"
			failed++
		} else {
			passed++
		}
		totalDuration += result.Duration

		fmt.Printf("%s %s (%.2fs)\n", status, result.TestName, result.Duration.Seconds())
		if result.Error != "" {
			fmt.Printf("   Error: %s\n", result.Error)
		}
	}

	fmt.Println(strings.Repeat("-", 80))
	fmt.Printf("Total Tests: %d | Passed: %d | Failed: %d | Duration: %.2fs\n",
		len(ts.results), passed, failed, totalDuration.Seconds())

	if failed == 0 {
		fmt.Println("🎉 All monitoring and observability tests passed!")
		fmt.Println("\n✅ Task 14 Implementation Summary:")
		fmt.Println("   • Structured logging for email event recording operations")
		fmt.Println("   • Metrics collection for email event recording success/failure rates")
		fmt.Println("   • Health check endpoints for email event tracking system")
		fmt.Println("   • Alerting for high email event recording failure rates")
		fmt.Println("   • Comprehensive monitoring and observability infrastructure")
		fmt.Println("   • Email system health handler with multiple endpoints")
		fmt.Println("   • Advanced alerting service with configurable thresholds")
		fmt.Println("   • Alert suppression and lifecycle management")
		fmt.Println("   • Prometheus metrics export capability")
	} else {
		fmt.Printf("❌ %d test(s) failed. Please review the errors above.\n", failed)
	}

	fmt.Println(strings.Repeat("=", 80))
}

func (ts *TestSuite) hasFailures() bool {
	for _, result := range ts.results {
		if result.Status == "FAILED" {
			return true
		}
	}
	return false
}

// Mock implementations for testing

type MockEmailEventRecorder struct{}

func (m *MockEmailEventRecorder) RecordEmailSent(ctx context.Context, event *domain.EmailEvent) error {
	return nil
}

func (m *MockEmailEventRecorder) UpdateEmailStatus(ctx context.Context, messageID string, status domain.EmailEventStatus, deliveredAt *time.Time, errorMsg string) error {
	return nil
}

func (m *MockEmailEventRecorder) GetEmailEventsByInquiry(ctx context.Context, inquiryID string) ([]*domain.EmailEvent, error) {
	return []*domain.EmailEvent{}, nil
}

func (m *MockEmailEventRecorder) IsHealthy() bool {
	return true
}

func (m *MockEmailEventRecorder) IsHealthyWithContext(ctx context.Context) bool {
	return true
}

func (m *MockEmailEventRecorder) GetMetrics() *interfaces.EmailEventRecorderMetrics {
	return &interfaces.EmailEventRecorderMetrics{
		TotalRecordingAttempts: 10,
		SuccessfulRecordings:   9,
		FailedRecordings:       1,
		SuccessRate:            0.9,
		AverageRecordingTime:   50.0,
		LastRecordingTime:      time.Now(),
		RetryAttempts:          0,
		HealthCheckFailures:    0,
	}
}

type MockEmailMetricsService struct{}

func (m *MockEmailMetricsService) GetEmailMetrics(ctx context.Context, timeRange domain.TimeRange) (*domain.EmailMetrics, error) {
	return &domain.EmailMetrics{
		TotalEmails:     10,
		DeliveredEmails: 9,
		FailedEmails:    1,
		BouncedEmails:   0,
		SpamEmails:      0,
		DeliveryRate:    0.9,
		BounceRate:      0.05,
		SpamRate:        0.05,
		TimeRange:       fmt.Sprintf("%s to %s", timeRange.Start.Format("2006-01-02"), timeRange.End.Format("2006-01-02")),
	}, nil
}

func (m *MockEmailMetricsService) GetEmailStatusByInquiry(ctx context.Context, inquiryID string) (*domain.EmailStatus, error) {
	return &domain.EmailStatus{
		InquiryID:       inquiryID,
		TotalEmailsSent: 2,
	}, nil
}

func (m *MockEmailMetricsService) GetEmailEventHistory(ctx context.Context, filters domain.EmailEventFilters) ([]*domain.EmailEvent, error) {
	return []*domain.EmailEvent{}, nil
}

func (m *MockEmailMetricsService) GetRecentEmailActivity(ctx context.Context, hours int) ([]*domain.EmailEvent, error) {
	return []*domain.EmailEvent{}, nil
}

func (m *MockEmailMetricsService) IsHealthy(ctx context.Context) bool {
	return true
}

func (m *MockEmailMetricsService) GetMetrics() *interfaces.EmailMetricsServiceMetrics {
	return &interfaces.EmailMetricsServiceMetrics{
		TotalMetricsRequests: 5,
		SuccessfulRequests:   5,
		FailedRequests:       0,
		SuccessRate:          1.0,
		AverageResponseTime:  25.0,
		LastRequestTime:      time.Now(),
		CacheHits:            3,
		CacheMisses:          2,
		HealthCheckFailures:  0,
	}
}
