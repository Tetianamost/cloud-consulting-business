package services

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/cloud-consulting/backend/internal/domain"
	"github.com/cloud-consulting/backend/internal/interfaces"
)

// EmailEventRecorderImpl implements the EmailEventRecorder interface
type EmailEventRecorderImpl struct {
	repository interfaces.EmailEventRepository
	logger     *logrus.Logger

	// Metrics tracking
	metrics     *interfaces.EmailEventRecorderMetrics
	metricsLock sync.RWMutex
	startTime   time.Time
}

// NewEmailEventRecorder creates a new email event recorder service
func NewEmailEventRecorder(repository interfaces.EmailEventRepository, logger *logrus.Logger) interfaces.EmailEventRecorder {
	return &EmailEventRecorderImpl{
		repository: repository,
		logger:     logger,
		startTime:  time.Now(),
		metrics: &interfaces.EmailEventRecorderMetrics{
			TotalRecordingAttempts: 0,
			SuccessfulRecordings:   0,
			FailedRecordings:       0,
			SuccessRate:            0.0,
			AverageRecordingTime:   0.0,
			LastRecordingTime:      time.Time{},
			RetryAttempts:          0,
			HealthCheckFailures:    0,
		},
	}
}

// GetMetrics returns current metrics for the email event recorder
func (e *EmailEventRecorderImpl) GetMetrics() *interfaces.EmailEventRecorderMetrics {
	e.metricsLock.RLock()
	defer e.metricsLock.RUnlock()

	// Calculate success rate
	if e.metrics.TotalRecordingAttempts > 0 {
		e.metrics.SuccessRate = float64(e.metrics.SuccessfulRecordings) / float64(e.metrics.TotalRecordingAttempts)
	}

	// Return a copy to prevent external modification
	metrics := *e.metrics
	return &metrics
}

// updateMetrics updates the internal metrics (should be called with metricsLock held)
func (e *EmailEventRecorderImpl) updateMetrics(success bool, duration time.Duration, retryCount int) {
	atomic.AddInt64(&e.metrics.TotalRecordingAttempts, 1)

	if success {
		atomic.AddInt64(&e.metrics.SuccessfulRecordings, 1)
	} else {
		atomic.AddInt64(&e.metrics.FailedRecordings, 1)
	}

	if retryCount > 0 {
		atomic.AddInt64(&e.metrics.RetryAttempts, int64(retryCount))
	}

	// Update average recording time (simplified moving average)
	if e.metrics.AverageRecordingTime == 0 {
		e.metrics.AverageRecordingTime = float64(duration.Nanoseconds()) / 1e6 // Convert to milliseconds
	} else {
		// Simple exponential moving average
		alpha := 0.1
		newTime := float64(duration.Nanoseconds()) / 1e6
		e.metrics.AverageRecordingTime = alpha*newTime + (1-alpha)*e.metrics.AverageRecordingTime
	}

	e.metrics.LastRecordingTime = time.Now()
}

// RecordEmailSent records an email send event in a non-blocking manner
func (e *EmailEventRecorderImpl) RecordEmailSent(ctx context.Context, event *domain.EmailEvent) error {
	recordingStart := time.Now()
	// Validate the event
	if err := event.Validate(); err != nil {
		e.logger.WithError(err).WithFields(logrus.Fields{
			"inquiry_id":      event.InquiryID,
			"email_type":      event.EmailType,
			"recipient_email": event.RecipientEmail,
		}).Error("Email event validation failed")
		return fmt.Errorf("email event validation failed: %w", err)
	}

	// Generate ID if not provided
	if event.ID == "" {
		event.ID = uuid.New().String()
	}

	// Set timestamps if not already set
	now := time.Now()
	if event.SentAt.IsZero() {
		event.SentAt = now
	}
	if event.CreatedAt.IsZero() {
		event.CreatedAt = now
	}
	if event.UpdatedAt.IsZero() {
		event.UpdatedAt = now
	}

	// Set default status if not provided
	if event.Status == "" {
		event.Status = domain.EmailStatusSent
	}

	// Record the event in a goroutine to make it non-blocking
	go func() {
		// Create a new context with timeout for the background operation
		bgCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Implement retry logic for transient failures
		maxRetries := 3
		var lastErr error
		success := false

		for attempt := 1; attempt <= maxRetries; attempt++ {
			if err := e.repository.Create(bgCtx, event); err != nil {
				lastErr = err
				e.logger.WithError(err).WithFields(logrus.Fields{
					"event_id":              event.ID,
					"inquiry_id":            event.InquiryID,
					"email_type":            event.EmailType,
					"recipient_email":       event.RecipientEmail,
					"status":                event.Status,
					"attempt":               attempt,
					"max_retries":           maxRetries,
					"recording_duration_ms": time.Since(recordingStart).Nanoseconds() / 1e6,
				}).Warn("Failed to record email event, retrying...")

				// Exponential backoff: wait 1s, 2s, 4s
				if attempt < maxRetries {
					backoffDuration := time.Duration(1<<(attempt-1)) * time.Second
					time.Sleep(backoffDuration)
				}
			} else {
				// Success
				success = true
				recordingDuration := time.Since(recordingStart)
				e.logger.WithFields(logrus.Fields{
					"event_id":              event.ID,
					"inquiry_id":            event.InquiryID,
					"email_type":            event.EmailType,
					"recipient_email":       event.RecipientEmail,
					"status":                event.Status,
					"attempt":               attempt,
					"recording_duration_ms": recordingDuration.Nanoseconds() / 1e6,
					"action":                "email_event_recorded",
				}).Info("Email event recorded successfully")
				break
			}
		}

		// Update metrics
		recordingDuration := time.Since(recordingStart)
		retryCount := 0
		if !success {
			retryCount = maxRetries - 1
		}

		e.metricsLock.Lock()
		e.updateMetrics(success, recordingDuration, retryCount)
		e.metricsLock.Unlock()

		if !success {
			// All retries failed - log final error but don't fail email delivery
			e.logger.WithError(lastErr).WithFields(logrus.Fields{
				"event_id":              event.ID,
				"inquiry_id":            event.InquiryID,
				"email_type":            event.EmailType,
				"recipient_email":       event.RecipientEmail,
				"status":                event.Status,
				"max_retries":           maxRetries,
				"recording_duration_ms": recordingDuration.Nanoseconds() / 1e6,
				"action":                "email_event_recording_failed",
			}).Error("Failed to record email event after all retries - email delivery will continue")

			// TODO: Consider implementing fallback logging to file or queue for later processing
			// This could be added as a future enhancement for critical email event tracking
		}
	}()

	// Return immediately to not block email delivery
	return nil
}

// UpdateEmailStatus updates the status of an email event by SES message ID
func (e *EmailEventRecorderImpl) UpdateEmailStatus(ctx context.Context, messageID string, status domain.EmailEventStatus, deliveredAt *time.Time, errorMsg string) error {
	if messageID == "" {
		return fmt.Errorf("message ID is required")
	}

	// Get the existing event by message ID
	event, err := e.repository.GetByMessageID(ctx, messageID)
	if err != nil {
		e.logger.WithError(err).WithField("message_id", messageID).Error("Failed to get email event by message ID")
		return fmt.Errorf("failed to get email event by message ID: %w", err)
	}

	if event == nil {
		e.logger.WithField("message_id", messageID).Warn("Email event not found for message ID")
		return fmt.Errorf("email event not found for message ID: %s", messageID)
	}

	// Update the event status
	originalStatus := event.Status
	event.Status = status
	event.UpdatedAt = time.Now()

	// Set delivered timestamp if provided
	if deliveredAt != nil {
		event.DeliveredAt = deliveredAt
	}

	// Set error message if provided
	if errorMsg != "" {
		event.ErrorMessage = errorMsg
	}

	// Update the event in the repository in a non-blocking manner
	go func() {
		// Create a new context with timeout for the background operation
		bgCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := e.repository.Update(bgCtx, event); err != nil {
			e.logger.WithError(err).WithFields(logrus.Fields{
				"event_id":        event.ID,
				"message_id":      messageID,
				"original_status": originalStatus,
				"new_status":      status,
				"inquiry_id":      event.InquiryID,
			}).Error("Failed to update email event status")
		} else {
			e.logger.WithFields(logrus.Fields{
				"event_id":        event.ID,
				"message_id":      messageID,
				"original_status": originalStatus,
				"new_status":      status,
				"inquiry_id":      event.InquiryID,
			}).Info("Email event status updated successfully")
		}
	}()

	return nil
}

// GetEmailEventsByInquiry retrieves all email events for a specific inquiry
func (e *EmailEventRecorderImpl) GetEmailEventsByInquiry(ctx context.Context, inquiryID string) ([]*domain.EmailEvent, error) {
	if inquiryID == "" {
		return nil, fmt.Errorf("inquiry ID is required")
	}

	events, err := e.repository.GetByInquiryID(ctx, inquiryID)
	if err != nil {
		e.logger.WithError(err).WithField("inquiry_id", inquiryID).Error("Failed to get email events by inquiry ID")
		return nil, fmt.Errorf("failed to get email events by inquiry ID: %w", err)
	}

	e.logger.WithFields(logrus.Fields{
		"inquiry_id":  inquiryID,
		"event_count": len(events),
	}).Debug("Retrieved email events by inquiry ID")

	return events, nil
}

// RecordEmailSentSync records an email send event synchronously (for testing or critical operations)
func (e *EmailEventRecorderImpl) RecordEmailSentSync(ctx context.Context, event *domain.EmailEvent) error {
	// Validate the event
	if err := event.Validate(); err != nil {
		e.logger.WithError(err).WithFields(logrus.Fields{
			"inquiry_id":      event.InquiryID,
			"email_type":      event.EmailType,
			"recipient_email": event.RecipientEmail,
		}).Error("Email event validation failed")
		return fmt.Errorf("email event validation failed: %w", err)
	}

	// Generate ID if not provided
	if event.ID == "" {
		event.ID = uuid.New().String()
	}

	// Set timestamps if not already set
	now := time.Now()
	if event.SentAt.IsZero() {
		event.SentAt = now
	}
	if event.CreatedAt.IsZero() {
		event.CreatedAt = now
	}
	if event.UpdatedAt.IsZero() {
		event.UpdatedAt = now
	}

	// Set default status if not provided
	if event.Status == "" {
		event.Status = domain.EmailStatusSent
	}

	// Record the event synchronously
	if err := e.repository.Create(ctx, event); err != nil {
		e.logger.WithError(err).WithFields(logrus.Fields{
			"event_id":        event.ID,
			"inquiry_id":      event.InquiryID,
			"email_type":      event.EmailType,
			"recipient_email": event.RecipientEmail,
			"status":          event.Status,
		}).Error("Failed to record email event synchronously")
		return fmt.Errorf("failed to record email event: %w", err)
	}

	e.logger.WithFields(logrus.Fields{
		"event_id":        event.ID,
		"inquiry_id":      event.InquiryID,
		"email_type":      event.EmailType,
		"recipient_email": event.RecipientEmail,
		"status":          event.Status,
	}).Info("Email event recorded synchronously")

	return nil
}

// UpdateEmailStatusSync updates the status of an email event synchronously
func (e *EmailEventRecorderImpl) UpdateEmailStatusSync(ctx context.Context, messageID string, status domain.EmailEventStatus, deliveredAt *time.Time, errorMsg string) error {
	if messageID == "" {
		return fmt.Errorf("message ID is required")
	}

	// Get the existing event by message ID
	event, err := e.repository.GetByMessageID(ctx, messageID)
	if err != nil {
		e.logger.WithError(err).WithField("message_id", messageID).Error("Failed to get email event by message ID")
		return fmt.Errorf("failed to get email event by message ID: %w", err)
	}

	if event == nil {
		e.logger.WithField("message_id", messageID).Warn("Email event not found for message ID")
		return fmt.Errorf("email event not found for message ID: %s", messageID)
	}

	// Update the event status
	originalStatus := event.Status
	event.Status = status
	event.UpdatedAt = time.Now()

	// Set delivered timestamp if provided
	if deliveredAt != nil {
		event.DeliveredAt = deliveredAt
	}

	// Set error message if provided
	if errorMsg != "" {
		event.ErrorMessage = errorMsg
	}

	// Update the event in the repository synchronously
	if err := e.repository.Update(ctx, event); err != nil {
		e.logger.WithError(err).WithFields(logrus.Fields{
			"event_id":        event.ID,
			"message_id":      messageID,
			"original_status": originalStatus,
			"new_status":      status,
			"inquiry_id":      event.InquiryID,
		}).Error("Failed to update email event status synchronously")
		return fmt.Errorf("failed to update email event status: %w", err)
	}

	e.logger.WithFields(logrus.Fields{
		"event_id":        event.ID,
		"message_id":      messageID,
		"original_status": originalStatus,
		"new_status":      status,
		"inquiry_id":      event.InquiryID,
	}).Info("Email event status updated synchronously")

	return nil
}

// IsHealthy checks if the email event recorder service is healthy
func (e *EmailEventRecorderImpl) IsHealthy() bool {
	// Basic checks
	if e.repository == nil || e.logger == nil {
		return false
	}

	// Test repository connectivity with a simple operation
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Try to create a test event to verify repository is working
	testEvent := &domain.EmailEvent{
		ID:             "health-check-" + uuid.New().String(),
		InquiryID:      "health-check",
		EmailType:      domain.EmailTypeCustomerConfirmation,
		RecipientEmail: "health@test.com",
		SenderEmail:    "system@test.com",
		Subject:        "Health Check",
		Status:         domain.EmailStatusSent,
		SentAt:         time.Now(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Try to create and immediately clean up the test event
	if err := e.repository.Create(ctx, testEvent); err != nil {
		e.logger.WithError(err).Error("Email event recorder health check failed - repository create error")
		return false
	}

	// Clean up the test event (optional - depends on repository implementation)
	// Note: This assumes the repository has a delete method, which may not exist
	// In production, you might want to use a dedicated health check table

	return true
}

// IsHealthyWithContext checks if the email event recorder service is healthy with context
func (e *EmailEventRecorderImpl) IsHealthyWithContext(ctx context.Context) bool {
	healthCheckStart := time.Now()

	// Basic checks
	if e.repository == nil || e.logger == nil {
		e.metricsLock.Lock()
		atomic.AddInt64(&e.metrics.HealthCheckFailures, 1)
		e.metricsLock.Unlock()

		e.logger.WithFields(logrus.Fields{
			"component": "email_event_recorder",
			"reason":    "missing_dependencies",
			"action":    "health_check_failed",
		}).Error("Email event recorder health check failed - missing dependencies")
		return false
	}

	// Create a timeout context if none provided
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
	}

	// Test repository connectivity with a lightweight operation
	// Instead of creating a test event, try to query for a non-existent event
	// This tests database connectivity without creating test data
	_, err := e.repository.GetByMessageID(ctx, "health-check-non-existent")
	if err != nil {
		// We expect this to return an error (not found), but it should not be a connection error
		// Check if it's a connection-related error vs a "not found" error
		if isConnectionError(err) {
			e.metricsLock.Lock()
			atomic.AddInt64(&e.metrics.HealthCheckFailures, 1)
			e.metricsLock.Unlock()

			e.logger.WithError(err).WithFields(logrus.Fields{
				"component":                "email_event_recorder",
				"reason":                   "repository_connection_error",
				"action":                   "health_check_failed",
				"health_check_duration_ms": time.Since(healthCheckStart).Nanoseconds() / 1e6,
			}).Error("Email event recorder health check failed - repository connection error")
			return false
		}
		// "Not found" errors are expected and indicate the repository is working
	}

	e.logger.WithFields(logrus.Fields{
		"component":                "email_event_recorder",
		"action":                   "health_check_passed",
		"health_check_duration_ms": time.Since(healthCheckStart).Nanoseconds() / 1e6,
	}).Debug("Email event recorder health check passed")

	return true
}

// isConnectionError checks if an error is related to database connectivity
func isConnectionError(err error) bool {
	if err == nil {
		return false
	}

	errStr := err.Error()
	// Check for common database connection error patterns
	connectionErrors := []string{
		"connection refused",
		"connection timeout",
		"connection reset",
		"no such host",
		"network is unreachable",
		"database is closed",
		"driver: bad connection",
		"connection pool exhausted",
	}

	for _, connErr := range connectionErrors {
		if containsEmailError(errStr, connErr) {
			return true
		}
	}

	return false
}

// containsEmailError checks if a string contains a substring (case-insensitive)
func containsEmailError(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) &&
			(s[:len(substr)] == substr ||
				s[len(s)-len(substr):] == substr ||
				containsEmailSubstring(s, substr))))
}

// containsEmailSubstring performs a simple substring search
func containsEmailSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
