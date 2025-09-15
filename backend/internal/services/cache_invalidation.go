package services

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/cloud-consulting/backend/internal/storage"
)

// CacheInvalidationService manages cache invalidation and consistency
type CacheInvalidationService struct {
	cache   *storage.RedisCache
	monitor *CacheMonitor
	logger  *logrus.Logger

	// Invalidation tracking
	invalidationQueue chan InvalidationRequest
	batchSize         int
	batchTimeout      time.Duration

	// Consistency checking
	consistencyChecks map[string]time.Time
	consistencyMutex  sync.RWMutex

	// Configuration
	maxRetries     int
	retryDelay     time.Duration
	consistencyTTL time.Duration
}

// InvalidationRequest represents a cache invalidation request
type InvalidationRequest struct {
	Type      InvalidationType `json:"type"`
	Key       string           `json:"key"`
	Pattern   string           `json:"pattern,omitempty"`
	UserID    string           `json:"user_id,omitempty"`
	SessionID string           `json:"session_id,omitempty"`
	Timestamp time.Time        `json:"timestamp"`
	Retries   int              `json:"retries"`
}

// InvalidationType represents the type of cache invalidation
type InvalidationType string

const (
	InvalidationTypeSession         InvalidationType = "session"
	InvalidationTypeUserSessions    InvalidationType = "user_sessions"
	InvalidationTypeSessionMessages InvalidationType = "session_messages"
	InvalidationTypePattern         InvalidationType = "pattern"
	InvalidationTypeAll             InvalidationType = "all"
)

// NewCacheInvalidationService creates a new cache invalidation service
func NewCacheInvalidationService(
	cache *storage.RedisCache,
	monitor *CacheMonitor,
	logger *logrus.Logger,
) *CacheInvalidationService {
	service := &CacheInvalidationService{
		cache:             cache,
		monitor:           monitor,
		logger:            logger,
		invalidationQueue: make(chan InvalidationRequest, 1000),
		batchSize:         10,
		batchTimeout:      100 * time.Millisecond,
		consistencyChecks: make(map[string]time.Time),
		maxRetries:        3,
		retryDelay:        1 * time.Second,
		consistencyTTL:    5 * time.Minute,
	}

	return service
}

// Start starts the cache invalidation service
func (s *CacheInvalidationService) Start(ctx context.Context) {
	s.logger.Info("Starting cache invalidation service")

	// Start batch processor
	go s.processBatchInvalidations(ctx)

	// Start consistency checker
	go s.runConsistencyChecker(ctx)

	s.logger.Info("Cache invalidation service started")
}

// InvalidateSession invalidates a specific session cache
func (s *CacheInvalidationService) InvalidateSession(ctx context.Context, sessionID string) error {
	request := InvalidationRequest{
		Type:      InvalidationTypeSession,
		SessionID: sessionID,
		Key:       fmt.Sprintf("chat:session:%s", sessionID),
		Timestamp: time.Now(),
	}

	return s.queueInvalidation(request)
}

// InvalidateUserSessions invalidates user sessions cache
func (s *CacheInvalidationService) InvalidateUserSessions(ctx context.Context, userID string) error {
	request := InvalidationRequest{
		Type:      InvalidationTypeUserSessions,
		UserID:    userID,
		Key:       fmt.Sprintf("chat:user_sessions:%s", userID),
		Timestamp: time.Now(),
	}

	return s.queueInvalidation(request)
}

// InvalidateSessionMessages invalidates session messages cache
func (s *CacheInvalidationService) InvalidateSessionMessages(ctx context.Context, sessionID string) error {
	request := InvalidationRequest{
		Type:      InvalidationTypeSessionMessages,
		SessionID: sessionID,
		Key:       fmt.Sprintf("chat:session_messages:%s", sessionID),
		Timestamp: time.Now(),
	}

	return s.queueInvalidation(request)
}

// InvalidatePattern invalidates cache entries matching a pattern
func (s *CacheInvalidationService) InvalidatePattern(ctx context.Context, pattern string) error {
	request := InvalidationRequest{
		Type:      InvalidationTypePattern,
		Pattern:   pattern,
		Timestamp: time.Now(),
	}

	return s.queueInvalidation(request)
}

// InvalidateAll invalidates all chat-related cache entries
func (s *CacheInvalidationService) InvalidateAll(ctx context.Context) error {
	request := InvalidationRequest{
		Type:      InvalidationTypeAll,
		Pattern:   "chat:*",
		Timestamp: time.Now(),
	}

	return s.queueInvalidation(request)
}

// queueInvalidation queues an invalidation request for batch processing
func (s *CacheInvalidationService) queueInvalidation(request InvalidationRequest) error {
	select {
	case s.invalidationQueue <- request:
		s.logger.WithFields(logrus.Fields{
			"type":       request.Type,
			"key":        request.Key,
			"session_id": request.SessionID,
			"user_id":    request.UserID,
		}).Debug("Queued cache invalidation request")
		return nil
	default:
		s.logger.WithFields(logrus.Fields{
			"type":       request.Type,
			"key":        request.Key,
			"session_id": request.SessionID,
			"user_id":    request.UserID,
		}).Warn("Cache invalidation queue is full, dropping request")
		return fmt.Errorf("invalidation queue is full")
	}
}

// processBatchInvalidations processes invalidation requests in batches
func (s *CacheInvalidationService) processBatchInvalidations(ctx context.Context) {
	ticker := time.NewTicker(s.batchTimeout)
	defer ticker.Stop()

	var batch []InvalidationRequest

	for {
		select {
		case <-ctx.Done():
			// Process remaining batch before shutting down
			if len(batch) > 0 {
				s.processBatch(ctx, batch)
			}
			s.logger.Info("Cache invalidation processor stopped")
			return

		case request := <-s.invalidationQueue:
			batch = append(batch, request)

			// Process batch if it reaches the batch size
			if len(batch) >= s.batchSize {
				s.processBatch(ctx, batch)
				batch = batch[:0] // Reset batch
			}

		case <-ticker.C:
			// Process batch on timeout if it has any items
			if len(batch) > 0 {
				s.processBatch(ctx, batch)
				batch = batch[:0] // Reset batch
			}
		}
	}
}

// processBatch processes a batch of invalidation requests
func (s *CacheInvalidationService) processBatch(ctx context.Context, batch []InvalidationRequest) {
	startTime := time.Now()

	s.logger.WithField("batch_size", len(batch)).Debug("Processing cache invalidation batch")

	for _, request := range batch {
		if err := s.processInvalidationRequest(ctx, request); err != nil {
			s.logger.WithError(err).WithFields(logrus.Fields{
				"type":       request.Type,
				"key":        request.Key,
				"session_id": request.SessionID,
				"user_id":    request.UserID,
				"retries":    request.Retries,
			}).Error("Failed to process invalidation request")

			// Retry if under retry limit
			if request.Retries < s.maxRetries {
				request.Retries++
				time.AfterFunc(s.retryDelay, func() {
					s.queueInvalidation(request)
				})
			}
		}
	}

	processingTime := time.Since(startTime)
	s.monitor.RecordCacheInvalidation(processingTime)

	s.logger.WithFields(logrus.Fields{
		"batch_size":      len(batch),
		"processing_time": processingTime,
	}).Debug("Completed cache invalidation batch")
}

// processInvalidationRequest processes a single invalidation request
func (s *CacheInvalidationService) processInvalidationRequest(ctx context.Context, request InvalidationRequest) error {
	switch request.Type {
	case InvalidationTypeSession:
		return s.cache.DeleteSession(ctx, request.SessionID)

	case InvalidationTypeUserSessions:
		return s.cache.InvalidateUserSessions(ctx, request.UserID)

	case InvalidationTypeSessionMessages:
		return s.cache.InvalidateSessionMessages(ctx, request.SessionID)

	case InvalidationTypePattern:
		return s.invalidateByPattern(ctx, request.Pattern)

	case InvalidationTypeAll:
		return s.invalidateByPattern(ctx, "chat:*")

	default:
		return fmt.Errorf("unknown invalidation type: %s", request.Type)
	}
}

// invalidateByPattern invalidates cache entries matching a pattern
func (s *CacheInvalidationService) invalidateByPattern(ctx context.Context, pattern string) error {
	// Note: This is a simplified implementation
	// In a production environment, you might want to use Redis SCAN command
	// to avoid blocking the Redis server with KEYS command

	s.logger.WithField("pattern", pattern).Debug("Invalidating cache entries by pattern")

	// For now, we'll implement specific patterns we know about
	switch pattern {
	case "chat:*":
		// Invalidate all chat-related cache entries
		// This is a simplified approach - in production, you'd want to be more selective
		s.logger.Warn("Invalidating all chat cache entries - this may impact performance")
		return nil // Placeholder - implement actual pattern-based invalidation

	default:
		s.logger.WithField("pattern", pattern).Warn("Pattern-based invalidation not fully implemented")
		return nil
	}
}

// CheckConsistency checks cache consistency for a specific key
func (s *CacheInvalidationService) CheckConsistency(ctx context.Context, key string) error {
	s.consistencyMutex.RLock()
	lastCheck, exists := s.consistencyChecks[key]
	s.consistencyMutex.RUnlock()

	// Skip if recently checked
	if exists && time.Since(lastCheck) < s.consistencyTTL {
		return nil
	}

	// Update last check time
	s.consistencyMutex.Lock()
	s.consistencyChecks[key] = time.Now()
	s.consistencyMutex.Unlock()

	// Perform consistency check
	// This is a placeholder for actual consistency checking logic
	s.logger.WithField("key", key).Debug("Performing cache consistency check")

	return nil
}

// runConsistencyChecker runs periodic consistency checks
func (s *CacheInvalidationService) runConsistencyChecker(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	s.logger.Info("Started cache consistency checker")

	for {
		select {
		case <-ctx.Done():
			s.logger.Info("Cache consistency checker stopped")
			return

		case <-ticker.C:
			s.performConsistencyMaintenance()
		}
	}
}

// performConsistencyMaintenance performs periodic consistency maintenance
func (s *CacheInvalidationService) performConsistencyMaintenance() {
	s.consistencyMutex.Lock()
	defer s.consistencyMutex.Unlock()

	// Clean up old consistency check records
	cutoff := time.Now().Add(-s.consistencyTTL)
	for key, lastCheck := range s.consistencyChecks {
		if lastCheck.Before(cutoff) {
			delete(s.consistencyChecks, key)
		}
	}

	s.logger.WithField("active_checks", len(s.consistencyChecks)).Debug("Performed consistency maintenance")
}

// GetInvalidationStats returns invalidation statistics
func (s *CacheInvalidationService) GetInvalidationStats() *InvalidationStats {
	return &InvalidationStats{
		QueueSize:       len(s.invalidationQueue),
		QueueCapacity:   cap(s.invalidationQueue),
		ActiveChecks:    s.getActiveChecksCount(),
		BatchSize:       s.batchSize,
		BatchTimeout:    s.batchTimeout,
		MaxRetries:      s.maxRetries,
		RetryDelay:      s.retryDelay,
		ConsistencyTTL:  s.consistencyTTL,
		LastMaintenance: time.Now(), // Placeholder
	}
}

// getActiveChecksCount returns the number of active consistency checks
func (s *CacheInvalidationService) getActiveChecksCount() int {
	s.consistencyMutex.RLock()
	defer s.consistencyMutex.RUnlock()
	return len(s.consistencyChecks)
}

// InvalidationStats represents invalidation service statistics
type InvalidationStats struct {
	QueueSize       int           `json:"queue_size"`
	QueueCapacity   int           `json:"queue_capacity"`
	ActiveChecks    int           `json:"active_checks"`
	BatchSize       int           `json:"batch_size"`
	BatchTimeout    time.Duration `json:"batch_timeout"`
	MaxRetries      int           `json:"max_retries"`
	RetryDelay      time.Duration `json:"retry_delay"`
	ConsistencyTTL  time.Duration `json:"consistency_ttl"`
	LastMaintenance time.Time     `json:"last_maintenance"`
}

// CacheConsistencyReport represents a cache consistency report
type CacheConsistencyReport struct {
	CheckedKeys      []string  `json:"checked_keys"`
	InconsistentKeys []string  `json:"inconsistent_keys"`
	TotalChecks      int       `json:"total_checks"`
	FailedChecks     int       `json:"failed_checks"`
	GeneratedAt      time.Time `json:"generated_at"`
}

// GenerateConsistencyReport generates a cache consistency report
func (s *CacheInvalidationService) GenerateConsistencyReport(ctx context.Context) *CacheConsistencyReport {
	s.consistencyMutex.RLock()
	defer s.consistencyMutex.RUnlock()

	var checkedKeys []string
	for key := range s.consistencyChecks {
		checkedKeys = append(checkedKeys, key)
	}

	report := &CacheConsistencyReport{
		CheckedKeys:      checkedKeys,
		InconsistentKeys: []string{}, // Placeholder - implement actual inconsistency detection
		TotalChecks:      len(checkedKeys),
		FailedChecks:     0, // Placeholder
		GeneratedAt:      time.Now(),
	}

	return report
}
