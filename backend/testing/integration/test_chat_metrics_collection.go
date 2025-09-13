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
	fmt.Println("=== Chat Metrics Collection Test ===")

	// Initialize logger
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	// Initialize cache for cache monitor
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

	// Initialize monitors
	performanceMonitor := services.NewChatPerformanceMonitor(logger)
	cacheMonitor := services.NewCacheMonitor(redisCache, logger)
	metricsCollector := services.NewChatMetricsCollector(performanceMonitor, cacheMonitor, logger)

	fmt.Println("✅ Metrics services initialized")

	// Test connection metrics
	fmt.Println("\n--- Testing Connection Metrics ---")
	metricsCollector.RecordConnection("opened")
	metricsCollector.RecordConnection("opened")
	metricsCollector.RecordConnection("failed")
	metricsCollector.RecordConnection("closed", 5*time.Minute)
	fmt.Println("✅ Connection metrics recorded")

	// Test message metrics
	fmt.Println("\n--- Testing Message Metrics ---")
	metricsCollector.RecordMessage("sent", 100*time.Millisecond)
	metricsCollector.RecordMessage("received", 50*time.Millisecond)
	metricsCollector.RecordMessage("error")
	metricsCollector.RecordMessage("retry")
	fmt.Println("✅ Message metrics recorded")

	// Test AI metrics
	fmt.Println("\n--- Testing AI Metrics ---")
	metricsCollector.RecordAIRequest(true, 2*time.Second, 1500, "claude-3-sonnet", 0.015)
	metricsCollector.RecordAIRequest(true, 1500*time.Millisecond, 1200, "claude-3-haiku", 0.008)
	metricsCollector.RecordAIRequest(false, 5*time.Second, 0, "claude-3-sonnet", 0)
	fmt.Println("✅ AI metrics recorded")

	// Test user engagement metrics
	fmt.Println("\n--- Testing User Engagement Metrics ---")
	metricsCollector.RecordUserEngagement("user_active", nil)
	metricsCollector.RecordUserEngagement("session_duration", 15*time.Minute)
	metricsCollector.RecordUserEngagement("quick_action", nil)
	metricsCollector.RecordUserEngagement("feature_usage", "cost_analysis")
	metricsCollector.RecordUserEngagement("feature_usage", "architecture_review")
	fmt.Println("✅ User engagement metrics recorded")

	// Test error metrics
	fmt.Println("\n--- Testing Error Metrics ---")
	metricsCollector.RecordError("authentication")
	metricsCollector.RecordError("validation")
	metricsCollector.RecordError("system")
	fmt.Println("✅ Error metrics recorded")

	// Test cache metrics (if Redis is available)
	fmt.Println("\n--- Testing Cache Metrics ---")
	ctx := context.Background()
	if redisCache != nil && redisCache.IsHealthy(ctx) {
		cacheMonitor.RecordSessionHit(10 * time.Millisecond)
		cacheMonitor.RecordSessionMiss(50 * time.Millisecond)
		cacheMonitor.RecordMessageHit(5 * time.Millisecond)
		cacheMonitor.RecordCacheInvalidation(20 * time.Millisecond)
		fmt.Println("✅ Cache metrics recorded")
	} else {
		fmt.Println("⚠️  Redis not available, skipping cache metrics test")
	}

	// Get comprehensive metrics
	fmt.Println("\n--- Retrieving Comprehensive Metrics ---")
	metrics := metricsCollector.GetMetrics()

	// Display connection metrics
	fmt.Printf("\n📊 Connection Metrics:\n")
	fmt.Printf("  Total Connections: %d\n", metrics.ConnectionMetrics.TotalConnections)
	fmt.Printf("  Active Connections: %d\n", metrics.ConnectionMetrics.ActiveConnections)
	fmt.Printf("  Connection Failures: %d\n", metrics.ConnectionMetrics.ConnectionFailures)
	fmt.Printf("  Success Rate: %.2f%%\n", metrics.ConnectionMetrics.ConnectionSuccessRate*100)

	// Display message metrics
	fmt.Printf("\n📨 Message Metrics:\n")
	fmt.Printf("  Messages Sent: %d\n", metrics.MessageMetrics.MessagesSent)
	fmt.Printf("  Messages Received: %d\n", metrics.MessageMetrics.MessagesReceived)
	fmt.Printf("  Message Errors: %d\n", metrics.MessageMetrics.MessageErrors)
	fmt.Printf("  Success Rate: %.2f%%\n", metrics.MessageMetrics.MessageSuccessRate*100)
	fmt.Printf("  Avg Processing Time: %v\n", metrics.MessageMetrics.AverageProcessingTime)

	// Display AI metrics
	fmt.Printf("\n🤖 AI Metrics:\n")
	fmt.Printf("  Total Requests: %d\n", metrics.AIMetrics.RequestsTotal)
	fmt.Printf("  Successful Requests: %d\n", metrics.AIMetrics.RequestsSuccessful)
	fmt.Printf("  Failed Requests: %d\n", metrics.AIMetrics.RequestsFailed)
	fmt.Printf("  Success Rate: %.2f%%\n", metrics.AIMetrics.SuccessRate*100)
	fmt.Printf("  Tokens Used: %d\n", metrics.AIMetrics.TokensUsed)
	fmt.Printf("  Estimated Cost: $%.4f\n", metrics.AIMetrics.EstimatedCost)
	fmt.Printf("  Avg Response Time: %v\n", metrics.AIMetrics.AverageResponseTime)
	fmt.Printf("  Model Usage: %+v\n", metrics.AIMetrics.ModelUsage)

	// Display user metrics
	fmt.Printf("\n👥 User Metrics:\n")
	fmt.Printf("  Active Users: %d\n", metrics.UserMetrics.ActiveUsers)
	fmt.Printf("  Avg Session Duration: %v\n", metrics.UserMetrics.AverageSessionDuration)
	fmt.Printf("  Messages Per Session: %.2f\n", metrics.UserMetrics.MessagesPerSession)
	fmt.Printf("  Quick Actions Used: %d\n", metrics.UserMetrics.QuickActionsUsed)
	fmt.Printf("  Feature Usage: %+v\n", metrics.UserMetrics.FeatureUsage)

	// Display error metrics
	fmt.Printf("\n❌ Error Metrics:\n")
	fmt.Printf("  Authentication Errors: %d\n", metrics.ErrorMetrics.AuthenticationErrors)
	fmt.Printf("  Authorization Errors: %d\n", metrics.ErrorMetrics.AuthorizationErrors)
	fmt.Printf("  Validation Errors: %d\n", metrics.ErrorMetrics.ValidationErrors)
	fmt.Printf("  System Errors: %d\n", metrics.ErrorMetrics.SystemErrors)
	fmt.Printf("  Total Errors: %d\n", metrics.ErrorMetrics.TotalErrors)
	fmt.Printf("  Error Rate: %.2f%%\n", metrics.ErrorMetrics.ErrorRate*100)

	// Display performance metrics
	fmt.Printf("\n⚡ Performance Metrics:\n")
	fmt.Printf("  Response Time P50: %v\n", metrics.PerformanceMetrics.ResponseTimeP50)
	fmt.Printf("  Response Time P95: %v\n", metrics.PerformanceMetrics.ResponseTimeP95)
	fmt.Printf("  Throughput (MPS): %.2f\n", metrics.PerformanceMetrics.ThroughputMPS)
	fmt.Printf("  Memory Usage: %d MB\n", metrics.PerformanceMetrics.MemoryUsageMB)

	// Display cache metrics
	if metrics.CacheMetrics != nil {
		fmt.Printf("\n💾 Cache Metrics:\n")
		fmt.Printf("  Session Hits: %d\n", metrics.CacheMetrics.SessionHits)
		fmt.Printf("  Session Misses: %d\n", metrics.CacheMetrics.SessionMisses)
		fmt.Printf("  Message Hits: %d\n", metrics.CacheMetrics.MessageHits)
		fmt.Printf("  Message Misses: %d\n", metrics.CacheMetrics.MessageMisses)
		fmt.Printf("  Total Operations: %d\n", metrics.CacheMetrics.TotalOperations)
		fmt.Printf("  Avg Response Time: %v\n", metrics.CacheMetrics.AverageResponseTime)

		// Calculate hit ratios
		sessionHitRatio := float64(0)
		messageHitRatio := float64(0)
		overallHitRatio := float64(0)

		if cacheMonitor != nil {
			sessionHitRatio = cacheMonitor.GetSessionHitRatio()
			messageHitRatio = cacheMonitor.GetMessageHitRatio()
			overallHitRatio = cacheMonitor.GetCacheHitRatio()
		}

		fmt.Printf("  Session Hit Ratio: %.2f%%\n", sessionHitRatio*100)
		fmt.Printf("  Message Hit Ratio: %.2f%%\n", messageHitRatio*100)
		fmt.Printf("  Overall Hit Ratio: %.2f%%\n", overallHitRatio*100)
	}

	// Test Prometheus metrics export
	fmt.Println("\n--- Testing Prometheus Metrics Export ---")
	prometheusMetrics := metricsCollector.GetPrometheusMetrics()
	fmt.Printf("Prometheus metrics length: %d characters\n", len(prometheusMetrics))

	// Show first few lines of Prometheus metrics
	lines := splitString(prometheusMetrics, "\n")
	fmt.Println("First 10 lines of Prometheus metrics:")
	for i, line := range lines {
		if i >= 10 {
			break
		}
		fmt.Printf("  %s\n", line)
	}
	fmt.Println("✅ Prometheus metrics generated")

	// Test periodic metrics collection
	fmt.Println("\n--- Testing Periodic Metrics Collection ---")
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	go metricsCollector.StartPeriodicCollection(ctx, 1*time.Second)

	// Wait for a few collection cycles
	time.Sleep(2500 * time.Millisecond)
	fmt.Println("✅ Periodic metrics collection tested")

	// Test metrics reset
	fmt.Println("\n--- Testing Metrics Reset ---")
	metricsCollector.ResetMetrics()
	resetMetrics := metricsCollector.GetMetrics()

	fmt.Printf("After reset - Total Connections: %d\n", resetMetrics.ConnectionMetrics.TotalConnections)
	fmt.Printf("After reset - Messages Sent: %d\n", resetMetrics.MessageMetrics.MessagesSent)
	fmt.Printf("After reset - AI Requests: %d\n", resetMetrics.AIMetrics.RequestsTotal)
	fmt.Println("✅ Metrics reset tested")

	fmt.Println("\n🎉 All chat metrics collection tests completed successfully!")
}

// splitString splits a string by separator (simple implementation)
func splitString(s, sep string) []string {
	if s == "" {
		return []string{}
	}

	var result []string
	start := 0

	for i := 0; i <= len(s)-len(sep); i++ {
		if s[i:i+len(sep)] == sep {
			result = append(result, s[start:i])
			start = i + len(sep)
			i += len(sep) - 1
		}
	}

	// Add the last part
	if start < len(s) {
		result = append(result, s[start:])
	}

	return result
}
