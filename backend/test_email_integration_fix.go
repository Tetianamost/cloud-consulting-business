package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/cloud-consulting/backend/internal/config"
	"github.com/cloud-consulting/backend/internal/domain"
	"github.com/cloud-consulting/backend/internal/services"
	"github.com/sirupsen/logrus"
)

func main() {
	fmt.Println("🔧 Testing Email Integration Fix...")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create logger
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	// Check SES configuration
	fmt.Printf("📧 SES Configuration:\n")
	fmt.Printf("   Access Key ID: %s\n", maskString(cfg.SES.AccessKeyID))
	fmt.Printf("   Secret Key: %s\n", maskString(cfg.SES.SecretAccessKey))
	fmt.Printf("   Region: %s\n", cfg.SES.Region)
	fmt.Printf("   Sender Email: %s\n", cfg.SES.SenderEmail)
	fmt.Printf("   Reply To Email: %s\n", cfg.SES.ReplyToEmail)

	// Test email service initialization
	fmt.Println("\n🚀 Testing Email Service Initialization...")
	emailService, err := services.NewEmailServiceWithSES(cfg.SES, logger)
	if err != nil {
		log.Fatalf("❌ Failed to initialize email service: %v", err)
	}
	fmt.Println("✅ Email service initialized successfully!")

	// Test email service health
	fmt.Println("\n🏥 Testing Email Service Health...")
	if emailService.IsHealthy() {
		fmt.Println("✅ Email service is healthy!")
	} else {
		fmt.Println("❌ Email service is not healthy!")
		return
	}

	// Create a test inquiry
	fmt.Println("\n📝 Creating Test Inquiry...")
	inquiry := &domain.Inquiry{
		ID:        "test-email-fix-" + fmt.Sprintf("%d", time.Now().Unix()),
		Name:      "Test Customer",
		Email:     "test@example.com", // This won't actually send since it's not verified
		Company:   "Test Company",
		Phone:     "555-0123",
		Services:  []string{"optimization"},
		Message:   "Test inquiry for email integration fix",
		Status:    domain.InquiryStatusNew,
		Priority:  domain.InquiryPriorityMedium,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Create a test report
	report := &domain.Report{
		ID:        "test-report-" + fmt.Sprintf("%d", time.Now().Unix()),
		Type:      domain.ReportTypeConsulting,
		Status:    domain.ReportStatusCompleted,
		Content:   "This is a test report for email integration testing.",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	inquiry.Reports = []*domain.Report{report}

	ctx := context.Background()

	// Test customer confirmation email
	fmt.Println("\n📧 Testing Customer Confirmation Email...")
	err = emailService.SendCustomerConfirmation(ctx, inquiry)
	if err != nil {
		fmt.Printf("❌ Failed to send customer confirmation: %v\n", err)
	} else {
		fmt.Println("✅ Customer confirmation email sent successfully!")
	}

	// Test internal report email
	fmt.Println("\n📧 Testing Internal Report Email...")
	err = emailService.SendReportEmail(ctx, inquiry, report)
	if err != nil {
		fmt.Printf("❌ Failed to send internal report email: %v\n", err)
	} else {
		fmt.Println("✅ Internal report email sent successfully!")
	}

	fmt.Println("\n🎉 Email integration test completed!")
	fmt.Println("\n📝 Note: Emails are sent to AWS SES. Check your SES console for delivery status.")
	fmt.Println("   Customer emails go to: test@example.com (may bounce if not verified)")
	fmt.Println("   Internal emails go to: info@cloudpartner.pro")
}

func maskString(s string) string {
	if len(s) <= 4 {
		return "****"
	}
	return s[:4] + "****" + s[len(s)-4:]
}
