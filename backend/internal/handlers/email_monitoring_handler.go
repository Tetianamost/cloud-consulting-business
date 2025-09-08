package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/cloud-consulting/backend/internal/services"
)

// EmailMonitoringHandler handles email monitoring and observability endpoints
type EmailMonitoringHandler struct {
	emailMonitoringService *services.EmailMonitoringService
	logger                 *logrus.Logger
}

// NewEmailMonitoringHandler creates a new email monitoring handler
func NewEmailMonitoringHandler(
	emailMonitoringService *services.EmailMonitoringService,
	logger *logrus.Logger,
) *EmailMonitoringHandler {
	return &EmailMonitoringHandler{
		emailMonitoringService: emailMonitoringService,
		logger:                 logger,
	}
}

// GetEmailSystemHealth handles GET /health/email-system
func (h *EmailMonitoringHandler) GetEmailSystemHealth(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	h.logger.WithField("component", "email_monitoring_handler").Debug("Getting email system health status")

	healthStatus := h.emailMonitoringService.PerformHealthCheck(ctx)

	statusCode := http.StatusOK
	if healthStatus.OverallStatus == "unhealthy" {
		statusCode = http.StatusServiceUnavailable
	}

	response := gin.H{
		"status":    healthStatus.OverallStatus,
		"timestamp": healthStatus.LastChecked.Format(time.RFC3339),
		"health":    healthStatus,
	}

	c.JSON(statusCode, response)
}

// GetEmailSystemMetrics handles GET /metrics/email-system
func (h *EmailMonitoringHandler) GetEmailSystemMetrics(c *gin.Context) {
	h.logger.WithField("component", "email_monitoring_handler").Debug("Getting email system metrics")

	metrics := h.emailMonitoringService.GetSystemMetrics()

	response := gin.H{
		"success":   true,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"metrics":   metrics,
	}

	c.JSON(http.StatusOK, response)
}

// GetEmailAlerts handles GET /alerts/email-system
func (h *EmailMonitoringHandler) GetEmailAlerts(c *gin.Context) {
	// Get hours parameter (default to 24 hours)
	hoursStr := c.DefaultQuery("hours", "24")
	hours, err := strconv.Atoi(hoursStr)
	if err != nil || hours < 1 || hours > 168 { // Max 1 week
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid hours parameter. Must be between 1 and 168.",
		})
		return
	}

	h.logger.WithFields(logrus.Fields{
		"component": "email_monitoring_handler",
		"hours":     hours,
	}).Debug("Getting recent email system alerts")

	alerts := h.emailMonitoringService.GetRecentAlerts(hours)

	response := gin.H{
		"success":   true,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"hours":     hours,
		"alerts":    alerts,
	}

	c.JSON(http.StatusOK, response)
}

// UpdateAlertConfig handles PUT /config/email-alerts
func (h *EmailMonitoringHandler) UpdateAlertConfig(c *gin.Context) {
	var config services.EmailAlertConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		h.logger.WithError(err).Error("Failed to parse alert configuration")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid alert configuration format",
		})
		return
	}

	// Validate configuration
	if config.RecordingFailureThreshold < 0 || config.RecordingFailureThreshold > 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Recording failure threshold must be between 0 and 1",
		})
		return
	}

	if config.MetricsFailureThreshold < 0 || config.MetricsFailureThreshold > 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Metrics failure threshold must be between 0 and 1",
		})
		return
	}

	if config.HealthCheckFailureThreshold < 1 || config.HealthCheckFailureThreshold > 100 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Health check failure threshold must be between 1 and 100",
		})
		return
	}

	if config.AlertSuppressionWindow < time.Minute || config.AlertSuppressionWindow > 24*time.Hour {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Alert suppression window must be between 1 minute and 24 hours",
		})
		return
	}

	h.emailMonitoringService.UpdateAlertConfig(&config)

	h.logger.WithFields(logrus.Fields{
		"component":                      "email_monitoring_handler",
		"recording_failure_threshold":    config.RecordingFailureThreshold,
		"metrics_failure_threshold":      config.MetricsFailureThreshold,
		"health_check_failure_threshold": config.HealthCheckFailureThreshold,
		"alert_suppression_window":       config.AlertSuppressionWindow,
	}).Info("Email alert configuration updated")

	response := gin.H{
		"success": true,
		"message": "Alert configuration updated successfully",
		"config":  config,
	}

	c.JSON(http.StatusOK, response)
}

// GetAlertConfig handles GET /config/email-alerts
func (h *EmailMonitoringHandler) GetAlertConfig(c *gin.Context) {
	config := h.emailMonitoringService.GetAlertConfig()

	response := gin.H{
		"success": true,
		"config":  config,
	}

	c.JSON(http.StatusOK, response)
}

// TriggerHealthCheck handles POST /health/email-system/check
func (h *EmailMonitoringHandler) TriggerHealthCheck(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 15*time.Second)
	defer cancel()

	h.logger.WithField("component", "email_monitoring_handler").Info("Manual email system health check triggered")

	healthStatus := h.emailMonitoringService.PerformHealthCheck(ctx)

	statusCode := http.StatusOK
	if healthStatus.OverallStatus == "unhealthy" {
		statusCode = http.StatusServiceUnavailable
	}

	response := gin.H{
		"success":   true,
		"message":   "Health check completed",
		"timestamp": healthStatus.LastChecked.Format(time.RFC3339),
		"health":    healthStatus,
	}

	c.JSON(statusCode, response)
}

// GetEmailSystemStatus handles GET /status/email-system
func (h *EmailMonitoringHandler) GetEmailSystemStatus(c *gin.Context) {
	h.logger.WithField("component", "email_monitoring_handler").Debug("Getting comprehensive email system status")

	// Get current health status
	healthStatus := h.emailMonitoringService.GetHealthStatus()

	// Get current metrics
	metrics := h.emailMonitoringService.GetSystemMetrics()

	// Get alert configuration
	alertConfig := h.emailMonitoringService.GetAlertConfig()

	// Get recent alerts (last 24 hours)
	recentAlerts := h.emailMonitoringService.GetRecentAlerts(24)

	statusCode := http.StatusOK
	if healthStatus.OverallStatus == "unhealthy" {
		statusCode = http.StatusServiceUnavailable
	}

	response := gin.H{
		"success":        true,
		"timestamp":      time.Now().UTC().Format(time.RFC3339),
		"overall_status": healthStatus.OverallStatus,
		"health":         healthStatus,
		"metrics":        metrics,
		"alert_config":   alertConfig,
		"recent_alerts":  recentAlerts,
		"monitoring": gin.H{
			"health_check_interval": metrics.HealthCheckInterval.String(),
			"system_uptime":         time.Since(time.Now().Add(-metrics.SystemUptime)).String(),
			"monitoring_active":     h.emailMonitoringService.IsHealthy(),
		},
	}

	c.JSON(statusCode, response)
}

// GetPrometheusMetrics handles GET /metrics/prometheus/email-system
func (h *EmailMonitoringHandler) GetPrometheusMetrics(c *gin.Context) {
	metrics := h.emailMonitoringService.GetSystemMetrics()

	// Generate Prometheus-style metrics
	prometheusMetrics := h.generatePrometheusMetrics(metrics)

	c.Header("Content-Type", "text/plain; version=0.0.4; charset=utf-8")
	c.String(http.StatusOK, prometheusMetrics)
}

// generatePrometheusMetrics converts internal metrics to Prometheus format
func (h *EmailMonitoringHandler) generatePrometheusMetrics(metrics *services.EmailMonitoringMetrics) string {
	var output string

	// Email event recorder metrics
	if metrics.RecorderMetrics != nil {
		output += "# HELP email_event_recorder_total_attempts Total number of email event recording attempts\n"
		output += "# TYPE email_event_recorder_total_attempts counter\n"
		output += "email_event_recorder_total_attempts " + strconv.FormatInt(metrics.RecorderMetrics.TotalRecordingAttempts, 10) + "\n"

		output += "# HELP email_event_recorder_successful_recordings Total number of successful email event recordings\n"
		output += "# TYPE email_event_recorder_successful_recordings counter\n"
		output += "email_event_recorder_successful_recordings " + strconv.FormatInt(metrics.RecorderMetrics.SuccessfulRecordings, 10) + "\n"

		output += "# HELP email_event_recorder_failed_recordings Total number of failed email event recordings\n"
		output += "# TYPE email_event_recorder_failed_recordings counter\n"
		output += "email_event_recorder_failed_recordings " + strconv.FormatInt(metrics.RecorderMetrics.FailedRecordings, 10) + "\n"

		output += "# HELP email_event_recorder_success_rate Success rate of email event recordings (0-1)\n"
		output += "# TYPE email_event_recorder_success_rate gauge\n"
		output += "email_event_recorder_success_rate " + strconv.FormatFloat(metrics.RecorderMetrics.SuccessRate, 'f', 4, 64) + "\n"

		output += "# HELP email_event_recorder_average_recording_time_ms Average time to record an email event in milliseconds\n"
		output += "# TYPE email_event_recorder_average_recording_time_ms gauge\n"
		output += "email_event_recorder_average_recording_time_ms " + strconv.FormatFloat(metrics.RecorderMetrics.AverageRecordingTime, 'f', 2, 64) + "\n"

		output += "# HELP email_event_recorder_retry_attempts Total number of retry attempts\n"
		output += "# TYPE email_event_recorder_retry_attempts counter\n"
		output += "email_event_recorder_retry_attempts " + strconv.FormatInt(metrics.RecorderMetrics.RetryAttempts, 10) + "\n"

		output += "# HELP email_event_recorder_health_check_failures Total number of health check failures\n"
		output += "# TYPE email_event_recorder_health_check_failures counter\n"
		output += "email_event_recorder_health_check_failures " + strconv.FormatInt(metrics.RecorderMetrics.HealthCheckFailures, 10) + "\n"
	}

	// Email metrics service metrics
	if metrics.MetricsServiceMetrics != nil {
		output += "# HELP email_metrics_service_total_requests Total number of metrics requests\n"
		output += "# TYPE email_metrics_service_total_requests counter\n"
		output += "email_metrics_service_total_requests " + strconv.FormatInt(metrics.MetricsServiceMetrics.TotalMetricsRequests, 10) + "\n"

		output += "# HELP email_metrics_service_successful_requests Total number of successful metrics requests\n"
		output += "# TYPE email_metrics_service_successful_requests counter\n"
		output += "email_metrics_service_successful_requests " + strconv.FormatInt(metrics.MetricsServiceMetrics.SuccessfulRequests, 10) + "\n"

		output += "# HELP email_metrics_service_failed_requests Total number of failed metrics requests\n"
		output += "# TYPE email_metrics_service_failed_requests counter\n"
		output += "email_metrics_service_failed_requests " + strconv.FormatInt(metrics.MetricsServiceMetrics.FailedRequests, 10) + "\n"

		output += "# HELP email_metrics_service_success_rate Success rate of metrics requests (0-1)\n"
		output += "# TYPE email_metrics_service_success_rate gauge\n"
		output += "email_metrics_service_success_rate " + strconv.FormatFloat(metrics.MetricsServiceMetrics.SuccessRate, 'f', 4, 64) + "\n"

		output += "# HELP email_metrics_service_average_response_time_ms Average response time for metrics requests in milliseconds\n"
		output += "# TYPE email_metrics_service_average_response_time_ms gauge\n"
		output += "email_metrics_service_average_response_time_ms " + strconv.FormatFloat(metrics.MetricsServiceMetrics.AverageResponseTime, 'f', 2, 64) + "\n"

		output += "# HELP email_metrics_service_health_check_failures Total number of health check failures\n"
		output += "# TYPE email_metrics_service_health_check_failures counter\n"
		output += "email_metrics_service_health_check_failures " + strconv.FormatInt(metrics.MetricsServiceMetrics.HealthCheckFailures, 10) + "\n"
	}

	// System-level metrics
	output += "# HELP email_system_alerts_triggered Total number of alerts triggered\n"
	output += "# TYPE email_system_alerts_triggered counter\n"
	output += "email_system_alerts_triggered " + strconv.FormatInt(metrics.AlertsTriggered, 10) + "\n"

	output += "# HELP email_system_alert_suppression_count Total number of suppressed alerts\n"
	output += "# TYPE email_system_alert_suppression_count counter\n"
	output += "email_system_alert_suppression_count " + strconv.FormatInt(metrics.AlertSuppressionCount, 10) + "\n"

	output += "# HELP email_system_average_processing_time_ms Average email processing time in milliseconds\n"
	output += "# TYPE email_system_average_processing_time_ms gauge\n"
	output += "email_system_average_processing_time_ms " + strconv.FormatFloat(metrics.AverageEmailProcessingTime, 'f', 2, 64) + "\n"

	output += "# HELP email_system_throughput_per_hour Email throughput per hour\n"
	output += "# TYPE email_system_throughput_per_hour gauge\n"
	output += "email_system_throughput_per_hour " + strconv.FormatFloat(metrics.EmailThroughputPerHour, 'f', 2, 64) + "\n"

	return output
}
