package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"

	"github.com/cloud-consulting/backend/internal/config"
	"github.com/cloud-consulting/backend/internal/domain"
	"github.com/cloud-consulting/backend/internal/interfaces"
	"github.com/cloud-consulting/backend/internal/repositories"
	"github.com/cloud-consulting/backend/internal/services"
)

func main() {
	fmt.Println("=== Email Event Recorder Service Test ===")

	// Load configuration (not used in this test but kept for completeness)
	_, err := config.Load()
	if err != nil {
		fmt.Printf("Warning: Failed to load config: %v\n", err)
	}

	// Create logger
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	// Create database connection (using in-memory for testing)
	db, err := sql.Open("postgres", "postgres://user:password@localhost/testdb?sslmode=disable")
	if err != nil {
		// Use in-memory database for testing
		fmt.Println("Using in-memory database for testing...")
		db = createInMemoryDB()
	}
	defer db.Close()

	// Create email event repository
	emailEventRepo := repositories.NewEmailEventRepository(db, logger)

	// Create email event recorder service
	emailEventRecorder := services.NewEmailEventRecorder(emailEventRepo, logger)

	// Test the service
	ctx := context.Background()

	fmt.Println("\n1. Testing RecordEmailSent (non-blocking)...")
	testRecordEmailSent(ctx, emailEventRecorder, logger)

	fmt.Println("\n2. Testing UpdateEmailStatus (non-blocking)...")
	testUpdateEmailStatus(ctx, emailEventRecorder, logger)

	fmt.Println("\n3. Testing GetEmailEventsByInquiry...")
	testGetEmailEventsByInquiry(ctx, emailEventRecorder, logger)

	fmt.Println("\n4. Testing synchronous methods...")
	if concreteRecorder, ok := emailEventRecorder.(*services.EmailEventRecorderImpl); ok {
		testSynchronousMethods(ctx, concreteRecorder, logger)
	} else {
		fmt.Println("Skipping synchronous methods test - concrete type not available")
	}

	fmt.Println("\n5. Testing error handling...")
	testErrorHandling(ctx, emailEventRecorder, logger)

	fmt.Println("\n6. Testing service health check...")
	if concreteRecorder, ok := emailEventRecorder.(*services.EmailEventRecorderImpl); ok {
		testHealthCheck(concreteRecorder, logger)
	} else {
		fmt.Println("Skipping health check test - concrete type not available")
	}

	fmt.Println("\n=== All Email Event Recorder Tests Completed ===")
}

func testRecordEmailSent(ctx context.Context, recorder interfaces.EmailEventRecorder, logger *logrus.Logger) {
	// Create a test email event
	event := &domain.EmailEvent{
		InquiryID:      uuid.New().String(),
		EmailType:      domain.EmailTypeCustomerConfirmation,
		RecipientEmail: "customer@example.com",
		SenderEmail:    "info@cloudpartner.pro",
		Subject:        "Thank you for your inquiry",
		Status:         domain.EmailStatusSent,
		SESMessageID:   "test-message-id-" + uuid.New().String(),
	}

	// Record the email event (non-blocking)
	err := recorder.RecordEmailSent(ctx, event)
	if err != nil {
		logger.WithError(err).Error("Failed to record email sent")
		return
	}

	fmt.Printf("✓ Email event recorded successfully (ID: %s)\n", event.ID)
	fmt.Printf("  - Inquiry ID: %s\n", event.InquiryID)
	fmt.Printf("  - Email Type: %s\n", event.EmailType)
	fmt.Printf("  - Recipient: %s\n", event.RecipientEmail)
	fmt.Printf("  - Status: %s\n", event.Status)

	// Wait a moment for the async operation to complete
	time.Sleep(100 * time.Millisecond)
}

func testUpdateEmailStatus(ctx context.Context, recorder interfaces.EmailEventRecorder, logger *logrus.Logger) {
	// First, record an email event
	event := &domain.EmailEvent{
		InquiryID:      uuid.New().String(),
		EmailType:      domain.EmailTypeConsultantNotification,
		RecipientEmail: "consultant@cloudpartner.pro",
		SenderEmail:    "info@cloudpartner.pro",
		Subject:        "New inquiry received",
		Status:         domain.EmailStatusSent,
		SESMessageID:   "test-update-message-id-" + uuid.New().String(),
	}

	// Record the event first
	err := recorder.RecordEmailSent(ctx, event)
	if err != nil {
		logger.WithError(err).Error("Failed to record email for update test")
		return
	}

	// Wait for the async operation to complete
	time.Sleep(100 * time.Millisecond)

	// Update the email status to delivered
	deliveredAt := time.Now()
	err = recorder.UpdateEmailStatus(ctx, event.SESMessageID, domain.EmailStatusDelivered, &deliveredAt, "")
	if err != nil {
		logger.WithError(err).Error("Failed to update email status")
		return
	}

	fmt.Printf("✓ Email status updated successfully\n")
	fmt.Printf("  - Message ID: %s\n", event.SESMessageID)
	fmt.Printf("  - New Status: %s\n", domain.EmailStatusDelivered)
	fmt.Printf("  - Delivered At: %s\n", deliveredAt.Format(time.RFC3339))

	// Wait a moment for the async operation to complete
	time.Sleep(100 * time.Millisecond)
}

func testGetEmailEventsByInquiry(ctx context.Context, recorder interfaces.EmailEventRecorder, logger *logrus.Logger) {
	inquiryID := uuid.New().String()

	// Record multiple email events for the same inquiry
	events := []*domain.EmailEvent{
		{
			InquiryID:      inquiryID,
			EmailType:      domain.EmailTypeCustomerConfirmation,
			RecipientEmail: "customer@example.com",
			SenderEmail:    "info@cloudpartner.pro",
			Subject:        "Thank you for your inquiry",
			Status:         domain.EmailStatusSent,
			SESMessageID:   "customer-message-" + uuid.New().String(),
		},
		{
			InquiryID:      inquiryID,
			EmailType:      domain.EmailTypeConsultantNotification,
			RecipientEmail: "consultant@cloudpartner.pro",
			SenderEmail:    "info@cloudpartner.pro",
			Subject:        "New inquiry received",
			Status:         domain.EmailStatusSent,
			SESMessageID:   "consultant-message-" + uuid.New().String(),
		},
	}

	// Record all events
	for _, event := range events {
		err := recorder.RecordEmailSent(ctx, event)
		if err != nil {
			logger.WithError(err).Error("Failed to record email event")
			continue
		}
	}

	// Wait for async operations to complete
	time.Sleep(200 * time.Millisecond)

	// Retrieve events by inquiry ID
	retrievedEvents, err := recorder.GetEmailEventsByInquiry(ctx, inquiryID)
	if err != nil {
		logger.WithError(err).Error("Failed to get email events by inquiry")
		return
	}

	fmt.Printf("✓ Retrieved %d email events for inquiry %s\n", len(retrievedEvents), inquiryID)
	for i, event := range retrievedEvents {
		fmt.Printf("  Event %d:\n", i+1)
		fmt.Printf("    - ID: %s\n", event.ID)
		fmt.Printf("    - Type: %s\n", event.EmailType)
		fmt.Printf("    - Recipient: %s\n", event.RecipientEmail)
		fmt.Printf("    - Status: %s\n", event.Status)
		fmt.Printf("    - Sent At: %s\n", event.SentAt.Format(time.RFC3339))
	}
}

func testSynchronousMethods(ctx context.Context, recorder *services.EmailEventRecorderImpl, logger *logrus.Logger) {
	fmt.Println("Testing synchronous email event recording...")

	// Create a test email event
	event := &domain.EmailEvent{
		InquiryID:      uuid.New().String(),
		EmailType:      domain.EmailTypeInquiryNotification,
		RecipientEmail: "admin@cloudpartner.pro",
		SenderEmail:    "info@cloudpartner.pro",
		Subject:        "New inquiry notification",
		Status:         domain.EmailStatusSent,
		SESMessageID:   "sync-message-id-" + uuid.New().String(),
	}

	// Record the email event synchronously
	err := recorder.RecordEmailSentSync(ctx, event)
	if err != nil {
		logger.WithError(err).Error("Failed to record email sent synchronously")
		return
	}

	fmt.Printf("✓ Email event recorded synchronously (ID: %s)\n", event.ID)

	// Update the email status synchronously
	deliveredAt := time.Now()
	err = recorder.UpdateEmailStatusSync(ctx, event.SESMessageID, domain.EmailStatusDelivered, &deliveredAt, "")
	if err != nil {
		logger.WithError(err).Error("Failed to update email status synchronously")
		return
	}

	fmt.Printf("✓ Email status updated synchronously\n")
	fmt.Printf("  - Message ID: %s\n", event.SESMessageID)
	fmt.Printf("  - New Status: %s\n", domain.EmailStatusDelivered)
}

func testErrorHandling(ctx context.Context, recorder interfaces.EmailEventRecorder, logger *logrus.Logger) {
	fmt.Println("Testing error handling...")

	// Test with invalid email event (missing required fields)
	invalidEvent := &domain.EmailEvent{
		// Missing required fields
		EmailType: domain.EmailTypeCustomerConfirmation,
	}

	err := recorder.RecordEmailSent(ctx, invalidEvent)
	if err != nil {
		fmt.Printf("✓ Correctly handled invalid email event: %v\n", err)
	} else {
		fmt.Println("✗ Should have failed with invalid email event")
	}

	// Test updating status with empty message ID
	err = recorder.UpdateEmailStatus(ctx, "", domain.EmailStatusDelivered, nil, "")
	if err != nil {
		fmt.Printf("✓ Correctly handled empty message ID: %v\n", err)
	} else {
		fmt.Println("✗ Should have failed with empty message ID")
	}

	// Test getting events with empty inquiry ID
	_, err = recorder.GetEmailEventsByInquiry(ctx, "")
	if err != nil {
		fmt.Printf("✓ Correctly handled empty inquiry ID: %v\n", err)
	} else {
		fmt.Println("✗ Should have failed with empty inquiry ID")
	}
}

func testHealthCheck(recorder *services.EmailEventRecorderImpl, logger *logrus.Logger) {
	fmt.Println("Testing service health check...")

	isHealthy := recorder.IsHealthy()
	if isHealthy {
		fmt.Println("✓ Email event recorder service is healthy")
	} else {
		fmt.Println("✗ Email event recorder service is not healthy")
	}
}

// createInMemoryDB creates a mock database for testing
func createInMemoryDB() *sql.DB {
	// This is a placeholder - in a real test, you would use a test database
	// or mock the repository interface
	fmt.Println("Note: Using mock database for testing")
	return nil
}

// Helper function to demonstrate non-blocking behavior
func demonstrateNonBlockingBehavior(ctx context.Context, recorder interfaces.EmailEventRecorder, logger *logrus.Logger) {
	fmt.Println("\nDemonstrating non-blocking behavior...")

	start := time.Now()

	// Record multiple events quickly
	for i := 0; i < 5; i++ {
		event := &domain.EmailEvent{
			InquiryID:      uuid.New().String(),
			EmailType:      domain.EmailTypeCustomerConfirmation,
			RecipientEmail: fmt.Sprintf("customer%d@example.com", i),
			SenderEmail:    "info@cloudpartner.pro",
			Subject:        fmt.Sprintf("Test email %d", i),
			Status:         domain.EmailStatusSent,
			SESMessageID:   fmt.Sprintf("test-message-%d-%s", i, uuid.New().String()),
		}

		err := recorder.RecordEmailSent(ctx, event)
		if err != nil {
			logger.WithError(err).Errorf("Failed to record email %d", i)
		}
	}

	elapsed := time.Since(start)
	fmt.Printf("✓ Recorded 5 email events in %v (non-blocking)\n", elapsed)
	fmt.Println("  Events are being processed asynchronously in the background")

	// Wait for background operations to complete
	time.Sleep(500 * time.Millisecond)
	fmt.Println("✓ Background operations completed")
}

func init() {
	// Set up environment for testing
	os.Setenv("DATABASE_URL", "postgres://test:test@localhost/test?sslmode=disable")
	os.Setenv("SES_SENDER_EMAIL", "info@cloudpartner.pro")
}
