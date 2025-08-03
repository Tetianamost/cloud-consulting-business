package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/cloud-consulting/backend/internal/domain"
	"github.com/cloud-consulting/backend/internal/repositories"
	"github.com/cloud-consulting/backend/internal/services"
	"github.com/cloud-consulting/backend/internal/storage"
)

func main() {
	// Set up logger
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	ctx := context.Background()

	// Test database connection pooling
	fmt.Println("=== Testing Database Connection Pooling ===")
	if err := testDatabasePool(logger); err != nil {
		log.Fatalf("Database pool test failed: %v", err)
	}
	fmt.Println("✓ Database connection pooling test passed")

	// Test Redis caching
	fmt.Println("\n=== Testing Redis Caching ===")
	if err := testRedisCache(ctx, logger); err != nil {
		log.Fatalf("Redis cache test failed: %v", err)
	}
	fmt.Println("✓ Redis caching test passed")

	// Test enhanced repositories
	fmt.Println("\n=== Testing Enhanced Repositories ===")
	if err := testEnhancedRepositories(ctx, logger); err != nil {
		log.Fatalf("Enhanced repositories test failed: %v", err)
	}
	fmt.Println("✓ Enhanced repositories test passed")

	// Test cache monitoring
	fmt.Println("\n=== Testing Cache Monitoring ===")
	if err := testCacheMonitoring(ctx, logger); err != nil {
		log.Fatalf("Cache monitoring test failed: %v", err)
	}
	fmt.Println("✓ Cache monitoring test passed")

	// Test cache invalidation
	fmt.Println("\n=== Testing Cache Invalidation ===")
	if err := testCacheInvalidation(ctx, logger); err != nil {
		log.Fatalf("Cache invalidation test failed: %v", err)
	}
	fmt.Println("✓ Cache invalidation test passed")

	fmt.Println("\n🎉 All enhanced repository and caching tests passed!")
}

func testDatabasePool(logger *logrus.Logger) error {
	// Create database configuration
	config := &storage.DatabaseConfig{
		Host:            getEnv("DB_HOST", "localhost"),
		Port:            5432,
		Database:        getEnv("DB_NAME", "consulting"),
		Username:        getEnv("DB_USER", "postgres"),
		Password:        getEnv("DB_PASSWORD", "password"),
		SSLMode:         "disable",
		MaxOpenConns:    10,
		MaxIdleConns:    5,
		ConnMaxLifetime: 5 * time.Minute,
		ConnMaxIdleTime: 1 * time.Minute,
	}

	// Create database pool
	pool, err := storage.NewDatabasePool(config, logger)
	if err != nil {
		return fmt.Errorf("failed to create database pool: %w", err)
	}
	defer pool.Close()

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// Test health check
	if !pool.IsHealthy(ctx) {
		return fmt.Errorf("database health check failed")
	}

	// Test connection pool statistics
	stats := pool.GetStats()
	logger.WithFields(logrus.Fields{
		"open_connections": stats.OpenConnections,
		"in_use":           stats.InUse,
		"idle":             stats.Idle,
	}).Info("Database pool statistics")

	// Test transaction
	err = pool.WithTx(ctx, func(tx *sql.Tx) error {
		// Simple test query
		var count int
		err := tx.QueryRowContext(ctx, "SELECT COUNT(*) FROM chat_sessions").Scan(&count)
		if err != nil {
			return fmt.Errorf("failed to execute test query: %w", err)
		}
		logger.WithField("session_count", count).Info("Test transaction executed successfully")
		return nil
	})
	if err != nil {
		return fmt.Errorf("transaction test failed: %w", err)
	}

	return nil
}

func testRedisCache(ctx context.Context, logger *logrus.Logger) error {
	// Create Redis configuration
	config := &storage.RedisCacheConfig{
		Host:               getEnv("REDIS_HOST", "localhost"),
		Port:               6379,
		Password:           getEnv("REDIS_PASSWORD", ""),
		Database:           0,
		MaxRetries:         3,
		PoolSize:           10,
		MinIdleConns:       2,
		DialTimeout:        5 * time.Second,
		ReadTimeout:        3 * time.Second,
		WriteTimeout:       3 * time.Second,
		PoolTimeout:        4 * time.Second,
		IdleTimeout:        5 * time.Minute,
		IdleCheckFrequency: 1 * time.Minute,
		SessionTTL:         30 * time.Minute,
		MessageTTL:         15 * time.Minute,
	}

	// Create Redis cache
	cache, err := storage.NewRedisCache(config, logger)
	if err != nil {
		return fmt.Errorf("failed to create Redis cache: %w", err)
	}
	defer cache.Close()

	// Test connection
	if err := cache.Ping(ctx); err != nil {
		return fmt.Errorf("failed to ping Redis: %w", err)
	}

	// Test health check
	if !cache.IsHealthy(ctx) {
		return fmt.Errorf("Redis health check failed")
	}

	// Test session caching
	testSession := &domain.ChatSession{
		ID:           uuid.New().String(),
		UserID:       "test_user",
		ClientName:   "Test Client",
		Context:      "Test context",
		Status:       domain.SessionStatusActive,
		Metadata:     map[string]interface{}{"test": "data"},
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		LastActivity: time.Now(),
	}

	// Set session in cache
	if err := cache.SetSession(ctx, testSession); err != nil {
		return fmt.Errorf("failed to set session in cache: %w", err)
	}

	// Get session from cache
	cachedSession, err := cache.GetSession(ctx, testSession.ID)
	if err != nil {
		return fmt.Errorf("failed to get session from cache: %w", err)
	}

	if cachedSession == nil {
		return fmt.Errorf("session not found in cache")
	}

	if cachedSession.ID != testSession.ID {
		return fmt.Errorf("cached session ID mismatch: expected %s, got %s", testSession.ID, cachedSession.ID)
	}

	// Test message caching
	testMessages := []*domain.ChatMessage{
		{
			ID:        uuid.New().String(),
			SessionID: testSession.ID,
			Type:      domain.MessageTypeUser,
			Content:   "Test message 1",
			Metadata:  map[string]interface{}{"test": "data1"},
			Status:    domain.MessageStatusDelivered,
			CreatedAt: time.Now(),
		},
		{
			ID:        uuid.New().String(),
			SessionID: testSession.ID,
			Type:      domain.MessageTypeAssistant,
			Content:   "Test response 1",
			Metadata:  map[string]interface{}{"test": "data2"},
			Status:    domain.MessageStatusDelivered,
			CreatedAt: time.Now().Add(1 * time.Second),
		},
	}

	// Set messages in cache
	if err := cache.SetSessionMessages(ctx, testSession.ID, testMessages); err != nil {
		return fmt.Errorf("failed to set messages in cache: %w", err)
	}

	// Get messages from cache
	cachedMessages, err := cache.GetSessionMessages(ctx, testSession.ID)
	if err != nil {
		return fmt.Errorf("failed to get messages from cache: %w", err)
	}

	if len(cachedMessages) != len(testMessages) {
		return fmt.Errorf("cached messages count mismatch: expected %d, got %d", len(testMessages), len(cachedMessages))
	}

	// Test cache deletion
	if err := cache.DeleteSession(ctx, testSession.ID); err != nil {
		return fmt.Errorf("failed to delete session from cache: %w", err)
	}

	// Verify deletion
	deletedSession, err := cache.GetSession(ctx, testSession.ID)
	if err != nil {
		return fmt.Errorf("error checking deleted session: %w", err)
	}

	if deletedSession != nil {
		return fmt.Errorf("session was not deleted from cache")
	}

	// Test cache statistics
	stats, err := cache.GetStats(ctx)
	if err != nil {
		return fmt.Errorf("failed to get cache stats: %w", err)
	}

	logger.WithFields(logrus.Fields{
		"is_healthy": stats.IsHealthy,
	}).Info("Redis cache statistics")

	return nil
}

func testEnhancedRepositories(ctx context.Context, logger *logrus.Logger) error {
	// Set up database pool
	dbConfig := &storage.DatabaseConfig{
		Host:            getEnv("DB_HOST", "localhost"),
		Port:            5432,
		Database:        getEnv("DB_NAME", "consulting"),
		Username:        getEnv("DB_USER", "postgres"),
		Password:        getEnv("DB_PASSWORD", "password"),
		SSLMode:         "disable",
		MaxOpenConns:    10,
		MaxIdleConns:    5,
		ConnMaxLifetime: 5 * time.Minute,
		ConnMaxIdleTime: 1 * time.Minute,
	}

	pool, err := storage.NewDatabasePool(dbConfig, logger)
	if err != nil {
		return fmt.Errorf("failed to create database pool: %w", err)
	}
	defer pool.Close()

	// Set up Redis cache
	cacheConfig := storage.DefaultRedisCacheConfig()
	cacheConfig.Host = getEnv("REDIS_HOST", "localhost")
	cacheConfig.Password = getEnv("REDIS_PASSWORD", "")

	cache, err := storage.NewRedisCache(cacheConfig, logger)
	if err != nil {
		// If Redis is not available, test without cache
		logger.Warn("Redis not available, testing repositories without cache")
		cache = nil
	} else {
		defer cache.Close()
	}

	// Create enhanced repositories
	sessionRepo := repositories.NewEnhancedChatSessionRepository(pool, cache, logger)
	messageRepo := repositories.NewEnhancedChatMessageRepository(pool, cache, logger)

	// Test session repository
	testSession := &domain.ChatSession{
		ID:           uuid.New().String(),
		UserID:       "test_enhanced_user",
		ClientName:   "Enhanced Test Client",
		Context:      "Enhanced test context",
		Status:       domain.SessionStatusActive,
		Metadata:     map[string]interface{}{"enhanced": "test"},
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		LastActivity: time.Now(),
	}

	// Create session
	if err := sessionRepo.Create(ctx, testSession); err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}

	// Get session (should hit cache if available)
	retrievedSession, err := sessionRepo.GetByID(ctx, testSession.ID)
	if err != nil {
		return fmt.Errorf("failed to get session: %w", err)
	}

	if retrievedSession == nil {
		return fmt.Errorf("session not found")
	}

	if retrievedSession.ID != testSession.ID {
		return fmt.Errorf("session ID mismatch")
	}

	// Test message repository
	testMessage := &domain.ChatMessage{
		ID:        uuid.New().String(),
		SessionID: testSession.ID,
		Type:      domain.MessageTypeUser,
		Content:   "Enhanced test message",
		Metadata:  map[string]interface{}{"enhanced": "message"},
		Status:    domain.MessageStatusDelivered,
		CreatedAt: time.Now(),
	}

	// Create message
	if err := messageRepo.Create(ctx, testMessage); err != nil {
		return fmt.Errorf("failed to create message: %w", err)
	}

	// Get messages for session
	messages, err := messageRepo.GetBySessionID(ctx, testSession.ID, 10, 0)
	if err != nil {
		return fmt.Errorf("failed to get messages: %w", err)
	}

	if len(messages) == 0 {
		return fmt.Errorf("no messages found")
	}

	// Test session update (should invalidate cache)
	testSession.Context = "Updated context"
	testSession.UpdatedAt = time.Now()
	if err := sessionRepo.Update(ctx, testSession); err != nil {
		return fmt.Errorf("failed to update session: %w", err)
	}

	// Get updated session
	updatedSession, err := sessionRepo.GetByID(ctx, testSession.ID)
	if err != nil {
		return fmt.Errorf("failed to get updated session: %w", err)
	}

	if updatedSession.Context != "Updated context" {
		return fmt.Errorf("session update not reflected")
	}

	// Test user sessions retrieval
	userSessions, err := sessionRepo.GetByUserID(ctx, testSession.UserID)
	if err != nil {
		return fmt.Errorf("failed to get user sessions: %w", err)
	}

	if len(userSessions) == 0 {
		return fmt.Errorf("no user sessions found")
	}

	// Clean up test data
	if err := messageRepo.DeleteBySessionID(ctx, testSession.ID); err != nil {
		return fmt.Errorf("failed to delete messages: %w", err)
	}

	if err := sessionRepo.Delete(ctx, testSession.ID); err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	logger.Info("Enhanced repositories test completed successfully")
	return nil
}

func testCacheMonitoring(ctx context.Context, logger *logrus.Logger) error {
	// Set up Redis cache
	cacheConfig := storage.DefaultRedisCacheConfig()
	cacheConfig.Host = getEnv("REDIS_HOST", "localhost")
	cacheConfig.Password = getEnv("REDIS_PASSWORD", "")

	cache, err := storage.NewRedisCache(cacheConfig, logger)
	if err != nil {
		logger.Warn("Redis not available, skipping cache monitoring test")
		return nil
	}
	defer cache.Close()

	// Create cache monitor
	monitor := services.NewCacheMonitor(cache, logger)

	// Simulate cache operations
	monitor.RecordSessionHit(10 * time.Millisecond)
	monitor.RecordSessionMiss(50 * time.Millisecond)
	monitor.RecordMessageHit(5 * time.Millisecond)
	monitor.RecordSessionSet(20 * time.Millisecond)

	// Get metrics
	metrics := monitor.GetMetrics()
	if metrics.SessionHits != 1 {
		return fmt.Errorf("expected 1 session hit, got %d", metrics.SessionHits)
	}

	if metrics.SessionMisses != 1 {
		return fmt.Errorf("expected 1 session miss, got %d", metrics.SessionMisses)
	}

	// Test hit ratios
	sessionHitRatio := monitor.GetSessionHitRatio()
	if sessionHitRatio != 0.5 { // 1 hit out of 2 gets
		return fmt.Errorf("expected session hit ratio 0.5, got %f", sessionHitRatio)
	}

	// Test health status
	healthStatus := monitor.GetHealthStatus(ctx)
	if healthStatus.Status == "" {
		return fmt.Errorf("health status not set")
	}

	// Test performance report
	report := monitor.GeneratePerformanceReport(ctx)
	if report.Metrics == nil {
		return fmt.Errorf("performance report metrics not set")
	}

	if len(report.Recommendations) == 0 {
		return fmt.Errorf("performance report recommendations not set")
	}

	// Log metrics
	monitor.LogMetrics()

	logger.Info("Cache monitoring test completed successfully")
	return nil
}

func testCacheInvalidation(ctx context.Context, logger *logrus.Logger) error {
	// Set up Redis cache
	cacheConfig := storage.DefaultRedisCacheConfig()
	cacheConfig.Host = getEnv("REDIS_HOST", "localhost")
	cacheConfig.Password = getEnv("REDIS_PASSWORD", "")

	cache, err := storage.NewRedisCache(cacheConfig, logger)
	if err != nil {
		logger.Warn("Redis not available, skipping cache invalidation test")
		return nil
	}
	defer cache.Close()

	// Create cache monitor and invalidation service
	monitor := services.NewCacheMonitor(cache, logger)
	invalidationService := services.NewCacheInvalidationService(cache, monitor, logger)

	// Start invalidation service
	serviceCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	invalidationService.Start(serviceCtx)

	// Test session invalidation
	sessionID := uuid.New().String()
	if err := invalidationService.InvalidateSession(ctx, sessionID); err != nil {
		return fmt.Errorf("failed to invalidate session: %w", err)
	}

	// Test user sessions invalidation
	userID := "test_user"
	if err := invalidationService.InvalidateUserSessions(ctx, userID); err != nil {
		return fmt.Errorf("failed to invalidate user sessions: %w", err)
	}

	// Test session messages invalidation
	if err := invalidationService.InvalidateSessionMessages(ctx, sessionID); err != nil {
		return fmt.Errorf("failed to invalidate session messages: %w", err)
	}

	// Wait a bit for batch processing
	time.Sleep(200 * time.Millisecond)

	// Test invalidation stats
	stats := invalidationService.GetInvalidationStats()
	if stats.QueueCapacity == 0 {
		return fmt.Errorf("invalidation stats not properly initialized")
	}

	// Test consistency report
	report := invalidationService.GenerateConsistencyReport(ctx)
	if report.GeneratedAt.IsZero() {
		return fmt.Errorf("consistency report not properly generated")
	}

	logger.Info("Cache invalidation test completed successfully")
	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
