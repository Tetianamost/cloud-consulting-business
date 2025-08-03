package services

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// ChatStructuredLogger provides structured logging with correlation IDs for chat system
type ChatStructuredLogger struct {
	logger *logrus.Logger
}

// LogLevel represents log levels
type LogLevel string

const (
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
	LogLevelFatal LogLevel = "fatal"
)

// LogContext represents structured log context
type LogContext struct {
	CorrelationID string                 `json:"correlation_id"`
	UserID        string                 `json:"user_id,omitempty"`
	SessionID     string                 `json:"session_id,omitempty"`
	ConnectionID  string                 `json:"connection_id,omitempty"`
	MessageID     string                 `json:"message_id,omitempty"`
	Component     string                 `json:"component"`
	Operation     string                 `json:"operation"`
	Duration      time.Duration          `json:"duration,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
	Timestamp     time.Time              `json:"timestamp"`
}

// SecurityEvent represents a security-related log event
type SecurityEvent struct {
	EventType     string                 `json:"event_type"`
	Severity      string                 `json:"severity"`
	UserID        string                 `json:"user_id,omitempty"`
	IPAddress     string                 `json:"ip_address,omitempty"`
	UserAgent     string                 `json:"user_agent,omitempty"`
	Details       map[string]interface{} `json:"details"`
	Timestamp     time.Time              `json:"timestamp"`
	CorrelationID string                 `json:"correlation_id"`
}

// PerformanceEvent represents a performance-related log event
type PerformanceEvent struct {
	Operation     string                 `json:"operation"`
	Duration      time.Duration          `json:"duration"`
	Success       bool                   `json:"success"`
	ErrorMessage  string                 `json:"error_message,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
	Timestamp     time.Time              `json:"timestamp"`
	CorrelationID string                 `json:"correlation_id"`
}

// ErrorEvent represents an error log event
type ErrorEvent struct {
	ErrorType     string                 `json:"error_type"`
	ErrorMessage  string                 `json:"error_message"`
	StackTrace    string                 `json:"stack_trace,omitempty"`
	Component     string                 `json:"component"`
	Operation     string                 `json:"operation"`
	UserID        string                 `json:"user_id,omitempty"`
	SessionID     string                 `json:"session_id,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
	Timestamp     time.Time              `json:"timestamp"`
	CorrelationID string                 `json:"correlation_id"`
}

// NewChatStructuredLogger creates a new structured logger
func NewChatStructuredLogger(logger *logrus.Logger) *ChatStructuredLogger {
	// Configure logger for structured logging
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339Nano,
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "timestamp",
			logrus.FieldKeyLevel: "level",
			logrus.FieldKeyMsg:   "message",
		},
	})

	return &ChatStructuredLogger{
		logger: logger,
	}
}

// GenerateCorrelationID generates a new correlation ID
func (l *ChatStructuredLogger) GenerateCorrelationID() string {
	return uuid.New().String()
}

// WithContext creates a log context with correlation ID
func (l *ChatStructuredLogger) WithContext(ctx context.Context) *LogContext {
	correlationID := l.getCorrelationIDFromContext(ctx)
	if correlationID == "" {
		correlationID = l.GenerateCorrelationID()
	}

	return &LogContext{
		CorrelationID: correlationID,
		Timestamp:     time.Now(),
		Metadata:      make(map[string]interface{}),
	}
}

// WithCorrelationID creates a log context with specific correlation ID
func (l *ChatStructuredLogger) WithCorrelationID(correlationID string) *LogContext {
	return &LogContext{
		CorrelationID: correlationID,
		Timestamp:     time.Now(),
		Metadata:      make(map[string]interface{}),
	}
}

// getCorrelationIDFromContext extracts correlation ID from context
func (l *ChatStructuredLogger) getCorrelationIDFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	if correlationID, ok := ctx.Value("correlation_id").(string); ok {
		return correlationID
	}

	return ""
}

// SetUserID sets the user ID in the log context
func (lc *LogContext) SetUserID(userID string) *LogContext {
	lc.UserID = userID
	return lc
}

// SetSessionID sets the session ID in the log context
func (lc *LogContext) SetSessionID(sessionID string) *LogContext {
	lc.SessionID = sessionID
	return lc
}

// SetConnectionID sets the connection ID in the log context
func (lc *LogContext) SetConnectionID(connectionID string) *LogContext {
	lc.ConnectionID = connectionID
	return lc
}

// SetMessageID sets the message ID in the log context
func (lc *LogContext) SetMessageID(messageID string) *LogContext {
	lc.MessageID = messageID
	return lc
}

// SetComponent sets the component name in the log context
func (lc *LogContext) SetComponent(component string) *LogContext {
	lc.Component = component
	return lc
}

// SetOperation sets the operation name in the log context
func (lc *LogContext) SetOperation(operation string) *LogContext {
	lc.Operation = operation
	return lc
}

// SetDuration sets the operation duration in the log context
func (lc *LogContext) SetDuration(duration time.Duration) *LogContext {
	lc.Duration = duration
	return lc
}

// AddMetadata adds metadata to the log context
func (lc *LogContext) AddMetadata(key string, value interface{}) *LogContext {
	if lc.Metadata == nil {
		lc.Metadata = make(map[string]interface{})
	}
	lc.Metadata[key] = value
	return lc
}

// toLogrusFields converts LogContext to logrus fields
func (lc *LogContext) toLogrusFields() logrus.Fields {
	fields := logrus.Fields{
		"correlation_id": lc.CorrelationID,
		"component":      lc.Component,
		"operation":      lc.Operation,
		"timestamp":      lc.Timestamp,
	}

	if lc.UserID != "" {
		fields["user_id"] = lc.UserID
	}
	if lc.SessionID != "" {
		fields["session_id"] = lc.SessionID
	}
	if lc.ConnectionID != "" {
		fields["connection_id"] = lc.ConnectionID
	}
	if lc.MessageID != "" {
		fields["message_id"] = lc.MessageID
	}
	if lc.Duration > 0 {
		fields["duration_ms"] = lc.Duration.Milliseconds()
	}

	// Add metadata
	for k, v := range lc.Metadata {
		fields[k] = v
	}

	return fields
}

// Debug logs a debug message
func (l *ChatStructuredLogger) Debug(logCtx *LogContext, message string) {
	l.logger.WithFields(logCtx.toLogrusFields()).Debug(message)
}

// Info logs an info message
func (l *ChatStructuredLogger) Info(logCtx *LogContext, message string) {
	l.logger.WithFields(logCtx.toLogrusFields()).Info(message)
}

// Warn logs a warning message
func (l *ChatStructuredLogger) Warn(logCtx *LogContext, message string) {
	l.logger.WithFields(logCtx.toLogrusFields()).Warn(message)
}

// Error logs an error message
func (l *ChatStructuredLogger) Error(logCtx *LogContext, message string, err error) {
	fields := logCtx.toLogrusFields()
	if err != nil {
		fields["error"] = err.Error()
		fields["error_type"] = fmt.Sprintf("%T", err)
	}
	l.logger.WithFields(fields).Error(message)
}

// Fatal logs a fatal message
func (l *ChatStructuredLogger) Fatal(logCtx *LogContext, message string, err error) {
	fields := logCtx.toLogrusFields()
	if err != nil {
		fields["error"] = err.Error()
		fields["error_type"] = fmt.Sprintf("%T", err)
	}
	l.logger.WithFields(fields).Fatal(message)
}

// LogSecurityEvent logs a security-related event
func (l *ChatStructuredLogger) LogSecurityEvent(event SecurityEvent) {
	fields := logrus.Fields{
		"event_category": "security",
		"event_type":     event.EventType,
		"severity":       event.Severity,
		"correlation_id": event.CorrelationID,
		"timestamp":      event.Timestamp,
	}

	if event.UserID != "" {
		fields["user_id"] = event.UserID
	}
	if event.IPAddress != "" {
		fields["ip_address"] = event.IPAddress
	}
	if event.UserAgent != "" {
		fields["user_agent"] = event.UserAgent
	}

	// Add event details
	for k, v := range event.Details {
		fields[k] = v
	}

	message := fmt.Sprintf("Security event: %s", event.EventType)

	switch event.Severity {
	case "critical", "high":
		l.logger.WithFields(fields).Error(message)
	case "medium":
		l.logger.WithFields(fields).Warn(message)
	default:
		l.logger.WithFields(fields).Info(message)
	}
}

// LogPerformanceEvent logs a performance-related event
func (l *ChatStructuredLogger) LogPerformanceEvent(event PerformanceEvent) {
	fields := logrus.Fields{
		"event_category": "performance",
		"operation":      event.Operation,
		"duration_ms":    event.Duration.Milliseconds(),
		"success":        event.Success,
		"correlation_id": event.CorrelationID,
		"timestamp":      event.Timestamp,
	}

	if event.ErrorMessage != "" {
		fields["error_message"] = event.ErrorMessage
	}

	// Add metadata
	for k, v := range event.Metadata {
		fields[k] = v
	}

	message := fmt.Sprintf("Performance event: %s (%.2fms)", event.Operation, float64(event.Duration.Nanoseconds())/1e6)

	if event.Success {
		l.logger.WithFields(fields).Info(message)
	} else {
		l.logger.WithFields(fields).Warn(message)
	}
}

// LogErrorEvent logs a detailed error event
func (l *ChatStructuredLogger) LogErrorEvent(event ErrorEvent) {
	fields := logrus.Fields{
		"event_category": "error",
		"error_type":     event.ErrorType,
		"error_message":  event.ErrorMessage,
		"component":      event.Component,
		"operation":      event.Operation,
		"correlation_id": event.CorrelationID,
		"timestamp":      event.Timestamp,
	}

	if event.StackTrace != "" {
		fields["stack_trace"] = event.StackTrace
	}
	if event.UserID != "" {
		fields["user_id"] = event.UserID
	}
	if event.SessionID != "" {
		fields["session_id"] = event.SessionID
	}

	// Add metadata
	for k, v := range event.Metadata {
		fields[k] = v
	}

	message := fmt.Sprintf("Error in %s.%s: %s", event.Component, event.Operation, event.ErrorMessage)
	l.logger.WithFields(fields).Error(message)
}

// LogConnectionEvent logs WebSocket connection events
func (l *ChatStructuredLogger) LogConnectionEvent(correlationID, userID, connectionID, eventType string, metadata map[string]interface{}) {
	logCtx := l.WithCorrelationID(correlationID).
		SetUserID(userID).
		SetConnectionID(connectionID).
		SetComponent("websocket").
		SetOperation("connection")

	for k, v := range metadata {
		logCtx.AddMetadata(k, v)
	}

	message := fmt.Sprintf("WebSocket connection %s", eventType)

	switch eventType {
	case "opened", "authenticated":
		l.Info(logCtx, message)
	case "closed":
		l.Info(logCtx, message)
	case "failed", "error":
		l.Error(logCtx, message, nil)
	default:
		l.Debug(logCtx, message)
	}
}

// LogMessageEvent logs chat message events
func (l *ChatStructuredLogger) LogMessageEvent(correlationID, userID, sessionID, messageID, eventType string, duration time.Duration, metadata map[string]interface{}) {
	logCtx := l.WithCorrelationID(correlationID).
		SetUserID(userID).
		SetSessionID(sessionID).
		SetMessageID(messageID).
		SetComponent("chat").
		SetOperation("message").
		SetDuration(duration)

	for k, v := range metadata {
		logCtx.AddMetadata(k, v)
	}

	message := fmt.Sprintf("Chat message %s", eventType)

	switch eventType {
	case "sent", "received", "processed":
		l.Info(logCtx, message)
	case "failed", "error":
		l.Error(logCtx, message, nil)
	case "retry":
		l.Warn(logCtx, message)
	default:
		l.Debug(logCtx, message)
	}
}

// LogAIEvent logs AI service events
func (l *ChatStructuredLogger) LogAIEvent(correlationID, userID, sessionID, model, eventType string, duration time.Duration, tokensUsed int64, cost float64, metadata map[string]interface{}) {
	logCtx := l.WithCorrelationID(correlationID).
		SetUserID(userID).
		SetSessionID(sessionID).
		SetComponent("ai").
		SetOperation("generate").
		SetDuration(duration).
		AddMetadata("model", model).
		AddMetadata("tokens_used", tokensUsed).
		AddMetadata("cost_usd", cost)

	for k, v := range metadata {
		logCtx.AddMetadata(k, v)
	}

	message := fmt.Sprintf("AI %s request using %s", eventType, model)

	switch eventType {
	case "success":
		l.Info(logCtx, message)
	case "failed", "error":
		l.Error(logCtx, message, nil)
	case "timeout":
		l.Warn(logCtx, message)
	default:
		l.Debug(logCtx, message)
	}
}

// getStackTrace returns the current stack trace
func getStackTrace() string {
	buf := make([]byte, 4096)
	n := runtime.Stack(buf, false)
	return string(buf[:n])
}

// CreateSecurityEvent creates a security event
func (l *ChatStructuredLogger) CreateSecurityEvent(eventType, severity, userID, ipAddress, userAgent, correlationID string, details map[string]interface{}) SecurityEvent {
	return SecurityEvent{
		EventType:     eventType,
		Severity:      severity,
		UserID:        userID,
		IPAddress:     ipAddress,
		UserAgent:     userAgent,
		Details:       details,
		Timestamp:     time.Now(),
		CorrelationID: correlationID,
	}
}

// CreatePerformanceEvent creates a performance event
func (l *ChatStructuredLogger) CreatePerformanceEvent(operation string, duration time.Duration, success bool, errorMessage, correlationID string, metadata map[string]interface{}) PerformanceEvent {
	return PerformanceEvent{
		Operation:     operation,
		Duration:      duration,
		Success:       success,
		ErrorMessage:  errorMessage,
		Metadata:      metadata,
		Timestamp:     time.Now(),
		CorrelationID: correlationID,
	}
}

// CreateErrorEvent creates an error event
func (l *ChatStructuredLogger) CreateErrorEvent(errorType, errorMessage, component, operation, userID, sessionID, correlationID string, err error, metadata map[string]interface{}) ErrorEvent {
	stackTrace := ""
	if err != nil {
		stackTrace = getStackTrace()
	}

	return ErrorEvent{
		ErrorType:     errorType,
		ErrorMessage:  errorMessage,
		StackTrace:    stackTrace,
		Component:     component,
		Operation:     operation,
		UserID:        userID,
		SessionID:     sessionID,
		Metadata:      metadata,
		Timestamp:     time.Now(),
		CorrelationID: correlationID,
	}
}

// LogAggregationQuery represents a query for log aggregation
type LogAggregationQuery struct {
	StartTime     time.Time              `json:"start_time"`
	EndTime       time.Time              `json:"end_time"`
	Level         LogLevel               `json:"level,omitempty"`
	Component     string                 `json:"component,omitempty"`
	Operation     string                 `json:"operation,omitempty"`
	UserID        string                 `json:"user_id,omitempty"`
	SessionID     string                 `json:"session_id,omitempty"`
	CorrelationID string                 `json:"correlation_id,omitempty"`
	EventCategory string                 `json:"event_category,omitempty"`
	Filters       map[string]interface{} `json:"filters,omitempty"`
	Limit         int                    `json:"limit"`
	Offset        int                    `json:"offset"`
}

// LogSearchResult represents search results
type LogSearchResult struct {
	Logs       []map[string]interface{} `json:"logs"`
	TotalCount int64                    `json:"total_count"`
	Query      LogAggregationQuery      `json:"query"`
	Timestamp  time.Time                `json:"timestamp"`
}

// SearchLogs provides a placeholder for log search functionality
// In a real implementation, this would query a log aggregation system like ELK stack
func (l *ChatStructuredLogger) SearchLogs(query LogAggregationQuery) (*LogSearchResult, error) {
	// This is a placeholder implementation
	// In production, this would integrate with Elasticsearch, Splunk, or similar

	result := &LogSearchResult{
		Logs:       []map[string]interface{}{},
		TotalCount: 0,
		Query:      query,
		Timestamp:  time.Now(),
	}

	// Log the search query for audit purposes
	logCtx := l.WithCorrelationID(l.GenerateCorrelationID()).
		SetComponent("logging").
		SetOperation("search").
		AddMetadata("query", query)

	l.Info(logCtx, "Log search query executed")

	return result, nil
}

// GetLogStatistics returns log statistics for monitoring
func (l *ChatStructuredLogger) GetLogStatistics(startTime, endTime time.Time) map[string]interface{} {
	// This is a placeholder implementation
	// In production, this would query log aggregation system for statistics

	stats := map[string]interface{}{
		"period": map[string]interface{}{
			"start": startTime,
			"end":   endTime,
		},
		"counts": map[string]int64{
			"total": 0,
			"debug": 0,
			"info":  0,
			"warn":  0,
			"error": 0,
			"fatal": 0,
		},
		"categories": map[string]int64{
			"security":    0,
			"performance": 0,
			"error":       0,
			"connection":  0,
			"message":     0,
			"ai":          0,
		},
		"top_errors": []map[string]interface{}{},
		"timestamp":  time.Now(),
	}

	return stats
}
