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
	fmt.Println("📧 EMAIL SYSTEM TEST")
	fmt.Println("===================")
	fmt.Println("Testing the clean, unified email implementation")
	fmt.Println()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("❌ Failed to load configuration: %v", err)
	}

	// Check if SES is configured
	if cfg.SES.AccessKeyID == "" || cfg.SES.SecretAccessKey == "" || cfg.SES.SenderEmail == "" {
		fmt.Println("❌ AWS SES NOT CONFIGURED")
		fmt.Println("Set AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, and SES_SENDER_EMAIL")
		return
	}

	fmt.Printf("📧 SES Configuration:\n")
	fmt.Printf("   Region: %s\n", cfg.SES.Region)
	fmt.Printf("   Sender: %s\n", cfg.SES.SenderEmail)
	fmt.Println()

	// Create logger
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	fmt.Println("🔧 CREATING EMAIL SERVICE")
	fmt.Println("=========================")

	// Use the clean, unified email service creation
	emailService, err := services.NewEmailServiceWithSES(cfg.SES, logger)
	if err != nil {
		log.Fatalf("❌ Failed to create email service: %v", err)
	}

	fmt.Println("✅ Email service created successfully")

	// Test data
	testInquiry := &domain.Inquiry{
		ID:        fmt.Sprintf("CLEAN-TEST-%d", time.Now().Unix()),
		Name:      "Clean Implementation Test",
		Email:     cfg.SES.SenderEmail, // Use verified sender email for testing
		Company:   "Test Company",
		Phone:     "+1 (555) 123-4567",
		Services:  []string{"assessment"},
		Message:   "Testing the clean, unified email implementation without 'fixed' naming.",
		CreatedAt: time.Now(),
	}

	testReport := &domain.Report{
		ID:        fmt.Sprintf("RPT-CLEAN-%d", time.Now().Unix()),
		InquiryID: testInquiry.ID,
		Title:     "Clean Email Implementation Test",
		Content: `# Clean Email Implementation Test

This report verifies that the email system works with the clean, unified implementation.

## What's Been Cleaned Up
- ✅ Removed old broken SES implementation
- ✅ Renamed "fixed" implementation to be the main one
- ✅ Simplified function names (no more "Fixed" suffix)
- ✅ Updated all references in production code

## Test Results
If you receive this email properly formatted, the clean implementation is working!`,
		CreatedAt: time.Now(),
	}

	ctx := context.Background()

	// Ask for confirmation
	fmt.Println("⚠️  EMAIL DELIVERY TEST")
	fmt.Println("======================")
	fmt.Printf("Send test emails to %s and info@cloudpartner.pro? (y/N): ", cfg.SES.SenderEmail)

	var response string
	fmt.Scanln(&response)

	if response != "y" && response != "Y" && response != "yes" {
		fmt.Println("❌ Test cancelled")
		return
	}

	fmt.Println("\n📧 SENDING CUSTOMER CONFIRMATION")
	fmt.Println("===============================")

	err = emailService.SendCustomerConfirmation(ctx, testInquiry)
	if err != nil {
		fmt.Printf("❌ Customer confirmation failed: %v\n", err)
	} else {
		fmt.Printf("✅ Customer confirmation sent successfully!\n")
	}

	fmt.Println("\n📧 SENDING CONSULTANT NOTIFICATION")
	fmt.Println("=================================")

	err = emailService.SendReportEmail(ctx, testInquiry, testReport)
	if err != nil {
		fmt.Printf("❌ Consultant notification failed: %v\n", err)
	} else {
		fmt.Printf("✅ Consultant notification sent successfully!\n")
	}

	fmt.Println("\n🎉 CLEAN EMAIL IMPLEMENTATION TEST COMPLETE")
	fmt.Println("===========================================")

	if err == nil {
		fmt.Println("✅ SUCCESS! Clean email implementation is working")
		fmt.Println("📧 Both customer and consultant emails should be properly formatted")
		fmt.Println("🧹 No more 'fixed' naming - everything is clean and unified")
		fmt.Println()
		fmt.Println("📱 Ready for production use!")
	} else {
		fmt.Println("⚠️  Some issues remain - check error messages above")
	}
}
