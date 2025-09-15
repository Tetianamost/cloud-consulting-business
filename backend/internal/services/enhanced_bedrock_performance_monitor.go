package services

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/cloud-consulting/backend/internal/interfaces"
	"github.com/sirupsen/logrus"
)

// EnhancedBedrockPerformanceMonitor provides comprehensive performance monitoring and alerting
type EnhancedBedrockPerformanceMonitor struct {
	logger *logrus.Logger
	mu     sync.RWMutex

	// Performance metrics
	requestMetrics       *interfaces.SystemRequestMetrics
	responseTimeMetrics  *interfaces.SystemResponseTimeMetrics
	cacheMetrics         *interfaces.SystemCacheMetrics
	loadBalancingMetrics *interfaces.LoadBalancingMetrics
	systemMetrics        *interfaces.SystemMetrics

	// Alerting configuration
	alertThresholds *interfaces.AlertThresholds
	alertHandlers   map[string]interfaces.AlertHandler

	// Monitoring configuration
	monitoringInterval time.Duration
	metricsRetention   time.Duration
	alertingEnabled    bool
}

// NewEnhancedBedrockPerformanceMonitor creates a new enhanced performance monitor
func NewEnhancedBedrockPerformanceMonitor(logger *logrus.Logger) *EnhancedBedrockPerformanceMonitor {
	monitor := &EnhancedBedrockPerformanceMonitor{
		logger:         logger,
		requestMetrics: &interfaces.SystemRequestMetrics{LastUpdated: time.Now()},
		responseTimeMetrics: &interfaces.SystemResponseTimeMetrics{
			ResponseTimeHistory: make([]time.Duration, 0, 1000),
			LastUpdated:         time.Now(),
		},
		cacheMetrics:       &interfaces.SystemCacheMetrics{LastUpdated: time.Now()},
		systemMetrics:      &interfaces.SystemMetrics{LastUpdated: time.Now()},
		alertHandlers:      make(map[string]interfaces.AlertHandler),
		monitoringInterval: 30 * time.Second,
		metricsRetention:   24 * time.Hour,
		alertingEnabled:    true,
	}

	// Set default alert thresholds
	monitor.alertThresholds = &interfaces.AlertThresholds{
		MaxResponseTime:       5 * time.Second,
		MinCacheHitRate:       0.7,
		MaxErrorRate:          0.05,
		MaxConcurrentRequests: 100,
		MaxCPUUsage:           80.0,
		MaxMemoryUsage:        85.0,
	}

	// Initialize default alert handlers
	monitor.initializeAlertHandlers()

	return monitor
}

// RecordRequest records a request metric
func (m *EnhancedBedrockPerformanceMonitor) RecordRequest(success bool, responseTime time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.requestMetrics.TotalRequests++

	if success {
		m.requestMetrics.SuccessfulRequests++
	} else {
		m.requestMetrics.FailedRequests++
	}

	// Update response time metrics
	m.updateResponseTimeMetrics(responseTime)

	// Update requests per second (simple moving average)
	now := time.Now()
	timeSinceLastUpdate := now.Sub(m.requestMetrics.LastUpdated)
	if timeSinceLastUpdate > 0 {
		m.requestMetrics.RequestsPerSecond = 1.0 / timeSinceLastUpdate.Seconds()
	}
	m.requestMetrics.LastUpdated = now

	// Check for alerts
	if m.alertingEnabled {
		m.checkAlerts()
	}
}

// RecordCacheMetrics records cache performance metrics
func (m *EnhancedBedrockPerformanceMonitor) RecordCacheMetrics(hits, misses, size, evictions int64, avgAge time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.cacheMetrics.CacheHits = hits
	m.cacheMetrics.CacheMisses = misses
	m.cacheMetrics.CacheSize = size
	m.cacheMetrics.CacheEvictions = evictions
	m.cacheMetrics.AverageCacheAge = avgAge

	// Calculate hit rate
	total := hits + misses
	if total > 0 {
		m.cacheMetrics.CacheHitRate = float64(hits) / float64(total)
	}

	m.cacheMetrics.LastUpdated = time.Now()
}

// RecordSystemMetrics records system performance metrics
func (m *EnhancedBedrockPerformanceMonitor) RecordSystemMetrics(cpu, memory float64, goroutines int, heapSize int64, gcPause time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.systemMetrics.CPUUsage = cpu
	m.systemMetrics.MemoryUsage = memory
	m.systemMetrics.GoroutineCount = goroutines
	m.systemMetrics.HeapSize = heapSize
	m.systemMetrics.GCPauseTime = gcPause
	m.systemMetrics.LastUpdated = time.Now()
}

// GetPerformanceReport returns a comprehensive performance report
func (m *EnhancedBedrockPerformanceMonitor) GetPerformanceReport() *interfaces.SystemPerformanceReport {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return &interfaces.SystemPerformanceReport{
		RequestMetrics:      *m.requestMetrics,
		ResponseTimeMetrics: *m.responseTimeMetrics,
		CacheMetrics:        *m.cacheMetrics,
		SystemMetrics:       *m.systemMetrics,
		AlertThresholds:     *m.alertThresholds,
		GeneratedAt:         time.Now(),
	}
}

// StartMonitoring starts the performance monitoring loop
func (m *EnhancedBedrockPerformanceMonitor) StartMonitoring(ctx context.Context) {
	ticker := time.NewTicker(m.monitoringInterval)
	defer ticker.Stop()

	m.logger.Info("Started enhanced Bedrock performance monitoring")

	for {
		select {
		case <-ctx.Done():
			m.logger.Info("Stopped enhanced Bedrock performance monitoring")
			return
		case <-ticker.C:
			m.collectSystemMetrics()
			m.performHealthCheck()
			m.cleanupOldMetrics()
		}
	}
}

// RegisterAlertHandler registers a custom alert handler
func (m *EnhancedBedrockPerformanceMonitor) RegisterAlertHandler(name string, handler interfaces.AlertHandler) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.alertHandlers[name] = handler
	m.logger.WithField("handler_name", name).Info("Registered performance alert handler")
}

// SetAlertThresholds updates alert thresholds
func (m *EnhancedBedrockPerformanceMonitor) SetAlertThresholds(thresholds *interfaces.AlertThresholds) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.alertThresholds = thresholds
	m.logger.Info("Updated performance alert thresholds")
}

// Private helper methods

func (m *EnhancedBedrockPerformanceMonitor) updateResponseTimeMetrics(responseTime time.Duration) {
	// Add to history
	m.responseTimeMetrics.ResponseTimeHistory = append(m.responseTimeMetrics.ResponseTimeHistory, responseTime)

	// Keep only recent history
	if len(m.responseTimeMetrics.ResponseTimeHistory) > 1000 {
		m.responseTimeMetrics.ResponseTimeHistory = m.responseTimeMetrics.ResponseTimeHistory[1:]
	}

	// Update min/max
	if m.responseTimeMetrics.MinResponseTime == 0 || responseTime < m.responseTimeMetrics.MinResponseTime {
		m.responseTimeMetrics.MinResponseTime = responseTime
	}
	if responseTime > m.responseTimeMetrics.MaxResponseTime {
		m.responseTimeMetrics.MaxResponseTime = responseTime
	}

	// Calculate percentiles from history
	m.calculatePercentiles()

	m.responseTimeMetrics.LastUpdated = time.Now()
}

func (m *EnhancedBedrockPerformanceMonitor) calculatePercentiles() {
	history := m.responseTimeMetrics.ResponseTimeHistory
	if len(history) == 0 {
		return
	}

	// Simple percentile calculation (in production, use proper sorting)
	total := time.Duration(0)
	for _, rt := range history {
		total += rt
	}

	m.responseTimeMetrics.AverageResponseTime = total / time.Duration(len(history))

	// For simplicity, using average as approximation for percentiles
	// In production, implement proper percentile calculation
	m.responseTimeMetrics.MedianResponseTime = m.responseTimeMetrics.AverageResponseTime
	m.responseTimeMetrics.P95ResponseTime = time.Duration(float64(m.responseTimeMetrics.AverageResponseTime) * 1.5)
	m.responseTimeMetrics.P99ResponseTime = time.Duration(float64(m.responseTimeMetrics.AverageResponseTime) * 2.0)
}

func (m *EnhancedBedrockPerformanceMonitor) checkAlerts() {
	// Check response time alert
	if m.responseTimeMetrics.AverageResponseTime > m.alertThresholds.MaxResponseTime {
		m.triggerAlert(&interfaces.PerformanceAlert{
			Type:      interfaces.PerformanceAlertTypeResponseTime,
			Severity:  interfaces.PerformanceAlertSeverityWarning,
			Message:   "Average response time exceeds threshold",
			Metric:    "average_response_time",
			Value:     m.responseTimeMetrics.AverageResponseTime,
			Threshold: m.alertThresholds.MaxResponseTime,
			Timestamp: time.Now(),
		})
	}

	// Check cache hit rate alert
	if m.cacheMetrics.CacheHitRate < m.alertThresholds.MinCacheHitRate && m.cacheMetrics.CacheHits+m.cacheMetrics.CacheMisses > 100 {
		m.triggerAlert(&interfaces.PerformanceAlert{
			Type:      interfaces.PerformanceAlertTypeCacheHitRate,
			Severity:  interfaces.PerformanceAlertSeverityWarning,
			Message:   "Cache hit rate below threshold",
			Metric:    "cache_hit_rate",
			Value:     m.cacheMetrics.CacheHitRate,
			Threshold: m.alertThresholds.MinCacheHitRate,
			Timestamp: time.Now(),
		})
	}

	// Check error rate alert
	errorRate := float64(m.requestMetrics.FailedRequests) / float64(m.requestMetrics.TotalRequests)
	if errorRate > m.alertThresholds.MaxErrorRate && m.requestMetrics.TotalRequests > 100 {
		m.triggerAlert(&interfaces.PerformanceAlert{
			Type:      interfaces.PerformanceAlertTypeErrorRate,
			Severity:  interfaces.PerformanceAlertSeverityCritical,
			Message:   "Error rate exceeds threshold",
			Metric:    "error_rate",
			Value:     errorRate,
			Threshold: m.alertThresholds.MaxErrorRate,
			Timestamp: time.Now(),
		})
	}

	// Check system resource alerts
	if m.systemMetrics.CPUUsage > m.alertThresholds.MaxCPUUsage {
		m.triggerAlert(&interfaces.PerformanceAlert{
			Type:      interfaces.PerformanceAlertTypeSystemResource,
			Severity:  interfaces.PerformanceAlertSeverityWarning,
			Message:   "CPU usage exceeds threshold",
			Metric:    "cpu_usage",
			Value:     m.systemMetrics.CPUUsage,
			Threshold: m.alertThresholds.MaxCPUUsage,
			Timestamp: time.Now(),
		})
	}

	if m.systemMetrics.MemoryUsage > m.alertThresholds.MaxMemoryUsage {
		m.triggerAlert(&interfaces.PerformanceAlert{
			Type:      interfaces.PerformanceAlertTypeSystemResource,
			Severity:  interfaces.PerformanceAlertSeverityWarning,
			Message:   "Memory usage exceeds threshold",
			Metric:    "memory_usage",
			Value:     m.systemMetrics.MemoryUsage,
			Threshold: m.alertThresholds.MaxMemoryUsage,
			Timestamp: time.Now(),
		})
	}
}

func (m *EnhancedBedrockPerformanceMonitor) triggerAlert(alert *interfaces.PerformanceAlert) {
	alert.ID = generateAlertID()

	// Log the alert
	m.logger.WithFields(logrus.Fields{
		"alert_id":  alert.ID,
		"type":      alert.Type,
		"severity":  alert.Severity,
		"metric":    alert.Metric,
		"value":     alert.Value,
		"threshold": alert.Threshold,
	}).Warn("Performance alert triggered")

	// Execute alert handlers
	for name, handler := range m.alertHandlers {
		go func(handlerName string, h interfaces.AlertHandler) {
			if err := h(context.Background(), alert); err != nil {
				m.logger.WithError(err).WithField("handler", handlerName).Error("Alert handler failed")
			}
		}(name, handler)
	}
}

func (m *EnhancedBedrockPerformanceMonitor) collectSystemMetrics() {
	// In production, collect actual system metrics
	// For now, using placeholder values
	m.RecordSystemMetrics(
		50.0,             // CPU usage
		60.0,             // Memory usage
		100,              // Goroutine count
		1024*1024*50,     // Heap size (50MB)
		time.Millisecond, // GC pause time
	)
}

func (m *EnhancedBedrockPerformanceMonitor) performHealthCheck() {
	// Perform comprehensive health check
	report := m.GetPerformanceReport()

	m.logger.WithFields(logrus.Fields{
		"total_requests":    report.RequestMetrics.TotalRequests,
		"avg_response_time": report.ResponseTimeMetrics.AverageResponseTime.String(),
		"cache_hit_rate":    report.CacheMetrics.CacheHitRate,
		"cpu_usage":         report.SystemMetrics.CPUUsage,
		"memory_usage":      report.SystemMetrics.MemoryUsage,
	}).Info("Performance health check completed")
}

func (m *EnhancedBedrockPerformanceMonitor) cleanupOldMetrics() {
	// Clean up old metrics data
	cutoff := time.Now().Add(-m.metricsRetention)

	// Clean up response time history
	if len(m.responseTimeMetrics.ResponseTimeHistory) > 100 {
		// Keep only recent 100 entries for memory efficiency
		m.responseTimeMetrics.ResponseTimeHistory = m.responseTimeMetrics.ResponseTimeHistory[len(m.responseTimeMetrics.ResponseTimeHistory)-100:]
	}

	m.logger.WithField("cutoff", cutoff).Debug("Cleaned up old performance metrics")
}

func (m *EnhancedBedrockPerformanceMonitor) initializeAlertHandlers() {
	// Default log alert handler
	m.alertHandlers["log"] = func(ctx context.Context, alert *interfaces.PerformanceAlert) error {
		m.logger.WithFields(logrus.Fields{
			"alert_id":  alert.ID,
			"type":      alert.Type,
			"severity":  alert.Severity,
			"message":   alert.Message,
			"metric":    alert.Metric,
			"value":     alert.Value,
			"threshold": alert.Threshold,
		}).Error("Performance alert")
		return nil
	}
}

func generateAlertID() string {
	return fmt.Sprintf("alert-%d", time.Now().UnixNano())
}
