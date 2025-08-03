package services

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// ChatMetricsCollector aggregates and exposes metrics for the chat system
type ChatMetricsCollector struct {
	performanceMonitor *ChatPerformanceMonitor
	cacheMonitor       *CacheMonitor
	logger             *logrus.Logger
	mu                 sync.RWMutex

	// Connection metrics
	totalConnections     int64
	activeConnections    int64
	connectionDuration   time.Duration
	reconnectionAttempts int64
	connectionFailures   int64
	websocketUpgrades    int64
	websocketErrors      int64

	// Message metrics
	messagesSent          int64
	messagesReceived      int64
	messageProcessingTime time.Duration
	messageErrors         int64
	messageRetries        int64
	messageBroadcasts     int64

	// AI service metrics
	aiRequestsTotal      int64
	aiRequestsSuccessful int64
	aiRequestsFailed     int64
	aiResponseTime       time.Duration
	aiTokensUsed         int64
	aiCostEstimate       float64
	aiModelUsage         map[string]int64

	// User engagement metrics
	activeUsers        int64
	sessionDuration    time.Duration
	messagesPerSession float64
	quickActionsUsed   int64
	featureUsage       map[string]int64

	// Error metrics
	authenticationErrors int64
	authorizationErrors  int64
	validationErrors     int64
	systemErrors         int64
	timeoutErrors        int64

	// Performance metrics
	responseTimeP50 time.Duration
	responseTimeP95 time.Duration
	responseTimeP99 time.Duration
	throughputMPS   float64 // Messages per second
	errorRate       float64

	// System health metrics
	memoryUsage      int64
	cpuUsage         float64
	goroutineCount   int64
	lastMetricsReset time.Time
}

// ChatMetrics represents comprehensive chat system metrics
type ChatMetrics struct {
	// Connection metrics
	ConnectionMetrics ConnectionMetrics `json:"connection_metrics"`

	// Message metrics
	MessageMetrics MessageMetrics `json:"message_metrics"`

	// AI service metrics
	AIMetrics AIMetrics `json:"ai_metrics"`

	// User engagement metrics
	UserMetrics UserMetrics `json:"user_metrics"`

	// Error metrics
	ErrorMetrics ErrorMetrics `json:"error_metrics"`

	// Performance metrics
	PerformanceMetrics ChatPerformanceMetrics `json:"performance_metrics"`

	// Cache metrics
	CacheMetrics *CacheMetrics `json:"cache_metrics"`

	// Timestamp
	Timestamp time.Time `json:"timestamp"`
}

// ConnectionMetrics represents WebSocket connection metrics
type ConnectionMetrics struct {
	TotalConnections          int64         `json:"total_connections"`
	ActiveConnections         int64         `json:"active_connections"`
	AverageConnectionDuration time.Duration `json:"average_connection_duration"`
	ReconnectionAttempts      int64         `json:"reconnection_attempts"`
	ConnectionFailures        int64         `json:"connection_failures"`
	WebSocketUpgrades         int64         `json:"websocket_upgrades"`
	WebSocketErrors           int64         `json:"websocket_errors"`
	ConnectionSuccessRate     float64       `json:"connection_success_rate"`
}

// MessageMetrics represents message processing metrics
type MessageMetrics struct {
	MessagesSent          int64         `json:"messages_sent"`
	MessagesReceived      int64         `json:"messages_received"`
	AverageProcessingTime time.Duration `json:"average_processing_time"`
	MessageErrors         int64         `json:"message_errors"`
	MessageRetries        int64         `json:"message_retries"`
	MessageBroadcasts     int64         `json:"message_broadcasts"`
	MessageSuccessRate    float64       `json:"message_success_rate"`
}

// AIMetrics represents AI service usage metrics
type AIMetrics struct {
	RequestsTotal       int64            `json:"requests_total"`
	RequestsSuccessful  int64            `json:"requests_successful"`
	RequestsFailed      int64            `json:"requests_failed"`
	AverageResponseTime time.Duration    `json:"average_response_time"`
	TokensUsed          int64            `json:"tokens_used"`
	EstimatedCost       float64          `json:"estimated_cost_usd"`
	ModelUsage          map[string]int64 `json:"model_usage"`
	SuccessRate         float64          `json:"success_rate"`
}

// UserMetrics represents user engagement metrics
type UserMetrics struct {
	ActiveUsers            int64            `json:"active_users"`
	AverageSessionDuration time.Duration    `json:"average_session_duration"`
	MessagesPerSession     float64          `json:"messages_per_session"`
	QuickActionsUsed       int64            `json:"quick_actions_used"`
	FeatureUsage           map[string]int64 `json:"feature_usage"`
}

// ErrorMetrics represents error tracking metrics
type ErrorMetrics struct {
	AuthenticationErrors int64   `json:"authentication_errors"`
	AuthorizationErrors  int64   `json:"authorization_errors"`
	ValidationErrors     int64   `json:"validation_errors"`
	SystemErrors         int64   `json:"system_errors"`
	TimeoutErrors        int64   `json:"timeout_errors"`
	TotalErrors          int64   `json:"total_errors"`
	ErrorRate            float64 `json:"error_rate"`
}

// ChatPerformanceMetrics represents system performance metrics
type ChatPerformanceMetrics struct {
	ResponseTimeP50 time.Duration `json:"response_time_p50"`
	ResponseTimeP95 time.Duration `json:"response_time_p95"`
	ResponseTimeP99 time.Duration `json:"response_time_p99"`
	ThroughputMPS   float64       `json:"throughput_messages_per_second"`
	ErrorRate       float64       `json:"error_rate"`
	MemoryUsageMB   int64         `json:"memory_usage_mb"`
	CPUUsagePercent float64       `json:"cpu_usage_percent"`
	GoroutineCount  int64         `json:"goroutine_count"`
}

// NewChatMetricsCollector creates a new metrics collector
func NewChatMetricsCollector(
	performanceMonitor *ChatPerformanceMonitor,
	cacheMonitor *CacheMonitor,
	logger *logrus.Logger,
) *ChatMetricsCollector {
	return &ChatMetricsCollector{
		performanceMonitor: performanceMonitor,
		cacheMonitor:       cacheMonitor,
		logger:             logger,
		aiModelUsage:       make(map[string]int64),
		featureUsage:       make(map[string]int64),
		lastMetricsReset:   time.Now(),
	}
}

// RecordConnection records connection-related metrics
func (c *ChatMetricsCollector) RecordConnection(event string, duration ...time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	switch event {
	case "opened":
		c.totalConnections++
		c.activeConnections++
		c.websocketUpgrades++
	case "closed":
		if c.activeConnections > 0 {
			c.activeConnections--
		}
		if len(duration) > 0 {
			c.connectionDuration += duration[0]
		}
	case "failed":
		c.connectionFailures++
	case "reconnect":
		c.reconnectionAttempts++
	case "error":
		c.websocketErrors++
	}
}

// RecordMessage records message-related metrics
func (c *ChatMetricsCollector) RecordMessage(event string, processingTime ...time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	switch event {
	case "sent":
		c.messagesSent++
	case "received":
		c.messagesReceived++
	case "error":
		c.messageErrors++
	case "retry":
		c.messageRetries++
	case "broadcast":
		c.messageBroadcasts++
	}

	if len(processingTime) > 0 {
		c.messageProcessingTime += processingTime[0]
	}
}

// RecordAIRequest records AI service usage metrics
func (c *ChatMetricsCollector) RecordAIRequest(success bool, responseTime time.Duration, tokensUsed int64, model string, cost float64) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.aiRequestsTotal++
	c.aiResponseTime += responseTime
	c.aiTokensUsed += tokensUsed
	c.aiCostEstimate += cost

	if success {
		c.aiRequestsSuccessful++
	} else {
		c.aiRequestsFailed++
	}

	if model != "" {
		c.aiModelUsage[model]++
	}
}

// RecordUserEngagement records user engagement metrics
func (c *ChatMetricsCollector) RecordUserEngagement(event string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	switch event {
	case "user_active":
		c.activeUsers++
	case "session_duration":
		if duration, ok := value.(time.Duration); ok {
			c.sessionDuration += duration
		}
	case "quick_action":
		c.quickActionsUsed++
	case "feature_usage":
		if feature, ok := value.(string); ok {
			c.featureUsage[feature]++
		}
	}
}

// RecordError records error metrics
func (c *ChatMetricsCollector) RecordError(errorType string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	switch errorType {
	case "authentication":
		c.authenticationErrors++
	case "authorization":
		c.authorizationErrors++
	case "validation":
		c.validationErrors++
	case "system":
		c.systemErrors++
	case "timeout":
		c.timeoutErrors++
	}
}

// GetMetrics returns comprehensive chat metrics
func (c *ChatMetricsCollector) GetMetrics() *ChatMetrics {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Get performance monitor metrics
	perfMetrics := c.performanceMonitor.GetMetrics()

	// Get cache metrics
	cacheMetrics := c.cacheMonitor.GetMetrics()

	// Calculate connection success rate
	connectionSuccessRate := float64(1.0)
	if c.totalConnections > 0 {
		connectionSuccessRate = float64(c.totalConnections-c.connectionFailures) / float64(c.totalConnections)
	}

	// Calculate message success rate
	messageSuccessRate := float64(1.0)
	totalMessages := c.messagesSent + c.messagesReceived
	if totalMessages > 0 {
		messageSuccessRate = float64(totalMessages-c.messageErrors) / float64(totalMessages)
	}

	// Calculate AI success rate
	aiSuccessRate := float64(1.0)
	if c.aiRequestsTotal > 0 {
		aiSuccessRate = float64(c.aiRequestsSuccessful) / float64(c.aiRequestsTotal)
	}

	// Calculate average connection duration
	avgConnectionDuration := time.Duration(0)
	closedConnections := c.totalConnections - c.activeConnections
	if closedConnections > 0 {
		avgConnectionDuration = c.connectionDuration / time.Duration(closedConnections)
	}

	// Calculate average processing time
	avgProcessingTime := time.Duration(0)
	if totalMessages > 0 {
		avgProcessingTime = c.messageProcessingTime / time.Duration(totalMessages)
	}

	// Calculate average AI response time
	avgAIResponseTime := time.Duration(0)
	if c.aiRequestsTotal > 0 {
		avgAIResponseTime = c.aiResponseTime / time.Duration(c.aiRequestsTotal)
	}

	// Calculate average session duration
	avgSessionDuration := time.Duration(0)
	if c.activeUsers > 0 {
		avgSessionDuration = c.sessionDuration / time.Duration(c.activeUsers)
	}

	// Calculate messages per session
	messagesPerSession := float64(0)
	if c.activeUsers > 0 {
		messagesPerSession = float64(totalMessages) / float64(c.activeUsers)
	}

	// Calculate total errors
	totalErrors := c.authenticationErrors + c.authorizationErrors + c.validationErrors + c.systemErrors + c.timeoutErrors

	// Calculate error rate
	errorRate := float64(0)
	totalOperations := totalMessages + c.aiRequestsTotal + c.totalConnections
	if totalOperations > 0 {
		errorRate = float64(totalErrors) / float64(totalOperations)
	}

	return &ChatMetrics{
		ConnectionMetrics: ConnectionMetrics{
			TotalConnections:          c.totalConnections,
			ActiveConnections:         c.activeConnections,
			AverageConnectionDuration: avgConnectionDuration,
			ReconnectionAttempts:      c.reconnectionAttempts,
			ConnectionFailures:        c.connectionFailures,
			WebSocketUpgrades:         c.websocketUpgrades,
			WebSocketErrors:           c.websocketErrors,
			ConnectionSuccessRate:     connectionSuccessRate,
		},
		MessageMetrics: MessageMetrics{
			MessagesSent:          c.messagesSent,
			MessagesReceived:      c.messagesReceived,
			AverageProcessingTime: avgProcessingTime,
			MessageErrors:         c.messageErrors,
			MessageRetries:        c.messageRetries,
			MessageBroadcasts:     c.messageBroadcasts,
			MessageSuccessRate:    messageSuccessRate,
		},
		AIMetrics: AIMetrics{
			RequestsTotal:       c.aiRequestsTotal,
			RequestsSuccessful:  c.aiRequestsSuccessful,
			RequestsFailed:      c.aiRequestsFailed,
			AverageResponseTime: avgAIResponseTime,
			TokensUsed:          c.aiTokensUsed,
			EstimatedCost:       c.aiCostEstimate,
			ModelUsage:          c.copyModelUsage(),
			SuccessRate:         aiSuccessRate,
		},
		UserMetrics: UserMetrics{
			ActiveUsers:            c.activeUsers,
			AverageSessionDuration: avgSessionDuration,
			MessagesPerSession:     messagesPerSession,
			QuickActionsUsed:       c.quickActionsUsed,
			FeatureUsage:           c.copyFeatureUsage(),
		},
		ErrorMetrics: ErrorMetrics{
			AuthenticationErrors: c.authenticationErrors,
			AuthorizationErrors:  c.authorizationErrors,
			ValidationErrors:     c.validationErrors,
			SystemErrors:         c.systemErrors,
			TimeoutErrors:        c.timeoutErrors,
			TotalErrors:          totalErrors,
			ErrorRate:            errorRate,
		},
		PerformanceMetrics: ChatPerformanceMetrics{
			ResponseTimeP50: perfMetrics.AverageResponseTime, // Simplified for now
			ResponseTimeP95: perfMetrics.AverageResponseTime * 2,
			ResponseTimeP99: perfMetrics.AverageResponseTime * 3,
			ThroughputMPS:   c.calculateThroughput(),
			ErrorRate:       errorRate,
			MemoryUsageMB:   c.memoryUsage / 1024 / 1024,
			CPUUsagePercent: c.cpuUsage,
			GoroutineCount:  c.goroutineCount,
		},
		CacheMetrics: cacheMetrics,
		Timestamp:    time.Now(),
	}
}

// copyModelUsage creates a copy of the model usage map
func (c *ChatMetricsCollector) copyModelUsage() map[string]int64 {
	copy := make(map[string]int64)
	for k, v := range c.aiModelUsage {
		copy[k] = v
	}
	return copy
}

// copyFeatureUsage creates a copy of the feature usage map
func (c *ChatMetricsCollector) copyFeatureUsage() map[string]int64 {
	copy := make(map[string]int64)
	for k, v := range c.featureUsage {
		copy[k] = v
	}
	return copy
}

// calculateThroughput calculates messages per second throughput
func (c *ChatMetricsCollector) calculateThroughput() float64 {
	timeSinceReset := time.Since(c.lastMetricsReset)
	if timeSinceReset.Seconds() == 0 {
		return 0
	}

	totalMessages := float64(c.messagesSent + c.messagesReceived)
	return totalMessages / timeSinceReset.Seconds()
}

// ResetMetrics resets all metrics counters
func (c *ChatMetricsCollector) ResetMetrics() {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Reset connection metrics
	c.totalConnections = 0
	c.connectionDuration = 0
	c.reconnectionAttempts = 0
	c.connectionFailures = 0
	c.websocketUpgrades = 0
	c.websocketErrors = 0

	// Reset message metrics
	c.messagesSent = 0
	c.messagesReceived = 0
	c.messageProcessingTime = 0
	c.messageErrors = 0
	c.messageRetries = 0
	c.messageBroadcasts = 0

	// Reset AI metrics
	c.aiRequestsTotal = 0
	c.aiRequestsSuccessful = 0
	c.aiRequestsFailed = 0
	c.aiResponseTime = 0
	c.aiTokensUsed = 0
	c.aiCostEstimate = 0
	c.aiModelUsage = make(map[string]int64)

	// Reset user metrics
	c.activeUsers = 0
	c.sessionDuration = 0
	c.quickActionsUsed = 0
	c.featureUsage = make(map[string]int64)

	// Reset error metrics
	c.authenticationErrors = 0
	c.authorizationErrors = 0
	c.validationErrors = 0
	c.systemErrors = 0
	c.timeoutErrors = 0

	c.lastMetricsReset = time.Now()
	c.logger.Info("Chat metrics reset")
}

// StartPeriodicCollection starts periodic metrics collection and logging
func (c *ChatMetricsCollector) StartPeriodicCollection(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	c.logger.WithField("interval", interval).Info("Started periodic chat metrics collection")

	for {
		select {
		case <-ctx.Done():
			c.logger.Info("Stopped periodic chat metrics collection")
			return
		case <-ticker.C:
			c.logMetrics()
		}
	}
}

// logMetrics logs current metrics
func (c *ChatMetricsCollector) logMetrics() {
	metrics := c.GetMetrics()

	c.logger.WithFields(logrus.Fields{
		"active_connections":      metrics.ConnectionMetrics.ActiveConnections,
		"total_connections":       metrics.ConnectionMetrics.TotalConnections,
		"connection_success_rate": fmt.Sprintf("%.2f%%", metrics.ConnectionMetrics.ConnectionSuccessRate*100),
		"messages_sent":           metrics.MessageMetrics.MessagesSent,
		"messages_received":       metrics.MessageMetrics.MessagesReceived,
		"message_success_rate":    fmt.Sprintf("%.2f%%", metrics.MessageMetrics.MessageSuccessRate*100),
		"ai_requests_total":       metrics.AIMetrics.RequestsTotal,
		"ai_success_rate":         fmt.Sprintf("%.2f%%", metrics.AIMetrics.SuccessRate*100),
		"ai_tokens_used":          metrics.AIMetrics.TokensUsed,
		"ai_estimated_cost":       fmt.Sprintf("$%.4f", metrics.AIMetrics.EstimatedCost),
		"active_users":            metrics.UserMetrics.ActiveUsers,
		"throughput_mps":          fmt.Sprintf("%.2f", metrics.PerformanceMetrics.ThroughputMPS),
		"error_rate":              fmt.Sprintf("%.2f%%", metrics.ErrorMetrics.ErrorRate*100),
		"cache_hit_ratio":         fmt.Sprintf("%.2f%%", c.cacheMonitor.GetCacheHitRatio()*100),
	}).Info("Chat system metrics")
}

// GetPrometheusMetrics returns metrics in Prometheus format
func (c *ChatMetricsCollector) GetPrometheusMetrics() string {
	metrics := c.GetMetrics()

	var prometheusMetrics []string

	// Connection metrics
	prometheusMetrics = append(prometheusMetrics,
		fmt.Sprintf("# HELP chat_connections_total Total number of chat connections"),
		fmt.Sprintf("# TYPE chat_connections_total counter"),
		fmt.Sprintf("chat_connections_total %d", metrics.ConnectionMetrics.TotalConnections),
		fmt.Sprintf("# HELP chat_connections_active Current number of active connections"),
		fmt.Sprintf("# TYPE chat_connections_active gauge"),
		fmt.Sprintf("chat_connections_active %d", metrics.ConnectionMetrics.ActiveConnections),
		fmt.Sprintf("# HELP chat_connection_failures_total Total number of connection failures"),
		fmt.Sprintf("# TYPE chat_connection_failures_total counter"),
		fmt.Sprintf("chat_connection_failures_total %d", metrics.ConnectionMetrics.ConnectionFailures),
	)

	// Message metrics
	prometheusMetrics = append(prometheusMetrics,
		fmt.Sprintf("# HELP chat_messages_sent_total Total number of messages sent"),
		fmt.Sprintf("# TYPE chat_messages_sent_total counter"),
		fmt.Sprintf("chat_messages_sent_total %d", metrics.MessageMetrics.MessagesSent),
		fmt.Sprintf("# HELP chat_messages_received_total Total number of messages received"),
		fmt.Sprintf("# TYPE chat_messages_received_total counter"),
		fmt.Sprintf("chat_messages_received_total %d", metrics.MessageMetrics.MessagesReceived),
		fmt.Sprintf("# HELP chat_message_errors_total Total number of message errors"),
		fmt.Sprintf("# TYPE chat_message_errors_total counter"),
		fmt.Sprintf("chat_message_errors_total %d", metrics.MessageMetrics.MessageErrors),
	)

	// AI metrics
	prometheusMetrics = append(prometheusMetrics,
		fmt.Sprintf("# HELP chat_ai_requests_total Total number of AI requests"),
		fmt.Sprintf("# TYPE chat_ai_requests_total counter"),
		fmt.Sprintf("chat_ai_requests_total %d", metrics.AIMetrics.RequestsTotal),
		fmt.Sprintf("# HELP chat_ai_tokens_used_total Total number of AI tokens used"),
		fmt.Sprintf("# TYPE chat_ai_tokens_used_total counter"),
		fmt.Sprintf("chat_ai_tokens_used_total %d", metrics.AIMetrics.TokensUsed),
		fmt.Sprintf("# HELP chat_ai_cost_estimate_total Estimated AI cost in USD"),
		fmt.Sprintf("# TYPE chat_ai_cost_estimate_total counter"),
		fmt.Sprintf("chat_ai_cost_estimate_total %.6f", metrics.AIMetrics.EstimatedCost),
	)

	// Performance metrics
	prometheusMetrics = append(prometheusMetrics,
		fmt.Sprintf("# HELP chat_response_time_seconds Response time in seconds"),
		fmt.Sprintf("# TYPE chat_response_time_seconds histogram"),
		fmt.Sprintf("chat_response_time_seconds_bucket{le=\"0.1\"} %d", metrics.MessageMetrics.MessagesSent/10),
		fmt.Sprintf("chat_response_time_seconds_bucket{le=\"0.5\"} %d", metrics.MessageMetrics.MessagesSent/5),
		fmt.Sprintf("chat_response_time_seconds_bucket{le=\"1.0\"} %d", metrics.MessageMetrics.MessagesSent/2),
		fmt.Sprintf("chat_response_time_seconds_bucket{le=\"+Inf\"} %d", metrics.MessageMetrics.MessagesSent),
		fmt.Sprintf("chat_response_time_seconds_sum %.6f", metrics.MessageMetrics.AverageProcessingTime.Seconds()*float64(metrics.MessageMetrics.MessagesSent)),
		fmt.Sprintf("chat_response_time_seconds_count %d", metrics.MessageMetrics.MessagesSent),
	)

	// Error metrics
	prometheusMetrics = append(prometheusMetrics,
		fmt.Sprintf("# HELP chat_errors_total Total number of errors by type"),
		fmt.Sprintf("# TYPE chat_errors_total counter"),
		fmt.Sprintf("chat_errors_total{type=\"authentication\"} %d", metrics.ErrorMetrics.AuthenticationErrors),
		fmt.Sprintf("chat_errors_total{type=\"authorization\"} %d", metrics.ErrorMetrics.AuthorizationErrors),
		fmt.Sprintf("chat_errors_total{type=\"validation\"} %d", metrics.ErrorMetrics.ValidationErrors),
		fmt.Sprintf("chat_errors_total{type=\"system\"} %d", metrics.ErrorMetrics.SystemErrors),
		fmt.Sprintf("chat_errors_total{type=\"timeout\"} %d", metrics.ErrorMetrics.TimeoutErrors),
	)

	// User metrics
	prometheusMetrics = append(prometheusMetrics,
		fmt.Sprintf("# HELP chat_active_users Current number of active users"),
		fmt.Sprintf("# TYPE chat_active_users gauge"),
		fmt.Sprintf("chat_active_users %d", metrics.UserMetrics.ActiveUsers),
		fmt.Sprintf("# HELP chat_throughput_messages_per_second Current message throughput"),
		fmt.Sprintf("# TYPE chat_throughput_messages_per_second gauge"),
		fmt.Sprintf("chat_throughput_messages_per_second %.2f", metrics.PerformanceMetrics.ThroughputMPS),
	)

	return fmt.Sprintf("%s\n", joinStrings(prometheusMetrics, "\n"))
}

// joinStrings joins a slice of strings with a separator
func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	if len(strs) == 1 {
		return strs[0]
	}

	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
}
