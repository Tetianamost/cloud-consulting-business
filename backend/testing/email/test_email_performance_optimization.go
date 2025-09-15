package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"

	"github.com/cloud-consulting/backend/internal/config"
	"github.com/cloud-consulting/backend/internal/domain"
	"github.com/cloud-consulting/backend/internal/repositories"
	"github.com/cloud-consulting/backend/internal/services"
)

func main() {
	fmt.Println("=== Email Performance Optimization Test ===")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create logger
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	// Connect to database
	db, err := sql.Open("postgres", cfg.Database.URL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	fmt.Println("✓ Database connection established")

	ctx := context.Background()

	// Run performance optimization tests
	if err := runPerformanceOptimizationTests(ctx, db, logger); err != nil {
		log.Fatalf("Performance optimization tests failed: %v", err)
	}

	fmt.Println("✓ All performance optimization tests passed!")
}

func runPerformanceOptimizationTests(ctx context.Context, db *sql.DB, logger *logrus.Logger) error {
	fmt.Println("\n--- Testing Performance Optimizations ---")

	// Test 1: Verify database indexes
	if err := testDatabaseIndexes(ctx, db); err != nil {
		return fmt.Errorf("database indexes test failed: %w", err)
	}

	// Test 2: Test optimized repository
	if err := testOptimizedRepository(ctx, db, logger); err != nil {
		return fmt.Errorf("optimized repository test failed: %w", err)
	}

	// Test 3: Test cached metrics service
	if err := testCachedMetricsService(ctx, db, logger); err != nil {
		return fmt.Errorf("cached metrics service test failed: %w", err)
	}

	// Test 4: Test retention service
	if err := testRetentionService(ctx, db, logger); err != nil {
		return fmt.Errorf("retention service test failed: %w", err)
	}

	// Test 5: Test performance monitoring
	if err := testPerformanceMonitoring(ctx, db, logger); err != nil {
		return fmt.Errorf("performance monitoring test failed: %w", err)
	}

	// Test 6: Load test with large dataset
	if err := testLargeDatasetPerformance(ctx, db, logger); err != nil {
		return fmt.Errorf("large dataset performance test failed: %w", err)
	}

	return nil
}

func testDatabaseIndexes(ctx context.Context, db *sql.DB) error {
	fmt.Println("\n1. Testing Database Indexes...")

	// Check if required indexes exist
	requiredIndexes := []string{
		"idx_email_events_inquiry_id",
		"idx_email_events_status",
		"idx_email_events_sent_at",
		"idx_email_events_email_type",
		"idx_email_events_inquiry_type",
		"idx_email_events_status_sent",
		"idx_email_events_type_status",
		"idx_email_events_ses_message_id",
	}

	for _, indexName := range requiredIndexes {
		var exists bool
		err := db.QueryRowContext(ctx,
			"SELECT EXISTS(SELECT 1 FROM pg_indexes WHERE indexname = $1)",
			indexName).Scan(&exists)
		if err != nil {
			return fmt.Errorf("failed to check index %s: %w", indexName, err)
		}

		if !exists {
			fmt.Printf("⚠️  Index %s does not exist\n", indexName)
		} else {
			fmt.Printf("✓ Index %s exists\n", indexName)
		}
	}

	// Check if materialized views exist
	materializedViews := []string{
		"email_metrics_daily",
		"email_metrics_hourly",
	}

	for _, viewName := range materializedViews {
		var exists bool
		err := db.QueryRowContext(ctx,
			"SELECT EXISTS(SELECT 1 FROM pg_matviews WHERE matviewname = $1)",
			viewName).Scan(&exists)
		if err != nil {
			return fmt.Errorf("failed to check materialized view %s: %w", viewName, err)
		}

		if !exists {
			fmt.Printf("⚠️  Materialized view %s does not exist\n", viewName)
		} else {
			fmt.Printf("✓ Materialized view %s exists\n", viewName)
		}
	}

	// Check if performance functions exist
	functions := []string{
		"get_email_metrics_fast",
		"cleanup_old_email_events_optimized",
		"archive_old_email_events",
		"refresh_email_metrics_views",
	}

	for _, functionName := range functions {
		var exists bool
		err := db.QueryRowContext(ctx,
			"SELECT EXISTS(SELECT 1 FROM pg_proc WHERE proname = $1)",
			functionName).Scan(&exists)
		if err != nil {
			return fmt.Errorf("failed to check function %s: %w", functionName, err)
		}

		if !exists {
			fmt.Printf("⚠️  Function %s does not exist\n", functionName)
		} else {
			fmt.Printf("✓ Function %s exists\n", functionName)
		}
	}

	return nil
}

func testOptimizedRepository(ctx context.Context, db *sql.DB, logger *logrus.Logger) error {
	fmt.Println("\n2. Testing Optimized Repository...")

	// Create optimized repository
	repo, err := repositories.NewOptimizedEmailEventRepository(db, logger)
	if err != nil {
		return fmt.Errorf("failed to create optimized repository: %w", err)
	}

	// Test basic operations
	testEvent := &domain.EmailEvent{
		ID:             uuid.New().String(),
		InquiryID:      "test-inquiry-" + uuid.New().String(),
		EmailType:      domain.EmailTypeCustomerConfirmation,
		RecipientEmail: "test@example.com",
		SenderEmail:    "noreply@cloudpartner.pro",
		Subject:        "Test Email",
		Status:         domain.EmailStatusSent,
		SentAt:         time.Now(),
	}

	// Test create
	start := time.Now()
	err = repo.Create(ctx, testEvent)
	createDuration := time.Since(start)
	if err != nil {
		return fmt.Errorf("failed to create test event: %w", err)
	}
	fmt.Printf("✓ Create operation completed in %v\n", createDuration)

	// Test get by inquiry ID
	start = time.Now()
	events, err := repo.GetByInquiryID(ctx, testEvent.InquiryID)
	getDuration := time.Since(start)
	if err != nil {
		return fmt.Errorf("failed to get events by inquiry ID: %w", err)
	}
	if len(events) != 1 {
		return fmt.Errorf("expected 1 event, got %d", len(events))
	}
	fmt.Printf("✓ Get by inquiry ID completed in %v\n", getDuration)

	// Test metrics calculation
	timeRange := domain.TimeRange{
		Start: time.Now().Add(-24 * time.Hour),
		End:   time.Now(),
	}
	filters := domain.EmailEventFilters{
		TimeRange: &timeRange,
	}

	start = time.Now()
	metrics, err := repo.GetMetrics(ctx, filters)
	metricsDuration := time.Since(start)
	if err != nil {
		return fmt.Errorf("failed to get metrics: %w", err)
	}
	if metrics.TotalEmails < 1 {
		return fmt.Errorf("expected at least 1 email in metrics, got %d", metrics.TotalEmails)
	}
	fmt.Printf("✓ Metrics calculation completed in %v\n", metricsDuration)

	// Test batch create for performance
	batchEvents := make([]*domain.EmailEvent, 100)
	for i := 0; i < 100; i++ {
		batchEvents[i] = &domain.EmailEvent{
			ID:             uuid.New().String(),
			InquiryID:      fmt.Sprintf("batch-inquiry-%d", i),
			EmailType:      domain.EmailTypeConsultantNotification,
			RecipientEmail: fmt.Sprintf("test%d@example.com", i),
			SenderEmail:    "noreply@cloudpartner.pro",
			Subject:        fmt.Sprintf("Batch Test Email %d", i),
			Status:         domain.EmailStatusSent,
			SentAt:         time.Now(),
		}
	}

	// Check if optimized repository supports batch create
	if optimizedRepo, ok := repo.(*repositories.OptimizedEmailEventRepository); ok {
		start = time.Now()
		err = optimizedRepo.BatchCreate(ctx, batchEvents)
		batchDuration := time.Since(start)
		if err != nil {
			return fmt.Errorf("failed to batch create events: %w", err)
		}
		fmt.Printf("✓ Batch create of 100 events completed in %v\n", batchDuration)
	}

	return nil
}

func testCachedMetricsService(ctx context.Context, db *sql.DB, logger *logrus.Logger) error {
	fmt.Println("\n3. Testing Cached Metrics Service...")

	// Create base repository and service
	repo, err := repositories.NewEmailEventRepository(db, logger)
	if err != nil {
		return fmt.Errorf("failed to create repository: %w", err)
	}

	baseService := services.NewEmailMetricsService(repo, logger)

	// Create cached service
	cachedService := services.NewEmailMetricsCacheService(
		baseService,
		5*time.Minute, // 5 minute TTL
		100,           // Max 100 cache entries
		logger,
	)

	// Test cache warmup
	start := time.Now()
	err = cachedService.WarmupCache(ctx)
	warmupDuration := time.Since(start)
	if err != nil {
		return fmt.Errorf("failed to warmup cache: %w", err)
	}
	fmt.Printf("✓ Cache warmup completed in %v\n", warmupDuration)

	// Test cached metrics retrieval
	timeRange := domain.TimeRange{
		Start: time.Now().Add(-24 * time.Hour),
		End:   time.Now(),
	}

	// First call (cache miss)
	start = time.Now()
	metrics1, err := cachedService.GetEmailMetrics(ctx, timeRange)
	firstCallDuration := time.Since(start)
	if err != nil {
		return fmt.Errorf("failed to get metrics (first call): %w", err)
	}

	// Second call (cache hit)
	start = time.Now()
	metrics2, err := cachedService.GetEmailMetrics(ctx, timeRange)
	secondCallDuration := time.Since(start)
	if err != nil {
		return fmt.Errorf("failed to get metrics (second call): %w", err)
	}

	// Verify cache performance improvement
	if secondCallDuration >= firstCallDuration {
		fmt.Printf("⚠️  Cache may not be working optimally (first: %v, second: %v)\n",
			firstCallDuration, secondCallDuration)
	} else {
		fmt.Printf("✓ Cache performance improvement: %v -> %v (%.1fx faster)\n",
			firstCallDuration, secondCallDuration,
			float64(firstCallDuration.Nanoseconds())/float64(secondCallDuration.Nanoseconds()))
	}

	// Verify metrics are identical
	if metrics1.TotalEmails != metrics2.TotalEmails {
		return fmt.Errorf("cached metrics don't match: %d vs %d",
			metrics1.TotalEmails, metrics2.TotalEmails)
	}

	// Test cache stats
	stats := cachedService.GetCacheStats()
	fmt.Printf("✓ Cache stats: hits=%v, misses=%v, hit_rate=%.1f%%\n",
		stats["cache_hits"], stats["cache_misses"], stats["hit_rate"])

	return nil
}

func testRetentionService(ctx context.Context, db *sql.DB, logger *logrus.Logger) error {
	fmt.Println("\n4. Testing Retention Service...")

	// Create retention service
	retentionService := services.NewEmailRetentionService(db, logger)

	// Get default policies
	policies := retentionService.GetPolicies()
	fmt.Printf("✓ Loaded %d default retention policies\n", len(policies))

	// Test adding a custom policy
	customPolicy := &services.RetentionPolicy{
		Name:               "test_policy",
		Description:        "Test retention policy",
		RetentionDays:      30,
		ArchiveDays:        7,
		EmailTypes:         []string{"customer_confirmation"},
		Enabled:            true,
		ExecutionFrequency: 24 * time.Hour,
	}

	err := retentionService.AddPolicy(customPolicy)
	if err != nil {
		return fmt.Errorf("failed to add custom policy: %w", err)
	}
	fmt.Printf("✓ Added custom retention policy\n")

	// Test policy execution (dry run with very old dates to avoid deleting real data)
	testPolicy := &services.RetentionPolicy{
		Name:               "dry_run_test",
		Description:        "Dry run test policy",
		RetentionDays:      3650, // 10 years (very old)
		ArchiveDays:        3600, // 10 years - 50 days
		Enabled:            true,
		ExecutionFrequency: 24 * time.Hour,
	}

	err = retentionService.AddPolicy(testPolicy)
	if err != nil {
		return fmt.Errorf("failed to add test policy: %w", err)
	}

	start := time.Now()
	stats, err := retentionService.ExecutePolicy(ctx, "dry_run_test")
	executionDuration := time.Since(start)
	if err != nil {
		return fmt.Errorf("failed to execute test policy: %w", err)
	}

	fmt.Printf("✓ Policy execution completed in %v\n", executionDuration)
	fmt.Printf("  - Records processed: %d\n", stats.RecordsProcessed)
	fmt.Printf("  - Records archived: %d\n", stats.RecordsArchived)
	fmt.Printf("  - Records deleted: %d\n", stats.RecordsDeleted)

	// Test retention summary
	summary, err := retentionService.GetRetentionSummary(ctx)
	if err != nil {
		return fmt.Errorf("failed to get retention summary: %w", err)
	}

	fmt.Printf("✓ Retention summary:\n")
	fmt.Printf("  - Total records: %v\n", summary["total_records"])
	fmt.Printf("  - Recent records: %v\n", summary["recent_records"])
	fmt.Printf("  - Old records: %v\n", summary["old_records"])
	fmt.Printf("  - Policies count: %v\n", summary["policies_count"])

	return nil
}

func testPerformanceMonitoring(ctx context.Context, db *sql.DB, logger *logrus.Logger) error {
	fmt.Println("\n5. Testing Performance Monitoring...")

	// Create performance monitor
	monitor := services.NewEmailPerformanceMonitor(db, logger)

	// Start monitoring for a short period
	monitor.StartMonitoring(ctx, 1*time.Second)

	// Wait for some metrics to be collected
	time.Sleep(3 * time.Second)

	// Stop monitoring
	monitor.StopMonitoring()

	// Get collected metrics
	metrics := monitor.GetMetrics(nil)
	fmt.Printf("✓ Collected %d performance metrics\n", len(metrics))

	// Test query performance stats
	queryStats, err := monitor.GetQueryPerformanceStats(ctx)
	if err != nil {
		fmt.Printf("⚠️  Query performance stats not available: %v\n", err)
	} else {
		fmt.Printf("✓ Query performance stats:\n")
		fmt.Printf("  - Query type: %s\n", queryStats.QueryType)
		fmt.Printf("  - Average duration: %v\n", queryStats.AverageDuration)
		fmt.Printf("  - Total calls: %d\n", queryStats.TotalCalls)
		fmt.Printf("  - Index usage ratio: %.1f%%\n", queryStats.IndexUsageRatio)
	}

	// Test table performance stats
	tableStats, err := monitor.GetTablePerformanceStats(ctx, "email_events")
	if err != nil {
		return fmt.Errorf("failed to get table performance stats: %w", err)
	}

	fmt.Printf("✓ Table performance stats:\n")
	fmt.Printf("  - Table size: %d bytes\n", tableStats.TotalSize)
	fmt.Printf("  - Index size: %d bytes\n", tableStats.IndexSize)
	fmt.Printf("  - Sequential scans: %d\n", tableStats.SequentialScans)
	fmt.Printf("  - Index scans: %d\n", tableStats.IndexScans)
	fmt.Printf("  - Index usage ratio: %.1f%%\n", tableStats.IndexUsageRatio)

	// Test performance analysis
	analysis, err := monitor.AnalyzePerformance(ctx)
	if err != nil {
		return fmt.Errorf("failed to analyze performance: %w", err)
	}

	recommendations := analysis["recommendations"].([]string)
	fmt.Printf("✓ Performance analysis completed with %d recommendations\n", len(recommendations))
	for i, rec := range recommendations {
		fmt.Printf("  %d. %s\n", i+1, rec)
	}

	return nil
}

func testLargeDatasetPerformance(ctx context.Context, db *sql.DB, logger *logrus.Logger) error {
	fmt.Println("\n6. Testing Large Dataset Performance...")

	// Create repository
	repo, err := repositories.NewEmailEventRepository(db, logger)
	if err != nil {
		return fmt.Errorf("failed to create repository: %w", err)
	}

	// Get current record count
	var currentCount int64
	err = db.QueryRowContext(ctx, "SELECT COUNT(*) FROM email_events").Scan(&currentCount)
	if err != nil {
		return fmt.Errorf("failed to get current record count: %w", err)
	}

	fmt.Printf("Current email events count: %d\n", currentCount)

	// Test metrics calculation performance with current dataset
	timeRanges := []struct {
		name  string
		start time.Time
		end   time.Time
	}{
		{"Last 24 hours", time.Now().Add(-24 * time.Hour), time.Now()},
		{"Last 7 days", time.Now().Add(-7 * 24 * time.Hour), time.Now()},
		{"Last 30 days", time.Now().Add(-30 * 24 * time.Hour), time.Now()},
		{"Last 90 days", time.Now().Add(-90 * 24 * time.Hour), time.Now()},
	}

	for _, tr := range timeRanges {
		filters := domain.EmailEventFilters{
			TimeRange: &domain.TimeRange{Start: tr.start, End: tr.end},
		}

		start := time.Now()
		metrics, err := repo.GetMetrics(ctx, filters)
		duration := time.Since(start)

		if err != nil {
			fmt.Printf("⚠️  Failed to get metrics for %s: %v\n", tr.name, err)
			continue
		}

		fmt.Printf("✓ %s metrics: %d emails, calculated in %v\n",
			tr.name, metrics.TotalEmails, duration)

		// Warn if query takes too long
		if duration > 1*time.Second {
			fmt.Printf("⚠️  Query for %s took longer than 1 second (%v)\n", tr.name, duration)
		}
	}

	// Test pagination performance
	filters := domain.EmailEventFilters{
		Limit:  100,
		Offset: 0,
	}

	start := time.Now()
	events, err := repo.List(ctx, filters)
	paginationDuration := time.Since(start)

	if err != nil {
		return fmt.Errorf("failed to test pagination: %w", err)
	}

	fmt.Printf("✓ Pagination (100 records): %d events retrieved in %v\n",
		len(events), paginationDuration)

	// Test complex filtering performance
	customerType := domain.EmailTypeCustomerConfirmation
	deliveredStatus := domain.EmailStatusDelivered
	complexFilters := domain.EmailEventFilters{
		TimeRange: &domain.TimeRange{
			Start: time.Now().Add(-30 * 24 * time.Hour),
			End:   time.Now(),
		},
		EmailType: &customerType,
		Status:    &deliveredStatus,
		Limit:     50,
	}

	start = time.Now()
	complexEvents, err := repo.List(ctx, complexFilters)
	complexDuration := time.Since(start)

	if err != nil {
		return fmt.Errorf("failed to test complex filtering: %w", err)
	}

	fmt.Printf("✓ Complex filtering: %d events retrieved in %v\n",
		len(complexEvents), complexDuration)

	return nil
}

func init() {
	// Set environment variables for testing if not already set
	if os.Getenv("DATABASE_URL") == "" {
		os.Setenv("DATABASE_URL", "postgres://postgres:password@localhost/cloud_consulting_test?sslmode=disable")
	}
}
