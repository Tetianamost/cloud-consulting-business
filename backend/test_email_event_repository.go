package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"

	"github.com/cloud-consulting/backend/internal/domain"
	"github.com/cloud-consulting/backend/internal/interfaces"
	"github.com/cloud-consulting/backend/internal/repositories"
)

func main() {
	fmt.Println("=== Email Event Repository Test ===")

	// Create logger
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	// Create in-memory database for testing
	db, err := sql.Open("postgres", "postgres://user:password@localhost/testdb?sslmode=disable")
	if err != nil {
		// Fallback to in-memory SQLite for testing
		db, err = sql.Open("sqlite3", ":memory:")
		if err != nil {
			log.Fatalf("Failed to create test database: %v", err)
		}
	}
	defer db.Close()

	// Create the email_events table for testing
	err = createTestTable(db)
	if err != nil {
		log.Fatalf("Failed to create test table: %v", err)
	}

	// Create repository
	repo := repositories.NewEmailEventRepository(db, logger)

	// Run tests
	ctx := context.Background()

	fmt.Println("\n1. Testing Create operation...")
	testCreate(ctx, repo)

	fmt.Println("\n2. Testing GetByInquiryID operation...")
	testGetByInquiryID(ctx, repo)

	fmt.Println("\n3. Testing GetByMessageID operation...")
	testGetByMessageID(ctx, repo)

	fmt.Println("\n4. Testing Update operation...")
	testUpdate(ctx, repo)

	fmt.Println("\n5. Testing GetMetrics operation...")
	testGetMetrics(ctx, repo)

	fmt.Println("\n6. Testing List operation with filters...")
	testList(ctx, repo)

	fmt.Println("\n=== All tests completed successfully! ===")
}

func createTestTable(db *sql.DB) error {
	// Create a simplified version of the email_events table for testing
	query := `
	CREATE TABLE IF NOT EXISTS email_events (
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
	)`

	_, err := db.Exec(query)
	return err
}

func testCreate(ctx context.Context, repo interfaces.EmailEventRepository) {
	event := &domain.EmailEvent{
		ID:             uuid.New().String(),
		InquiryID:      "inquiry-123",
		EmailType:      domain.EmailTypeCustomerConfirmation,
		RecipientEmail: "customer@example.com",
		SenderEmail:    "info@cloudpartner.pro",
		Subject:        "Thank you for your inquiry",
		Status:         domain.EmailStatusSent,
		SentAt:         time.Now(),
		SESMessageID:   "ses-msg-123",
	}

	err := repo.Create(ctx, event)
	if err != nil {
		log.Fatalf("Failed to create email event: %v", err)
	}

	fmt.Printf("✓ Created email event: %s\n", event.ID)
}

func testGetByInquiryID(ctx context.Context, repo interfaces.EmailEventRepository) {
	// Create test events
	inquiryID := "inquiry-456"
	events := []*domain.EmailEvent{
		{
			ID:             uuid.New().String(),
			InquiryID:      inquiryID,
			EmailType:      domain.EmailTypeCustomerConfirmation,
			RecipientEmail: "customer@example.com",
			SenderEmail:    "info@cloudpartner.pro",
			Subject:        "Customer confirmation",
			Status:         domain.EmailStatusSent,
			SentAt:         time.Now(),
		},
		{
			ID:             uuid.New().String(),
			InquiryID:      inquiryID,
			EmailType:      domain.EmailTypeConsultantNotification,
			RecipientEmail: "consultant@cloudpartner.pro",
			SenderEmail:    "info@cloudpartner.pro",
			Subject:        "New inquiry received",
			Status:         domain.EmailStatusDelivered,
			SentAt:         time.Now(),
			DeliveredAt:    timePtr(time.Now().Add(5 * time.Minute)),
		},
	}

	// Create events
	for _, event := range events {
		err := repo.Create(ctx, event)
		if err != nil {
			log.Fatalf("Failed to create test event: %v", err)
		}
	}

	// Retrieve events by inquiry ID
	retrievedEvents, err := repo.GetByInquiryID(ctx, inquiryID)
	if err != nil {
		log.Fatalf("Failed to get events by inquiry ID: %v", err)
	}

	if len(retrievedEvents) != 2 {
		log.Fatalf("Expected 2 events, got %d", len(retrievedEvents))
	}

	fmt.Printf("✓ Retrieved %d events for inquiry %s\n", len(retrievedEvents), inquiryID)
}

func testGetByMessageID(ctx context.Context, repo interfaces.EmailEventRepository) {
	messageID := "ses-msg-789"
	event := &domain.EmailEvent{
		ID:             uuid.New().String(),
		InquiryID:      "inquiry-789",
		EmailType:      domain.EmailTypeInquiryNotification,
		RecipientEmail: "admin@cloudpartner.pro",
		SenderEmail:    "info@cloudpartner.pro",
		Subject:        "Inquiry notification",
		Status:         domain.EmailStatusSent,
		SentAt:         time.Now(),
		SESMessageID:   messageID,
	}

	// Create event
	err := repo.Create(ctx, event)
	if err != nil {
		log.Fatalf("Failed to create test event: %v", err)
	}

	// Retrieve by message ID
	retrievedEvent, err := repo.GetByMessageID(ctx, messageID)
	if err != nil {
		log.Fatalf("Failed to get event by message ID: %v", err)
	}

	if retrievedEvent == nil {
		log.Fatalf("Expected event, got nil")
	}

	if retrievedEvent.SESMessageID != messageID {
		log.Fatalf("Expected message ID %s, got %s", messageID, retrievedEvent.SESMessageID)
	}

	fmt.Printf("✓ Retrieved event by message ID: %s\n", messageID)
}

func testUpdate(ctx context.Context, repo interfaces.EmailEventRepository) {
	// Create initial event
	event := &domain.EmailEvent{
		ID:             uuid.New().String(),
		InquiryID:      "inquiry-update",
		EmailType:      domain.EmailTypeCustomerConfirmation,
		RecipientEmail: "update@example.com",
		SenderEmail:    "info@cloudpartner.pro",
		Subject:        "Update test",
		Status:         domain.EmailStatusSent,
		SentAt:         time.Now(),
		SESMessageID:   "ses-update-123",
	}

	err := repo.Create(ctx, event)
	if err != nil {
		log.Fatalf("Failed to create event for update test: %v", err)
	}

	// Update the event
	event.Status = domain.EmailStatusDelivered
	event.DeliveredAt = timePtr(time.Now())
	event.ErrorMessage = ""

	err = repo.Update(ctx, event)
	if err != nil {
		log.Fatalf("Failed to update event: %v", err)
	}

	// Verify update
	updatedEvent, err := repo.GetByMessageID(ctx, event.SESMessageID)
	if err != nil {
		log.Fatalf("Failed to get updated event: %v", err)
	}

	if updatedEvent.Status != domain.EmailStatusDelivered {
		log.Fatalf("Expected status %s, got %s", domain.EmailStatusDelivered, updatedEvent.Status)
	}

	fmt.Printf("✓ Updated event status to: %s\n", updatedEvent.Status)
}

func testGetMetrics(ctx context.Context, repo interfaces.EmailEventRepository) {
	// Create test events with different statuses
	inquiryID := "inquiry-metrics"
	testEvents := []*domain.EmailEvent{
		{
			ID:             uuid.New().String(),
			InquiryID:      inquiryID,
			EmailType:      domain.EmailTypeCustomerConfirmation,
			RecipientEmail: "metrics1@example.com",
			SenderEmail:    "info@cloudpartner.pro",
			Status:         domain.EmailStatusDelivered,
			SentAt:         time.Now(),
		},
		{
			ID:             uuid.New().String(),
			InquiryID:      inquiryID,
			EmailType:      domain.EmailTypeConsultantNotification,
			RecipientEmail: "metrics2@example.com",
			SenderEmail:    "info@cloudpartner.pro",
			Status:         domain.EmailStatusFailed,
			SentAt:         time.Now(),
			ErrorMessage:   "Test failure",
		},
		{
			ID:             uuid.New().String(),
			InquiryID:      inquiryID,
			EmailType:      domain.EmailTypeInquiryNotification,
			RecipientEmail: "metrics3@example.com",
			SenderEmail:    "info@cloudpartner.pro",
			Status:         domain.EmailStatusBounced,
			SentAt:         time.Now(),
			BounceType:     "permanent",
		},
	}

	// Create test events
	for _, event := range testEvents {
		err := repo.Create(ctx, event)
		if err != nil {
			log.Fatalf("Failed to create metrics test event: %v", err)
		}
	}

	// Get metrics
	filters := domain.EmailEventFilters{
		TimeRange: &domain.TimeRange{
			Start: time.Now().Add(-1 * time.Hour),
			End:   time.Now().Add(1 * time.Hour),
		},
	}

	metrics, err := repo.GetMetrics(ctx, filters)
	if err != nil {
		log.Fatalf("Failed to get metrics: %v", err)
	}

	fmt.Printf("✓ Metrics calculated - Total: %d, Delivered: %d, Failed: %d, Bounced: %d\n",
		metrics.TotalEmails, metrics.DeliveredEmails, metrics.FailedEmails, metrics.BouncedEmails)
	fmt.Printf("  Delivery Rate: %.2f%%, Bounce Rate: %.2f%%\n",
		metrics.DeliveryRate, metrics.BounceRate)
}

func testList(ctx context.Context, repo interfaces.EmailEventRepository) {
	// Test listing with filters
	filters := domain.EmailEventFilters{
		EmailType: emailTypePtr(domain.EmailTypeCustomerConfirmation),
		Limit:     10,
		Offset:    0,
	}

	events, err := repo.List(ctx, filters)
	if err != nil {
		log.Fatalf("Failed to list events: %v", err)
	}

	fmt.Printf("✓ Listed %d customer confirmation events\n", len(events))

	// Test listing with status filter
	statusFilters := domain.EmailEventFilters{
		Status: emailStatusPtr(domain.EmailStatusDelivered),
		Limit:  5,
	}

	deliveredEvents, err := repo.List(ctx, statusFilters)
	if err != nil {
		log.Fatalf("Failed to list delivered events: %v", err)
	}

	fmt.Printf("✓ Listed %d delivered events\n", len(deliveredEvents))
}

// Helper functions
func timePtr(t time.Time) *time.Time {
	return &t
}

func emailTypePtr(et domain.EmailEventType) *domain.EmailEventType {
	return &et
}

func emailStatusPtr(es domain.EmailEventStatus) *domain.EmailEventStatus {
	return &es
}
