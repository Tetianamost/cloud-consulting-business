package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// AlertSeverity represents alert severity levels
type AlertSeverity string

const (
	AlertSeverityLow      AlertSeverity = "low"
	AlertSeverityMedium   AlertSeverity = "medium"
	AlertSeverityHigh     AlertSeverity = "high"
	AlertSeverityCritical AlertSeverity = "critical"
)

// AlertType represents different types of alerts
type AlertType string

const (
	AlertTypeError       AlertType = "error"
	AlertTypePerformance AlertType = "performance"
	AlertTypeSecurity    AlertType = "security"
	AlertTypeConnection  AlertType = "connection"
	AlertTypeAI          AlertType = "ai"
	AlertTypeSystem      AlertType = "system"
)

// Alert represents an alert event
type Alert struct {
	ID            string                 `json:"id"`
	Type          AlertType              `json:"type"`
	Severity      AlertSeverity          `json:"severity"`
	Title         string                 `json:"title"`
	Description   string                 `json:"description"`
	Component     string                 `json:"component"`
	Metric        string                 `json:"metric,omitempty"`
	Value         interface{}            `json:"value,omitempty"`
	Threshold     interface{}            `json:"threshold,omitempty"`
	CorrelationID string                 `json:"correlation_id,omitempty"`
	UserID        string                 `json:"user_id,omitempty"`
	SessionID     string                 `json:"session_id,omitempty"`
	Metadata      map[string]interface{} `json:"metadata"`
	Timestamp     time.Time              `json:"timestamp"`
	Resolved      bool                   `json:"resolved"`
	ResolvedAt    *time.Time             `json:"resolved_at,omitempty"`
}

// AlertRule represents a rule for triggering alerts
type AlertRule struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Type        AlertType              `json:"type"`
	Severity    AlertSeverity          `json:"severity"`
	Metric      string                 `json:"metric"`
	Condition   string                 `json:"condition"` // "gt", "lt", "eq", "gte", "lte"
	Threshold   interface{}            `json:"threshold"`
	Duration    time.Duration          `json:"duration"` // How long condition must be true
	Cooldown    time.Duration          `json:"cooldown"` // Minimum time between alerts
	Enabled     bool                   `json:"enabled"`
	Description string                 `json:"description"`
	Metadata    map[string]interface{} `json:"metadata"`
	LastFired   *time.Time             `json:"last_fired,omitempty"`
}

// AlertChannel represents a notification channel
type AlertChannel struct {
	ID      string                 `json:"id"`
	Name    string                 `json:"name"`
	Type    string                 `json:"type"` // "webhook", "email", "slack", "teams"
	Config  map[string]interface{} `json:"config"`
	Enabled bool                   `json:"enabled"`
	Filters []AlertFilter          `json:"filters,omitempty"`
}

// AlertFilter represents filters for alert channels
type AlertFilter struct {
	Field    string      `json:"field"`
	Operator string      `json:"operator"` // "eq", "ne", "contains", "in"
	Value    interface{} `json:"value"`
}

// WebhookConfig represents webhook configuration
type WebhookConfig struct {
	URL     string            `json:"url"`
	Method  string            `json:"method"`
	Headers map[string]string `json:"headers"`
	Timeout time.Duration     `json:"timeout"`
}

// ChatAlertingService manages alerts for the chat system
type ChatAlertingService struct {
	logger           *logrus.Logger
	structuredLogger *ChatStructuredLogger
	metricsCollector *ChatMetricsCollector

	// Alert management
	alerts     map[string]*Alert
	alertRules map[string]*AlertRule
	channels   map[string]*AlertChannel

	// State management
	alertHistory []Alert
	mu           sync.RWMutex

	// Background processing
	ctx    context.Context
	cancel context.CancelFunc
}

// NewChatAlertingService creates a new alerting service
func NewChatAlertingService(
	logger *logrus.Logger,
	structuredLogger *ChatStructuredLogger,
	metricsCollector *ChatMetricsCollector,
) *ChatAlertingService {
	ctx, cancel := context.WithCancel(context.Background())

	service := &ChatAlertingService{
		logger:           logger,
		structuredLogger: structuredLogger,
		metricsCollector: metricsCollector,
		alerts:           make(map[string]*Alert),
		alertRules:       make(map[string]*AlertRule),
		channels:         make(map[string]*AlertChannel),
		alertHistory:     make([]Alert, 0),
		ctx:              ctx,
		cancel:           cancel,
	}

	// Initialize default alert rules
	service.initializeDefaultRules()

	// Initialize default channels
	service.initializeDefaultChannels()

	return service
}

// initializeDefaultRules sets up default alert rules
func (s *ChatAlertingService) initializeDefaultRules() {
	defaultRules := []*AlertRule{
		{
			ID:          "high_error_rate",
			Name:        "High Error Rate",
			Type:        AlertTypeError,
			Severity:    AlertSeverityHigh,
			Metric:      "error_rate",
			Condition:   "gt",
			Threshold:   0.05, // 5%
			Duration:    2 * time.Minute,
			Cooldown:    10 * time.Minute,
			Enabled:     true,
			Description: "Error rate is above 5%",
		},
		{
			ID:          "low_connection_success_rate",
			Name:        "Low Connection Success Rate",
			Type:        AlertTypeConnection,
			Severity:    AlertSeverityMedium,
			Metric:      "connection_success_rate",
			Condition:   "lt",
			Threshold:   0.95, // 95%
			Duration:    5 * time.Minute,
			Cooldown:    15 * time.Minute,
			Enabled:     true,
			Description: "Connection success rate is below 95%",
		},
		{
			ID:          "high_response_time",
			Name:        "High Response Time",
			Type:        AlertTypePerformance,
			Severity:    AlertSeverityMedium,
			Metric:      "response_time_p95",
			Condition:   "gt",
			Threshold:   int64(2000), // 2 seconds in milliseconds
			Duration:    3 * time.Minute,
			Cooldown:    10 * time.Minute,
			Enabled:     true,
			Description: "95th percentile response time is above 2 seconds",
		},
		{
			ID:          "low_cache_hit_rate",
			Name:        "Low Cache Hit Rate",
			Type:        AlertTypePerformance,
			Severity:    AlertSeverityLow,
			Metric:      "cache_hit_rate",
			Condition:   "lt",
			Threshold:   0.7, // 70%
			Duration:    10 * time.Minute,
			Cooldown:    30 * time.Minute,
			Enabled:     true,
			Description: "Cache hit rate is below 70%",
		},
		{
			ID:          "high_ai_cost",
			Name:        "High AI Usage Cost",
			Type:        AlertTypeAI,
			Severity:    AlertSeverityMedium,
			Metric:      "ai_cost_daily",
			Condition:   "gt",
			Threshold:   100.0, // $100
			Duration:    1 * time.Minute,
			Cooldown:    6 * time.Hour,
			Enabled:     true,
			Description: "Daily AI usage cost exceeds $100",
		},
		{
			ID:          "authentication_failures",
			Name:        "Multiple Authentication Failures",
			Type:        AlertTypeSecurity,
			Severity:    AlertSeverityHigh,
			Metric:      "auth_failures_per_minute",
			Condition:   "gt",
			Threshold:   10,
			Duration:    1 * time.Minute,
			Cooldown:    5 * time.Minute,
			Enabled:     true,
			Description: "More than 10 authentication failures per minute",
		},
	}

	for _, rule := range defaultRules {
		s.alertRules[rule.ID] = rule
	}
}

// initializeDefaultChannels sets up default notification channels
func (s *ChatAlertingService) initializeDefaultChannels() {
	// Default webhook channel (can be configured via environment variables)
	webhookChannel := &AlertChannel{
		ID:      "default_webhook",
		Name:    "Default Webhook",
		Type:    "webhook",
		Enabled: false, // Disabled by default, enable when configured
		Config: map[string]interface{}{
			"url":     "",
			"method":  "POST",
			"timeout": "30s",
			"headers": map[string]string{
				"Content-Type": "application/json",
			},
		},
	}

	s.channels[webhookChannel.ID] = webhookChannel
}

// CreateAlert creates a new alert
func (s *ChatAlertingService) CreateAlert(alertType AlertType, severity AlertSeverity, title, description, component string, metadata map[string]interface{}) *Alert {
	s.mu.Lock()
	defer s.mu.Unlock()

	alert := &Alert{
		ID:          s.generateAlertID(),
		Type:        alertType,
		Severity:    severity,
		Title:       title,
		Description: description,
		Component:   component,
		Metadata:    metadata,
		Timestamp:   time.Now(),
		Resolved:    false,
	}

	if metadata != nil {
		if correlationID, ok := metadata["correlation_id"].(string); ok {
			alert.CorrelationID = correlationID
		}
		if userID, ok := metadata["user_id"].(string); ok {
			alert.UserID = userID
		}
		if sessionID, ok := metadata["session_id"].(string); ok {
			alert.SessionID = sessionID
		}
	}

	s.alerts[alert.ID] = alert
	s.alertHistory = append(s.alertHistory, *alert)

	// Log the alert creation
	logCtx := s.structuredLogger.WithCorrelationID(alert.CorrelationID).
		SetComponent("alerting").
		SetOperation("create_alert").
		AddMetadata("alert_id", alert.ID).
		AddMetadata("alert_type", string(alertType)).
		AddMetadata("severity", string(severity))

	s.structuredLogger.Info(logCtx, fmt.Sprintf("Alert created: %s", title))

	// Send notifications
	go s.sendAlertNotifications(alert)

	return alert
}

// ResolveAlert marks an alert as resolved
func (s *ChatAlertingService) ResolveAlert(alertID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	alert, exists := s.alerts[alertID]
	if !exists {
		return fmt.Errorf("alert not found: %s", alertID)
	}

	if alert.Resolved {
		return fmt.Errorf("alert already resolved: %s", alertID)
	}

	now := time.Now()
	alert.Resolved = true
	alert.ResolvedAt = &now

	// Log the alert resolution
	logCtx := s.structuredLogger.WithCorrelationID(alert.CorrelationID).
		SetComponent("alerting").
		SetOperation("resolve_alert").
		AddMetadata("alert_id", alert.ID).
		AddMetadata("resolution_time", now.Sub(alert.Timestamp))

	s.structuredLogger.Info(logCtx, fmt.Sprintf("Alert resolved: %s", alert.Title))

	return nil
}

// CheckMetricsForAlerts checks current metrics against alert rules
func (s *ChatAlertingService) CheckMetricsForAlerts() {
	if s.metricsCollector == nil {
		return
	}

	metrics := s.metricsCollector.GetMetrics()
	now := time.Now()

	s.mu.RLock()
	rules := make([]*AlertRule, 0, len(s.alertRules))
	for _, rule := range s.alertRules {
		if rule.Enabled {
			rules = append(rules, rule)
		}
	}
	s.mu.RUnlock()

	for _, rule := range rules {
		// Check cooldown
		if rule.LastFired != nil && now.Sub(*rule.LastFired) < rule.Cooldown {
			continue
		}

		// Get metric value
		var metricValue interface{}
		var shouldAlert bool

		switch rule.Metric {
		case "error_rate":
			metricValue = metrics.ErrorMetrics.ErrorRate
			shouldAlert = s.evaluateCondition(metricValue, rule.Condition, rule.Threshold)
		case "connection_success_rate":
			metricValue = metrics.ConnectionMetrics.ConnectionSuccessRate
			shouldAlert = s.evaluateCondition(metricValue, rule.Condition, rule.Threshold)
		case "response_time_p95":
			metricValue = metrics.PerformanceMetrics.ResponseTimeP95.Milliseconds()
			shouldAlert = s.evaluateCondition(metricValue, rule.Condition, rule.Threshold)
		case "cache_hit_rate":
			if metrics.CacheMetrics != nil {
				metricValue = s.metricsCollector.cacheMonitor.GetCacheHitRatio()
				shouldAlert = s.evaluateCondition(metricValue, rule.Condition, rule.Threshold)
			}
		case "ai_cost_daily":
			metricValue = metrics.AIMetrics.EstimatedCost
			shouldAlert = s.evaluateCondition(metricValue, rule.Condition, rule.Threshold)
		}

		if shouldAlert {
			// Create alert
			metadata := map[string]interface{}{
				"rule_id":   rule.ID,
				"metric":    rule.Metric,
				"value":     metricValue,
				"threshold": rule.Threshold,
				"condition": rule.Condition,
			}

			alert := s.CreateAlert(
				rule.Type,
				rule.Severity,
				rule.Name,
				fmt.Sprintf("%s: %v (threshold: %v)", rule.Description, metricValue, rule.Threshold),
				"metrics",
				metadata,
			)

			alert.Metric = rule.Metric
			alert.Value = metricValue
			alert.Threshold = rule.Threshold

			// Update rule last fired time
			s.mu.Lock()
			rule.LastFired = &now
			s.mu.Unlock()
		}
	}
}

// evaluateCondition evaluates a condition against a threshold
func (s *ChatAlertingService) evaluateCondition(value interface{}, condition string, threshold interface{}) bool {
	switch condition {
	case "gt":
		return s.compareValues(value, threshold) > 0
	case "gte":
		return s.compareValues(value, threshold) >= 0
	case "lt":
		return s.compareValues(value, threshold) < 0
	case "lte":
		return s.compareValues(value, threshold) <= 0
	case "eq":
		return s.compareValues(value, threshold) == 0
	case "ne":
		return s.compareValues(value, threshold) != 0
	}
	return false
}

// compareValues compares two values and returns -1, 0, or 1
func (s *ChatAlertingService) compareValues(a, b interface{}) int {
	// Convert to float64 for comparison
	aFloat, aOk := s.toFloat64(a)
	bFloat, bOk := s.toFloat64(b)

	if !aOk || !bOk {
		return 0 // Cannot compare
	}

	if aFloat < bFloat {
		return -1
	} else if aFloat > bFloat {
		return 1
	}
	return 0
}

// toFloat64 converts various numeric types to float64
func (s *ChatAlertingService) toFloat64(v interface{}) (float64, bool) {
	switch val := v.(type) {
	case float64:
		return val, true
	case float32:
		return float64(val), true
	case int:
		return float64(val), true
	case int32:
		return float64(val), true
	case int64:
		return float64(val), true
	case uint:
		return float64(val), true
	case uint32:
		return float64(val), true
	case uint64:
		return float64(val), true
	default:
		return 0, false
	}
}

// sendAlertNotifications sends notifications for an alert
func (s *ChatAlertingService) sendAlertNotifications(alert *Alert) {
	s.mu.RLock()
	channels := make([]*AlertChannel, 0, len(s.channels))
	for _, channel := range s.channels {
		if channel.Enabled && s.shouldSendToChannel(alert, channel) {
			channels = append(channels, channel)
		}
	}
	s.mu.RUnlock()

	for _, channel := range channels {
		go s.sendToChannel(alert, channel)
	}
}

// shouldSendToChannel checks if an alert should be sent to a specific channel
func (s *ChatAlertingService) shouldSendToChannel(alert *Alert, channel *AlertChannel) bool {
	for _, filter := range channel.Filters {
		if !s.evaluateFilter(alert, filter) {
			return false
		}
	}
	return true
}

// evaluateFilter evaluates a filter against an alert
func (s *ChatAlertingService) evaluateFilter(alert *Alert, filter AlertFilter) bool {
	var fieldValue interface{}

	switch filter.Field {
	case "type":
		fieldValue = string(alert.Type)
	case "severity":
		fieldValue = string(alert.Severity)
	case "component":
		fieldValue = alert.Component
	default:
		if alert.Metadata != nil {
			fieldValue = alert.Metadata[filter.Field]
		}
	}

	switch filter.Operator {
	case "eq":
		return fieldValue == filter.Value
	case "ne":
		return fieldValue != filter.Value
	case "contains":
		if str, ok := fieldValue.(string); ok {
			if filterStr, ok := filter.Value.(string); ok {
				return contains(str, filterStr)
			}
		}
	case "in":
		if slice, ok := filter.Value.([]interface{}); ok {
			for _, item := range slice {
				if fieldValue == item {
					return true
				}
			}
		}
	}

	return false
}

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(substr) <= len(s) && s[:len(substr)] == substr) ||
		(len(s) > len(substr) && containsHelper(s, substr)))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// sendToChannel sends an alert to a specific channel
func (s *ChatAlertingService) sendToChannel(alert *Alert, channel *AlertChannel) {
	switch channel.Type {
	case "webhook":
		s.sendWebhookNotification(alert, channel)
	case "email":
		s.sendEmailNotification(alert, channel)
	default:
		s.logger.WithFields(logrus.Fields{
			"channel_type": channel.Type,
			"channel_id":   channel.ID,
			"alert_id":     alert.ID,
		}).Warn("Unsupported channel type")
	}
}

// sendWebhookNotification sends a webhook notification
func (s *ChatAlertingService) sendWebhookNotification(alert *Alert, channel *AlertChannel) {
	config := channel.Config
	url, ok := config["url"].(string)
	if !ok || url == "" {
		s.logger.WithField("channel_id", channel.ID).Error("Webhook URL not configured")
		return
	}

	method := "POST"
	if m, ok := config["method"].(string); ok {
		method = m
	}

	timeout := 30 * time.Second
	if t, ok := config["timeout"].(string); ok {
		if parsed, err := time.ParseDuration(t); err == nil {
			timeout = parsed
		}
	}

	// Prepare payload
	payload := map[string]interface{}{
		"alert":     alert,
		"channel":   channel.Name,
		"timestamp": time.Now(),
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		s.logger.WithError(err).WithField("alert_id", alert.ID).Error("Failed to marshal webhook payload")
		return
	}

	// Create HTTP request
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(jsonData))
	if err != nil {
		s.logger.WithError(err).WithField("alert_id", alert.ID).Error("Failed to create webhook request")
		return
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	if headers, ok := config["headers"].(map[string]string); ok {
		for key, value := range headers {
			req.Header.Set(key, value)
		}
	}

	// Send request
	client := &http.Client{Timeout: timeout}
	resp, err := client.Do(req)
	if err != nil {
		s.logger.WithError(err).WithFields(logrus.Fields{
			"alert_id":   alert.ID,
			"channel_id": channel.ID,
			"url":        url,
		}).Error("Failed to send webhook notification")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		s.logger.WithFields(logrus.Fields{
			"alert_id":    alert.ID,
			"channel_id":  channel.ID,
			"status_code": resp.StatusCode,
		}).Info("Webhook notification sent successfully")
	} else {
		s.logger.WithFields(logrus.Fields{
			"alert_id":    alert.ID,
			"channel_id":  channel.ID,
			"status_code": resp.StatusCode,
		}).Error("Webhook notification failed")
	}
}

// sendEmailNotification sends an email notification (placeholder)
func (s *ChatAlertingService) sendEmailNotification(alert *Alert, channel *AlertChannel) {
	// This is a placeholder for email notifications
	// In a real implementation, this would integrate with an email service
	s.logger.WithFields(logrus.Fields{
		"alert_id":   alert.ID,
		"channel_id": channel.ID,
	}).Info("Email notification would be sent here")
}

// GetActiveAlerts returns all active (unresolved) alerts
func (s *ChatAlertingService) GetActiveAlerts() []*Alert {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var activeAlerts []*Alert
	for _, alert := range s.alerts {
		if !alert.Resolved {
			activeAlerts = append(activeAlerts, alert)
		}
	}

	return activeAlerts
}

// GetAlertHistory returns alert history with pagination
func (s *ChatAlertingService) GetAlertHistory(limit, offset int) []Alert {
	s.mu.RLock()
	defer s.mu.RUnlock()

	total := len(s.alertHistory)
	if offset >= total {
		return []Alert{}
	}

	end := offset + limit
	if end > total {
		end = total
	}

	// Return in reverse chronological order (newest first)
	result := make([]Alert, end-offset)
	for i := 0; i < end-offset; i++ {
		result[i] = s.alertHistory[total-1-offset-i]
	}

	return result
}

// StartPeriodicChecks starts periodic metric checks for alerts
func (s *ChatAlertingService) StartPeriodicChecks(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	s.logger.WithField("interval", interval).Info("Started periodic alert checks")

	for {
		select {
		case <-s.ctx.Done():
			s.logger.Info("Stopped periodic alert checks")
			return
		case <-ticker.C:
			s.CheckMetricsForAlerts()
		}
	}
}

// Stop stops the alerting service
func (s *ChatAlertingService) Stop() {
	s.cancel()
}

// generateAlertID generates a unique alert ID
func (s *ChatAlertingService) generateAlertID() string {
	return fmt.Sprintf("alert_%d", time.Now().UnixNano())
}

// GetAlertRules returns all alert rules
func (s *ChatAlertingService) GetAlertRules() []*AlertRule {
	s.mu.RLock()
	defer s.mu.RUnlock()

	rules := make([]*AlertRule, 0, len(s.alertRules))
	for _, rule := range s.alertRules {
		rules = append(rules, rule)
	}

	return rules
}

// GetAlertChannels returns all alert channels
func (s *ChatAlertingService) GetAlertChannels() []*AlertChannel {
	s.mu.RLock()
	defer s.mu.RUnlock()

	channels := make([]*AlertChannel, 0, len(s.channels))
	for _, channel := range s.channels {
		channels = append(channels, channel)
	}

	return channels
}

// UpdateAlertRule updates an existing alert rule
func (s *ChatAlertingService) UpdateAlertRule(rule *AlertRule) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.alertRules[rule.ID]; !exists {
		return fmt.Errorf("alert rule not found: %s", rule.ID)
	}

	s.alertRules[rule.ID] = rule
	return nil
}

// UpdateAlertChannel updates an existing alert channel
func (s *ChatAlertingService) UpdateAlertChannel(channel *AlertChannel) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.channels[channel.ID]; !exists {
		return fmt.Errorf("alert channel not found: %s", channel.ID)
	}

	s.channels[channel.ID] = channel
	return nil
}
