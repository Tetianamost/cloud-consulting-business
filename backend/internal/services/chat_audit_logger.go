package services

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/cloud-consulting/backend/internal/interfaces"
)

// ChatAuditLogger implements the ChatAuditLogger interface
type ChatAuditLogger struct {
	logger    *logrus.Logger
	auditLogs []*interfaces.AuditLog
	mutex     sync.RWMutex
	maxLogs   int
}

// NewChatAuditLogger creates a new chat audit logger
func NewChatAuditLogger(logger *logrus.Logger) *ChatAuditLogger {
	return &ChatAuditLogger{
		logger:    logger,
		auditLogs: make([]*interfaces.AuditLog, 0),
		maxLogs:   10000, // Keep last 10k audit logs in memory
	}
}

// LogLogin logs a login event
func (a *ChatAuditLogger) LogLogin(ctx context.Context, userID string, success bool, metadata map[string]interface{}) error {
	result := "success"
	if !success {
		result = "failure"
	}

	auditLog := &interfaces.AuditLog{
		ID:        uuid.New().String(),
		UserID:    userID,
		Action:    "login",
		Resource:  "authentication",
		Result:    result,
		Timestamp: time.Now(),
		Metadata:  metadata,
	}

	// Extract IP and User Agent from metadata if available
	if ip, ok := metadata["ip_address"].(string); ok {
		auditLog.IPAddress = ip
	}
	if ua, ok := metadata["user_agent"].(string); ok {
		auditLog.UserAgent = ua
	}

	a.addAuditLog(auditLog)

	a.logger.WithFields(logrus.Fields{
		"audit_id":   auditLog.ID,
		"user_id":    userID,
		"action":     "login",
		"result":     result,
		"ip_address": auditLog.IPAddress,
	}).Info("Login attempt logged")

	return nil
}

// LogLogout logs a logout event
func (a *ChatAuditLogger) LogLogout(ctx context.Context, userID string, metadata map[string]interface{}) error {
	auditLog := &interfaces.AuditLog{
		ID:        uuid.New().String(),
		UserID:    userID,
		Action:    "logout",
		Resource:  "authentication",
		Result:    "success",
		Timestamp: time.Now(),
		Metadata:  metadata,
	}

	// Extract IP and User Agent from metadata if available
	if ip, ok := metadata["ip_address"].(string); ok {
		auditLog.IPAddress = ip
	}
	if ua, ok := metadata["user_agent"].(string); ok {
		auditLog.UserAgent = ua
	}

	a.addAuditLog(auditLog)

	a.logger.WithFields(logrus.Fields{
		"audit_id": auditLog.ID,
		"user_id":  userID,
		"action":   "logout",
	}).Info("Logout logged")

	return nil
}

// LogTokenRefresh logs a token refresh event
func (a *ChatAuditLogger) LogTokenRefresh(ctx context.Context, userID string, metadata map[string]interface{}) error {
	auditLog := &interfaces.AuditLog{
		ID:        uuid.New().String(),
		UserID:    userID,
		Action:    "token_refresh",
		Resource:  "authentication",
		Result:    "success",
		Timestamp: time.Now(),
		Metadata:  metadata,
	}

	// Extract IP and User Agent from metadata if available
	if ip, ok := metadata["ip_address"].(string); ok {
		auditLog.IPAddress = ip
	}
	if ua, ok := metadata["user_agent"].(string); ok {
		auditLog.UserAgent = ua
	}

	a.addAuditLog(auditLog)

	a.logger.WithFields(logrus.Fields{
		"audit_id": auditLog.ID,
		"user_id":  userID,
		"action":   "token_refresh",
	}).Debug("Token refresh logged")

	return nil
}

// LogSessionCreated logs a session creation event
func (a *ChatAuditLogger) LogSessionCreated(ctx context.Context, userID, sessionID string, metadata map[string]interface{}) error {
	auditLog := &interfaces.AuditLog{
		ID:        uuid.New().String(),
		UserID:    userID,
		Action:    "session_created",
		Resource:  "chat_session",
		Result:    "success",
		Timestamp: time.Now(),
		Metadata:  metadata,
		SessionID: sessionID,
	}

	// Extract IP and User Agent from metadata if available
	if ip, ok := metadata["ip_address"].(string); ok {
		auditLog.IPAddress = ip
	}
	if ua, ok := metadata["user_agent"].(string); ok {
		auditLog.UserAgent = ua
	}

	a.addAuditLog(auditLog)

	a.logger.WithFields(logrus.Fields{
		"audit_id":   auditLog.ID,
		"user_id":    userID,
		"session_id": sessionID,
		"action":     "session_created",
	}).Info("Session creation logged")

	return nil
}

// LogSessionAccessed logs a session access event
func (a *ChatAuditLogger) LogSessionAccessed(ctx context.Context, userID, sessionID string, metadata map[string]interface{}) error {
	auditLog := &interfaces.AuditLog{
		ID:        uuid.New().String(),
		UserID:    userID,
		Action:    "session_accessed",
		Resource:  "chat_session",
		Result:    "success",
		Timestamp: time.Now(),
		Metadata:  metadata,
		SessionID: sessionID,
	}

	// Extract IP and User Agent from metadata if available
	if ip, ok := metadata["ip_address"].(string); ok {
		auditLog.IPAddress = ip
	}
	if ua, ok := metadata["user_agent"].(string); ok {
		auditLog.UserAgent = ua
	}

	a.addAuditLog(auditLog)

	a.logger.WithFields(logrus.Fields{
		"audit_id":   auditLog.ID,
		"user_id":    userID,
		"session_id": sessionID,
		"action":     "session_accessed",
	}).Debug("Session access logged")

	return nil
}

// LogSessionDeleted logs a session deletion event
func (a *ChatAuditLogger) LogSessionDeleted(ctx context.Context, userID, sessionID string, metadata map[string]interface{}) error {
	auditLog := &interfaces.AuditLog{
		ID:        uuid.New().String(),
		UserID:    userID,
		Action:    "session_deleted",
		Resource:  "chat_session",
		Result:    "success",
		Timestamp: time.Now(),
		Metadata:  metadata,
		SessionID: sessionID,
	}

	// Extract IP and User Agent from metadata if available
	if ip, ok := metadata["ip_address"].(string); ok {
		auditLog.IPAddress = ip
	}
	if ua, ok := metadata["user_agent"].(string); ok {
		auditLog.UserAgent = ua
	}

	a.addAuditLog(auditLog)

	a.logger.WithFields(logrus.Fields{
		"audit_id":   auditLog.ID,
		"user_id":    userID,
		"session_id": sessionID,
		"action":     "session_deleted",
	}).Info("Session deletion logged")

	return nil
}

// LogMessageSent logs a message sent event
func (a *ChatAuditLogger) LogMessageSent(ctx context.Context, userID, sessionID, messageID string, metadata map[string]interface{}) error {
	auditLog := &interfaces.AuditLog{
		ID:        uuid.New().String(),
		UserID:    userID,
		Action:    "message_sent",
		Resource:  "chat_message",
		Result:    "success",
		Timestamp: time.Now(),
		Metadata:  metadata,
		SessionID: sessionID,
		MessageID: messageID,
	}

	// Extract IP and User Agent from metadata if available
	if ip, ok := metadata["ip_address"].(string); ok {
		auditLog.IPAddress = ip
	}
	if ua, ok := metadata["user_agent"].(string); ok {
		auditLog.UserAgent = ua
	}

	a.addAuditLog(auditLog)

	a.logger.WithFields(logrus.Fields{
		"audit_id":   auditLog.ID,
		"user_id":    userID,
		"session_id": sessionID,
		"message_id": messageID,
		"action":     "message_sent",
	}).Debug("Message sent logged")

	return nil
}

// LogMessage logs a generic message event
func (a *ChatAuditLogger) LogMessage(ctx context.Context, userID, sessionID, action string, metadata map[string]interface{}) error {
	auditLog := &interfaces.AuditLog{
		ID:        uuid.New().String(),
		UserID:    userID,
		Action:    action,
		Resource:  "chat_message",
		Result:    "success",
		Timestamp: time.Now(),
		Metadata:  metadata,
		SessionID: sessionID,
	}

	// Extract message ID from metadata if available
	if messageID, ok := metadata["message_id"].(string); ok {
		auditLog.MessageID = messageID
	}

	// Extract IP and User Agent from metadata if available
	if ip, ok := metadata["ip_address"].(string); ok {
		auditLog.IPAddress = ip
	}
	if ua, ok := metadata["user_agent"].(string); ok {
		auditLog.UserAgent = ua
	}

	a.addAuditLog(auditLog)

	a.logger.WithFields(logrus.Fields{
		"audit_id":   auditLog.ID,
		"user_id":    userID,
		"session_id": sessionID,
		"action":     action,
	}).Debug("Message event logged")

	return nil
}

// LogMessageAccessed logs a message access event
func (a *ChatAuditLogger) LogMessageAccessed(ctx context.Context, userID, messageID string, metadata map[string]interface{}) error {
	auditLog := &interfaces.AuditLog{
		ID:        uuid.New().String(),
		UserID:    userID,
		Action:    "message_accessed",
		Resource:  "chat_message",
		Result:    "success",
		Timestamp: time.Now(),
		Metadata:  metadata,
		MessageID: messageID,
	}

	// Extract IP and User Agent from metadata if available
	if ip, ok := metadata["ip_address"].(string); ok {
		auditLog.IPAddress = ip
	}
	if ua, ok := metadata["user_agent"].(string); ok {
		auditLog.UserAgent = ua
	}

	a.addAuditLog(auditLog)

	a.logger.WithFields(logrus.Fields{
		"audit_id":   auditLog.ID,
		"user_id":    userID,
		"message_id": messageID,
		"action":     "message_accessed",
	}).Debug("Message access logged")

	return nil
}

// LogSecurityViolation logs a security violation
func (a *ChatAuditLogger) LogSecurityViolation(ctx context.Context, userID string, violation *interfaces.SecurityViolation) error {
	auditLog := &interfaces.AuditLog{
		ID:        uuid.New().String(),
		UserID:    userID,
		Action:    "security_violation",
		Resource:  violation.Resource,
		Result:    "violation",
		IPAddress: violation.IPAddress,
		UserAgent: violation.UserAgent,
		Timestamp: time.Now(),
		Metadata: map[string]interface{}{
			"violation_type": violation.Type,
			"severity":       violation.Severity,
			"description":    violation.Description,
			"action":         violation.Action,
		},
	}

	// Merge violation metadata
	for k, v := range violation.Metadata {
		auditLog.Metadata[k] = v
	}

	a.addAuditLog(auditLog)

	a.logger.WithFields(logrus.Fields{
		"audit_id":       auditLog.ID,
		"user_id":        userID,
		"violation_type": violation.Type,
		"severity":       violation.Severity,
		"resource":       violation.Resource,
		"action":         violation.Action,
		"ip_address":     violation.IPAddress,
	}).Warn("Security violation logged")

	return nil
}

// LogRateLimitExceeded logs a rate limit exceeded event
func (a *ChatAuditLogger) LogRateLimitExceeded(ctx context.Context, userID string, action string, metadata map[string]interface{}) error {
	auditLog := &interfaces.AuditLog{
		ID:        uuid.New().String(),
		UserID:    userID,
		Action:    "rate_limit_exceeded",
		Resource:  "rate_limiter",
		Result:    "blocked",
		Timestamp: time.Now(),
		Metadata: map[string]interface{}{
			"blocked_action": action,
		},
	}

	// Merge provided metadata
	for k, v := range metadata {
		auditLog.Metadata[k] = v
	}

	// Extract IP and User Agent from metadata if available
	if ip, ok := metadata["ip_address"].(string); ok {
		auditLog.IPAddress = ip
	}
	if ua, ok := metadata["user_agent"].(string); ok {
		auditLog.UserAgent = ua
	}

	a.addAuditLog(auditLog)

	a.logger.WithFields(logrus.Fields{
		"audit_id":       auditLog.ID,
		"user_id":        userID,
		"blocked_action": action,
		"ip_address":     auditLog.IPAddress,
	}).Warn("Rate limit exceeded logged")

	return nil
}

// LogUnauthorizedAccess logs an unauthorized access attempt
func (a *ChatAuditLogger) LogUnauthorizedAccess(ctx context.Context, userID string, resource string, metadata map[string]interface{}) error {
	auditLog := &interfaces.AuditLog{
		ID:        uuid.New().String(),
		UserID:    userID,
		Action:    "unauthorized_access",
		Resource:  resource,
		Result:    "denied",
		Timestamp: time.Now(),
		Metadata:  metadata,
	}

	// Extract IP and User Agent from metadata if available
	if ip, ok := metadata["ip_address"].(string); ok {
		auditLog.IPAddress = ip
	}
	if ua, ok := metadata["user_agent"].(string); ok {
		auditLog.UserAgent = ua
	}

	a.addAuditLog(auditLog)

	a.logger.WithFields(logrus.Fields{
		"audit_id":   auditLog.ID,
		"user_id":    userID,
		"resource":   resource,
		"ip_address": auditLog.IPAddress,
	}).Warn("Unauthorized access logged")

	return nil
}

// GetAuditLogs retrieves audit logs based on filters
func (a *ChatAuditLogger) GetAuditLogs(ctx context.Context, filters *interfaces.AuditLogFilters) ([]*interfaces.AuditLog, error) {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	var filteredLogs []*interfaces.AuditLog

	for _, log := range a.auditLogs {
		if a.matchesFilters(log, filters) {
			filteredLogs = append(filteredLogs, log)
		}
	}

	// Apply limit and offset
	if filters.Offset > 0 && filters.Offset < len(filteredLogs) {
		filteredLogs = filteredLogs[filters.Offset:]
	}

	if filters.Limit > 0 && filters.Limit < len(filteredLogs) {
		filteredLogs = filteredLogs[:filters.Limit]
	}

	return filteredLogs, nil
}

// GetUserAuditLogs retrieves audit logs for a specific user
func (a *ChatAuditLogger) GetUserAuditLogs(ctx context.Context, userID string, limit int) ([]*interfaces.AuditLog, error) {
	filters := &interfaces.AuditLogFilters{
		UserID: userID,
		Limit:  limit,
	}

	return a.GetAuditLogs(ctx, filters)
}

// Helper methods

// addAuditLog adds an audit log entry
func (a *ChatAuditLogger) addAuditLog(log *interfaces.AuditLog) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	// Add the log
	a.auditLogs = append(a.auditLogs, log)

	// Trim logs if we exceed the maximum
	if len(a.auditLogs) > a.maxLogs {
		// Remove the oldest logs (keep the most recent maxLogs entries)
		a.auditLogs = a.auditLogs[len(a.auditLogs)-a.maxLogs:]
	}
}

// matchesFilters checks if an audit log matches the given filters
func (a *ChatAuditLogger) matchesFilters(log *interfaces.AuditLog, filters *interfaces.AuditLogFilters) bool {
	if filters == nil {
		return true
	}

	// Check user ID filter
	if filters.UserID != "" && log.UserID != filters.UserID {
		return false
	}

	// Check action filter
	if filters.Action != "" && log.Action != filters.Action {
		return false
	}

	// Check resource filter
	if filters.Resource != "" && log.Resource != filters.Resource {
		return false
	}

	// Check time range filters
	if !filters.StartTime.IsZero() && log.Timestamp.Before(filters.StartTime) {
		return false
	}

	if !filters.EndTime.IsZero() && log.Timestamp.After(filters.EndTime) {
		return false
	}

	return true
}

// GetAuditLogCount returns the total number of audit logs
func (a *ChatAuditLogger) GetAuditLogCount() int {
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	return len(a.auditLogs)
}

// ClearAuditLogs clears all audit logs (for testing purposes)
func (a *ChatAuditLogger) ClearAuditLogs() {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.auditLogs = make([]*interfaces.AuditLog, 0)
	a.logger.Info("Audit logs cleared")
}

// SetMaxLogs sets the maximum number of logs to keep in memory
func (a *ChatAuditLogger) SetMaxLogs(maxLogs int) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	a.maxLogs = maxLogs

	// Trim existing logs if necessary
	if len(a.auditLogs) > maxLogs {
		a.auditLogs = a.auditLogs[len(a.auditLogs)-maxLogs:]
	}

	a.logger.WithField("max_logs", maxLogs).Info("Maximum audit logs updated")
}
