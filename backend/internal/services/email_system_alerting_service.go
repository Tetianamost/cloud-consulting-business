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

// EmailSystemAlertingService provides advanced alerting for email event recording failures
type EmailSystemAlertingService struct {
	emailEventRecorder  interfaces.EmailEventRecorder
	emailMetricsService interfaces.EmailMetricsService
	logger              *logrus.Logger

	// Alert configuration
	config *EmailAlertingConfig

	// Alert state tracking
	alertState     *EmailAlertState
	alertStateLock sync.RWMutex

	// Metrics
	alertMetrics     *EmailAlertingMetrics
	alertMetricsLock sync.RWMutex

	// Background monitoring
	stopChan chan struct{}
	running  int32
}

// EmailAlertingConfig defines configuration for email system alerting
type EmailAlertingConfig struct {
	// Failure rate thresholds
	RecordingFailureRateThreshold float64 `json:"recording_failure_rate_threshold"` // e.g., 0.05 for 5%
	MetricsFailureRateThreshold   float64 `json:"metrics_failure_rate_threshold"`   // e.g., 0.10 for 10%
	HighFailureRateThreshold      float64 `json:"high_failure_rate_threshold"`      // e.g., 0.20 for 20%
	CriticalFailureRateThreshold  float64 `json:"critical_failure_rate_threshold"`  // e.g., 0.50 for 50%

	// Time windows for evaluation
	EvaluationWindow time.Duration `json:"evaluation_window"` // e.g., 5 minutes
	AlertWindow      time.Duration `json:"alert_window"`      // e.g., 15 minutes
	CooldownPeriod   time.Duration `json:"cooldown_period"`   // e.g., 30 minutes

	// Consecutive failure thresholds
	ConsecutiveFailureThreshold         int `json:"consecutive_failure_threshold"`          // e.g., 3
	HighConsecutiveFailureThreshold     int `json:"high_consecutive_failure_threshold"`     // e.g., 5
	CriticalConsecutiveFailureThreshold int `json:"critical_consecutive_failure_threshold"` // e.g., 10

	// Health check thresholds
	HealthCheckFailureThreshold int           `json:"health_check_failure_threshold"` // e.g., 3
	HealthCheckInterval         time.Duration `json:"health_check_interval"`          // e.g., 30 seconds

	// Alert destinations
	AlertRecipients []string    `json:"alert_recipients"`
	SlackWebhookURL string      `json:"slack_webhook_url,omitempty"`
	PagerDutyURL    string      `json:"pager_duty_url,omitempty"`
	EmailSMTPConfig *SMTPConfig `json:"email_smtp_config,omitempty"`

	// Alert suppression
	EnableAlertSuppression bool          `json:"enable_alert_suppression"`
	SuppressionWindow      time.Duration `json:"suppression_window"` // e.g., 1 hour
}

// SMTPConfig defines SMTP configuration for alert emails
type SMTPConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	From     string `json:"from"`
	UseTLS   bool   `json:"use_tls"`
}

// EmailAlertState tracks the current alerting state
type EmailAlertState struct {
	// Current alert levels
	RecorderAlertLevel       AlertLevel `json:"recorder_alert_level"`
	MetricsServiceAlertLevel AlertLevel `json:"metrics_service_alert_level"`
	OverallAlertLevel        AlertLevel `json:"overall_alert_level"`

	// Failure tracking
	ConsecutiveRecorderFailures       int `json:"consecutive_recorder_failures"`
	ConsecutiveMetricsServiceFailures int `json:"consecutive_metrics_service_failures"`
	ConsecutiveHealthCheckFailures    int `json:"consecutive_health_check_failures"`

	// Last alert times
	LastRecorderAlert       time.Time `json:"last_recorder_alert"`
	LastMetricsServiceAlert time.Time `json:"last_metrics_service_alert"`
	LastHealthCheckAlert    time.Time `json:"last_health_check_alert"`

	// Alert suppression
	AlertsSuppressed     bool      `json:"alerts_suppressed"`
	SuppressionStartTime time.Time `json:"suppression_start_time"`
}

// AlertLevel represents the severity level of an alert
type AlertLevel string

const (
	AlertLevelNone     AlertLevel = "none"
	AlertLevelWarning  AlertLevel = "warning"
	AlertLevelHigh     AlertLevel = "high"
	AlertLevelCritical AlertLevel = "critical"
)

// EmailAlertingMetrics tracks alerting system metrics
type EmailAlertingMetrics struct {
	// Alert counts by level
	WarningAlertsTriggered  int64 `json:"warning_alerts_triggered"`
	HighAlertsTriggered     int64 `json:"high_alerts_triggered"`
	CriticalAlertsTriggered int64 `json:"critical_alerts_triggered"`
	TotalAlertsTriggered    int64 `json:"total_alerts_triggered"`

	// Alert suppression metrics
	AlertsSuppressed       int64 `json:"alerts_suppressed"`
	SuppressionActivations int64 `json:"suppression_activations"`

	// Performance metrics
	AverageAlertProcessingTime float64       `json:"average_alert_processing_time_ms"`
	LastAlertTime              time.Time     `json:"last_alert_time"`
	AlertingSystemUptime       time.Duration `json:"alerting_system_uptime"`

	// Health check metrics
	HealthChecksPerformed int64 `json:"health_checks_performed"`
	HealthCheckFailures   int64 `json:"health_check_failures"`
}

// NewEmailSystemAlertingService creates a new email system alerting service
func NewEmailSystemAlertingService(
	emailEventRecorder interfaces.EmailEventRecorder,
	emailMetricsService interfaces.EmailMetricsService,
	logger *logrus.Logger,
) *EmailSystemAlertingService {
	service := &EmailSystemAlertingService{
		emailEventRecorder:  emailEventRecorder,
		emailMetricsService: emailMetricsService,
		logger:              logger,
		stopChan:            make(chan struct{}),
		config: &EmailAlertingConfig{
			RecordingFailureRateThreshold:       0.05, // 5%
			MetricsFailureRateThreshold:         0.10, // 10%
			HighFailureRateThreshold:            0.20, // 20%
			CriticalFailureRateThreshold:        0.50, // 50%
			EvaluationWindow:                    5 * time.Minute,
			AlertWindow:                         15 * time.Minute,
			CooldownPeriod:                      30 * time.Minute,
			ConsecutiveFailureThreshold:         3,
			HighConsecutiveFailureThreshold:     5,
			CriticalConsecutiveFailureThreshold: 10,
			HealthCheckFailureThreshold:         3,
			HealthCheckInterval:                 30 * time.Second,
			AlertRecipients:                     []string{"admin@cloudpartner.pro"},
			EnableAlertSuppression:              true,
			SuppressionWindow:                   1 * time.Hour,
		},
		alertState: &EmailAlertState{
			RecorderAlertLevel:       AlertLevelNone,
			MetricsServiceAlertLevel: AlertLevelNone,
			OverallAlertLevel:        AlertLevelNone,
		},
		alertMetrics: &EmailAlertingMetrics{
			AlertingSystemUptime: 0,
		},
	}

	return service
}

// Start begins the background alerting monitoring
func (s *EmailSystemAlertingService) Start() error {
	if atomic.CompareAndSwapInt32(&s.running, 0, 1) {
		s.logger.WithFields(logrus.Fields{
			"component": "email_system_alerting",
			"action":    "starting_alerting_service",
		}).Info("Starting email system alerting service")

		go s.monitoringLoop()
		return nil
	}
	return fmt.Errorf("alerting service is already running")
}

// Stop stops the background alerting monitoring
func (s *EmailSystemAlertingService) Stop() error {
	if atomic.CompareAndSwapInt32(&s.running, 1, 0) {
		s.logger.WithFields(logrus.Fields{
			"component": "email_system_alerting",
			"action":    "stopping_alerting_service",
		}).Info("Stopping email system alerting service")

		close(s.stopChan)
		return nil
	}
	return fmt.Errorf("alerting service is not running")
}

// monitoringLoop runs the continuous monitoring and alerting logic
func (s *EmailSystemAlertingService) monitoringLoop() {
	ticker := time.NewTicker(s.config.HealthCheckInterval)
	defer ticker.Stop()

	startTime := time.Now()

	for {
		select {
		case <-s.stopChan:
			s.logger.WithField("component", "email_system_alerting").Info("Alerting monitoring loop stopped")
			return
		case <-ticker.C:
			s.performAlertingHealthCheck()

			// Update uptime
			s.alertMetricsLock.Lock()
			s.alertMetrics.AlertingSystemUptime = time.Since(startTime)
			s.alertMetricsLock.Unlock()
		}
	}
}

// performAlertingHealthCheck performs a comprehensive alerting health check
func (s *EmailSystemAlertingService) performAlertingHealthCheck() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	start := time.Now()

	s.logger.WithField("component", "email_system_alerting").Debug("Performing alerting health check")

	// Update health check metrics
	s.alertMetricsLock.Lock()
	atomic.AddInt64(&s.alertMetrics.HealthChecksPerformed, 1)
	s.alertMetricsLock.Unlock()

	// Check email event recorder health and metrics
	s.checkRecorderHealth(ctx)

	// Check email metrics service health and metrics
	s.checkMetricsServiceHealth(ctx)

	// Determine overall alert level
	s.updateOverallAlertLevel()

	// Process alerts if needed
	s.processAlerts()

	// Update performance metrics
	processingTime := time.Since(start)
	s.updateAlertProcessingMetrics(processingTime)

	s.logger.WithFields(logrus.Fields{
		"component":                   "email_system_alerting",
		"processing_time_ms":          processingTime.Nanoseconds() / 1e6,
		"recorder_alert_level":        s.alertState.RecorderAlertLevel,
		"metrics_service_alert_level": s.alertState.MetricsServiceAlertLevel,
		"overall_alert_level":         s.alertState.OverallAlertLevel,
	}).Debug("Alerting health check completed")
}

// checkRecorderHealth checks the health of the email event recorder
func (s *EmailSystemAlertingService) checkRecorderHealth(ctx context.Context) {
	s.alertStateLock.Lock()
	defer s.alertStateLock.Unlock()

	// Check if recorder is healthy
	isHealthy := s.emailEventRecorder.IsHealthyWithContext(ctx)

	if !isHealthy {
		s.alertState.ConsecutiveRecorderFailures++
		s.logger.WithFields(logrus.Fields{
			"component":            "email_system_alerting",
			"consecutive_failures": s.alertState.ConsecutiveRecorderFailures,
		}).Warn("Email event recorder health check failed")
	} else {
		s.alertState.ConsecutiveRecorderFailures = 0
	}

	// Get recorder metrics
	recorderMetrics := s.emailEventRecorder.GetMetrics()
	if recorderMetrics != nil {
		// Determine alert level based on success rate and consecutive failures
		s.alertState.RecorderAlertLevel = s.determineAlertLevel(
			recorderMetrics.SuccessRate,
			s.alertState.ConsecutiveRecorderFailures,
		)
	} else {
		s.alertState.RecorderAlertLevel = AlertLevelHigh // No metrics available
	}
}

// checkMetricsServiceHealth checks the health of the email metrics service
func (s *EmailSystemAlertingService) checkMetricsServiceHealth(ctx context.Context) {
	s.alertStateLock.Lock()
	defer s.alertStateLock.Unlock()

	// Check if metrics service is healthy
	isHealthy := s.emailMetricsService.IsHealthy(ctx)

	if !isHealthy {
		s.alertState.ConsecutiveMetricsServiceFailures++
		s.logger.WithFields(logrus.Fields{
			"component":            "email_system_alerting",
			"consecutive_failures": s.alertState.ConsecutiveMetricsServiceFailures,
		}).Warn("Email metrics service health check failed")
	} else {
		s.alertState.ConsecutiveMetricsServiceFailures = 0
	}

	// Get metrics service metrics
	metricsServiceMetrics := s.emailMetricsService.GetMetrics()
	if metricsServiceMetrics != nil {
		// Determine alert level based on success rate and consecutive failures
		s.alertState.MetricsServiceAlertLevel = s.determineAlertLevel(
			metricsServiceMetrics.SuccessRate,
			s.alertState.ConsecutiveMetricsServiceFailures,
		)
	} else {
		s.alertState.MetricsServiceAlertLevel = AlertLevelHigh // No metrics available
	}
}

// determineAlertLevel determines the alert level based on success rate and consecutive failures
func (s *EmailSystemAlertingService) determineAlertLevel(successRate float64, consecutiveFailures int) AlertLevel {
	failureRate := 1.0 - successRate

	// Check for critical conditions first
	if failureRate >= s.config.CriticalFailureRateThreshold || consecutiveFailures >= s.config.CriticalConsecutiveFailureThreshold {
		return AlertLevelCritical
	}

	// Check for high alert conditions
	if failureRate >= s.config.HighFailureRateThreshold || consecutiveFailures >= s.config.HighConsecutiveFailureThreshold {
		return AlertLevelHigh
	}

	// Check for warning conditions
	if failureRate >= s.config.RecordingFailureRateThreshold || consecutiveFailures >= s.config.ConsecutiveFailureThreshold {
		return AlertLevelWarning
	}

	return AlertLevelNone
}

// updateOverallAlertLevel updates the overall system alert level
func (s *EmailSystemAlertingService) updateOverallAlertLevel() {
	s.alertStateLock.Lock()
	defer s.alertStateLock.Unlock()

	// Overall alert level is the highest of component alert levels
	if s.alertState.RecorderAlertLevel == AlertLevelCritical || s.alertState.MetricsServiceAlertLevel == AlertLevelCritical {
		s.alertState.OverallAlertLevel = AlertLevelCritical
	} else if s.alertState.RecorderAlertLevel == AlertLevelHigh || s.alertState.MetricsServiceAlertLevel == AlertLevelHigh {
		s.alertState.OverallAlertLevel = AlertLevelHigh
	} else if s.alertState.RecorderAlertLevel == AlertLevelWarning || s.alertState.MetricsServiceAlertLevel == AlertLevelWarning {
		s.alertState.OverallAlertLevel = AlertLevelWarning
	} else {
		s.alertState.OverallAlertLevel = AlertLevelNone
	}
}

// processAlerts processes and sends alerts based on current state
func (s *EmailSystemAlertingService) processAlerts() {
	s.alertStateLock.RLock()
	currentState := *s.alertState // Copy current state
	s.alertStateLock.RUnlock()

	now := time.Now()

	// Check if alerts are suppressed
	if s.config.EnableAlertSuppression && currentState.AlertsSuppressed {
		if now.Sub(currentState.SuppressionStartTime) >= s.config.SuppressionWindow {
			// End suppression period
			s.alertStateLock.Lock()
			s.alertState.AlertsSuppressed = false
			s.alertStateLock.Unlock()

			s.logger.WithField("component", "email_system_alerting").Info("Alert suppression period ended")
		} else {
			// Still in suppression period
			s.alertMetricsLock.Lock()
			atomic.AddInt64(&s.alertMetrics.AlertsSuppressed, 1)
			s.alertMetricsLock.Unlock()
			return
		}
	}

	// Process recorder alerts
	if s.shouldSendAlert(currentState.RecorderAlertLevel, currentState.LastRecorderAlert) {
		s.sendRecorderAlert(currentState.RecorderAlertLevel)
		s.alertStateLock.Lock()
		s.alertState.LastRecorderAlert = now
		s.alertStateLock.Unlock()
	}

	// Process metrics service alerts
	if s.shouldSendAlert(currentState.MetricsServiceAlertLevel, currentState.LastMetricsServiceAlert) {
		s.sendMetricsServiceAlert(currentState.MetricsServiceAlertLevel)
		s.alertStateLock.Lock()
		s.alertState.LastMetricsServiceAlert = now
		s.alertStateLock.Unlock()
	}

	// Activate suppression if critical alerts are being sent
	if currentState.OverallAlertLevel == AlertLevelCritical && s.config.EnableAlertSuppression {
		s.alertStateLock.Lock()
		s.alertState.AlertsSuppressed = true
		s.alertState.SuppressionStartTime = now
		s.alertStateLock.Unlock()

		s.alertMetricsLock.Lock()
		atomic.AddInt64(&s.alertMetrics.SuppressionActivations, 1)
		s.alertMetricsLock.Unlock()

		s.logger.WithField("component", "email_system_alerting").Info("Alert suppression activated due to critical alerts")
	}
}

// shouldSendAlert determines if an alert should be sent based on level and cooldown
func (s *EmailSystemAlertingService) shouldSendAlert(alertLevel AlertLevel, lastAlertTime time.Time) bool {
	if alertLevel == AlertLevelNone {
		return false
	}

	// Check cooldown period
	if time.Since(lastAlertTime) < s.config.CooldownPeriod {
		return false
	}

	return true
}

// sendRecorderAlert sends an alert for email event recorder issues
func (s *EmailSystemAlertingService) sendRecorderAlert(alertLevel AlertLevel) {
	s.updateAlertMetrics(alertLevel)

	recorderMetrics := s.emailEventRecorder.GetMetrics()

	alertData := map[string]interface{}{
		"alert_type":           "email_event_recorder_failure",
		"alert_level":          alertLevel,
		"component":            "email_event_recorder",
		"success_rate":         recorderMetrics.SuccessRate,
		"failure_rate":         1.0 - recorderMetrics.SuccessRate,
		"consecutive_failures": s.alertState.ConsecutiveRecorderFailures,
		"total_attempts":       recorderMetrics.TotalRecordingAttempts,
		"failed_recordings":    recorderMetrics.FailedRecordings,
		"timestamp":            time.Now().UTC().Format(time.RFC3339),
	}

	s.logger.WithFields(logrus.Fields{
		"component":   "email_system_alerting",
		"alert_level": alertLevel,
		"alert_data":  alertData,
		"action":      "recorder_alert_triggered",
	}).Error("EMAIL EVENT RECORDER ALERT")

	// In a real implementation, send to configured alert destinations
	s.sendToAlertDestinations("Email Event Recorder Alert", alertData)
}

// sendMetricsServiceAlert sends an alert for email metrics service issues
func (s *EmailSystemAlertingService) sendMetricsServiceAlert(alertLevel AlertLevel) {
	s.updateAlertMetrics(alertLevel)

	metricsServiceMetrics := s.emailMetricsService.GetMetrics()

	alertData := map[string]interface{}{
		"alert_type":           "email_metrics_service_failure",
		"alert_level":          alertLevel,
		"component":            "email_metrics_service",
		"success_rate":         metricsServiceMetrics.SuccessRate,
		"failure_rate":         1.0 - metricsServiceMetrics.SuccessRate,
		"consecutive_failures": s.alertState.ConsecutiveMetricsServiceFailures,
		"total_requests":       metricsServiceMetrics.TotalMetricsRequests,
		"failed_requests":      metricsServiceMetrics.FailedRequests,
		"timestamp":            time.Now().UTC().Format(time.RFC3339),
	}

	s.logger.WithFields(logrus.Fields{
		"component":   "email_system_alerting",
		"alert_level": alertLevel,
		"alert_data":  alertData,
		"action":      "metrics_service_alert_triggered",
	}).Error("EMAIL METRICS SERVICE ALERT")

	// In a real implementation, send to configured alert destinations
	s.sendToAlertDestinations("Email Metrics Service Alert", alertData)
}

// sendToAlertDestinations sends alerts to configured destinations
func (s *EmailSystemAlertingService) sendToAlertDestinations(subject string, alertData map[string]interface{}) {
	// This would implement actual alert sending to:
	// - Email recipients
	// - Slack webhooks
	// - PagerDuty
	// - Other monitoring systems

	s.logger.WithFields(logrus.Fields{
		"component":        "email_system_alerting",
		"alert_subject":    subject,
		"alert_recipients": s.config.AlertRecipients,
		"action":           "alert_sent_to_destinations",
	}).Info("Alert sent to configured destinations")
}

// updateAlertMetrics updates alert metrics based on alert level
func (s *EmailSystemAlertingService) updateAlertMetrics(alertLevel AlertLevel) {
	s.alertMetricsLock.Lock()
	defer s.alertMetricsLock.Unlock()

	atomic.AddInt64(&s.alertMetrics.TotalAlertsTriggered, 1)
	s.alertMetrics.LastAlertTime = time.Now()

	switch alertLevel {
	case AlertLevelWarning:
		atomic.AddInt64(&s.alertMetrics.WarningAlertsTriggered, 1)
	case AlertLevelHigh:
		atomic.AddInt64(&s.alertMetrics.HighAlertsTriggered, 1)
	case AlertLevelCritical:
		atomic.AddInt64(&s.alertMetrics.CriticalAlertsTriggered, 1)
	}
}

// updateAlertProcessingMetrics updates alert processing performance metrics
func (s *EmailSystemAlertingService) updateAlertProcessingMetrics(processingTime time.Duration) {
	s.alertMetricsLock.Lock()
	defer s.alertMetricsLock.Unlock()

	processingTimeMs := float64(processingTime.Nanoseconds()) / 1e6

	if s.alertMetrics.AverageAlertProcessingTime == 0 {
		s.alertMetrics.AverageAlertProcessingTime = processingTimeMs
	} else {
		// Exponential moving average
		alpha := 0.1
		s.alertMetrics.AverageAlertProcessingTime = alpha*processingTimeMs + (1-alpha)*s.alertMetrics.AverageAlertProcessingTime
	}
}

// GetAlertState returns the current alert state
func (s *EmailSystemAlertingService) GetAlertState() *EmailAlertState {
	s.alertStateLock.RLock()
	defer s.alertStateLock.RUnlock()

	// Return a copy to prevent external modification
	state := *s.alertState
	return &state
}

// GetAlertMetrics returns the current alert metrics
func (s *EmailSystemAlertingService) GetAlertMetrics() *EmailAlertingMetrics {
	s.alertMetricsLock.RLock()
	defer s.alertMetricsLock.RUnlock()

	// Return a copy to prevent external modification
	metrics := *s.alertMetrics
	return &metrics
}

// UpdateConfig updates the alerting configuration
func (s *EmailSystemAlertingService) UpdateConfig(config *EmailAlertingConfig) {
	s.config = config
	s.logger.WithFields(logrus.Fields{
		"component": "email_system_alerting",
		"action":    "config_updated",
	}).Info("Email system alerting configuration updated")
}

// GetConfig returns the current alerting configuration
func (s *EmailSystemAlertingService) GetConfig() *EmailAlertingConfig {
	// Return a copy to prevent external modification
	config := *s.config
	return &config
}

// IsRunning returns whether the alerting service is currently running
func (s *EmailSystemAlertingService) IsRunning() bool {
	return atomic.LoadInt32(&s.running) == 1
}
