package services

import (
	"context"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// ChatPerformanceMonitor tracks performance metrics for chat operations
type ChatPerformanceMonitor struct {
	logger *logrus.Logger
	mu     sync.RWMutex

	// Metrics
	messagesSent      int64
	messagesReceived  int64
	totalResponseTime time.Duration
	responseCount     int64
	connectionCount   int64
	activeConnections int64
	cacheHits         int64
	cacheMisses       int64

	// Performance thresholds
	maxResponseTime time.Duration
	maxConnections  int64

	// Monitoring intervals
	lastReset          time.Time
	monitoringInterval time.Duration
}

// PerformanceMetrics represents current performance metrics
type PerformanceMetrics struct {
	MessagesSent        int64         `json:"messages_sent"`
	MessagesReceived    int64         `json:"messages_received"`
	AverageResponseTime time.Duration `json:"average_response_time"`
	ActiveConnections   int64         `json:"active_connections"`
	TotalConnections    int64         `json:"total_connections"`
	CacheHitRate        float64       `json:"cache_hit_rate"`
	Timestamp           time.Time     `json:"timestamp"`
}

// NewChatPerformanceMonitor creates a new performance monitor
func NewChatPerformanceMonitor(logger *logrus.Logger) *ChatPerformanceMonitor {
	return &ChatPerformanceMonitor{
		logger:             logger,
		maxResponseTime:    5 * time.Second,
		maxConnections:     1000,
		lastReset:          time.Now(),
		monitoringInterval: 1 * time.Minute,
	}
}

// RecordMessageSent records a sent message
func (m *ChatPerformanceMonitor) RecordMessageSent() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.messagesSent++
}

// RecordMessageReceived records a received message
func (m *ChatPerformanceMonitor) RecordMessageReceived() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.messagesReceived++
}

// RecordResponseTime records the time taken to generate a response
func (m *ChatPerformanceMonitor) RecordResponseTime(duration time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.totalResponseTime += duration
	m.responseCount++

	// Log slow responses
	if duration > m.maxResponseTime {
		m.logger.WithFields(logrus.Fields{
			"response_time": duration.String(),
			"threshold":     m.maxResponseTime.String(),
		}).Warn("Slow response detected")
	}
}

// RecordConnectionOpened records a new connection
func (m *ChatPerformanceMonitor) RecordConnectionOpened() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.connectionCount++
	m.activeConnections++

	// Log high connection count
	if m.activeConnections > m.maxConnections {
		m.logger.WithFields(logrus.Fields{
			"active_connections": m.activeConnections,
			"max_connections":    m.maxConnections,
		}).Warn("High connection count detected")
	}
}

// RecordConnectionClosed records a closed connection
func (m *ChatPerformanceMonitor) RecordConnectionClosed() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.activeConnections > 0 {
		m.activeConnections--
	}
}

// RecordCacheHit records a cache hit
func (m *ChatPerformanceMonitor) RecordCacheHit() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.cacheHits++
}

// RecordCacheMiss records a cache miss
func (m *ChatPerformanceMonitor) RecordCacheMiss() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.cacheMisses++
}

// GetMetrics returns current performance metrics
func (m *ChatPerformanceMonitor) GetMetrics() PerformanceMetrics {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var avgResponseTime time.Duration
	if m.responseCount > 0 {
		avgResponseTime = m.totalResponseTime / time.Duration(m.responseCount)
	}

	var cacheHitRate float64
	totalCacheRequests := m.cacheHits + m.cacheMisses
	if totalCacheRequests > 0 {
		cacheHitRate = float64(m.cacheHits) / float64(totalCacheRequests)
	}

	return PerformanceMetrics{
		MessagesSent:        m.messagesSent,
		MessagesReceived:    m.messagesReceived,
		AverageResponseTime: avgResponseTime,
		ActiveConnections:   m.activeConnections,
		TotalConnections:    m.connectionCount,
		CacheHitRate:        cacheHitRate,
		Timestamp:           time.Now(),
	}
}

// ResetMetrics resets all metrics (typically called periodically)
func (m *ChatPerformanceMonitor) ResetMetrics() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.messagesSent = 0
	m.messagesReceived = 0
	m.totalResponseTime = 0
	m.responseCount = 0
	m.connectionCount = 0
	m.cacheHits = 0
	m.cacheMisses = 0
	m.lastReset = time.Now()

	m.logger.Info("Performance metrics reset")
}

// StartMonitoring starts periodic monitoring and logging
func (m *ChatPerformanceMonitor) StartMonitoring(ctx context.Context) {
	ticker := time.NewTicker(m.monitoringInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			m.logMetrics()
		}
	}
}

// logMetrics logs current performance metrics
func (m *ChatPerformanceMonitor) logMetrics() {
	metrics := m.GetMetrics()

	m.logger.WithFields(logrus.Fields{
		"messages_sent":         metrics.MessagesSent,
		"messages_received":     metrics.MessagesReceived,
		"average_response_time": metrics.AverageResponseTime.String(),
		"active_connections":    metrics.ActiveConnections,
		"total_connections":     metrics.TotalConnections,
		"cache_hit_rate":        metrics.CacheHitRate,
	}).Info("Chat performance metrics")
}

// GetHealthStatus returns the health status based on performance metrics
func (m *ChatPerformanceMonitor) GetHealthStatus() map[string]interface{} {
	metrics := m.GetMetrics()

	status := map[string]interface{}{
		"healthy": true,
		"metrics": metrics,
		"issues":  []string{},
	}

	issues := []string{}

	// Check response time
	if metrics.AverageResponseTime > m.maxResponseTime {
		issues = append(issues, "High average response time")
		status["healthy"] = false
	}

	// Check connection count
	if metrics.ActiveConnections > m.maxConnections {
		issues = append(issues, "High active connection count")
		status["healthy"] = false
	}

	// Check cache hit rate
	if metrics.CacheHitRate < 0.5 && (m.cacheHits+m.cacheMisses) > 100 {
		issues = append(issues, "Low cache hit rate")
	}

	status["issues"] = issues
	return status
}

// SetThresholds updates performance thresholds
func (m *ChatPerformanceMonitor) SetThresholds(maxResponseTime time.Duration, maxConnections int64) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.maxResponseTime = maxResponseTime
	m.maxConnections = maxConnections

	m.logger.WithFields(logrus.Fields{
		"max_response_time": maxResponseTime.String(),
		"max_connections":   maxConnections,
	}).Info("Performance thresholds updated")
}
