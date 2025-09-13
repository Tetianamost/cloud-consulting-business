// Package shared provides test fixtures and data builders for testing
package shared

import (
	"time"

	"github.com/google/uuid"
)

// Generic test data structures that don't depend on internal packages

// TestInquiry represents a test inquiry data structure
type TestInquiry struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Company     string    `json:"company"`
	Phone       string    `json:"phone"`
	Services    []string  `json:"services"`
	Description string    `json:"description"`
	Budget      string    `json:"budget"`
	Timeline    string    `json:"timeline"`
	Status      string    `json:"status"`
	Priority    string    `json:"priority"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TestReport represents a test report data structure
type TestReport struct {
	ID          string    `json:"id"`
	InquiryID   string    `json:"inquiry_id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	Summary     string    `json:"summary"`
	Status      string    `json:"status"`
	GeneratedAt time.Time `json:"generated_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TestChatSession represents a test chat session data structure
type TestChatSession struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

// TestChatMessage represents a test chat message data structure
type TestChatMessage struct {
	ID        string    `json:"id"`
	SessionID string    `json:"session_id"`
	Content   string    `json:"content"`
	Type      string    `json:"type"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// TestEmailEvent represents a test email event data structure
type TestEmailEvent struct {
	ID          string    `json:"id"`
	InquiryID   string    `json:"inquiry_id"`
	EventType   string    `json:"event_type"`
	Recipient   string    `json:"recipient"`
	Subject     string    `json:"subject"`
	Status      string    `json:"status"`
	Timestamp   time.Time `json:"timestamp"`
	MessageID   string    `json:"message_id"`
	ErrorReason string    `json:"error_reason"`
}

// TestBedrockResponse represents a test Bedrock response data structure
type TestBedrockResponse struct {
	Content    string  `json:"content"`
	Confidence float64 `json:"confidence"`
	TokensUsed int     `json:"tokens_used"`
	Model      string  `json:"model"`
	RequestID  string  `json:"request_id"`
}

// TestBedrockOptions represents test Bedrock options data structure
type TestBedrockOptions struct {
	MaxTokens   int     `json:"max_tokens"`
	Temperature float64 `json:"temperature"`
	TopP        float64 `json:"top_p"`
	Model       string  `json:"model"`
}

// TestEmailMessage represents a test email message data structure
type TestEmailMessage struct {
	From     string   `json:"from"`
	To       []string `json:"to"`
	Subject  string   `json:"subject"`
	HTMLBody string   `json:"html_body"`
	TextBody string   `json:"text_body"`
	ReplyTo  string   `json:"reply_to"`
}

// TestSESQuota represents a test SES quota data structure
type TestSESQuota struct {
	Max24HourSend   float64 `json:"max_24_hour_send"`
	MaxSendRate     float64 `json:"max_send_rate"`
	SentLast24Hours float64 `json:"sent_last_24_hours"`
}

// TestCreateInquiryRequest represents a test create inquiry request
type TestCreateInquiryRequest struct {
	Name        string   `json:"name"`
	Email       string   `json:"email"`
	Company     string   `json:"company"`
	Phone       string   `json:"phone"`
	Services    []string `json:"services"`
	Description string   `json:"description"`
	Budget      string   `json:"budget"`
	Timeline    string   `json:"timeline"`
}

// TestDataBuilder provides methods to build test data objects
type TestDataBuilder struct{}

// NewTestDataBuilder creates a new test data builder
func NewTestDataBuilder() *TestDataBuilder {
	return &TestDataBuilder{}
}

// BuildTestInquiry creates a test inquiry with default values
func (b *TestDataBuilder) BuildTestInquiry() *TestInquiry {
	return &TestInquiry{
		ID:          uuid.New().String(),
		Name:        "John Doe",
		Email:       "john.doe@example.com",
		Company:     "Test Company",
		Phone:       "+1-555-0123",
		Services:    []string{"assessment", "migration"},
		Description: "Test inquiry description",
		Budget:      "10000-50000",
		Timeline:    "3-6 months",
		Status:      "pending",
		Priority:    "normal",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// BuildTestInquiryWithOptions creates a test inquiry with custom options
func (b *TestDataBuilder) BuildTestInquiryWithOptions(options map[string]interface{}) *TestInquiry {
	inquiry := b.BuildTestInquiry()

	if name, ok := options["name"].(string); ok {
		inquiry.Name = name
	}
	if email, ok := options["email"].(string); ok {
		inquiry.Email = email
	}
	if company, ok := options["company"].(string); ok {
		inquiry.Company = company
	}
	if services, ok := options["services"].([]string); ok {
		inquiry.Services = services
	}
	if status, ok := options["status"].(string); ok {
		inquiry.Status = status
	}
	if priority, ok := options["priority"].(string); ok {
		inquiry.Priority = priority
	}

	return inquiry
}

// BuildTestReport creates a test report with default values
func (b *TestDataBuilder) BuildTestReport() *TestReport {
	return &TestReport{
		ID:          uuid.New().String(),
		InquiryID:   uuid.New().String(),
		Title:       "Test Cloud Assessment Report",
		Content:     "This is a test report content with detailed analysis.",
		Summary:     "Test report summary",
		Status:      "completed",
		GeneratedAt: time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// BuildTestReportWithInquiry creates a test report linked to a specific inquiry
func (b *TestDataBuilder) BuildTestReportWithInquiry(inquiryID string) *TestReport {
	report := b.BuildTestReport()
	report.InquiryID = inquiryID
	return report
}

// BuildTestChatSession creates a test chat session with default values
func (b *TestDataBuilder) BuildTestChatSession() *TestChatSession {
	return &TestChatSession{
		ID:        uuid.New().String(),
		UserID:    uuid.New().String(),
		Status:    "active",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		ExpiresAt: time.Now().Add(30 * time.Minute),
	}
}

// BuildTestChatMessage creates a test chat message with default values
func (b *TestDataBuilder) BuildTestChatMessage() *TestChatMessage {
	return &TestChatMessage{
		ID:        uuid.New().String(),
		SessionID: uuid.New().String(),
		Content:   "Hello, this is a test message",
		Type:      "user",
		Status:    "delivered",
		CreatedAt: time.Now(),
	}
}

// BuildTestChatMessageWithSession creates a test chat message linked to a specific session
func (b *TestDataBuilder) BuildTestChatMessageWithSession(sessionID string) *TestChatMessage {
	message := b.BuildTestChatMessage()
	message.SessionID = sessionID
	return message
}

// BuildTestEmailEvent creates a test email event with default values
func (b *TestDataBuilder) BuildTestEmailEvent() *TestEmailEvent {
	return &TestEmailEvent{
		ID:          uuid.New().String(),
		InquiryID:   uuid.New().String(),
		EventType:   "sent",
		Recipient:   "test@example.com",
		Subject:     "Test Email Subject",
		Status:      "delivered",
		Timestamp:   time.Now(),
		MessageID:   uuid.New().String(),
		ErrorReason: "",
	}
}

// BuildTestBedrockResponse creates a test Bedrock response with default values
func (b *TestDataBuilder) BuildTestBedrockResponse() *TestBedrockResponse {
	return &TestBedrockResponse{
		Content:    "This is a test AI response with detailed analysis and recommendations.",
		Confidence: 0.95,
		TokensUsed: 150,
		Model:      "anthropic.claude-3-sonnet-20240229-v1:0",
		RequestID:  uuid.New().String(),
	}
}

// BuildTestBedrockOptions creates test Bedrock options with default values
func (b *TestDataBuilder) BuildTestBedrockOptions() *TestBedrockOptions {
	return &TestBedrockOptions{
		MaxTokens:   1000,
		Temperature: 0.7,
		TopP:        0.9,
		Model:       "anthropic.claude-3-sonnet-20240229-v1:0",
	}
}

// BuildTestEmailMessage creates a test email message with default values
func (b *TestDataBuilder) BuildTestEmailMessage() *TestEmailMessage {
	return &TestEmailMessage{
		From:     "info@cloudpartner.pro",
		To:       []string{"test@example.com"},
		Subject:  "Test Email Subject",
		HTMLBody: "<html><body><h1>Test Email</h1><p>This is a test email.</p></body></html>",
		TextBody: "Test Email\n\nThis is a test email.",
		ReplyTo:  "info@cloudpartner.pro",
	}
}

// BuildTestSESQuota creates a test SES quota with default values
func (b *TestDataBuilder) BuildTestSESQuota() *TestSESQuota {
	return &TestSESQuota{
		Max24HourSend:   200.0,
		MaxSendRate:     14.0,
		SentLast24Hours: 5.0,
	}
}

// BuildTestCreateInquiryRequest creates a test create inquiry request
func (b *TestDataBuilder) BuildTestCreateInquiryRequest() *TestCreateInquiryRequest {
	return &TestCreateInquiryRequest{
		Name:        "John Doe",
		Email:       "john.doe@example.com",
		Company:     "Test Company",
		Phone:       "+1-555-0123",
		Services:    []string{"assessment", "migration"},
		Description: "Test inquiry description",
		Budget:      "10000-50000",
		Timeline:    "3-6 months",
	}
}

// TestFixtures provides a collection of test data
type TestFixtures struct {
	Inquiries    []*TestInquiry
	Reports      []*TestReport
	ChatSessions []*TestChatSession
	ChatMessages []*TestChatMessage
	EmailEvents  []*TestEmailEvent
}

// LoadTestFixtures creates a set of test fixtures
func LoadTestFixtures() *TestFixtures {
	builder := NewTestDataBuilder()

	// Create test inquiries
	inquiry1 := builder.BuildTestInquiry()
	inquiry2 := builder.BuildTestInquiryWithOptions(map[string]interface{}{
		"name":     "Jane Smith",
		"email":    "jane.smith@example.com",
		"company":  "Another Company",
		"priority": "high",
	})

	// Create test reports
	report1 := builder.BuildTestReportWithInquiry(inquiry1.ID)
	report2 := builder.BuildTestReportWithInquiry(inquiry2.ID)

	// Create test chat sessions
	session1 := builder.BuildTestChatSession()
	session2 := builder.BuildTestChatSession()

	// Create test chat messages
	message1 := builder.BuildTestChatMessageWithSession(session1.ID)
	message2 := builder.BuildTestChatMessageWithSession(session1.ID)
	message2.Type = "ai"
	message2.Content = "Hello! How can I help you with your cloud consulting needs?"

	// Create test email events
	event1 := builder.BuildTestEmailEvent()
	event1.InquiryID = inquiry1.ID
	event2 := builder.BuildTestEmailEvent()
	event2.InquiryID = inquiry2.ID
	event2.EventType = "delivered"

	return &TestFixtures{
		Inquiries:    []*TestInquiry{inquiry1, inquiry2},
		Reports:      []*TestReport{report1, report2},
		ChatSessions: []*TestChatSession{session1, session2},
		ChatMessages: []*TestChatMessage{message1, message2},
		EmailEvents:  []*TestEmailEvent{event1, event2},
	}
}

// GetTestInquiryByID returns a test inquiry by ID from fixtures
func (f *TestFixtures) GetTestInquiryByID(id string) *TestInquiry {
	for _, inquiry := range f.Inquiries {
		if inquiry.ID == id {
			return inquiry
		}
	}
	return nil
}

// GetTestReportByInquiryID returns a test report by inquiry ID from fixtures
func (f *TestFixtures) GetTestReportByInquiryID(inquiryID string) *TestReport {
	for _, report := range f.Reports {
		if report.InquiryID == inquiryID {
			return report
		}
	}
	return nil
}

// GetTestChatMessagesBySessionID returns test chat messages by session ID from fixtures
func (f *TestFixtures) GetTestChatMessagesBySessionID(sessionID string) []*TestChatMessage {
	var messages []*TestChatMessage
	for _, message := range f.ChatMessages {
		if message.SessionID == sessionID {
			messages = append(messages, message)
		}
	}
	return messages
}
