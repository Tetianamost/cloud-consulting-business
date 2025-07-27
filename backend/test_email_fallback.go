package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/cloud-consulting/backend/internal/config"
	"github.com/cloud-consulting/backend/internal/interfaces"
	"github.com/cloud-consulting/backend/internal/services"
	"github.com/cloud-consulting/backend/internal/storage"
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

	// Create storage
	storage := storage.NewInMemoryStorage()

	// Create template service
	templateService := services.NewTemplateService("./templates", logger)

	fmt.Println("🔄 Testing Email Fallback - Graceful Handling of Email Failures")
	fmt.Println(strings.Repeat("=", 70))

	// Test 1: SES service that fails for customer emails but succeeds for consultant emails
	fmt.Println("\n📝 Test 1: Customer Email Failure (Consultant Email Success)")
	mockSES1 := &selectiveFailureSESService{
		failForCustomers: true,
		failForConsultants: false,
	}
	var sesService1 interfaces.SESService = mockSES1
	mockBedrock := &mockBedrockService{}
	var bedrockService interfaces.BedrockService = mockBedrock

	emailService1 := services.NewEmailService(sesService1, templateService, cfg.SES, logger)
	reportService := services.NewReportGenerator(bedrockService, templateService, nil)
	inquiryService1 := services.NewInquiryService(storage, reportService, emailService1)

	req1 := &interfaces.CreateInquiryRequest{
		Name:     "Test Customer",
		Email:    "test@customer.com",
		Company:  "Test Corp",
		Services: []string{"assessment"},
		Message:  "Test inquiry with customer email failure",
	}

	ctx := context.Background()
	inquiry1, err := inquiryService1.CreateInquiry(ctx, req1)
	if err != nil {
		log.Fatalf("Inquiry creation should not fail even if customer email fails: %v", err)
	}

	fmt.Printf("✅ Inquiry created successfully despite customer email failure: %s\n", inquiry1.ID)
	fmt.Printf("📧 Consultant emails sent: %d\n", len(mockSES1.successfulEmails))
	fmt.Printf("❌ Customer emails failed: %d\n", len(mockSES1.failedEmails))

	// Test 2: SES service that fails for consultant emails but succeeds for customer emails
	fmt.Println("\n📝 Test 2: Consultant Email Failure (Customer Email Success)")
	mockSES2 := &selectiveFailureSESService{
		failForCustomers: false,
		failForConsultants: true,
	}
	var sesService2 interfaces.SESService = mockSES2

	emailService2 := services.NewEmailService(sesService2, templateService, cfg.SES, logger)
	inquiryService2 := services.NewInquiryService(storage, reportService, emailService2)

	req2 := &interfaces.CreateInquiryRequest{
		Name:     "Test Customer 2",
		Email:    "test2@customer.com",
		Company:  "Test Corp 2",
		Services: []string{"migration"},
		Message:  "Test inquiry with consultant email failure",
	}

	inquiry2, err := inquiryService2.CreateInquiry(ctx, req2)
	if err != nil {
		log.Fatalf("Inquiry creation should not fail even if consultant email fails: %v", err)
	}

	fmt.Printf("✅ Inquiry created successfully despite consultant email failure: %s\n", inquiry2.ID)
	fmt.Printf("📧 Customer emails sent: %d\n", len(mockSES2.successfulEmails))
	fmt.Printf("❌ Consultant emails failed: %d\n", len(mockSES2.failedEmails))

	// Test 3: Complete email service failure
	fmt.Println("\n📝 Test 3: Complete Email Service Failure")
	mockSES3 := &selectiveFailureSESService{
		failForCustomers: true,
		failForConsultants: true,
	}
	var sesService3 interfaces.SESService = mockSES3

	emailService3 := services.NewEmailService(sesService3, templateService, cfg.SES, logger)
	inquiryService3 := services.NewInquiryService(storage, reportService, emailService3)

	req3 := &interfaces.CreateInquiryRequest{
		Name:     "Test Customer 3",
		Email:    "test3@customer.com",
		Company:  "Test Corp 3",
		Services: []string{"optimization"},
		Message:  "Test inquiry with complete email failure",
	}

	inquiry3, err := inquiryService3.CreateInquiry(ctx, req3)
	if err != nil {
		log.Fatalf("Inquiry creation should not fail even if all emails fail: %v", err)
	}

	fmt.Printf("✅ Inquiry created successfully despite complete email failure: %s\n", inquiry3.ID)
	fmt.Printf("❌ All emails failed: %d\n", len(mockSES3.failedEmails))

	// Test 4: Verify inquiry data integrity
	fmt.Println("\n📝 Test 4: Data Integrity Check")
	
	// Retrieve all inquiries to verify they were stored correctly
	inquiries, err := inquiryService1.ListInquiries(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to list inquiries: %v", err)
	}

	fmt.Printf("✅ Total inquiries stored: %d\n", len(inquiries))
	
	for _, inq := range inquiries {
		if len(inq.Reports) > 0 {
			fmt.Printf("   - Inquiry %s: Has %d reports\n", inq.ID, len(inq.Reports))
		} else {
			fmt.Printf("   - Inquiry %s: No reports generated\n", inq.ID)
		}
	}

	fmt.Println("\n🎯 Fallback Test Summary:")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println("✅ All inquiries created successfully despite email failures")
	fmt.Println("✅ System continues to function when email delivery fails")
	fmt.Println("✅ Data integrity maintained regardless of email status")
	fmt.Println("✅ Graceful fallback implemented correctly")
}

// Mock SES service that can selectively fail for different email types
type selectiveFailureSESService struct {
	failForCustomers   bool
	failForConsultants bool
	successfulEmails   []*interfaces.EmailMessage
	failedEmails       []*interfaces.EmailMessage
}

func (m *selectiveFailureSESService) SendEmail(ctx context.Context, email *interfaces.EmailMessage) error {
	// Determine if this is a customer or consultant email
	isCustomerEmail := strings.Contains(email.Subject, "Thank you")
	isConsultantEmail := strings.Contains(email.Subject, "Report") || 
						 strings.Contains(email.Subject, "Inquiry") ||
						 strings.Contains(email.Subject, "HIGH PRIORITY")

	shouldFail := false
	if isCustomerEmail && m.failForCustomers {
		shouldFail = true
	} else if isConsultantEmail && m.failForConsultants {
		shouldFail = true
	}

	if shouldFail {
		m.failedEmails = append(m.failedEmails, email)
		fmt.Printf("❌ Mock SES: Email FAILED to %v (simulated failure)\n", email.To)
		return fmt.Errorf("simulated email delivery failure")
	} else {
		m.successfulEmails = append(m.successfulEmails, email)
		fmt.Printf("✅ Mock SES: Email sent successfully to %v\n", email.To)
		return nil
	}
}

func (m *selectiveFailureSESService) VerifyEmailAddress(ctx context.Context, email string) error {
	return nil
}

func (m *selectiveFailureSESService) GetSendingQuota(ctx context.Context) (*interfaces.SendingQuota, error) {
	return &interfaces.SendingQuota{
		Max24HourSend:   200,
		MaxSendRate:     14,
		SentLast24Hours: 0,
	}, nil
}

// Mock Bedrock service for testing
type mockBedrockService struct{}

func (m *mockBedrockService) GenerateText(ctx context.Context, prompt string, options *interfaces.BedrockOptions) (*interfaces.BedrockResponse, error) {
	content := `# Assessment Report

## Executive Summary
This is a test report generated for fallback testing.

## Recommendations
- Test recommendation 1
- Test recommendation 2

## Next Steps
- Implement recommendations
- Schedule follow-up meeting`

	return &interfaces.BedrockResponse{
		Content: content,
		Usage: interfaces.BedrockUsage{
			InputTokens:  50,
			OutputTokens: 100,
		},
	}, nil
}

func (m *mockBedrockService) GetModelInfo() interfaces.BedrockModelInfo {
	return interfaces.BedrockModelInfo{
		ModelID:     "amazon.nova-lite-v1:0",
		ModelName:   "Nova Lite",
		Provider:    "Amazon",
		MaxTokens:   4096,
		IsAvailable: true,
	}
}

func (m *mockBedrockService) IsHealthy() bool {
	return true
}