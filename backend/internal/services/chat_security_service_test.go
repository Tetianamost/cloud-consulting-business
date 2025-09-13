package services

import (
	"context"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/cloud-consulting/backend/internal/interfaces"
)

// MockChatRateLimiter for testing
type MockChatRateLimiter struct {
	mock.Mock
}

func (m *MockChatRateLimiter) CheckRateLimit(ctx context.Context, userID string, action string) (*interfaces.RateLimitResult, error) {
	args := m.Called(ctx, userID, action)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*interfaces.RateLimitResult), args.Error(1)
}

func (m *MockChatRateLimiter) ResetRateLimit(ctx context.Context, userID string, action string) error {
	args := m.Called(ctx, userID, action)
	return args.Error(0)
}

func (m *MockChatRateLimiter) GetRateLimitStatus(ctx context.Context, userID string, action string) (*interfaces.RateLimitStatus, error) {
	args := m.Called(ctx, userID, action)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*interfaces.RateLimitStatus), args.Error(1)
}

// MockChatAuditLogger for testing
type MockChatAuditLogger struct {
	mock.Mock
}

func (m *MockChatAuditLogger) LogLogin(ctx context.Context, userID string, success bool, metadata map[string]interface{}) error {
	args := m.Called(ctx, userID, success, metadata)
	return args.Error(0)
}

func (m *MockChatAuditLogger) GetAuditLogs(ctx context.Context, filters *interfaces.AuditLogFilters) ([]*interfaces.AuditLog, error) {
	args := m.Called(ctx, filters)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*interfaces.AuditLog), args.Error(1)
}

func (m *MockChatAuditLogger) GetUserAuditLogs(ctx context.Context, userID string, limit int) ([]*interfaces.AuditLog, error) {
	args := m.Called(ctx, userID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*interfaces.AuditLog), args.Error(1)
}

func (m *MockChatAuditLogger) LogLogout(ctx context.Context, userID string, metadata map[string]interface{}) error {
	args := m.Called(ctx, userID, metadata)
	return args.Error(0)
}

func (m *MockChatAuditLogger) LogMessageSent(ctx context.Context, userID string, sessionID string, messageID string, metadata map[string]interface{}) error {
	args := m.Called(ctx, userID, sessionID, messageID, metadata)
	return args.Error(0)
}

func (m *MockChatAuditLogger) LogSessionCreated(ctx context.Context, userID string, sessionID string, metadata map[string]interface{}) error {
	args := m.Called(ctx, userID, sessionID, metadata)
	return args.Error(0)
}

func (m *MockChatAuditLogger) LogSessionDeleted(ctx context.Context, userID string, sessionID string, metadata map[string]interface{}) error {
	args := m.Called(ctx, userID, sessionID, metadata)
	return args.Error(0)
}

func (m *MockChatAuditLogger) LogRateLimitExceeded(ctx context.Context, userID string, action string, metadata map[string]interface{}) error {
	args := m.Called(ctx, userID, action, metadata)
	return args.Error(0)
}

func (m *MockChatAuditLogger) LogSecurityEvent(ctx context.Context, userID string, eventType string, severity string, metadata map[string]interface{}) error {
	args := m.Called(ctx, userID, eventType, severity, metadata)
	return args.Error(0)
}

func (m *MockChatAuditLogger) LogTokenRefresh(ctx context.Context, userID string, metadata map[string]interface{}) error {
	args := m.Called(ctx, userID, metadata)
	return args.Error(0)
}

func (m *MockChatAuditLogger) LogSessionAccessed(ctx context.Context, userID, sessionID string, metadata map[string]interface{}) error {
	args := m.Called(ctx, userID, sessionID, metadata)
	return args.Error(0)
}

func (m *MockChatAuditLogger) LogMessageAccessed(ctx context.Context, userID, messageID string, metadata map[string]interface{}) error {
	args := m.Called(ctx, userID, messageID, metadata)
	return args.Error(0)
}

func (m *MockChatAuditLogger) LogMessage(ctx context.Context, userID, sessionID, action string, metadata map[string]interface{}) error {
	args := m.Called(ctx, userID, sessionID, action, metadata)
	return args.Error(0)
}

func (m *MockChatAuditLogger) LogSecurityViolation(ctx context.Context, userID string, violation *interfaces.SecurityViolation) error {
	args := m.Called(ctx, userID, violation)
	return args.Error(0)
}

func (m *MockChatAuditLogger) LogUnauthorizedAccess(ctx context.Context, userID string, resource string, metadata map[string]interface{}) error {
	args := m.Called(ctx, userID, resource, metadata)
	return args.Error(0)
}

// Test helper functions
func createTestChatSecurityService() (*ChatSecurityService, *MockChatAuditLogger) {
	mockAuditLogger := &MockChatAuditLogger{}
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel) // Reduce log noise in tests

	// Create a real rate limiter for testing
	rateLimiter := NewChatRateLimiter(logger)
	service := NewChatSecurityService(logger, "test-jwt-secret", rateLimiter, mockAuditLogger)
	return service, mockAuditLogger
}

// Test CheckRateLimit
func TestChatSecurityService_CheckRateLimit_Success(t *testing.T) {
	service, _ := createTestChatSecurityService()
	ctx := context.Background()

	result, err := service.CheckRateLimit(ctx, "test-user", "message")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.Allowed) // Should be allowed for first request
}

func TestChatSecurityService_CheckRateLimit_Exceeded(t *testing.T) {
	service, _ := createTestChatSecurityService()
	ctx := context.Background()

	// Send many requests to exceed rate limit
	userID := "test-user-exceeded"
	for i := 0; i < 70; i++ { // Exceed the 60 messages per minute limit
		service.CheckRateLimit(ctx, userID, "message")
	}

	result, err := service.CheckRateLimit(ctx, userID, "message")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.Allowed) // Should be blocked after exceeding limit
}

func TestChatSecurityService_CheckRateLimit_Error(t *testing.T) {
	ctx := context.Background()

	// Test with nil rate limiter by creating service without rate limiter
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)
	mockAuditLogger := &MockChatAuditLogger{}
	serviceWithoutRateLimiter := NewChatSecurityService(logger, "test-jwt-secret", nil, mockAuditLogger)

	result, err := serviceWithoutRateLimiter.CheckRateLimit(ctx, "test-user", "message")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.Allowed) // Should allow when no rate limiter is configured
}

// Test ValidateMessageContent
func TestChatSecurityService_ValidateMessageContent_Valid(t *testing.T) {
	service, _ := createTestChatSecurityService()

	tests := []string{
		"Normal message content",
		"Message with numbers 123",
		"Message with special chars !@#$%",
		"Multi-line\nmessage\ncontent",
	}

	for _, content := range tests {
		err := service.ValidateMessageContent(content)
		assert.NoError(t, err, "Content: %s", content)
	}
}

func TestChatSecurityService_ValidateMessageContent_Invalid(t *testing.T) {
	service, _ := createTestChatSecurityService()

	tests := []struct {
		content string
		reason  string
	}{
		{"", "empty content"},
		{string(make([]byte, 10001)), "content too long"},
		{"<script>alert('xss')</script>", "contains script tag"},
		{"javascript:alert('xss')", "contains javascript protocol"},
		{"<iframe src='evil.com'></iframe>", "contains iframe tag"},
		{"SELECT * FROM users", "potential SQL injection"},
		{"DROP TABLE users", "potential SQL injection"},
	}

	for _, test := range tests {
		err := service.ValidateMessageContent(test.content)
		assert.Error(t, err, "Content should be invalid: %s (%s)", test.content, test.reason)
	}
}

// Test SanitizeMessageContent
func TestChatSecurityService_SanitizeMessageContent(t *testing.T) {
	service, _ := createTestChatSecurityService()

	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "Normal message",
			expected: "Normal message",
		},
		{
			input:    "<script>alert('xss')</script>Hello",
			expected: "Hello",
		},
		{
			input:    "Hello <b>world</b>",
			expected: "Hello world",
		},
		{
			input:    "  Multiple   spaces  ",
			expected: "Multiple spaces",
		},
		{
			input:    "<div>Content</div>",
			expected: "Content",
		},
		{
			input:    "Line 1\n\nLine 2\n\n\nLine 3",
			expected: "Line 1\nLine 2\nLine 3",
		},
	}

	for _, test := range tests {
		result := service.SanitizeMessageContent(test.input)
		assert.Equal(t, test.expected, result, "Input: %s", test.input)
	}
}

// Test ValidateSessionData
func TestChatSecurityService_ValidateSessionData_Valid(t *testing.T) {
	service, _ := createTestChatSecurityService()

	validContexts := []map[string]interface{}{
		{
			"client_name":     "Valid Client Name",
			"meeting_type":    "consultation",
			"project_context": "Valid project description",
		},
		{
			"client_name":   "Client with Numbers 123",
			"service_types": []string{"migration", "optimization"},
		},
		{
			"custom_fields": map[string]string{
				"priority": "high",
				"budget":   "100k",
			},
		},
	}

	for _, context := range validContexts {
		err := service.ValidateSessionData(context)
		assert.NoError(t, err, "Context: %v", context)
	}
}

// Test content filtering functionality
func TestChatSecurityService_FilterContent_Normal(t *testing.T) {
	service, _ := createTestChatSecurityService()

	normalMessages := []string{
		"What is AWS EC2?",
		"How much does S3 cost?",
		"Can you help with migration planning?",
		"What are the best practices for security?",
	}

	for _, message := range normalMessages {
		result, err := service.FilterContent(message)
		assert.NoError(t, err)
		assert.True(t, result.IsAllowed, "Message should be allowed: %s", message)
		assert.Equal(t, message, result.FilteredText)
	}
}

func TestChatSecurityService_FilterContent_Suspicious(t *testing.T) {
	service, _ := createTestChatSecurityService()

	suspiciousTests := []struct {
		message string
		reason  string
	}{
		{
			message: "My credit card number is 4111-1111-1111-1111",
			reason:  "contains potential credit card",
		},
		{
			message: "My SSN is 123-45-6789",
			reason:  "contains potential SSN",
		},
		{
			message: "SELECT * FROM users WHERE password = 'admin'",
			reason:  "contains SQL injection patterns",
		},
		{
			message: "../../etc/passwd",
			reason:  "contains path traversal patterns",
		},
		{
			message: "<script>alert('xss')</script>",
			reason:  "contains XSS patterns",
		},
	}

	for _, test := range suspiciousTests {
		result, err := service.FilterContent(test.message)
		assert.NoError(t, err)
		assert.False(t, result.IsAllowed, "Message should not be allowed: %s", test.message)
		assert.NotEmpty(t, result.Violations)
	}
}

// Test LogSecurityEvent
func TestChatSecurityService_LogSecurityEvent_Success(t *testing.T) {
	service, _ := createTestChatSecurityService()
	ctx := context.Background()

	event := &interfaces.SecurityEvent{
		ID:          "test-event-1",
		Type:        "suspicious_message",
		UserID:      "test-user",
		Resource:    "chat_message",
		Action:      "send_message",
		Result:      "blocked",
		IPAddress:   "192.168.1.1",
		UserAgent:   "Test Agent",
		Timestamp:   time.Now(),
		Metadata:    map[string]interface{}{"reason": "content_violation"},
		Severity:    "medium",
		Description: "Message blocked due to suspicious content",
	}

	err := service.LogSecurityEvent(ctx, event)

	assert.NoError(t, err)
}

// Test EncryptSensitiveData
func TestChatSecurityService_EncryptSensitiveData(t *testing.T) {
	service, _ := createTestChatSecurityService()

	originalData := "sensitive information"

	encryptedData, err := service.EncryptSensitiveData(originalData)
	assert.NoError(t, err)
	assert.NotEqual(t, originalData, encryptedData)
	assert.NotEmpty(t, encryptedData)
}

func TestChatSecurityService_EncryptSensitiveData_EmptyData(t *testing.T) {
	service, _ := createTestChatSecurityService()

	encryptedData, err := service.EncryptSensitiveData("")
	assert.NoError(t, err)
	assert.NotEmpty(t, encryptedData) // Should still return encrypted empty string
}

// Test DecryptSensitiveData
func TestChatSecurityService_DecryptSensitiveData(t *testing.T) {
	service, _ := createTestChatSecurityService()

	originalData := "sensitive information"

	// Encrypt first
	encryptedData, err := service.EncryptSensitiveData(originalData)
	assert.NoError(t, err)

	// Then decrypt
	decryptedData, err := service.DecryptSensitiveData(encryptedData)
	assert.NoError(t, err)
	assert.Equal(t, originalData, decryptedData)
}

func TestChatSecurityService_DecryptSensitiveData_InvalidData(t *testing.T) {
	service, _ := createTestChatSecurityService()

	// Try to decrypt invalid data
	_, err := service.DecryptSensitiveData("invalid-encrypted-data")
	assert.Error(t, err)
}

// Test HashSensitiveData
func TestChatSecurityService_HashSensitiveData(t *testing.T) {
	service, _ := createTestChatSecurityService()

	data := "sensitive information"
	hash1 := service.HashSensitiveData(data)
	hash2 := service.HashSensitiveData(data)

	assert.NotEmpty(t, hash1)
	assert.Equal(t, hash1, hash2, "Same data should produce same hash")

	differentData := "different sensitive information"
	hash3 := service.HashSensitiveData(differentData)
	assert.NotEqual(t, hash1, hash3, "Different data should produce different hashes")
}

// Test security event logging integration
func TestChatSecurityService_SecurityEventLogging(t *testing.T) {
	service, _ := createTestChatSecurityService()
	ctx := context.Background()

	metadata := map[string]interface{}{
		"ip_address": "192.168.1.1",
		"message":    "suspicious content",
	}

	// Test different severity levels
	severities := []string{"low", "medium", "high", "critical"}

	for _, severity := range severities {
		event := &interfaces.SecurityEvent{
			ID:          "test-event-1",
			Type:        "test_event",
			UserID:      "test-user",
			Resource:    "chat_message",
			Action:      "send_message",
			Result:      "blocked",
			IPAddress:   "192.168.1.1",
			UserAgent:   "Test Agent",
			Timestamp:   time.Now(),
			Metadata:    metadata,
			Severity:    severity,
			Description: "Test security event",
		}

		err := service.LogSecurityEvent(ctx, event)
		assert.NoError(t, err)
	}

}

// Benchmark tests
func BenchmarkChatSecurityService_ValidateMessageContent(b *testing.B) {
	service, _ := createTestChatSecurityService()
	message := "This is a normal message that should pass validation checks"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.ValidateMessageContent(message)
	}
}

func BenchmarkChatSecurityService_SanitizeMessageContent(b *testing.B) {
	service, _ := createTestChatSecurityService()
	message := "This is a message with <script>alert('test')</script> some HTML tags <b>bold</b> and extra   spaces"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.SanitizeMessageContent(message)
	}
}

func BenchmarkChatSecurityService_FilterContent(b *testing.B) {
	service, _ := createTestChatSecurityService()
	message := "This is a normal message asking about AWS services and best practices"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.FilterContent(message)
	}
}
