package shared

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewTestDataBuilder(t *testing.T) {
	builder := NewTestDataBuilder()

	assert.NotNil(t, builder)
}

func TestBuildTestInquiry(t *testing.T) {
	builder := NewTestDataBuilder()
	inquiry := builder.BuildTestInquiry()

	assert.NotNil(t, inquiry)
	assert.NotEmpty(t, inquiry.ID)
	assert.Equal(t, "John Doe", inquiry.Name)
	assert.Equal(t, "john.doe@example.com", inquiry.Email)
	assert.Equal(t, "Test Company", inquiry.Company)
	assert.Equal(t, "pending", inquiry.Status)
	assert.Equal(t, "normal", inquiry.Priority)
	assert.Len(t, inquiry.Services, 2)
	assert.Contains(t, inquiry.Services, "assessment")
	assert.Contains(t, inquiry.Services, "migration")
}

func TestBuildTestInquiryWithOptions(t *testing.T) {
	builder := NewTestDataBuilder()
	options := map[string]interface{}{
		"name":     "Jane Smith",
		"email":    "jane.smith@example.com",
		"priority": "high",
		"services": []string{"optimization"},
	}

	inquiry := builder.BuildTestInquiryWithOptions(options)

	assert.NotNil(t, inquiry)
	assert.Equal(t, "Jane Smith", inquiry.Name)
	assert.Equal(t, "jane.smith@example.com", inquiry.Email)
	assert.Equal(t, "high", inquiry.Priority)
	assert.Len(t, inquiry.Services, 1)
	assert.Contains(t, inquiry.Services, "optimization")
}

func TestBuildTestReport(t *testing.T) {
	builder := NewTestDataBuilder()
	report := builder.BuildTestReport()

	assert.NotNil(t, report)
	assert.NotEmpty(t, report.ID)
	assert.NotEmpty(t, report.InquiryID)
	assert.Equal(t, "Test Cloud Assessment Report", report.Title)
	assert.Equal(t, "completed", report.Status)
	assert.NotEmpty(t, report.Content)
	assert.NotEmpty(t, report.Summary)
}

func TestBuildTestReportWithInquiry(t *testing.T) {
	builder := NewTestDataBuilder()
	inquiryID := "test-inquiry-123"

	report := builder.BuildTestReportWithInquiry(inquiryID)

	assert.NotNil(t, report)
	assert.Equal(t, inquiryID, report.InquiryID)
}

func TestBuildTestChatSession(t *testing.T) {
	builder := NewTestDataBuilder()
	session := builder.BuildTestChatSession()

	assert.NotNil(t, session)
	assert.NotEmpty(t, session.ID)
	assert.NotEmpty(t, session.UserID)
	assert.Equal(t, "active", session.Status)
	assert.True(t, session.ExpiresAt.After(time.Now()))
}

func TestBuildTestChatMessage(t *testing.T) {
	builder := NewTestDataBuilder()
	message := builder.BuildTestChatMessage()

	assert.NotNil(t, message)
	assert.NotEmpty(t, message.ID)
	assert.NotEmpty(t, message.SessionID)
	assert.Equal(t, "Hello, this is a test message", message.Content)
	assert.Equal(t, "user", message.Type)
	assert.Equal(t, "delivered", message.Status)
}

func TestBuildTestChatMessageWithSession(t *testing.T) {
	builder := NewTestDataBuilder()
	sessionID := "test-session-123"

	message := builder.BuildTestChatMessageWithSession(sessionID)

	assert.NotNil(t, message)
	assert.Equal(t, sessionID, message.SessionID)
}

func TestBuildTestEmailEvent(t *testing.T) {
	builder := NewTestDataBuilder()
	event := builder.BuildTestEmailEvent()

	assert.NotNil(t, event)
	assert.NotEmpty(t, event.ID)
	assert.NotEmpty(t, event.InquiryID)
	assert.Equal(t, "sent", event.EventType)
	assert.Equal(t, "test@example.com", event.Recipient)
	assert.Equal(t, "delivered", event.Status)
	assert.NotEmpty(t, event.MessageID)
}

func TestBuildTestBedrockResponse(t *testing.T) {
	builder := NewTestDataBuilder()
	response := builder.BuildTestBedrockResponse()

	assert.NotNil(t, response)
	assert.NotEmpty(t, response.Content)
	assert.Equal(t, 0.95, response.Confidence)
	assert.Equal(t, 150, response.TokensUsed)
	assert.Equal(t, "anthropic.claude-3-sonnet-20240229-v1:0", response.Model)
	assert.NotEmpty(t, response.RequestID)
}

func TestBuildTestBedrockOptions(t *testing.T) {
	builder := NewTestDataBuilder()
	options := builder.BuildTestBedrockOptions()

	assert.NotNil(t, options)
	assert.Equal(t, 1000, options.MaxTokens)
	assert.Equal(t, 0.7, options.Temperature)
	assert.Equal(t, 0.9, options.TopP)
	assert.Equal(t, "anthropic.claude-3-sonnet-20240229-v1:0", options.Model)
}

func TestBuildTestEmailMessage(t *testing.T) {
	builder := NewTestDataBuilder()
	message := builder.BuildTestEmailMessage()

	assert.NotNil(t, message)
	assert.Equal(t, "info@cloudpartner.pro", message.From)
	assert.Len(t, message.To, 1)
	assert.Contains(t, message.To, "test@example.com")
	assert.Equal(t, "Test Email Subject", message.Subject)
	assert.NotEmpty(t, message.HTMLBody)
	assert.NotEmpty(t, message.TextBody)
	assert.Equal(t, "info@cloudpartner.pro", message.ReplyTo)
}

func TestBuildTestSESQuota(t *testing.T) {
	builder := NewTestDataBuilder()
	quota := builder.BuildTestSESQuota()

	assert.NotNil(t, quota)
	assert.Equal(t, 200.0, quota.Max24HourSend)
	assert.Equal(t, 14.0, quota.MaxSendRate)
	assert.Equal(t, 5.0, quota.SentLast24Hours)
}

func TestBuildTestCreateInquiryRequest(t *testing.T) {
	builder := NewTestDataBuilder()
	request := builder.BuildTestCreateInquiryRequest()

	assert.NotNil(t, request)
	assert.Equal(t, "John Doe", request.Name)
	assert.Equal(t, "john.doe@example.com", request.Email)
	assert.Equal(t, "Test Company", request.Company)
	assert.Len(t, request.Services, 2)
	assert.Contains(t, request.Services, "assessment")
	assert.Contains(t, request.Services, "migration")
}

func TestLoadTestFixtures(t *testing.T) {
	fixtures := LoadTestFixtures()

	assert.NotNil(t, fixtures)
	assert.Len(t, fixtures.Inquiries, 2)
	assert.Len(t, fixtures.Reports, 2)
	assert.Len(t, fixtures.ChatSessions, 2)
	assert.Len(t, fixtures.ChatMessages, 2)
	assert.Len(t, fixtures.EmailEvents, 2)

	// Test first inquiry
	inquiry1 := fixtures.Inquiries[0]
	assert.Equal(t, "John Doe", inquiry1.Name)

	// Test second inquiry with custom options
	inquiry2 := fixtures.Inquiries[1]
	assert.Equal(t, "Jane Smith", inquiry2.Name)
	assert.Equal(t, "high", inquiry2.Priority)

	// Test that reports are linked to inquiries
	report1 := fixtures.Reports[0]
	assert.Equal(t, inquiry1.ID, report1.InquiryID)

	report2 := fixtures.Reports[1]
	assert.Equal(t, inquiry2.ID, report2.InquiryID)
}

func TestTestFixturesGetMethods(t *testing.T) {
	fixtures := LoadTestFixtures()

	// Test GetTestInquiryByID
	inquiry1 := fixtures.Inquiries[0]
	foundInquiry := fixtures.GetTestInquiryByID(inquiry1.ID)
	assert.NotNil(t, foundInquiry)
	assert.Equal(t, inquiry1.ID, foundInquiry.ID)

	// Test GetTestInquiryByID with non-existent ID
	notFound := fixtures.GetTestInquiryByID("non-existent")
	assert.Nil(t, notFound)

	// Test GetTestReportByInquiryID
	inquiry2 := fixtures.Inquiries[1]
	foundReport := fixtures.GetTestReportByInquiryID(inquiry2.ID)
	assert.NotNil(t, foundReport)
	assert.Equal(t, inquiry2.ID, foundReport.InquiryID)

	// Test GetTestChatMessagesBySessionID
	session1 := fixtures.ChatSessions[0]
	messages := fixtures.GetTestChatMessagesBySessionID(session1.ID)
	assert.Len(t, messages, 2) // Both messages are linked to session1
}
