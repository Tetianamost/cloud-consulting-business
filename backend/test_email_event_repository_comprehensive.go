package main

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cloud-consulting/backend/internal/domain"
	"github.com/cloud-consulting/backend/internal/interfaces"
	"github.com/cloud-consulting/backend/internal/repositories"
)

// TestEmailEventRepository provides comprehensive unit tests for email event repository operations
func TestEmailEventRepository(t *testing.T) {
	// Setup test database and repository
	db, repo := setupTestRepository(t)
	defer db.Close()

	ctx := context.Background()

	t.Run("Create", func(t *testing.T) {
		testCreateEmailEvent(t, ctx, repo)
	})

	t.Run("Update", func(t *testing.T) {
		testUpdateEmailEvent(t, ctx, repo)
	})

	t.Run("GetByInquiryID", func(t *testing.T) {
		testGetByInquiryID(t, ctx, repo)
	})

	t.Run("GetByMessageID", func(t *testing.T) {
		testGetByMessageID(t, ctx, repo)
	})

	t.Run("GetMetrics", func(t *testing.T) {
		testGetMetrics(t, ctx, repo)
	})

	t.Run("List", func(t *testing.T) {
		testListEmailEvents(t, ctx, repo)
	})

	t.Run("ErrorHandling", func(t *testing.T) {
		testErrorHandling(t, ctx, repo)
	})

	t.Run("EdgeCases", func(t *testing.T) {
		testEdgeCases(t, ctx, repo)
	})
}

func setupTestRepository(t *testing.T) (*sql.DB, interfaces.EmailEventRepository) {
	// Create logger
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel) // Reduce noise in tests

	// Create test database connection
	db, err := sql.Open("postgres", "postgres://test:test@localhost/test_email_events?sslmode=disable")
	if err != nil {
		// Fallback to in-memory SQLite for CI/CD environments
		t.Skip("PostgreSQL not available, skipping database tests")
	}

	// Test connection
	if err := db.Ping(); err != nil {
		t.Skip("Database connection failed, skipping database tests")
	}

	// Create test table
	createTestEmailEventsTable(t, db)

	// Create repository
	repo := repositories.NewEmailEventRepository(db, logger)
	require.NotNil(t, repo)

	return db, repo
}

func createTestEmailEventsTable(t *testing.T, db *sql.DB) {
	query := `
	DROP TABLE IF EXISTS email_events;
	CREATE TABLE email_events (
		id TEXT PRIMARY KEY,
		inquiry_id TEXT NOT NULL,
		email_type TEXT NOT NULL,
		recipient_email TEXT NOT NULL,
		sender_email TEXT NOT NULL,
		subject TEXT,
		status TEXT NOT NULL DEFAULT 'sent',
		sent_at TIMESTAMP NOT NULL,
		delivered_at TIMESTAMP,
		error_message TEXT,
		bounce_type TEXT,
		ses_message_id TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	
	CREATE INDEX IF NOT EXISTS idx_email_events_inquiry_id ON email_events(inquiry_id);
	CREATE INDEX IF NOT EXISTS idx_email_events_status ON email_events(status);
	CREATE INDEX IF NOT EXISTS idx_email_events_sent_at ON email_events(sent_at);
	CREATE INDEX IF NOT EXISTS idx_email_events_email_type ON email_events(email_type);
	CREATE INDEX IF NOT EXISTS idx_email_events_ses_message_id ON email_events(ses_message_id);
	`

	_, err := db.Exec(query)
	require.NoError(t, err, "Failed to create test table")
}

func testCreateEmailEvent(t *testing.T, ctx context.Context, repo interfaces.EmailEventRepository) {
	t.Run("ValidEvent", func(t *testing.T) {
		event := createTestEmailEvent()

		err := repo.Create(ctx, event)
		assert.NoError(t, err)
		assert.NotEmpty(t, event.ID)
		assert.False(t, event.CreatedAt.IsZero())
		assert.False(t, event.UpdatedAt.IsZero())
	})

	t.Run("EventWithAllFields", func(t *testing.T) {
		now := time.Now()
		deliveredAt := now.Add(5 * time.Minute)

		event := &domain.EmailEvent{
			ID:             uuid.New().String(),
			InquiryID:      "inquiry-complete-" + uuid.New().String(),
			EmailType:      domain.EmailTypeConsultantNotification,
			RecipientEmail: "consultant@cloudpartner.pro",
			SenderEmail:    "info@cloudpartner.pro",
			Subject:        "Complete test email",
			Status:         domain.EmailStatusDelivered,
			SentAt:         now,
			DeliveredAt:    &deliveredAt,
			ErrorMessage:   "",
			BounceType:     "",
			SESMessageID:   "ses-complete-" + uuid.New().String(),
			CreatedAt:      now,
			UpdatedAt:      now,
		}

		err := repo.Create(ctx, event)
		assert.NoError(t, err)
	})

	t.Run("DuplicateID", func(t *testing.T) {
		event1 := createTestEmailEvent()
		event2 := createTestEmailEvent()
		event2.ID = event1.ID // Same ID

		err := repo.Create(ctx, event1)
		assert.NoError(t, err)

		err = repo.Create(ctx, event2)
		assert.Error(t, err, "Should fail with duplicate ID")
	})
}

func testUpdateEmailEvent(t *testing.T, ctx context.Context, repo interfaces.EmailEventRepository) {
	t.Run("ValidUpdate", func(t *testing.T) {
		// Create initial event
		event := createTestEmailEvent()
		err := repo.Create(ctx, event)
		require.NoError(t, err)

		// Update the event
		originalStatus := event.Status
		event.Status = domain.EmailStatusDelivered
		deliveredAt := time.Now()
		event.DeliveredAt = &deliveredAt
		event.ErrorMessage = ""

		err = repo.Update(ctx, event)
		assert.NoError(t, err)
		assert.NotEqual(t, originalStatus, event.Status)
		assert.NotNil(t, event.DeliveredAt)
	})

	t.Run("UpdateNonExistentEvent", func(t *testing.T) {
		event := createTestEmailEvent()
		event.ID = "non-existent-" + uuid.New().String()

		err := repo.Update(ctx, event)
		assert.Error(t, err, "Should fail when updating non-existent event")
	})

	t.Run("UpdateToFailedStatus", func(t *testing.T) {
		// Create initial event
		event := createTestEmailEvent()
		err := repo.Create(ctx, event)
		require.NoError(t, err)

		// Update to failed status
		event.Status = domain.EmailStatusFailed
		event.ErrorMessage = "SMTP connection failed"

		err = repo.Update(ctx, event)
		assert.NoError(t, err)
		assert.Equal(t, domain.EmailStatusFailed, event.Status)
		assert.Equal(t, "SMTP connection failed", event.ErrorMessage)
	})

	t.Run("UpdateToBouncedStatus", func(t *testing.T) {
		// Create initial event
		event := createTestEmailEvent()
		err := repo.Create(ctx, event)
		require.NoError(t, err)

		// Update to bounced status
		event.Status = domain.EmailStatusBounced
		event.BounceType = "permanent"
		event.ErrorMessage = "Email address does not exist"

		err = repo.Update(ctx, event)
		assert.NoError(t, err)
		assert.Equal(t, domain.EmailStatusBounced, event.Status)
		assert.Equal(t, "permanent", event.BounceType)
		assert.Equal(t, "Email address does not exist", event.ErrorMessage)
	})
}

func testGetByInquiryID(t *testing.T, ctx context.Context, repo interfaces.EmailEventRepository) {
	inquiryID := "inquiry-get-test-" + uuid.New().String()

	t.Run("MultipleEventsForInquiry", func(t *testing.T) {
		// Create multiple events for the same inquiry
		events := []*domain.EmailEvent{
			{
				ID:             uuid.New().String(),
				InquiryID:      inquiryID,
				EmailType:      domain.EmailTypeCustomerConfirmation,
				RecipientEmail: "customer@example.com",
				SenderEmail:    "info@cloudpartner.pro",
				Subject:        "Customer confirmation",
				Status:         domain.EmailStatusSent,
				SentAt:         time.Now().Add(-2 * time.Hour),
				SESMessageID:   "ses-customer-" + uuid.New().String(),
			},
			{
				ID:             uuid.New().String(),
				InquiryID:      inquiryID,
				EmailType:      domain.EmailTypeConsultantNotification,
				RecipientEmail: "consultant@cloudpartner.pro",
				SenderEmail:    "info@cloudpartner.pro",
				Subject:        "New inquiry received",
				Status:         domain.EmailStatusDelivered,
				SentAt:         time.Now().Add(-1 * time.Hour),
				SESMessageID:   "ses-consultant-" + uuid.New().String(),
			},
		}

		// Create events
		for _, event := range events {
			err := repo.Create(ctx, event)
			require.NoError(t, err)
		}

		// Retrieve events
		retrievedEvents, err := repo.GetByInquiryID(ctx, inquiryID)
		assert.NoError(t, err)
		assert.Len(t, retrievedEvents, 2)

		// Verify events are ordered by sent_at DESC (most recent first)
		assert.True(t, retrievedEvents[0].SentAt.After(retrievedEvents[1].SentAt) ||
			retrievedEvents[0].SentAt.Equal(retrievedEvents[1].SentAt))
	})

	t.Run("NoEventsForInquiry", func(t *testing.T) {
		nonExistentInquiryID := "non-existent-" + uuid.New().String()

		events, err := repo.GetByInquiryID(ctx, nonExistentInquiryID)
		assert.NoError(t, err)
		assert.Empty(t, events)
	})

	t.Run("EmptyInquiryID", func(t *testing.T) {
		events, err := repo.GetByInquiryID(ctx, "")
		assert.NoError(t, err)
		assert.Empty(t, events)
	})
}

func testGetByMessageID(t *testing.T, ctx context.Context, repo interfaces.EmailEventRepository) {
	t.Run("ExistingMessageID", func(t *testing.T) {
		messageID := "ses-get-test-" + uuid.New().String()
		event := createTestEmailEvent()
		event.SESMessageID = messageID

		err := repo.Create(ctx, event)
		require.NoError(t, err)

		retrievedEvent, err := repo.GetByMessageID(ctx, messageID)
		assert.NoError(t, err)
		assert.NotNil(t, retrievedEvent)
		assert.Equal(t, messageID, retrievedEvent.SESMessageID)
		assert.Equal(t, event.ID, retrievedEvent.ID)
	})

	t.Run("NonExistentMessageID", func(t *testing.T) {
		nonExistentMessageID := "non-existent-" + uuid.New().String()

		event, err := repo.GetByMessageID(ctx, nonExistentMessageID)
		assert.NoError(t, err)
		assert.Nil(t, event)
	})

	t.Run("EmptyMessageID", func(t *testing.T) {
		event, err := repo.GetByMessageID(ctx, "")
		assert.NoError(t, err)
		assert.Nil(t, event)
	})
}

func testGetMetrics(t *testing.T, ctx context.Context, repo interfaces.EmailEventRepository) {
	// Create test events with different statuses
	now := time.Now()
	inquiryID := "inquiry-metrics-" + uuid.New().String()

	testEvents := []*domain.EmailEvent{
		{
			ID:             uuid.New().String(),
			InquiryID:      inquiryID,
			EmailType:      domain.EmailTypeCustomerConfirmation,
			RecipientEmail: "delivered@example.com",
			SenderEmail:    "info@cloudpartner.pro",
			Status:         domain.EmailStatusDelivered,
			SentAt:         now.Add(-1 * time.Hour),
			DeliveredAt:    timePtr(now.Add(-50 * time.Minute)),
		},
		{
			ID:             uuid.New().String(),
			InquiryID:      inquiryID,
			EmailType:      domain.EmailTypeConsultantNotification,
			RecipientEmail: "failed@example.com",
			SenderEmail:    "info@cloudpartner.pro",
			Status:         domain.EmailStatusFailed,
			SentAt:         now.Add(-2 * time.Hour),
			ErrorMessage:   "SMTP connection failed",
		},
		{
			ID:             uuid.New().String(),
			InquiryID:      inquiryID,
			EmailType:      domain.EmailTypeInquiryNotification,
			RecipientEmail: "bounced@example.com",
			SenderEmail:    "info@cloudpartner.pro",
			Status:         domain.EmailStatusBounced,
			SentAt:         now.Add(-3 * time.Hour),
			BounceType:     "permanent",
			ErrorMessage:   "Email address does not exist",
		},
		{
			ID:             uuid.New().String(),
			InquiryID:      inquiryID,
			EmailType:      domain.EmailTypeCustomerConfirmation,
			RecipientEmail: "spam@example.com",
			SenderEmail:    "info@cloudpartner.pro",
			Status:         domain.EmailStatusSpam,
			SentAt:         now.Add(-4 * time.Hour),
			ErrorMessage:   "Marked as spam",
		},
		{
			ID:             uuid.New().String(),
			InquiryID:      inquiryID,
			EmailType:      domain.EmailTypeConsultantNotification,
			RecipientEmail: "sent@example.com",
			SenderEmail:    "info@cloudpartner.pro",
			Status:         domain.EmailStatusSent,
			SentAt:         now.Add(-30 * time.Minute),
		},
	}

	// Create all test events
	for _, event := range testEvents {
		err := repo.Create(ctx, event)
		require.NoError(t, err)
	}

	t.Run("AllTimeMetrics", func(t *testing.T) {
		filters := domain.EmailEventFilters{
			TimeRange: &domain.TimeRange{
				Start: now.Add(-5 * time.Hour),
				End:   now,
			},
		}

		metrics, err := repo.GetMetrics(ctx, filters)
		assert.NoError(t, err)
		assert.NotNil(t, metrics)

		// Verify counts
		assert.Equal(t, int64(5), metrics.TotalEmails)
		assert.Equal(t, int64(1), metrics.DeliveredEmails)
		assert.Equal(t, int64(1), metrics.FailedEmails)
		assert.Equal(t, int64(1), metrics.BouncedEmails)
		assert.Equal(t, int64(1), metrics.SpamEmails)

		// Verify rates
		assert.Equal(t, 20.0, metrics.DeliveryRate) // 1/5 * 100
		assert.Equal(t, 20.0, metrics.BounceRate)   // 1/5 * 100
		assert.Equal(t, 20.0, metrics.SpamRate)     // 1/5 * 100
	})

	t.Run("FilterByEmailType", func(t *testing.T) {
		customerEmailType := domain.EmailTypeCustomerConfirmation
		filters := domain.EmailEventFilters{
			TimeRange: &domain.TimeRange{
				Start: now.Add(-5 * time.Hour),
				End:   now,
			},
			EmailType: &customerEmailType,
		}

		metrics, err := repo.GetMetrics(ctx, filters)
		assert.NoError(t, err)
		assert.Equal(t, int64(2), metrics.TotalEmails) // 2 customer confirmation emails
		assert.Equal(t, int64(1), metrics.DeliveredEmails)
		assert.Equal(t, int64(1), metrics.SpamEmails)
	})

	t.Run("FilterByStatus", func(t *testing.T) {
		deliveredStatus := domain.EmailStatusDelivered
		filters := domain.EmailEventFilters{
			TimeRange: &domain.TimeRange{
				Start: now.Add(-5 * time.Hour),
				End:   now,
			},
			Status: &deliveredStatus,
		}

		metrics, err := repo.GetMetrics(ctx, filters)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), metrics.TotalEmails)
		assert.Equal(t, int64(1), metrics.DeliveredEmails)
		assert.Equal(t, 100.0, metrics.DeliveryRate)
	})

	t.Run("FilterByInquiryID", func(t *testing.T) {
		filters := domain.EmailEventFilters{
			TimeRange: &domain.TimeRange{
				Start: now.Add(-5 * time.Hour),
				End:   now,
			},
			InquiryID: &inquiryID,
		}

		metrics, err := repo.GetMetrics(ctx, filters)
		assert.NoError(t, err)
		assert.Equal(t, int64(5), metrics.TotalEmails)
	})

	t.Run("EmptyTimeRange", func(t *testing.T) {
		filters := domain.EmailEventFilters{
			TimeRange: &domain.TimeRange{
				Start: now.Add(1 * time.Hour), // Future time range
				End:   now.Add(2 * time.Hour),
			},
		}

		metrics, err := repo.GetMetrics(ctx, filters)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), metrics.TotalEmails)
		assert.Equal(t, 0.0, metrics.DeliveryRate)
	})
}

func testListEmailEvents(t *testing.T, ctx context.Context, repo interfaces.EmailEventRepository) {
	// Create test events for listing
	now := time.Now()
	inquiryID := "inquiry-list-" + uuid.New().String()

	testEvents := []*domain.EmailEvent{
		{
			ID:             uuid.New().String(),
			InquiryID:      inquiryID,
			EmailType:      domain.EmailTypeCustomerConfirmation,
			RecipientEmail: "list1@example.com",
			SenderEmail:    "info@cloudpartner.pro",
			Status:         domain.EmailStatusDelivered,
			SentAt:         now.Add(-1 * time.Hour),
		},
		{
			ID:             uuid.New().String(),
			InquiryID:      inquiryID,
			EmailType:      domain.EmailTypeConsultantNotification,
			RecipientEmail: "list2@example.com",
			SenderEmail:    "info@cloudpartner.pro",
			Status:         domain.EmailStatusFailed,
			SentAt:         now.Add(-2 * time.Hour),
		},
		{
			ID:             uuid.New().String(),
			InquiryID:      inquiryID + "-other",
			EmailType:      domain.EmailTypeInquiryNotification,
			RecipientEmail: "list3@example.com",
			SenderEmail:    "info@cloudpartner.pro",
			Status:         domain.EmailStatusSent,
			SentAt:         now.Add(-3 * time.Hour),
		},
	}

	// Create all test events
	for _, event := range testEvents {
		err := repo.Create(ctx, event)
		require.NoError(t, err)
	}

	t.Run("ListAll", func(t *testing.T) {
		filters := domain.EmailEventFilters{
			Limit:  10,
			Offset: 0,
		}

		events, err := repo.List(ctx, filters)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(events), 3) // At least our test events

		// Verify ordering (most recent first)
		for i := 1; i < len(events); i++ {
			assert.True(t, events[i-1].SentAt.After(events[i].SentAt) ||
				events[i-1].SentAt.Equal(events[i].SentAt))
		}
	})

	t.Run("ListWithPagination", func(t *testing.T) {
		filters := domain.EmailEventFilters{
			Limit:  2,
			Offset: 0,
		}

		events, err := repo.List(ctx, filters)
		assert.NoError(t, err)
		assert.LessOrEqual(t, len(events), 2)

		// Test second page
		filters.Offset = 2
		moreEvents, err := repo.List(ctx, filters)
		assert.NoError(t, err)

		// Verify no overlap
		if len(events) > 0 && len(moreEvents) > 0 {
			assert.NotEqual(t, events[0].ID, moreEvents[0].ID)
		}
	})

	t.Run("ListWithFilters", func(t *testing.T) {
		customerEmailType := domain.EmailTypeCustomerConfirmation
		filters := domain.EmailEventFilters{
			EmailType: &customerEmailType,
			Limit:     10,
			Offset:    0,
		}

		events, err := repo.List(ctx, filters)
		assert.NoError(t, err)

		// Verify all events are customer confirmation type
		for _, event := range events {
			assert.Equal(t, domain.EmailTypeCustomerConfirmation, event.EmailType)
		}
	})

	t.Run("ListWithTimeRange", func(t *testing.T) {
		filters := domain.EmailEventFilters{
			TimeRange: &domain.TimeRange{
				Start: now.Add(-2*time.Hour - 30*time.Minute),
				End:   now.Add(-30 * time.Minute),
			},
			Limit:  10,
			Offset: 0,
		}

		events, err := repo.List(ctx, filters)
		assert.NoError(t, err)

		// Verify all events are within time range
		for _, event := range events {
			assert.True(t, event.SentAt.After(filters.TimeRange.Start) ||
				event.SentAt.Equal(filters.TimeRange.Start))
			assert.True(t, event.SentAt.Before(filters.TimeRange.End) ||
				event.SentAt.Equal(filters.TimeRange.End))
		}
	})
}

func testErrorHandling(t *testing.T, ctx context.Context, repo interfaces.EmailEventRepository) {
	t.Run("CreateWithNilEvent", func(t *testing.T) {
		err := repo.Create(ctx, nil)
		assert.Error(t, err)
	})

	t.Run("UpdateWithNilEvent", func(t *testing.T) {
		err := repo.Update(ctx, nil)
		assert.Error(t, err)
	})

	t.Run("GetMetricsWithNilFilters", func(t *testing.T) {
		// This should work with default filters
		metrics, err := repo.GetMetrics(ctx, domain.EmailEventFilters{})
		assert.NoError(t, err)
		assert.NotNil(t, metrics)
	})

	t.Run("ListWithNilFilters", func(t *testing.T) {
		// This should work with default filters
		events, err := repo.List(ctx, domain.EmailEventFilters{})
		assert.NoError(t, err)
		assert.NotNil(t, events)
	})
}

func testEdgeCases(t *testing.T, ctx context.Context, repo interfaces.EmailEventRepository) {
	t.Run("EventWithEmptyOptionalFields", func(t *testing.T) {
		event := &domain.EmailEvent{
			ID:             uuid.New().String(),
			InquiryID:      "inquiry-empty-" + uuid.New().String(),
			EmailType:      domain.EmailTypeCustomerConfirmation,
			RecipientEmail: "empty@example.com",
			SenderEmail:    "info@cloudpartner.pro",
			Subject:        "", // Empty subject
			Status:         domain.EmailStatusSent,
			SentAt:         time.Now(),
			DeliveredAt:    nil, // Nil delivered at
			ErrorMessage:   "",  // Empty error message
			BounceType:     "",  // Empty bounce type
			SESMessageID:   "",  // Empty SES message ID
		}

		err := repo.Create(ctx, event)
		assert.NoError(t, err)

		// Verify retrieval
		retrievedEvent, err := repo.GetByInquiryID(ctx, event.InquiryID)
		assert.NoError(t, err)
		assert.Len(t, retrievedEvent, 1)
		assert.Equal(t, event.ID, retrievedEvent[0].ID)
	})

	t.Run("EventWithVeryLongFields", func(t *testing.T) {
		longString := string(make([]byte, 1000)) // Very long string
		for i := range longString {
			longString = longString[:i] + "a" + longString[i+1:]
		}

		event := &domain.EmailEvent{
			ID:             uuid.New().String(),
			InquiryID:      "inquiry-long-" + uuid.New().String(),
			EmailType:      domain.EmailTypeCustomerConfirmation,
			RecipientEmail: "long@example.com",
			SenderEmail:    "info@cloudpartner.pro",
			Subject:        longString[:500], // Truncate to reasonable length
			Status:         domain.EmailStatusSent,
			SentAt:         time.Now(),
			ErrorMessage:   longString[:500], // Truncate to reasonable length
		}

		err := repo.Create(ctx, event)
		assert.NoError(t, err)
	})

	t.Run("ConcurrentOperations", func(t *testing.T) {
		// Test concurrent creates
		const numGoroutines = 10
		errors := make(chan error, numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			go func(index int) {
				event := &domain.EmailEvent{
					ID:             fmt.Sprintf("concurrent-%d-%s", index, uuid.New().String()),
					InquiryID:      fmt.Sprintf("inquiry-concurrent-%d", index),
					EmailType:      domain.EmailTypeCustomerConfirmation,
					RecipientEmail: fmt.Sprintf("concurrent%d@example.com", index),
					SenderEmail:    "info@cloudpartner.pro",
					Status:         domain.EmailStatusSent,
					SentAt:         time.Now(),
				}

				err := repo.Create(ctx, event)
				errors <- err
			}(i)
		}

		// Collect results
		for i := 0; i < numGoroutines; i++ {
			err := <-errors
			assert.NoError(t, err, "Concurrent operation %d failed", i)
		}
	})
}

// Helper functions

func createTestEmailEvent() *domain.EmailEvent {
	return &domain.EmailEvent{
		ID:             uuid.New().String(),
		InquiryID:      "inquiry-" + uuid.New().String(),
		EmailType:      domain.EmailTypeCustomerConfirmation,
		RecipientEmail: "test@example.com",
		SenderEmail:    "info@cloudpartner.pro",
		Subject:        "Test email",
		Status:         domain.EmailStatusSent,
		SentAt:         time.Now(),
		SESMessageID:   "ses-" + uuid.New().String(),
	}
}

func timePtr(t time.Time) *time.Time {
	return &t
}

// Main function for running tests standalone
func main() {
	fmt.Println("=== Comprehensive Email Event Repository Tests ===")

	// Note: This would normally be run with `go test` command
	// This main function is for demonstration purposes

	fmt.Println("Run with: go test -v ./test_email_event_repository_comprehensive.go")
	fmt.Println("Or integrate into your test suite")
}
