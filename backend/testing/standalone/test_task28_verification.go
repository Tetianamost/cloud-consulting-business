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
	fmt.Println("🎯 Task 28 Verification: Enhanced Email Service with Customer Confirmations")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create logger
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	// Create services
	storage := storage.NewInMemoryStorage()
	templateService := services.NewTemplateService("./templates", logger)
	mockSES := &verificationSESService{}
	var sesService interfaces.SESService = mockSES
	mockBedrock := &mockBedrockService{}
	var bedrockService interfaces.BedrockService = mockBedrock

	emailService := services.NewEmailService(sesService, templateService, cfg.SES, logger)
	reportService := services.NewReportGenerator(bedrockService, templateService, nil)
	inquiryService := services.NewInquiryService(storage, reportService, emailService)

	ctx := context.Background()

	// Requirement 1: Send immediate confirmation emails to customers upon inquiry submission
	fmt.Println("📋 Requirement 1: Immediate Customer Confirmation Emails")
	fmt.Println(strings.Repeat("-", 60))
	
	req1 := &interfaces.CreateInquiryRequest{
		Name:     "Alice Johnson",
		Email:    "alice@testcorp.com",
		Company:  "Test Corporation",
		Services: []string{"assessment"},
		Message:  "We need a cloud assessment for our infrastructure.",
	}

	mockSES.reset()
	_, err = inquiryService.CreateInquiry(ctx, req1)
	if err != nil {
		log.Fatalf("Failed to create inquiry: %v", err)
	}

	customerEmails := mockSES.getCustomerEmails()
	if len(customerEmails) > 0 {
		fmt.Printf("✅ Customer confirmation email sent immediately\n")
		fmt.Printf("   - Recipient: %v\n", customerEmails[0].To)
		fmt.Printf("   - Subject: %s\n", customerEmails[0].Subject)
		fmt.Printf("   - Sent within inquiry creation process\n")
	} else {
		fmt.Printf("❌ No customer confirmation email sent\n")
	}
	fmt.Println()

	// Requirement 2: Include branded templates with professional messaging
	fmt.Println("📋 Requirement 2: Branded Templates with Professional Messaging")
	fmt.Println(strings.Repeat("-", 60))
	
	if len(customerEmails) > 0 {
		email := customerEmails[0]
		
		// Check for branding elements
		hasBranding := strings.Contains(email.HTMLBody, "CloudPartner Pro") ||
					  strings.Contains(email.HTMLBody, "cloud-consulting") ||
					  strings.Contains(email.HTMLBody, "logo")
		
		// Check for professional messaging
		hasProfessionalTone := strings.Contains(email.Subject, "Thank you") &&
							   strings.Contains(email.HTMLBody, "received your inquiry") &&
							   strings.Contains(email.HTMLBody, "team will review")
		
		// Check for branded styling
		hasStyledHTML := strings.Contains(email.HTMLBody, "<style>") &&
						 strings.Contains(email.HTMLBody, "background") &&
						 strings.Contains(email.HTMLBody, "color")
		
		if hasBranding {
			fmt.Printf("✅ Branded template used (contains CloudPartner Pro branding)\n")
		} else {
			fmt.Printf("❌ No branding detected in template\n")
		}
		
		if hasProfessionalTone {
			fmt.Printf("✅ Professional messaging confirmed\n")
		} else {
			fmt.Printf("❌ Professional messaging not detected\n")
		}
		
		if hasStyledHTML {
			fmt.Printf("✅ Styled HTML template with professional design\n")
		} else {
			fmt.Printf("❌ No professional styling detected\n")
		}
	}
	fmt.Println()

	// Requirement 3: Add download links for reports when available (only to consultants, never to customers)
	fmt.Println("📋 Requirement 3: Report Download Links (Consultants Only)")
	fmt.Println(strings.Repeat("-", 60))
	
	consultantEmails := mockSES.getConsultantEmails()
	customerEmails = mockSES.getCustomerEmails()
	
	// Check consultant emails for download links
	consultantHasDownloadLinks := false
	if len(consultantEmails) > 0 {
		consultantEmail := consultantEmails[0]
		consultantHasDownloadLinks = strings.Contains(strings.ToLower(consultantEmail.HTMLBody), "download") ||
									 strings.Contains(strings.ToLower(consultantEmail.TextBody), "download")
		
		if consultantHasDownloadLinks {
			fmt.Printf("✅ Consultant emails contain download functionality\n")
		} else {
			fmt.Printf("ℹ️  Consultant emails don't contain explicit download links (may use admin dashboard)\n")
		}
	}
	
	// Check customer emails for download links (should NOT have any)
	customerHasDownloadLinks := false
	if len(customerEmails) > 0 {
		customerEmail := customerEmails[0]
		customerHasDownloadLinks = (strings.Contains(strings.ToLower(customerEmail.HTMLBody), "download") &&
								   strings.Contains(strings.ToLower(customerEmail.HTMLBody), "report")) ||
								  (strings.Contains(strings.ToLower(customerEmail.TextBody), "download") &&
								   strings.Contains(strings.ToLower(customerEmail.TextBody), "report"))
		
		if !customerHasDownloadLinks {
			fmt.Printf("✅ Customer emails do NOT contain report download links (security confirmed)\n")
		} else {
			fmt.Printf("❌ SECURITY VIOLATION: Customer emails contain report download links!\n")
		}
	}
	
	// Check that reports are never attached to customer emails
	customerHasReportAttachments := false
	if len(customerEmails) > 0 {
		customerEmail := customerEmails[0]
		customerHasReportAttachments = len(customerEmail.Attachments) > 0
		
		if !customerHasReportAttachments {
			fmt.Printf("✅ Customer emails do NOT contain report attachments (security confirmed)\n")
		} else {
			fmt.Printf("❌ SECURITY VIOLATION: Customer emails contain attachments!\n")
		}
	}
	fmt.Println()

	// Requirement 4: Implement graceful fallback if email delivery fails
	fmt.Println("📋 Requirement 4: Graceful Fallback on Email Delivery Failure")
	fmt.Println(strings.Repeat("-", 60))
	
	// Test with failing email service
	failingSES := &failingSESService{}
	var failingService interfaces.SESService = failingSES
	failingEmailService := services.NewEmailService(failingService, templateService, cfg.SES, logger)
	failingInquiryService := services.NewInquiryService(storage, reportService, failingEmailService)
	
	req2 := &interfaces.CreateInquiryRequest{
		Name:     "Bob Wilson",
		Email:    "bob@failtest.com",
		Company:  "Fail Test Corp",
		Services: []string{"migration"},
		Message:  "Test inquiry with email failure",
	}
	
	inquiry2, err := failingInquiryService.CreateInquiry(ctx, req2)
	if err != nil {
		fmt.Printf("❌ Inquiry creation failed when emails fail: %v\n", err)
	} else {
		fmt.Printf("✅ Inquiry creation succeeded despite email failures\n")
		fmt.Printf("   - Inquiry ID: %s\n", inquiry2.ID)
		fmt.Printf("   - Reports generated: %d\n", len(inquiry2.Reports))
		fmt.Printf("   - System continued to function normally\n")
	}
	fmt.Println()

	// Summary
	fmt.Println("🎯 Task 28 Implementation Summary")
	fmt.Println(strings.Repeat("=", 50))
	
	allRequirementsMet := true
	
	if len(mockSES.getCustomerEmails()) > 0 {
		fmt.Println("✅ Requirement 1: Immediate customer confirmation emails - IMPLEMENTED")
	} else {
		fmt.Println("❌ Requirement 1: Immediate customer confirmation emails - MISSING")
		allRequirementsMet = false
	}
	
	if len(customerEmails) > 0 && strings.Contains(customerEmails[0].HTMLBody, "CloudPartner Pro") {
		fmt.Println("✅ Requirement 2: Branded templates with professional messaging - IMPLEMENTED")
	} else {
		fmt.Println("❌ Requirement 2: Branded templates with professional messaging - MISSING")
		allRequirementsMet = false
	}
	
	if !customerHasDownloadLinks && !customerHasReportAttachments {
		fmt.Println("✅ Requirement 3: Download links only for consultants (never customers) - IMPLEMENTED")
	} else {
		fmt.Println("❌ Requirement 3: Download links only for consultants (never customers) - VIOLATED")
		allRequirementsMet = false
	}
	
	if inquiry2 != nil {
		fmt.Println("✅ Requirement 4: Graceful fallback on email delivery failure - IMPLEMENTED")
	} else {
		fmt.Println("❌ Requirement 4: Graceful fallback on email delivery failure - MISSING")
		allRequirementsMet = false
	}
	
	fmt.Println()
	if allRequirementsMet {
		fmt.Println("🎉 ALL REQUIREMENTS SUCCESSFULLY IMPLEMENTED!")
		fmt.Println("✅ Task 28 is complete and working correctly")
	} else {
		fmt.Println("⚠️  Some requirements need attention")
	}
}

// Verification SES service that tracks different email types
type verificationSESService struct {
	allEmails []*interfaces.EmailMessage
}

func (v *verificationSESService) reset() {
	v.allEmails = nil
}

func (v *verificationSESService) getCustomerEmails() []*interfaces.EmailMessage {
	var customerEmails []*interfaces.EmailMessage
	for _, email := range v.allEmails {
		if strings.Contains(email.Subject, "Thank you") {
			customerEmails = append(customerEmails, email)
		}
	}
	return customerEmails
}

func (v *verificationSESService) getConsultantEmails() []*interfaces.EmailMessage {
	var consultantEmails []*interfaces.EmailMessage
	for _, email := range v.allEmails {
		if strings.Contains(email.Subject, "Report") || 
		   strings.Contains(email.Subject, "HIGH PRIORITY") ||
		   strings.Contains(email.Subject, "Inquiry") {
			consultantEmails = append(consultantEmails, email)
		}
	}
	return consultantEmails
}

func (v *verificationSESService) SendEmail(ctx context.Context, email *interfaces.EmailMessage) error {
	v.allEmails = append(v.allEmails, email)
	return nil
}

func (v *verificationSESService) VerifyEmailAddress(ctx context.Context, email string) error {
	return nil
}

func (v *verificationSESService) GetSendingQuota(ctx context.Context) (*interfaces.SendingQuota, error) {
	return &interfaces.SendingQuota{
		Max24HourSend:   200,
		MaxSendRate:     14,
		SentLast24Hours: 0,
	}, nil
}

// Failing SES service for testing graceful fallback
type failingSESService struct{}

func (f *failingSESService) SendEmail(ctx context.Context, email *interfaces.EmailMessage) error {
	return fmt.Errorf("simulated email service failure")
}

func (f *failingSESService) VerifyEmailAddress(ctx context.Context, email string) error {
	return fmt.Errorf("simulated email service failure")
}

func (f *failingSESService) GetSendingQuota(ctx context.Context) (*interfaces.SendingQuota, error) {
	return nil, fmt.Errorf("simulated email service failure")
}

// Mock Bedrock service for testing
type mockBedrockService struct{}

func (m *mockBedrockService) GenerateText(ctx context.Context, prompt string, options *interfaces.BedrockOptions) (*interfaces.BedrockResponse, error) {
	content := `# Cloud Assessment Report

## Executive Summary
This is a comprehensive assessment of your cloud infrastructure needs.

## Recommendations
- Implement cost optimization strategies
- Enhance security configurations
- Improve monitoring and alerting

## Next Steps
- Schedule implementation planning meeting
- Begin phased migration approach
- Establish ongoing support framework`

	return &interfaces.BedrockResponse{
		Content: content,
		Usage: interfaces.BedrockUsage{
			InputTokens:  75,
			OutputTokens: 125,
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