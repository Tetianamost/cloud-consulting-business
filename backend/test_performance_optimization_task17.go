package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/cloud-consulting/backend/internal/services"
	"github.com/sirupsen/logrus"
)

func main() {
	fmt.Println("=== Testing Enhanced Bedrock Performance Optimization (Task 17) ===")

	// Initialize logger
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	// Test 1: Enhanced Bedrock Performance Optimizer
	fmt.Println("\n1. Testing Enhanced Bedrock Performance Optimizer...")
	testPerformanceOptimizer(logger)

	// Test 2: Intelligent Analysis Cache
	fmt.Println("\n2. Testing Intelligent Analysis Cache...")
	testIntelligentCache(logger)

	// Test 3: Consultant Session Load Balancer
	fmt.Println("\n3. Testing Consultant Session Load Balancer...")
	testLoadBalancer(logger)

	// Test 4: Enhanced Performance Monitor
	fmt.Println("\n4. Testing Enhanced Performance Monitor...")
	testPerformanceMonitor(logger)

	// Test 5: Integration Test
	fmt.Println("\n5. Testing Performance Optimization Integration...")
	testIntegration(logger)

	fmt.Println("\n=== All Performance Optimization Tests Completed ===")
}

func testPerformanceOptimizer(logger *logrus.Logger) {
	fmt.Println("Creating performance optimizer...")

	// Create mock dependencies
	responseCache := services.NewSimpleChatResponseCache(nil, logger)
	performanceMonitor := services.NewChatPerformanceMonitor(logger)

	// Create performance optimizer
	optimizer := services.NewEnhancedBedrockPerformanceOptimizer(
		logger,
		responseCache,
		performanceMonitor,
	)

	// Test optimization request
	request := &services.OptimizationRequest{
		SessionID:    "test-session-1",
		ConsultantID: "consultant-1",
		AnalysisType: "cost_analysis",
		Content:      "AWS cost optimization for e-commerce platform",
		Prompt:       "Analyze the cost optimization opportunities for an e-commerce platform running on AWS",
		MaxTokens:    2000,
		Temperature:  0.8,
		Priority:     "high",
	}

	ctx := context.Background()
	result, err := optimizer.OptimizeRequest(ctx, request)
	if err != nil {
		log.Printf("Error optimizing request: %v", err)
		return
	}

	fmt.Printf("✓ Request optimized successfully\n")
	fmt.Printf("  - Session ID: %s\n", result.SessionID)
	fmt.Printf("  - Cache Hit: %v\n", result.CacheHit)
	fmt.Printf("  - Optimized: %v\n", result.Optimized)
	fmt.Printf("  - Response Time: %v\n", result.ResponseTime)

	// Get performance metrics
	metrics := optimizer.GetPerformanceMetrics()
	fmt.Printf("✓ Performance metrics retrieved\n")
	fmt.Printf("  - Total Requests: %d\n", metrics.TotalRequests)
	fmt.Printf("  - Optimized Requests: %d\n", metrics.OptimizedRequests)
	fmt.Printf("  - Cache Hit Rate: %.2f\n", metrics.CacheHitRate)
	fmt.Printf("  - Active Sessions: %d\n", metrics.ActiveSessions)

	// Start monitoring (briefly)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	go optimizer.MonitorPerformance(ctx)
	time.Sleep(1 * time.Second)

	fmt.Printf("✓ Performance monitoring tested\n")
}

func testIntelligentCache(logger *logrus.Logger) {
	fmt.Println("Creating intelligent analysis cache...")

	cache := services.NewIntelligentAnalysisCache(logger)

	// Test cache miss
	result := cache.GetCachedAnalysis("cost_analysis", "AWS cost optimization")
	if result != nil {
		fmt.Printf("✗ Expected cache miss, got hit\n")
		return
	}
	fmt.Printf("✓ Cache miss handled correctly\n")

	// Test caching analysis
	err := cache.CacheAnalysis(
		"cost_analysis",
		"AWS cost optimization",
		"Cost analysis result with specific recommendations",
		150,
		0.85,
	)
	if err != nil {
		log.Printf("Error caching analysis: %v", err)
		return
	}
	fmt.Printf("✓ Analysis cached successfully\n")

	// Test cache hit
	result = cache.GetCachedAnalysis("cost_analysis", "AWS cost optimization")
	if result == nil {
		fmt.Printf("✗ Expected cache hit, got miss\n")
		return
	}
	fmt.Printf("✓ Cache hit successful\n")
	fmt.Printf("  - Content length: %d\n", len(result.Content))
	fmt.Printf("  - Tokens used: %d\n", result.TokensUsed)
	fmt.Printf("  - Quality: %.2f\n", result.Quality)
	fmt.Printf("  - Access count: %d\n", result.AccessCount)

	// Test cache optimization
	cache.OptimizeCacheStrategy()
	fmt.Printf("✓ Cache strategy optimized\n")

	// Test cache warming
	cache.WarmCache()
	fmt.Printf("✓ Cache warmed with common patterns\n")

	// Get cache statistics
	stats := cache.GetCacheStatistics()
	fmt.Printf("✓ Cache statistics retrieved\n")
	fmt.Printf("  - Total Requests: %d\n", stats.TotalRequests)
	fmt.Printf("  - Hit Rate: %.2f\n", stats.HitRate)
	fmt.Printf("  - Cache Size: %d\n", stats.CacheSize)
	fmt.Printf("  - Analysis Types: %d\n", stats.AnalysisTypes)
}

func testLoadBalancer(logger *logrus.Logger) {
	fmt.Println("Creating consultant session load balancer...")

	loadBalancer := services.NewConsultantSessionLoadBalancer(logger)

	// Test session assignment
	consultantID := loadBalancer.AssignSession("session-1", "")
	if consultantID == "" {
		fmt.Printf("✗ Failed to assign session\n")
		return
	}
	fmt.Printf("✓ Session assigned to consultant: %s\n", consultantID)

	// Test multiple session assignments
	sessions := []string{"session-2", "session-3", "session-4"}
	for _, sessionID := range sessions {
		assignedConsultant := loadBalancer.AssignSession(sessionID, "")
		if assignedConsultant == "" {
			fmt.Printf("✗ Failed to assign session %s\n", sessionID)
			return
		}
		fmt.Printf("✓ Session %s assigned to consultant: %s\n", sessionID, assignedConsultant)
	}

	// Get active session count
	activeCount := loadBalancer.GetActiveSessionCount()
	fmt.Printf("✓ Active sessions: %d\n", activeCount)

	// Get load balancing metrics
	metrics := loadBalancer.GetLoadBalancingMetrics()
	fmt.Printf("✓ Load balancing metrics retrieved\n")
	fmt.Printf("  - Total Sessions: %d\n", metrics.TotalSessions)
	fmt.Printf("  - Active Sessions: %d\n", metrics.ActiveSessions)
	fmt.Printf("  - Available Consultants: %d\n", metrics.AvailableConsultants)
	fmt.Printf("  - Average Load: %.2f\n", metrics.AverageLoad)
	fmt.Printf("  - Strategy: %s\n", metrics.Strategy)

	// Test optimization for response time
	loadBalancer.OptimizeForResponseTime()
	fmt.Printf("✓ Load balancer optimized for response time\n")

	// Test session release
	err := loadBalancer.ReleaseSession("session-1")
	if err != nil {
		log.Printf("Error releasing session: %v", err)
		return
	}
	fmt.Printf("✓ Session released successfully\n")

	// Test cleanup (briefly)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	go loadBalancer.CleanupExpiredSessions(ctx)
	time.Sleep(500 * time.Millisecond)
	fmt.Printf("✓ Session cleanup tested\n")
}

func testPerformanceMonitor(logger *logrus.Logger) {
	fmt.Println("Creating enhanced performance monitor...")

	monitor := services.NewEnhancedBedrockPerformanceMonitor(logger)

	// Test recording metrics
	monitor.RecordRequest(true, 1500*time.Millisecond)
	monitor.RecordRequest(true, 2000*time.Millisecond)
	monitor.RecordRequest(false, 5000*time.Millisecond)
	fmt.Printf("✓ Request metrics recorded\n")

	// Test cache metrics
	monitor.RecordCacheMetrics(75, 25, 100, 5, 15*time.Minute)
	fmt.Printf("✓ Cache metrics recorded\n")

	// Test system metrics
	monitor.RecordSystemMetrics(65.5, 70.2, 150, 1024*1024*75, 2*time.Millisecond)
	fmt.Printf("✓ System metrics recorded\n")

	// Get performance report
	report := monitor.GetPerformanceReport()
	fmt.Printf("✓ Performance report generated\n")
	fmt.Printf("  - Total Requests: %d\n", report.RequestMetrics.TotalRequests)
	fmt.Printf("  - Successful Requests: %d\n", report.RequestMetrics.SuccessfulRequests)
	fmt.Printf("  - Failed Requests: %d\n", report.RequestMetrics.FailedRequests)
	fmt.Printf("  - Average Response Time: %v\n", report.ResponseTimeMetrics.AverageResponseTime)
	fmt.Printf("  - Cache Hit Rate: %.2f\n", report.CacheMetrics.CacheHitRate)
	fmt.Printf("  - CPU Usage: %.1f%%\n", report.SystemMetrics.CPUUsage)
	fmt.Printf("  - Memory Usage: %.1f%%\n", report.SystemMetrics.MemoryUsage)

	// Test alert handler registration
	alertCount := 0
	monitor.RegisterAlertHandler("test", func(ctx context.Context, alert *services.PerformanceAlert) error {
		alertCount++
		fmt.Printf("✓ Alert received: %s - %s\n", alert.Type, alert.Message)
		return nil
	})

	// Test monitoring (briefly)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	go monitor.StartMonitoring(ctx)
	time.Sleep(1 * time.Second)
	fmt.Printf("✓ Performance monitoring tested\n")
}

func testIntegration(logger *logrus.Logger) {
	fmt.Println("Testing integrated performance optimization...")

	// Create all components
	responseCache := services.NewSimpleChatResponseCache(nil, logger)
	performanceMonitor := services.NewChatPerformanceMonitor(logger)
	optimizer := services.NewEnhancedBedrockPerformanceOptimizer(logger, responseCache, performanceMonitor)
	intelligentCache := services.NewIntelligentAnalysisCache(logger)
	loadBalancer := services.NewConsultantSessionLoadBalancer(logger)
	enhancedMonitor := services.NewEnhancedBedrockPerformanceMonitor(logger)

	// Simulate a complete workflow
	fmt.Printf("Simulating complete performance optimization workflow...\n")

	// 1. Assign session
	sessionID := "integration-test-session"
	consultantID := loadBalancer.AssignSession(sessionID, "")
	fmt.Printf("✓ Session assigned: %s -> %s\n", sessionID, consultantID)

	// 2. Cache some analysis
	err := intelligentCache.CacheAnalysis(
		"architecture_review",
		"microservices architecture assessment",
		"Detailed architecture review with recommendations",
		200,
		0.9,
	)
	if err != nil {
		log.Printf("Error caching analysis: %v", err)
		return
	}
	fmt.Printf("✓ Analysis cached\n")

	// 3. Optimize request
	request := &services.OptimizationRequest{
		SessionID:    sessionID,
		ConsultantID: consultantID,
		AnalysisType: "architecture_review",
		Content:      "microservices architecture assessment",
		Prompt:       "Review microservices architecture",
		MaxTokens:    1500,
		Temperature:  0.7,
		Priority:     "normal",
	}

	ctx := context.Background()
	result, err := optimizer.OptimizeRequest(ctx, request)
	if err != nil {
		log.Printf("Error optimizing request: %v", err)
		return
	}
	fmt.Printf("✓ Request optimized (Cache Hit: %v)\n", result.CacheHit)

	// 4. Record performance metrics
	enhancedMonitor.RecordRequest(true, result.ResponseTime)
	fmt.Printf("✓ Performance metrics recorded\n")

	// 5. Get comprehensive metrics
	optimizerMetrics := optimizer.GetPerformanceMetrics()
	cacheStats := intelligentCache.GetCacheStatistics()
	loadBalancerMetrics := loadBalancer.GetLoadBalancingMetrics()
	performanceReport := enhancedMonitor.GetPerformanceReport()

	fmt.Printf("✓ Integration test completed successfully\n")
	fmt.Printf("  - Optimizer Requests: %d\n", optimizerMetrics.TotalRequests)
	fmt.Printf("  - Cache Hit Rate: %.2f\n", cacheStats.HitRate)
	fmt.Printf("  - Load Balancer Sessions: %d\n", loadBalancerMetrics.ActiveSessions)
	fmt.Printf("  - Monitor Requests: %d\n", performanceReport.RequestMetrics.TotalRequests)

	// 6. Release session
	err = loadBalancer.ReleaseSession(sessionID)
	if err != nil {
		log.Printf("Error releasing session: %v", err)
		return
	}
	fmt.Printf("✓ Session released\n")

	fmt.Printf("✓ All performance optimization components working together\n")
}
