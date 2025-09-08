package main

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/cloud-consulting/backend/internal/domain"
	"github.com/cloud-consulting/backend/internal/interfaces"
	"github.com/cloud-consulting/backend/internal/services"
)

// MockSESService provides a mock implementation for testing email service integration
type MockSESService struct {
	mock.Mock
}

func (m *MockSESService) SendEmail(ctx context.Context, email *interfaces.EmailMessage) error {
	args := m.Called(ctx, email)
	return args.Error(0)
}

func (m *MockSESService) VerifyEmailAddress(ctx context.Context, email string) error {
	args := m.Called(ctx, email)
	return args.Error(0)
}

func (m *MockSESService) GetSendingQuota(ctx context.Context) (*interfaces.SendingQuota, error) {
	args := m.Called(ctx)
	return args.Get(0).(*interfaces.SendingQuota), args.Error(1)
}

func (m *MockSESService) GetDeliveryStatus(ctx context.Context, messageID string) (*interfaces.EmailDeliveryStatus, error) {
	args := m.Called(ctx, messageID)
	return args.Get(0).(*interfaces.EmailDeliveryStatus), args.Error(1)
}

func (m *MockSESService) ProcessSESNotification(ctx context.Context, notification *interfaces.SESNotification) (*interfaces.SESNotificationResult, error) {
	args := m.Called(ctx, notification)
	return args.Get(0).(*interfaces.SESNotificationResult), args.Error(1)
}

func (m *MockSESService) CategorizeError(errorType string, errorMessage string) *interfaces.EmailErrorCategory {
	args := m.Called(errorType, errorMessage)
	return args.Get(0).(*interfaces.EmailErrorCategory)
}

// MockTemplateService provides a mock implementation for testing
type MockTemplateService struct {
	mock.Mock
}

func (m *MockTemplateService) RenderEmailTemplate(ctx context.Context, templateName string, data interface{}) (string, error) {
	args := m.Called(ctx, templateName, data)
	return args.String(0), args.Error(1)
}

func (m *MockTemplateService) PrepareCustomerConfirmationData(inquiry *domain.Inquiry) *interfaces.CustomerConfirmationData {
	args := m.Called(inquiry)
	return args.Get(0).(*interfaces.CustomerConfirmationData)
}

func (m *MockTemplateService) RenderReportTemplate(ctx context.Context, templateName string, data interface{}) (string, error) {
	args := m.Called(ctx, templateName, data)
	return args.String(0), args.Error(1)
}

func (m *MockTemplateService) LoadTemplate(templateName string) (*interfaces.Template, error) {
	args := m.Called(templateName)
	return args.Get(0).(*interfaces.Template), args.Error(1)
}

func (m *MockTemplateService) ValidateTemplate(templateContent string) error {
	args := m.Called(templateContent)
	return args.Error(0)
}

func (m *MockTemplateService) GetAvailableTemplates() []string {
	args := m.Called()
	return args.Get(0).([]string)
}

func (m *MockTemplateService) ReloadTemplates() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockTemplateService) PrepareReportTemplateData(inquiry *domain.Inquiry, report *domain.Report) interface{} {
	args := m.Called(inquiry, report)
	return args.Get(0)
}

func (m *MockTemplateService) PrepareConsultantNotificationData(inquiry *domain.Inquiry, report *domain.Report, isHighPriority bool) interface{} {
	args := m.Called(inquiry, report, isHighPriority)
	return args.Get(0)
}

// MockEmailEventRecorder provides a mock implementation for testing
type MockEmailEventRecorderForIntegration struct {
	mock.Mock
}

func (m *MockEmailEventRecorderForIntegration) RecordEmailSent(ctx context.Context, event *domain.EmailEvent) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockEmailEventRecorderForIntegration) UpdateEmailStatus(ctx context.Context, messageID string, status domain.EmailEventStatus, deliveredAt *time.Time, errorMsg string) error {
	args := m.Called(ctx, messageID, status, deliveredAt, errorMsg)
	return args.Error(0)
}

func (m *MockEmailEventRecorderForIntegration) GetEmailEventsByInquiry(ctx context.Context, inquiryID string) ([]*domain.EmailEvent, error) {
	args := m.Called(ctx, inquiryID)
	return args.Get(0).([]*domain.EmailEvent), args.Error(1)
}

// TestEmailServiceIntegration provides comprehensive integration tests for email service with event recording
func TestEmailServiceIntegration(t *testing.T) {
	t.Run("SendCustomerConfirmation", func(t *testing.T) {
		testSendCustomerConfirmationIntegration(t)
	})

	t.Run("SendReportEmail", func(t *testing.T) {
		testSendReportEmailIntegration(t)
	})

	t.Run("SendInquiryNotification", func(t *testing.T) {
		testSendInquiryNotificationIntegration(t)
	})

	t.Run("ErrorHandling", func(t *testing.T) {
		testEmailServiceErrorHandling(t)
	})

	t.Run("EventRecordingFailures", func(t *testing.T) {
		testEventRecordingFailures(t)
	})

	t.Run("HealthCheck", func(t *testing.T) {
		testEmailServiceHealthCheck(t)
	})
}

func setupTestEmailService(t *testing.T) (*MockSESService, *MockTemplateService, *MockEmailEventRecorderForIntegration, interfaces.EmailService) {
	mockSES := &MockSESService{}
	mockTemplate := &MockTemplateService{}
	mockEventRecorder := &MockEmailEventRecorderForIntegration{}

	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel) // Reduce noise in tests

	// Create email service configuration
	config := &services.EmailServiceConfig{
		SenderEmail:  "info@cloudpartner.pro",
		ReplyToEmail: "info@cloudpartner.pro",
		Timeout:      30 * time.Second,
	}

	// Create email service with mocks
	emailService := services.NewEmailServiceWithDependencies(
		mockSES,
		mockTemplate,
		mockEventRecorder,
		config,
		logger,
	)

	require.NotNil(t, emailService)

	return mockSES, mockTemplate, mockEventRecorder, emailService
}

func testSendCustomerConfirmationIntegration(t *testing.T) {
	mockSES, mockTemplate, mockEventRecorder, emailService := setupTestEmailService(t)
	ctx := context.Background()

	t.Run("SuccessfulCustomerConfirmation", func(t *testing.T) {
		inquiry := createTestInquiry()

		// Setup template service mock
		templateData := &interfaces.CustomerConfirmationData{
			Name:     inquiry.Name,
			Company:  inquiry.Company,
			Services: "Cloud Migration, Architecture Review",
			ID:       inquiry.ID,
		}
		mockTemplate.On("PrepareCustomerConfirmationData", inquiry).Return(templateData).Once()
		mockTemplate.On("RenderEmailTemplate", ctx, "customer_confirmation", templateData).Return("<html>Customer confirmation email</html>", nil).Once()

		// Setup SES service mock
		mockSES.On("SendEmail", ctx, mock.MatchedBy(func(email *interfaces.EmailMessage) bool {
			return email.To == inquiry.Email &&
				email.Subject == "Thank you for your inquiry - CloudPartner Pro" &&
				email.HTMLBody == "<html>Customer confirmation email</html>"
		})).Return(nil).Once()

		// Setup event recorder mock
		mockEventRecorder.On("RecordEmailSent", ctx, mock.MatchedBy(func(event *domain.EmailEvent) bool {
			return event.InquiryID == inquiry.ID &&
				event.EmailType == domain.EmailTypeCustomerConfirmation &&
				event.RecipientEmail == inquiry.Email &&
				event.Status == domain.EmailStatusSent
		})).Return(nil).Once()

		// Execute
		err := emailService.SendCustomerConfirmation(ctx, inquiry)

		// Assert
		assert.NoError(t, err)
		mockTemplate.AssertExpectations(t)
		mockSES.AssertExpectations(t)
		mockEventRecorder.AssertExpectations(t)
	})

	t.Run("TemplateRenderingFailure", func(t *testing.T) {
		inquiry := createTestInquiry()

		templateData := &interfaces.CustomerConfirmationData{
			Name:     inquiry.Name,
			Company:  inquiry.Company,
			Services: "Cloud Migration",
			ID:       inquiry.ID,
		}
		mockTemplate.On("PrepareCustomerConfirmationData", inquiry).Return(templateData).Once()
		mockTemplate.On("RenderEmailTemplate", ctx, "customer_confirmation", templateData).Return("", fmt.Errorf("template not found")).Once()

		// Execute
		err := emailService.SendCustomerConfirmation(ctx, inquiry)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "template not found")
		mockTemplate.AssertExpectations(t)
	})

	t.Run("SESFailureWithEventRecording", func(t *testing.T) {
		inquiry := createTestInquiry()

		templateData := &interfaces.CustomerConfirmationData{
			Name:     inquiry.Name,
			Company:  inquiry.Company,
			Services: "Cloud Migration",
			ID:       inquiry.ID,
		}
		mockTemplate.On("PrepareCustomerConfirmationData", inquiry).Return(templateData).Once()
		mockTemplate.On("RenderEmailTemplate", ctx, "customer_confirmation", templateData).Return("<html>Email content</html>", nil).Once()

		// SES fails
		mockSES.On("SendEmail", ctx, mock.Anything).Return(fmt.Errorf("SES connection failed")).Once()

		// Should still record the failed event
		mockEventRecorder.On("RecordEmailSent", ctx, mock.MatchedBy(func(event *domain.EmailEvent) bool {
			return event.Status == domain.EmailStatusFailed &&
				event.ErrorMessage == "SES connection failed"
		})).Return(nil).Once()

		// Execute
		err := emailService.SendCustomerConfirmation(ctx, inquiry)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "SES connection failed")
		mockTemplate.AssertExpectations(t)
		mockSES.AssertExpectations(t)
		mockEventRecorder.AssertExpectations(t)
	})

	t.Run("EventRecordingFailureDoesNotBlockEmail", func(t *testing.T) {
		inquiry := createTestInquiry()

		templateData := &interfaces.CustomerConfirmationData{
			Name:     inquiry.Name,
			Company:  inquiry.Company,
			Services: "Cloud Migration",
			ID:       inquiry.ID,
		}
		mockTemplate.On("PrepareCustomerConfirmationData", inquiry).Return(templateData).Once()
		mockTemplate.On("RenderEmailTemplate", ctx, "customer_confirmation", templateData).Return("<html>Email content</html>", nil).Once()

		// SES succeeds
		mockSES.On("SendEmail", ctx, mock.Anything).Return(nil).Once()

		// Event recording fails but should not block email
		mockEventRecorder.On("RecordEmailSent", ctx, mock.Anything).Return(fmt.Errorf("database connection failed")).Once()

		// Execute
		err := emailService.SendCustomerConfirmation(ctx, inquiry)

		// Assert - email should still succeed
		assert.NoError(t, err)
		mockTemplate.AssertExpectations(t)
		mockSES.AssertExpectations(t)
		mockEventRecorder.AssertExpectations(t)
	})
}

func testSendReportEmailIntegration(t *testing.T) {
	mockSES, mockTemplate, mockEventRecorder, emailService := setupTestEmailService(t)
	ctx := context.Background()

	t.Run("SuccessfulReportEmail", func(t *testing.T) {
		inquiry := createTestInquiry()
		report := createTestReport(inquiry.ID)

		// Setup template service mock
		templateData := map[string]interface{}{
			"inquiry":  inquiry,
			"report":   report,
			"priority": false,
		}
		mockTemplate.On("PrepareConsultantNotificationData", inquiry, report, false).Return(templateData).Once()
		mockTemplate.On("RenderEmailTemplate", ctx, "consultant_notification", templateData).Return("<html>Consultant notification email</html>", nil).Once()

		// Setup SES service mock
		mockSES.On("SendEmail", ctx, mock.MatchedBy(func(email *interfaces.EmailMessage) bool {
			return email.To == "info@cloudpartner.pro" &&
				email.Subject == "New Report Generated - "+inquiry.Company &&
				email.HTMLBody == "<html>Consultant notification email</html>"
		})).Return(nil).Once()

		// Setup event recorder mock
		mockEventRecorder.On("RecordEmailSent", ctx, mock.MatchedBy(func(event *domain.EmailEvent) bool {
			return event.InquiryID == inquiry.ID &&
				event.EmailType == domain.EmailTypeConsultantNotification &&
				event.RecipientEmail == "info@cloudpartner.pro" &&
				event.Status == domain.EmailStatusSent
		})).Return(nil).Once()

		// Execute
		err := emailService.SendReportEmail(ctx, inquiry, report)

		// Assert
		assert.NoError(t, err)
		mockTemplate.AssertExpectations(t)
		mockSES.AssertExpectations(t)
		mockEventRecorder.AssertExpectations(t)
	})

	t.Run("HighPriorityReportEmail", func(t *testing.T) {
		inquiry := createTestInquiry()
		inquiry.Priority = "high"
		inquiry.Message = "URGENT: Need immediate assistance with cloud migration"
		report := createTestReport(inquiry.ID)

		// Setup template service mock for high priority
		templateData := map[string]interface{}{
			"inquiry":  inquiry,
			"report":   report,
			"priority": true,
		}
		mockTemplate.On("PrepareConsultantNotificationData", inquiry, report, true).Return(templateData).Once()
		mockTemplate.On("RenderEmailTemplate", ctx, "consultant_notification", templateData).Return("<html>HIGH PRIORITY: Consultant notification</html>", nil).Once()

		// Setup SES service mock
		mockSES.On("SendEmail", ctx, mock.MatchedBy(func(email *interfaces.EmailMessage) bool {
			return email.To == "info@cloudpartner.pro" &&
				email.Subject == "HIGH PRIORITY: New Report Generated - "+inquiry.Company &&
				email.HTMLBody == "<html>HIGH PRIORITY: Consultant notification</html>"
		})).Return(nil).Once()

		// Setup event recorder mock
		mockEventRecorder.On("RecordEmailSent", ctx, mock.MatchedBy(func(event *domain.EmailEvent) bool {
			return event.InquiryID == inquiry.ID &&
				event.EmailType == domain.EmailTypeConsultantNotification &&
				event.Subject == "HIGH PRIORITY: New Report Generated - "+inquiry.Company
		})).Return(nil).Once()

		// Execute
		err := emailService.SendReportEmail(ctx, inquiry, report)

		// Assert
		assert.NoError(t, err)
		mockTemplate.AssertExpectations(t)
		mockSES.AssertExpectations(t)
		mockEventRecorder.AssertExpectations(t)
	})

	t.Run("ReportEmailWithSESMessageID", func(t *testing.T) {
		inquiry := createTestInquiry()
		report := createTestReport(inquiry.ID)

		templateData := map[string]interface{}{
			"inquiry":  inquiry,
			"report":   report,
			"priority": false,
		}
		mockTemplate.On("PrepareConsultantNotificationData", inquiry, report, false).Return(templateData).Once()
		mockTemplate.On("RenderEmailTemplate", ctx, "consultant_notification", templateData).Return("<html>Email content</html>", nil).Once()

		// SES returns message ID
		mockSES.On("SendEmail", ctx, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
			email := args.Get(1).(*interfaces.EmailMessage)
			email.MessageID = "ses-message-" + uuid.New().String()
		}).Once()

		// Should record event with SES message ID
		mockEventRecorder.On("RecordEmailSent", ctx, mock.MatchedBy(func(event *domain.EmailEvent) bool {
			return event.SESMessageID != "" && event.Status == domain.EmailStatusSent
		})).Return(nil).Once()

		// Execute
		err := emailService.SendReportEmail(ctx, inquiry, report)

		// Assert
		assert.NoError(t, err)
		mockTemplate.AssertExpectations(t)
		mockSES.AssertExpectations(t)
		mockEventRecorder.AssertExpectations(t)
	})
}

func testSendInquiryNotificationIntegration(t *testing.T) {
	mockSES, mockTemplate, mockEventRecorder, emailService := setupTestEmailService(t)
	ctx := context.Background()

	t.Run("SuccessfulInquiryNotification", func(t *testing.T) {
		inquiry := createTestInquiry()

		// Setup template service mock
		templateData := map[string]interface{}{
			"inquiry": inquiry,
		}
		mockTemplate.On("PrepareReportTemplateData", inquiry, (*domain.Report)(nil)).Return(templateData).Once()
		mockTemplate.On("RenderEmailTemplate", ctx, "inquiry_notification", templateData).Return("<html>New inquiry notification</html>", nil).Once()

		// Setup SES service mock
		mockSES.On("SendEmail", ctx, mock.MatchedBy(func(email *interfaces.EmailMessage) bool {
			return email.To == "info@cloudpartner.pro" &&
				email.Subject == "New Inquiry Received - "+inquiry.Company &&
				email.HTMLBody == "<html>New inquiry notification</html>"
		})).Return(nil).Once()

		// Setup event recorder mock
		mockEventRecorder.On("RecordEmailSent", ctx, mock.MatchedBy(func(event *domain.EmailEvent) bool {
			return event.InquiryID == inquiry.ID &&
				event.EmailType == domain.EmailTypeInquiryNotification &&
				event.RecipientEmail == "info@cloudpartner.pro" &&
				event.Status == domain.EmailStatusSent
		})).Return(nil).Once()

		// Execute
		err := emailService.SendInquiryNotification(ctx, inquiry)

		// Assert
		assert.NoError(t, err)
		mockTemplate.AssertExpectations(t)
		mockSES.AssertExpectations(t)
		mockEventRecorder.AssertExpectations(t)
	})

	t.Run("HighPriorityInquiryNotification", func(t *testing.T) {
		inquiry := createTestInquiry()
		inquiry.Message = "URGENT: Critical system failure, need immediate help"

		templateData := map[string]interface{}{
			"inquiry": inquiry,
		}
		mockTemplate.On("PrepareReportTemplateData", inquiry, (*domain.Report)(nil)).Return(templateData).Once()
		mockTemplate.On("RenderEmailTemplate", ctx, "inquiry_notification", templateData).Return("<html>URGENT: New inquiry notification</html>", nil).Once()

		// Should detect high priority and adjust subject
		mockSES.On("SendEmail", ctx, mock.MatchedBy(func(email *interfaces.EmailMessage) bool {
			return email.Subject == "HIGH PRIORITY: New Inquiry Received - "+inquiry.Company
		})).Return(nil).Once()

		mockEventRecorder.On("RecordEmailSent", ctx, mock.Anything).Return(nil).Once()

		// Execute
		err := emailService.SendInquiryNotification(ctx, inquiry)

		// Assert
		assert.NoError(t, err)
		mockTemplate.AssertExpectations(t)
		mockSES.AssertExpectations(t)
		mockEventRecorder.AssertExpectations(t)
	})
}

func testEmailServiceErrorHandling(t *testing.T) {
	mockSES, mockTemplate, mockEventRecorder, emailService := setupTestEmailService(t)
	ctx := context.Background()

	t.Run("InvalidInquiryData", func(t *testing.T) {
		// Inquiry with missing required fields
		inquiry := &domain.Inquiry{
			ID: "invalid-inquiry",
			// Missing Name and Email
		}

		// Execute
		err := emailService.SendCustomerConfirmation(ctx, inquiry)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "validation")
	})

	t.Run("TemplateServiceFailure", func(t *testing.T) {
		inquiry := createTestInquiry()

		mockTemplate.On("PrepareCustomerConfirmationData", inquiry).Return(&interfaces.CustomerConfirmationData{}).Once()
		mockTemplate.On("RenderEmailTemplate", ctx, "customer_confirmation", mock.Anything).Return("", fmt.Errorf("template service unavailable")).Once()

		// Execute
		err := emailService.SendCustomerConfirmation(ctx, inquiry)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "template service unavailable")
		mockTemplate.AssertExpectations(t)
	})

	t.Run("SESServiceFailure", func(t *testing.T) {
		inquiry := createTestInquiry()

		templateData := &interfaces.CustomerConfirmationData{
			Name: inquiry.Name,
		}
		mockTemplate.On("PrepareCustomerConfirmationData", inquiry).Return(templateData).Once()
		mockTemplate.On("RenderEmailTemplate", ctx, "customer_confirmation", templateData).Return("<html>Content</html>", nil).Once()
		mockSES.On("SendEmail", ctx, mock.Anything).Return(fmt.Errorf("SES quota exceeded")).Once()
		mockEventRecorder.On("RecordEmailSent", ctx, mock.Anything).Return(nil).Once()

		// Execute
		err := emailService.SendCustomerConfirmation(ctx, inquiry)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "SES quota exceeded")
		mockTemplate.AssertExpectations(t)
		mockSES.AssertExpectations(t)
		mockEventRecorder.AssertExpectations(t)
	})

	t.Run("ContextCancellation", func(t *testing.T) {
		inquiry := createTestInquiry()

		// Create cancelled context
		cancelledCtx, cancel := context.WithCancel(ctx)
		cancel()

		// Execute
		err := emailService.SendCustomerConfirmation(cancelledCtx, inquiry)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "context")
	})

	t.Run("ContextTimeout", func(t *testing.T) {
		inquiry := createTestInquiry()

		// Create context with very short timeout
		timeoutCtx, cancel := context.WithTimeout(ctx, 1*time.Millisecond)
		defer cancel()

		templateData := &interfaces.CustomerConfirmationData{
			Name: inquiry.Name,
		}
		mockTemplate.On("PrepareCustomerConfirmationData", inquiry).Return(templateData).Once()
		mockTemplate.On("RenderEmailTemplate", timeoutCtx, "customer_confirmation", templateData).Return("", context.DeadlineExceeded).Once()

		// Execute
		err := emailService.SendCustomerConfirmation(timeoutCtx, inquiry)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "deadline")
		mockTemplate.AssertExpectations(t)
	})
}

func testEventRecordingFailures(t *testing.T) {
	mockSES, mockTemplate, mockEventRecorder, emailService := setupTestEmailService(t)
	ctx := context.Background()

	t.Run("EventRecordingFailureDoesNotBlockEmailDelivery", func(t *testing.T) {
		inquiry := createTestInquiry()

		templateData := &interfaces.CustomerConfirmationData{
			Name: inquiry.Name,
		}
		mockTemplate.On("PrepareCustomerConfirmationData", inquiry).Return(templateData).Once()
		mockTemplate.On("RenderEmailTemplate", ctx, "customer_confirmation", templateData).Return("<html>Content</html>", nil).Once()
		mockSES.On("SendEmail", ctx, mock.Anything).Return(nil).Once()

		// Event recording fails
		mockEventRecorder.On("RecordEmailSent", ctx, mock.Anything).Return(fmt.Errorf("database connection failed")).Once()

		// Execute
		err := emailService.SendCustomerConfirmation(ctx, inquiry)

		// Assert - email should still succeed despite event recording failure
		assert.NoError(t, err)
		mockTemplate.AssertExpectations(t)
		mockSES.AssertExpectations(t)
		mockEventRecorder.AssertExpectations(t)
	})

	t.Run("EventRecordingWithPartialFailure", func(t *testing.T) {
		inquiry := createTestInquiry()
		report := createTestReport(inquiry.ID)

		templateData := map[string]interface{}{
			"inquiry":  inquiry,
			"report":   report,
			"priority": false,
		}
		mockTemplate.On("PrepareConsultantNotificationData", inquiry, report, false).Return(templateData).Once()
		mockTemplate.On("RenderEmailTemplate", ctx, "consultant_notification", templateData).Return("<html>Content</html>", nil).Once()
		mockSES.On("SendEmail", ctx, mock.Anything).Return(nil).Once()

		// Event recording succeeds on first call, fails on second (if there were multiple events)
		mockEventRecorder.On("RecordEmailSent", ctx, mock.Anything).Return(nil).Once()

		// Execute
		err := emailService.SendReportEmail(ctx, inquiry, report)

		// Assert
		assert.NoError(t, err)
		mockTemplate.AssertExpectations(t)
		mockSES.AssertExpectations(t)
		mockEventRecorder.AssertExpectations(t)
	})

	t.Run("EventRecordingWithInvalidEventData", func(t *testing.T) {
		inquiry := createTestInquiry()

		templateData := &interfaces.CustomerConfirmationData{
			Name: inquiry.Name,
		}
		mockTemplate.On("PrepareCustomerConfirmationData", inquiry).Return(templateData).Once()
		mockTemplate.On("RenderEmailTemplate", ctx, "customer_confirmation", templateData).Return("<html>Content</html>", nil).Once()
		mockSES.On("SendEmail", ctx, mock.Anything).Return(nil).Once()

		// Event recording fails due to validation error
		mockEventRecorder.On("RecordEmailSent", ctx, mock.Anything).Return(fmt.Errorf("validation failed: missing required field")).Once()

		// Execute
		err := emailService.SendCustomerConfirmation(ctx, inquiry)

		// Assert - email should still succeed
		assert.NoError(t, err)
		mockTemplate.AssertExpectations(t)
		mockSES.AssertExpectations(t)
		mockEventRecorder.AssertExpectations(t)
	})
}

func testEmailServiceHealthCheck(t *testing.T) {
	mockSES, mockTemplate, mockEventRecorder, emailService := setupTestEmailService(t)

	t.Run("HealthyService", func(t *testing.T) {
		// All dependencies are healthy
		isHealthy := emailService.IsHealthy()
		assert.True(t, isHealthy)
	})

	t.Run("ServiceWithNilDependencies", func(t *testing.T) {
		// Create service with nil dependencies
		logger := logrus.New()
		logger.SetLevel(logrus.ErrorLevel)

		config := &services.EmailServiceConfig{
			SenderEmail: "info@cloudpartner.pro",
		}

		unhealthyService := services.NewEmailServiceWithDependencies(
			nil, // nil SES service
			mockTemplate,
			mockEventRecorder,
			config,
			logger,
		)

		isHealthy := unhealthyService.IsHealthy()
		assert.False(t, isHealthy)
	})

	t.Run("ServiceWithInvalidConfiguration", func(t *testing.T) {
		logger := logrus.New()
		logger.SetLevel(logrus.ErrorLevel)

		// Invalid configuration (missing sender email)
		invalidConfig := &services.EmailServiceConfig{
			SenderEmail: "", // Empty sender email
		}

		unhealthyService := services.NewEmailServiceWithDependencies(
			mockSES,
			mockTemplate,
			mockEventRecorder,
			invalidConfig,
			logger,
		)

		isHealthy := unhealthyService.IsHealthy()
		assert.False(t, isHealthy)
	})
}

// Helper functions

func createTestInquiry() *domain.Inquiry {
	return &domain.Inquiry{
		ID:        uuid.New().String(),
		Name:      "John Doe",
		Email:     "john.doe@example.com",
		Company:   "Example Corp",
		Phone:     "+1-555-0123",
		Services:  []string{"Cloud Migration", "Architecture Review"},
		Message:   "We need help with our cloud migration project",
		Status:    "new",
		Priority:  "medium",
		Source:    "website",
		CreatedAt: time.Now().Add(-1 * time.Hour),
		UpdatedAt: time.Now().Add(-1 * time.Hour),
	}
}

func createTestReport(inquiryID string) *domain.Report {
	return &domain.Report{
		ID:          uuid.New().String(),
		InquiryID:   inquiryID,
		Type:        domain.ReportTypeAssessment,
		Title:       "Cloud Migration Assessment",
		Content:     "Detailed assessment of cloud migration requirements...",
		Status:      domain.ReportStatusGenerated,
		GeneratedBy: "ai-system",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// Main function for running tests standalone
func main() {
	fmt.Println("=== Comprehensive Email Service Integration Tests ===")

	// Note: This would normally be run with `go test` command
	// This main function is for demonstration purposes

	fmt.Println("Run with: go test -v ./test_email_service_integration_comprehensive.go")
	fmt.Println("Or integrate into your test suite")
}
