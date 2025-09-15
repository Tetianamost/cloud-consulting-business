package services

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/cloud-consulting/backend/internal/domain"
	"github.com/cloud-consulting/backend/internal/interfaces"
)

// CacheEntry represents a cached email metrics entry
type CacheEntry struct {
	Data      *domain.EmailMetrics `json:"data"`
	ExpiresAt time.Time            `json:"expires_at"`
	CreatedAt time.Time            `json:"created_at"`
}

// EmailMetricsCacheService provides caching layer for email metrics
type EmailMetricsCacheService struct {
	baseService  interfaces.EmailMetricsService
	cache        map[string]*CacheEntry
	cacheMutex   sync.RWMutex
	defaultTTL   time.Duration
	maxCacheSize int
	logger       *logrus.Logger

	// Performance metrics
	cacheHits      int64
	cacheMisses    int64
	cacheEvictions int64
	lastCleanup    time.Time
}

// NewEmailMetricsCacheService creates a new cached email metrics service
func NewEmailMetricsCacheService(
	baseService interfaces.EmailMetricsService,
	defaultTTL time.Duration,
	maxCacheSize int,
	logger *logrus.Logger,
) *EmailMetricsCacheService {
	service := &EmailMetricsCacheService{
		baseService:  baseService,
		cache:        make(map[string]*CacheEntry),
		defaultTTL:   defaultTTL,
		maxCacheSize: maxCacheSize,
		logger:       logger,
		lastCleanup:  time.Now(),
	}

	// Start background cleanup goroutine
	go service.cleanupExpiredEntries()

	return service
}

// GetEmailMetrics returns cached email metrics or calculates and caches them
func (s *EmailMetricsCacheService) GetEmailMetrics(ctx context.Context, timeRange domain.TimeRange) (*domain.EmailMetrics, error) {
	cacheKey := s.generateCacheKey(timeRange, nil, nil)

	// Try to get from cache first
	if cachedMetrics := s.getFromCache(cacheKey); cachedMetrics != nil {
		s.cacheHits++
		s.logger.WithFields(logrus.Fields{
			"cache_key":  cacheKey,
			"cache_hits": s.cacheHits,
			"action":     "cache_hit",
		}).Debug("Email metrics served from cache")
		return cachedMetrics, nil
	}

	// Cache miss - get from base service
	s.cacheMisses++
	metrics, err := s.baseService.GetEmailMetrics(ctx, timeRange)
	if err != nil {
		s.logger.WithError(err).WithField("cache_key", cacheKey).Error("Failed to get email metrics from base service")
		return nil, err
	}

	// Cache the result
	s.setInCache(cacheKey, metrics, s.defaultTTL)

	s.logger.WithFields(logrus.Fields{
		"cache_key":    cacheKey,
		"cache_misses": s.cacheMisses,
		"action":       "cache_miss_and_store",
	}).Debug("Email metrics calculated and cached")

	return metrics, nil
}

// GetEmailMetricsByType returns cached email metrics by type
func (s *EmailMetricsCacheService) GetEmailMetricsByType(ctx context.Context, timeRange domain.TimeRange) (map[domain.EmailEventType]*domain.EmailMetrics, error) {
	cacheKey := s.generateCacheKey(timeRange, nil, nil) + "_by_type"

	// Try to get from cache first
	s.cacheMutex.RLock()
	if entry, exists := s.cache[cacheKey]; exists && time.Now().Before(entry.ExpiresAt) {
		s.cacheMutex.RUnlock()
		s.cacheHits++

		// Deserialize the cached data
		var cachedData map[domain.EmailEventType]*domain.EmailMetrics
		if dataBytes, err := json.Marshal(entry.Data); err == nil {
			if err := json.Unmarshal(dataBytes, &cachedData); err == nil {
				s.logger.WithField("cache_key", cacheKey).Debug("Email metrics by type served from cache")
				return cachedData, nil
			}
		}
	}
	s.cacheMutex.RUnlock()

	// Cache miss - get from base service
	s.cacheMisses++

	// Check if base service supports this method
	if baseServiceWithType, ok := s.baseService.(*EmailMetricsServiceImpl); ok {
		metrics, err := baseServiceWithType.GetEmailMetricsByType(ctx, timeRange)
		if err != nil {
			return nil, err
		}

		// Cache the result (serialize as JSON for complex types)
		if dataBytes, err := json.Marshal(metrics); err == nil {
			var genericData interface{}
			if err := json.Unmarshal(dataBytes, &genericData); err == nil {
				s.setInCache(cacheKey, &domain.EmailMetrics{TimeRange: fmt.Sprintf("by_type_%s", timeRange.Start.Format("2006-01-02"))}, s.defaultTTL)
			}
		}

		return metrics, nil
	}

	return nil, fmt.Errorf("base service does not support GetEmailMetricsByType")
}

// GetEmailStatusByInquiry returns email status for an inquiry (not cached due to real-time nature)
func (s *EmailMetricsCacheService) GetEmailStatusByInquiry(ctx context.Context, inquiryID string) (*domain.EmailStatus, error) {
	// Don't cache inquiry-specific status as it changes frequently
	return s.baseService.GetEmailStatusByInquiry(ctx, inquiryID)
}

// GetEmailEventHistory returns email event history (with limited caching for recent queries)
func (s *EmailMetricsCacheService) GetEmailEventHistory(ctx context.Context, filters domain.EmailEventFilters) ([]*domain.EmailEvent, error) {
	// Only cache if it's a recent time range query without specific inquiry ID
	if filters.InquiryID == nil && filters.TimeRange != nil {
		timeDiff := time.Now().Sub(filters.TimeRange.End)
		if timeDiff > time.Hour { // Only cache queries for data older than 1 hour
			cacheKey := s.generateCacheKey(*filters.TimeRange, filters.EmailType, filters.Status) + "_history"

			// Try cache first
			s.cacheMutex.RLock()
			if entry, exists := s.cache[cacheKey]; exists && time.Now().Before(entry.ExpiresAt) {
				s.cacheMutex.RUnlock()
				s.cacheHits++

				// For history, we need to deserialize differently
				// This is a workaround - in production, you'd want a more sophisticated caching mechanism
				s.logger.WithField("cache_key", cacheKey).Debug("Email event history served from cache")
			}
			s.cacheMutex.RUnlock()
		}
	}

	// Get from base service (most history queries are not cached due to complexity)
	return s.baseService.GetEmailEventHistory(ctx, filters)
}

// IsHealthy checks if the cached service is healthy
func (s *EmailMetricsCacheService) IsHealthy(ctx context.Context) bool {
	// Check base service health
	if !s.baseService.IsHealthy(ctx) {
		return false
	}

	// Check cache health
	s.cacheMutex.RLock()
	cacheSize := len(s.cache)
	s.cacheMutex.RUnlock()

	// Cache is unhealthy if it's too large or cleanup hasn't run recently
	if cacheSize > s.maxCacheSize*2 {
		s.logger.WithField("cache_size", cacheSize).Warn("Cache size exceeds healthy limits")
		return false
	}

	if time.Since(s.lastCleanup) > time.Hour {
		s.logger.WithField("last_cleanup", s.lastCleanup).Warn("Cache cleanup hasn't run recently")
		return false
	}

	return true
}

// GetCacheStats returns cache performance statistics
func (s *EmailMetricsCacheService) GetCacheStats() map[string]interface{} {
	s.cacheMutex.RLock()
	defer s.cacheMutex.RUnlock()

	totalRequests := s.cacheHits + s.cacheMisses
	hitRate := float64(0)
	if totalRequests > 0 {
		hitRate = float64(s.cacheHits) / float64(totalRequests) * 100
	}

	return map[string]interface{}{
		"cache_hits":      s.cacheHits,
		"cache_misses":    s.cacheMisses,
		"cache_evictions": s.cacheEvictions,
		"hit_rate":        hitRate,
		"cache_size":      len(s.cache),
		"max_cache_size":  s.maxCacheSize,
		"last_cleanup":    s.lastCleanup,
	}
}

// ClearCache clears all cached entries
func (s *EmailMetricsCacheService) ClearCache() {
	s.cacheMutex.Lock()
	defer s.cacheMutex.Unlock()

	oldSize := len(s.cache)
	s.cache = make(map[string]*CacheEntry)
	s.cacheEvictions += int64(oldSize)

	s.logger.WithField("cleared_entries", oldSize).Info("Cache cleared")
}

// generateCacheKey creates a cache key from time range and filters
func (s *EmailMetricsCacheService) generateCacheKey(timeRange domain.TimeRange, emailType *domain.EmailEventType, status *domain.EmailEventStatus) string {
	key := fmt.Sprintf("metrics_%s_%s",
		timeRange.Start.Format("2006-01-02T15:04"),
		timeRange.End.Format("2006-01-02T15:04"))

	if emailType != nil {
		key += fmt.Sprintf("_type_%s", *emailType)
	}

	if status != nil {
		key += fmt.Sprintf("_status_%s", *status)
	}

	return key
}

// getFromCache retrieves an entry from cache if it exists and is not expired
func (s *EmailMetricsCacheService) getFromCache(key string) *domain.EmailMetrics {
	s.cacheMutex.RLock()
	defer s.cacheMutex.RUnlock()

	entry, exists := s.cache[key]
	if !exists {
		return nil
	}

	if time.Now().After(entry.ExpiresAt) {
		// Entry is expired, but don't remove it here to avoid write lock
		return nil
	}

	return entry.Data
}

// setInCache stores an entry in cache with TTL
func (s *EmailMetricsCacheService) setInCache(key string, data *domain.EmailMetrics, ttl time.Duration) {
	s.cacheMutex.Lock()
	defer s.cacheMutex.Unlock()

	// Check if cache is full and evict oldest entry
	if len(s.cache) >= s.maxCacheSize {
		s.evictOldestEntry()
	}

	s.cache[key] = &CacheEntry{
		Data:      data,
		ExpiresAt: time.Now().Add(ttl),
		CreatedAt: time.Now(),
	}
}

// evictOldestEntry removes the oldest cache entry (should be called with write lock held)
func (s *EmailMetricsCacheService) evictOldestEntry() {
	var oldestKey string
	var oldestTime time.Time

	for key, entry := range s.cache {
		if oldestKey == "" || entry.CreatedAt.Before(oldestTime) {
			oldestKey = key
			oldestTime = entry.CreatedAt
		}
	}

	if oldestKey != "" {
		delete(s.cache, oldestKey)
		s.cacheEvictions++
	}
}

// cleanupExpiredEntries runs in background to remove expired cache entries
func (s *EmailMetricsCacheService) cleanupExpiredEntries() {
	ticker := time.NewTicker(15 * time.Minute) // Cleanup every 15 minutes
	defer ticker.Stop()

	for range ticker.C {
		s.performCleanup()
	}
}

// performCleanup removes expired entries from cache
func (s *EmailMetricsCacheService) performCleanup() {
	s.cacheMutex.Lock()
	defer s.cacheMutex.Unlock()

	now := time.Now()
	expiredKeys := make([]string, 0)

	// Find expired keys
	for key, entry := range s.cache {
		if now.After(entry.ExpiresAt) {
			expiredKeys = append(expiredKeys, key)
		}
	}

	// Remove expired keys
	for _, key := range expiredKeys {
		delete(s.cache, key)
		s.cacheEvictions++
	}

	s.lastCleanup = now

	if len(expiredKeys) > 0 {
		s.logger.WithFields(logrus.Fields{
			"expired_entries": len(expiredKeys),
			"cache_size":      len(s.cache),
		}).Debug("Cleaned up expired cache entries")
	}
}

// WarmupCache pre-loads commonly requested metrics into cache
func (s *EmailMetricsCacheService) WarmupCache(ctx context.Context) error {
	s.logger.Info("Starting cache warmup")

	now := time.Now()
	commonTimeRanges := []domain.TimeRange{
		{Start: now.Add(-24 * time.Hour), End: now},      // Last 24 hours
		{Start: now.Add(-7 * 24 * time.Hour), End: now},  // Last 7 days
		{Start: now.Add(-30 * 24 * time.Hour), End: now}, // Last 30 days
	}

	for _, timeRange := range commonTimeRanges {
		_, err := s.GetEmailMetrics(ctx, timeRange)
		if err != nil {
			s.logger.WithError(err).WithFields(logrus.Fields{
				"start_time": timeRange.Start,
				"end_time":   timeRange.End,
			}).Error("Failed to warmup cache for time range")
			return fmt.Errorf("cache warmup failed: %w", err)
		}
	}

	s.logger.WithField("warmed_ranges", len(commonTimeRanges)).Info("Cache warmup completed")
	return nil
}
