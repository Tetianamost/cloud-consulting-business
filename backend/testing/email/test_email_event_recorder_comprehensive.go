package main

import (
	"context"
	"fmt"
	"sync"
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

// MockEmailEventRepository provides a mock implementation for testing
type MockEmailEventRepository struct {
	mock.Mock
	events map[string]*domain.EmailEvent
	mutex  sync.RWMutex
}

func NewMockEmailEventRepository() *MockEmailEventRepository {
	return &MockEmailEventRepository{
		events: make(map[string]*domain.EmailEvent),
	}
}

func (m *MockEmailEventRepository) Create(ctx context.Context, event *domain.EmailEvent) error {
	args := m.Called(ctx, event)

	if args.Error(0) == nil {
		m.mutex.Lock()
		defer m.mutex.Unlock()

		// Generate ID if not set
		if event.ID == "" {
			event.ID = uuid.New().String()
		}

		// Set timestamps if not set
		now := time.Now()
		if event.CreatedAt.IsZero() {
			event.CreatedAt = now
		}
		if event.UpdatedAt.IsZero() {
			event.UpdatedAt = now
		}
		if event.SentAt.IsZero() {
			event.SentAt = now
		}

		// Store the event
		m.events[event.ID] = event
	}

	return args.Error(0)
}

func (m *MockEmailEventRepository) Update(ctx context.Context, event *domain.EmailEvent) error {
	args := m.Called(ctx, event)

	if args.Error(0) == nil {
		m.mutex.Lock()
		defer m.mutex.Unlock()

		if _, exists := m.events[event.ID]; exists {
			event.UpdatedAt = time.Now()
			m.events[event.ID] = event
		} else {
			return fmt.Errorf("email event not found")
		}
	}

	return args.Error(0)
}

func (m *MockEmailEventRepository) GetByInquiryID(ctx context.Context, inquiryID string) ([]*domain.EmailEvent, error) {
	args := m.Called(ctx, inquiryID)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	m.mutex.RLock()
	defer m.mutex.RUnlock()

	var events []*domain.EmailEvent
	for _, event := range m.events {
		if event.InquiryID == inquiryID {
			events = append(events, event)
		}
	}

	return events, nil
}

func (m *MockEmailEventRepository) GetByMessageID(ctx context.Context, messageID string) (*domain.EmailEvent, error) {
	args := m.Called(ctx, messageID)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	m.mutex.RLock()
	defer m.mutex.RUnlock()

	for _, event := range m.events {
		if event.SESMessageID == messageID {
			return event, nil
		}
	}

	return nil, nil
}

func (m *MockEmailEventRepository) GetMetrics(ctx context.Context, filters domain.EmailEventFilters) (*domain.EmailMetrics, error) {
	args := m.Called(ctx, filters)
	return args.Get(0).(*domain.EmailMetrics), args.Error(1)
}

func (m *MockEmailEventRepository) List(ctx context.Context, filters domain.EmailEventFilters) ([]*domain.EmailEvent, error) {
	args := m.Called(ctx, filters)
	return args.Get(0).([]*domain.EmailEvent), args.Error(1)
}

// TestEmailEventRecorder provides comprehensive unit tests for email event recorder service
func TestEmailEventRecorder(t *testing.T) {
	t.Run("RecordEmailSent", func(t *testing.T) {
		testRecordEmailSent(t)
	})

	t.Run("UpdateEmailStatus", func(t *testing.T) {
		testUpdateEmailStatus(t)
	})

	t.Run("GetEmailEventsByInquiry", func(t *testing.T) {
		testGetEmailEventsByInquiry(t)
	})

	t.Run("SynchronousMethods", func(t *testing.T) {
		testSynchronousMethods(t)
	})

	t.Run("ErrorHandling", func(t *testing.T) {
		testErrorHandling(t)
	})

	t.Run("NonBlockingBehavior", func(t *testing.T) {
		testNonBlockingBehavior(t)
	})

	t.Run("RetryLogic", func(t *testing.T) {
		testRetryLogic(t)
	})

	t.Run("HealthCheck", func(t *testing.T) {
		testHealthCheck(t)
	})

	t.Run("ConcurrentOperations", func(t *testing.T) {
		testConcurrentOperations(t)
	})
}

func setupTestRecorder(t *testing.T) (*MockEmailEventRepository, interfaces.EmailEventRecorder) {
	mockRepo := NewMockEmailEventRepository()
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel) // Reduce noise in tests

	recorder := services.NewEmailEventRecorder(mockRepo, logger)
	require.NotNil(t, recorder)

	return mockRepo, recorder
}

func testRecordEmailSent(t *testing.T) {
	mockRepo, recorder := setupTestRecorder(t)
	ctx := context.Background()

	t.Run("ValidEvent", func(t *testing.T) {
		event := createTestEmailEvent()

		mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(e *domain.EmailEvent) bool {
			return e.InquiryID == event.InquiryID && e.EmailType == event.EmailType
		})).Return(nil).Once()

		err := recorder.RecordEmailSent(ctx, event)
		assert.NoError(t, err)
		assert.NotEmpty(t, event.ID)
		assert.False(t, event.SentAt.IsZero())
		assert.False(t, event.CreatedAt.IsZero())
		assert.False(t, event.UpdatedAt.IsZero())
		assert.Equal(t, domain.EmailStatusSent, event.Status)

		// Wait for async operation to complete
		time.Sleep(100 * time.Millisecond)
		mockRepo.AssertExpectations(t)
	})

	t.Run("EventWithExistingID", func(t *testing.T) {
		event := createTestEmailEvent()
		existingID := uuid.New().String()
		event.ID = existingID

		mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(e *domain.EmailEvent) bool {
			return e.ID == existingID
		})).Return(nil).Once()

		err := recorder.RecordEmailSent(ctx, event)
		assert.NoError(t, err)
		assert.Equal(t, existingID, event.ID) // Should preserve existing ID

		time.Sleep(100 * time.Millisecond)
		mockRepo.AssertExpectations(t)
	})

	t.Run("EventWithExistingTimestamps", func(t *testing.T) {
		event := createTestEmailEvent()
		existingTime := time.Now().Add(-1 * time.Hour)
		event.SentAt = existingTime
		event.CreatedAt = existingTime
		event.UpdatedAt = existingTime

		mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(e *domain.EmailEvent) bool {
			return e.SentAt.Equal(existingTime)
		})).Return(nil).Once()

		err := recorder.RecordEmailSent(ctx, event)
		assert.NoError(t, err)
		assert.Equal(t, existingTime, event.SentAt) // Should preserve existing timestamps

		time.Sleep(100 * time.Millisecond)
		mockRepo.AssertExpectations(t)
	})

	t.Run("EventWithExistingStatus", func(t *testing.T) {
		event := createTestEmailEvent()
		event.Status = domain.EmailStatusDelivered

		mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(e *domain.EmailEvent) bool {
			return e.Status == domain.EmailStatusDelivered
		})).Return(nil).Once()

		err := recorder.RecordEmailSent(ctx, event)
		assert.NoError(t, err)
		assert.Equal(t, domain.EmailStatusDelivered, event.Status) // Should preserve existing status

		time.Sleep(100 * time.Millisecond)
		mockRepo.AssertExpectations(t)
	})

	t.Run("InvalidEvent", func(t *testing.T) {
		event := &domain.EmailEvent{
			// Missing required fields
			EmailType: domain.EmailTypeCustomerConfirmation,
		}

		err := recorder.RecordEmailSent(ctx, event)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "validation failed")
	})
}

func testUpdateEmailStatus(t *testing.T) {
	mockRepo, recorder := setupTestRecorder(t)
	ctx := context.Background()

	t.Run("ValidUpdate", func(t *testing.T) {
		messageID := "ses-update-" + uuid.New().String()
		existingEvent := createTestEmailEvent()
		existingEvent.SESMessageID = messageID

		mockRepo.On("GetByMessageID", ctx, messageID).Return(existingEvent, nil).Once()
		mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(e *domain.EmailEvent) bool {
			return e.Status == domain.EmailStatusDelivered && e.SESMessageID == messageID
		})).Return(nil).Once()

		deliveredAt := time.Now()
		err := recorder.UpdateEmailStatus(ctx, messageID, domain.EmailStatusDelivered, &deliveredAt, "")
		assert.NoError(t, err)

		time.Sleep(100 * time.Millisecond)
		mockRepo.AssertExpectations(t)
	})

	t.Run("UpdateWithErrorMessage", func(t *testing.T) {
		messageID := "ses-error-" + uuid.New().String()
		existingEvent := createTestEmailEvent()
		existingEvent.SESMessageID = messageID

		mockRepo.On("GetByMessageID", ctx, messageID).Return(existingEvent, nil).Once()
		mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(e *domain.EmailEvent) bool {
			return e.Status == domain.EmailStatusFailed && e.ErrorMessage == "SMTP connection failed"
		})).Return(nil).Once()

		err := recorder.UpdateEmailStatus(ctx, messageID, domain.EmailStatusFailed, nil, "SMTP connection failed")
		assert.NoError(t, err)

		time.Sleep(100 * time.Millisecond)
		mockRepo.AssertExpectations(t)
	})

	t.Run("NonExistentMessageID", func(t *testing.T) {
		messageID := "non-existent-" + uuid.New().String()

		mockRepo.On("GetByMessageID", ctx, messageID).Return(nil, nil).Once()

		err := recorder.UpdateEmailStatus(ctx, messageID, domain.EmailStatusDelivered, nil, "")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")

		mockRepo.AssertExpectations(t)
	})

	t.Run("EmptyMessageID", func(t *testing.T) {
		err := recorder.UpdateEmailStatus(ctx, "", domain.EmailStatusDelivered, nil, "")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "message ID is required")
	})

	t.Run("RepositoryError", func(t *testing.T) {
		messageID := "ses-repo-error-" + uuid.New().String()

		mockRepo.On("GetByMessageID", ctx, messageID).Return(nil, fmt.Errorf("database error")).Once()

		err := recorder.UpdateEmailStatus(ctx, messageID, domain.EmailStatusDelivered, nil, "")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database error")

		mockRepo.AssertExpectations(t)
	})
}

func testGetEmailEventsByInquiry(t *testing.T) {
	mockRepo, recorder := setupTestRecorder(t)
	ctx := context.Background()

	t.Run("ValidInquiryID", func(t *testing.T) {
		inquiryID := "inquiry-" + uuid.New().String()
		expectedEvents := []*domain.EmailEvent{
			createTestEmailEvent(),
			createTestEmailEvent(),
		}
		expectedEvents[0].InquiryID = inquiryID
		expectedEvents[1].InquiryID = inquiryID

		mockRepo.On("GetByInquiryID", ctx, inquiryID).Return(expectedEvents, nil).Once()

		events, err := recorder.GetEmailEventsByInquiry(ctx, inquiryID)
		assert.NoError(t, err)
		assert.Len(t, events, 2)
		assert.Equal(t, inquiryID, events[0].InquiryID)
		assert.Equal(t, inquiryID, events[1].InquiryID)

		mockRepo.AssertExpectations(t)
	})

	t.Run("EmptyInquiryID", func(t *testing.T) {
		events, err := recorder.GetEmailEventsByInquiry(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, events)
		assert.Contains(t, err.Error(), "inquiry ID is required")
	})

	t.Run("NoEventsFound", func(t *testing.T) {
		inquiryID := "inquiry-empty-" + uuid.New().String()

		mockRepo.On("GetByInquiryID", ctx, inquiryID).Return([]*domain.EmailEvent{}, nil).Once()

		events, err := recorder.GetEmailEventsByInquiry(ctx, inquiryID)
		assert.NoError(t, err)
		assert.Empty(t, events)

		mockRepo.AssertExpectations(t)
	})

	t.Run("RepositoryError", func(t *testing.T) {
		inquiryID := "inquiry-error-" + uuid.New().String()

		mockRepo.On("GetByInquiryID", ctx, inquiryID).Return(nil, fmt.Errorf("database error")).Once()

		events, err := recorder.GetEmailEventsByInquiry(ctx, inquiryID)
		assert.Error(t, err)
		assert.Nil(t, events)
		assert.Contains(t, err.Error(), "database error")

		mockRepo.AssertExpectations(t)
	})
}

func testSynchronousMethods(t *testing.T) {
	mockRepo, recorder := setupTestRecorder(t)
	ctx := context.Background()

	// Cast to concrete type to access synchronous methods
	concreteRecorder, ok := recorder.(*services.EmailEventRecorderImpl)
	require.True(t, ok, "Recorder should be concrete implementation")

	t.Run("RecordEmailSentSync", func(t *testing.T) {
		event := createTestEmailEvent()

		mockRepo.On("Create", ctx, mock.MatchedBy(func(e *domain.EmailEvent) bool {
			return e.InquiryID == event.InquiryID
		})).Return(nil).Once()

		err := concreteRecorder.RecordEmailSentSync(ctx, event)
		assert.NoError(t, err)
		assert.NotEmpty(t, event.ID)

		mockRepo.AssertExpectations(t)
	})

	t.Run("RecordEmailSentSyncWithError", func(t *testing.T) {
		event := createTestEmailEvent()

		mockRepo.On("Create", ctx, mock.Anything).Return(fmt.Errorf("database error")).Once()

		err := concreteRecorder.RecordEmailSentSync(ctx, event)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database error")

		mockRepo.AssertExpectations(t)
	})

	t.Run("UpdateEmailStatusSync", func(t *testing.T) {
		messageID := "ses-sync-" + uuid.New().String()
		existingEvent := createTestEmailEvent()
		existingEvent.SESMessageID = messageID

		mockRepo.On("GetByMessageID", ctx, messageID).Return(existingEvent, nil).Once()
		mockRepo.On("Update", ctx, mock.MatchedBy(func(e *domain.EmailEvent) bool {
			return e.Status == domain.EmailStatusDelivered
		})).Return(nil).Once()

		deliveredAt := time.Now()
		err := concreteRecorder.UpdateEmailStatusSync(ctx, messageID, domain.EmailStatusDelivered, &deliveredAt, "")
		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("UpdateEmailStatusSyncWithError", func(t *testing.T) {
		messageID := "ses-sync-error-" + uuid.New().String()
		existingEvent := createTestEmailEvent()
		existingEvent.SESMessageID = messageID

		mockRepo.On("GetByMessageID", ctx, messageID).Return(existingEvent, nil).Once()
		mockRepo.On("Update", ctx, mock.Anything).Return(fmt.Errorf("database error")).Once()

		err := concreteRecorder.UpdateEmailStatusSync(ctx, messageID, domain.EmailStatusDelivered, nil, "")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database error")

		mockRepo.AssertExpectations(t)
	})
}

func testErrorHandling(t *testing.T) {
	mockRepo, recorder := setupTestRecorder(t)
	ctx := context.Background()

	t.Run("RecordEmailSentWithNilEvent", func(t *testing.T) {
		err := recorder.RecordEmailSent(ctx, nil)
		assert.Error(t, err)
	})

	t.Run("RecordEmailSentWithInvalidEvent", func(t *testing.T) {
		event := &domain.EmailEvent{
			// Missing required fields
			ID: uuid.New().String(),
		}

		err := recorder.RecordEmailSent(ctx, event)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "validation failed")
	})

	t.Run("UpdateEmailStatusWithRepositoryError", func(t *testing.T) {
		messageID := "ses-repo-error-" + uuid.New().String()

		mockRepo.On("GetByMessageID", mock.Anything, messageID).Return(nil, fmt.Errorf("connection failed")).Once()

		err := recorder.UpdateEmailStatus(ctx, messageID, domain.EmailStatusDelivered, nil, "")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "connection failed")

		mockRepo.AssertExpectations(t)
	})
}

func testNonBlockingBehavior(t *testing.T) {
	mockRepo, recorder := setupTestRecorder(t)
	ctx := context.Background()

	t.Run("RecordEmailSentIsNonBlocking", func(t *testing.T) {
		// Set up mock to simulate slow repository operation
		mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
			time.Sleep(200 * time.Millisecond) // Simulate slow operation
		}).Times(5)

		start := time.Now()

		// Record multiple events quickly
		for i := 0; i < 5; i++ {
			event := createTestEmailEvent()
			event.InquiryID = fmt.Sprintf("inquiry-nonblocking-%d", i)

			err := recorder.RecordEmailSent(ctx, event)
			assert.NoError(t, err)
		}

		elapsed := time.Since(start)

		// Should complete quickly (non-blocking)
		assert.Less(t, elapsed, 100*time.Millisecond, "RecordEmailSent should be non-blocking")

		// Wait for background operations to complete
		time.Sleep(1 * time.Second)
		mockRepo.AssertExpectations(t)
	})

	t.Run("UpdateEmailStatusIsNonBlocking", func(t *testing.T) {
		messageID := "ses-nonblocking-" + uuid.New().String()
		existingEvent := createTestEmailEvent()
		existingEvent.SESMessageID = messageID

		mockRepo.On("GetByMessageID", mock.Anything, messageID).Return(existingEvent, nil).Once()
		mockRepo.On("Update", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
			time.Sleep(200 * time.Millisecond) // Simulate slow operation
		}).Once()

		start := time.Now()

		err := recorder.UpdateEmailStatus(ctx, messageID, domain.EmailStatusDelivered, nil, "")
		assert.NoError(t, err)

		elapsed := time.Since(start)

		// Should complete quickly (non-blocking)
		assert.Less(t, elapsed, 100*time.Millisecond, "UpdateEmailStatus should be non-blocking")

		// Wait for background operations to complete
		time.Sleep(500 * time.Millisecond)
		mockRepo.AssertExpectations(t)
	})
}

func testRetryLogic(t *testing.T) {
	mockRepo, recorder := setupTestRecorder(t)
	ctx := context.Background()

	t.Run("RetryOnTransientFailure", func(t *testing.T) {
		event := createTestEmailEvent()

		// First two calls fail, third succeeds
		mockRepo.On("Create", mock.Anything, mock.Anything).Return(fmt.Errorf("connection timeout")).Twice()
		mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil).Once()

		err := recorder.RecordEmailSent(ctx, event)
		assert.NoError(t, err) // Should not return error (non-blocking)

		// Wait for retries to complete
		time.Sleep(10 * time.Second) // Allow time for exponential backoff
		mockRepo.AssertExpectations(t)
	})

	t.Run("MaxRetriesExceeded", func(t *testing.T) {
		event := createTestEmailEvent()

		// All attempts fail
		mockRepo.On("Create", mock.Anything, mock.Anything).Return(fmt.Errorf("persistent error")).Times(3)

		err := recorder.RecordEmailSent(ctx, event)
		assert.NoError(t, err) // Should not return error (non-blocking)

		// Wait for all retries to complete
		time.Sleep(15 * time.Second) // Allow time for all retries with exponential backoff
		mockRepo.AssertExpectations(t)
	})
}

func testHealthCheck(t *testing.T) {
	mockRepo, recorder := setupTestRecorder(t)

	// Cast to concrete type to access health check methods
	concreteRecorder, ok := recorder.(*services.EmailEventRecorderImpl)
	require.True(t, ok, "Recorder should be concrete implementation")

	t.Run("HealthyService", func(t *testing.T) {
		// Mock successful health check
		mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(e *domain.EmailEvent) bool {
			return e.InquiryID == "health-check"
		})).Return(nil).Once()

		isHealthy := concreteRecorder.IsHealthy()
		assert.True(t, isHealthy)

		mockRepo.AssertExpectations(t)
	})

	t.Run("UnhealthyService", func(t *testing.T) {
		// Mock failed health check
		mockRepo.On("Create", mock.Anything, mock.Anything).Return(fmt.Errorf("database connection failed")).Once()

		isHealthy := concreteRecorder.IsHealthy()
		assert.False(t, isHealthy)

		mockRepo.AssertExpectations(t)
	})

	t.Run("HealthyWithContext", func(t *testing.T) {
		ctx := context.Background()

		// Mock successful health check
		mockRepo.On("GetByMessageID", ctx, "health-check-non-existent").Return(nil, nil).Once()

		isHealthy := concreteRecorder.IsHealthyWithContext(ctx)
		assert.True(t, isHealthy)

		mockRepo.AssertExpectations(t)
	})

	t.Run("UnhealthyWithContext", func(t *testing.T) {
		ctx := context.Background()

		// Mock connection error
		mockRepo.On("GetByMessageID", ctx, "health-check-non-existent").Return(nil, fmt.Errorf("connection refused")).Once()

		isHealthy := concreteRecorder.IsHealthyWithContext(ctx)
		assert.False(t, isHealthy)

		mockRepo.AssertExpectations(t)
	})
}

func testConcurrentOperations(t *testing.T) {
	mockRepo, recorder := setupTestRecorder(t)
	ctx := context.Background()

	t.Run("ConcurrentRecordEmailSent", func(t *testing.T) {
		const numGoroutines = 10

		// Set up mocks for concurrent operations
		mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil).Times(numGoroutines)

		var wg sync.WaitGroup
		errors := make(chan error, numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(index int) {
				defer wg.Done()

				event := createTestEmailEvent()
				event.InquiryID = fmt.Sprintf("inquiry-concurrent-%d", index)

				err := recorder.RecordEmailSent(ctx, event)
				errors <- err
			}(i)
		}

		wg.Wait()
		close(errors)

		// Check all operations succeeded
		for err := range errors {
			assert.NoError(t, err)
		}

		// Wait for async operations to complete
		time.Sleep(500 * time.Millisecond)
		mockRepo.AssertExpectations(t)
	})

	t.Run("ConcurrentUpdateEmailStatus", func(t *testing.T) {
		const numGoroutines = 5

		// Create events first
		events := make([]*domain.EmailEvent, numGoroutines)
		for i := 0; i < numGoroutines; i++ {
			events[i] = createTestEmailEvent()
			events[i].SESMessageID = fmt.Sprintf("ses-concurrent-%d", i)

			mockRepo.On("GetByMessageID", mock.Anything, events[i].SESMessageID).Return(events[i], nil).Once()
			mockRepo.On("Update", mock.Anything, mock.Anything).Return(nil).Once()
		}

		var wg sync.WaitGroup
		errors := make(chan error, numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(index int) {
				defer wg.Done()

				err := recorder.UpdateEmailStatus(ctx, events[index].SESMessageID, domain.EmailStatusDelivered, nil, "")
				errors <- err
			}(i)
		}

		wg.Wait()
		close(errors)

		// Check all operations succeeded
		for err := range errors {
			assert.NoError(t, err)
		}

		// Wait for async operations to complete
		time.Sleep(500 * time.Millisecond)
		mockRepo.AssertExpectations(t)
	})
}

// Helper functions

func createTestEmailEvent() *domain.EmailEvent {
	return &domain.EmailEvent{
		InquiryID:      "inquiry-" + uuid.New().String(),
		EmailType:      domain.EmailTypeCustomerConfirmation,
		RecipientEmail: "test@example.com",
		SenderEmail:    "info@cloudpartner.pro",
		Subject:        "Test email",
		Status:         domain.EmailStatusSent,
		SESMessageID:   "ses-" + uuid.New().String(),
	}
}

// Main function for running tests standalone
func main() {
	fmt.Println("=== Comprehensive Email Event Recorder Tests ===")

	// Note: This would normally be run with `go test` command
	// This main function is for demonstration purposes

	fmt.Println("Run with: go test -v ./test_email_event_recorder_comprehensive.go")
	fmt.Println("Or integrate into your test suite")
}
