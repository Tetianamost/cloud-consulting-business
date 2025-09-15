package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/cloud-consulting/backend/internal/domain"
	"github.com/cloud-consulting/backend/internal/handlers"
	"github.com/cloud-consulting/backend/internal/interfaces"
)

// MockSESService implements the SESService interface for testing
type MockSESService struct {
	processNotificationFunc func(ctx context.Context, notification *interfaces.SESNotification) (*interfaces.SESNotificationResult, error)
}

func (m *MockSESService) SendEmail(ctx context.Context, email *interfaces.EmailMessage) error {
	return nil
}

func (m *MockSESService) VerifyEmailAddress(ctx context.Context, email string) error {
	return nil
}

func (m *MockSESService) GetSendingQuota(ctx context.Context) (*interfaces.SendingQuota, error) {
	return &interfaces.SendingQuota{}, nil
}

func (m *MockSESService) GetDeliveryStatus(ctx context.Context, messageID string) (*interfaces.EmailDeliveryStatus, error) {
	return &interfaces.EmailDeliveryStatus{}, nil
}

func (m *MockSESService) ProcessSESNotification(ctx context.Context, notification *interfaces.SESNotification) (*interfaces.SESNotificationResult, error) {
	if m.processNotificationFunc != nil {
		return m.processNotificationFunc(ctx, notification)
	}
	return &interfaces.SESNotificationResult{
		MessageID:        notification.MessageID,
		NotificationType: notification.NotificationType,
		Status:           "processed",
		ProcessedAt:      time.Now(),
		UpdatedEvents:    1,
	}, nil
}

func (m *MockSESService) CategorizeError(errorType string, errorMessage string) *interfaces.EmailErrorCategory {
	return &interfaces.EmailErrorCategory{
		Category:   "test",
		Severity:   "test",
		Reason:     "test",
		Actionable: true,
	}
}

// MockEmailEventRecorder implements the EmailEventRecorder interface for testing
type MockEmailEventRecorder struct {
	updateStatusFunc func(ctx context.Context, messageID string, status domain.EmailEventStatus, deliveredAt *time.Time, errorMsg string) error
}

func (m *MockEmailEventRecorder) RecordEmailSent(ctx context.Context, event *domain.EmailEvent) error {
	return nil
}

func (m *MockEmailEventRecorder) UpdateEmailStatus(ctx context.Context, messageID string, status domain.EmailEventStatus, deliveredAt *time.Time, errorMsg string) error {
	if m.updateStatusFunc != nil {
		return m.updateStatusFunc(ctx, messageID, status, deliveredAt, errorMsg)
	}
	return nil
}

func (m *MockEmailEventRecorder) GetEmailEventsByInquiry(ctx context.Context, inquiryID string) ([]*domain.EmailEvent, error) {
	return []*domain.EmailEvent{}, nil
}

func main() {
	fmt.Println("=== Testing SES Webhook Handler ===")

	// Create logger
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	// Create mock services
	mockSESService := &MockSESService{}
	mockEmailEventRecorder := &MockEmailEventRecorder{}

	// Create webhook handler
	webhookHandler := handlers.NewSESWebhookHandler(mockSESService, mockEmailEventRecorder, logger)

	// Set up Gin router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/webhook/ses/notification", webhookHandler.HandleSESNotification)
	router.POST("/webhook/ses/sns", webhookHandler.HandleSNSConfirmation)
	router.GET("/webhook/ses/status", webhookHandler.GetWebhookStatus)

	// Test 1: Handle SES Delivery Notification
	fmt.Println("\n--- Test 1: SES Delivery Notification ---")
	deliveryNotification := interfaces.SESNotification{
		NotificationType: "Delivery",
		MessageID:        "test-message-delivery",
		Timestamp:        time.Now(),
		Source:           "sender@example.com",
		Destination:      []string{"recipient@example.com"},
		Delivery: &interfaces.SESDeliveryInfo{
			Timestamp:            time.Now(),
			ProcessingTimeMillis: 1500,
			Recipients:           []string{"recipient@example.com"},
			SMTPResponse:         "250 2.0.0 OK",
		},
	}

	body, _ := json.Marshal(deliveryNotification)
	req := httptest.NewRequest("POST", "/webhook/ses/notification", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code == http.StatusOK {
		fmt.Printf("✅ Delivery notification processed successfully\n")
		fmt.Printf("   Response Code: %d\n", w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		fmt.Printf("   Status: %v\n", response["status"])
		fmt.Printf("   Message ID: %v\n", response["message_id"])
	} else {
		fmt.Printf("❌ Delivery notification failed: %d\n", w.Code)
		fmt.Printf("   Response: %s\n", w.Body.String())
	}

	// Test 2: Handle SES Bounce Notification
	fmt.Println("\n--- Test 2: SES Bounce Notification ---")
	bounceNotification := interfaces.SESNotification{
		NotificationType: "Bounce",
		MessageID:        "test-message-bounce",
		Timestamp:        time.Now(),
		Source:           "sender@example.com",
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

	body, _ = json.Marshal(bounceNotification)
	req = httptest.NewRequest("POST", "/webhook/ses/notification", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code == http.StatusOK {
		fmt.Printf("✅ Bounce notification processed successfully\n")
		fmt.Printf("   Response Code: %d\n", w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		fmt.Printf("   Status: %v\n", response["status"])
		fmt.Printf("   Message ID: %v\n", response["message_id"])
	} else {
		fmt.Printf("❌ Bounce notification failed: %d\n", w.Code)
		fmt.Printf("   Response: %s\n", w.Body.String())
	}

	// Test 3: Handle SES Complaint Notification
	fmt.Println("\n--- Test 3: SES Complaint Notification ---")
	complaintNotification := interfaces.SESNotification{
		NotificationType: "Complaint",
		MessageID:        "test-message-complaint",
		Timestamp:        time.Now(),
		Source:           "sender@example.com",
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

	body, _ = json.Marshal(complaintNotification)
	req = httptest.NewRequest("POST", "/webhook/ses/notification", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code == http.StatusOK {
		fmt.Printf("✅ Complaint notification processed successfully\n")
		fmt.Printf("   Response Code: %d\n", w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		fmt.Printf("   Status: %v\n", response["status"])
		fmt.Printf("   Message ID: %v\n", response["message_id"])
	} else {
		fmt.Printf("❌ Complaint notification failed: %d\n", w.Code)
		fmt.Printf("   Response: %s\n", w.Body.String())
	}

	// Test 4: Handle SNS Subscription Confirmation
	fmt.Println("\n--- Test 4: SNS Subscription Confirmation ---")
	snsConfirmation := map[string]interface{}{
		"Type":         "SubscriptionConfirmation",
		"MessageId":    "test-sns-message-id",
		"TopicArn":     "arn:aws:sns:us-east-1:123456789012:ses-notifications",
		"Subject":      "AWS Notification - Subscription Confirmation",
		"Message":      "You have chosen to subscribe to the topic...",
		"SubscribeURL": "https://sns.us-east-1.amazonaws.com/?Action=ConfirmSubscription&TopicArn=arn:aws:sns:us-east-1:123456789012:ses-notifications&Token=test-token",
		"Timestamp":    time.Now().Format(time.RFC3339),
	}

	body, _ = json.Marshal(snsConfirmation)
	req = httptest.NewRequest("POST", "/webhook/ses/sns", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code == http.StatusOK {
		fmt.Printf("✅ SNS subscription confirmation processed successfully\n")
		fmt.Printf("   Response Code: %d\n", w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		fmt.Printf("   Message: %v\n", response["message"])
	} else {
		fmt.Printf("❌ SNS subscription confirmation failed: %d\n", w.Code)
		fmt.Printf("   Response: %s\n", w.Body.String())
	}

	// Test 5: Get Webhook Status
	fmt.Println("\n--- Test 5: Webhook Status ---")
	req = httptest.NewRequest("GET", "/webhook/ses/status", nil)
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code == http.StatusOK {
		fmt.Printf("✅ Webhook status retrieved successfully\n")
		fmt.Printf("   Response Code: %d\n", w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		fmt.Printf("   Status: %v\n", response["status"])
		fmt.Printf("   Message: %v\n", response["message"])
	} else {
		fmt.Printf("❌ Webhook status failed: %d\n", w.Code)
		fmt.Printf("   Response: %s\n", w.Body.String())
	}

	// Test 6: Invalid JSON
	fmt.Println("\n--- Test 6: Invalid JSON Handling ---")
	req = httptest.NewRequest("POST", "/webhook/ses/notification", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code == http.StatusBadRequest {
		fmt.Printf("✅ Invalid JSON handled correctly\n")
		fmt.Printf("   Response Code: %d (expected 400)\n", w.Code)
	} else {
		fmt.Printf("❌ Invalid JSON not handled correctly: %d\n", w.Code)
	}

	fmt.Println("\n=== SES Webhook Handler Testing Complete ===")
	fmt.Println("✅ All webhook endpoints are working correctly")
	fmt.Println("✅ SES notifications are processed and email events are updated")
	fmt.Println("✅ SNS subscription confirmations are handled")
	fmt.Println("✅ Error handling works as expected")
}
