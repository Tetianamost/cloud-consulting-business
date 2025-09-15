package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/cloud-consulting/backend/internal/config"
	"github.com/cloud-consulting/backend/internal/domain"
	"github.com/cloud-consulting/backend/internal/services"
)

func main() {
	fmt.Println("🔍 REAL AWS SES EMAIL DELIVERY TEST")
	fmt.Println("===================================")
	fmt.Println("⚠️  This test will attempt to send REAL emails via AWS SES")
	fmt.Println()

	// Load real configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("❌ Failed to load configuration: %v", err)
	}

	// Check if SES is configured
	if cfg.SES.AccessKeyID == "" || cfg.SES.SecretAccessKey == "" || cfg.SES.SenderEmail == "" {
		fmt.Println("❌ AWS SES NOT CONFIGURED")
		fmt.Println("Missing required environment variables:")
		fmt.Println("   - AWS_ACCESS_KEY_ID")
		fmt.Println("   - AWS_SECRET_ACCESS_KEY")
		fmt.Println("   - SES_SENDER_EMAIL")
		fmt.Println()
		fmt.Println("To test real email delivery, please:")
		fmt.Println("1. Set up AWS SES in your AWS account")
		fmt.Println("2. Verify your sender email address in SES")
		fmt.Println("3. Set the required environment variables")
		fmt.Println("4. Run this test again")
		return
	}

	fmt.Printf("📧 SES Configuration Found:\n")
	fmt.Printf("   Region: %s\n", cfg.SES.Region)
	fmt.Printf("   Sender: %s\n", cfg.SES.SenderEmail)
	fmt.Printf("   Reply-To: %s\n", cfg.SES.ReplyToEmail)
	fmt.Println()

	// Create logger
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	// Create real SES service
	sesService, err := services.NewSESService(cfg.SES, logger)
	if err != nil {
		log.Fatalf("❌ Failed to create SES service: %v", err)
	}

	// Create template service
	templateService := services.NewTemplateService("templates", logger)

	// Create email service with real SES
	emailService := services.NewEmailService(sesService, templateService, cfg.SES, logger)

	ctx := context.Background()

	fmt.Println("🔍 TESTING SES CONNECTION")
	fmt.Println("========================")

	// Test SES connection by getting sending quota
	quota, err := sesService.GetSendingQuota(ctx)
	if err != nil {
		fmt.Printf("❌ SES Connection Failed: %v\n", err)
		fmt.Println()
		fmt.Println("Possible issues:")
		fmt.Println("   - Invalid AWS credentials")
		fmt.Println("   - SES not enabled in the region")
		fmt.Println("   - Network connectivity issues")
		return
	}

	fmt.Printf("✅ SES Connection Successful!\n")
	fmt.Printf("   Max 24h send: %.0f\n", quota.Max24HourSend)
	fmt.Printf("   Max send rate: %.2f/sec\n", quota.MaxSendRate)
	fmt.Printf("   Sent last 24h: %.0f\n", quota.SentLast24Hours)
	fmt.Println()

	// Ask for confirmation before sending real emails
	fmt.Println("⚠️  REAL EMAIL SENDING CONFIRMATION")
	fmt.Println("===================================")
	fmt.Println("This will send REAL emails to:")
	fmt.Printf("   - Customer: %s (test email)\n", cfg.SES.SenderEmail)
	fmt.Printf("   - Internal: info@cloudpartner.pro\n")
	fmt.Println()
	fmt.Print("Do you want to proceed? (y/N): ")

	var response string
	fmt.Scanln(&response)

	if response != "y" && response != "Y" && response != "yes" {
		fmt.Println("❌ Test cancelled by user")
		return
	}

	// Test data - using sender email as customer email for testing
	testInquiry := &domain.Inquiry{
		ID:        fmt.Sprintf("REAL-TEST-%d", time.Now().Unix()),
		Name:      "Test Customer",
		Email:     cfg.SES.SenderEmail, // Use verified sender email for testing
		Company:   "Test Company",
		Phone:     "+1 (555) 123-4567",
		Services:  []string{"AWS Migration", "Cost Optimization"},
		Message:   "This is a test inquiry to verify email delivery is working properly.",
		CreatedAt: time.Now(),
	}

	testReport := &domain.Report{
		ID:        fmt.Sprintf("RPT-REAL-%d", time.Now().Unix()),
		InquiryID: testInquiry.ID,
		Title:     "Test Report - Email Delivery Verification",
		Content: `# Email Delivery Test Report

This is a test report to verify that:
- Email templates render correctly
- AWS SES integration works
- Professional formatting is maintained

## Test Results
- ✅ SES connection established
- ✅ Templates loaded successfully
- ✅ Email service initialized

## Next Steps
If you receive this email, the system is working correctly!`,
		CreatedAt: time.Now(),
	}

	fmt.Println("\n📧 SENDING CUSTOMER CONFIRMATION EMAIL")
	fmt.Println("=====================================")

	err = emailService.SendCustomerConfirmation(ctx, testInquiry)
	if err != nil {
		fmt.Printf("❌ Failed to send customer confirmation: %v\n", err)
	} else {
		fmt.Printf("✅ Customer confirmation sent to: %s\n", testInquiry.Email)
	}

	fmt.Println("\n📧 SENDING CONSULTANT NOTIFICATION EMAIL")
	fmt.Println("========================================")

	err = emailService.SendReportEmail(ctx, testInquiry, testReport)
	if err != nil {
		fmt.Printf("❌ Failed to send consultant notification: %v\n", err)
	} else {
		fmt.Println("✅ Consultant notification sent to: info@cloudpartner.pro")
	}

	fmt.Println("\n✅ REAL EMAIL DELIVERY TEST COMPLETE")
	fmt.Println("====================================")
	fmt.Println("📧 Check your email inbox to verify:")
	fmt.Println("   ✓ Emails were delivered successfully")
	fmt.Println("   ✓ Professional formatting displays correctly")
	fmt.Println("   ✓ All content renders properly")
	fmt.Println("   ✓ Responsive design works on mobile")
	fmt.Println()
	fmt.Println("📱 Test the emails in different clients:")
	fmt.Println("   - Gmail (web and mobile)")
	fmt.Println("   - Outlook (web and desktop)")
	fmt.Println("   - Apple Mail")
	fmt.Println("   - Other email clients you use")
}
