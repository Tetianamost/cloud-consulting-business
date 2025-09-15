package services

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/cloud-consulting/backend/internal/interfaces"
)

// EmailMonitoringService provides monitoring and observability for email event system
type EmailMonitoringService struct {
	emailEventRecorder  interfaces.EmailEventRecorder
	emailMetricsService interfaces.EmailMetricsService
	logger              *logrus.Logger

	// Metrics tracking
	metrics     *EmailMonitoringMetrics
	metricsLock sync.RWMutex

	// Alerting configuration
	alertConfig *EmailAlertConfig

	// Health check configuration
	healthCheckInterval time.Duration
	lastHealthCheck     time.Time
	healthStatus        *EmailSystemHealthStatus
	healthLock          sync.RWMutex
}

// EmailMonitoringMetrics aggregates metrics from all email system components
type EmailMonitoringMetrics struct {
	// Overall system metrics
	SystemUptime        time.Duration `json:"system_uptime"`
	LastHealthCheck     time.Time     `json:"last_health_check"`
	HealthCheckInterval time.Duration `json:"health_check_interval"`

	// Email event recorder metrics
	RecorderMetrics *interfaces.EmailEventRecorderMetrics `json:"recorder_metrics"`

	// Email metrics service metrics
	MetricsServiceMetrics *interfaces.EmailMetricsServiceMetrics `json:"metrics_service_metrics"`

	// Alert metrics
	AlertsTriggered       int64     `json:"alerts_triggered"`
	LastAlertTime         time.Time `json:"last_alert_time"`
	AlertSuppressionCount int64     `json:"alert_suppression_count"`

	// Performance metrics
	AverageEmailProcessingTime float64 `json:"average_email_processing_time_ms"`
	EmailThroughputPerHour     float64 `json:"email_throughput_per_hour"`
}

// EmailAlertConfig defines configuration for email system alerting
type EmailAlertConfig struct {
	// Failure rate thresholds
	RecordingFailureThreshold float64 `json:"recording_failure_threshold"` // e.g., 0.05 for 5%
	MetricsFailureThreshold   float64 `json:"metrics_failure_threshold"`   // e.g., 0.10 for 10%

	// Health check thresholds
	HealthCheckFailureThreshold int `json:"health_check_failure_threshold"` // consecutive failures

	// Alert suppression
	AlertSuppressionWindow time.Duration `json:"alert_suppression_window"` // e.g., 30 minutes

	// Notification settings
	AlertRecipients []string `json:"alert_recipients"`
	SlackWebhookURL string   `json:"slack_webhook_url,omitempty"`
}

// EmailSystemHealthStatus represents the overall health of the email system
type EmailSystemHealthStatus struct {
	OverallStatus         string        `json:"overall_status"` // "healthy", "degraded", "unhealthy"
	LastChecked           time.Time     `json:"last_checked"`
	RecorderHealthy       bool          `json:"recorder_healthy"`
	MetricsServiceHealthy bool          `json:"metrics_service_healthy"`
	DatabaseHealthy       bool          `json:"database_healthy"`
	ConsecutiveFailures   int           `json:"consecutive_failures"`
	HealthCheckDuration   time.Duration `json:"health_check_duration"`
	Issues                []string      `json:"issues,omitempty"`
}

// NewEmailMonitoringService creates a new email monitoring service
func NewEmailMonitoringService(
	emailEventRecorder interfaces.EmailEventRecorder,
	emailMetricsService interfaces.EmailMetricsService,
	logger *logrus.Logger,
) *EmailMonitoringService {
	service := &EmailMonitoringService{
		emailEventRecorder:  emailEventRecorder,
		emailMetricsService: emailMetricsService,
		logger:              logger,
		healthCheckInterval: 30 * time.Second, // Default 30 seconds
		metrics: &EmailMonitoringMetrics{
			SystemUptime:        0,
			LastHealthCheck:     time.Now(),
			HealthCheckInterval: 30 * time.Second,
		},
		alertConfig: &EmailAlertConfig{
			RecordingFailureThreshold:   0.05, // 5%
			MetricsFailureThreshold:     0.10, // 10%
			HealthCheckFailureThreshold: 3,    // 3 consecutive failures
			AlertSuppressionWindow:      30 * time.Minute,
			AlertRecipients:             []string{"admin@cloudpartner.pro"},
		},
		healthStatus: &EmailSystemHealthStatus{
			OverallStatus: "unknown",
			LastChecked:   time.Now(),
		},
	}

	// Start background health monitoring
	go service.startHealthMonitoring()

	return service
}

// GetSystemMetrics returns comprehensive email system metrics
func (s *EmailMonitoringService) GetSystemMetrics() *EmailMonitoringMetrics {
	s.metricsLock.RLock()
	defer s.metricsLock.RUnlock()

	// Update metrics from components
	if s.emailEventRecorder != nil {
		s.metrics.RecorderMetrics = s.emailEventRecorder.GetMetrics()
	}

	if s.emailMetricsService != nil {
		s.metrics.MetricsServiceMetrics = s.emailMetricsService.GetMetrics()
	}

	// Calculate derived metrics
	s.calculateDerivedMetrics()

	return s.metrics
}

// GetHealthStatus returns the current health status of the email system
func (s *EmailMonitoringService) GetHealthStatus() *EmailSystemHealthStatus {
	s.healthLock.RLock()
	defer s.healthLock.RUnlock()

	// Return a copy to prevent external modification
	status := *s.healthStatus
	return &status
}

// PerformHealthCheck performs a comprehensive health check of the email system
func (s *EmailMonitoringService) PerformHealthCheck(ctx context.Context) *EmailSystemHealthStatus {
	start := time.Now()

	s.logger.WithField("component", "email_monitoring").Debug("Starting email system health check")

	status := &EmailSystemHealthStatus{
		LastChecked: start,
		Issues:      []string{},
	}

	// Check email event recorder health
	if s.emailEventRecorder != nil {
		status.RecorderHealthy = s.emailEventRecorder.IsHealthyWithContext(ctx)
		if !status.RecorderHealthy {
			status.Issues = append(status.Issues, "Email event recorder is unhealthy")
		}
	} else {
		status.RecorderHealthy = false
		status.Issues = append(status.Issues, "Email event recorder not configured")
	}

	// Check email metrics service health
	if s.emailMetricsService != nil {
		status.MetricsServiceHealthy = s.emailMetricsService.IsHealthy(ctx)
		if !status.MetricsServiceHealthy {
			status.Issues = append(status.Issues, "Email metrics service is unhealthy")
		}
	} else {
		status.MetricsServiceHealthy = false
		status.Issues = append(status.Issues, "Email metrics service not configured")
	}

	// Determine overall status
	if status.RecorderHealthy && status.MetricsServiceHealthy {
		status.OverallStatus = "healthy"
	} else if status.RecorderHealthy || status.MetricsServiceHealthy {
		status.OverallStatus = "degraded"
	} else {
		status.OverallStatus = "unhealthy"
	}

	status.HealthCheckDuration = time.Since(start)

	// Update consecutive failures counter
	s.healthLock.Lock()
	if status.OverallStatus == "unhealthy" {
		s.healthStatus.ConsecutiveFailures++
	} else {
		s.healthStatus.ConsecutiveFailures = 0
	}
	status.ConsecutiveFailures = s.healthStatus.ConsecutiveFailures
	s.healthStatus = status
	s.healthLock.Unlock()

	// Log health check results
	s.logger.WithFields(logrus.Fields{
		"component":               "email_monitoring",
		"overall_status":          status.OverallStatus,
		"recorder_healthy":        status.RecorderHealthy,
		"metrics_service_healthy": status.MetricsServiceHealthy,
		"consecutive_failures":    status.ConsecutiveFailures,
		"health_check_duration":   status.HealthCheckDuration,
		"issues_count":            len(status.Issues),
	}).Info("Email system health check completed")

	// Check if alerting is needed
	s.checkAndTriggerAlerts(status)

	return status
}

// startHealthMonitoring starts the background health monitoring goroutine
func (s *EmailMonitoringService) startHealthMonitoring() {
	ticker := time.NewTicker(s.healthCheckInterval)
	defer ticker.Stop()

	s.logger.WithFields(logrus.Fields{
		"component": "email_monitoring",
		"interval":  s.healthCheckInterval,
	}).Info("Starting email system health monitoring")

	for {
		select {
		case <-ticker.C:
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			s.PerformHealthCheck(ctx)
			cancel()
		}
	}
}

// checkAndTriggerAlerts checks if alerts should be triggered based on current status
func (s *EmailMonitoringService) checkAndTriggerAlerts(status *EmailSystemHealthStatus) {
	s.metricsLock.Lock()
	defer s.metricsLock.Unlock()

	now := time.Now()

	// Check if we're in alert suppression window
	if s.metrics.LastAlertTime.Add(s.alertConfig.AlertSuppressionWindow).After(now) {
		atomic.AddInt64(&s.metrics.AlertSuppressionCount, 1)
		return
	}

	shouldAlert := false
	alertReasons := []string{}

	// Check consecutive health check failures
	if status.ConsecutiveFailures >= s.alertConfig.HealthCheckFailureThreshold {
		shouldAlert = true
		alertReasons = append(alertReasons, fmt.Sprintf("Email system unhealthy for %d consecutive checks", status.ConsecutiveFailures))
	}

	// Check recording failure rate
	if s.metrics.RecorderMetrics != nil && s.metrics.RecorderMetrics.SuccessRate < (1.0-s.alertConfig.RecordingFailureThreshold) {
		shouldAlert = true
		alertReasons = append(alertReasons, fmt.Sprintf("Email recording failure rate: %.2f%%", (1.0-s.metrics.RecorderMetrics.SuccessRate)*100))
	}

	// Check metrics service failure rate
	if s.metrics.MetricsServiceMetrics != nil && s.metrics.MetricsServiceMetrics.SuccessRate < (1.0-s.alertConfig.MetricsFailureThreshold) {
		shouldAlert = true
		alertReasons = append(alertReasons, fmt.Sprintf("Email metrics service failure rate: %.2f%%", (1.0-s.metrics.MetricsServiceMetrics.SuccessRate)*100))
	}

	if shouldAlert {
		s.triggerAlert(alertReasons, status)
		s.metrics.LastAlertTime = now
		atomic.AddInt64(&s.metrics.AlertsTriggered, 1)
	}
}

// triggerAlert sends alerts to configured recipients
func (s *EmailMonitoringService) triggerAlert(reasons []string, status *EmailSystemHealthStatus) {
	s.logger.WithFields(logrus.Fields{
		"component":      "email_monitoring",
		"alert_reasons":  reasons,
		"overall_status": status.OverallStatus,
		"issues":         status.Issues,
	}).Error("Email system alert triggered")

	// In a real implementation, you would send alerts via:
	// - Email notifications to alert recipients
	// - Slack webhook notifications
	// - PagerDuty integration
	// - Metrics to monitoring systems (Prometheus, etc.)

	// For now, we'll log structured alert information
	alertData := map[string]interface{}{
		"alert_type": "email_system_failure",
		"severity":   s.determineSeverity(status),
		"reasons":    reasons,
		"status":     status,
		"timestamp":  time.Now().UTC().Format(time.RFC3339),
		"alert_id":   fmt.Sprintf("email-alert-%d", time.Now().Unix()),
	}

	s.logger.WithFields(logrus.Fields{
		"component":  "email_monitoring",
		"alert_data": alertData,
		"action":     "alert_triggered",
	}).Error("EMAIL SYSTEM ALERT")
}

// determineSeverity determines the severity level of an alert
func (s *EmailMonitoringService) determineSeverity(status *EmailSystemHealthStatus) string {
	switch status.OverallStatus {
	case "unhealthy":
		if status.ConsecutiveFailures >= 10 {
			return "critical"
		}
		return "high"
	case "degraded":
		return "medium"
	default:
		return "low"
	}
}

// calculateDerivedMetrics calculates derived metrics from component metrics
func (s *EmailMonitoringService) calculateDerivedMetrics() {
	// Calculate average email processing time
	if s.metrics.RecorderMetrics != nil && s.metrics.MetricsServiceMetrics != nil {
		totalTime := s.metrics.RecorderMetrics.AverageRecordingTime + s.metrics.MetricsServiceMetrics.AverageResponseTime
		s.metrics.AverageEmailProcessingTime = totalTime / 2
	}

	// Calculate email throughput (simplified - would need actual email counts)
	if s.metrics.RecorderMetrics != nil {
		// This is a simplified calculation - in reality you'd track emails per time period
		hoursRunning := time.Since(s.metrics.LastHealthCheck).Hours()
		if hoursRunning > 0 {
			s.metrics.EmailThroughputPerHour = float64(s.metrics.RecorderMetrics.SuccessfulRecordings) / hoursRunning
		}
	}
}

// UpdateAlertConfig updates the alert configuration
func (s *EmailMonitoringService) UpdateAlertConfig(config *EmailAlertConfig) {
	s.metricsLock.Lock()
	defer s.metricsLock.Unlock()

	s.alertConfig = config

	s.logger.WithFields(logrus.Fields{
		"component":                      "email_monitoring",
		"recording_failure_threshold":    config.RecordingFailureThreshold,
		"metrics_failure_threshold":      config.MetricsFailureThreshold,
		"health_check_failure_threshold": config.HealthCheckFailureThreshold,
		"alert_suppression_window":       config.AlertSuppressionWindow,
	}).Info("Email monitoring alert configuration updated")
}

// GetAlertConfig returns the current alert configuration
func (s *EmailMonitoringService) GetAlertConfig() *EmailAlertConfig {
	s.metricsLock.RLock()
	defer s.metricsLock.RUnlock()

	// Return a copy to prevent external modification
	config := *s.alertConfig
	return &config
}

// IsHealthy returns whether the email monitoring service itself is healthy
func (s *EmailMonitoringService) IsHealthy() bool {
	s.healthLock.RLock()
	defer s.healthLock.RUnlock()

	// Check if health checks are running (last check should be recent)
	if time.Since(s.healthStatus.LastChecked) > s.healthCheckInterval*2 {
		return false
	}

	// Check if the monitored services are healthy
	return s.healthStatus.OverallStatus != "unhealthy"
}

// SetHealthCheckInterval updates the health check interval
func (s *EmailMonitoringService) SetHealthCheckInterval(interval time.Duration) {
	s.healthCheckInterval = interval
	s.metrics.HealthCheckInterval = interval

	s.logger.WithFields(logrus.Fields{
		"component": "email_monitoring",
		"interval":  interval,
	}).Info("Email monitoring health check interval updated")
}

// GetRecentAlerts returns recent alert information
func (s *EmailMonitoringService) GetRecentAlerts(hours int) []map[string]interface{} {
	// In a real implementation, this would query a persistent store of alerts
	// For now, return basic information
	s.metricsLock.RLock()
	defer s.metricsLock.RUnlock()

	alerts := []map[string]interface{}{}

	if !s.metrics.LastAlertTime.IsZero() && time.Since(s.metrics.LastAlertTime) < time.Duration(hours)*time.Hour {
		alerts = append(alerts, map[string]interface{}{
			"timestamp":   s.metrics.LastAlertTime.Format(time.RFC3339),
			"type":        "email_system_failure",
			"severity":    s.determineSeverity(s.healthStatus),
			"description": "Email system health check failure",
		})
	}

	return alerts
}
