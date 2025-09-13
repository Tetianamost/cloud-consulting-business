package main

import (
	"context"
	"fmt"
	"os"
	"strings"
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
}

func (m *MockEmailEventRecorder) RecordEmailSent(ctx context.Context, event *domain.EmailEvent) error {
	m.events = append(m.events, *event)
	fmt.Printf("✅ Mock: Recorded email event - Type: %s, Status: %s, Recipient: %s\n",
		event.EmailType, event.Status, event.RecipientEmail)
	return nil
}

func (m *MockEmailEventRecorder) UpdateEmailStatus(ctx context.Context, messageID string, status domain.EmailEventStatus, deliveredAt *time.Time, errorMsg string) error {
	fmt.Printf("✅ Mock: Updated email status - MessageID: %s, Status: %s\n", messageID, status)
	return nil
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

func main() {
	fmt.Println("🧪 Testing Email Service Factory Implementation")
	fmt.Println(strings.Repeat("=", 60))

	// Create logger
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	// Test 1: Configuration Validation
	fmt.Println("\n📋 Test 1: Configuration Validation")
	fmt.Println(strings.Repeat("-", 40))

	// Test invalid configuration
	invalidCfg := config.SESConfig{}
	err := services.ValidateSESConfig(invalidCfg)
	if err != nil {
		fmt.Printf("✅ Invalid config correctly rejected: %v\n", err)
	} else {
		fmt.Printf("❌ Invalid config should have been rejected\n")
	}

	// Test valid configuration
	validCfg := config.SESConfig{
		AccessKeyID:     "test-access-key",
		SecretAccessKey: "test-secret-key",
		Region:          "us-east-1",
		SenderEmail:     "test@example.com",
		ReplyToEmail:    "reply@example.com",
		Timeout:         30,
	}

	err = services.ValidateSESConfig(validCfg)
	if err != nil {
		fmt.Printf("❌ Valid config rejected: %v\n", err)
	} else {
		fmt.Printf("✅ Valid config accepted\n")
	}

	// Test 2: Email Service Factory without Event Recording
	fmt.Println("\n📋 Test 2: Email Service Factory without Event Recording")
	fmt.Println(strings.Repeat("-", 40))

	emailService, err := services.NewEmailService(validCfg, nil, logger)
	if err != nil {
		fmt.Printf("❌ Failed to create email service without event recording: %v\n", err)
	} else {
		fmt.Printf("✅ Email service created successfully without event recording\n")
		fmt.Printf("   Service healthy: %v\n", emailService.IsHealthy())
	}

	// Test 3: Email Service Factory with Mock Event Recording
	fmt.Println("\n📋 Test 3: Email Service Factory with Mock Event Recording")
	fmt.Println(strings.Repeat("-", 40))

	mockEventRecorder := &MockEmailEventRecorder{
		events: make([]domain.EmailEvent, 0),
	}

	emailServiceWithEvents, err := services.NewEmailService(validCfg, mockEventRecorder, logger)
	if err != nil {
		fmt.Printf("❌ Failed to create email service with event recording: %v\n", err)
	} else {
		fmt.Printf("✅ Email service created successfully with event recording\n")
		fmt.Printf("   Service healthy: %v\n", emailServiceWithEvents.IsHealthy())
	}

	// Test 4: Database Configuration Validation
	fmt.Println("\n📋 Test 4: Database Configuration Validation")
	fmt.Println(strings.Repeat("-", 40))

	// Test configuration with email events enabled but no database
	cfg1 := &config.Config{
		Database: config.DatabaseConfig{
			EnableEmailEvents: true,
			URL:               "", // Missing database URL
		},
		SES: validCfg,
	}

	err = cfg1.ValidateEmailEventTracking()
	if err != nil {
		fmt.Printf("✅ Configuration correctly rejected (missing database): %v\n", err)
	} else {
		fmt.Printf("❌ Configuration should have been rejected\n")
	}

	// Test valid configuration
	cfg2 := &config.Config{
		Database: config.DatabaseConfig{
			EnableEmailEvents: true,
			URL:               "postgres://test:test@localhost/test",
		},
		SES: validCfg,
	}

	err = cfg2.ValidateEmailEventTracking()
	if err != nil {
		fmt.Printf("❌ Valid configuration rejected: %v\n", err)
	} else {
		fmt.Printf("✅ Valid configuration accepted\n")
		fmt.Printf("   Email event tracking enabled: %v\n", cfg2.IsEmailEventTrackingEnabled())
	}

	// Test 5: Real Database Integration (if available)
	fmt.Println("\n📋 Test 5: Real Database Integration (if available)")
	fmt.Println(strings.Repeat("-", 40))

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL != "" {
		fmt.Printf("Database URL found, testing real integration...\n")

		dbConfig := config.DatabaseConfig{
			URL:               databaseURL,
			EnableEmailEvents: true,
		}

		dbConnection, err := storage.NewDatabaseConnection(&dbConfig, logger)
		if err != nil {
			fmt.Printf("❌ Failed to connect to database: %v\n", err)
		} else {
			fmt.Printf("✅ Database connection successful\n")

			// Create real email event recorder
			emailEventRepo := repositories.NewEmailEventRepository(dbConnection.GetDB(), logger)
			realEventRecorder := services.NewEmailEventRecorder(emailEventRepo, logger)

			// Test email service with real event recording
			realEmailService, err := services.NewEmailService(validCfg, realEventRecorder, logger)
			if err != nil {
				fmt.Printf("❌ Failed to create email service with real event recording: %v\n", err)
			} else {
				fmt.Printf("✅ Email service with real event recording created successfully\n")
			}

			dbConnection.Close()
		}
	} else {
		fmt.Printf("⚠️  No DATABASE_URL found, skipping real database integration test\n")
		fmt.Printf("   Set DATABASE_URL environment variable to test with real database\n")
	}

	// Test 6: Email Validation
	fmt.Println("\n📋 Test 6: Email Validation")
	fmt.Println(strings.Repeat("-", 40))

	testEmails := []struct {
		email string
		valid bool
	}{
		{"test@example.com", true},
		{"user.name@domain.co.uk", true},
		{"invalid-email", false},
		{"@domain.com", false},
		{"user@", false},
		{"", false},
		{"user@@domain.com", false},
	}

	for _, test := range testEmails {
		isValid := services.IsValidEmail(test.email)
		status := "✅"
		if isValid != test.valid {
			status = "❌"
		}
		fmt.Printf("%s Email '%s' - Expected: %v, Got: %v\n", status, test.email, test.valid, isValid)
	}

	fmt.Println("\n🎉 Email Service Factory Testing Complete!")
	fmt.Println(strings.Repeat("=", 60))
}
