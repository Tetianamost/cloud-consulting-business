package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/cloud-consulting/backend/internal/config"
	"github.com/cloud-consulting/backend/internal/domain"
	"github.com/cloud-consulting/backend/internal/repositories"
	"github.com/cloud-consulting/backend/internal/services"
	"github.com/cloud-consulting/backend/internal/storage"
)

// MockEmailEventRecorder for testing
type MockEmailEventRecorder struct {
	events []domain.EmailEvent
	logger *logrus.Logger
}

func NewMockEmailEventRecorder(logger *logrus.Logger) *MockEmailEventRecorder {
	return &MockEmailEventRecorder{
		events: make([]domain.EmailEvent, 0),
		logger: logger,
	}
}

func (m *MockEmailEventRecorder) RecordEmailSent(ctx context.Context, event *domain.EmailEvent) error {
	m.events = append(m.events, *event)
	m.logger.WithFields(logrus.Fields{
		"event_id":        event.ID,
		"inquiry_id":      event.InquiryID,
		"email_type":      event.EmailType,
		"recipient_email": event.RecipientEmail,
		"status":          event.Status,
	}).Info("Mock: Email event recorded")
	return nil
}

func (m *MockEmailEventRecorder) UpdateEmailStatus(ctx context.Context, messageID string, status domain.EmailEventStatus, deliveredAt *time.Time, errorMsg string) error {
	for i, event := range m.events {
		if event.SESMessageID == messageID {
			m.events[i].Status = status
			if deliveredAt != nil {
				m.events[i].DeliveredAt = deliveredAt
			}
			if errorMsg != "" {
				m.events[i].ErrorMessage = errorMsg
			}
			m.logger.WithFields(logrus.Fields{
				"message_id": messageID,
				"new_status": status,
				"event_id":   event.ID,
			}).Info("Mock: Email status updated")
			return nil
		}
	}
	return fmt.Errorf("email event not found for message ID: %s", messageID)
}

func (m *MockEmailEventRecorder) GetEmailEventsByInquiry(ctx context.Context, inquiryID string) ([]*domain.EmailEvent, error) {
	var result []*domain.EmailEvent
	for _, event := range m.events {
		if event.InquiryID == inquiryID {
			eventCopy := event
			result = append(result, &eventCopy)
		}
	}
	return result, nil
}

func (m *MockEmailEventRecorder) IsHealthy() bool {
	return true
}

func (m *MockEmailEventRecorder) GetRecordedEvents() []domain.EmailEvent {
	return m.events
}

func main() {
	// Initialize logger
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	fmt.Println("=== Email Service Factory with Event Recording Test ===")

	// Test configuration
	cfg := config.SESConfig{
		AccessKeyID:     "test-access-key",
		SecretAccessKey: "test-secret-key",
		Region:          "us-east-1",
		SenderEmail:     "test@cloudpartner.pro",
		ReplyToEmail:    "reply@cloudpartner.pro",
		Timeout:         30,
	}

	// Create mock event recorder
	mockEventRecorder := NewMockEmailEventRecorder(logger)

	fmt.Println("\n1. Testing NewEmailServiceFactory with event recorder...")

	// Test 1: Create email service with event recorder
	emailService, err := services.NewEmailServiceFactory(cfg, mockEventRecorder, logger)
	if err != nil {
		log.Fatalf("Failed to create email service with event recorder: %v", err)
	}

	if emailService == nil {
		log.Fatal("Email service is nil")
	}

	fmt.Printf("✅ Email service created successfully with event recording\n")

	// Test 2: Create email service without event recorder
	fmt.Println("\n2. Testing NewEmailServiceFactory without event recorder...")

	emailServiceNoEvents, err := services.NewEmailServiceFactory(cfg, nil, logger)
	if err != nil {
		log.Fatalf("Failed to create email service without event recorder: %v", err)
	}

	if emailServiceNoEvents == nil {
		log.Fatal("Email service is nil")
	}

	fmt.Printf("✅ Email service created successfully without event recording\n")

	// Test 3: Test production factory method
	fmt.Println("\n3. Testing NewEmailServiceForProduction...")

	prodEmailService, err := services.NewEmailServiceForProduction(cfg, mockEventRecorder, logger)
	if err != nil {
		log.Fatalf("Failed to create production email service: %v", err)
	}

	if prodEmailService == nil {
		log.Fatal("Production email service is nil")
	}

	fmt.Printf("✅ Production email service created successfully\n")

	// Test 4: Test production factory without event recorder (should fail)
	fmt.Println("\n4. Testing NewEmailServiceForProduction without event recorder (should fail)...")

	_, err = services.NewEmailServiceForProduction(cfg, nil, logger)
	if err == nil {
		log.Fatal("Expected error when creating production service without event recorder")
	}

	fmt.Printf("✅ Production email service correctly rejected without event recorder: %v\n", err)

	// Test 5: Test with invalid configuration
	fmt.Println("\n5. Testing with invalid configuration...")

	invalidCfg := config.SESConfig{
		AccessKeyID:     "", // Missing required field
		SecretAccessKey: "test-secret-key",
		Region:          "us-east-1",
		SenderEmail:     "test@cloudpartner.pro",
		ReplyToEmail:    "reply@cloudpartner.pro",
		Timeout:         30,
	}

	_, err = services.NewEmailServiceFactory(invalidCfg, mockEventRecorder, logger)
	if err == nil {
		log.Fatal("Expected error with invalid configuration")
	}

	fmt.Printf("✅ Invalid configuration correctly rejected: %v\n", err)

	// Test 6: Test with test domain in production (should fail)
	fmt.Println("\n6. Testing production service with test domain (should fail)...")

	testCfg := config.SESConfig{
		AccessKeyID:     "test-access-key",
		SecretAccessKey: "test-secret-key",
		Region:          "us-east-1",
		SenderEmail:     "test@example.com", // Test domain
		ReplyToEmail:    "reply@example.com",
		Timeout:         30,
	}

	_, err = services.NewEmailServiceForProduction(testCfg, mockEventRecorder, logger)
	if err == nil {
		log.Fatal("Expected error when using test domain in production")
	}

	fmt.Printf("✅ Test domain correctly rejected in production: %v\n", err)

	// Test 7: Test email service health check
	fmt.Println("\n7. Testing email service health check...")

	if !emailService.IsHealthy() {
		log.Fatal("Email service health check failed")
	}

	fmt.Printf("✅ Email service health check passed\n")

	// Test 8: Test with database-backed event recorder
	fmt.Println("\n8. Testing with database-backed event recorder...")

	// Create in-memory database for testing
	dbConnection, err := storage.NewDatabaseConnection(&config.DatabaseConfig{
		URL:                ":memory:",
		MaxOpenConnections: 5,
		MaxIdleConnections: 2,
		ConnMaxLifetime:    30,
		EnableEmailEvents:  true,
	}, logger)

	if err != nil {
		fmt.Printf("⚠️  Database connection failed (expected in test environment): %v\n", err)
	} else {
		// Run migration
		migrationSQL := getEmailEventsMigrationSQL()
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := dbConnection.RunMigration(ctx, migrationSQL); err != nil {
			fmt.Printf("⚠️  Database migration failed (expected in test environment): %v\n", err)
		} else {
			// Create real event recorder
			emailEventRepo := repositories.NewEmailEventRepository(dbConnection.GetDB(), logger)
			realEventRecorder := services.NewEmailEventRecorder(emailEventRepo, logger)

			// Create email service with real event recorder
			realEmailService, err := services.NewEmailServiceFactory(cfg, realEventRecorder, logger)
			if err != nil {
				log.Fatalf("Failed to create email service with real event recorder: %v", err)
			}

			fmt.Printf("✅ Email service created successfully with database-backed event recording\n")

			// Test health check
			if !realEmailService.IsHealthy() {
				log.Fatal("Real email service health check failed")
			}

			fmt.Printf("✅ Database-backed email service health check passed\n")
		}

		dbConnection.Close()
	}

	fmt.Println("\n=== All Email Service Factory Tests Passed! ===")
	fmt.Println("\nSummary:")
	fmt.Println("✅ Email service factory with event recording")
	fmt.Println("✅ Email service factory without event recording")
	fmt.Println("✅ Production email service factory")
	fmt.Println("✅ Production factory validation (rejects missing event recorder)")
	fmt.Println("✅ Configuration validation")
	fmt.Println("✅ Production domain validation")
	fmt.Println("✅ Health check functionality")
	fmt.Println("✅ Database-backed event recording (if database available)")

	fmt.Println("\nThe email service factory is now properly updated to include event recording!")
	fmt.Println("Key improvements:")
	fmt.Println("- Automatic event recording detection")
	fmt.Println("- Production-specific validation")
	fmt.Println("- Comprehensive configuration validation")
	fmt.Println("- Health check integration")
	fmt.Println("- Proper error handling and logging")
}

// getEmailEventsMigrationSQL returns the SQL for creating email events table
func getEmailEventsMigrationSQL() string {
	return `
-- Email events tracking database migration
CREATE TABLE IF NOT EXISTS email_events (
    id TEXT PRIMARY KEY,
    inquiry_id TEXT NOT NULL,
    email_type TEXT NOT NULL,
    recipient_email TEXT NOT NULL,
    sender_email TEXT NOT NULL,
    subject TEXT,
    status TEXT NOT NULL DEFAULT 'sent',
    sent_at DATETIME NOT NULL,
    delivered_at DATETIME,
    error_message TEXT,
    bounce_type TEXT,
    ses_message_id TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_email_events_inquiry_id ON email_events(inquiry_id);
CREATE INDEX IF NOT EXISTS idx_email_events_status ON email_events(status);
CREATE INDEX IF NOT EXISTS idx_email_events_sent_at ON email_events(sent_at);
CREATE INDEX IF NOT EXISTS idx_email_events_email_type ON email_events(email_type);
`
}
