package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/cloud-consulting/backend/internal/services"
)

// ChatMetricsHandler handles chat metrics HTTP requests
type ChatMetricsHandler struct {
	metricsCollector   *services.ChatMetricsCollector
	performanceMonitor *services.ChatPerformanceMonitor
	cacheMonitor       *services.CacheMonitor
	logger             *logrus.Logger
}

// NewChatMetricsHandler creates a new chat metrics handler
func NewChatMetricsHandler(
	metricsCollector *services.ChatMetricsCollector,
	performanceMonitor *services.ChatPerformanceMonitor,
	cacheMonitor *services.CacheMonitor,
	logger *logrus.Logger,
) *ChatMetricsHandler {
	return &ChatMetricsHandler{
		metricsCollector:   metricsCollector,
		performanceMonitor: performanceMonitor,
		cacheMonitor:       cacheMonitor,
		logger:             logger,
	}
}

// GetChatMetrics handles GET /api/v1/admin/chat/metrics
func (h *ChatMetricsHandler) GetChatMetrics(c *gin.Context) {
	metrics := h.metricsCollector.GetMetrics()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    metrics,
	})
}

// GetConnectionMetrics handles GET /api/v1/admin/chat/metrics/connections
func (h *ChatMetricsHandler) GetConnectionMetrics(c *gin.Context) {
	metrics := h.metricsCollector.GetMetrics()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    metrics.ConnectionMetrics,
	})
}

// GetMessageMetrics handles GET /api/v1/admin/chat/metrics/messages
func (h *ChatMetricsHandler) GetMessageMetrics(c *gin.Context) {
	metrics := h.metricsCollector.GetMetrics()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    metrics.MessageMetrics,
	})
}

// GetAIMetrics handles GET /api/v1/admin/chat/metrics/ai
func (h *ChatMetricsHandler) GetAIMetrics(c *gin.Context) {
	metrics := h.metricsCollector.GetMetrics()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    metrics.AIMetrics,
	})
}

// GetUserMetrics handles GET /api/v1/admin/chat/metrics/users
func (h *ChatMetricsHandler) GetUserMetrics(c *gin.Context) {
	metrics := h.metricsCollector.GetMetrics()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    metrics.UserMetrics,
	})
}

// GetErrorMetrics handles GET /api/v1/admin/chat/metrics/errors
func (h *ChatMetricsHandler) GetErrorMetrics(c *gin.Context) {
	metrics := h.metricsCollector.GetMetrics()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    metrics.ErrorMetrics,
	})
}

// GetPerformanceMetrics handles GET /api/v1/admin/chat/metrics/performance
func (h *ChatMetricsHandler) GetPerformanceMetrics(c *gin.Context) {
	metrics := h.metricsCollector.GetMetrics()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    metrics.PerformanceMetrics,
	})
}

// GetCacheMetrics handles GET /api/v1/admin/chat/metrics/cache
func (h *ChatMetricsHandler) GetCacheMetrics(c *gin.Context) {
	metrics := h.metricsCollector.GetMetrics()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    metrics.CacheMetrics,
	})
}

// GetPrometheusMetrics handles GET /metrics (Prometheus endpoint)
func (h *ChatMetricsHandler) GetPrometheusMetrics(c *gin.Context) {
	prometheusMetrics := h.metricsCollector.GetPrometheusMetrics()

	c.Header("Content-Type", "text/plain; version=0.0.4; charset=utf-8")
	c.String(http.StatusOK, prometheusMetrics)
}

// GetHealthStatus handles GET /api/v1/admin/chat/health
func (h *ChatMetricsHandler) GetHealthStatus(c *gin.Context) {
	// Get performance health status
	perfHealth := h.performanceMonitor.GetHealthStatus()

	// Get cache health status
	cacheHealth := h.cacheMonitor.GetHealthStatus(c.Request.Context())

	// Get overall metrics
	metrics := h.metricsCollector.GetMetrics()

	// Determine overall health
	overallHealthy := true
	issues := []string{}

	// Check performance health
	if !perfHealth["healthy"].(bool) {
		overallHealthy = false
		if perfIssues, ok := perfHealth["issues"].([]string); ok {
			issues = append(issues, perfIssues...)
		}
	}

	// Check cache health
	if !cacheHealth.IsHealthy {
		overallHealthy = false
		issues = append(issues, "Cache: "+cacheHealth.Message)
	}

	// Check error rates
	if metrics.ErrorMetrics.ErrorRate > 0.1 { // More than 10% error rate
		overallHealthy = false
		issues = append(issues, "High error rate detected")
	}

	// Check connection health
	if metrics.ConnectionMetrics.ConnectionSuccessRate < 0.9 { // Less than 90% success rate
		overallHealthy = false
		issues = append(issues, "Low connection success rate")
	}

	status := "healthy"
	if !overallHealthy {
		status = "unhealthy"
	}

	healthStatus := gin.H{
		"status":    status,
		"healthy":   overallHealthy,
		"timestamp": time.Now(),
		"issues":    issues,
		"components": gin.H{
			"performance": perfHealth,
			"cache":       cacheHealth,
			"metrics": gin.H{
				"error_rate":              metrics.ErrorMetrics.ErrorRate,
				"connection_success_rate": metrics.ConnectionMetrics.ConnectionSuccessRate,
				"ai_success_rate":         metrics.AIMetrics.SuccessRate,
				"message_success_rate":    metrics.MessageMetrics.MessageSuccessRate,
			},
		},
	}

	statusCode := http.StatusOK
	if !overallHealthy {
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, gin.H{
		"success": overallHealthy,
		"data":    healthStatus,
	})
}

// ResetMetrics handles POST /api/v1/admin/chat/metrics/reset
func (h *ChatMetricsHandler) ResetMetrics(c *gin.Context) {
	// Reset all metrics
	h.metricsCollector.ResetMetrics()
	h.performanceMonitor.ResetMetrics()
	h.cacheMonitor.ResetMetrics()

	h.logger.Info("Chat metrics reset via admin API")

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"message":   "All chat metrics have been reset",
		"timestamp": time.Now(),
	})
}

// GetMetricsHistory handles GET /api/v1/admin/chat/metrics/history
func (h *ChatMetricsHandler) GetMetricsHistory(c *gin.Context) {
	// Parse query parameters
	hoursStr := c.DefaultQuery("hours", "24")
	hours, err := strconv.Atoi(hoursStr)
	if err != nil || hours < 1 || hours > 168 { // Max 1 week
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid hours parameter. Must be between 1 and 168",
		})
		return
	}

	intervalStr := c.DefaultQuery("interval", "1h")
	_, err = time.ParseDuration(intervalStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid interval parameter. Use format like '1h', '30m', '15m'",
		})
		return
	}

	// For now, return current metrics as historical data
	// In a real implementation, this would query a time-series database
	currentMetrics := h.metricsCollector.GetMetrics()

	// Generate mock historical data points
	now := time.Now()
	dataPoints := []gin.H{}

	for i := hours; i >= 0; i-- {
		timestamp := now.Add(-time.Duration(i) * time.Hour)

		// Create a data point with slight variations from current metrics
		dataPoint := gin.H{
			"timestamp": timestamp,
			"metrics": gin.H{
				"active_connections":   currentMetrics.ConnectionMetrics.ActiveConnections,
				"messages_per_hour":    currentMetrics.MessageMetrics.MessagesSent + currentMetrics.MessageMetrics.MessagesReceived,
				"ai_requests_per_hour": currentMetrics.AIMetrics.RequestsTotal,
				"error_rate":           currentMetrics.ErrorMetrics.ErrorRate,
				"cache_hit_rate":       h.cacheMonitor.GetCacheHitRatio(),
				"response_time_ms":     currentMetrics.PerformanceMetrics.ResponseTimeP50.Milliseconds(),
			},
		}
		dataPoints = append(dataPoints, dataPoint)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"period":      hours,
			"interval":    intervalStr,
			"data_points": dataPoints,
		},
	})
}

// GetMetricsSummary handles GET /api/v1/admin/chat/metrics/summary
func (h *ChatMetricsHandler) GetMetricsSummary(c *gin.Context) {
	metrics := h.metricsCollector.GetMetrics()

	summary := gin.H{
		"overview": gin.H{
			"active_connections": metrics.ConnectionMetrics.ActiveConnections,
			"total_connections":  metrics.ConnectionMetrics.TotalConnections,
			"active_users":       metrics.UserMetrics.ActiveUsers,
			"messages_today":     metrics.MessageMetrics.MessagesSent + metrics.MessageMetrics.MessagesReceived,
			"ai_requests_today":  metrics.AIMetrics.RequestsTotal,
			"estimated_ai_cost":  metrics.AIMetrics.EstimatedCost,
		},
		"health": gin.H{
			"connection_success_rate": metrics.ConnectionMetrics.ConnectionSuccessRate,
			"message_success_rate":    metrics.MessageMetrics.MessageSuccessRate,
			"ai_success_rate":         metrics.AIMetrics.SuccessRate,
			"cache_hit_rate":          h.cacheMonitor.GetCacheHitRatio(),
			"error_rate":              metrics.ErrorMetrics.ErrorRate,
		},
		"performance": gin.H{
			"avg_response_time_ms": metrics.PerformanceMetrics.ResponseTimeP50.Milliseconds(),
			"throughput_mps":       metrics.PerformanceMetrics.ThroughputMPS,
			"memory_usage_mb":      metrics.PerformanceMetrics.MemoryUsageMB,
			"cpu_usage_percent":    metrics.PerformanceMetrics.CPUUsagePercent,
		},
		"usage": gin.H{
			"messages_per_session":     metrics.UserMetrics.MessagesPerSession,
			"avg_session_duration_min": metrics.UserMetrics.AverageSessionDuration.Minutes(),
			"quick_actions_used":       metrics.UserMetrics.QuickActionsUsed,
			"feature_usage":            metrics.UserMetrics.FeatureUsage,
		},
		"timestamp": metrics.Timestamp,
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    summary,
	})
}

// GetAlertsStatus handles GET /api/v1/admin/chat/alerts
func (h *ChatMetricsHandler) GetAlertsStatus(c *gin.Context) {
	metrics := h.metricsCollector.GetMetrics()

	alerts := []gin.H{}

	// Check for various alert conditions

	// High error rate alert
	if metrics.ErrorMetrics.ErrorRate > 0.05 { // More than 5%
		alerts = append(alerts, gin.H{
			"type":      "error_rate",
			"severity":  "warning",
			"message":   "High error rate detected",
			"value":     metrics.ErrorMetrics.ErrorRate,
			"threshold": 0.05,
			"timestamp": time.Now(),
		})
	}

	// Low connection success rate alert
	if metrics.ConnectionMetrics.ConnectionSuccessRate < 0.95 { // Less than 95%
		alerts = append(alerts, gin.H{
			"type":      "connection_success",
			"severity":  "warning",
			"message":   "Low connection success rate",
			"value":     metrics.ConnectionMetrics.ConnectionSuccessRate,
			"threshold": 0.95,
			"timestamp": time.Now(),
		})
	}

	// High response time alert
	if metrics.PerformanceMetrics.ResponseTimeP95 > 2*time.Second {
		alerts = append(alerts, gin.H{
			"type":      "response_time",
			"severity":  "warning",
			"message":   "High response time detected",
			"value":     metrics.PerformanceMetrics.ResponseTimeP95.Milliseconds(),
			"threshold": 2000,
			"timestamp": time.Now(),
		})
	}

	// Low cache hit rate alert
	cacheHitRate := h.cacheMonitor.GetCacheHitRatio()
	if cacheHitRate < 0.7 { // Less than 70%
		alerts = append(alerts, gin.H{
			"type":      "cache_hit_rate",
			"severity":  "info",
			"message":   "Low cache hit rate",
			"value":     cacheHitRate,
			"threshold": 0.7,
			"timestamp": time.Now(),
		})
	}

	// High AI cost alert
	if metrics.AIMetrics.EstimatedCost > 100.0 { // More than $100
		alerts = append(alerts, gin.H{
			"type":      "ai_cost",
			"severity":  "warning",
			"message":   "High AI usage cost",
			"value":     metrics.AIMetrics.EstimatedCost,
			"threshold": 100.0,
			"timestamp": time.Now(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"alerts":      alerts,
			"alert_count": len(alerts),
			"timestamp":   time.Now(),
		},
	})
}
