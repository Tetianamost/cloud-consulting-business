package main

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/cloud-consulting/backend/internal/config"
	"github.com/cloud-consulting/backend/internal/domain"
	"github.com/cloud-consulting/backend/internal/services"
)

func main() {
	fmt.Println("🔧 EMAIL FIX VERIFICATION TEST")
	fmt.Println("==============================")
	fmt.Println("Testing the MIME boundary fix for email delivery")
	fmt.Println()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("❌ Failed to load configuration: %v\n", err)
		return
	}

	// Create logger
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	// Create SES service
	sesService, err := services.NewSESService(cfg.SES, logger)
	if err != nil {
		fmt.Printf("❌ Failed to create SES service: %v\n", err)
		return
	}

	// Create template service
	templateService := services.NewTemplateService("templates", logger)

	// Create email service
	emailService := services.NewEmailService(sesService, templateService, cfg.SES, logger)

	// Test data
	testInquiry := &domain.Inquiry{
		ID:        fmt.Sprintf("EMAIL-FIX-TEST-%d", time.Now().Unix()),
		Name:      "Email Fix Test",
		Email:     cfg.SES.SenderEmail, // Use verified sender email for testing
		Company:   "Test Company",
		Phone:     "+1 (555) 123-4567",
		Services:  []string{"Email Testing"},
		Message:   "This is a test to verify the email MIME boundary fix is working correctly.",
		CreatedAt: time.Now(),
	}

	ctx := context.Background()

	fmt.Println("📧 TESTING CUSTOMER CONFIRMATION EMAIL (Simple API)")
	fmt.Println("==================================================")

	err = emailService.SendCustomerConfirmation(ctx, testInquiry)
	if err != nil {
		fmt.Printf("❌ Customer confirmation failed: %v\n", err)
	} else {
		fmt.Printf("✅ Customer confirmation sent successfully to: %s\n", testInquiry.Email)
	}

	fmt.Println("\n📧 TESTING INQUIRY NOTIFICATION EMAIL (Simple API)")
	fmt.Println("=================================================")

	err = emailService.SendInquiryNotification(ctx, testInquiry)
	if err != nil {
		fmt.Printf("❌ Inquiry notification failed: %v\n", err)
	} else {
		fmt.Println("✅ Inquiry notification sent successfully to: info@cloudpartner.pro")
	}

	fmt.Println("\n🎯 EMAIL FIX VERIFICATION COMPLETE")
	fmt.Println("==================================")
	fmt.Println("If both emails sent successfully, the MIME boundary issue is fixed!")
	fmt.Println("Check your email inbox to verify delivery.")
}
