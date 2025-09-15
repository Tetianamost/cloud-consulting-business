// Package shared provides common mock implementations for testing
package shared

import (
	"context"
	"database/sql"

	"github.com/stretchr/testify/mock"
)

// Generic interfaces for testing (to avoid internal package dependencies)

// MockBedrockService provides a common mock implementation for Bedrock AI service
type MockBedrockService struct {
	mock.Mock
}

func (m *MockBedrockService) GenerateText(ctx context.Context, prompt string, options interface{}) (interface{}, error) {
	args := m.Called(ctx, prompt, options)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0), args.Error(1)
}

func (m *MockBedrockService) GenerateResponse(ctx context.Context, prompt string) (interface{}, error) {
	args := m.Called(ctx, prompt)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0), args.Error(1)
}

func (m *MockBedrockService) IsHealthy() bool {
	args := m.Called()
	return args.Bool(0)
}

// MockEmailService provides a common mock implementation for email service
type MockEmailService struct {
	mock.Mock
}

func (m *MockEmailService) SendReportEmail(ctx context.Context, inquiry interface{}, report interface{}) error {
	args := m.Called(ctx, inquiry, report)
	return args.Error(0)
}

func (m *MockEmailService) SendCustomerConfirmation(ctx context.Context, inquiry interface{}) error {
	args := m.Called(ctx, inquiry)
	return args.Error(0)
}

func (m *MockEmailService) SendInquiryNotification(ctx context.Context, inquiry interface{}) error {
	args := m.Called(ctx, inquiry)
	return args.Error(0)
}

func (m *MockEmailService) IsHealthy() bool {
	args := m.Called()
	return args.Bool(0)
}

// MockEmailMetricsService provides a common mock implementation for email metrics
type MockEmailMetricsService struct {
	mock.Mock
}

func (m *MockEmailMetricsService) GetEmailMetrics(ctx context.Context, timeRange interface{}) (interface{}, error) {
	args := m.Called(ctx, timeRange)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0), args.Error(1)
}

func (m *MockEmailMetricsService) GetEmailStatusByInquiry(ctx context.Context, inquiryID string) (interface{}, error) {
	args := m.Called(ctx, inquiryID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0), args.Error(1)
}

func (m *MockEmailMetricsService) GetEmailEventHistory(ctx context.Context, filters interface{}) ([]interface{}, error) {
	args := m.Called(ctx, filters)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]interface{}), args.Error(1)
}

func (m *MockEmailMetricsService) IsHealthy(ctx context.Context) bool {
	args := m.Called(ctx)
	return args.Bool(0)
}

// MockInquiryService provides a common mock implementation for inquiry service
type MockInquiryService struct {
	mock.Mock
}

func (m *MockInquiryService) CreateInquiry(ctx context.Context, req interface{}) (interface{}, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0), args.Error(1)
}

func (m *MockInquiryService) GetInquiry(ctx context.Context, id string) (interface{}, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0), args.Error(1)
}

func (m *MockInquiryService) ListInquiries(ctx context.Context, filters interface{}) ([]interface{}, error) {
	args := m.Called(ctx, filters)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]interface{}), args.Error(1)
}

func (m *MockInquiryService) UpdateInquiryStatus(ctx context.Context, id string, status string) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockInquiryService) AssignConsultant(ctx context.Context, id string, consultantID string) error {
	args := m.Called(ctx, id, consultantID)
	return args.Error(0)
}

func (m *MockInquiryService) GetInquiryCount(ctx context.Context, filters interface{}) (int64, error) {
	args := m.Called(ctx, filters)
	return args.Get(0).(int64), args.Error(1)
}

// MockReportService provides a common mock implementation for report service
type MockReportService struct {
	mock.Mock
}

func (m *MockReportService) GenerateReport(ctx context.Context, inquiry interface{}) (interface{}, error) {
	args := m.Called(ctx, inquiry)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0), args.Error(1)
}

func (m *MockReportService) GetReport(ctx context.Context, id string) (interface{}, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0), args.Error(1)
}

func (m *MockReportService) ListReports(ctx context.Context, filters interface{}) ([]interface{}, error) {
	args := m.Called(ctx, filters)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]interface{}), args.Error(1)
}

func (m *MockReportService) UpdateReportStatus(ctx context.Context, id string, status string) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

// MockTemplateService provides a common mock implementation for template service
type MockTemplateService struct {
	mock.Mock
}

func (m *MockTemplateService) RenderReportTemplate(ctx context.Context, templateName string, data interface{}) (string, error) {
	args := m.Called(ctx, templateName, data)
	return args.String(0), args.Error(1)
}

func (m *MockTemplateService) RenderEmailTemplate(ctx context.Context, templateName string, data interface{}) (string, error) {
	args := m.Called(ctx, templateName, data)
	return args.String(0), args.Error(1)
}

func (m *MockTemplateService) PrepareCustomerConfirmationData(inquiry interface{}) interface{} {
	args := m.Called(inquiry)
	return args.Get(0)
}

func (m *MockTemplateService) PrepareConsultantNotificationData(inquiry interface{}, report interface{}, isHighPriority bool) interface{} {
	args := m.Called(inquiry, report, isHighPriority)
	return args.Get(0)
}

// MockSESService provides a common mock implementation for SES service
type MockSESService struct {
	mock.Mock
	SentEmails []interface{}
}

func (m *MockSESService) SendEmail(ctx context.Context, message interface{}) error {
	args := m.Called(ctx, message)
	if args.Error(0) == nil {
		m.SentEmails = append(m.SentEmails, message)
	}
	return args.Error(0)
}

func (m *MockSESService) VerifyEmailAddress(ctx context.Context, email string) error {
	args := m.Called(ctx, email)
	return args.Error(0)
}

func (m *MockSESService) GetSendingQuota(ctx context.Context) (interface{}, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0), args.Error(1)
}

func (m *MockSESService) IsHealthy() bool {
	args := m.Called()
	return args.Bool(0)
}

// MockChatService provides a common mock implementation for chat service
type MockChatService struct {
	mock.Mock
}

func (m *MockChatService) SendMessage(ctx context.Context, sessionID string, message string) (interface{}, error) {
	args := m.Called(ctx, sessionID, message)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0), args.Error(1)
}

func (m *MockChatService) GetHistory(ctx context.Context, sessionID string) ([]interface{}, error) {
	args := m.Called(ctx, sessionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]interface{}), args.Error(1)
}

func (m *MockChatService) CreateSession(ctx context.Context, userID string) (interface{}, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0), args.Error(1)
}

func (m *MockChatService) GetSession(ctx context.Context, sessionID string) (interface{}, error) {
	args := m.Called(ctx, sessionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0), args.Error(1)
}

// MockDatabaseService provides a common mock implementation for database operations
type MockDatabaseService struct {
	mock.Mock
	Data map[string]interface{}
}

func (m *MockDatabaseService) Query(ctx context.Context, query string, args ...interface{}) (interface{}, error) {
	mockArgs := m.Called(ctx, query, args)
	if mockArgs.Get(0) == nil {
		return nil, mockArgs.Error(1)
	}
	return mockArgs.Get(0), mockArgs.Error(1)
}

func (m *MockDatabaseService) QueryRow(ctx context.Context, query string, args ...interface{}) interface{} {
	mockArgs := m.Called(ctx, query, args)
	return mockArgs.Get(0)
}

func (m *MockDatabaseService) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	mockArgs := m.Called(ctx, query, args)
	if mockArgs.Get(0) == nil {
		return nil, mockArgs.Error(1)
	}
	return mockArgs.Get(0).(sql.Result), mockArgs.Error(1)
}

func (m *MockDatabaseService) BeginTx(ctx context.Context) (interface{}, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0), args.Error(1)
}

func (m *MockDatabaseService) IsHealthy(ctx context.Context) bool {
	args := m.Called(ctx)
	return args.Bool(0)
}

// MockRows provides a mock implementation for database rows
type MockRows struct {
	mock.Mock
}

func (m *MockRows) Next() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockRows) Scan(dest ...interface{}) error {
	args := m.Called(dest)
	return args.Error(0)
}

func (m *MockRows) Close() {
	m.Called()
}

// MockRow provides a mock implementation for database row
type MockRow struct {
	mock.Mock
}

func (m *MockRow) Scan(dest ...interface{}) error {
	args := m.Called(dest)
	return args.Error(0)
}

// MockResult provides a mock implementation for database result
type MockResult struct {
	mock.Mock
}

func (m *MockResult) LastInsertId() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockResult) RowsAffected() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

// MockTx provides a mock implementation for database transaction
type MockTx struct {
	mock.Mock
}

func (m *MockTx) Query(ctx context.Context, query string, args ...interface{}) (interface{}, error) {
	mockArgs := m.Called(ctx, query, args)
	if mockArgs.Get(0) == nil {
		return nil, mockArgs.Error(1)
	}
	return mockArgs.Get(0), mockArgs.Error(1)
}

func (m *MockTx) QueryRow(ctx context.Context, query string, args ...interface{}) interface{} {
	mockArgs := m.Called(ctx, query, args)
	return mockArgs.Get(0)
}

func (m *MockTx) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	mockArgs := m.Called(ctx, query, args)
	if mockArgs.Get(0) == nil {
		return nil, mockArgs.Error(1)
	}
	return mockArgs.Get(0).(sql.Result), mockArgs.Error(1)
}

func (m *MockTx) Commit() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockTx) Rollback() error {
	args := m.Called()
	return args.Error(0)
}

// MockKnowledgeBase provides a common mock implementation for knowledge base
type MockKnowledgeBase struct {
	mock.Mock
}

func (m *MockKnowledgeBase) GetServiceOfferings(ctx context.Context) ([]interface{}, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]interface{}), args.Error(1)
}

func (m *MockKnowledgeBase) GetBestPractices(ctx context.Context, category string) ([]interface{}, error) {
	args := m.Called(ctx, category)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]interface{}), args.Error(1)
}

func (m *MockKnowledgeBase) SearchDocumentation(ctx context.Context, query string) ([]interface{}, error) {
	args := m.Called(ctx, query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]interface{}), args.Error(1)
}

// MockPromptArchitect provides a common mock implementation for prompt architect
type MockPromptArchitect struct {
	mock.Mock
}

func (m *MockPromptArchitect) BuildReportPrompt(ctx context.Context, inquiry interface{}, options interface{}) (string, error) {
	args := m.Called(ctx, inquiry, options)
	return args.String(0), args.Error(1)
}

func (m *MockPromptArchitect) BuildAnalysisPrompt(ctx context.Context, data interface{}, options interface{}) (string, error) {
	args := m.Called(ctx, data, options)
	return args.String(0), args.Error(1)
}

func (m *MockPromptArchitect) ValidatePrompt(ctx context.Context, prompt string) error {
	args := m.Called(ctx, prompt)
	return args.Error(0)
}
