package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/cloud-consulting/backend/internal/config"
	"github.com/cloud-consulting/backend/internal/interfaces"
	"github.com/cloud-consulting/backend/internal/services"
)

func main() {
	fmt.Println("=== Testing SES Delivery Status and Notification Processing ===")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create logger
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	// Create SES service
	sesService, err := services.NewSESService(cfg.SES, logger)
	if err != nil {
		log.Fatalf("Failed to create SES service: %v", err)
	}

	ctx := context.Background()

	// Test 1: Send a test email and capture message ID
	fmt.Println("\n--- Test 1: Send Email and Capture Message ID ---")
	testEmail := &interfaces.EmailMessage{
		From:     cfg.SES.SenderEmail,
		To:       []string{cfg.SES.SenderEmail}, // Send to self for testing
		Subject:  "SES Delivery Status Test",
		HTMLBody: "<h1>Test Email</h1><p>This is a test email to verify message ID capture and delivery status tracking.</p>",
		TextBody: "Test Email\n\nThis is a test email to verify message ID capture and delivery status tracking.",
	}

	err = sesService.SendEmail(ctx, testEmail)
	if err != nil {
		log.Printf("Failed to send test email: %v", err)
	} else {
		fmt.Printf("✅ Email sent successfully\n")
		fmt.Printf("📧 Message ID: %s\n", testEmail.MessageID)

		if testEmail.MessageID != "" {
			fmt.Printf("✅ Message ID captured successfully\n")
		} else {
			fmt.Printf("❌ Message ID not captured\n")
		}
	}

	// Test 2: Get delivery status (placeholder implementation)
	if testEmail.MessageID != "" {
		fmt.Println("\n--- Test 2: Get Delivery Status ---")
		status, err := sesService.GetDeliveryStatus(ctx, testEmail.MessageID)
		if err != nil {
			log.Printf("Failed to get delivery status: %v", err)
		} else {
			fmt.Printf("✅ Delivery status retrieved\n")
			fmt.Printf("📊 Status: %s\n", status.Status)
			fmt.Printf("🕐 Timestamp: %s\n", status.Timestamp.Format(time.RFC3339))
		}
	}

	// Test 3: Process SES Notifications
	fmt.Println("\n--- Test 3: Process SES Notifications ---")

	// Test bounce notification
	bounceNotification := &interfaces.SESNotification{
		NotificationType: "Bounce",
		MessageID:        "test-message-id-bounce",
		Timestamp:        time.Now(),
		Source:           cfg.SES.SenderEmail,
		Destination:      []string{"bounce@example.com"},
		Bounce: &interfaces.SESBounceInfo{
			BounceType:    "Permanent",
			BounceSubType: "General",
			BouncedRecipients: []interfaces.SESBouncedRecipient{
				{
					EmailAddress:   "bounce@example.com",
					Action:         "failed",
					Status:         "5.1.1",
					DiagnosticCode: "smtp; 550 5.1.1 User unknown",
				},
			},
			Timestamp:  time.Now(),
			FeedbackID: "bounce-feedback-id",
		},
	}

	result, err := sesService.ProcessSESNotification(ctx, bounceNotification)
	if err != nil {
		log.Printf("Failed to process bounce notification: %v", err)
	} else {
		fmt.Printf("✅ Bounce notification processed\n")
		fmt.Printf("📊 Result Status: %s\n", result.Status)
		fmt.Printf("🕐 Processed At: %s\n", result.ProcessedAt.Format(time.RFC3339))
	}

	// Test complaint notification
	complaintNotification := &interfaces.SESNotification{
		NotificationType: "Complaint",
		MessageID:        "test-message-id-complaint",
		Timestamp:        time.Now(),
		Source:           cfg.SES.SenderEmail,
		Destination:      []string{"complaint@example.com"},
		Complaint: &interfaces.SESComplaintInfo{
			ComplainedRecipients: []interfaces.SESComplainedRecipient{
				{
					EmailAddress: "complaint@example.com",
				},
			},
			Timestamp:             time.Now(),
			FeedbackID:            "complaint-feedback-id",
			ComplaintFeedbackType: "abuse",
		},
	}

	result, err = sesService.ProcessSESNotification(ctx, complaintNotification)
	if err != nil {
		log.Printf("Failed to process complaint notification: %v", err)
	} else {
		fmt.Printf("✅ Complaint notification processed\n")
		fmt.Printf("📊 Result Status: %s\n", result.Status)
		fmt.Printf("🕐 Processed At: %s\n", result.ProcessedAt.Format(time.RFC3339))
	}

	// Test delivery notification
	deliveryNotification := &interfaces.SESNotification{
		NotificationType: "Delivery",
		MessageID:        "test-message-id-delivery",
		Timestamp:        time.Now(),
		Source:           cfg.SES.SenderEmail,
		Destination:      []string{"delivered@example.com"},
		Delivery: &interfaces.SESDeliveryInfo{
			Timestamp:            time.Now(),
			ProcessingTimeMillis: 1500,
			Recipients:           []string{"delivered@example.com"},
			SMTPResponse:         "250 2.0.0 OK",
		},
	}

	result, err = sesService.ProcessSESNotification(ctx, deliveryNotification)
	if err != nil {
		log.Printf("Failed to process delivery notification: %v", err)
	} else {
		fmt.Printf("✅ Delivery notification processed\n")
		fmt.Printf("📊 Result Status: %s\n", result.Status)
		fmt.Printf("🕐 Processed At: %s\n", result.ProcessedAt.Format(time.RFC3339))
	}

	// Test 4: Error Categorization
	fmt.Println("\n--- Test 4: Error Categorization ---")

	testErrors := []struct {
		errorType    string
		errorMessage string
		description  string
	}{
		{
			errorType:    "Bounce",
			errorMessage: "550 5.1.1 User unknown",
			description:  "Permanent bounce - user unknown",
		},
		{
			errorType:    "Bounce",
			errorMessage: "452 4.2.2 Mailbox full",
			description:  "Temporary bounce - mailbox full",
		},
		{
			errorType:    "Complaint",
			errorMessage: "Recipient marked as spam",
			description:  "Spam complaint",
		},
		{
			errorType:    "Throttling",
			errorMessage: "Rate limit exceeded",
			description:  "Rate limiting error",
		},
		{
			errorType:    "Authentication",
			errorMessage: "Invalid credentials",
			description:  "Configuration error",
		},
	}

	for i, testError := range testErrors {
		fmt.Printf("\nTest 4.%d: %s\n", i+1, testError.description)
		category := sesService.CategorizeError(testError.errorType, testError.errorMessage)

		fmt.Printf("  📂 Category: %s\n", category.Category)
		fmt.Printf("  ⚠️  Severity: %s\n", category.Severity)
		fmt.Printf("  📝 Reason: %s\n", category.Reason)
		fmt.Printf("  🔧 Actionable: %t\n", category.Actionable)
		if category.RetryAfter != nil {
			fmt.Printf("  ⏰ Retry After: %d seconds\n", *category.RetryAfter)
		}
	}

	fmt.Println("\n=== SES Delivery Status and Notification Testing Complete ===")
}
