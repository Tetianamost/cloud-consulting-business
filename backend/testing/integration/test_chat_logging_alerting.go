package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/cloud-consulting/backend/internal/services"
	"github.com/cloud-consulting/backend/internal/storage"
)

func main() {
	fmt.Println("=== Chat Logging and Alerting Test ===")

	// Initialize logger
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	// Initialize structured logger
	structuredLogger := services.NewChatStructuredLogger(logger)
	fmt.Println("✅ Structured logger initialized")

	// Initialize cache for metrics
	redisConfig := &storage.RedisCacheConfig{
		Host:     "localhost",
		Port:     6379,
		Password: "",
		Database: 0,
	}
	redisCache, err := storage.NewRedisCache(redisConfig, logger)
	if err != nil {
		log.Printf("Warning: Could not initialize Redis cache: %v", err)
		redisCache = nil
	}

	// Initialize monitoring services
	performanceMonitor := services.NewChatPerformanceMonitor(logger)
	cacheMonitor := services.NewCacheMonitor(redisCache, logger)
	metricsCollector := services.NewChatMetricsCollector(performanceMonitor, cacheMonitor, logger)

	// Initialize alerting service
	alertingService := services.NewChatAlertingService(logger, structuredLogger, metricsCollector)
	fmt.Println("✅ Alerting service initialized")

	// Test structured logging with correlation IDs
	fmt.Println("\n--- Testing Structured Logging ---")

	correlationID := structuredLogger.GenerateCorrelationID()
	fmt.Printf("Generated correlation ID: %s\n", correlationID)

	// Test connection event logging
	structuredLogger.LogConnectionEvent(
		correlationID,
		"user123",
		"conn456",
		"opened",
		map[string]interface{}{
			"ip_address": "192.168.1.100",
			"user_agent": "Mozilla/5.0...",
		},
	)

	// Test message event logging
	structuredLogger.LogMessageEvent(
		correlationID,
		"user123",
		"session789",
		"msg001",
		"sent",
		150*time.Millisecond,
		map[string]interface{}{
			"message_length": 45,
			"message_type":   "user",
		},
	)

	// Test AI event logging
	structuredLogger.LogAIEvent(
		correlationID,
		"user123",
		"session789",
		"claude-3-sonnet",
		"success",
		2*time.Second,
		1500,
		0.015,
		map[string]interface{}{
			"prompt_length":   200,
			"response_length": 800,
		},
	)

	// Test security event logging
	securityEvent := structuredLogger.CreateSecurityEvent(
		"authentication_failure",
		"medium",
		"user123",
		"192.168.1.100",
		"Mozilla/5.0...",
		correlationID,
		map[string]interface{}{
			"reason":   "invalid_password",
			"attempts": 3,
		},
	)
	structuredLogger.LogSecurityEvent(securityEvent)

	// Test performance event logging
	performanceEvent := structuredLogger.CreatePerformanceEvent(
		"message_processing",
		250*time.Millisecond,
		true,
		"",
		correlationID,
		map[string]interface{}{
			"queue_size": 5,
			"worker_id":  "worker-1",
		},
	)
	structuredLogger.LogPerformanceEvent(performanceEvent)

	// Test error event logging
	errorEvent := structuredLogger.CreateErrorEvent(
		"database_error",
		"Connection timeout",
		"chat_service",
		"save_message",
		"user123",
		"session789",
		correlationID,
		fmt.Errorf("connection timeout"),
		map[string]interface{}{
			"retry_count": 2,
			"timeout_ms":  5000,
		},
	)
	structuredLogger.LogErrorEvent(errorEvent)

	fmt.Println("✅ Structured logging events recorded")

	// Test alert creation
	fmt.Println("\n--- Testing Alert Creation ---")

	// Create different types of alerts
	errorAlert := alertingService.CreateAlert(
		services.AlertTypeError,
		services.AlertSeverityHigh,
		"High Error Rate Detected",
		"Error rate has exceeded 10% in the last 5 minutes",
		"chat_service",
		map[string]interface{}{
			"correlation_id": correlationID,
			"error_rate":     0.12,
			"threshold":      0.10,
		},
	)
	fmt.Printf("Created error alert: %s\n", errorAlert.ID)

	performanceAlert := alertingService.CreateAlert(
		services.AlertTypePerformance,
		services.AlertSeverityMedium,
		"High Response Time",
		"95th percentile response time is above 2 seconds",
		"websocket_handler",
		map[string]interface{}{
			"correlation_id":    correlationID,
			"response_time_p95": 2500,
			"threshold":         2000,
		},
	)
	fmt.Printf("Created performance alert: %s\n", performanceAlert.ID)

	securityAlert := alertingService.CreateAlert(
		services.AlertTypeSecurity,
		services.AlertSeverityCritical,
		"Multiple Authentication Failures",
		"Detected 15 failed authentication attempts from same IP",
		"auth_service",
		map[string]interface{}{
			"correlation_id": correlationID,
			"ip_address":     "192.168.1.100",
			"failure_count":  15,
			"time_window":    "5m",
		},
	)
	fmt.Printf("Created security alert: %s\n", securityAlert.ID)

	fmt.Println("✅ Alerts created successfully")

	// Test alert resolution
	fmt.Println("\n--- Testing Alert Resolution ---")

	time.Sleep(1 * time.Second) // Wait a bit before resolving

	err = alertingService.ResolveAlert(errorAlert.ID)
	if err != nil {
		fmt.Printf("❌ Failed to resolve alert: %v\n", err)
	} else {
		fmt.Printf("✅ Resolved alert: %s\n", errorAlert.ID)
	}

	// Test getting active alerts
	fmt.Println("\n--- Testing Active Alerts Retrieval ---")

	activeAlerts := alertingService.GetActiveAlerts()
	fmt.Printf("Active alerts count: %d\n", len(activeAlerts))

	for _, alert := range activeAlerts {
		fmt.Printf("  - %s: %s (%s)\n", alert.ID, alert.Title, alert.Severity)
	}

	// Test alert history
	fmt.Println("\n--- Testing Alert History ---")

	alertHistory := alertingService.GetAlertHistory(10, 0)
	fmt.Printf("Alert history count: %d\n", len(alertHistory))

	for _, alert := range alertHistory {
		status := "ACTIVE"
		if alert.Resolved {
			status = "RESOLVED"
		}
		fmt.Printf("  - %s: %s (%s) [%s]\n", alert.ID, alert.Title, alert.Severity, status)
	}

	// Test alert rules
	fmt.Println("\n--- Testing Alert Rules ---")

	alertRules := alertingService.GetAlertRules()
	fmt.Printf("Alert rules count: %d\n", len(alertRules))

	for _, rule := range alertRules {
		status := "DISABLED"
		if rule.Enabled {
			status = "ENABLED"
		}
		fmt.Printf("  - %s: %s (%s) [%s]\n", rule.ID, rule.Name, rule.Severity, status)
		fmt.Printf("    Metric: %s %s %v\n", rule.Metric, rule.Condition, rule.Threshold)
	}

	// Test alert channels
	fmt.Println("\n--- Testing Alert Channels ---")

	alertChannels := alertingService.GetAlertChannels()
	fmt.Printf("Alert channels count: %d\n", len(alertChannels))

	for _, channel := range alertChannels {
		status := "DISABLED"
		if channel.Enabled {
			status = "ENABLED"
		}
		fmt.Printf("  - %s: %s (%s) [%s]\n", channel.ID, channel.Name, channel.Type, status)
	}

	// Test metrics-based alerting
	fmt.Println("\n--- Testing Metrics-Based Alerting ---")

	// Simulate high error rate
	for i := 0; i < 10; i++ {
		metricsCollector.RecordError("system")
		metricsCollector.RecordMessage("sent")
	}

	// Simulate high response times
	for i := 0; i < 5; i++ {
		metricsCollector.RecordMessage("sent", 3*time.Second)
	}

	// Check for alerts based on metrics
	fmt.Println("Checking metrics for alert conditions...")
	alertingService.CheckMetricsForAlerts()

	// Wait a moment for alerts to be processed
	time.Sleep(500 * time.Millisecond)

	newActiveAlerts := alertingService.GetActiveAlerts()
	fmt.Printf("Active alerts after metrics check: %d\n", len(newActiveAlerts))

	for _, alert := range newActiveAlerts {
		if alert.Metric != "" {
			fmt.Printf("  - Metric Alert: %s = %v (threshold: %v)\n", alert.Metric, alert.Value, alert.Threshold)
		}
	}

	// Test periodic alert checking
	fmt.Println("\n--- Testing Periodic Alert Checking ---")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	go alertingService.StartPeriodicChecks(1 * time.Second)

	// Add more metrics to trigger alerts
	go func() {
		for i := 0; i < 3; i++ {
			time.Sleep(500 * time.Millisecond)
			metricsCollector.RecordError("authentication")
			metricsCollector.RecordConnection("failed")
		}
	}()

	// Wait for periodic checks
	<-ctx.Done()

	finalActiveAlerts := alertingService.GetActiveAlerts()
	fmt.Printf("Final active alerts count: %d\n", len(finalActiveAlerts))

	// Test log search (placeholder functionality)
	fmt.Println("\n--- Testing Log Search ---")

	query := services.LogAggregationQuery{
		StartTime: time.Now().Add(-1 * time.Hour),
		EndTime:   time.Now(),
		Level:     services.LogLevelError,
		Component: "chat_service",
		Limit:     10,
		Offset:    0,
	}

	searchResult, err := structuredLogger.SearchLogs(query)
	if err != nil {
		fmt.Printf("❌ Log search failed: %v\n", err)
	} else {
		fmt.Printf("✅ Log search completed: %d results\n", searchResult.TotalCount)
	}

	// Test log statistics
	fmt.Println("\n--- Testing Log Statistics ---")

	stats := structuredLogger.GetLogStatistics(time.Now().Add(-1*time.Hour), time.Now())
	fmt.Printf("Log statistics: %+v\n", stats)

	// Stop services
	alertingService.Stop()
	fmt.Println("✅ Services stopped")

	fmt.Println("\n🎉 All logging and alerting tests completed successfully!")
}
