package services

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/cloud-consulting/backend/internal/storage"
)

// CacheMonitor provides monitoring and metrics for the Redis cache
type CacheMonitor struct {
	cache   *storage.RedisCache
	logger  *logrus.Logger
	metrics *CacheMetrics
	mu      sync.RWMutex
}

// CacheMetrics holds cache performance metrics
type CacheMetrics struct {
	// Hit/Miss statistics
	SessionHits   int64 `json:"session_hits"`
	SessionMisses int64 `json:"session_misses"`
	MessageHits   int64 `json:"message_hits"`
	MessageMisses int64 `json:"message_misses"`

	// Operation counts
	SessionSets        int64 `json:"session_sets"`
	SessionGets        int64 `json:"session_gets"`
	SessionDeletes     int64 `json:"session_deletes"`
	MessageSets        int64 `json:"message_sets"`
	MessageGets        int64 `json:"message_gets"`
	MessageDeletes     int64 `json:"message_deletes"`
	CacheInvalidations int64 `json:"cache_invalidations"`

	// Performance metrics
	AverageResponseTime time.Duration `json:"average_response_time"`
	MaxResponseTime     time.Duration `json:"max_response_time"`
	MinResponseTime     time.Duration `json:"min_response_time"`
	TotalOperations     int64         `json:"total_operations"`

	// Error tracking
	ConnectionErrors    int64 `json:"connection_errors"`
	TimeoutErrors       int64 `json:"timeout_errors"`
	SerializationErrors int64 `json:"serialization_errors"`

	// Timestamps
	LastReset   time.Time `json:"last_reset"`
	LastUpdated time.Time `json:"last_updated"`
}

// NewCacheMonitor creates a new cache monitor
func NewCacheMonitor(cache *storage.RedisCache, logger *logrus.Logger) *CacheMonitor {
	return &CacheMonitor{
		cache:  cache,
		logger: logger,
		metrics: &CacheMetrics{
			LastReset:       time.Now(),
			LastUpdated:     time.Now(),
			MinResponseTime: time.Hour, // Initialize to high value
		},
	}
}

// RecordSessionHit records a cache hit for session data
func (m *CacheMonitor) RecordSessionHit(responseTime time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.metrics.SessionHits++
	m.metrics.SessionGets++
	m.updateResponseTimeMetrics(responseTime)
	m.metrics.LastUpdated = time.Now()
}

// RecordSessionMiss records a cache miss for session data
func (m *CacheMonitor) RecordSessionMiss(responseTime time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.metrics.SessionMisses++
	m.metrics.SessionGets++
	m.updateResponseTimeMetrics(responseTime)
	m.metrics.LastUpdated = time.Now()
}

// RecordMessageHit records a cache hit for message data
func (m *CacheMonitor) RecordMessageHit(responseTime time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.metrics.MessageHits++
	m.metrics.MessageGets++
	m.updateResponseTimeMetrics(responseTime)
	m.metrics.LastUpdated = time.Now()
}

// RecordMessageMiss records a cache miss for message data
func (m *CacheMonitor) RecordMessageMiss(responseTime time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.metrics.MessageMisses++
	m.metrics.MessageGets++
	m.updateResponseTimeMetrics(responseTime)
	m.metrics.LastUpdated = time.Now()
}

// RecordSessionSet records a session cache set operation
func (m *CacheMonitor) RecordSessionSet(responseTime time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.metrics.SessionSets++
	m.updateResponseTimeMetrics(responseTime)
	m.metrics.LastUpdated = time.Now()
}

// RecordMessageSet records a message cache set operation
func (m *CacheMonitor) RecordMessageSet(responseTime time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.metrics.MessageSets++
	m.updateResponseTimeMetrics(responseTime)
	m.metrics.LastUpdated = time.Now()
}

// RecordSessionDelete records a session cache delete operation
func (m *CacheMonitor) RecordSessionDelete(responseTime time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.metrics.SessionDeletes++
	m.updateResponseTimeMetrics(responseTime)
	m.metrics.LastUpdated = time.Now()
}

// RecordMessageDelete records a message cache delete operation
func (m *CacheMonitor) RecordMessageDelete(responseTime time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.metrics.MessageDeletes++
	m.updateResponseTimeMetrics(responseTime)
	m.metrics.LastUpdated = time.Now()
}

// RecordCacheInvalidation records a cache invalidation operation
func (m *CacheMonitor) RecordCacheInvalidation(responseTime time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.metrics.CacheInvalidations++
	m.updateResponseTimeMetrics(responseTime)
	m.metrics.LastUpdated = time.Now()
}

// RecordConnectionError records a cache connection error
func (m *CacheMonitor) RecordConnectionError() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.metrics.ConnectionErrors++
	m.metrics.LastUpdated = time.Now()
}

// RecordTimeoutError records a cache timeout error
func (m *CacheMonitor) RecordTimeoutError() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.metrics.TimeoutErrors++
	m.metrics.LastUpdated = time.Now()
}

// RecordSerializationError records a serialization error
func (m *CacheMonitor) RecordSerializationError() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.metrics.SerializationErrors++
	m.metrics.LastUpdated = time.Now()
}

// updateResponseTimeMetrics updates response time statistics (must be called with lock held)
func (m *CacheMonitor) updateResponseTimeMetrics(responseTime time.Duration) {
	m.metrics.TotalOperations++

	// Update average response time
	if m.metrics.TotalOperations == 1 {
		m.metrics.AverageResponseTime = responseTime
	} else {
		// Calculate running average
		totalTime := time.Duration(int64(m.metrics.AverageResponseTime) * (m.metrics.TotalOperations - 1))
		m.metrics.AverageResponseTime = (totalTime + responseTime) / time.Duration(m.metrics.TotalOperations)
	}

	// Update max response time
	if responseTime > m.metrics.MaxResponseTime {
		m.metrics.MaxResponseTime = responseTime
	}

	// Update min response time
	if responseTime < m.metrics.MinResponseTime {
		m.metrics.MinResponseTime = responseTime
	}
}

// GetMetrics returns a copy of the current metrics
func (m *CacheMonitor) GetMetrics() *CacheMetrics {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Create a copy to avoid race conditions
	metricsCopy := *m.metrics
	return &metricsCopy
}

// GetCacheHitRatio returns the overall cache hit ratio
func (m *CacheMonitor) GetCacheHitRatio() float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()

	totalHits := m.metrics.SessionHits + m.metrics.MessageHits
	totalRequests := m.metrics.SessionGets + m.metrics.MessageGets

	if totalRequests == 0 {
		return 0.0
	}

	return float64(totalHits) / float64(totalRequests)
}

// GetSessionHitRatio returns the session cache hit ratio
func (m *CacheMonitor) GetSessionHitRatio() float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.metrics.SessionGets == 0 {
		return 0.0
	}

	return float64(m.metrics.SessionHits) / float64(m.metrics.SessionGets)
}

// GetMessageHitRatio returns the message cache hit ratio
func (m *CacheMonitor) GetMessageHitRatio() float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.metrics.MessageGets == 0 {
		return 0.0
	}

	return float64(m.metrics.MessageHits) / float64(m.metrics.MessageGets)
}

// GetErrorRate returns the overall error rate
func (m *CacheMonitor) GetErrorRate() float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()

	totalErrors := m.metrics.ConnectionErrors + m.metrics.TimeoutErrors + m.metrics.SerializationErrors
	totalOperations := m.metrics.TotalOperations

	if totalOperations == 0 {
		return 0.0
	}

	return float64(totalErrors) / float64(totalOperations)
}

// ResetMetrics resets all metrics to zero
func (m *CacheMonitor) ResetMetrics() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.metrics = &CacheMetrics{
		LastReset:       time.Now(),
		LastUpdated:     time.Now(),
		MinResponseTime: time.Hour, // Initialize to high value
	}

	m.logger.Info("Cache metrics reset")
}

// LogMetrics logs current cache metrics
func (m *CacheMonitor) LogMetrics() {
	metrics := m.GetMetrics()

	m.logger.WithFields(logrus.Fields{
		"session_hit_ratio":    m.GetSessionHitRatio(),
		"message_hit_ratio":    m.GetMessageHitRatio(),
		"overall_hit_ratio":    m.GetCacheHitRatio(),
		"error_rate":           m.GetErrorRate(),
		"avg_response_time_ms": metrics.AverageResponseTime.Milliseconds(),
		"max_response_time_ms": metrics.MaxResponseTime.Milliseconds(),
		"min_response_time_ms": metrics.MinResponseTime.Milliseconds(),
		"total_operations":     metrics.TotalOperations,
		"session_hits":         metrics.SessionHits,
		"session_misses":       metrics.SessionMisses,
		"message_hits":         metrics.MessageHits,
		"message_misses":       metrics.MessageMisses,
		"connection_errors":    metrics.ConnectionErrors,
		"timeout_errors":       metrics.TimeoutErrors,
		"serialization_errors": metrics.SerializationErrors,
	}).Info("Cache performance metrics")
}

// GetHealthStatus returns the health status of the cache
func (m *CacheMonitor) GetHealthStatus(ctx context.Context) *CacheHealthStatus {
	startTime := time.Now()
	isHealthy := m.cache.IsHealthy(ctx)
	responseTime := time.Since(startTime)

	errorRate := m.GetErrorRate()
	hitRatio := m.GetCacheHitRatio()

	status := &CacheHealthStatus{
		IsHealthy:    isHealthy,
		ResponseTime: responseTime,
		ErrorRate:    errorRate,
		HitRatio:     hitRatio,
		LastChecked:  time.Now(),
	}

	// Determine overall health based on multiple factors
	if !isHealthy {
		status.Status = "UNHEALTHY"
		status.Message = "Cache connection failed"
	} else if errorRate > 0.1 { // More than 10% error rate
		status.Status = "DEGRADED"
		status.Message = fmt.Sprintf("High error rate: %.2f%%", errorRate*100)
	} else if responseTime > 100*time.Millisecond {
		status.Status = "DEGRADED"
		status.Message = fmt.Sprintf("High response time: %v", responseTime)
	} else {
		status.Status = "HEALTHY"
		status.Message = "Cache operating normally"
	}

	return status
}

// StartPeriodicLogging starts periodic logging of cache metrics
func (m *CacheMonitor) StartPeriodicLogging(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	m.logger.WithField("interval", interval).Info("Started periodic cache metrics logging")

	for {
		select {
		case <-ctx.Done():
			m.logger.Info("Stopped periodic cache metrics logging")
			return
		case <-ticker.C:
			m.LogMetrics()
		}
	}
}

// CacheHealthStatus represents the health status of the cache
type CacheHealthStatus struct {
	Status       string        `json:"status"`
	Message      string        `json:"message"`
	IsHealthy    bool          `json:"is_healthy"`
	ResponseTime time.Duration `json:"response_time"`
	ErrorRate    float64       `json:"error_rate"`
	HitRatio     float64       `json:"hit_ratio"`
	LastChecked  time.Time     `json:"last_checked"`
}

// CachePerformanceReport generates a comprehensive performance report
type CachePerformanceReport struct {
	Metrics         *CacheMetrics      `json:"metrics"`
	HealthStatus    *CacheHealthStatus `json:"health_status"`
	Recommendations []string           `json:"recommendations"`
	GeneratedAt     time.Time          `json:"generated_at"`
}

// GeneratePerformanceReport generates a comprehensive cache performance report
func (m *CacheMonitor) GeneratePerformanceReport(ctx context.Context) *CachePerformanceReport {
	metrics := m.GetMetrics()
	healthStatus := m.GetHealthStatus(ctx)

	report := &CachePerformanceReport{
		Metrics:      metrics,
		HealthStatus: healthStatus,
		GeneratedAt:  time.Now(),
	}

	// Generate recommendations based on metrics
	var recommendations []string

	hitRatio := m.GetCacheHitRatio()
	if hitRatio < 0.7 {
		recommendations = append(recommendations,
			fmt.Sprintf("Cache hit ratio is low (%.2f%%). Consider increasing cache TTL or reviewing caching strategy.", hitRatio*100))
	}

	errorRate := m.GetErrorRate()
	if errorRate > 0.05 {
		recommendations = append(recommendations,
			fmt.Sprintf("Error rate is high (%.2f%%). Check Redis connection and configuration.", errorRate*100))
	}

	if metrics.AverageResponseTime > 50*time.Millisecond {
		recommendations = append(recommendations,
			fmt.Sprintf("Average response time is high (%v). Consider optimizing Redis configuration or network.", metrics.AverageResponseTime))
	}

	if metrics.ConnectionErrors > 0 {
		recommendations = append(recommendations,
			"Connection errors detected. Check Redis server availability and network connectivity.")
	}

	if metrics.TimeoutErrors > 0 {
		recommendations = append(recommendations,
			"Timeout errors detected. Consider increasing Redis timeout settings.")
	}

	if len(recommendations) == 0 {
		recommendations = append(recommendations, "Cache is performing well. No immediate optimizations needed.")
	}

	report.Recommendations = recommendations
	return report
}
