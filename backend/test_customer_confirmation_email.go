package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/cloud-consulting/backend/internal/config"
	"github.com/cloud-consulting/backend/internal/domain"
	"github.com/cloud-consulting/backend/internal/interfaces"
	"github.com/cloud-consulting/backend/internal/services"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create logger
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	// Create template service
	templateService := services.NewTemplateService("./templates", logger)

	// Create SES service (mock for testing)
	mockSES := &mockSESService{}
	var sesService interfaces.SESService = mockSES

	// Create email service
	emailService := services.NewEmailService(sesService, templateService, cfg.SES, logger)

	// Create test inquiry
	inquiry := &domain.Inquiry{
		ID:       "test-inquiry-123",
		Name:     "John Doe",
		Email:    "john.doe@example.com",
		Company:  "Acme Corporation",
		Phone:    "+1-555-123-4567",
		Services: []string{"assessment", "migration"},
		Message:  "We need help migrating our legacy systems to AWS. This is urgent as we have a deadline next month.",
	}

	fmt.Println("Testing customer confirmation email...")
	fmt.Printf("Inquiry ID: %s\n", inquiry.ID)
	fmt.Printf("Customer: %s (%s)\n", inquiry.Name, inquiry.Email)
	fmt.Printf("Company: %s\n", inquiry.Company)
	fmt.Printf("Services: %v\n", inquiry.Services)
	fmt.Println()

	// Test customer confirmation email
	ctx := context.Background()
	err = emailService.SendCustomerConfirmation(ctx, inquiry)
	if err != nil {
		log.Fatalf("Failed to send customer confirmation email: %v", err)
	}

	fmt.Println("✅ Customer confirmation email sent successfully!")
	fmt.Println()

	// Display what was sent
	if len(mockSES.sentEmails) > 0 {
		email := mockSES.sentEmails[0]
		fmt.Println("📧 Email Details:")
		fmt.Printf("From: %s\n", email.From)
		fmt.Printf("To: %v\n", email.To)
		fmt.Printf("Subject: %s\n", email.Subject)
		fmt.Println()
		fmt.Println("📝 HTML Content Preview (first 500 chars):")
		if len(email.HTMLBody) > 500 {
			fmt.Printf("%s...\n", email.HTMLBody[:500])
		} else {
			fmt.Printf("%s\n", email.HTMLBody)
		}
		fmt.Println()
		fmt.Println("📄 Text Content Preview (first 300 chars):")
		if len(email.TextBody) > 300 {
			fmt.Printf("%s...\n", email.TextBody[:300])
		} else {
			fmt.Printf("%s\n", email.TextBody)
		}
	}

	// Test with high priority inquiry
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("Testing high priority inquiry...")
	
	highPriorityInquiry := &domain.Inquiry{
		ID:       "test-urgent-456",
		Name:     "Jane Smith",
		Email:    "jane.smith@urgentcorp.com",
		Company:  "Urgent Corp",
		Phone:    "+1-555-987-6543",
		Services: []string{"architecture_review"},
		Message:  "URGENT: We need immediate help with our cloud architecture. We have a critical meeting tomorrow and need recommendations ASAP!",
	}

	// Reset mock
	mockSES.sentEmails = nil

	err = emailService.SendCustomerConfirmation(ctx, highPriorityInquiry)
	if err != nil {
		log.Fatalf("Failed to send high priority customer confirmation email: %v", err)
	}

	fmt.Println("✅ High priority customer confirmation email sent successfully!")
	
	if len(mockSES.sentEmails) > 0 {
		email := mockSES.sentEmails[0]
		fmt.Printf("Subject: %s\n", email.Subject)
		fmt.Printf("Recipients: %v\n", email.To)
	}
}

// Mock SES service for testing
type mockSESService struct {
	sentEmails []*interfaces.EmailMessage
}

func (m *mockSESService) SendEmail(ctx context.Context, email *interfaces.EmailMessage) error {
	m.sentEmails = append(m.sentEmails, email)
	fmt.Printf("📤 Mock SES: Email sent to %v\n", email.To)
	return nil
}

func (m *mockSESService) VerifyEmailAddress(ctx context.Context, email string) error {
	return nil
}

func (m *mockSESService) GetSendingQuota(ctx context.Context) (*interfaces.SendingQuota, error) {
	return &interfaces.SendingQuota{
		Max24HourSend:   200,
		MaxSendRate:     14,
		SentLast24Hours: 0,
	}, nil
}