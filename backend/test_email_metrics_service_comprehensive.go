package main

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/cloud-consulting/backend/internal/domain"
	"github.com/cloud-consulting/backend/internal/interfaces"
	"github.com/cloud-consulting/backend/internal/services"
)

// MockEmailEventRepositoryForMetrics provides a mock implementation for testing email metrics service
type MockEmailEventRepositoryForMetrics struct {
	mock.Mock
	events []*domain.EmailEvent
}

func NewMockEmailEventRepositoryForMetrics() *MockEmailEventRepositoryForMetrics {
	return &MockEmailEventRepositoryForMetrics{
		events: make([]*domain.EmailEvent, 0),
	}
}

func (m *MockEmailEventRepositoryForMetrics) Create(ctx context.Context, event *domain.EmailEvent) error {
	args := m.Called(ctx, event)
	if args.Error(0) == nil {
		m.events = append(m.events, event)
	}
	return args.Error(0)
}

func (m *MockEmailEventRepositoryForMetrics) Update(ctx context.Context, event *domain.EmailEvent) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockEmailEventRepositoryForMetrics) GetByInquiryID(ctx context.Context, inquiryID string) ([]*domain.EmailEvent, error) {
	args := m.Called(ctx, inquiryID)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	var events []*domain.EmailEvent
	for _, event := range m.events {
		if event.InquiryID == inquiryID {
			events = append(events, event)
		}
	}

	return events, nil
}

func (m *MockEmailEventRepositoryForMetrics) GetByMessageID(ctx context.Context, messageID string) (*domain.EmailEvent, error) {
	args := m.Called(ctx, messageID)
	return args.Get(0).(*domain.EmailEvent), args.Error(1)
}

func (m *MockEmailEventRepositoryForMetrics) GetMetrics(ctx context.Context, filters domain.EmailEventFilters) (*domain.EmailMetrics, error) {
	args := m.Called(ctx, filters)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	// Calculate metrics from stored events
	metrics := &domain.EmailMetrics{}

	for _, event := range m.events {
		// Apply filters
		if filters.TimeRange != nil {
			if event.SentAt.Before(filters.TimeRange.Start) || event.SentAt.After(filters.TimeRange.End) {
				continue
			}
		}
		if filters.EmailType != nil && event.EmailType != *filters.EmailType {
			continue
		}
		if filters.Status != nil && event.Status != *filters.Status {
			continue
		}
		if filters.InquiryID != nil && event.InquiryID != *filters.InquiryID {
			continue
		}

		metrics.TotalEmails++
		switch event.Status {
		case domain.EmailStatusDelivered:
			metrics.DeliveredEmails++
		case domain.EmailStatusFailed:
			metrics.FailedEmails++
		case domain.EmailStatusBounced:
			metrics.BouncedEmails++
		case domain.EmailStatusSpam:
			metrics.SpamEmails++
		}
	}

	// Calculate rates
	if metrics.TotalEmails > 0 {
		metrics.DeliveryRate = float64(metrics.DeliveredEmails) / float64(metrics.TotalEmails) * 100
		metrics.BounceRate = float64(metrics.BouncedEmails) / float64(metrics.TotalEmails) * 100
		metrics.SpamRate = float64(metrics.SpamEmails) / float64(metrics.TotalEmails) * 100
	}

	// Set time range description
	if filters.TimeRange != nil {
		metrics.TimeRange = fmt.Sprintf("%s to %s",
			filters.TimeRange.Start.Format("2006-01-02"),
			filters.TimeRange.End.Format("2006-01-02"))
	} else {
		metrics.TimeRange = "All time"
	}

	return metrics, nil
}

func (m *MockEmailEventRepositoryForMetrics) List(ctx context.Context, filters domain.EmailEventFilters) ([]*domain.EmailEvent, error) {
	args := m.Called(ctx, filters)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	var events []*domain.EmailEvent

	for _, event := range m.events {
		// Apply filters
		if filters.TimeRange != nil {
			if event.SentAt.Before(filters.TimeRange.Start) || event.SentAt.After(filters.TimeRange.End) {
				continue
			}
		}
		if filters.EmailType != nil && event.EmailType != *filters.EmailType {
			continue
		}
		if filters.Status != nil && event.Status != *filters.Status {
			continue
		}
		if filters.InquiryID != nil && event.InquiryID != *filters.InquiryID {
			continue
		}

		events = append(events, event)
	}

	// Apply pagination
	start := filters.Offset
	end := start + filters.Limit

	if start >= len(events) {
		return []*domain.EmailEvent{}, nil
	}

	if end > len(events) {
		end = len(events)
	}

	if filters.Limit > 0 {
		return events[start:end], nil
	}

	return events, nil
}

// Add test events to the mock repository
func (m *MockEmailEventRepositoryForMetrics) AddTestEvents(events []*domain.EmailEvent) {
	m.events = append(m.events, events...)
}

// TestEmailMetricsService provides comprehensive unit tests for email metrics service
func TestEmailMetricsService(t *testing.T) {
	t.Run("GetEmailMetrics", func(t *testing.T) {
		testGetEmailMetrics(t)
	})

	t.Run("GetEmailStatusByInquiry", func(t *testing.T) {
		testGetEmailStatusByInquiry(t)
	})

	t.Run("GetEmailEventHistory", func(t *testing.T) {
		testGetEmailEventHistory(t)
	})

	t.Run("FilterValidation", func(t *testing.T) {
		testFilterValidation(t)
	})

	t.Run("ExtendedMethods", func(t *testing.T) {
		testExtendedMethods(t)
	})

	t.Run("ErrorHandling", func(t *testing.T) {
		testErrorHandling(t)
	})

	t.Run("EdgeCases", func(t *testing.T) {
		testEdgeCases(t)
	})

	t.Run("HealthCheck", func(t *testing.T) {
		testHealthCheck(t)
	})
}

func setupTestMetricsService(t *testing.T) (*MockEmailEventRepositoryForMetrics, interfaces.EmailMetricsService) {
	mockRepo := NewMockEmailEventRepositoryForMetrics()
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel) // Reduce noise in tests

	service := services.NewEmailMetricsService(mockRepo, logger)
	require.NotNil(t, service)

	return mockRepo, service
}

func testGetEmailMetrics(t *testing.T) {
	mockRepo, service := setupTestMetricsService(t)
	ctx := context.Background()

	// Create test events
	now := time.Now()
	testEvents := createTestEmailEvents(now)
	mockRepo.AddTestEvents(testEvents)

	t.Run("AllTimeMetrics", func(t *testing.T) {
		timeRange := domain.TimeRange{
			Start: now.Add(-5 * time.Hour),
			End:   now,
		}

		mockRepo.On("GetMetrics", ctx, mock.MatchedBy(func(filters domain.EmailEventFilters) bool {
			return filters.TimeRange != nil &&
				filters.TimeRange.Start.Equal(timeRange.Start) &&
				filters.TimeRange.End.Equal(timeRange.End)
		})).Return(&domain.EmailMetrics{
			TotalEmails:     5,
			DeliveredEmails: 2,
			FailedEmails:    1,
			BouncedEmails:   1,
			SpamEmails:      1,
			DeliveryRate:    40.0,
			BounceRate:      20.0,
			SpamRate:        20.0,
			TimeRange:       "All time",
		}, nil).Once()

		metrics, err := service.GetEmailMetrics(ctx, timeRange)
		assert.NoError(t, err)
		assert.NotNil(t, metrics)
		assert.Equal(t, int64(5), metrics.TotalEmails)
		assert.Equal(t, int64(2), metrics.DeliveredEmails)
		assert.Equal(t, int64(1), metrics.FailedEmails)
		assert.Equal(t, int64(1), metrics.BouncedEmails)
		assert.Equal(t, int64(1), metrics.SpamEmails)
		assert.Equal(t, 40.0, metrics.DeliveryRate)
		assert.Equal(t, 20.0, metrics.BounceRate)
		assert.Equal(t, 20.0, metrics.SpamRate)

		mockRepo.AssertExpectations(t)
	})

	t.Run("EmptyTimeRange", func(t *testing.T) {
		timeRange := domain.TimeRange{
			Start: now.Add(1 * time.Hour), // Future time range
			End:   now.Add(2 * time.Hour),
		}

		mockRepo.On("GetMetrics", ctx, mock.Anything).Return(&domain.EmailMetrics{
			TotalEmails:     0,
			DeliveredEmails: 0,
			FailedEmails:    0,
			BouncedEmails:   0,
			SpamEmails:      0,
			DeliveryRate:    0.0,
			BounceRate:      0.0,
			SpamRate:        0.0,
			TimeRange:       "Empty range",
		}, nil).Once()

		metrics, err := service.GetEmailMetrics(ctx, timeRange)
		assert.NoError(t, err)
		assert.NotNil(t, metrics)
		assert.Equal(t, int64(0), metrics.TotalEmails)
		assert.Equal(t, 0.0, metrics.DeliveryRate)

		mockRepo.AssertExpectations(t)
	})

	t.Run("RepositoryError", func(t *testing.T) {
		timeRange := domain.TimeRange{
			Start: now.Add(-1 * time.Hour),
			End:   now,
		}

		mockRepo.On("GetMetrics", ctx, mock.Anything).Return(nil, fmt.Errorf("database error")).Once()

		metrics, err := service.GetEmailMetrics(ctx, timeRange)
		assert.Error(t, err)
		assert.Nil(t, metrics)
		assert.Contains(t, err.Error(), "database error")

		mockRepo.AssertExpectations(t)
	})
}

func testGetEmailStatusByInquiry(t *testing.T) {
	mockRepo, service := setupTestMetricsService(t)
	ctx := context.Background()

	t.Run("InquiryWithMultipleEmails", func(t *testing.T) {
		inquiryID := "inquiry-multi-" + uuid.New().String()
		now := time.Now()

		events := []*domain.EmailEvent{
			{
				ID:             uuid.New().String(),
				InquiryID:      inquiryID,
				EmailType:      domain.EmailTypeCustomerConfirmation,
				RecipientEmail: "customer@example.com",
				SenderEmail:    "info@cloudpartner.pro",
				Status:         domain.EmailStatusDelivered,
				SentAt:         now.Add(-2 * time.Hour),
				DeliveredAt:    timePtr(now.Add(-2*time.Hour + 30*time.Minute)),
			},
			{
				ID:             uuid.New().String(),
				InquiryID:      inquiryID,
				EmailType:      domain.EmailTypeConsultantNotification,
				RecipientEmail: "consultant@cloudpartner.pro",
				SenderEmail:    "info@cloudpartner.pro",
				Status:         domain.EmailStatusDelivered,
				SentAt:         now.Add(-1 * time.Hour),
				DeliveredAt:    timePtr(now.Add(-1*time.Hour + 15*time.Minute)),
			},
		}

		mockRepo.On("GetByInquiryID", ctx, inquiryID).Return(events, nil).Once()

		status, err := service.GetEmailStatusByInquiry(ctx, inquiryID)
		assert.NoError(t, err)
		assert.NotNil(t, status)
		assert.Equal(t, inquiryID, status.InquiryID)
		assert.Equal(t, 2, status.TotalEmailsSent)
		assert.NotNil(t, status.CustomerEmail)
		assert.NotNil(t, status.ConsultantEmail)
		assert.Nil(t, status.InquiryNotification)
		assert.NotNil(t, status.LastEmailSent)

		// Verify most recent email is the consultant notification
		assert.True(t, status.LastEmailSent.Equal(events[1].SentAt))

		mockRepo.AssertExpectations(t)
	})

	t.Run("InquiryWithNoEmails", func(t *testing.T) {
		inquiryID := "inquiry-empty-" + uuid.New().String()

		mockRepo.On("GetByInquiryID", ctx, inquiryID).Return([]*domain.EmailEvent{}, nil).Once()

		status, err := service.GetEmailStatusByInquiry(ctx, inquiryID)
		assert.NoError(t, err)
		assert.Nil(t, status)

		mockRepo.AssertExpectations(t)
	})

	t.Run("InquiryWithFailedEmails", func(t *testing.T) {
		inquiryID := "inquiry-failed-" + uuid.New().String()
		now := time.Now()

		events := []*domain.EmailEvent{
			{
				ID:             uuid.New().String(),
				InquiryID:      inquiryID,
				EmailType:      domain.EmailTypeCustomerConfirmation,
				RecipientEmail: "customer@example.com",
				SenderEmail:    "info@cloudpartner.pro",
				Status:         domain.EmailStatusFailed,
				SentAt:         now.Add(-1 * time.Hour),
				ErrorMessage:   "SMTP connection failed",
			},
		}

		mockRepo.On("GetByInquiryID", ctx, inquiryID).Return(events, nil).Once()

		status, err := service.GetEmailStatusByInquiry(ctx, inquiryID)
		assert.NoError(t, err)
		assert.NotNil(t, status)
		assert.Equal(t, 1, status.TotalEmailsSent)
		assert.NotNil(t, status.CustomerEmail)
		assert.Equal(t, domain.EmailStatusFailed, status.CustomerEmail.Status)
		assert.Equal(t, "SMTP connection failed", status.CustomerEmail.ErrorMessage)

		mockRepo.AssertExpectations(t)
	})

	t.Run("RepositoryError", func(t *testing.T) {
		inquiryID := "inquiry-error-" + uuid.New().String()

		mockRepo.On("GetByInquiryID", ctx, inquiryID).Return(nil, fmt.Errorf("database error")).Once()

		status, err := service.GetEmailStatusByInquiry(ctx, inquiryID)
		assert.Error(t, err)
		assert.Nil(t, status)
		assert.Contains(t, err.Error(), "database error")

		mockRepo.AssertExpectations(t)
	})
}

func testGetEmailEventHistory(t *testing.T) {
	mockRepo, service := setupTestMetricsService(t)
	ctx := context.Background()

	t.Run("ValidFilters", func(t *testing.T) {
		filters := domain.EmailEventFilters{
			Limit:  10,
			Offset: 0,
		}

		expectedEvents := createTestEmailEvents(time.Now())
		mockRepo.On("List", ctx, filters).Return(expectedEvents, nil).Once()

		events, err := service.GetEmailEventHistory(ctx, filters)
		assert.NoError(t, err)
		assert.Len(t, events, len(expectedEvents))

		mockRepo.AssertExpectations(t)
	})

	t.Run("FilterByEmailType", func(t *testing.T) {
		customerEmailType := domain.EmailTypeCustomerConfirmation
		filters := domain.EmailEventFilters{
			EmailType: &customerEmailType,
			Limit:     10,
			Offset:    0,
		}

		expectedEvents := []*domain.EmailEvent{
			{
				ID:             uuid.New().String(),
				InquiryID:      "inquiry-filter",
				EmailType:      domain.EmailTypeCustomerConfirmation,
				RecipientEmail: "customer@example.com",
				SenderEmail:    "info@cloudpartner.pro",
				Status:         domain.EmailStatusDelivered,
				SentAt:         time.Now(),
			},
		}

		mockRepo.On("List", ctx, filters).Return(expectedEvents, nil).Once()

		events, err := service.GetEmailEventHistory(ctx, filters)
		assert.NoError(t, err)
		assert.Len(t, events, 1)
		assert.Equal(t, domain.EmailTypeCustomerConfirmation, events[0].EmailType)

		mockRepo.AssertExpectations(t)
	})

	t.Run("FilterByStatus", func(t *testing.T) {
		deliveredStatus := domain.EmailStatusDelivered
		filters := domain.EmailEventFilters{
			Status: &deliveredStatus,
			Limit:  10,
			Offset: 0,
		}

		expectedEvents := []*domain.EmailEvent{
			{
				ID:             uuid.New().String(),
				InquiryID:      "inquiry-delivered",
				EmailType:      domain.EmailTypeCustomerConfirmation,
				RecipientEmail: "delivered@example.com",
				SenderEmail:    "info@cloudpartner.pro",
				Status:         domain.EmailStatusDelivered,
				SentAt:         time.Now(),
			},
		}

		mockRepo.On("List", ctx, filters).Return(expectedEvents, nil).Once()

		events, err := service.GetEmailEventHistory(ctx, filters)
		assert.NoError(t, err)
		assert.Len(t, events, 1)
		assert.Equal(t, domain.EmailStatusDelivered, events[0].Status)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Pagination", func(t *testing.T) {
		filters := domain.EmailEventFilters{
			Limit:  5,
			Offset: 10,
		}

		expectedEvents := createTestEmailEvents(time.Now())[:3] // Return fewer events for second page
		mockRepo.On("List", ctx, filters).Return(expectedEvents, nil).Once()

		events, err := service.GetEmailEventHistory(ctx, filters)
		assert.NoError(t, err)
		assert.Len(t, events, 3)

		mockRepo.AssertExpectations(t)
	})

	t.Run("RepositoryError", func(t *testing.T) {
		filters := domain.EmailEventFilters{
			Limit:  10,
			Offset: 0,
		}

		mockRepo.On("List", ctx, filters).Return(nil, fmt.Errorf("database error")).Once()

		events, err := service.GetEmailEventHistory(ctx, filters)
		assert.Error(t, err)
		assert.Nil(t, events)
		assert.Contains(t, err.Error(), "database error")

		mockRepo.AssertExpectations(t)
	})
}

func testFilterValidation(t *testing.T) {
	_, service := setupTestMetricsService(t)
	ctx := context.Background()

	t.Run("InvalidTimeRange", func(t *testing.T) {
		filters := domain.EmailEventFilters{
			TimeRange: &domain.TimeRange{
				Start: time.Now(),
				End:   time.Now().Add(-1 * time.Hour), // End before start
			},
			Limit:  10,
			Offset: 0,
		}

		events, err := service.GetEmailEventHistory(ctx, filters)
		assert.Error(t, err)
		assert.Nil(t, events)
		assert.Contains(t, err.Error(), "start time cannot be after end time")
	})

	t.Run("FutureTimeRange", func(t *testing.T) {
		filters := domain.EmailEventFilters{
			TimeRange: &domain.TimeRange{
				Start: time.Now().Add(25 * time.Hour), // More than 24 hours in future
				End:   time.Now().Add(26 * time.Hour),
			},
			Limit:  10,
			Offset: 0,
		}

		events, err := service.GetEmailEventHistory(ctx, filters)
		assert.Error(t, err)
		assert.Nil(t, events)
		assert.Contains(t, err.Error(), "cannot be more than 24 hours in the future")
	})

	t.Run("NegativeLimit", func(t *testing.T) {
		filters := domain.EmailEventFilters{
			Limit:  -1,
			Offset: 0,
		}

		events, err := service.GetEmailEventHistory(ctx, filters)
		assert.Error(t, err)
		assert.Nil(t, events)
		assert.Contains(t, err.Error(), "limit cannot be negative")
	})

	t.Run("ExcessiveLimit", func(t *testing.T) {
		filters := domain.EmailEventFilters{
			Limit:  1001, // Exceeds maximum
			Offset: 0,
		}

		events, err := service.GetEmailEventHistory(ctx, filters)
		assert.Error(t, err)
		assert.Nil(t, events)
		assert.Contains(t, err.Error(), "limit cannot exceed 1000")
	})

	t.Run("NegativeOffset", func(t *testing.T) {
		filters := domain.EmailEventFilters{
			Limit:  10,
			Offset: -1,
		}

		events, err := service.GetEmailEventHistory(ctx, filters)
		assert.Error(t, err)
		assert.Nil(t, events)
		assert.Contains(t, err.Error(), "offset cannot be negative")
	})

	t.Run("InvalidEmailType", func(t *testing.T) {
		invalidEmailType := domain.EmailEventType("invalid_type")
		filters := domain.EmailEventFilters{
			EmailType: &invalidEmailType,
			Limit:     10,
			Offset:    0,
		}

		events, err := service.GetEmailEventHistory(ctx, filters)
		assert.Error(t, err)
		assert.Nil(t, events)
		assert.Contains(t, err.Error(), "invalid email type")
	})

	t.Run("InvalidStatus", func(t *testing.T) {
		invalidStatus := domain.EmailEventStatus("invalid_status")
		filters := domain.EmailEventFilters{
			Status: &invalidStatus,
			Limit:  10,
			Offset: 0,
		}

		events, err := service.GetEmailEventHistory(ctx, filters)
		assert.Error(t, err)
		assert.Nil(t, events)
		assert.Contains(t, err.Error(), "invalid email status")
	})
}

func testExtendedMethods(t *testing.T) {
	mockRepo, service := setupTestMetricsService(t)
	ctx := context.Background()

	// Cast to concrete type to access extended methods
	concreteService, ok := service.(*services.EmailMetricsServiceImpl)
	require.True(t, ok, "Service should be concrete implementation")

	t.Run("GetEmailMetricsByType", func(t *testing.T) {
		now := time.Now()
		timeRange := domain.TimeRange{
			Start: now.Add(-1 * time.Hour),
			End:   now,
		}

		// Mock metrics for each email type
		customerType := domain.EmailTypeCustomerConfirmation
		consultantType := domain.EmailTypeConsultantNotification
		inquiryType := domain.EmailTypeInquiryNotification

		mockRepo.On("GetMetrics", ctx, mock.MatchedBy(func(filters domain.EmailEventFilters) bool {
			return filters.EmailType != nil && *filters.EmailType == customerType
		})).Return(&domain.EmailMetrics{
			TotalEmails:     2,
			DeliveredEmails: 1,
			FailedEmails:    1,
			DeliveryRate:    50.0,
		}, nil).Once()

		mockRepo.On("GetMetrics", ctx, mock.MatchedBy(func(filters domain.EmailEventFilters) bool {
			return filters.EmailType != nil && *filters.EmailType == consultantType
		})).Return(&domain.EmailMetrics{
			TotalEmails:     3,
			DeliveredEmails: 3,
			FailedEmails:    0,
			DeliveryRate:    100.0,
		}, nil).Once()

		mockRepo.On("GetMetrics", ctx, mock.MatchedBy(func(filters domain.EmailEventFilters) bool {
			return filters.EmailType != nil && *filters.EmailType == inquiryType
		})).Return(&domain.EmailMetrics{
			TotalEmails:     1,
			DeliveredEmails: 0,
			FailedEmails:    1,
			DeliveryRate:    0.0,
		}, nil).Once()

		metricsByType, err := concreteService.GetEmailMetricsByType(ctx, timeRange)
		assert.NoError(t, err)
		assert.Len(t, metricsByType, 3)
		assert.Contains(t, metricsByType, customerType)
		assert.Contains(t, metricsByType, consultantType)
		assert.Contains(t, metricsByType, inquiryType)

		assert.Equal(t, int64(2), metricsByType[customerType].TotalEmails)
		assert.Equal(t, int64(3), metricsByType[consultantType].TotalEmails)
		assert.Equal(t, int64(1), metricsByType[inquiryType].TotalEmails)

		mockRepo.AssertExpectations(t)
	})

	t.Run("GetRecentEmailActivity", func(t *testing.T) {
		hours := 24
		expectedEvents := createTestEmailEvents(time.Now())

		mockRepo.On("List", ctx, mock.MatchedBy(func(filters domain.EmailEventFilters) bool {
			return filters.TimeRange != nil && filters.Limit == 100
		})).Return(expectedEvents, nil).Once()

		events, err := concreteService.GetRecentEmailActivity(ctx, hours)
		assert.NoError(t, err)
		assert.Len(t, events, len(expectedEvents))

		mockRepo.AssertExpectations(t)
	})

	t.Run("GetRecentEmailActivityInvalidHours", func(t *testing.T) {
		// Test with invalid hours (0)
		events, err := concreteService.GetRecentEmailActivity(ctx, 0)
		assert.Error(t, err)
		assert.Nil(t, events)
		assert.Contains(t, err.Error(), "hours must be between 1 and 168")

		// Test with excessive hours (more than 1 week)
		events, err = concreteService.GetRecentEmailActivity(ctx, 200)
		assert.Error(t, err)
		assert.Nil(t, events)
		assert.Contains(t, err.Error(), "hours must be between 1 and 168")
	})
}

func testErrorHandling(t *testing.T) {
	mockRepo, service := setupTestMetricsService(t)
	ctx := context.Background()

	t.Run("GetEmailMetricsWithRepositoryError", func(t *testing.T) {
		timeRange := domain.TimeRange{
			Start: time.Now().Add(-1 * time.Hour),
			End:   time.Now(),
		}

		mockRepo.On("GetMetrics", ctx, mock.Anything).Return(nil, fmt.Errorf("connection failed")).Once()

		metrics, err := service.GetEmailMetrics(ctx, timeRange)
		assert.Error(t, err)
		assert.Nil(t, metrics)
		assert.Contains(t, err.Error(), "connection failed")

		mockRepo.AssertExpectations(t)
	})

	t.Run("GetEmailStatusByInquiryWithRepositoryError", func(t *testing.T) {
		inquiryID := "inquiry-error"

		mockRepo.On("GetByInquiryID", ctx, inquiryID).Return(nil, fmt.Errorf("database timeout")).Once()

		status, err := service.GetEmailStatusByInquiry(ctx, inquiryID)
		assert.Error(t, err)
		assert.Nil(t, status)
		assert.Contains(t, err.Error(), "database timeout")

		mockRepo.AssertExpectations(t)
	})
}

func testEdgeCases(t *testing.T) {
	mockRepo, service := setupTestMetricsService(t)
	ctx := context.Background()

	t.Run("EmptyInquiryID", func(t *testing.T) {
		// This should be handled gracefully by the repository
		mockRepo.On("GetByInquiryID", ctx, "").Return([]*domain.EmailEvent{}, nil).Once()

		status, err := service.GetEmailStatusByInquiry(ctx, "")
		assert.NoError(t, err)
		assert.Nil(t, status)

		mockRepo.AssertExpectations(t)
	})

	t.Run("ZeroLimitFilter", func(t *testing.T) {
		filters := domain.EmailEventFilters{
			Limit:  0,
			Offset: 0,
		}

		expectedEvents := createTestEmailEvents(time.Now())
		mockRepo.On("List", ctx, filters).Return(expectedEvents, nil).Once()

		events, err := service.GetEmailEventHistory(ctx, filters)
		assert.NoError(t, err)
		assert.NotNil(t, events)

		mockRepo.AssertExpectations(t)
	})

	t.Run("InquiryWithDuplicateEmailTypes", func(t *testing.T) {
		inquiryID := "inquiry-duplicate-" + uuid.New().String()
		now := time.Now()

		// Multiple customer confirmation emails (should keep most recent)
		events := []*domain.EmailEvent{
			{
				ID:             uuid.New().String(),
				InquiryID:      inquiryID,
				EmailType:      domain.EmailTypeCustomerConfirmation,
				RecipientEmail: "customer@example.com",
				SenderEmail:    "info@cloudpartner.pro",
				Status:         domain.EmailStatusSent,
				SentAt:         now.Add(-2 * time.Hour),
			},
			{
				ID:             uuid.New().String(),
				InquiryID:      inquiryID,
				EmailType:      domain.EmailTypeCustomerConfirmation,
				RecipientEmail: "customer@example.com",
				SenderEmail:    "info@cloudpartner.pro",
				Status:         domain.EmailStatusDelivered,
				SentAt:         now.Add(-1 * time.Hour), // More recent
			},
		}

		mockRepo.On("GetByInquiryID", ctx, inquiryID).Return(events, nil).Once()

		status, err := service.GetEmailStatusByInquiry(ctx, inquiryID)
		assert.NoError(t, err)
		assert.NotNil(t, status)
		assert.Equal(t, 2, status.TotalEmailsSent)
		assert.NotNil(t, status.CustomerEmail)

		// Should have the most recent customer email (delivered status)
		assert.Equal(t, domain.EmailStatusDelivered, status.CustomerEmail.Status)
		assert.True(t, status.CustomerEmail.SentAt.Equal(events[1].SentAt))

		mockRepo.AssertExpectations(t)
	})
}

func testHealthCheck(t *testing.T) {
	mockRepo, service := setupTestMetricsService(t)
	ctx := context.Background()

	// Cast to concrete type to access health check method
	concreteService, ok := service.(*services.EmailMetricsServiceImpl)
	require.True(t, ok, "Service should be concrete implementation")

	t.Run("HealthyService", func(t *testing.T) {
		mockRepo.On("GetMetrics", ctx, mock.Anything).Return(&domain.EmailMetrics{
			TotalEmails: 0,
		}, nil).Once()

		isHealthy := concreteService.IsHealthy(ctx)
		assert.True(t, isHealthy)

		mockRepo.AssertExpectations(t)
	})

	t.Run("UnhealthyService", func(t *testing.T) {
		mockRepo.On("GetMetrics", ctx, mock.Anything).Return(nil, fmt.Errorf("database connection failed")).Once()

		isHealthy := concreteService.IsHealthy(ctx)
		assert.False(t, isHealthy)

		mockRepo.AssertExpectations(t)
	})
}

// Helper functions

func createTestEmailEvents(now time.Time) []*domain.EmailEvent {
	return []*domain.EmailEvent{
		{
			ID:             uuid.New().String(),
			InquiryID:      "inquiry-1",
			EmailType:      domain.EmailTypeCustomerConfirmation,
			RecipientEmail: "customer1@example.com",
			SenderEmail:    "info@cloudpartner.pro",
			Status:         domain.EmailStatusDelivered,
			SentAt:         now.Add(-1 * time.Hour),
			DeliveredAt:    timePtr(now.Add(-50 * time.Minute)),
		},
		{
			ID:             uuid.New().String(),
			InquiryID:      "inquiry-1",
			EmailType:      domain.EmailTypeConsultantNotification,
			RecipientEmail: "consultant@cloudpartner.pro",
			SenderEmail:    "info@cloudpartner.pro",
			Status:         domain.EmailStatusDelivered,
			SentAt:         now.Add(-2 * time.Hour),
			DeliveredAt:    timePtr(now.Add(-110 * time.Minute)),
		},
		{
			ID:             uuid.New().String(),
			InquiryID:      "inquiry-2",
			EmailType:      domain.EmailTypeCustomerConfirmation,
			RecipientEmail: "customer2@example.com",
			SenderEmail:    "info@cloudpartner.pro",
			Status:         domain.EmailStatusFailed,
			SentAt:         now.Add(-3 * time.Hour),
			ErrorMessage:   "SMTP connection failed",
		},
		{
			ID:             uuid.New().String(),
			InquiryID:      "inquiry-3",
			EmailType:      domain.EmailTypeInquiryNotification,
			RecipientEmail: "admin@cloudpartner.pro",
			SenderEmail:    "info@cloudpartner.pro",
			Status:         domain.EmailStatusBounced,
			SentAt:         now.Add(-4 * time.Hour),
			BounceType:     "permanent",
			ErrorMessage:   "Email address does not exist",
		},
		{
			ID:             uuid.New().String(),
			InquiryID:      "inquiry-4",
			EmailType:      domain.EmailTypeCustomerConfirmation,
			RecipientEmail: "spam@example.com",
			SenderEmail:    "info@cloudpartner.pro",
			Status:         domain.EmailStatusSpam,
			SentAt:         now.Add(-5 * time.Hour),
			ErrorMessage:   "Marked as spam",
		},
	}
}

func timePtr(t time.Time) *time.Time {
	return &t
}

// Main function for running tests standalone
func main() {
	fmt.Println("=== Comprehensive Email Metrics Service Tests ===")

	// Note: This would normally be run with `go test` command
	// This main function is for demonstration purposes

	fmt.Println("Run with: go test -v ./test_email_metrics_service_comprehensive.go")
	fmt.Println("Or integrate into your test suite")
}
