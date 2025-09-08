package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/cloud-consulting/backend/internal/services"
)

// EmailSystemHealthHandler provides comprehensive health check endpoints for the email event system
type EmailSystemHealthHandler struct {
	emailMonitoringService *services.EmailMonitoringService
	logger                 *logrus.Logger
}

// NewEmailSystemHealthHandler creates a new email system health handler
func NewEmailSystemHealthHandler(
	emailMonitoringService *services.EmailMonitoringService,
	logger *logrus.Logger,
) *EmailSystemHealthHandler {
	return &EmailSystemHealthHandler{
		emailMonitoringService: emailMonitoringService,
		logger:                 logger,
	}
}

// GetEmailSystemHealthCheck handles GET /health/email-system
// This is the main health check endpoint for the email event tracking system
func (h *EmailSystemHealthHandler) GetEmailSystemHealthCheck(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 15*time.Second)
	defer cancel()

	h.logger.WithFields(logrus.Fields{
		"component": "email_system_health_handler",
		"endpoint":  "/health/email-system",
		"action":    "health_check_requested",
	}).Info("Email system health check requested")

	// Perform comprehensive health check
	healthStatus := h.emailMonitoringService.PerformHealthCheck(ctx)

	// Determine HTTP status code based on health
	statusCode := http.StatusOK
	if healthStatus.OverallStatus == "unhealthy" {
		statusCode = http.StatusServiceUnavailable
	} else if healthStatus.OverallStatus == "degraded" {
		statusCode = http.StatusPartialContent // 206 for degraded service
	}

	// Get current metrics for additional context
	metrics := h.emailMonitoringService.GetSystemMetrics()

	response := gin.H{
		"status":    healthStatus.OverallStatus,
		"timestamp": healthStatus.LastChecked.Format(time.RFC3339),
		"health": gin.H{
			"overall_status":           healthStatus.OverallStatus,
			"recorder_healthy":         healthStatus.RecorderHealthy,
			"metrics_service_healthy":  healthStatus.MetricsServiceHealthy,
			"database_healthy":         healthStatus.DatabaseHealthy,
			"consecutive_failures":     healthStatus.ConsecutiveFailures,
			"health_check_duration_ms": healthStatus.HealthCheckDuration.Nanoseconds() / 1e6,
			"issues":                   healthStatus.Issues,
		},
		"metrics": gin.H{
			"system_uptime_seconds":         metrics.SystemUptime.Seconds(),
			"health_check_interval_seconds": metrics.HealthCheckInterval.Seconds(),
			"alerts_triggered":              metrics.AlertsTriggered,
			"alert_suppression_count":       metrics.AlertSuppressionCount,
			"average_processing_time_ms":    metrics.AverageEmailProcessingTime,
			"email_throughput_per_hour":     metrics.EmailThroughputPerHour,
		},
		"components": gin.H{
			"email_event_recorder":  h.getRecorderHealthInfo(metrics),
			"email_metrics_service": h.getMetricsServiceHealthInfo(metrics),
		},
	}

	h.logger.WithFields(logrus.Fields{
		"component":               "email_system_health_handler",
		"overall_status":          healthStatus.OverallStatus,
		"recorder_healthy":        healthStatus.RecorderHealthy,
		"metrics_service_healthy": healthStatus.MetricsServiceHealthy,
		"consecutive_failures":    healthStatus.ConsecutiveFailures,
		"issues_count":            len(healthStatus.Issues),
		"response_code":           statusCode,
		"action":                  "health_check_completed",
	}).Info("Email system health check completed")

	c.JSON(statusCode, response)
}

// GetEmailSystemLiveness handles GET /health/email-system/liveness
// Simple liveness probe for Kubernetes
func (h *EmailSystemHealthHandler) GetEmailSystemLiveness(c *gin.Context) {
	// Basic liveness check - just verify the monitoring service is running
	isAlive := h.emailMonitoringService.IsHealthy()

	statusCode := http.StatusOK
	if !isAlive {
		statusCode = http.StatusServiceUnavailable
	}

	response := gin.H{
		"status":    map[bool]string{true: "alive", false: "dead"}[isAlive],
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"service":   "email-event-tracking-system",
	}

	h.logger.WithFields(logrus.Fields{
		"component":     "email_system_health_handler",
		"endpoint":      "/health/email-system/liveness",
		"is_alive":      isAlive,
		"response_code": statusCode,
		"action":        "liveness_check",
	}).Debug("Email system liveness check")

	c.JSON(statusCode, response)
}

// GetEmailSystemReadiness handles GET /health/email-system/readiness
// Readiness probe for Kubernetes - checks if system is ready to handle requests
func (h *EmailSystemHealthHandler) GetEmailSystemReadiness(c *gin.Context) {
	// Get current health status
	healthStatus := h.emailMonitoringService.GetHealthStatus()

	// System is ready if it's healthy or degraded (but not unhealthy)
	isReady := healthStatus.OverallStatus != "unhealthy"

	statusCode := http.StatusOK
	if !isReady {
		statusCode = http.StatusServiceUnavailable
	}

	response := gin.H{
		"status":    map[bool]string{true: "ready", false: "not_ready"}[isReady],
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"service":   "email-event-tracking-system",
		"details": gin.H{
			"overall_status":          healthStatus.OverallStatus,
			"recorder_healthy":        healthStatus.RecorderHealthy,
			"metrics_service_healthy": healthStatus.MetricsServiceHealthy,
			"consecutive_failures":    healthStatus.ConsecutiveFailures,
		},
	}

	h.logger.WithFields(logrus.Fields{
		"component":               "email_system_health_handler",
		"endpoint":                "/health/email-system/readiness",
		"is_ready":                isReady,
		"overall_status":          healthStatus.OverallStatus,
		"recorder_healthy":        healthStatus.RecorderHealthy,
		"metrics_service_healthy": healthStatus.MetricsServiceHealthy,
		"response_code":           statusCode,
		"action":                  "readiness_check",
	}).Debug("Email system readiness check")

	c.JSON(statusCode, response)
}

// GetEmailSystemDeepHealthCheck handles GET /health/email-system/deep
// Comprehensive deep health check with detailed component analysis
func (h *EmailSystemHealthHandler) GetEmailSystemDeepHealthCheck(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	h.logger.WithFields(logrus.Fields{
		"component": "email_system_health_handler",
		"endpoint":  "/health/email-system/deep",
		"action":    "deep_health_check_requested",
	}).Info("Deep email system health check requested")

	start := time.Now()

	// Perform comprehensive health check
	healthStatus := h.emailMonitoringService.PerformHealthCheck(ctx)

	// Get detailed metrics
	metrics := h.emailMonitoringService.GetSystemMetrics()

	// Get alert configuration
	alertConfig := h.emailMonitoringService.GetAlertConfig()

	// Get recent alerts
	recentAlerts := h.emailMonitoringService.GetRecentAlerts(24)

	// Determine status code
	statusCode := http.StatusOK
	if healthStatus.OverallStatus == "unhealthy" {
		statusCode = http.StatusServiceUnavailable
	} else if healthStatus.OverallStatus == "degraded" {
		statusCode = http.StatusPartialContent
	}

	deepCheckDuration := time.Since(start)

	response := gin.H{
		"status":                 healthStatus.OverallStatus,
		"timestamp":              time.Now().UTC().Format(time.RFC3339),
		"deep_check_duration_ms": deepCheckDuration.Nanoseconds() / 1e6,
		"health": gin.H{
			"overall_status":           healthStatus.OverallStatus,
			"recorder_healthy":         healthStatus.RecorderHealthy,
			"metrics_service_healthy":  healthStatus.MetricsServiceHealthy,
			"database_healthy":         healthStatus.DatabaseHealthy,
			"consecutive_failures":     healthStatus.ConsecutiveFailures,
			"health_check_duration_ms": healthStatus.HealthCheckDuration.Nanoseconds() / 1e6,
			"issues":                   healthStatus.Issues,
			"last_checked":             healthStatus.LastChecked.Format(time.RFC3339),
		},
		"metrics": gin.H{
			"system": gin.H{
				"uptime_seconds":                metrics.SystemUptime.Seconds(),
				"health_check_interval_seconds": metrics.HealthCheckInterval.Seconds(),
				"alerts_triggered":              metrics.AlertsTriggered,
				"alert_suppression_count":       metrics.AlertSuppressionCount,
				"average_processing_time_ms":    metrics.AverageEmailProcessingTime,
				"email_throughput_per_hour":     metrics.EmailThroughputPerHour,
				"last_health_check":             metrics.LastHealthCheck.Format(time.RFC3339),
			},
			"recorder":        h.getDetailedRecorderMetrics(metrics),
			"metrics_service": h.getDetailedMetricsServiceMetrics(metrics),
		},
		"alerting": gin.H{
			"config": gin.H{
				"recording_failure_threshold":    alertConfig.RecordingFailureThreshold,
				"metrics_failure_threshold":      alertConfig.MetricsFailureThreshold,
				"health_check_failure_threshold": alertConfig.HealthCheckFailureThreshold,
				"alert_suppression_window":       alertConfig.AlertSuppressionWindow.String(),
				"alert_recipients":               alertConfig.AlertRecipients,
			},
			"recent_alerts_24h": len(recentAlerts),
			"alerts":            recentAlerts,
		},
		"recommendations": h.generateHealthRecommendations(healthStatus, metrics),
	}

	h.logger.WithFields(logrus.Fields{
		"component":               "email_system_health_handler",
		"overall_status":          healthStatus.OverallStatus,
		"deep_check_duration_ms":  deepCheckDuration.Nanoseconds() / 1e6,
		"recorder_healthy":        healthStatus.RecorderHealthy,
		"metrics_service_healthy": healthStatus.MetricsServiceHealthy,
		"consecutive_failures":    healthStatus.ConsecutiveFailures,
		"issues_count":            len(healthStatus.Issues),
		"recent_alerts_count":     len(recentAlerts),
		"response_code":           statusCode,
		"action":                  "deep_health_check_completed",
	}).Info("Deep email system health check completed")

	c.JSON(statusCode, response)
}

// Helper methods for building health response data

func (h *EmailSystemHealthHandler) getRecorderHealthInfo(metrics *services.EmailMonitoringMetrics) gin.H {
	if metrics.RecorderMetrics == nil {
		return gin.H{
			"status": "unknown",
			"error":  "metrics not available",
		}
	}

	return gin.H{
		"status":                    h.determineComponentStatus(metrics.RecorderMetrics.SuccessRate),
		"total_recording_attempts":  metrics.RecorderMetrics.TotalRecordingAttempts,
		"successful_recordings":     metrics.RecorderMetrics.SuccessfulRecordings,
		"failed_recordings":         metrics.RecorderMetrics.FailedRecordings,
		"success_rate":              metrics.RecorderMetrics.SuccessRate,
		"average_recording_time_ms": metrics.RecorderMetrics.AverageRecordingTime,
		"health_check_failures":     metrics.RecorderMetrics.HealthCheckFailures,
	}
}

func (h *EmailSystemHealthHandler) getMetricsServiceHealthInfo(metrics *services.EmailMonitoringMetrics) gin.H {
	if metrics.MetricsServiceMetrics == nil {
		return gin.H{
			"status": "unknown",
			"error":  "metrics not available",
		}
	}

	return gin.H{
		"status":                   h.determineComponentStatus(metrics.MetricsServiceMetrics.SuccessRate),
		"total_metrics_requests":   metrics.MetricsServiceMetrics.TotalMetricsRequests,
		"successful_requests":      metrics.MetricsServiceMetrics.SuccessfulRequests,
		"failed_requests":          metrics.MetricsServiceMetrics.FailedRequests,
		"success_rate":             metrics.MetricsServiceMetrics.SuccessRate,
		"average_response_time_ms": metrics.MetricsServiceMetrics.AverageResponseTime,
		"health_check_failures":    metrics.MetricsServiceMetrics.HealthCheckFailures,
	}
}

func (h *EmailSystemHealthHandler) getDetailedRecorderMetrics(metrics *services.EmailMonitoringMetrics) gin.H {
	if metrics.RecorderMetrics == nil {
		return gin.H{"status": "unavailable"}
	}

	return gin.H{
		"status":                    h.determineComponentStatus(metrics.RecorderMetrics.SuccessRate),
		"total_recording_attempts":  metrics.RecorderMetrics.TotalRecordingAttempts,
		"successful_recordings":     metrics.RecorderMetrics.SuccessfulRecordings,
		"failed_recordings":         metrics.RecorderMetrics.FailedRecordings,
		"success_rate":              metrics.RecorderMetrics.SuccessRate,
		"average_recording_time_ms": metrics.RecorderMetrics.AverageRecordingTime,
		"last_recording_time":       metrics.RecorderMetrics.LastRecordingTime.Format(time.RFC3339),
		"retry_attempts":            metrics.RecorderMetrics.RetryAttempts,
		"health_check_failures":     metrics.RecorderMetrics.HealthCheckFailures,
	}
}

func (h *EmailSystemHealthHandler) getDetailedMetricsServiceMetrics(metrics *services.EmailMonitoringMetrics) gin.H {
	if metrics.MetricsServiceMetrics == nil {
		return gin.H{"status": "unavailable"}
	}

	return gin.H{
		"status":                   h.determineComponentStatus(metrics.MetricsServiceMetrics.SuccessRate),
		"total_metrics_requests":   metrics.MetricsServiceMetrics.TotalMetricsRequests,
		"successful_requests":      metrics.MetricsServiceMetrics.SuccessfulRequests,
		"failed_requests":          metrics.MetricsServiceMetrics.FailedRequests,
		"success_rate":             metrics.MetricsServiceMetrics.SuccessRate,
		"average_response_time_ms": metrics.MetricsServiceMetrics.AverageResponseTime,
		"last_request_time":        metrics.MetricsServiceMetrics.LastRequestTime.Format(time.RFC3339),
		"cache_hits":               metrics.MetricsServiceMetrics.CacheHits,
		"cache_misses":             metrics.MetricsServiceMetrics.CacheMisses,
		"health_check_failures":    metrics.MetricsServiceMetrics.HealthCheckFailures,
	}
}

func (h *EmailSystemHealthHandler) determineComponentStatus(successRate float64) string {
	if successRate >= 0.95 {
		return "healthy"
	} else if successRate >= 0.80 {
		return "degraded"
	} else {
		return "unhealthy"
	}
}

func (h *EmailSystemHealthHandler) generateHealthRecommendations(
	healthStatus *services.EmailSystemHealthStatus,
	metrics *services.EmailMonitoringMetrics,
) []string {
	recommendations := []string{}

	// Check for consecutive failures
	if healthStatus.ConsecutiveFailures >= 3 {
		recommendations = append(recommendations, "High consecutive failure count detected. Consider investigating database connectivity and system resources.")
	}

	// Check recorder success rate
	if metrics.RecorderMetrics != nil && metrics.RecorderMetrics.SuccessRate < 0.90 {
		recommendations = append(recommendations, "Email event recorder success rate is below 90%. Check database performance and connection pool settings.")
	}

	// Check metrics service success rate
	if metrics.MetricsServiceMetrics != nil && metrics.MetricsServiceMetrics.SuccessRate < 0.90 {
		recommendations = append(recommendations, "Email metrics service success rate is below 90%. Review query performance and database indexes.")
	}

	// Check response times
	if metrics.RecorderMetrics != nil && metrics.RecorderMetrics.AverageRecordingTime > 1000 {
		recommendations = append(recommendations, "Email event recording time is high (>1s). Consider optimizing database writes and connection pooling.")
	}

	if metrics.MetricsServiceMetrics != nil && metrics.MetricsServiceMetrics.AverageResponseTime > 500 {
		recommendations = append(recommendations, "Email metrics service response time is high (>500ms). Consider adding caching or optimizing queries.")
	}

	// Check alert frequency
	if metrics.AlertsTriggered > 10 {
		recommendations = append(recommendations, "High number of alerts triggered. Review alert thresholds and investigate underlying issues.")
	}

	// Check health check failures
	if metrics.RecorderMetrics != nil && metrics.RecorderMetrics.HealthCheckFailures > 5 {
		recommendations = append(recommendations, "Multiple email event recorder health check failures. Verify database connectivity and service configuration.")
	}

	if metrics.MetricsServiceMetrics != nil && metrics.MetricsServiceMetrics.HealthCheckFailures > 5 {
		recommendations = append(recommendations, "Multiple email metrics service health check failures. Check database performance and query optimization.")
	}

	// Default recommendation if system is healthy
	if len(recommendations) == 0 && healthStatus.OverallStatus == "healthy" {
		recommendations = append(recommendations, "Email event tracking system is operating normally. Continue monitoring for any performance degradation.")
	}

	return recommendations
}
