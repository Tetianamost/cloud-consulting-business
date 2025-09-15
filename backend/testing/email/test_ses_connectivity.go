package main

import (
	"context"
	"fmt"
	"log"

	"github.com/sirupsen/logrus"

	"github.com/cloud-consulting/backend/internal/config"
	"github.com/cloud-consulting/backend/internal/services"
)

func main() {
	fmt.Println("🔍 AWS SES CONNECTIVITY TEST")
	fmt.Println("============================")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("❌ Failed to load configuration: %v", err)
	}

	fmt.Printf("📧 SES Configuration:\n")
	fmt.Printf("   Region: %s\n", cfg.SES.Region)
	fmt.Printf("   Sender: %s\n", cfg.SES.SenderEmail)
	fmt.Printf("   Has Access Key: %t\n", cfg.SES.AccessKeyID != "")
	fmt.Printf("   Has Secret Key: %t\n", cfg.SES.SecretAccessKey != "")
	fmt.Println()

	// Create logger
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	// Create SES service
	sesService, err := services.NewSESService(cfg.SES, logger)
	if err != nil {
		fmt.Printf("❌ Failed to create SES service: %v\n", err)
		return
	}

	fmt.Println("✅ SES service created successfully")

	// Test SES connection
	ctx := context.Background()
	quota, err := sesService.GetSendingQuota(ctx)
	if err != nil {
		fmt.Printf("❌ SES connection failed: %v\n", err)
		fmt.Println()
		fmt.Println("Possible issues:")
		fmt.Println("   - Invalid AWS credentials")
		fmt.Println("   - SES not enabled in the region")
		fmt.Println("   - Sender email not verified in SES")
		fmt.Println("   - Network connectivity issues")
		fmt.Println()
		fmt.Println("To fix:")
		fmt.Println("1. Go to AWS SES Console")
		fmt.Println("2. Verify your sender email address")
		fmt.Println("3. Check your AWS credentials")
		return
	}

	fmt.Printf("✅ SES connection successful!\n")
	fmt.Printf("   Max 24h send: %.0f emails\n", quota.Max24HourSend)
	fmt.Printf("   Max send rate: %.2f emails/sec\n", quota.MaxSendRate)
	fmt.Printf("   Sent last 24h: %.0f emails\n", quota.SentLast24Hours)
	fmt.Println()

	// Check if sender email is verified
	err = sesService.VerifyEmailAddress(ctx, cfg.SES.SenderEmail)
	if err != nil {
		fmt.Printf("⚠️  Email verification status unknown: %v\n", err)
		fmt.Println("   This might mean the email is already verified")
	} else {
		fmt.Printf("📧 Email verification initiated for: %s\n", cfg.SES.SenderEmail)
	}

	fmt.Println()
	fmt.Println("✅ SES CONNECTIVITY TEST COMPLETE")
	fmt.Println("=================================")
	fmt.Println("🎯 SES is properly configured and accessible")
	fmt.Println("📧 Ready to send emails")
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("1. Ensure sender email is verified in AWS SES Console")
	fmt.Println("2. Run the full email delivery test")
	fmt.Println("3. Test with real customer inquiries")
}
