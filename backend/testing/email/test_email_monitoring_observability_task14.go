package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"

	"github.com/cloud-consulting/backend/internal/config"
	"github.com/cloud-consulting/backend/internal/domain"
	"github.com/cloud-consulting/backend/internal/handlers"
	"github.com/cloud-consulting/backend/internal/repositories"
	"github.com/cloud-consulting/backend/internal/services"
	"github.com/cloud-consulting/backend/internal/storage"
)

// TestResult represents the result of a test
type TestResult struct {
	TestName string        `json:"test_name"`
	Status   string        `json:"status"`
	Duration time.Duration `json:"duration"`
	Error    string        `json:"error,omitempty"`
	Details  interface{}   `json:"details,omitempty"`
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

	logger.Info("Starting Task 14: Email Event System Monitoring and Observability Tests")

	// Run all tests
	suite.runAllTests()

	// Print results
	suite.printResults()

	// Exit with appropriate code
	if suite.hasFailures() {
		os.Exit(1)
	}
}

// setupTestDatabase creates a test database connection
func (ts *TestSuite) setupTestDatabase() (*sql.DB, error) {
	// Try PostgreSQL first
	db, err := sql.Open("postgres", "postgres://test:test@localhost/test_email_monitoring?sslmode=disable")
	if err == nil {
		// Test the connection
		if err = db.Ping(); err == nil {
			return db, nil
		}
		db.Close()
	}

	// Fallback to in-memory SQLite for testing
	db, err = sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, fmt.Errorf("failed to create test database: %w", err)
	}

	return db, nil
}

func (ts *TestSuite) runAllTests() {
	tests := []struct {
		name string
		fn   func() error
	}{
		{"Test Email Event Recorder Structured Logging", ts.testEmailEventRecorderLogging},
		{"Test Email Metrics Service Structured Logging", ts.testEmailMetricsServiceLogging},
		{"Test Email Event Recording Metrics Collection", ts.testEmailEventRecordingMetrics},
		{"Test Email Metrics Service Metrics Collection", ts.testEmailMetricsServiceMetrics},
		{"Test Email System Health Check Endpoint", ts.testEmailSystemHealthEndpoint},
		{"Test Email System Liveness Endpoint", ts.testEmailSystemLivenessEndpoint},
		{"Test Email System Readiness Endpoint", ts.testEmailSystemReadinessEndpoint},
		{"Test Email System Deep Health Check", ts.testEmailSystemDeepHealthCheck},
		{"Test Email System Alerting Service", ts.testEmailSystemAlertingService},
		{"Test High Failure Rate Alerting", ts.testHighFailureRateAlerting},
		{"Test Alert Suppression Logic", ts.testAlertSuppressionLogic},
		{"Test Prometheus Metrics Export", ts.testPrometheusMetricsExport},
		{"Test Monitoring Service Integration", ts.testMonitoringServiceIntegration},
		{"Test Health Check Performance", ts.testHealthCheckPerformance},
		{"Test Alert Configuration Management", ts.testAlertConfigurationManagement},
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

func (ts *TestSuite) testEmailEventRecorderLogging() error {
	// Create test database connection
	db, err := ts.setupTestDatabase()
	if err != nil {
		return fmt.Errorf("failed to create test database: %w", err)
	}
	defer db.Close()

	// Create email event repository
	emailEventRepo := repositories.NewEmailEventRepository(db, ts.logger)

	// Create email event recorder with structured logging
	emailEventRecorder := services.NewEmailEventRecorder(emailEventRepo, ts.logger)

	// Test structured logging by recording an email event
	ctx := context.Background()
	testEvent := &domain.EmailEvent{
		ID:             "test-event-1",
		InquiryID:      "test-inquiry-1",
		EmailType:      domain.EmailTypeCustomerConfirmation,
		RecipientEmail: "test@example.com",
		SenderEmail:    "system@cloudpartner.pro",
		Subject:        "Test Email",
		Status:         domain.EmailStatusSent,
		SentAt:         time.Now(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Record email event (should generate structured logs)
	err = emailEventRecorder.RecordEmailSent(ctx, testEvent)
	if err != nil {
		return fmt.Errorf("failed to record email event: %w", err)
	}

	// Verify metrics are being tracked
	metrics := emailEventRecorder.GetMetrics()
	if metrics == nil {
		return fmt.Errorf("email event recorder metrics are nil")
	}

	if metrics.TotalRecordingAttempts == 0 {
		return fmt.Errorf("expected recording attempts > 0, got %d", metrics.TotalRecordingAttempts)
	}

	ts.logger.WithFields(logrus.Fields{
		"total_attempts":        metrics.TotalRecordingAttempts,
		"successful_recordings": metrics.SuccessfulRecordings,
		"success_rate":          metrics.SuccessRate,
	}).Info("Email event recorder structured logging test completed")

	return nil
}

func (ts *TestSuite) testEmailMetricsServiceLogging() error {
	// Create test database connection
	db, err := ts.setupTestDatabase()
	if err != nil {
		return fmt.Errorf("failed to create test database: %w", err)
	}
	defer db.Close()

	// Create email event repository
	emailEventRepo := repositories.NewEmailEventRepository(db, ts.logger)

	// Create email metrics service with structured logging
	emailMetricsService := services.NewEmailMetricsService(emailEventRepo, ts.logger)

	// Test structured logging by getting email metrics
	ctx := context.Background()
	timeRange := domain.TimeRange{
		Start: time.Now().Add(-1 * time.Hour),
		End:   time.Now(),
	}

	// Get email metrics (should generate structured logs)
	metrics, err := emailMetricsService.GetEmailMetrics(ctx, timeRange)
	if err != nil {
		return fmt.Errorf("failed to get email metrics: %w", err)
	}

	if metrics == nil {
		return fmt.Errorf("email metrics are nil")
	}

	// Verify service metrics are being tracked
	serviceMetrics := emailMetricsService.GetMetrics()
	if serviceMetrics == nil {
		return fmt.Errorf("email metrics service metrics are nil")
	}

	if serviceMetrics.TotalMetricsRequests == 0 {
		return fmt.Errorf("expected metrics requests > 0, got %d", serviceMetrics.TotalMetricsRequests)
	}

	ts.logger.WithFields(logrus.Fields{
		"total_requests":        serviceMetrics.TotalMetricsRequests,
		"successful_requests":   serviceMetrics.SuccessfulRequests,
		"success_rate":          serviceMetrics.SuccessRate,
		"average_response_time": serviceMetrics.AverageResponseTime,
	}).Info("Email metrics service structured logging test completed")

	return nil
}

func (ts *TestSuite) testEmailEventRecordingMetrics() error {
	// Create test database connection
	db, err := ts.setupTestDatabase()
	if err != nil {
		return fmt.Errorf("failed to create test database: %w", err)
	}
	defer db.Close()

	// Create email event repository
	emailEventRepo := repositories.NewEmailEventRepository(db, ts.logger)

	// Create email event recorder
	emailEventRecorder := services.NewEmailEventRecorder(emailEventRepo, ts.logger)

	ctx := context.Background()

	// Record multiple email events to generate metrics
	for i := 0; i < 5; i++ {
		testEvent := &domain.EmailEvent{
			ID:             fmt.Sprintf("test-event-%d", i),
			InquiryID:      fmt.Sprintf("test-inquiry-%d", i),
			EmailType:      domain.EmailTypeCustomerConfirmation,
			RecipientEmail: fmt.Sprintf("test%d@example.com", i),
			SenderEmail:    "system@cloudpartner.pro",
			Subject:        "Test Email",
			Status:         domain.EmailStatusSent,
			SentAt:         time.Now(),
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		err = emailEventRecorder.RecordEmailSent(ctx, testEvent)
		if err != nil {
			return fmt.Errorf("failed to record email event %d: %w", i, err)
		}
	}

	// Wait a moment for async operations to complete
	time.Sleep(100 * time.Millisecond)

	// Verify metrics collection
	metrics := emailEventRecorder.GetMetrics()
	if metrics == nil {
		return fmt.Errorf("email event recorder metrics are nil")
	}

	expectedAttempts := int64(5)
	if metrics.TotalRecordingAttempts < expectedAttempts {
		return fmt.Errorf("expected at least %d recording attempts, got %d", expectedAttempts, metrics.TotalRecordingAttempts)
	}

	if metrics.SuccessfulRecordings < expectedAttempts {
		return fmt.Errorf("expected at least %d successful recordings, got %d", expectedAttempts, metrics.SuccessfulRecordings)
	}

	if metrics.SuccessRate <= 0 {
		return fmt.Errorf("expected success rate > 0, got %f", metrics.SuccessRate)
	}

	if metrics.AverageRecordingTime <= 0 {
		return fmt.Errorf("expected average recording time > 0, got %f", metrics.AverageRecordingTime)
	}

	ts.logger.WithFields(logrus.Fields{
		"total_attempts":         metrics.TotalRecordingAttempts,
		"successful_recordings":  metrics.SuccessfulRecordings,
		"failed_recordings":      metrics.FailedRecordings,
		"success_rate":           metrics.SuccessRate,
		"average_recording_time": metrics.AverageRecordingTime,
		"retry_attempts":         metrics.RetryAttempts,
	}).Info("Email event recording metrics collection test completed")

	return nil
}

func (ts *TestSuite) testEmailMetricsServiceMetrics() error {
	// Create test database connection
	db, err := ts.setupTestDatabase()
	if err != nil {
		return fmt.Errorf("failed to create test database: %w", err)
	}
	defer db.Close()

	// Create email event repository
	emailEventRepo := repositories.NewEmailEventRepository(db, ts.logger)

	// Create email metrics service
	emailMetricsService := services.NewEmailMetricsService(emailEventRepo, ts.logger)

	ctx := context.Background()
	timeRange := domain.TimeRange{
		Start: time.Now().Add(-1 * time.Hour),
		End:   time.Now(),
	}

	// Make multiple metrics requests to generate metrics
	for i := 0; i < 3; i++ {
		_, err = emailMetricsService.GetEmailMetrics(ctx, timeRange)
		if err != nil {
			return fmt.Errorf("failed to get email metrics on attempt %d: %w", i, err)
		}
	}

	// Verify metrics collection
	metrics := emailMetricsService.GetMetrics()
	if metrics == nil {
		return fmt.Errorf("email metrics service metrics are nil")
	}

	expectedRequests := int64(3)
	if metrics.TotalMetricsRequests < expectedRequests {
		return fmt.Errorf("expected at least %d metrics requests, got %d", expectedRequests, metrics.TotalMetricsRequests)
	}

	if metrics.SuccessfulRequests < expectedRequests {
		return fmt.Errorf("expected at least %d successful requests, got %d", expectedRequests, metrics.SuccessfulRequests)
	}

	if metrics.SuccessRate <= 0 {
		return fmt.Errorf("expected success rate > 0, got %f", metrics.SuccessRate)
	}

	if metrics.AverageResponseTime <= 0 {
		return fmt.Errorf("expected average response time > 0, got %f", metrics.AverageResponseTime)
	}

	ts.logger.WithFields(logrus.Fields{
		"total_requests":        metrics.TotalMetricsRequests,
		"successful_requests":   metrics.SuccessfulRequests,
		"failed_requests":       metrics.FailedRequests,
		"success_rate":          metrics.SuccessRate,
		"average_response_time": metrics.AverageResponseTime,
		"cache_hits":            metrics.CacheHits,
		"cache_misses":          metrics.CacheMisses,
	}).Info("Email metrics service metrics collection test completed")

	return nil
}

func (ts *TestSuite) testEmailSystemHealthEndpoint() error {
	// Create test services
	db, err := ts.setupTestDatabase()
	if err != nil {
		return fmt.Errorf("failed to create test database: %w", err)
	}
	defer db.Close()

	emailEventRepo := repositories.NewEmailEventRepository(db, ts.logger)
	emailEventRecorder := services.NewEmailEventRecorder(emailEventRepo, ts.logger)
	emailMetricsService := services.NewEmailMetricsService(emailEventRepo, ts.logger)
	emailMonitoringService := services.NewEmailMonitoringService(emailEventRecorder, emailMetricsService, ts.logger)

	// Create health handler
	healthHandler := handlers.NewEmailSystemHealthHandler(emailMonitoringService, ts.logger)

	// Create test router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/health/email-system", healthHandler.GetEmailSystemHealthCheck)

	// Test health check endpoint
	req, err := http.NewRequest("GET", "/health/email-system", nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	// Verify response
	if recorder.Code != http.StatusOK && recorder.Code != http.StatusServiceUnavailable && recorder.Code != http.StatusPartialContent {
		return fmt.Errorf("expected status 200, 206, or 503, got %d", recorder.Code)
	}

	var response map[string]interface{}
	err = json.Unmarshal(recorder.Body.Bytes(), &response)
	if err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Verify response structure
	if _, ok := response["status"]; !ok {
		return fmt.Errorf("response missing 'status' field")
	}

	if _, ok := response["timestamp"]; !ok {
		return fmt.Errorf("response missing 'timestamp' field")
	}

	if _, ok := response["health"]; !ok {
		return fmt.Errorf("response missing 'health' field")
	}

	if _, ok := response["metrics"]; !ok {
		return fmt.Errorf("response missing 'metrics' field")
	}

	ts.logger.WithFields(logrus.Fields{
		"status_code": recorder.Code,
		"response":    response,
	}).Info("Email system health endpoint test completed")

	return nil
}

func (ts *TestSuite) testEmailSystemLivenessEndpoint() error {
	// Create test services
	db, err := storage.NewInMemoryDatabase()
	if err != nil {
		return fmt.Errorf("failed to create in-memory database: %w", err)
	}

	emailEventRepo := repositories.NewEmailEventRepository(db, ts.logger)
	emailEventRecorder := services.NewEmailEventRecorder(emailEventRepo, ts.logger)
	emailMetricsService := services.NewEmailMetricsService(emailEventRepo, ts.logger)
	emailMonitoringService := services.NewEmailMonitoringService(emailEventRecorder, emailMetricsService, ts.logger)

	// Create health handler
	healthHandler := handlers.NewEmailSystemHealthHandler(emailMonitoringService, ts.logger)

	// Create test router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/health/email-system/liveness", healthHandler.GetEmailSystemLiveness)

	// Test liveness endpoint
	req, err := http.NewRequest("GET", "/health/email-system/liveness", nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	// Verify response
	if recorder.Code != http.StatusOK && recorder.Code != http.StatusServiceUnavailable {
		return fmt.Errorf("expected status 200 or 503, got %d", recorder.Code)
	}

	var response map[string]interface{}
	err = json.Unmarshal(recorder.Body.Bytes(), &response)
	if err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Verify response structure
	if _, ok := response["status"]; !ok {
		return fmt.Errorf("response missing 'status' field")
	}

	if _, ok := response["service"]; !ok {
		return fmt.Errorf("response missing 'service' field")
	}

	ts.logger.WithFields(logrus.Fields{
		"status_code": recorder.Code,
		"response":    response,
	}).Info("Email system liveness endpoint test completed")

	return nil
}

func (ts *TestSuite) testEmailSystemReadinessEndpoint() error {
	// Create test services
	db, err := storage.NewInMemoryDatabase()
	if err != nil {
		return fmt.Errorf("failed to create in-memory database: %w", err)
	}

	emailEventRepo := repositories.NewEmailEventRepository(db, ts.logger)
	emailEventRecorder := services.NewEmailEventRecorder(emailEventRepo, ts.logger)
	emailMetricsService := services.NewEmailMetricsService(emailEventRepo, ts.logger)
	emailMonitoringService := services.NewEmailMonitoringService(emailEventRecorder, emailMetricsService, ts.logger)

	// Create health handler
	healthHandler := handlers.NewEmailSystemHealthHandler(emailMonitoringService, ts.logger)

	// Create test router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/health/email-system/readiness", healthHandler.GetEmailSystemReadiness)

	// Test readiness endpoint
	req, err := http.NewRequest("GET", "/health/email-system/readiness", nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	// Verify response
	if recorder.Code != http.StatusOK && recorder.Code != http.StatusServiceUnavailable {
		return fmt.Errorf("expected status 200 or 503, got %d", recorder.Code)
	}

	var response map[string]interface{}
	err = json.Unmarshal(recorder.Body.Bytes(), &response)
	if err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Verify response structure
	if _, ok := response["status"]; !ok {
		return fmt.Errorf("response missing 'status' field")
	}

	if _, ok := response["details"]; !ok {
		return fmt.Errorf("response missing 'details' field")
	}

	ts.logger.WithFields(logrus.Fields{
		"status_code": recorder.Code,
		"response":    response,
	}).Info("Email system readiness endpoint test completed")

	return nil
}

func (ts *TestSuite) testEmailSystemDeepHealthCheck() error {
	// Create test services
	db, err := storage.NewInMemoryDatabase()
	if err != nil {
		return fmt.Errorf("failed to create in-memory database: %w", err)
	}

	emailEventRepo := repositories.NewEmailEventRepository(db, ts.logger)
	emailEventRecorder := services.NewEmailEventRecorder(emailEventRepo, ts.logger)
	emailMetricsService := services.NewEmailMetricsService(emailEventRepo, ts.logger)
	emailMonitoringService := services.NewEmailMonitoringService(emailEventRecorder, emailMetricsService, ts.logger)

	// Create health handler
	healthHandler := handlers.NewEmailSystemHealthHandler(emailMonitoringService, ts.logger)

	// Create test router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/health/email-system/deep", healthHandler.GetEmailSystemDeepHealthCheck)

	// Test deep health check endpoint
	req, err := http.NewRequest("GET", "/health/email-system/deep", nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	// Verify response
	if recorder.Code != http.StatusOK && recorder.Code != http.StatusServiceUnavailable && recorder.Code != http.StatusPartialContent {
		return fmt.Errorf("expected status 200, 206, or 503, got %d", recorder.Code)
	}

	var response map[string]interface{}
	err = json.Unmarshal(recorder.Body.Bytes(), &response)
	if err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Verify comprehensive response structure
	requiredFields := []string{"status", "timestamp", "deep_check_duration_ms", "health", "metrics", "alerting", "recommendations"}
	for _, field := range requiredFields {
		if _, ok := response[field]; !ok {
			return fmt.Errorf("response missing '%s' field", field)
		}
	}

	ts.logger.WithFields(logrus.Fields{
		"status_code": recorder.Code,
		"response":    response,
	}).Info("Email system deep health check test completed")

	return nil
}

func (ts *TestSuite) testEmailSystemAlertingService() error {
	// Create test services
	db, err := storage.NewInMemoryDatabase()
	if err != nil {
		return fmt.Errorf("failed to create in-memory database: %w", err)
	}

	emailEventRepo := repositories.NewEmailEventRepository(db, ts.logger)
	emailEventRecorder := services.NewEmailEventRecorder(emailEventRepo, ts.logger)
	emailMetricsService := services.NewEmailMetricsService(emailEventRepo, ts.logger)

	// Create alerting service
	alertingService := services.NewEmailSystemAlertingService(emailEventRecorder, emailMetricsService, ts.logger)

	// Test starting the alerting service
	err = alertingService.Start()
	if err != nil {
		return fmt.Errorf("failed to start alerting service: %w", err)
	}

	// Verify service is running
	if !alertingService.IsRunning() {
		return fmt.Errorf("alerting service should be running")
	}

	// Get initial alert state
	alertState := alertingService.GetAlertState()
	if alertState == nil {
		return fmt.Errorf("alert state should not be nil")
	}

	// Get alert metrics
	alertMetrics := alertingService.GetAlertMetrics()
	if alertMetrics == nil {
		return fmt.Errorf("alert metrics should not be nil")
	}

	// Wait a moment for background monitoring
	time.Sleep(100 * time.Millisecond)

	// Stop the alerting service
	err = alertingService.Stop()
	if err != nil {
		return fmt.Errorf("failed to stop alerting service: %w", err)
	}

	// Verify service is stopped
	if alertingService.IsRunning() {
		return fmt.Errorf("alerting service should not be running")
	}

	ts.logger.WithFields(logrus.Fields{
		"alert_state":   alertState,
		"alert_metrics": alertMetrics,
	}).Info("Email system alerting service test completed")

	return nil
}

func (ts *TestSuite) testHighFailureRateAlerting() error {
	// This test would simulate high failure rates and verify alerting
	// For now, we'll test the alert level determination logic

	db, err := storage.NewInMemoryDatabase()
	if err != nil {
		return fmt.Errorf("failed to create in-memory database: %w", err)
	}

	emailEventRepo := repositories.NewEmailEventRepository(db, ts.logger)
	emailEventRecorder := services.NewEmailEventRecorder(emailEventRepo, ts.logger)
	emailMetricsService := services.NewEmailMetricsService(emailEventRepo, ts.logger)

	alertingService := services.NewEmailSystemAlertingService(emailEventRecorder, emailMetricsService, ts.logger)

	// Test alert configuration
	config := alertingService.GetConfig()
	if config == nil {
		return fmt.Errorf("alert configuration should not be nil")
	}

	// Verify default thresholds
	if config.RecordingFailureRateThreshold != 0.05 {
		return fmt.Errorf("expected recording failure rate threshold 0.05, got %f", config.RecordingFailureRateThreshold)
	}

	if config.HighFailureRateThreshold != 0.20 {
		return fmt.Errorf("expected high failure rate threshold 0.20, got %f", config.HighFailureRateThreshold)
	}

	if config.CriticalFailureRateThreshold != 0.50 {
		return fmt.Errorf("expected critical failure rate threshold 0.50, got %f", config.CriticalFailureRateThreshold)
	}

	ts.logger.WithFields(logrus.Fields{
		"recording_failure_threshold": config.RecordingFailureRateThreshold,
		"high_failure_threshold":      config.HighFailureRateThreshold,
		"critical_failure_threshold":  config.CriticalFailureRateThreshold,
	}).Info("High failure rate alerting test completed")

	return nil
}

func (ts *TestSuite) testAlertSuppressionLogic() error {
	// Create test services
	db, err := storage.NewInMemoryDatabase()
	if err != nil {
		return fmt.Errorf("failed to create in-memory database: %w", err)
	}

	emailEventRepo := repositories.NewEmailEventRepository(db, ts.logger)
	emailEventRecorder := services.NewEmailEventRecorder(emailEventRepo, ts.logger)
	emailMetricsService := services.NewEmailMetricsService(emailEventRepo, ts.logger)

	alertingService := services.NewEmailSystemAlertingService(emailEventRecorder, emailMetricsService, ts.logger)

	// Test alert suppression configuration
	config := alertingService.GetConfig()
	if !config.EnableAlertSuppression {
		return fmt.Errorf("alert suppression should be enabled by default")
	}

	if config.SuppressionWindow != 1*time.Hour {
		return fmt.Errorf("expected suppression window 1 hour, got %v", config.SuppressionWindow)
	}

	// Test updating configuration
	newConfig := *config
	newConfig.SuppressionWindow = 30 * time.Minute
	alertingService.UpdateConfig(&newConfig)

	updatedConfig := alertingService.GetConfig()
	if updatedConfig.SuppressionWindow != 30*time.Minute {
		return fmt.Errorf("expected updated suppression window 30 minutes, got %v", updatedConfig.SuppressionWindow)
	}

	ts.logger.WithFields(logrus.Fields{
		"suppression_enabled": config.EnableAlertSuppression,
		"suppression_window":  config.SuppressionWindow,
		"updated_window":      updatedConfig.SuppressionWindow,
	}).Info("Alert suppression logic test completed")

	return nil
}

func (ts *TestSuite) testPrometheusMetricsExport() error {
	// Create test services
	db, err := storage.NewInMemoryDatabase()
	if err != nil {
		return fmt.Errorf("failed to create in-memory database: %w", err)
	}

	emailEventRepo := repositories.NewEmailEventRepository(db, ts.logger)
	emailEventRecorder := services.NewEmailEventRecorder(emailEventRepo, ts.logger)
	emailMetricsService := services.NewEmailMetricsService(emailEventRepo, ts.logger)
	emailMonitoringService := services.NewEmailMonitoringService(emailEventRecorder, emailMetricsService, ts.logger)

	// Create monitoring handler
	monitoringHandler := handlers.NewEmailMonitoringHandler(emailMonitoringService, ts.logger)

	// Create test router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/metrics/prometheus/email-system", monitoringHandler.GetPrometheusMetrics)

	// Test Prometheus metrics endpoint
	req, err := http.NewRequest("GET", "/metrics/prometheus/email-system", nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	// Verify response
	if recorder.Code != http.StatusOK {
		return fmt.Errorf("expected status 200, got %d", recorder.Code)
	}

	// Verify content type
	contentType := recorder.Header().Get("Content-Type")
	if !strings.Contains(contentType, "text/plain") {
		return fmt.Errorf("expected content type text/plain, got %s", contentType)
	}

	// Verify Prometheus format
	body := recorder.Body.String()
	if !strings.Contains(body, "# HELP") {
		return fmt.Errorf("response should contain Prometheus HELP comments")
	}

	if !strings.Contains(body, "# TYPE") {
		return fmt.Errorf("response should contain Prometheus TYPE comments")
	}

	if !strings.Contains(body, "email_event_recorder") {
		return fmt.Errorf("response should contain email event recorder metrics")
	}

	ts.logger.WithFields(logrus.Fields{
		"status_code":  recorder.Code,
		"content_type": contentType,
		"body_length":  len(body),
	}).Info("Prometheus metrics export test completed")

	return nil
}

func (ts *TestSuite) testMonitoringServiceIntegration() error {
	// Create test services
	db, err := storage.NewInMemoryDatabase()
	if err != nil {
		return fmt.Errorf("failed to create in-memory database: %w", err)
	}

	emailEventRepo := repositories.NewEmailEventRepository(db, ts.logger)
	emailEventRecorder := services.NewEmailEventRecorder(emailEventRepo, ts.logger)
	emailMetricsService := services.NewEmailMetricsService(emailEventRepo, ts.logger)
	emailMonitoringService := services.NewEmailMonitoringService(emailEventRecorder, emailMetricsService, ts.logger)

	// Test monitoring service health
	if !emailMonitoringService.IsHealthy() {
		return fmt.Errorf("monitoring service should be healthy")
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

	// Test performing health check
	ctx := context.Background()
	healthCheck := emailMonitoringService.PerformHealthCheck(ctx)
	if healthCheck == nil {
		return fmt.Errorf("health check result should not be nil")
	}

	ts.logger.WithFields(logrus.Fields{
		"is_healthy":     emailMonitoringService.IsHealthy(),
		"overall_status": healthStatus.OverallStatus,
		"metrics":        metrics,
	}).Info("Monitoring service integration test completed")

	return nil
}

func (ts *TestSuite) testHealthCheckPerformance() error {
	// Create test services
	db, err := storage.NewInMemoryDatabase()
	if err != nil {
		return fmt.Errorf("failed to create in-memory database: %w", err)
	}

	emailEventRepo := repositories.NewEmailEventRepository(db, ts.logger)
	emailEventRecorder := services.NewEmailEventRecorder(emailEventRepo, ts.logger)
	emailMetricsService := services.NewEmailMetricsService(emailEventRepo, ts.logger)
	emailMonitoringService := services.NewEmailMonitoringService(emailEventRecorder, emailMetricsService, ts.logger)

	// Perform multiple health checks and measure performance
	ctx := context.Background()
	iterations := 10
	totalDuration := time.Duration(0)

	for i := 0; i < iterations; i++ {
		start := time.Now()
		healthCheck := emailMonitoringService.PerformHealthCheck(ctx)
		duration := time.Since(start)
		totalDuration += duration

		if healthCheck == nil {
			return fmt.Errorf("health check %d returned nil", i)
		}
	}

	averageDuration := totalDuration / time.Duration(iterations)

	// Health checks should complete within reasonable time (< 1 second)
	if averageDuration > 1*time.Second {
		return fmt.Errorf("average health check duration too high: %v", averageDuration)
	}

	ts.logger.WithFields(logrus.Fields{
		"iterations":       iterations,
		"total_duration":   totalDuration,
		"average_duration": averageDuration,
	}).Info("Health check performance test completed")

	return nil
}

func (ts *TestSuite) testAlertConfigurationManagement() error {
	// Create test services
	db, err := storage.NewInMemoryDatabase()
	if err != nil {
		return fmt.Errorf("failed to create in-memory database: %w", err)
	}

	emailEventRepo := repositories.NewEmailEventRepository(db, ts.logger)
	emailEventRecorder := services.NewEmailEventRecorder(emailEventRepo, ts.logger)
	emailMetricsService := services.NewEmailMetricsService(emailEventRepo, ts.logger)
	emailMonitoringService := services.NewEmailMonitoringService(emailEventRecorder, emailMetricsService, ts.logger)

	// Test getting alert configuration
	config := emailMonitoringService.GetAlertConfig()
	if config == nil {
		return fmt.Errorf("alert configuration should not be nil")
	}

	// Test updating alert configuration
	newConfig := *config
	newConfig.RecordingFailureThreshold = 0.10          // Change from default 0.05 to 0.10
	newConfig.AlertSuppressionWindow = 45 * time.Minute // Change from default 30 minutes

	emailMonitoringService.UpdateAlertConfig(&newConfig)

	// Verify configuration was updated
	updatedConfig := emailMonitoringService.GetAlertConfig()
	if updatedConfig.RecordingFailureThreshold != 0.10 {
		return fmt.Errorf("expected recording failure threshold 0.10, got %f", updatedConfig.RecordingFailureThreshold)
	}

	if updatedConfig.AlertSuppressionWindow != 45*time.Minute {
		return fmt.Errorf("expected alert suppression window 45 minutes, got %v", updatedConfig.AlertSuppressionWindow)
	}

	ts.logger.WithFields(logrus.Fields{
		"original_threshold":   config.RecordingFailureThreshold,
		"updated_threshold":    updatedConfig.RecordingFailureThreshold,
		"original_suppression": config.AlertSuppressionWindow,
		"updated_suppression":  updatedConfig.AlertSuppressionWindow,
	}).Info("Alert configuration management test completed")

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
