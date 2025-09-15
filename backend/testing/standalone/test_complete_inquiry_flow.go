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

	// Create mock services
	mockSES := &mockSESService{}
	var sesService interfaces.SESService = mockSES
	mockBedrock := &mockBedrockService{}
	var bedrockService interfaces.BedrockService = mockBedrock

	// Create email service
	emailService := services.NewEmailService(sesService, templateService, cfg.SES, logger)

	// Create report service (using the correct constructor)
	reportService := services.NewReportGenerator(bedrockService, templateService, nil)

	// Create inquiry service
	inquiryService := services.NewInquiryService(storage, reportService, emailService)

	// Create test inquiry request
	req := &interfaces.CreateInquiryRequest{
		Name:     "Alice Johnson",
		Email:    "alice.johnson@testcorp.com",
		Company:  "Test Corporation",
		Phone:    "+1-555-999-8888",
		Services: []string{"assessment", "optimization"},
		Message:  "We need help optimizing our current AWS infrastructure. We're looking to reduce costs and improve performance. This is urgent as we have budget reviews next week.",
	}

	fmt.Println("Testing complete inquiry creation flow...")
	fmt.Printf("Customer: %s (%s)\n", req.Name, req.Email)
	fmt.Printf("Company: %s\n", req.Company)
	fmt.Printf("Services: %v\n", req.Services)
	fmt.Printf("Message: %s\n", req.Message)
	fmt.Println()

	// Create inquiry
	ctx := context.Background()
	inquiry, err := inquiryService.CreateInquiry(ctx, req)
	if err != nil {
		log.Fatalf("Failed to create inquiry: %v", err)
	}

	fmt.Printf("✅ Inquiry created successfully! ID: %s\n", inquiry.ID)
	fmt.Println()

	// Check what emails were sent
	fmt.Println("📧 Emails sent during inquiry creation:")
	for i, email := range mockSES.sentEmails {
		fmt.Printf("%d. To: %v\n", i+1, email.To)
		fmt.Printf("   Subject: %s\n", email.Subject)
		fmt.Printf("   Type: %s\n", getEmailType(email))
		fmt.Println()
	}

	// Check if reports were generated
	if len(inquiry.Reports) > 0 {
		fmt.Printf("📋 Report generated: %s\n", inquiry.Reports[0].ID)
		fmt.Printf("Report content preview (first 200 chars):\n%s...\n", 
			truncateString(inquiry.Reports[0].Content, 200))
	} else {
		fmt.Println("⚠️  No reports were generated")
	}

	fmt.Println()
	fmt.Println("✅ Complete inquiry flow test completed!")
}

func getEmailType(email *interfaces.EmailMessage) string {
	if strings.Contains(email.Subject, "Thank you") {
		return "Customer Confirmation"
	}
	if strings.Contains(email.Subject, "Report Generated") || strings.Contains(email.Subject, "Report") {
		return "Consultant Report Notification"
	}
	if strings.Contains(email.Subject, "Inquiry") {
		return "Consultant Inquiry Notification"
	}
	return "Unknown"
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
	fmt.Printf("📤 Mock SES: Email sent to %v with subject: %s\n", email.To, email.Subject)
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
	// Generate a mock report based on the prompt
	content := `# Cloud Infrastructure Assessment Report

## Executive Summary
Based on your inquiry regarding AWS infrastructure optimization, we have identified several key areas for improvement that can help reduce costs and enhance performance.

## Current State Assessment
Your current infrastructure appears to be running on standard configurations that may not be optimized for your specific workload patterns.

## Recommendations
1. **Cost Optimization**: Implement Reserved Instances for predictable workloads
2. **Performance Enhancement**: Consider upgrading to newer instance types
3. **Monitoring**: Implement comprehensive CloudWatch monitoring
4. **Security**: Review and update security groups and IAM policies

## Next Steps
1. Schedule a detailed technical review meeting
2. Conduct a comprehensive infrastructure audit
3. Develop a phased optimization plan
4. Implement changes with minimal downtime

## Priority Level
**HIGH PRIORITY** - Budget reviews are time-sensitive and require immediate attention.

This assessment provides a foundation for your infrastructure optimization project. Our team is ready to assist with implementation.`

	return &interfaces.BedrockResponse{
		Content: content,
		Usage: interfaces.BedrockUsage{
			InputTokens:  100,
			OutputTokens: 250,
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