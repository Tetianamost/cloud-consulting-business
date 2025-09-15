package services

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/cloud-consulting/backend/internal/domain"
	"github.com/cloud-consulting/backend/internal/interfaces"
)

// EmailMetricsServiceImpl implements the EmailMetricsService interface
type EmailMetricsServiceImpl struct {
	emailEventRepo interfaces.EmailEventRepository
	logger         *logrus.Logger

	// Metrics tracking
	metrics     *interfaces.EmailMetricsServiceMetrics
	metricsLock sync.RWMutex
	startTime   time.Time
}

// NewEmailMetricsService creates a new email metrics service
func NewEmailMetricsService(emailEventRepo interfaces.EmailEventRepository, logger *logrus.Logger) interfaces.EmailMetricsService {
	return &EmailMetricsServiceImpl{
		emailEventRepo: emailEventRepo,
		logger:         logger,
		startTime:      time.Now(),
		metrics: &interfaces.EmailMetricsServiceMetrics{
			TotalMetricsRequests: 0,
			SuccessfulRequests:   0,
			FailedRequests:       0,
			SuccessRate:          0.0,
			AverageResponseTime:  0.0,
			LastRequestTime:      time.Time{},
			CacheHits:            0,
			CacheMisses:          0,
			HealthCheckFailures:  0,
		},
	}
}

// GetMetrics returns current metrics for the email metrics service
func (s *EmailMetricsServiceImpl) GetMetrics() *interfaces.EmailMetricsServiceMetrics {
	s.metricsLock.RLock()
	defer s.metricsLock.RUnlock()

	// Calculate success rate
	if s.metrics.TotalMetricsRequests > 0 {
		s.metrics.SuccessRate = float64(s.metrics.SuccessfulRequests) / float64(s.metrics.TotalMetricsRequests)
	}

	// Return a copy to prevent external modification
	metrics := *s.metrics
	return &metrics
}

// updateMetrics updates the internal metrics (should be called with metricsLock held)
func (s *EmailMetricsServiceImpl) updateMetrics(success bool, duration time.Duration) {
	atomic.AddInt64(&s.metrics.TotalMetricsRequests, 1)

	if success {
		atomic.AddInt64(&s.metrics.SuccessfulRequests, 1)
	} else {
		atomic.AddInt64(&s.metrics.FailedRequests, 1)
	}

	// Update average response time (simplified moving average)
	if s.metrics.AverageResponseTime == 0 {
		s.metrics.AverageResponseTime = float64(duration.Nanoseconds()) / 1e6 // Convert to milliseconds
	} else {
		// Simple exponential moving average
		alpha := 0.1
		newTime := float64(duration.Nanoseconds()) / 1e6
		s.metrics.AverageResponseTime = alpha*newTime + (1-alpha)*s.metrics.AverageResponseTime
	}

	s.metrics.LastRequestTime = time.Now()
}

// GetEmailMetrics calculates and returns email metrics for the specified time range
func (s *EmailMetricsServiceImpl) GetEmailMetrics(ctx context.Context, timeRange domain.TimeRange) (*domain.EmailMetrics, error) {
	requestStart := time.Now()

	s.logger.WithFields(logrus.Fields{
		"component":  "email_metrics_service",
		"start_time": timeRange.Start.Format(time.RFC3339),
		"end_time":   timeRange.End.Format(time.RFC3339),
		"action":     "get_email_metrics",
	}).Debug("Calculating email metrics for time range")

	// Create filters for the time range
	filters := domain.EmailEventFilters{
		TimeRange: &timeRange,
		Limit:     0, // No limit for metrics calculation
		Offset:    0,
	}

	// Get metrics from repository
	metrics, err := s.emailEventRepo.GetMetrics(ctx, filters)
	requestDuration := time.Since(requestStart)

	// Update internal metrics
	s.metricsLock.Lock()
	s.updateMetrics(err == nil, requestDuration)
	s.metricsLock.Unlock()

	if err != nil {
		s.logger.WithError(err).WithFields(logrus.Fields{
			"component":   "email_metrics_service",
			"start_time":  timeRange.Start.Format(time.RFC3339),
			"end_time":    timeRange.End.Format(time.RFC3339),
			"action":      "get_email_metrics_failed",
			"duration_ms": requestDuration.Nanoseconds() / 1e6,
		}).Error("Failed to calculate email metrics")
		return nil, fmt.Errorf("failed to calculate email metrics: %w", err)
	}

	s.logger.WithFields(logrus.Fields{
		"component":        "email_metrics_service",
		"total_emails":     metrics.TotalEmails,
		"delivered_emails": metrics.DeliveredEmails,
		"failed_emails":    metrics.FailedEmails,
		"delivery_rate":    metrics.DeliveryRate,
		"bounce_rate":      metrics.BounceRate,
		"spam_rate":        metrics.SpamRate,
		"time_range":       metrics.TimeRange,
		"action":           "get_email_metrics_success",
		"duration_ms":      requestDuration.Nanoseconds() / 1e6,
	}).Info("Successfully calculated email metrics")

	return metrics, nil
}

// GetEmailStatusByInquiry returns the email status for a specific inquiry
func (s *EmailMetricsServiceImpl) GetEmailStatusByInquiry(ctx context.Context, inquiryID string) (*domain.EmailStatus, error) {
	s.logger.WithField("inquiry_id", inquiryID).Debug("Getting email status for inquiry")

	// Get all email events for the inquiry
	events, err := s.emailEventRepo.GetByInquiryID(ctx, inquiryID)
	if err != nil {
		s.logger.WithError(err).WithField("inquiry_id", inquiryID).Error("Failed to get email events for inquiry")
		return nil, fmt.Errorf("failed to get email events for inquiry: %w", err)
	}

	if len(events) == 0 {
		s.logger.WithField("inquiry_id", inquiryID).Debug("No email events found for inquiry")
		return nil, nil
	}

	// Build email status from events
	emailStatus := &domain.EmailStatus{
		InquiryID:       inquiryID,
		TotalEmailsSent: len(events),
	}

	var lastEmailTime *time.Time

	// Categorize events by type and find the most recent of each type
	for _, event := range events {
		// Update last email sent time
		if lastEmailTime == nil || event.SentAt.After(*lastEmailTime) {
			lastEmailTime = &event.SentAt
		}

		// Categorize by email type and keep the most recent of each type
		switch event.EmailType {
		case domain.EmailTypeCustomerConfirmation:
			if emailStatus.CustomerEmail == nil || event.SentAt.After(emailStatus.CustomerEmail.SentAt) {
				emailStatus.CustomerEmail = event
			}
		case domain.EmailTypeConsultantNotification:
			if emailStatus.ConsultantEmail == nil || event.SentAt.After(emailStatus.ConsultantEmail.SentAt) {
				emailStatus.ConsultantEmail = event
			}
		case domain.EmailTypeInquiryNotification:
			if emailStatus.InquiryNotification == nil || event.SentAt.After(emailStatus.InquiryNotification.SentAt) {
				emailStatus.InquiryNotification = event
			}
		}
	}

	emailStatus.LastEmailSent = lastEmailTime

	s.logger.WithFields(logrus.Fields{
		"inquiry_id":               inquiryID,
		"total_emails_sent":        emailStatus.TotalEmailsSent,
		"has_customer_email":       emailStatus.CustomerEmail != nil,
		"has_consultant_email":     emailStatus.ConsultantEmail != nil,
		"has_inquiry_notification": emailStatus.InquiryNotification != nil,
		"last_email_sent":          lastEmailTime,
	}).Debug("Successfully built email status for inquiry")

	return emailStatus, nil
}

// GetEmailEventHistory returns email event history based on filters
func (s *EmailMetricsServiceImpl) GetEmailEventHistory(ctx context.Context, filters domain.EmailEventFilters) ([]*domain.EmailEvent, error) {
	s.logger.WithFields(logrus.Fields{
		"email_type": filters.EmailType,
		"status":     filters.Status,
		"inquiry_id": filters.InquiryID,
		"limit":      filters.Limit,
		"offset":     filters.Offset,
	}).Debug("Getting email event history")

	// Validate filters
	if err := s.validateFilters(filters); err != nil {
		s.logger.WithError(err).Error("Invalid filters provided")
		return nil, fmt.Errorf("invalid filters: %w", err)
	}

	// Get events from repository
	events, err := s.emailEventRepo.List(ctx, filters)
	if err != nil {
		s.logger.WithError(err).Error("Failed to get email event history")
		return nil, fmt.Errorf("failed to get email event history: %w", err)
	}

	s.logger.WithFields(logrus.Fields{
		"event_count": len(events),
		"limit":       filters.Limit,
		"offset":      filters.Offset,
	}).Debug("Successfully retrieved email event history")

	return events, nil
}

// validateFilters validates the email event filters
func (s *EmailMetricsServiceImpl) validateFilters(filters domain.EmailEventFilters) error {
	// Validate time range
	if filters.TimeRange != nil {
		if filters.TimeRange.Start.After(filters.TimeRange.End) {
			return fmt.Errorf("start time cannot be after end time")
		}

		// Check if time range is too far in the future
		now := time.Now()
		if filters.TimeRange.Start.After(now.Add(24 * time.Hour)) {
			return fmt.Errorf("start time cannot be more than 24 hours in the future")
		}
	}

	// Validate limit
	if filters.Limit < 0 {
		return fmt.Errorf("limit cannot be negative")
	}
	if filters.Limit > 1000 {
		return fmt.Errorf("limit cannot exceed 1000")
	}

	// Validate offset
	if filters.Offset < 0 {
		return fmt.Errorf("offset cannot be negative")
	}

	// Validate email type
	if filters.EmailType != nil {
		switch *filters.EmailType {
		case domain.EmailTypeCustomerConfirmation, domain.EmailTypeConsultantNotification, domain.EmailTypeInquiryNotification:
			// Valid email types
		default:
			return fmt.Errorf("invalid email type: %s", *filters.EmailType)
		}
	}

	// Validate status
	if filters.Status != nil {
		switch *filters.Status {
		case domain.EmailStatusSent, domain.EmailStatusDelivered, domain.EmailStatusFailed, domain.EmailStatusBounced, domain.EmailStatusSpam:
			// Valid statuses
		default:
			return fmt.Errorf("invalid email status: %s", *filters.Status)
		}
	}

	return nil
}

// GetEmailMetricsByType returns email metrics broken down by email type
func (s *EmailMetricsServiceImpl) GetEmailMetricsByType(ctx context.Context, timeRange domain.TimeRange) (map[domain.EmailEventType]*domain.EmailMetrics, error) {
	s.logger.WithFields(logrus.Fields{
		"start_time": timeRange.Start.Format(time.RFC3339),
		"end_time":   timeRange.End.Format(time.RFC3339),
	}).Debug("Calculating email metrics by type")

	emailTypes := []domain.EmailEventType{
		domain.EmailTypeCustomerConfirmation,
		domain.EmailTypeConsultantNotification,
		domain.EmailTypeInquiryNotification,
	}

	result := make(map[domain.EmailEventType]*domain.EmailMetrics)

	for _, emailType := range emailTypes {
		filters := domain.EmailEventFilters{
			TimeRange: &timeRange,
			EmailType: &emailType,
			Limit:     0,
			Offset:    0,
		}

		metrics, err := s.emailEventRepo.GetMetrics(ctx, filters)
		if err != nil {
			s.logger.WithError(err).WithFields(logrus.Fields{
				"email_type": emailType,
				"start_time": timeRange.Start.Format(time.RFC3339),
				"end_time":   timeRange.End.Format(time.RFC3339),
			}).Error("Failed to calculate email metrics for type")
			return nil, fmt.Errorf("failed to calculate email metrics for type %s: %w", emailType, err)
		}

		result[emailType] = metrics
	}

	s.logger.WithFields(logrus.Fields{
		"customer_confirmation_emails":   result[domain.EmailTypeCustomerConfirmation].TotalEmails,
		"consultant_notification_emails": result[domain.EmailTypeConsultantNotification].TotalEmails,
		"inquiry_notification_emails":    result[domain.EmailTypeInquiryNotification].TotalEmails,
	}).Info("Successfully calculated email metrics by type")

	return result, nil
}

// GetRecentEmailActivity returns recent email activity for monitoring
func (s *EmailMetricsServiceImpl) GetRecentEmailActivity(ctx context.Context, hours int) ([]*domain.EmailEvent, error) {
	if hours <= 0 || hours > 168 { // Max 1 week
		return nil, fmt.Errorf("hours must be between 1 and 168")
	}

	s.logger.WithField("hours", hours).Debug("Getting recent email activity")

	// Calculate time range
	now := time.Now()
	startTime := now.Add(-time.Duration(hours) * time.Hour)
	timeRange := domain.TimeRange{
		Start: startTime,
		End:   now,
	}

	filters := domain.EmailEventFilters{
		TimeRange: &timeRange,
		Limit:     100, // Limit recent activity to 100 events
		Offset:    0,
	}

	events, err := s.emailEventRepo.List(ctx, filters)
	if err != nil {
		s.logger.WithError(err).WithField("hours", hours).Error("Failed to get recent email activity")
		return nil, fmt.Errorf("failed to get recent email activity: %w", err)
	}

	s.logger.WithFields(logrus.Fields{
		"hours":       hours,
		"event_count": len(events),
	}).Debug("Successfully retrieved recent email activity")

	return events, nil
}

// IsHealthy checks if the email metrics service is healthy
func (s *EmailMetricsServiceImpl) IsHealthy(ctx context.Context) bool {
	healthCheckStart := time.Now()

	// Test basic functionality by trying to get metrics for the last hour
	now := time.Now()
	timeRange := domain.TimeRange{
		Start: now.Add(-1 * time.Hour),
		End:   now,
	}

	_, err := s.GetEmailMetrics(ctx, timeRange)
	healthCheckDuration := time.Since(healthCheckStart)

	if err != nil {
		s.metricsLock.Lock()
		atomic.AddInt64(&s.metrics.HealthCheckFailures, 1)
		s.metricsLock.Unlock()

		s.logger.WithError(err).WithFields(logrus.Fields{
			"component":                "email_metrics_service",
			"action":                   "health_check_failed",
			"health_check_duration_ms": healthCheckDuration.Nanoseconds() / 1e6,
		}).Error("Email metrics service health check failed")
		return false
	}

	s.logger.WithFields(logrus.Fields{
		"component":                "email_metrics_service",
		"action":                   "health_check_passed",
		"health_check_duration_ms": healthCheckDuration.Nanoseconds() / 1e6,
	}).Debug("Email metrics service health check passed")

	return true
}
