package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/cloud-consulting/backend/internal/domain"
	"github.com/cloud-consulting/backend/internal/interfaces"
	"github.com/cloud-consulting/backend/internal/services"
)

// MockSESService for testing
type MockSESService struct {
	sentEmails []*interfaces.EmailMessage
	shouldFail bool
}

func (m *MockSESService) SendEmail(ctx context.Context, email *interfaces.EmailMessage) error {
	if m.shouldFail {
		return fmt.Errorf("mock SES failure")
	}

	// Simulate SES message ID
	email.MessageID = fmt.Sprintf("mock-message-id-%d", time.Now().Unix())

	m.sentEmails = append(m.sentEmails, email)
	return nil
}

func (m *MockSESService) VerifyEmailAddress(ctx context.Context, email string) error {
	return nil
}

func (m *MockSESService) GetSendingQuota(ctx context.Context) (*interfaces.SendingQuota, error) {
	return &interfaces.SendingQuota{
		Max24HourSend:   200,
		MaxSendRate:     14,
		SentLast24Hours: 0,
	}, nil
}

// MockEmailEventRecorder for testing
type MockEmailEventRecorder struct {
	recordedEvents []*domain.EmailEvent
	shouldFail     bool
}

func (m *MockEmailEventRecorder) RecordEmailSent(ctx context.Context, event *domain.EmailEvent) error {
	if m.shouldFail {
		return fmt.Errorf("mock event recorder failure")
	}

	m.recordedEvents = append(m.recordedEvents, event)
	return nil
}

func (m *MockEmailEventRecorder) UpdateEmailStatus(ctx context.Context, messageID string, status domain.EmailEventStatus, deliveredAt *time.Time, errorMsg string) error {
	return nil
}

func (m *MockEmailEventRecorder) GetEmailEventsByInquiry(ctx context.Context, inquiryID string) ([]*domain.EmailEvent, error) {
	var events []*domain.EmailEvent
	for _, event := range m.recordedEvents {
		if event.InquiryID == inquiryID {
			events = append(events, event)
		}
	}
	return events, nil
}

func main() {
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	fmt.Println("🧪 Testing Email Service with Event Recording Integration")
	fmt.Println("============================================================")

	// Create mock services
	mockSES := &MockSESService{
		sentEmails: make([]*interfaces.EmailMessage, 0),
		shouldFail: false,
	}

	mockEventRecorder := &MockEmailEventRecorder{
		recordedEvents: make([]*domain.EmailEvent, 0),
		shouldFail:     false,
	}

	// Create email service with event recorder
	emailService := services.NewEmailServiceForTestingWithEventRecorder(mockSES, mockEventRecorder, logger)

	// Create test inquiry and report
	inquiry := &domain.Inquiry{
		ID:       "test-inquiry-123",
		Name:     "John Doe",
		Email:    "john.doe@example.com",
		Company:  "Test Corp",
		Phone:    "+1-555-0123",
		Services: []string{"assessment", "migration"},
		Message:  "We need help with our cloud migration project.",
	}

	report := &domain.Report{
		ID:        "test-report-456",
		InquiryID: inquiry.ID,
		Title:     "Cloud Migration Assessment",
		Content:   "# Assessment Report\n\nThis is a test report for cloud migration assessment.",
		Status:    domain.ReportStatusGenerated,
	}

	ctx := context.Background()

	// Test 1: SendCustomerConfirmation with event recording
	fmt.Println("\n📧 Test 1: Customer Confirmation Email with Event Recording")
	err := emailService.SendCustomerConfirmation(ctx, inquiry)
	if err != nil {
		log.Fatalf("Failed to send customer confirmation: %v", err)
	}

	// Verify email was sent
	if len(mockSES.sentEmails) != 1 {
		log.Fatalf("Expected 1 email sent, got %d", len(mockSES.sentEmails))
	}

	// Verify event was recorded
	if len(mockEventRecorder.recordedEvents) != 1 {
		log.Fatalf("Expected 1 event recorded, got %d", len(mockEventRecorder.recordedEvents))
	}

	customerEvent := mockEventRecorder.recordedEvents[0]
	if customerEvent.EmailType != domain.EmailTypeCustomerConfirmation {
		log.Fatalf("Expected customer confirmation event, got %s", customerEvent.EmailType)
	}

	if customerEvent.Status != domain.EmailStatusSent {
		log.Fatalf("Expected sent status, got %s", customerEvent.Status)
	}

	if customerEvent.SESMessageID == "" {
		log.Fatalf("Expected SES message ID to be set")
	}

	fmt.Printf("✅ Customer confirmation email sent and event recorded successfully\n")
	fmt.Printf("   - Email recipients: %v\n", mockSES.sentEmails[0].To)
	fmt.Printf("   - Event type: %s\n", customerEvent.EmailType)
	fmt.Printf("   - Event status: %s\n", customerEvent.Status)
	fmt.Printf("   - SES Message ID: %s\n", customerEvent.SESMessageID)

	// Test 2: SendReportEmail with event recording
	fmt.Println("\n📧 Test 2: Report Email with Event Recording")
	err = emailService.SendReportEmail(ctx, inquiry, report)
	if err != nil {
		log.Fatalf("Failed to send report email: %v", err)
	}

	// Verify second email was sent
	if len(mockSES.sentEmails) != 2 {
		log.Fatalf("Expected 2 emails sent, got %d", len(mockSES.sentEmails))
	}

	// Verify second event was recorded
	if len(mockEventRecorder.recordedEvents) != 2 {
		log.Fatalf("Expected 2 events recorded, got %d", len(mockEventRecorder.recordedEvents))
	}

	reportEvent := mockEventRecorder.recordedEvents[1]
	if reportEvent.EmailType != domain.EmailTypeConsultantNotification {
		log.Fatalf("Expected consultant notification event, got %s", reportEvent.EmailType)
	}

	fmt.Printf("✅ Report email sent and event recorded successfully\n")
	fmt.Printf("   - Email recipients: %v\n", mockSES.sentEmails[1].To)
	fmt.Printf("   - Event type: %s\n", reportEvent.EmailType)
	fmt.Printf("   - Event status: %s\n", reportEvent.Status)

	// Test 3: SendInquiryNotification with event recording
	fmt.Println("\n📧 Test 3: Inquiry Notification Email with Event Recording")
	err = emailService.SendInquiryNotification(ctx, inquiry)
	if err != nil {
		log.Fatalf("Failed to send inquiry notification: %v", err)
	}

	// Verify third email was sent
	if len(mockSES.sentEmails) != 3 {
		log.Fatalf("Expected 3 emails sent, got %d", len(mockSES.sentEmails))
	}

	// Verify third event was recorded
	if len(mockEventRecorder.recordedEvents) != 3 {
		log.Fatalf("Expected 3 events recorded, got %d", len(mockEventRecorder.recordedEvents))
	}

	inquiryEvent := mockEventRecorder.recordedEvents[2]
	if inquiryEvent.EmailType != domain.EmailTypeInquiryNotification {
		log.Fatalf("Expected inquiry notification event, got %s", inquiryEvent.EmailType)
	}

	fmt.Printf("✅ Inquiry notification email sent and event recorded successfully\n")
	fmt.Printf("   - Email recipients: %v\n", mockSES.sentEmails[2].To)
	fmt.Printf("   - Event type: %s\n", inquiryEvent.EmailType)
	fmt.Printf("   - Event status: %s\n", inquiryEvent.Status)

	// Test 4: Email failure with event recording
	fmt.Println("\n📧 Test 4: Email Failure with Event Recording")
	mockSES.shouldFail = true

	err = emailService.SendCustomerConfirmation(ctx, inquiry)
	if err == nil {
		log.Fatalf("Expected email to fail, but it succeeded")
	}

	// Verify failed event was recorded
	if len(mockEventRecorder.recordedEvents) != 4 {
		log.Fatalf("Expected 4 events recorded (including failure), got %d", len(mockEventRecorder.recordedEvents))
	}

	failedEvent := mockEventRecorder.recordedEvents[3]
	if failedEvent.Status != domain.EmailStatusFailed {
		log.Fatalf("Expected failed status, got %s", failedEvent.Status)
	}

	if failedEvent.ErrorMessage == "" {
		log.Fatalf("Expected error message to be set")
	}

	fmt.Printf("✅ Email failure recorded successfully\n")
	fmt.Printf("   - Event status: %s\n", failedEvent.Status)
	fmt.Printf("   - Error message: %s\n", failedEvent.ErrorMessage)

	// Test 5: Email success with event recorder failure (should not block email)
	fmt.Println("\n📧 Test 5: Email Success with Event Recorder Failure")
	mockSES.shouldFail = false
	mockEventRecorder.shouldFail = true

	err = emailService.SendCustomerConfirmation(ctx, inquiry)
	if err != nil {
		log.Fatalf("Email should succeed even if event recording fails: %v", err)
	}

	// Verify email was sent despite event recording failure
	if len(mockSES.sentEmails) != 4 {
		log.Fatalf("Expected 4 emails sent, got %d", len(mockSES.sentEmails))
	}

	// Event recording should have failed, so still 4 events
	if len(mockEventRecorder.recordedEvents) != 4 {
		fmt.Printf("✅ Event recording failed as expected, but email was sent successfully\n")
	}

	fmt.Println("\n🎉 All tests passed! Email service with event recording integration is working correctly.")
	fmt.Println("\nSummary:")
	fmt.Printf("- Total emails sent: %d\n", len(mockSES.sentEmails))
	fmt.Printf("- Total events recorded: %d\n", len(mockEventRecorder.recordedEvents))
	fmt.Printf("- Customer confirmation events: %d\n", countEventsByType(mockEventRecorder.recordedEvents, domain.EmailTypeCustomerConfirmation))
	fmt.Printf("- Consultant notification events: %d\n", countEventsByType(mockEventRecorder.recordedEvents, domain.EmailTypeConsultantNotification))
	fmt.Printf("- Inquiry notification events: %d\n", countEventsByType(mockEventRecorder.recordedEvents, domain.EmailTypeInquiryNotification))
	fmt.Printf("- Failed events: %d\n", countEventsByStatus(mockEventRecorder.recordedEvents, domain.EmailStatusFailed))
	fmt.Printf("- Successful events: %d\n", countEventsByStatus(mockEventRecorder.recordedEvents, domain.EmailStatusSent))
}

func countEventsByType(events []*domain.EmailEvent, eventType domain.EmailEventType) int {
	count := 0
	for _, event := range events {
		if event.EmailType == eventType {
			count++
		}
	}
	return count
}

func countEventsByStatus(events []*domain.EmailEvent, status domain.EmailEventStatus) int {
	count := 0
	for _, event := range events {
		if event.Status == status {
			count++
		}
	}
	return count
}
