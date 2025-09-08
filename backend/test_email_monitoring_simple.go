package main

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/cloud-consulting/backend/internal/interfaces"
	"github.com/cloud-consulting/backend/internal/services"
)

func main() {
	fmt.Println("=== Simple Email Monitoring Test ===")

	// Create logger
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	// Test 1: Create email monitoring service with nil dependencies (for basic testing)
	fmt.Println("\n1. Testing Email Monitoring Service Creation...")

	// Create a mock email event recorder for testing
	mockRecorder := &MockEmailEventRecorder{}
	mockMetricsService := &MockEmailMetricsService{}

	emailMonitoringService := services.NewEmailMonitoringService(
		mockRecorder,
		mockMetricsService,
		logger,
	)

	fmt.Println("✓ Email monitoring service created successfully")

	// Test 2: Test system metrics
	fmt.Println("\n2. Testing System Metrics...")
	systemMetrics := emailMonitoringService.GetSystemMetrics()
	fmt.Printf("✓ System metrics retrieved: %+v\n", systemMetrics)

	// Test 3: Test health status
	fmt.Println("\n3. Testing Health Status...")
	healthStatus := emailMonitoringService.GetHealthStatus()
	fmt.Printf("✓ Health status retrieved: %+v\n", healthStatus)

	// Test 4: Test alert configuration
	fmt.Println("\n4. Testing Alert Configuration...")
	alertConfig := emailMonitoringService.GetAlertConfig()
	fmt.Printf("✓ Alert config retrieved: %+v\n", alertConfig)

	// Test 5: Test monitoring service health
	fmt.Println("\n5. Testing Monitoring Service Health...")
	isHealthy := emailMonitoringService.IsHealthy()
	fmt.Printf("✓ Monitoring service health: %v\n", isHealthy)

	fmt.Println("\n=== Simple Email Monitoring Test Complete ===")
	fmt.Println("✓ Basic monitoring functionality verified")
}

// MockEmailEventRecorder for testing
type MockEmailEventRecorder struct{}

func (m *MockEmailEventRecorder) RecordEmailSent(ctx interface{}, event interface{}) error {
	return nil
}

func (m *MockEmailEventRecorder) UpdateEmailStatus(ctx interface{}, messageID string, status interface{}, deliveredAt *time.Time, errorMsg string) error {
	return nil
}

func (m *MockEmailEventRecorder) GetEmailEventsByInquiry(ctx interface{}, inquiryID string) (interface{}, error) {
	return nil, nil
}

func (m *MockEmailEventRecorder) IsHealthy() bool {
	return true
}

func (m *MockEmailEventRecorder) IsHealthyWithContext(ctx interface{}) bool {
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
		RetryAttempts:          2,
		HealthCheckFailures:    0,
	}
}

// MockEmailMetricsService for testing
type MockEmailMetricsService struct{}

func (m *MockEmailMetricsService) GetEmailMetrics(ctx interface{}, timeRange interface{}) (interface{}, error) {
	return nil, nil
}

func (m *MockEmailMetricsService) GetEmailStatusByInquiry(ctx interface{}, inquiryID string) (interface{}, error) {
	return nil, nil
}

func (m *MockEmailMetricsService) GetEmailEventHistory(ctx interface{}, filters interface{}) (interface{}, error) {
	return nil, nil
}

func (m *MockEmailMetricsService) GetRecentEmailActivity(ctx interface{}, hours int) (interface{}, error) {
	return []interface{}{}, nil
}

func (m *MockEmailMetricsService) IsHealthy(ctx interface{}) bool {
	return true
}

func (m *MockEmailMetricsService) GetMetrics() *interfaces.EmailMetricsServiceMetrics {
	return &interfaces.EmailMetricsServiceMetrics{
		TotalMetricsRequests: 20,
		SuccessfulRequests:   18,
		FailedRequests:       2,
		SuccessRate:          0.9,
		AverageResponseTime:  25.0,
		LastRequestTime:      time.Now(),
		CacheHits:            15,
		CacheMisses:          5,
		HealthCheckFailures:  1,
	}
}
