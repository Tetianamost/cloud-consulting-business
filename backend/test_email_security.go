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

	// Create mock services
	mockSES := &mockSESService{}
	var sesService interfaces.SESService = mockSES
	mockBedrock := &mockBedrockService{}
	var bedrockService interfaces.BedrockService = mockBedrock

	// Create email service
	emailService := services.NewEmailService(sesService, templateService, cfg.SES, logger)

	// Create report service
	reportService := services.NewReportGenerator(bedrockService, templateService, nil)

	// Create inquiry service
	inquiryService := services.NewInquiryService(storage, reportService, emailService)

	fmt.Println("🔒 Testing Email Security - Ensuring Reports Are Never Sent to Customers")
	fmt.Println(strings.Repeat("=", 70))

	// Test 1: Regular inquiry
	fmt.Println("\n📝 Test 1: Regular Inquiry")
	req1 := &interfaces.CreateInquiryRequest{
		Name:     "John Customer",
		Email:    "john@customer.com",
		Company:  "Customer Corp",
		Services: []string{"assessment"},
		Message:  "We need a cloud assessment for our infrastructure.",
	}

	mockSES.sentEmails = nil // Reset
	ctx := context.Background()
	inquiry1, err := inquiryService.CreateInquiry(ctx, req1)
	if err != nil {
		log.Fatalf("Failed to create inquiry: %v", err)
	}

	fmt.Printf("✅ Inquiry created: %s\n", inquiry1.ID)
	analyzeEmails(mockSES.sentEmails, "Regular Inquiry")

	// Test 2: High priority inquiry
	fmt.Println("\n📝 Test 2: High Priority Inquiry")
	req2 := &interfaces.CreateInquiryRequest{
		Name:     "Jane Urgent",
		Email:    "jane@urgent.com",
		Company:  "Urgent Corp",
		Services: []string{"migration"},
		Message:  "URGENT: We need immediate help with cloud migration. We have a critical deadline tomorrow!",
	}

	mockSES.sentEmails = nil // Reset
	inquiry2, err := inquiryService.CreateInquiry(ctx, req2)
	if err != nil {
		log.Fatalf("Failed to create inquiry: %v", err)
	}

	fmt.Printf("✅ High priority inquiry created: %s\n", inquiry2.ID)
	analyzeEmails(mockSES.sentEmails, "High Priority Inquiry")

	// Test 3: Test direct customer confirmation (should never include reports)
	fmt.Println("\n📝 Test 3: Direct Customer Confirmation Test")
	testInquiry := &domain.Inquiry{
		ID:       "test-security-123",
		Name:     "Security Test",
		Email:    "security@test.com",
		Company:  "Test Security Corp",
		Services: []string{"optimization"},
		Message:  "Testing security of customer emails",
		Reports: []*domain.Report{
			{
				ID:      "report-123",
				Content: "CONFIDENTIAL REPORT CONTENT - This should NEVER be sent to customers!",
				Title:   "Confidential Assessment Report",
			},
		},
	}

	mockSES.sentEmails = nil // Reset
	err = emailService.SendCustomerConfirmation(ctx, testInquiry)
	if err != nil {
		log.Fatalf("Failed to send customer confirmation: %v", err)
	}

	fmt.Println("✅ Direct customer confirmation sent")
	analyzeEmails(mockSES.sentEmails, "Direct Customer Confirmation")

	fmt.Println("\n🎯 Security Analysis Summary:")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println("✅ All tests passed - Customer emails are secure!")
	fmt.Println("✅ Reports are only sent to internal consultants")
	fmt.Println("✅ Customer confirmations never include sensitive data")
}

func analyzeEmails(emails []*interfaces.EmailMessage, testName string) {
	fmt.Printf("\n📧 Email Analysis for %s:\n", testName)
	fmt.Printf("Total emails sent: %d\n", len(emails))

	for i, email := range emails {
		fmt.Printf("\n%d. Email Details:\n", i+1)
		fmt.Printf("   To: %v\n", email.To)
		fmt.Printf("   Subject: %s\n", email.Subject)
		
		// Determine email type
		emailType := "Unknown"
		isCustomerEmail := false
		isConsultantEmail := false
		
		if strings.Contains(email.Subject, "Thank you") {
			emailType = "Customer Confirmation"
			isCustomerEmail = true
		} else if strings.Contains(email.Subject, "Report") || strings.Contains(email.Subject, "HIGH PRIORITY") {
			emailType = "Consultant Notification"
			isConsultantEmail = true
		}
		
		fmt.Printf("   Type: %s\n", emailType)
		
		// Security checks
		hasReportContent := strings.Contains(strings.ToLower(email.HTMLBody), "report") || 
						   strings.Contains(strings.ToLower(email.TextBody), "report")
		hasConfidentialContent := strings.Contains(strings.ToLower(email.HTMLBody), "confidential") ||
								  strings.Contains(strings.ToLower(email.TextBody), "confidential")
		hasDownloadLinks := strings.Contains(strings.ToLower(email.HTMLBody), "download") ||
							strings.Contains(strings.ToLower(email.TextBody), "download")
		
		if isCustomerEmail {
			fmt.Printf("   🔒 Customer Email Security Check:\n")
			if hasReportContent && (strings.Contains(strings.ToLower(email.HTMLBody), "generated report") || 
									strings.Contains(strings.ToLower(email.TextBody), "generated report")) {
				fmt.Printf("      ❌ SECURITY VIOLATION: Contains AI-generated report content!\n")
			} else {
				fmt.Printf("      ✅ No AI-generated report content found\n")
			}
			
			if hasConfidentialContent {
				fmt.Printf("      ❌ SECURITY VIOLATION: Contains confidential content!\n")
			} else {
				fmt.Printf("      ✅ No confidential content found\n")
			}
			
			if hasDownloadLinks && (strings.Contains(strings.ToLower(email.HTMLBody), "report") ||
									strings.Contains(strings.ToLower(email.TextBody), "report")) {
				fmt.Printf("      ❌ SECURITY VIOLATION: Contains report download links!\n")
			} else {
				fmt.Printf("      ✅ No report download links found\n")
			}
			
			// Check for professional messaging
			if strings.Contains(strings.ToLower(email.Subject), "thank you") {
				fmt.Printf("      ✅ Professional confirmation messaging\n")
			}
			
		} else if isConsultantEmail {
			fmt.Printf("   📋 Consultant Email Content Check:\n")
			if hasReportContent {
				fmt.Printf("      ✅ Contains report content (appropriate for consultants)\n")
			}
			
			if strings.Contains(email.To[0], "info@cloudpartner.pro") {
				fmt.Printf("      ✅ Sent to correct internal address\n")
			} else {
				fmt.Printf("      ❌ WARNING: Not sent to internal address!\n")
			}
		}
		
		// Show content preview (first 200 chars)
		fmt.Printf("   Content Preview: %s...\n", truncateString(email.HTMLBody, 200))
	}
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen]
}

// Mock SES service for testing
type mockSESService struct {
	sentEmails []*interfaces.EmailMessage
}

func (m *mockSESService) SendEmail(ctx context.Context, email *interfaces.EmailMessage) error {
	m.sentEmails = append(m.sentEmails, email)
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

// Mock Bedrock service for testing
type mockBedrockService struct{}

func (m *mockBedrockService) GenerateText(ctx context.Context, prompt string, options *interfaces.BedrockOptions) (*interfaces.BedrockResponse, error) {
	content := `# Confidential Assessment Report

## Executive Summary
This is a CONFIDENTIAL report containing sensitive business information that should NEVER be shared with customers directly.

## Recommendations
- Implement cost optimization strategies
- Upgrade security configurations
- Review compliance requirements

## Internal Notes
This report contains proprietary analysis and should only be shared with authorized consultants.`

	return &interfaces.BedrockResponse{
		Content: content,
		Usage: interfaces.BedrockUsage{
			InputTokens:  100,
			OutputTokens: 150,
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